//! `osprey-rs` — the Rust front-end CLI.
//!
//! Strangler-fig stage: today it parses Osprey and dumps the AST or reports
//! syntax errors. As the port advances (`osprey-types`, `osprey-codegen`) this
//! grows the `--llvm` / `--compile` / `--run` surface of the Go `cli.go`.

use std::process::ExitCode;

fn main() -> ExitCode {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 2 {
        eprintln!("usage: osprey-rs <file.osp> [--ast | --check]");
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
        "--check" => {
            println!(
                "{path}: ok ({} statements)",
                parsed.program.statements.len()
            );
        }
        _ => {
            println!("{:#?}", parsed.program);
        }
    }
    ExitCode::SUCCESS
}
