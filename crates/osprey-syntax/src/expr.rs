//! Expression lowering: every [`Expr`] form — literals, operators, call shapes
//! and named arguments, match/handler arms, and string interpolation.

use crate::lower::Lowerer;
use osprey_ast::{
    Expr, FieldAssignment, HandlerArm, InterpolatedPart, MapEntry, MatchArm, NamedArgument,
    Pattern, Stmt,
};
use tree_sitter::Node;

impl Lowerer<'_> {
    pub(crate) fn lower_expr_field(&self, node: Node<'_>, field: &str) -> Expr {
        match node.child_by_field_name(field) {
            Some(n) => self.lower_expr(n),
            None => Expr::Bool(false), // unreachable for well-formed trees
        }
    }

    /// Lower an expression node, transparently unwrapping the `expression` and
    /// `primary_expression` wrapper nodes tree-sitter inserts.
    pub(crate) fn lower_expr(&self, node: Node<'_>) -> Expr {
        match node.kind() {
            "expression" | "primary_expression" => match self.first_named(node) {
                Some(inner) => self.lower_expr(inner),
                None => Expr::Bool(false),
            },
            "binary_expression" => Expr::Binary {
                op: self.field_text(node, "operator"),
                left: Box::new(self.lower_expr_field(node, "left")),
                right: Box::new(self.lower_expr_field(node, "right")),
            },
            "unary_expression" => Expr::Unary {
                op: self.field_text(node, "operator"),
                operand: Box::new(self.lower_expr_field(node, "operand")),
            },
            // `x |> f` desugars to `f(x)` and `x |> f(a, …)` to `f(x, a, …)` —
            // the piped value becomes the callee's first positional argument, so
            // both the type checker and codegen see an ordinary call.
            "pipe_expression" => pipe_into(
                self.lower_expr_field(node, "left"),
                self.lower_expr_field(node, "right"),
            ),
            "call_expression" => self.lower_call(node),
            "match_expression" => Expr::Match {
                value: Box::new(self.lower_expr_field(node, "value")),
                arms: self.lower_arms(node),
            },
            "select_expression" => Expr::Select {
                arms: self.lower_arms(node),
            },
            "handler_expression" => Expr::Handler {
                effect: self.field_text(node, "effect"),
                arms: self.lower_handler_arms(node),
                body: Box::new(self.lower_expr_field(node, "body")),
            },
            "perform_expression" => {
                let (arguments, named_arguments) = self.lower_arg_list(node);
                Expr::Perform {
                    effect: self.field_text(node, "effect"),
                    operation: self.field_text(node, "operation"),
                    arguments,
                    named_arguments,
                }
            }
            "spawn_expression" => Expr::Spawn(Box::new(self.lower_inner_expr(node))),
            "yield_expression" => {
                Expr::Yield(self.first_named(node).map(|n| Box::new(self.lower_expr(n))))
            }
            "await_call" => Expr::Await(Box::new(self.lower_inner_expr(node))),
            "recv_call" => Expr::Recv(Box::new(self.lower_inner_expr(node))),
            "send_call" => {
                let mut cursor = node.walk();
                let mut exprs = node
                    .named_children(&mut cursor)
                    .filter(|c| c.kind() == "expression");
                let channel = exprs
                    .next()
                    .map_or(Expr::Bool(false), |n| self.lower_expr(n));
                let value = exprs
                    .next()
                    .map_or(Expr::Bool(false), |n| self.lower_expr(n));
                Expr::Send {
                    channel: Box::new(channel),
                    value: Box::new(value),
                }
            }
            "lambda_expression" => Expr::Lambda {
                parameters: self.lower_params(node.child_by_field_name("parameters").or_else(
                    || {
                        self.named_of_kind(node, "parameter_list")
                            .into_iter()
                            .next()
                    },
                )),
                return_type: None,
                body: Box::new(self.lower_expr_field(node, "body")),
            },
            "type_constructor" => Expr::TypeConstructor {
                name: self.field_text(node, "name"),
                type_args: Vec::new(),
                fields: self.lower_field_assignments(node),
            },
            "update_expression" => Expr::Update {
                record: self.field_text(node, "record"),
                fields: self.lower_field_assignments(node),
            },
            "object_literal" => Expr::Object(self.lower_field_assignments(node)),
            "block" => self.lower_block(node),
            "literal" => self.lower_literal(node),
            "identifier" => Expr::Identifier(self.text(node)),
            "ternary_expression" => self.lower_ternary(node),
            _ => Expr::Bool(false),
        }
    }

    /// `cond ? then : else` desugars to `match cond { true => then  false => else }`
    /// (and the Elvis form `cond ?: else` reuses the condition as the `then`),
    /// so the existing boolean-match lowering carries the runtime semantics.
    fn lower_ternary(&self, node: Node<'_>) -> Expr {
        let condition = self.lower_expr_field(node, "condition");
        // Structural form `cond { f1, f2 } ? then : else`: bind each field from
        // `cond` and evaluate `then` — a record/object always carries its declared
        // fields, so the structural check succeeds. The then/else are the two
        // positional `expression` children after `condition`.
        if let Some(fp) = self.first_child_of_kind(node, "field_pattern") {
            let exprs = self.named_of_kind(node, "expression");
            let then_expr = exprs
                .get(1)
                .map_or(Expr::Bool(false), |n| self.lower_expr(*n));
            let statements = self
                .texts_of_kind(fp, "identifier")
                .into_iter()
                .map(|f| Stmt::Let {
                    name: f.clone(),
                    mutable: false,
                    ty: None,
                    value: Expr::FieldAccess {
                        target: Box::new(condition.clone()),
                        field: f,
                    },
                    position: None,
                })
                .collect();
            return Expr::Block {
                statements,
                value: Some(Box::new(then_expr)),
            };
        }
        let else_expr = self.lower_expr_field(node, "else");
        let then_expr = match node.child_by_field_name("then") {
            Some(n) => self.lower_expr(n),
            None => condition.clone(), // Elvis `?:`
        };
        Expr::Match {
            value: Box::new(condition),
            arms: vec![
                MatchArm {
                    pattern: Pattern::Literal(Box::new(Expr::Bool(true))),
                    body: then_expr,
                },
                MatchArm {
                    pattern: Pattern::Literal(Box::new(Expr::Bool(false))),
                    body: else_expr,
                },
            ],
        }
    }

    fn lower_inner_expr(&self, node: Node<'_>) -> Expr {
        match self.first_named(node) {
            Some(n) => self.lower_expr(n),
            None => Expr::Bool(false),
        }
    }

    fn lower_call(&self, node: Node<'_>) -> Expr {
        let callee = self.lower_expr_field(node, "callee");
        if let Some(member) = node.child_by_field_name("member") {
            return Expr::FieldAccess {
                target: Box::new(callee),
                field: self.text(member),
            };
        }
        if let Some(index) = node.child_by_field_name("index") {
            return Expr::Index {
                target: Box::new(callee),
                index: Box::new(self.lower_expr(index)),
            };
        }
        // function/method call. UFCS: `x.f(a, …)` is sugar for `f(x, a, …)`, so a
        // field-access callee lowers to an ordinary call with the receiver as the
        // first positional argument — keeping method calls invisible downstream.
        let (mut arguments, named_arguments) = self.lower_arg_list(node);
        match callee {
            Expr::FieldAccess { target, field } => {
                arguments.insert(0, *target);
                Expr::Call {
                    function: Box::new(Expr::Identifier(field)),
                    arguments,
                    named_arguments,
                }
            }
            other => Expr::Call {
                function: Box::new(other),
                arguments,
                named_arguments,
            },
        }
    }

    /// Collect positional + named args from an `argument_list` child. Named
    /// arguments live in a single direct `named_argument_list` child; a *direct*
    /// lookup is essential — descending would steal the named arguments of a
    /// nested call (`print(cc(c1: .., c2: ..))` must not hoist c1/c2 onto print).
    fn lower_arg_list(&self, node: Node<'_>) -> (Vec<Expr>, Vec<NamedArgument>) {
        let Some(list) = self.named_of_kind(node, "argument_list").into_iter().next() else {
            return (Vec::new(), Vec::new());
        };
        if let Some(nal) = self.first_child_of_kind(list, "named_argument_list") {
            let nodes = self.named_of_kind(nal, "named_argument");
            let named = self.lower_name_value(&nodes, |name, value| NamedArgument { name, value });
            return (Vec::new(), named);
        }
        let mut cursor = list.walk();
        let positional = list
            .named_children(&mut cursor)
            .filter(|c| c.kind() == "expression")
            .map(|c| self.lower_expr(c))
            .collect();
        (positional, Vec::new())
    }

    fn lower_arms(&self, node: Node<'_>) -> Vec<MatchArm> {
        self.named_of_kind(node, "match_arm")
            .iter()
            .chain(self.named_of_kind(node, "select_arm").iter())
            .map(|arm| MatchArm {
                pattern: arm
                    .child_by_field_name("pattern")
                    .map_or(Pattern::Wildcard, |p| self.lower_pattern(p)),
                body: self.lower_expr_field(*arm, "body"),
            })
            .collect()
    }

    fn lower_handler_arms(&self, node: Node<'_>) -> Vec<HandlerArm> {
        self.named_of_kind(node, "handler_arm")
            .iter()
            .map(|arm| HandlerArm {
                operation: self.field_text(*arm, "operation"),
                params: self
                    .first_child_of_kind(*arm, "handler_params")
                    .map(|hp| self.texts_of_kind(hp, "identifier"))
                    .unwrap_or_default(),
                body: self.lower_expr_field(*arm, "body"),
            })
            .collect()
    }

    fn lower_field_assignments(&self, node: Node<'_>) -> Vec<FieldAssignment> {
        let nodes = self.descendants_of_kind(node, "field_assignment");
        self.lower_name_value(&nodes, |name, value| FieldAssignment { name, value })
    }

    /// Map each node carrying a `name` field and a `value` expression into an AST
    /// node built by `ctor` — the shared shape of named arguments and record
    /// field assignments.
    fn lower_name_value<T>(&self, nodes: &[Node<'_>], ctor: impl Fn(String, Expr) -> T) -> Vec<T> {
        nodes
            .iter()
            .map(|n| {
                ctor(
                    self.field_text(*n, "name"),
                    self.lower_expr_field(*n, "value"),
                )
            })
            .collect()
    }

    fn lower_block(&self, node: Node<'_>) -> Expr {
        let mut statements = Vec::new();
        let mut value = None;
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            match child.kind() {
                "statement" => {
                    if let Some(s) = self.first_named(child).and_then(|n| self.lower_stmt(n)) {
                        statements.push(s);
                    }
                }
                "expression" => value = Some(Box::new(self.lower_expr(child))),
                _ => {}
            }
        }
        // A block evaluates to its last expression. The grammar sometimes emits
        // that trailing expression as an `expression_statement`; recover it as
        // the block value so the type of `{ ...; r }` is the type of `r`.
        if value.is_none() {
            if let Some(Stmt::Expr(_)) = statements.last() {
                if let Some(Stmt::Expr(e)) = statements.pop() {
                    value = Some(Box::new(e));
                }
            }
        }
        Expr::Block { statements, value }
    }

    pub(crate) fn lower_literal(&self, node: Node<'_>) -> Expr {
        let Some(inner) = self.first_named(node) else {
            return Expr::Bool(false);
        };
        match inner.kind() {
            "integer" => Expr::Integer(self.text(inner).parse().unwrap_or(0)),
            "float" => Expr::Float(self.text(inner).parse().unwrap_or(0.0)),
            "boolean" => Expr::Bool(self.text(inner) == "true"),
            "string" => {
                // Equal-length `string`/`interpolated_string` token matches let the
                // plain `string` rule win, so a `"...${e}..."` can arrive tagged as
                // `string`. Detect the `${` marker here and interpolate either way.
                let raw = self.text(inner);
                if raw.contains("${") {
                    Expr::InterpolatedStr(Self::lower_interpolation(&raw))
                } else {
                    Expr::Str(unquote(&raw))
                }
            }
            "interpolated_string" => {
                Expr::InterpolatedStr(Self::lower_interpolation(&self.text(inner)))
            }
            "list_literal" => Expr::List(self.exprs_of_kind(inner, "expression")),
            "map_literal" => Expr::Map(
                self.named_of_kind(inner, "map_entry")
                    .iter()
                    .map(|me| MapEntry {
                        key: self.lower_expr_field(*me, "key"),
                        value: self.lower_expr_field(*me, "value"),
                    })
                    .collect(),
            ),
            _ => Expr::Bool(false),
        }
    }

    /// Split a `"text ${expr} more"` literal into [`InterpolatedPart`]s, parsing
    /// each embedded expression as an Osprey fragment.
    fn lower_interpolation(raw: &str) -> Vec<InterpolatedPart> {
        let inner = unquote(raw);
        let bytes = inner.as_bytes();
        let mut parts = Vec::new();
        let mut text_start = 0usize;
        let mut i = 0usize;
        while i < bytes.len() {
            if bytes.get(i) == Some(&b'$') && bytes.get(i + 1) == Some(&b'{') {
                if i > text_start {
                    if let Some(text) = inner.get(text_start..i) {
                        parts.push(InterpolatedPart::Text(text.to_string()));
                    }
                }
                // Find the `}` that closes this `${`, honouring nested braces so
                // `${match x { a => 1 b => 2 }}` captures the whole match.
                let mut depth = 1i32;
                let mut j = i + 2;
                while let Some(byte) = bytes.get(j) {
                    match byte {
                        b'{' => depth += 1,
                        b'}' => {
                            depth -= 1;
                            if depth == 0 {
                                break;
                            }
                        }
                        _ => {}
                    }
                    j += 1;
                }
                if let Some(frag) = inner.get(i + 2..j) {
                    parts.push(InterpolatedPart::Expr(parse_fragment(frag)));
                }
                i = j + 1;
                text_start = i;
            } else {
                i += 1;
            }
        }
        if let Some(text) = inner.get(text_start..) {
            if !text.is_empty() {
                parts.push(InterpolatedPart::Text(text.to_string()));
            }
        }
        parts
    }
}

/// Strip surrounding quotes and resolve backslash escapes in one pass (so a
/// literal `\\` can never be re-interpreted): `\n` `\r` `\t` newline/CR/tab,
/// `\e` the ANSI ESC (0x1B, used by the terminal-color helpers), `\0` NUL,
/// `\"` and `\\` the literals. An unrecognised escape is kept verbatim.
fn unquote(s: &str) -> String {
    let trimmed = s
        .strip_prefix('"')
        .and_then(|x| x.strip_suffix('"'))
        .unwrap_or(s);
    let mut out = String::with_capacity(trimmed.len());
    let mut chars = trimmed.chars();
    while let Some(c) = chars.next() {
        if c != '\\' {
            out.push(c);
            continue;
        }
        match chars.next() {
            Some('n') => out.push('\n'),
            Some('r') => out.push('\r'),
            Some('t') => out.push('\t'),
            Some('e') => out.push('\u{1b}'),
            Some('0') => out.push('\0'),
            Some('"') => out.push('"'),
            // An escaped backslash, or a trailing lone backslash at end of input.
            Some('\\') | None => out.push('\\'),
            Some(other) => {
                out.push('\\');
                out.push(other);
            }
        }
    }
    out
}

/// Parse an interpolation fragment (`${ ... }` contents) into a single [`Expr`].
fn parse_fragment(frag: &str) -> Expr {
    let parsed = crate::parse_program(&format!("let __frag__ = {frag}\n"));
    match parsed.program.statements.into_iter().next() {
        Some(Stmt::Let { value, .. }) => value,
        _ => Expr::Identifier(frag.trim().to_string()),
    }
}

/// Fold a piped value into its right-hand callee: `x |> f(a, …)` becomes
/// `f(x, a, …)` (the piped value is prepended as the first positional
/// argument). A bare callee `x |> f` becomes `f(x)`. Producing a plain
/// [`Expr::Call`] keeps pipes invisible to every later stage.
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
        callee => Expr::Call {
            function: Box::new(callee),
            arguments: vec![left],
            named_arguments: Vec::new(),
        },
    }
}
