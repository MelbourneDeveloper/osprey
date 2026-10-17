//! Python concrete syntax lowers directly to the shared Osprey AST.
//! Implements [FLAVOR-PYTHON-BOUNDARY], [FLAVOR-PYTHON-SELECT],
//! [FLAVOR-PYTHON-CST], and [FLAVOR-PYTHON-FUNCTIONS].
use osprey_syntax::{parse_program_with_flavor, resolve_flavor, Flavor, Parsed};

#[expect(
    clippy::expect_used,
    reason = "the fixture cannot run until Python selection exists"
)]
fn python() -> Flavor {
    "python"
        .parse()
        .expect("the Python flavor must be selectable")
}

fn parse(source: &str) -> Parsed {
    parse_program_with_flavor(source, python())
}

fn equal(default: &str, source: &str) {
    let expected = parse_program_with_flavor(default, Flavor::Default);
    let actual = parse(source);
    assert!(expected.errors.is_empty(), "{:?}", expected.errors);
    assert!(actual.errors.is_empty(), "{:?}", actual.errors);
    assert_eq!(
        osprey_ast::canonical::without_positions(&actual.program),
        osprey_ast::canonical::without_positions(&expected.program)
    );
}

#[test]
fn python_routing_preserves_explicit_selection_and_marker_conflicts() {
    assert_eq!(python().to_string(), "python");
    assert_eq!(resolve_flavor(None, "model.py", ""), Ok(python()));
    let header = "#!/usr/bin/env osprey\n# coding: utf-8\n# osprey: flavor=python\nx = 1\n";
    assert_eq!(resolve_flavor(None, "model.txt", header), Ok(python()));
    assert!(resolve_flavor(None, "model.osp", header).is_err());
    assert_eq!(
        resolve_flavor(Some(python()), "model.osp", header),
        Ok(python())
    );
}

#[test]
fn python_flat_functions_calls_lambdas_and_bindings_equal_default_ast() {
    equal(
        "fn scale(value, factor) = value * factor\nlet twice = fn(value) => scale(value: value, factor: 2)\nlet result = twice(21)\n",
        "def scale(value, factor):\n    return value * factor\ntwice = lambda value: scale(value=value, factor=2)\nresult = twice(21)\n",
    );
}

#[test]
fn python_block_locals_and_numeric_negation_equal_default_ast() {
    equal(
        "fn answer(x) = { let offset = -2\n x + offset }\nlet smallest = -9223372036854775808\n",
        "def answer(x):\n    offset = -2\n    return x + offset\nsmallest = -9223372036854775808\n",
    );
}

#[test]
fn python_annotations_and_inferred_mutable_cells_equal_default_ast() {
    equal(
        "fn next(value: int) -> int = value\nmut counter: int = 1\ncounter = 2\n",
        "def next(value: int) -> int:\n    return value\ncounter: int = 1\ncounter = 2\n",
    );
}

#[test]
fn python_boolean_and_conditional_spelling_equal_default_ast() {
    equal(
        "fn pick(a, b, yes) = if yes && !false { a } else { b }\n",
        "def pick(a, b, yes):\n    return a if yes and not False else b\n",
    );
}

#[test]
fn python_tail_if_branches_equal_default_matches() {
    equal(
        "fn absolute(x) = if x < 0 { -x } else { x }\n",
        "def absolute(x):\n    if x < 0:\n        return -x\n    else:\n        return x\n",
    );
}

#[test]
fn python_effect_decorators_preserve_operation_metadata() {
    equal(
        "effect Ask<T> { control once replayable ask: fn(T) -> T }\nstatic effect Tensor<B> { size: fn(B) -> int }\n",
        "@effect\nclass Ask[T]:\n    @control\n    @once\n    @replayable\n    def ask(value: T) -> T: ...\n@static_effect\nclass Tensor[B]:\n    def size(buffer: B) -> int: ...\n",
    );
}

#[test]
fn python_requires_preserves_closed_and_open_effect_rows() {
    equal(
        "fn quiet(x) ![] = x\nfn logging(x) ![Log<int> | e] = x\n",
        "@requires()\ndef quiet(x):\n    return x\n@requires(Log[int], tail=\"e\")\ndef logging(x):\n    return x\n",
    );
}

#[test]
fn python_perform_resume_and_handler_values_equal_default_ast() {
    equal(
        "effect Ask { control ask: fn() -> int }\nlet Answer = handler Ask { ask => resume(7) }\nfn work() = perform Ask.ask()\nlet result = Answer(work)\n",
        "@effect\nclass Ask:\n    @control\n    def ask() -> int: ...\n@handler(Ask)\nclass Answer:\n    def ask():\n        return resume(7)\ndef work():\n    return perform(Ask.ask())\nresult = Answer(work)\n",
    );
}

#[test]
fn python_with_handler_is_a_lexical_zero_argument_callback() {
    equal(
        "fn work() = console(fn() => { let message = \"ok\"\n print(message) })\n",
        "def work():\n    with console:\n        message = \"ok\"\n        return print(message)\n",
    );
}

#[test]
fn python_typing_aliases_remain_lexically_scoped() {
    equal(
        "fn inner() = {}\nfn outer(value: Items<int>) = value\n",
        "def inner():\n    from typing import List as Items\n    pass\ndef outer(value: Items[int]):\n    return value\n",
    );
    assert!(
        !parse("from collections.abc import Final\n")
            .errors
            .is_empty(),
        "accepted a typing alias imported from outside `typing`"
    );
}

#[test]
fn python_multiline_string_tabs_are_content() {
    equal(
        "let text = \"line\\n\\tcontent\\n\"\n",
        "text = \"\"\"line\n\tcontent\n\"\"\"\n",
    );
    assert!(
        !parse("def f():\n\treturn 1\n").errors.is_empty(),
        "accepted a tab-indented body"
    );
}

#[test]
fn python_result_and_list_patterns_equal_default_ast() {
    equal(
        "fn read(text) = match parseInt(text) { Success { value } => value\n Error { message } => 0 }\nfn first(xs) = match xs { [head, ...tail] => head\n [] => 0 }\n",
        "def read(text):\n    match parseInt(text):\n        case Success(value=value):\n            return value\n        case Error(message=message):\n            return 0\ndef first(xs):\n    match xs:\n        case [head, *tail]:\n            return head\n        case []:\n            return 0\n",
    );
}

#[test]
fn python_fiber_intrinsics_and_updates_equal_default_ast() {
    equal(
        "fn run(work, channel, value) = { let task = spawn work()\n send(channel, value)\n yield\n let received = recv(channel)\n await(task) }\nlet newer = point { x: 2 }\n",
        "def run(work, channel, value):\n    task = spawn(work)\n    send(channel, value)\n    yield_fiber()\n    received = recv(channel)\n    return await_fiber(task)\nnewer = update(point, x=2)\n",
    );
}

#[test]
fn python_partial_tail_returns_are_rejected_even_for_unit() {
    for source in [
        "def f(x):\n    if x:\n        return None\n",
        "def f(x):\n    if x:\n        return\n    else:\n        pass\n",
        "def f(x):\n    match x:\n        case True:\n            return\n        case False:\n            pass\n",
    ] {
        assert!(!parse(source).errors.is_empty(), "accepted partial return: {source}");
    }
}

#[test]
fn python_intrinsics_reject_invalid_arity_and_context_managers() {
    for source in [
        "perform()\n",
        "resume(1, 2)\n",
        "spawn()\n",
        "await_fiber()\n",
        "yield_fiber(1, 2)\n",
        "send(1)\n",
        "recv()\n",
        "update(1, x=2)\n",
        "with handler as value:\n    pass\n",
        "with first, second:\n    pass\n",
    ] {
        assert!(
            !parse(source).errors.is_empty(),
            "accepted invalid intrinsic: {source}"
        );
    }
}

#[test]
fn python_file_directives_imports_and_exports_equal_default_ast() {
    equal(
        "namespace gpt2;\nimport core::Tensor::{build as make}\nimport core::Model as Model\nmodule PythonModel { export let WIDTH = 768\nexport fn forward(x) = x\nfn privateHelper(x) = x }\n",
        "# osprey: namespace=gpt2\n# osprey: module=PythonModel\nfrom core.Tensor import build as make\nimport core.Model as Model\n__all__ = [\"WIDTH\"]\nWIDTH = 768\n@export\ndef forward(x):\n    return x\ndef privateHelper(x):\n    return x\n",
    );
}

#[test]
fn imported_python_qualifiers_lower_to_paths_and_respect_local_shadowing() {
    equal(
        "import core::Model as Model\nfn load() = Model::load()\nfn read(Model) = Model.value\n",
        "import core.Model as Model\ndef load():\n    return Model.load()\ndef read(Model):\n    return Model.value\n",
    );
}

#[test]
fn python_qualified_record_annotations_and_effect_mentions_equal_default_ast() {
    equal(
        "import data::Values as V\nfn keep(value: V::Box<int>) -> V::Box<int> ![V::Read<int>] = value\n",
        "import data.Values as V\n@requires(V.Read[int])\ndef keep(value: V.Box[int]) -> V.Box[int]:\n    return value\n",
    );
}

#[test]
fn python_docstrings_lower_to_the_shared_documentation_model() {
    equal(
        "//! File documentation.\n/// A value.\ntype Value = { number: int }\n/// Returns its argument.\nfn identity(x) = x\n",
        "\"\"\"File documentation.\"\"\"\nclass Value:\n    \"\"\"A value.\"\"\"\n    number: int\ndef identity(x):\n    \"\"\"Returns its argument.\"\"\"\n    return x\n",
    );
}

#[test]
fn python_effect_handler_and_foreign_signature_docs_preserve_export_metadata() {
    equal(
        "namespace app;\nmodule Api {\n/// A capability.\nexport effect Log {\n/// Writes a message.\nwrite: fn(string) -> Unit }\n/// A console.\nlet Console = handler Log { write message => print(message) }\n/// A foreign call.\nexport extern fn foreign(value: int) -> int\n}\n",
        "# osprey: namespace=app\n# osprey: module=Api\n@export\n@effect\nclass Log:\n    \"\"\"A capability.\"\"\"\n    def write(message: str) -> None:\n        \"\"\"Writes a message.\"\"\"\n        ...\n@handler(Log)\nclass Console:\n    \"\"\"A console.\"\"\"\n    def write(message):\n        return print(message)\n@extern\n@export\ndef foreign(value: int) -> int:\n    \"\"\"A foreign call.\"\"\"\n    ...\n",
    );
}

#[test]
fn python_file_metadata_rejects_late_duplicate_or_invalid_declarations() {
    for source in [
        "x = 1\n# osprey: namespace=late\n",
        "# osprey: namespace=one\n# osprey: namespace=two\n",
        "# osprey: flavor=python\n# osprey: flavor=python\n",
        "# osprey: module=Alone\n",
        "# osprey: namespace=bad-name\n",
        "# osprey: namespace=one\n# osprey: module=Two\n__all__ = [\"missing\"]\n",
        "# osprey: namespace=one\n# osprey: module=Two\n__all__ = [\"x\", \"x\"]\nx = 1\n",
        "@export\ndef alone():\n    return 1\n",
        "from .relative import thing\n",
        "from logical import *\n",
    ] {
        assert!(
            !parse(source).errors.is_empty(),
            "accepted invalid metadata: {source}"
        );
    }
}

#[test]
fn python_builtin_aliases_lists_maps_and_division_equal_default_ast() {
    equal(
        "let xs = [1, 2]\nlet n = length(xs)\nlet values = range(0, n)\nlet q = intDiv(-7, 3)\nlet mapping = {\"x\": 2}\n",
        "xs = [1, 2]\nn = len(xs)\nvalues = range(n)\nq = -7 // 3\nmapping = {\"x\": 2}\n",
    );
}

#[test]
fn python_fstrings_retain_expression_structure_and_source_positions() {
    let parsed = parse("def label(value):\n    return f\"λ={value + 1}\"\n");
    assert!(parsed.errors.is_empty(), "{:?}", parsed.errors);
    equal(
        "let text = \"value=${1 + 2}\"\n",
        "text = f\"value={1 + 2}\"\n",
    );
    let dump = format!("{:#?}", parsed.program);
    assert!(dump.contains("line: 2"), "{dump}");
    assert!(dump.contains("column: 17"), "{dump}");
}

#[test]
fn python_rejects_unsupported_semantics_at_original_source_positions() {
    for source in [
        "x = ...\n",
        "x = (1, 2)\n",
        "while True:\n    pass\n",
        "def f(x=1):\n    return x\n",
        "def f(*xs):\n    return xs\n",
        "from typing import Final\nx: Final = 1\nx = 2\n",
        "x = [1]\nx[0] = 2\n",
        "def f(x):\n    return x\n    print(x)\n",
        "@unknown\ndef f():\n    pass\n",
        "x = 1 < 2 < 3\n",
        "x = 9223372036854775808\n",
        "x = 1e999\n",
        "async def f():\n    return 1\n",
    ] {
        let parsed = parse(source);
        assert!(!parsed.errors.is_empty(), "silently accepted {source}");
        assert!(parsed.errors.iter().all(|error| error.position.line > 0));
    }
}

#[test]
fn python_reports_malformed_indentation_and_incomplete_cst_nodes() {
    for source in [
        "def f(:\n",
        "def f():\n\treturn 1\n",
        "def f():\n    x = 1\n  return x\n",
    ] {
        let parsed = parse(source);
        assert!(!parsed.errors.is_empty(), "silently accepted {source}");
    }
}

#[test]
fn python_closed_record_classes_and_generic_constructors_equal_default_ast() {
    // [FLAVOR-PYTHON-CLASS]
    equal(
        "type Box<T> = { value: T }\nlet box = Box<int> { value: 42 }\n",
        "class Box[T]:\n    value: T\nbox = Box[int](value=42)\n",
    );
}

#[test]
fn python_class_function_signatures_declare_stored_fields() {
    // [FLAVOR-PYTHON-CLASS]
    equal(
        "type Price = { value: int, discount: (int) -> int }\n",
        "class Price:\n    value: int\n    def discount(amount: int) -> int: ...\n",
    );
}

#[test]
fn python_extern_signatures_equal_the_checked_c_boundary() {
    // [FLAVOR-PYTHON-FUNCTIONS] The interop agent reuses this exact boundary.
    equal(
        "extern fn bridge(operation: string, argument: int) -> string\n",
        "@extern\ndef bridge(operation: str, argument: int) -> str: ...\n",
    );
}

#[test]
fn python_classes_and_externs_reject_python_object_or_dynamic_decorator_behavior() {
    for source in [
        "class A(Base):\n    value: int\n",
        "class A:\n    value: int = 1\n",
        "class A:\n    def method(self):\n        return 1\n",
        "class A:\n    def method(x): ...\n",
        "@extern\ndef f(x): ...\n",
        "@extern\ndef f(x: int) -> int:\n    return x\n",
        "@python(\"math\", \"exp\")\ndef exp(x: float) -> float: ...\n",
    ] {
        assert!(
            !parse(source).errors.is_empty(),
            "silently accepted {source}"
        );
    }
}

#[test]
fn typed_python_standard_imports_and_callable_annotations_equal_default_ast() {
    // [FLAVOR-PYTHON-TYPING] PEP 484 and PEP 585 vocabulary is compile-time syntax.
    equal("fn apply(f: fn(int) -> int, value: int) -> int = f(value)\n",
        "from typing import Callable\ndef apply(f: Callable[[int], int], value: int) -> int:\n    return f(value)\n");
    equal(
        "fn size(xs: List<int>) -> int = length(xs)\n",
        "from typing import List as Items\ndef size(xs: Items[int]) -> int:\n    return len(xs)\n",
    );
}

#[test]
fn typed_python_none_spells_unit_and_pep695_aliases_remain_osprey_types() {
    // [FLAVOR-PYTHON-TYPING] PEP 484 None and PEP 695 aliases, no null object.
    equal(
        "fn noop() -> Unit = {}\ntype Label = string\n",
        "def noop() -> None:\n    return None\ntype Label = str\n",
    );
}

#[test]
fn typed_python_reassignment_and_final_use_existing_binding_kinds() {
    equal(
        "mut count = 1\ncount = 2\nlet limit: int = 3\n",
        "from typing import Final\ncount = 1\ncount = 2\nlimit: Final[int] = 3\n",
    );
    let invalid = parse("from typing import Final\nlimit: Final[int] = 3\nlimit = 4\n");
    assert!(
        invalid
            .errors
            .iter()
            .any(|error| error.message.contains("limit") && error.message.contains("Final")),
        "{:?}",
        invalid.errors
    );
}
