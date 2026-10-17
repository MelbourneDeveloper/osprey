//! Execute the Python package against the compiler built by Cargo.

#[test]
#[cfg(any(target_os = "linux", target_os = "macos"))]
fn python_interoperability() {
    let root = std::path::Path::new(env!("CARGO_MANIFEST_DIR")).join("../..");
    let output = std::process::Command::new("python3")
        .args(["-m", "unittest", "discover", "-s", "python/tests", "-v"])
        .current_dir(&root)
        .env("PYTHONPATH", root.join("python"))
        .env("OSPREY", env!("CARGO_BIN_EXE_osprey"))
        .output()
        .expect("run Python interop tests");
    assert!(
        output.status.success(),
        "{}\n{}",
        String::from_utf8_lossy(&output.stdout),
        String::from_utf8_lossy(&output.stderr)
    );
}

#[test]
#[cfg(not(any(target_os = "linux", target_os = "macos")))]
fn python_target_rejects_unsupported_host() {
    let root = std::path::Path::new(env!("CARGO_MANIFEST_DIR")).join("../..");
    let output = std::process::Command::new(env!("CARGO_BIN_EXE_osprey"))
        .arg(root.join("examples/python/interop.osp"))
        .args(["--check", "--target=python"])
        .output()
        .expect("check Python target");
    assert!(!output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Linux or macOS"));
}
