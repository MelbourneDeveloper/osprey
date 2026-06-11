//! Statement, type, and pattern lowering: declarations (`fn`, `let`, `type`,
//! `effect`, `extern`, `module`), type expressions, and match patterns.

use osprey_ast::{
    EffectOperation, Expr, ExternParameter, Parameter, Pattern, Position, Program, Stmt, TypeExpr,
    TypeField, TypeVariant,
};
use tree_sitter::Node;

/// Holds the source bytes so node text can be sliced during lowering.
#[derive(Debug)]
pub struct Lowerer<'a> {
    src: &'a [u8],
}

impl<'a> Lowerer<'a> {
    /// Creates a lowerer over the given source bytes.
    #[must_use]
    pub fn new(src: &'a [u8]) -> Self {
        Lowerer { src }
    }

    pub(crate) fn text(&self, node: Node<'_>) -> String {
        node.utf8_text(self.src).unwrap_or("").to_string()
    }

    #[expect(
        clippy::unused_self,
        reason = "kept for Lowerer method-call ergonomics"
    )]
    pub(crate) fn pos(&self, node: Node<'_>) -> Position {
        let p = node.start_position();
        Position {
            line: u32::try_from(p.row).unwrap_or(u32::MAX).saturating_add(1),
            column: u32::try_from(p.column).unwrap_or(u32::MAX),
        }
    }

    /// First *named* child (skips anonymous tokens). Used to unwrap the wrapper
    /// nodes tree-sitter inserts (`statement`, `expression`, `primary_expression`).
    #[expect(
        clippy::unused_self,
        reason = "kept for Lowerer method-call ergonomics"
    )]
    pub(crate) fn first_named<'t>(&self, node: Node<'t>) -> Option<Node<'t>> {
        let mut cursor = node.walk();
        let found = node.named_children(&mut cursor).next();
        found
    }

    /// First named child of a given kind.
    #[expect(
        clippy::unused_self,
        reason = "kept for Lowerer method-call ergonomics"
    )]
    pub(crate) fn first_child_of_kind<'t>(&self, node: Node<'t>, kind: &str) -> Option<Node<'t>> {
        let mut cursor = node.walk();
        let found = node.named_children(&mut cursor).find(|c| c.kind() == kind);
        found
    }

    /// Lowers the root `source_file` node into a full program AST.
    #[must_use]
    pub fn lower_program(&self, root: Node<'_>) -> Program {
        let mut statements = Vec::new();
        let mut cursor = root.walk();
        for child in root.named_children(&mut cursor) {
            if child.kind() == "statement" {
                if let Some(stmt) = self.first_named(child).and_then(|n| self.lower_stmt(n)) {
                    statements.push(stmt);
                }
            }
        }
        Program { statements }
    }

    pub(crate) fn lower_stmt(&self, node: Node<'_>) -> Option<Stmt> {
        Some(match node.kind() {
            "import_statement" => Stmt::Import {
                module: self.texts_of_kind(node, "identifier"),
            },
            "let_declaration" => Stmt::Let {
                name: self.field_text(node, "name"),
                mutable: node
                    .child_by_field_name("keyword")
                    .is_some_and(|n| self.text(n) == "mut"),
                ty: node.child_by_field_name("type").map(|n| self.lower_type(n)),
                value: self.lower_expr_field(node, "value"),
                position: Some(self.pos(node)),
            },
            "assignment" => Stmt::Assignment {
                name: self.field_text(node, "name"),
                value: self.lower_expr_field(node, "value"),
                position: Some(self.pos(node)),
            },
            "function_declaration" => Stmt::Function {
                name: self.field_text(node, "name"),
                parameters: self.lower_params(node.child_by_field_name("parameters")),
                return_type: node
                    .child_by_field_name("return_type")
                    .map(|n| self.lower_type(n)),
                effects: self.lower_effects(node.child_by_field_name("effects")),
                body: self.lower_expr_field(node, "body"),
                position: Some(self.pos(node)),
            },
            "extern_declaration" => Stmt::Extern {
                name: self.field_text(node, "name"),
                parameters: self.lower_extern_params(node.child_by_field_name("parameters")),
                return_type: node
                    .child_by_field_name("return_type")
                    .map(|n| self.lower_type(n)),
                position: Some(self.pos(node)),
            },
            "type_declaration" => self.lower_type_decl(node),
            "effect_declaration" => Stmt::Effect {
                name: self.field_text(node, "name"),
                operations: self.lower_operations(node),
                position: Some(self.pos(node)),
            },
            "module_declaration" => Stmt::Module {
                name: self.field_text(node, "name"),
                body: self
                    .named_of_kind(node, "module_statement")
                    .iter()
                    .filter_map(|n| self.first_named(*n))
                    .filter_map(|n| self.lower_stmt(n))
                    .collect(),
            },
            "expression_statement" => Stmt::Expr(self.lower_expr(self.first_named(node)?)),
            _ => return None,
        })
    }

    fn lower_type_decl(&self, node: Node<'_>) -> Stmt {
        let def = node.child_by_field_name("definition");
        let variants = match def.map(|d| (d.kind(), d)) {
            Some(("union_type", d)) => self
                .named_of_kind(d, "variant")
                .iter()
                .map(|v| self.lower_variant(*v))
                .collect(),
            Some(("record_type", d)) => vec![TypeVariant {
                name: self.field_text(node, "name"),
                fields: self.lower_field_decls(d),
            }],
            _ => Vec::new(),
        };
        Stmt::Type {
            name: self.field_text(node, "name"),
            type_params: node
                .child_by_field_name("type_parameters")
                .map(|tp| self.texts_of_kind(tp, "identifier"))
                .unwrap_or_default(),
            variants,
            validation_func: self
                .first_child_of_kind(node, "type_validation")
                .and_then(|tv| self.first_named(tv))
                .map(|n| self.text(n)),
            position: Some(self.pos(node)),
        }
    }

    fn lower_variant(&self, node: Node<'_>) -> TypeVariant {
        TypeVariant {
            name: self.field_text(node, "name"),
            fields: node
                .child(node.child_count().saturating_sub(1))
                .filter(|_| node.child_count() > 1)
                .map(|_| self.lower_field_decls(node))
                .unwrap_or_default(),
        }
    }

    fn lower_field_decls(&self, node: Node<'_>) -> Vec<TypeField> {
        let mut out = Vec::new();
        for fd in self.descendants_of_kind(node, "field_declaration") {
            out.push(TypeField {
                name: self.field_text(fd, "name"),
                ty: fd
                    .child_by_field_name("type")
                    .map(|n| self.lower_type(n).name)
                    .unwrap_or_default(),
                constraint: None,
            });
        }
        out
    }

    fn lower_operations(&self, node: Node<'_>) -> Vec<EffectOperation> {
        self.named_of_kind(node, "operation_declaration")
            .iter()
            .map(|op| EffectOperation {
                name: self.field_text(*op, "name"),
                ty: op
                    .child_by_field_name("type")
                    .map(|n| self.text(n))
                    .unwrap_or_default(),
                parameters: Vec::new(),
                return_type: String::new(),
            })
            .collect()
    }

    pub(crate) fn lower_params(&self, list: Option<Node<'_>>) -> Vec<Parameter> {
        let Some(list) = list else { return Vec::new() };
        self.named_of_kind(list, "parameter")
            .iter()
            .map(|p| Parameter {
                name: self.field_text(*p, "name"),
                ty: p.child_by_field_name("type").map(|n| self.lower_type(n)),
            })
            .collect()
    }

    fn lower_extern_params(&self, list: Option<Node<'_>>) -> Vec<ExternParameter> {
        let Some(list) = list else { return Vec::new() };
        self.named_of_kind(list, "extern_parameter")
            .iter()
            .map(|p| ExternParameter {
                name: self.field_text(*p, "name"),
                ty: p
                    .child_by_field_name("type")
                    .map_or_else(|| TypeExpr::named(""), |n| self.lower_type(n)),
            })
            .collect()
    }

    fn lower_effects(&self, effects: Option<Node<'_>>) -> Vec<String> {
        let Some(effects) = effects else {
            return Vec::new();
        };
        self.descendants_of_kind(effects, "identifier")
            .iter()
            .map(|n| self.text(*n))
            .collect()
    }

    /// Lower a `_type` node (function/generic/array/identifier).
    pub(crate) fn lower_type(&self, node: Node<'_>) -> TypeExpr {
        match node.kind() {
            "type_identifier" => TypeExpr::named(
                self.first_named(node)
                    .map(|n| self.text(n))
                    .unwrap_or_default(),
            ),
            "generic_type" => {
                let mut t = TypeExpr::named(self.field_text(node, "name"));
                if let Some(list) = self.first_child_of_kind(node, "type_list") {
                    t.generic_params = self.lower_type_list(list);
                }
                t
            }
            "array_type" => {
                let mut t = TypeExpr::named(self.field_text(node, "name"));
                t.is_array = true;
                t.array_element = self.descendants_type_in(node).map(Box::new);
                t
            }
            "function_type" => {
                let types = self
                    .first_child_of_kind(node, "type_list")
                    .map(|l| self.lower_type_list(l))
                    .unwrap_or_default();
                // last bare type child is the return type
                let ret = self.last_type_child(node);
                TypeExpr {
                    name: "fn".into(),
                    generic_params: Vec::new(),
                    is_array: false,
                    array_element: None,
                    is_function: true,
                    parameter_types: types,
                    return_type: ret.map(Box::new),
                    position: Some(self.pos(node)),
                }
            }
            _ => TypeExpr::named(self.text(node)),
        }
    }

    fn lower_type_list(&self, list: Node<'_>) -> Vec<TypeExpr> {
        let mut out = Vec::new();
        let mut cursor = list.walk();
        for child in list.named_children(&mut cursor) {
            if Self::is_type_kind(child.kind()) {
                out.push(self.lower_type(child));
            }
        }
        out
    }

    fn descendants_type_in(&self, node: Node<'_>) -> Option<TypeExpr> {
        let mut cursor = node.walk();
        let found = node
            .named_children(&mut cursor)
            .find(|c| Self::is_type_kind(c.kind()));
        found.map(|c| self.lower_type(c))
    }

    fn last_type_child(&self, node: Node<'_>) -> Option<TypeExpr> {
        let mut cursor = node.walk();
        let found = node
            .named_children(&mut cursor)
            .filter(|c| Self::is_type_kind(c.kind()))
            .last();
        found.map(|c| self.lower_type(c))
    }

    fn is_type_kind(kind: &str) -> bool {
        matches!(
            kind,
            "type_identifier" | "generic_type" | "array_type" | "function_type"
        )
    }

    // ---- Patterns ----
    pub(crate) fn lower_pattern(&self, node: Node<'_>) -> Pattern {
        match node.kind() {
            "pattern" => {
                if let Some(inner) = self.first_named(node) {
                    return self.lower_pattern_inner(node, inner);
                }
                // bare `_` wildcard has no named child
                Pattern::Wildcard
            }
            _ => self.lower_pattern_inner(node, node),
        }
    }

    fn lower_pattern_inner(&self, pat: Node<'_>, inner: Node<'_>) -> Pattern {
        match inner.kind() {
            // A `-N` / `+N` pattern carries the sign in the `operator` field
            // (grammar: `seq(operator: choice('-','+'), literal)`); fold it into
            // the literal so `-5` matches `-5`, not `5`.
            "literal" => {
                let lit = self.lower_literal(inner);
                let negated = pat
                    .child_by_field_name("operator")
                    .is_some_and(|op| self.text(op) == "-");
                Pattern::Literal(Box::new(if negated { negate_literal(lit) } else { lit }))
            }
            "field_pattern" => Pattern::Structural {
                fields: self.field_pattern_names(inner),
            },
            "identifier" => {
                // Could be: constructor `Ctor { fields }`, type-annotated, sub-patterns,
                // or a bare binding. Inspect siblings of the name field.
                let name = self.text(inner);
                if let Some(fp) = self.first_child_of_kind(pat, "field_pattern") {
                    let fields = self.field_pattern_names(fp);
                    if pat.child_by_field_name("type").is_some()
                        || Self::has_colon_before_brace(pat)
                    {
                        return Pattern::Constructor {
                            name,
                            fields,
                            sub_patterns: vec![],
                        };
                    }
                    return Pattern::Constructor {
                        name,
                        fields,
                        sub_patterns: vec![],
                    };
                }
                if let Some(ty) = pat.child_by_field_name("type") {
                    return Pattern::TypeAnnotated {
                        name,
                        ty: self.lower_type(ty),
                    };
                }
                let subs: Vec<Pattern> = self
                    .named_of_kind(pat, "pattern")
                    .iter()
                    .map(|p| self.lower_pattern(*p))
                    .collect();
                if !subs.is_empty() {
                    return Pattern::Constructor {
                        name,
                        fields: vec![],
                        sub_patterns: subs,
                    };
                }
                Pattern::Binding(name)
            }
            _ => Pattern::Wildcard,
        }
    }

    fn has_colon_before_brace(_node: Node<'_>) -> bool {
        false
    }

    fn field_pattern_names(&self, fp: Node<'_>) -> Vec<String> {
        self.texts_of_kind(fp, "identifier")
    }

    // ---- small node helpers ----
    pub(crate) fn field_text(&self, node: Node<'_>, field: &str) -> String {
        node.child_by_field_name(field)
            .map(|n| self.text(n))
            .unwrap_or_default()
    }

    #[expect(
        clippy::unused_self,
        reason = "kept for Lowerer method-call ergonomics"
    )]
    pub(crate) fn named_of_kind<'t>(&self, node: Node<'t>, kind: &str) -> Vec<Node<'t>> {
        let mut out = Vec::new();
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            if child.kind() == kind {
                out.push(child);
            }
        }
        out
    }

    /// The source text of every named child of `node` of the given `kind`.
    pub(crate) fn texts_of_kind(&self, node: Node<'_>, kind: &str) -> Vec<String> {
        self.named_of_kind(node, kind)
            .iter()
            .map(|n| self.text(*n))
            .collect()
    }

    /// The lowered expression of every named child of `node` of the given `kind`.
    pub(crate) fn exprs_of_kind(&self, node: Node<'_>, kind: &str) -> Vec<Expr> {
        self.named_of_kind(node, kind)
            .iter()
            .map(|e| self.lower_expr(*e))
            .collect()
    }

    /// Recursive search for all descendants of a kind (for nested wrappers).
    #[expect(
        clippy::self_only_used_in_recursion,
        reason = "kept for Lowerer method-call ergonomics"
    )]
    pub(crate) fn descendants_of_kind<'t>(&self, node: Node<'t>, kind: &str) -> Vec<Node<'t>> {
        let mut out = Vec::new();
        let mut cursor = node.walk();
        for child in node.named_children(&mut cursor) {
            if child.kind() == kind {
                out.push(child);
            } else {
                out.extend(self.descendants_of_kind(child, kind));
            }
        }
        out
    }
}

/// Negate a numeric literal for a `-N` pattern; non-numeric literals pass through.
fn negate_literal(e: Expr) -> Expr {
    match e {
        Expr::Integer(n) => Expr::Integer(-n),
        Expr::Float(f) => Expr::Float(-f),
        other => other,
    }
}
