//! Program/function/statement orchestration — the top-level walk over the
//! module: emit each user function (parameter and return types taken from
//! inference), then synthesize `main` from either a user `main` or the trailing
//! top-level statements.

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, Parameter, Program, Stmt};

/// Compile a whole program to an LLVM IR module (text), driven by the inferred
/// types in [`osprey_types::ProgramTypes`].
///
/// # Errors
///
/// Returns `Err` if any function body, top-level statement, or `main`
/// expression contains a construct that cannot be lowered to LLVM IR.
pub fn compile_program(program: &Program) -> Result<String> {
    let prog = osprey_types::infer_program(program);
    let mut cg = Codegen::with_types(prog);

    // Pre-pass: record parameter names so named-argument calls can be ordered,
    // and parse `effect` operation signatures for `handle`/`perform`.
    for stmt in &program.statements {
        match stmt {
            Stmt::Function {
                name,
                parameters,
                body,
                ..
            } => {
                let _ = cg.fn_params.insert(
                    name.clone(),
                    parameters.iter().map(|p| p.name.clone()).collect(),
                );
                // A polymorphic function is specialised by inlining at each call
                // site, so keep its body reachable.
                if cg.is_generic_fn(name) {
                    let _ = cg
                        .fn_defs
                        .insert(name.clone(), (parameters.clone(), body.clone()));
                }
            }
            Stmt::Effect { name, operations } => {
                for op in operations {
                    cg.register_effect_op(
                        format!("{name}.{}", op.name),
                        crate::effects::parse_op_sig(&op.ty),
                    );
                }
            }
            _ => {}
        }
    }

    let mut top_level: Vec<&Stmt> = Vec::new();
    let mut user_main: Option<&Expr> = None;
    for stmt in &program.statements {
        match stmt {
            Stmt::Function { name, body, .. } if name == "main" => user_main = Some(body),
            // A generic function is specialised by inlining at each call site
            // (recorded in `fn_defs`), so it is not emitted as a monomorphic def.
            Stmt::Function { name, .. } if cg.fn_defs.contains_key(name) => {}
            Stmt::Function {
                name,
                parameters,
                body,
                ..
            } => gen_function(&mut cg, name, parameters, body)?,
            Stmt::Let { .. } | Stmt::Assignment { .. } | Stmt::Expr(_) => top_level.push(stmt),
            _ => {}
        }
    }

    cg.begin_function();
    if let Some(body) = user_main {
        let _ = gen_expr(&mut cg, body)?;
    } else {
        for stmt in &top_level {
            gen_local_stmt(&mut cg, stmt)?;
        }
    }
    cg.emit("ret i32 0");
    cg.finish_function(LType::I32.as_str(), "main", &[]);

    Ok(cg.render())
}

fn gen_function(cg: &mut Codegen, name: &str, parameters: &[Parameter], body: &Expr) -> Result<()> {
    let param_sig = cg
        .fn_param_sig(name)
        .unwrap_or_else(|| vec![(LType::I64, None); parameters.len()]);

    cg.begin_function();
    // Record any function-typed parameters so a call through one lowers to an
    // indirect call (the higher-order `f(x)` in `fn apply(f, x) = f(x)`).
    let fn_ptr_params: Vec<(String, crate::builder::FnSig)> = cg
        .prog
        .param_types(name)
        .map(|ptys| {
            parameters
                .iter()
                .zip(ptys)
                .filter_map(|(p, t)| Codegen::fn_value_sig(t).map(|s| (p.name.clone(), s)))
                .collect()
        })
        .unwrap_or_default();
    for (n, s) in fn_ptr_params {
        let _ = cg.fn_ptr_locals.insert(n, s);
    }
    let mut params = Vec::new();
    for (p, (pty, owner)) in parameters.iter().zip(param_sig.iter()) {
        let v = Value::new(format!("%{}", p.name), *pty).with_owner(owner.clone());
        cg.bind(p.name.clone(), v);
        params.push((*pty, p.name.clone()));
    }
    let body_val = gen_expr(cg, body)?;
    let ret = coerce_return(cg, name, body_val)?;
    cg.emit(format!("ret {} {}", ret.llvm_ty(), ret.operand));
    cg.finish_function(&ret.llvm_ty(), name, &params);
    Ok(())
}

/// Coerce a function body value to its declared return type. A `Result<T, E>`
/// return wraps a bare body into a Success block (or passes an existing Result
/// through); everything else coerces to the inferred scalar return type.
fn coerce_return(cg: &mut Codegen, name: &str, body: Value) -> Result<Value> {
    if let Some(inner) = cg.fn_ret_result_inner(name) {
        if body.result_inner.is_some() {
            return Ok(body);
        }
        return crate::result::make_ok(cg, body, inner);
    }
    let ret_ty = cg.fn_ret_ltype(name).unwrap_or(LType::I64);
    crate::cast::coerce_to(cg, body, ret_ty)
}

pub(crate) fn gen_local_stmt(cg: &mut Codegen, stmt: &Stmt) -> Result<()> {
    match stmt {
        // An immutable `let` keeps a Result wrapper (so `let v = 21 * 2;
        // toString(v)` shows `Success(42)`); a `mut` reassignment auto-unwraps it
        // (the `mut` auto-unwrap rule: the cell holds the success payload).
        Stmt::Let { name, value, .. } => gen_bind(cg, name, value, false),
        Stmt::Assignment { name, value, .. } => gen_bind(cg, name, value, true),
        Stmt::Expr(e) => {
            let _ = gen_expr(cg, e)?;
            Ok(())
        }
        _ => Err(CodegenError::unsupported("statement in block/main")),
    }
}

/// Bind `name` to `value`. A lambda is recorded for inline application at its
/// call sites rather than evaluated (the backend lowers no closures). When
/// `unwrap` is set (a mutable assignment), a Result value is unwrapped to its
/// success payload before binding.
fn gen_bind(cg: &mut Codegen, name: &str, value: &Expr, unwrap: bool) -> Result<()> {
    if let Expr::Lambda {
        parameters, body, ..
    } = value
    {
        let _ = cg
            .lambdas
            .insert(name.to_string(), (parameters.clone(), (**body).clone()));
        return Ok(());
    }
    let v = gen_expr(cg, value)?;
    let v = if unwrap {
        crate::result::unwrap(cg, v)
    } else {
        v
    };
    cg.bind(name.to_string(), v);
    Ok(())
}
