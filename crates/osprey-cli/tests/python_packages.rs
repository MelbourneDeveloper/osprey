//! The ported Python packages ([PORT-SUITE] in
//! docs/specs/ecosystem/0043-PythonFlavorPorts.md). Every `packages/<name>/tests`
//! project is a Python-flavor program whose assertions compare the port with
//! values the real `CPython` package computed. It must print exactly its
//! committed `expected_output.txt` — one `ok` line per assertion and a
//! `N passed, 0 failed` total — under the default, GC and ARC memory backends,
//! and an ARC run must exit with no live objects.

use std::path::{Path, PathBuf};
use std::process::Command;

fn repo() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("../..")
}

/// Every package directory that holds a `tests/osprey.toml`, sorted.
fn suites() -> Vec<PathBuf> {
    let mut found: Vec<PathBuf> = std::fs::read_dir(repo().join("packages"))
        .into_iter()
        .flatten()
        .flatten()
        .map(|entry| entry.path().join("tests"))
        .filter(|tests| tests.join("osprey.toml").is_file())
        .collect();
    found.sort();
    found
}

/// `(succeeded, stdout, stderr)` of the suite run under `memory`.
fn run(suite: &Path, memory: &str) -> (bool, String, String) {
    let output = Command::new(env!("CARGO_BIN_EXE_osprey"))
        .current_dir(repo())
        .arg(suite)
        .args(["--run", "--quiet", &format!("--memory={memory}")])
        .env("OSPREY_ARC_DEBUG", "1")
        .output();
    output.map_or_else(
        |error| (false, String::new(), error.to_string()),
        |out| {
            (
                out.status.success(),
                String::from_utf8_lossy(&out.stdout).into_owned(),
                String::from_utf8_lossy(&out.stderr).into_owned(),
            )
        },
    )
}

/// The live-object count of the single ARC exit sentinel, `None` when the
/// sentinel is absent or repeated: an unverified run never reads as leak-free.
fn arc_live(stderr: &str) -> Option<String> {
    let counts: Vec<&str> = stderr
        .lines()
        .filter_map(|line| line.strip_prefix("[osp-arc] exit: "))
        .filter_map(|rest| rest.split(' ').next())
        .collect();
    match counts.as_slice() {
        [count] => Some((*count).to_string()),
        _ => None,
    }
}

/// The failures of one suite under one backend, as readable lines.
fn failures(suite: &Path, memory: &str) -> Vec<String> {
    let name = suite.display();
    let golden = std::fs::read_to_string(suite.join("expected_output.txt")).unwrap_or_default();
    let (ok, stdout, stderr) = run(suite, memory);
    let mut problems = Vec::new();
    if !ok {
        problems.push(format!(
            "{name} [{memory}] exited unsuccessfully:\n{stderr}"
        ));
    }
    if golden.is_empty() || stdout != golden {
        problems.push(format!(
            "{name} [{memory}] output differs from expected_output.txt:\n{stdout}"
        ));
    }
    if memory == "arc" && arc_live(&stderr).as_deref() != Some("0") {
        problems.push(format!("{name} [arc] did not exit leak-free:\n{stderr}"));
    }
    problems
}

#[test]
fn every_ported_package_matches_its_cpython_oracle_under_every_memory_backend() {
    let suites = suites();
    assert!(!suites.is_empty(), "no packages/*/tests suites were found");
    let problems: Vec<String> = suites
        .iter()
        .flat_map(|suite| {
            ["default", "gc", "arc"]
                .into_iter()
                .flat_map(move |memory| failures(suite, memory))
        })
        .collect();
    assert!(problems.is_empty(), "{}", problems.join("\n"));
}

/// The string items of the one-line manifest array `key = [...]`.
fn manifest_list(manifest: &str, key: &str) -> Vec<String> {
    manifest
        .lines()
        .find_map(|line| line.strip_prefix(&format!("{key} = [")))
        .and_then(|rest| rest.strip_suffix(']'))
        .map(|items| {
            items
                .split(',')
                .map(|item| item.trim().trim_matches('"').to_string())
                .filter(|item| !item.is_empty())
                .collect()
        })
        .unwrap_or_default()
}

/// Every `.py` file beneath `directory`.
fn python_files(directory: &Path) -> Vec<PathBuf> {
    let mut found = Vec::new();
    let mut pending = vec![directory.to_path_buf()];
    while let Some(next) = pending.pop() {
        for path in std::fs::read_dir(next)
            .into_iter()
            .flatten()
            .flatten()
            .map(|e| e.path())
        {
            if path.is_dir() {
                pending.push(path);
            } else if path.extension().is_some_and(|extension| extension == "py") {
                found.push(path);
            }
        }
    }
    found
}

/// [PACKAGE-MANIFEST]: with `[package]` present, every `python_sources` path
/// lies inside one of the manifest's `source_roots`.
fn assert_sources_rooted(package: &Path, manifest: &str) {
    let roots = manifest_list(manifest, "source_roots");
    for source in manifest_list(manifest, "python_sources") {
        assert!(
            roots
                .iter()
                .any(|root| source.starts_with(&format!("{root}/"))),
            "{}: python source `{source}` is outside source_roots {roots:?}",
            package.display()
        );
    }
}

/// [PACKAGE-SOURCE-CANONICAL]: the registry refuses a release holding any
/// `.py` file whose first line is not the Python-flavor marker, so `CPython`
/// code (an oracle, a tool) can never ship inside a port.
fn assert_python_files_marked(package: &Path) {
    for file in python_files(package) {
        let text = std::fs::read_to_string(&file).unwrap_or_default();
        assert_eq!(
            text.lines().next(),
            Some("# osprey: flavor=python"),
            "{} must begin with `# osprey: flavor=python`",
            file.display()
        );
    }
}

#[test]
fn every_ported_package_is_a_publishable_python_flavor_package() {
    for suite in suites() {
        let package = suite.parent().map(Path::to_path_buf).unwrap_or_default();
        let manifest = std::fs::read_to_string(package.join("osprey.toml")).unwrap_or_default();
        assert!(
            manifest.contains("[package]") && manifest.contains("name = \"@osprey/"),
            "{} must declare an @osprey [package]",
            package.display()
        );
        assert!(
            package.join("README.md").is_file(),
            "{} needs a README.md naming the ported package and its divergences",
            package.display()
        );
        assert_sources_rooted(&package, &manifest);
        assert_python_files_marked(&package);
    }
}
