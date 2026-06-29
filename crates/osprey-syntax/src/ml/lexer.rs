//! The ML-flavor lexer: a hand-written scanner that turns source text into a
//! flat [`Token`] stream, then derives the layout markers (`Indent`, `Dedent`,
//! `Newline`) from the offside rule ([FLAVOR-ML-LAYOUT]).
//!
//! Two phases keep each piece small and testable: [`scan`] produces content
//! tokens with positions (no layout), and [`insert_layout`] walks those tokens
//! and inserts the layout markers from each line's first-token column, with
//! bracket depth suppressing layout inside `( … )`.
//!
//! ESCAPE HATCH: if this hand-written layout frontend becomes onerous or
//! accrues parsing bugs we cannot tame, we fall back to a `tree-sitter-osprey-ml`
//! grammar with an external INDENT/DEDENT/NEWLINE scanner.c — the boundary law
//! makes the parser mechanism a flavor-internal swap (docs/specs/0023).

use super::token::{keyword_or_ident, TokKind, Token};
use crate::SyntaxError;
use osprey_ast::Position;

/// Lex `source` into a layout-resolved token stream terminated by
/// [`TokKind::Eof`], plus any lexical errors.
pub(crate) fn lex(source: &str) -> (Vec<Token>, Vec<SyntaxError>) {
    let mut scanner = Scanner::new(source);
    let (content, mut errors) = scanner.scan();
    let (tokens, layout_errors) = insert_layout(content);
    errors.extend(layout_errors);
    (tokens, errors)
}

/// Phase-1 scanner over the raw characters.
struct Scanner {
    chars: Vec<char>,
    i: usize,
    line: u32,
    col: u32,
    errors: Vec<SyntaxError>,
}

impl Scanner {
    fn new(source: &str) -> Self {
        Scanner {
            chars: source.chars().collect(),
            i: 0,
            line: 1,
            col: 0,
            errors: Vec::new(),
        }
    }

    fn pos(&self) -> Position {
        Position {
            line: self.line,
            column: self.col,
        }
    }

    fn peek(&self, ahead: usize) -> Option<char> {
        self.chars.get(self.i + ahead).copied()
    }

    fn bump(&mut self) -> Option<char> {
        let c = self.chars.get(self.i).copied()?;
        self.i += 1;
        if c == '\n' {
            self.line += 1;
            self.col = 0;
        } else {
            self.col += 1;
        }
        Some(c)
    }

    fn error(&mut self, pos: Position, message: impl Into<String>) {
        self.errors.push(SyntaxError {
            message: message.into(),
            position: pos,
        });
    }

    /// Skip inline whitespace, newlines (layout is position-derived later), and
    /// `// …` line comments.
    fn skip_trivia(&mut self) {
        while let Some(c) = self.peek(0) {
            match c {
                ' ' | '\t' | '\r' | '\n' => {
                    self.bump();
                }
                '/' if self.peek(1) == Some('/') => {
                    while !matches!(self.peek(0), Some('\n') | None) {
                        self.bump();
                    }
                }
                _ => break,
            }
        }
    }

    fn scan(&mut self) -> (Vec<Token>, Vec<SyntaxError>) {
        let mut out = Vec::new();
        loop {
            self.skip_trivia();
            if self.i >= self.chars.len() {
                break;
            }
            let pos = self.pos();
            if let Some(kind) = self.scan_token(pos) {
                out.push(Token { kind, pos });
            }
        }
        (out, std::mem::take(&mut self.errors))
    }

    fn scan_token(&mut self, pos: Position) -> Option<TokKind> {
        let c = self.peek(0)?;
        match c {
            '0'..='9' => Some(self.scan_number(pos)),
            '"' => Some(self.scan_string(pos)),
            c if c.is_alphabetic() || c == '_' => Some(self.scan_ident()),
            _ => self.scan_operator(pos),
        }
    }

    fn scan_number(&mut self, pos: Position) -> TokKind {
        let start = self.i;
        while matches!(self.peek(0), Some('0'..='9')) {
            self.bump();
        }
        let is_float = self.peek(0) == Some('.') && matches!(self.peek(1), Some('0'..='9'));
        if is_float {
            self.bump();
            while matches!(self.peek(0), Some('0'..='9')) {
                self.bump();
            }
        }
        let text: String = self.chars[start..self.i].iter().collect();
        if is_float {
            text.parse::<f64>().map_or_else(
                |_| {
                    self.error(pos, format!("invalid float literal '{text}'"));
                    TokKind::Float(0.0)
                },
                TokKind::Float,
            )
        } else {
            text.parse::<i64>().map_or_else(
                |_| {
                    self.error(pos, format!("invalid integer literal '{text}'"));
                    TokKind::Int(0)
                },
                TokKind::Int,
            )
        }
    }

    fn scan_string(&mut self, pos: Position) -> TokKind {
        self.bump(); // opening quote
        let mut raw = String::new();
        loop {
            match self.peek(0) {
                None | Some('\n') => {
                    self.error(pos, "unterminated string literal");
                    break;
                }
                Some('"') => {
                    self.bump();
                    break;
                }
                Some('\\') => {
                    self.bump();
                    if let Some(escaped) = self.bump() {
                        raw.push('\\');
                        raw.push(escaped);
                    }
                }
                Some(c) => {
                    self.bump();
                    raw.push(c);
                }
            }
        }
        TokKind::Str(raw)
    }

    fn scan_ident(&mut self) -> TokKind {
        let start = self.i;
        while matches!(self.peek(0), Some(c) if c.is_alphanumeric() || c == '_') {
            self.bump();
        }
        let text: String = self.chars[start..self.i].iter().collect();
        keyword_or_ident(&text)
    }

    fn scan_operator(&mut self, pos: Position) -> Option<TokKind> {
        let c = self.peek(0)?;
        let next = self.peek(1);
        if let Some(kind) = two_char_operator(c, next) {
            self.bump();
            self.bump();
            return Some(kind);
        }
        let kind = single_char_operator(c);
        self.bump();
        match kind {
            Some(kind) => Some(kind),
            None => {
                self.error(pos, format!("unexpected character '{c}'"));
                None
            }
        }
    }
}

/// Match a two-character operator/punctuation lexeme.
fn two_char_operator(c: char, next: Option<char>) -> Option<TokKind> {
    let next = next?;
    let kind = match (c, next) {
        (':', '=') => TokKind::ColonEq,
        ('-', '>') => TokKind::Arrow,
        ('=', '>') => TokKind::FatArrow,
        ('=', '=') => TokKind::Op("==".to_owned()),
        ('!', '=') => TokKind::Op("!=".to_owned()),
        ('<', '=') => TokKind::Op("<=".to_owned()),
        ('>', '=') => TokKind::Op(">=".to_owned()),
        ('&', '&') => TokKind::Op("&&".to_owned()),
        ('|', '|') => TokKind::Op("||".to_owned()),
        ('|', '>') => TokKind::Op("|>".to_owned()),
        _ => return None,
    };
    Some(kind)
}

/// Match a single-character operator/punctuation lexeme.
fn single_char_operator(c: char) -> Option<TokKind> {
    let kind = match c {
        '=' => TokKind::Eq,
        ':' => TokKind::Colon,
        '\\' => TokKind::Backslash,
        '(' => TokKind::LParen,
        ')' => TokKind::RParen,
        ',' => TokKind::Comma,
        '.' => TokKind::Dot,
        '+' | '-' | '*' | '/' | '%' | '<' | '>' | '!' => TokKind::Op(c.to_string()),
        _ => return None,
    };
    Some(kind)
}

/// Phase 2: insert `Indent`/`Dedent`/`Newline` from each line's first-token
/// column. Layout is suppressed while bracket depth is non-zero so a
/// parenthesised expression may span lines. Implements [FLAVOR-ML-LAYOUT].
fn insert_layout(content: Vec<Token>) -> (Vec<Token>, Vec<SyntaxError>) {
    let mut out = Vec::new();
    let mut errors = Vec::new();
    let mut stack = vec![0u32];
    let mut depth = 0i32;
    let mut prev_line = 0u32;
    let mut started = false;
    for tok in content {
        if depth == 0 && tok.pos.line != prev_line {
            emit_layout(&mut out, &mut stack, &mut errors, tok.pos, started);
        }
        match tok.kind {
            TokKind::LParen => depth += 1,
            TokKind::RParen => depth = (depth - 1).max(0),
            _ => {}
        }
        prev_line = tok.pos.line;
        out.push(tok);
        started = true;
    }
    let close = out.last().map_or(Position::default(), |t| t.pos);
    while stack.last().copied().unwrap_or(0) > 0 {
        stack.pop();
        out.push(layout_tok(TokKind::Dedent, close));
    }
    out.push(layout_tok(TokKind::Eof, close));
    (out, errors)
}

/// Compare a logical line's indentation against the stack, pushing one `Indent`,
/// a run of `Dedent`s, or a separating `Newline`.
fn emit_layout(
    out: &mut Vec<Token>,
    stack: &mut Vec<u32>,
    errors: &mut Vec<SyntaxError>,
    pos: Position,
    started: bool,
) {
    let col = pos.column;
    let top = stack.last().copied().unwrap_or(0);
    if col > top {
        stack.push(col);
        out.push(layout_tok(TokKind::Indent, pos));
        return;
    }
    while col < stack.last().copied().unwrap_or(0) {
        stack.pop();
        out.push(layout_tok(TokKind::Dedent, pos));
    }
    if col == stack.last().copied().unwrap_or(0) {
        if started {
            out.push(layout_tok(TokKind::Newline, pos));
        }
    } else {
        errors.push(SyntaxError {
            message: "inconsistent indentation does not match any enclosing block".to_owned(),
            position: pos,
        });
        stack.push(col);
        out.push(layout_tok(TokKind::Indent, pos));
    }
}

fn layout_tok(kind: TokKind, pos: Position) -> Token {
    Token { kind, pos }
}

#[cfg(test)]
#[expect(
    clippy::indexing_slicing,
    reason = "test assertions: an out-of-bounds index is a test failure, not a panic"
)]
mod tests {
    use super::*;

    fn kinds(source: &str) -> Vec<TokKind> {
        let (tokens, errors) = lex(source);
        assert!(errors.is_empty(), "lex errors: {errors:?}");
        tokens.into_iter().map(|t| t.kind).collect()
    }

    #[test]
    fn lexes_binding_with_no_layout() {
        let k = kinds("x = 42\n");
        assert_eq!(
            k,
            vec![
                TokKind::Ident("x".to_owned()),
                TokKind::Eq,
                TokKind::Int(42),
                TokKind::Eof,
            ]
        );
    }

    #[test]
    fn separates_top_level_lines_with_newline() {
        let k = kinds("a = 1\nb = 2\n");
        let newlines = k.iter().filter(|t| **t == TokKind::Newline).count();
        assert_eq!(newlines, 1, "one separator between two top-level bindings");
    }

    #[test]
    fn indents_and_dedents_a_block() {
        let k = kinds("f =\n    g\nh = 1\n");
        assert!(k.contains(&TokKind::Indent), "block opens with Indent");
        assert!(k.contains(&TokKind::Dedent), "block closes with Dedent");
        // The Dedent must precede the sibling `h` binding.
        let dedent = k.iter().position(|t| *t == TokKind::Dedent).unwrap();
        let h = k
            .iter()
            .position(|t| *t == TokKind::Ident("h".to_owned()))
            .unwrap();
        assert!(dedent < h);
    }

    #[test]
    fn suppresses_layout_inside_parentheses() {
        // A line break inside parens must not start a new statement.
        let k = kinds("x = (1 +\n2)\n");
        assert!(!k.contains(&TokKind::Indent), "no layout inside parens");
    }

    #[test]
    fn ignores_blank_and_comment_lines() {
        let k = kinds("a = 1\n\n// note\nb = 2\n");
        let newlines = k.iter().filter(|t| **t == TokKind::Newline).count();
        assert_eq!(newlines, 1, "blank/comment lines are not separators");
    }

    #[test]
    fn lexes_curried_application_and_operators() {
        let k = kinds("r = add 1 2 == 3\n");
        assert!(k.contains(&TokKind::Op("==".to_owned())));
        assert!(k.contains(&TokKind::Int(1)) && k.contains(&TokKind::Int(2)));
    }

    #[test]
    fn reports_unterminated_string() {
        let (_, errors) = lex("x = \"oops\n");
        assert!(errors.iter().any(|e| e.message.contains("unterminated")));
    }
}
