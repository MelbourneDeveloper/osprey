//! Compile the C runtime with the *same* hardening flags the Makefile uses
//! (`-D_FORTIFY_SOURCE=2 -fstack-protector-strong`, warnings-as-errors) via the
//! `cc` crate. No C is rewritten — the C sources in `compiler/runtime` stay the
//! single implementation.
//!
//! Only the self-contained, dependency-free units are compiled here (the generic
//! FFI pointer cells). The concurrency/HTTP units (pthreads/OpenSSL) link the
//! same way and are added as their crates come online.

use std::path::PathBuf;

fn main() {
    let runtime = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("../../compiler/runtime")
        .canonicalize()
        .expect("locate compiler/runtime");

    let ffi = runtime.join("ffi_runtime.c");
    println!("cargo:rerun-if-changed={}", ffi.display());

    let mut build = cc::Build::new();
    build
        .file(&ffi)
        .opt_level(2)
        .define("_FORTIFY_SOURCE", "2")
        .flag_if_supported("-fstack-protector-strong")
        .flag_if_supported("-std=c11")
        .warnings(true);
    build.compile("osprey_runtime_ffi");
}
