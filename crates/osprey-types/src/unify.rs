//! Unification — the heart of the inferencer. Ports `TypeInferer.Unify` and its
//! category-specific helpers from `type_inference.go`, including the
//! Osprey-specific rules:
//!   * `any` unifies with anything;
//!   * a bare collection name (`List`, `Map`, `Fiber`, `Channel`) unifies with
//!     its parameterized form (the Go "wildcard generic" rule);
//!   * structural record unification by field name+type;
//!   * Result auto-unwrap at assignment sites (spec 0004), via
//!     [`unify_assignable`].

use crate::ctx::InferCtx;
use crate::error::TypeError;
use crate::ty::{names, Type, VarId};

/// Unify two types, recording the solution in `ctx`. Errors are structural; a
/// failing call may have applied partial bindings, so callers that want to
/// "try" a unification should pre-check shapes rather than relying on rollback.
pub fn unify(ctx: &mut InferCtx, a: &Type, b: &Type) -> Result<(), TypeError> {
    let a = ctx.prune(a);
    let b = ctx.prune(b);
    match (&a, &b) {
        (Type::Var(x), Type::Var(y)) if x == y => Ok(()),
        (Type::Var(x), _) => bind_var(ctx, *x, &b),
        (_, Type::Var(y)) => bind_var(ctx, *y, &a),

        // `any` is the universal wildcard.
        _ if a.is_named(names::ANY) || b.is_named(names::ANY) => Ok(()),

        (Type::Con { name: n1, args: a1 }, Type::Con { name: n2, args: a2 }) => {
            unify_con(ctx, n1, a1, n2, a2, &a, &b)
        }

        (
            Type::Fun {
                params: p1,
                ret: r1,
            },
            Type::Fun {
                params: p2,
                ret: r2,
            },
        ) => unify_fun(ctx, p1, r1, p2, r2),

        (Type::Record { fields: f1, .. }, Type::Record { fields: f2, .. }) => {
            unify_record(ctx, f1, f2, &a, &b)
        }

        // A nominal nullary constructor and a record of the same name describe
        // the same type — a record-typed annotation (`Con "Point"`) meeting a
        // constructed record value (`Record "Point"{..}`).
        (Type::Con { name: n, args }, Type::Record { name: rn, .. })
        | (Type::Record { name: rn, .. }, Type::Con { name: n, args })
            if args.is_empty() && n == rn =>
        {
            Ok(())
        }

        (
            Type::Union {
                name: n1,
                variants: v1,
            },
            Type::Union {
                name: n2,
                variants: v2,
            },
        ) => {
            if n1 != n2 || v1.len() != v2.len() {
                return Err(TypeError::mismatch(&a, &b));
            }
            unify_seq(ctx, v1, v2, &a, &b)
        }

        _ => Err(TypeError::mismatch(&a, &b)),
    }
}

/// Osprey's Result auto-unwrap rule: a `Result<T, E>` value is assignable where
/// a concrete (non-Result) `T` is expected. Used at call arguments, return
/// positions and annotated lets, never inside plain `unify`.
pub fn unify_assignable(
    ctx: &mut InferCtx,
    expected: &Type,
    actual: &Type,
) -> Result<(), TypeError> {
    let expected = ctx.prune(expected);
    let actual = ctx.prune(actual);
    // Unwrap: a `Result<T, E>` value satisfies a concrete `T`.
    if !matches!(expected, Type::Var(_)) && !expected.is_named(names::RESULT) {
        if let Type::Con { name, args } = &actual {
            if name == names::RESULT {
                if let Some(inner) = args.first() {
                    return unify(ctx, &expected, inner);
                }
            }
        }
    }
    // Wrap: a bare `T` value satisfies a `Result<T, E>` return (implicit
    // `Success`), e.g. `fn f() -> Result<bool, E> = x > 0`.
    if let Type::Con { name, args } = &expected {
        if name == names::RESULT
            && !matches!(actual, Type::Var(_))
            && !actual.is_named(names::RESULT)
        {
            if let Some(inner) = args.first() {
                return unify(ctx, inner, &actual);
            }
        }
    }
    unify(ctx, &expected, &actual)
}

fn bind_var(ctx: &mut InferCtx, id: VarId, t: &Type) -> Result<(), TypeError> {
    if let Type::Var(v) = t {
        if *v == id {
            return Ok(());
        }
    }
    if ctx.occurs(id, t) {
        return Err(TypeError::recursive(&Type::Var(id), t));
    }
    ctx.bind(id, t.clone());
    Ok(())
}

fn unify_con(
    ctx: &mut InferCtx,
    n1: &str,
    a1: &[Type],
    n2: &str,
    a2: &[Type],
    a: &Type,
    b: &Type,
) -> Result<(), TypeError> {
    if n1 == n2 && a1.len() == a2.len() {
        return unify_seq(ctx, a1, a2, a, b);
    }
    // A bare constructor name unifies with its applied form (`Fiber` ~
    // `Fiber<int>`, `Box` ~ `Box<int>`) — the Go "wildcard generic" rule,
    // generalized to every nominal type so a bare-named annotation accepts a
    // parameterized value.
    if n1 == n2 && (a1.is_empty() || a2.is_empty()) {
        return Ok(());
    }
    Err(TypeError::mismatch(a, b))
}

fn unify_fun(
    ctx: &mut InferCtx,
    p1: &[Type],
    r1: &Type,
    p2: &[Type],
    r2: &Type,
) -> Result<(), TypeError> {
    if p1.len() != p2.len() {
        return Err(TypeError::new(format!(
            "function arity mismatch: {} vs {} parameters",
            p1.len(),
            p2.len()
        )));
    }
    for (x, y) in p1.iter().zip(p2) {
        unify(ctx, x, y)?;
    }
    unify_fn_return(ctx, r1, r2)
}

/// Function return positions inherit the Result auto-unwrap rule symmetrically:
/// a lambda whose body is `Result<int, E>` satisfies a `(..) -> int` slot, and
/// vice-versa (Go `unifyFunctionTypes`, lines 1541-1551).
fn unify_fn_return(ctx: &mut InferCtx, r1: &Type, r2: &Type) -> Result<(), TypeError> {
    let p1 = ctx.prune(r1);
    let p2 = ctx.prune(r2);
    if !matches!(p1, Type::Var(_)) && !p1.is_named(names::RESULT) {
        if let Type::Con { name, args } = &p2 {
            if name == names::RESULT {
                if let Some(inner) = args.first() {
                    return unify(ctx, &p1, inner);
                }
            }
        }
    }
    if !matches!(p2, Type::Var(_)) && !p2.is_named(names::RESULT) {
        if let Type::Con { name, args } = &p1 {
            if name == names::RESULT {
                if let Some(inner) = args.first() {
                    return unify(ctx, inner, &p2);
                }
            }
        }
    }
    unify(ctx, &p1, &p2)
}

fn unify_record(
    ctx: &mut InferCtx,
    f1: &std::collections::BTreeMap<String, Type>,
    f2: &std::collections::BTreeMap<String, Type>,
    a: &Type,
    b: &Type,
) -> Result<(), TypeError> {
    if f1.len() != f2.len() {
        return Err(TypeError::mismatch(a, b));
    }
    for (name, t1) in f1 {
        match f2.get(name) {
            Some(t2) => unify(ctx, t1, t2)?,
            None => {
                return Err(TypeError::new(format!(
                    "record field mismatch: {a} has no field `{name}` to match {b}"
                )))
            }
        }
    }
    Ok(())
}

fn unify_seq(
    ctx: &mut InferCtx,
    xs: &[Type],
    ys: &[Type],
    a: &Type,
    b: &Type,
) -> Result<(), TypeError> {
    if xs.len() != ys.len() {
        return Err(TypeError::mismatch(a, b));
    }
    for (x, y) in xs.iter().zip(ys) {
        unify(ctx, x, y)?;
    }
    Ok(())
}

#[cfg(test)]
#[expect(
    unused_results,
    reason = "tests drive unification for its side effects and discard the returned types"
)]
mod tests {
    use super::*;

    #[test]
    fn unifies_var_with_concrete() {
        let mut c = InferCtx::new();
        let v = c.fresh();
        unify(&mut c, &v, &Type::int()).unwrap();
        assert_eq!(c.apply(&v), Type::int());
    }

    #[test]
    fn rejects_distinct_primitives() {
        let mut c = InferCtx::new();
        assert!(unify(&mut c, &Type::int(), &Type::string()).is_err());
    }

    #[test]
    fn any_unifies_with_anything() {
        let mut c = InferCtx::new();
        unify(&mut c, &Type::any(), &Type::int()).unwrap();
        unify(&mut c, &Type::list(Type::string()), &Type::any()).unwrap();
    }

    #[test]
    fn bare_collection_unifies_with_parameterized() {
        let mut c = InferCtx::new();
        unify(
            &mut c,
            &Type::prim("Fiber"),
            &Type::con("Fiber", vec![Type::int()]),
        )
        .unwrap();
    }

    #[test]
    fn result_auto_unwraps_at_assignment() {
        let mut c = InferCtx::new();
        let r = Type::result(Type::int(), Type::prim("MathError"));
        unify_assignable(&mut c, &Type::int(), &r).unwrap();
        // But a bare unify keeps them distinct.
        assert!(unify(&mut c, &Type::int(), &r).is_err());
    }

    #[test]
    fn structural_records_ignore_field_order() {
        use std::collections::BTreeMap;
        let mut c = InferCtx::new();
        let mut f1 = BTreeMap::new();
        f1.insert("x".to_string(), Type::int());
        f1.insert("y".to_string(), Type::int());
        let r1 = Type::Record {
            name: "A".into(),
            fields: f1.clone(),
        };
        let r2 = Type::Record {
            name: "B".into(),
            fields: f1,
        };
        unify(&mut c, &r1, &r2).unwrap();
    }

    #[test]
    fn occurs_check_blocks_infinite_type() {
        let mut c = InferCtx::new();
        let v = c.fresh();
        assert!(unify(&mut c, &v, &Type::list(v.clone())).is_err());
    }
}
