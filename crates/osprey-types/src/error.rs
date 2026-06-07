//! Type errors. Mirrors the wrapped-error style of `compiler/internal/codegen/
//! errors.go` but as plain located messages — the CLI prints `file:line:col:
//! message`, identical to the syntax-error path.

use crate::ty::Type;
use osprey_ast::Position;

/// A single type error with an optional source position.
#[derive(Debug, Clone, PartialEq)]
pub struct TypeError {
    pub message: String,
    pub position: Option<Position>,
}

impl TypeError {
    pub fn new(message: impl Into<String>) -> TypeError {
        TypeError {
            message: message.into(),
            position: None,
        }
    }

    pub fn at(message: impl Into<String>, position: Position) -> TypeError {
        TypeError {
            message: message.into(),
            position: Some(position),
        }
    }

    /// Attach a position if one is known and none is set yet.
    pub fn with_pos(mut self, position: Option<Position>) -> TypeError {
        if self.position.is_none() {
            self.position = position;
        }
        self
    }

    pub fn mismatch(a: &Type, b: &Type) -> TypeError {
        TypeError::new(format!("type mismatch: cannot unify {a} with {b}"))
    }

    pub fn recursive(a: &Type, b: &Type) -> TypeError {
        TypeError::new(format!("recursive type: {a} occurs in {b}"))
    }
}
