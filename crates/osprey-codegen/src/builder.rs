//! The emitter state: a growing LLVM module (external declarations, string
//! globals, finished functions) plus the in-progress function (SSA counter,
//! current basic block, lexical scopes). Low-level helpers here only *emit*
//! text; the AST-walking lives in `lower.rs`.

use crate::llty::{LType, Value};
use crate::types::{ltype_of, ltype_of_name};
use osprey_types::{ProgramTypes, Type};
use std::collections::{BTreeSet, HashMap, HashSet};
use std::fmt::Write as _;

/// Accumulates a whole module while lowering one function at a time.
pub struct Codegen {
    /// `declare` lines, de-duplicated and stably ordered.
    externs: BTreeSet<String>,
    /// Global constant definitions (string literals).
    globals: Vec<String>,
    /// Rendered `define` blocks.
    funcs: Vec<String>,
    glob_count: usize,

    // ---- current function state ----
    reg_count: usize,
    label_count: usize,
    cur_lines: Vec<String>,
    cur_block: String,
    scopes: Vec<HashMap<String, Value>>,

    /// Declared parameter names per function, for named-argument ordering.
    pub(crate) fn_params: HashMap<String, Vec<String>>,
    /// Resolved signatures, constructor layouts and union tags from inference.
    pub(crate) prog: ProgramTypes,
    /// Stream-fusion pipeline: pending `map`/`filter` stages recorded by those
    /// builtins and replayed (in source order) when `forEach`/`fold` consumes
    /// the iterator. Cleared after each consumer. Ports `pendingIterOps`.
    pub(crate) pending_iter_ops: Vec<crate::iter::IterOp>,
    /// Let-bound lambdas, stored for inline application at their call sites
    /// (`let f = fn(x) => …` then `f(y)`), since the backend lowers no closures.
    pub(crate) lambdas: HashMap<String, (Vec<osprey_ast::Parameter>, osprey_ast::Expr)>,
    /// Whether the fiber-result global table has been emitted yet.
    pub(crate) fiber_table_emitted: bool,
    /// Parsed `effect` operation signatures, keyed `"Effect.operation"`.
    pub(crate) effect_ops: HashMap<String, crate::effects::OpSig>,
    /// Monotonic id giving each emitted handler function a unique name.
    pub(crate) handler_count: usize,
    /// Synthetic layouts of anonymous object literals (`{ a: 1, b: "x" }`),
    /// keyed by the generated owner name carried on the handle, so field access
    /// can recover the ordered `(field, LType)` slots.
    pub(crate) obj_layouts: HashMap<String, Vec<(String, LType)>>,
    /// Monotonic id giving each object literal a unique synthetic owner name.
    pub(crate) obj_count: usize,
    /// User function `(parameters, body)` defs, for inlining a *generic*
    /// function at each call site so its type variables monomorphize to the
    /// concrete argument types there (the Go backend emits a mangled copy per
    /// instantiation; inlining achieves the same without name mangling).
    pub(crate) fn_defs: HashMap<String, (Vec<osprey_ast::Parameter>, osprey_ast::Expr)>,
    /// Generic functions currently being inlined — a re-entry guard so a
    /// (mutually) recursive generic call falls back to a direct call instead of
    /// inlining forever.
    pub(crate) inlining: HashSet<String>,
    /// Function-typed locals in the current function (a higher-order parameter
    /// `f: (int) -> int`): name → its signature ([`FnSig`]), so a call `f(x)`
    /// lowers to an indirect call through the `i8*` handle.
    pub(crate) fn_ptr_locals: HashMap<String, FnSig>,
    /// While inlining a generic function, a function-valued parameter bound to a
    /// callee *by name* (`apply(f: toString, …)`): the parameter redirects to
    /// that real callee, so `f(x)` in the body becomes `toString(x)`. This keeps
    /// a builtin or another generic function callable through the parameter.
    pub(crate) call_aliases: HashMap<String, String>,
}

/// A function value's lowered signature: parameter [`LType`]s, the return
/// [`LType`], and (when it returns `Result<T, _>`) the success inner type.
pub(crate) type FnSig = (Vec<LType>, LType, Option<LType>);

/// Saved emission state of a suspended function (see [`Codegen::enter_nested_fn`]).
pub(crate) struct SavedFn {
    lines: Vec<String>,
    block: String,
    regs: usize,
    labels: usize,
    scopes: Vec<HashMap<String, Value>>,
}

impl Codegen {
    pub fn new() -> Codegen {
        Codegen::with_types(ProgramTypes::default())
    }

    /// Build with the inferred program types that drive parameter/return/value
    /// typing.
    pub fn with_types(prog: ProgramTypes) -> Codegen {
        Codegen {
            externs: BTreeSet::new(),
            globals: Vec::new(),
            funcs: Vec::new(),
            glob_count: 0,
            reg_count: 0,
            label_count: 0,
            cur_lines: Vec::new(),
            cur_block: String::from("entry"),
            scopes: Vec::new(),
            fn_params: HashMap::new(),
            prog,
            pending_iter_ops: Vec::new(),
            lambdas: HashMap::new(),
            fiber_table_emitted: false,
            effect_ops: HashMap::new(),
            handler_count: 0,
            obj_layouts: HashMap::new(),
            obj_count: 0,
            fn_defs: HashMap::new(),
            inlining: HashSet::new(),
            fn_ptr_locals: HashMap::new(),
            call_aliases: HashMap::new(),
        }
    }

    /// Whether `name` is a user function whose inferred signature still contains
    /// a type variable (in a parameter or the return) — i.e. it is polymorphic
    /// and must be specialised to the concrete call-site types.
    pub(crate) fn is_generic_fn(&self, name: &str) -> bool {
        let Some((params, ret)) = self.prog.functions.get(name) else {
            return false;
        };
        params.iter().chain(std::iter::once(ret)).any(has_type_var)
    }

    /// The declared type parameters of a constructor's owner (`["T"]` for
    /// `Generic<T>`), used to spot a generic field whose LLVM type is fixed per
    /// construction rather than by the (placeholder) written type.
    pub(crate) fn ctor_type_params(&self, name: &str) -> Vec<String> {
        self.prog
            .ctors
            .get(name)
            .map(|c| c.type_params.clone())
            .unwrap_or_default()
    }

    /// The lowered [`FnSig`] of a function-typed value `ty` (a higher-order
    /// parameter), for the indirect-call bitcast — `None` if `ty` is not a
    /// function.
    pub(crate) fn fn_value_sig(ty: &Type) -> Option<FnSig> {
        match ty {
            Type::Fun { params, ret } => Some((
                params.iter().map(ltype_of).collect(),
                ltype_of(ret),
                crate::types::result_inner(ret),
            )),
            _ => None,
        }
    }

    /// Register an anonymous object literal's ordered field layout and return the
    /// synthetic owner name to tag its handle with.
    pub(crate) fn register_obj_layout(&mut self, fields: Vec<(String, LType)>) -> String {
        let name = format!("__obj_{}", self.obj_count);
        self.obj_count += 1;
        let _ = self.obj_layouts.insert(name.clone(), fields);
        name
    }

    /// The struct spelling and ordered fields of an owner — a real constructor or
    /// a synthetic object literal — for unified field access.
    pub(crate) fn record_layout(&self, owner: &str) -> Option<(String, Vec<(String, LType)>)> {
        if let Some(fields) = self.obj_layouts.get(owner) {
            let mut parts = vec!["i64".to_string()];
            parts.extend(fields.iter().map(|(_, lt)| lt.as_str().to_string()));
            return Some((format!("{{ {} }}", parts.join(", ")), fields.clone()));
        }
        let view = self.ctor_layout(owner)?;
        Some((self.ctor_struct_ty(owner)?, view.fields))
    }

    /// A fresh, module-unique handler-function id.
    pub(crate) fn next_handler_id(&mut self) -> usize {
        let id = self.handler_count;
        self.handler_count += 1;
        id
    }

    /// Append a module-level global definition (e.g. the fiber-result table).
    pub(crate) fn add_global_def(&mut self, def: impl Into<String>) {
        self.globals.push(def.into());
    }

    /// Suspend the in-progress function and start a fresh one (a handler function
    /// emitted while lowering its enclosing `handle`). Returns the saved state to
    /// hand back to [`Codegen::exit_nested_fn`]. The new function gets its own
    /// SSA/label counters and an isolated scope stack (handlers capture nothing).
    pub(crate) fn enter_nested_fn(&mut self) -> SavedFn {
        let saved = SavedFn {
            lines: std::mem::take(&mut self.cur_lines),
            block: std::mem::replace(&mut self.cur_block, String::from("entry")),
            regs: self.reg_count,
            labels: self.label_count,
            scopes: std::mem::take(&mut self.scopes),
        };
        self.reg_count = 0;
        self.label_count = 0;
        self.cur_lines = vec!["entry:".to_string()];
        self.scopes = vec![HashMap::new()];
        saved
    }

    /// Finish the nested function (append it) and resume the suspended one.
    pub(crate) fn exit_nested_fn(
        &mut self,
        saved: SavedFn,
        ret: &str,
        name: &str,
        params: &[(LType, String)],
    ) {
        self.finish_function(ret, name, params);
        self.cur_lines = saved.lines;
        self.cur_block = saved.block;
        self.reg_count = saved.regs;
        self.label_count = saved.labels;
        self.scopes = saved.scopes;
    }

    /// Register an `effect` operation's parsed signature for `handle`/`perform`.
    pub(crate) fn register_effect_op(&mut self, key: String, sig: crate::effects::OpSig) {
        let _ = self.effect_ops.insert(key, sig);
    }

    /// The parsed signature of `Effect.operation`, if declared.
    pub(crate) fn effect_op(&self, key: &str) -> Option<crate::effects::OpSig> {
        self.effect_ops.get(key).cloned()
    }

    // ---- inferred typing ----

    /// The LLVM return type of a user/runtime function, from inference.
    pub(crate) fn fn_ret_ltype(&self, name: &str) -> Option<LType> {
        self.prog.return_type(name).map(ltype_of)
    }

    /// The LLVM parameter types of a user function, from inference.
    pub(crate) fn fn_param_ltypes(&self, name: &str) -> Option<Vec<LType>> {
        self.prog
            .param_types(name)
            .map(|ps| ps.iter().map(ltype_of).collect())
    }

    /// The `(LType, owner)` parameter signature — `owner` tags record/union
    /// parameters so their fields are reachable inside the body.
    pub(crate) fn fn_param_sig(&self, name: &str) -> Option<Vec<(LType, Option<String>)>> {
        self.prog.param_types(name).map(|ps| {
            ps.iter()
                .map(|t| (ltype_of(t), crate::types::owner_name(t)))
                .collect()
        })
    }

    /// The owner type name of a function's return value, if it is a record/union.
    pub(crate) fn fn_ret_owner(&self, name: &str) -> Option<String> {
        self.prog
            .return_type(name)
            .and_then(crate::types::owner_name)
    }

    /// The inner [`LType`] when a function is declared to return `Result<T, E>`
    /// — the success payload's LLVM type — so calls and returns carry the
    /// `{ T, i8 }*` Result block rather than a bare `T`.
    pub(crate) fn fn_ret_result_inner(&self, name: &str) -> Option<LType> {
        crate::types::result_inner(self.prog.return_type(name)?)
    }

    /// The full heap layout of a constructor: owning type, whether it is a
    /// record, the discriminant tag (variant index within its union; 0 for a
    /// record), and ordered `(field, LType)` pairs.
    pub(crate) fn ctor_layout(&self, name: &str) -> Option<CtorView> {
        let c = self.prog.ctors.get(name)?;
        let tag = i64::try_from(
            self.prog
                .unions
                .get(&c.owner)
                .and_then(|vs| vs.iter().position(|v| v == name))
                .unwrap_or(0),
        )
        .unwrap_or(0);
        let fields = c
            .fields
            .iter()
            .map(|(f, t)| (f.clone(), ltype_of_name(t)))
            .collect();
        Some(CtorView {
            owner: c.owner.clone(),
            owner_is_record: c.owner_is_record,
            tag,
            fields,
        })
    }

    /// Resolve a field name to an owning constructor when the target's static
    /// type is unknown — polymorphic field access inside a generic accessor like
    /// `fn getFirst(p) = p.first`, where `p` infers to a type variable. Mirrors
    /// Go's `generateStructFieldAccessFallback`. Prefers a layout whose field
    /// type is a concrete scalar (so the load type and `toString` match the
    /// runtime value), breaking ties by owner name for deterministic output.
    pub(crate) fn find_field_owner(&self, field: &str) -> Option<String> {
        let mut candidates: Vec<(&String, LType)> = self
            .prog
            .ctors
            .iter()
            .filter_map(|(name, c)| {
                c.fields
                    .iter()
                    .find(|(f, _)| f == field)
                    .map(|(_, t)| (name, ltype_of_name(t)))
            })
            .collect();
        candidates.sort_by(|a, b| a.0.cmp(b.0));
        candidates
            .iter()
            .find(|(_, lt)| *lt != LType::Ptr)
            .or_else(|| candidates.first())
            .map(|(name, _)| (*name).clone())
    }

    /// The LLVM struct spelling for a constructor's heap block: `{ i64, f0, … }`
    /// — a leading `i64` discriminant tag followed by each field's LLVM type.
    pub(crate) fn ctor_struct_ty(&self, name: &str) -> Option<String> {
        let view = self.ctor_layout(name)?;
        let mut parts = vec!["i64".to_string()];
        for (_, lt) in &view.fields {
            parts.push(lt.as_str().to_string());
        }
        Some(format!("{{ {} }}", parts.join(", ")))
    }

    /// Whether a name is a known constructor.
    pub(crate) fn is_ctor(&self, name: &str) -> bool {
        self.prog.ctors.contains_key(name)
    }

    /// The owner type name to tag a loaded aggregate field with: the field's
    /// written type when that type is itself a known record/union, else `None`
    /// (scalars carry no owner).
    pub(crate) fn ctor_field_written(&self, owner: &str, field: &str) -> Option<String> {
        let written = self
            .prog
            .ctors
            .get(owner)?
            .fields
            .iter()
            .find(|(f, _)| f == field)
            .map(|(_, t)| t.clone())?;
        let head = written.split(['<', '[']).next().unwrap_or(&written).trim();
        if self.prog.ctors.contains_key(head) || self.prog.unions.contains_key(head) {
            Some(head.to_string())
        } else {
            None
        }
    }

    /// The variant constructor names of a union owner, in tag order.
    pub(crate) fn union_variants(&self, owner: &str) -> Option<&[String]> {
        self.prog.unions.get(owner).map(std::vec::Vec::as_slice)
    }

    // ---- SSA + block naming (function-local) ----

    pub(crate) fn fresh_reg(&mut self) -> String {
        let r = format!("%r{}", self.reg_count);
        self.reg_count += 1;
        r
    }

    pub(crate) fn fresh_label(&mut self) -> String {
        let l = format!("L{}", self.label_count);
        self.label_count += 1;
        l
    }

    pub(crate) fn cur_block(&self) -> &str {
        &self.cur_block
    }

    // ---- emission ----

    pub(crate) fn emit(&mut self, line: impl Into<String>) {
        self.cur_lines.push(format!("  {}", line.into()));
    }

    /// Emit `r = {rhs}` to a fresh SSA register and return `r` — the ubiquitous
    /// "name the result of one instruction" step (`zext …`, `icmp …`, `fneg …`).
    pub(crate) fn emit_reg(&mut self, rhs: impl std::fmt::Display) -> String {
        let r = self.fresh_reg();
        self.emit(format!("{r} = {rhs}"));
        r
    }

    /// Start a new basic block and make it current (its label becomes the
    /// predecessor recorded for any `phi` that follows).
    pub(crate) fn start_block(&mut self, label: &str) {
        self.cur_lines.push(format!("{label}:"));
        self.cur_block = label.to_string();
    }

    /// Snapshot the current block label, then branch to `end` — the predecessor
    /// a `phi` at `end` reads back. Closes a one-arm path of a Result/match split.
    pub(crate) fn snapshot_to(&mut self, end: &str) -> String {
        let block = self.cur_block.clone();
        self.emit(format!("br label %{end}"));
        block
    }

    pub(crate) fn add_extern(&mut self, decl: impl Into<String>) {
        let _ = self.externs.insert(decl.into());
    }

    /// Intern a string literal as a private global and return an `i8*` pointing
    /// at its first byte.
    pub(crate) fn string_constant(&mut self, text: &str) -> Value {
        let (escaped, len) = escape_c_string(text);
        let name = format!("@.str.{}", self.glob_count);
        self.glob_count += 1;
        self.globals.push(format!(
            "{name} = private unnamed_addr constant [{len} x i8] c\"{escaped}\""
        ));
        let reg = self.fresh_reg();
        self.emit(format!(
            "{reg} = getelementptr [{len} x i8], [{len} x i8]* {name}, i64 0, i64 0"
        ));
        Value::new(reg, LType::Str)
    }

    // ---- scopes ----

    pub(crate) fn push_scope(&mut self) {
        self.scopes.push(HashMap::new());
    }

    pub(crate) fn pop_scope(&mut self) {
        let _ = self.scopes.pop();
    }

    pub(crate) fn bind(&mut self, name: impl Into<String>, value: Value) {
        if let Some(scope) = self.scopes.last_mut() {
            let _ = scope.insert(name.into(), value);
        }
    }

    pub(crate) fn lookup(&self, name: &str) -> Option<Value> {
        self.scopes.iter().rev().find_map(|s| s.get(name).cloned())
    }

    // ---- function framing ----

    /// Reset per-function state and open a fresh `entry` block + scope.
    pub(crate) fn begin_function(&mut self) {
        self.reg_count = 0;
        self.label_count = 0;
        self.cur_lines.clear();
        self.cur_block = String::from("entry");
        self.fn_ptr_locals.clear();
        self.push_scope();
        self.cur_lines.push("entry:".to_string());
    }

    /// Render the in-progress function and append it to the module. `ret` is the
    /// already-rendered LLVM return type (`i64`, `{ i1, i8 }*`, …).
    pub(crate) fn finish_function(&mut self, ret: &str, name: &str, params: &[(LType, String)]) {
        let param_list = params
            .iter()
            .map(|(ty, n)| format!("{ty} %{n}"))
            .collect::<Vec<_>>()
            .join(", ");
        let body = std::mem::take(&mut self.cur_lines).join("\n");
        self.funcs
            .push(format!("define {ret} @{name}({param_list}) {{\n{body}\n}}"));
        self.pop_scope();
    }

    /// Assemble the final module text: header, externals, globals, functions.
    pub(crate) fn render(&self) -> String {
        let mut out = String::from("; Generated by osprey-rs (Rust LLVM-text backend)\n\n");
        for decl in &self.externs {
            out.push_str(decl);
            out.push('\n');
        }
        out.push('\n');
        for g in &self.globals {
            out.push_str(g);
            out.push('\n');
        }
        out.push('\n');
        out.push_str(&self.funcs.join("\n\n"));
        out.push('\n');
        out
    }

    /// Allocate a heap block sized for the LLVM struct type `struct_ty`, via the
    /// portable `getelementptr null, 1` sizeof trick, and return the typed
    /// pointer register (`{TY}*`).
    pub(crate) fn malloc_struct(&mut self, struct_ty: &str) -> String {
        self.add_extern("declare i8* @malloc(i64)");
        let szp = self.fresh_reg();
        self.emit(format!(
            "{szp} = getelementptr {struct_ty}, {struct_ty}* null, i64 1"
        ));
        let sz = self.fresh_reg();
        self.emit(format!("{sz} = ptrtoint {struct_ty}* {szp} to i64"));
        let raw = self.fresh_reg();
        self.emit(format!("{raw} = call i8* @malloc(i64 {sz})"));
        let obj = self.fresh_reg();
        self.emit(format!("{obj} = bitcast i8* {raw} to {struct_ty}*"));
        obj
    }
}

/// Whether a (fully substituted) inferred type still mentions a type variable —
/// the mark of a polymorphic signature the backend must specialise per use.
fn has_type_var(ty: &Type) -> bool {
    match ty {
        Type::Var(_) => true,
        Type::Con { args, .. } => args.iter().any(has_type_var),
        Type::Fun { params, ret } => params.iter().any(has_type_var) || has_type_var(ret),
        Type::Record { fields, .. } => fields.values().any(has_type_var),
        Type::Union { variants, .. } => variants.iter().any(has_type_var),
    }
}

/// The resolved heap layout of a constructor.
pub(crate) struct CtorView {
    pub owner: String,
    pub owner_is_record: bool,
    pub tag: i64,
    pub fields: Vec<(String, LType)>,
}

impl Default for Codegen {
    fn default() -> Self {
        Codegen::new()
    }
}

/// Escape a Rust string into an LLVM `c"..."` body, returning the escaped text
/// and the byte length **including** the trailing NUL. Bytes outside printable
/// ASCII (and `"`/`\`) are emitted as `\HH`.
fn escape_c_string(text: &str) -> (String, usize) {
    let mut out = String::new();
    let bytes = text.as_bytes();
    for &b in bytes {
        match b {
            b'\\' => out.push_str("\\5C"),
            b'"' => out.push_str("\\22"),
            0x20..=0x7e => out.push(char::from(b)),
            _ => {
                let _ = write!(out, "\\{b:02X}");
            }
        }
    }
    out.push_str("\\00");
    (out, bytes.len() + 1)
}
