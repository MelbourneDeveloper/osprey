//! Native end-to-end pins for the file-backed GPU loaders, [GPU-BUFFER-TENSOR]
//! and [GPU-BUFFER-LOAD], and for [GPU-MATMUL] over a loaded tensor.
//!
//! The corpus harness cannot read a committed file — wasm runs in a fresh
//! preopened sandbox and the mobile hosts chdir into a scratch directory — so
//! `fixtures/transformer.safetensors` is read here by an absolute path spliced
//! into both flavors of `fixtures/tensor_fixture`. The fixture is
//! `python3 -c "import json,struct;h={'w':{'dtype':'F32','shape':[2,2],'data_offsets':[0,16]},'h.1.b':{'dtype':'F16','shape':[2],'data_offsets':[16,20]}};j=json.dumps(h,separators=(',',':')).encode();open('transformer.safetensors','wb').write(struct.pack('<Q',len(j))+j+struct.pack('<4f',0.5,-1.0,2.0,3.0)+struct.pack('<2e',1.0,-2.5))"`:
//! a 116-byte header, so the data section starts at byte 124.

use std::path::{Path, PathBuf};
use std::process::Command;

/// `w` is `[0.5, -1.0, 2.0, 3.0]` as a 2×2, `h.1.b` the binary16 pair, the
/// miss names the tensor exactly, `gpuLoad` at the data offset sees `w`'s
/// bytes, and `w · w` in k-order is `[-1.75, -3.5, 7.0, 7.0]`.
const EXPECTED: &str = "w 0.5 -1.0 2.0 3.0\n\
    bias 1.0 -2.5\n\
    missing gpuLoadTensor: h.11.b: no tensor of that name in the header\n\
    raw 0.5 -1.0 2.0 3.0\n\
    square -1.75 -3.5 7.0 7.0\n";

const TWINS: [(&str, &str); 2] = [
    ("osp", include_str!("fixtures/tensor_fixture.osp")),
    ("ospml", include_str!("fixtures/tensor_fixture.ospml")),
];

fn manifest() -> &'static Path {
    Path::new(env!("CARGO_MANIFEST_DIR"))
}

/// Write `template` with the fixture's absolute path spliced in; a failed
/// write yields `None` and the test reports it through its own assertion.
fn program(ext: &str, template: &str) -> Option<PathBuf> {
    let fixture = manifest()
        .join("tests")
        .join("fixtures")
        .join("transformer.safetensors");
    let source = template.replace("@FIXTURE@", &fixture.display().to_string());
    let path = std::env::temp_dir().join(format!(
        "osprey_tensor_fixture_{}.{ext}",
        std::process::id()
    ));
    std::fs::write(&path, source).ok().map(|()| path)
}

/// `(exit code, stdout, stderr)` of `osprey <path> --run --quiet`, rooted at
/// the repo so the runtime archives resolve.
fn run(path: &Path) -> (Option<i32>, String, String) {
    let output = Command::new(env!("CARGO_BIN_EXE_osprey"))
        .current_dir(manifest().join("..").join(".."))
        .arg(path)
        .args(["--run", "--quiet"])
        .output();
    match output {
        Ok(out) => (
            out.status.code(),
            String::from_utf8_lossy(&out.stdout).into_owned(),
            String::from_utf8_lossy(&out.stderr).into_owned(),
        ),
        Err(_) => (None, String::new(), String::new()),
    }
}

#[test]
fn both_flavors_read_the_safetensors_fixture_exactly() {
    for (ext, template) in TWINS {
        let path = program(ext, template);
        assert!(path.is_some(), "{ext}: the program could not be written");
        let (code, stdout, stderr) = path.as_deref().map(run).unwrap_or_default();
        assert_eq!(code, Some(0), "{ext}: stderr: {stderr}");
        assert_eq!(stdout, EXPECTED, "{ext}");
    }
}
