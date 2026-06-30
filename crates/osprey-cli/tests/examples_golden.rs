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

macro_rules! tested_examples {
    ($($name:ident => $path:literal,)+) => {
        const REGISTERED_EXAMPLES: &[&str] = &[$($path,)+];

        $(
            #[test]
            fn $name() {
                assert_example_matches($path);
            }
        )+
    };
}

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

tested_examples! {
    basics_blocks_block_statements_basic_osp => "examples/tested/basics/blocks/block_statements_basic.osp",
    basics_blocks_block_statements_basic_ospml => "examples/tested/basics/blocks/block_statements_basic.ospml",
    basics_cursor_codepoint_roundtrip_osp => "examples/tested/basics/cursor/codepoint_roundtrip.osp",
    basics_cursor_codepoint_roundtrip_ospml => "examples/tested/basics/cursor/codepoint_roundtrip.ospml",
    basics_cursor_kv_parser_osp => "examples/tested/basics/cursor/kv_parser.osp",
    basics_cursor_kv_parser_ospml => "examples/tested/basics/cursor/kv_parser.ospml",
    basics_cursor_token_scan_osp => "examples/tested/basics/cursor/token_scan.osp",
    basics_cursor_token_scan_ospml => "examples/tested/basics/cursor/token_scan.ospml",
    basics_cursor_utf8_walk_osp => "examples/tested/basics/cursor/utf8_walk.osp",
    basics_cursor_utf8_walk_ospml => "examples/tested/basics/cursor/utf8_walk.ospml",
    basics_errors_error_messages_osp => "examples/tested/basics/errors/error_messages.osp",
    basics_errors_error_messages_ospml => "examples/tested/basics/errors/error_messages.ospml",
    basics_errors_validation_pipeline_osp => "examples/tested/basics/errors/validation_pipeline.osp",
    basics_errors_validation_pipeline_ospml => "examples/tested/basics/errors/validation_pipeline.ospml",
    basics_feature_omnibus_osp => "examples/tested/basics/feature_omnibus.osp",
    basics_feature_omnibus_ospml => "examples/tested/basics/feature_omnibus.ospml",
    basics_field_access_comprehensive_osp => "examples/tested/basics/field_access_comprehensive.osp",
    basics_field_access_comprehensive_ospml => "examples/tested/basics/field_access_comprehensive.ospml",
    basics_files_file_io_json_workflow_osp => "examples/tested/basics/files/file_io_json_workflow.osp",
    basics_files_file_io_json_workflow_ospml => "examples/tested/basics/files/file_io_json_workflow.ospml",
    basics_function_composition_test_osp => "examples/tested/basics/function_composition_test.osp",
    basics_functional_functional_showcase_osp => "examples/tested/basics/functional/functional_showcase.osp",
    basics_functional_functional_showcase_ospml => "examples/tested/basics/functional/functional_showcase.ospml",
    basics_games_adventure_game_osp => "examples/tested/basics/games/adventure_game.osp",
    basics_games_adventure_game_ospml => "examples/tested/basics/games/adventure_game.ospml",
    basics_games_space_trader_osp => "examples/tested/basics/games/space_trader.osp",
    basics_games_space_trader_ospml => "examples/tested/basics/games/space_trader.ospml",
    basics_knownbugs_bug1_spawn_record_osp => "examples/tested/basics/knownbugs/bug1_spawn_record.osp",
    basics_knownbugs_bug1_spawn_record_ospml => "examples/tested/basics/knownbugs/bug1_spawn_record.ospml",
    basics_knownbugs_bug2_string_union_payload_osp => "examples/tested/basics/knownbugs/bug2_string_union_payload.osp",
    basics_knownbugs_bug2_string_union_payload_ospml => "examples/tested/basics/knownbugs/bug2_string_union_payload.ospml",
    basics_knownbugs_bug3_map_built_index_osp => "examples/tested/basics/knownbugs/bug3_map_built_index.osp",
    basics_knownbugs_bug3_map_built_index_ospml => "examples/tested/basics/knownbugs/bug3_map_built_index.ospml",
    basics_knownbugs_bug4_union_return_arg_osp => "examples/tested/basics/knownbugs/bug4_union_return_arg.osp",
    basics_knownbugs_bug4_union_return_arg_ospml => "examples/tested/basics/knownbugs/bug4_union_return_arg.ospml",
    basics_lists_list_basics_osp => "examples/tested/basics/lists/list_basics.osp",
    basics_lists_list_basics_ospml => "examples/tested/basics/lists/list_basics.ospml",
    basics_lists_map_basics_osp => "examples/tested/basics/lists/map_basics.osp",
    basics_lists_map_basics_ospml => "examples/tested/basics/lists/map_basics.ospml",
    basics_math_comprehensive_math_osp => "examples/tested/basics/math/comprehensive_math.osp",
    basics_math_comprehensive_math_ospml => "examples/tested/basics/math/comprehensive_math.ospml",
    basics_operators_boolean_consolidated_osp => "examples/tested/basics/operators/boolean_consolidated.osp",
    basics_operators_boolean_consolidated_ospml => "examples/tested/basics/operators/boolean_consolidated.ospml",
    basics_osprey_mega_showcase_osp => "examples/tested/basics/osprey_mega_showcase.osp",
    basics_osprey_mega_showcase_ospml => "examples/tested/basics/osprey_mega_showcase.ospml",
    basics_pattern_matching_pattern_matching_complete_osp => "examples/tested/basics/pattern_matching/pattern_matching_complete.osp",
    basics_processes_async_process_management_osp => "examples/tested/basics/processes/async_process_management.osp",
    basics_processes_async_process_management_ospml => "examples/tested/basics/processes/async_process_management.ospml",
    basics_processes_callback_stdout_demo_osp => "examples/tested/basics/processes/callback_stdout_demo.osp",
    basics_processes_callback_stdout_demo_ospml => "examples/tested/basics/processes/callback_stdout_demo.ospml",
    basics_strings_string_edge_cases_osp => "examples/tested/basics/strings/string_edge_cases.osp",
    basics_strings_string_edge_cases_ospml => "examples/tested/basics/strings/string_edge_cases.ospml",
    basics_strings_string_pipeline_osp => "examples/tested/basics/strings/string_pipeline.osp",
    basics_strings_string_pipeline_ospml => "examples/tested/basics/strings/string_pipeline.ospml",
    basics_types_any_type_comprehensive_osp => "examples/tested/basics/types/any_type_comprehensive.osp",
    basics_types_any_type_comprehensive_ospml => "examples/tested/basics/types/any_type_comprehensive.ospml",
    basics_types_pure_hindley_milner_test_osp => "examples/tested/basics/types/pure_hindley_milner_test.osp",
    basics_types_pure_hindley_milner_test_ospml => "examples/tested/basics/types/pure_hindley_milner_test.ospml",
    basics_types_record_update_basic_osp => "examples/tested/basics/types/record_update_basic.osp",
    basics_types_record_update_basic_ospml => "examples/tested/basics/types/record_update_basic.ospml",
    basics_types_recursive_unions_osp => "examples/tested/basics/types/recursive_unions.osp",
    basics_types_recursive_unions_ospml => "examples/tested/basics/types/recursive_unions.ospml",
    basics_types_type_equality_comprehensive_osp => "examples/tested/basics/types/type_equality_comprehensive.osp",
    basics_types_type_equality_comprehensive_ospml => "examples/tested/basics/types/type_equality_comprehensive.ospml",
    basics_types_user_defined_unions_osp => "examples/tested/basics/types/user_defined_unions.osp",
    basics_types_user_defined_unions_ospml => "examples/tested/basics/types/user_defined_unions.ospml",
    basics_validation_proper_validation_test_osp => "examples/tested/basics/validation/proper_validation_test.osp",
    basics_validation_proper_validation_test_ospml => "examples/tested/basics/validation/proper_validation_test.ospml",
    basics_website_website_examples_osp => "examples/tested/basics/website/website_examples.osp",
    basics_website_website_examples_ospml => "examples/tested/basics/website/website_examples.ospml",
    db_database_effect_osp => "examples/tested/db/database_effect.osp",
    db_database_effect_ospml => "examples/tested/db/database_effect.ospml",
    db_sqlite_basics_osp => "examples/tested/db/sqlite_basics.osp",
    db_sqlite_basics_ospml => "examples/tested/db/sqlite_basics.ospml",
    effects_algebraic_effects_comprehensive_osp => "examples/tested/effects/algebraic_effects_comprehensive.osp",
    effects_algebraic_effects_comprehensive_ospml => "examples/tested/effects/algebraic_effects_comprehensive.ospml",
    effects_fiber_effects_osp => "examples/tested/effects/fiber_effects.osp",
    effects_fiber_effects_ospml => "examples/tested/effects/fiber_effects.ospml",
    effects_handler_scoping_osp => "examples/tested/effects/handler_scoping.osp",
    effects_handler_scoping_ospml => "examples/tested/effects/handler_scoping.ospml",
    effects_http_state_levels_osp => "examples/tested/effects/http_state_levels.osp",
    effects_http_state_levels_ospml => "examples/tested/effects/http_state_levels.ospml",
    effects_resume_abort_early_exit_osp => "examples/tested/effects/resume_abort_early_exit.osp",
    effects_resume_abort_early_exit_ospml => "examples/tested/effects/resume_abort_early_exit.ospml",
    effects_resume_lifo_audit_osp => "examples/tested/effects/resume_lifo_audit.osp",
    effects_resume_lifo_audit_ospml => "examples/tested/effects/resume_lifo_audit.ospml",
    effects_resume_outer_handler_bridge_osp => "examples/tested/effects/resume_outer_handler_bridge.osp",
    effects_resume_outer_handler_bridge_ospml => "examples/tested/effects/resume_outer_handler_bridge.ospml",
    effects_resume_unit_markers_osp => "examples/tested/effects/resume_unit_markers.osp",
    effects_resume_unit_markers_ospml => "examples/tested/effects/resume_unit_markers.ospml",
    effects_resume_value_rewrite_osp => "examples/tested/effects/resume_value_rewrite.osp",
    effects_resume_value_rewrite_ospml => "examples/tested/effects/resume_value_rewrite.ospml",
    fiber_fiber_showcase_osp => "examples/tested/fiber/fiber_showcase.osp",
    fiber_fiber_showcase_ospml => "examples/tested/fiber/fiber_showcase.ospml",
    http_http_client_example_osp => "examples/tested/http/http_client_example.osp",
    http_http_client_example_ospml => "examples/tested/http/http_client_example.ospml",
    http_http_create_client_osp => "examples/tested/http/http_create_client.osp",
    http_http_create_client_ospml => "examples/tested/http/http_create_client.ospml",
    http_http_response_handle_osp => "examples/tested/http/http_response_handle.osp",
    http_http_response_handle_ospml => "examples/tested/http/http_response_handle.ospml",
    http_http_server_example_osp => "examples/tested/http/http_server_example.osp",
    http_http_server_example_ospml => "examples/tested/http/http_server_example.ospml",
    http_tui_repo_table_osp => "examples/tested/http/tui_repo_table.osp",
    http_tui_repo_table_ospml => "examples/tested/http/tui_repo_table.ospml",
    ml_arith_osp => "examples/tested/ml/arith.osp",
    ml_arith_ospml => "examples/tested/ml/arith.ospml",
    ml_booleans_osp => "examples/tested/ml/booleans.osp",
    ml_booleans_ospml => "examples/tested/ml/booleans.ospml",
    ml_closures_osp => "examples/tested/ml/closures.osp",
    ml_closures_ospml => "examples/tested/ml/closures.ospml",
    ml_curry_partial_osp => "examples/tested/ml/curry_partial.osp",
    ml_curry_partial_ospml => "examples/tested/ml/curry_partial.ospml",
    ml_curry_tour_osp => "examples/tested/ml/curry_tour.osp",
    ml_curry_tour_ospml => "examples/tested/ml/curry_tour.ospml",
    ml_hello_osp => "examples/tested/ml/hello.osp",
    ml_hello_ospml => "examples/tested/ml/hello.ospml",
    ml_hof_osp => "examples/tested/ml/hof.osp",
    ml_hof_ospml => "examples/tested/ml/hof.ospml",
    ml_match_tour_osp => "examples/tested/ml/match_tour.osp",
    ml_match_tour_ospml => "examples/tested/ml/match_tour.ospml",
    ml_matchbool_osp => "examples/tested/ml/matchbool.osp",
    ml_matchbool_ospml => "examples/tested/ml/matchbool.ospml",
    ml_matchint_osp => "examples/tested/ml/matchint.osp",
    ml_matchint_ospml => "examples/tested/ml/matchint.ospml",
    ml_mixed_osp => "examples/tested/ml/mixed.osp",
    ml_mixed_ospml => "examples/tested/ml/mixed.ospml",
    ml_mutation_osp => "examples/tested/ml/mutation.osp",
    ml_mutation_ospml => "examples/tested/ml/mutation.ospml",
    ml_nested_calls_osp => "examples/tested/ml/nested_calls.osp",
    ml_nested_calls_ospml => "examples/tested/ml/nested_calls.ospml",
    ml_pipechain_osp => "examples/tested/ml/pipechain.osp",
    ml_pipechain_ospml => "examples/tested/ml/pipechain.ospml",
    ml_recursion_osp => "examples/tested/ml/recursion.osp",
    ml_recursion_ospml => "examples/tested/ml/recursion.ospml",
    ml_results_state_hof_osp => "examples/tested/ml/results_state_hof.osp",
    ml_results_state_hof_ospml => "examples/tested/ml/results_state_hof.ospml",
    ml_strings_osp => "examples/tested/ml/strings.osp",
    ml_strings_ospml => "examples/tested/ml/strings.ospml",
}
