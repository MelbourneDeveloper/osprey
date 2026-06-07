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
    if args.len() < 2 {
        eprintln!("usage: osprey-rs <file.osp> [--ast | --check | --llvm | --run]");
        return ExitCode::from(2);
    }

    if args[1] == "--version" {
        println!("osprey-rs 0.0.0-dev");
        return ExitCode::SUCCESS;
    }

    let path = &args[1];
    let mode = args.get(2).map(String::as_str).unwrap_or("--ast");

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
        "--run" => run_program(path, &parsed.program),
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

/// Compile the program to LLVM IR, hand it to clang, and run the resulting
/// binary — the `--run` end-to-end path.
fn run_program(path: &str, program: &osprey_ast::Program) -> ExitCode {
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

    match Command::new("clang")
        .arg(&ll)
        .arg("-o")
        .arg(&exe)
        .arg("-Wno-override-module")
        .status()
    {
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
        Ok(s) => ExitCode::from(s.code().unwrap_or(0) as u8),
        Err(e) => {
            eprintln!("error: could not run {}: {e}", exe.display());
            ExitCode::FAILURE
        }
    }
}
