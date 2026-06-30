//! ML-flavor tokens. The lexer ([`super::lexer`]) emits a flat stream of these,
//! including the layout markers [`TokKind::Indent`], [`TokKind::Dedent`], and
//! [`TokKind::Newline`] derived from the offside rule ([FLAVOR-ML-LAYOUT]).

use osprey_ast::Position;

/// A lexed token with its source position.
#[derive(Debug, Clone, PartialEq)]
pub(crate) struct Token {
    /// What kind of token this is, with any payload.
    pub kind: TokKind,
    /// 1-based line / 0-based column where the token starts.
    pub pos: Position,
    /// Whether this token immediately follows the previous content token with
    /// no intervening whitespace/comment. Disambiguates `xs[0]` (a *glued*
    /// postfix index) from `f [0]` (whitespace application to a list literal) —
    /// the only place ML whitespace-application overlaps bracket syntax
    /// ([FLAVOR-ML-INDEX]).
    pub glued: bool,
}

/// The kind (and payload) of an ML token.
#[derive(Debug, Clone, PartialEq)]
pub(crate) enum TokKind {
    /// Integer literal.
    Int(i64),
    /// Float literal.
    Float(f64),
    /// String literal body (raw, including `${...}` spans, escapes unresolved).
    Str(String),
    /// Identifier or keyword-as-name (lowercase var / uppercase constructor).
    Ident(String),
    /// `mut`.
    KwMut,
    /// `true`.
    KwTrue,
    /// `false`.
    KwFalse,
    /// `match`.
    KwMatch,
    /// `type` — introduces a union/enum/record type declaration ([FLAVOR-ML-TYPE]).
    KwType,
    /// `extern` — introduces an external (FFI) function declaration ([FLAVOR-ML-EXTERN]).
    KwExtern,
    /// `spawn` — starts a fiber, evaluating its block/expr concurrently ([FLAVOR-ML-SPAWN]).
    KwSpawn,
    /// A reserved word reserved for a not-yet-implemented construct (`effect`,
    /// `handler`, `handle`, `do`, `perform`). Carries its spelling so the
    /// parser can report a precise "not yet supported" diagnostic.
    Reserved(String),
    /// `=`.
    Eq,
    /// `:=`.
    ColonEq,
    /// `:`.
    Colon,
    /// `->`.
    Arrow,
    /// `=>`.
    FatArrow,
    /// `\` (lambda head).
    Backslash,
    /// `(`.
    LParen,
    /// `)`.
    RParen,
    /// `[`.
    LBracket,
    /// `]`.
    RBracket,
    /// `,`.
    Comma,
    /// `.`.
    Dot,
    /// A binary/unary operator spelled exactly as it lowers (`+`, `==`, `&&`, …).
    Op(String),
    /// Significant end-of-line within a layout region.
    Newline,
    /// Start of a more-indented region.
    Indent,
    /// Return to a less-indented region.
    Dedent,
    /// End of input.
    Eof,
}

/// Map a bare identifier spelling to its keyword/reserved kind, or treat it as
/// an ordinary identifier.
pub(crate) fn keyword_or_ident(text: &str) -> TokKind {
    match text {
        "mut" => TokKind::KwMut,
        "true" => TokKind::KwTrue,
        "false" => TokKind::KwFalse,
        "match" => TokKind::KwMatch,
        "type" => TokKind::KwType,
        "extern" => TokKind::KwExtern,
        "spawn" => TokKind::KwSpawn,
        "effect" | "handler" | "handle" | "do" | "perform" => {
            TokKind::Reserved(text.to_owned())
        }
        _ => TokKind::Ident(text.to_owned()),
    }
}
