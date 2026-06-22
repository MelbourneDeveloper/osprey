//! CST -> AST lowering: an explicit recursive descent over tree-sitter named
//! nodes (no visitor plumbing, exhaustive matching).
//!
//! `parse_program` is the public entry: source text in, [`Program`] out (plus any
//! syntax errors discovered by tree-sitter). Errors are collected, never fatal:
//! the front-end never panics on bad input and always produces a best-effort AST.

use osprey_ast::{Position, Program};
use tree_sitter::{Node, Parser, Tree};

mod expr;
mod lower;

pub use lower::Lowerer;

/// A syntax error located in the source (an ERROR/MISSING node from tree-sitter).
#[derive(Debug, Clone, PartialEq)]
pub struct SyntaxError {
    /// Human-readable description of what went wrong at this location.
    pub message: String,
    /// Source location (line/column) where the error was detected.
    pub position: Position,
}

/// The result of lowering: the program plus any syntax errors. Errors being
/// non-empty does not prevent producing a best-effort tree.
#[derive(Debug, Clone, PartialEq)]
pub struct Parsed {
    /// The lowered program (best-effort even when errors are present).
    pub program: Program,
    /// Syntax errors discovered while parsing; empty on a clean parse.
    pub errors: Vec<SyntaxError>,
}

/// Parse Osprey source into a typed [`Program`].
#[must_use]
pub fn parse_program(source: &str) -> Parsed {
    let Some(tree) = parse_tree(source) else {
        return Parsed {
            program: Program {
                statements: Vec::new(),
            },
            errors: vec![SyntaxError {
                message: "failed to initialize Osprey grammar".to_owned(),
                position: Position { line: 1, column: 0 },
            }],
        };
    };
    let root = tree.root_node();
    let lowerer = Lowerer::new(source.as_bytes());
    let program = lowerer.lower_program(root);
    let mut errors = Vec::new();
    collect_errors(root, source.as_bytes(), &mut errors);
    Parsed { program, errors }
}

/// Run only the tree-sitter parse (used by tooling that wants the raw CST).
///
/// Returns [`None`] if the embedded Osprey grammar cannot be loaded or
/// tree-sitter declines to produce a tree (neither happens for a valid build).
#[must_use]
pub fn parse_tree(source: &str) -> Option<Tree> {
    let mut parser = Parser::new();
    parser
        .set_language(&tree_sitter_osprey::LANGUAGE.into())
        .ok()?;
    parser.parse(source, None)
}

fn collect_errors(node: Node<'_>, src: &[u8], out: &mut Vec<SyntaxError>) {
    if node.is_error() || node.is_missing() {
        let p = node.start_position();
        out.push(SyntaxError {
            message: if node.is_missing() {
                format!("missing {}", node.kind())
            } else {
                format!("syntax error near {:?}", node.utf8_text(src).unwrap_or(""))
            },
            position: Position {
                line: u32::try_from(p.row).unwrap_or(u32::MAX).saturating_add(1),
                column: u32::try_from(p.column).unwrap_or(u32::MAX),
            },
        });
    }
    let mut cursor = node.walk();
    for child in node.children(&mut cursor) {
        collect_errors(child, src, out);
    }
}

#[cfg(test)]
#[expect(
    clippy::indexing_slicing,
    reason = "test assertions: an out-of-bounds index is a test failure, not a production panic"
)]
mod tests {
    use super::*;
    use osprey_ast::{Expr, Pattern, Stmt};

    fn one(src: &str) -> Stmt {
        let parsed = parse_program(src);
        assert!(parsed.errors.is_empty(), "errors: {:?}", parsed.errors);
        assert_eq!(parsed.program.statements.len(), 1);
        parsed.program.statements.into_iter().next().unwrap()
    }

    #[test]
    fn lowers_let() {
        match one("let x = 42\n") {
            Stmt::Let {
                name,
                value,
                mutable,
                ..
            } => {
                assert_eq!(name, "x");
                assert!(!mutable);
                assert_eq!(value, Expr::Integer(42));
            }
            s => panic!("expected let, got {s:?}"),
        }
    }

    #[test]
    fn lowers_function_with_binary_body() {
        match one("fn add(a: int, b: int) -> int = a + b\n") {
            Stmt::Function {
                name,
                parameters,
                return_type,
                body,
                ..
            } => {
                assert_eq!(name, "add");
                assert_eq!(parameters.len(), 2);
                assert_eq!(parameters[0].name, "a");
                assert_eq!(return_type.unwrap().name, "int");
                match body {
                    Expr::Binary { op, .. } => assert_eq!(op, "+"),
                    b => panic!("expected binary, got {b:?}"),
                }
            }
            s => panic!("expected function, got {s:?}"),
        }
    }

    #[test]
    fn lowers_union_type() {
        match one("type Color = Red | Green | Blue\n") {
            Stmt::Type { name, variants, .. } => {
                assert_eq!(name, "Color");
                assert_eq!(variants.len(), 3);
                assert_eq!(variants[2].name, "Blue");
            }
            s => panic!("expected type, got {s:?}"),
        }
    }

    #[test]
    fn lowers_extern_with_ptr() {
        match one("extern fn sqlite3_open(filename: string, ppDb: Ptr) -> int\n") {
            Stmt::Extern {
                name,
                parameters,
                return_type,
                ..
            } => {
                assert_eq!(name, "sqlite3_open");
                assert_eq!(parameters.len(), 2);
                assert_eq!(parameters[1].ty.name, "Ptr");
                assert_eq!(return_type.unwrap().name, "int");
            }
            s => panic!("expected extern, got {s:?}"),
        }
    }

    #[test]
    fn lowers_match() {
        match one("let r = match x {\n  Ok { value } => value\n  _ => 0\n}\n") {
            Stmt::Let {
                value: Expr::Match { arms, .. },
                ..
            } => {
                assert_eq!(arms.len(), 2);
                assert!(matches!(arms[1].pattern, Pattern::Wildcard));
            }
            s => panic!("expected let-match, got {s:?}"),
        }
    }

    #[test]
    fn lowers_effect_and_perform() {
        let parsed = parse_program(
            "effect Log { info: fn(string) -> Unit }\nfn go() = perform Log.info(msg: \"hi\")\n",
        );
        assert!(parsed.errors.is_empty(), "{:?}", parsed.errors);
        assert!(matches!(parsed.program.statements[0], Stmt::Effect { .. }));
    }

    #[test]
    fn reports_syntax_error() {
        let parsed = parse_program("fn (= \n");
        assert!(!parsed.errors.is_empty());
    }

    #[test]
    fn reports_missing_node_error() {
        // `type T =` with no variant name forces tree-sitter to insert a MISSING
        // identifier; collect_errors reports it via the is_missing format branch.
        let parsed = parse_program("type T =\n");
        assert!(
            parsed
                .errors
                .iter()
                .any(|e| e.message.starts_with("missing")),
            "expected a missing-node error, got {:?}",
            parsed.errors
        );
        // The error carries a 1-based line.
        assert!(parsed.errors[0].position.line >= 1);
    }
}
