//! Iterator builtins: integer `range`, the stream-fused higher-order operations
//! (`map`/`filter`/`forEach`/`fold`) and the eager list operations
//! (`forEachList`/`mapList`/`filterList`/`foldList`). Ports
//! `iterator_generation.go` + the list-loop generators in
//! `collection_codegen.go`. A range is a stack `{ i64, i64 }` (start, end);
//! `map`/`filter` record a pending stage and pass the range through; the
//! consuming `forEach`/`fold` emits one counted loop replaying those stages.
//! Implements [BUILTIN-ITER-*].

use crate::builder::Codegen;
use crate::conv::{as_i64, box_to_i64};
use crate::error::{CodegenError, Result};
use crate::expr::{call_with_values, gen_expr};
use crate::llty::{LType, Value};
use osprey_ast::{Expr, NamedArgument};

const RANGE_TY: &str = "{ i64, i64 }";
const RANGE_OWNER: &str = "Range";

/// A recorded stream-fusion stage.
#[derive(Clone)]
pub(crate) struct IterOp {
    pub map: bool, // true = map (transform), false = filter (predicate)
    pub fn_name: String,
}

/// Dispatch an iterator builtin by name, or `None` if `name` is not one.
pub(crate) fn gen(
    cg: &mut Codegen,
    name: &str,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Option<Value>> {
    let v = match name {
        "range" => range(cg, args)?,
        "map" => record(cg, args, true)?,
        "filter" => record(cg, args, false)?,
        "forEach" => for_each(cg, args)?,
        "fold" => fold(cg, args)?,
        "forEachList" => for_each_list(cg, args)?,
        "mapList" => list_builder(cg, args, false)?,
        "filterList" => list_builder(cg, args, true)?,
        "foldList" => fold_list(cg, args)?,
        _ => return Ok(None),
    };
    Ok(Some(v))
}

/// The callback function name from an iterator argument (a bare identifier).
fn callback_name(e: &Expr) -> Result<String> {
    match e {
        Expr::Identifier(n) => Ok(n.clone()),
        _ => Err(CodegenError::unsupported(
            "iterator callback must be a named function",
        )),
    }
}

fn nth(args: &[Expr], i: usize) -> Result<&Expr> {
    args.get(i)
        .ok_or_else(|| CodegenError::invalid("iterator builtin: missing argument"))
}

/// `range(start, end)` → a `{ start, end }` block (half-open, step 1).
fn range(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let s = gen_expr(cg, nth(args, 0)?)?;
    let s = as_i64(cg, s)?;
    let e = gen_expr(cg, nth(args, 1)?)?;
    let e = as_i64(cg, e)?;
    let obj = cg.malloc_struct(RANGE_TY);
    crate::aggregate::store_field(cg, RANGE_TY, &obj, 0, LType::I64, &s.operand);
    crate::aggregate::store_field(cg, RANGE_TY, &obj, 1, LType::I64, &e.operand);
    Ok(Value::handle(obj, RANGE_OWNER))
}

/// `map`/`filter`: record a pending stage and return the iterator unchanged.
fn record(cg: &mut Codegen, args: &[Expr], is_map: bool) -> Result<Value> {
    let iter = gen_expr(cg, nth(args, 0)?)?;
    let fn_name = callback_name(nth(args, 1)?)?;
    cg.pending_iter_ops.push(IterOp {
        map: is_map,
        fn_name,
    });
    Ok(iter)
}

/// Load a range block's `(start, end)` bounds.
fn bounds(cg: &mut Codegen, range: &Value) -> (String, String) {
    let s = crate::aggregate::load_field(cg, RANGE_TY, &range.operand, 0, LType::I64);
    let e = crate::aggregate::load_field(cg, RANGE_TY, &range.operand, 1, LType::I64);
    (s, e)
}

/// Replay the pending map/filter stages on element `v` in the current block,
/// branching to `skip` when a filter rejects it. Returns the transformed value.
fn replay(cg: &mut Codegen, v: Value, skip: &str) -> Result<Value> {
    let ops = std::mem::take(&mut cg.pending_iter_ops);
    let mut cur = v;
    for op in &ops {
        if op.map {
            cur = call_with_values(cg, &op.fn_name, vec![cur])?;
            cur = crate::result::unwrap(cg, cur);
        } else {
            let pred = call_with_values(cg, &op.fn_name, vec![cur.clone()])?;
            let pred = crate::result::unwrap(cg, pred);
            let pb = as_i64(cg, pred)?;
            let nz = cg.fresh_reg();
            cg.emit(format!("{nz} = icmp ne i64 {}, 0", pb.operand));
            let pass = cg.fresh_label();
            cg.emit(format!("br i1 {nz}, label %{pass}, label %{skip}"));
            cg.start_block(&pass);
        }
    }
    Ok(cur)
}

/// A counted loop over a range's half-open `[start, end)`, the range analogue of
/// [`ListLoop`]. After [`open_range_loop`] the builder sits in the body with the
/// current index in `i`; [`close_range_loop`] emits the increment + back-edge.
struct RangeLoop {
    i: String,
    ctr: String,
    cond: String,
    incr: String,
    endl: String,
}

fn open_range_loop(cg: &mut Codegen, start: &str, end: &str) -> RangeLoop {
    let ctr = cg.fresh_reg();
    cg.emit(format!("{ctr} = alloca i64"));
    cg.emit(format!("store i64 {start}, i64* {ctr}"));
    let cond = cg.fresh_label();
    let body = cg.fresh_label();
    let incr = cg.fresh_label();
    let endl = cg.fresh_label();
    cg.emit(format!("br label %{cond}"));

    cg.start_block(&cond);
    let i = cg.fresh_reg();
    cg.emit(format!("{i} = load i64, i64* {ctr}"));
    let more = cg.fresh_reg();
    cg.emit(format!("{more} = icmp slt i64 {i}, {end}"));
    cg.emit(format!("br i1 {more}, label %{body}, label %{endl}"));

    cg.start_block(&body);
    RangeLoop {
        i,
        ctr,
        cond,
        incr,
        endl,
    }
}

fn close_range_loop(cg: &mut Codegen, lp: &RangeLoop) {
    cg.emit(format!("br label %{}", lp.incr));
    cg.start_block(&lp.incr);
    let next = cg.fresh_reg();
    cg.emit(format!("{next} = add i64 {}, 1", lp.i));
    cg.emit(format!("store i64 {next}, i64* {}", lp.ctr));
    cg.emit(format!("br label %{}", lp.cond));
    cg.start_block(&lp.endl);
}

/// `forEach(iterator, fn)` — counted loop applying `fn` to each (fused) element.
fn for_each(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let range = gen_expr(cg, nth(args, 0)?)?;
    let consumer = callback_name(nth(args, 1)?)?;
    let (start, end) = bounds(cg, &range);

    let lp = open_range_loop(cg, &start, &end);
    let elem = replay(cg, Value::new(lp.i.clone(), LType::I64), &lp.incr)?;
    let _ = call_with_values(cg, &consumer, vec![elem])?;
    close_range_loop(cg, &lp);
    Ok(range)
}

/// `fold(iterator, initial, fn)` — counted loop accumulating `fn(acc, elem)`.
fn fold(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let range = gen_expr(cg, nth(args, 0)?)?;
    let initial = gen_expr(cg, nth(args, 1)?)?;
    let initial = crate::result::unwrap(cg, initial);
    let initial = box_to_i64(cg, initial);
    let combine = callback_name(nth(args, 2)?)?;
    let (start, end) = bounds(cg, &range);

    let acc = cg.fresh_reg();
    cg.emit(format!("{acc} = alloca i64"));
    cg.emit(format!("store i64 {}, i64* {acc}", initial.operand));

    let lp = open_range_loop(cg, &start, &end);
    let elem = replay(cg, Value::new(lp.i.clone(), LType::I64), &lp.incr)?;
    let a = cg.fresh_reg();
    cg.emit(format!("{a} = load i64, i64* {acc}"));
    let new = call_with_values(cg, &combine, vec![Value::new(a, LType::I64), elem])?;
    let new = crate::result::unwrap(cg, new);
    let new = box_to_i64(cg, new);
    cg.emit(format!("store i64 {}, i64* {acc}", new.operand));
    close_range_loop(cg, &lp);

    let out = cg.fresh_reg();
    cg.emit(format!("{out} = load i64, i64* {acc}"));
    Ok(Value::new(out, LType::I64))
}

/// Open a counted loop over a runtime list handle `l`, returning
/// `(index-reg, element-i64-reg, len, cond-label, body-label, incr-label, end-label)`
/// with the builder positioned in the body. The caller fills the body, then
/// emits the back-edge via [`close_list_loop`].
struct ListLoop {
    idx: String,
    elem: String,
    incr: String,
    cond: String,
    endl: String,
}

fn open_list_loop(cg: &mut Codegen, l: &Value) -> ListLoop {
    cg.add_extern("declare i64 @osprey_list_length(i8*)");
    cg.add_extern("declare i64 @osprey_list_get(i8*, i64)");
    let len = cg.fresh_reg();
    cg.emit(format!(
        "{len} = call i64 @osprey_list_length(i8* {})",
        l.operand
    ));
    let idx = cg.fresh_reg();
    cg.emit(format!("{idx} = alloca i64"));
    cg.emit(format!("store i64 0, i64* {idx}"));
    let cond = cg.fresh_label();
    let body = cg.fresh_label();
    let incr = cg.fresh_label();
    let endl = cg.fresh_label();
    cg.emit(format!("br label %{cond}"));

    cg.start_block(&cond);
    let i = cg.fresh_reg();
    cg.emit(format!("{i} = load i64, i64* {idx}"));
    let more = cg.fresh_reg();
    cg.emit(format!("{more} = icmp slt i64 {i}, {len}"));
    cg.emit(format!("br i1 {more}, label %{body}, label %{endl}"));

    cg.start_block(&body);
    let elem = cg.fresh_reg();
    cg.emit(format!(
        "{elem} = call i64 @osprey_list_get(i8* {}, i64 {i})",
        l.operand
    ));
    ListLoop {
        idx,
        elem,
        incr,
        cond,
        endl,
    }
}

fn close_list_loop(cg: &mut Codegen, lp: &ListLoop) {
    cg.emit(format!("br label %{}", lp.incr));
    cg.start_block(&lp.incr);
    let i = cg.fresh_reg();
    cg.emit(format!("{i} = load i64, i64* {}", lp.idx));
    let next = cg.fresh_reg();
    cg.emit(format!("{next} = add i64 {i}, 1"));
    cg.emit(format!("store i64 {next}, i64* {}", lp.idx));
    cg.emit(format!("br label %{}", lp.cond));
    cg.start_block(&lp.endl);
}

/// The `i`-th positional argument as a list handle.
fn list_arg(cg: &mut Codegen, args: &[Expr], i: usize) -> Result<Value> {
    let v = gen_expr(cg, nth(args, i)?)?;
    crate::cast::coerce_to(cg, v, LType::Ptr)
}

/// `forEachList(list, fn)` — call `fn` on each element in order.
fn for_each_list(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let l = list_arg(cg, args, 0)?;
    let consumer = callback_name(nth(args, 1)?)?;
    let lp = open_list_loop(cg, &l);
    let _ = call_with_values(cg, &consumer, vec![Value::new(lp.elem.clone(), LType::I64)])?;
    close_list_loop(cg, &lp);
    Ok(l)
}

/// `mapList`/`filterList` — build a new list via the runtime list builder.
fn list_builder(cg: &mut Codegen, args: &[Expr], filter: bool) -> Result<Value> {
    let l = list_arg(cg, args, 0)?;
    let f = callback_name(nth(args, 1)?)?;
    cg.add_extern("declare i8* @osprey_list_builder_new()");
    cg.add_extern("declare void @osprey_list_builder_push(i8*, i64)");
    cg.add_extern("declare i8* @osprey_list_builder_seal(i8*)");
    let bld = cg.fresh_reg();
    cg.emit(format!("{bld} = call i8* @osprey_list_builder_new()"));
    let lp = open_list_loop(cg, &l);
    let elem = Value::new(lp.elem.clone(), LType::I64);
    if filter {
        let pred = call_with_values(cg, &f, vec![elem.clone()])?;
        let pred = crate::result::unwrap(cg, pred);
        let pb = as_i64(cg, pred)?;
        let nz = cg.fresh_reg();
        cg.emit(format!("{nz} = icmp ne i64 {}, 0", pb.operand));
        let push = cg.fresh_label();
        let skip = cg.fresh_label();
        cg.emit(format!("br i1 {nz}, label %{push}, label %{skip}"));
        cg.start_block(&push);
        cg.emit(format!(
            "call void @osprey_list_builder_push(i8* {bld}, i64 {})",
            lp.elem
        ));
        cg.emit(format!("br label %{skip}"));
        cg.start_block(&skip);
    } else {
        let mapped = call_with_values(cg, &f, vec![elem])?;
        let mapped = crate::result::unwrap(cg, mapped);
        let boxed = box_to_i64(cg, mapped);
        cg.emit(format!(
            "call void @osprey_list_builder_push(i8* {bld}, i64 {})",
            boxed.operand
        ));
    }
    close_list_loop(cg, &lp);
    let sealed = cg.fresh_reg();
    cg.emit(format!(
        "{sealed} = call i8* @osprey_list_builder_seal(i8* {bld})"
    ));
    Ok(Value::handle(sealed, crate::collections::LIST_OWNER))
}

/// `foldList(list, initial, fn)` — reduce a list with `fn(acc, elem)`.
fn fold_list(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let l = list_arg(cg, args, 0)?;
    let initial = gen_expr(cg, nth(args, 1)?)?;
    let initial = crate::result::unwrap(cg, initial);
    let initial = box_to_i64(cg, initial);
    let combine = callback_name(nth(args, 2)?)?;
    let acc = cg.fresh_reg();
    cg.emit(format!("{acc} = alloca i64"));
    cg.emit(format!("store i64 {}, i64* {acc}", initial.operand));
    let lp = open_list_loop(cg, &l);
    let a = cg.fresh_reg();
    cg.emit(format!("{a} = load i64, i64* {acc}"));
    let new = call_with_values(
        cg,
        &combine,
        vec![
            Value::new(a, LType::I64),
            Value::new(lp.elem.clone(), LType::I64),
        ],
    )?;
    let new = crate::result::unwrap(cg, new);
    let new = box_to_i64(cg, new);
    cg.emit(format!("store i64 {}, i64* {acc}", new.operand));
    close_list_loop(cg, &lp);
    let out = cg.fresh_reg();
    cg.emit(format!("{out} = load i64, i64* {acc}"));
    Ok(Value::new(out, LType::I64))
}
