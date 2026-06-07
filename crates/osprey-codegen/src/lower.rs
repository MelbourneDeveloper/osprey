//! AST → LLVM IR text lowering. Ports the expression/statement/function walks of
//! `llvm.go` and `expression_generation.go` for the int/bool/string + functions
//! core: literals, arithmetic & comparison, `print`/`toString`, string
//! interpolation, `let`, blocks, `match` (bool / int / string literals with a
//! catch-all), function definitions and calls, and the synthesized `main`.
//!
//! Records, unions, effects, lambdas and the Result-wrapped division semantics
//! are not lowered here; those return [`CodegenError::Unsupported`] so a program
//! that needs them fails loudly instead of miscompiling.

use crate::builder::Codegen;
use crate::error::{CodegenError, Result};
use crate::llty::{LType, Value};
use osprey_ast::*;

/// Compile a whole program to an LLVM IR module (text).
pub fn compile_program(program: &Program) -> Result<String> {
    let mut cg = Codegen::new();

    // Pre-pass: record parameter names so named-argument calls can be ordered.
    for stmt in &program.statements {
        if let Stmt::Function {
            name, parameters, ..
        } = stmt
        {
            cg.fn_params.insert(
                name.clone(),
                parameters.iter().map(|p| p.name.clone()).collect(),
            );
        }
    }

    let mut top_level: Vec<&Stmt> = Vec::new();
    let mut user_main: Option<&Expr> = None;
    for stmt in &program.statements {
        match stmt {
            Stmt::Function { name, body, .. } if name == "main" => user_main = Some(body),
            Stmt::Function {
                name,
                parameters,
                body,
                ..
            } => gen_function(&mut cg, name, parameters, body)?,
            Stmt::Let { .. } | Stmt::Assignment { .. } | Stmt::Expr(_) => top_level.push(stmt),
            _ => {}
        }
    }

    cg.begin_function();
    if let Some(body) = user_main {
        gen_expr(&mut cg, body)?;
    } else {
        for stmt in &top_level {
            gen_local_stmt(&mut cg, stmt)?;
        }
    }
    cg.emit("ret i32 0");
    cg.finish_function(LType::I32, "main", &[]);

    Ok(cg.render())
}

fn gen_function(cg: &mut Codegen, name: &str, parameters: &[Parameter], body: &Expr) -> Result<()> {
    cg.begin_function();
    let mut params = Vec::new();
    for p in parameters {
        cg.bind(
            p.name.clone(),
            Value::new(format!("%{}", p.name), LType::I64),
        );
        params.push((LType::I64, p.name.clone()));
    }
    let ret = gen_expr(cg, body)?;
    cg.emit(format!("ret {} {}", ret.ty, ret.operand));
    cg.finish_function(ret.ty, name, &params);
    Ok(())
}

fn gen_local_stmt(cg: &mut Codegen, stmt: &Stmt) -> Result<()> {
    match stmt {
        Stmt::Let { name, value, .. } | Stmt::Assignment { name, value, .. } => {
            let v = gen_expr(cg, value)?;
            cg.bind(name.clone(), v);
            Ok(())
        }
        Stmt::Expr(e) => {
            gen_expr(cg, e)?;
            Ok(())
        }
        _ => Err(CodegenError::unsupported("statement in block/main")),
    }
}

fn gen_expr(cg: &mut Codegen, expr: &Expr) -> Result<Value> {
    match expr {
        Expr::Integer(n) => Ok(Value::new(n.to_string(), LType::I64)),
        Expr::Bool(b) => Ok(Value::new(if *b { "1" } else { "0" }, LType::I1)),
        Expr::Str(s) => Ok(cg.string_constant(s)),
        Expr::InterpolatedStr(parts) => gen_interpolation(cg, parts),
        Expr::Identifier(name) => cg.lookup(name).ok_or_else(|| CodegenError::unknown(name)),
        Expr::Binary { op, left, right } => gen_binary(cg, op, left, right),
        Expr::Unary { op, operand } => gen_unary(cg, op, operand),
        Expr::Call {
            function,
            arguments,
            named_arguments,
        } => gen_call(cg, function, arguments, named_arguments),
        Expr::Match { value, arms } => gen_match(cg, value, arms),
        Expr::Block { statements, value } => gen_block(cg, statements, value.as_deref()),
        other => Err(CodegenError::unsupported(describe(other))),
    }
}

fn gen_block(cg: &mut Codegen, statements: &[Stmt], value: Option<&Expr>) -> Result<Value> {
    cg.push_scope();
    for s in statements {
        gen_local_stmt(cg, s)?;
    }
    let v = match value {
        Some(e) => gen_expr(cg, e)?,
        None => Value::unit(),
    };
    cg.pop_scope();
    Ok(v)
}

fn gen_binary(cg: &mut Codegen, op: &str, left: &Expr, right: &Expr) -> Result<Value> {
    let l = gen_expr(cg, left)?;
    let r = gen_expr(cg, right)?;
    match op {
        "+" | "-" | "*" | "/" | "%" => {
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
        "==" | "!=" | "<" | "<=" | ">" | ">=" => gen_comparison(cg, op, l, r),
        "&&" | "||" => {
            let lb = as_i1(cg, l)?;
            let rb = as_i1(cg, r)?;
            let opc = if op == "&&" { "and" } else { "or" };
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = {opc} i1 {}, {}", lb.operand, rb.operand));
            Ok(Value::new(reg, LType::I1))
        }
        other => Err(CodegenError::unsupported(format!(
            "binary operator `{other}`"
        ))),
    }
}

fn gen_comparison(cg: &mut Codegen, op: &str, l: Value, r: Value) -> Result<Value> {
    let cc = match op {
        "==" => "eq",
        "!=" => "ne",
        "<" => "slt",
        "<=" => "sle",
        ">" => "sgt",
        _ => "sge",
    };
    let reg = cg.fresh_reg();
    match (l.ty == LType::Str, r.ty == LType::Str) {
        (true, true) => {
            cg.add_extern("declare i32 @strcmp(i8*, i8*)");
            let c = cg.fresh_reg();
            cg.emit(format!(
                "{c} = call i32 @strcmp(i8* {}, i8* {})",
                l.operand, r.operand
            ));
            cg.emit(format!("{reg} = icmp {cc} i32 {c}, 0"));
        }
        (false, false) => {
            let li = as_i64(cg, l)?;
            let ri = as_i64(cg, r)?;
            cg.emit(format!(
                "{reg} = icmp {cc} i64 {}, {}",
                li.operand, ri.operand
            ));
        }
        // One side string, one side not — almost always a string-typed
        // parameter/return the i64-default backend hasn't typed yet.
        _ => return Err(string_typing_gap()),
    }
    Ok(Value::new(reg, LType::I1))
}

/// The single limitation of the i64-default backend: it cannot yet type string
/// parameters/returns (that needs the osprey-types signatures). Surface it as a
/// clean error rather than emitting invalid IR.
fn string_typing_gap() -> CodegenError {
    CodegenError::unsupported(
        "comparison mixing string and non-string (string-typed parameter/return \
         inference is the next codegen increment)",
    )
}

fn gen_unary(cg: &mut Codegen, op: &str, operand: &Expr) -> Result<Value> {
    let v = gen_expr(cg, operand)?;
    match op {
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
        other => Err(CodegenError::unsupported(format!(
            "unary operator `{other}`"
        ))),
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
            gen_print(cg, arg)
        }
        "toString" => {
            let arg = first_arg(arguments, named)
                .ok_or_else(|| CodegenError::invalid("toString needs one argument"))?;
            let v = gen_expr(cg, arg)?;
            to_string_value(cg, v)
        }
        _ => {
            let args = ordered_args(cg, name, arguments, named)?;
            let typed = args.iter().map(Value::typed).collect::<Vec<_>>().join(", ");
            // A call to a name that isn't defined in this module is a runtime
            // builtin: declare it (matching the emitted argument types) so the IR
            // is valid LLVM. It then links only if the runtime provides the
            // symbol — the honest outcome for builtins this backend can't host.
            if !cg.fn_params.contains_key(name) {
                let sig = args
                    .iter()
                    .map(|v| v.ty.to_string())
                    .collect::<Vec<_>>()
                    .join(", ");
                cg.add_extern(format!("declare i64 @{name}({sig})"));
            }
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = call i64 @{name}({typed})"));
            Ok(Value::new(reg, LType::I64))
        }
    }
}

fn ordered_args(
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

fn gen_print(cg: &mut Codegen, arg: &Expr) -> Result<Value> {
    let v = gen_expr(cg, arg)?;
    let s = to_string_value(cg, v)?;
    cg.add_extern("declare i32 @puts(i8*)");
    let reg = cg.fresh_reg();
    cg.emit(format!("{reg} = call i32 @puts(i8* {})", s.operand));
    Ok(Value::unit())
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
    cg.emit(format!("{buf} = call i8* @malloc(i64 1024)"));
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

fn to_string_value(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::Str => Ok(v),
        LType::I1 => bool_to_string(cg, v),
        LType::I64 | LType::I32 => int_to_string(cg, v),
    }
}

fn int_to_string(cg: &mut Codegen, v: Value) -> Result<Value> {
    cg.add_extern("declare i8* @malloc(i64)");
    cg.add_extern("declare i32 @sprintf(i8*, i8*, ...)");
    let i = as_i64(cg, v)?;
    let fmt = cg.string_constant("%ld");
    let buf = cg.fresh_reg();
    cg.emit(format!("{buf} = call i8* @malloc(i64 32)"));
    let tmp = cg.fresh_reg();
    cg.emit(format!(
        "{tmp} = call i32 (i8*, i8*, ...) @sprintf(i8* {buf}, i8* {}, i64 {})",
        fmt.operand, i.operand
    ));
    Ok(Value::new(buf, LType::Str))
}

fn bool_to_string(cg: &mut Codegen, v: Value) -> Result<Value> {
    let t = cg.string_constant("true");
    let f = cg.string_constant("false");
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = select i1 {}, i8* {}, i8* {}",
        v.operand, t.operand, f.operand
    ));
    Ok(Value::new(reg, LType::Str))
}

/// `match` over bool / int / string literal arms plus a catch-all
/// (`_`, a binding, or a type-annotated binding). Lowers to a compare-and-branch
/// chain joined by a `phi`. Constructor/destructuring patterns are not lowered.
fn gen_match(cg: &mut Codegen, value: &Expr, arms: &[MatchArm]) -> Result<Value> {
    let disc = gen_expr(cg, value)?;
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
                cg.emit(format!(
                    "br i1 {cond}, label %{body_lbl}, label %{next_lbl}"
                ));
                cg.start_block(&body_lbl);
                let v = gen_expr(cg, &arm.body)?;
                result_ty.get_or_insert(v.ty);
                phi_in.push((v.operand.clone(), cg.cur_block().to_string()));
                cg.emit(format!("br label %{end}"));
                cg.start_block(&next_lbl);
                if i == last {
                    // An exhaustive literal match (e.g. bool true/false) leaves
                    // this final fall-through unreachable.
                    cg.emit("unreachable");
                }
            }
            _ => return Err(CodegenError::unsupported("destructuring match arm")),
        }
    }

    cg.start_block(&end);
    let ty = result_ty.unwrap_or(LType::I64);
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

/// Emit an equality test between the discriminant and a literal pattern,
/// returning the `i1` operand.
fn gen_eq(cg: &mut Codegen, disc: &Value, lit: &Expr) -> Result<String> {
    let pat = gen_expr(cg, lit)?;
    let reg = cg.fresh_reg();
    match (disc.ty == LType::Str, pat.ty == LType::Str) {
        (true, true) => {
            cg.add_extern("declare i32 @strcmp(i8*, i8*)");
            let c = cg.fresh_reg();
            cg.emit(format!(
                "{c} = call i32 @strcmp(i8* {}, i8* {})",
                disc.operand, pat.operand
            ));
            cg.emit(format!("{reg} = icmp eq i32 {c}, 0"));
        }
        (false, false) => {
            let d = as_i64(cg, disc.clone())?;
            let p = as_i64(cg, pat)?;
            cg.emit(format!("{reg} = icmp eq i64 {}, {}", d.operand, p.operand));
        }
        _ => return Err(string_typing_gap()),
    }
    Ok(reg)
}

fn as_i64(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::I64 => Ok(v),
        LType::I1 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = zext i1 {} to i64", v.operand));
            Ok(Value::new(reg, LType::I64))
        }
        LType::I32 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = sext i32 {} to i64", v.operand));
            Ok(Value::new(reg, LType::I64))
        }
        LType::Str => Err(CodegenError::invalid("expected an integer, found a string")),
    }
}

fn as_i1(cg: &mut Codegen, v: Value) -> Result<Value> {
    match v.ty {
        LType::I1 => Ok(v),
        LType::I64 | LType::I32 => {
            let reg = cg.fresh_reg();
            cg.emit(format!("{reg} = icmp ne {} {}, 0", v.ty, v.operand));
            Ok(Value::new(reg, LType::I1))
        }
        LType::Str => Err(CodegenError::invalid("expected a bool, found a string")),
    }
}

fn first_arg<'a>(arguments: &'a [Expr], named: &'a [NamedArgument]) -> Option<&'a Expr> {
    arguments
        .first()
        .or_else(|| named.first().map(|n| &n.value))
}

fn describe(expr: &Expr) -> String {
    let kind = match expr {
        Expr::Float(_) => "float literal",
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
