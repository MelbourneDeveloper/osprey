//! Emission of the C-runtime / libc calls that back Osprey's built-ins:
//! `toString` per type, `print`, and the numeric→string conversions. Float
//! formatting is delegated to `osp_float_to_string` (linked from
//! `libfiber_runtime`) so whole-valued floats keep their visible `.0`, exactly
//! as `llvm.go`'s `generateFloatToString` does.

use crate::builder::Codegen;
use crate::conv::as_i64;
use crate::error::Result;
use crate::llty::{LType, Value};

/// Convert any value to its `i8*` string form (`toString` / interpolation /
/// `print`). Strings pass through; the rest go through libc `sprintf` or the
/// float runtime.
pub(crate) fn to_string_value(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::Str | LType::Ptr => Ok(Value::new(v.operand, LType::Str)),
        LType::I1 => Ok(bool_to_string(cg, v)),
        LType::Double => Ok(float_to_string(cg, v)),
        LType::I64 | LType::I32 => int_to_string(cg, v),
    }
}

/// `print(x)` → `puts(toString(x))`; yields Unit.
pub(crate) fn gen_print(cg: &mut Codegen, v: Value) -> Result<Value> {
    let s = to_string_value(cg, v)?;
    cg.add_extern("declare i32 @puts(i8*)");
    let reg = cg.fresh_reg();
    cg.emit(format!("{reg} = call i32 @puts(i8* {})", s.operand));
    Ok(Value::unit())
}

pub(crate) fn int_to_string(cg: &mut Codegen, v: Value) -> Result<Value> {
    cg.add_extern("declare i8* @malloc(i64)");
    cg.add_extern("declare i32 @sprintf(i8*, i8*, ...)");
    let i = as_i64(cg, v)?;
    let fmt = cg.string_constant("%ld");
    let buf = cg.fresh_reg();
    cg.emit(format!("{buf} = call i8* @malloc(i64 32)"));
    let tmp = cg.fresh_reg();
    cg.emit(format!(
        "{tmp} = call i32 (i8*, i8*, ...) @sprintf(i8* {buf}, i8* {}, i64 {})",
        fmt.operand, i.operand
    ));
    Ok(Value::new(buf, LType::Str))
}

/// Whole-valued floats must print with a trailing `.0`; the runtime handles
/// that (and NaN/inf) — see `runtime/string_runtime.c`.
pub(crate) fn float_to_string(cg: &mut Codegen, v: Value) -> Value {
    cg.add_extern("declare i8* @osp_float_to_string(double)");
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = call i8* @osp_float_to_string(double {})",
        v.operand
    ));
    Value::new(reg, LType::Str)
}

pub(crate) fn bool_to_string(cg: &mut Codegen, v: Value) -> Value {
    let t = cg.string_constant("true");
    let f = cg.string_constant("false");
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = select i1 {}, i8* {}, i8* {}",
        v.operand, t.operand, f.operand
    ));
    Value::new(reg, LType::Str)
}
