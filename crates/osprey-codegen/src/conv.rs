//! Numeric/boolean coercions between LLVM machine types. These bridge the few
//! places where an operand arrives wider/narrower than the instruction wants
//! (a bool used as an int, an int used as a condition, an int promoted to a
//! double for mixed arithmetic).

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::llty::{LType, Value};

/// Coerce to `i64`.
pub(crate) fn as_i64(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::I64 => Ok(v),
        LType::I1 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = zext i1 {} to i64", v.operand));
            Ok(Value::new(reg, LType::I64))
        }
        LType::I32 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = sext i32 {} to i64", v.operand));
            Ok(Value::new(reg, LType::I64))
        }
        LType::Double => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = fptosi double {} to i64", v.operand));
            Ok(Value::new(reg, LType::I64))
        }
        LType::Str | LType::Ptr => Err(CodegenError::invalid(
            "expected an integer, found a string/handle",
        )),
    }
}

/// Coerce to `i1` (truthiness: non-zero).
pub(crate) fn as_i1(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::I1 => Ok(v),
        LType::I64 | LType::I32 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = icmp ne {} {}, 0", v.ty, v.operand));
            Ok(Value::new(reg, LType::I1))
        }
        LType::Double | LType::Str | LType::Ptr => Err(CodegenError::invalid("expected a bool")),
    }
}

/// Widen any value to the uniform `i64` collection-element ABI: pointers
/// `ptrtoint`, narrow ints `zext`, `double` `bitcast`. Ports `boxToI64`.
pub(crate) fn box_to_i64(cg: &mut Codegen, v: Value) -> Value {
    match v.ty {
        LType::I64 => v,
        LType::Str | LType::Ptr => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = ptrtoint {} {} to i64", v.ty, v.operand));
            Value::new(reg, LType::I64)
        }
        LType::I1 | LType::I32 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = zext {} {} to i64", v.ty, v.operand));
            Value::new(reg, LType::I64)
        }
        LType::Double => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = bitcast double {} to i64", v.operand));
            Value::new(reg, LType::I64)
        }
    }
}

/// Coerce to `double` (promoting an integer operand for mixed arithmetic).
pub(crate) fn as_double(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::Double => Ok(v),
        LType::I64 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = sitofp i64 {} to double", v.operand));
            Ok(Value::new(reg, LType::Double))
        }
        LType::I1 => {
            let i = as_i64(cg, v)?;
            as_double(cg, i)
        }
        LType::I32 | LType::Str | LType::Ptr => Err(CodegenError::invalid("expected a number")),
    }
}
