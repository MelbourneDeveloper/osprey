//! Expression lowering (port of `builder_expressions.go` / `builder_calls.go` /
//! `builder_interpolation.go`).

use crate::lower::Lowerer;
use osprey_ast::*;
use tree_sitter::Node;

impl<'a> Lowerer<'a> {
    pub(crate) fn lower_expr_field(&self, node: Node, field: &str) -> Expr {
        match node.child_by_field_name(field) {
            Some(n) => self.lower_expr(n),
            None => Expr::Bool(false), // unreachable for well-formed trees
        }
    }

    /// Lower an expression node, transparently unwrapping the `expression` and
    /// `primary_expression` wrapper nodes tree-sitter inserts.
    pub(crate) fn lower_expr(&self, node: Node) -> Expr {
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
            "pipe_expression" => Expr::Pipe {
                left: Box::new(self.lower_expr_field(node, "left")),
                right: Box::new(self.lower_expr_field(node, "right")),
            },
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
                let exprs: Vec<Node> = node
                    .named_children(&mut cursor)
                    .filter(|c| c.kind() == "expression")
                    .collect();
                Expr::Send {
                    channel: Box::new(self.lower_expr(exprs[0])),
                    value: Box::new(self.lower_expr(exprs[1])),
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
            "ternary_expression" => {
                // Lower the condition; full ternary semantics handled in a later pass.
                self.lower_expr_field(node, "condition")
            }
            _ => Expr::Bool(false),
        }
    }

    fn lower_inner_expr(&self, node: Node) -> Expr {
        match self.first_named(node) {
            Some(n) => self.lower_expr(n),
            None => Expr::Bool(false),
        }
    }

    fn lower_call(&self, node: Node) -> Expr {
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
        // function/method call
        let (arguments, named_arguments) = self.lower_arg_list(node);
        match callee {
            Expr::FieldAccess { target, field } => Expr::MethodCall {
                target,
                method: field,
                arguments,
                named_arguments,
            },
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
    fn lower_arg_list(&self, node: Node) -> (Vec<Expr>, Vec<NamedArgument>) {
        let Some(list) = self.named_of_kind(node, "argument_list").into_iter().next() else {
            return (Vec::new(), Vec::new());
        };
        if let Some(nal) = self.first_child_of_kind(list, "named_argument_list") {
            let named = self
                .named_of_kind(nal, "named_argument")
                .iter()
                .map(|na| NamedArgument {
                    name: self.field_text(*na, "name"),
                    value: self.lower_expr_field(*na, "value"),
                })
                .collect();
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

    fn lower_arms(&self, node: Node) -> Vec<MatchArm> {
        self.named_of_kind(node, "match_arm")
            .iter()
            .chain(self.named_of_kind(node, "select_arm").iter())
            .map(|arm| MatchArm {
                pattern: self.lower_pattern(arm.child_by_field_name("pattern").unwrap()),
                body: self.lower_expr_field(*arm, "body"),
            })
            .collect()
    }

    fn lower_handler_arms(&self, node: Node) -> Vec<HandlerArm> {
        self.named_of_kind(node, "handler_arm")
            .iter()
            .map(|arm| HandlerArm {
                operation: self.field_text(*arm, "operation"),
                params: self
                    .first_child_of_kind(*arm, "handler_params")
                    .map(|hp| {
                        self.named_of_kind(hp, "identifier")
                            .iter()
                            .map(|n| self.text(*n))
                            .collect()
                    })
                    .unwrap_or_default(),
                body: self.lower_expr_field(*arm, "body"),
            })
            .collect()
    }

    fn lower_field_assignments(&self, node: Node) -> Vec<FieldAssignment> {
        self.descendants_of_kind(node, "field_assignment")
            .iter()
            .map(|fa| FieldAssignment {
                name: self.field_text(*fa, "name"),
                value: self.lower_expr_field(*fa, "value"),
            })
            .collect()
    }

    fn lower_block(&self, node: Node) -> Expr {
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

    pub(crate) fn lower_literal(&self, node: Node) -> Expr {
        let inner = match self.first_named(node) {
            Some(n) => n,
            None => return Expr::Bool(false),
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
                    Expr::InterpolatedStr(self.lower_interpolation(&raw))
                } else {
                    Expr::Str(unquote(&raw))
                }
            }
            "interpolated_string" => {
                Expr::InterpolatedStr(self.lower_interpolation(&self.text(inner)))
            }
            "list_literal" => Expr::List(
                self.named_of_kind(inner, "expression")
                    .iter()
                    .map(|e| self.lower_expr(*e))
                    .collect(),
            ),
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
    /// each embedded expression as an Osprey fragment (port of
    /// `builder_interpolation.go`).
    fn lower_interpolation(&self, raw: &str) -> Vec<InterpolatedPart> {
        let inner = unquote(raw);
        let mut parts = Vec::new();
        let mut text = String::new();
        let bytes = inner.as_bytes();
        let mut i = 0;
        while i < bytes.len() {
            if bytes[i] == b'$' && i + 1 < bytes.len() && bytes[i + 1] == b'{' {
                if !text.is_empty() {
                    parts.push(InterpolatedPart::Text(std::mem::take(&mut text)));
                }
                let mut j = i + 2;
                while j < bytes.len() && bytes[j] != b'}' {
                    j += 1;
                }
                let frag = &inner[i + 2..j];
                parts.push(InterpolatedPart::Expr(parse_fragment(frag)));
                i = j + 1;
            } else {
                text.push(bytes[i] as char);
                i += 1;
            }
        }
        if !text.is_empty() {
            parts.push(InterpolatedPart::Text(text));
        }
        parts
    }
}

/// Strip surrounding quotes and apply basic escape unescaping.
fn unquote(s: &str) -> String {
    let trimmed = s
        .strip_prefix('"')
        .and_then(|x| x.strip_suffix('"'))
        .unwrap_or(s);
    trimmed
        .replace("\\n", "\n")
        .replace("\\t", "\t")
        .replace("\\\"", "\"")
        .replace("\\\\", "\\")
}

/// Parse an interpolation fragment (`${ ... }` contents) into a single [`Expr`].
fn parse_fragment(frag: &str) -> Expr {
    let parsed = crate::parse_program(&format!("let __frag__ = {frag}\n"));
    match parsed.program.statements.into_iter().next() {
        Some(Stmt::Let { value, .. }) => value,
        _ => Expr::Identifier(frag.trim().to_string()),
    }
}
