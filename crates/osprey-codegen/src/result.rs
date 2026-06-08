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
    let struct_ty = format!("{{ {inner}, i8 }}");
    let obj = cg.malloc_struct(&struct_ty);
    crate::aggregate::store_field(cg, &struct_ty, obj.as_str(), 0, inner, &v.operand);
    let dp = cg.fresh_reg();
    cg.emit(format!(
        "{dp} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 1"
    ));
    cg.emit(format!("store i8 {disc}, i8* {dp}"));
    Ok(Value::result(obj, inner))
}

/// A Success result wrapping `value` (disc 0).
pub(crate) fn make_ok(cg: &mut Codegen, value: Value, inner: LType) -> Result<Value> {
    make_result(cg, value, inner, "0")
}

/// Load a Result block's `i8` discriminant operand.
pub(crate) fn load_disc(cg: &mut Codegen, v: &Value) -> String {
    let struct_ty = v
        .result_struct_ty()
        .expect("load_disc on a non-Result value");
    let dp = cg.fresh_reg();
    cg.emit(format!(
        "{dp} = getelementptr {struct_ty}, {struct_ty}* {}, i32 0, i32 1",
        v.operand
    ));
    let d = cg.fresh_reg();
    cg.emit(format!("{d} = load i8, i8* {dp}"));
    d
}

/// Load a Result block's success payload as its inner [`LType`].
pub(crate) fn load_value(cg: &mut Codegen, v: &Value) -> Value {
    let inner = v.result_inner.expect("load_value on a non-Result value");
    let struct_ty = v.result_struct_ty().unwrap();
    let loaded = crate::aggregate::load_field(cg, &struct_ty, v.operand.as_str(), 0, inner);
    Value::new(loaded, inner)
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
