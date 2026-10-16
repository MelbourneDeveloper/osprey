//! The statement-bearing shapes of a kernel body: a `fold` over a `range` as
//! an in-order loop, and `gpuGet ?: default` as a bounds-checked read
//! [GPU-BACKEND-DEVICE] [GPU-GET].

use super::dialect::{callee, describe};
use super::emit::{Emitted, Emitter, Local};
use super::Ty;
use osprey_ast::Expr;

impl Emitter<'_> {
    /// `fold(range(start, end), seed, combine)` as an in-order loop.
    pub(super) fn fold(&mut self, arguments: &[Expr]) -> Emitted {
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

    /// `gpuGet(buffer, index) ?: fallback` over a captured buffer: the word at
    /// a valid index, typed as the fallback is, else the fallback [GPU-GET].
    pub(super) fn gather(&mut self, buffer: &Expr, index: &Expr, fallback: &Expr) -> Emitted {
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
            ty.value_of(&format!("{words}[{at}]"))
        ));
        self.line("} else {".into());
        self.indented(fallback_lines);
        self.line(format!("  {result} = {fallback_value};"));
        self.line("}".into());
        Ok((ty, result))
    }
}
