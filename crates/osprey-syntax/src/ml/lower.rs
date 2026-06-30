//! The ML **lowerer**: CST ([`super::cst`]) → canonical [`osprey_ast::Program`].
//! This is the *only* place ML surface syntax is normalised into the shared
//! core, and it is where the boundary law is enforced — the output is canonical
//! AST that no later phase can distinguish from Default-flavor output
//! ([FLAVOR-BOUNDARY], [FLAVOR-LOWER-CONTRACT], docs/specs/0023).
//!
//! What this module canonicalises (and the parser deliberately does not):
//! - **Uncurried application** ([FLAVOR-CURRY]): ML is a thin *syntactic skin*
//!   over the canonical AST, not a curry-by-default calculus. A multi-parameter
//!   binding `f a b = …` lowers to one multi-parameter [`Stmt::Function`] (NOT
//!   nested lambdas), whitespace application `f a b` collapses its spine to one
//!   multi-argument [`Expr::Call`], and a lambda `\a b => …` becomes one
//!   multi-parameter [`Expr::Lambda`] — each byte-identical to the Default
//!   flavor's `fn f(a, b)`, `f(a, b)`, and `(a, b) => …`. This is what makes the
//!   IR-equivalence guarantee ([FLAVOR-IR-EQUIV]) hold for idiomatic, uncurried
//!   Default examples; partial application is written explicitly with a lambda
//!   in both flavors, never implied.
//! - **Pipes**: `x |> f` desugars to a call, exactly as the Default lowerer does.
//! - **Records / blocks / interpolation**: surface nodes map to
//!   [`Expr::TypeConstructor`], [`Expr::Block`], and [`Expr::InterpolatedStr`].

use super::cst::{
    MlArm, MlExpr, MlExternParam, MlField, MlItem, MlParam, MlPattern, MlType, MlTypeField,
    MlVariant,
};
use crate::strings::{lower_interpolation, unquote};
use osprey_ast::{
    Expr, ExternParameter, FieldAssignment, MatchArm, Parameter, Pattern, Position, Program, Stmt,
    TypeExpr, TypeField, TypeVariant,
};

/// Lower a parsed ML CST into the canonical program.
pub(crate) fn lower(items: Vec<MlItem>) -> Program {
    Program {
        statements: lower_items(items),
    }
}

/// Lower a run of items, pairing each type signature with the binding of the
/// same name that immediately follows it. An orphaned signature (no matching
/// binding next) is dropped. Used at top level and inside layout blocks so
/// local signed functions work too.
fn lower_items(items: Vec<MlItem>) -> Vec<Stmt> {
    let mut out = Vec::new();
    let mut pending: Option<(String, MlType)> = None;
    for item in items {
        match item {
            MlItem::Signature { name, ty } => pending = Some((name, ty)),
            MlItem::Binding {
                mutable,
                name,
                params,
                body,
                pos,
            } => {
                let sig = pending
                    .take()
                    .filter(|(signed, _)| *signed == name)
                    .map(|(_, ty)| ty);
                out.push(lower_binding(
                    mutable,
                    name,
                    params,
                    body,
                    pos,
                    sig.as_ref(),
                ));
            }
            MlItem::Assign { name, value, pos } => {
                pending = None;
                out.push(Stmt::Assignment {
                    name,
                    value: lower_expr(value),
                    position: Some(pos),
                });
            }
            MlItem::Type {
                name,
                type_params,
                variants,
                pos,
            } => {
                pending = None;
                out.push(Stmt::Type {
                    name,
                    type_params,
                    variants: variants.into_iter().map(lower_variant).collect(),
                    validation_func: None,
                    position: Some(pos),
                });
            }
            MlItem::Extern {
                name,
                params,
                return_type,
                pos,
            } => {
                pending = None;
                out.push(Stmt::Extern {
                    name,
                    parameters: params.into_iter().map(lower_extern_param).collect(),
                    return_type: return_type.as_ref().and_then(type_expr),
                    position: Some(pos),
                });
            }
            MlItem::Expr { value, pos } => {
                pending = None;
                out.push(Stmt::Expr {
                    value: lower_expr(value),
                    position: Some(pos),
                });
            }
        }
    }
    out
}

/// Lower one `extern` parameter to a canonical [`ExternParameter`], threading its
/// declared type through the shared [`type_expr`] path so it is byte-identical to
/// the Default flavor's extern parameter ([FLAVOR-ML-EXTERN]). A type with no
/// canonical [`TypeExpr`] form (a tuple) falls back to its rendered surface name.
fn lower_extern_param(param: MlExternParam) -> ExternParameter {
    let ty = type_expr(&param.ty).unwrap_or_else(|| TypeExpr::named(render_type(&param.ty)));
    ExternParameter {
        name: param.name,
        ty,
    }
}

/// Lower one CST variant to a canonical [`TypeVariant`], rendering each field's
/// type to the same surface string the Default flavor stores ([FLAVOR-ML-TYPE]).
fn lower_variant(variant: MlVariant) -> TypeVariant {
    TypeVariant {
        name: variant.name,
        fields: variant.fields.into_iter().map(lower_type_field).collect(),
    }
}

/// Lower one `field : type` line, with `constraint: None` (ML has no `where`
/// clause on type fields yet) — byte-identical to the Default field shape.
fn lower_type_field(field: MlTypeField) -> TypeField {
    TypeField {
        name: field.name,
        ty: render_type(&field.ty),
        constraint: None,
    }
}

/// Render an [`MlType`] to the surface type string the Default flavor stores in
/// [`TypeField::ty`]: a bare name as itself, an application as `Head<a, b>`.
fn render_type(ty: &MlType) -> String {
    match ty {
        MlType::Name(name) => name.clone(),
        MlType::App { head, args } => {
            let rendered = args.iter().map(render_type).collect::<Vec<_>>().join(", ");
            format!("{head}<{rendered}>")
        }
        MlType::Arrow { from, to } => format!("{} -> {}", render_type(from), render_type(to)),
        MlType::Tuple(parts) => {
            let rendered = parts.iter().map(render_type).collect::<Vec<_>>().join(", ");
            format!("({rendered})")
        }
    }
}

/// A binding with no parameters is a `let` (its signature becomes the binding's
/// type); one with parameters is a function ([FLAVOR-CURRY]). The unit marker
/// `()` yields a zero-parameter function, matching `fn f() = …`.
fn lower_binding(
    mutable: bool,
    name: String,
    params: Vec<MlParam>,
    body: MlExpr,
    pos: Position,
    sig: Option<&MlType>,
) -> Stmt {
    let body = lower_expr(body);
    // An empty surface parameter list is a value binding; a non-empty one (even
    // the lone unit marker `()`) is a function. `()` binds no canonical
    // parameter, so `f () = e` is a zero-parameter function like `fn f() = e`.
    if params.is_empty() {
        return Stmt::Let {
            name,
            mutable,
            ty: sig.and_then(type_expr),
            value: body,
            doc: None,
            position: Some(pos),
        };
    }
    let (parameters, body, return_type) = build_function(params, body, sig, pos);
    Stmt::Function {
        name,
        parameters,
        return_type,
        effects: Vec::new(),
        body,
        doc: None,
        position: Some(pos),
    }
}

/// Build an **uncurried** function, threading the signature's argument/result
/// types through positionally so the result is byte-identical to the Default
/// flavor's `fn f(a, b) -> r` ([FLAVOR-CURRY]). Every named surface parameter
/// becomes a real canonical parameter; the unit marker `()` consumes one
/// signature slot but binds no parameter (so `f () = e : Unit -> int` is a
/// zero-parameter function returning `int`). The return type is the arrow tail
/// left after all surface parameters are consumed.
fn build_function(
    params: Vec<MlParam>,
    body: Expr,
    sig: Option<&MlType>,
    pos: Position,
) -> (Vec<Parameter>, Expr, Option<TypeExpr>) {
    let _ = pos;
    let spine = sig.map(arrow_spine).unwrap_or_default();
    let consumed = params.len();
    let parameters = params
        .into_iter()
        .enumerate()
        .filter_map(|(i, p)| match p {
            MlParam::Named(name) => Some(Parameter {
                name,
                ty: spine.get(i).and_then(type_expr),
            }),
            MlParam::Unit => None,
        })
        .collect();
    let return_type = arrow_of(spine.get(consumed..).unwrap_or(&[]));
    (parameters, body, return_type)
}

/// Flatten the top-level arrow spine of a type: `a -> b -> c` ⇒ `[a, b, c]`,
/// `(a, b) -> c` ⇒ `[(a,b), c]`, a non-arrow ⇒ a single-element list.
fn arrow_spine(ty: &MlType) -> Vec<MlType> {
    match ty {
        MlType::Arrow { from, to } => {
            let mut spine = vec![(**from).clone()];
            spine.extend(arrow_spine(to));
            spine
        }
        other => vec![other.clone()],
    }
}

/// Rebuild a right-associative function type from an arrow-spine slice: `[]` ⇒
/// no type, `[t]` ⇒ `t`, `[a, b, …]` ⇒ `a -> (b -> …)`.
fn arrow_of(slice: &[MlType]) -> Option<TypeExpr> {
    match slice {
        [] => None,
        [single] => type_expr(single),
        [first, rest @ ..] => Some(TypeExpr {
            name: "fn".to_owned(),
            generic_params: Vec::new(),
            is_array: false,
            array_element: None,
            is_function: true,
            parameter_types: vec![type_expr(first)?],
            return_type: Some(Box::new(arrow_of(rest)?)),
            position: None,
        }),
    }
}

/// Convert an ML type to a canonical [`TypeExpr`]. A tuple type has no canonical
/// `TypeExpr` form, so it (and anything containing one) yields `None` — leaving
/// that position to inference rather than annotating it wrongly.
fn type_expr(ty: &MlType) -> Option<TypeExpr> {
    match ty {
        MlType::Name(name) => Some(TypeExpr::named(name.clone())),
        MlType::App { head, args } => {
            let generic_params = args.iter().map(type_expr).collect::<Option<Vec<_>>>()?;
            Some(TypeExpr {
                generic_params,
                ..TypeExpr::named(head.clone())
            })
        }
        MlType::Arrow { from, to } => Some(TypeExpr {
            name: "fn".to_owned(),
            generic_params: Vec::new(),
            is_array: false,
            array_element: None,
            is_function: true,
            parameter_types: vec![type_expr(from)?],
            return_type: Some(Box::new(type_expr(to)?)),
            position: None,
        }),
        MlType::Tuple(_) => None,
    }
}

/// Keep only the named parameters; the unit marker `()` carries no name and so
/// contributes nothing (a `Unit -> T` function has zero canonical parameters).
fn named_params(params: Vec<MlParam>) -> Vec<Parameter> {
    params
        .into_iter()
        .filter_map(|p| match p {
            MlParam::Named(name) => Some(Parameter { name, ty: None }),
            MlParam::Unit => None,
        })
        .collect()
}

/// Lower one CST expression to a canonical [`Expr`].
fn lower_expr(expr: MlExpr) -> Expr {
    match expr {
        MlExpr::Int(n) => Expr::Integer(n),
        MlExpr::Float(f) => Expr::Float(f),
        MlExpr::Bool(b) => Expr::Bool(b),
        MlExpr::Str(raw) => lower_string(&raw),
        MlExpr::Ident(name) => Expr::Identifier(name),
        MlExpr::Paren(inner) => lower_expr(*inner),
        MlExpr::Unary { op, operand } => Expr::Unary {
            op,
            operand: Box::new(lower_expr(*operand)),
        },
        MlExpr::Binary { op, left, right } => lower_binary(&op, *left, *right),
        MlExpr::App { func, arg } => lower_application(*func, *arg),
        MlExpr::UnitApp { func } => call(lower_expr(*func), Vec::new()),
        MlExpr::List(items) => Expr::List(items.into_iter().map(lower_expr).collect()),
        MlExpr::Index { target, index } => Expr::Index {
            target: Box::new(lower_expr(*target)),
            index: Box::new(lower_expr(*index)),
        },
        MlExpr::Field { target, name } => Expr::FieldAccess {
            target: Box::new(lower_expr(*target)),
            field: name,
        },
        MlExpr::Lambda { params, body, pos } => Expr::Lambda {
            parameters: named_params(params),
            return_type: None,
            body: Box::new(lower_expr(*body)),
            position: Some(pos),
        },
        MlExpr::Match { scrutinee, arms } => Expr::Match {
            value: Box::new(lower_expr(*scrutinee)),
            arms: arms.into_iter().map(lower_arm).collect(),
        },
        MlExpr::Record { name, fields } => Expr::TypeConstructor {
            name,
            type_args: Vec::new(),
            fields: fields.into_iter().map(lower_field).collect(),
        },
        MlExpr::Block { items, value } => lower_block(items, value),
        MlExpr::Spawn(body) => Expr::Spawn(Box::new(lower_expr(*body))),
    }
}

/// `|>` desugars to a call (the pipe is invisible downstream); every other
/// operator is a canonical [`Expr::Binary`].
fn lower_binary(op: &str, left: MlExpr, right: MlExpr) -> Expr {
    let left = lower_expr(left);
    let right = lower_expr(right);
    if op == "|>" {
        return pipe_into(left, right);
    }
    Expr::Binary {
        op: op.to_owned(),
        left: Box::new(left),
        right: Box::new(right),
    }
}

/// A block lowers to [`Expr::Block`]; a block that is a single trailing value
/// with no statements unwraps to that value, so it is structurally identical to
/// the Default inline body.
fn lower_block(items: Vec<MlItem>, value: Option<Box<MlExpr>>) -> Expr {
    let statements = lower_items(items);
    let value = value.map(|v| Box::new(lower_expr(*v)));
    match (statements.is_empty(), value) {
        (true, Some(value)) => *value,
        (_, value) => Expr::Block { statements, value },
    }
}

fn lower_arm(arm: MlArm) -> MatchArm {
    MatchArm {
        pattern: lower_pattern(arm.pattern),
        body: lower_expr(arm.body),
    }
}

fn lower_pattern(pattern: MlPattern) -> Pattern {
    match pattern {
        MlPattern::Wildcard => Pattern::Wildcard,
        MlPattern::Int(n) => Pattern::Literal(Box::new(Expr::Integer(n))),
        MlPattern::Str(raw) => Pattern::Literal(Box::new(lower_string(&raw))),
        MlPattern::Bool(b) => Pattern::Literal(Box::new(Expr::Bool(b))),
        MlPattern::Bind(name) => Pattern::Binding(name),
        MlPattern::Ctor { name, fields } => Pattern::Constructor {
            name,
            fields,
            sub_patterns: Vec::new(),
        },
    }
}

fn lower_field(field: MlField) -> FieldAssignment {
    FieldAssignment {
        name: field.name,
        value: lower_expr(field.value),
    }
}

/// `x |> f a` → `f x a`: prepend the piped value as the first argument of the
/// right-hand call, or wrap a bare callee in a one-argument call.
fn pipe_into(left: Expr, right: Expr) -> Expr {
    match right {
        Expr::Call {
            function,
            mut arguments,
            named_arguments,
        } => {
            arguments.insert(0, left);
            Expr::Call {
                function,
                arguments,
                named_arguments,
            }
        }
        callee => call(callee, vec![left]),
    }
}

/// Build a single positional [`Expr::Call`] node.
fn call(function: Expr, arguments: Vec<Expr>) -> Expr {
    Expr::Call {
        function: Box::new(function),
        arguments,
        named_arguments: Vec::new(),
    }
}

/// Collapse a whitespace-application spine `((f a) b) c` into one multi-argument
/// [`Expr::Call`] `f(a, b, c)` ([FLAVOR-CURRY]) — the uncurried skin: byte-for-
/// byte the AST the Default flavor emits for `f(a, b, c)`. `func`/`arg` are the
/// outermost `App`'s halves; nested `App` heads are peeled left-to-right.
fn lower_application(func: MlExpr, arg: MlExpr) -> Expr {
    let mut args = vec![lower_expr(arg)];
    let mut head = func;
    while let MlExpr::App { func, arg } = head {
        args.push(lower_expr(*arg));
        head = *func;
    }
    args.reverse();
    call(lower_expr(head), args)
}

/// Lower a raw string token to a plain or interpolated string expression,
/// reusing the Default frontend's escape/`${…}` handling with an ML fragment
/// parser ([FLAVOR-FRONTEND]).
fn lower_string(raw: &str) -> Expr {
    if raw.contains("${") {
        Expr::InterpolatedStr(lower_interpolation(raw, parse_fragment))
    } else {
        Expr::Str(unquote(raw))
    }
}

/// Parse a `${…}` fragment as an ML expression (`${toString id}` is ML
/// application), threading the flavor through interpolation re-entry.
fn parse_fragment(frag: &str) -> Expr {
    let (items, _) = super::parser::parse(&format!("__frag__ = {frag}\n"));
    match items.into_iter().next() {
        Some(MlItem::Binding { body, .. }) => lower_expr(body),
        _ => Expr::Identifier(frag.trim().to_owned()),
    }
}

#[cfg(test)]
#[expect(
    clippy::indexing_slicing,
    reason = "test assertions: an out-of-bounds index is a test failure, not a panic"
)]
mod tests {
    use super::super::parse_ml;
    use osprey_ast::{Expr, InterpolatedPart, Pattern, Stmt};

    fn stmts(src: &str) -> Vec<Stmt> {
        let parsed = parse_ml(src);
        assert!(parsed.errors.is_empty(), "ml errors: {:?}", parsed.errors);
        parsed.program.statements
    }

    fn one(src: &str) -> Stmt {
        let mut s = stmts(src);
        assert_eq!(s.len(), 1, "expected exactly one statement: {s:?}");
        s.pop().unwrap()
    }

    #[test]
    fn value_binding_lowers_to_let() {
        match one("answer = 42\n") {
            Stmt::Let {
                name,
                mutable,
                value,
                ..
            } => {
                assert_eq!(name, "answer");
                assert!(!mutable);
                assert_eq!(value, Expr::Integer(42));
            }
            s => panic!("expected let, got {s:?}"),
        }
    }

    #[test]
    fn mut_and_assignment_lower_distinctly() {
        let s = stmts("mut requests = 0\nrequests := requests + 1\n");
        assert!(matches!(s[0], Stmt::Let { mutable: true, .. }));
        assert!(matches!(s[1], Stmt::Assignment { ref name, .. } if name == "requests"));
    }

    #[test]
    fn multi_param_function_is_one_uncurried_function() {
        // `add x y = x + y` ⇒ a single TWO-parameter function with the operator
        // straight in its body — byte-identical to the Default `fn add(x, y)`
        // ([FLAVOR-CURRY], the uncurried skin), deliberately NOT a one-parameter
        // head returning a lambda.
        match one("add x y = x + y\n") {
            Stmt::Function {
                name,
                parameters,
                body,
                ..
            } => {
                assert_eq!(name, "add");
                assert_eq!(parameters.len(), 2);
                assert_eq!(parameters[0].name, "x");
                assert_eq!(parameters[1].name, "y");
                assert!(matches!(body, Expr::Binary { ref op, .. } if op == "+"));
            }
            s => panic!("expected function, got {s:?}"),
        }
    }

    #[test]
    fn single_param_function_has_no_extra_lambda() {
        match one("inc x = x + 1\n") {
            Stmt::Function {
                parameters, body, ..
            } => {
                assert_eq!(parameters.len(), 1);
                assert!(matches!(body, Expr::Binary { .. }));
            }
            s => panic!("expected function, got {s:?}"),
        }
    }

    #[test]
    fn unit_function_has_zero_parameters() {
        // `f () = body` is a zero-parameter function, like the Default `fn f()`.
        match one("greet () = 1\n") {
            Stmt::Function { parameters, .. } => assert!(parameters.is_empty()),
            s => panic!("expected function, got {s:?}"),
        }
    }

    #[test]
    fn whitespace_application_is_one_multi_arg_call() {
        // `add 1 2` ⇒ Call(add, [1, 2]) — the spine collapses to one
        // multi-argument call, byte-identical to the Default `add(1, 2)`.
        match one("r = add 1 2\n") {
            Stmt::Let {
                value:
                    Expr::Call {
                        function,
                        arguments,
                        ..
                    },
                ..
            } => {
                assert_eq!(*function, Expr::Identifier("add".to_owned()));
                assert_eq!(arguments, vec![Expr::Integer(1), Expr::Integer(2)]);
            }
            s => panic!("expected let-call, got {s:?}"),
        }
    }

    #[test]
    fn application_binds_tighter_than_operators() {
        // `add 1 2 == 3` ⇒ (add 1 2) == 3.
        match one("r = add 1 2 == 3\n") {
            Stmt::Let {
                value: Expr::Binary { op, left, right },
                ..
            } => {
                assert_eq!(op, "==");
                assert!(matches!(*left, Expr::Call { .. }));
                assert_eq!(*right, Expr::Integer(3));
            }
            s => panic!("expected comparison, got {s:?}"),
        }
    }

    #[test]
    fn unit_application_is_zero_arg_call() {
        match one("r = make ()\n") {
            Stmt::Let {
                value:
                    Expr::Call {
                        function,
                        arguments,
                        ..
                    },
                ..
            } => {
                assert!(arguments.is_empty());
                assert_eq!(*function, Expr::Identifier("make".to_owned()));
            }
            s => panic!("expected zero-arg call, got {s:?}"),
        }
    }

    #[test]
    fn match_lowers_constructor_and_wildcard_arms() {
        match one("r =\n    match x\n        Success value => value\n        _ => 0\n") {
            Stmt::Let {
                value: Expr::Match { arms, .. },
                ..
            } => {
                assert_eq!(arms.len(), 2);
                match &arms[0].pattern {
                    Pattern::Constructor { name, fields, .. } => {
                        assert_eq!(name, "Success");
                        assert_eq!(fields, &vec!["value".to_owned()]);
                    }
                    p => panic!("expected constructor pattern, got {p:?}"),
                }
                assert!(matches!(arms[1].pattern, Pattern::Wildcard));
            }
            s => panic!("expected match, got {s:?}"),
        }
    }

    #[test]
    fn lambda_is_multi_param_and_pipe_desugars_to_call() {
        // `\x y => x + y` ⇒ one two-parameter lambda, byte-identical to the
        // Default `(x, y) => x + y` — not a nested one-parameter chain.
        match one("f = \\x y => x + y\n") {
            Stmt::Let {
                value: Expr::Lambda {
                    parameters, body, ..
                },
                ..
            } => {
                assert_eq!(parameters.len(), 2);
                assert_eq!(parameters[0].name, "x");
                assert_eq!(parameters[1].name, "y");
                assert!(matches!(*body, Expr::Binary { ref op, .. } if op == "+"));
            }
            s => panic!("expected lambda, got {s:?}"),
        }
        // `x |> f` becomes `f(x)` — no Pipe node survives, matching Default.
        match one("r = x |> f\n") {
            Stmt::Let {
                value:
                    Expr::Call {
                        function,
                        arguments,
                        ..
                    },
                ..
            } => {
                assert_eq!(*function, Expr::Identifier("f".to_owned()));
                assert_eq!(arguments, vec![Expr::Identifier("x".to_owned())]);
            }
            s => panic!("expected piped call, got {s:?}"),
        }
    }

    #[test]
    fn record_block_lowers_to_type_constructor() {
        let src = "p =\n    Point\n        x = 1\n        y = 2\n";
        match one(src) {
            Stmt::Let {
                value: Expr::TypeConstructor { name, fields, .. },
                ..
            } => {
                assert_eq!(name, "Point");
                assert_eq!(fields.len(), 2);
                assert_eq!(fields[0].name, "x");
            }
            s => panic!("expected type constructor, got {s:?}"),
        }
    }

    #[test]
    fn inline_record_lowers_to_type_constructor() {
        // `Ok(value = "x")` in expression position is an inline record literal —
        // it lowers to the SAME `Expr::TypeConstructor` the layout form and the
        // Default `Ok { value: "x" }` produce ([FLAVOR-ML-RECORD]).
        match one("r = Ok(value = \"x\")\n") {
            Stmt::Let {
                value: Expr::TypeConstructor { name, fields, .. },
                ..
            } => {
                assert_eq!(name, "Ok");
                assert_eq!(fields.len(), 1);
                assert_eq!(fields[0].name, "value");
                assert_eq!(fields[0].value, Expr::Str("x".to_owned()));
            }
            s => panic!("expected type constructor, got {s:?}"),
        }
    }

    #[test]
    fn spawn_inline_expr_lowers_to_spawn() {
        // `spawn f x` lowers to `Expr::Spawn` wrapping the call, byte-identical
        // to the Default `spawn f(x)` ([FLAVOR-ML-SPAWN]).
        match one("r = spawn task 1\n") {
            Stmt::Let {
                value: Expr::Spawn(inner),
                ..
            } => assert!(
                matches!(*inner, Expr::Call { .. }),
                "spawn body should be the call, got {inner:?}"
            ),
            s => panic!("expected spawn, got {s:?}"),
        }
    }

    #[test]
    fn spawn_block_lowers_to_spawn_block() {
        // `spawn` + an indented block lowers to `Expr::Spawn` wrapping the block.
        match one("r = spawn\n    x = 1\n    task x\n") {
            Stmt::Let {
                value: Expr::Spawn(inner),
                ..
            } => assert!(
                matches!(*inner, Expr::Block { .. }),
                "spawn block body should be a Block, got {inner:?}"
            ),
            s => panic!("expected spawn, got {s:?}"),
        }
    }

    #[test]
    fn interpolation_parses_fragment_as_ml_application() {
        // `${toString id}` is ML whitespace application inside the fragment.
        match one("r = \"n=${toString id}\"\n") {
            Stmt::Let {
                value: Expr::InterpolatedStr(parts),
                ..
            } => {
                assert!(matches!(parts[0], InterpolatedPart::Text(ref t) if t == "n="));
                assert!(matches!(
                    parts[1],
                    InterpolatedPart::Expr(Expr::Call { .. })
                ));
            }
            s => panic!("expected interpolated string, got {s:?}"),
        }
    }

    #[test]
    fn block_body_with_statements_keeps_block_with_trailing_value() {
        let src = "f x =\n    y = x + 1\n    y + 2\n";
        match one(src) {
            Stmt::Function {
                body: Expr::Block { statements, value },
                ..
            } => {
                assert_eq!(statements.len(), 1);
                assert!(value.is_some());
            }
            s => panic!("expected block body, got {s:?}"),
        }
    }

    #[test]
    fn name_binding_is_an_immutable_let_at_top_level_and_in_a_block() {
        // The cross-flavor guarantee `name = expr` must satisfy: it lowers to the
        // SAME node a Default `let name = expr` does — `Stmt::Let { mutable:
        // false }` — both at the top level and inside a layout block. This is the
        // structural precondition for byte-identical IR with the Default twin
        // ([FLAVOR-CURRY], [FLAVOR-IR-EQUIV]); only `mut`+`:=` produces an
        // `Assignment`, never a bare `=`.
        assert!(matches!(
            one("answer = 41 + 1\n"),
            Stmt::Let { mutable: false, .. }
        ));
        // Same binding, this time the first statement of a function block.
        match one("main () =\n    answer = 41 + 1\n    answer\n") {
            Stmt::Function {
                body: Expr::Block { statements, .. },
                ..
            } => assert!(
                matches!(statements.first(), Some(Stmt::Let { mutable: false, name, .. }) if name == "answer"),
                "block-local `name = expr` must be an immutable Let, got {statements:?}"
            ),
            s => panic!("expected function with block body, got {s:?}"),
        }
    }

    #[test]
    fn union_type_lowers_to_canonical_type_stmt() {
        // The ML layout union must lower to the SAME `Stmt::Type` the Default
        // `type Outcome = Ok { value: string } | Err { message: string }` emits:
        // two payload-carrying variants with `validation_func: None` and each
        // field `constraint: None` ([FLAVOR-ML-TYPE], [FLAVOR-IR-EQUIV]).
        let src = "type Outcome =\n    Ok\n        value : string\n    Err\n        message : string\n";
        match one(src) {
            Stmt::Type {
                name,
                type_params,
                variants,
                validation_func,
                ..
            } => {
                assert_eq!(name, "Outcome");
                assert!(type_params.is_empty());
                assert!(validation_func.is_none());
                assert_eq!(variants.len(), 2);
                assert_eq!(variants[0].name, "Ok");
                assert_eq!(variants[0].fields.len(), 1);
                assert_eq!(variants[0].fields[0].name, "value");
                assert_eq!(variants[0].fields[0].ty, "string");
                assert!(variants[0].fields[0].constraint.is_none());
                assert_eq!(variants[1].name, "Err");
                assert_eq!(variants[1].fields[0].name, "message");
            }
            s => panic!("expected type, got {s:?}"),
        }
    }

    #[test]
    fn enum_type_lowers_to_fieldless_variants() {
        match one("type Status =\n    Active\n    Inactive\n") {
            Stmt::Type { variants, .. } => {
                assert_eq!(variants.len(), 2);
                assert_eq!(variants[0].name, "Active");
                assert!(variants[0].fields.is_empty());
                assert_eq!(variants[1].name, "Inactive");
                assert!(variants[1].fields.is_empty());
            }
            s => panic!("expected type, got {s:?}"),
        }
    }

    #[test]
    fn record_type_lowers_to_single_variant_named_after_type() {
        // A lowercase first field marks the record form; its lone variant takes
        // the type's own name, exactly as Default's `type Point = { x, y }` does.
        match one("type Point =\n    x : int\n    y : int\n") {
            Stmt::Type {
                name, variants, ..
            } => {
                assert_eq!(name, "Point");
                assert_eq!(variants.len(), 1);
                assert_eq!(variants[0].name, "Point");
                assert_eq!(variants[0].fields.len(), 2);
                assert_eq!(variants[0].fields[0].name, "x");
                assert_eq!(variants[0].fields[0].ty, "int");
            }
            s => panic!("expected type, got {s:?}"),
        }
    }

    #[test]
    fn extern_lowers_to_canonical_extern_stmt() {
        // `extern name (p : T) (q : U) -> R` lowers to the SAME `Stmt::Extern`
        // the Default `extern fn name(p: T, q: U) -> R` emits — typed parameters
        // in order plus a return type ([FLAVOR-ML-EXTERN], [FLAVOR-IR-EQUIV]).
        match one("extern sqlite3_open (filename : string) (ppDb : Ptr) -> int\n") {
            Stmt::Extern {
                name,
                parameters,
                return_type,
                ..
            } => {
                assert_eq!(name, "sqlite3_open");
                assert_eq!(parameters.len(), 2);
                assert_eq!(parameters[0].name, "filename");
                assert_eq!(parameters[0].ty.name, "string");
                assert_eq!(parameters[1].name, "ppDb");
                assert_eq!(parameters[1].ty.name, "Ptr");
                assert_eq!(return_type.map(|t| t.name), Some("int".to_owned()));
            }
            s => panic!("expected extern, got {s:?}"),
        }
    }

    #[test]
    fn reserved_effect_word_reports_a_clear_error() {
        let parsed = parse_ml("effect Db\n    add : string => int\n");
        assert!(parsed
            .errors
            .iter()
            .any(|e| e.message.contains("not yet supported")));
    }
}
