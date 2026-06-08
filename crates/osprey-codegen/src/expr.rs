//! Expression lowering — the type-driven walk ported from
//! `expression_generation.go`. Every node returns a [`Value`] carrying its LLVM
//! type, seeded by inference (`osprey-types`) for the things a local walk cannot
//! know: function parameter and return types. Unsupported nodes fail loudly via
//! [`CodegenError::Unsupported`] rather than miscompiling.

use crate::builder::Codegen;
use crate::conv::{as_double, as_i1, as_i64};
use crate::error::{CodegenError, Result};
use crate::llty::{LType, Value};
use crate::pattern::gen_match;
use crate::runtime::{gen_print, to_string_value};
use osprey_ast::*;

pub(crate) fn gen_expr(cg: &mut Codegen, expr: &Expr) -> Result<Value> {
    match expr {
        Expr::Integer(n) => Ok(Value::new(n.to_string(), LType::I64)),
        Expr::Float(f) => Ok(Value::new(fmt_double(*f), LType::Double)),
        Expr::Bool(b) => Ok(Value::new(if *b { "1" } else { "0" }, LType::I1)),
        Expr::Str(s) => Ok(cg.string_constant(s)),
        Expr::InterpolatedStr(parts) => gen_interpolation(cg, parts),
        Expr::Identifier(name) => match cg.lookup(name) {
            Some(v) => Ok(v),
            // A bare name that is a nullary constructor (`Active`, `Red`, …) is a
            // zero-field variant value.
            None if cg.is_ctor(name) => crate::aggregate::gen_constructor(cg, name, &[]),
            None => Err(CodegenError::unknown(name)),
        },
        Expr::Binary { op, left, right } => gen_binary(cg, op, left, right),
        Expr::Unary { op, operand } => gen_unary(cg, op, operand),
        Expr::Call {
            function,
            arguments,
            named_arguments,
        } => gen_call(cg, function, arguments, named_arguments),
        Expr::Match { value, arms } => gen_match(cg, value, arms),
        Expr::Block { statements, value } => gen_block(cg, statements, value.as_deref()),
        Expr::TypeConstructor { name, fields, .. } => {
            crate::aggregate::gen_constructor(cg, name, fields)
        }
        Expr::Update { record, fields } => crate::aggregate::gen_update(cg, record, fields),
        Expr::FieldAccess { target, field } => {
            crate::aggregate::gen_field_access(cg, target, field)
        }
        other => Err(CodegenError::unsupported(describe(other))),
    }
}

/// LLVM requires a decimal point or exponent in a `double` literal; render a
/// whole number as `N.0`.
fn fmt_double(f: f64) -> String {
    if f.is_finite() && f.fract() == 0.0 {
        format!("{f:.1}")
    } else {
        // Hex float is the exact, locale-free spelling LLVM accepts.
        format!("0x{:016X}", f.to_bits())
    }
}

fn gen_block(cg: &mut Codegen, statements: &[Stmt], value: Option<&Expr>) -> Result<Value> {
    cg.push_scope();
    for s in statements {
        crate::lower::gen_local_stmt(cg, s)?;
    }
    let v = match value {
        Some(e) => gen_expr(cg, e)?,
        None => Value::unit(),
    };
    cg.pop_scope();
    Ok(v)
}

fn gen_binary(cg: &mut Codegen, op: &str, left: &Expr, right: &Expr) -> Result<Value> {
    // Logical operators are control flow over booleans; keep them lazy-safe by
    // evaluating both sides (the lowered programs have pure operands).
    if op == "&&" || op == "||" {
        let l = gen_expr(cg, left)?;
        let r = gen_expr(cg, right)?;
        let lb = as_i1(cg, l)?;
        let rb = as_i1(cg, r)?;
        let opc = if op == "&&" { "and" } else { "or" };
        let reg = cg.fresh_reg();
        cg.emit(format!("{reg} = {opc} i1 {}, {}", lb.operand, rb.operand));
        return Ok(Value::new(reg, LType::I1));
    }

    let l = gen_expr(cg, left)?;
    let r = gen_expr(cg, right)?;
    match op {
        "+" | "-" | "*" | "/" | "%" => gen_arith(cg, op, l, r),
        "==" | "!=" | "<" | "<=" | ">" | ">=" => gen_comparison(cg, op, l, r),
        other => Err(CodegenError::unsupported(format!("binary operator `{other}`"))),
    }
}

/// Arithmetic. Float if either operand is a float (the other is promoted),
/// otherwise integer. Division ALWAYS returns float (the Osprey spec — see
/// `generateDivisionWithZeroCheck`); modulo stays integer. The Result<…,
/// MathError> wrapper the type system tracks is auto-unwrapped at value sites.
fn gen_arith(cg: &mut Codegen, op: &str, l: Value, r: Value) -> Result<Value> {
    // `+` with a string operand is concatenation (libc strlen/strcpy/strcat),
    // matching `generateStringConcatenation`.
    if op == "+" && (l.ty == LType::Str || r.ty == LType::Str) {
        return gen_str_concat(cg, l, r);
    }
    if op == "/" {
        let ld = as_double(cg, l)?;
        let rd = as_double(cg, r)?;
        let reg = cg.fresh_reg();
        cg.emit(format!("{reg} = fdiv double {}, {}", ld.operand, rd.operand));
        return Ok(Value::new(reg, LType::Double));
    }
    if l.ty == LType::Double || r.ty == LType::Double {
        let ld = as_double(cg, l)?;
        let rd = as_double(cg, r)?;
        let opc = match op {
            "+" => "fadd",
            "-" => "fsub",
            "*" => "fmul",
            "/" => "fdiv",
            _ => "frem",
        };
        let reg = cg.fresh_reg();
        cg.emit(format!("{reg} = {opc} double {}, {}", ld.operand, rd.operand));
        return Ok(Value::new(reg, LType::Double));
    }
    let li = as_i64(cg, l)?;
    let ri = as_i64(cg, r)?;
    let opc = match op {
        "+" => "add",
        "-" => "sub",
        "*" => "mul",
        "/" => "sdiv",
        _ => "srem",
    };
    let reg = cg.fresh_reg();
    cg.emit(format!("{reg} = {opc} i64 {}, {}", li.operand, ri.operand));
    Ok(Value::new(reg, LType::I64))
}

/// String concatenation: `malloc(strlen a + strlen b + 1)` then `strcpy`+`strcat`
/// (libc), promoting a non-string operand through `toString` first.
fn gen_str_concat(cg: &mut Codegen, l: Value, r: Value) -> Result<Value> {
    let ls = to_string_value(cg, l)?;
    let rs = to_string_value(cg, r)?;
    cg.add_extern("declare i64 @strlen(i8*)");
    cg.add_extern("declare i8* @malloc(i64)");
    cg.add_extern("declare i8* @strcpy(i8*, i8*)");
    cg.add_extern("declare i8* @strcat(i8*, i8*)");
    let ll = cg.fresh_reg();
    cg.emit(format!("{ll} = call i64 @strlen(i8* {})", ls.operand));
    let rl = cg.fresh_reg();
    cg.emit(format!("{rl} = call i64 @strlen(i8* {})", rs.operand));
    let sum = cg.fresh_reg();
    cg.emit(format!("{sum} = add i64 {ll}, {rl}"));
    let total = cg.fresh_reg();
    cg.emit(format!("{total} = add i64 {sum}, 1"));
    let buf = cg.fresh_reg();
    cg.emit(format!("{buf} = call i8* @malloc(i64 {total})"));
    let cp = cg.fresh_reg();
    cg.emit(format!("{cp} = call i8* @strcpy(i8* {buf}, i8* {})", ls.operand));
    let ct = cg.fresh_reg();
    cg.emit(format!("{ct} = call i8* @strcat(i8* {buf}, i8* {})", rs.operand));
    Ok(Value::new(buf, LType::Str))
}

fn gen_comparison(cg: &mut Codegen, op: &str, l: Value, r: Value) -> Result<Value> {
    let reg = cg.fresh_reg();
    let is_str = |t: LType| t == LType::Str || t == LType::Ptr;
    if is_str(l.ty) && is_str(r.ty) {
        let cc = match op {
            "==" => "eq",
            "!=" => "ne",
            "<" => "slt",
            "<=" => "sle",
            ">" => "sgt",
            _ => "sge",
        };
        cg.add_extern("declare i32 @strcmp(i8*, i8*)");
        let c = cg.fresh_reg();
        cg.emit(format!(
            "{c} = call i32 @strcmp(i8* {}, i8* {})",
            l.operand, r.operand
        ));
        cg.emit(format!("{reg} = icmp {cc} i32 {c}, 0"));
        return Ok(Value::new(reg, LType::I1));
    }
    if l.ty == LType::Double || r.ty == LType::Double {
        let cc = match op {
            "==" => "oeq",
            "!=" => "one",
            "<" => "olt",
            "<=" => "ole",
            ">" => "ogt",
            _ => "oge",
        };
        let ld = as_double(cg, l)?;
        let rd = as_double(cg, r)?;
        cg.emit(format!(
            "{reg} = fcmp {cc} double {}, {}",
            ld.operand, rd.operand
        ));
        return Ok(Value::new(reg, LType::I1));
    }
    let cc = match op {
        "==" => "eq",
        "!=" => "ne",
        "<" => "slt",
        "<=" => "sle",
        ">" => "sgt",
        _ => "sge",
    };
    let li = as_i64(cg, l)?;
    let ri = as_i64(cg, r)?;
    cg.emit(format!("{reg} = icmp {cc} i64 {}, {}", li.operand, ri.operand));
    Ok(Value::new(reg, LType::I1))
}

fn gen_unary(cg: &mut Codegen, op: &str, operand: &Expr) -> Result<Value> {
    let v = gen_expr(cg, operand)?;
    match op {
        "-" if v.ty == LType::Double => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = fneg double {}", v.operand));
            Ok(Value::new(reg, LType::Double))
        }
        "-" => {
            let i = as_i64(cg, v)?;
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = sub i64 0, {}", i.operand));
            Ok(Value::new(reg, LType::I64))
        }
        "!" | "not" => {
            let b = as_i1(cg, v)?;
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = xor i1 {}, true", b.operand));
            Ok(Value::new(reg, LType::I1))
        }
        other => Err(CodegenError::unsupported(format!("unary operator `{other}`"))),
    }
}

fn gen_call(
    cg: &mut Codegen,
    function: &Expr,
    arguments: &[Expr],
    named: &[NamedArgument],
) -> Result<Value> {
    let Expr::Identifier(name) = function else {
        return Err(CodegenError::unsupported("indirect / higher-order call"));
    };
    match name.as_str() {
        "print" => {
            let arg = first_arg(arguments, named)
                .ok_or_else(|| CodegenError::invalid("print needs one argument"))?;
            let v = gen_expr(cg, arg)?;
            gen_print(cg, v)
        }
        "toString" => {
            let arg = first_arg(arguments, named)
                .ok_or_else(|| CodegenError::invalid("toString needs one argument"))?;
            let v = gen_expr(cg, arg)?;
            to_string_value(cg, v)
        }
        _ => gen_user_call(cg, name, arguments, named),
    }
}

/// A call to a user-defined or runtime function. Parameter types come from
/// inference (so a string/float/bool parameter is passed in its real LLVM
/// type), as does the return type.
fn gen_user_call(
    cg: &mut Codegen,
    name: &str,
    arguments: &[Expr],
    named: &[NamedArgument],
) -> Result<Value> {
    let args = ordered_args(cg, name, arguments, named)?;
    // Coerce each argument to the declared parameter type where known.
    let coerced = match cg.fn_param_ltypes(name) {
        Some(ptys) if ptys.len() == args.len() => args
            .into_iter()
            .zip(ptys)
            .map(|(a, want)| crate::cast::coerce_to(cg, a, want))
            .collect::<Result<Vec<_>>>()?,
        _ => args,
    };
    let typed = coerced
        .iter()
        .map(Value::typed)
        .collect::<Vec<_>>()
        .join(", ");
    let ret = cg.fn_ret_ltype(name).unwrap_or(LType::I64);
    // A name with no user definition is a runtime builtin: declare it so the IR
    // is valid; it links only if the runtime provides the symbol.
    if !cg.fn_params.contains_key(name) {
        let sig = coerced
            .iter()
            .map(|v| v.ty.to_string())
            .collect::<Vec<_>>()
            .join(", ");
        cg.add_extern(format!("declare {ret} @{name}({sig})"));
    }
    let reg = cg.fresh_reg();
    cg.emit(format!("{reg} = call {ret} @{name}({typed})"));
    Ok(Value::new(reg, ret).with_owner(cg.fn_ret_owner(name)))
}

pub(crate) fn ordered_args(
    cg: &mut Codegen,
    name: &str,
    arguments: &[Expr],
    named: &[NamedArgument],
) -> Result<Vec<Value>> {
    if !named.is_empty() {
        if let Some(pnames) = cg.fn_params.get(name).cloned() {
            let mut out = Vec::new();
            for pn in &pnames {
                if let Some(na) = named.iter().find(|a| &a.name == pn) {
                    out.push(gen_expr(cg, &na.value)?);
                }
            }
            if out.len() == named.len() {
                return Ok(out);
            }
        }
        return named.iter().map(|na| gen_expr(cg, &na.value)).collect();
    }
    arguments.iter().map(|a| gen_expr(cg, a)).collect()
}

fn gen_interpolation(cg: &mut Codegen, parts: &[InterpolatedPart]) -> Result<Value> {
    let mut fmt = String::new();
    let mut args: Vec<String> = Vec::new();
    for part in parts {
        match part {
            InterpolatedPart::Text(t) => fmt.push_str(&t.replace('%', "%%")),
            InterpolatedPart::Expr(e) => {
                let v = gen_expr(cg, e)?;
                let s = to_string_value(cg, v)?;
                fmt.push_str("%s");
                args.push(format!("i8* {}", s.operand));
            }
        }
    }
    let fmtv = cg.string_constant(&fmt);
    cg.add_extern("declare i8* @malloc(i64)");
    cg.add_extern("declare i32 @sprintf(i8*, i8*, ...)");
    let buf = cg.fresh_reg();
    cg.emit(format!("{buf} = call i8* @malloc(i64 4096)"));
    let tmp = cg.fresh_reg();
    let extra = if args.is_empty() {
        String::new()
    } else {
        format!(", {}", args.join(", "))
    };
    cg.emit(format!(
        "{tmp} = call i32 (i8*, i8*, ...) @sprintf(i8* {buf}, i8* {}{extra})",
        fmtv.operand
    ));
    Ok(Value::new(buf, LType::Str))
}

fn first_arg<'a>(arguments: &'a [Expr], named: &'a [NamedArgument]) -> Option<&'a Expr> {
    arguments.first().or_else(|| named.first().map(|n| &n.value))
}

pub(crate) fn describe(expr: &Expr) -> String {
    let kind = match expr {
        Expr::List(_) => "list literal",
        Expr::Map(_) => "map literal",
        Expr::Object(_) => "object literal",
        Expr::Pipe { .. } => "pipe expression",
        Expr::FieldAccess { .. } => "field access",
        Expr::MethodCall { .. } => "method call",
        Expr::Index { .. } => "index expression",
        Expr::Lambda { .. } => "lambda",
        Expr::TypeConstructor { .. } => "type constructor",
        Expr::Update { .. } => "record update",
        Expr::Spawn(_) => "spawn",
        Expr::Await(_) => "await",
        Expr::Perform { .. } => "perform",
        Expr::Handler { .. } => "handler",
        _ => "expression",
    };
    kind.to_string()
}
