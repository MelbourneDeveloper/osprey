//! Editor-integration output: `--symbols` (document outline as JSON) and
//! `--hover` (built-in signature as markdown). The JSON shape is the contract
//! the VS Code extension's language server parses (`getSymbolsFromCompiler` in
//! `vscode-extension/server/src/server.ts`): an array of
//! `{name, kind, type, line, column, signature?, parameters?, returnType?}`
//! objects with 1-based `line` and `column`.

use osprey_ast::{ExternParameter, Parameter, Position, Program, Stmt, TypeExpr};
use std::fmt::Write as _;

/// One outline entry, rendered to a JSON object by [`sym_json`].
struct Sym {
    name: String,
    kind: &'static str,
    ty: String,
    position: Option<Position>,
    signature: Option<String>,
    parameters: Vec<(String, String)>,
    return_type: Option<String>,
}

/// The whole document outline as a JSON array (the `--symbols` payload).
pub(crate) fn symbols_json(program: &Program) -> String {
    let mut syms = Vec::new();
    collect(&program.statements, &mut syms);
    let rendered: Vec<String> = syms.iter().map(sym_json).collect();
    format!("[{}]", rendered.join(","))
}

/// Markdown hover text for a built-in function name (the `--hover` payload),
/// or `None` when the name is not a built-in.
pub(crate) fn builtin_hover(name: &str) -> Option<String> {
    osprey_types::builtin_signature(name).map(|sig| format!("```osprey\n{sig}\n```"))
}

/// Walk top-level statements (recursing into modules) into outline entries.
fn collect(stmts: &[Stmt], out: &mut Vec<Sym>) {
    for stmt in stmts {
        match stmt {
            Stmt::Module { body, .. } => collect(body, out),
            other => out.extend(sym_of(other)),
        }
    }
}

/// The outline entry for one declaration, or `None` for non-declarations.
fn sym_of(stmt: &Stmt) -> Option<Sym> {
    match stmt {
        Stmt::Function {
            name,
            parameters,
            return_type,
            position,
            ..
        } => Some(fn_sym(
            name,
            param_pairs(parameters),
            return_type.as_ref(),
            *position,
        )),
        Stmt::Extern {
            name,
            parameters,
            return_type,
            position,
        } => Some(fn_sym(
            name,
            extern_pairs(parameters),
            return_type.as_ref(),
            *position,
        )),
        Stmt::Let {
            name, ty, position, ..
        } => Some(let_sym(name, ty.as_ref(), *position)),
        Stmt::Type { name, position, .. } => Some(decl_sym(name, "type", *position)),
        Stmt::Effect { name, position, .. } => Some(decl_sym(name, "effect", *position)),
        _ => None,
    }
}

/// A function/extern entry: `type` and `signature` are the rendered signature.
fn fn_sym(
    name: &str,
    parameters: Vec<(String, String)>,
    return_type: Option<&TypeExpr>,
    position: Option<Position>,
) -> Sym {
    let ret = return_type.map_or_else(|| String::from("Unit"), render_type);
    let shown: Vec<String> = parameters
        .iter()
        .map(|(n, t)| {
            if t.is_empty() {
                n.clone()
            } else {
                format!("{n}: {t}")
            }
        })
        .collect();
    let signature = format!("fn {name}({}) -> {ret}", shown.join(", "));
    Sym {
        name: name.into(),
        kind: "function",
        ty: signature.clone(),
        position,
        signature: Some(signature),
        parameters,
        return_type: Some(ret),
    }
}

/// A `let` binding entry; `type` is the written annotation when present.
fn let_sym(name: &str, ty: Option<&TypeExpr>, position: Option<Position>) -> Sym {
    Sym {
        name: name.into(),
        kind: "variable",
        ty: ty.map(render_type).unwrap_or_default(),
        position,
        signature: None,
        parameters: Vec::new(),
        return_type: None,
    }
}

/// A `type`/`effect` declaration entry (both surface as outline kind `type`).
fn decl_sym(name: &str, ty: &str, position: Option<Position>) -> Sym {
    Sym {
        name: name.into(),
        kind: "type",
        ty: ty.into(),
        position,
        signature: None,
        parameters: Vec::new(),
        return_type: None,
    }
}

/// `(name, rendered type or "")` pairs for ordinary parameters.
fn param_pairs(params: &[Parameter]) -> Vec<(String, String)> {
    params
        .iter()
        .map(|p| {
            (
                p.name.clone(),
                p.ty.as_ref().map(render_type).unwrap_or_default(),
            )
        })
        .collect()
}

/// `(name, rendered type)` pairs for extern parameters (always typed).
fn extern_pairs(params: &[ExternParameter]) -> Vec<(String, String)> {
    params
        .iter()
        .map(|p| (p.name.clone(), render_type(&p.ty)))
        .collect()
}

/// Render a written type expression back to source-ish text.
fn render_type(t: &TypeExpr) -> String {
    if t.is_function {
        let ps: Vec<String> = t.parameter_types.iter().map(render_type).collect();
        let ret = t
            .return_type
            .as_deref()
            .map_or_else(|| String::from("Unit"), render_type);
        return format!("fn({}) -> {ret}", ps.join(", "));
    }
    if t.is_array {
        return t
            .array_element
            .as_deref()
            .map_or_else(|| String::from("[]"), |e| format!("[{}]", render_type(e)));
    }
    if t.generic_params.is_empty() {
        return t.name.clone();
    }
    let gs: Vec<String> = t.generic_params.iter().map(render_type).collect();
    format!("{}<{}>", t.name, gs.join(", "))
}

/// Render one entry as a JSON object. The AST column is 0-based; the wire
/// format is 1-based, so it is shifted here.
fn sym_json(s: &Sym) -> String {
    let (line, column) = s
        .position
        .map_or((1, 1), |p| (p.line, p.column.saturating_add(1)));
    let mut o = format!(
        "{{\"name\":{},\"kind\":{},\"type\":{},\"line\":{line},\"column\":{column}",
        json_str(&s.name),
        json_str(s.kind),
        json_str(&s.ty)
    );
    if let Some(sig) = &s.signature {
        let _ = write!(o, ",\"signature\":{}", json_str(sig));
    }
    if !s.parameters.is_empty() {
        let _ = write!(o, ",\"parameters\":{}", params_json(&s.parameters));
    }
    if let Some(ret) = &s.return_type {
        let _ = write!(o, ",\"returnType\":{}", json_str(ret));
    }
    o.push('}');
    o
}

/// Render `(name, type)` pairs as a JSON array of `{name, type}` objects.
fn params_json(params: &[(String, String)]) -> String {
    let items: Vec<String> = params
        .iter()
        .map(|(n, t)| format!("{{\"name\":{},\"type\":{}}}", json_str(n), json_str(t)))
        .collect();
    format!("[{}]", items.join(","))
}

/// A JSON string literal: quoted, with the mandatory escapes.
fn json_str(s: &str) -> String {
    let mut out = String::with_capacity(s.len().saturating_add(2));
    out.push('"');
    for c in s.chars() {
        match c {
            '"' => out.push_str("\\\""),
            '\\' => out.push_str("\\\\"),
            '\n' => out.push_str("\\n"),
            '\r' => out.push_str("\\r"),
            '\t' => out.push_str("\\t"),
            c if u32::from(c) < 0x20 => {
                let _ = write!(out, "\\u{:04x}", u32::from(c));
            }
            c => out.push(c),
        }
    }
    out.push('"');
    out
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn outline_covers_every_declaration_form() {
        let parsed = osprey_syntax::parse_program(
            "type Shade = Light | Dark\n\
             effect Log { info: fn(string) -> Unit }\n\
             extern fn puts(s: string) -> int\n\
             let limit: int = 10\n\
             fn multiply(a: int, b: int) -> int = a * b\n\
             fn main() -> Unit = print(multiply(a: limit, b: 2))\n",
        );
        assert!(parsed.errors.is_empty(), "{:?}", parsed.errors);
        let json = symbols_json(&parsed.program);
        for frag in [
            "\"name\":\"Shade\",\"kind\":\"type\",\"type\":\"type\",\"line\":1,\"column\":1",
            "\"name\":\"Log\",\"kind\":\"type\",\"type\":\"effect\",\"line\":2",
            "\"name\":\"puts\",\"kind\":\"function\"",
            "\"signature\":\"fn puts(s: string) -> int\"",
            "\"name\":\"limit\",\"kind\":\"variable\",\"type\":\"int\",\"line\":4",
            "\"name\":\"multiply\",\"kind\":\"function\"",
            "\"signature\":\"fn multiply(a: int, b: int) -> int\"",
            "\"parameters\":[{\"name\":\"a\",\"type\":\"int\"},{\"name\":\"b\",\"type\":\"int\"}]",
            "\"returnType\":\"int\"",
            "\"name\":\"main\",\"kind\":\"function\",\"type\":\"fn main() -> Unit\",\"line\":6",
        ] {
            assert!(json.contains(frag), "missing {frag} in {json}");
        }
    }

    #[test]
    fn hover_renders_builtin_signature_and_rejects_unknowns() {
        let md = builtin_hover("print");
        assert!(
            md.as_deref().is_some_and(|m| m.contains("print : ")),
            "{md:?}"
        );
        assert!(builtin_hover("notARealBuiltin").is_none());
    }

    #[test]
    fn json_strings_escape_quotes_and_control_chars() {
        assert_eq!(json_str("a\"b\\c\nd"), "\"a\\\"b\\\\c\\nd\"");
        assert_eq!(json_str("\u{1}"), "\"\\u0001\"");
    }
}
