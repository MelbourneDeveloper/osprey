//! Cross-flavor equivalence ([FLAVOR-TEST] / [FLAVOR-CURRY],
//! docs/specs/0023-LanguageFlavors.md): ML is an **uncurried syntactic skin**,
//! so a multi-parameter ML binding `add x y = …` lowers to the SAME canonical
//! AST (modulo source positions) as the Default *multi-parameter* `fn add(x, y)`
//! — and it must NOT collapse to the Default *explicit-curry* `fn add(x) =
//! fn(y) => …`. That distinction is the currying boundary the flavors meet at:
//! ML curries only through explicit lambdas, exactly as Default does.

use osprey_syntax::{parse_program_with_flavor, Flavor};

/// The canonical AST as a debug string with every source `Position { … }`
/// payload scrubbed, so structural equality ignores spans.
fn canonical(src: &str, flavor: Flavor) -> String {
    let parsed = parse_program_with_flavor(src, flavor);
    assert!(
        parsed.errors.is_empty(),
        "unexpected {flavor} syntax errors: {:?}",
        parsed.errors
    );
    scrub_positions(&format!("{:?}", parsed.program))
}

/// Drop every `Position { line: N, column: M }` from a debug string. `Position`
/// has no nested braces, so the next `}` always closes it.
fn scrub_positions(debug: &str) -> String {
    let mut out = String::with_capacity(debug.len());
    let mut rest = debug;
    while let Some(idx) = rest.find("Position {") {
        out.push_str(&rest[..idx]);
        rest = &rest[idx..];
        match rest.find('}') {
            Some(close) => rest = &rest[close + 1..],
            None => break,
        }
    }
    out.push_str(rest);
    out
}

#[test]
fn ml_multiparam_equals_default_multiparam() {
    // Both lower to one two-parameter Function { params: [x, y], body: x + y } —
    // ML is an uncurried skin, so `add x y = …` is the Default `fn add(x, y)`.
    let default_multi = "fn add(x, y) = x + y\n";
    let ml_multi = "add x y = x + y\n";
    assert_eq!(
        canonical(default_multi, Flavor::Default),
        canonical(ml_multi, Flavor::Ml),
        "ML multi-param must equal Default multi-param at the canonical AST"
    );
}

#[test]
fn ml_multiparam_differs_from_default_explicit_curry() {
    // Default `fn add(x) = fn(y) => …` is a one-parameter Function returning a
    // Lambda — a different node than the ML uncurried two-parameter Function.
    // ML never auto-curries; explicit currying needs an explicit lambda.
    let default_curry = "fn add(x) = fn(y) => x + y\n";
    let ml_multi = "add x y = x + y\n";
    assert_ne!(
        canonical(default_curry, Flavor::Default),
        canonical(ml_multi, Flavor::Ml),
        "Default explicit-curry must NOT equal ML uncurried multi-param"
    );
}
