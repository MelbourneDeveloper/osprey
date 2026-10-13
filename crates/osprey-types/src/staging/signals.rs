//! Immutable reactive snapshots. Implements [STAGE-SIGNALS-CACHE], [STAGE-SIGNALS-REBUILD].
use super::StaticEffects;
use osprey_ast::{
    mutate, walk_program, AstVisitor, Expr, FieldAssignment, MatchArm, NamedArgument, Parameter,
    Pattern, Program, Stmt,
};

pub(super) fn is_intrinsic(name: &str) -> bool {
    matches!(
        name,
        "signalDependencies" | "signalSnapshot" | "signalRefresh"
    )
}

/// Freeze checker-derived metadata before discharge removes signal operations.
/// Traverse source and clone in the same preorder; pointer identities never
/// cross the clone boundary. Implements [STAGE-SIGNALS-EXACT].
pub(super) fn materialize_signal_queries(
    program: &Program,
    types: &crate::ProgramTypes,
) -> Program {
    let mut effects = StaticEffects::default();
    walk_program(program, &mut effects);
    let mut queries = SignalQueries {
        types,
        effects,
        values: Vec::new(),
    };
    walk_program(program, &mut queries);
    let mut values = queries.values.into_iter();
    let mut resolved = program.clone();
    for statement in &mut resolved.statements {
        mutate::statement_children_mut(statement, &mut |expression| {
            materialize_query(expression, &mut values);
        });
    }
    resolved
}

struct SignalQueries<'a> {
    types: &'a crate::ProgramTypes,
    effects: StaticEffects,
    values: Vec<Option<Vec<String>>>,
}

impl AstVisitor for SignalQueries<'_> {
    fn expression(&mut self, expression: &Expr) {
        let query = match expression {
            Expr::Call { function, .. } => self
                .types
                .signal_queries
                .get(&std::ptr::from_ref(function.as_ref()).addr()),
            _ => None,
        };
        self.values
            .push(query.map(|row| self.effects.dependencies(row)));
    }
}

fn materialize_query(
    expression: &mut Expr,
    values: &mut impl Iterator<Item = Option<Vec<String>>>,
) {
    let query = values.next().flatten();
    mutate::children_mut(expression, &mut |child| materialize_query(child, values));
    if let Some(operations) = query {
        *expression = expand(expression, operations);
    }
}

fn expand(expression: &Expr, operations: Vec<String>) -> Expr {
    let Expr::Call {
        function,
        arguments,
        named_arguments,
    } = expression
    else {
        return expression.clone();
    };
    let Some(name) = callee_name(function) else {
        return expression.clone();
    };
    let mut statements = bound_arguments(name, arguments, named_arguments);
    statements.push(dependency_binding(operations));
    let value = match name {
        "signalSnapshot" => snapshot(Expr::Integer(1)),
        "signalRefresh" => refresh(),
        _ => ident("$signalDependencies"),
    };
    Expr::Block {
        statements,
        value: Some(Box::new(value)),
        position: None,
    }
}

fn callee_name(function: &Expr) -> Option<&str> {
    match function {
        Expr::Identifier(name) => Some(name),
        Expr::TypeApply { function, .. } => callee_name(function),
        _ => None,
    }
}

fn bound_arguments(name: &str, arguments: &[Expr], named: &[NamedArgument]) -> Vec<Stmt> {
    let slots = if name == "signalRefresh" {
        vec!["$signalCached", "$signalChanged", "$signalWork"]
    } else {
        vec!["$signalWork"]
    };
    let positional = slots
        .iter()
        .zip(arguments)
        .map(|(name, value)| binding(name, value.clone()));
    let named = named.iter().map(|argument| {
        binding(
            match argument.name.as_str() {
                "cached" => "$signalCached",
                "changed" => "$signalChanged",
                _ => "$signalWork",
            },
            argument.value.clone(),
        )
    });
    positional.chain(named).collect()
}

fn dependency_binding(operations: Vec<String>) -> Stmt {
    let empty = operations.is_empty();
    let mut declaration = binding(
        "$signalDependencies",
        Expr::List(operations.into_iter().map(Expr::Str).collect(), None),
    );
    if empty {
        if let Stmt::Let { ty, .. } = &mut declaration {
            *ty = Some(string_list_type());
        }
    }
    declaration
}

fn refresh() -> Expr {
    let dirty = binary(
        "||",
        binary(
            "!=",
            joined(ident("$signalDependencies")),
            joined(field("dependencies")),
        ),
        binary(
            "||",
            call(
                "listContains",
                vec![
                    ident("$signalDependencies"),
                    Expr::Str("<unknown>".to_owned()),
                ],
            ),
            binary(
                "||",
                call(
                    "listContains",
                    vec![ident("$signalChanged"), Expr::Str("<unknown>".to_owned())],
                ),
                affected(),
            ),
        ),
    );
    Expr::Match {
        value: Box::new(dirty),
        arms: vec![
            arm(
                true,
                snapshot(call("satAdd", vec![field("revision"), Expr::Integer(1)])),
            ),
            arm(false, ident("$signalCached")),
        ],
    }
}

fn affected() -> Expr {
    let predicate = Expr::Lambda {
        return_type: None,
        parameters: vec!["$signalDirty", "$signalId"]
            .into_iter()
            .map(|name| Parameter {
                name: name.to_owned(),
                ty: None,
                inline_constraint: false,
            })
            .collect(),
        body: Box::new(binary("||", ident("$signalDirty"), changed_entry())),
        position: None,
    };
    call(
        "fold",
        vec![
            call(
                "range",
                vec![
                    Expr::Integer(0),
                    call("length", vec![ident("$signalDependencies")]),
                ],
            ),
            Expr::Bool(false),
            predicate,
        ],
    )
}

fn snapshot(revision: Expr) -> Expr {
    Expr::Object(vec![
        assignment("dependencies", ident("$signalDependencies")),
        assignment("value", call("$signalWork", vec![])),
        assignment("revision", revision),
    ])
}

fn assignment(name: &str, value: Expr) -> FieldAssignment {
    FieldAssignment {
        name: name.to_owned(),
        value,
    }
}

fn binding(name: &str, value: Expr) -> Stmt {
    Stmt::Let {
        name: name.to_owned(),
        mutable: false,
        ty: None,
        value,
        doc: None,
        position: None,
    }
}

fn ident(name: &str) -> Expr {
    Expr::Identifier(name.to_owned())
}

fn field(name: &str) -> Expr {
    Expr::FieldAccess {
        target: Box::new(ident("$signalCached")),
        field: name.to_owned(),
    }
}

fn call(name: &str, arguments: Vec<Expr>) -> Expr {
    Expr::Call {
        function: Box::new(ident(name)),
        arguments,
        named_arguments: Vec::new(),
    }
}

fn binary(op: &str, left: Expr, right: Expr) -> Expr {
    Expr::Binary {
        op: op.to_owned(),
        left: Box::new(left),
        right: Box::new(right),
        position: None,
    }
}

fn arm(value: bool, body: Expr) -> MatchArm {
    MatchArm {
        pattern: Pattern::Literal(Box::new(Expr::Bool(value))),
        body,
    }
}

fn string_list_type() -> osprey_ast::TypeExpr {
    let mut ty = osprey_ast::TypeExpr::named("List");
    ty.generic_params
        .push(osprey_ast::TypeExpr::named("string"));
    ty
}

fn changed_entry() -> Expr {
    Expr::Match {
        value: Box::new(call(
            "listGet",
            vec![ident("$signalDependencies"), ident("$signalId")],
        )),
        arms: vec![
            MatchArm {
                pattern: Pattern::Constructor {
                    name: "Success".to_owned(),
                    fields: vec!["$signalValue".to_owned()],
                    sub_patterns: Vec::new(),
                },
                body: call(
                    "listContains",
                    vec![ident("$signalChanged"), ident("$signalValue")],
                ),
            },
            MatchArm {
                pattern: Pattern::Constructor {
                    name: "Error".to_owned(),
                    fields: Vec::new(),
                    sub_patterns: Vec::new(),
                },
                body: Expr::Bool(true),
            },
        ],
    }
}

fn joined(dependencies: Expr) -> Expr {
    call("join", vec![dependencies, Expr::Str("\n".to_owned())])
}

// Read argument labels from the same documentation/signature table as tooling.
pub(crate) fn parameter_names() -> std::collections::HashMap<String, Vec<String>> {
    crate::builtin_names()
        .into_iter()
        .filter(|name| is_intrinsic(name))
        .filter_map(|name| {
            crate::builtin_docs::parameter_names(&name).map(|parameters| (name, parameters))
        })
        .collect()
}
