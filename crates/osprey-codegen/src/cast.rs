//! Value coercion to a wanted LLVM type (numeric promotion at call/return/store
//! boundaries). String/handle targets take the operand as-is, only re-tagging
//! the LLVM type.

use crate::builder::Codegen;
use crate::conv::{as_double, as_i1, as_i64};
use crate::error::Result;
use crate::llty::{LType, Value};

/// Coerce a value to the wanted type, preserving its aggregate owner tag.
pub(crate) fn coerce_to(cg: &mut Codegen, v: Value, want: LType) -> Result<Value> {
    if v.ty == want {
        return Ok(v);
    }
    let owner = v.osp_ty.clone();
    let out = match want {
        LType::Double => as_double(cg, v)?,
        LType::I64 => as_i64(cg, v)?,
        LType::I1 => as_i1(cg, v)?,
        // String/handle targets keep the pointer operand; just retag.
        _ => Value::new(v.operand, want),
    };
    Ok(out.with_owner(owner))
}
