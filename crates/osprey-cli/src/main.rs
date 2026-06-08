//! `osprey-rs` — the Rust front-end CLI.
//!
//! Strangler-fig stage: it parses Osprey, dumps the AST (`--ast`), runs
//! Hindley-Milner type inference (`--check`), emits LLVM IR (`--llvm`), or
//! compiles-and-runs via clang (`--run`). As the port advances this grows the
//! rest of the Go `cli.go` surface.

use std::path::Path;
use std::process::{Command, ExitCode};

fn main() -> ExitCode {
    let args: Vec<String> = std::env::args().collect();
    let Some(path) = args.get(1) else {
        eprintln!("usage: osprey-rs <file.osp> [--ast | --check | --llvm | --run]");
        return ExitCode::from(2);
    };

    if path == "--version" {
        println!("osprey-rs 0.0.0-dev");
        return ExitCode::SUCCESS;
    }

    let mode = args.get(2).map_or("--ast", String::as_str);

    let source = match std::fs::read_to_string(path) {
        Ok(s) => s,
        Err(e) => {
            eprintln!("error: cannot read {path}: {e}");
            return ExitCode::from(2);
        }
    };

    let parsed = osprey_syntax::parse_program(&source);
    if !parsed.errors.is_empty() {
        for err in &parsed.errors {
            eprintln!(
                "{path}:{}:{}: {}",
                err.position.line, err.position.column, err.message
            );
        }
        return ExitCode::FAILURE;
    }

    match mode {
        "--check" => run_check(path, &parsed.program),
        "--llvm" => match osprey_codegen::compile_program(&parsed.program) {
            Ok(ir) => {
                print!("{ir}");
                ExitCode::SUCCESS
            }
            Err(e) => {
                eprintln!("{path}: {e}");
                ExitCode::FAILURE
            }
        },
        "--run" => run_program(path, &parsed.program, &source),
        _ => {
            println!("{:#?}", parsed.program);
            ExitCode::SUCCESS
        }
    }
}

fn run_check(path: &str, program: &osprey_ast::Program) -> ExitCode {
    let errors = osprey_types::check_program(program);
    if errors.is_empty() {
        println!("{path}: ok ({} statements)", program.statements.len());
        return ExitCode::SUCCESS;
    }
    for e in &errors {
        match e.position {
            Some(p) => eprintln!("{path}:{}:{}: {}", p.line, p.column, e.message),
            None => eprintln!("{path}: {}", e.message),
        }
    }
    ExitCode::FAILURE
}

/// Compile the program to LLVM IR, hand it to clang together with the prebuilt
/// C runtime, and run the resulting binary — the `--run` end-to-end path.
fn run_program(path: &str, program: &osprey_ast::Program, source: &str) -> ExitCode {
    let ir = match osprey_codegen::compile_program(program) {
        Ok(ir) => ir,
        Err(e) => {
            eprintln!("{path}: {e}");
            return ExitCode::FAILURE;
        }
    };

    let stem = Path::new(path)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("osprey_rs");
    let dir = std::env::temp_dir();
    let ll = dir.join(format!("{stem}.ll"));
    let exe = dir.join(format!("{stem}.out"));

    if let Err(e) = std::fs::write(&ll, ir.as_bytes()) {
        eprintln!("error: cannot write IR to {}: {e}", ll.display());
        return ExitCode::FAILURE;
    }

    let mut cmd = Command::new("clang");
    let _ = cmd
        .arg(&ll)
        .arg("-o")
        .arg(&exe)
        .arg("-Wno-override-module")
        .args(link_args(&ir, source));

    match cmd.status() {
        Ok(s) if s.success() => {}
        Ok(_) => {
            eprintln!("error: clang failed to compile {}", ll.display());
            return ExitCode::FAILURE;
        }
        Err(e) => {
            eprintln!("error: could not invoke clang: {e}");
            return ExitCode::FAILURE;
        }
    }

    match Command::new(&exe).status() {
        Ok(s) => ExitCode::from(u8::try_from(s.code().unwrap_or(0)).unwrap_or(1)),
        Err(e) => {
            eprintln!("error: could not run {}: {e}", exe.display());
            ExitCode::FAILURE
        }
    }
}

/// Assemble the link arguments: the prebuilt C runtime static library (the HTTP
/// superset when the program touches HTTP/WebSocket, else the fiber runtime),
/// OpenSSL for HTTP, and any `// @link:` / `// @linkdir:` FFI directives — the
/// same surface `jit_executor.go` builds.
fn link_args(ir: &str, source: &str) -> Vec<String> {
    let mut args: Vec<String> = Vec::new();
    let uses_http = ir.contains("@http") || ir.contains("@websocket");

    let lib = if uses_http {
        "libhttp_runtime.a"
    } else {
        "libfiber_runtime.a"
    };
    if let Some(p) = find_runtime_lib(lib) {
        args.push(p);
    } else if let Some(p) = find_runtime_lib("libfiber_runtime.a") {
        args.push(p);
    }

    if uses_http {
        args.extend(openssl_flags());
    }

    // FFI directives: `// @link: sqlite3` -> `-lsqlite3`, `// @linkdir: P` -> `-LP`.
    for line in source.lines() {
        let t = line.trim();
        if let Some(rest) = t
            .strip_prefix("// @link:")
            .or_else(|| t.strip_prefix("//@link:"))
        {
            args.push(format!("-l{}", rest.trim()));
        } else if let Some(rest) = t
            .strip_prefix("// @linkdir:")
            .or_else(|| t.strip_prefix("//@linkdir:"))
        {
            args.push(format!("-L{}", rest.trim()));
        }
    }
    args
}

/// Search the conventional install/build locations for a runtime static lib.
fn find_runtime_lib(lib: &str) -> Option<String> {
    let mut roots = vec![
        format!("compiler/bin/{lib}"),
        format!("bin/{lib}"),
        format!("../bin/{lib}"),
        format!("../../bin/{lib}"),
        format!("/usr/local/lib/{lib}"),
    ];
    if let Ok(wd) = std::env::current_dir() {
        roots.push(wd.join("compiler/bin").join(lib).display().to_string());
        roots.push(wd.join("bin").join(lib).display().to_string());
    }
    roots.into_iter().find(|p| Path::new(p).exists())
}

/// OpenSSL link flags, mirroring `addOpenSSLFlags`.
fn openssl_flags() -> Vec<String> {
    for dir in [
        "/opt/homebrew/opt/openssl@3/lib",
        "/opt/homebrew/lib",
        "/usr/local/opt/openssl@3/lib",
        "/usr/local/lib",
    ] {
        if Path::new(dir).join("libssl.dylib").exists() {
            return vec![format!("-L{dir}"), "-lssl".into(), "-lcrypto".into()];
        }
    }
    vec!["-lssl".into(), "-lcrypto".into()]
}
