# 🚀 ALGEBRAIC EFFECTS EXAMPLES 🔥

This directory contains examples demonstrating **PRIMO ALGEBRAIC EFFECTS** in Osprey! 

Algebraic effects are a revolutionary approach to handling side effects in functional programming languages. They allow you to:

- **Suspend computation** at any point with `perform`
- **Handle effects functionally** with structured handlers  
- **Resume execution** exactly where you left off
- **Compose effects** cleanly without callback hell
- **Separate pure and effectful code** at the type level

## Examples

### `algebraic_effects.osp` 
Basic effect declaration - shows how to declare effects and pure functions.

### `algebraic_effects_perform.osp`
Demonstrates `perform` expressions and unhandled effect detection.

### `algebraic_effects_complete.osp` 
Complete example with multiple effects showing the full system in action.

### `algebraic_effects_with_handlers.osp`
Shows future handler syntax (commented out until fully implemented).

## Current Implementation Status

✅ **IMPLEMENTED:**
- Effect declarations (`effect EffectName { ... }`)
- Effect annotations on function types (`fn foo() -> T !Effect`)
- Perform expressions (`perform Effect.operation(args)`)
- Unhandled effect detection and runtime errors
- CPS infrastructure for future handlers
- Effect registry and type tracking

🚧 **COMING SOON:**
- Handler expressions (`with handler Effect { ... } do { ... }`)
- Resume operations (`resume(value)`)
- Effect set composition (`![Effect1, Effect2]`)
- Full CPS transformation with continuation capture
- Zero-overhead effect compilation

## How It Works

Algebraic effects work by:

1. **Effect Declaration**: Define what operations an effect supports
2. **Effect Performance**: Use `perform` to suspend computation and invoke effect operations  
3. **Effect Handling**: Catch and handle effect operations with custom logic
4. **Continuation Resumption**: Resume computation from where it was suspended

This creates a powerful abstraction for handling:
- **State management** (get/set operations)
- **I/O operations** (read/write with error handling)
- **Async/concurrency** (async/await without promises)
- **Exception handling** (throw/catch as effects)
- **Non-determinism** (choice/fail for backtracking)

## PRIMO Features

- **Type-safe**: Effects are tracked in the type system
- **Composable**: Effects can be combined and nested
- **Zero-overhead**: Compiles to efficient LLVM IR  
- **Structured**: No callback hell or monad stacks
- **Direct-style**: Write code that looks synchronous but handles effects

**ALGEBRAIC EFFECTS = THE FUTURE OF PROGRAMMING! 🔥** 