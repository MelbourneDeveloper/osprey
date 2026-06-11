//! Numeric/boolean coercions between LLVM machine types. These bridge the few
//! places where an operand arrives wider/narrower than the instruction wants
//! (a bool used as an int, an int used as a condition, an int promoted to a
//! double for mixed arithmetic).

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::llty::{LType, Value};

/// Coerce to `i64`.
pub(crate) fn as_i64(cg: &mut Codegen, v: Value) -> Result<Value> {
    let reg = match v.ty {
        LType::I64 => return Ok(v),
        LType::I1 => cg.emit_reg(format!("zext i1 {} to i64", v.operand)),
        LType::I32 => cg.emit_reg(format!("sext i32 {} to i64", v.operand)),
        LType::Double => cg.emit_reg(format!("fptosi double {} to i64", v.operand)),
        LType::Str | LType::Ptr => {
            return Err(CodegenError::invalid(
                "expected an integer, found a string/handle",
            ))
        }
    };
    Ok(Value::new(reg, LType::I64))
}

/// Coerce to `i1` (truthiness: non-zero).
pub(crate) fn as_i1(cg: &mut Codegen, v: Value) -> Result<Value> {
    let reg = match v.ty {
        LType::I1 => return Ok(v),
        LType::I64 | LType::I32 => cg.emit_reg(format!("icmp ne {} {}, 0", v.ty, v.operand)),
        LType::Double | LType::Str | LType::Ptr => {
            return Err(CodegenError::invalid("expected a bool"))
        }
    };
    Ok(Value::new(reg, LType::I1))
}

/// Widen any value to the uniform `i64` collection-element ABI: pointers
/// `ptrtoint`, narrow ints `zext`, `double` `bitcast`.
pub(crate) fn box_to_i64(cg: &mut Codegen, v: Value) -> Value {
    let reg = match v.ty {
        LType::I64 => return v,
        LType::Str | LType::Ptr => cg.emit_reg(format!("ptrtoint {} {} to i64", v.ty, v.operand)),
        LType::I1 | LType::I32 => cg.emit_reg(format!("zext {} {} to i64", v.ty, v.operand)),
        LType::Double => cg.emit_reg(format!("bitcast double {} to i64", v.operand)),
    };
    Value::new(reg, LType::I64)
}

/// Coerce to `double` (promoting an integer operand for mixed arithmetic).
pub(crate) fn as_double(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::Double => Ok(v),
        LType::I64 => Ok(Value::new(
            cg.emit_reg(format!("sitofp i64 {} to double", v.operand)),
            LType::Double,
        )),
        LType::I1 => {
            let i = as_i64(cg, v)?;
            as_double(cg, i)
        }
        LType::I32 | LType::Str | LType::Ptr => Err(CodegenError::invalid("expected a number")),
    }
}
