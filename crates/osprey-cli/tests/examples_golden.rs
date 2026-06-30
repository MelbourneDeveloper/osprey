//! Per-fixture golden tests for the published `examples/tested` corpus.
//!
//! The shell harness remains the compact CI oracle, but every example also has
//! a named Rust test here so VS Code's Rust test tree can discover and run the
//! fixtures one by one.

use std::collections::BTreeSet;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::Command;
use std::sync::{Mutex, MutexGuard, OnceLock};

#[derive(Debug)]
struct Out {
    code: Option<i32>,
    stdout: String,
    stderr: String,
}

fn repo_root() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("..").join("..")
}

fn example_lock() -> MutexGuard<'static, ()> {
    let lock = EXAMPLE_LOCK.get_or_init(|| Mutex::new(()));
    match lock.lock() {
        Ok(guard) => guard,
        Err(poisoned) => poisoned.into_inner(),
    }
}

fn osprey() -> Command {
    let mut cmd = Command::new(env!("CARGO_BIN_EXE_osprey"));
    let _ = cmd.current_dir(repo_root());
    cmd
}

fn run_example(source: &Path) -> Result<Out, String> {
    let mut cmd = osprey();
    let _ = cmd.arg(source).arg("--run");
    cmd.output()
        .map(|out| Out {
            code: out.status.code(),
            stdout: String::from_utf8_lossy(&out.stdout).into_owned(),
            stderr: String::from_utf8_lossy(&out.stderr).into_owned(),
        })
        .map_err(|e| format!("spawn failed for {}: {e}", source.display()))
}

fn read_text(path: &Path) -> Result<String, String> {
    fs::read_to_string(path).map_err(|e| format!("cannot read {}: {e}", path.display()))
}

fn path_with_suffix(path: &Path, suffix: &str) -> PathBuf {
    PathBuf::from(format!("{}{}", path.display(), suffix))
}

fn source_base(source: &Path) -> PathBuf {
    let mut base = source.to_path_buf();
    assert!(
        base.set_extension(""),
        "example path has no extension: {}",
        source.display()
    );
    base
}

fn uname_s() -> &'static str {
    match std::env::consts::OS {
        "macos" => "Darwin",
        "linux" => "Linux",
        "windows" => "Windows_NT",
        other => other,
    }
}

fn expected_candidates(source: &Path) -> Vec<PathBuf> {
    let os = uname_s();
    let base = source_base(source);
    vec![
        path_with_suffix(source, ".expectedoutput"),
        path_with_suffix(source, &format!(".expectedoutput.{os}")),
        path_with_suffix(&base, ".osp.expectedoutput"),
        path_with_suffix(&base, &format!(".osp.expectedoutput.{os}")),
        path_with_suffix(&base, ".expectedoutput"),
    ]
}

fn expected_output_path(source: &Path) -> Result<PathBuf, String> {
    for candidate in expected_candidates(source) {
        if candidate.is_file() {
            return Ok(candidate);
        }
    }
    Err(format!("missing expected output for {}", source.display()))
}

fn check_example_matches(rel_source: &str) -> Result<(), String> {
    let _guard = example_lock();
    let source = repo_root().join(rel_source);
    let expected_path = expected_output_path(&source)?;
    let expected = read_text(&expected_path)?;
    let actual = run_example(&source)?;

    if actual.code == Some(0) && actual.stdout.trim() == expected.trim() {
        return Ok(());
    }

    Err(format!(
        "{rel_source}\nstatus={:?}\nexpected file={}\n--- expected ---\n{}\n--- actual ---\n{}\n--- stderr ---\n{}",
        actual.code,
        expected_path.display(),
        expected.trim(),
        actual.stdout.trim(),
        actual.stderr.trim()
    ))
}

fn assert_example_matches(rel_source: &str) {
    if let Err(e) = check_example_matches(rel_source) {
        assert!(e.is_empty(), "{e}");
    }
}

static EXAMPLE_LOCK: OnceLock<Mutex<()>> = OnceLock::new();

fn collect_sources(dir: &Path, out: &mut Vec<PathBuf>) {
    let Ok(entries) = fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect_sources(&path, out);
        } else if matches!(
            path.extension().and_then(|ext| ext.to_str()),
            Some("osp" | "ospml")
        ) {
            out.push(path);
        }
    }
}

fn tested_example_sources() -> Vec<PathBuf> {
    let mut out = Vec::new();
    collect_sources(&repo_root().join("examples/tested"), &mut out);
    out.sort();
    out
}

fn repo_relative(path: &Path) -> String {
    let root = repo_root();
    match path.strip_prefix(&root) {
        Ok(rel) => rel.to_string_lossy().replace('\\', "/"),
        Err(_) => path.to_string_lossy().replace('\\', "/"),
    }
}

#[test]
fn all_tested_examples_are_registered_as_individual_tests() {
    let discovered = tested_example_sources()
        .iter()
        .map(|path| repo_relative(path))
        .collect::<BTreeSet<_>>();
    let registered = REGISTERED_EXAMPLES
        .iter()
        .map(|path| (*path).to_string())
        .collect::<BTreeSet<_>>();

    let missing = discovered
        .difference(&registered)
        .cloned()
        .collect::<Vec<_>>();
    let stale = registered
        .difference(&discovered)
        .cloned()
        .collect::<Vec<_>>();

    assert!(
        missing.is_empty() && stale.is_empty(),
        "every examples/tested fixture must have a named Rust test\nmissing:\n{}\nstale:\n{}",
        missing.join("\n"),
        stale.join("\n")
    );
}

const REGISTERED_EXAMPLES: &[&str] = &[
    "examples/tested/basics/blocks/block_statements_basic.osp",
    "examples/tested/basics/blocks/block_statements_basic.ospml",
    "examples/tested/basics/cursor/codepoint_roundtrip.osp",
    "examples/tested/basics/cursor/codepoint_roundtrip.ospml",
    "examples/tested/basics/cursor/kv_parser.osp",
    "examples/tested/basics/cursor/kv_parser.ospml",
    "examples/tested/basics/cursor/token_scan.osp",
    "examples/tested/basics/cursor/token_scan.ospml",
    "examples/tested/basics/cursor/utf8_walk.osp",
    "examples/tested/basics/cursor/utf8_walk.ospml",
    "examples/tested/basics/errors/error_messages.osp",
    "examples/tested/basics/errors/error_messages.ospml",
    "examples/tested/basics/errors/validation_pipeline.osp",
    "examples/tested/basics/errors/validation_pipeline.ospml",
    "examples/tested/basics/feature_omnibus.osp",
    "examples/tested/basics/feature_omnibus.ospml",
    "examples/tested/basics/field_access_comprehensive.osp",
    "examples/tested/basics/field_access_comprehensive.ospml",
    "examples/tested/basics/files/file_io_json_workflow.osp",
    "examples/tested/basics/files/file_io_json_workflow.ospml",
    "examples/tested/basics/function_composition_test.osp",
    "examples/tested/basics/functional/functional_showcase.osp",
    "examples/tested/basics/functional/functional_showcase.ospml",
    "examples/tested/basics/games/adventure_game.osp",
    "examples/tested/basics/games/adventure_game.ospml",
    "examples/tested/basics/games/space_trader.osp",
    "examples/tested/basics/games/space_trader.ospml",
    "examples/tested/basics/knownbugs/bug1_spawn_record.osp",
    "examples/tested/basics/knownbugs/bug1_spawn_record.ospml",
    "examples/tested/basics/knownbugs/bug2_string_union_payload.osp",
    "examples/tested/basics/knownbugs/bug2_string_union_payload.ospml",
    "examples/tested/basics/knownbugs/bug3_map_built_index.osp",
    "examples/tested/basics/knownbugs/bug3_map_built_index.ospml",
    "examples/tested/basics/knownbugs/bug4_union_return_arg.osp",
    "examples/tested/basics/knownbugs/bug4_union_return_arg.ospml",
    "examples/tested/basics/lists/list_basics.osp",
    "examples/tested/basics/lists/list_basics.ospml",
    "examples/tested/basics/lists/map_basics.osp",
    "examples/tested/basics/lists/map_basics.ospml",
    "examples/tested/basics/math/comprehensive_math.osp",
    "examples/tested/basics/math/comprehensive_math.ospml",
    "examples/tested/basics/operators/boolean_consolidated.osp",
    "examples/tested/basics/operators/boolean_consolidated.ospml",
    "examples/tested/basics/osprey_mega_showcase.osp",
    "examples/tested/basics/osprey_mega_showcase.ospml",
    "examples/tested/basics/pattern_matching/pattern_matching_complete.osp",
    "examples/tested/basics/processes/async_process_management.osp",
    "examples/tested/basics/processes/async_process_management.ospml",
    "examples/tested/basics/processes/callback_stdout_demo.osp",
    "examples/tested/basics/processes/callback_stdout_demo.ospml",
    "examples/tested/basics/strings/string_edge_cases.osp",
    "examples/tested/basics/strings/string_edge_cases.ospml",
    "examples/tested/basics/strings/string_pipeline.osp",
    "examples/tested/basics/strings/string_pipeline.ospml",
    "examples/tested/basics/types/any_type_comprehensive.osp",
    "examples/tested/basics/types/any_type_comprehensive.ospml",
    "examples/tested/basics/types/pure_hindley_milner_test.osp",
    "examples/tested/basics/types/pure_hindley_milner_test.ospml",
    "examples/tested/basics/types/record_update_basic.osp",
    "examples/tested/basics/types/record_update_basic.ospml",
    "examples/tested/basics/types/recursive_unions.osp",
    "examples/tested/basics/types/recursive_unions.ospml",
    "examples/tested/basics/types/type_equality_comprehensive.osp",
    "examples/tested/basics/types/type_equality_comprehensive.ospml",
    "examples/tested/basics/types/user_defined_unions.osp",
    "examples/tested/basics/types/user_defined_unions.ospml",
    "examples/tested/basics/validation/proper_validation_test.osp",
    "examples/tested/basics/validation/proper_validation_test.ospml",
    "examples/tested/basics/website/website_examples.osp",
    "examples/tested/basics/website/website_examples.ospml",
    "examples/tested/db/database_effect.osp",
    "examples/tested/db/database_effect.ospml",
    "examples/tested/db/sqlite_basics.osp",
    "examples/tested/db/sqlite_basics.ospml",
    "examples/tested/effects/algebraic_effects_comprehensive.osp",
    "examples/tested/effects/algebraic_effects_comprehensive.ospml",
    "examples/tested/effects/fiber_effects.osp",
    "examples/tested/effects/fiber_effects.ospml",
    "examples/tested/effects/handler_scoping.osp",
    "examples/tested/effects/handler_scoping.ospml",
    "examples/tested/effects/http_state_levels.osp",
    "examples/tested/effects/http_state_levels.ospml",
    "examples/tested/effects/resume_abort_early_exit.osp",
    "examples/tested/effects/resume_abort_early_exit.ospml",
    "examples/tested/effects/resume_lifo_audit.osp",
    "examples/tested/effects/resume_lifo_audit.ospml",
    "examples/tested/effects/resume_outer_handler_bridge.osp",
    "examples/tested/effects/resume_outer_handler_bridge.ospml",
    "examples/tested/effects/resume_unit_markers.osp",
    "examples/tested/effects/resume_unit_markers.ospml",
    "examples/tested/effects/resume_value_rewrite.osp",
    "examples/tested/effects/resume_value_rewrite.ospml",
    "examples/tested/fiber/fiber_showcase.osp",
    "examples/tested/fiber/fiber_showcase.ospml",
    "examples/tested/http/http_client_example.osp",
    "examples/tested/http/http_client_example.ospml",
    "examples/tested/http/http_create_client.osp",
    "examples/tested/http/http_create_client.ospml",
    "examples/tested/http/http_response_handle.osp",
    "examples/tested/http/http_response_handle.ospml",
    "examples/tested/http/http_server_example.osp",
    "examples/tested/http/http_server_example.ospml",
    "examples/tested/http/tui_repo_table.osp",
    "examples/tested/http/tui_repo_table.ospml",
    "examples/tested/ml/arith.osp",
    "examples/tested/ml/arith.ospml",
    "examples/tested/ml/booleans.osp",
    "examples/tested/ml/booleans.ospml",
    "examples/tested/ml/closures.osp",
    "examples/tested/ml/closures.ospml",
    "examples/tested/ml/curry_partial.osp",
    "examples/tested/ml/curry_partial.ospml",
    "examples/tested/ml/curry_tour.osp",
    "examples/tested/ml/curry_tour.ospml",
    "examples/tested/ml/hello.osp",
    "examples/tested/ml/hello.ospml",
    "examples/tested/ml/hof.osp",
    "examples/tested/ml/hof.ospml",
    "examples/tested/ml/match_tour.osp",
    "examples/tested/ml/match_tour.ospml",
    "examples/tested/ml/matchbool.osp",
    "examples/tested/ml/matchbool.ospml",
    "examples/tested/ml/matchint.osp",
    "examples/tested/ml/matchint.ospml",
    "examples/tested/ml/mixed.osp",
    "examples/tested/ml/mixed.ospml",
    "examples/tested/ml/mutation.osp",
    "examples/tested/ml/mutation.ospml",
    "examples/tested/ml/nested_calls.osp",
    "examples/tested/ml/nested_calls.ospml",
    "examples/tested/ml/pipechain.osp",
    "examples/tested/ml/pipechain.ospml",
    "examples/tested/ml/recursion.osp",
    "examples/tested/ml/recursion.ospml",
    "examples/tested/ml/results_state_hof.osp",
    "examples/tested/ml/results_state_hof.ospml",
    "examples/tested/ml/strings.osp",
    "examples/tested/ml/strings.ospml",
];

#[test]
fn basics_blocks_block_statements_basic_osp() {
    assert_example_matches("examples/tested/basics/blocks/block_statements_basic.osp");
}

#[test]
fn basics_blocks_block_statements_basic_ospml() {
    assert_example_matches("examples/tested/basics/blocks/block_statements_basic.ospml");
}

#[test]
fn basics_cursor_codepoint_roundtrip_osp() {
    assert_example_matches("examples/tested/basics/cursor/codepoint_roundtrip.osp");
}

#[test]
fn basics_cursor_codepoint_roundtrip_ospml() {
    assert_example_matches("examples/tested/basics/cursor/codepoint_roundtrip.ospml");
}

#[test]
fn basics_cursor_kv_parser_osp() {
    assert_example_matches("examples/tested/basics/cursor/kv_parser.osp");
}

#[test]
fn basics_cursor_kv_parser_ospml() {
    assert_example_matches("examples/tested/basics/cursor/kv_parser.ospml");
}

#[test]
fn basics_cursor_token_scan_osp() {
    assert_example_matches("examples/tested/basics/cursor/token_scan.osp");
}

#[test]
fn basics_cursor_token_scan_ospml() {
    assert_example_matches("examples/tested/basics/cursor/token_scan.ospml");
}

#[test]
fn basics_cursor_utf8_walk_osp() {
    assert_example_matches("examples/tested/basics/cursor/utf8_walk.osp");
}

#[test]
fn basics_cursor_utf8_walk_ospml() {
    assert_example_matches("examples/tested/basics/cursor/utf8_walk.ospml");
}

#[test]
fn basics_errors_error_messages_osp() {
    assert_example_matches("examples/tested/basics/errors/error_messages.osp");
}

#[test]
fn basics_errors_error_messages_ospml() {
    assert_example_matches("examples/tested/basics/errors/error_messages.ospml");
}

#[test]
fn basics_errors_validation_pipeline_osp() {
    assert_example_matches("examples/tested/basics/errors/validation_pipeline.osp");
}

#[test]
fn basics_errors_validation_pipeline_ospml() {
    assert_example_matches("examples/tested/basics/errors/validation_pipeline.ospml");
}

#[test]
fn basics_feature_omnibus_osp() {
    assert_example_matches("examples/tested/basics/feature_omnibus.osp");
}

#[test]
fn basics_feature_omnibus_ospml() {
    assert_example_matches("examples/tested/basics/feature_omnibus.ospml");
}

#[test]
fn basics_field_access_comprehensive_osp() {
    assert_example_matches("examples/tested/basics/field_access_comprehensive.osp");
}

#[test]
fn basics_field_access_comprehensive_ospml() {
    assert_example_matches("examples/tested/basics/field_access_comprehensive.ospml");
}

#[test]
fn basics_files_file_io_json_workflow_osp() {
    assert_example_matches("examples/tested/basics/files/file_io_json_workflow.osp");
}

#[test]
fn basics_files_file_io_json_workflow_ospml() {
    assert_example_matches("examples/tested/basics/files/file_io_json_workflow.ospml");
}

#[test]
fn basics_function_composition_test_osp() {
    assert_example_matches("examples/tested/basics/function_composition_test.osp");
}

#[test]
fn basics_functional_functional_showcase_osp() {
    assert_example_matches("examples/tested/basics/functional/functional_showcase.osp");
}

#[test]
fn basics_functional_functional_showcase_ospml() {
    assert_example_matches("examples/tested/basics/functional/functional_showcase.ospml");
}

#[test]
fn basics_games_adventure_game_osp() {
    assert_example_matches("examples/tested/basics/games/adventure_game.osp");
}

#[test]
fn basics_games_adventure_game_ospml() {
    assert_example_matches("examples/tested/basics/games/adventure_game.ospml");
}

#[test]
fn basics_games_space_trader_osp() {
    assert_example_matches("examples/tested/basics/games/space_trader.osp");
}

#[test]
fn basics_games_space_trader_ospml() {
    assert_example_matches("examples/tested/basics/games/space_trader.ospml");
}

#[test]
fn basics_knownbugs_bug1_spawn_record_osp() {
    assert_example_matches("examples/tested/basics/knownbugs/bug1_spawn_record.osp");
}

#[test]
fn basics_knownbugs_bug1_spawn_record_ospml() {
    assert_example_matches("examples/tested/basics/knownbugs/bug1_spawn_record.ospml");
}

#[test]
fn basics_knownbugs_bug2_string_union_payload_osp() {
    assert_example_matches("examples/tested/basics/knownbugs/bug2_string_union_payload.osp");
}

#[test]
fn basics_knownbugs_bug2_string_union_payload_ospml() {
    assert_example_matches("examples/tested/basics/knownbugs/bug2_string_union_payload.ospml");
}

#[test]
fn basics_knownbugs_bug3_map_built_index_osp() {
    assert_example_matches("examples/tested/basics/knownbugs/bug3_map_built_index.osp");
}

#[test]
fn basics_knownbugs_bug3_map_built_index_ospml() {
    assert_example_matches("examples/tested/basics/knownbugs/bug3_map_built_index.ospml");
}

#[test]
fn basics_knownbugs_bug4_union_return_arg_osp() {
    assert_example_matches("examples/tested/basics/knownbugs/bug4_union_return_arg.osp");
}

#[test]
fn basics_knownbugs_bug4_union_return_arg_ospml() {
    assert_example_matches("examples/tested/basics/knownbugs/bug4_union_return_arg.ospml");
}

#[test]
fn basics_lists_list_basics_osp() {
    assert_example_matches("examples/tested/basics/lists/list_basics.osp");
}

#[test]
fn basics_lists_list_basics_ospml() {
    assert_example_matches("examples/tested/basics/lists/list_basics.ospml");
}

#[test]
fn basics_lists_map_basics_osp() {
    assert_example_matches("examples/tested/basics/lists/map_basics.osp");
}

#[test]
fn basics_lists_map_basics_ospml() {
    assert_example_matches("examples/tested/basics/lists/map_basics.ospml");
}

#[test]
fn basics_math_comprehensive_math_osp() {
    assert_example_matches("examples/tested/basics/math/comprehensive_math.osp");
}

#[test]
fn basics_math_comprehensive_math_ospml() {
    assert_example_matches("examples/tested/basics/math/comprehensive_math.ospml");
}

#[test]
fn basics_operators_boolean_consolidated_osp() {
    assert_example_matches("examples/tested/basics/operators/boolean_consolidated.osp");
}

#[test]
fn basics_operators_boolean_consolidated_ospml() {
    assert_example_matches("examples/tested/basics/operators/boolean_consolidated.ospml");
}

#[test]
fn basics_osprey_mega_showcase_osp() {
    assert_example_matches("examples/tested/basics/osprey_mega_showcase.osp");
}

#[test]
fn basics_osprey_mega_showcase_ospml() {
    assert_example_matches("examples/tested/basics/osprey_mega_showcase.ospml");
}

#[test]
fn basics_pattern_matching_pattern_matching_complete_osp() {
    assert_example_matches("examples/tested/basics/pattern_matching/pattern_matching_complete.osp");
}

#[test]
fn basics_processes_async_process_management_osp() {
    assert_example_matches("examples/tested/basics/processes/async_process_management.osp");
}

#[test]
fn basics_processes_async_process_management_ospml() {
    assert_example_matches("examples/tested/basics/processes/async_process_management.ospml");
}

#[test]
fn basics_processes_callback_stdout_demo_osp() {
    assert_example_matches("examples/tested/basics/processes/callback_stdout_demo.osp");
}

#[test]
fn basics_processes_callback_stdout_demo_ospml() {
    assert_example_matches("examples/tested/basics/processes/callback_stdout_demo.ospml");
}

#[test]
fn basics_strings_string_edge_cases_osp() {
    assert_example_matches("examples/tested/basics/strings/string_edge_cases.osp");
}

#[test]
fn basics_strings_string_edge_cases_ospml() {
    assert_example_matches("examples/tested/basics/strings/string_edge_cases.ospml");
}

#[test]
fn basics_strings_string_pipeline_osp() {
    assert_example_matches("examples/tested/basics/strings/string_pipeline.osp");
}

#[test]
fn basics_strings_string_pipeline_ospml() {
    assert_example_matches("examples/tested/basics/strings/string_pipeline.ospml");
}

#[test]
fn basics_types_any_type_comprehensive_osp() {
    assert_example_matches("examples/tested/basics/types/any_type_comprehensive.osp");
}

#[test]
fn basics_types_any_type_comprehensive_ospml() {
    assert_example_matches("examples/tested/basics/types/any_type_comprehensive.ospml");
}

#[test]
fn basics_types_pure_hindley_milner_test_osp() {
    assert_example_matches("examples/tested/basics/types/pure_hindley_milner_test.osp");
}

#[test]
fn basics_types_pure_hindley_milner_test_ospml() {
    assert_example_matches("examples/tested/basics/types/pure_hindley_milner_test.ospml");
}

#[test]
fn basics_types_record_update_basic_osp() {
    assert_example_matches("examples/tested/basics/types/record_update_basic.osp");
}

#[test]
fn basics_types_record_update_basic_ospml() {
    assert_example_matches("examples/tested/basics/types/record_update_basic.ospml");
}

#[test]
fn basics_types_recursive_unions_osp() {
    assert_example_matches("examples/tested/basics/types/recursive_unions.osp");
}

#[test]
fn basics_types_recursive_unions_ospml() {
    assert_example_matches("examples/tested/basics/types/recursive_unions.ospml");
}

#[test]
fn basics_types_type_equality_comprehensive_osp() {
    assert_example_matches("examples/tested/basics/types/type_equality_comprehensive.osp");
}

#[test]
fn basics_types_type_equality_comprehensive_ospml() {
    assert_example_matches("examples/tested/basics/types/type_equality_comprehensive.ospml");
}

#[test]
fn basics_types_user_defined_unions_osp() {
    assert_example_matches("examples/tested/basics/types/user_defined_unions.osp");
}

#[test]
fn basics_types_user_defined_unions_ospml() {
    assert_example_matches("examples/tested/basics/types/user_defined_unions.ospml");
}

#[test]
fn basics_validation_proper_validation_test_osp() {
    assert_example_matches("examples/tested/basics/validation/proper_validation_test.osp");
}

#[test]
fn basics_validation_proper_validation_test_ospml() {
    assert_example_matches("examples/tested/basics/validation/proper_validation_test.ospml");
}

#[test]
fn basics_website_website_examples_osp() {
    assert_example_matches("examples/tested/basics/website/website_examples.osp");
}

#[test]
fn basics_website_website_examples_ospml() {
    assert_example_matches("examples/tested/basics/website/website_examples.ospml");
}

#[test]
fn db_database_effect_osp() {
    assert_example_matches("examples/tested/db/database_effect.osp");
}

#[test]
fn db_database_effect_ospml() {
    assert_example_matches("examples/tested/db/database_effect.ospml");
}

#[test]
fn db_sqlite_basics_osp() {
    assert_example_matches("examples/tested/db/sqlite_basics.osp");
}

#[test]
fn db_sqlite_basics_ospml() {
    assert_example_matches("examples/tested/db/sqlite_basics.ospml");
}

#[test]
fn effects_algebraic_effects_comprehensive_osp() {
    assert_example_matches("examples/tested/effects/algebraic_effects_comprehensive.osp");
}

#[test]
fn effects_algebraic_effects_comprehensive_ospml() {
    assert_example_matches("examples/tested/effects/algebraic_effects_comprehensive.ospml");
}

#[test]
fn effects_fiber_effects_osp() {
    assert_example_matches("examples/tested/effects/fiber_effects.osp");
}

#[test]
fn effects_fiber_effects_ospml() {
    assert_example_matches("examples/tested/effects/fiber_effects.ospml");
}

#[test]
fn effects_handler_scoping_osp() {
    assert_example_matches("examples/tested/effects/handler_scoping.osp");
}

#[test]
fn effects_handler_scoping_ospml() {
    assert_example_matches("examples/tested/effects/handler_scoping.ospml");
}

#[test]
fn effects_http_state_levels_osp() {
    assert_example_matches("examples/tested/effects/http_state_levels.osp");
}

#[test]
fn effects_http_state_levels_ospml() {
    assert_example_matches("examples/tested/effects/http_state_levels.ospml");
}

#[test]
fn effects_resume_abort_early_exit_osp() {
    assert_example_matches("examples/tested/effects/resume_abort_early_exit.osp");
}

#[test]
fn effects_resume_abort_early_exit_ospml() {
    assert_example_matches("examples/tested/effects/resume_abort_early_exit.ospml");
}

#[test]
fn effects_resume_lifo_audit_osp() {
    assert_example_matches("examples/tested/effects/resume_lifo_audit.osp");
}

#[test]
fn effects_resume_lifo_audit_ospml() {
    assert_example_matches("examples/tested/effects/resume_lifo_audit.ospml");
}

#[test]
fn effects_resume_outer_handler_bridge_osp() {
    assert_example_matches("examples/tested/effects/resume_outer_handler_bridge.osp");
}

#[test]
fn effects_resume_outer_handler_bridge_ospml() {
    assert_example_matches("examples/tested/effects/resume_outer_handler_bridge.ospml");
}

#[test]
fn effects_resume_unit_markers_osp() {
    assert_example_matches("examples/tested/effects/resume_unit_markers.osp");
}

#[test]
fn effects_resume_unit_markers_ospml() {
    assert_example_matches("examples/tested/effects/resume_unit_markers.ospml");
}

#[test]
fn effects_resume_value_rewrite_osp() {
    assert_example_matches("examples/tested/effects/resume_value_rewrite.osp");
}

#[test]
fn effects_resume_value_rewrite_ospml() {
    assert_example_matches("examples/tested/effects/resume_value_rewrite.ospml");
}

#[test]
fn fiber_fiber_showcase_osp() {
    assert_example_matches("examples/tested/fiber/fiber_showcase.osp");
}

#[test]
fn fiber_fiber_showcase_ospml() {
    assert_example_matches("examples/tested/fiber/fiber_showcase.ospml");
}

#[test]
fn http_http_client_example_osp() {
    assert_example_matches("examples/tested/http/http_client_example.osp");
}

#[test]
fn http_http_client_example_ospml() {
    assert_example_matches("examples/tested/http/http_client_example.ospml");
}

#[test]
fn http_http_create_client_osp() {
    assert_example_matches("examples/tested/http/http_create_client.osp");
}

#[test]
fn http_http_create_client_ospml() {
    assert_example_matches("examples/tested/http/http_create_client.ospml");
}

#[test]
fn http_http_response_handle_osp() {
    assert_example_matches("examples/tested/http/http_response_handle.osp");
}

#[test]
fn http_http_response_handle_ospml() {
    assert_example_matches("examples/tested/http/http_response_handle.ospml");
}

#[test]
fn http_http_server_example_osp() {
    assert_example_matches("examples/tested/http/http_server_example.osp");
}

#[test]
fn http_http_server_example_ospml() {
    assert_example_matches("examples/tested/http/http_server_example.ospml");
}

#[test]
fn http_tui_repo_table_osp() {
    assert_example_matches("examples/tested/http/tui_repo_table.osp");
}

#[test]
fn http_tui_repo_table_ospml() {
    assert_example_matches("examples/tested/http/tui_repo_table.ospml");
}

#[test]
fn ml_arith_osp() {
    assert_example_matches("examples/tested/ml/arith.osp");
}

#[test]
fn ml_arith_ospml() {
    assert_example_matches("examples/tested/ml/arith.ospml");
}

#[test]
fn ml_booleans_osp() {
    assert_example_matches("examples/tested/ml/booleans.osp");
}

#[test]
fn ml_booleans_ospml() {
    assert_example_matches("examples/tested/ml/booleans.ospml");
}

#[test]
fn ml_closures_osp() {
    assert_example_matches("examples/tested/ml/closures.osp");
}

#[test]
fn ml_closures_ospml() {
    assert_example_matches("examples/tested/ml/closures.ospml");
}

#[test]
fn ml_curry_partial_osp() {
    assert_example_matches("examples/tested/ml/curry_partial.osp");
}

#[test]
fn ml_curry_partial_ospml() {
    assert_example_matches("examples/tested/ml/curry_partial.ospml");
}

#[test]
fn ml_curry_tour_osp() {
    assert_example_matches("examples/tested/ml/curry_tour.osp");
}

#[test]
fn ml_curry_tour_ospml() {
    assert_example_matches("examples/tested/ml/curry_tour.ospml");
}

#[test]
fn ml_hello_osp() {
    assert_example_matches("examples/tested/ml/hello.osp");
}

#[test]
fn ml_hello_ospml() {
    assert_example_matches("examples/tested/ml/hello.ospml");
}

#[test]
fn ml_hof_osp() {
    assert_example_matches("examples/tested/ml/hof.osp");
}

#[test]
fn ml_hof_ospml() {
    assert_example_matches("examples/tested/ml/hof.ospml");
}

#[test]
fn ml_match_tour_osp() {
    assert_example_matches("examples/tested/ml/match_tour.osp");
}

#[test]
fn ml_match_tour_ospml() {
    assert_example_matches("examples/tested/ml/match_tour.ospml");
}

#[test]
fn ml_matchbool_osp() {
    assert_example_matches("examples/tested/ml/matchbool.osp");
}

#[test]
fn ml_matchbool_ospml() {
    assert_example_matches("examples/tested/ml/matchbool.ospml");
}

#[test]
fn ml_matchint_osp() {
    assert_example_matches("examples/tested/ml/matchint.osp");
}

#[test]
fn ml_matchint_ospml() {
    assert_example_matches("examples/tested/ml/matchint.ospml");
}

#[test]
fn ml_mixed_osp() {
    assert_example_matches("examples/tested/ml/mixed.osp");
}

#[test]
fn ml_mixed_ospml() {
    assert_example_matches("examples/tested/ml/mixed.ospml");
}

#[test]
fn ml_mutation_osp() {
    assert_example_matches("examples/tested/ml/mutation.osp");
}

#[test]
fn ml_mutation_ospml() {
    assert_example_matches("examples/tested/ml/mutation.ospml");
}

#[test]
fn ml_nested_calls_osp() {
    assert_example_matches("examples/tested/ml/nested_calls.osp");
}

#[test]
fn ml_nested_calls_ospml() {
    assert_example_matches("examples/tested/ml/nested_calls.ospml");
}

#[test]
fn ml_pipechain_osp() {
    assert_example_matches("examples/tested/ml/pipechain.osp");
}

#[test]
fn ml_pipechain_ospml() {
    assert_example_matches("examples/tested/ml/pipechain.ospml");
}

#[test]
fn ml_recursion_osp() {
    assert_example_matches("examples/tested/ml/recursion.osp");
}

#[test]
fn ml_recursion_ospml() {
    assert_example_matches("examples/tested/ml/recursion.ospml");
}

#[test]
fn ml_results_state_hof_osp() {
    assert_example_matches("examples/tested/ml/results_state_hof.osp");
}

#[test]
fn ml_results_state_hof_ospml() {
    assert_example_matches("examples/tested/ml/results_state_hof.ospml");
}

#[test]
fn ml_strings_osp() {
    assert_example_matches("examples/tested/ml/strings.osp");
}

#[test]
fn ml_strings_ospml() {
    assert_example_matches("examples/tested/ml/strings.ospml");
}
