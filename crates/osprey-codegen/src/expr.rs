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
use osprey_ast::{Expr, InterpolatedPart, NamedArgument, Parameter, Stmt};

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
            // A bare top-level function name used as a value is a callback (passed
            // to `spawnProcess`/`httpListen`): take its address as an `i8*`.
            None if cg.fn_params.contains_key(name) => Ok(fn_pointer(cg, name)),
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
        Expr::Object(fields) => crate::aggregate::gen_object(cg, fields),
        Expr::List(elements) => crate::listlit::gen_list(cg, elements),
        Expr::Map(entries) => crate::collections::gen_map_literal(cg, entries),
        Expr::Index { target, index } => crate::listlit::gen_index(cg, target, index),
        Expr::Spawn(e) => crate::fiber::gen_spawn(cg, e),
        Expr::Await(e) => crate::fiber::gen_await(cg, e),
        Expr::Yield(e) => crate::fiber::gen_yield(cg, e.as_deref()),
        Expr::Send { channel, value } => crate::fiber::gen_send(cg, channel, value),
        Expr::Recv(e) => crate::fiber::gen_recv(cg, e),
        Expr::Select { arms } => crate::fiber::gen_select(cg, arms),
        Expr::Perform {
            effect,
            operation,
            arguments,
            ..
        } => crate::effects::gen_perform(cg, effect, operation, arguments),
        Expr::Handler { effect, arms, body } => crate::effects::gen_handler(cg, effect, arms, body),
        other => Err(CodegenError::unsupported(describe(other))),
    }
}

/// A bare top-level function name used as a runtime value (a callback handed to
/// `spawnProcess`/`httpListen`): emit its address bitcast to `i8*`. The source
/// type of the bitcast is the function's exact emitted signature — built the
/// same way `gen_function`/`coerce_return` spelled its `define` — so the cast is
/// well-typed; the C runtime calls back through its own function-pointer cast.
/// Mirrors the handler-pointer bitcast in `effects::gen_perform`.
fn fn_pointer(cg: &mut Codegen, name: &str) -> Value {
    let fty = fn_ptr_type(cg, name);
    let reg = cg.emit_reg(format!("bitcast {fty} @{name} to i8*"));
    Value::new(reg, LType::Ptr)
}

/// The LLVM function-pointer type spelling for a top-level function, e.g.
/// `i64 (i64, i64, i8*)*` — return type (a `{ T, i8 }*` Result block, or the
/// inferred scalar; `Unit` rides as `i64`) then its parameter type list.
fn fn_ptr_type(cg: &Codegen, name: &str) -> String {
    let params = cg
        .fn_param_ltypes(name)
        .unwrap_or_default()
        .iter()
        .map(|t| t.as_str())
        .collect::<Vec<_>>()
        .join(", ");
    let ret = match cg.fn_ret_result_inner(name) {
        Some(inner) => format!("{{ {inner}, i8 }}*"),
        None => cg.fn_ret_ltype(name).unwrap_or(LType::I64).to_string(),
    };
    format!("{ret} ({params})*")
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
    // A block does NOT open a new scope: the Go backend uses a flat per-function
    // symbol table, so a nested `let` rebinds (and leaks) the name in the
    // enclosing scope. Replicating that keeps shadowing byte-compatible — e.g.
    // block_statements' inner `let outer` is visible to the outer `outer + inner`.
    for s in statements {
        crate::lower::gen_local_stmt(cg, s)?;
    }
    match value {
        Some(e) => gen_expr(cg, e),
        None => Ok(Value::unit()),
    }
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

    // Operands auto-unwrap a Result to its success payload before arithmetic or
    // comparison (mirrors Go's `unwrapIfResult`).
    let l = gen_expr(cg, left)?;
    let l = crate::result::unwrap(cg, l);
    let r = gen_expr(cg, right)?;
    let r = crate::result::unwrap(cg, r);
    match op {
        "+" | "-" | "*" | "/" | "%" => gen_arith(cg, op, l, r),
        "==" | "!=" | "<" | "<=" | ">" | ">=" => gen_comparison(cg, op, l, r),
        other => Err(CodegenError::unsupported(format!(
            "binary operator `{other}`"
        ))),
    }
}

/// Arithmetic. Float if either operand is a float (the other is promoted),
/// otherwise integer. Division ALWAYS returns float (the Osprey spec — see
/// `generateDivisionWithZeroCheck`); modulo stays integer. The Result<…,
/// `MathError`> wrapper the type system tracks is auto-unwrapped at value sites.
fn gen_arith(cg: &mut Codegen, op: &str, l: Value, r: Value) -> Result<Value> {
    // `+` on list handles is concatenation (`a + b` ≡ `listConcat(a, b)`); on
    // map handles it is a right-biased merge (`a + b` ≡ `mapMerge(a, b)`).
    let is_list = |v: &Value| v.osp_ty.as_deref() == Some(crate::collections::LIST_OWNER);
    let is_map = |v: &Value| v.osp_ty.as_deref() == Some(crate::collections::MAP_OWNER);
    if op == "+" && (is_list(&l) || is_list(&r)) {
        return Ok(crate::collections::concat_handles(cg, &l, &r));
    }
    if op == "+" && (is_map(&l) || is_map(&r)) {
        return Ok(crate::collections::merge_handles(cg, &l, &r));
    }
    // `+` with a string operand is concatenation (libc strlen/strcpy/strcat),
    // matching `generateStringConcatenation`.
    if op == "+" && (l.ty == LType::Str || r.ty == LType::Str) {
        return gen_str_concat(cg, l, r);
    }
    // Numeric arithmetic infers `Result<…, MathError>` in Go (`createSuccessResult`
    // / `createSuccessResultFloat`): `/` always float, the rest follow operand
    // type. The Success wrapper auto-unwraps at value sites (interpolation,
    // comparison, args), but `toString`/`print` show it as `Success(n)`.
    if op == "/" {
        return gen_division(cg, l, r);
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
        let reg = cg.emit_reg(format!("{opc} double {}, {}", ld.operand, rd.operand));
        return crate::result::make_ok(cg, Value::new(reg, LType::Double), LType::Double);
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
    let reg = cg.emit_reg(format!("{opc} i64 {}, {}", li.operand, ri.operand));
    crate::result::make_ok(cg, Value::new(reg, LType::I64), LType::I64)
}

/// Division — always float, with a runtime divide-by-zero check (Go's
/// `generateDivisionWithZeroCheck`): a zero divisor yields `Error`
/// (`Result<float, MathError>` disc 1), else `Success(quotient)`.
fn gen_division(cg: &mut Codegen, l: Value, r: Value) -> Result<Value> {
    use crate::result::make_result;
    let ld = as_double(cg, l)?;
    let rd = as_double(cg, r)?;
    let isz = cg.fresh_reg();
    cg.emit(format!("{isz} = fcmp oeq double {}, 0.0", rd.operand));
    let zero_bb = cg.fresh_label();
    let nonzero_bb = cg.fresh_label();
    let end = cg.fresh_label();
    cg.emit(format!(
        "br i1 {isz}, label %{zero_bb}, label %{nonzero_bb}"
    ));

    cg.start_block(&nonzero_bb);
    let q = cg.fresh_reg();
    cg.emit(format!("{q} = fdiv double {}, {}", ld.operand, rd.operand));
    let ok = make_result(cg, Value::new(q, LType::Double), LType::Double, "0")?;
    let okb = cg.snapshot_to(&end);

    cg.start_block(&zero_bb);
    let err = make_result(cg, Value::new("0.0", LType::Double), LType::Double, "1")?;
    let errb = cg.snapshot_to(&end);

    cg.start_block(&end);
    let reg = cg.fresh_reg();
    cg.emit(format!(
        "{reg} = phi {{ double, i8 }}* [ {}, %{okb} ], [ {}, %{errb} ]",
        ok.operand, err.operand
    ));
    Ok(Value::result(reg, LType::Double))
}

/// String concatenation: `malloc(strlen a + strlen b + 1)` then `strcpy`+`strcat`
/// (libc), promoting a non-string operand through `toString` first.
fn gen_str_concat(cg: &mut Codegen, l: Value, r: Value) -> Result<Value> {
    let ls = to_string_value(cg, l)?;
    let rs = to_string_value(cg, r)?;
    let ll = cg.call("i64", "strlen", "i8*", &[&ls.operand]);
    let rl = cg.call("i64", "strlen", "i8*", &[&rs.operand]);
    let sum = cg.emit_reg(format!("add i64 {ll}, {rl}"));
    let total = cg.emit_reg(format!("add i64 {sum}, 1"));
    let buf = cg.call("i8*", "malloc", "i64", &[&total]);
    let _ = cg.call("i8*", "strcpy", "i8*, i8*", &[&buf, &ls.operand]);
    let _ = cg.call("i8*", "strcat", "i8*, i8*", &[&buf, &rs.operand]);
    Ok(Value::new(buf, LType::Str))
}

/// The LLVM condition code for a comparison `op`. `float` picks the ordered
/// `fcmp` codes (`oeq`, `olt`, …); otherwise the signed-integer / `icmp` codes
/// (`eq`, `slt`, …) — also used on a `strcmp` result.
fn cmp_code(op: &str, float: bool) -> &'static str {
    match (op, float) {
        ("==", false) => "eq",
        ("!=", false) => "ne",
        ("<", false) => "slt",
        ("<=", false) => "sle",
        (">", false) => "sgt",
        (_, false) => "sge",
        ("==", true) => "oeq",
        ("!=", true) => "one",
        ("<", true) => "olt",
        ("<=", true) => "ole",
        (">", true) => "ogt",
        (_, true) => "oge",
    }
}

fn gen_comparison(cg: &mut Codegen, op: &str, l: Value, r: Value) -> Result<Value> {
    let reg = cg.fresh_reg();
    let is_str = |t: LType| t == LType::Str || t == LType::Ptr;
    if is_str(l.ty) && is_str(r.ty) {
        let c = cg.call("i32", "strcmp", "i8*, i8*", &[&l.operand, &r.operand]);
        cg.emit(format!("{reg} = icmp {} i32 {c}, 0", cmp_code(op, false)));
        return Ok(Value::new(reg, LType::I1));
    }
    if l.ty == LType::Double || r.ty == LType::Double {
        let ld = as_double(cg, l)?;
        let rd = as_double(cg, r)?;
        cg.emit(format!(
            "{reg} = fcmp {} double {}, {}",
            cmp_code(op, true),
            ld.operand,
            rd.operand
        ));
        return Ok(Value::new(reg, LType::I1));
    }
    let cc = cmp_code(op, false);
    let li = as_i64(cg, l)?;
    let ri = as_i64(cg, r)?;
    cg.emit(format!(
        "{reg} = icmp {cc} i64 {}, {}",
        li.operand, ri.operand
    ));
    Ok(Value::new(reg, LType::I1))
}

fn gen_unary(cg: &mut Codegen, op: &str, operand: &Expr) -> Result<Value> {
    let v = gen_expr(cg, operand)?;
    match op {
        "-" if v.ty == LType::Double => Ok(Value::new(
            cg.emit_reg(format!("fneg double {}", v.operand)),
            LType::Double,
        )),
        "-" => {
            let i = as_i64(cg, v)?;
            Ok(Value::new(
                cg.emit_reg(format!("sub i64 0, {}", i.operand)),
                LType::I64,
            ))
        }
        "!" | "not" => {
            let b = as_i1(cg, v)?;
            Ok(Value::new(
                cg.emit_reg(format!("xor i1 {}, true", b.operand)),
                LType::I1,
            ))
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
    // A directly-applied lambda (`x |> fn(y) => …`, `(fn(y) => …)(x)`) is
    // beta-reduced inline.
    if let Expr::Lambda {
        parameters, body, ..
    } = function
    {
        return apply_lambda(cg, parameters, body, arguments);
    }
    let Expr::Identifier(ident) = function else {
        return Err(CodegenError::unsupported("indirect / higher-order call"));
    };
    // A function-valued parameter (bound while inlining a generic function)
    // redirects to its real callee, so `f(x)` becomes `toString(x)` / `addOne(x)`.
    let name: String = cg
        .call_aliases
        .get(ident)
        .cloned()
        .unwrap_or_else(|| ident.clone());
    let name = name.as_str();
    // A let-bound lambda is inlined at its call site (no closures lowered).
    if let Some((params, body)) = cg.lambdas.get(name).cloned() {
        return apply_lambda(cg, &params, &body, arguments);
    }
    // A call through a function-typed local (`f(x)` where `f: (int) -> int` is a
    // higher-order parameter) is an indirect call through its `i8*` handle.
    if let Some(v) = crate::genfn::try_indirect(cg, name, arguments, named)? {
        return Ok(v);
    }
    match name {
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
        // Runtime builtins take precedence over a same-named user function: the
        // names below are reserved. Each dispatcher returns `None` when the name
        // is not its builtin, so the chain falls through to a user call.
        _ => {
            if let Some(v) = crate::strings::gen(cg, name, arguments, named)? {
                return Ok(v);
            }
            if let Some(v) = crate::collections::gen(cg, name, arguments, named)? {
                return Ok(v);
            }
            if let Some(v) = crate::iter::gen(cg, name, arguments, named)? {
                return Ok(v);
            }
            if let Some(v) = crate::fiber::gen_builtin(cg, name, arguments)? {
                return Ok(v);
            }
            if let Some(v) = crate::extern_call::gen(cg, name, arguments, named)? {
                return Ok(v);
            }
            // A generic user function is specialised by inlining its body with
            // the concrete argument types at this call site.
            if let Some(v) = crate::genfn::try_inline(cg, name, arguments, named)? {
                return Ok(v);
            }
            gen_user_call(cg, name, arguments, named)
        }
    }
}

/// Beta-reduce a lambda at its application site: bind each parameter to its
/// argument, lower the body in a fresh scope, then unwrap a `Result` return.
/// A lambda's inferred return type is its body's success payload, so applying
/// `fn(x) => x * 2` yields a plain `int` — matching how Go unwraps a call's
/// `Result<…, MathError>` at a non-Result return boundary (`maybeWrapInResult`).
fn apply_lambda(
    cg: &mut Codegen,
    parameters: &[Parameter],
    body: &Expr,
    arguments: &[Expr],
) -> Result<Value> {
    cg.push_scope();
    for (p, a) in parameters.iter().zip(arguments) {
        let av = gen_expr(cg, a)?;
        cg.bind(p.name.clone(), av);
    }
    let v = gen_expr(cg, body);
    cg.pop_scope();
    Ok(crate::result::unwrap(cg, v?))
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
    call_with_values(cg, name, args)
}

/// Call `name` with already-evaluated argument values — the shared tail of
/// `gen_user_call` and the iterator callbacks. Coerces each argument to the
/// inferred parameter type, declares unknown (runtime) callees, and tags a
/// `Result`-returning callee's value.
pub(crate) fn call_with_values(cg: &mut Codegen, name: &str, args: Vec<Value>) -> Result<Value> {
    // `print` as a first-class callback maps to the print intrinsic.
    if name == "print" {
        let v = args.into_iter().next().unwrap_or_else(Value::unit);
        return gen_print(cg, v);
    }
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
    // A function declared `-> Result<T, E>` hands back a `{ T, i8 }*` block.
    if let Some(inner) = cg.fn_ret_result_inner(name) {
        let rty = format!("{{ {inner}, i8 }}*");
        let reg = emit_user_call(cg, name, &rty, &coerced, &typed);
        return Ok(Value::result(reg, inner));
    }
    let ret = cg.fn_ret_ltype(name).unwrap_or(LType::I64);
    let reg = emit_user_call(cg, name, ret.as_str(), &coerced, &typed);
    Ok(Value::new(reg, ret).with_owner(cg.fn_ret_owner(name)))
}

/// Emit a call to `name` returning LLVM type `rty`. A name with no user
/// definition is a runtime builtin, so synthesize its `declare` (param types
/// from `coerced`) — the IR stays valid and links only if the symbol exists.
fn emit_user_call(
    cg: &mut Codegen,
    name: &str,
    rty: &str,
    coerced: &[Value],
    typed: &str,
) -> String {
    if !cg.fn_params.contains_key(name) {
        let sig = coerced
            .iter()
            .map(Value::llvm_ty)
            .collect::<Vec<_>>()
            .join(", ");
        cg.add_extern(format!("declare {rty} @{name}({sig})"));
    }
    cg.emit_reg(format!("call {rty} @{name}({typed})"))
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
                // `${expr}` unwraps a Result to its payload before formatting
                // (Go's `generateInterpolatedString` calls `unwrapIfResult`), so
                // `${21 * 2}` prints `42`, not `Success(42)`.
                let v = gen_expr(cg, e)?;
                let v = crate::result::unwrap(cg, v);
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
    arguments
        .first()
        .or_else(|| named.first().map(|n| &n.value))
}

/// A call's argument expressions in call order — positional, or named in written
/// order — for callees with a fixed parameter list (runtime builtins, indirect
/// calls) that bind by position rather than reordering by parameter name.
pub(crate) fn arg_exprs<'a>(args: &'a [Expr], named: &'a [NamedArgument]) -> Vec<&'a Expr> {
    if named.is_empty() {
        args.iter().collect()
    } else {
        named.iter().map(|n| &n.value).collect()
    }
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
