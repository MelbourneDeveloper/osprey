//! The type checker driver: a two-pass walk over a [`Program`]. Pass one
//! collects every top-level declaration (types + their constructors, effects,
//! externs, function signatures) so forward references and recursion resolve.
//! Pass two infers each function body and top-level statement, unifying against
//! the declared signatures. Ports the orchestration in `type_inference.go`
//! around `InferType` / `ResolveAllEnvironmentTypes`.

use crate::builtins::base_env;
use crate::convert::{parse_fn_sig, type_expr_to_type, type_name_to_type};
use crate::ctx::InferCtx;
use crate::env::{generalize, TypeEnv};
use crate::error::TypeError;
use crate::ty::{Scheme, Type};
use crate::unify::{unify, unify_assignable};
use osprey_ast::{
    EffectOperation, Expr, ExternParameter, Parameter, Position, Program, Stmt, TypeExpr,
    TypeVariant,
};
use std::collections::HashMap;

/// A constructor (record builder, union variant, or built-in `Success`/`Error`).
pub(crate) struct CtorInfo {
    pub owner: String,
    pub owner_is_record: bool,
    pub type_params: Vec<String>,
    /// (field name, field type as written).
    pub fields: Vec<(String, String)>,
}

/// A constructor instantiated against fresh type arguments:
/// (owner type arguments, instantiated `(field, type)` pairs, owner name,
/// whether the owner is a record).
pub(crate) type CtorInstance = (Vec<Type>, Vec<(String, Type)>, String, bool);

/// All cross-cutting declaration tables, plus the inference context.
pub struct Checker {
    pub(crate) ctx: InferCtx,
    pub(crate) errors: Vec<TypeError>,
    pub(crate) ctors: HashMap<String, CtorInfo>,
    /// Effect name -> operation name -> return type.
    pub(crate) effects: HashMap<String, HashMap<String, Type>>,
    /// Union/Result type name -> its variant constructor names (exhaustiveness).
    pub(crate) union_variants: HashMap<String, Vec<String>>,
    /// Function/extern name -> declared parameter names (for named arguments).
    pub(crate) fn_params: HashMap<String, Vec<String>>,
    /// Function name -> the exact (params, ret) types created in pass one, so
    /// body inference reuses the very same variables the signature exported.
    pub(crate) fn_sigs: HashMap<String, (Vec<Type>, Type)>,
}

impl Checker {
    fn new() -> Checker {
        let mut c = Checker {
            ctx: InferCtx::new(),
            errors: Vec::new(),
            ctors: HashMap::new(),
            effects: HashMap::new(),
            union_variants: HashMap::new(),
            fn_params: HashMap::new(),
            fn_sigs: HashMap::new(),
        };
        c.register_result_ctors();
        c
    }

    /// Built-in `Result` constructors `Success { value: T }` / `Error { message: E }`.
    fn register_result_ctors(&mut self) {
        let _ = self.ctors.insert(
            "Success".into(),
            CtorInfo {
                owner: "Result".into(),
                owner_is_record: false,
                type_params: vec!["T".into(), "E".into()],
                fields: vec![("value".into(), "T".into())],
            },
        );
        // `Error { message: <string> }` builds the E side of a `Result<T, E>`;
        // the message is a concrete string, leaving E free to unify with the
        // declared error type (e.g. the nominal `Error`), not pinned to string.
        let _ = self.ctors.insert(
            "Error".into(),
            CtorInfo {
                owner: "Result".into(),
                owner_is_record: false,
                type_params: vec!["T".into(), "E".into()],
                fields: vec![("message".into(), "string".into())],
            },
        );
        let _ = self
            .union_variants
            .insert("Result".into(), vec!["Success".into(), "Error".into()]);
        // Built-in HttpResponse record returned by HTTP request handlers.
        let _ = self.ctors.insert(
            "HttpResponse".into(),
            CtorInfo {
                owner: "HttpResponse".into(),
                owner_is_record: true,
                type_params: Vec::new(),
                fields: vec![
                    ("status".into(), "int".into()),
                    ("headers".into(), "string".into()),
                    ("contentType".into(), "string".into()),
                    ("body".into(), "string".into()),
                ],
            },
        );
    }

    fn record_err(&mut self, e: TypeError, pos: Option<Position>) {
        self.errors.push(e.with_pos(pos));
    }

    /// Unify and record any failure. Shared by the expr/pattern modules.
    pub(crate) fn push_unify(&mut self, a: &Type, b: &Type) {
        if let Err(e) = unify(&mut self.ctx, a, b) {
            self.errors.push(e);
        }
    }

    /// Assignment-site unification (Result auto-unwrap), recording failures.
    pub(crate) fn push_assign(&mut self, expected: &Type, actual: &Type) {
        if let Err(e) = unify_assignable(&mut self.ctx, expected, actual) {
            self.errors.push(e);
        }
    }

    /// Check a statement appearing inside a block expression, threading new
    /// bindings into the block's local scope.
    pub(crate) fn infer_block_stmt(&mut self, s: &Stmt, env: &mut TypeEnv) {
        match s {
            Stmt::Let {
                name,
                ty,
                value,
                position,
                ..
            } => self.check_let(name, ty.as_ref(), value, env, *position),
            Stmt::Assignment {
                name,
                value,
                position,
            } => self.check_assignment(name, value, env, *position),
            Stmt::Expr(e) => {
                let _ = self.infer_expr(e, env);
            }
            _ => {}
        }
    }

    /// Pass one: fill the declaration tables and the base environment.
    fn collect(&mut self, program: &Program, env: &mut TypeEnv) {
        for stmt in &program.statements {
            match stmt {
                Stmt::Type {
                    name,
                    type_params,
                    variants,
                    ..
                } => self.collect_type(name, type_params, variants),
                Stmt::Effect { name, operations } => self.collect_effect(name, operations),
                Stmt::Extern {
                    name,
                    parameters,
                    return_type,
                } => self.collect_extern(name, parameters, return_type.as_ref(), env),
                Stmt::Function {
                    name,
                    parameters,
                    return_type,
                    ..
                } => self.collect_function(name, parameters, return_type.as_ref(), env),
                _ => {}
            }
        }
    }

    fn collect_type(&mut self, name: &str, type_params: &[String], variants: &[TypeVariant]) {
        let is_record = match variants.first() {
            Some(first) => variants.len() == 1 && first.name == name,
            None => false,
        };
        if !is_record {
            let _ = self.union_variants.insert(
                name.to_string(),
                variants.iter().map(|v| v.name.clone()).collect(),
            );
        }
        for v in variants {
            let fields = v
                .fields
                .iter()
                .map(|f| (f.name.clone(), f.ty.clone()))
                .collect();
            let _ = self.ctors.insert(
                v.name.clone(),
                CtorInfo {
                    owner: name.to_string(),
                    owner_is_record: is_record,
                    type_params: type_params.to_vec(),
                    fields,
                },
            );
        }
    }

    fn collect_effect(&mut self, name: &str, operations: &[EffectOperation]) {
        let mut ops = HashMap::new();
        for op in operations {
            let (_, ret) = parse_fn_sig(&op.ty, &HashMap::new());
            let _ = ops.insert(op.name.clone(), ret);
        }
        let _ = self.effects.insert(name.to_string(), ops);
    }

    fn collect_extern(
        &mut self,
        name: &str,
        parameters: &[ExternParameter],
        return_type: Option<&TypeExpr>,
        env: &mut TypeEnv,
    ) {
        let empty = HashMap::new();
        let params: Vec<Type> = parameters
            .iter()
            .map(|p| type_expr_to_type(&p.ty, &empty))
            .collect();
        let ret = return_type.map_or_else(Type::unit, |r| type_expr_to_type(r, &empty));
        let _ = self.fn_params.insert(
            name.to_string(),
            parameters.iter().map(|p| p.name.clone()).collect(),
        );
        env.insert(name, Scheme::mono(Type::fun(params, ret)));
    }

    fn collect_function(
        &mut self,
        name: &str,
        parameters: &[Parameter],
        return_type: Option<&TypeExpr>,
        env: &mut TypeEnv,
    ) {
        let empty = HashMap::new();
        let params: Vec<Type> = parameters
            .iter()
            .map(|p| match &p.ty {
                Some(te) => type_expr_to_type(te, &empty),
                None => self.ctx.fresh(),
            })
            .collect();
        let ret = match return_type {
            Some(te) => type_expr_to_type(te, &empty),
            None => self.ctx.fresh(),
        };
        let _ = self.fn_params.insert(
            name.to_string(),
            parameters.iter().map(|p| p.name.clone()).collect(),
        );
        let _ = self
            .fn_sigs
            .insert(name.to_string(), (params.clone(), ret.clone()));
        env.insert(name, Scheme::mono(Type::fun(params, ret)));
    }

    /// Pass two: infer bodies and run top-level statements.
    fn check(&mut self, program: &Program, env: &mut TypeEnv) {
        for stmt in &program.statements {
            match stmt {
                Stmt::Function {
                    name,
                    parameters,
                    body,
                    position,
                    ..
                } => self.check_function(name, parameters, body, env, *position),
                Stmt::Module { body, .. } => {
                    let mut inner = env.child();
                    let prog = Program {
                        statements: body.clone(),
                    };
                    self.check(&prog, &mut inner);
                }
                Stmt::Let {
                    name,
                    ty,
                    value,
                    position,
                    ..
                } => self.check_let(name, ty.as_ref(), value, env, *position),
                Stmt::Assignment {
                    name,
                    value,
                    position,
                } => self.check_assignment(name, value, env, *position),
                Stmt::Expr(e) => {
                    let _ = self.infer_expr(e, env);
                }
                _ => {}
            }
        }
    }

    fn check_function(
        &mut self,
        name: &str,
        parameters: &[Parameter],
        body: &Expr,
        env: &mut TypeEnv,
        pos: Option<Position>,
    ) {
        let (params, ret) = match self.fn_sigs.get(name) {
            Some(sig) => sig.clone(),
            None => return,
        };
        let mut local = env.child();
        for (p, ty) in parameters.iter().zip(&params) {
            local.insert(p.name.clone(), Scheme::mono(ty.clone()));
        }
        let body_ty = self.infer_expr(body, &local);
        if let Err(e) = unify_assignable(&mut self.ctx, &ret, &body_ty) {
            self.record_err(
                TypeError::new(format!("function `{name}` body: {}", e.message)),
                pos,
            );
        }
        // Generalize the now-constrained signature so later call sites can use
        // the function polymorphically (HM let-generalization for top-level fns).
        // Remove the function's own monomorphic entry first, else its signature
        // variables would count as "free in the environment" and nothing would
        // generalize.
        let fun_ty = Type::fun(params, ret);
        env.remove(name);
        let scheme = generalize(&mut self.ctx, env, &fun_ty);
        env.insert(name, scheme);
    }

    fn check_let(
        &mut self,
        name: &str,
        ty: Option<&TypeExpr>,
        value: &Expr,
        env: &mut TypeEnv,
        pos: Option<Position>,
    ) {
        let value_ty = self.infer_expr(value, env);
        if let Some(te) = ty {
            let annotated = type_expr_to_type(te, &HashMap::new());
            if let Err(e) = unify_assignable(&mut self.ctx, &annotated, &value_ty) {
                self.record_err(TypeError::new(format!("let `{name}`: {}", e.message)), pos);
            }
        }
        let scheme = generalize(&mut self.ctx, env, &value_ty);
        env.insert(name, scheme);
    }

    fn check_assignment(
        &mut self,
        name: &str,
        value: &Expr,
        env: &mut TypeEnv,
        pos: Option<Position>,
    ) {
        let value_ty = self.infer_expr(value, env);
        match env.get(name).cloned() {
            Some(scheme) => {
                let existing = crate::env::instantiate(&mut self.ctx, &scheme);
                if let Err(e) = unify_assignable(&mut self.ctx, &existing, &value_ty) {
                    self.record_err(
                        TypeError::new(format!("assignment to `{name}`: {}", e.message)),
                        pos,
                    );
                }
            }
            None => self.record_err(
                TypeError::new(format!("assignment to undeclared `{name}`")),
                pos,
            ),
        }
    }

    /// Build the instantiated field types of a constructor against fresh type
    /// arguments. Returns (per-type-param fresh var, declared field map).
    pub(crate) fn ctor_instance(&mut self, name: &str) -> Option<CtorInstance> {
        let info = self.ctors.get(name)?;
        let owner = info.owner.clone();
        let is_record = info.owner_is_record;
        let params = info.type_params.clone();
        let raw_fields = info.fields.clone();
        let mut pmap = HashMap::new();
        let mut args = Vec::new();
        for p in &params {
            let v = self.ctx.fresh();
            let _ = pmap.insert(p.clone(), v.clone());
            args.push(v);
        }
        let fields = raw_fields
            .iter()
            .map(|(fname, fty)| (fname.clone(), type_name_to_type(fty, &pmap)))
            .collect();
        Some((args, fields, owner, is_record))
    }
}

/// Type-check a program. Returns every type error found (empty ⇒ well-typed).
#[must_use]
pub fn check_program(program: &Program) -> Vec<TypeError> {
    let mut checker = Checker::new();
    let mut env = base_env();
    checker.collect(program, &mut env);
    checker.check(program, &mut env);
    checker.errors
}

/// Run inference and publish the resolved signatures, constructor layouts and
/// union tags for the code generator. Type errors are intentionally dropped
/// here — codegen runs after `check_program` has gated correctness — so the
/// backend always receives the best-effort resolved shape of every declaration.
#[must_use]
pub fn infer_program(program: &Program) -> crate::info::ProgramTypes {
    use crate::info::{CtorLayout, ProgramTypes};
    let mut checker = Checker::new();
    let mut env = base_env();
    checker.collect(program, &mut env);
    checker.check(program, &mut env);

    let functions = checker
        .fn_sigs
        .iter()
        .map(|(name, (params, ret))| {
            let rp = params.iter().map(|t| checker.ctx.apply(t)).collect();
            let rr = checker.ctx.apply(ret);
            (name.clone(), (rp, rr))
        })
        .collect();
    let ctors = checker
        .ctors
        .iter()
        .map(|(name, info)| {
            (
                name.clone(),
                CtorLayout {
                    owner: info.owner.clone(),
                    owner_is_record: info.owner_is_record,
                    fields: info.fields.clone(),
                },
            )
        })
        .collect();
    let unions = checker.union_variants.clone();
    ProgramTypes {
        functions,
        ctors,
        unions,
    }
}
