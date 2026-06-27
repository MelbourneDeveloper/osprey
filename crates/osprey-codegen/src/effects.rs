//! Algebraic effects: `effect` declarations, `handle … in …` and `perform`.
//! Each `handle` arm becomes a top-level handler function; entering the
//! `handle` pushes those functions onto the C runtime's handler stack
//! (`__osprey_handler_push`, keyed by effect+operation name) and leaving pops
//! them, so a `perform` in any (even forward-referenced) function resolves the
//! innermost active handler dynamically via `__osprey_handler_lookup` and an
//! indirect call. The example handlers never `resume`, so an arm is an ordinary
//! function returning the operation's result.

use crate::builder::{CellSlot, Codegen};
use crate::cast::coerce_to;
use crate::error::Result;
use crate::expr::gen_expr;
use crate::freevars::free_idents;
use crate::llty::{LType, Value};
use crate::types::{ltype_of, result_inner};
use osprey_ast::{Expr, HandlerArm, MatchArm, NamedArgument, Stmt};
use std::collections::{BTreeSet, HashSet};

/// A parsed effect-operation signature: parameter types, the result LLVM type,
/// and (when the result is `Result<T, _>`) the success inner type.
#[derive(Clone)]
pub(crate) struct OpSig {
    pub params: Vec<LType>,
    pub ret: LType,
    pub ret_result_inner: Option<LType>,
}

impl OpSig {
    /// The handler function's LLVM return-type spelling (the Result block
    /// pointer for a Result result, else the plain type).
    fn ret_ty(&self) -> String {
        crate::llty::ret_spelling(self.ret, self.ret_result_inner)
    }

    /// The handler function-pointer type. Every arm takes a hidden leading
    /// `i8* env` (its captured cells + values), e.g. `i64 (i8*, i64)*`.
    fn fn_ptr_ty(&self) -> String {
        let mut parts = vec!["i8*".to_string()];
        parts.extend(self.params.iter().map(LType::to_string));
        format!("{} ({})*", self.ret_ty(), parts.join(", "))
    }
}

/// One binding shared by every arm of a single `handle` region, captured into
/// the region's environment.
enum ArmCap {
    /// A handler-captured mutable: the env carries the heap cell's `i8*` pointer
    /// so arms `load`/`store` the same slot — handler-owned state. `ptr` is the
    /// cell pointer in the enclosing scope (a `{pointee}*` operand).
    Cell {
        name: String,
        ptr: String,
        pointee: LType,
        osp_ty: Option<String>,
    },
    /// Any other free variable: captured by value, closure-style.
    Val { name: String, val: Value },
}

impl ArmCap {
    /// The env-slot LLVM type: a cell travels as its `i8*` pointer, a value as
    /// its own travelling type.
    fn slot_ty(&self) -> String {
        match self {
            ArmCap::Cell { .. } => "i8*".to_string(),
            ArmCap::Val { val, .. } => val.llvm_ty(),
        }
    }
}

/// Build an [`OpSig`] from inference's resolved operation signature — the one
/// source of truth for effect types (no string re-parsing in the backend).
pub(crate) fn op_sig_of(op: &osprey_types::OpType) -> OpSig {
    let inner = result_inner(&op.ret);
    let ret = if inner.is_some() {
        LType::Ptr
    } else {
        ltype_of(&op.ret)
    };
    OpSig {
        params: op.params.iter().map(ltype_of).collect(),
        ret,
        ret_result_inner: inner,
    }
}

fn declare_stack(cg: &mut Codegen) {
    cg.add_extern("declare i32 @__osprey_handler_push(i8*, i8*, i8*, i8*)");
    cg.add_extern("declare i32 @__osprey_handler_pop()");
    cg.add_extern("declare i8* @__osprey_handler_lookup(i8*, i8*)");
    cg.add_extern("declare i8* @__osprey_handler_lookup_env(i8*, i8*)");
}

/// Mutable locals that an effect handler arm captures from an enclosing scope —
/// the set promoted to shared heap cells so a plain `mut` becomes a reference
/// cell the handler owns (`get`/`set` arms and the outer scope share one slot).
/// `cell_vars = {mutable bindings} ∩ {names a handler arm references freely}`.
/// Implements [EFFECTS-HANDLER-STATE].
pub(crate) fn captured_mut_vars(body: &Expr) -> HashSet<String> {
    let (mut muts, mut captured) = (BTreeSet::new(), BTreeSet::new());
    scan_expr(body, &mut muts, &mut captured);
    muts.intersection(&captured).cloned().collect()
}

/// As [`captured_mut_vars`] but over the trailing top-level statements that
/// synthesize `main` when there is no user `main`.
pub(crate) fn captured_mut_vars_in_stmts(stmts: &[&Stmt]) -> HashSet<String> {
    let (mut muts, mut captured) = (BTreeSet::new(), BTreeSet::new());
    for s in stmts {
        scan_stmt(s, &mut muts, &mut captured);
    }
    muts.intersection(&captured).cloned().collect()
}

// A purpose-built AST walk (parallel in *shape* to `freevars::walk`, but
// collecting two different sets: every mutable binding and the names handler
// arms reference freely). The handler-arm free idents themselves come from
// `free_idents` so that one definition of "what does this close over" stays in
// `freevars`; only the find-the-handlers/muts traversal lives here.
fn scan_expr(e: &Expr, muts: &mut BTreeSet<String>, captured: &mut BTreeSet<String>) {
    match e {
        Expr::Handler { arms, body, .. } => {
            for arm in arms {
                let mut free = BTreeSet::new();
                free_idents(&arm.body, &mut free);
                captured.extend(free.into_iter().filter(|n| !arm.params.contains(n)));
                scan_expr(&arm.body, muts, captured);
            }
            scan_expr(body, muts, captured);
        }
        Expr::Block { statements, value } => {
            for s in statements {
                scan_stmt(s, muts, captured);
            }
            if let Some(v) = value {
                scan_expr(v, muts, captured);
            }
        }
        Expr::Match { value, arms } => {
            scan_expr(value, muts, captured);
            scan_arms(arms, muts, captured);
        }
        Expr::Select { arms } => scan_arms(arms, muts, captured),
        _ => scan_children(e, muts, captured),
    }
}

fn scan_arms(arms: &[MatchArm], muts: &mut BTreeSet<String>, captured: &mut BTreeSet<String>) {
    for arm in arms {
        scan_expr(&arm.body, muts, captured);
    }
}

/// Recurse into every child expression of `e` (the variants that are not
/// special-cased in [`scan_expr`]), so a handler/mut nested anywhere is found.
fn scan_children(e: &Expr, muts: &mut BTreeSet<String>, captured: &mut BTreeSet<String>) {
    match e {
        Expr::InterpolatedStr(parts) => {
            for p in parts {
                if let osprey_ast::InterpolatedPart::Expr(x) = p {
                    scan_expr(x, muts, captured);
                }
            }
        }
        Expr::List(xs) => scan_all(xs, muts, captured),
        Expr::Map(es) => {
            for en in es {
                scan_expr(&en.key, muts, captured);
                scan_expr(&en.value, muts, captured);
            }
        }
        Expr::Object(fs)
        | Expr::TypeConstructor { fields: fs, .. }
        | Expr::Update { fields: fs, .. } => {
            for f in fs {
                scan_expr(&f.value, muts, captured);
            }
        }
        Expr::Binary { left, right, .. } | Expr::Pipe { left, right } => {
            scan_expr(left, muts, captured);
            scan_expr(right, muts, captured);
        }
        Expr::Unary { operand, .. } => scan_expr(operand, muts, captured),
        Expr::Call {
            function,
            arguments,
            named_arguments,
        } => {
            scan_expr(function, muts, captured);
            scan_all(arguments, muts, captured);
            scan_named(named_arguments, muts, captured);
        }
        Expr::MethodCall {
            target,
            arguments,
            named_arguments,
            ..
        } => {
            scan_expr(target, muts, captured);
            scan_all(arguments, muts, captured);
            scan_named(named_arguments, muts, captured);
        }
        Expr::FieldAccess { target, .. } => scan_expr(target, muts, captured),
        Expr::Index { target, index } => {
            scan_expr(target, muts, captured);
            scan_expr(index, muts, captured);
        }
        Expr::Lambda { body, .. } | Expr::Spawn(body) | Expr::Await(body) | Expr::Recv(body) => {
            scan_expr(body, muts, captured);
        }
        Expr::Yield(Some(x)) => scan_expr(x, muts, captured),
        Expr::Send { channel, value } => {
            scan_expr(channel, muts, captured);
            scan_expr(value, muts, captured);
        }
        Expr::Perform {
            arguments,
            named_arguments,
            ..
        } => {
            scan_all(arguments, muts, captured);
            scan_named(named_arguments, muts, captured);
        }
        _ => {}
    }
}

fn scan_all(xs: &[Expr], muts: &mut BTreeSet<String>, captured: &mut BTreeSet<String>) {
    for x in xs {
        scan_expr(x, muts, captured);
    }
}

fn scan_named(
    named: &[NamedArgument],
    muts: &mut BTreeSet<String>,
    captured: &mut BTreeSet<String>,
) {
    for n in named {
        scan_expr(&n.value, muts, captured);
    }
}

fn scan_stmt(s: &Stmt, muts: &mut BTreeSet<String>, captured: &mut BTreeSet<String>) {
    match s {
        Stmt::Let {
            name,
            value,
            mutable,
            ..
        } => {
            if *mutable {
                let _ = muts.insert(name.clone());
            }
            scan_expr(value, muts, captured);
        }
        Stmt::Assignment { name, value, .. } => {
            let _ = muts.insert(name.clone());
            scan_expr(value, muts, captured);
        }
        Stmt::Expr { value, .. } => scan_expr(value, muts, captured),
        _ => {}
    }
}

/// `handle Effect arm… in body` — capture the region's environment (the cells
/// and values its arms reference), emit a handler function per arm bound to that
/// env, push them on the runtime stack for the duration of `body`, then pop.
pub(crate) fn gen_handler(
    cg: &mut Codegen,
    effect: &str,
    arms: &[HandlerArm],
    body: &Expr,
) -> Result<Value> {
    declare_stack(cg);
    let caps = capture_list(cg, arms);
    let (env, env_ty) = build_env(cg, &caps);
    for arm in arms {
        let key = format!("{effect}.{}", arm.operation);
        let sig = cg.effect_op(&key).unwrap_or(OpSig {
            params: vec![LType::I64; arm.params.len()],
            ret: LType::I64,
            ret_result_inner: None,
        });
        let id = cg.next_handler_id();
        let fn_name = format!("__handler_{effect}_{}_{id}", arm.operation);
        emit_handler_fn(cg, &fn_name, arm, &sig, &caps, &env_ty)?;
        let eff_s = cg.string_constant(effect);
        let op_s = cg.string_constant(&arm.operation);
        let fp = cg.fresh_reg();
        cg.emit(format!(
            "{fp} = bitcast {} @{fn_name} to i8*",
            sig.fn_ptr_ty()
        ));
        let r = cg.fresh_reg();
        cg.emit(format!(
            "{r} = call i32 @__osprey_handler_push(i8* {}, i8* {}, i8* {fp}, i8* {env})",
            eff_s.operand, op_s.operand
        ));
    }

    let result = gen_expr(cg, body)?;

    for _ in arms {
        let r = cg.fresh_reg();
        cg.emit(format!("{r} = call i32 @__osprey_handler_pop()"));
    }
    Ok(result)
}

/// The bindings every arm of this region captures, in stable (sorted) order: a
/// handler-captured mutable becomes a shared [`ArmCap::Cell`]; any other bound
/// free variable is captured by value. Names that resolve to nothing in scope
/// (top-level functions, constructors) need no capture — the arm resolves them
/// directly.
fn capture_list(cg: &Codegen, arms: &[HandlerArm]) -> Vec<ArmCap> {
    let mut names = BTreeSet::new();
    for arm in arms {
        let mut free = BTreeSet::new();
        free_idents(&arm.body, &mut free);
        names.extend(free.into_iter().filter(|n| !arm.params.contains(n)));
    }
    names
        .into_iter()
        .filter_map(|name| {
            if let Some(slot) = cg.cell_slots.get(&name) {
                Some(ArmCap::Cell {
                    name,
                    ptr: slot.ptr.clone(),
                    pointee: slot.pointee,
                    osp_ty: slot.osp_ty.clone(),
                })
            } else {
                cg.lookup(&name).map(|val| ArmCap::Val { name, val })
            }
        })
        .collect()
}

/// Allocate the region's environment cell and store each capture into it,
/// returning its `i8*` handle and the struct type. A capture-free region uses a
/// `null` env (the arms ignore it).
fn build_env(cg: &mut Codegen, caps: &[ArmCap]) -> (String, String) {
    if caps.is_empty() {
        return ("null".to_string(), String::new());
    }
    let env_ty = format!(
        "{{ {} }}",
        caps.iter()
            .map(ArmCap::slot_ty)
            .collect::<Vec<_>>()
            .join(", ")
    );
    let cell = cg.malloc_struct(&env_ty);
    for (i, c) in caps.iter().enumerate() {
        let slot_ty = c.slot_ty();
        let p = cg.emit_reg(format!(
            "getelementptr {env_ty}, {env_ty}* {cell}, i32 0, i32 {i}"
        ));
        let operand = store_operand(cg, c);
        cg.emit(format!("store {slot_ty} {operand}, {slot_ty}* {p}"));
    }
    let env = cg.emit_reg(format!("bitcast {env_ty}* {cell} to i8*"));
    (env, env_ty)
}

/// The operand stored into the env slot for a capture: a cell's `i8*` pointer
/// (the heap slot, shared so arms mutate the same location), or a value's
/// operand.
fn store_operand(cg: &mut Codegen, c: &ArmCap) -> String {
    match c {
        ArmCap::Cell { ptr, pointee, .. } => {
            let ty = pointee.as_str();
            cg.emit_reg(format!("bitcast {ty}* {ptr} to i8*"))
        }
        ArmCap::Val { val, .. } => val.operand.clone(),
    }
}

/// Emit a top-level handler function for one arm: a hidden leading `i8* %__env`
/// it reloads its captures from, then the operation's own parameters; its body
/// is the arm body coerced to the operation's result.
fn emit_handler_fn(
    cg: &mut Codegen,
    name: &str,
    arm: &HandlerArm,
    sig: &OpSig,
    caps: &[ArmCap],
    env_ty: &str,
) -> Result<()> {
    let saved = cg.enter_nested_fn();
    let mut params = vec![(LType::Ptr, String::from("__env"))];
    reload_env(cg, caps, env_ty);
    for (i, pname) in arm.params.iter().enumerate() {
        let pty = sig.params.get(i).copied().unwrap_or(LType::I64);
        cg.bind(pname.clone(), Value::new(format!("%{pname}"), pty));
        params.push((pty, pname.clone()));
    }
    let body = gen_expr(cg, &arm.body)?;
    let ret = if let Some(inner) = sig.ret_result_inner {
        if body.result_inner.is_some() {
            body
        } else {
            crate::result::make_ok(cg, body, inner)?
        }
    } else {
        coerce_to(cg, body, sig.ret)?
    };
    cg.emit(format!("ret {} {}", ret.llvm_ty(), ret.operand));
    let ret_ty = sig.ret_ty();
    cg.exit_nested_fn(saved, &ret_ty, name, &params);
    Ok(())
}

/// Inside an arm function: cast `%__env` back to the region's struct and rebuild
/// each capture — a [`ArmCap::Cell`] as a live cell slot (so reads `load` and
/// reassignments `store` the shared heap location), a value by binding its
/// reloaded register.
fn reload_env(cg: &mut Codegen, caps: &[ArmCap], env_ty: &str) {
    if caps.is_empty() {
        return;
    }
    let env = cg.emit_reg(format!("bitcast i8* %__env to {env_ty}*"));
    for (i, c) in caps.iter().enumerate() {
        let slot_ty = c.slot_ty();
        let p = cg.emit_reg(format!(
            "getelementptr {env_ty}, {env_ty}* {env}, i32 0, i32 {i}"
        ));
        let loaded = cg.emit_reg(format!("load {slot_ty}, {slot_ty}* {p}"));
        match c {
            ArmCap::Cell {
                name,
                pointee,
                osp_ty,
                ..
            } => {
                let ptr = cg.emit_reg(format!("bitcast i8* {loaded} to {}*", pointee.as_str()));
                let _ = cg.cell_slots.insert(
                    name.clone(),
                    CellSlot {
                        ptr,
                        pointee: *pointee,
                        osp_ty: osp_ty.clone(),
                    },
                );
            }
            ArmCap::Val { name, val } => {
                let mut v = val.clone();
                v.operand = loaded;
                cg.bind(name.clone(), v);
            }
        }
    }
}

/// `perform Effect.op(args)` — look up the active handler and call it.
pub(crate) fn gen_perform(
    cg: &mut Codegen,
    effect: &str,
    operation: &str,
    args: &[Expr],
) -> Result<Value> {
    declare_stack(cg);
    let key = format!("{effect}.{operation}");
    let sig = cg.effect_op(&key).unwrap_or(OpSig {
        params: vec![LType::I64; args.len()],
        ret: LType::I64,
        ret_result_inner: None,
    });

    // Evaluate + coerce arguments to the operation's parameter types.
    let mut typed = Vec::new();
    for (i, a) in args.iter().enumerate() {
        let v = gen_expr(cg, a)?;
        let want = sig.params.get(i).copied().unwrap_or(LType::I64);
        let v = coerce_to(cg, v, want)?;
        typed.push(v.typed());
    }

    let eff_s = cg.string_constant(effect);
    let op_s = cg.string_constant(operation);
    let raw = cg.fresh_reg();
    cg.emit(format!(
        "{raw} = call i8* @__osprey_handler_lookup(i8* {}, i8* {})",
        eff_s.operand, op_s.operand
    ));
    let env = cg.fresh_reg();
    cg.emit(format!(
        "{env} = call i8* @__osprey_handler_lookup_env(i8* {}, i8* {})",
        eff_s.operand, op_s.operand
    ));
    let fp = cg.fresh_reg();
    cg.emit(format!("{fp} = bitcast i8* {raw} to {}", sig.fn_ptr_ty()));
    let ret_ty = sig.ret_ty();
    let r = cg.fresh_reg();
    let mut call_args = vec![format!("i8* {env}")];
    call_args.extend(typed);
    cg.emit(format!(
        "{r} = call {ret_ty} {fp}({})",
        call_args.join(", ")
    ));
    Ok(match sig.ret_result_inner {
        Some(inner) => Value::result(r, inner),
        None => Value::new(r, sig.ret),
    })
}
