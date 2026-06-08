//! Records & union variants. Each constructed value is a heap block laid out as
//! `{ i64 tag, fields… }` (the leading tag is the variant index within its
//! union, `0` for a record), handed around as an `i8*` handle that carries its
//! Osprey owner type so field access and `match` can recover the layout. Ports
//! the record/union construction + field-access paths of
//! `expression_generation.go`.

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, FieldAssignment};

/// `Type { field: value, … }` — allocate the heap block, write the tag and each
/// declared field (in layout order), and return the owner-tagged handle.
pub(crate) fn gen_constructor(
    cg: &mut Codegen,
    name: &str,
    fields: &[FieldAssignment],
) -> Result<Value> {
    // A `name { … }` where `name` is a bound variable is a record *update*, not
    // a constructor (the parser cannot tell them apart).
    if !cg.is_ctor(name) {
        if cg.lookup(name).is_some() {
            return gen_update(cg, name, fields);
        }
        return Err(CodegenError::unknown(name));
    }
    // `Success { value: x }` / `Error { message: m }` build the Result ABI block
    // `{ inner, i8 disc }` directly (disc 0 = Success), not a generic record —
    // so they interoperate with `match`, `toString` and effect handlers that
    // return `Result<…>` (e.g. an `input => Success { value: … }` handler arm).
    // The field disambiguates from a same-named *nullary* union variant (e.g.
    // `type TaskResult = Success | …`), which takes the ordinary union path.
    if (name == "Success" || name == "Error") && !fields.is_empty() {
        return gen_result_ctor(cg, name, fields);
    }
    let view = cg
        .ctor_layout(name)
        .ok_or_else(|| CodegenError::unknown(name))?;
    let struct_ty = cg
        .ctor_struct_ty(name)
        .ok_or_else(|| CodegenError::unknown(name))?;
    let obj = cg.malloc_struct(&struct_ty);

    // tag
    let tagp = cg.fresh_reg();
    cg.emit(format!(
        "{tagp} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 0"
    ));
    cg.emit(format!("store i64 {}, i64* {tagp}", view.tag));

    // fields, in declared order
    for (i, (fname, fty)) in view.fields.iter().enumerate() {
        let fa = fields.iter().find(|f| &f.name == fname).ok_or_else(|| {
            CodegenError::invalid(format!("missing field `{fname}` for `{name}`"))
        })?;
        let v = gen_expr(cg, &fa.value)?;
        let v = crate::cast::coerce_to(cg, v, *fty)?;
        store_field(cg, &struct_ty, obj.as_str(), i + 1, *fty, &v.operand);
    }

    let handle = cg.fresh_reg();
    cg.emit(format!("{handle} = bitcast {struct_ty}* {obj} to i8*"));
    Ok(Value::handle(handle, view.owner))
}

/// `{ field: value, … }` — an anonymous object literal: the same `{ i64 tag,
/// fields… }` heap block as a named record, with a synthetic layout registered so
/// field access can recover the slots.
pub(crate) fn gen_object(cg: &mut Codegen, fields: &[FieldAssignment]) -> Result<Value> {
    let mut vals = Vec::with_capacity(fields.len());
    for fa in fields {
        let v = gen_expr(cg, &fa.value)?;
        vals.push((fa.name.clone(), v));
    }
    let mut parts = vec!["i64".to_string()];
    parts.extend(vals.iter().map(|(_, v)| v.ty.as_str().to_string()));
    let struct_ty = format!("{{ {} }}", parts.join(", "));
    let layout: Vec<(String, LType)> = vals.iter().map(|(n, v)| (n.clone(), v.ty)).collect();
    let owner = cg.register_obj_layout(layout);

    let obj = cg.malloc_struct(&struct_ty);
    let tagp = cg.fresh_reg();
    cg.emit(format!(
        "{tagp} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 0"
    ));
    cg.emit(format!("store i64 0, i64* {tagp}"));
    for (i, (_, v)) in vals.iter().enumerate() {
        store_field(cg, &struct_ty, obj.as_str(), i + 1, v.ty, &v.operand);
    }
    let handle = cg.fresh_reg();
    cg.emit(format!("{handle} = bitcast {struct_ty}* {obj} to i8*"));
    Ok(Value::handle(handle, owner))
}

/// Build a `Success`/`Error` value in the Result ABI: the single field becomes
/// the `{ inner, i8 }` block's payload, with disc `0` (Success) or `1` (Error).
fn gen_result_ctor(cg: &mut Codegen, name: &str, fields: &[FieldAssignment]) -> Result<Value> {
    let fa = fields
        .first()
        .ok_or_else(|| CodegenError::invalid(format!("`{name}` needs one field")))?;
    let v = gen_expr(cg, &fa.value)?;
    let inner = v.ty;
    let disc = if name == "Success" { "0" } else { "1" };
    crate::result::make_result(cg, v, inner, disc)
}

/// `record { field: newValue }` — copy every field of `record` into a fresh
/// block, overriding the named ones.
pub(crate) fn gen_update(
    cg: &mut Codegen,
    record: &str,
    fields: &[FieldAssignment],
) -> Result<Value> {
    let base = cg
        .lookup(record)
        .ok_or_else(|| CodegenError::unknown(record))?;
    let owner = base
        .osp_ty
        .clone()
        .ok_or_else(|| CodegenError::invalid(format!("`{record}` is not a record")))?;
    let view = cg
        .ctor_layout(&owner)
        .ok_or_else(|| CodegenError::unknown(&owner))?;
    let struct_ty = cg
        .ctor_struct_ty(&owner)
        .ok_or_else(|| CodegenError::unknown(&owner))?;

    let src = cg.fresh_reg();
    cg.emit(format!(
        "{src} = bitcast i8* {} to {struct_ty}*",
        base.operand
    ));
    let obj = cg.malloc_struct(&struct_ty);

    let tagp = cg.fresh_reg();
    cg.emit(format!(
        "{tagp} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 0"
    ));
    cg.emit(format!("store i64 {}, i64* {tagp}", view.tag));

    for (i, (fname, fty)) in view.fields.iter().enumerate() {
        let val = match fields.iter().find(|f| &f.name == fname) {
            Some(fa) => {
                let v = gen_expr(cg, &fa.value)?;
                crate::cast::coerce_to(cg, v, *fty)?.operand
            }
            None => load_field(cg, &struct_ty, src.as_str(), i + 1, *fty),
        };
        store_field(cg, &struct_ty, obj.as_str(), i + 1, *fty, &val);
    }

    let handle = cg.fresh_reg();
    cg.emit(format!("{handle} = bitcast {struct_ty}* {obj} to i8*"));
    Ok(Value::handle(handle, owner))
}

/// `obj.field` — recover the record layout from the handle's owner type and
/// load the field.
pub(crate) fn gen_field_access(cg: &mut Codegen, target: &Expr, field: &str) -> Result<Value> {
    let tv = gen_expr(cg, target)?;
    // Use the statically-known owner (a named record or an anonymous object
    // literal) when it actually declares `field`; otherwise (a generic accessor
    // whose parameter infers to a type variable) resolve the field by name across
    // known layouts — Go's polymorphic field-access fallback.
    let known = tv.osp_ty.clone().filter(|o| {
        cg.record_layout(o)
            .is_some_and(|(_, fs)| fs.iter().any(|(f, _)| f == field))
    });
    let owner = known
        .or_else(|| cg.find_field_owner(field))
        .ok_or_else(|| CodegenError::invalid(format!("field `{field}` on a non-record")))?;
    let (struct_ty, fields_layout) = cg
        .record_layout(&owner)
        .ok_or_else(|| CodegenError::unknown(&owner))?;
    let (idx, fty) = fields_layout
        .iter()
        .enumerate()
        .find_map(|(i, (f, t))| (f == field).then_some((i, *t)))
        .ok_or_else(|| CodegenError::invalid(format!("`{owner}` has no field `{field}`")))?;

    let src = cg.fresh_reg();
    cg.emit(format!(
        "{src} = bitcast i8* {} to {struct_ty}*",
        tv.operand
    ));
    let loaded = load_field(cg, &struct_ty, src.as_str(), idx + 1, fty);
    let owner = cg.ctor_field_written(&owner, field);
    Ok(Value::new(loaded, fty).with_owner(owner))
}

/// Store `val` (LLVM type `fty`) into the `idx`-th element of a `{TY}*` block.
pub(crate) fn store_field(
    cg: &mut Codegen,
    struct_ty: &str,
    obj: &str,
    idx: usize,
    fty: LType,
    val: &str,
) {
    let p = cg.fresh_reg();
    cg.emit(format!(
        "{p} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 {idx}"
    ));
    cg.emit(format!("store {fty} {val}, {fty}* {p}"));
}

/// Load the `idx`-th element of a `{TY}*` block, returning the value register.
pub(crate) fn load_field(
    cg: &mut Codegen,
    struct_ty: &str,
    obj: &str,
    idx: usize,
    fty: LType,
) -> String {
    let p = cg.fresh_reg();
    cg.emit(format!(
        "{p} = getelementptr {struct_ty}, {struct_ty}* {obj}, i32 0, i32 {idx}"
    ));
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = load {fty}, {fty}* {p}"));
    r
}
