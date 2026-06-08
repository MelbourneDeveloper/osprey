//! `match` lowering. Three shapes, dispatched on the arm patterns:
//!   * literal arms (bool/int/float/string) + catch-all — a compare/branch chain;
//!   * `Success`/`Error` arms — Result discrimination (the unwrapped-scalar
//!     fallback `disc >= 0` ⇒ Success, matching `generateResultMatchCondition`);
//!   * user-union variant arms — tag comparison against the heap block's leading
//!     discriminant, binding the variant's fields.
//! Ports the literal + Result + union paths of `generateMatchExpression`.

use crate::builder::Codegen;
use crate::conv::as_i64;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use osprey_ast::*;

pub(crate) fn gen_match(cg: &mut Codegen, value: &Expr, arms: &[MatchArm]) -> Result<Value> {
    let disc = gen_expr(cg, value)?;
    if arms.iter().any(|a| is_result_arm(&a.pattern)) {
        return gen_result_match(cg, disc, arms);
    }
    if let Some(owner) = union_owner(cg, arms) {
        return gen_union_match(cg, disc, arms, &owner);
    }
    gen_literal_match(cg, disc, arms)
}

fn is_result_arm(p: &Pattern) -> bool {
    matches!(p, Pattern::Constructor { name, .. } if name == "Success" || name == "Error")
}

/// The constructor name a pattern selects, if any: an explicit `Ctor { … }` or a
/// bare `Ctor` (a nullary variant lowers to a `Binding` indistinguishable from a
/// capture until we know the constructor table).
fn pattern_ctor<'a>(cg: &Codegen, p: &'a Pattern) -> Option<(&'a str, &'a [String])> {
    match p {
        Pattern::Constructor { name, fields, .. } => Some((name, fields)),
        Pattern::Binding(name) if cg.is_ctor(name) => Some((name, &[])),
        _ => None,
    }
}

/// If any arm is a user-union variant constructor, the union's owner name.
fn union_owner(cg: &Codegen, arms: &[MatchArm]) -> Option<String> {
    for a in arms {
        if let Some((name, _)) = pattern_ctor(cg, &a.pattern) {
            if let Some(view) = cg.ctor_layout(name) {
                if !view.owner_is_record && cg.union_variants(&view.owner).is_some() {
                    return Some(view.owner);
                }
            }
        }
    }
    None
}

/// Result match over an unwrapped scalar discriminant: `disc >= 0` ⇒ Success
/// (binds the success field to the scalar), else Error.
fn gen_result_match(cg: &mut Codegen, disc: Value, arms: &[MatchArm]) -> Result<Value> {
    let success = arms.iter().find(|a| matches!(&a.pattern,
        Pattern::Constructor { name, .. } if name == "Success"));
    let error = arms.iter().find(|a| matches!(&a.pattern,
        Pattern::Constructor { name, .. } if name == "Error"));

    let di = as_i64(cg, disc)?;
    let cond = cg.fresh_reg();
    cg.emit(format!("{cond} = icmp sge i64 {}, 0", di.operand));
    let sl = cg.fresh_label();
    let el = cg.fresh_label();
    let end = cg.fresh_label();
    cg.emit(format!("br i1 {cond}, label %{sl}, label %{el}"));

    let mut phi_in: Vec<(String, String)> = Vec::new();
    let mut result_ty = LType::I64;

    cg.start_block(&sl);
    if let Some(arm) = success {
        if let Pattern::Constructor { fields, .. } = &arm.pattern {
            if let Some(f) = fields.first() {
                cg.bind(f.clone(), Value::new(di.operand.clone(), LType::I64));
            }
        }
        let v = gen_expr(cg, &arm.body)?;
        result_ty = v.ty;
        phi_in.push((v.operand, cg.cur_block().to_string()));
    }
    cg.emit(format!("br label %{end}"));

    cg.start_block(&el);
    if let Some(arm) = error {
        if let Pattern::Constructor { fields, .. } = &arm.pattern {
            if let Some(f) = fields.first() {
                // No real message in the scalar-fallback path: bind an empty
                // string stand-in so the (dead) Error body still type-checks.
                let empty = cg.string_constant("");
                cg.bind(f.clone(), empty);
            }
        }
        let v = gen_expr(cg, &arm.body)?;
        result_ty = v.ty;
        phi_in.push((v.operand, cg.cur_block().to_string()));
    }
    cg.emit(format!("br label %{end}"));

    cg.start_block(&end);
    finish_phi(cg, &phi_in, result_ty)
}

/// User-union match: read the leading tag of the heap block and branch per
/// variant, binding that variant's fields.
fn gen_union_match(
    cg: &mut Codegen,
    disc: Value,
    arms: &[MatchArm],
    owner: &str,
) -> Result<Value> {
    // Load the discriminant tag (every variant block starts with `{ i64 tag, … }`).
    let tagp = cg.fresh_reg();
    cg.emit(format!(
        "{tagp} = bitcast i8* {} to i64*",
        disc.operand
    ));
    let tag = cg.fresh_reg();
    cg.emit(format!("{tag} = load i64, i64* {tagp}"));

    let end = cg.fresh_label();
    let mut phi_in: Vec<(String, String)> = Vec::new();
    let mut result_ty: Option<LType> = None;
    let variants = cg.union_variants(owner).unwrap_or(&[]).to_vec();

    for arm in arms {
        if let Some((name, fields)) = pattern_ctor(cg, &arm.pattern) {
            let name = name.to_string();
            let fields = fields.to_vec();
            let vtag = variants.iter().position(|v| *v == name).unwrap_or(0) as i64;
            let cond = cg.fresh_reg();
            cg.emit(format!("{cond} = icmp eq i64 {tag}, {vtag}"));
            let body_lbl = cg.fresh_label();
            let next_lbl = cg.fresh_label();
            cg.emit(format!("br i1 {cond}, label %{body_lbl}, label %{next_lbl}"));
            cg.start_block(&body_lbl);
            bind_variant_fields(cg, &disc, &name, &fields);
            let v = gen_expr(cg, &arm.body)?;
            result_ty.get_or_insert(v.ty);
            phi_in.push((v.operand, cg.cur_block().to_string()));
            cg.emit(format!("br label %{end}"));
            cg.start_block(&next_lbl);
        } else {
            match &arm.pattern {
                Pattern::Wildcard | Pattern::Binding(_) | Pattern::TypeAnnotated { .. } => {
                    if let Pattern::Binding(n) | Pattern::TypeAnnotated { name: n, .. } =
                        &arm.pattern
                    {
                        cg.bind(n.clone(), disc.clone());
                    }
                    let v = gen_expr(cg, &arm.body)?;
                    result_ty.get_or_insert(v.ty);
                    phi_in.push((v.operand, cg.cur_block().to_string()));
                    cg.emit(format!("br label %{end}"));
                    break;
                }
                _ => return Err(CodegenError::unsupported("structural union arm")),
            }
        }
    }
    // A non-exhaustive fall-through is unreachable by construction.
    cg.emit("unreachable");
    cg.start_block(&end);
    finish_phi(cg, &phi_in, result_ty.unwrap_or(LType::I64))
}

/// Bind a matched variant's fields (in declared order) from the heap block.
fn bind_variant_fields(cg: &mut Codegen, disc: &Value, variant: &str, pat_fields: &[String]) {
    let Some(view) = cg.ctor_layout(variant) else {
        return;
    };
    let struct_ty = match cg.ctor_struct_ty(variant) {
        Some(s) => s,
        None => return,
    };
    if view.fields.is_empty() || pat_fields.is_empty() {
        return;
    }
    let src = cg.fresh_reg();
    cg.emit(format!(
        "{src} = bitcast i8* {} to {struct_ty}*",
        disc.operand
    ));
    for (i, (declared, fty)) in view.fields.iter().enumerate() {
        let Some(bind_name) = pat_fields.get(i) else {
            break;
        };
        let loaded = crate::aggregate::load_field(cg, &struct_ty, src.as_str(), i + 1, *fty);
        let owner = cg.ctor_field_written(variant, declared);
        cg.bind(bind_name.clone(), Value::new(loaded, *fty).with_owner(owner));
    }
}

/// Literal/catch-all match: compare-and-branch chain joined by a `phi`.
fn gen_literal_match(cg: &mut Codegen, disc: Value, arms: &[MatchArm]) -> Result<Value> {
    let end = cg.fresh_label();
    let mut phi_in: Vec<(String, String)> = Vec::new();
    let mut result_ty: Option<LType> = None;
    let last = arms.len().saturating_sub(1);

    for (i, arm) in arms.iter().enumerate() {
        match &arm.pattern {
            Pattern::Wildcard | Pattern::Binding(_) | Pattern::TypeAnnotated { .. } => {
                bind_catch_all(cg, &arm.pattern, &disc);
                let v = gen_expr(cg, &arm.body)?;
                result_ty.get_or_insert(v.ty);
                phi_in.push((v.operand.clone(), cg.cur_block().to_string()));
                cg.emit(format!("br label %{end}"));
                break;
            }
            Pattern::Literal(lit) => {
                let cond = gen_eq(cg, &disc, lit)?;
                let body_lbl = cg.fresh_label();
                let next_lbl = cg.fresh_label();
                cg.emit(format!("br i1 {cond}, label %{body_lbl}, label %{next_lbl}"));
                cg.start_block(&body_lbl);
                let v = gen_expr(cg, &arm.body)?;
                result_ty.get_or_insert(v.ty);
                phi_in.push((v.operand.clone(), cg.cur_block().to_string()));
                cg.emit(format!("br label %{end}"));
                cg.start_block(&next_lbl);
                if i == last {
                    cg.emit("unreachable");
                }
            }
            _ => return Err(CodegenError::unsupported("destructuring match arm")),
        }
    }

    cg.start_block(&end);
    finish_phi(cg, &phi_in, result_ty.unwrap_or(LType::I64))
}

fn finish_phi(cg: &mut Codegen, phi_in: &[(String, String)], ty: LType) -> Result<Value> {
    if phi_in.is_empty() {
        return Ok(Value::unit());
    }
    let incoming = phi_in
        .iter()
        .map(|(op, blk)| format!("[ {op}, %{blk} ]"))
        .collect::<Vec<_>>()
        .join(", ");
    let reg = cg.fresh_reg();
    cg.emit(format!("{reg} = phi {ty} {incoming}"));
    Ok(Value::new(reg, ty))
}

fn bind_catch_all(cg: &mut Codegen, pattern: &Pattern, disc: &Value) {
    match pattern {
        Pattern::Binding(name) => cg.bind(name.clone(), disc.clone()),
        Pattern::TypeAnnotated { name, .. } => cg.bind(name.clone(), disc.clone()),
        _ => {}
    }
}

/// Equality test between the discriminant and a literal pattern → the `i1`
/// operand.
fn gen_eq(cg: &mut Codegen, disc: &Value, lit: &Expr) -> Result<String> {
    let pat = gen_expr(cg, lit)?;
    let reg = cg.fresh_reg();
    let is_str = |t: LType| t == LType::Str || t == LType::Ptr;
    if is_str(disc.ty) && is_str(pat.ty) {
        cg.add_extern("declare i32 @strcmp(i8*, i8*)");
        let c = cg.fresh_reg();
        cg.emit(format!(
            "{c} = call i32 @strcmp(i8* {}, i8* {})",
            disc.operand, pat.operand
        ));
        cg.emit(format!("{reg} = icmp eq i32 {c}, 0"));
    } else if disc.ty == LType::Double || pat.ty == LType::Double {
        cg.emit(format!(
            "{reg} = fcmp oeq double {}, {}",
            disc.operand, pat.operand
        ));
    } else {
        let d = as_i64(cg, disc.clone())?;
        let p = as_i64(cg, pat)?;
        cg.emit(format!("{reg} = icmp eq i64 {}, {}", d.operand, p.operand));
    }
    Ok(reg)
}
