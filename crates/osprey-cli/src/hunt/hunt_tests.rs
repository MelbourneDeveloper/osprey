//! End-to-end tests of `osprey hunt` against a stand-in `klee` executable
//! [SYMEX-CLI] [SYMEX-FINDINGS] [SYMEX-REPLAY]. The stand-in copies a
//! prepared KLEE output directory into `--output-dir`, so the whole pipeline
//! (compile, link, run, scan, replay, report, exit code) runs without the
//! engine, and the real engine is covered by `tests/hunt/` under `make`.

#![expect(
    clippy::indexing_slicing,
    reason = "test assertions: a missing JSON key is a test failure, not a production panic"
)]

use super::toolchain::{KLEE_ENV, RUNTIME_ENV};
use super::*;
use std::sync::{Mutex, MutexGuard};

/// Environment variables are process-wide; tests that set them take turns.
static ENV_LOCK: Mutex<()> = Mutex::new(());

fn lock() -> MutexGuard<'static, ()> {
    ENV_LOCK
        .lock()
        .unwrap_or_else(std::sync::PoisonError::into_inner)
}

const PTR_ERR: &str = "Error: memory error: out of bound pointer\n\
File: list_runtime.c\nLine: 88\nassembly.ll line: 4012\nState: 7\nStack: \n\
\t#000000031 in osp_list_get(ptr, i64) at list_runtime.c:88\n\
\t#000000024 in main() at PROGRAM:1\n";

/// A scratch directory holding the program, the fake engine and its canned
/// output; removed on drop.
struct Scene {
    root: PathBuf,
}

impl Scene {
    fn new(name: &str, program: &str) -> Self {
        let root =
            std::env::temp_dir().join(format!("osprey-hunt-e2e-{name}-{}", std::process::id()));
        let _ = std::fs::remove_dir_all(&root);
        std::fs::create_dir_all(root.join("canned")).expect("scene dir");
        std::fs::write(root.join("main.osp"), program).expect("program");
        std::fs::write(root.join("runtime.bc"), "BC").expect("runtime");
        Self { root }
    }

    fn program(&self) -> String {
        self.root.join("main.osp").display().to_string()
    }

    fn canned(&self, name: &str, bytes: &[u8]) {
        std::fs::write(self.root.join("canned").join(name), bytes).expect("canned file");
    }

    fn canned_finding(&self) {
        let err = PTR_ERR.replace("PROGRAM", &self.program());
        self.canned("test000001.ptr.err", err.as_bytes());
        let ktest = ktest::encode(&[("stdin", b"9\nzz"), ("stdin.len", &2i32.to_le_bytes())]);
        self.canned("test000001.ktest", &ktest);
        self.canned("test000002.early", b"max-time\n");
        self.canned("info", b"Elapsed: 00:00:02\nKLEE: done: explored paths = 4\nKLEE: done: completed paths = 3\nKLEE: done: total instructions = 99\n");
    }

    /// A `klee` that copies the canned directory to `--output-dir` and
    /// exits with `code`.
    fn fake_klee(&self, code: u8) -> PathBuf {
        let script = self.root.join("klee");
        let body = format!(
            "#!/bin/sh\nfor a in \"$@\"; do case \"$a\" in --output-dir=*) out=\"${{a#--output-dir=}}\";; esac; done\n\
             mkdir -p \"$out\" && cp {}/canned/* \"$out\"/ 2>/dev/null\necho engine-said-this >&2\nexit {code}\n",
            self.root.display()
        );
        std::fs::write(&script, body).expect("fake klee");
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            std::fs::set_permissions(&script, std::fs::Permissions::from_mode(0o755))
                .expect("chmod");
        }
        script
    }

    fn with_toolchain<T>(&self, klee: &Path, run: impl FnOnce() -> T) -> T {
        let _guard = lock();
        std::env::set_var(KLEE_ENV, klee);
        std::env::set_var(RUNTIME_ENV, self.root.join("runtime.bc"));
        let result = run();
        std::env::remove_var(KLEE_ENV);
        std::env::remove_var(RUNTIME_ENV);
        result
    }
}

impl Drop for Scene {
    fn drop(&mut self) {
        let _ = std::fs::remove_dir_all(&self.root);
    }
}

fn options(scene: &Scene, extra: &[&str]) -> Options {
    let mut args = vec![scene.program()];
    args.extend(extra.iter().map(ToString::to_string));
    options::parse(&args).expect("options")
}

#[cfg(unix)]
#[test]
fn a_finding_without_replay_fails_and_is_reported_as_json() {
    let scene = Scene::new("json", "print(\"hi\")\n");
    scene.canned_finding();
    let klee = scene.fake_klee(0);
    let out = scene.root.join("kout");
    let opts = options(
        &scene,
        &["--no-replay", "--json", "--out", &out.display().to_string()],
    );
    let (rendered, code) = scene.with_toolchain(&klee, || hunt(&opts)).expect("hunt");
    assert_eq!(code, ExitCode::FAILURE);
    let value: serde_json::Value = serde_json::from_str(&rendered).expect("json");
    assert_eq!(value["engine"], "KLEE");
    assert_eq!(value["explored"]["paths"], 4);
    assert_eq!(value["findings"][0]["class"], "memory");
    assert_eq!(value["findings"][0]["line"], 1);
    assert_eq!(value["findings"][0]["stdin"], "9\\n");
    assert_eq!(value["findings"][0]["replay"], "not-run");
    assert_eq!(value["incomplete"][0]["reason"], "max-time");
    assert!(out.join("info").exists(), "--out keeps KLEE's directory");
}

#[cfg(unix)]
#[test]
fn a_finding_the_native_build_does_not_reproduce_is_klee_only() {
    let scene = Scene::new("replay", "print(\"hi\")\n");
    scene.canned_finding();
    let klee = scene.fake_klee(0);
    let opts = options(&scene, &["--quiet"]);
    let (rendered, code) = scene.with_toolchain(&klee, || hunt(&opts)).expect("hunt");
    assert_eq!(code, ExitCode::SUCCESS, "{rendered}");
    assert!(
        rendered.contains("symbolic execution by KLEE"),
        "{rendered}"
    );
    assert!(rendered.contains("replay   klee-only"), "{rendered}");
    assert!(
        rendered.contains("explored 4 paths (3 complete), 99 instructions, 00:00:02"),
        "{rendered}"
    );
}

#[cfg(unix)]
#[test]
fn a_clean_run_reports_no_findings() {
    let scene = Scene::new("clean", "print(\"hi\")\n");
    scene.canned("info", b"KLEE: done: explored paths = 1\n");
    let klee = scene.fake_klee(0);
    let opts = options(&scene, &["--quiet"]);
    let (rendered, code) = scene.with_toolchain(&klee, || hunt(&opts)).expect("hunt");
    assert_eq!(code, ExitCode::SUCCESS);
    assert!(rendered.contains("findings: none"), "{rendered}");
}

#[cfg(unix)]
#[test]
fn a_failing_engine_is_a_tool_error() {
    let scene = Scene::new("engine", "print(\"hi\")\n");
    let klee = scene.fake_klee(3);
    let opts = options(&scene, &["--quiet"]);
    let code = scene
        .with_toolchain(&klee, || hunt(&opts))
        .expect_err("engine failure");
    assert_eq!(code, ExitCode::from(2));
}

#[test]
fn a_missing_engine_is_a_tool_error() {
    let scene = Scene::new("missing", "print(\"hi\")\n");
    let opts = options(&scene, &[]);
    let missing = scene.root.join("no-klee-here");
    let code = scene
        .with_toolchain(&missing, || hunt(&opts))
        .expect_err("missing klee");
    assert_eq!(code, ExitCode::from(2));
}

#[cfg(unix)]
#[test]
fn a_program_that_does_not_type_check_fails_before_the_engine_runs() {
    let scene = Scene::new("typeerr", "let x: int = \"text\"\n");
    let klee = scene.fake_klee(0);
    let opts = options(&scene, &["--quiet"]);
    let code = scene
        .with_toolchain(&klee, || hunt(&opts))
        .expect_err("type error");
    assert_eq!(code, ExitCode::FAILURE);
}

#[test]
fn run_rejects_bad_arguments_with_usage() {
    assert_eq!(run(&[]), ExitCode::from(2));
    assert_eq!(run(&["--bogus".to_string()]), ExitCode::from(2));
    assert_eq!(run(&["--help".to_string()]), ExitCode::SUCCESS);
    assert_eq!(run(&["-h".to_string()]), ExitCode::SUCCESS);
}

#[test]
fn only_a_klee_output_directory_is_cleared() {
    let dir = std::env::temp_dir().join(format!("osprey-hunt-clear-{}", std::process::id()));
    let _ = std::fs::remove_dir_all(&dir);
    assert!(clear_run_directory(&dir).is_ok(), "absent is fine");
    std::fs::create_dir_all(&dir).expect("dir");
    assert!(clear_run_directory(&dir).is_ok(), "empty is fine");
    std::fs::create_dir_all(&dir).expect("dir");
    std::fs::write(dir.join("precious.txt"), "x").expect("file");
    let error = clear_run_directory(&dir).expect_err("unrelated directory");
    assert!(error.contains("refusing to clear"));
    std::fs::write(dir.join("info"), "KLEE").expect("info");
    assert!(clear_run_directory(&dir).is_ok() && !dir.exists());
}

/// A program the native build accepts lowers for the engine too: the
/// symbolic switch changes only the unhandled-effect guard [SYMEX-CODEGEN].
#[test]
fn a_fiber_program_lowers_symbolically() {
    let path = concat!(
        env!("CARGO_MANIFEST_DIR"),
        "/../../tests/core/collections/nested_generic_collections_fibers.test.osp"
    )
    .to_string();
    let cli = parse_args(std::slice::from_ref(&path)).expect("args");
    let input = load_input(&cli).expect("load");
    let symbolic = module::symbolic_module(&input, 64).expect("symbolic build");
    assert!(symbolic.contains("target triple = \"x86_64-unknown-linux-gnu\""));
    assert!(symbolic.contains("@osp_hunt_input_bytes = global i64 64"));
}
