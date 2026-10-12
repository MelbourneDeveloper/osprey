//! The Osprey-facing C ABI — slot count, two fixed-point scales, uniform block,
//! six exports — is declared ONCE in `ospgfx_abi.h` and included by both
//! bridges. The two shader libraries cannot include it (each bridge compiles
//! its shader from source at run time with no include path), so they mirror
//! the uniform layout by hand. These tests pin the header as the canonical
//! declaration, the two mirrors against it, and the bridges to including it
//! rather than restating it: a disagreement compiles cleanly and then shifts
//! `res` and `time` under every scene on one platform.
use super::{read_required, repo_root};
use std::path::PathBuf;

const ABI_HEADER: &str = "ospgfx_abi.h";

/// Slots 0..=7 carry motion; 8..=20 carry the authored grade. `base/base.osp`
/// and both shader libraries spell this out, and the header sizes the array.
const LOOK_SLOT_LAST: usize = 20;

/// The whole Osprey-facing contract, in the order `base/base.osp` declares it.
const BRIDGE_EXPORTS: [&str; 6] = [
    "osp_gfx_open",
    "osp_gfx_set",
    "osp_gfx_set_milli",
    "osp_gfx_draw",
    "osp_gfx_ticks",
    "osp_gfx_close",
];

/// Declarations a bridge may no longer carry itself, now that the header does.
const HEADER_OWNED: [&str; 5] = [
    "#define OSP_GFX_SLOTS",
    "typedef struct OspGfxUniforms",
    "#define OSP_GFX_FIXED_ONE",
    "#define OSP_GFX_MILLI_ONE",
    "#define OSP_GFX_MAX_EXTENT",
];

fn graphics() -> PathBuf {
    repo_root().join("examples/graphics")
}

/// The integer following `needle` in `source`, e.g. the `24` in
/// `#define OSP_GFX_SLOTS 24`.
fn number_after(source: &str, needle: &str) -> Option<usize> {
    source.split(needle).nth(1).and_then(|rest| {
        let digits: String = rest
            .trim_start()
            .chars()
            .take_while(char::is_ascii_digit)
            .collect();
        digits.parse().ok()
    })
}

/// The Windows bridge is three files — the shared ABI header's D3D half, the
/// one-time bring-up and the exported entry points — and it is the union of
/// them that has to match ospgfx.m, so it is compared as one text.
fn d3d12_bridge() -> String {
    ["ospgfx_d3d12.h", "ospgfx_d3d12_setup.c", "ospgfx_d3d12.c"]
        .iter()
        .map(|file| read_required(&graphics().join(file)))
        .collect::<Vec<_>>()
        .join("\n")
}

/// Whether `source` declares the uniform block's four fields, with `res`
/// spelled as that dialect spells a float pair.
fn declares_uniform_block(source: &str, res: &str) -> bool {
    source.contains("float slot[OSP_GFX_SLOTS];")
        && source.contains(res)
        && source.contains("float time;")
        && source.contains("float pad;")
}

#[test]
fn the_header_is_the_one_declaration_and_both_shader_libraries_mirror_it() {
    let header = read_required(&graphics().join(ABI_HEADER));
    let slots = number_after(&header, "#define OSP_GFX_SLOTS");
    assert!(
        slots.unwrap_or_default() > LOOK_SLOT_LAST,
        "OSP_GFX_SLOTS in {ABI_HEADER} is too small for look slot {LOOK_SLOT_LAST}"
    );
    assert!(
        declares_uniform_block(&header, "float res[2];")
            && header.contains("#define OSP_GFX_FIXED_ONE 4096.0f")
            && header.contains("#define OSP_GFX_MILLI_ONE 1000.0f"),
        "{ABI_HEADER} must declare the uniform block and both fixed-point scales"
    );
    for name in ["base.metal", "base.hlsl"] {
        let source = read_required(&graphics().join(name));
        assert_eq!(
            number_after(&source, "#define OSP_GFX_SLOTS"),
            slots,
            "OSP_GFX_SLOTS in {name} disagrees with {ABI_HEADER}"
        );
        assert!(
            declares_uniform_block(&source, "float2 res;"),
            "{name}'s uniform layout lost a field or stopped sizing itself from OSP_GFX_SLOTS"
        );
    }
}

#[test]
fn both_bridges_include_the_header_and_restate_none_of_it() {
    let metal = read_required(&graphics().join("ospgfx.m"));
    let d3d = d3d12_bridge();
    for (name, bridge) in [("ospgfx.m", &metal), ("the D3D12 bridge", &d3d)] {
        assert!(
            bridge.contains(&format!("#include \"{ABI_HEADER}\"")),
            "{name} must include {ABI_HEADER}"
        );
        for restated in HEADER_OWNED {
            assert!(
                !bridge.contains(restated),
                "{name} restates `{restated}` instead of taking it from {ABI_HEADER}"
            );
        }
    }
}

#[test]
fn both_bridges_present_one_c_abi_to_the_shared_host() {
    let metal = read_required(&graphics().join("ospgfx.m"));
    let d3d = d3d12_bridge();
    // Osprey's `extern fn` declarations are written once and linked against
    // whichever bridge the platform built: the same six symbols, a fragment
    // entry as a parameter, and a diagnostic that names the entry it could not
    // find. A symbol missing from one side is a link error there and nothing
    // at all on the other.
    for (name, bridge) in [("ospgfx.m", &metal), ("the D3D12 bridge", &d3d)] {
        for export in BRIDGE_EXPORTS {
            assert!(
                bridge.contains(&format!("{export}(")),
                "{name} must keep exporting {export}"
            );
        }
        for shared in [
            "const char *fragmentName)",
            r#""ospgfx: shader needs osp_vertex and %s\n", entry"#,
        ] {
            assert!(bridge.contains(shared), "{name} must keep `{shared}`");
        }
    }
    let d3d_abi = read_required(&graphics().join("ospgfx_d3d12.c"));
    for export in BRIDGE_EXPORTS {
        assert!(
            d3d_abi.contains(&format!("OSP_GFX_API int64_t {export}("))
                || d3d_abi.contains(&format!("OSP_GFX_API void *{export}(")),
            "the D3D12 bridge must mark {export} for export from its DLL"
        );
    }
}

#[test]
fn the_host_names_one_shader_path_and_agrees_on_the_slot_ranges() {
    let host = read_required(&graphics().join("base/base.osp"));
    // The host's own view of the layout: where the grade starts and ends.
    assert!(
        host.contains("fn motionCount() = 8")
            && host.contains("fn lookFirst() = 8")
            && host.contains(&format!("fn lookLast() = {LOOK_SLOT_LAST}")),
        "the shared host must agree with the shaders about the slot ranges"
    );
    // The scenes name one shader path for every platform; the extension swap in
    // the Windows bridge is what keeps that true without editing a scene.
    assert!(
        d3d12_bridge().contains(r#"#define OSP_GFX_SHADER_EXT ".hlsl""#)
            && host.contains(r#"fn shaderPath() = "examples/graphics/base.metal""#),
        "the host must name one shader path and the D3D12 bridge must resolve it"
    );
}
