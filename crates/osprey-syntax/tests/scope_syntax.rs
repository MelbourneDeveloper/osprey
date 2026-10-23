//! `finally`, `mask` and `instance` parse to the same AST in both flavors.
//! Implements [EFFECTS-FINALIZATION], [EFFECTS-MASK], [EFFECTS-INSTANCE] and
//! [FLAVOR-ML-EFFECT-SCOPE].

use osprey_syntax::{parse_program_with_flavor, Flavor};

fn canonical(source: &str, flavor: Flavor) -> String {
    let parsed = parse_program_with_flavor(source, flavor);
    assert!(
        parsed.errors.is_empty(),
        "{flavor} syntax errors: {:?}",
        parsed.errors
    );
    osprey_ast::canonical::without_positions(&parsed.program)
}

fn rejection(source: &str, flavor: Flavor) -> String {
    let parsed = parse_program_with_flavor(source, flavor);
    parsed
        .errors
        .iter()
        .map(|error| error.message.clone())
        .collect::<Vec<_>>()
        .join("\n")
}

const DEFAULT: &str = concat!(
    "effect Log { info: fn(string) -> Unit }\n",
    "let audit = instance Log\n",
    "let answer = {\n",
    "    handle Log {\n",
    "        info m => print(m)\n",
    "        finally => print(\"closed\")\n",
    "    }\n",
    "    mask Log { perform audit.info(\"one\") }\n",
    "}\n",
);

const ML: &str = concat!(
    "effect Log\n    info : string => Unit\n",
    "audit = instance Log\n",
    "answer =\n",
    "    handle Log\n",
    "        info m => print m\n",
    "        finally => print \"closed\"\n",
    "    mask Log\n",
    "        perform audit.info \"one\"\n",
);

#[test]
fn scope_forms_lower_to_the_same_ast_in_both_flavors() {
    let default = canonical(DEFAULT, Flavor::Default);
    assert!(default.contains("Mask"), "{default}");
    assert!(default.contains("Instance"), "{default}");
    assert_eq!(default, canonical(ML, Flavor::Ml));
}

#[test]
fn an_ml_finally_clause_takes_no_parameter() {
    let source = ML.replace("finally =>", "finally x =>");
    assert!(rejection(&source, Flavor::Ml).contains("one finally clause with no parameter"));
}

#[test]
fn an_ml_mask_needs_an_indented_body() {
    let source = "effect Log\n    info : string => Unit\nanswer = mask Log\n";
    assert!(rejection(source, Flavor::Ml).contains("expected an indented body after `mask Effect`"));
}
