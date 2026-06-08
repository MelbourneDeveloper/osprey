//! List<T> and Map<K,V> builtins backed by the C runtime (`osprey_list_*` /
//! `osprey_map_*` in `libfiber_runtime`). Ports `collection_codegen.go`.
//! Element values cross the boundary as a uniform `i64`; pointers are
//! `ptrtoint`-boxed. List/Map handles are `i8*` tagged with their owner so the
//! `+` operator and `toString` can tell them from records. Implements
//! [TYPE-LIST-OPS], [TYPE-MAP-OPS].

use crate::builder::Codegen;
use crate::cast::coerce_to;
use crate::conv::box_to_i64;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use crate::result::make_result;
use osprey_ast::{Expr, NamedArgument};

/// The owner tag carried by runtime list / map handles.
pub(crate) const LIST_OWNER: &str = "List";
pub(crate) const MAP_OWNER: &str = "Map";

/// Dispatch a collection builtin by name, or `None` if `name` is not one.
pub(crate) fn gen(
    cg: &mut Codegen,
    name: &str,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Option<Value>> {
    let v = match name {
        "List" => list_empty(cg)?,
        "listLength" => one_list_i64(cg, "osprey_list_length", args)?,
        "listAppend" => list_box2(cg, "osprey_list_append", args)?,
        "listPrepend" => list_box2(cg, "osprey_list_prepend", args)?,
        "listDrop" => list_box2(cg, "osprey_list_drop", args)?,
        "listConcat" => list_concat(cg, args)?,
        "listReverse" => one_list_handle(cg, "osprey_list_reverse", args)?,
        "listGet" => list_get(cg, args)?,
        "listContains" => list_contains(cg, args)?,
        "Map" => map_empty(cg)?,
        "mapLength" => one_list_i64(cg, "osprey_map_length", args)?,
        "mapSet" => map_set(cg, args)?,
        "mapGet" => map_get(cg, args)?,
        "mapContains" => map_contains(cg, args)?,
        "mapRemove" => map_remove(cg, args)?,
        "mapMerge" => map_merge(cg, args)?,
        "mapKeys" => map_to_list(cg, args, true)?,
        "mapValues" => map_to_list(cg, args, false)?,
        _ => return Ok(None),
    };
    Ok(Some(v))
}

fn declare(cg: &mut Codegen, cname: &str, ret: &str, params: &str) {
    cg.add_extern(format!("declare {ret} @{cname}({params})"));
}

/// The `i`-th positional argument as an opaque `i8*` collection handle.
fn handle_arg(cg: &mut Codegen, args: &[Expr], i: usize) -> Result<Value> {
    let e = args
        .get(i)
        .ok_or_else(|| CodegenError::invalid("collection builtin: missing argument"))?;
    let v = gen_expr(cg, e)?;
    coerce_to(cg, v, LType::Ptr)
}

/// The `i`-th positional argument, boxed to the uniform `i64` element ABI.
fn boxed_arg(cg: &mut Codegen, args: &[Expr], i: usize) -> Result<Value> {
    let e = args
        .get(i)
        .ok_or_else(|| CodegenError::invalid("collection builtin: missing argument"))?;
    let v = gen_expr(cg, e)?;
    let v = crate::result::unwrap(cg, v);
    Ok(box_to_i64(cg, v))
}

fn list_empty(cg: &mut Codegen) -> Result<Value> {
    declare(cg, "osprey_list_empty", "i8*", "");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @osprey_list_empty()"));
    Ok(Value::handle(r, LIST_OWNER))
}

fn map_empty(cg: &mut Codegen) -> Result<Value> {
    declare(cg, "osprey_map_empty", "i8*", "i32");
    let r = cg.fresh_reg();
    // OSPREY_KEY_STRING = 1 (Map() defaults to string keys).
    cg.emit(format!("{r} = call i8* @osprey_map_empty(i32 1)"));
    Ok(Value::handle(r, MAP_OWNER))
}

/// `f(handle) -> int`.
fn one_list_i64(cg: &mut Codegen, cname: &str, args: &[Expr]) -> Result<Value> {
    let h = handle_arg(cg, args, 0)?;
    declare(cg, cname, "i64", "i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i64 @{cname}(i8* {})", h.operand));
    Ok(Value::new(r, LType::I64))
}

/// `f(handle) -> handle`.
fn one_list_handle(cg: &mut Codegen, cname: &str, args: &[Expr]) -> Result<Value> {
    let h = handle_arg(cg, args, 0)?;
    declare(cg, cname, "i8*", "i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @{cname}(i8* {})", h.operand));
    Ok(Value::handle(r, LIST_OWNER))
}

/// `f(handle, boxed) -> handle` (append / prepend / drop).
fn list_box2(cg: &mut Codegen, cname: &str, args: &[Expr]) -> Result<Value> {
    let h = handle_arg(cg, args, 0)?;
    let x = boxed_arg(cg, args, 1)?;
    declare(cg, cname, "i8*", "i8*, i64");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @{cname}(i8* {}, i64 {})", h.operand, x.operand));
    Ok(Value::handle(r, LIST_OWNER))
}

/// `listConcat(a, b) -> List` — also the lowering of `a + b` on lists.
fn list_concat(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let a = handle_arg(cg, args, 0)?;
    let b = handle_arg(cg, args, 1)?;
    concat_handles(cg, &a, &b)
}

/// Emit `osprey_list_concat` on two already-evaluated list handles.
pub(crate) fn concat_handles(cg: &mut Codegen, a: &Value, b: &Value) -> Result<Value> {
    declare(cg, "osprey_list_concat", "i8*", "i8*, i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @osprey_list_concat(i8* {}, i8* {})", a.operand, b.operand));
    Ok(Value::handle(r, LIST_OWNER))
}

/// `listGet(l, i) -> Result<T, _>` gated on `osprey_list_in_bounds`.
fn list_get(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let l = handle_arg(cg, args, 0)?;
    let i = boxed_arg(cg, args, 1)?;
    declare(cg, "osprey_list_in_bounds", "i32", "i8*, i64");
    declare(cg, "osprey_list_get", "i64", "i8*, i64");
    let inb = cg.fresh_reg();
    cg.emit(format!("{inb} = call i32 @osprey_list_in_bounds(i8* {}, i64 {})", l.operand, i.operand));
    let val = cg.fresh_reg();
    cg.emit(format!("{val} = call i64 @osprey_list_get(i8* {}, i64 {})", l.operand, i.operand));
    let oob = cg.fresh_reg();
    cg.emit(format!("{oob} = icmp eq i32 {inb}, 0"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {oob}, i8 1, i8 0"));
    make_result(cg, Value::new(val, LType::I64), LType::I64, &disc)
}

/// `listContains(l, x) -> bool`: linear scan, content-equality for strings.
fn list_contains(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let l = handle_arg(cg, args, 0)?;
    let needle_e = args
        .get(1)
        .ok_or_else(|| CodegenError::invalid("listContains: missing argument"))?;
    let needle = gen_expr(cg, needle_e)?;
    let needle = crate::result::unwrap(cg, needle);
    let is_str = needle.ty == LType::Str;
    let boxed = box_to_i64(cg, needle.clone());

    declare(cg, "osprey_list_length", "i64", "i8*");
    declare(cg, "osprey_list_get", "i64", "i8*, i64");
    let len = cg.fresh_reg();
    cg.emit(format!("{len} = call i64 @osprey_list_length(i8* {})", l.operand));
    let idx = cg.fresh_reg();
    cg.emit(format!("{idx} = alloca i64"));
    cg.emit(format!("store i64 0, i64* {idx}"));
    let res = cg.fresh_reg();
    cg.emit(format!("{res} = alloca i1"));
    cg.emit(format!("store i1 0, i1* {res}"));

    let loop_l = cg.fresh_label();
    let body = cg.fresh_label();
    let found = cg.fresh_label();
    let cont = cg.fresh_label();
    let done = cg.fresh_label();
    cg.emit(format!("br label %{loop_l}"));

    cg.start_block(&loop_l);
    let i = cg.fresh_reg();
    cg.emit(format!("{i} = load i64, i64* {idx}"));
    let more = cg.fresh_reg();
    cg.emit(format!("{more} = icmp slt i64 {i}, {len}"));
    cg.emit(format!("br i1 {more}, label %{body}, label %{done}"));

    cg.start_block(&body);
    let elem = cg.fresh_reg();
    cg.emit(format!("{elem} = call i64 @osprey_list_get(i8* {}, i64 {i})", l.operand));
    let eq = cg.fresh_reg();
    if is_str {
        cg.add_extern("declare i32 @strcmp(i8*, i8*)");
        let ep = cg.fresh_reg();
        cg.emit(format!("{ep} = inttoptr i64 {elem} to i8*"));
        let c = cg.fresh_reg();
        cg.emit(format!("{c} = call i32 @strcmp(i8* {ep}, i8* {})", needle.operand));
        cg.emit(format!("{eq} = icmp eq i32 {c}, 0"));
    } else {
        cg.emit(format!("{eq} = icmp eq i64 {elem}, {}", boxed.operand));
    }
    cg.emit(format!("br i1 {eq}, label %{found}, label %{cont}"));

    cg.start_block(&found);
    cg.emit(format!("store i1 1, i1* {res}"));
    cg.emit(format!("br label %{done}"));

    cg.start_block(&cont);
    let next = cg.fresh_reg();
    cg.emit(format!("{next} = add i64 {i}, 1"));
    cg.emit(format!("store i64 {next}, i64* {idx}"));
    cg.emit(format!("br label %{loop_l}"));

    cg.start_block(&done);
    let out = cg.fresh_reg();
    cg.emit(format!("{out} = load i1, i1* {res}"));
    Ok(Value::new(out, LType::I1))
}

/// `mapSet(m, k, v) -> Map`.
fn map_set(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let m = handle_arg(cg, args, 0)?;
    let k = boxed_arg(cg, args, 1)?;
    let v = boxed_arg(cg, args, 2)?;
    declare(cg, "osprey_map_set", "i8*", "i8*, i64, i64");
    let r = cg.fresh_reg();
    cg.emit(format!(
        "{r} = call i8* @osprey_map_set(i8* {}, i64 {}, i64 {})",
        m.operand, k.operand, v.operand
    ));
    Ok(Value::handle(r, MAP_OWNER))
}

/// `mapRemove(m, k) -> Map`.
fn map_remove(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let m = handle_arg(cg, args, 0)?;
    let k = boxed_arg(cg, args, 1)?;
    declare(cg, "osprey_map_remove", "i8*", "i8*, i64");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @osprey_map_remove(i8* {}, i64 {})", m.operand, k.operand));
    Ok(Value::handle(r, MAP_OWNER))
}

/// `mapContains(m, k) -> bool`.
fn map_contains(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let m = handle_arg(cg, args, 0)?;
    let k = boxed_arg(cg, args, 1)?;
    declare(cg, "osprey_map_contains", "i32", "i8*, i64");
    let raw = cg.fresh_reg();
    cg.emit(format!("{raw} = call i32 @osprey_map_contains(i8* {}, i64 {})", m.operand, k.operand));
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = icmp ne i32 {raw}, 0"));
    Ok(Value::new(r, LType::I1))
}

/// `mapGet(m, k) -> Result<V, _>` gated on `osprey_map_contains`.
fn map_get(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let m = handle_arg(cg, args, 0)?;
    let k = boxed_arg(cg, args, 1)?;
    runtime_map_get(cg, &m, &k)
}

/// `mapMerge(a, b) -> Map` (right-biased) — also the lowering of `a + b` on maps.
fn map_merge(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let a = handle_arg(cg, args, 0)?;
    let b = handle_arg(cg, args, 1)?;
    merge_handles(cg, &a, &b)
}

/// Emit `osprey_map_merge` on two already-evaluated map handles.
pub(crate) fn merge_handles(cg: &mut Codegen, a: &Value, b: &Value) -> Result<Value> {
    declare(cg, "osprey_map_merge", "i8*", "i8*, i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @osprey_map_merge(i8* {}, i8* {})", a.operand, b.operand));
    Ok(Value::handle(r, MAP_OWNER))
}

/// `mapKeys`/`mapValues` → a `List` built by iterating the map.
fn map_to_list(cg: &mut Codegen, args: &[Expr], take_key: bool) -> Result<Value> {
    let m = handle_arg(cg, args, 0)?;
    declare(cg, "osprey_map_iter_new", "i8*", "i8*");
    declare(cg, "osprey_map_iter_next", "i32", "i8*, i64*, i64*");
    declare(cg, "osprey_list_builder_new", "i8*", "");
    declare(cg, "osprey_list_builder_push", "void", "i8*, i64");
    declare(cg, "osprey_list_builder_seal", "i8*", "i8*");
    let bld = cg.fresh_reg();
    cg.emit(format!("{bld} = call i8* @osprey_list_builder_new()"));
    let iter = cg.fresh_reg();
    cg.emit(format!("{iter} = call i8* @osprey_map_iter_new(i8* {})", m.operand));
    let kp = cg.fresh_reg();
    cg.emit(format!("{kp} = alloca i64"));
    let vp = cg.fresh_reg();
    cg.emit(format!("{vp} = alloca i64"));

    let cond = cg.fresh_label();
    let body = cg.fresh_label();
    let endl = cg.fresh_label();
    cg.emit(format!("br label %{cond}"));

    cg.start_block(&cond);
    let has = cg.fresh_reg();
    cg.emit(format!("{has} = call i32 @osprey_map_iter_next(i8* {iter}, i64* {kp}, i64* {vp})"));
    let more = cg.fresh_reg();
    cg.emit(format!("{more} = icmp ne i32 {has}, 0"));
    cg.emit(format!("br i1 {more}, label %{body}, label %{endl}"));

    cg.start_block(&body);
    let slot = if take_key { &kp } else { &vp };
    let elem = cg.fresh_reg();
    cg.emit(format!("{elem} = load i64, i64* {slot}"));
    cg.emit(format!("call void @osprey_list_builder_push(i8* {bld}, i64 {elem})"));
    cg.emit(format!("br label %{cond}"));

    cg.start_block(&endl);
    let sealed = cg.fresh_reg();
    cg.emit(format!("{sealed} = call i8* @osprey_list_builder_seal(i8* {bld})"));
    Ok(Value::handle(sealed, LIST_OWNER))
}

/// `{ k: v, … }` — build a runtime map (string keys) via the map builder.
pub(crate) fn gen_map_literal(cg: &mut Codegen, entries: &[osprey_ast::MapEntry]) -> Result<Value> {
    declare(cg, "osprey_map_builder_new", "i8*", "i32");
    declare(cg, "osprey_map_builder_put", "void", "i8*, i64, i64");
    declare(cg, "osprey_map_builder_seal", "i8*", "i8*");
    let bld = cg.fresh_reg();
    // OSPREY_KEY_STRING = 1.
    cg.emit(format!("{bld} = call i8* @osprey_map_builder_new(i32 1)"));
    for e in entries {
        let k = gen_expr(cg, &e.key)?;
        let k = crate::result::unwrap(cg, k);
        let k = box_to_i64(cg, k);
        let v = gen_expr(cg, &e.value)?;
        let v = crate::result::unwrap(cg, v);
        let v = box_to_i64(cg, v);
        cg.emit(format!(
            "call void @osprey_map_builder_put(i8* {bld}, i64 {}, i64 {})",
            k.operand, v.operand
        ));
    }
    let sealed = cg.fresh_reg();
    cg.emit(format!("{sealed} = call i8* @osprey_map_builder_seal(i8* {bld})"));
    Ok(Value::handle(sealed, MAP_OWNER))
}

/// Shared runtime map lookup → `Result<i64, _>` (also used by `m[key]` indexing).
pub(crate) fn runtime_map_get(cg: &mut Codegen, m: &Value, k: &Value) -> Result<Value> {
    declare(cg, "osprey_map_contains", "i32", "i8*, i64");
    declare(cg, "osprey_map_get", "i64", "i8*, i64");
    let has = cg.fresh_reg();
    cg.emit(format!("{has} = call i32 @osprey_map_contains(i8* {}, i64 {})", m.operand, k.operand));
    let got = cg.fresh_reg();
    cg.emit(format!("{got} = call i64 @osprey_map_get(i8* {}, i64 {})", m.operand, k.operand));
    let miss = cg.fresh_reg();
    cg.emit(format!("{miss} = icmp eq i32 {has}, 0"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {miss}, i8 1, i8 0"));
    make_result(cg, Value::new(got, LType::I64), LType::I64, &disc)
}
