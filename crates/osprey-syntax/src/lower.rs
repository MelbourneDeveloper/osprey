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
                // Keep the full source text (`List<Self>`, `(int) -> bool`) —
                // taking just the lowered head name would collapse a function
                // type to `fn` and a generic to its constructor.
                ty: fd
                    .child_by_field_name("type")
                    .map(|n| self.text(n))
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

    pub(crate) fn last_type_child(&self, node: Node<'_>) -> Option<TypeExpr> {
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
            // (grammar: `seq(operator: choice('-','+'), integer|float)`); fold it
            // into the literal so `-5` matches `-5`, not `5`. Scalar literals now
            // appear unwrapped (no `literal` node) so `[…]` stays a list_pattern.
            "integer" | "float" | "boolean" | "string" | "interpolated_string" => {
                let lit = self.lower_literal_node(inner);
                let negated = pat
                    .child_by_field_name("operator")
                    .is_some_and(|op| self.text(op) == "-");
                Pattern::Literal(Box::new(if negated { negate_literal(lit) } else { lit }))
            }
            "list_pattern" => self.lower_list_pattern(inner),
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

    /// Build a [`Pattern::List`] from a `list_pattern` node: the `element` fields
    /// (each a `pattern`) become the fixed-prefix patterns in source order, and
    /// the `rest` field (an identifier) becomes the optional tail binder.
    fn lower_list_pattern(&self, node: Node<'_>) -> Pattern {
        let elements = self
            .named_of_kind(node, "pattern")
            .iter()
            .map(|p| self.lower_pattern(*p))
            .collect();
        let rest = node.child_by_field_name("rest").map(|r| self.text(r));
        Pattern::List { elements, rest }
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

#[cfg(test)]
#[expect(
    clippy::indexing_slicing,
    reason = "test assertions: an out-of-bounds index is a test failure, not a production panic"
)]
mod tests {
    use crate::{parse_program, parse_tree};
    use osprey_ast::{Expr, Pattern, Stmt};
    use tree_sitter::Node;

    fn stmts(src: &str) -> Vec<Stmt> {
        let parsed = parse_program(src);
        assert!(parsed.errors.is_empty(), "errors: {:?}", parsed.errors);
        parsed.program.statements
    }

    fn one(src: &str) -> Stmt {
        let mut s = stmts(src);
        assert_eq!(s.len(), 1, "expected one stmt for {src:?}");
        s.pop().unwrap()
    }

    /// Find the first descendant node of a given kind anywhere in the tree.
    fn find_kind<'t>(node: Node<'t>, kind: &str) -> Option<Node<'t>> {
        if node.kind() == kind {
            return Some(node);
        }
        let mut cursor = node.walk();
        let children: Vec<Node<'t>> = node.children(&mut cursor).collect();
        children.into_iter().find_map(|c| find_kind(c, kind))
    }

    #[test]
    fn lowers_import_and_module() {
        // `import` exercises lower_stmt's Import arm (texts_of_kind identifiers).
        match one("import std.io.file\n") {
            Stmt::Import { module } => assert_eq!(module, vec!["std", "io", "file"]),
            s => panic!("expected import, got {s:?}"),
        }
        // A module body re-enters lower_stmt for nested declarations.
        match one("module M {\n  let x = 1\n  fn f() = x\n}\n") {
            Stmt::Module { name, body } => {
                assert_eq!(name, "M");
                assert_eq!(body.len(), 2);
                assert!(matches!(body[0], Stmt::Let { .. }));
                assert!(matches!(body[1], Stmt::Function { .. }));
            }
            s => panic!("expected module, got {s:?}"),
        }
    }

    #[test]
    fn lowers_record_type_and_array_and_function_types() {
        // record_type definition (lower_type_decl record arm + lower_field_decls)
        match one("type Point = {\n  x: int,\n  y: int\n}\n") {
            Stmt::Type { name, variants, .. } => {
                assert_eq!(name, "Point");
                assert_eq!(variants.len(), 1);
                assert_eq!(variants[0].fields.len(), 2);
                assert_eq!(variants[0].fields[0].name, "x");
            }
            s => panic!("expected record type, got {s:?}"),
        }
        // array_type `Item[int]` (lower_type array_type arm + descendants_type_in),
        // a function type, and a generic type — all in one signature.
        match one(
            "fn f(xs: Item[int], g: fn(int) -> bool, m: Map<string, int>) -> Item[int] = xs\n",
        ) {
            Stmt::Function {
                parameters,
                return_type,
                ..
            } => {
                let arr = parameters[0].ty.as_ref().unwrap();
                assert!(arr.is_array);
                assert_eq!(arr.array_element.as_ref().unwrap().name, "int");
                let func = parameters[1].ty.as_ref().unwrap();
                assert!(func.is_function);
                assert_eq!(func.return_type.as_ref().unwrap().name, "bool");
                let gen = parameters[2].ty.as_ref().unwrap();
                assert_eq!(gen.generic_params.len(), 2);
                assert!(return_type.unwrap().is_array);
            }
            s => panic!("expected function, got {s:?}"),
        }
    }

    /// The single match arm's pattern for `match x { <arm> => 0  _ => 1 }`.
    fn first_pattern(arm: &str) -> Pattern {
        let src = format!("let r = match x {{ {arm} => 0  _ => 1 }}\n");
        match one(&src) {
            Stmt::Let {
                value: Expr::Match { mut arms, .. },
                ..
            } => arms.swap_remove(0).pattern,
            s => panic!("expected match, got {s:?}"),
        }
    }

    #[test]
    fn lowers_constructor_type_annotated_negative_and_type_params() {
        // Sub-pattern constructor `Some(inner)` -> identifier arm -> sub_patterns.
        assert!(matches!(
            first_pattern("Some(inner)"),
            Pattern::Constructor { sub_patterns, .. } if sub_patterns.len() == 1
        ));
        // `n: int` -> TypeAnnotated.
        assert!(matches!(
            first_pattern("n: int"),
            Pattern::TypeAnnotated { ref name, .. } if name == "n"
        ));
        // `-1.5` -> negated float literal (drives negate_literal's Float arm).
        assert!(matches!(
            first_pattern("-1.5"),
            Pattern::Literal(b) if matches!(*b, Expr::Float(f) if f < 0.0)
        ));
        // Generic type params on a type declaration (type_parameters field).
        match one("type Foo<T> = Bar | Baz\n") {
            Stmt::Type {
                type_params,
                variants,
                ..
            } => {
                assert_eq!(type_params, vec!["T"]);
                assert_eq!(variants.len(), 2);
            }
            s => panic!("expected type, got {s:?}"),
        }
    }

    #[test]
    fn negate_literal_passes_through_non_numeric() {
        // negate_literal flips numerics and returns non-numeric literals as-is.
        assert_eq!(super::negate_literal(Expr::Integer(3)), Expr::Integer(-3));
        assert_eq!(super::negate_literal(Expr::Float(2.0)), Expr::Float(-2.0));
        assert_eq!(
            super::negate_literal(Expr::Str("x".into())),
            Expr::Str("x".into())
        );
    }

    #[test]
    fn lowers_assignment_effects_structural_and_list_patterns() {
        // Reassignment statement (lower_stmt Assignment arm).
        match one("x = 5\n") {
            Stmt::Assignment { name, value, .. } => {
                assert_eq!(name, "x");
                assert_eq!(value, Expr::Integer(5));
            }
            s => panic!("expected assignment, got {s:?}"),
        }
        // Function effect clause `! [Log, State]` (lower_effects descendants).
        match one("fn act() ! [Log, State] = 1\n") {
            Stmt::Function { effects, .. } => assert_eq!(effects, vec!["Log", "State"]),
            s => panic!("expected function, got {s:?}"),
        }
        // Bare structural `{ name, age }` and a fixed-length list `[a, b]`.
        assert!(matches!(
            first_pattern("{ name, age }"),
            Pattern::Structural { fields } if fields == vec!["name", "age"]
        ));
        assert!(matches!(
            first_pattern("[a, b]"),
            Pattern::List { elements, rest: None } if elements.len() == 2
        ));
    }

    #[test]
    fn defensive_fallthrough_arms() {
        // Drive lower_stmt / lower_type / lower_pattern on a node kind none of
        // their match arms expect (a `line_comment`), hitting their `_` fallbacks.
        let src = "// hi\nlet x = 1\n";
        let tree = parse_tree(src).unwrap();
        let lw = super::Lowerer::new(src.as_bytes());
        let comment = find_kind(tree.root_node(), "line_comment").unwrap();

        assert!(lw.lower_stmt(comment).is_none()); // `_ => return None`
        assert_eq!(lw.lower_type(comment).name, lw.text(comment)); // `_ => named(text)`
        assert!(matches!(lw.lower_pattern(comment), Pattern::Wildcard)); // `_` -> inner `_`
    }
}
