//! The small slice of the LLVM type system the code generator emits. Osprey
//! primitives map to `i64` (int), `i1` (bool) and `i8*` (string); `i32` is the
//! C `main` return and `puts`/`sprintf` result. Records, unions and effect
//! closures are not yet lowered — those callers fail loudly rather than guess.

use std::fmt;

/// An LLVM first-class type the emitter knows how to name and move around.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum LType {
    /// 64-bit integer — Osprey `int`.
    I64,
    /// 1-bit integer — Osprey `bool`.
    I1,
    /// `i8*` — an Osprey `string` (NUL-terminated C string).
    Str,
    /// 32-bit integer — `main` return / libc call results.
    I32,
}

impl LType {
    /// The textual LLVM spelling.
    pub fn as_str(self) -> &'static str {
        match self {
            LType::I64 => "i64",
            LType::I1 => "i1",
            LType::Str => "i8*",
            LType::I32 => "i32",
        }
    }
}

impl fmt::Display for LType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.as_str())
    }
}

/// An SSA value: a textual operand (`%3`, a literal like `42`, or a `getelementptr`
/// result) paired with its LLVM type.
#[derive(Debug, Clone)]
pub struct Value {
    pub operand: String,
    pub ty: LType,
}

impl Value {
    pub fn new(operand: impl Into<String>, ty: LType) -> Value {
        Value {
            operand: operand.into(),
            ty,
        }
    }

    /// The canonical Unit value — Osprey `Unit` carries no data, so it is the
    /// `i64 0` placeholder a side-effecting expression yields.
    pub fn unit() -> Value {
        Value::new("0", LType::I64)
    }

    /// Render as a typed operand, e.g. `i64 %3` — the form arguments and `ret`
    /// take.
    pub fn typed(&self) -> String {
        format!("{} {}", self.ty, self.operand)
    }
}
