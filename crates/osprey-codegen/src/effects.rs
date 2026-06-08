//! Algebraic effects: `effect` declarations, `handle … in …` and `perform`.
//! Ports `effects_generation.go`. Each `handle` arm becomes a top-level handler
//! function; entering the `handle` pushes those functions onto the C runtime's
//! handler stack (`__osprey_handler_push`, keyed by effect+operation name) and
//! leaving pops them, so a `perform` in any (even forward-referenced) function
//! resolves the innermost active handler dynamically via
//! `__osprey_handler_lookup` and an indirect call. The example handlers never
//! `resume`, so an arm is an ordinary function returning the operation's result.

use crate::builder::Codegen;
use crate::cast::coerce_to;
use crate::error::Result;
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use crate::types::ltype_of_name;
use osprey_ast::{Expr, HandlerArm};

/// A parsed effect-operation signature: parameter types, the result LLVM type,
/// and (when the result is `Result<T, _>`) the success inner type.
#[derive(Clone)]
pub(crate) struct OpSig {
    pub params: Vec<LType>,
    pub ret: LType,
    pub ret_result_inner: Option<LType>,
}

impl OpSig {
    /// The handler function's LLVM return-type spelling (`{ T, i8 }*` for a
    /// Result result, else the plain type).
    fn ret_ty(&self) -> String {
        match self.ret_result_inner {
            Some(inner) => format!("{{ {inner}, i8 }}*"),
            None => self.ret.to_string(),
        }
    }

    /// The handler function-pointer type, e.g. `i64 (i8*, i64)*`.
    fn fn_ptr_ty(&self) -> String {
        let params = self
            .params
            .iter()
            .map(LType::to_string)
            .collect::<Vec<_>>()
            .join(", ");
        format!("{} ({params})*", self.ret_ty())
    }
}

/// Parse an effect operation's written type (`fn(T1, T2) -> R`) into an [`OpSig`].
pub(crate) fn parse_op_sig(ty: &str) -> OpSig {
    let inner = ty.trim();
    let open = inner.find('(');
    let close = inner.rfind(')');
    let (params, ret) = match (open, close) {
        (Some(o), Some(c)) if c > o => {
            let params = split_top(&inner[o + 1..c]);
            let ret = inner[c + 1..]
                .trim_start()
                .trim_start_matches("->")
                .trim()
                .to_string();
            (params, ret)
        }
        _ => (Vec::new(), inner.to_string()),
    };
    let param_tys = params.iter().map(|p| ltype_of_name(p)).collect();
    let (ret_lty, result_inner) = parse_ret(&ret);
    OpSig {
        params: param_tys,
        ret: ret_lty,
        ret_result_inner: result_inner,
    }
}

/// Map a return-type spelling to `(carried LLVM type, Result success inner)`.
/// `Unit` carries an `i64 0`; a `Result<T, _>` carries the `{ T, i8 }*` block.
fn parse_ret(ret: &str) -> (LType, Option<LType>) {
    let open = ret.find('<').or_else(|| ret.find('['));
    let head = match open {
        Some(i) => ret[..i].trim(),
        None => ret.trim(),
    };
    if head != "Result" {
        return (ltype_of_name(head), None);
    }
    // Result<T, E> — the success inner type T is the first generic argument.
    let inner = open
        .map(|i| {
            let body = ret[i + 1..]
                .trim_end()
                .trim_end_matches('>')
                .trim_end_matches(']');
            split_top(body)
        })
        .and_then(|args| args.into_iter().next())
        .map_or(LType::I64, |a| ltype_of_name(&a));
    (LType::Ptr, Some(inner))
}

/// Split a comma-separated type list at the top bracket level.
fn split_top(s: &str) -> Vec<String> {
    let mut out = Vec::new();
    let mut depth = 0i32;
    let mut cur = String::new();
    for ch in s.chars() {
        match ch {
            '<' | '[' | '(' => {
                depth += 1;
                cur.push(ch);
            }
            '>' | ']' | ')' => {
                depth -= 1;
                cur.push(ch);
            }
            ',' if depth == 0 => {
                if !cur.trim().is_empty() {
                    out.push(cur.trim().to_string());
                }
                cur.clear();
            }
            _ => cur.push(ch),
        }
    }
    if !cur.trim().is_empty() {
        out.push(cur.trim().to_string());
    }
    out
}

fn declare_stack(cg: &mut Codegen) {
    cg.add_extern("declare i32 @__osprey_handler_push(i8*, i8*, i8*)");
    cg.add_extern("declare i32 @__osprey_handler_pop()");
    cg.add_extern("declare i8* @__osprey_handler_lookup(i8*, i8*)");
}

/// `handle Effect arm… in body` — emit a handler function per arm, push them on
/// the runtime stack for the duration of `body`, then pop.
pub(crate) fn gen_handler(
    cg: &mut Codegen,
    effect: &str,
    arms: &[HandlerArm],
    body: &Expr,
) -> Result<Value> {
    declare_stack(cg);
    let mut pushed: Vec<(String, OpSig)> = Vec::new();
    for arm in arms {
        let key = format!("{effect}.{}", arm.operation);
        let sig = cg.effect_op(&key).unwrap_or(OpSig {
            params: vec![LType::I64; arm.params.len()],
            ret: LType::I64,
            ret_result_inner: None,
        });
        let id = cg.next_handler_id();
        let fn_name = format!("__handler_{effect}_{}_{id}", arm.operation);
        emit_handler_fn(cg, &fn_name, arm, &sig)?;
        let eff_s = cg.string_constant(effect);
        let op_s = cg.string_constant(&arm.operation);
        let fp = cg.fresh_reg();
        cg.emit(format!(
            "{fp} = bitcast {} @{fn_name} to i8*",
            sig.fn_ptr_ty()
        ));
        let r = cg.fresh_reg();
        cg.emit(format!(
            "{r} = call i32 @__osprey_handler_push(i8* {}, i8* {}, i8* {fp})",
            eff_s.operand, op_s.operand
        ));
        pushed.push((arm.operation.clone(), sig));
    }

    let result = gen_expr(cg, body)?;

    for _ in &pushed {
        let r = cg.fresh_reg();
        cg.emit(format!("{r} = call i32 @__osprey_handler_pop()"));
    }
    Ok(result)
}

/// Emit a top-level handler function for one arm: its parameters are the
/// operation's, its body the arm body coerced to the operation's result.
fn emit_handler_fn(cg: &mut Codegen, name: &str, arm: &HandlerArm, sig: &OpSig) -> Result<()> {
    let saved = cg.enter_nested_fn();
    let mut params = Vec::new();
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
    let fp = cg.fresh_reg();
    cg.emit(format!("{fp} = bitcast i8* {raw} to {}", sig.fn_ptr_ty()));
    let ret_ty = sig.ret_ty();
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call {ret_ty} {fp}({})", typed.join(", ")));
    Ok(match sig.ret_result_inner {
        Some(inner) => Value::result(r, inner),
        None => Value::new(r, sig.ret),
    })
}
