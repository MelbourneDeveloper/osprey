//! Native replay [SYMEX-REPLAY] (docs/specs/tooling/SymbolicExecution.md).
//! Each finding's recorded input is fed to the ordinary native build on this
//! host: stdin from the `stdin` object, each file object written under a
//! scratch working directory at the relative path the program opened. A
//! crash or abort is `reproduced`; a clean exit is `klee-only`.

use super::findings::{Class, Finding, Replay};
use super::process::wait_bounded;
use crate::native::{build_input, child_exit_code, scratch_stem, NativeOptions};
use crate::project::CompilationInput;
use osprey_debug::BuildKind;
use std::io::Write;
use std::path::{Path, PathBuf};
use std::process::{Command, Stdio};
use std::time::Duration;

/// A program that never exits (a server, a prompt loop) cannot reproduce a
/// finding; its replay is skipped after this long.
const REPLAY_TIMEOUT: Duration = Duration::from_secs(30);

pub(crate) fn replay(input: &CompilationInput, findings: &mut [Finding]) -> Result<(), String> {
    if !findings.iter().any(replayable) {
        return Ok(());
    }
    let stem = scratch_stem(input.display_path());
    let exe = std::env::temp_dir().join(format!("{stem}-hunt.out"));
    let workdir = std::env::temp_dir().join(format!("{stem}-hunt-replay"));
    build_input(
        input,
        &exe,
        NativeOptions::new("default", BuildKind::Release),
    )
    .map_err(|code| format!("the native build for replay failed (exit {code:?})"))?;
    let outcome = replay_all(&exe, &workdir, findings);
    let _ = std::fs::remove_file(&exe);
    let _ = std::fs::remove_dir_all(&workdir);
    outcome
}

fn replayable(finding: &Finding) -> bool {
    finding.class != Class::Unmodelled
}

fn replay_all(exe: &Path, workdir: &Path, findings: &mut [Finding]) -> Result<(), String> {
    for finding in findings.iter_mut() {
        finding.replay = if replayable(finding) {
            run_one(exe, workdir, finding)?
        } else {
            Replay::Skipped("an unmodelled call stops the path under KLEE only".to_string())
        };
    }
    Ok(())
}

pub(crate) fn run_one(exe: &Path, workdir: &Path, finding: &Finding) -> Result<Replay, String> {
    let _ = std::fs::remove_dir_all(workdir);
    std::fs::create_dir_all(workdir)
        .map_err(|error| format!("cannot create {}: {error}", workdir.display()))?;
    if let Some(reason) = write_files(workdir, &finding.files) {
        return Ok(Replay::Skipped(reason));
    }
    let mut child = Command::new(exe)
        .current_dir(workdir)
        .stdin(Stdio::piped())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn()
        .map_err(|error| format!("could not run {}: {error}", exe.display()))?;
    if let Some(mut stdin) = child.stdin.take() {
        // A program that crashes before reading closes the pipe early; that
        // write error is the finding, not a replay failure.
        let _ = stdin.write_all(&finding.stdin);
    }
    match wait_bounded(&mut child, REPLAY_TIMEOUT)? {
        Some(status) => Ok(outcome(child_exit_code(status))),
        None => Ok(Replay::Skipped(format!(
            "the native program did not exit within {}s",
            REPLAY_TIMEOUT.as_secs()
        ))),
    }
}

/// Writes every file object; a path outside the scratch directory is never
/// written, and names the reason the finding is skipped.
fn write_files(workdir: &Path, files: &[(String, Vec<u8>)]) -> Option<String> {
    for (path, bytes) in files {
        let relative = Path::new(path);
        if relative.is_absolute()
            || relative
                .components()
                .any(|c| c == std::path::Component::ParentDir)
        {
            return Some(format!("file {path} is outside the replay directory"));
        }
        let target: PathBuf = workdir.join(relative);
        if let Some(parent) = target.parent() {
            let _ = std::fs::create_dir_all(parent);
        }
        if let Err(error) = std::fs::write(&target, bytes) {
            return Some(format!("cannot write {path}: {error}"));
        }
    }
    None
}

/// `child_exit_code` folds a signal death into `128 + signal`; only a signal
/// reproduces a finding, since a non-zero `main` is a program's own choice.
pub(crate) fn outcome(code: u8) -> Replay {
    match code.checked_sub(128) {
        Some(signal) if signal > 0 => Replay::Reproduced(signal_name(signal)),
        _ => Replay::KleeOnly,
    }
}

fn signal_name(signal: u8) -> String {
    let bus = if cfg!(target_os = "macos") { 10 } else { 7 };
    match signal {
        4 => "SIGILL".to_string(),
        5 => "SIGTRAP".to_string(),
        6 => "SIGABRT".to_string(),
        8 => "SIGFPE".to_string(),
        11 => "SIGSEGV".to_string(),
        n if n == bus => "SIGBUS".to_string(),
        n => format!("signal {n}"),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    #[cfg(unix)]
    use std::os::unix::fs::PermissionsExt;

    fn finding(files: Vec<(String, Vec<u8>)>) -> Finding {
        Finding {
            test: "test000001".into(),
            class: Class::Memory,
            suffix: "ptr".into(),
            message: "boom".into(),
            location: None,
            runtime_frame: None,
            stdin: b"9\n".to_vec(),
            files,
            replay: Replay::NotRun,
        }
    }

    #[test]
    fn only_a_signal_reproduces() {
        assert_eq!(outcome(0), Replay::KleeOnly);
        assert_eq!(outcome(3), Replay::KleeOnly);
        assert_eq!(outcome(128), Replay::KleeOnly);
        assert_eq!(outcome(128 + 11), Replay::Reproduced("SIGSEGV".into()));
        assert_eq!(outcome(128 + 6), Replay::Reproduced("SIGABRT".into()));
        assert_eq!(outcome(128 + 8), Replay::Reproduced("SIGFPE".into()));
        assert_eq!(signal_name(4), "SIGILL");
        assert_eq!(signal_name(5), "SIGTRAP");
        assert_eq!(
            signal_name(if cfg!(target_os = "macos") { 10 } else { 7 }),
            "SIGBUS"
        );
        assert_eq!(signal_name(30), "signal 30");
    }

    #[cfg(unix)]
    #[test]
    fn stdin_and_files_reach_the_native_program() {
        let workdir =
            std::env::temp_dir().join(format!("osprey-hunt-replay-{}", std::process::id()));
        let script =
            std::env::temp_dir().join(format!("osprey-hunt-replay-{}.sh", std::process::id()));
        std::fs::write(&script, "#!/bin/sh\nread n\n[ \"$n\" = 9 ] && [ \"$(cat data.txt)\" = payload ] && kill -SEGV $$\nexit 0\n").expect("script");
        std::fs::set_permissions(&script, std::fs::Permissions::from_mode(0o755)).expect("chmod");
        let crash = finding(vec![("data.txt".into(), b"payload".to_vec())]);
        assert_eq!(
            run_one(&script, &workdir, &crash).expect("replay"),
            Replay::Reproduced("SIGSEGV".into())
        );
        let clean = finding(vec![("data.txt".into(), b"other".to_vec())]);
        assert_eq!(
            run_one(&script, &workdir, &clean).expect("replay"),
            Replay::KleeOnly
        );
        let escape = finding(vec![("../escape.txt".into(), Vec::new())]);
        assert!(
            matches!(run_one(&script, &workdir, &escape).expect("replay"), Replay::Skipped(reason) if reason.contains("outside"))
        );
        assert!(!workdir.join("../escape.txt").exists());
        let _ = std::fs::remove_dir_all(&workdir);
        let _ = std::fs::remove_file(&script);
    }

    #[test]
    fn a_missing_executable_is_an_error() {
        let workdir =
            std::env::temp_dir().join(format!("osprey-hunt-replay-missing-{}", std::process::id()));
        let error = run_one(
            Path::new("/nonexistent/program"),
            &workdir,
            &finding(Vec::new()),
        )
        .expect_err("missing");
        assert!(error.contains("could not run"));
        let _ = std::fs::remove_dir_all(&workdir);
    }
}
