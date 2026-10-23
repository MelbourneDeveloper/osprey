//! The live-attach fixture both live-profiler suites run [PROF-ATTACH-ENDPOINT]:
//! compiled once with debug info by this build's `osprey`, then started per
//! test with a private `TMPDIR`, so its endpoint lives where only that test
//! looks.

use std::path::{Path, PathBuf};
use std::process::{Child, Command, Stdio};
use std::sync::OnceLock;
use std::time::{Duration, Instant};

/// The fixture's hot function.
pub(crate) const HOT_FUNCTION: &str = "spin";
/// How long a test waits for the fixture's endpoint to appear.
const STARTUP: Duration = Duration::from_secs(20);

/// The fixture binary, compiled once for every test in this file. Under
/// target/, so a run leaves nothing in the system temp directory.
pub(crate) fn fixture_exe() -> &'static Path {
    static EXE: OnceLock<PathBuf> = OnceLock::new();
    EXE.get_or_init(|| {
        let exe = Path::new(env!("CARGO_TARGET_TMPDIR")).join("live_attach");
        let source = Path::new(env!("CARGO_MANIFEST_DIR")).join("tests/fixtures/live_attach.osp");
        let status = Command::new(env!("CARGO_BIN_EXE_osprey"))
            .current_dir(Path::new(env!("CARGO_MANIFEST_DIR")).join("../.."))
            .arg(&source)
            .args(["--compile", "--profile", "--quiet", "-o"])
            .arg(&exe)
            .status()
            .expect("osprey runs");
        assert!(status.success(), "the fixture compiles with debug info");
        exe
    })
}

/// The running fixture with its own `TMPDIR`, so its endpoint lives in a
/// directory only this test reads. Killed, and its directory removed, on drop.
pub(crate) struct Fixture {
    child: Child,
    pub(crate) pid: u32,
    pub(crate) tmp: PathBuf,
    pub(crate) socket_dir: PathBuf,
}

impl Fixture {
    pub(crate) fn start(name: &str) -> Self {
        // A Unix socket path is at most 103 bytes on macOS, and a coverage
        // build's target directory is already most of that; the system temp
        // directory leaves room for `osprey-diagnostics-<uid>/<pid>.sock`.
        let tmp = std::env::temp_dir().join(format!("{name}-{}", std::process::id()));
        drop(std::fs::remove_dir_all(&tmp));
        std::fs::create_dir_all(&tmp).expect("scratch directory");
        let child = Command::new(fixture_exe())
            .env("TMPDIR", &tmp)
            .stdin(Stdio::null())
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .spawn()
            .expect("the fixture starts");
        let pid = child.id();
        let socket_dir = wait_for_endpoint(&tmp, pid);
        Self {
            child,
            pid,
            tmp,
            socket_dir,
        }
    }

    pub(crate) fn kill(&mut self) {
        drop(self.child.kill());
        drop(self.child.wait());
    }
}

impl Drop for Fixture {
    fn drop(&mut self) {
        self.kill();
        drop(std::fs::remove_dir_all(&self.tmp));
    }
}

/// `<tmp>/osprey-diagnostics-<uid>`, once it holds `<pid>.sock`.
fn wait_for_endpoint(tmp: &Path, pid: u32) -> PathBuf {
    let deadline = Instant::now() + STARTUP;
    while Instant::now() < deadline {
        let found = std::fs::read_dir(tmp)
            .into_iter()
            .flatten()
            .flatten()
            .map(|entry| entry.path())
            .find(|dir| dir.join(format!("{pid}.sock")).exists());
        if let Some(dir) = found {
            return dir;
        }
        std::thread::sleep(Duration::from_millis(20));
    }
    panic!("no endpoint for PID {pid} under {}", tmp.display());
}
