//! `[SYMEX-CODEGEN]`: the module a symbolic-execution engine explores.

use super::*;

const DYNAMIC_EFFECT: &str = "effect Wide { control op: fn(int, string) -> int }\n\
     fn body() -> int !Wide = perform Wide.op(1, \"x\")\n\
     fn main() -> int { let r = {\n\
         handle Wide {\n\
             op a b => resume(a)\n\
         }\n\
         body()\n\
     }\n\
       print(\"r=${toString(r)}\")\n\
       0 }\n";

fn symbolic_module(src: &str) -> String {
    let parsed = parse_program(src);
    assert!(
        parsed.errors.is_empty(),
        "syntax errors: {:?}",
        parsed.errors
    );
    compile_program_symbolic(
        &parsed.program,
        DebugSource {
            filename: "hunt.osp".to_string(),
            directory: "/tmp".to_string(),
        },
    )
    .expect("symbolic codegen should succeed")
}

#[test]
fn klee_mode_aborts_on_the_unhandled_effect_guard() {
    let plain = module(DYNAMIC_EFFECT);
    let symbolic = symbolic_module(DYNAMIC_EFFECT);
    assert!(
        plain.contains("call void @exit(i32 1)") && !plain.contains("@abort"),
        "a native build keeps exit on the guard path"
    );
    assert!(
        symbolic.contains("call void @abort()") && !symbolic.contains("call void @exit"),
        "a symbolic build aborts on the guard path"
    );
}

#[test]
fn klee_mode_carries_debug_metadata_without_debug_flag() {
    let symbolic = symbolic_module("print(\"hello\")\n");
    assert!(
        symbolic.contains("!DICompileUnit") && symbolic.contains("!DILocation"),
        "{symbolic}"
    );
    assert!(!symbolic.contains("optnone"));
}
