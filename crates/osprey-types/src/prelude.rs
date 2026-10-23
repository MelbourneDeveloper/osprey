//! The built-in source prelude: Osprey declarations every program may use
//! without declaring them. The prelude is attached on demand — only to a
//! program that names one of its declarations — so a program that does not
//! is checked and compiled exactly as written. Implements [BUILTIN-JSON-TYPED].

use osprey_ast::{walk_program, AstVisitor, Expr, Position, Program, Stmt};
use std::borrow::Cow;
use std::collections::BTreeSet;
use std::sync::OnceLock;

const JSON: &str = include_str!("prelude/json.osp");

/// Prelude positions start here, past any assembled project, so no prelude
/// inference site shares a key with a source one and coverage can leave the
/// prelude's lines out of a program's report.
pub const PRELUDE_LINE_BASE: u32 = 1 << 24;

/// Whether a source line belongs to the prelude rather than the program.
#[must_use]
pub fn is_prelude_line(line: u32) -> bool {
    line > PRELUDE_LINE_BASE
}

/// The parsed, rebased prelude. Its source is a constant of this crate and a
/// unit test proves it parses clean, so the parse never reports errors here.
fn prelude() -> &'static Program {
    static PRELUDE: OnceLock<Program> = OnceLock::new();
    PRELUDE.get_or_init(|| {
        let mut program = osprey_syntax::parse_program(JSON).program;
        osprey_ast::span::offset_program(&mut program, PRELUDE_LINE_BASE);
        program
    })
}

/// The program with the prelude in front of it when the program needs it.
/// A program that declares a type or constructor the prelude declares keeps
/// its own and gets no prelude; one that declares a prelude function name
/// keeps its own function in place of that one. Attaching is idempotent.
pub(crate) fn with_prelude(program: &Program) -> Cow<'_, Program> {
    let prelude = prelude();
    let own = Names::declared(program);
    let offered = Names::declared(prelude);
    let needed =
        own.data.is_disjoint(&offered.data) && !attached(program) && mentions(program, &offered);
    if !needed {
        return Cow::Borrowed(program);
    }
    let statements = prelude
        .statements
        .iter()
        .filter(|statement| {
            function_name(statement).is_none_or(|name| !own.functions.contains(name))
        })
        .chain(&program.statements)
        .cloned()
        .collect();
    let mut attached = Program {
        statements,
        doc: program.doc.clone(),
    };
    // The program was parsed without the prelude in view, so its parser left
    // every `JsonInt(1)` a call; the declarations are in view now.
    let arities = osprey_ast::positional::declared_arities(&prelude.statements);
    osprey_ast::positional::fold_constructions(&mut attached, &arities);
    Cow::Owned(attached)
}

fn attached(program: &Program) -> bool {
    program
        .statements
        .first()
        .and_then(type_position)
        .is_some_and(|position| is_prelude_line(position.line))
}

fn type_position(statement: &Stmt) -> Option<Position> {
    match statement {
        Stmt::Type { position, .. } => *position,
        _ => None,
    }
}

fn function_name(statement: &Stmt) -> Option<&str> {
    match statement {
        Stmt::Function { name, .. } => Some(name),
        _ => None,
    }
}

/// Whether any expression in `program` names a type, constructor or public
/// function the prelude offers.
fn mentions(program: &Program, offered: &Names) -> bool {
    let mut uses = Uses::default();
    walk_program(program, &mut uses);
    uses.names.iter().any(|name| {
        offered.data.contains(name) || (offered.functions.contains(name) && !name.starts_with("__"))
    })
}

#[derive(Default)]
struct Names {
    /// Declared type and constructor names.
    data: BTreeSet<String>,
    /// Declared function names.
    functions: BTreeSet<String>,
}

impl Names {
    fn declared(program: &Program) -> Self {
        let mut names = Self::default();
        walk_program(program, &mut names);
        names
    }
}

impl AstVisitor for Names {
    fn statement(&mut self, statement: &Stmt) {
        match statement {
            Stmt::Type { name, variants, .. } => {
                let _ = self.data.insert(name.clone());
                self.data
                    .extend(variants.iter().map(|variant| variant.name.clone()));
            }
            Stmt::Function { name, .. } => {
                let _ = self.functions.insert(name.clone());
            }
            _ => {}
        }
    }
}

#[derive(Default)]
struct Uses {
    names: BTreeSet<String>,
}

impl AstVisitor for Uses {
    fn expression(&mut self, expression: &Expr) {
        match expression {
            Expr::Identifier(name) | Expr::TypeConstructor { name, .. } => {
                let _ = self.names.insert(name.clone());
            }
            _ => {}
        }
    }
}

#[cfg(test)]
mod tests;
