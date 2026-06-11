//! LLVM IR (text) code generation for Osprey.
//!
//! The backend walks the AST and prints LLVM assembly that clang compiles and
//! links against libc and the prebuilt C runtime archives in `compiler/bin/`
//! (`libfiber_runtime.a` / `libhttp_runtime.a`). Two anchors define correct
//! output: the C runtime ABI (those archives' symbols and conventions) and the
//! golden outputs in `compiler/examples/tested`, exercised end-to-end by
//! `crates/diff_examples.sh`. Constructs the backend does not lower return
//! [`CodegenError::Unsupported`] — it never emits a placeholder.
//!
//! Public surface: [`compile_program`] turns a parsed [`osprey_ast::Program`]
//! into a module string.

mod aggregate;
mod builder;
mod call;
mod cast;
mod collections;
mod conv;
mod effects;
mod error;
mod expr;
mod extern_call;
mod fiber;
mod genfn;
mod iter;
mod listlit;
mod llty;
mod loops;
mod lower;
mod pattern;
mod result;
mod runtime;
mod strings;
mod types;

pub use error::{CodegenError, Result};
pub use llty::{LType, Value};
pub use lower::compile_program;

#[cfg(test)]
mod tests {
    use super::*;
    use osprey_syntax::parse_program;

    fn module(src: &str) -> String {
        let parsed = parse_program(src);
        assert!(
            parsed.errors.is_empty(),
            "syntax errors: {:?}",
            parsed.errors
        );
        compile_program(&parsed.program).expect("codegen should succeed")
    }

    #[test]
    fn emits_main_and_puts_for_hello() {
        let ir = module("print(\"hello\")\n");
        assert!(ir.contains("define i32 @main()"));
        assert!(ir.contains("declare i32 @puts(i8*)"));
        assert!(ir.contains("call i32 @puts"));
        assert!(ir.contains("hello\\00"));
        assert!(ir.trim_end().ends_with('}'));
    }

    #[test]
    fn emits_arithmetic_function() {
        // A monomorphic (annotated) function is emitted as a real definition and
        // called directly; a generic one would instead inline at its call sites.
        let ir = module("fn add(a: int, b: int) -> int = a + b\nlet r = add(2, 3)\n");
        assert!(ir.contains("define i64 @add(i64 %a, i64 %b)"));
        assert!(ir.contains("add i64 %a, %b"));
        assert!(ir.contains("call i64 @add(i64 2, i64 3)"));
    }

    #[test]
    fn generic_function_inlines_at_call_site() {
        // A polymorphic function is specialised by inlining, so no monomorphic
        // definition is emitted; the call computes directly at the use site.
        let ir = module("fn identity(x) = x\nlet a = identity(7)\nprint(\"v=${a}\")\n");
        assert!(!ir.contains("@identity"));
    }

    #[test]
    fn spawn_lowers_to_a_real_fiber_with_spilled_captures() {
        // `spawn` lifts its expression into a no-arg thunk handed to the C
        // runtime's `fiber_spawn`; the local it closes over is spilled through
        // a per-spawn module global (store at the spawn site, reload in the
        // thunk), and `await` maps to `fiber_await`.
        let ir = module(
            "fn work(n: int) -> int = n * 2\n\
             fn main() -> Unit = {\n\
               let x = 21\n\
               let f = spawn work(x)\n\
               print(\"got ${await(f)}\")\n\
             }\n",
        );
        assert!(ir.contains("call i64 @fiber_spawn(i64 ()* @__fiber_closure_"));
        assert!(ir.contains("@__fiber_cap_") && ir.contains("_x = global i64 0"));
        assert!(ir.contains("call i64 @fiber_await(i64"));
    }

    #[test]
    fn inline_lambda_argument_is_lifted_to_a_function_pointer() {
        // An inline lambda flowing into a function-typed parameter is lifted to a
        // top-level `@__lambda_*` function and passed as an `i8*` code pointer,
        // not evaluated as a value — so the indirect call inside `apply` reaches
        // it the same way a bare function name would.
        let ir = module(
            "fn apply(value: int, f: (int) -> int) -> int = f(value)\n\
             let r = apply(value: 10, f: fn(x: int) => x + 1)\n\
             print(\"r=${r}\")\n",
        );
        assert!(ir.contains("define i64 @__lambda_0(i64 %x)"));
        assert!(ir.contains("@__lambda_0 to i8*"));
    }

    #[test]
    fn interpolation_uses_sprintf() {
        let ir = module("let x = 7\nprint(\"x=${x}\")\n");
        assert!(ir.contains("@sprintf"));
        assert!(ir.contains("malloc"));
    }

    #[test]
    fn match_lowers_to_phi() {
        let ir = module("fn pick(a: int, b: int) -> int = match a < b { true => a false => b }\n");
        assert!(ir.contains("icmp"));
        assert!(ir.contains("br i1"));
        assert!(ir.contains("phi i64"));
    }

    #[test]
    fn named_arguments_are_ordered_by_declaration() {
        // Call sites pass b before a; the emitted call must follow declared order.
        // `sub`'s `a - b` body infers `Result<int, MathError>`, so the call's
        // return type is `{ i64, i8 }*`; what matters here is the argument order.
        let ir = module("fn sub(a, b) = a - b\nlet r = sub(b: 1, a: 9)\n");
        assert!(ir.contains("@sub(i64 9, i64 1)"));
    }

    #[test]
    fn unsupported_construct_fails_loudly() {
        // A bare lambda used as a runtime value is not lowered (the backend lowers
        // no closures) — it must fail loudly, never silently (CLAUDE.md: no
        // placeholders, fail hard).
        let parsed = parse_program("print(fn(x) => x)\n");
        let err = compile_program(&parsed.program).unwrap_err();
        assert!(matches!(err, CodegenError::Unsupported(_)));
    }
}
