//! LLVM IR (text) code generation for Osprey — a Rust port of the textual-IR
//! emission in `compiler/internal/codegen/llvm.go`.
//!
//! The backend walks the AST and prints LLVM assembly that clang compiles and
//! links against libc (and, as more is ported, the C runtime in
//! `compiler/runtime/`). It covers the int/bool/string + functions core:
//! literals, arithmetic & comparison, `print`/`toString`, string interpolation,
//! `let`, blocks, `match` over literals, function definitions and calls, and a
//! synthesized `main`. Constructs it does not lower yet return
//! [`CodegenError::Unsupported`] — it never emits a placeholder.
//!
//! Public surface: [`compile_program`] turns a parsed [`osprey_ast::Program`]
//! into a module string.

mod aggregate;
mod builder;
mod cast;
mod collections;
mod conv;
mod effects;
mod error;
mod expr;
mod fiber;
mod iter;
mod listlit;
mod llty;
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
        let ir = module("fn add(a, b) = a + b\nlet r = add(2, 3)\n");
        assert!(ir.contains("define i64 @add(i64 %a, i64 %b)"));
        assert!(ir.contains("add i64 %a, %b"));
        assert!(ir.contains("call i64 @add(i64 2, i64 3)"));
    }

    #[test]
    fn interpolation_uses_sprintf() {
        let ir = module("let x = 7\nprint(\"x=${x}\")\n");
        assert!(ir.contains("@sprintf"));
        assert!(ir.contains("malloc"));
    }

    #[test]
    fn match_lowers_to_phi() {
        let ir = module("fn pick(a, b) = match a < b { true => a false => b }\n");
        assert!(ir.contains("icmp"));
        assert!(ir.contains("br i1"));
        assert!(ir.contains("phi i64"));
    }

    #[test]
    fn named_arguments_are_ordered_by_declaration() {
        // Call sites pass b before a; the emitted call must follow declared order.
        let ir = module("fn sub(a, b) = a - b\nlet r = sub(b: 1, a: 9)\n");
        assert!(ir.contains("call i64 @sub(i64 9, i64 1)"));
    }

    #[test]
    fn unsupported_construct_fails_loudly() {
        // `perform` is not lowered yet — it must fail loudly, never silently.
        let parsed = parse_program("perform Logger.log(\"hi\")\n");
        let err = compile_program(&parsed.program).unwrap_err();
        assert!(matches!(err, CodegenError::Unsupported(_)));
    }
}
