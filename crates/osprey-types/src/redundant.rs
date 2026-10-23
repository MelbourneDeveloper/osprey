//! Redundant type-annotation detection.
//!
//! Implements [TYPE-ANNOTATION-REDUNDANT]. An annotation is redundant exactly
//! when erasing it leaves the solved program unchanged, which is a property of
//! inference and not of the annotation's spelling: the same
//! `string -> int -> string` is redundant on one function and load-bearing on
//! the next. The rule is therefore decided by re-running the inferrer over an
//! erased copy of the program and comparing what it published.

use std::collections::{BTreeMap, HashMap, HashSet};

use osprey_ast::{Position, Program};

use crate::check::infer_checked;
use crate::redundant_sites::{erase, sites, Site, Slot};
use crate::ty::{Type, VarId};

/// The rule identifier this diagnostic is reported under, and the name a
/// per-rule severity setting will address it by ([TYPE-ANNOTATION-REDUNDANT]).
pub const REDUNDANT_ANNOTATION: &str = "redundant-annotation";

/// A diagnostic that is true of the program but is not a reason to reject it.
///
/// Warnings are reported alongside [`TypeError`](crate::TypeError)s and change
/// no exit code and no generated code.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct TypeWarning {
    /// Human-readable text, already naming the offending declaration.
    pub message: String,
    /// Where the annotated declaration was written, when recorded.
    pub position: Option<Position>,
    /// The rule that raised it.
    pub rule: &'static str,
}

/// The source construct whose redundant annotation can be removed.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum RedundantTarget {
    /// One standalone written type signature.
    Signature,
    /// An inline parameter annotation.
    Parameter {
        /// Source spelling of the parameter.
        name: String,
        /// Zero-based slot in the canonical parameter list.
        index: usize,
    },
    /// An inline return type annotation.
    Return,
    /// A value binding annotation.
    Binding,
}

/// A warning together with source identity for precise editor edits.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct RedundantAnnotation {
    /// The existing public diagnostic, unchanged by source edit support.
    pub warning: TypeWarning,
    /// Original annotation position, when the frontend recorded it.
    pub annotation_position: Option<Position>,
    /// The annotation's slot when only a declaration position is available.
    pub target: RedundantTarget,
}

/// Every type the inferrer publishes, under a stable label.
///
/// Function signatures are keyed by name; bindings, lambdas and list literals
/// are keyed by source position, which an erasure never moves. Comparing the
/// whole map is what makes the rule honest: an annotation that changed only a
/// `let`'s type is still an annotation that changed something.
type Published = BTreeMap<String, Type>;

/// Equal signatures alone do not prove that dotted calls choose the same
/// implementation. Preserve dispatch choices alongside the inferred types.
pub(crate) struct Snapshot {
    pub(crate) types: Published,
    methods: crate::methods::Targets,
}

/// Every redundant annotation in `program`, in source order.
///
/// An ill-typed program yields none: its inferred types are the checker's
/// best effort at a program that does not typecheck, so no conclusion drawn
/// from comparing them would be trustworthy.
#[must_use]
pub fn redundant_annotations(program: &Program) -> Vec<TypeWarning> {
    redundant_annotations_where(program, |_| true)
}

/// Select one jointly removable set for the entire program, then report the
/// requested source annotations. Reports for separate files share that set.
#[must_use]
pub fn redundant_annotations_where(
    program: &Program,
    include: impl Fn(Option<Position>) -> bool,
) -> Vec<TypeWarning> {
    redundant_annotation_sites_where(program, include)
        .into_iter()
        .map(|site| site.warning)
        .collect()
}

/// Source identities for the same jointly removable set as the warning API.
#[must_use]
pub fn redundant_annotation_sites(program: &Program) -> Vec<RedundantAnnotation> {
    redundant_annotation_sites_where(program, |_| true)
}

/// Select source identities without changing the redundancy oracle or its set.
#[must_use]
pub fn redundant_annotation_sites_where(
    program: &Program,
    include: impl Fn(Option<Position>) -> bool,
) -> Vec<RedundantAnnotation> {
    redundant_annotation_sites_for(program, include, &|_| false)
}

/// [`redundant_annotation_sites_where`], holding the signatures of the
/// functions `contract` names fixed. An exported module member's signature is
/// its API: whether another project still needs it is not a question one
/// program can answer, so it is never reported and never erased while the
/// rest is judged. Annotations inside its body are judged as usual.
#[must_use]
pub fn redundant_annotation_sites_for(
    program: &Program,
    include: impl Fn(Option<Position>) -> bool,
    contract: &dyn Fn(&str) -> bool,
) -> Vec<RedundantAnnotation> {
    let written: Vec<Site> = sites(program)
        .into_iter()
        .filter(|site| !site.slot.signs(contract))
        .collect();
    if !written.iter().any(|site| include(site.position)) {
        return Vec::new();
    }
    let Some(baseline) = published(program) else {
        return Vec::new();
    };
    let all: Vec<usize> = (0..written.len()).collect();
    let removable = if preserves(program, &baseline, &written, &all) {
        all
    } else {
        accumulate(program, &baseline, &written, judges())
    };
    removable
        .iter()
        .filter_map(|index| written.get(*index))
        .filter(|site| include(site.position))
        .map(annotation)
        .collect()
}

fn annotation(site: &Site) -> RedundantAnnotation {
    let target = match &site.slot {
        Slot::Signature { .. } => RedundantTarget::Signature,
        Slot::Param {
            parameter, index, ..
        } => RedundantTarget::Parameter {
            name: parameter.clone(),
            index: *index,
        },
        Slot::Return { .. } => RedundantTarget::Return,
        Slot::Binding { .. } => RedundantTarget::Binding,
    };
    RedundantAnnotation {
        warning: warn(site),
        annotation_position: site.source,
        target,
    }
}

/// Grow the reported set in source order, keeping only the annotations that
/// stay removable alongside everything already kept.
///
/// A reader acts on the whole list, so the whole list has to be safe. Judging
/// each annotation against the untouched program answers a question nobody
/// asked: mutually-pinning annotations each look removable while the others
/// hold the type down, and a reader who deletes all of them gets a different
/// program. Every set this returns has been solved for as a set.
///
/// Each candidate costs a whole-program solve, so `window` undecided
/// candidates are judged at once on separate threads, each against the set
/// kept so far. Erasing more annotations only removes constraints: it never
/// restores a type another erasure changed, nor repairs an error. So a
/// candidate rejected against the kept set stays rejected as the set grows,
/// and when the window's accepted candidates also preserve the program
/// together, the one-at-a-time search would have kept each of them in turn.
/// That combined set is solved once; when it fails, only its first member is
/// kept and the rest are judged again against the larger set.
pub(crate) fn accumulate(
    program: &Program,
    baseline: &Snapshot,
    written: &[Site],
    window: usize,
) -> Vec<usize> {
    let mut kept: Vec<usize> = Vec::new();
    let mut pending: std::collections::VecDeque<usize> = (0..written.len()).collect();
    while !pending.is_empty() {
        let width = pending.len().min(window.max(1));
        let candidates: Vec<usize> = pending.drain(..width).collect();
        let verdicts = judged(program, baseline, written, &kept, &candidates);
        let accepted: Vec<usize> = candidates
            .iter()
            .zip(verdicts)
            .filter_map(|(index, keep)| keep.then_some(*index))
            .collect();
        settle(
            program,
            baseline,
            written,
            &mut kept,
            accepted,
            &mut pending,
        );
    }
    kept
}

/// Keep `accepted` at once when the set survives with all of them; otherwise
/// keep the first, which was judged against exactly `kept`, and requeue the
/// rest ahead of everything later.
fn settle(
    program: &Program,
    baseline: &Snapshot,
    written: &[Site],
    kept: &mut Vec<usize>,
    accepted: Vec<usize>,
    pending: &mut std::collections::VecDeque<usize>,
) {
    let Some((&first, rest)) = accepted.split_first() else {
        return;
    };
    let together: Vec<usize> = kept.iter().chain(&accepted).copied().collect();
    if rest.is_empty() || preserves(program, baseline, written, &together) {
        kept.extend(accepted);
        return;
    }
    kept.push(first);
    for index in rest.iter().rev() {
        pending.push_front(*index);
    }
}

/// Inference recurses on the syntax tree; a spawned thread gets the main
/// thread's headroom, not the 2 MiB default.
const JUDGE_STACK: usize = 64 * 1024 * 1024;

/// Whether each candidate stays removable alongside `kept`, in order.
fn judged(
    program: &Program,
    baseline: &Snapshot,
    written: &[Site],
    kept: &[usize],
    candidates: &[usize],
) -> Vec<bool> {
    let trial = |index: usize| {
        let mut chosen = kept.to_vec();
        chosen.push(index);
        preserves(program, baseline, written, &chosen)
    };
    if let [index] = candidates {
        return vec![trial(*index)];
    }
    std::thread::scope(|scope| {
        let judges: Vec<_> = candidates
            .iter()
            .map(|&index| {
                let judge = std::thread::Builder::new()
                    .stack_size(JUDGE_STACK)
                    .spawn_scoped(scope, move || trial(index));
                (index, judge)
            })
            .collect();
        judges
            .into_iter()
            .map(|(index, judge)| verdict(judge).unwrap_or_else(|| trial(index)))
            .collect()
    })
}

/// A judge's answer, or `None` when the host could not start it and the
/// candidate must be judged here instead. A judge that panicked re-raises
/// here, as the solve would have on this thread.
fn verdict(judge: std::io::Result<std::thread::ScopedJoinHandle<'_, bool>>) -> Option<bool> {
    judge.ok().map(|handle| {
        handle
            .join()
            .unwrap_or_else(|panic| std::panic::resume_unwind(panic))
    })
}

/// Whether erasing every annotation in `chosen` — all of them at once — leaves
/// `program` typechecking and publishing the same types, up to renaming of
/// type variables.
///
/// One checked solve both validates the candidate and publishes its types.
fn preserves(program: &Program, baseline: &Snapshot, written: &[Site], chosen: &[usize]) -> bool {
    let erasing: HashSet<usize> = chosen
        .iter()
        .filter_map(|index| written.get(*index))
        .flat_map(|site| site.members.iter().copied())
        .collect();
    let candidate = erase(program, &mut |slot| !erasing.contains(&slot));
    published(&candidate).is_some_and(|candidate| {
        baseline.methods == candidate.methods && same_types(&baseline.types, &candidate.types)
    })
}

/// One judge per core the host offers.
fn judges() -> usize {
    std::thread::available_parallelism().map_or(1, std::num::NonZeroUsize::get)
}

/// Everything the inferrer resolved, labelled so two runs compare entry by
/// entry.
pub(crate) fn published(program: &Program) -> Option<Snapshot> {
    let types = infer_checked(program).ok()?;
    let functions = types.functions.iter().map(|(name, (params, ret))| {
        let signature = Type::Fun {
            params: params.clone(),
            ret: Box::new(ret.clone()),
        };
        (format!("fn {name}"), signature)
    });
    let mut published: Published = functions
        .chain(sited("let", &types.lets))
        .chain(sited("lambda", &types.lambdas))
        .chain(sited("list", &types.lists))
        .collect();
    for (function, binders) in &types.declared_params {
        for (name, ty) in binders {
            let _ = published.insert(format!("binder {function} {name}"), ty.clone());
        }
    }
    for (owner, obligations) in &types.obligations {
        for (index, (name, ty)) in obligations
            .iter()
            .filter(|(name, ty)| {
                !name.starts_with(crate::builtin_constraints::NUMERIC_OPERAND_PREFIX)
                    || !crate::builtin_constraints::is_numeric_scalar(ty)
            })
            .enumerate()
        {
            let _ = published.insert(format!("constraint {owner} {index} {name}"), ty.clone());
        }
    }
    for (position, site) in &types.performs {
        publish_operation(
            &mut published,
            &format!("perform {position:?}"),
            &site.op,
            &site.effect_args,
        );
    }
    for (position, site) in &types.handler_ops {
        for (operation, op) in &site.ops {
            publish_operation(
                &mut published,
                &format!("handler {position:?} {operation}"),
                op,
                &site.effect_args,
            );
        }
    }
    Some(Snapshot {
        types: published,
        methods: types.methods,
    })
}

/// A callback rewrite deletes only its wrapper lambdas. Other inferred sites,
/// method choices and entry effect requirements must remain valid unchanged.
pub(crate) fn preserves_callbacks(
    baseline: &Snapshot,
    candidate: &Program,
    removed: &[Position],
) -> bool {
    let Some(candidate) = published(candidate) else {
        return false;
    };
    let mut types = baseline.types.clone();
    for position in removed {
        let _ = types.remove(&format!("lambda {}:{}", position.line, position.column));
    }
    ordered_methods(&baseline.methods) == ordered_methods(&candidate.methods)
        && same_types(&types, &candidate.types)
}

fn ordered_methods(methods: &crate::methods::Targets) -> Vec<&crate::methods::Target> {
    let ordered: BTreeMap<_, _> = methods.iter().collect();
    ordered.into_values().collect()
}

fn publish_operation(
    published: &mut Published,
    label: &str,
    operation: &crate::info::OpType,
    arguments: &[Type],
) {
    let _ = published.insert(
        label.to_owned(),
        Type::fun(operation.params.clone(), operation.ret.clone()),
    );
    for (index, argument) in arguments.iter().enumerate() {
        let _ = published.insert(format!("{label} argument {index}"), argument.clone());
    }
}

/// Label one position-keyed table so its entries join the same flat map.
fn sited<'a>(
    kind: &'a str,
    table: &'a crate::info::SiteTable,
) -> impl Iterator<Item = (String, Type)> + 'a {
    table
        .iter()
        .map(move |((line, column), ty)| (format!("{kind} {line}:{column}"), ty.clone()))
}

/// The diagnostic for one redundant annotation.
///
/// The reported type is the written one: redundancy *means* inference derives
/// exactly that, so printing it tells the reader the deletion is a no-op.
fn warn(site: &Site) -> TypeWarning {
    let written = &site.written;
    let message = match &site.slot {
        Slot::Signature { owner } => format!(
            "redundant type signature on `{owner}`: inference derives `{written}` without it"
        ),
        Slot::Param {
            owner, parameter, ..
        } => format!(
            "redundant type annotation on parameter `{parameter}` of `{owner}`: inference derives `{written}` without it"
        ),
        Slot::Return { owner } => format!(
            "redundant return type annotation on `{owner}`: inference derives `{written}` without it"
        ),
        Slot::Binding { name } => format!(
            "redundant type annotation on `{name}`: inference derives `{written}` without it"
        ),
    };
    TypeWarning {
        message: osprey_ast::symbol::demangle_message(&message).into_owned(),
        position: site.position,
        rule: REDUNDANT_ANNOTATION,
    }
}

/// Whether two runs published the same types up to a consistent renaming of
/// inference variables.
///
/// The renaming is global rather than per entry: a shared variable is real
/// sharing between two signatures, and an erasure that breaks it changed the
/// program's types even though each signature alone still looks the same.
fn same_types(left: &Published, right: &Published) -> bool {
    let mut renaming = Renaming::default();
    left.len() == right.len()
        && left
            .iter()
            .zip(right.iter())
            .all(|((left_label, left_ty), (right_label, right_ty))| {
                left_label == right_label && equivalent(left_ty, right_ty, &mut renaming)
            })
}

/// A partial bijection between the two runs' inference variables.
#[derive(Default)]
struct Renaming {
    /// Left variable → the right variable it was first paired with.
    forward: HashMap<VarId, VarId>,
    /// Right variable → the left variable it was first paired with.
    backward: HashMap<VarId, VarId>,
}

impl Renaming {
    /// Pair `left` with `right`, failing if either is already paired elsewhere.
    fn pair(&mut self, left: VarId, right: VarId) -> bool {
        matches(&mut self.forward, left, right) && matches(&mut self.backward, right, left)
    }
}

/// Whether `key` already maps to `value` in `map`, recording it if it is new.
fn matches(map: &mut HashMap<VarId, VarId>, key: VarId, value: VarId) -> bool {
    match map.insert(key, value) {
        Some(existing) => existing == value,
        None => true,
    }
}

/// Whether two types are equal up to the variable renaming built so far.
fn equivalent(left: &Type, right: &Type, renaming: &mut Renaming) -> bool {
    match (left, right) {
        (Type::Var(l), Type::Var(r)) => renaming.pair(*l, *r),
        (Type::Con { name: l, args: la }, Type::Con { name: r, args: ra })
        | (
            Type::Union {
                name: l,
                variants: la,
            },
            Type::Union {
                name: r,
                variants: ra,
            },
        ) => l == r && all_equivalent(la, ra, renaming),
        (
            Type::Fun {
                params: lp,
                ret: lr,
            },
            Type::Fun {
                params: rp,
                ret: rr,
            },
        ) => all_equivalent(lp, rp, renaming) && equivalent(lr, rr, renaming),
        (
            Type::Record {
                name: l,
                fields: lf,
            },
            Type::Record {
                name: r,
                fields: rf,
            },
        ) => l == r && lf.keys().eq(rf.keys()) && all_equivalent_values(lf, rf, renaming),
        _ => false,
    }
}

/// Whether two type sequences are pairwise equivalent under one renaming.
fn all_equivalent(left: &[Type], right: &[Type], renaming: &mut Renaming) -> bool {
    left.len() == right.len()
        && left
            .iter()
            .zip(right.iter())
            .all(|(l, r)| equivalent(l, r, renaming))
}

/// Whether two records' field types are pairwise equivalent. Fields are keyed
/// in a `BTreeMap`, so iteration already aligns them by name.
fn all_equivalent_values(
    left: &BTreeMap<String, Type>,
    right: &BTreeMap<String, Type>,
    renaming: &mut Renaming,
) -> bool {
    left.values()
        .zip(right.values())
        .all(|(l, r)| equivalent(l, r, renaming))
}
