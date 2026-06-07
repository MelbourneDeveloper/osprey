//! Hindley-Milner type inference for Osprey — a Rust port of
//! `compiler/internal/codegen/type_inference.go` (and the slices of
//! `builtin_registry.go` / `match_validation.go` it depends on).
//!
//! The pipeline is the textbook one: a [`ty::Type`] language, an index-addressed
//! union-find substitution ([`ctx::InferCtx`]), [`unify`](unify::unify)
//! with the Osprey-specific rules (`any`, bare-collection generics, structural
//! records, Result auto-unwrap), let-polymorphism ([`env`]), and a two-pass
//! [`check::check_program`] driver over the AST.
//!
//! Public surface: [`check_program`] takes a parsed [`osprey_ast::Program`] and
//! returns the list of [`TypeError`]s (empty ⇒ well-typed).

mod builtins;
mod check;
mod convert;
mod ctx;
mod env;
mod error;
mod expr;
mod pattern;
mod ty;
mod unify;

pub use check::check_program;
pub use error::TypeError;
pub use ty::{names, Scheme, Type, VarId};

#[cfg(test)]
mod tests {
    use super::*;
    use osprey_syntax::parse_program;

    /// Parse + type-check a snippet, asserting it is well-typed.
    fn ok(src: &str) {
        let parsed = parse_program(src);
        assert!(
            parsed.errors.is_empty(),
            "syntax errors: {:?}",
            parsed.errors
        );
        let errs = check_program(&parsed.program);
        assert!(errs.is_empty(), "unexpected type errors: {errs:?}");
    }

    /// Parse + type-check, asserting at least one type error is reported.
    fn bad(src: &str) -> Vec<TypeError> {
        let parsed = parse_program(src);
        assert!(
            parsed.errors.is_empty(),
            "syntax errors: {:?}",
            parsed.errors
        );
        let errs = check_program(&parsed.program);
        assert!(!errs.is_empty(), "expected a type error, got none");
        errs
    }

    #[test]
    fn checks_arithmetic_and_let() {
        ok("fn inc(x: int) -> int = x + 1\nlet y = inc(41)\n");
    }

    #[test]
    fn string_concatenation_infers_string() {
        ok("fn greet(name: string) -> string = \"hi \" + name\n");
    }

    #[test]
    fn lambda_param_is_inferred_from_use() {
        // `s` has no annotation; `s + \"!\"` forces it to string.
        ok("let exclaim = fn(s) => s + \"!\"\nlet r = exclaim(\"hi\")\n");
    }

    #[test]
    fn records_field_access_and_update() {
        ok("type Point = { x: int, y: int }\n\
            let p = Point { x: 1, y: 2 }\n\
            let q = p { x: 10 }\n\
            fn px(pt: Point) -> int = pt.x\n");
    }

    #[test]
    fn result_pattern_binds_payload_type() {
        ok("fn unwrap(r: Result<int, Error>) -> int = match r {\n\
              Success { value } => value\n\
              Error { message } => 0\n\
            }\n");
    }

    #[test]
    fn generic_union_flows_type_argument() {
        ok("type Box<T> = Empty | Full { value: T }\n\
            let b = Full { value: 7 }\n\
            let s = match b {\n\
              Full { value } => toString(value)\n\
              Empty => \"empty\"\n\
            }\n");
    }

    #[test]
    fn higher_order_function_application() {
        ok(
            "fn applyFn(value: int, func: (int) -> int) -> int = func(value)\n\
            fn double(x: int) -> int = x * 2\n\
            let r = applyFn(value: 10, func: double)\n",
        );
    }

    #[test]
    fn reports_type_mismatch_in_call() {
        bad("fn inc(x: int) -> int = x + 1\nlet r = inc(\"not an int\")\n");
    }

    #[test]
    fn reports_non_exhaustive_bool_match() {
        let errs = bad("fn f(b: bool) -> int = match b { true => 1 }\n");
        assert!(errs.iter().any(|e| e.message.contains("non-exhaustive")));
    }

    #[test]
    fn reports_unknown_identifier() {
        let errs = bad("let x = totallyUndefinedThing\n");
        assert!(errs
            .iter()
            .any(|e| e.message.contains("unknown identifier")));
    }
}
