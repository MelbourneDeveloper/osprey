//! Algebraic effects: `effect` declarations, `handle … in …` and `perform`.
//! Each `handle` arm becomes a top-level handler function; entering the
//! `handle` pushes those functions onto the C runtime's handler stack
//! (`__osprey_handler_push`, keyed by effect+operation name) and leaving pops
//! them, so a `perform` in any (even forward-referenced) function resolves the
//! innermost active handler dynamically via `__osprey_handler_lookup` and an
//! indirect call. The example handlers never `resume`, so an arm is an ordinary
//! function returning the operation's result.

use crate::builder::{CellSlot, Codegen, ResumeCodegenContext};
use crate::cast::coerce_to;
use crate::conv::{box_to_i64, unbox_from_i64};
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

fn declare_coro(cg: &mut Codegen) {
    cg.add_extern("declare i8* @__osprey_coro_new(i8*)");
    cg.add_extern("declare void @__osprey_coro_start(i8*, i64 (i8*)*, i8*, i8*)");
    cg.add_extern("declare i64 @__osprey_coro_suspend(i8*, i64, i64*, i64)");
    cg.add_extern("declare i64 @__osprey_coro_resume(i8*, i64)");
    cg.add_extern("declare i64 @__osprey_coro_done(i8*)");
    cg.add_extern("declare i64 @__osprey_coro_op(i8*)");
    cg.add_extern("declare i64 @__osprey_coro_arg(i8*, i64)");
    cg.add_extern("declare i64 @__osprey_coro_result(i8*)");
    cg.add_extern("declare void @__osprey_coro_abort(i8*)");
    cg.add_extern("declare void @__osprey_coro_free(i8*)");
    cg.add_extern("declare i8* @__osprey_handler_snapshot()");
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
        Expr::Resume(Some(value)) => scan_expr(value, muts, captured),
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
    if arms.iter().any(|arm| contains_resume(&arm.body)) {
        return gen_resuming_handler(cg, effect, arms, body);
    }
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

fn contains_resume(e: &Expr) -> bool {
    match e {
        Expr::Resume(_) => true,
        Expr::InterpolatedStr(parts) => parts.iter().any(
            |p| matches!(p, osprey_ast::InterpolatedPart::Expr(inner) if contains_resume(inner)),
        ),
        Expr::List(xs) => xs.iter().any(contains_resume),
        Expr::Map(entries) => entries
            .iter()
            .any(|entry| contains_resume(&entry.key) || contains_resume(&entry.value)),
        Expr::Object(fields)
        | Expr::TypeConstructor { fields, .. }
        | Expr::Update { fields, .. } => fields.iter().any(|f| contains_resume(&f.value)),
        Expr::Binary { left, right, .. } | Expr::Pipe { left, right } => {
            contains_resume(left) || contains_resume(right)
        }
        Expr::Unary { operand, .. } => contains_resume(operand),
        Expr::Call {
            function,
            arguments,
            named_arguments,
        } => {
            contains_resume(function)
                || arguments.iter().any(contains_resume)
                || named_arguments.iter().any(|n| contains_resume(&n.value))
        }
        Expr::MethodCall {
            target,
            arguments,
            named_arguments,
            ..
        } => {
            contains_resume(target)
                || arguments.iter().any(contains_resume)
                || named_arguments.iter().any(|n| contains_resume(&n.value))
        }
        Expr::FieldAccess { target, .. } => contains_resume(target),
        Expr::Index { target, index } => contains_resume(target) || contains_resume(index),
        Expr::Lambda { body, .. } | Expr::Spawn(body) | Expr::Await(body) | Expr::Recv(body) => {
            contains_resume(body)
        }
        Expr::Yield(Some(value)) => contains_resume(value),
        Expr::Send { channel, value } => contains_resume(channel) || contains_resume(value),
        Expr::Match { value, arms } => {
            contains_resume(value) || arms.iter().any(|arm| contains_resume(&arm.body))
        }
        Expr::Block { statements, value } => {
            statements.iter().any(stmt_contains_resume)
                || value.as_deref().is_some_and(contains_resume)
        }
        Expr::Select { arms } => arms.iter().any(|arm| contains_resume(&arm.body)),
        Expr::Perform {
            arguments,
            named_arguments,
            ..
        } => {
            arguments.iter().any(contains_resume)
                || named_arguments.iter().any(|n| contains_resume(&n.value))
        }
        // A nested handler owns its own `resume`; do not mark the outer handler
        // as a resuming region because of it.
        Expr::Handler { body, .. } => contains_resume(body),
        _ => false,
    }
}

fn stmt_contains_resume(stmt: &Stmt) -> bool {
    match stmt {
        Stmt::Let { value, .. } | Stmt::Assignment { value, .. } | Stmt::Expr { value, .. } => {
            contains_resume(value)
        }
        _ => false,
    }
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
    caps_from_names(cg, names)
}

fn capture_list_resuming(cg: &Codegen, arms: &[HandlerArm], body: &Expr) -> Vec<ArmCap> {
    let mut names = BTreeSet::new();
    for arm in arms {
        let mut free = BTreeSet::new();
        free_idents(&arm.body, &mut free);
        names.extend(free.into_iter().filter(|n| !arm.params.contains(n)));
    }
    let mut body_free = BTreeSet::new();
    free_idents(body, &mut body_free);
    names.extend(body_free);
    caps_from_names(cg, names)
}

fn caps_from_names(cg: &Codegen, names: BTreeSet<String>) -> Vec<ArmCap> {
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

#[derive(Clone)]
struct DriveArm {
    op_id: usize,
    operation: String,
    sig: OpSig,
    arm_fn: String,
}

/// `handle` region whose arms contain explicit `resume`: the handled body runs
/// on a body thread and each `perform` suspends into this host-side dispatcher.
fn gen_resuming_handler(
    cg: &mut Codegen,
    effect: &str,
    arms: &[HandlerArm],
    body: &Expr,
) -> Result<Value> {
    declare_stack(cg);
    declare_coro(cg);

    let caps = capture_list_resuming(cg, arms, body);
    let (env, env_ty) = build_env(cg, &caps);
    let id = cg.next_handler_id();
    let body_fn = format!("__resume_body_{effect}_{id}");
    let drive_fn = format!("__resume_drive_{effect}_{id}");

    let answer_ty = emit_resuming_body_fn(cg, &body_fn, body, &caps, &env_ty)?;
    let mut drive_arms = Vec::new();
    for (op_id, arm) in arms.iter().enumerate() {
        let key = format!("{effect}.{}", arm.operation);
        let sig = cg.effect_op(&key).unwrap_or(OpSig {
            params: vec![LType::I64; arm.params.len()],
            ret: LType::I64,
            ret_result_inner: None,
        });
        let suspend_fn = format!("__resume_suspend_{effect}_{}_{id}_{op_id}", arm.operation);
        let arm_fn = format!("__resume_arm_{effect}_{}_{id}_{op_id}", arm.operation);
        emit_suspend_fn(cg, &suspend_fn, op_id, &sig);
        emit_resuming_arm_fn(
            cg,
            arm,
            &ArmFnSpec {
                name: &arm_fn,
                drive_fn: &drive_fn,
                answer_ty,
                sig: &sig,
                caps: &caps,
                env_ty: &env_ty,
            },
        )?;
        drive_arms.push(DriveArm {
            op_id,
            operation: arm.operation.clone(),
            sig,
            arm_fn,
        });
    }
    emit_drive_fn(cg, &drive_fn, &drive_arms);

    let coro = cg.call("i8*", "__osprey_coro_new", "i8*", &[&env]);
    for arm in &drive_arms {
        let suspend_fn = format!(
            "__resume_suspend_{effect}_{}_{id}_{}",
            arm.operation, arm.op_id
        );
        let eff_s = cg.string_constant(effect);
        let op_s = cg.string_constant(&arm.operation);
        let fp = cg.emit_reg(format!(
            "bitcast {} @{suspend_fn} to i8*",
            arm.sig.fn_ptr_ty()
        ));
        let _ = cg.call(
            "i32",
            "__osprey_handler_push",
            "i8*, i8*, i8*, i8*",
            &[&eff_s.operand, &op_s.operand, &fp, &coro],
        );
    }

    let snap = cg.call("i8*", "__osprey_handler_snapshot", "", &[]);
    cg.call_void(
        "__osprey_coro_start",
        "i8*, i64 (i8*)*, i8*, i8*",
        &[&coro, &format!("@{body_fn}"), &env, &snap],
    );
    let boxed = cg.emit_reg(format!("call i64 @{drive_fn}(i8* {env}, i8* {coro})"));

    for _ in arms {
        let _ = cg.call("i32", "__osprey_handler_pop", "", &[]);
    }
    cg.call_void("__osprey_coro_free", "i8*", &[&coro]);
    Ok(unbox_from_i64(cg, &boxed, answer_ty))
}

fn emit_resuming_body_fn(
    cg: &mut Codegen,
    name: &str,
    body: &Expr,
    caps: &[ArmCap],
    env_ty: &str,
) -> Result<LType> {
    let saved = cg.enter_nested_fn();
    reload_env(cg, caps, env_ty);
    let body_raw = gen_expr(cg, body)?;
    let body = crate::result::unwrap(cg, body_raw);
    let answer_ty = body.ty;
    let boxed = box_codegen_value(cg, body);
    cg.emit(format!("ret i64 {}", boxed.operand));
    cg.exit_nested_fn(saved, "i64", name, &[(LType::Ptr, String::from("__env"))]);
    Ok(answer_ty)
}

fn emit_suspend_fn(cg: &mut Codegen, name: &str, op_id: usize, sig: &OpSig) {
    let saved = cg.enter_nested_fn();
    let mut params = vec![(LType::Ptr, String::from("__coro"))];
    for (i, pty) in sig.params.iter().copied().enumerate() {
        params.push((pty, format!("__arg{i}")));
    }

    let args_ptr = if sig.params.is_empty() {
        String::from("null")
    } else {
        let arr_ty = format!("[{} x i64]", sig.params.len());
        let arr = cg.emit_reg(format!("alloca {arr_ty}"));
        for (i, pty) in sig.params.iter().copied().enumerate() {
            let value = Value::new(format!("%__arg{i}"), pty);
            let boxed = box_codegen_value(cg, value);
            let slot = cg.emit_reg(format!(
                "getelementptr {arr_ty}, {arr_ty}* {arr}, i64 0, i64 {i}"
            ));
            cg.emit(format!("store i64 {}, i64* {slot}", boxed.operand));
        }
        cg.emit_reg(format!(
            "getelementptr {arr_ty}, {arr_ty}* {arr}, i64 0, i64 0"
        ))
    };
    let raw = cg.call(
        "i64",
        "__osprey_coro_suspend",
        "i8*, i64, i64*, i64",
        &[
            "%__coro",
            &op_id.to_string(),
            &args_ptr,
            &sig.params.len().to_string(),
        ],
    );
    let ret = unbox_coro_value(cg, &raw, sig.ret, sig.ret_result_inner);
    cg.emit(format!("ret {} {}", ret.llvm_ty(), ret.operand));
    let ret_ty = sig.ret_ty();
    cg.exit_nested_fn(saved, &ret_ty, name, &params);
}

struct ArmFnSpec<'a> {
    name: &'a str,
    drive_fn: &'a str,
    answer_ty: LType,
    sig: &'a OpSig,
    caps: &'a [ArmCap],
    env_ty: &'a str,
}

fn emit_resuming_arm_fn(cg: &mut Codegen, arm: &HandlerArm, spec: &ArmFnSpec<'_>) -> Result<()> {
    let saved = cg.enter_nested_fn();
    reload_env(cg, spec.caps, spec.env_ty);
    let mut params = vec![
        (LType::Ptr, String::from("__env")),
        (LType::Ptr, String::from("__coro")),
    ];
    for (i, pname) in arm.params.iter().enumerate() {
        let pty = spec.sig.params.get(i).copied().unwrap_or(LType::I64);
        cg.bind(pname.clone(), Value::new(format!("%{pname}"), pty));
        params.push((pty, pname.clone()));
    }
    cg.resume_ctx = Some(ResumeCodegenContext {
        env: String::from("%__env"),
        coro: String::from("%__coro"),
        drive_fn: spec.drive_fn.to_string(),
        answer_ty: spec.answer_ty,
    });
    let body_raw = gen_expr(cg, &arm.body)?;
    let body = coerce_to(cg, body_raw, spec.answer_ty)?;
    let boxed = box_codegen_value(cg, body);
    cg.emit(format!("ret i64 {}", boxed.operand));
    cg.exit_nested_fn(saved, "i64", spec.name, &params);
    Ok(())
}

fn emit_drive_fn(cg: &mut Codegen, name: &str, arms: &[DriveArm]) {
    let saved = cg.enter_nested_fn();
    let params = vec![
        (LType::Ptr, String::from("__env")),
        (LType::Ptr, String::from("__coro")),
    ];
    let done = cg.call("i64", "__osprey_coro_done", "i8*", &["%__coro"]);
    let done_cond = cg.emit_reg(format!("icmp ne i64 {done}, 0"));
    let done_lbl = cg.fresh_label();
    let dispatch_lbl = cg.fresh_label();
    cg.emit(format!(
        "br i1 {done_cond}, label %{done_lbl}, label %{dispatch_lbl}"
    ));

    cg.start_block(&done_lbl);
    let result = cg.call("i64", "__osprey_coro_result", "i8*", &["%__coro"]);
    cg.emit(format!("ret i64 {result}"));

    cg.start_block(&dispatch_lbl);
    let op = cg.call("i64", "__osprey_coro_op", "i8*", &["%__coro"]);
    let miss_lbl = cg.fresh_label();
    let check_labels: Vec<String> = arms.iter().map(|_| cg.fresh_label()).collect();
    let arm_labels: Vec<String> = arms.iter().map(|_| cg.fresh_label()).collect();
    if let Some(first) = check_labels.first() {
        cg.emit(format!("br label %{first}"));
    } else {
        cg.emit(format!("br label %{miss_lbl}"));
    }

    for (i, ((arm, check_label), arm_label)) in
        arms.iter().zip(&check_labels).zip(&arm_labels).enumerate()
    {
        cg.start_block(check_label);
        let cmp = cg.emit_reg(format!("icmp eq i64 {op}, {}", arm.op_id));
        let next = check_labels.get(i + 1).unwrap_or(&miss_lbl);
        cg.emit(format!("br i1 {cmp}, label %{arm_label}, label %{next}"));
    }

    for (arm, arm_label) in arms.iter().zip(&arm_labels) {
        cg.start_block(arm_label);
        let mut args = vec![String::from("i8* %__env"), String::from("i8* %__coro")];
        for (idx, pty) in arm.sig.params.iter().copied().enumerate() {
            let raw = cg.call(
                "i64",
                "__osprey_coro_arg",
                "i8*, i64",
                &["%__coro", &idx.to_string()],
            );
            let value = unbox_from_i64(cg, &raw, pty);
            args.push(value.typed());
        }
        let arm_result = cg.emit_reg(format!("call i64 @{}({})", arm.arm_fn, args.join(", ")));
        let done_after = cg.call("i64", "__osprey_coro_done", "i8*", &["%__coro"]);
        let done_after_cond = cg.emit_reg(format!("icmp ne i64 {done_after}, 0"));
        let abort_lbl = cg.fresh_label();
        let return_lbl = cg.fresh_label();
        cg.emit(format!(
            "br i1 {done_after_cond}, label %{return_lbl}, label %{abort_lbl}"
        ));
        cg.start_block(&abort_lbl);
        cg.call_void("__osprey_coro_abort", "i8*", &["%__coro"]);
        cg.emit(format!("br label %{return_lbl}"));
        cg.start_block(&return_lbl);
        cg.emit(format!("ret i64 {arm_result}"));
    }

    cg.start_block(&miss_lbl);
    cg.call_void("__osprey_coro_abort", "i8*", &["%__coro"]);
    cg.emit("ret i64 0");
    cg.exit_nested_fn(saved, "i64", name, &params);
}

pub(crate) fn gen_resume(cg: &mut Codegen, value: Option<&Expr>) -> Result<Value> {
    declare_coro(cg);
    let Some(ctx) = cg.resume_ctx.clone() else {
        return Err(crate::error::CodegenError::invalid(
            "`resume` outside a handler arm",
        ));
    };
    let raw_value = match value {
        Some(expr) => gen_expr(cg, expr)?,
        None => Value::unit(),
    };
    let raw_value = crate::result::unwrap(cg, raw_value);
    let boxed_value = box_codegen_value(cg, raw_value);
    let resumed = cg.call(
        "i64",
        "__osprey_coro_resume",
        "i8*, i64",
        &[&ctx.coro, &boxed_value.operand],
    );
    let done = cg.call("i64", "__osprey_coro_done", "i8*", &[&ctx.coro]);
    let done_cond = cg.emit_reg(format!("icmp ne i64 {done}, 0"));
    let done_lbl = cg.fresh_label();
    let more_lbl = cg.fresh_label();
    let end_lbl = cg.fresh_label();
    cg.emit(format!(
        "br i1 {done_cond}, label %{done_lbl}, label %{more_lbl}"
    ));

    cg.start_block(&done_lbl);
    let done_pred = cg.snapshot_to(&end_lbl);

    cg.start_block(&more_lbl);
    let nested = cg.emit_reg(format!(
        "call i64 @{}(i8* {}, i8* {})",
        ctx.drive_fn, ctx.env, ctx.coro
    ));
    let more_pred = cg.snapshot_to(&end_lbl);

    cg.start_block(&end_lbl);
    let phi = cg.emit_reg(format!(
        "phi i64 [ {resumed}, %{done_pred} ], [ {nested}, %{more_pred} ]"
    ));
    Ok(unbox_from_i64(cg, &phi, ctx.answer_ty))
}

fn box_codegen_value(cg: &mut Codegen, value: Value) -> Value {
    if value.result_inner.is_some() {
        let ptr = cg.emit_reg(format!(
            "bitcast {} {} to i8*",
            value.llvm_ty(),
            value.operand
        ));
        return box_to_i64(cg, Value::new(ptr, LType::Ptr));
    }
    box_to_i64(cg, value)
}

fn unbox_coro_value(cg: &mut Codegen, raw: &str, ty: LType, result_inner: Option<LType>) -> Value {
    if let Some(inner) = result_inner {
        let ptr = cg.emit_reg(format!("inttoptr i64 {raw} to i8*"));
        let struct_ty = crate::llty::result_struct_ty(inner);
        let typed = cg.emit_reg(format!("bitcast i8* {ptr} to {struct_ty}*"));
        return Value::result(typed, inner);
    }
    unbox_from_i64(cg, raw, ty)
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
