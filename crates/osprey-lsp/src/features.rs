//! Feature computations over a document's source text.
//!
//! Each entry point parses with [`osprey_syntax`] and answers one editor
//! feature, returning the neutral [`crate::model`] types the server maps to the
//! wire protocol. Navigation is AST-driven (declarations carry real positions);
//! find-references falls back to whole-word scanning for occurrences.

use lspkit_vfs::PositionEncoding;

use crate::analysis::{builtin_hover, collect_symbols, SymbolInfo, SymbolKind};
use crate::model::{CompletionItem, CompletionKind, Location, SignatureInfo, Span};
use crate::text::{occurrences, prefix_to, word_at, Occurrence};

/// Hover markdown for the identifier at `(line, character)`.
#[must_use]
pub fn hover(text: &str, line: u32, character: u32, enc: PositionEncoding) -> Option<String> {
    let word = word_under(text, line, character, enc)?;
    let parsed = osprey_syntax::parse_program(text);
    collect_symbols(&parsed.program)
        .iter()
        .find(|s| s.name == word)
        .map(hover_markdown)
        .or_else(|| builtin_hover(&word))
}

fn hover_markdown(s: &SymbolInfo) -> String {
    match &s.signature {
        Some(sig) => format!("```osprey\n{sig}\n```"),
        None => format!("```osprey\n{}: {}\n```", s.name, s.ty),
    }
}

/// Definition location(s) for the identifier at `(line, character)`.
#[must_use]
pub fn definition(
    text: &str,
    uri: &str,
    line: u32,
    character: u32,
    enc: PositionEncoding,
) -> Vec<Location> {
    let Some(word) = word_under(text, line, character, enc) else {
        return Vec::new();
    };
    declarations(text, &word, enc)
        .into_iter()
        .map(|o| located(uri, (o.line, o.start, o.line, o.end)))
        .collect()
}

/// All references to the identifier at `(line, character)`.
#[must_use]
pub fn references(
    text: &str,
    uri: &str,
    line: u32,
    character: u32,
    enc: PositionEncoding,
    include_declaration: bool,
) -> Vec<Location> {
    let Some(word) = word_under(text, line, character, enc) else {
        return Vec::new();
    };
    let decls: Vec<(u32, u32)> = declarations(text, &word, enc)
        .iter()
        .map(|o| (o.line, o.start))
        .collect();
    occurrences(text, &word, enc)
        .into_iter()
        .filter(|o| include_declaration || !decls.contains(&(o.line, o.start)))
        .map(|o| located(uri, (o.line, o.start, o.line, o.end)))
        .collect()
}

/// Signature help for the call enclosing `(line, character)`.
#[must_use]
pub fn signature_help(
    text: &str,
    line: u32,
    character: u32,
    enc: PositionEncoding,
) -> Option<SignatureInfo> {
    let line_str = nth_line(text, line)?;
    let (name, active) = enclosing_call(prefix_to(line_str, character, enc))?;
    let parsed = osprey_syntax::parse_program(text);
    let sym = collect_symbols(&parsed.program)
        .into_iter()
        .find(|s| s.name == name && s.kind == SymbolKind::Function)?;
    let params: Vec<String> = sym.parameters.iter().map(param_label).collect();
    let last = u32::try_from(params.len().saturating_sub(1)).unwrap_or(0);
    Some(SignatureInfo {
        label: sym.signature.unwrap_or(sym.name),
        parameters: params,
        active_parameter: active.min(last),
    })
}

/// Completion items: keywords plus the document's own declarations.
#[must_use]
pub fn completion(text: &str) -> Vec<CompletionItem> {
    let parsed = osprey_syntax::parse_program(text);
    keyword_items()
        .into_iter()
        .chain(collect_symbols(&parsed.program).iter().map(symbol_item))
        .collect()
}

fn word_under(text: &str, line: u32, character: u32, enc: PositionEncoding) -> Option<String> {
    word_at(nth_line(text, line)?, character, enc).map(|w| w.word)
}

fn nth_line(text: &str, line: u32) -> Option<&str> {
    usize::try_from(line).ok().and_then(|i| text.lines().nth(i))
}

fn located(uri: &str, span: Span) -> Location {
    Location {
        uri: uri.to_owned(),
        span,
    }
}

/// The identifier occurrence of each declaration of `name`.
///
/// A declaration's recorded position points at its keyword (`fn`/`type`/`let`),
/// not the name, so this finds the first whole-word occurrence of `name` on each
/// declaration line — the location editors expect for go-to-definition.
fn declarations(text: &str, name: &str, enc: PositionEncoding) -> Vec<Occurrence> {
    let parsed = osprey_syntax::parse_program(text);
    let occs = occurrences(text, name, enc);
    collect_symbols(&parsed.program)
        .iter()
        .filter(|s| s.name == name)
        .filter_map(|s| s.position.map(|p| p.line.saturating_sub(1)))
        .filter_map(|line| occs.iter().find(|o| o.line == line).cloned())
        .collect()
}

fn param_label((name, ty): &(String, String)) -> String {
    if ty.is_empty() {
        name.clone()
    } else {
        format!("{name}: {ty}")
    }
}

fn symbol_item(s: &SymbolInfo) -> CompletionItem {
    let kind = match s.kind {
        SymbolKind::Function => CompletionKind::Function,
        SymbolKind::Variable => CompletionKind::Variable,
        SymbolKind::Type => CompletionKind::Type,
    };
    CompletionItem {
        label: s.name.clone(),
        kind,
        detail: Some(s.ty.clone()),
        insert_text: None,
    }
}

/// The fixed keyword/snippet completions (superset of the old TS server's six).
fn keyword_items() -> Vec<CompletionItem> {
    const KEYWORDS: [(&str, &str, &str); 6] = [
        (
            "fn",
            "Function declaration",
            "fn ${1:name}(${2:params}) = ${3:body}",
        ),
        ("let", "Variable declaration", "let ${1:name} = ${2:value}"),
        (
            "mut",
            "Mutable variable declaration",
            "mut ${1:name} = ${2:value}",
        ),
        (
            "match",
            "Pattern matching",
            "match ${1:expr} {\n\t${2:pattern} => ${3:result}\n}",
        ),
        (
            "type",
            "Type declaration",
            "type ${1:Name} = ${2:Variant} | ${3:Variant}",
        ),
        (
            "effect",
            "Effect declaration",
            "effect ${1:Name} {\n\t${2:op}: ${3:fn() -> Unit}\n}",
        ),
    ];
    KEYWORDS
        .iter()
        .map(|(label, detail, snippet)| CompletionItem {
            label: (*label).to_owned(),
            kind: CompletionKind::Keyword,
            detail: Some((*detail).to_owned()),
            insert_text: Some((*snippet).to_owned()),
        })
        .collect()
}

/// Parse `before` (the line text up to the cursor) and return the name of the
/// innermost still-open call and the active (comma-separated) argument index.
fn enclosing_call(before: &str) -> Option<(String, u32)> {
    let mut names: Vec<String> = Vec::new();
    let mut commas: Vec<u32> = Vec::new();
    let mut current = String::new();
    let mut last = String::new();
    for c in before.chars() {
        if c.is_alphanumeric() || c == '_' {
            current.push(c);
            continue;
        }
        if !current.is_empty() {
            last = std::mem::take(&mut current);
        }
        step_call(c, &mut names, &mut commas, &mut last);
    }
    let name = names.last().filter(|n| !n.is_empty())?;
    Some((name.clone(), commas.last().copied().unwrap_or(0)))
}

fn step_call(c: char, names: &mut Vec<String>, commas: &mut Vec<u32>, last: &mut String) {
    match c {
        '(' => {
            names.push(std::mem::take(last));
            commas.push(0);
        }
        ')' => {
            let _ = names.pop();
            let _ = commas.pop();
        }
        ',' => {
            if let Some(top) = commas.last_mut() {
                *top = top.saturating_add(1);
            }
        }
        _ => {}
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    const U16: PositionEncoding = PositionEncoding::Utf16;
    const SRC: &str = "fn add(a: int, b: int) -> int = a + b\nlet total = add(1, 2)\n";

    #[test]
    fn hover_uses_signature_for_functions_and_builtins() {
        assert!(hover(SRC, 1, 12, U16).is_some_and(|m| m.contains("fn add(a: int, b: int) -> int")));
        assert!(hover("fn main() = print(1)\n", 0, 13, U16).is_some_and(|m| m.contains("print")));
    }

    #[test]
    fn definition_points_at_the_declaration() {
        let defs = definition(SRC, "file:///a.osp", 1, 12, U16);
        let first = defs.first().expect("definition");
        assert_eq!(first.span.0, 0, "{defs:?}");
    }

    #[test]
    fn references_can_exclude_the_declaration() {
        let with = references(SRC, "file:///a.osp", 0, 3, U16, true);
        let without = references(SRC, "file:///a.osp", 0, 3, U16, false);
        assert_eq!(with.len(), 2);
        assert_eq!(without.len(), 1);
    }

    #[test]
    fn signature_help_tracks_the_active_parameter() {
        // Line 1 is `let total = add(1, 2)`; char 19 is over the second argument.
        let sig = signature_help(SRC, 1, 19, U16).expect("sig");
        assert_eq!(sig.active_parameter, 1, "{sig:?}");
        assert_eq!(sig.parameters.len(), 2);
    }

    #[test]
    fn completion_includes_keywords_and_declarations() {
        let items = completion(SRC);
        assert!(items
            .iter()
            .any(|i| i.label == "fn" && i.kind == CompletionKind::Keyword));
        assert!(items
            .iter()
            .any(|i| i.label == "add" && i.kind == CompletionKind::Function));
    }
}
