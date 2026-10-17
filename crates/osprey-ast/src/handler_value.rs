//! Callable handler values share the ordinary closure and lexical handler core.

use crate::{Expr, HandlerArm, Parameter, Position, Stage, Stmt};

// `$` cannot occur in a source identifier in either flavor.
const ACTION: &str = "$handler_action";

/// Build a reusable handler taking a zero-argument computation.
/// Implements [EFFECTS-HANDLER-VALUE] and [FLAVOR-HANDLER-VALUE].
#[must_use]
pub fn handler_value(
    stage: Stage,
    effect: String,
    arms: Vec<HandlerArm>,
    return_clause: Option<Box<Expr>>,
    position: Option<Position>,
) -> Expr {
    Expr::Lambda {
        parameters: vec![action_parameter()],
        return_type: None,
        body: Box::new(handled_action(stage, effect, arms, return_clause, position)),
        position,
    }
}

/// A named handler function over a zero-argument computation: the module
/// idiom `fn name(work) = { handle E { … }\n work() }`, which a module may
/// export where a module `let` must be a compile-time constant.
/// Implements [FLAVOR-PYTHON-HANDLER].
#[must_use]
pub fn handler_function(
    name: String,
    stage: Stage,
    effect: String,
    arms: Vec<HandlerArm>,
    position: Option<Position>,
) -> Stmt {
    Stmt::Function {
        name,
        type_params: Vec::new(),
        parameters: vec![action_parameter()],
        return_type: None,
        effects: Vec::new(),
        effect_tail: None,
        effect_row_present: false,
        body: handled_action(stage, effect, arms, None, position),
        doc: None,
        position,
    }
}

fn action_parameter() -> Parameter {
    Parameter {
        name: ACTION.into(),
        ty: None,
        inline_constraint: false,
    }
}

/// The handler installed over the invocation of its computation.
fn handled_action(
    stage: Stage,
    effect: String,
    arms: Vec<HandlerArm>,
    return_clause: Option<Box<Expr>>,
    position: Option<Position>,
) -> Expr {
    Expr::Handler {
        stage,
        effect,
        arms,
        return_clause,
        body: Box::new(Expr::Call {
            function: Box::new(Expr::Identifier(ACTION.into())),
            arguments: Vec::new(),
            named_arguments: Vec::new(),
        }),
        position,
    }
}
