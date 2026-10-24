//! Effect row mask tests: a request under `n` masks is discharged by the
//! `n+1`th answering activation, and saturation never discharges.
//! Implements [EFFECTS-MASK].
use super::*;

fn get() -> Summary {
    let mut summary = Summary::default();
    let _ = summary
        .required
        .insert(Requirement::new("Store", "get", Vec::new()));
    let _ = summary
        .required
        .insert(Requirement::new("Log", "note", Vec::new()));
    summary
}

fn handled(summary: Summary) -> Summary {
    let operations = HashSet::from(["get".to_owned()]);
    summary.without_operations("Store", &[], &operations, false)
}

fn has_store(summary: &Summary) -> bool {
    summary.required.iter().any(|r| r.effect == "Store")
}

#[test]
fn a_masked_request_passes_one_handler_and_meets_the_next() {
    let masked = get().masked("Store", &[]);
    let once = handled(masked);
    assert!(has_store(&once), "the first handler is skipped");
    assert!(!has_store(&handled(once)), "the second handler answers");
}

#[test]
fn a_mask_leaves_other_effects_alone() {
    let masked = handled(get().masked("Store", &[]));
    assert!(masked.required.iter().any(|r| r.effect == "Log"));
    assert!(masked
        .required
        .iter()
        .all(|r| r.effect != "Log" || r.masked == 0));
}

#[test]
fn a_saturated_request_is_never_discharged() {
    let deep = (0..MASK_LIMIT).fold(get(), |summary, _| summary.masked("Store", &[]));
    let answered = (0..MASK_LIMIT * 2).fold(deep, |summary, _| handled(summary));
    assert!(has_store(&answered));
}

#[test]
fn a_callback_use_replays_its_layers_innermost_first() {
    let mask = Scope::Masked {
        effect: "Store".to_owned(),
        arguments: Vec::new(),
    };
    let answer = Scope::Handled(BTreeMap::from([(
        Requirement::new("Store", "get", Vec::new()),
        1,
    )]));
    assert!(has_store(&get().within(&[mask.clone(), answer.clone()])));
    assert!(!has_store(&get().within(std::slice::from_ref(&answer))));
    let twice = Scope::Handled(BTreeMap::from([(
        Requirement::new("Store", "get", Vec::new()),
        2,
    )]));
    assert!(!has_store(&get().within(&[mask, twice])));
}
