//! Bindings scoped to a body: a function's parameters, its local bindings and
//! the parameters of handler arms. They resolve before any top-level
//! declaration, so `work` in `onGpu(work)` and `work` in `onLists(work)` stay
//! two names. Implements [LSP-DEFINITION-LOCAL] and [LSP-REFERENCES-LOCAL].

use lspkit_vfs::PositionEncoding;
use osprey_ast::{walk_program, AstVisitor, Expr, Program};

use super::{declaration_occurrence, symbol_matches};
use crate::analysis::{collect_all_symbols, collect_symbols, SymbolInfo, SymbolKind};
use crate::text::{occurrences, Occurrence};

/// A scoped binding and the last line its scope can reach: the line before the
/// next top-level declaration.
pub(super) struct Binding {
    pub(super) at: Occurrence,
    last_line: u32,
}

/// The nearest binding of `word` written at or before `cursor` (0-based line
/// and column) inside the enclosing function, or `None` when the name is
/// top-level or unbound.
pub(super) fn binding(
    text: &str,
    program: &Program,
    word: &str,
    cursor: (u32, u32),
    enc: PositionEncoding,
) -> Option<Binding> {
    let all = collect_all_symbols(program);
    let top = collect_symbols(program);
    let owner = enclosing_function(&all, cursor.0);
    let first_line = owner
        .and_then(|function| function.position)
        .map_or(0, |position| position.line.saturating_sub(1));
    let lookup = Lookup { text, word, enc };
    let at = lookup
        .candidates(program, &all, &top, owner)
        .into_iter()
        .filter(|found| found.line >= first_line && (found.line, found.start) <= cursor)
        .max_by_key(|found| (found.line, found.start))?;
    let last_line = next_declaration_line(&top, at.line).unwrap_or(u32::MAX);
    Some(Binding { at, last_line })
}

/// Every occurrence of the binding's name inside its scope.
pub(super) fn uses(
    text: &str,
    word: &str,
    binding: &Binding,
    include_declaration: bool,
    enc: PositionEncoding,
) -> Vec<Occurrence> {
    occurrences(text, word, enc)
        .into_iter()
        .filter(|found| (binding.at.line..=binding.last_line).contains(&found.line))
        .filter(|found| include_declaration || *found != binding.at)
        .collect()
}

/// The declared function whose body contains `line` (0-based): the nearest
/// declaration at or above it. A parameter is only in scope inside its own
/// body, so resolving without this would let one function's `x` answer for
/// another's.
pub(crate) fn enclosing_function(symbols: &[SymbolInfo], line: u32) -> Option<&SymbolInfo> {
    let cursor = line.saturating_add(1); // AST positions are 1-based lines.
    symbols
        .iter()
        .filter(|s| s.kind == SymbolKind::Function)
        .filter(|s| s.position.is_some_and(|p| p.line <= cursor))
        .max_by_key(|s| s.position.map_or(0, |p| p.line))
}

fn next_declaration_line(top: &[SymbolInfo], line: u32) -> Option<u32> {
    top.iter()
        .filter_map(|symbol| symbol.position)
        .map(|position| position.line.saturating_sub(1))
        .filter(|declared| *declared > line)
        .min()
        .map(|declared| declared.saturating_sub(1))
}

struct Lookup<'a> {
    text: &'a str,
    word: &'a str,
    enc: PositionEncoding,
}

impl Lookup<'_> {
    fn candidates(
        &self,
        program: &Program,
        all: &[SymbolInfo],
        top: &[SymbolInfo],
        owner: Option<&SymbolInfo>,
    ) -> Vec<Occurrence> {
        let locals = all
            .iter()
            .filter(|symbol| symbol_matches(symbol, self.word))
            .filter(|symbol| {
                !top.iter()
                    .any(|t| t.name == symbol.name && t.position == symbol.position)
            })
            .filter_map(|symbol| declaration_occurrence(self.text, symbol, self.enc));
        owner
            .and_then(|function| self.parameter(function))
            .into_iter()
            .chain(locals)
            .chain(self.arm_parameters(program))
            .collect()
    }

    /// The parameter `word` of `function`, written after the function's name.
    fn parameter(&self, function: &SymbolInfo) -> Option<Occurrence> {
        function
            .parameters
            .iter()
            .any(|(name, _)| name == self.word)
            .then(|| declaration_occurrence(self.text, function, self.enc))
            .flatten()
            .and_then(|name| self.after(&name))
    }

    /// Each handler arm parameter named `word`, written after its operation.
    fn arm_parameters(&self, program: &Program) -> Vec<Occurrence> {
        let mut arms = Arms {
            word: self.word,
            found: Vec::new(),
        };
        walk_program(program, &mut arms);
        arms.found
            .iter()
            .filter_map(|(operation, line)| {
                let name = occurrences(self.text, operation, self.enc)
                    .into_iter()
                    .find(|found| found.line == *line)?;
                self.after(&name)
            })
            .collect()
    }

    /// The first occurrence of `word` after `anchor` ends.
    fn after(&self, anchor: &Occurrence) -> Option<Occurrence> {
        occurrences(self.text, self.word, self.enc)
            .into_iter()
            .find(|found| (found.line, found.start) >= (anchor.line, anchor.end))
    }
}

/// The operation name and 0-based line of every handler arm binding `word`.
struct Arms<'a> {
    word: &'a str,
    found: Vec<(String, u32)>,
}

impl AstVisitor for Arms<'_> {
    fn expression(&mut self, expression: &Expr) {
        if let Expr::Handler { arms, .. } = expression {
            self.found.extend(
                arms.iter()
                    .filter(|arm| arm.params.iter().any(|param| param == self.word))
                    .filter_map(|arm| {
                        Some((arm.operation.clone(), arm.position?.line.checked_sub(1)?))
                    }),
            );
        }
    }
}
