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
/// float runtime. A `Result` formats as `Success(value)` / `Error(message)`.
pub(crate) fn to_string_value(cg: &mut Codegen, v: Value) -> Result<Value> {
    if v.result_inner.is_some() {
        return result_to_string(cg, &v);
    }
    match v.ty {
        LType::Str | LType::Ptr => Ok(Value::new(v.operand, LType::Str)),
        LType::I1 => Ok(bool_to_string(cg, &v)),
        LType::Double => Ok(float_to_string(cg, &v)),
        LType::I64 | LType::I32 => int_to_string(cg, v),
    }
}

/// Format a `Result` block as `Success(<value>)` or `Error(<message>)`, branching
/// on its discriminant — ports `convertResultToString`. A NULL string error
/// payload prints the bare `Error`.
fn result_to_string(cg: &mut Codegen, v: &Value) -> Result<Value> {
    let inner = v
        .result_inner
        .ok_or_else(|| crate::error::CodegenError::invalid("result_to_string on non-Result"))?;
    let (_sl, el, end) = crate::result::open_result_branch(cg, v);
    let val = crate::result::load_value(cg, v);
    let vs = to_string_value(cg, val)?;
    let succ = sprintf_wrap(cg, "Success(%s)", &vs.operand);
    let sb = cg.cur_block().to_string();
    cg.emit(format!("br label %{end}"));

    cg.start_block(&el);
    let err = if inner == LType::Str {
        // A runtime-stored NULL payload prints a bare "Error".
        let payload = crate::result::load_value(cg, v);
        let isnull = cg.fresh_reg();
        cg.emit(format!("{isnull} = icmp eq i8* {}, null", payload.operand));
        let fl = cg.fresh_label();
        let nl = cg.fresh_label();
        let jl = cg.fresh_label();
        cg.emit(format!("br i1 {isnull}, label %{nl}, label %{fl}"));
        cg.start_block(&fl);
        let with = sprintf_wrap(cg, "Error(%s)", &payload.operand);
        let fb = cg.cur_block().to_string();
        cg.emit(format!("br label %{jl}"));
        cg.start_block(&nl);
        let bare = cg.string_constant("Error");
        cg.emit(format!("br label %{jl}"));
        cg.start_block(&jl);
        let phi = cg.fresh_reg();
        cg.emit(format!(
            "{phi} = phi i8* [ {with}, %{fb} ], [ {}, %{nl} ]",
            bare.operand
        ));
        phi
    } else {
        cg.string_constant("Error").operand
    };
    let eb = cg.cur_block().to_string();
    cg.emit(format!("br label %{end}"));

    cg.start_block(&end);
    let phi = cg.fresh_reg();
    cg.emit(format!(
        "{phi} = phi i8* [ {succ}, %{sb} ], [ {err}, %{eb} ]"
    ));
    Ok(Value::new(phi, LType::Str))
}

/// `malloc(64)` + `sprintf(buf, fmt, arg)` for a single `%s` substitution,
/// returning the buffer.
fn sprintf_wrap(cg: &mut Codegen, fmt: &str, arg: &str) -> String {
    cg.add_extern("declare i8* @malloc(i64)");
    cg.add_extern("declare i32 @sprintf(i8*, i8*, ...)");
    let fmtv = cg.string_constant(fmt);
    let buf = cg.fresh_reg();
    cg.emit(format!("{buf} = call i8* @malloc(i64 64)"));
    let tmp = cg.fresh_reg();
    cg.emit(format!(
        "{tmp} = call i32 (i8*, i8*, ...) @sprintf(i8* {buf}, i8* {}, i8* {arg})",
        fmtv.operand
    ));
    buf
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
pub(crate) fn float_to_string(cg: &mut Codegen, v: &Value) -> Value {
    cg.add_extern("declare i8* @osp_float_to_string(double)");
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = call i8* @osp_float_to_string(double {})",
        v.operand
    ));
    Value::new(reg, LType::Str)
}

pub(crate) fn bool_to_string(cg: &mut Codegen, v: &Value) -> Value {
    let t = cg.string_constant("true");
    let f = cg.string_constant("false");
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = select i1 {}, i8* {}, i8* {}",
        v.operand, t.operand, f.operand
    ));
    Value::new(reg, LType::Str)
}
