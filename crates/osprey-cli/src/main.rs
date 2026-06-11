//! `osprey` — the Osprey compiler's command-line front end.
//!
//! Modes: report type errors (`--check`, the default — the editor's
//! diagnostics path), dump the AST (`--ast`), emit LLVM IR (`--llvm`), build
//! an executable (`--compile`), compile-and-run via clang (`--run`), emit the
//! document outline as JSON (`--symbols`), or print a built-in's signature as
//! markdown (`--hover <name>`). Every compiling mode gates on Hindley-Milner
//! type inference first — an ill-typed program never reaches codegen — and on
//! the capability sandbox (`--sandbox`, `--no-http`, `--no-websocket`,
//! `--no-fs`, `--no-ffi`). `--quiet` suppresses non-essential output.

mod lsp;
mod sandbox;

use sandbox::Policy;
use std::path::{Path, PathBuf};
use std::process::{Command, ExitCode};

const USAGE: &str = "usage: osprey <file.osp> [--check | --ast | --llvm | --compile | --run | \
--symbols] [--quiet] [--sandbox | --no-http | --no-websocket | --no-fs | --no-ffi]\n\
       osprey --hover <name>";

/// The parsed invocation: source path, mode flag, and behaviour switches.
struct Cli {
    path: String,
    mode: String,
    quiet: bool,
    policy: Policy,
}

fn main() -> ExitCode {
    let args: Vec<String> = std::env::args().skip(1).collect();
    if args.first().map(String::as_str) == Some("--version") {
        // [SWR-VERSION-BUILD-STAMPING] the real version is stamped from the git
        // tag at release-build time via OSPREY_VERSION; source stays 0.0.0-dev.
        // [SWR-VERSION-CLI-OUTPUT] `--json` emits the manifest form the VS Code
        // extension version-checks at activation.
        let version = option_env!("OSPREY_VERSION").unwrap_or("0.0.0-dev");
        if args.iter().any(|a| a == "--json") {
            println!(
                "{{\"manifestVersion\":1,\"name\":\"osprey\",\"version\":\"{version}\",\
\"kind\":\"cli\",\"product\":\"osprey\"}}"
            );
        } else {
            println!("osprey {version}");
        }
        return ExitCode::SUCCESS;
    }
    let cli = match parse_args(&args) {
        Ok(cli) => cli,
        Err(msg) => {
            eprintln!("{msg}");
            return ExitCode::from(2);
        }
    };
    if cli.mode == "--hover" {
        // The positional is a built-in NAME, not a file. Unknown names print
        // nothing (the editor simply shows no hover) and still exit 0.
        if let Some(md) = lsp::builtin_hover(&cli.path) {
            println!("{md}");
        }
        return ExitCode::SUCCESS;
    }
    run(&cli)
}

/// Parse the argument list: the first non-flag is the source path; mode flags
/// select the action (last one wins); the rest toggle behaviour.
fn parse_args(args: &[String]) -> Result<Cli, String> {
    let mut path = None;
    let mut mode = String::from("--check");
    let mut quiet = false;
    let mut policy = Policy::allow_all();
    for a in args {
        match a.as_str() {
            "--ast" | "--check" | "--llvm" | "--compile" | "--run" | "--symbols" | "--hover" => {
                mode.clone_from(a);
            }
            "--quiet" => quiet = true,
            "--sandbox" => policy = Policy::sandbox(),
            "--no-http" => policy.http = false,
            "--no-websocket" => policy.websocket = false,
            "--no-fs" => policy.fs = false,
            "--no-ffi" => policy.ffi = false,
            flag if flag.starts_with("--") => return Err(format!("unknown flag {flag}\n{USAGE}")),
            _ if path.is_none() => path = Some(a.clone()),
            other => return Err(format!("unexpected argument {other}\n{USAGE}")),
        }
    }
    match path {
        Some(path) => Ok(Cli {
            path,
            mode,
            quiet,
            policy,
        }),
        None => Err(USAGE.to_string()),
    }
}

/// Parse, gate (syntax → sandbox → types), and dispatch the selected mode.
fn run(cli: &Cli) -> ExitCode {
    let path = &cli.path;
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
    let violations = sandbox::violations(&parsed.program, cli.policy);
    if !violations.is_empty() {
        for v in &violations {
            eprintln!("{path}: {v}");
        }
        return ExitCode::FAILURE;
    }
    dispatch(cli, &parsed.program, &source)
}

/// Route the type-gated modes: an ill-typed program never reaches codegen.
fn dispatch(cli: &Cli, program: &osprey_ast::Program, source: &str) -> ExitCode {
    let path = &cli.path;
    match cli.mode.as_str() {
        "--check" => run_check(cli, program),
        // The outline must work for ill-typed (but parsable) files, so
        // `--symbols` deliberately skips the type gate.
        "--symbols" => {
            println!("{}", lsp::symbols_json(program));
            ExitCode::SUCCESS
        }
        "--llvm" | "--run" | "--compile" if report_type_errors(path, program) > 0 => {
            ExitCode::FAILURE
        }
        "--llvm" => match osprey_codegen::compile_program(program) {
            Ok(ir) => {
                print!("{ir}");
                ExitCode::SUCCESS
            }
            Err(e) => {
                eprintln!("{path}: {e}");
                ExitCode::FAILURE
            }
        },
        "--run" => run_program(path, program, source),
        "--compile" => compile_program_to_disk(cli, program, source),
        _ => {
            println!("{program:#?}");
            ExitCode::SUCCESS
        }
    }
}

/// Type-check `program`, print every error in `file:line:col: message` form,
/// and return how many there were. The shared gate for every compiling mode.
fn report_type_errors(path: &str, program: &osprey_ast::Program) -> usize {
    let errors = osprey_types::check_program(program);
    for e in &errors {
        match e.position {
            Some(p) => eprintln!("{path}:{}:{}: {}", p.line, p.column, e.message),
            None => eprintln!("{path}: {}", e.message),
        }
    }
    errors.len()
}

fn run_check(cli: &Cli, program: &osprey_ast::Program) -> ExitCode {
    if report_type_errors(&cli.path, program) == 0 {
        if !cli.quiet {
            println!("{}: ok ({} statements)", cli.path, program.statements.len());
        }
        return ExitCode::SUCCESS;
    }
    ExitCode::FAILURE
}

/// `--compile`: build an executable named after the source file, in the
/// current directory.
fn compile_program_to_disk(cli: &Cli, program: &osprey_ast::Program, source: &str) -> ExitCode {
    let exe = PathBuf::from(stem_of(&cli.path));
    match build_executable(&cli.path, program, source, &exe) {
        Ok(()) => {
            if !cli.quiet {
                println!("{}", exe.display());
            }
            ExitCode::SUCCESS
        }
        Err(code) => code,
    }
}

/// Compile to a temp executable and run it — the `--run` end-to-end path.
fn run_program(path: &str, program: &osprey_ast::Program, source: &str) -> ExitCode {
    let exe = std::env::temp_dir().join(format!("{}.out", stem_of(path)));
    if let Err(code) = build_executable(path, program, source, &exe) {
        return code;
    }
    match Command::new(&exe).status() {
        Ok(s) => ExitCode::from(child_exit_code(s)),
        Err(e) => {
            eprintln!("error: could not run {}: {e}", exe.display());
            ExitCode::FAILURE
        }
    }
}

/// Lower to LLVM IR and hand it to clang together with the prebuilt C runtime,
/// producing `exe`.
fn build_executable(
    path: &str,
    program: &osprey_ast::Program,
    source: &str,
    exe: &Path,
) -> Result<(), ExitCode> {
    let ir = match osprey_codegen::compile_program(program) {
        Ok(ir) => ir,
        Err(e) => {
            eprintln!("{path}: {e}");
            return Err(ExitCode::FAILURE);
        }
    };
    let ll = std::env::temp_dir().join(format!("{}.ll", stem_of(path)));
    if let Err(e) = std::fs::write(&ll, ir.as_bytes()) {
        eprintln!("error: cannot write IR to {}: {e}", ll.display());
        return Err(ExitCode::FAILURE);
    }
    let mut cmd = Command::new("clang");
    let _ = cmd
        .arg(&ll)
        .arg("-o")
        .arg(exe)
        .arg("-Wno-override-module")
        .args(link_args(&ir, source));
    match cmd.status() {
        Ok(s) if s.success() => Ok(()),
        Ok(_) => {
            eprintln!("error: clang failed to compile {}", ll.display());
            Err(ExitCode::FAILURE)
        }
        Err(e) => {
            eprintln!("error: could not invoke clang: {e}");
            Err(ExitCode::FAILURE)
        }
    }
}

/// The source file's stem (`demo` for `examples/demo.osp`).
fn stem_of(path: &str) -> String {
    Path::new(path)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("osprey_out")
        .to_string()
}

/// The exit code to propagate for a finished child: its own code when it exited
/// normally, else (Unix) `128 + signal` for a signal death — so a segfaulting
/// program is NOT masked as success (`status.code()` is `None` for a signal).
fn child_exit_code(status: std::process::ExitStatus) -> u8 {
    if let Some(code) = status.code() {
        return u8::try_from(code).unwrap_or(1);
    }
    #[cfg(unix)]
    {
        use std::os::unix::process::ExitStatusExt;
        if let Some(sig) = status.signal() {
            return 128u8.saturating_add(u8::try_from(sig).unwrap_or(0));
        }
    }
    1
}

/// Assemble the link arguments — everything a compiled binary needs beyond
/// libc: the prebuilt C runtime static library (the HTTP superset when the
/// program touches HTTP/WebSocket, else the fiber runtime), OpenSSL for HTTP,
/// and any `// @link:` / `// @linkdir:` FFI directives (e.g. `-lsqlite3`).
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

    // Windows (MinGW UCRT64): the C runtime's fibers are winpthreads-backed, so
    // `pthread_*` must be linked explicitly — unlike Linux/macOS where libc /
    // libSystem provide them implicitly. Must come AFTER the archive that
    // references them. Compiled out on Unix.
    #[cfg(windows)]
    {
        args.push("-lpthread".to_string());
    }

    // FFI directives: `// @link: sqlite3` -> `-lsqlite3`, `// @linkdir: P` -> `-LP`.
    for line in source.lines() {
        if let Some(lib) = directive(line, "link") {
            args.push(format!("-l{lib}"));
        } else if let Some(dir) = directive(line, "linkdir") {
            args.push(format!("-L{dir}"));
        }
    }
    args
}

/// The trimmed value of a `// @<key>:` FFI directive line (accepting the
/// space-less `//@<key>:` spelling too), or `None` if `line` is not one.
fn directive<'a>(line: &'a str, key: &str) -> Option<&'a str> {
    let t = line.trim();
    t.strip_prefix(&format!("// @{key}:"))
        .or_else(|| t.strip_prefix(&format!("//@{key}:")))
        .map(str::trim)
}

/// Search the conventional install/build locations for a runtime static lib:
/// the working directory's repo layout, then next to the `osprey` executable
/// (the release-tarball layout, and `target/release` two levels under the repo
/// root), then the system lib dir.
fn find_runtime_lib(lib: &str) -> Option<String> {
    let mut roots = vec![
        format!("compiler/bin/{lib}"),
        format!("compiler/lib/{lib}"),
        format!("bin/{lib}"),
        format!("../bin/{lib}"),
        format!("../../bin/{lib}"),
    ];
    if let Some(dir) = std::env::current_exe()
        .ok()
        .and_then(|e| e.parent().map(std::path::Path::to_path_buf))
    {
        roots.push(dir.join(lib).display().to_string());
        for up in ["../../compiler/lib", "../../compiler/bin"] {
            roots.push(dir.join(up).join(lib).display().to_string());
        }
    }
    roots.push(format!("/usr/local/lib/{lib}"));
    roots.into_iter().find(|p| Path::new(p).exists())
}

/// OpenSSL link flags, searching the conventional Homebrew/system lib dirs.
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
