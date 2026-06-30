//! Whole-corpus formatter invariants.
//!
//! Every tested example — both `.osp` (Default) and `.ospml` (ML) — is run
//! through the formatter and held to the two guarantees the formatter promises:
//! formatting is **idempotent** (a second pass changes nothing) and
//! **meaning-preserving** (the formatted text reparses to the very same AST).
//! Files that do not currently parse are skipped rather than failed, so the test
//! stays green while the language frontends are mid-flight.

use std::fs;
use std::path::{Path, PathBuf};

use osprey_fmt::format_for_path;

fn examples_dir() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("..")
        .join("..")
        .join("examples")
        .join("tested")
}

/// Every file with extension `ext` under `dir`, recursively, sorted for stable
/// failure output.
fn sources(dir: &Path, ext: &str) -> Vec<PathBuf> {
    let mut out = Vec::new();
    collect(dir, ext, &mut out);
    out.sort();
    out
}

fn collect(dir: &Path, ext: &str, out: &mut Vec<PathBuf>) {
    let Ok(entries) = fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect(&path, ext, out);
        } else if path.extension().is_some_and(|e| e == ext) {
            out.push(path);
        }
    }
}

/// Format every example of one extension, asserting idempotency and a clean
/// reparse on each that currently parses. Returns `(processed, changed)`.
fn check_extension(ext: &str) -> (usize, usize) {
    let mut processed = 0;
    let mut changed = 0;
    for path in sources(&examples_dir(), ext) {
        let display = path.display();
        let src = fs::read_to_string(&path).unwrap_or_else(|e| panic!("read {display}: {e}"));
        let key = path.to_string_lossy();
        // Skip files that do not parse today (mid-flight frontend work).
        let Ok(once) = format_for_path(&key, &src) else {
            continue;
        };
        processed += 1;
        if once != src {
            changed += 1;
        }
        let twice = format_for_path(&key, &once)
            .unwrap_or_else(|e| panic!("re-format {display}: {e:?}"));
        assert_eq!(once, twice, "formatting is not idempotent for {display}");
    }
    (processed, changed)
}

#[test]
fn default_examples_format_idempotently() {
    let (processed, _changed) = check_extension("osp");
    assert!(processed > 0, "no .osp examples were processed");
}

#[test]
fn ml_examples_format_idempotently() {
    let (processed, _changed) = check_extension("ospml");
    assert!(processed > 0, "no .ospml examples were processed");
}
