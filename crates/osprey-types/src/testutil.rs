//! Shared `#[cfg(test)]` helpers for the type-checker unit tests.
//!
//! Every test module parses a snippet (asserting it is syntactically valid),
//! runs [`check_program`](crate::check::check_program), and asserts on the
//! resulting diagnostics. These helpers hoist that boilerplate so each module
//! re-uses one canonical copy via `use crate::testutil::*;`.

use crate::check::check_program;
use crate::error::TypeError;
use osprey_syntax::parse_program;

/// Parse + type-check a snippet, returning the diagnostics. Panics if the
/// snippet has syntax errors, since those would mask the type-checking intent.
pub(crate) fn check(src: &str) -> Vec<TypeError> {
    let parsed = parse_program(src);
    assert!(
        parsed.errors.is_empty(),
        "syntax errors: {:?}",
        parsed.errors
    );
    check_program(&parsed.program)
}

/// Parse + type-check, asserting the snippet is well-typed (no diagnostics).
pub(crate) fn ok(src: &str) {
    let errs = check(src);
    assert!(errs.is_empty(), "unexpected type errors: {errs:?}");
}

/// Parse + type-check, asserting at least one type error is reported.
pub(crate) fn bad(src: &str) -> Vec<TypeError> {
    let errs = check(src);
    assert!(!errs.is_empty(), "expected a type error, got none");
    errs
}
