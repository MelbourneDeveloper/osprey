//! Position-aware text helpers shared by hover, definition and references.
//!
//! LSP positions are `(line, character)` with `character` counted in the
//! negotiated [`PositionEncoding`]. Osprey identifiers are ASCII, but these
//! helpers honour the encoding so multi-byte text still lines up.

use lspkit_vfs::PositionEncoding;

/// An identifier found under a cursor, with its span in the line.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct WordSpan {
    /// The identifier text.
    pub word: String,
    /// Start character offset within the line (negotiated encoding).
    pub start: u32,
    /// End character offset within the line (negotiated encoding).
    pub end: u32,
}

/// A whole-word occurrence of an identifier within a document.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Occurrence {
    /// Zero-based line.
    pub line: u32,
    /// Start character offset within the line (negotiated encoding).
    pub start: u32,
    /// End character offset within the line (negotiated encoding).
    pub end: u32,
}

fn is_ident(c: char) -> bool {
    c.is_alphanumeric() || c == '_'
}

/// Width of `c` in `encoding`'s character units.
#[must_use]
pub fn char_width(c: char, encoding: PositionEncoding) -> u32 {
    match encoding {
        PositionEncoding::Utf8 => u32::try_from(c.len_utf8()).unwrap_or(1),
        PositionEncoding::Utf16 => u32::try_from(c.len_utf16()).unwrap_or(1),
        // UTF-32 counts code points; `PositionEncoding` is `#[non_exhaustive]`,
        // so anything new also counts as one unit per `char`.
        _ => 1,
    }
}

/// Length of `s` in `encoding`'s character units.
#[must_use]
pub fn measure(s: &str, encoding: PositionEncoding) -> u32 {
    s.chars().map(|c| char_width(c, encoding)).sum()
}

/// The prefix of `line` up to (but excluding) character offset `character`.
#[must_use]
pub fn prefix_to(line: &str, character: u32, encoding: PositionEncoding) -> &str {
    let mut offset = 0u32;
    for (idx, c) in line.char_indices() {
        if offset >= character {
            return line.get(..idx).unwrap_or(line);
        }
        offset = offset.saturating_add(char_width(c, encoding));
    }
    line
}

/// The identifier under `character` within `line`, or `None` if the cursor is
/// not over an identifier character.
#[must_use]
pub fn word_at(line: &str, character: u32, encoding: PositionEncoding) -> Option<WordSpan> {
    let mut offset = 0u32;
    let mut run_start = 0u32;
    let mut run = String::new();
    let mut found: Option<WordSpan> = None;
    for c in line.chars() {
        let w = char_width(c, encoding);
        if is_ident(c) {
            if run.is_empty() {
                run_start = offset;
            }
            run.push(c);
        } else {
            found = found.or_else(|| take_if_covers(&run, run_start, offset, character));
            run.clear();
        }
        offset = offset.saturating_add(w);
    }
    found.or_else(|| take_if_covers(&run, run_start, offset, character))
}

/// Promote an accumulated identifier run to a [`WordSpan`] when `character`
/// falls within `[start, end]` (inclusive of the trailing edge so a cursor
/// resting just after the word still resolves it).
fn take_if_covers(run: &str, start: u32, end: u32, character: u32) -> Option<WordSpan> {
    if run.is_empty() || character < start || character > end {
        return None;
    }
    Some(WordSpan {
        word: run.to_owned(),
        start,
        end,
    })
}

/// Every whole-word occurrence of `name` across `text`.
#[must_use]
pub fn occurrences(text: &str, name: &str, encoding: PositionEncoding) -> Vec<Occurrence> {
    text.lines()
        .enumerate()
        .flat_map(|(idx, line)| line_occurrences(line, idx, name, encoding))
        .collect()
}

fn line_occurrences(
    line: &str,
    line_idx: usize,
    name: &str,
    enc: PositionEncoding,
) -> Vec<Occurrence> {
    let line_no = u32::try_from(line_idx).unwrap_or(u32::MAX);
    let mut offset = 0u32;
    let mut run_start = 0u32;
    let mut run = String::new();
    let mut out = Vec::new();
    for c in line.chars() {
        let w = char_width(c, enc);
        if is_ident(c) {
            if run.is_empty() {
                run_start = offset;
            }
            run.push(c);
        } else {
            push_match(&mut out, &run, run_start, offset, line_no, name);
            run.clear();
        }
        offset = offset.saturating_add(w);
    }
    push_match(&mut out, &run, run_start, offset, line_no, name);
    out
}

fn push_match(out: &mut Vec<Occurrence>, run: &str, start: u32, end: u32, line: u32, name: &str) {
    if run == name {
        out.push(Occurrence { line, start, end });
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    const U16: PositionEncoding = PositionEncoding::Utf16;

    #[test]
    fn word_at_finds_identifier_under_and_after_cursor() {
        let line = "let total = a + b";
        assert_eq!(
            word_at(line, 5, U16).map(|w| w.word),
            Some("total".to_owned())
        );
        // Trailing edge of the word still resolves it.
        assert_eq!(
            word_at(line, 9, U16).map(|w| w.word),
            Some("total".to_owned())
        );
        assert_eq!(word_at(line, 16, U16).map(|w| w.word), Some("b".to_owned()));
    }

    #[test]
    fn word_at_returns_none_off_identifier() {
        assert_eq!(word_at("a + b", 2, U16), None);
        assert_eq!(word_at("", 0, U16), None);
    }

    #[test]
    fn occurrences_are_whole_word_only() {
        let text = "fn add(a) = a\nlet adder = add(adding)\n";
        let hits = occurrences(text, "add", U16);
        let lines: Vec<u32> = hits.iter().map(|h| h.line).collect();
        assert_eq!(lines, vec![0, 1], "{hits:?}");
    }
}
