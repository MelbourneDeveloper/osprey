//! Type errors as plain located messages — the CLI prints `file:line:col:
//! message`, identical to the syntax-error path.

use crate::ty::Type;
use osprey_ast::Position;

/// A single type error with an optional source position.
#[derive(Debug, Clone, PartialEq)]
pub struct TypeError {
    /// Human-readable description of what went wrong.
    pub message: String,
    /// Source location the error refers to, when known.
    pub position: Option<Position>,
}

impl TypeError {
    /// Create an error with a message but no associated position.
    pub fn new(message: impl Into<String>) -> TypeError {
        TypeError {
            message: message.into(),
            position: None,
        }
    }

    /// Create an error with a message anchored to a source position.
    pub fn at(message: impl Into<String>, position: Position) -> TypeError {
        TypeError {
            message: message.into(),
            position: Some(position),
        }
    }

    /// Attach a position if one is known and none is set yet.
    #[must_use]
    pub fn with_pos(mut self, position: Option<Position>) -> TypeError {
        if self.position.is_none() {
            self.position = position;
        }
        self
    }

    /// Build an error for two types that fail to unify.
    #[must_use]
    pub fn mismatch(a: &Type, b: &Type) -> TypeError {
        TypeError::new(format!("type mismatch: cannot unify {a} with {b}"))
    }

    /// Build an error for a type that recursively contains itself (occurs check).
    #[must_use]
    pub fn recursive(a: &Type, b: &Type) -> TypeError {
        TypeError::new(format!("recursive type: {a} occurs in {b}"))
    }
}
