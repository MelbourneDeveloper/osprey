//! Bridges Hindley-Milner inference ([`osprey_types`]) to the backend's LLVM
//! type lattice. The Go compiler interleaves type resolution with emission via
//! its `typeInferer`; here inference is a finished table ([`ProgramTypes`]) and
//! this module just maps an inferred [`Type`] to the [`LType`] the value travels
//! as. Unresolved/polymorphic variables degrade to `i64`, matching the C
//! runtime's uniform machine-word representation for generic values.

use crate::llty::LType;
use osprey_types::{names, Type};

/// Map an inferred type to the LLVM type a runtime value of it travels as.
pub fn ltype_of(ty: &Type) -> LType {
    match ty {
        Type::Con { name, args } => ltype_of_con(name, args),
        // A function reference is a code pointer; values never hold one directly
        // in the lowered programs (calls are direct), so treat as a handle.
        Type::Fun { .. } => LType::Ptr,
        // Records, unions and bare variables are runtime handles / machine words.
        Type::Record { .. } | Type::Union { .. } => LType::Ptr,
        Type::Var(_) => LType::I64,
    }
}

fn ltype_of_con(name: &str, args: &[Type]) -> LType {
    match name {
        names::INT => LType::I64,
        names::FLOAT => LType::Double,
        names::STRING => LType::Str,
        names::BOOL => LType::I1,
        names::UNIT => LType::I64,
        names::ANY => LType::I64,
        // Result<T, E> at a value site carries its unwrapped success value (the
        // auto-unwrap the type checker already applied), so it travels as T.
        names::RESULT => args.first().map(ltype_of).unwrap_or(LType::I64),
        // Collections, fibers, channels, pointers — opaque runtime handles.
        names::LIST | names::MAP | names::FIBER | names::CHANNEL | names::PTR => LType::Ptr,
        // A nullary user type name (nominal record/union referenced by name).
        _ => LType::Ptr,
    }
}

/// The Osprey owner type name to tag an aggregate value with, if `ty` is a
/// nominal record/union (so field access / match can recover its layout).
/// Scalars, collections and `Result` (auto-unwrapped at value sites) carry no
/// owner.
pub fn owner_name(ty: &Type) -> Option<String> {
    match ty {
        Type::Record { name, .. } | Type::Union { name, .. } => Some(name.clone()),
        Type::Con { name, .. } => match name.as_str() {
            names::INT
            | names::FLOAT
            | names::STRING
            | names::BOOL
            | names::UNIT
            | names::ANY
            | names::RESULT
            | names::LIST
            | names::MAP
            | names::FIBER
            | names::CHANNEL
            | names::PTR => None,
            other => Some(other.to_string()),
        },
        _ => None,
    }
}

/// When `ty` is `Result<T, E>`, the inner success type `T` as an [`LType`].
/// Used to carry the `{ T, i8 }*` Result block across call/return boundaries.
pub fn result_inner(ty: &Type) -> Option<LType> {
    match ty {
        Type::Con { name, args } if name == names::RESULT => args.first().map(ltype_of),
        _ => None,
    }
}

/// Map a field/parameter type *as written* (the string forms stored in
/// constructor layouts: `int`, `string`, `float`, `bool`, `Point`, …) to an
/// `LType`. Anything not a known scalar is a runtime handle.
pub fn ltype_of_name(written: &str) -> LType {
    // A generic application like `List<int>` keeps only its head for layout.
    let head = written.split(['<', '[']).next().unwrap_or(written).trim();
    match head {
        names::INT => LType::I64,
        names::FLOAT => LType::Double,
        names::STRING => LType::Str,
        names::BOOL => LType::I1,
        names::UNIT => LType::I64,
        names::ANY => LType::I64,
        _ => LType::Ptr,
    }
}
