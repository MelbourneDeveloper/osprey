//! The typing environment and the two operations that make HM polymorphic:
//! `instantiate` (fresh-rename a scheme's bound vars at each use) and
//! `generalize` (quantify a let-binding over the vars not free in the
//! environment). Ports `TypeEnv`, `TypeInferer.Instantiate` and
//! `TypeInferer.Generalize`.

use crate::ctx::InferCtx;
use crate::ty::{Scheme, Type, VarId};
use std::collections::{BTreeSet, HashMap};

/// Maps names to their type schemes. Cloned to form child scopes (lambda
/// bodies, match arms) — the same value-semantics the Go `TypeEnv.Clone` gave.
#[derive(Debug, Clone, Default)]
pub struct TypeEnv {
    vars: HashMap<String, Scheme>,
}

impl TypeEnv {
    pub fn new() -> TypeEnv {
        TypeEnv::default()
    }

    pub fn get(&self, name: &str) -> Option<&Scheme> {
        self.vars.get(name)
    }

    pub fn insert(&mut self, name: impl Into<String>, scheme: Scheme) {
        let _ = self.vars.insert(name.into(), scheme);
    }

    pub fn remove(&mut self, name: &str) {
        let _ = self.vars.remove(name);
    }

    /// A fresh child scope (a clone — bindings added to the child don't leak).
    pub fn child(&self) -> TypeEnv {
        self.clone()
    }

    /// The free variables of the whole environment — the vars `generalize` must
    /// *not* quantify, because an outer scope may still constrain them.
    pub fn free_vars(&self, ctx: &mut InferCtx) -> BTreeSet<VarId> {
        let mut out = BTreeSet::new();
        for scheme in self.vars.values() {
            let mut fv = BTreeSet::new();
            ctx.free_vars(&scheme.ty, &mut fv);
            for q in &scheme.vars {
                let _ = fv.remove(q);
            }
            out.extend(fv);
        }
        out
    }
}

/// Instantiate a scheme: replace each quantified variable with a fresh one.
pub fn instantiate(ctx: &mut InferCtx, scheme: &Scheme) -> Type {
    if scheme.vars.is_empty() {
        return scheme.ty.clone();
    }
    let map: HashMap<VarId, Type> = scheme.vars.iter().map(|v| (*v, ctx.fresh())).collect();
    subst_vars(&scheme.ty, &map)
}

/// Generalize a type over the variables free in it but not in the environment.
pub fn generalize(ctx: &mut InferCtx, env: &TypeEnv, ty: &Type) -> Scheme {
    let ty = ctx.apply(ty);
    let env_fv = env.free_vars(ctx);
    let mut ty_fv = BTreeSet::new();
    ctx.free_vars(&ty, &mut ty_fv);
    let vars: Vec<VarId> = ty_fv.difference(&env_fv).copied().collect();
    Scheme { vars, ty }
}

fn subst_vars(t: &Type, map: &HashMap<VarId, Type>) -> Type {
    match t {
        Type::Var(v) => map.get(v).cloned().unwrap_or_else(|| t.clone()),
        Type::Con { name, args } => Type::Con {
            name: name.clone(),
            args: args.iter().map(|a| subst_vars(a, map)).collect(),
        },
        Type::Fun { params, ret } => Type::Fun {
            params: params.iter().map(|p| subst_vars(p, map)).collect(),
            ret: Box::new(subst_vars(ret, map)),
        },
        Type::Record { name, fields } => Type::Record {
            name: name.clone(),
            fields: fields
                .iter()
                .map(|(k, v)| (k.clone(), subst_vars(v, map)))
                .collect(),
        },
        Type::Union { name, variants } => Type::Union {
            name: name.clone(),
            variants: variants.iter().map(|v| subst_vars(v, map)).collect(),
        },
    }
}

#[cfg(test)]
#[expect(
    clippy::indexing_slicing,
    reason = "test assertions: an out-of-bounds index is a test failure, not a production panic"
)]
mod tests {
    use super::*;
    use crate::unify::unify;

    #[test]
    fn instantiation_gives_fresh_independent_vars() {
        // id : forall t0. (t0) -> t0
        let mut ctx = InferCtx::new();
        let scheme = Scheme::poly(vec![0], Type::fun(vec![Type::Var(0)], Type::Var(0)));

        // Two uses must not share a variable: id(1) and id("x") both type-check.
        let a = instantiate(&mut ctx, &scheme);
        let b = instantiate(&mut ctx, &scheme);
        if let (Type::Fun { params: pa, .. }, Type::Fun { params: pb, .. }) = (&a, &b) {
            unify(&mut ctx, &pa[0], &Type::int()).unwrap();
            unify(&mut ctx, &pb[0], &Type::string()).unwrap();
        } else {
            panic!("expected function types");
        }
    }

    #[test]
    fn generalize_skips_env_bound_vars() {
        let mut ctx = InferCtx::new();
        let mut env = TypeEnv::new();
        // An outer binding pins t0, so a let over (t0 -> t1) may only quantify t1.
        env.insert("outer", Scheme::mono(Type::Var(0)));
        let scheme = generalize(&mut ctx, &env, &Type::fun(vec![Type::Var(0)], Type::Var(1)));
        assert_eq!(scheme.vars, vec![1]);
    }
}
