//! The shared, flavor-neutral line scanner.
//!
//! Both formatters reindent a file line-by-line. To do that safely they must
//! understand *where the code is*: string literals (including `${…}`
//! interpolation, which may itself contain nested strings and braces) and `//`
//! line comments are copied verbatim, while runs of interior whitespace between
//! real tokens collapse to a single space. The scanner never inserts or deletes
//! a token — it only normalises whitespace and measures bracket nesting — so the
//! reparse guard in [`crate::format_source`] can prove the result is equivalent.

/// One physical source line, scanned into the facts a reindenter needs.
#[derive(Debug, Clone)]
pub(crate) struct Line {
    /// Width of the original leading whitespace (spaces and tabs each count 1).
    pub leading_ws: usize,
    /// The line's content with the leading indent removed, interior whitespace
    /// collapsed to single spaces, and trailing whitespace trimmed.
    pub content: String,
    /// Net change in `{}`/`()`/`[]` nesting over the line (opens minus closes),
    /// counting only brackets outside strings and comments.
    pub open_delta: i32,
    /// How many closing brackets the content opens with (a line that *starts* by
    /// closing a block dedents itself by this many levels).
    pub leading_closers: i32,
}

impl Line {
    /// Whether the line carries no code (empty after normalisation).
    pub(crate) fn is_blank(&self) -> bool {
        self.content.is_empty()
    }

    /// Whether the line is nothing but a `//` comment.
    pub(crate) fn is_comment_only(&self) -> bool {
        self.content.starts_with("//")
    }
}

/// Scan one raw line (a trailing `\r` is treated as part of the line ending).
pub(crate) fn scan_line(raw: &str) -> Line {
    let chars: Vec<char> = raw.trim_end_matches('\r').chars().collect();
    let leading_ws = chars.iter().take_while(|c| **c == ' ' || **c == '\t').count();
    let (content, open_delta) = normalize(&chars[leading_ws..]);
    let leading_closers = content
        .chars()
        .take_while(|c| matches!(c, '}' | ')' | ']'))
        .count();
    Line {
        leading_ws,
        content,
        open_delta,
        leading_closers: i32::try_from(leading_closers).unwrap_or(0),
    }
}

/// Normalise the post-indent part of a line: collapse interior whitespace,
/// trim the end, copy strings/comments verbatim, and tally bracket nesting.
fn normalize(chars: &[char]) -> (String, i32) {
    let mut out = String::with_capacity(chars.len());
    let mut delta = 0i32;
    let mut i = 0;
    while let Some(&c) = chars.get(i) {
        if c == '/' && chars.get(i + 1) == Some(&'/') {
            out.extend(&chars[i..]);
            break;
        }
        if c == '"' {
            i = copy_string(chars, i, &mut out);
        } else if c == ' ' || c == '\t' {
            out.push(' ');
            i += chars[i..].iter().take_while(|x| **x == ' ' || **x == '\t').count();
        } else {
            delta += bracket_delta(c);
            out.push(c);
            i += 1;
        }
    }
    (out.trim_end().to_string(), delta)
}

/// The nesting contribution of a single character: `+1` to open a block-forming
/// bracket, `-1` to close one, `0` otherwise. Angle brackets are deliberately
/// excluded — `<`/`>` are comparison operators far more often than generics.
fn bracket_delta(c: char) -> i32 {
    match c {
        '{' | '(' | '[' => 1,
        '}' | ')' | ']' => -1,
        _ => 0,
    }
}

/// Copy a string literal verbatim starting at the opening quote `chars[start]`,
/// returning the index just past the closing quote (or end of line if
/// unterminated). Escapes are preserved and `${…}` interpolation is delegated.
fn copy_string(chars: &[char], start: usize, out: &mut String) -> usize {
    out.push('"');
    let mut i = start + 1;
    while let Some(&c) = chars.get(i) {
        match c {
            '\\' => {
                out.push('\\');
                if let Some(&next) = chars.get(i + 1) {
                    out.push(next);
                }
                i += 2;
            }
            '"' => {
                out.push('"');
                return i + 1;
            }
            '$' if chars.get(i + 1) == Some(&'{') => i = copy_interpolation(chars, i, out),
            other => {
                out.push(other);
                i += 1;
            }
        }
    }
    i
}

/// Copy a `${…}` interpolation verbatim starting at the `$` (`chars[start]`),
/// honouring nested braces and nested strings so the matching `}` is found even
/// inside an embedded string. Returns the index just past that `}`.
fn copy_interpolation(chars: &[char], start: usize, out: &mut String) -> usize {
    out.push('$');
    out.push('{');
    let mut i = start + 2;
    let mut depth = 1i32;
    while depth > 0 {
        match chars.get(i) {
            None => break,
            Some('"') => {
                i = copy_string(chars, i, out);
            }
            Some(&c) => {
                depth += bracket_delta_brace(c);
                out.push(c);
                i += 1;
            }
        }
    }
    i
}

/// Like [`bracket_delta`] but only braces matter when tracking interpolation
/// depth — parentheses and indexing inside `${…}` never close the interpolation.
fn bracket_delta_brace(c: char) -> i32 {
    match c {
        '{' => 1,
        '}' => -1,
        _ => 0,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn collapses_interior_whitespace_and_trims_ends() {
        let line = scan_line("    let   x =    1   ");
        assert_eq!(line.leading_ws, 4);
        assert_eq!(line.content, "let x = 1");
        assert_eq!(line.open_delta, 0);
    }

    #[test]
    fn counts_block_brackets_outside_strings() {
        let open = scan_line("fn main() = {");
        assert_eq!(open.open_delta, 1);
        assert_eq!(open.leading_closers, 0);
        let close = scan_line("}");
        assert_eq!(close.open_delta, -1);
        assert_eq!(close.leading_closers, 1);
    }

    #[test]
    fn strings_and_comments_are_copied_verbatim() {
        // The `{`/`}` and `//` inside the string must not move brackets, and the
        // double spaces inside the literal must survive.
        let line = scan_line(r#"print("a  {b}  // not a comment")"#);
        assert_eq!(line.content, r#"print("a  {b}  // not a comment")"#);
        assert_eq!(line.open_delta, 0);
    }

    #[test]
    fn interpolation_with_nested_string_and_braces_is_balanced() {
        let line = scan_line(r#"print("v=${ f("}") } end")"#);
        // The `}` inside the nested string does not close the interpolation, and
        // the whole literal is preserved untouched.
        assert_eq!(line.content, r#"print("v=${ f("}") } end")"#);
        assert_eq!(line.open_delta, 0);
    }

    #[test]
    fn trailing_comment_keeps_its_spacing() {
        let line = scan_line("x = 1   // note   with   spaces");
        assert_eq!(line.content, "x = 1 // note   with   spaces");
    }

    #[test]
    fn blank_and_comment_classification() {
        assert!(scan_line("    ").is_blank());
        assert!(scan_line("   // hi").is_comment_only());
        assert!(!scan_line("x // hi").is_comment_only());
    }
}
