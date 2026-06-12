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
mod closure;
mod collections;
mod conv;
mod effects;
mod error;
mod expr;
mod extern_call;
mod fiber;
mod freevars;
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

/// Every identifier referenced anywhere in `program` — function bodies, lets,
/// nested modules. The CLI's capability sandbox uses this to detect gated
/// builtins (`httpGet`, `readFile`, …) without compiling.
#[must_use]
pub fn referenced_idents(program: &osprey_ast::Program) -> std::collections::BTreeSet<String> {
    let mut out = std::collections::BTreeSet::new();
    for s in &program.statements {
        stmt_idents(s, &mut out);
    }
    out
}

fn stmt_idents(s: &osprey_ast::Stmt, out: &mut std::collections::BTreeSet<String>) {
    use osprey_ast::Stmt;
    match s {
        Stmt::Let { value, .. } | Stmt::Assignment { value, .. } => {
            freevars::free_idents(value, out);
        }
        Stmt::Expr(e) | Stmt::Function { body: e, .. } => freevars::free_idents(e, out),
        Stmt::Module { body, .. } => {
            for inner in body {
                stmt_idents(inner, out);
            }
        }
        _ => {}
    }
}

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
    fn spawn_lowers_to_a_per_instance_closure_cell() {
        // `spawn` lowers its expression as a zero-parameter closure: the thunk
        // takes its heap cell as env (so two in-flight spawns from one site
        // never alias captures) and goes to `fiber_spawn_env`; `await` maps to
        // `fiber_await`. No module globals are involved.
        let ir = module(
            "fn work(n: int) -> int = n * 2\n\
             fn main() -> Unit = {\n\
               let x = 21\n\
               let f = spawn work(x)\n\
               print(\"got ${await(f)}\")\n\
             }\n",
        );
        assert!(ir.contains("call i64 @fiber_spawn_env(i64 (i8*)* @__fiber_thunk_"));
        assert!(ir.contains("define i64 @__fiber_thunk_0(i8* %__env)"));
        assert!(!ir.contains("@__fiber_cap_"));
        assert!(ir.contains("call i64 @fiber_await(i64"));
    }

    #[test]
    fn inline_lambda_argument_becomes_a_closure_cell() {
        // An inline lambda flowing into a function-typed parameter becomes a
        // closure cell `{ fnptr, captures… }`: the emitted function takes a
        // hidden `i8* %__env`, and the indirect call inside `apply` loads the
        // fnptr from the cell and passes the cell back as the env.
        let ir = module(
            "fn apply(value: int, f: (int) -> int) -> int = f(value)\n\
             let r = apply(value: 10, f: fn(x: int) => x + 1)\n\
             print(\"r=${r}\")\n",
        );
        assert!(ir.contains("define i64 @__closure_fn_0(i8* %__env, i64 %x)"));
        assert!(ir.contains("@__closure_cell_0 = private unnamed_addr constant { i8* }"));
        assert!(ir.contains("call i64 %"));
    }

    #[test]
    fn escaping_closure_captures_its_makers_state() {
        // The headline closure case [TYPE-FN-CLOSURE]: a returned lambda
        // capturing its maker's parameter stays callable — the capture is
        // stored in a malloc'd cell and reloaded from `%__env` inside the
        // lifted function.
        let ir = module(
            "fn makeAdder(n: int) -> (int) -> int = fn(x: int) => x + n\n\
             fn main() -> Unit = {\n\
               let add5 = makeAdder(5)\n\
               print(\"r=${add5(3)}\")\n\
             }\n",
        );
        assert!(ir.contains("define i8* @makeAdder(i64 %n)"));
        assert!(ir.contains("bitcast i8* %__env to { i8*, i64 }*"));
        assert!(ir.contains("call i8* @malloc"));
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
        // A construct the backend cannot lower must fail loudly, never
        // silently (CLAUDE.md: no placeholders, fail hard). A method call on a
        // value reaches codegen only through the UFCS rewrite, so a synthetic
        // raw MethodCall node is unsupported.
        let program = osprey_ast::Program {
            statements: vec![osprey_ast::Stmt::Expr(osprey_ast::Expr::MethodCall {
                target: Box::new(osprey_ast::Expr::Integer(1)),
                method: String::from("frobnicate"),
                arguments: Vec::new(),
                named_arguments: Vec::new(),
            })],
        };
        let err = compile_program(&program).unwrap_err();
        assert!(matches!(err, CodegenError::Unsupported(_)));
    }
}
