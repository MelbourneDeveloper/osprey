//! Fibers, channels, `yield` and `select`, lowered to the same C fiber runtime
//! every compiled Osprey program links (`fiber_runtime.c` in
//! `libfiber_runtime.a`). `spawn e` lowers `e` as a zero-parameter closure
//! (`crate::closure`): the thunk takes the closure cell as its env, reloads the
//! captures `e` closes over, and is handed to `fiber_spawn_env` together with
//! its per-spawn heap cell — so two in-flight spawns from one site never share
//! capture state (the runtime restores the spawner's effect-handler snapshot
//! inside the fiber, so `perform` works there). `await`/`fiberDone` map to
//! `fiber_await`/`fiber_done`, the non-blocking probe a foreground loop can
//! animate against while the fiber works. Channels are
//! `channel_create`/`channel_send`/`channel_recv`; channel ids and fiber ids
//! draw from the runtime's one shared counter. `yield e` performs the runtime's
//! cooperative hand-off (`fiber_yield`) and evaluates to its operand; `select`
//! takes its first arm (the deterministic examples drive arm readiness by
//! `send`/`recv` order).

use crate::builder::Codegen;
use crate::conv::{as_i64, box_to_i64};
use crate::error::Result;
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, MatchArm};

/// The thunk's value-ABI signature: no parameters, returns the boxed `i64`
/// fiber result.
const THUNK_SIG: (LType, Option<LType>) = (LType::I64, None);

/// `spawn e` — lower `e` as a zero-parameter closure and start it on a real
/// fiber via `fiber_spawn_env(thunk, cell)`.
pub(crate) fn gen_spawn(cg: &mut Codegen, e: &Expr) -> Result<Value> {
    let id = cg.next_lambda_id();
    let thunk = format!("__fiber_thunk_{id}");
    let caps = crate::closure::capture_list(cg, &[], e);
    let cell_ty = crate::closure::cell_struct_ty(&caps);
    let saved = cg.enter_nested_fn();
    crate::closure::reload_captures(cg, &cell_ty, &caps);
    let body = thunk_body(cg, e);
    cg.exit_nested_fn(saved, "i64", &thunk, &[(LType::Ptr, String::from("__env"))]);
    body?;
    let sig = (Vec::new(), THUNK_SIG.0, THUNK_SIG.1);
    let cell = crate::closure::cell_value(cg, id, &thunk, &cell_ty, &caps, &sig);
    let r = cg.call(
        "i64",
        "fiber_spawn_env",
        "i64 (i8*)*, i8*",
        &[&format!("@{thunk}"), &cell.operand],
    );
    Ok(Value::new(r, LType::I64))
}

fn thunk_body(cg: &mut Codegen, e: &Expr) -> Result<()> {
    let v = gen_expr(cg, e)?;
    let v = crate::result::unwrap(cg, v);
    let b = box_to_i64(cg, v);
    cg.emit(format!("ret i64 {}", b.operand));
    Ok(())
}

/// `await(fiber)` — block on the C runtime until the fiber completes.
pub(crate) fn gen_await(cg: &mut Codegen, e: &Expr) -> Result<Value> {
    let f = gen_expr(cg, e)?;
    let id = as_i64(cg, f)?;
    let r = cg.call("i64", "fiber_await", "i64", &[&id.operand]);
    Ok(Value::new(r, LType::I64))
}

/// `yield e` / `yield` — drive the runtime's cooperative hand-off, then evaluate
/// to the operand. `fiber_yield` donates the CPU to the scheduler in concurrent
/// mode (a no-op under deterministic execution) and forwards its argument
/// unchanged, so the Osprey value is preserved with its original type.
pub(crate) fn gen_yield(cg: &mut Codegen, e: Option<&Expr>) -> Result<Value> {
    let value = match e {
        Some(inner) => gen_expr(cg, inner)?,
        None => Value::unit(),
    };
    let boxed = box_to_i64(cg, value.clone());
    let _ = cg.call("i64", "fiber_yield", "i64", &[&boxed.operand]);
    Ok(value)
}

/// `send(channel, value)` — `channel_send` on the C runtime (blocks when full).
pub(crate) fn gen_send(cg: &mut Codegen, channel: &Expr, value: &Expr) -> Result<Value> {
    let ch = gen_expr(cg, channel)?;
    let id = as_i64(cg, ch)?;
    let v = gen_expr(cg, value)?;
    let v = box_to_i64(cg, v);
    let r = cg.call(
        "i64",
        "channel_send",
        "i64, i64",
        &[&id.operand, &v.operand],
    );
    Ok(Value::new(r, LType::I64))
}

/// `recv(channel)` — `channel_recv` on the C runtime (blocks when empty).
pub(crate) fn gen_recv(cg: &mut Codegen, channel: &Expr) -> Result<Value> {
    let ch = gen_expr(cg, channel)?;
    let id = as_i64(cg, ch)?;
    let r = cg.call("i64", "channel_recv", "i64", &[&id.operand]);
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
        // `Channel(capacity)` — a real C-runtime channel; its id comes from the
        // same counter as fiber ids.
        "Channel" => {
            let cap = match args.first() {
                Some(a) => {
                    let v = gen_expr(cg, a)?;
                    as_i64(cg, v)?.operand
                }
                None => String::from("0"),
            };
            let r = cg.call("i64", "channel_create", "i64", &[&cap]);
            Value::new(r, LType::I64)
        }
        // `fiber_yield(v)` called as an ordinary function shares `yield`'s
        // lowering — the same runtime hand-off, forwarding `v`.
        "fiber_yield" => gen_yield(cg, args.first())?,
        // `fiberDone(f)` — the C runtime's non-blocking completion probe.
        "fiberDone" => {
            let Some(a) = args.first() else {
                return Err(crate::error::CodegenError::invalid(
                    "fiberDone needs a fiber argument",
                ));
            };
            let v = gen_expr(cg, a)?;
            let id = as_i64(cg, v)?;
            let r = cg.call("i64", "fiber_done", "i64", &[&id.operand]);
            Value::new(r, LType::I64)
        }
        _ => return Ok(None),
    };
    Ok(Some(v))
}
