use super::*;
use osprey_syntax::parse_program;

fn program(source: &str) -> Program {
    let parsed = parse_program(source);
    assert!(parsed.errors.is_empty(), "syntax: {:?}", parsed.errors);
    parsed.program
}

fn attached_names(source: &str) -> Vec<String> {
    let source = program(source);
    with_prelude(&source)
        .statements
        .iter()
        .filter_map(|statement| match statement {
            Stmt::Type { name, .. } | Stmt::Function { name, .. } => Some(name.clone()),
            _ => None,
        })
        .collect()
}

#[test]
fn the_prelude_source_parses_clean_and_lives_past_every_source_line() {
    let parsed = parse_program(JSON);
    assert!(
        parsed.errors.is_empty(),
        "prelude syntax: {:?}",
        parsed.errors
    );
    let rebased = prelude();
    assert_eq!(rebased.statements.len(), parsed.program.statements.len());
    assert!(attached(rebased), "the prelude must recognise itself");
    let lines: Vec<u32> = rebased
        .statements
        .iter()
        .filter_map(type_position)
        .map(|p| p.line)
        .collect();
    assert!(
        !lines.is_empty() && lines.iter().all(|line| is_prelude_line(*line)),
        "{lines:?}"
    );
    assert!(!is_prelude_line(PRELUDE_LINE_BASE) && is_prelude_line(PRELUDE_LINE_BASE + 1));
}

#[test]
fn a_program_that_names_nothing_from_the_prelude_is_left_exactly_as_written() {
    let source = program("fn jsonish(x) = x\nprint(toString(jsonish(1)))\n");
    assert!(matches!(with_prelude(&source), Cow::Borrowed(_)));
    // Naming an internal primitive does not pull the prelude in either.
    let internal = program("print(toString(__jsonNodeKind(1, 0)))\n");
    assert!(matches!(with_prelude(&internal), Cow::Borrowed(_)));
}

#[test]
fn naming_any_offered_declaration_attaches_the_whole_prelude_once() {
    for source in [
        "print(jsonEncode(JsonNull) ?: \"\")\n",
        "let decoded = jsonDecode(\"[]\")\n",
        "let member = JsonMember(\"k\", JsonInt(1))\n",
        "fn f(v) = jsonText(v)\n",
    ] {
        let names = attached_names(source);
        for offered in [
            "Json",
            "JsonMember",
            "jsonEncode",
            "jsonDecode",
            "jsonQuote",
            "jsonField",
            "jsonIndex",
            "jsonText",
            "jsonInt",
            "jsonFloat",
            "jsonBool",
        ] {
            assert_eq!(
                names.iter().filter(|name| *name == offered).count(),
                1,
                "{offered} in {source}"
            );
        }
        let once = with_prelude(&program(source)).into_owned();
        assert!(
            matches!(with_prelude(&once), Cow::Borrowed(_)),
            "attaching twice: {source}"
        );
    }
}

#[test]
fn a_program_declaring_a_prelude_type_or_constructor_keeps_its_own_and_gets_none() {
    for source in [
        "type JsonMember = JsonMember(string, int)\nprint(jsonEncode(JsonNull) ?: \"\")\n",
        "type Json = { raw: string }\nlet j = Json { raw: \"x\" }\n",
        "type Mine = JsonNull | Other\nlet m = JsonNull\n",
    ] {
        let source = program(source);
        assert!(matches!(with_prelude(&source), Cow::Borrowed(_)));
    }
}

#[test]
fn a_program_function_replaces_only_the_prelude_function_it_names() {
    let names = attached_names("fn jsonField(a, b) = a + b\nprint(jsonEncode(JsonNull) ?: \"\")\n");
    assert_eq!(names.iter().filter(|name| *name == "jsonField").count(), 1);
    assert_eq!(
        names.last().map(String::as_str),
        Some("jsonField"),
        "the program's own copy stays"
    );
    assert!(
        names.iter().any(|name| name == "jsonQuote")
            && names.iter().any(|name| name == "__jsonMember")
    );
}

#[test]
fn the_prelude_type_checks_clean_beside_a_program_that_uses_all_of_it() {
    let source = program(
        "let doc = JsonObject([JsonMember(\"a\", JsonArray([JsonInt(1), JsonFloat(2.5), JsonBool(true), JsonNull]))])\n\
         let text = jsonEncode(doc) ?: \"\"\n\
         let back = jsonDecode(text) ?: JsonNull\n\
         let first = jsonField(back, \"a\") ?: JsonNull\n\
         let one = jsonIndex(first, 0) ?: JsonNull\n\
         print(\"${text} ${jsonQuote(\"q\") ?: \"\"} ${toString(jsonInt(one) ?: 0)} ${toString(jsonFloat(one) ?: 0.0)} ${toString(jsonBool(JsonBool(true)) ?: false)} ${jsonText(JsonString(\"s\")) ?: \"\"}\")\n",
    );
    assert_eq!(crate::check_program(&source), Vec::new());
}
