//! The Osprey type language — a Rust port of the `Type` interface hierarchy in
//! `compiler/internal/codegen/type_inference.go`.
//!
//! Go modelled types as an interface with six implementations (primitive,
//! concrete, generic, function, record, union, type-variable). The standard
//! Hindley-Milner representation collapses primitive/concrete/generic into one
//! *type-constructor application* (`Con`) — smaller, and exactly what
//! unification operates on. Names are kept identical to the Go `constants.go`
//! spellings so inferred types render the same for a differential test.

use std::collections::BTreeMap;
use std::fmt;

/// Identifier for an inference type variable (`TypeVar.id` in Go).
pub type VarId = u32;

/// Canonical type-constructor names. These MUST match `constants.go`.
pub mod names {
    pub const INT: &str = "int";
    pub const FLOAT: &str = "float";
    pub const STRING: &str = "string";
    pub const BOOL: &str = "bool";
    pub const ANY: &str = "any";
    pub const UNIT: &str = "Unit";
    pub const RESULT: &str = "Result";
    pub const MATH_ERROR: &str = "MathError";
    pub const LIST: &str = "List";
    pub const MAP: &str = "Map";
    pub const FIBER: &str = "Fiber";
    pub const CHANNEL: &str = "Channel";
    pub const PTR: &str = "Ptr";
}

/// A type in the Osprey type system.
#[derive(Debug, Clone, PartialEq)]
pub enum Type {
    /// An inference variable, rendered `t<id>`.
    Var(VarId),
    /// A named constructor applied to zero+ arguments. Zero args ⇒ a
    /// primitive/nullary type (`int`, `Unit`); arguments make it generic
    /// (`List<t>`, `Result<t, e>`).
    Con { name: String, args: Vec<Type> },
    /// A function `(p0, p1, ...) -> ret`.
    Fun { params: Vec<Type>, ret: Box<Type> },
    /// A structural record — equality is by field name+type, never field order
    /// (the HM-correctness fix the Go source calls out).
    Record {
        name: String,
        fields: BTreeMap<String, Type>,
    },
    /// A nominal sum type whose variants are nullary `Con`s or `Record`s.
    Union { name: String, variants: Vec<Type> },
}

impl Type {
    /// A constructor application, e.g. `Type::con("List", vec![Type::int()])`.
    pub fn con(name: impl Into<String>, args: Vec<Type>) -> Type {
        Type::Con {
            name: name.into(),
            args,
        }
    }
    /// A nullary named type (`int`, `Unit`, a bare user type).
    pub fn prim(name: impl Into<String>) -> Type {
        Type::con(name, Vec::new())
    }
    pub fn int() -> Type {
        Type::prim(names::INT)
    }
    pub fn float() -> Type {
        Type::prim(names::FLOAT)
    }
    pub fn string() -> Type {
        Type::prim(names::STRING)
    }
    pub fn bool() -> Type {
        Type::prim(names::BOOL)
    }
    pub fn unit() -> Type {
        Type::prim(names::UNIT)
    }
    pub fn any() -> Type {
        Type::prim(names::ANY)
    }
    pub fn ptr() -> Type {
        Type::prim(names::PTR)
    }
    pub fn fun(params: Vec<Type>, ret: Type) -> Type {
        Type::Fun {
            params,
            ret: Box::new(ret),
        }
    }
    /// `Result<ok, err>`.
    pub fn result(ok: Type, err: Type) -> Type {
        Type::con(names::RESULT, vec![ok, err])
    }
    pub fn list(elem: Type) -> Type {
        Type::con(names::LIST, vec![elem])
    }
    pub fn map(key: Type, value: Type) -> Type {
        Type::con(names::MAP, vec![key, value])
    }

    /// True if this is a nullary-or-applied constructor with the given name.
    pub fn is_named(&self, n: &str) -> bool {
        matches!(self, Type::Con { name, .. } if name == n)
    }
}

/// A polymorphic type scheme `forall vars. ty` (`TypeScheme` in Go), the engine
/// of let-polymorphism: generalize at bindings, instantiate at uses.
#[derive(Debug, Clone, PartialEq)]
pub struct Scheme {
    pub vars: Vec<VarId>,
    pub ty: Type,
}

impl Scheme {
    /// A monomorphic scheme — no quantified variables.
    pub fn mono(ty: Type) -> Scheme {
        Scheme {
            vars: Vec::new(),
            ty,
        }
    }
    /// A polymorphic scheme over the given variables.
    pub fn poly(vars: Vec<VarId>, ty: Type) -> Scheme {
        Scheme { vars, ty }
    }
}

impl fmt::Display for Type {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Type::Var(id) => write!(f, "t{id}"),
            Type::Con { name, args } if args.is_empty() => write!(f, "{name}"),
            Type::Con { name, args } => {
                write!(f, "{name}<")?;
                write_seq(f, args)?;
                write!(f, ">")
            }
            Type::Fun { params, ret } => {
                write!(f, "(")?;
                write_seq(f, params)?;
                write!(f, ") -> {ret}")
            }
            Type::Record { fields, .. } => {
                write!(f, "{{ ")?;
                for (i, (k, v)) in fields.iter().enumerate() {
                    if i > 0 {
                        write!(f, ", ")?;
                    }
                    write!(f, "{k}: {v}")?;
                }
                write!(f, " }}")
            }
            Type::Union { name, .. } => write!(f, "{name}"),
        }
    }
}

fn write_seq(f: &mut fmt::Formatter<'_>, items: &[Type]) -> fmt::Result {
    for (i, t) in items.iter().enumerate() {
        if i > 0 {
            write!(f, ", ")?;
        }
        write!(f, "{t}")?;
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn renders_primitives_and_generics() {
        assert_eq!(Type::int().to_string(), "int");
        assert_eq!(Type::list(Type::string()).to_string(), "List<string>");
        assert_eq!(
            Type::result(Type::int(), Type::prim("MathError")).to_string(),
            "Result<int, MathError>"
        );
        assert_eq!(
            Type::fun(vec![Type::int(), Type::int()], Type::bool()).to_string(),
            "(int, int) -> bool"
        );
        assert_eq!(Type::Var(3).to_string(), "t3");
    }
}
