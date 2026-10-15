//! One kernel body from its AST to the device dialect [GPU-BACKEND-DEVICE].
//!
//! Every expression becomes a pure C expression over `long` (int), `ulong`
//! (a float's bits) and `bool`; a construct that needs statements — a `let`,
//! a `match`, a `fold` over a `range`, a bounds-checked `gpuGet` — hoists them
//! into the enclosing block ahead of a temporary. Helper functions are inlined
//! at the AST level, so the kernel is one self-contained entry point.
//! Integer arithmetic goes through the library's checked helpers, which record
//! a fault the runtime aborts on ([GPU-DEVICE-FAULT]); float arithmetic through
//! the binary64 library. Anything else declines with the reason.

use super::{DeviceBody, Ty, Uniform};
use crate::builder::Codegen;
use osprey_ast::{Expr, MatchArm, NamedArgument, Parameter, Pattern, Stmt};
use std::collections::HashMap;

type Emitted = Result<(Ty, String), String>;

/// A name in scope: a scalar held in a C variable, or a buffer uniform with
/// its pointer and length variables.
#[derive(Clone)]
enum Local {
    Scalar(Ty, String),
    Buffer(String, String),
}

struct Emitter<'a> {
    cg: &'a Codegen,
    /// The statement block being written, innermost last.
    blocks: Vec<Vec<String>>,
    temps: usize,
    scopes: Vec<HashMap<String, Local>>,
    depth: usize,
}

/// Helper inlining deeper than this is treated as recursion.
const MAX_INLINE_DEPTH: usize = 24;
const RESULT_OK: &str = "Success";
const RESULT_ERR: &str = "Error";

pub(super) fn is_literal(expression: &Expr) -> bool {
    match expression {
        Expr::Integer(_) | Expr::Float(_) | Expr::Bool(_) => true,
        Expr::Unary { op, operand } => op == "-" && is_literal(operand),
        _ => false,
    }
}

/// Compile `body` with the kernel's parameters bound to `p0`, `p1`, … and
/// its captures to `u0`, `u1`, ….
pub(super) fn body(
    cg: &Codegen,
    parameters: &[Parameter],
    body: &Expr,
    params: &[Ty],
    captures: &[(String, Uniform)],
) -> Result<DeviceBody, String> {
    if parameters.len() != params.len() {
        return Err("the kernel's arity disagrees with the combinator's slots".into());
    }
    let mut emitter = Emitter {
        cg,
        blocks: vec![Vec::new()],
        temps: 0,
        scopes: vec![HashMap::new()],
        depth: 0,
    };
    for (i, (parameter, ty)) in parameters.iter().zip(params).enumerate() {
        emitter.bind(&parameter.name, Local::Scalar(*ty, format!("p{i}")));
    }
    for (i, (name, uniform)) in captures.iter().enumerate() {
        emitter.bind(name, uniform_local(i, uniform));
    }
    let (result_ty, result) = emitter.expression(body)?;
    Ok(DeviceBody {
        lines: emitter.blocks.pop().unwrap_or_default(),
        result,
        result_ty,
        params: params.to_vec(),
        uniforms: captures
            .iter()
            .map(|(_, uniform)| uniform.clone())
            .collect(),
    })
}

fn uniform_local(i: usize, uniform: &Uniform) -> Local {
    match uniform {
        Uniform::Scalar(ty) => Local::Scalar(*ty, format!("u{i}_v")),
        Uniform::Buffer => Local::Buffer(format!("u{i}"), format!("u{i}_len")),
    }
}

fn int_literal(n: i64) -> String {
    match n {
        i64::MIN => "LONG_MIN_WORD".to_string(),
        n if n < 0 => format!("(-(long){})", -n),
        n => format!("((long){n})"),
    }
}

fn float_literal(f: f64) -> String {
    format!("0x{:016x}ul", f.to_bits())
}

fn describe(expression: &Expr) -> &'static str {
    match expression {
        Expr::Lambda { .. } => "a function value",
        Expr::Str(_) | Expr::InterpolatedStr(_) => "a string",
        Expr::List(..) => "a list",
        Expr::FieldAccess { .. } => "a field access",
        Expr::TypeConstructor { .. } | Expr::Object(_) | Expr::Update { .. } => "a record",
        Expr::Perform { .. } => "an effect operation",
        _ => "this expression",
    }
}

/// The callee behind an explicit type application.
pub(super) fn callee(function: &Expr) -> &Expr {
    match function {
        Expr::TypeApply { function, .. } => callee(function),
        other => other,
    }
}

/// `gpuGet(buffer, index) ?: fallback`, as the parser desugars it.
fn elvis_gather<'a>(
    value: &'a Expr,
    arms: &'a [MatchArm],
) -> Option<(&'a Expr, &'a Expr, &'a Expr)> {
    let Expr::Call {
        function,
        arguments,
        ..
    } = value
    else {
        return None;
    };
    let [buffer, index] = arguments.as_slice() else {
        return None;
    };
    let [ok, err] = arms else {
        return None;
    };
    let is_get = matches!(callee(function), Expr::Identifier(name) if name == "gpuGet");
    let is_ok = matches!((&ok.pattern, &ok.body), (Pattern::Constructor { name, fields, .. }, Expr::Identifier(bound))
        if name == RESULT_OK && fields.first() == Some(bound));
    let is_err = matches!(&err.pattern, Pattern::Constructor { name, .. } if name == RESULT_ERR);
    (is_get && is_ok && is_err).then_some((buffer, index, &err.body))
}

/// A two-arm boolean match, in either order; a wildcard second arm is the
/// other case.
fn bool_arms(arms: &[MatchArm]) -> Option<(&Expr, &Expr)> {
    let [first, second] = arms else {
        return None;
    };
    let literal = |pattern: &Pattern| match pattern {
        Pattern::Literal(value) => match value.as_ref() {
            Expr::Bool(b) => Some(*b),
            _ => None,
        },
        _ => None,
    };
    let (Some(a), b) = (literal(&first.pattern), literal(&second.pattern)) else {
        return None;
    };
    let other_is_rest = b == Some(!a) || matches!(second.pattern, Pattern::Wildcard);
    other_is_rest.then_some(if a {
        (&first.body, &second.body)
    } else {
        (&second.body, &first.body)
    })
}

impl Emitter<'_> {
    fn line(&mut self, text: String) {
        if let Some(block) = self.blocks.last_mut() {
            block.push(text);
        }
    }

    fn temp(&mut self) -> String {
        self.temps += 1;
        format!("t{}", self.temps)
    }

    fn declare(&mut self, ty: Ty, value: &str) -> String {
        let name = self.temp();
        self.line(format!("{} {name} = {value};", ty.spelling()));
        name
    }

    fn bind(&mut self, name: &str, local: Local) {
        if let Some(scope) = self.scopes.last_mut() {
            let _ = scope.insert(name.to_string(), local);
        }
    }

    fn lookup(&self, name: &str) -> Option<Local> {
        self.scopes
            .iter()
            .rev()
            .find_map(|scope| scope.get(name).cloned())
    }

    /// Run `emit` in its own block and scope; its statements come back with
    /// its value, for the caller to place under a condition or a loop.
    fn nested(
        &mut self,
        emit: impl FnOnce(&mut Self) -> Emitted,
    ) -> Result<(Vec<String>, Ty, String), String> {
        self.blocks.push(Vec::new());
        self.scopes.push(HashMap::new());
        let result = emit(self);
        let _ = self.scopes.pop();
        let lines = self.blocks.pop().unwrap_or_default();
        let (ty, value) = result?;
        Ok((lines, ty, value))
    }

    fn indented(&mut self, lines: Vec<String>) {
        for line in lines {
            self.line(format!("  {line}"));
        }
    }

    fn expression(&mut self, expression: &Expr) -> Emitted {
        match expression {
            Expr::Integer(n) => Ok((Ty::Int, int_literal(*n))),
            Expr::Float(f) => Ok((Ty::Float, float_literal(*f))),
            Expr::Bool(b) => Ok((Ty::Bool, if *b { "true" } else { "false" }.into())),
            Expr::Identifier(name) => self.identifier(name),
            Expr::Binary {
                op, left, right, ..
            } => self.binary(op, left, right),
            Expr::Unary { op, operand } => self.unary(op, operand),
            Expr::Call {
                function,
                arguments,
                named_arguments,
            } => self.call(callee(function), arguments, named_arguments),
            Expr::Pipe { left, right } => self.pipe(left, right),
            Expr::Match { value, arms } => self.matching(value, arms),
            Expr::Block {
                statements, value, ..
            } => self.block(statements, value.as_deref()),
            Expr::TypeApply { function, .. } => self.expression(function),
            other => Err(format!("{} has no device form", describe(other))),
        }
    }

    fn identifier(&mut self, name: &str) -> Emitted {
        match self.lookup(name) {
            Some(Local::Scalar(ty, value)) => Ok((ty, value)),
            Some(Local::Buffer(..)) => Err(format!("`{name}` is a buffer, not a kernel value")),
            None => match self.cg.device.constants.get(name).cloned() {
                Some(literal) => self.expression(&literal),
                None => Err(format!("`{name}` is not a value a device kernel can reach")),
            },
        }
    }

    fn binary(&mut self, op: &str, left: &Expr, right: &Expr) -> Emitted {
        if op == "&&" || op == "||" {
            return self.short_circuit(op, left, right);
        }
        let (left_ty, l) = self.expression(left)?;
        let (right_ty, r) = self.expression(right)?;
        match (left_ty, right_ty) {
            (Ty::Int, Ty::Int) => int_binary(op, &l, &r),
            (Ty::Float, Ty::Float) => float_binary(op, &l, &r),
            (Ty::Bool, Ty::Bool) if op == "==" || op == "!=" => {
                Ok((Ty::Bool, format!("({l} {op} {r})")))
            }
            _ => Err(format!("`{op}` over {left_ty:?} and {right_ty:?} operands")),
        }
    }

    /// `&&`/`||` keep the host's short circuit, so a fallible operation on the
    /// right never runs when the left decides the answer.
    fn short_circuit(&mut self, op: &str, left: &Expr, right: &Expr) -> Emitted {
        let (ty, condition) = self.expression(left)?;
        if ty != Ty::Bool {
            return Err(format!("`{op}` over a {ty:?} operand"));
        }
        if op == "&&" {
            self.branch(&condition, right, &Expr::Bool(false))
        } else {
            self.branch(&condition, &Expr::Bool(true), right)
        }
    }

    fn unary(&mut self, op: &str, operand: &Expr) -> Emitted {
        let (ty, value) = self.expression(operand)?;
        match (op, ty) {
            ("-", Ty::Int) => Ok((Ty::Int, format!("osprey_neg({value}, fault)"))),
            ("-", Ty::Float) => Ok((Ty::Float, format!("neg64({value})"))),
            ("!" | "not", Ty::Bool) => Ok((Ty::Bool, format!("(!{value})"))),
            _ => Err(format!("unary `{op}` over a {ty:?} operand")),
        }
    }

    fn call(&mut self, function: &Expr, arguments: &[Expr], named: &[NamedArgument]) -> Emitted {
        if !named.is_empty() {
            return Err("named arguments have no device form".into());
        }
        let Expr::Identifier(name) = function else {
            return Err("a computed callee has no device form".into());
        };
        if self.lookup(name).is_some() {
            return Err(format!(
                "`{name}` is a kernel value, not a function the device can call"
            ));
        }
        match name.as_str() {
            "fold" => self.fold(arguments),
            "gpuLength" => self.length(arguments),
            _ if self.cg.device.defs.contains_key(name) => self.inline(name, arguments),
            _ => self.builtin(name, arguments),
        }
    }

    fn builtin(&mut self, name: &str, arguments: &[Expr]) -> Emitted {
        let values = arguments
            .iter()
            .map(|argument| self.expression(argument))
            .collect::<Result<Vec<_>, _>>()?;
        let unary = |ty: Ty, text: String| Ok((ty, text));
        match (name, values.as_slice()) {
            ("intDiv", [(Ty::Int, a), (Ty::Int, b)]) => {
                unary(Ty::Int, format!("osprey_div({a}, {b}, fault)"))
            }
            ("toFloat", [(Ty::Int, a)]) => unary(Ty::Float, format!("from_long({a})")),
            ("sqrt" | "exp" | "ln" | "tanh", [(Ty::Float, a)]) => {
                unary(Ty::Float, format!("{name}64({a})"))
            }
            ("abs", [(Ty::Int, a)]) => unary(Ty::Int, format!("osprey_abs({a}, fault)")),
            ("abs", [(Ty::Float, a)]) => unary(Ty::Float, format!("abs64({a})")),
            ("wrapAdd", [(Ty::Int, a), (Ty::Int, b)]) => {
                unary(Ty::Int, format!("(long)((ulong)({a}) + (ulong)({b}))"))
            }
            ("wrapSub", [(Ty::Int, a), (Ty::Int, b)]) => {
                unary(Ty::Int, format!("(long)((ulong)({a}) - (ulong)({b}))"))
            }
            ("wrapMul", [(Ty::Int, a), (Ty::Int, b)]) => {
                unary(Ty::Int, format!("(long)((ulong)({a}) * (ulong)({b}))"))
            }
            _ => Err(format!("`{name}` has no device form")),
        }
    }

    /// Inline a helper: each argument is bound once, in order, then the body
    /// is emitted in a scope holding the parameters.
    fn inline(&mut self, name: &str, arguments: &[Expr]) -> Emitted {
        let Some((parameters, body)) = self.cg.device.defs.get(name).cloned() else {
            return Err(format!("`{name}` is not a function the device can inline"));
        };
        if parameters.len() != arguments.len() {
            return Err(format!(
                "`{name}` is called with {} of {} arguments",
                arguments.len(),
                parameters.len()
            ));
        }
        if self.depth >= MAX_INLINE_DEPTH {
            return Err(format!("`{name}` recurses, which a device kernel cannot"));
        }
        let mut bindings = Vec::with_capacity(parameters.len());
        for (parameter, argument) in parameters.iter().zip(arguments) {
            bindings.push((parameter.name.clone(), self.argument(argument)?));
        }
        self.apply(&body, bindings)
    }

    /// Emit `body` under `bindings` one inlining level deeper.
    fn apply(&mut self, body: &Expr, bindings: Vec<(String, Local)>) -> Emitted {
        self.depth += 1;
        self.scopes.push(bindings.into_iter().collect());
        let result = self.expression(body);
        let _ = self.scopes.pop();
        self.depth -= 1;
        result
    }

    /// A call argument as a local: a buffer passes by alias, a scalar is
    /// evaluated once into a temporary.
    fn argument(&mut self, argument: &Expr) -> Result<Local, String> {
        if let Expr::Identifier(name) = argument {
            if let Some(buffer @ Local::Buffer(..)) = self.lookup(name) {
                return Ok(buffer);
            }
        }
        let (ty, value) = self.expression(argument)?;
        Ok(Local::Scalar(ty, self.declare(ty, &value)))
    }

    fn length(&mut self, arguments: &[Expr]) -> Emitted {
        match arguments {
            [Expr::Identifier(name)] => match self.lookup(name) {
                Some(Local::Buffer(_, len)) => Ok((Ty::Int, format!("(long){len}"))),
                _ => Err(format!(
                    "gpuLength of `{name}`, which is not a captured buffer"
                )),
            },
            _ => Err("gpuLength of a computed buffer has no device form".into()),
        }
    }

    /// `fold(range(start, end), seed, combine)` as an in-order loop.
    fn fold(&mut self, arguments: &[Expr]) -> Emitted {
        let [source, seed, combine] = arguments else {
            return Err("fold needs a range, a seed and a combine".into());
        };
        let (start, end) = self.range_bounds(source)?;
        let (seed_ty, seed_value) = self.expression(seed)?;
        let acc = self.declare(seed_ty, &seed_value);
        let index = self.temp();
        let step = [
            Local::Scalar(seed_ty, acc.clone()),
            Local::Scalar(Ty::Int, index.clone()),
        ];
        let (lines, ty, value) = self.nested(|emitter| emitter.callback(combine, &step))?;
        if ty != seed_ty {
            return Err(format!(
                "a fold over a {seed_ty:?} seed combines to a {ty:?}"
            ));
        }
        self.line(format!(
            "for (long {index} = {start}; {index} < {end}; {index}++) {{"
        ));
        self.indented(lines);
        self.line(format!("  {acc} = {value};"));
        self.line("}".into());
        Ok((seed_ty, acc))
    }

    /// The bounds of `range(start, end)`, each evaluated once before the loop.
    fn range_bounds(&mut self, source: &Expr) -> Result<(String, String), String> {
        let Expr::Call {
            function,
            arguments,
            ..
        } = source
        else {
            return Err("a fold over something other than a range has no device form".into());
        };
        let ([start, end], true) = (
            arguments.as_slice(),
            matches!(callee(function), Expr::Identifier(name) if name == "range"),
        ) else {
            return Err(
                "a fold over something other than `range(start, end)` has no device form".into(),
            );
        };
        let (start_ty, start) = self.expression(start)?;
        let (end_ty, end) = self.expression(end)?;
        if start_ty != Ty::Int || end_ty != Ty::Int {
            return Err("range bounds that are not ints".into());
        }
        Ok((self.declare(Ty::Int, &start), self.declare(Ty::Int, &end)))
    }

    /// Apply a callback — a lambda or a named helper — to already-bound locals.
    fn callback(&mut self, callback: &Expr, values: &[Local]) -> Emitted {
        let (parameters, body) =
            match callee(callback) {
                Expr::Lambda {
                    parameters, body, ..
                } => (parameters.clone(), body.as_ref().clone()),
                Expr::Identifier(name) => {
                    self.cg.device.defs.get(name).cloned().ok_or_else(|| {
                        format!("`{name}` is not a function the device can inline")
                    })?
                }
                other => {
                    return Err(format!(
                        "{} as a callback has no device form",
                        describe(other)
                    ))
                }
            };
        if parameters.len() != values.len() {
            return Err("a callback whose arity disagrees with its combinator".into());
        }
        let bindings = parameters
            .iter()
            .zip(values)
            .map(|(parameter, local)| (parameter.name.clone(), local.clone()))
            .collect();
        self.apply(&body, bindings)
    }

    fn pipe(&mut self, left: &Expr, right: &Expr) -> Emitted {
        match callee(right) {
            Expr::Call {
                function,
                arguments,
                named_arguments,
            } => {
                let mut all = vec![left.clone()];
                all.extend(arguments.iter().cloned());
                self.call(callee(function), &all, named_arguments)
            }
            function @ Expr::Identifier(_) => self.call(function, std::slice::from_ref(left), &[]),
            other => Err(format!(
                "a pipe into {} has no device form",
                describe(other)
            )),
        }
    }

    fn matching(&mut self, value: &Expr, arms: &[MatchArm]) -> Emitted {
        if let Some((buffer, index, fallback)) = elvis_gather(value, arms) {
            return self.gather(buffer, index, fallback);
        }
        let Some((then, otherwise)) = bool_arms(arms) else {
            return Err("a match other than `true`/`false` arms or `gpuGet(...) ?: default` has no device form".into());
        };
        let (ty, condition) = self.expression(value)?;
        if ty != Ty::Bool {
            return Err(format!("a match over a {ty:?} scrutinee with boolean arms"));
        }
        self.branch(&condition, then, otherwise)
    }

    /// `if (condition) then else otherwise` into one temporary.
    fn branch(&mut self, condition: &str, then: &Expr, otherwise: &Expr) -> Emitted {
        let (then_lines, then_ty, then_value) = self.nested(|emitter| emitter.expression(then))?;
        let (else_lines, else_ty, else_value) =
            self.nested(|emitter| emitter.expression(otherwise))?;
        if then_ty != else_ty {
            return Err(format!("arms of a {then_ty:?} and a {else_ty:?}"));
        }
        let result = self.temp();
        self.line(format!("{} {result};", then_ty.spelling()));
        self.line(format!("if ({condition}) {{"));
        self.indented(then_lines);
        self.line(format!("  {result} = {then_value};"));
        self.line("} else {".into());
        self.indented(else_lines);
        self.line(format!("  {result} = {else_value};"));
        self.line("}".into());
        Ok((then_ty, result))
    }

    /// `gpuGet(buffer, index) ?: fallback` over a captured buffer: the word at
    /// a valid index, typed as the fallback is, else the fallback [GPU-GET].
    fn gather(&mut self, buffer: &Expr, index: &Expr, fallback: &Expr) -> Emitted {
        let Expr::Identifier(name) = buffer else {
            return Err("gpuGet over a computed buffer has no device form".into());
        };
        let Some(Local::Buffer(words, len)) = self.lookup(name) else {
            return Err(format!(
                "gpuGet over `{name}`, which is not a captured buffer"
            ));
        };
        let (index_ty, index) = self.expression(index)?;
        if index_ty != Ty::Int {
            return Err("a gpuGet index that is not an int".into());
        }
        let (fallback_lines, ty, fallback_value) =
            self.nested(|emitter| emitter.expression(fallback))?;
        let (result, at) = (self.temp(), self.temp());
        self.line(format!("{} {result};", ty.spelling()));
        self.line(format!("long {at} = {index};"));
        self.line(format!("if ({at} >= 0 && (ulong){at} < {len}) {{"));
        self.line(format!(
            "  {result} = {};",
            ty.from_word(&format!("{words}[{at}]"))
        ));
        self.line("} else {".into());
        self.indented(fallback_lines);
        self.line(format!("  {result} = {fallback_value};"));
        self.line("}".into());
        Ok((ty, result))
    }

    fn block(&mut self, statements: &[Stmt], value: Option<&Expr>) -> Emitted {
        self.scopes.push(HashMap::new());
        let result = self.block_contents(statements, value);
        let _ = self.scopes.pop();
        result
    }

    fn block_contents(&mut self, statements: &[Stmt], value: Option<&Expr>) -> Emitted {
        for statement in statements {
            self.statement(statement)?;
        }
        match value {
            Some(value) => self.expression(value),
            None => Err("a block without a value has no device form".into()),
        }
    }

    fn statement(&mut self, statement: &Stmt) -> Result<(), String> {
        let Stmt::Let {
            name,
            mutable: false,
            value,
            ..
        } = statement
        else {
            return Err("a statement other than an immutable `let` has no device form".into());
        };
        let local = self.argument(value)?;
        self.bind(name, local);
        Ok(())
    }
}

fn int_binary(op: &str, l: &str, r: &str) -> Emitted {
    Ok(match op {
        "+" => (Ty::Int, format!("osprey_add({l}, {r}, fault)")),
        "-" => (Ty::Int, format!("osprey_sub({l}, {r}, fault)")),
        "*" => (Ty::Int, format!("osprey_mul({l}, {r}, fault)")),
        "%" => (Ty::Int, format!("osprey_rem({l}, {r}, fault)")),
        "/" => (
            Ty::Float,
            format!("osprey_fdiv(from_long({l}), from_long({r}), fault)"),
        ),
        "==" | "!=" | "<" | "<=" | ">" | ">=" => (Ty::Bool, format!("({l} {op} {r})")),
        _ => return Err(format!("`{op}` over ints has no device form")),
    })
}

fn float_binary(op: &str, l: &str, r: &str) -> Emitted {
    Ok(match op {
        "+" => (Ty::Float, format!("add64({l}, {r})")),
        "-" => (Ty::Float, format!("sub64({l}, {r})")),
        "*" => (Ty::Float, format!("multiply64({l}, {r})")),
        "/" => (Ty::Float, format!("osprey_fdiv({l}, {r}, fault)")),
        "==" => (Ty::Bool, format!("eq64({l}, {r})")),
        "!=" => (Ty::Bool, format!("(!eq64({l}, {r}))")),
        "<" => (Ty::Bool, format!("lt64({l}, {r})")),
        "<=" => (Ty::Bool, format!("le64({l}, {r})")),
        ">" => (Ty::Bool, format!("lt64({r}, {l})")),
        ">=" => (Ty::Bool, format!("le64({r}, {l})")),
        _ => return Err(format!("`{op}` over floats has no device form")),
    })
}
