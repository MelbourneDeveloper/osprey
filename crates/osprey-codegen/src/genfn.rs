//! Polymorphism lowering: specialise a generic user function at each call site
//! by inlining its body with the concrete argument types bound to its
//! parameters, and lower a call through a function-typed parameter (`f(x)` where
//! `f: (int) -> int`) to an indirect call. The Go backend emits a name-mangled
//! monomorphic copy per instantiation (`identity_i64_i64`, `applyInt_fn_i64_i64`);
//! inlining + indirect calls reach the same runtime behaviour without mangling.

use crate::builder::Codegen;
use crate::error::Result;
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::{Expr, NamedArgument, Parameter};

/// If `name` is a generic user function, inline its body with the call's
/// arguments bound to its parameters (so its type variables monomorphise to the
/// concrete argument types here) and return the result. A re-entry guard makes a
/// recursive generic call fall back to a direct call rather than inline forever.
pub(crate) fn try_inline(
    cg: &mut Codegen,
    name: &str,
    args: &[Expr],
    named: &[NamedArgument],
) -> Result<Option<Value>> {
    if cg.inlining.contains(name) {
        return Ok(None);
    }
    let Some((params, body)) = cg.fn_defs.get(name).cloned() else {
        return Ok(None);
    };
    // Pair each parameter with its argument expression (named by name, else
    // positional), then bind it as a value — or, when the argument is a bare
    // callee name, as a call alias so the parameter stays callable.
    let saved_aliases = cg.call_aliases.clone();
    cg.push_scope();
    let _ = cg.inlining.insert(name.to_string());
    let result = (|| {
        for (p, a) in pair_args(&params, args, named) {
            if let Some(callee) = alias_target(cg, a) {
                let _ = cg.call_aliases.insert(p.name.clone(), callee);
            } else {
                let v = gen_expr(cg, a)?;
                cg.bind(p.name.clone(), v);
            }
        }
        gen_expr(cg, &body)
    })();
    let _ = cg.inlining.remove(name);
    cg.pop_scope();
    cg.call_aliases = saved_aliases;
    result.map(Some)
}

/// Pair parameters with their argument expressions — named arguments matched by
/// name, otherwise positional.
fn pair_args<'a>(
    params: &'a [Parameter],
    args: &'a [Expr],
    named: &'a [NamedArgument],
) -> Vec<(&'a Parameter, &'a Expr)> {
    if named.is_empty() {
        params.iter().zip(args).collect()
    } else {
        params
            .iter()
            .filter_map(|p| {
                named
                    .iter()
                    .find(|n| n.name == p.name)
                    .map(|n| (p, &n.value))
            })
            .collect()
    }
}

/// When an argument is a bare name that is a callee (a function/builtin) rather
/// than a bound value or a nullary constructor, return that name so the
/// parameter can redirect calls to it.
fn alias_target(cg: &Codegen, arg: &Expr) -> Option<String> {
    match arg {
        Expr::Identifier(n) if cg.lookup(n).is_none() && !cg.is_ctor(n) => Some(n.clone()),
        _ => None,
    }
}

/// If `name` is a function-typed local (a higher-order parameter), lower `f(x)`
/// to an indirect call: bitcast the `i8*` handle back to the function-pointer
/// type and call it with the coerced arguments.
pub(crate) fn try_indirect(
    cg: &mut Codegen,
    name: &str,
    args: &[Expr],
    named: &[NamedArgument],
) -> Result<Option<Value>> {
    let Some((param_tys, ret_ty, ret_inner)) = cg.fn_ptr_locals.get(name).cloned() else {
        return Ok(None);
    };
    let Some(handle) = cg.lookup(name) else {
        return Ok(None);
    };
    let exprs = crate::expr::arg_exprs(args, named);
    let mut typed = Vec::with_capacity(exprs.len());
    for (want, e) in param_tys.iter().zip(exprs) {
        let v = gen_expr(cg, e)?;
        typed.push(crate::cast::coerce_to(cg, v, *want)?.typed());
    }
    let ret_spelling = match ret_inner {
        Some(inner) => format!("{{ {inner}, i8 }}*"),
        None => ret_ty.to_string(),
    };
    let params = param_tys
        .iter()
        .map(LType::to_string)
        .collect::<Vec<_>>()
        .join(", ");
    let fp = cg.emit_reg(format!(
        "bitcast i8* {} to {ret_spelling} ({params})*",
        handle.operand
    ));
    let r = cg.emit_reg(format!("call {ret_spelling} {fp}({})", typed.join(", ")));
    Ok(Some(match ret_inner {
        Some(inner) => Value::result(r, inner),
        None => Value::new(r, ret_ty),
    }))
}
