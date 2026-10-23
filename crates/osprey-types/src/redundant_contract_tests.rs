//! An exported module member's signature is its API: never reported, never
//! erased while the rest is judged ([TYPE-ANNOTATION-REDUNDANT]).

use crate::redundant::redundant_annotation_sites_for;
use osprey_syntax::{parse_program_with_flavor, Flavor};

const TWO: &str = "fn published(name: string) -> string = \"hi \" + name\n\
    fn internal(name: string) -> string = {\n  let shout: string = name + \"!\"\n  shout\n}\n\
    fn api(label: string) -> string = {\n  let inner: string = label + \"?\"\n  inner\n}\n";

fn reported(contract: &dyn Fn(&str) -> bool) -> Vec<String> {
    let program = parse_program_with_flavor(TWO, Flavor::Default).program;
    redundant_annotation_sites_for(&program, |_| true, contract)
        .into_iter()
        .map(|site| site.warning.message)
        .collect()
}

#[test]
fn an_exported_members_signature_is_never_reported_but_its_body_still_is() {
    let all = reported(&|_| false);
    let contracted = reported(&|name| name == "published" || name == "api");
    assert!(all.iter().any(|m| m.contains("`published`")), "{all:?}");
    assert!(
        contracted.iter().all(|m| !m.contains("of `published`")
            && !m.contains("on `published`")
            && !m.contains("of `api`")
            && !m.contains("on `api`")),
        "{contracted:?}"
    );
    assert!(
        contracted.iter().any(|m| m.contains("`internal`")),
        "{contracted:?}"
    );
    assert!(
        contracted.iter().any(|m| m.contains("`inner`")),
        "{contracted:?}"
    );
}
