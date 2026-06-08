//! The inference context: fresh-variable supply plus the substitution that
//! backs unification. Ports the `Substitution map[int]Type` + `prune` /
//! `occursCheck` / `applySubst` machinery of `type_inference.go`, but as an
//! index-addressed arena (`Vec<Option<Type>>`) — the textbook union-find layout,
//! with path compression on `prune` exactly as the Go did with `ti.subst[id]`.

use crate::ty::{Type, VarId};
use std::collections::BTreeSet;

/// Holds every type variable's binding. Variable ids are indices into `subst`.
#[derive(Debug, Default)]
pub struct InferCtx {
    subst: Vec<Option<Type>>,
}

impl InferCtx {
    /// Create an empty context with no allocated type variables.
    pub fn new() -> InferCtx {
        InferCtx::default()
    }

    /// Allocate a fresh, unbound type variable (`TypeInferer.Fresh`).
    pub fn fresh(&mut self) -> Type {
        let id = VarId::try_from(self.subst.len()).unwrap_or(VarId::MAX);
        self.subst.push(None);
        Type::Var(id)
    }

    /// Follow a variable to its representative, compressing the path. Only the
    /// outermost variable is resolved — nested types are left intact (use
    /// [`InferCtx::apply`] for a deep walk), matching the Go `prune`.
    pub fn prune(&mut self, t: &Type) -> Type {
        if let Type::Var(id) = t {
            let idx = usize::try_from(*id).unwrap_or(usize::MAX);
            if let Some(bound) = self.subst.get(idx).and_then(Option::clone) {
                let pruned = self.prune(&bound);
                if let Some(slot) = self.subst.get_mut(idx) {
                    *slot = Some(pruned.clone());
                }
                return pruned;
            }
        }
        t.clone()
    }

    /// Bind a variable to a type. The caller guarantees the occurs-check passed.
    pub fn bind(&mut self, id: VarId, t: Type) {
        let idx = usize::try_from(id).unwrap_or(usize::MAX);
        if let Some(slot) = self.subst.get_mut(idx) {
            *slot = Some(t);
        }
    }

    /// `occursCheck`: does variable `id` appear anywhere in `t`? Prevents the
    /// construction of infinite types like `t0 ~ List<t0>`.
    pub fn occurs(&mut self, id: VarId, t: &Type) -> bool {
        let t = self.prune(t);
        match &t {
            Type::Var(v) => *v == id,
            Type::Fun { params, ret } => {
                params.iter().any(|p| self.occurs(id, p)) || self.occurs(id, ret)
            }
            Type::Con { args, .. } => args.iter().any(|a| self.occurs(id, a)),
            Type::Record { fields, .. } => fields.values().any(|v| self.occurs(id, v)),
            Type::Union { variants, .. } => variants.iter().any(|v| self.occurs(id, v)),
        }
    }

    /// `applySubst`: fully resolve `t` against the current substitution. The
    /// occurs-check keeps the substitution acyclic, so this terminates.
    pub fn apply(&mut self, t: &Type) -> Type {
        let t = self.prune(t);
        match &t {
            Type::Var(_) => t,
            Type::Con { name, args } => Type::Con {
                name: name.clone(),
                args: args.iter().map(|a| self.apply(a)).collect(),
            },
            Type::Fun { params, ret } => Type::Fun {
                params: params.iter().map(|p| self.apply(p)).collect(),
                ret: Box::new(self.apply(ret)),
            },
            Type::Record { name, fields } => Type::Record {
                name: name.clone(),
                fields: fields
                    .iter()
                    .map(|(k, v)| (k.clone(), self.apply(v)))
                    .collect(),
            },
            Type::Union { name, variants } => Type::Union {
                name: name.clone(),
                variants: variants.iter().map(|v| self.apply(v)).collect(),
            },
        }
    }

    /// Collect the free (unbound) variables of `t` into `out` (`getFreeVars`).
    pub fn free_vars(&mut self, t: &Type, out: &mut BTreeSet<VarId>) {
        let t = self.prune(t);
        match &t {
            Type::Var(v) => {
                let _ = out.insert(*v);
            }
            Type::Fun { params, ret } => {
                for p in params {
                    self.free_vars(p, out);
                }
                self.free_vars(ret, out);
            }
            Type::Con { args, .. } => {
                for a in args {
                    self.free_vars(a, out);
                }
            }
            Type::Record { fields, .. } => {
                for v in fields.values() {
                    self.free_vars(v, out);
                }
            }
            Type::Union { variants, .. } => {
                for v in variants {
                    self.free_vars(v, out);
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn fresh_vars_are_distinct_and_unbound() {
        let mut c = InferCtx::new();
        let a = c.fresh();
        let b = c.fresh();
        assert_ne!(a, b);
        assert_eq!(c.prune(&a), a);
    }

    #[test]
    fn prune_follows_and_compresses_chains() {
        let mut c = InferCtx::new();
        let (a, b) = (c.fresh(), c.fresh());
        if let (Type::Var(ia), Type::Var(_ib)) = (&a, &b) {
            c.bind(*ia, b.clone());
            let bb = if let Type::Var(ib) = &b { *ib } else { 0 };
            c.bind(bb, Type::int());
        }
        assert_eq!(c.prune(&a), Type::int());
    }

    #[test]
    fn occurs_check_detects_cycles() {
        let mut c = InferCtx::new();
        let a = c.fresh();
        let id = if let Type::Var(id) = a { id } else { 0 };
        assert!(c.occurs(id, &Type::list(a.clone())));
        assert!(!c.occurs(id, &Type::list(Type::int())));
    }
}
