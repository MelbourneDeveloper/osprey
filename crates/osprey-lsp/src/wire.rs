//! JSON-RPC payload conversion.
//!
//! Incoming LSP params are read straight off `serde_json::Value` (panic-free
//! `.get`/`.as_*` accessors); outgoing results are built with `json!`. Keeping
//! the wire shape here lets the rest of the server speak the neutral
//! [`crate::model`] vocabulary.

use lspkit_server::{Diagnostic, Severity};
use lspkit_vfs::{Position, PositionEncoding, Range, TextEdit};
use serde_json::{json, Value};

use crate::analysis::SymbolInfo;
use crate::model::{CompletionItem, CompletionKind, Location, SignatureInfo, Span};
use crate::text::{measure, occurrences};

// LSP `SymbolKind` numeric codes.
const SYMBOL_CLASS: u8 = 5;
const SYMBOL_FUNCTION: u8 = 12;
const SYMBOL_VARIABLE: u8 = 13;
// LSP `CompletionItemKind` numeric codes.
const COMPLETION_FUNCTION: u8 = 3;
const COMPLETION_VARIABLE: u8 = 6;
const COMPLETION_CLASS: u8 = 7;
const COMPLETION_KEYWORD: u8 = 14;
// LSP `InsertTextFormat`: snippet.
const INSERT_SNIPPET: u8 = 2;

/// The document URI of a request's `textDocument`, if present.
#[must_use]
pub fn doc_uri(params: &Value) -> Option<String> {
    params
        .get("textDocument")
        .and_then(|d| d.get("uri"))
        .and_then(Value::as_str)
        .map(str::to_owned)
}

/// The `(line, character)` of a request's `position`, if present.
#[must_use]
pub fn position(params: &Value) -> Option<(u32, u32)> {
    let pos = params.get("position")?;
    Some((field_u32(pos, "line"), field_u32(pos, "character")))
}

/// Whether a references request asks to include the declaration.
#[must_use]
pub fn include_declaration(params: &Value) -> bool {
    params
        .get("context")
        .and_then(|c| c.get("includeDeclaration"))
        .and_then(Value::as_bool)
        .unwrap_or(false)
}

fn field_u32(value: &Value, key: &str) -> u32 {
    value
        .get(key)
        .and_then(Value::as_u64)
        .and_then(|n| u32::try_from(n).ok())
        .unwrap_or(0)
}

/// The full text of a `textDocument/didOpen`.
#[must_use]
pub fn open_text(params: &Value) -> Option<String> {
    params
        .get("textDocument")
        .and_then(|d| d.get("text"))
        .and_then(Value::as_str)
        .map(str::to_owned)
}

/// The document version of a `didOpen`/`didChange`, defaulting to 0.
#[must_use]
pub fn version(params: &Value) -> i32 {
    params
        .get("textDocument")
        .and_then(|d| d.get("version"))
        .and_then(Value::as_i64)
        .and_then(|n| i32::try_from(n).ok())
        .unwrap_or(0)
}

/// A `didChange` content change: either an incremental edit (`Ok`) or a
/// whole-document replacement (`Err(full_text)`).
#[must_use]
pub fn content_changes(params: &Value) -> Vec<Result<TextEdit, String>> {
    params
        .get("contentChanges")
        .and_then(Value::as_array)
        .map(|items| items.iter().map(change_event).collect())
        .unwrap_or_default()
}

fn change_event(change: &Value) -> Result<TextEdit, String> {
    let text = change
        .get("text")
        .and_then(Value::as_str)
        .unwrap_or_default()
        .to_owned();
    match change.get("range").and_then(range_of) {
        Some(range) => Ok(TextEdit::new(range, text)),
        None => Err(text),
    }
}

fn range_of(range: &Value) -> Option<Range> {
    let start = range.get("start")?;
    let end = range.get("end")?;
    Some(Range::new(
        Position::new(field_u32(start, "line"), field_u32(start, "character")),
        Position::new(field_u32(end, "line"), field_u32(end, "character")),
    ))
}

/// The `initialize` result advertising the server's capabilities.
#[must_use]
pub fn initialize_result(encoding: &str) -> Value {
    json!({
        "capabilities": {
            "positionEncoding": encoding,
            "textDocumentSync": 2,
            "hoverProvider": true,
            "definitionProvider": true,
            "referencesProvider": true,
            "documentSymbolProvider": true,
            "completionProvider": {
                "resolveProvider": false,
                "triggerCharacters": [".", ":", "$", "(", "|"]
            },
            "signatureHelpProvider": { "triggerCharacters": ["(", ","] }
        },
        "serverInfo": { "name": "osprey-lsp" }
    })
}

/// `textDocument/hover` result, or JSON `null`.
#[must_use]
pub fn hover_result(markdown: Option<String>) -> Value {
    markdown.map_or(
        Value::Null,
        |value| json!({ "contents": { "kind": "markdown", "value": value } }),
    )
}

fn range_json(span: Span) -> Value {
    let (sl, sc, el, ec) = span;
    json!({ "start": { "line": sl, "character": sc }, "end": { "line": el, "character": ec } })
}

fn location_json(loc: &Location) -> Value {
    json!({ "uri": loc.uri, "range": range_json(loc.span) })
}

/// `textDocument/definition` / `references` result: an array of locations.
#[must_use]
pub fn locations_result(locations: &[Location]) -> Value {
    Value::Array(locations.iter().map(location_json).collect())
}

/// `textDocument/documentSymbol` result: a flat list of `DocumentSymbol`s.
#[must_use]
pub fn symbols_result(symbols: &[SymbolInfo], text: &str, encoding: PositionEncoding) -> Value {
    Value::Array(
        symbols
            .iter()
            .map(|s| symbol_json(s, text, encoding))
            .collect(),
    )
}

fn symbol_json(s: &SymbolInfo, text: &str, encoding: PositionEncoding) -> Value {
    let span = identifier_span(s, text, encoding);
    let kind = match s.kind {
        crate::analysis::SymbolKind::Function => SYMBOL_FUNCTION,
        crate::analysis::SymbolKind::Variable => SYMBOL_VARIABLE,
        crate::analysis::SymbolKind::Type => SYMBOL_CLASS,
    };
    json!({
        "name": s.name,
        "detail": s.ty,
        "kind": kind,
        "range": range_json(span),
        "selectionRange": range_json(span)
    })
}

/// The span of a symbol's NAME. The parser records a declaration's position at
/// its keyword (`fn`/`let`/`type`), so scan the declaration line for the first
/// whole-word occurrence of the name; fall back to the keyword column.
fn identifier_span(s: &SymbolInfo, text: &str, encoding: PositionEncoding) -> Span {
    let line = s.position.map_or(0, |p| p.line.saturating_sub(1));
    occurrences(text, &s.name, encoding)
        .into_iter()
        .find(|o| o.line == line)
        .map_or_else(
            || {
                let col = s.position.map_or(0, |p| p.column);
                (
                    line,
                    col,
                    line,
                    col.saturating_add(measure(&s.name, encoding)),
                )
            },
            |o| (o.line, o.start, o.line, o.end),
        )
}

/// `textDocument/signatureHelp` result, or JSON `null`.
#[must_use]
pub fn signature_result(info: Option<SignatureInfo>) -> Value {
    info.map_or(Value::Null, |s| {
        let params: Vec<Value> = s.parameters.iter().map(|p| json!({ "label": p })).collect();
        json!({
            "signatures": [{ "label": s.label, "parameters": params }],
            "activeSignature": 0,
            "activeParameter": s.active_parameter
        })
    })
}

/// `textDocument/completion` result: an array of completion items.
#[must_use]
pub fn completion_result(items: &[CompletionItem]) -> Value {
    Value::Array(items.iter().map(completion_json).collect())
}

fn completion_json(item: &CompletionItem) -> Value {
    let kind = match item.kind {
        CompletionKind::Keyword => COMPLETION_KEYWORD,
        CompletionKind::Function => COMPLETION_FUNCTION,
        CompletionKind::Variable => COMPLETION_VARIABLE,
        CompletionKind::Type => COMPLETION_CLASS,
    };
    let mut obj = json!({ "label": item.label, "kind": kind });
    insert_opt(&mut obj, "detail", item.detail.as_deref().map(Value::from));
    if let Some(text) = &item.insert_text {
        insert_opt(&mut obj, "insertText", Some(Value::from(text.clone())));
        insert_opt(
            &mut obj,
            "insertTextFormat",
            Some(Value::from(INSERT_SNIPPET)),
        );
    }
    obj
}

fn insert_opt(obj: &mut Value, key: &str, value: Option<Value>) {
    if let (Some(map), Some(value)) = (obj.as_object_mut(), value) {
        let _ = map.insert(key.to_owned(), value);
    }
}

/// `textDocument/publishDiagnostics` params for `uri`.
#[must_use]
pub fn publish_diagnostics(uri: &str, diagnostics: &[Diagnostic]) -> Value {
    json!({
        "uri": uri,
        "diagnostics": diagnostics.iter().map(diagnostic_json).collect::<Vec<_>>()
    })
}

fn diagnostic_json(d: &Diagnostic) -> Value {
    let severity = match d.severity {
        Severity::Warning => 2,
        Severity::Information => 3,
        Severity::Hint => 4,
        // Error and any future (`#[non_exhaustive]`) severity map to error.
        _ => 1,
    };
    let mut obj = json!({
        "range": range_json(d.range),
        "severity": severity,
        "message": d.message
    });
    insert_opt(&mut obj, "source", d.source.as_deref().map(Value::from));
    insert_opt(&mut obj, "code", d.code.as_deref().map(Value::from));
    obj
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn position_and_uri_parse_from_params() {
        let params = json!({
            "textDocument": { "uri": "file:///a.osp" },
            "position": { "line": 3, "character": 7 }
        });
        assert_eq!(doc_uri(&params).as_deref(), Some("file:///a.osp"));
        assert_eq!(position(&params), Some((3, 7)));
    }

    #[test]
    fn content_changes_split_incremental_and_full() {
        let params = json!({ "contentChanges": [
            { "range": { "start": { "line": 0, "character": 0 }, "end": { "line": 0, "character": 1 } }, "text": "x" },
            { "text": "whole file" }
        ] });
        let changes = content_changes(&params);
        assert!(matches!(changes.first(), Some(Ok(_))));
        assert!(matches!(changes.get(1), Some(Err(t)) if t == "whole file"));
    }

    #[test]
    fn hover_and_diagnostics_render_expected_shape() {
        assert_eq!(hover_result(None), Value::Null);
        let hov = hover_result(Some("**x**".to_owned()));
        assert_eq!(
            hov.pointer("/contents/kind"),
            Some(&Value::from("markdown"))
        );
        let diag = Diagnostic::new(Severity::Error, "boom", (1, 2, 1, 5)).with_source("osprey");
        let published = publish_diagnostics("file:///a.osp", &[diag]);
        assert_eq!(
            published.pointer("/diagnostics/0/severity"),
            Some(&Value::from(1))
        );
        assert_eq!(
            published.pointer("/diagnostics/0/source"),
            Some(&Value::from("osprey"))
        );
    }

    #[test]
    fn document_symbol_range_lands_on_the_identifier_not_the_keyword() {
        let src = "fn add(a: int) -> int = a\n";
        let parsed = osprey_syntax::parse_program(src);
        let syms = crate::analysis::collect_symbols(&parsed.program);
        let value = symbols_result(&syms, src, PositionEncoding::Utf16);
        // `add` is at column 3; the `fn` keyword is at column 0.
        assert_eq!(value.pointer("/0/name"), Some(&Value::from("add")));
        assert_eq!(
            value.pointer("/0/range/start/character"),
            Some(&Value::from(3))
        );
        assert_eq!(
            value.pointer("/0/range/end/character"),
            Some(&Value::from(6))
        );
    }
}
