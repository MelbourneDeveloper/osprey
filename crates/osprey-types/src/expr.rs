//! Expression inference. One `infer_expr` dispatch covers every `ast::Expr`.
//! Where a type genuinely cannot be resolved (an opaque field access, an
//! unknown dynamic builtin) the inferencer yields a fresh variable rather than
//! a false error; the structured cases (calls, arithmetic, constructors,
//! lambdas, match) do real unification.

use crate::check::Checker;
use crate::convert::type_expr_to_type;
use crate::env::{instantiate, TypeEnv};
use crate::error::TypeError;
use crate::ty::{names, Type};
use crate::unify::unify;
use osprey_ast::{
    Expr, FieldAssignment, InterpolatedPart, NamedArgument, Parameter, Stmt, TypeExpr,
};
use std::collections::{BTreeMap, HashMap};

fn math_err() -> Type {
    Type::prim(names::MATH_ERROR)
}
fn res_math(ok: Type) -> Type {
    Type::result(ok, math_err())
}
fn generic_err() -> Type {
    Type::prim("Error")
}

impl Checker {
    pub(crate) fn infer_expr(&mut self, e: &Expr, env: &TypeEnv) -> Type {
        match e {
            Expr::Integer(_) => Type::int(),
            Expr::Float(_) => Type::float(),
            Expr::Str(_) => Type::string(),
            Expr::Bool(_) => Type::bool(),
            Expr::InterpolatedStr(parts) => {
                for p in parts {
                    if let InterpolatedPart::Expr(inner) = p {
                        let _ = self.infer_expr(inner, env);
                    }
                }
                Type::string()
            }
            Expr::Identifier(name) => self.lookup_ident(name, env),
            Expr::List(items) => {
                let elem = self.ctx.fresh();
                for it in items {
                    let t = self.infer_expr(it, env);
                    self.push_unify(&elem, &t);
                }
                Type::list(elem)
            }
            Expr::Map(entries) => self.infer_map(entries, env),
            Expr::Object(fields) => self.infer_object(fields, env),
            Expr::Binary { op, left, right } => self.infer_binary(op, left, right, env),
            Expr::Unary { op, operand } => {
                let t = self.infer_expr(operand, env);
                if op == "!" || op == "not" {
                    self.push_assign(&Type::bool(), &t);
                    Type::bool()
                } else {
                    // numeric negation keeps the operand type (int or float)
                    t
                }
            }
            Expr::Call {
                function,
                arguments,
                named_arguments,
            } => self.infer_call(function, arguments, named_arguments, env),
            Expr::Pipe { left, right } => self.infer_pipe(left, right, env),
            Expr::FieldAccess { target, field } => {
                let tt = self.infer_expr(target, env);
                let tp = self.ctx.prune(&tt);
                match &tp {
                    Type::Record { fields, .. } => fields
                        .get(field)
                        .cloned()
                        .unwrap_or_else(|| self.ctx.fresh()),
                    _ => self.ctx.fresh(),
                }
            }
            Expr::MethodCall {
                target,
                method,
                arguments,
                named_arguments,
            } => self.infer_method_call(target, method, arguments, named_arguments, env),
            Expr::Index { target, index } => self.infer_index(target, index, env),
            Expr::Lambda {
                parameters,
                return_type,
                body,
            } => self.infer_lambda(parameters, return_type.as_ref(), body, env),
            Expr::Match { value, arms } => self.infer_match(value, arms, env),
            Expr::Block { statements, value } => {
                self.infer_block(statements, value.as_deref(), env)
            }
            Expr::TypeConstructor { name, fields, .. } => self.infer_constructor(name, fields, env),
            Expr::Update { record, fields } => self.infer_update(record, fields, env),
            Expr::Spawn(inner) => {
                let t = self.infer_expr(inner, env);
                Type::con(names::FIBER, vec![t])
            }
            Expr::Await(inner) => self.infer_unwrap_con(inner, names::FIBER, env),
            Expr::Recv(channel) => self.infer_unwrap_con(channel, names::CHANNEL, env),
            Expr::Send { channel, value } => {
                let _ = self.infer_expr(channel, env);
                let _ = self.infer_expr(value, env);
                Type::unit()
            }
            Expr::Yield(inner) => {
                if let Some(inner) = inner {
                    let _ = self.infer_expr(inner, env);
                }
                Type::unit()
            }
            Expr::Select { arms } => self.infer_arm_bodies(arms, env),
            Expr::Perform {
                effect,
                operation,
                arguments,
                named_arguments,
            } => self.infer_perform(effect, operation, arguments, named_arguments, env),
            Expr::Handler { arms, body, .. } => self.infer_handler(arms, body, env),
        }
    }

    /// Infer a map literal: unify all keys to one type and all values to another.
    fn infer_map(&mut self, entries: &[osprey_ast::MapEntry], env: &TypeEnv) -> Type {
        let (k, v) = (self.ctx.fresh(), self.ctx.fresh());
        for entry in entries {
            let kt = self.infer_expr(&entry.key, env);
            let vt = self.infer_expr(&entry.value, env);
            self.push_unify(&k, &kt);
            self.push_unify(&v, &vt);
        }
        Type::map(k, v)
    }

    /// Infer an anonymous object literal as an unnamed record of its fields.
    fn infer_object(&mut self, fields: &[FieldAssignment], env: &TypeEnv) -> Type {
        let mut map = BTreeMap::new();
        for fa in fields {
            let t = self.infer_expr(&fa.value, env);
            let _ = map.insert(fa.name.clone(), t);
        }
        Type::Record {
            name: String::new(),
            fields: map,
        }
    }

    /// Unwrap `await`/`recv`: unify the inner type with `con<elem>` and yield `elem`.
    fn infer_unwrap_con(&mut self, inner: &Expr, con: &str, env: &TypeEnv) -> Type {
        let t = self.infer_expr(inner, env);
        let elem = self.ctx.fresh();
        self.push_unify(&t, &Type::con(con, vec![elem.clone()]));
        elem
    }

    /// Infer a `perform`: walk its arguments, then yield the operation result type.
    fn infer_perform(
        &mut self,
        effect: &str,
        operation: &str,
        arguments: &[Expr],
        named_arguments: &[NamedArgument],
        env: &TypeEnv,
    ) -> Type {
        for a in arguments {
            let _ = self.infer_expr(a, env);
        }
        for na in named_arguments {
            let _ = self.infer_expr(&na.value, env);
        }
        let ret = self
            .effects
            .get(effect)
            .and_then(|ops| ops.get(operation))
            .cloned();
        ret.unwrap_or_else(|| self.ctx.fresh())
    }

    /// Infer a `handle`: type each arm body in a child scope, then the handled body.
    fn infer_handler(
        &mut self,
        arms: &[osprey_ast::HandlerArm],
        body: &Expr,
        env: &TypeEnv,
    ) -> Type {
        for arm in arms {
            let mut local = env.child();
            for p in &arm.params {
                let fv = self.ctx.fresh();
                local.insert(p.clone(), crate::ty::Scheme::mono(fv));
            }
            let _ = self.infer_expr(&arm.body, &local);
        }
        self.infer_expr(body, env)
    }

    fn lookup_ident(&mut self, name: &str, env: &TypeEnv) -> Type {
        // A bare nullary constructor (`Red`, `Empty`) is a value of its owner type.
        if self.ctors.get(name).is_some_and(|i| i.fields.is_empty()) {
            if let Some((args, _f, owner, is_record)) = self.ctor_instance(name) {
                return if is_record {
                    Type::Record {
                        name: owner,
                        fields: BTreeMap::new(),
                    }
                } else {
                    Type::con(owner, args)
                };
            }
        }
        if let Some(scheme) = env.get(name).cloned() {
            return instantiate(&mut self.ctx, &scheme);
        }
        self.errors
            .push(TypeError::new(format!("unknown identifier `{name}`")));
        self.ctx.fresh()
    }

    fn infer_call(
        &mut self,
        function: &Expr,
        arguments: &[Expr],
        named: &[NamedArgument],
        env: &TypeEnv,
    ) -> Type {
        let (fname, ft) = match function {
            Expr::Identifier(n) => (Some(n.clone()), self.lookup_ident(n, env)),
            other => (None, self.infer_expr(other, env)),
        };
        let args = self.ordered_arg_types(fname.as_deref(), arguments, named, env);
        self.apply_fn(&ft, args)
    }

    fn infer_method_call(
        &mut self,
        target: &Expr,
        method: &str,
        arguments: &[Expr],
        named: &[NamedArgument],
        env: &TypeEnv,
    ) -> Type {
        // UFCS: `t.m(a)` is `m(t, a)`.
        let ft = self.lookup_ident(method, env);
        let mut args = vec![self.infer_expr(target, env)];
        for a in arguments {
            args.push(self.infer_expr(a, env));
        }
        for na in named {
            args.push(self.infer_expr(&na.value, env));
        }
        self.apply_fn(&ft, args)
    }

    /// Resolve call arguments to types, reordering named arguments to the
    /// declared parameter order when the callee is a known function.
    fn ordered_arg_types(
        &mut self,
        fname: Option<&str>,
        arguments: &[Expr],
        named: &[NamedArgument],
        env: &TypeEnv,
    ) -> Vec<Type> {
        if !named.is_empty() {
            if let Some(pnames) = fname.and_then(|n| self.fn_params.get(n).cloned()) {
                let mut out = Vec::new();
                for pn in &pnames {
                    if let Some(na) = named.iter().find(|a| &a.name == pn) {
                        out.push(self.infer_expr(&na.value, env));
                    }
                }
                if out.len() == named.len() {
                    return out;
                }
            }
            return named
                .iter()
                .map(|a| self.infer_expr(&a.value, env))
                .collect();
        }
        arguments.iter().map(|a| self.infer_expr(a, env)).collect()
    }

    fn apply_fn(&mut self, ft: &Type, args: Vec<Type>) -> Type {
        match self.ctx.prune(ft) {
            Type::Fun { params, ret } => {
                if params.len() != args.len() {
                    self.errors.push(TypeError::new(format!(
                        "call arity mismatch: expected {} argument(s), got {}",
                        params.len(),
                        args.len()
                    )));
                    return *ret;
                }
                for (p, a) in params.iter().zip(&args) {
                    self.push_assign(p, a);
                }
                *ret
            }
            ft @ Type::Var(_) => {
                let ret = self.ctx.fresh();
                let f = Type::fun(args, ret.clone());
                let _ = unify(&mut self.ctx, &ft, &f);
                ret
            }
            other => {
                self.errors.push(TypeError::new(format!(
                    "cannot call non-function `{other}`"
                )));
                self.ctx.fresh()
            }
        }
    }

    fn infer_pipe(&mut self, left: &Expr, right: &Expr, env: &TypeEnv) -> Type {
        if let Expr::Call {
            function,
            arguments,
            named_arguments,
        } = right
        {
            let mut args = Vec::with_capacity(arguments.len() + 1);
            args.push(left.clone());
            args.extend(arguments.iter().cloned());
            let call = Expr::Call {
                function: function.clone(),
                arguments: args,
                named_arguments: named_arguments.clone(),
            };
            self.infer_expr(&call, env)
        } else {
            let ft = self.infer_expr(right, env);
            let lt = self.infer_expr(left, env);
            self.apply_fn(&ft, vec![lt])
        }
    }

    fn infer_index(&mut self, target: &Expr, index: &Expr, env: &TypeEnv) -> Type {
        let tt = self.infer_expr(target, env);
        let _ = self.infer_expr(index, env);
        match self.ctx.prune(&tt) {
            Type::Con { name, args } if name == names::LIST && !args.is_empty() => {
                res_math_like(args.first().cloned().unwrap_or_else(|| self.ctx.fresh()))
            }
            Type::Con { name, args } if name == names::MAP && args.len() == 2 => {
                res_math_like(args.get(1).cloned().unwrap_or_else(|| self.ctx.fresh()))
            }
            t if t.is_named(names::STRING) => res_math_like(Type::string()),
            _ => {
                let fresh = self.ctx.fresh();
                res_math_like(fresh)
            }
        }
    }

    fn infer_lambda(
        &mut self,
        parameters: &[Parameter],
        return_type: Option<&TypeExpr>,
        body: &Expr,
        env: &TypeEnv,
    ) -> Type {
        let empty = HashMap::new();
        let mut local = env.child();
        let mut ptys = Vec::new();
        for p in parameters {
            let ty = match &p.ty {
                Some(te) => type_expr_to_type(te, &empty),
                None => self.ctx.fresh(),
            };
            local.insert(p.name.clone(), crate::ty::Scheme::mono(ty.clone()));
            ptys.push(ty);
        }
        let body_ty = self.infer_expr(body, &local);
        let ret = match return_type {
            Some(te) => {
                let r = type_expr_to_type(te, &empty);
                self.push_assign(&r, &body_ty);
                r
            }
            None => body_ty,
        };
        Type::fun(ptys, ret)
    }

    fn infer_block(&mut self, statements: &[Stmt], value: Option<&Expr>, env: &TypeEnv) -> Type {
        let mut local = env.child();
        for s in statements {
            self.infer_block_stmt(s, &mut local);
        }
        match value {
            Some(v) => self.infer_expr(v, &local),
            None => Type::unit(),
        }
    }

    fn infer_constructor(&mut self, name: &str, fields: &[FieldAssignment], env: &TypeEnv) -> Type {
        if let Some((args, declared, owner, is_record)) = self.ctor_instance(name) {
            let dmap: BTreeMap<String, Type> = declared.into_iter().collect();
            for fa in fields {
                let vt = self.infer_expr(&fa.value, env);
                if let Some(dt) = dmap.get(&fa.name) {
                    self.push_assign(&dt.clone(), &vt);
                }
            }
            if is_record {
                Type::Record {
                    name: owner,
                    fields: dmap,
                }
            } else {
                Type::con(owner, args)
            }
        } else {
            // The grammar lowers a record update `rec { f: v }` over a
            // lower-cased binding as a constructor; recover it as an update
            // when the name resolves to an in-scope record.
            if env.get(name).is_some() {
                return self.infer_update(name, fields, env);
            }
            for fa in fields {
                let _ = self.infer_expr(&fa.value, env);
            }
            self.errors
                .push(TypeError::new(format!("unknown constructor `{name}`")));
            self.ctx.fresh()
        }
    }

    fn infer_update(&mut self, record: &str, fields: &[FieldAssignment], env: &TypeEnv) -> Type {
        let base = self.lookup_ident(record, env);
        let base_p = self.ctx.prune(&base);
        if let Type::Record { fields: rf, .. } = &base_p {
            let rf = rf.clone();
            for fa in fields {
                let vt = self.infer_expr(&fa.value, env);
                if let Some(dt) = rf.get(&fa.name) {
                    self.push_assign(&dt.clone(), &vt);
                }
            }
        } else {
            for fa in fields {
                let _ = self.infer_expr(&fa.value, env);
            }
        }
        base_p
    }
}

fn res_math_like(ok: Type) -> Type {
    Type::result(ok, generic_err())
}

fn both_vars(l: &Type, r: &Type) -> bool {
    matches!(l, Type::Var(_)) && matches!(r, Type::Var(_))
}

/// Operator → result type. Lives free of `self` so the borrow checker is happy.
fn unwrap_result(t: &Type) -> Type {
    match t {
        Type::Con { name, args } if name == names::RESULT => {
            args.first().cloned().unwrap_or_else(|| t.clone())
        }
        _ => t.clone(),
    }
}

impl Checker {
    fn infer_binary(&mut self, op: &str, left: &Expr, right: &Expr, env: &TypeEnv) -> Type {
        let lt = self.infer_expr(left, env);
        let rt = self.infer_expr(right, env);
        match classify(op) {
            OpKind::Logical => {
                self.push_assign(&Type::bool(), &lt);
                self.push_assign(&Type::bool(), &rt);
                Type::bool()
            }
            OpKind::Comparison => {
                let lu = unwrap_result(&self.ctx.prune(&lt));
                let ru = unwrap_result(&self.ctx.prune(&rt));
                let _ = unify(&mut self.ctx, &lu, &ru);
                Type::bool()
            }
            OpKind::Arith => self.infer_arith(op, &lt, &rt),
        }
    }

    fn infer_arith(&mut self, op: &str, lt: &Type, rt: &Type) -> Type {
        let l = self.ctx.prune(lt);
        let r = self.ctx.prune(rt);
        match op {
            "%" => {
                self.push_assign(&Type::int(), lt);
                self.push_assign(&Type::int(), rt);
                res_math(Type::int())
            }
            "/" => res_math(Type::float()),
            "+" => {
                if l.is_named(names::STRING) || r.is_named(names::STRING) {
                    self.push_assign(&Type::string(), lt);
                    self.push_assign(&Type::string(), rt);
                    Type::string()
                } else if l.is_named(names::FLOAT) || r.is_named(names::FLOAT) {
                    res_math(Type::float())
                } else if l.is_named(names::LIST) {
                    let _ = unify(&mut self.ctx, lt, rt);
                    l
                } else if r.is_named(names::LIST) {
                    let _ = unify(&mut self.ctx, lt, rt);
                    r
                } else if l.is_named(names::MAP) || r.is_named(names::MAP) {
                    let _ = unify(&mut self.ctx, lt, rt);
                    if l.is_named(names::MAP) {
                        l
                    } else {
                        r
                    }
                } else if both_vars(&l, &r) {
                    // Both operands unconstrained: defer (`+` is overloaded over
                    // int/float/string/list). Tie them and yield a fresh result
                    // so usage context can pick the type.
                    let _ = unify(&mut self.ctx, lt, rt);
                    self.ctx.fresh()
                } else {
                    self.push_assign(&Type::int(), lt);
                    self.push_assign(&Type::int(), rt);
                    res_math(Type::int())
                }
            }
            // "-" and "*": unlike "+", these have no string/list overload, so
            // unconstrained operands default to int — `fn square(v) = v * v`
            // infers `(int) -> Result<int, MathError>`.
            _ => {
                if l.is_named(names::FLOAT) || r.is_named(names::FLOAT) {
                    res_math(Type::float())
                } else {
                    self.push_assign(&Type::int(), lt);
                    self.push_assign(&Type::int(), rt);
                    res_math(Type::int())
                }
            }
        }
    }
}

enum OpKind {
    Arith,
    Comparison,
    Logical,
}

fn classify(op: &str) -> OpKind {
    match op {
        "&&" | "||" => OpKind::Logical,
        "==" | "!=" | "<" | "<=" | ">" | ">=" => OpKind::Comparison,
        _ => OpKind::Arith,
    }
}
