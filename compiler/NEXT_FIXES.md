# Compiler Fixes Required

Refer to the main branch often where all tests are passing. There was major refactoring in this branch

## 1. Channel Expression Support

**Priority: HIGH** - Many fiber tests failing

Fix support for channel AST expressions in `codegen/generation.go`. It is already working in main where the tests pass fine

```
unsupported expression: *ast.ChannelCreateExpression
unsupported expression: *ast.ChannelSendExpression  
unsupported expression: *ast.ChannelRecvExpression
```

**Fix**: Fix cases in `generateExpression()` switch statement to call existing functions:
- `*ast.ChannelCreateExpression` → `g.generateChannelCreateExpression()`
- `*ast.ChannelSendExpression` → `g.generateChannelSendExpression()`
- `*ast.ChannelRecvExpression` → `g.generateChannelRecvExpression()`

## 2. Select Expression Support

**Priority: HIGH** - Fiber tests failing

```
unsupported expression: *ast.SelectExpression
```

**Fix**: Fix case in `generateExpression()` switch statement:
- `*ast.SelectExpression` → `g.generateSelectExpression()`

## 3. Module Access Expression Support

**Priority: MEDIUM** - Module isolation tests failing

```
unsupported expression: *ast.ModuleAccessExpression
```

**Fix**: Fix case in `generateExpression()` switch statement:
- `*ast.ModuleAccessExpression` → `g.generateModuleAccessExpression()`

## 4. HTTP Tests Broken In This Branch

**Priority: MEDIUM** - All HTTP tests failing

HTTP functions showing as undefined variables:
- `httpCreateServer`
- `httpCreateClient`
- `httpListen`
- `httpSend`

**Issue**: Functions exist in builtin registry but type inferer not registering them.

**Fix**: Check `type_inferer.go` - ensure HTTP functions are being added to type environment during initialization.

## 5. WebSocket Tests Broken

**Priority: MEDIUM** - WebSocket tests failing

WebSocket functions showing as undefined variables:
- `websocketCreateServer`
- `websocketConnect`

**Fix**: Same as HTTP - check type inferer registration.

## 6. Algebraic Effects Issues

**Priority: LOW** - Effects tests failing

```
unhandled effect: Unhandled effect 'Logger.log'
unhandled effect: Unhandled effect 'Error.fail'  
undefined variable 'newVal'
```

**Fix**: 
- Review effect handler implementation
- Check variable scoping in effect handlers
- Ensure all effects are properly handled or forwarded

## 7. Rust Interop

**Priority: LOW** - Rust integration failing

```
undefined variable 'rust_add'
```

**Fix**: Add Rust extern function declarations to builtin registry or implement proper foreign function interface.

## 8. Validation/Constraint Issues

**Priority: LOW** - Some validation tests failing

```
constraint validation failures
```

**Fix**: Review constraint validation implementation for edge cases.

---

## Testing Strategy

1. Fix channel expressions first - will unlock most fiber tests
2. Test with: `make test | grep -E "(Channel|Select|Module)"`
3. Then fix HTTP/WebSocket registration issues
4. Test with: `make test | grep -E "(HTTP|WebSocket)"`

## Files to Focus On

- `internal/codegen/generation.go` - Add missing expression cases
- `internal/codegen/type_inferer.go` - Fix function registration
- `internal/codegen/fiber_generation.go` - Already cleaned up
- `internal/codegen/builtin_registry.go` - Already has functions registered 