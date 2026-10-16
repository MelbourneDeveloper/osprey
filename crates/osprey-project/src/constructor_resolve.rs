//! Imported constructors become record construction after identity resolution.
//! Implements [FLAVOR-PYTHON-CLASS], [MODULES-RESOLUTION].
use crate::model::{DeclKind, SymbolKey};
use crate::resolve::{Context, Resolver};
use osprey_ast::{Expr, FieldAssignment, TypeExpr};

impl Resolver<'_> {
    pub(crate) fn normalize_constructor_call(&mut self, expression: &mut Expr, context: &Context) {
        let Some((name, type_args)) = self.constructor_call_head(expression) else {
            return;
        };
        let Expr::Call {
            arguments,
            named_arguments,
            ..
        } = expression
        else {
            return;
        };
        if !self.record_arguments(arguments, context) {
            return;
        }
        *expression = construction(name, type_args, std::mem::take(named_arguments));
    }

    fn constructor_call_head(&self, expression: &Expr) -> Option<(String, Vec<TypeExpr>)> {
        let Expr::Call { function, .. } = expression else {
            return None;
        };
        let (name, type_args) = constructor_head(function)?;
        self.declared_constructor(&name)
            .then_some((name, type_args))
    }

    fn record_arguments(&mut self, arguments: &[Expr], context: &Context) -> bool {
        if arguments.is_empty() {
            return true;
        }
        self.error(
            context.source,
            None,
            "record construction requires named fields",
        );
        false
    }

    fn declared_constructor(&self, name: &str) -> bool {
        let Some(source) = self.source_names.get(name) else {
            return false;
        };
        let segments = source.split("::").map(str::to_owned).collect::<Vec<_>>();
        let Some((namespace, path)) = segments.split_first() else {
            return false;
        };
        let key = SymbolKey::new(namespace, path.to_vec());
        self.graph
            .declarations
            .get(&key)
            .is_some_and(|info| matches!(info.kind, DeclKind::Type | DeclKind::Constructor))
    }
}

fn construction(
    name: String,
    type_args: Vec<TypeExpr>,
    arguments: Vec<osprey_ast::NamedArgument>,
) -> Expr {
    Expr::TypeConstructor {
        name,
        type_args,
        fields: arguments
            .into_iter()
            .map(|argument| FieldAssignment {
                name: argument.name,
                value: argument.value,
            })
            .collect(),
    }
}

fn constructor_head(function: &Expr) -> Option<(String, Vec<TypeExpr>)> {
    match function {
        Expr::Identifier(name) => Some((name.clone(), Vec::new())),
        Expr::TypeApply {
            function,
            type_args,
            ..
        } => {
            let Expr::Identifier(name) = function.as_ref() else {
                return None;
            };
            Some((name.clone(), type_args.clone()))
        }
        _ => None,
    }
}
