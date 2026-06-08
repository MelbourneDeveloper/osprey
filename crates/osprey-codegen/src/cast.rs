//! Value coercion to a wanted LLVM type (numeric promotion at call/return/store
//! boundaries). String/handle targets take the operand as-is, only re-tagging
//! the LLVM type.

use crate::builder::Codegen;
use crate::conv::{as_double, as_i1, as_i64};
use crate::error::Result;
use crate::llty::{LType, Value};

/// Coerce a value to the wanted type, preserving its aggregate owner tag. A
/// `Result` arriving at a value boundary (argument, field, return scalar)
/// auto-unwraps to its success payload first, per the type-system spec.
pub(crate) fn coerce_to(cg: &mut Codegen, v: Value, want: LType) -> Result<Value> {
    let v = crate::result::unwrap(cg, v);
    if v.ty == want {
        return Ok(v);
    }
    let owner = v.osp_ty.clone();
    let out = match want {
        LType::Double => as_double(cg, v)?,
        // A string/handle reaching an `i64` boundary is an `any`/generic value
        // travelling in the uniform machine-word representation — `ptrtoint`-box
        // it (the inverse `inttoptr` is the `Str`/`Ptr` arm below). Genuine
        // type mismatches are already rejected by the checker.
        LType::I64 if matches!(v.ty, LType::Str | LType::Ptr) => crate::conv::box_to_i64(cg, v),
        LType::I64 => as_i64(cg, v)?,
        LType::I1 => as_i1(cg, v)?,
        // A pointer target. A boxed `i64` element (the uniform collection ABA)
        // must be `inttoptr`-converted back to a handle; an existing pointer
        // just retags (both are `i8*`).
        LType::Str | LType::Ptr | LType::I32 => {
            if v.ty == LType::I64 && matches!(want, LType::Str | LType::Ptr) {
                let reg = cg.fresh_reg();
                cg.emit(format!("{reg} = inttoptr i64 {} to i8*", v.operand));
                Value::new(reg, want)
            } else {
                Value::new(v.operand, want)
            }
        }
    };
    Ok(out.with_owner(owner))
}
