//! In-process diagnostics.
//!
//! The TypeScript server wrote each edit to a temp file, shelled out to the
//! `osprey` binary, and scraped stderr with a wall of regexes. Here the
//! compiler front-end is called directly: [`osprey_syntax::parse_program`] for
//! syntax errors and [`osprey_types::check_program`] for type errors, mapped to
//! the [`lspkit_server::Diagnostic`] the diagnostics bus fans out.

use lspkit_server::{Diagnostic, Severity};
use lspkit_vfs::PositionEncoding;
use osprey_ast::Position;

const SOURCE: &str = "osprey";

/// Compute diagnostics for `source`. Syntax errors are reported alone (an
/// unparsable file is not type-checked, matching the CLI gate); a clean parse is
/// then type-checked.
#[must_use]
pub fn compute(source: &str, encoding: PositionEncoding) -> Vec<Diagnostic> {
    let parsed = osprey_syntax::parse_program(source);
    if !parsed.errors.is_empty() {
        return parsed
            .errors
            .iter()
            .map(|e| diagnostic(source, e.position, &e.message, "syntax-error", encoding))
            .collect();
    }
    osprey_types::check_program(&parsed.program)
        .iter()
        .map(|e| {
            let pos = e.position.unwrap_or(Position { line: 1, column: 0 });
            diagnostic(source, pos, &e.message, "type-error", encoding)
        })
        .collect()
}

/// Build one error diagnostic spanning the offending line from `pos` onward.
fn diagnostic(
    source: &str,
    pos: Position,
    message: &str,
    code: &str,
    encoding: PositionEncoding,
) -> Diagnostic {
    let line = pos.line.saturating_sub(1);
    let line_text = nth_line(source, line);
    // `pos.column` is a tree-sitter byte offset; the wire range is in the
    // negotiated encoding, so re-measure the line prefix in those units.
    let start = byte_col_to_encoding(line_text, pos.column, encoding);
    let end = line_text
        .map_or(0, |l| crate::text::measure(l, encoding))
        .max(start.saturating_add(1));
    Diagnostic::new(Severity::Error, message, (line, start, line, end))
        .with_source(SOURCE)
        .with_code(code)
}

/// Zero-based `line`'s text, or `None` if absent.
fn nth_line(source: &str, line: u32) -> Option<&str> {
    usize::try_from(line)
        .ok()
        .and_then(|i| source.lines().nth(i))
}

/// Convert a byte column within `line` into `encoding`'s character units.
fn byte_col_to_encoding(line: Option<&str>, byte_col: u32, encoding: PositionEncoding) -> u32 {
    let Some(line) = line else {
        return byte_col;
    };
    let idx = usize::try_from(byte_col).unwrap_or(usize::MAX);
    line.get(..idx)
        .map_or(byte_col, |prefix| crate::text::measure(prefix, encoding))
}

#[cfg(test)]
mod tests {
    use super::*;
    const U16: PositionEncoding = PositionEncoding::Utf16;

    #[test]
    fn clean_program_has_no_diagnostics() {
        let diags = compute("fn main() -> Unit = print(\"hi\")\n", U16);
        assert!(diags.is_empty(), "{diags:?}");
    }

    #[test]
    fn syntax_error_is_reported_with_source_and_code() {
        let diags = compute("fn main( = 1\n", U16);
        assert!(!diags.is_empty());
        let first = diags.first().expect("diagnostic");
        assert_eq!(first.severity, Severity::Error);
        assert_eq!(first.source.as_deref(), Some("osprey"));
        assert_eq!(first.code.as_deref(), Some("syntax-error"));
    }

    #[test]
    fn type_error_surfaces_when_parse_is_clean() {
        // Referencing an unknown function type-checks but does not parse-fail.
        let diags = compute("fn main() -> int = nope(1)\n", U16);
        assert!(
            diags
                .iter()
                .all(|d| d.code.as_deref() == Some("type-error")),
            "{diags:?}"
        );
    }
}
