//! Program/function/statement orchestration. Ports the top-level walk of
//! `llvm.go` + `program_generation.go`: emit each user function (parameter and
//! return types taken from inference), then synthesize `main` from either a
//! user `main` or the trailing top-level statements.

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::*;

/// Compile a whole program to an LLVM IR module (text), driven by the inferred
/// types in [`osprey_types::ProgramTypes`].
pub fn compile_program(program: &Program) -> Result<String> {
    let prog = osprey_types::infer_program(program);
    let mut cg = Codegen::with_types(prog);

    // Pre-pass: record parameter names so named-argument calls can be ordered,
    // and parse `effect` operation signatures for `handle`/`perform`.
    for stmt in &program.statements {
        match stmt {
            Stmt::Function {
                name, parameters, ..
            } => {
                cg.fn_params.insert(
                    name.clone(),
                    parameters.iter().map(|p| p.name.clone()).collect(),
                );
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
        gen_expr(&mut cg, body)?;
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
        // A let-bound lambda is recorded for inline application at its call
        // sites rather than evaluated to a value (the backend lowers no
        // closures).
        Stmt::Let { name, value, .. } | Stmt::Assignment { name, value, .. } => {
            if let Expr::Lambda { parameters, body, .. } = value {
                cg.lambdas
                    .insert(name.clone(), (parameters.clone(), (**body).clone()));
                return Ok(());
            }
            let v = gen_expr(cg, value)?;
            cg.bind(name.clone(), v);
            Ok(())
        }
        Stmt::Expr(e) => {
            gen_expr(cg, e)?;
            Ok(())
        }
        _ => Err(CodegenError::unsupported("statement in block/main")),
    }
}
