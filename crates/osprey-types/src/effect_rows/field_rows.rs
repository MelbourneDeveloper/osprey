//! Effect rows written on function-typed fields: `edit: fn(string) -> A ![]`
//! bounds every function a construction stores in the field, and a call
//! through the field — projected, pattern-bound or called as a method —
//! performs exactly the written row, whatever provenance the stored value kept
//! on its way there. Implements [EFFECTS-FIELD-ROWS].
use super::{
    requirement_name, Analyzer, Callable, CallableEnv, Expr, FieldAssignment, HashMap, Index,
    KnownCallable, Requirement, Summary, TypeError, Value,
};
use osprey_ast::{EffectRef, TypeField};
use std::collections::BTreeMap;

/// The declarations of every field name: how many there are, and the written
/// rows among them, keyed by constructor so a construction finds its own.
#[derive(Default)]
pub(super) struct FieldRows<'a> {
    declarations: HashMap<String, usize>,
    written: HashMap<String, Vec<(&'a str, &'a TypeField)>>,
}

/// A violated field row, keyed by construction site and field so the last
/// analysis of each site — the one over converged rows — decides it.
pub(crate) type FieldRowViolations = BTreeMap<(usize, String), TypeError>;

impl<'a> FieldRows<'a> {
    pub(super) fn collect(&mut self, constructor: &'a str, fields: &'a [TypeField]) {
        for field in fields {
            *self.declarations.entry(field.name.clone()).or_default() += 1;
            if field.effects.is_some() {
                self.written
                    .entry(field.name.clone())
                    .or_default()
                    .push((constructor, field));
            }
        }
    }

    fn of(&self, constructor: &str, field: &str) -> Option<&'a [EffectRef]> {
        self.written
            .get(field)?
            .iter()
            .find(|(owner, _)| *owner == constructor)
            .and_then(|(_, declared)| declared.effects.as_deref())
    }
}

impl Index<'_> {
    /// What a call through `field` invokes: the written row when every
    /// declaration of the field writes one, otherwise the projected value.
    pub(super) fn field_value(&self, field: &str, projected: Option<Value>) -> Option<Value> {
        self.trusted_field(field).or(projected)
    }

    fn trusted_field(&self, field: &str) -> Option<Value> {
        let written = self.field_rows.written.get(field)?;
        if self.field_rows.declarations.get(field) != Some(&written.len()) {
            return None;
        }
        let mut summary = Summary::default();
        for (_, declared) in written {
            summary.union(self.row_summary(declared.effects.as_deref()?)?);
        }
        let known = KnownCallable {
            parameters: Vec::new(),
            summary,
            returned: None,
        };
        Some(Value {
            callable: Some(Callable::Known(Box::new(known))),
            ..Value::default()
        })
    }

    /// Every operation the row's effects declare. `None` when an entry names
    /// no declared effect, or a generic one without its instance.
    fn row_summary(&self, row: &[EffectRef]) -> Option<Summary> {
        let mut summary = Summary::default();
        for effect in row {
            let arguments = self.row_arguments(effect)?;
            for operation in self.effect_operations.get(&effect.name)? {
                let _ = summary.required.insert(Requirement::new(
                    &effect.name,
                    operation,
                    arguments.clone(),
                ));
            }
        }
        Some(summary)
    }

    fn row_arguments(&self, effect: &EffectRef) -> Option<Vec<String>> {
        let arity = *self.effects.get(&effect.name)?;
        (effect.type_args.len() == arity).then(|| {
            effect
                .type_args
                .iter()
                .map(|argument| crate::convert::type_expr_to_type(argument, &HashMap::new()))
                .map(|argument| argument.to_string())
                .collect()
        })
    }

    /// A row that cannot bound its field: on a field that is not a function,
    /// or naming an effect that is not declared or is generic without its
    /// instance.
    pub(super) fn field_row_errors(&self) -> Vec<TypeError> {
        let mut errors = Vec::new();
        for (constructor, field) in self.field_rows.written.values().flatten() {
            let row = field.effects.as_deref().unwrap_or_default();
            let position = row.first().and_then(|effect| effect.position);
            let ty = crate::convert::type_name_to_type(&field.ty, &HashMap::new());
            if !matches!(ty, crate::ty::Type::Fun { .. }) {
                errors.push(declaration_error(
                    constructor,
                    field,
                    "an effect row needs a function type",
                    position,
                ));
            } else if row
                .iter()
                .any(|effect| self.row_arguments(effect).is_none())
            {
                errors.push(declaration_error(
                    constructor,
                    field,
                    "its effect row names an undeclared effect or a generic one without its type arguments",
                    position,
                ));
            }
        }
        errors
    }
}

fn declaration_error(
    constructor: &str,
    field: &TypeField,
    reason: &str,
    position: Option<osprey_ast::Position>,
) -> TypeError {
    TypeError::new(format!(
        "field `{}` of `{constructor}`: {reason}",
        field.name
    ))
    .with_pos(position)
}

impl Analyzer<'_> {
    /// Hold each function stored in a row-declared field to that row.
    pub(super) fn check_field_rows(
        &self,
        site: usize,
        constructor: &str,
        fields: &[FieldAssignment],
        scope: &[String],
        env: &CallableEnv,
    ) {
        for field in fields {
            let Some(row) = self.index.field_rows.of(constructor, &field.name) else {
                continue;
            };
            let violation = self.field_row_violation(constructor, field, row, scope, env);
            let key = (site, field.name.clone());
            let mut violations = self.instances.field_row_violations.borrow_mut();
            let _ = match violation {
                Some(error) => violations.insert(key, error),
                None => violations.remove(&key),
            };
        }
    }

    fn field_row_violation(
        &self,
        constructor: &str,
        field: &FieldAssignment,
        row: &[EffectRef],
        scope: &[String],
        env: &CallableEnv,
    ) -> Option<TypeError> {
        let stored = self
            .callable(&field.value, scope, env)
            .map(|callee| self.invoke_with_values(callee, &[]));
        let reason = match stored {
            None => Some("the stored function's effects cannot be proven".to_owned()),
            Some(summary) => outside_row(&summary, row),
        }?;
        // A stored lambda carries its own position; anything else has none.
        let position = match &field.value {
            Expr::Lambda { position, .. } => *position,
            _ => None,
        };
        Some(
            TypeError::new(format!(
                "field `{}` of `{constructor}` declares effect row {}, but {reason}",
                field.name,
                spelled(row)
            ))
            .with_pos(position),
        )
    }
}

/// Why `summary` exceeds `row`, or `None` when the row admits it.
fn outside_row(summary: &Summary, row: &[EffectRef]) -> Option<String> {
    let undeclared: Vec<_> = summary
        .required
        .iter()
        .filter(|requirement| !row.iter().any(|effect| effect.name == requirement.effect))
        .map(requirement_name)
        .collect();
    if !undeclared.is_empty() {
        return Some(format!(
            "the stored function performs {}",
            undeclared.join(", ")
        ));
    }
    if row.is_empty() && !summary.runtime_builtins.is_empty() {
        let builtins: Vec<_> = summary.runtime_builtins.iter().cloned().collect();
        return Some(format!("the stored function uses {}", builtins.join(", ")));
    }
    (!summary.parameter_uses.is_empty() || summary.unresolved_dynamic_call)
        .then(|| "the stored function calls a function whose effects cannot be proven".to_owned())
}

fn spelled(row: &[EffectRef]) -> String {
    let names: Vec<_> = row.iter().map(|effect| effect.name.as_str()).collect();
    format!("![{}]", names.join(", "))
}
