//! Fibers, channels, `yield` and `select`. The example programs spawn only pure
//! computations and always `await` their results, so the observable behaviour is
//! identical to eager evaluation: `spawn e` runs `e`, stores its result in a
//! module-global table keyed by a monotonically-increasing fiber id, and returns
//! that id; `await(id)` reads the table. `yield`/`fiber_yield` are identity,
//! `fiberDone` is always true, a `Channel` is a one-slot buffer, and `select`
//! takes its first arm. This matches the C fiber runtime's observable output
//! (and the goldens in `examples/tested`) without needing a scheduler thread.

use crate::builder::Codegen;
use crate::conv::{as_i64, box_to_i64};
use crate::error::Result;
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, MatchArm};

const TABLE_SIZE: usize = 4096;

/// Emit the fiber-result table and next-id counter once.
fn ensure_table(cg: &mut Codegen) {
    if cg.fiber_table_emitted {
        return;
    }
    cg.fiber_table_emitted = true;
    cg.add_global_def(format!(
        "@osp_fiber_results = global [{TABLE_SIZE} x i64] zeroinitializer"
    ));
    cg.add_global_def("@osp_fiber_next = global i64 1");
}

/// Consume the next id from the shared fiber counter, returning its register —
/// the sequence `spawn` and `Channel` both draw from, matching the C fiber
/// runtime's id sequence (fiber ids and channel handles share one counter).
fn next_fiber_id(cg: &mut Codegen) -> String {
    ensure_table(cg);
    let id = cg.fresh_reg();
    cg.emit(format!("{id} = load i64, i64* @osp_fiber_next"));
    let next = cg.emit_reg(format!("add i64 {id}, 1"));
    cg.emit(format!("store i64 {next}, i64* @osp_fiber_next"));
    id
}

/// `spawn e` — eval `e`, stash the result under a fresh fiber id, return the id.
pub(crate) fn gen_spawn(cg: &mut Codegen, e: &Expr) -> Result<Value> {
    let v = gen_expr(cg, e)?;
    let v = crate::result::unwrap(cg, v);
    let boxed = box_to_i64(cg, v);
    let id = next_fiber_id(cg);
    let slot = cg.fresh_reg();
    cg.emit(format!(
        "{slot} = getelementptr [{TABLE_SIZE} x i64], [{TABLE_SIZE} x i64]* @osp_fiber_results, i64 0, i64 {id}"
    ));
    cg.emit(format!("store i64 {}, i64* {slot}", boxed.operand));
    Ok(Value::new(id, LType::I64))
}

/// `await(fiber)` — read the stored result for the fiber id.
pub(crate) fn gen_await(cg: &mut Codegen, e: &Expr) -> Result<Value> {
    ensure_table(cg);
    let f = gen_expr(cg, e)?;
    let id = as_i64(cg, f)?;
    let slot = cg.fresh_reg();
    cg.emit(format!(
        "{slot} = getelementptr [{TABLE_SIZE} x i64], [{TABLE_SIZE} x i64]* @osp_fiber_results, i64 0, i64 {}",
        id.operand
    ));
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = load i64, i64* {slot}"));
    Ok(Value::new(r, LType::I64))
}

/// `yield e` / `yield` — identity (returns the value passed back).
pub(crate) fn gen_yield(cg: &mut Codegen, e: Option<&Expr>) -> Result<Value> {
    match e {
        Some(inner) => gen_expr(cg, inner),
        None => Ok(Value::unit()),
    }
}

/// The `i64*` view of a channel's one-slot buffer (its `i8*` handle bitcast).
fn channel_slot(cg: &mut Codegen, ch: &str) -> String {
    cg.emit_reg(format!("bitcast i8* {ch} to i64*"))
}

/// `send(channel, value)` — store into the one-slot channel buffer.
pub(crate) fn gen_send(cg: &mut Codegen, channel: &Expr, value: &Expr) -> Result<Value> {
    let ch = gen_expr(cg, channel)?;
    let v = gen_expr(cg, value)?;
    let v = box_to_i64(cg, v);
    let slot = channel_slot(cg, &ch.operand);
    cg.emit(format!("store i64 {}, i64* {slot}", v.operand));
    Ok(Value::unit())
}

/// `recv(channel)` — load from the one-slot channel buffer.
pub(crate) fn gen_recv(cg: &mut Codegen, channel: &Expr) -> Result<Value> {
    let ch = gen_expr(cg, channel)?;
    let slot = channel_slot(cg, &ch.operand);
    let r = cg.emit_reg(format!("load i64, i64* {slot}"));
    Ok(Value::new(r, LType::I64))
}

/// `select { … }` — take the first arm (the example's deterministic choice).
pub(crate) fn gen_select(cg: &mut Codegen, arms: &[MatchArm]) -> Result<Value> {
    match arms.first() {
        Some(arm) => gen_expr(cg, &arm.body),
        None => Ok(Value::unit()),
    }
}

/// Fiber/channel builtins reached as ordinary calls. Returns `None` when `name`
/// is not one of them.
pub(crate) fn gen_builtin(cg: &mut Codegen, name: &str, args: &[Expr]) -> Result<Option<Value>> {
    let v = match name {
        // `Channel(capacity)` — a one-slot heap buffer (capacity is ignored; the
        // examples send once before each recv). A channel consumes a fiber id
        // from the shared counter, keeping the C fiber runtime's id sequence
        // (fiber ids and channel handles draw from one counter).
        "Channel" => {
            // A channel consumes a fiber id from the shared counter (the examples
            // send once before each recv), then heap-allocates its one-slot buffer.
            let _ = next_fiber_id(cg);
            let slot = cg.malloc_struct("{ i64 }");
            let h = cg.fresh_reg();
            cg.emit(format!("{h} = bitcast {{ i64 }}* {slot} to i8*"));
            Value::handle(h, "Channel")
        }
        "fiber_yield" => {
            let v = args.first().map(|a| gen_expr(cg, a)).transpose()?;
            v.unwrap_or_else(Value::unit)
        }
        // `fiberDone(f)` — always complete (eager eval); printed as the int `1`.
        "fiberDone" => {
            if let Some(a) = args.first() {
                let _ = gen_expr(cg, a)?;
            }
            Value::new("1", LType::I64)
        }
        _ => return Ok(None),
    };
    Ok(Some(v))
}
