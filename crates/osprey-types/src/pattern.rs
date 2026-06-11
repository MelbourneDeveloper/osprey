//! Pattern inference and match exhaustiveness.
//!
//! Binding a constructor pattern unifies the discriminant with the
//! constructor's owner type, so the discriminant's type arguments flow into the
//! bound field types (`Success { value }` over `Result<int, E>` binds
//! `value : int`). Exhaustiveness is enforced only where the checker can decide
//! it with confidence — `bool` and known union/`Result` discriminants — and is
//! otherwise deferred to a catch-all, so it never reports a false positive.

use crate::check::Checker;
use crate::convert::type_expr_to_type;
use crate::env::TypeEnv;
use crate::error::TypeError;
use crate::ty::{names, Scheme, Type};
use osprey_ast::{Expr, MatchArm, Pattern};
use std::collections::{BTreeMap, HashMap, HashSet};

fn unwrap_result(t: &Type) -> Type {
    match t {
        Type::Con { name, args } if name == names::RESULT => match args.first() {
            Some(first) => first.clone(),
            None => t.clone(),
        },
        _ => t.clone(),
    }
}

impl Checker {
    pub(crate) fn infer_match(&mut self, value: &Expr, arms: &[MatchArm], env: &TypeEnv) -> Type {
        let disc = self.infer_expr(value, env);
        let result = self.ctx.fresh();
        for arm in arms {
            let mut local = env.child();
            self.bind_pattern(&arm.pattern, &disc, &mut local);
            let body_ty = self.infer_expr(&arm.body, &local);
            self.push_unify(&result, &body_ty);
        }
        self.check_exhaustive(&disc, arms);
        result
    }

    /// `select { pattern => body ... }` — same arm-typing as match without a
    /// concrete discriminant.
    pub(crate) fn infer_arm_bodies(&mut self, arms: &[MatchArm], env: &TypeEnv) -> Type {
        let result = self.ctx.fresh();
        for arm in arms {
            let mut local = env.child();
            let disc = self.ctx.fresh();
            self.bind_pattern(&arm.pattern, &disc, &mut local);
            let body_ty = self.infer_expr(&arm.body, &local);
            self.push_unify(&result, &body_ty);
        }
        result
    }

    fn bind_pattern(&mut self, pattern: &Pattern, disc: &Type, local: &mut TypeEnv) {
        match pattern {
            Pattern::Wildcard => {}
            Pattern::Binding(name) => self.bind_binding(name, disc, local),
            Pattern::Literal(expr) => {
                let lt = self.infer_expr(expr, local);
                let du = unwrap_result(&self.ctx.prune(disc));
                self.push_unify(&du, &lt);
            }
            Pattern::TypeAnnotated { name, ty } => {
                let t = type_expr_to_type(ty, &HashMap::new());
                local.insert(name.clone(), Scheme::mono(t));
            }
            Pattern::Structural { fields } => {
                let dp = self.ctx.prune(disc);
                for fname in fields {
                    let ft = match &dp {
                        Type::Record { fields: rf, .. } => {
                            rf.get(fname).cloned().unwrap_or_else(|| self.ctx.fresh())
                        }
                        _ => self.ctx.fresh(),
                    };
                    local.insert(fname.clone(), Scheme::mono(ft));
                }
            }
            Pattern::Constructor {
                name,
                fields,
                sub_patterns,
            } => self.bind_constructor(name, fields, sub_patterns, disc, local),
        }
    }

    /// A bare identifier pattern is either a nullary constructor (matches that
    /// variant) or a fresh variable binding.
    fn bind_binding(&mut self, name: &str, disc: &Type, local: &mut TypeEnv) {
        if self.ctors.get(name).is_some_and(|i| i.fields.is_empty()) {
            if let Some((args, _f, owner, is_record)) = self.ctor_instance(name) {
                let owner_ty = nullary_owner_ty(owner, args, is_record);
                self.push_unify(&owner_ty, disc);
                return;
            }
        }
        local.insert(name.to_string(), Scheme::mono(disc.clone()));
    }

    fn bind_constructor(
        &mut self,
        name: &str,
        fields: &[String],
        sub_patterns: &[Pattern],
        disc: &Type,
        local: &mut TypeEnv,
    ) {
        let Some((args, declared, owner, is_record)) = self.ctor_instance(name) else {
            for f in fields {
                let fv = self.ctx.fresh();
                local.insert(f.clone(), Scheme::mono(fv));
            }
            return;
        };
        // `Result` patterns (`Success`/`Error`) auto-wrap a non-Result
        // discriminant: `match a + b { Success { value } => .. }` over a `string`
        // binds `value : string` (the match auto-wrap rule: any value may be
        // matched as if wrapped in `Success`). This also lets validated record
        // constructions be matched without a real Result.
        if owner == "Result" {
            let dp = self.ctx.prune(disc);
            let ok = match &dp {
                Type::Con { name, args } if name == names::RESULT && !args.is_empty() => {
                    args.first().cloned().unwrap_or_else(|| dp.clone())
                }
                _ => dp.clone(),
            };
            for fname in fields {
                let ft = match fname.as_str() {
                    "value" => ok.clone(),
                    "message" => Type::string(),
                    _ => self.ctx.fresh(),
                };
                local.insert(fname.clone(), Scheme::mono(ft));
            }
            return;
        }
        let declared_map: BTreeMap<String, Type> = declared.iter().cloned().collect();
        // Tie the discriminant's type arguments to this constructor's owner.
        let owner_ty = if is_record {
            Type::Record {
                name: owner,
                fields: declared_map.clone(),
            }
        } else {
            Type::con(owner, args)
        };
        self.push_unify(&owner_ty, disc);

        // Named field destructure: `Ctor { a, b }`.
        for fname in fields {
            let ft = declared_map
                .get(fname)
                .cloned()
                .unwrap_or_else(|| self.ctx.fresh());
            local.insert(fname.clone(), Scheme::mono(ft));
        }
        // Positional destructure: `Ctor(p0, p1)`.
        for (i, sub) in sub_patterns.iter().enumerate() {
            let field_ty = declared.get(i).map(|(_, t)| t.clone());
            let target = field_ty.unwrap_or_else(|| self.ctx.fresh());
            self.bind_pattern(sub, &target, local);
        }
    }

    /// Enforce exhaustiveness where it is unambiguous: `bool` needs both
    /// constructors; a known union/`Result` needs every variant — unless a
    /// catch-all arm is present.
    fn check_exhaustive(&mut self, disc: &Type, arms: &[MatchArm]) {
        if arms.iter().any(|a| self.is_catch_all(&a.pattern)) {
            return;
        }
        let dp = self.ctx.apply(disc);
        match &dp {
            t if t.is_named(names::BOOL) => self.check_bool_exhaustive(arms),
            Type::Con { name, .. } if self.union_variants.contains_key(name) => {
                let all = self.union_variants.get(name).cloned().unwrap_or_default();
                let covered: HashSet<String> = arms
                    .iter()
                    .filter_map(|a| self.pattern_ctor_name(&a.pattern))
                    .collect();
                let missing: Vec<String> = all
                    .iter()
                    .filter(|v| !covered.contains(*v))
                    .cloned()
                    .collect();
                if !missing.is_empty() {
                    self.errors.push(TypeError::new(format!(
                        "non-exhaustive match on `{name}`: missing {}",
                        missing.join(", ")
                    )));
                }
            }
            _ => {}
        }
    }

    fn check_bool_exhaustive(&mut self, arms: &[MatchArm]) {
        let mut has_true = false;
        let mut has_false = false;
        for arm in arms {
            if let Pattern::Literal(expr) = &arm.pattern {
                if let Expr::Bool(b) = expr.as_ref() {
                    has_true |= *b;
                    has_false |= !*b;
                }
            }
        }
        if !(has_true && has_false) {
            self.errors.push(TypeError::new(
                "non-exhaustive match on `bool`: needs both true and false",
            ));
        }
    }

    fn is_catch_all(&self, pattern: &Pattern) -> bool {
        match pattern {
            Pattern::Wildcard | Pattern::TypeAnnotated { .. } => true,
            Pattern::Binding(name) => self.ctors.get(name).is_none_or(|i| !i.fields.is_empty()),
            _ => false,
        }
    }

    fn pattern_ctor_name(&self, pattern: &Pattern) -> Option<String> {
        match pattern {
            Pattern::Constructor { name, .. } => Some(name.clone()),
            Pattern::Binding(name) if self.ctors.get(name).is_some_and(|i| i.fields.is_empty()) => {
                Some(name.clone())
            }
            _ => None,
        }
    }
}

fn nullary_owner_ty(owner: String, args: Vec<Type>, is_record: bool) -> Type {
    if is_record {
        Type::Record {
            name: owner,
            fields: BTreeMap::new(),
        }
    } else {
        Type::con(owner, args)
    }
}
