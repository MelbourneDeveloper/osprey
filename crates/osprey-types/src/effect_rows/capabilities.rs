//! Capability names, resolved for the effect proof. A `perform`, `handler` or
//! `mask` that names a capability is, to the proof, the capability's effect at
//! an identity that includes the capability's site: the checker already
//! appended that identity to the published arguments, so renaming the target
//! to its effect lets ordinary discharge tell capabilities apart.
//! Implements [EFFECTS-INSTANCE].

use osprey_ast::mutate::{children_mut, statement_children_mut};
use osprey_ast::{Expr, Position, Program};
use std::collections::HashMap;

/// `program` with every capability target renamed to its effect, or `None`
/// when no site names a capability.
pub(super) fn resolved(
    program: &Program,
    capabilities: &HashMap<(u32, u32), String>,
) -> Option<Program> {
    if capabilities.is_empty() {
        return None;
    }
    let mut copy = program.clone();
    for statement in &mut copy.statements {
        statement_children_mut(statement, &mut |expression| {
            rename(expression, capabilities);
        });
    }
    Some(copy)
}

fn rename(expression: &mut Expr, capabilities: &HashMap<(u32, u32), String>) {
    if let Some((effect, position)) = target(expression) {
        if let Some(resolved) = position.and_then(|at| capabilities.get(&(at.line, at.column))) {
            effect.clone_from(resolved);
        }
    }
    children_mut(expression, &mut |child| rename(child, capabilities));
}

fn target(expression: &mut Expr) -> Option<(&mut String, Option<Position>)> {
    match expression {
        Expr::Perform {
            effect, position, ..
        }
        | Expr::Handler {
            effect, position, ..
        }
        | Expr::Mask {
            effect, position, ..
        } => Some((effect, *position)),
        _ => None,
    }
}
