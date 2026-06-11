//! Fibers, channels, `yield` and `select`, lowered to the same C fiber runtime
//! every compiled Osprey program links (`fiber_runtime.c` in
//! `libfiber_runtime.a`). `spawn e` lifts `e` into a no-arg `i64 ()` thunk and
//! hands it to `fiber_spawn` — a really-concurrent fiber (the runtime restores
//! the spawner's effect-handler snapshot inside it, so `perform` works there).
//! Locals the expression closes over are spilled through per-spawn module
//! globals — stored at the spawn site, reloaded inside the thunk — preserving
//! the fixed `i64 ()*` ABI for any capture set (captures are by value).
//! `await`/`fiberDone` map to `fiber_await`/`fiber_done`, the non-blocking
//! probe a foreground loop can animate against while the fiber works. Channels
//! are `channel_create`/`channel_send`/`channel_recv`; channel ids and fiber
//! ids draw from the runtime's one shared counter. `yield e` evaluates to its
//! operand and `select` takes its first arm (the deterministic examples drive
//! arm readiness by `send`/`recv` order).

use crate::builder::Codegen;
use crate::conv::{as_i64, box_to_i64};
use crate::error::Result;
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, FieldAssignment, InterpolatedPart, MatchArm, NamedArgument, Stmt};
use std::collections::BTreeSet;

/// One spilled capture: the local's name, its per-spawn global, and the value
/// whose metadata (Result block shape, owner tag) the reload must keep.
struct Capture {
    name: String,
    global: String,
    val: Value,
}

/// `spawn e` — lift `e` into a no-arg thunk and start it on a real fiber.
pub(crate) fn gen_spawn(cg: &mut Codegen, e: &Expr) -> Result<Value> {
    let id = cg.next_lambda_id();
    let thunk = format!("__fiber_closure_{id}");
    let caps = spill_captures(cg, id, e);
    emit_thunk(cg, &thunk, &caps, e)?;
    let r = cg.call("i64", "fiber_spawn", "i64 ()*", &[&format!("@{thunk}")]);
    Ok(Value::new(r, LType::I64))
}

/// Store every local `e` closes over into a fresh module global at the spawn
/// site, so the thunk (a separate function) can reload it.
fn spill_captures(cg: &mut Codegen, id: usize, e: &Expr) -> Vec<Capture> {
    let mut names = BTreeSet::new();
    free_idents(e, &mut names);
    let mut caps = Vec::new();
    for name in names {
        let Some(val) = cg.lookup(&name) else {
            continue;
        };
        let ty = val.llvm_ty();
        let global = format!("@__fiber_cap_{id}_{name}");
        cg.add_global_def(format!("{global} = global {ty} {}", zero_of(&ty)));
        cg.emit(format!("store {ty} {}, {ty}* {global}", val.operand));
        caps.push(Capture { name, global, val });
    }
    caps
}

/// The zero constant for a module global of LLVM type `ty`.
fn zero_of(ty: &str) -> &'static str {
    if ty.ends_with('*') {
        "null"
    } else if ty == "double" {
        "0.0"
    } else {
        "0"
    }
}

/// Emit the thunk: reload each capture, lower the body, and return it as the
/// uniform `i64` fiber result (`fiber_await` hands the same word back).
fn emit_thunk(cg: &mut Codegen, thunk: &str, caps: &[Capture], e: &Expr) -> Result<()> {
    let saved = cg.enter_nested_fn();
    for c in caps {
        let ty = c.val.llvm_ty();
        let r = cg.emit_reg(format!("load {ty}, {ty}* {}", c.global));
        let mut v = c.val.clone();
        v.operand = r;
        cg.bind(c.name.clone(), v);
    }
    let body = thunk_body(cg, e);
    cg.exit_nested_fn(saved, "i64", thunk, &[]);
    body
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

/// `yield e` / `yield` — identity (returns the value passed back).
pub(crate) fn gen_yield(cg: &mut Codegen, e: Option<&Expr>) -> Result<Value> {
    match e {
        Some(inner) => gen_expr(cg, inner),
        None => Ok(Value::unit()),
    }
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
        "fiber_yield" => {
            let v = args.first().map(|a| gen_expr(cg, a)).transpose()?;
            v.unwrap_or_else(Value::unit)
        }
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

// ---- free-identifier collection for spawn captures -------------------------

/// Collect every identifier referenced anywhere in `e`. Over-collection is
/// harmless: only names bound to a *value* at the spawn site spill (a function
/// name resolves through the call path, and a name re-bound inside the
/// expression simply shadows its reload).
fn free_idents(e: &Expr, out: &mut BTreeSet<String>) {
    match e {
        Expr::Integer(_) | Expr::Float(_) | Expr::Str(_) | Expr::Bool(_) => {}
        Expr::Identifier(n) => {
            let _ = out.insert(n.clone());
        }
        Expr::InterpolatedStr(parts) => {
            for p in parts {
                if let InterpolatedPart::Expr(inner) = p {
                    free_idents(inner, out);
                }
            }
        }
        Expr::List(xs) => walk_all(xs, out),
        Expr::Map(entries) => {
            for en in entries {
                free_idents(&en.key, out);
                free_idents(&en.value, out);
            }
        }
        Expr::Object(fields) => walk_fields(fields, out),
        Expr::Binary { left, right, .. } | Expr::Pipe { left, right } => {
            free_idents(left, out);
            free_idents(right, out);
        }
        Expr::Unary { operand, .. } => free_idents(operand, out),
        e2 => free_idents_rest(e2, out),
    }
}

/// Continuation of [`free_idents`] (kept in two halves so each stays small).
fn free_idents_rest(e: &Expr, out: &mut BTreeSet<String>) {
    match e {
        Expr::Call {
            function,
            arguments,
            named_arguments,
        } => {
            free_idents(function, out);
            walk_all(arguments, out);
            walk_named(named_arguments, out);
        }
        Expr::MethodCall {
            target,
            arguments,
            named_arguments,
            ..
        } => {
            free_idents(target, out);
            walk_all(arguments, out);
            walk_named(named_arguments, out);
        }
        Expr::FieldAccess { target, .. } => free_idents(target, out),
        Expr::Index { target, index } => {
            free_idents(target, out);
            free_idents(index, out);
        }
        Expr::Lambda { body, .. } => free_idents(body, out),
        Expr::Match { value, arms } => {
            free_idents(value, out);
            walk_arms(arms, out);
        }
        Expr::Block { statements, value } => {
            walk_stmts(statements, out);
            if let Some(v) = value {
                free_idents(v, out);
            }
        }
        Expr::TypeConstructor { fields, .. } => walk_fields(fields, out),
        Expr::Update { record, fields } => {
            let _ = out.insert(record.clone());
            walk_fields(fields, out);
        }
        e2 => free_idents_fiber(e2, out),
    }
}

/// Final third of the walker: fiber/effect forms (and the leaf-handled rest).
fn free_idents_fiber(e: &Expr, out: &mut BTreeSet<String>) {
    match e {
        Expr::Spawn(inner) | Expr::Await(inner) | Expr::Recv(inner) => free_idents(inner, out),
        Expr::Yield(inner) => {
            if let Some(i) = inner {
                free_idents(i, out);
            }
        }
        Expr::Send { channel, value } => {
            free_idents(channel, out);
            free_idents(value, out);
        }
        Expr::Select { arms } => walk_arms(arms, out),
        Expr::Perform {
            arguments,
            named_arguments,
            ..
        } => {
            walk_all(arguments, out);
            walk_named(named_arguments, out);
        }
        Expr::Handler { arms, body, .. } => {
            for arm in arms {
                free_idents(&arm.body, out);
            }
            free_idents(body, out);
        }
        // Every other variant is fully handled by the first two thirds.
        _ => {}
    }
}

fn walk_all(xs: &[Expr], out: &mut BTreeSet<String>) {
    for x in xs {
        free_idents(x, out);
    }
}

fn walk_named(named: &[NamedArgument], out: &mut BTreeSet<String>) {
    for n in named {
        free_idents(&n.value, out);
    }
}

fn walk_fields(fields: &[FieldAssignment], out: &mut BTreeSet<String>) {
    for f in fields {
        free_idents(&f.value, out);
    }
}

fn walk_arms(arms: &[MatchArm], out: &mut BTreeSet<String>) {
    for arm in arms {
        free_idents(&arm.body, out);
    }
}

fn walk_stmts(statements: &[Stmt], out: &mut BTreeSet<String>) {
    for s in statements {
        match s {
            Stmt::Let { value, .. } | Stmt::Assignment { value, .. } => free_idents(value, out),
            Stmt::Expr(e) => free_idents(e, out),
            _ => {}
        }
    }
}
