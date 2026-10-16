// Osprey -> macOS graphics bridge: a Cocoa window backed by a Metal layer,
// exposed as a flat C ABI for Osprey's FFI ([FFI-LINK-DIRECTIVES], [FFI-PTR] —
// docs/specs/compiler/0019-ForeignFunctionInterface.md).
//
// The division of labour is the same one every real-time renderer uses. Osprey
// is the host: it owns the window, the clock, and the scene state it recomputes
// every frame. The GPU is the rasteriser: a fragment shader loaded from disk
// shades every pixel from the parameter slots Osprey pushes into it.
//
// Osprey's own `gpu*` builtins are NOT this path — [GPU-BACKEND-HOST] is the
// only backend implemented, so they run on the CPU. This bridge is the seam
// where [GPU-BACKEND-DEVICE] would plug in once it exists.

#import <Cocoa/Cocoa.h>
#import <Metal/Metal.h>
#import <QuartzCore/CAMetalLayer.h>
#include <stdint.h>

// The slot count, the two fixed-point scales, the uniform block and the
// argument validation are the Osprey-facing ABI, declared once in
// ospgfx_abi.h and shared with the Direct3D 12 bridge. `Uniforms` in
// base.metal mirrors `OspGfxUniforms` by hand — the GPU cannot see the header
// — and crates/osprey-cli/tests/graphics_scenes.rs asserts the mirror agrees.
#include "ospgfx_abi.h"

@interface OspGfxContext : NSObject <NSWindowDelegate>
@property(nonatomic, strong) NSWindow *window;
@property(nonatomic, strong) CAMetalLayer *layer;
@property(nonatomic, strong) id<MTLCommandQueue> queue;
@property(nonatomic, strong) id<MTLRenderPipelineState> pipeline;
@property(nonatomic) BOOL closed;
@property(nonatomic) double openedAt;
@property(nonatomic) OspGfxUniforms uniforms;
@end

@implementation OspGfxContext
- (BOOL)windowShouldClose:(NSWindow *)sender {
  (void)sender;
  self.closed = YES;
  return YES;
}
@end

// Seconds from a monotonic clock, so a paused or reset wall clock cannot make
// an animation jump backwards.
static double osp_gfx_now(void) {
  return (double)clock_gettime_nsec_np(CLOCK_MONOTONIC) / 1e9;
}

// Bring up NSApplication once, as a regular foreground app so the window
// actually takes focus when launched from a terminal.
static void osp_gfx_boot_app(void) {
  static int booted = 0;
  if (booted) {
    return;
  }
  booted = 1;
  [NSApplication sharedApplication];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp finishLaunching];
  [NSApp activateIgnoringOtherApps:YES];
}

// Compile the Metal library Osprey named. A shader that will not compile is
// reported with its diagnostics and fails the open, rather than presenting a
// blank window that looks like a renderer bug.
static id<MTLLibrary> osp_gfx_library(id<MTLDevice> device,
                                      const char *shaderPath) {
  NSError *error = nil;
  NSString *path = [NSString stringWithUTF8String:shaderPath ? shaderPath : ""];
  NSString *source = [NSString stringWithContentsOfFile:path
                                               encoding:NSUTF8StringEncoding
                                                  error:&error];
  if (!source) {
    fprintf(stderr, "ospgfx: cannot read shader %s\n", shaderPath);
    return nil;
  }
  id<MTLLibrary> library = [device newLibraryWithSource:source
                                                options:nil
                                                  error:&error];
  if (!library) {
    fprintf(stderr, "ospgfx: shader failed to compile:\n%s\n",
            [[error localizedDescription] UTF8String]);
  }
  return library;
}

// Select one fragment entry out of that library. Several scenes share a single
// library and differ only in which entry they name, which is why the entry is a
// parameter: a run-time `newLibraryWithSource:` has no include path, so sharing
// a shader means sharing one file, not including a header from several.
static id<MTLRenderPipelineState> osp_gfx_pipeline(id<MTLDevice> device,
                                                   const char *shaderPath,
                                                   const char *fragmentName) {
  id<MTLLibrary> library = osp_gfx_library(device, shaderPath);
  if (!library) {
    return nil;
  }
  const char *entry = osp_gfx_fragment_entry(fragmentName);
  NSError *error = nil;
  MTLRenderPipelineDescriptor *desc = [[MTLRenderPipelineDescriptor alloc] init];
  desc.vertexFunction =
      [library newFunctionWithName:[NSString stringWithUTF8String:OSP_GFX_VERTEX_ENTRY]];
  desc.fragmentFunction =
      [library newFunctionWithName:[NSString stringWithUTF8String:entry]];
  desc.colorAttachments[0].pixelFormat = MTLPixelFormatBGRA8Unorm;
  if (!desc.vertexFunction || !desc.fragmentFunction) {
    fprintf(stderr, "ospgfx: shader needs osp_vertex and %s\n", entry);
    return nil;
  }
  id<MTLRenderPipelineState> pipeline =
      [device newRenderPipelineStateWithDescriptor:desc error:&error];
  if (!pipeline) {
    fprintf(stderr, "ospgfx: pipeline failed: %s\n",
            [[error localizedDescription] UTF8String]);
  }
  return pipeline;
}

static void osp_gfx_attach_layer(OspGfxContext *ctx, id<MTLDevice> device,
                                 NSRect frame) {
  NSView *view = [[NSView alloc] initWithFrame:frame];
  ctx.layer = [CAMetalLayer layer];
  ctx.layer.device = device;
  ctx.layer.pixelFormat = MTLPixelFormatBGRA8Unorm;
  ctx.layer.framebufferOnly = YES;
  [view setWantsLayer:YES];
  [view setLayer:ctx.layer];
  [ctx.window setContentView:view];
  CGFloat backing = [ctx.window backingScaleFactor];
  ctx.layer.contentsScale = backing;
  ctx.layer.drawableSize =
      CGSizeMake(frame.size.width * backing, frame.size.height * backing);
}

static void osp_gfx_make_window(OspGfxContext *ctx, id<MTLDevice> device,
                                NSRect frame, const char *title) {
  ctx.window = [[NSWindow alloc]
      initWithContentRect:frame
                styleMask:(NSWindowStyleMaskTitled | NSWindowStyleMaskClosable |
                           NSWindowStyleMaskMiniaturizable)
                  backing:NSBackingStoreBuffered
                    defer:NO];
  [ctx.window setTitle:[NSString stringWithUTF8String:title ? title : "Osprey"]];
  [ctx.window setDelegate:ctx];
  osp_gfx_attach_layer(ctx, device, frame);
  [ctx.window center];
  [ctx.window makeKeyAndOrderFront:nil];
}

// `scale` magnifies the design resolution into window points: 320x200 at
// scale 3 opens a 960x600 window, drawn at the display's full backing density.
// `fragmentName` picks one entry out of the shared library at `shaderPath`.
void *osp_gfx_open(int64_t width, int64_t height, int64_t scale,
                   const char *title, const char *shaderPath,
                   const char *fragmentName) {
  if (!osp_gfx_extent_ok(width, height)) {
    return NULL;
  }
  scale = osp_gfx_clamp_scale(scale);
  osp_gfx_boot_app();
  id<MTLDevice> device = MTLCreateSystemDefaultDevice();
  if (!device) {
    fprintf(stderr, "ospgfx: no Metal device\n");
    return NULL;
  }
  id<MTLRenderPipelineState> pipeline =
      osp_gfx_pipeline(device, shaderPath, fragmentName);
  if (!pipeline) {
    return NULL;
  }
  OspGfxContext *ctx = [[OspGfxContext alloc] init];
  ctx.openedAt = osp_gfx_now();
  ctx.queue = [device newCommandQueue];
  ctx.pipeline = pipeline;
  osp_gfx_make_window(ctx, device,
                      NSMakeRect(0, 0, (CGFloat)(width * scale),
                                 (CGFloat)(height * scale)),
                      title);
  return (__bridge_retained void *)ctx;
}

// Store one slot through the shared ABI helper. The uniforms are an ObjC
// property, so the struct is copied out, written and copied back; the slot
// range check and the division live in ospgfx_abi.h, once for both bridges.
static int64_t osp_gfx_store(void *handle, int64_t slot, int64_t value,
                             float scale) {
  OspGfxContext *ctx = (__bridge OspGfxContext *)handle;
  if (!ctx) {
    return 0;
  }
  OspGfxUniforms u = ctx.uniforms;
  int64_t stored = osp_gfx_slot_store(&u, slot, value, scale);
  ctx.uniforms = u;
  return stored;
}

// Push one moving scene parameter, in 4096ths; the shader sees a float.
int64_t osp_gfx_set(void *handle, int64_t slot, int64_t value) {
  return osp_gfx_store(handle, slot, value, OSP_GFX_FIXED_ONE);
}

// Push one authored tuning constant, in thousandths. Dividing an exact integer
// by 1000 rounds to the nearest float, which is the same float the shader would
// have held had the number stayed written in its source — so moving a constant
// out of the shader and into Osprey is a refactor, not a re-grade.
int64_t osp_gfx_set_milli(void *handle, int64_t slot, int64_t thousandths) {
  return osp_gfx_store(handle, slot, thousandths, OSP_GFX_MILLI_ONE);
}

static void osp_gfx_pump_events(void) {
  NSEvent *event;
  while ((event = [NSApp nextEventMatchingMask:NSEventMaskAny
                                     untilDate:nil
                                        inMode:NSDefaultRunLoopMode
                                       dequeue:YES])) {
    [NSApp sendEvent:event];
  }
}

// Resolution and elapsed time are the two uniforms the host fills in rather
// than Osprey: only the drawable knows how big it turned out at this display's
// backing density, and the clock has moved since the last frame was pushed.
static OspGfxUniforms osp_gfx_frame_uniforms(OspGfxContext *ctx, CGSize size) {
  OspGfxUniforms u = ctx.uniforms;
  u.res[0] = (float)size.width;
  u.res[1] = (float)size.height;
  u.time = (float)(osp_gfx_now() - ctx.openedAt);
  ctx.uniforms = u;
  return u;
}

// One full-screen triangle, shaded by the selected fragment entry from the
// slots Osprey pushed, presented on the drawable.
static void osp_gfx_encode(OspGfxContext *ctx, id<CAMetalDrawable> drawable,
                           OspGfxUniforms u) {
  MTLRenderPassDescriptor *pass = [MTLRenderPassDescriptor renderPassDescriptor];
  pass.colorAttachments[0].texture = drawable.texture;
  pass.colorAttachments[0].loadAction = MTLLoadActionClear;
  pass.colorAttachments[0].storeAction = MTLStoreActionStore;
  pass.colorAttachments[0].clearColor = MTLClearColorMake(0, 0, 0, 1);
  id<MTLCommandBuffer> commands = [ctx.queue commandBuffer];
  id<MTLRenderCommandEncoder> encoder =
      [commands renderCommandEncoderWithDescriptor:pass];
  [encoder setRenderPipelineState:ctx.pipeline];
  [encoder setFragmentBytes:&u length:sizeof(u) atIndex:0];
  [encoder drawPrimitives:MTLPrimitiveTypeTriangle vertexStart:0 vertexCount:3];
  [encoder endEncoding];
  [commands presentDrawable:drawable];
  [commands commit];
}

// Shade every pixel from the current slots and present. Returns 1 while the
// window is open, 0 once the user has closed it.
int64_t osp_gfx_draw(void *handle) {
  OspGfxContext *ctx = (__bridge OspGfxContext *)handle;
  if (!ctx || ctx.closed) {
    return 0;
  }
  @autoreleasepool {
    id<CAMetalDrawable> drawable = [ctx.layer nextDrawable];
    if (drawable) {
      osp_gfx_encode(ctx, drawable,
                     osp_gfx_frame_uniforms(ctx, ctx.layer.drawableSize));
    }
    osp_gfx_pump_events();
  }
  return ctx.closed ? 0 : 1;
}

// Milliseconds since the window opened. Animations advance on wall-clock time
// so they play at the same speed whatever frame rate the host sustains.
int64_t osp_gfx_ticks(void *handle) {
  OspGfxContext *ctx = (__bridge OspGfxContext *)handle;
  return ctx ? osp_gfx_elapsed_ms(ctx.openedAt, osp_gfx_now()) : 0;
}

int64_t osp_gfx_close(void *handle) {
  OspGfxContext *ctx = (__bridge_transfer OspGfxContext *)handle;
  if (!ctx) {
    return 0;
  }
  [ctx.window setDelegate:nil];
  [ctx.window orderOut:nil];
  return 1;
}
