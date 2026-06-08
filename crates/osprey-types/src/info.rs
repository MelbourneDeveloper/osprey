//! Resolved type information published for the code generator. Type inference
//! runs to completion (`type_inference.go`'s `ResolveAllEnvironmentTypes`
//! equivalent); the resulting concrete signatures, record layouts and union
//! tags are then frozen into the plain, substitution-free tables below so the
//! backend can drive codegen off real types instead of guessing `i64`.

use crate::ty::Type;
use std::collections::HashMap;

/// The declared shape of a record/variant constructor: ordered `(field, type)`
/// pairs written as type names (`int`, `string`, `Point`, …) plus its owning
/// type. Field type strings are kept verbatim so the backend maps them to its
/// own LLVM type lattice.
#[derive(Debug, Clone)]
pub struct CtorLayout {
    /// The type this constructor builds (`Point` for a record, `Shape` for a
    /// union variant).
    pub owner: String,
    /// Whether the owner is a single-variant record (vs. a union variant).
    pub owner_is_record: bool,
    /// Ordered `(field name, field type as written)`.
    pub fields: Vec<(String, String)>,
}

/// Everything the code generator needs from inference: per-function signatures
/// (fully resolved against the final substitution), constructor layouts and
/// union memberships.
#[derive(Debug, Clone, Default)]
pub struct ProgramTypes {
    /// Function/extern name → (resolved parameter types, resolved return type).
    pub functions: HashMap<String, (Vec<Type>, Type)>,
    /// Constructor name → its record/variant layout.
    pub ctors: HashMap<String, CtorLayout>,
    /// Union type name → ordered variant constructor names (tag order).
    pub unions: HashMap<String, Vec<String>>,
}

impl ProgramTypes {
    /// The resolved return type of a named function, if known.
    pub fn return_type(&self, name: &str) -> Option<&Type> {
        self.functions.get(name).map(|(_, ret)| ret)
    }

    /// The resolved parameter types of a named function, if known.
    pub fn param_types(&self, name: &str) -> Option<&[Type]> {
        self.functions.get(name).map(|(p, _)| p.as_slice())
    }
}
