//! End-to-end CLI tests that drive the real `osprey` binary.
//!
//! `main`, `run_lsp`, and the full compile -> link -> run pipeline are process
//! entry points the in-process twins (the `src` unit tests and
//! `tests/examples_compile.rs`) can never reach. Spawning the built binary does
//! reach them, and because `cargo llvm-cov` instruments `CARGO_BIN_EXE_osprey`
//! too, each child's coverage is merged back into the report — so these tests
//! count toward the per-crate gate. [TEST-RULES][COVERAGE-THRESHOLDS-JSON]

use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};

/// Repo root: `crates/osprey-cli` -> `../..`. The C runtime archives `osprey`
/// links at `--run`/`--compile` time live under `compiler/bin/` there, and
/// `find_runtime_lib` resolves them relative to the process cwd — so every
/// child runs from here.
fn repo_root() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("..").join("..")
}

/// A `Command` for the built `osprey`, rooted at the repo so the runtime
/// archives resolve.
fn osprey() -> Command {
    let mut cmd = Command::new(env!("CARGO_BIN_EXE_osprey"));
    let _ = cmd.current_dir(repo_root());
    cmd
}

/// Write `body` to a uniquely-named temp `.osp` (the name doubles as the file
/// stem, which `--compile` turns into the output executable's name).
fn temp_osp(name: &str, body: &str) -> PathBuf {
    let path = std::env::temp_dir().join(format!("osprey_cli_e2e_{name}.osp"));
    // A failed write surfaces as a downstream "cannot read"/parse failure that
    // trips the test's own assertions — so no panic is needed here.
    let _ = std::fs::write(&path, body);
    path
}

const HELLO: &str = "let g = \"hi\"\nprint(\"v=${g}\")\n";

/// The captured result of one invocation.
struct Out {
    code: Option<i32>,
    stdout: String,
    stderr: String,
}

fn finish(mut cmd: Command) -> Out {
    // A spawn failure becomes an empty, code-less `Out`; every caller asserts on
    // an expected code/stdout, so the failure reports loudly through them.
    match cmd.output() {
        Ok(out) => Out {
            code: out.status.code(),
            stdout: String::from_utf8_lossy(&out.stdout).into_owned(),
            stderr: String::from_utf8_lossy(&out.stderr).into_owned(),
        },
        Err(e) => Out {
            code: None,
            stdout: String::new(),
            stderr: format!("spawn failed: {e}"),
        },
    }
}

/// Run with literal args (no source file): `--version`, `--hover`, etc.
fn run_args(args: &[&str]) -> Out {
    let mut cmd = osprey();
    let _ = cmd.args(args);
    finish(cmd)
}

/// Run against a source `path` plus extra flags — the common compiling shape.
fn run_file(path: &Path, extra: &[&str]) -> Out {
    let mut cmd = osprey();
    let _ = cmd.arg(path).args(extra);
    finish(cmd)
}

/// Run a single `mode` against `path` with `OSPREY_CC` overridden — used to
/// drive the C-driver failure branches of `build_executable`.
fn run_file_cc(path: &Path, mode: &str, cc: &str) -> Out {
    let mut cmd = osprey();
    let _ = cmd.arg(path).arg(mode).env("OSPREY_CC", cc);
    finish(cmd)
}

/// Type-clean but codegen-rejected: a generic function used as a first-class
/// value. It passes the type gate, so every compiling mode reaches codegen and
/// fails there — exercising the `Err` arms `compile_program` feeds.
const GENERIC_AS_VALUE: &str = "fn id(x) = x\nlet f = id\nprint(\"ok\")\n";

#[test]
fn version_plain_and_json() {
    let plain = run_args(&["--version"]);
    assert_eq!(plain.code, Some(0));
    assert!(plain.stdout.contains("osprey"), "{}", plain.stdout);
    let json = run_args(&["--version", "--json"]);
    assert_eq!(json.code, Some(0));
    assert!(json.stdout.contains("\"kind\":\"cli\""), "{}", json.stdout);
}

#[test]
fn lsp_exits_cleanly_on_closed_stdin() {
    let status = osprey()
        .arg("lsp")
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
        .expect("spawn lsp");
    assert!(status.success(), "lsp should exit 0 at EOF");
}

#[test]
fn hover_prints_known_builtin_and_is_silent_for_unknown() {
    let known = run_args(&["--hover", "print"]);
    assert_eq!(known.code, Some(0));
    assert!(known.stdout.contains("print"), "{}", known.stdout);
    let unknown = run_args(&["--hover", "__definitely_not_a_builtin__"]);
    assert_eq!(unknown.code, Some(0));
    assert!(unknown.stdout.trim().is_empty(), "{}", unknown.stdout);
}

#[test]
fn unknown_flag_exits_two_with_usage() {
    let prog = temp_osp("flag", HELLO);
    let o = run_file(&prog, &["--bogus"]);
    assert_eq!(o.code, Some(2));
    assert!(o.stderr.contains("unknown flag --bogus"), "{}", o.stderr);
}

#[test]
fn check_ok_reports_statement_count() {
    let prog = temp_osp("ok", HELLO);
    let o = run_file(&prog, &[]); // default mode is --check
    assert_eq!(o.code, Some(0), "stderr={}", o.stderr);
    assert!(o.stdout.contains("ok"), "{}", o.stdout);
}

#[test]
fn check_missing_file_exits_two() {
    let o = run_args(&["/no/such/osprey_e2e_missing.osp"]);
    assert_eq!(o.code, Some(2));
    assert!(o.stderr.contains("cannot read"), "{}", o.stderr);
}

#[test]
fn check_parse_error_is_reported() {
    let prog = temp_osp("parse", "fn = = =\n");
    let o = run_file(&prog, &["--check"]);
    assert_ne!(o.code, Some(0));
    assert!(!o.stderr.is_empty());
}

#[test]
fn check_type_error_is_reported() {
    let prog = temp_osp("typed", "let y = 1 + \"oops\" - true\n");
    let o = run_file(&prog, &["--check"]);
    assert_ne!(o.code, Some(0));
    assert!(!o.stderr.is_empty(), "{}", o.stderr);
}

#[test]
fn llvm_emits_ir_and_rejects_ill_typed() {
    let ok = temp_osp("llok", HELLO);
    let good = run_file(&ok, &["--llvm"]);
    assert_eq!(good.code, Some(0), "stderr={}", good.stderr);
    assert!(good.stdout.contains("define"), "{}", good.stdout);
    let bad = temp_osp("llbad", "let y = 1 + true\n");
    let rejected = run_file(&bad, &["--llvm"]);
    assert_ne!(rejected.code, Some(0));
}

#[test]
fn ast_and_symbols_modes() {
    let prog = temp_osp("astsym", HELLO);
    let ast = run_file(&prog, &["--ast"]);
    assert_eq!(ast.code, Some(0));
    assert!(!ast.stdout.is_empty());
    let sym = run_file(&prog, &["--symbols"]);
    assert_eq!(sym.code, Some(0));
    assert!(sym.stdout.contains("\"name\""), "{}", sym.stdout);
}

#[test]
fn run_compiles_links_and_executes() {
    let prog = temp_osp("run", HELLO);
    let o = run_file(&prog, &["--run"]);
    assert_eq!(o.code, Some(0), "stderr={}", o.stderr);
    assert!(o.stdout.contains("v=hi"), "{}", o.stdout);
}

#[test]
fn compile_writes_executable_to_cwd() {
    let prog = temp_osp("compile", HELLO);
    // `compile_program_to_disk` names the output after the source stem, in cwd.
    let artifact = repo_root().join("osprey_cli_e2e_compile");
    let _ = std::fs::remove_file(&artifact);
    let o = run_file(&prog, &["--compile"]);
    let produced = artifact.exists();
    let _ = std::fs::remove_file(&artifact);
    assert_eq!(o.code, Some(0), "stderr={}", o.stderr);
    assert!(produced, "expected executable at {}", artifact.display());
    assert!(o.stdout.contains("osprey_cli_e2e_compile"), "{}", o.stdout);
}

#[test]
fn sandbox_blocks_filesystem_capability() {
    let prog = temp_osp("fs", "let c = readFile(\"x.txt\")\n");
    let o = run_file(&prog, &["--llvm", "--no-fs"]);
    assert_ne!(o.code, Some(0), "stdout={}", o.stdout);
    assert!(!o.stderr.is_empty());
}

#[test]
fn quiet_suppresses_the_ok_line() {
    let prog = temp_osp("quiet", HELLO);
    let o = run_file(&prog, &["--check", "--quiet"]);
    assert_eq!(o.code, Some(0), "stderr={}", o.stderr);
    assert!(o.stdout.trim().is_empty(), "{}", o.stdout);
}

#[test]
fn llvm_reports_a_codegen_error() {
    let prog = temp_osp("cgllvm", GENERIC_AS_VALUE);
    let o = run_file(&prog, &["--llvm"]);
    assert_ne!(o.code, Some(0));
    assert!(o.stderr.contains("codegen"), "{}", o.stderr);
}

#[test]
fn run_reports_a_codegen_error() {
    let prog = temp_osp("cgrun", GENERIC_AS_VALUE);
    let o = run_file(&prog, &["--run"]);
    assert_ne!(o.code, Some(0));
    assert!(o.stderr.contains("codegen"), "{}", o.stderr);
}

#[test]
fn compile_reports_a_failing_c_compiler() {
    // `false` runs and exits non-zero -> the "cc failed to compile" branch.
    let prog = temp_osp("ccfail", HELLO);
    let o = run_file_cc(&prog, "--compile", "false");
    let _ = std::fs::remove_file(repo_root().join("osprey_cli_e2e_ccfail"));
    assert_ne!(o.code, Some(0));
    assert!(!o.stderr.is_empty(), "{}", o.stderr);
}

#[test]
fn run_reports_an_uninvokable_c_compiler() {
    // A missing driver can't be spawned at all -> the "could not invoke" branch.
    let prog = temp_osp("ccmiss", HELLO);
    let o = run_file_cc(&prog, "--run", "osprey_no_such_cc_zzz");
    assert_ne!(o.code, Some(0));
    assert!(!o.stderr.is_empty(), "{}", o.stderr);
}
