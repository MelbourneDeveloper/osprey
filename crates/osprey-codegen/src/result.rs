//! The `Result<T, E>` ABI: a heap block `{ T value, i8 disc }` reached by
//! pointer, `disc == 0` ⇒ Success. Ports the `{value, i8}` layout the Go
//! backend builds in `getResultType` / `wrapInMathResult` and reads in
//! `generateResultMatchCondition` / `convertResultToString`. Runtime fallible
//! builtins (list/map get, string ops) and user functions declared
//! `-> Result<…>` both produce this shape, so match, `toString` and value-site
//! coercion handle exactly one representation.

use crate::builder::Codegen;
use crate::cast::coerce_to;
use crate::error::Result;
use crate::llty::{LType, Value};

/// Build a `Result` block with the given success `value` and an explicit `i8`
/// discriminant operand (`"0"` Success, `"1"` Error, or an `i8` register from a
/// `select`). The value is coerced to `inner` before storing.
pub(crate) fn make_result(
    cg: &mut Codegen,
    value: Value,
    inner: LType,
    disc: &str,
) -> Result<Value> {
    let v = coerce_to(cg, value, inner)?;
    let payload_owner = v.osp_ty.clone();
    let struct_ty = format!("{{ {inner}, i8 }}");
    let obj = cg.malloc_struct(&struct_ty);
    crate::aggregate::store_field(cg, &struct_ty, obj.as_str(), 0, inner, &v.operand);
    let dp = cg.fresh_reg();
    cg.emit(format!(
        "{dp} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 1"
    ));
    cg.emit(format!("store i8 {disc}, i8* {dp}"));
    Ok(Value::result(obj, inner).with_payload_owner(payload_owner))
}

/// A Success result wrapping `value` (disc 0).
pub(crate) fn make_ok(cg: &mut Codegen, value: Value, inner: LType) -> Result<Value> {
    make_result(cg, value, inner, "0")
}

/// Build a `Result` whose discriminant is Error when `is_err` (an `i1` operand)
/// holds — folding the ubiquitous `select i1 …, i8 1, i8 0` then [`make_result`]
/// that every fallible runtime builtin ends with.
pub(crate) fn make_result_if_err(
    cg: &mut Codegen,
    value: Value,
    inner: LType,
    is_err: &str,
) -> Result<Value> {
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {is_err}, i8 1, i8 0"));
    make_result(cg, value, inner, &disc)
}

/// `Result<i64, _>` from a runtime `i32` success flag (`0` ⇒ Error) guarding an
/// `i64` payload — the shared shape of `listGet` / `mapGet`.
pub(crate) fn result_from_flag(cg: &mut Codegen, flag: &str, value: &str) -> Result<Value> {
    let err = cg.emit_reg(format!("icmp eq i32 {flag}, 0"));
    make_result_if_err(cg, Value::new(value, LType::I64), LType::I64, &err)
}

/// `Result<int, _>` from a C `i64` whose negative values signal failure — the
/// uniform convention of the file/process/HTTP/JSON runtime (a negative handle,
/// byte count, status or process id is Error). The success value carried is the
/// result itself.
pub(crate) fn result_from_i64(cg: &mut Codegen, result: &str) -> Result<Value> {
    let err = cg.emit_reg(format!("icmp slt i64 {result}, 0"));
    make_result_if_err(cg, Value::new(result, LType::I64), LType::I64, &err)
}

/// `Result<string, _>` from a possibly-NULL C `char*` (`ptr` an `i8*` operand):
/// NULL ⇒ Error, else Success. When `err` is `Some(msg)`, the error payload is
/// that constant string so `toString` shows `Error(msg)` (e.g. `readFile`'s
/// `Error(File read error)`); otherwise the bare (null) pointer is kept and a
/// plain `Error` is shown.
pub(crate) fn result_from_nullable(cg: &mut Codegen, ptr: &str, err: Option<&str>) -> Result<Value> {
    let is_null = cg.emit_reg(format!("icmp eq i8* {ptr}, null"));
    let value = match err {
        Some(msg) => {
            let c = cg.string_constant(msg);
            cg.emit_reg(format!("select i1 {is_null}, i8* {}, i8* {ptr}", c.operand))
        }
        None => ptr.to_string(),
    };
    make_result_if_err(cg, Value::new(value, LType::Str), LType::Str, &is_null)
}

/// Branch on a Result's discriminant: load it, test `== 0` (Success), and emit
/// the conditional branch to fresh `(success, error, end)` labels — leaving the
/// builder positioned at the start of the `success` block. The shared preamble
/// of every "do one thing on Success, another on Error, `phi` the results" path.
pub(crate) fn open_result_branch(cg: &mut Codegen, v: &Value) -> (String, String, String) {
    let d = load_disc(cg, v);
    let is_succ = cg.emit_reg(format!("icmp eq i8 {d}, 0"));
    let sl = cg.fresh_label();
    let el = cg.fresh_label();
    let end = cg.fresh_label();
    cg.emit(format!("br i1 {is_succ}, label %{sl}, label %{el}"));
    cg.start_block(&sl);
    (sl, el, end)
}

/// Load a Result block's `i8` discriminant operand. Invariant: `v` is a Result
/// (callers gate on `result_inner.is_some()`); a non-Result yields the Error
/// discriminant `1` rather than panicking.
pub(crate) fn load_disc(cg: &mut Codegen, v: &Value) -> String {
    let Some(struct_ty) = v.result_struct_ty() else {
        return "1".to_string();
    };
    let dp = cg.fresh_reg();
    cg.emit(format!(
        "{dp} = getelementptr {struct_ty}, {struct_ty}* {}, i32 0, i32 1",
        v.operand
    ));
    let d = cg.fresh_reg();
    cg.emit(format!("{d} = load i8, i8* {dp}"));
    d
}

/// Load a Result block's success payload as its inner [`LType`]. Invariant: `v`
/// is a Result; a non-Result yields Unit rather than panicking.
pub(crate) fn load_value(cg: &mut Codegen, v: &Value) -> Value {
    let Some(inner) = v.result_inner else {
        return Value::unit();
    };
    let struct_ty = format!("{{ {inner}, i8 }}");
    let loaded = crate::aggregate::load_field(cg, &struct_ty, v.operand.as_str(), 0, inner);
    Value::new(loaded, inner).with_owner(v.payload_owner.clone())
}

/// Auto-unwrap a Result at a value site (arithmetic, `print`, an argument),
/// yielding its success payload; a non-Result value passes through. Mirrors Go's
/// `unwrapIfResult`.
pub(crate) fn unwrap(cg: &mut Codegen, v: Value) -> Value {
    if v.result_inner.is_some() {
        load_value(cg, &v)
    } else {
        v
    }
}
