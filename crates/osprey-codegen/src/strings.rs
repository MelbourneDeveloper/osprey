//! String builtins — thin wrappers over the C string runtime declared in
//! `runtime/string_runtime.h` (linked from `libfiber_runtime`). Ports
//! `string_functions.go` + the `length`/`contains`/`substring`/`parseInt`/
//! `join` generators in `core_functions.go`. Total operations return their bare
//! value; fallible ones return the uniform `{ value, i8 }*` Result block.
//! Implements [BUILTIN-STRING-*].

use crate::builder::Codegen;
use crate::cast::coerce_to;
use crate::error::{CodegenError, Result};
use crate::expr::gen_expr;
use crate::llty::{LType, Value};
use crate::result::make_result;
use osprey_ast::{Expr, NamedArgument};

/// Dispatch a string builtin by name, or `None` if `name` is not one.
pub(crate) fn gen(
    cg: &mut Codegen,
    name: &str,
    args: &[Expr],
    named: &[NamedArgument],
) -> Result<Option<Value>> {
    let v = match name {
        "length" => unary_i64(cg, "strlen", args, named)?,
        "isEmpty" => bool_from_i64(cg, "osp_string_is_empty", &[(0, LType::Str)], args, named)?,
        "contains" => contains(cg, args, named)?,
        "startsWith" => {
            bool_from_i64(cg, "osp_string_starts_with", &[(0, LType::Str), (1, LType::Str)], args, named)?
        }
        "endsWith" => {
            bool_from_i64(cg, "osp_string_ends_with", &[(0, LType::Str), (1, LType::Str)], args, named)?
        }
        "toUpperCase" => unary_str(cg, "osp_string_to_upper", args, named)?,
        "toLowerCase" => unary_str(cg, "osp_string_to_lower", args, named)?,
        "trim" => unary_str(cg, "osp_string_trim", args, named)?,
        "trimStart" => unary_str(cg, "osp_string_trim_start", args, named)?,
        "trimEnd" => unary_str(cg, "osp_string_trim_end", args, named)?,
        "reverse" => unary_str(cg, "osp_string_reverse", args, named)?,
        "take" => str_int_str(cg, "osp_string_take", args, named)?,
        "drop" => str_int_str(cg, "osp_string_drop", args, named)?,
        "indexOf" => index_of(cg, args, named)?,
        "substring" => substring(cg, args, named)?,
        "replace" => nullable_str(cg, "osp_string_replace", &[LType::Str, LType::Str, LType::Str], args, named)?,
        "repeat" => nullable_str(cg, "osp_string_repeat", &[LType::Str, LType::I64], args, named)?,
        "padStart" => {
            nullable_str(cg, "osp_string_pad_start", &[LType::Str, LType::I64, LType::Str], args, named)?
        }
        "padEnd" => {
            nullable_str(cg, "osp_string_pad_end", &[LType::Str, LType::I64, LType::Str], args, named)?
        }
        "parseInt" => parse_strict(cg, "osp_parse_int_strict", LType::I64, args, named)?,
        "parseFloat" => parse_strict(cg, "osp_parse_float_strict", LType::Double, args, named)?,
        "join" => join(cg, args, named)?,
        "lines" => string_list(cg, "osp_string_lines", args)?,
        "words" => string_list(cg, "osp_string_words", args)?,
        "split" => split(cg, args)?,
        _ => return Ok(None),
    };
    Ok(Some(v))
}

/// The `i`-th positional argument, evaluated and coerced to `want`.
fn arg(cg: &mut Codegen, args: &[Expr], i: usize, want: LType) -> Result<Value> {
    let e = args
        .get(i)
        .ok_or_else(|| CodegenError::invalid("string builtin: missing argument"))?;
    let v = gen_expr(cg, e)?;
    coerce_to(cg, v, want)
}

/// Declare a runtime function once (idempotent via the extern set).
fn declare(cg: &mut Codegen, cname: &str, ret: &str, params: &str) {
    cg.add_extern(format!("declare {ret} @{cname}({params})"));
}

/// `f(s: string) -> int`.
fn unary_i64(
    cg: &mut Codegen,
    cname: &str,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    declare(cg, cname, "i64", "i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i64 @{cname}(i8* {})", s.operand));
    Ok(Value::new(r, LType::I64))
}

/// `f(s: string) -> string`.
fn unary_str(
    cg: &mut Codegen,
    cname: &str,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    declare(cg, cname, "i8*", "i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @{cname}(i8* {})", s.operand));
    Ok(Value::new(r, LType::Str))
}

/// `f(s: string, n: int) -> string`.
fn str_int_str(
    cg: &mut Codegen,
    cname: &str,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    let n = arg(cg, args, 1, LType::I64)?;
    declare(cg, cname, "i8*", "i8*, i64");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @{cname}(i8* {}, i64 {})", s.operand, n.operand));
    Ok(Value::new(r, LType::Str))
}

/// A runtime predicate returning `i64` truthiness, narrowed to `i1`. `sig` lists
/// each argument index and the LLVM type it travels as.
fn bool_from_i64(
    cg: &mut Codegen,
    cname: &str,
    sig: &[(usize, LType)],
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let mut typed = Vec::new();
    let mut params = Vec::new();
    for (i, ty) in sig {
        let v = arg(cg, args, *i, *ty)?;
        typed.push(v.typed());
        params.push(ty.to_string());
    }
    declare(cg, cname, "i64", &params.join(", "));
    let raw = cg.fresh_reg();
    cg.emit(format!("{raw} = call i64 @{cname}({})", typed.join(", ")));
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = icmp ne i64 {raw}, 0"));
    Ok(Value::new(r, LType::I1))
}

/// `contains(s, needle) -> bool` via libc `strstr` (non-NULL ⇒ found).
fn contains(cg: &mut Codegen, args: &[Expr], _named: &[NamedArgument]) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    let needle = arg(cg, args, 1, LType::Str)?;
    cg.add_extern("declare i8* @strstr(i8*, i8*)");
    let hit = cg.fresh_reg();
    cg.emit(format!("{hit} = call i8* @strstr(i8* {}, i8* {})", s.operand, needle.operand));
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = icmp ne i8* {hit}, null"));
    Ok(Value::new(r, LType::I1))
}

/// `indexOf(s, needle) -> Result<int, _>` (`-1` ⇒ Error).
fn index_of(cg: &mut Codegen, args: &[Expr], _named: &[NamedArgument]) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    let needle = arg(cg, args, 1, LType::Str)?;
    declare(cg, "osp_string_index_of", "i64", "i8*, i8*");
    let idx = cg.fresh_reg();
    cg.emit(format!("{idx} = call i64 @osp_string_index_of(i8* {}, i8* {})", s.operand, needle.operand));
    let iserr = cg.fresh_reg();
    cg.emit(format!("{iserr} = icmp slt i64 {idx}, 0"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {iserr}, i8 1, i8 0"));
    let val = cg.fresh_reg();
    cg.emit(format!("{val} = select i1 {iserr}, i64 0, i64 {idx}"));
    make_result(cg, Value::new(val, LType::I64), LType::I64, &disc)
}

/// `substring(s, start, end) -> Result<string, _>` (NULL ⇒ Error).
fn substring(cg: &mut Codegen, args: &[Expr], _named: &[NamedArgument]) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    let start = arg(cg, args, 1, LType::I64)?;
    let end = arg(cg, args, 2, LType::I64)?;
    declare(cg, "osp_string_substring", "i8*", "i8*, i64, i64");
    let ptr = cg.fresh_reg();
    cg.emit(format!(
        "{ptr} = call i8* @osp_string_substring(i8* {}, i64 {}, i64 {})",
        s.operand, start.operand, end.operand
    ));
    result_from_nullable(cg, &ptr)
}

/// A fallible string transform returning a runtime `char*` that is NULL on
/// failure, wrapped into `Result<string, _>`. `argtys` lists each argument's
/// LLVM type in order.
fn nullable_str(
    cg: &mut Codegen,
    cname: &str,
    argtys: &[LType],
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let mut typed = Vec::new();
    let mut params = Vec::new();
    for (i, ty) in argtys.iter().enumerate() {
        let v = arg(cg, args, i, *ty)?;
        typed.push(v.typed());
        params.push(ty.to_string());
    }
    declare(cg, cname, "i8*", &params.join(", "));
    let ptr = cg.fresh_reg();
    cg.emit(format!("{ptr} = call i8* @{cname}({})", typed.join(", ")));
    result_from_nullable(cg, &ptr)
}

/// `Result<string, _>` from a possibly-NULL `char*` (`ptr` is an `i8*` register):
/// NULL ⇒ Error, else Success.
fn result_from_nullable(cg: &mut Codegen, ptr: &str) -> Result<Value> {
    let iserr = cg.fresh_reg();
    cg.emit(format!("{iserr} = icmp eq i8* {ptr}, null"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {iserr}, i8 1, i8 0"));
    make_result(cg, Value::new(ptr, LType::Str), LType::Str, &disc)
}

/// `parseInt`/`parseFloat`: strict parse writing through an out-slot, returning
/// `0` on success. `inner` is the parsed value's LLVM type.
fn parse_strict(
    cg: &mut Codegen,
    cname: &str,
    inner: LType,
    args: &[Expr],
    _named: &[NamedArgument],
) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    declare(cg, cname, "i64", &format!("i8*, {inner}*"));
    let slot = cg.fresh_reg();
    cg.emit(format!("{slot} = alloca {inner}"));
    let zero = if inner == LType::Double { "0.0" } else { "0" };
    cg.emit(format!("store {inner} {zero}, {inner}* {slot}"));
    let rc = cg.fresh_reg();
    cg.emit(format!("{rc} = call i64 @{cname}(i8* {}, {inner}* {slot})", s.operand));
    let parsed = cg.fresh_reg();
    cg.emit(format!("{parsed} = load {inner}, {inner}* {slot}"));
    let iserr = cg.fresh_reg();
    cg.emit(format!("{iserr} = icmp ne i64 {rc}, 0"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {iserr}, i8 1, i8 0"));
    make_result(cg, Value::new(parsed, inner), inner, &disc)
}

/// `lines`/`words`: split a string into a `List<string>`. The C `osp_string_list`
/// shares its leading `i64 length` with the runtime list, so list builtins read
/// it directly; tag the handle as a list.
fn string_list(cg: &mut Codegen, cname: &str, args: &[Expr]) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    declare(cg, cname, "i8*", "i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @{cname}(i8* {})", s.operand));
    Ok(Value::handle(r, crate::collections::LIST_OWNER))
}

/// `split(s, sep) -> Result<List<string>, _>` (NULL ⇒ Error, e.g. empty sep).
fn split(cg: &mut Codegen, args: &[Expr]) -> Result<Value> {
    let s = arg(cg, args, 0, LType::Str)?;
    let sep = arg(cg, args, 1, LType::Str)?;
    declare(cg, "osp_string_split", "i8*", "i8*, i8*");
    let ptr = cg.fresh_reg();
    cg.emit(format!("{ptr} = call i8* @osp_string_split(i8* {}, i8* {})", s.operand, sep.operand));
    let iserr = cg.fresh_reg();
    cg.emit(format!("{iserr} = icmp eq i8* {ptr}, null"));
    let disc = cg.fresh_reg();
    cg.emit(format!("{disc} = select i1 {iserr}, i8 1, i8 0"));
    make_result(cg, Value::handle(ptr, crate::collections::LIST_OWNER), LType::Ptr, &disc)
}

/// `join(list: List<string>, separator: string) -> string`.
fn join(cg: &mut Codegen, args: &[Expr], _named: &[NamedArgument]) -> Result<Value> {
    let list = arg(cg, args, 0, LType::Ptr)?;
    let sep = arg(cg, args, 1, LType::Str)?;
    declare(cg, "osp_string_join", "i8*", "i8*, i8*");
    let r = cg.fresh_reg();
    cg.emit(format!("{r} = call i8* @osp_string_join(i8* {}, i8* {})", list.operand, sep.operand));
    Ok(Value::new(r, LType::Str))
}
