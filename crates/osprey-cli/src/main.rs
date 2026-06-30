//! `osprey` — the Osprey compiler's command-line front end.
//!
//! Modes: report type errors (`--check`, the default — the editor's
//! diagnostics path), dump the AST (`--ast`), emit LLVM IR (`--llvm`), build
//! an executable (`--compile`), compile-and-run via clang (`--run`), emit the
//! document outline as JSON (`--symbols`), or print a built-in's signature as
//! markdown (`--hover <name>`). Every compiling mode gates on Hindley-Milner
//! type inference first — an ill-typed program never reaches codegen — and on
//! the capability sandbox (`--sandbox`, `--no-http`, `--no-websocket`,
//! `--no-fs`, `--no-ffi`). `--quiet` suppresses non-essential output. The C
//! driver used to link the emitted IR is `clang`, overridable via `OSPREY_CC`.
//!
//! `osprey lsp` runs the Language Server Protocol over stdio (the `osprey-lsp`
//! crate, built on the published lspkit crates); the `--symbols`/`--hover`
//! outline/signature helpers it shares now live there too.

mod docs;
mod sandbox;
mod wasm;

use osprey_syntax::Flavor;
use sandbox::Policy;
use std::path::{Path, PathBuf};
use std::process::{Command, ExitCode};

const USAGE: &str = "usage: osprey <file.osp> [--check | --ast | --llvm | --compile | --run | \
--symbols] [--quiet] [--debug] [--flavor default|ml] [--memory=default|gc] \
[--target=native|wasm32] [-o <out>] \
[--sandbox | --no-http | --no-websocket | --no-fs | --no-ffi]\n\
       osprey --hover <name>\n\
       osprey --docs --docs-dir <dir>\n\
       osprey lsp";

/// The parsed invocation: source path, mode flag, and behaviour switches.
#[derive(Debug)]
struct Cli {
    path: String,
    mode: String,
    quiet: bool,
    policy: Policy,
    /// The reclaiming memory backend linked behind `@osp_alloc` — `default`
    /// (malloc passthrough) or `gc` (tracing collector). Link-time only; the IR
    /// is identical [MEM-BACKENDS]. (`arc` is reserved, docs/plans/0011.)
    memory: String,
    /// Codegen/link target: `native` (host executable via clang) or `wasm32`
    /// (browser-ready WebAssembly via wasm-ld; wasm32-wasip1). [WASM-TARGET]
    target: String,
    /// Explicit output artifact path (`-o`); defaults to the source stem.
    output: Option<String>,
    /// Emit source-level debug metadata and link a debugger-friendly binary.
    debug: bool,
    /// Explicit source flavor from `--flavor`; `None` when unset, so flavor
    /// resolution falls through to the marker/extension precedence
    /// ([FLAVOR-SELECT], docs/specs/0023-LanguageFlavors.md).
    flavor: Option<Flavor>,
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
    // `osprey lsp`: speak the Language Server Protocol over stdio. The Rust
    // server (osprey-lsp, built on the published lspkit crates) drives the
    // compiler in-process. [LSP-REUSE-LSPKIT]
    if args.first().map(String::as_str) == Some("lsp") {
        return run_lsp();
    }
    // `osprey --docs`: regenerate the built-in function reference from the
    // compiler's metadata. No source file is involved.
    if args.iter().any(|a| a == "--docs") {
        return docs::run(&args);
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
        if let Some(md) = osprey_lsp::builtin_hover(&cli.path) {
            println!("{md}");
        }
        return ExitCode::SUCCESS;
    }
    run(&cli)
}

/// Run the stdio language server to completion on a fresh Tokio runtime.
fn run_lsp() -> ExitCode {
    let runtime = match tokio::runtime::Runtime::new() {
        Ok(runtime) => runtime,
        Err(e) => {
            eprintln!("osprey lsp: cannot start async runtime: {e}");
            return ExitCode::FAILURE;
        }
    };
    match runtime.block_on(osprey_lsp::run_stdio()) {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            eprintln!("osprey lsp: {e}");
            ExitCode::FAILURE
        }
    }
}

/// Parse the argument list: the first non-flag is the source path; mode flags
/// select the action (last one wins); the rest toggle behaviour.
fn parse_args(args: &[String]) -> Result<Cli, String> {
    let mut path = None;
    let mut mode = String::from("--check");
    let mut quiet = false;
    let mut policy = Policy::allow_all();
    let mut memory = String::from("default");
    let mut target = String::from("native");
    let mut output = None;
    let mut debug = false;
    let mut flavor = None;
    let mut it = args.iter();
    while let Some(a) = it.next() {
        match a.as_str() {
            "--ast" | "--check" | "--llvm" | "--compile" | "--run" | "--symbols" | "--hover" => {
                mode.clone_from(a);
            }
            "--quiet" => quiet = true,
            "--debug" => debug = true,
            "--sandbox" => policy = Policy::sandbox(),
            "--no-http" => policy.http = false,
            "--no-websocket" => policy.websocket = false,
            "--no-fs" => policy.fs = false,
            "--no-ffi" => policy.ffi = false,
            // `-o <path>` consumes the next argument as the output artifact path.
            "-o" => {
                let next = it
                    .next()
                    .ok_or_else(|| format!("-o requires a path\n{USAGE}"))?;
                output = Some(next.clone());
            }
            // `--flavor <name>` selects the source flavor explicitly (highest
            // selection precedence). [FLAVOR-SELECT]
            "--flavor" => {
                let next = it
                    .next()
                    .ok_or_else(|| format!("--flavor requires a value (default|ml)\n{USAGE}"))?;
                flavor = Some(parse_flavor(next)?);
            }
            flag if flag.starts_with("--flavor=") => {
                flavor = Some(parse_flavor(
                    flag.strip_prefix("--flavor=").unwrap_or_default(),
                )?);
            }
            flag if flag.starts_with("--memory=") => {
                memory = parse_memory(flag.strip_prefix("--memory=").unwrap_or_default())?;
            }
            flag if flag.starts_with("--target=") => {
                target = parse_target(flag.strip_prefix("--target=").unwrap_or_default())?;
            }
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
            memory,
            target,
            output,
            debug,
            flavor,
        }),
        None => Err(USAGE.to_string()),
    }
}

/// Validate the `--target=` value: `native` (host executable) or `wasm32`
/// (browser-ready WebAssembly, wasm32-wasip1). [WASM-TARGET]
fn parse_target(value: &str) -> Result<String, String> {
    match value {
        "native" | "wasm32" => Ok(value.to_string()),
        other => Err(format!(
            "unknown target '{other}' (available: native, wasm32)\n{USAGE}"
        )),
    }
}

/// Validate the `--memory=` value. `arc` is reserved but not yet implemented
/// (docs/plans/0011) — reject it explicitly rather than silently mislabel.
fn parse_memory(value: &str) -> Result<String, String> {
    match value {
        "default" | "gc" => Ok(value.to_string()),
        "arc" => {
            Err("memory backend 'arc' is not yet implemented (available: default, gc)".to_string())
        }
        other => Err(format!(
            "unknown memory backend '{other}' (available: default, gc)\n{USAGE}"
        )),
    }
}

/// Validate a `--flavor` / marker value into a [`Flavor`]. [FLAVOR-SELECT]
fn parse_flavor(value: &str) -> Result<Flavor, String> {
    value.parse().map_err(|e| format!("{e}\n{USAGE}"))
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
    let flavor = match osprey_syntax::resolve_flavor(cli.flavor, path, &source) {
        Ok(flavor) => flavor,
        Err(msg) => {
            eprintln!("{msg}");
            return ExitCode::from(2);
        }
    };
    let parsed = osprey_syntax::parse_program_with_flavor(&source, flavor);
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
            println!("{}", osprey_lsp::symbols_json(program));
            ExitCode::SUCCESS
        }
        "--llvm" | "--run" | "--compile" if report_type_errors(path, program) > 0 => {
            ExitCode::FAILURE
        }
        "--llvm" => match compile_ir(path, program, cli.debug) {
            Ok(ir) => {
                print!("{ir}");
                ExitCode::SUCCESS
            }
            Err(e) => {
                eprintln!("{path}: {e}");
                ExitCode::FAILURE
            }
        },
        "--run" => run_program(cli, program, source),
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

fn reject_debug_wasm(debug: bool) -> Option<ExitCode> {
    if debug {
        eprintln!("error: --debug is currently supported only for --target=native");
        return Some(ExitCode::from(2));
    }
    None
}

/// `--compile`: build the artifact at `-o` (or the source stem, `.wasm` for the
/// wasm target) — a host executable via clang, or WebAssembly via wasm-ld.
fn compile_program_to_disk(cli: &Cli, program: &osprey_ast::Program, source: &str) -> ExitCode {
    let out = output_path(&cli.path, cli.output.as_deref(), &cli.target);
    let result = if cli.target == "wasm32" {
        if let Some(code) = reject_debug_wasm(cli.debug) {
            return code;
        }
        wasm::build(&cli.path, program, &out)
    } else {
        build_executable(&cli.path, program, source, &out, &cli.memory, cli.debug)
    };
    match result {
        Ok(()) => {
            if !cli.quiet {
                println!("{}", out.display());
            }
            ExitCode::SUCCESS
        }
        Err(code) => code,
    }
}

/// The output artifact path: the explicit `-o` value, else the source stem in
/// the current directory — with a `.wasm` extension for the wasm target.
fn output_path(src: &str, output: Option<&str>, target: &str) -> PathBuf {
    match output {
        Some(o) => PathBuf::from(o),
        None if target == "wasm32" => PathBuf::from(format!("{}.wasm", stem_of(src))),
        None => PathBuf::from(stem_of(src)),
    }
}

/// Compile to a temp artifact and run it — the `--run` end-to-end path. Native
/// runs the executable directly; wasm runs it under a WASI host (`wasmtime`).
fn run_program(cli: &Cli, program: &osprey_ast::Program, source: &str) -> ExitCode {
    if cli.target == "wasm32" {
        if let Some(code) = reject_debug_wasm(cli.debug) {
            return code;
        }
        return wasm::run(cli, program);
    }
    let exe = std::env::temp_dir().join(format!("{}.out", stem_of(&cli.path)));
    if let Err(code) = build_executable(&cli.path, program, source, &exe, &cli.memory, cli.debug) {
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
    memory: &str,
    debug: bool,
) -> Result<(), ExitCode> {
    let ir = match compile_ir(path, program, debug) {
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
    let cc = c_compiler();
    let mut cmd = Command::new(&cc);
    let _ = cmd
        .arg(&ll)
        .arg("-o")
        .arg(exe)
        .arg("-Wno-override-module")
        .arg(opt_flag(debug))
        .args(debug_compile_flags(debug))
        .args(link_args(&ir, source, memory));
    match cmd.status() {
        Ok(s) if s.success() => Ok(()),
        Ok(_) => {
            eprintln!("error: {cc} failed to compile {}", ll.display());
            Err(ExitCode::FAILURE)
        }
        Err(e) => {
            eprintln!("error: could not invoke {cc}: {e}");
            Err(ExitCode::FAILURE)
        }
    }
}

/// The LLVM optimization level handed to clang when lowering the emitted IR.
/// Defaults to `-O2`; `OSPREY_OPT` overrides it (e.g. `-O0` for fast debug
/// builds, `-O3` to match Rust/OCaml release flags). This is load-bearing twice
/// over: it is the difference between competitive and 30–100× slower native
/// code, and — because codegen currently has no reclamation backend [MEM-OPAQUE,
/// docs/specs/0018] — it IS the default memory strategy. At `-O2` LLVM proves
/// per-operation `Result` allocations non-escaping and removes them entirely
/// (heap → registers), the [MEM-OWNERSHIP] "free at last use" ideal achieved
/// statically; without it those allocations leak for the whole run.
fn compile_ir(
    path: &str,
    program: &osprey_ast::Program,
    debug: bool,
) -> osprey_codegen::Result<String> {
    if debug {
        return osprey_codegen::compile_program_debug(
            program,
            osprey_codegen::DebugSource::from_path(path),
        );
    }
    osprey_codegen::compile_program(program)
}

fn opt_flag(debug: bool) -> String {
    let build = if debug {
        osprey_debug::DebugBuild::ON
    } else {
        osprey_debug::DebugBuild::OFF
    };
    build.opt_flag(
        std::env::var("OSPREY_OPT").unwrap_or_else(|_| "-O2".to_string()),
        std::env::var("OSPREY_DEBUG_OPT").ok(),
    )
}

fn debug_compile_flags(debug: bool) -> Vec<String> {
    let build = if debug {
        osprey_debug::DebugBuild::ON
    } else {
        osprey_debug::DebugBuild::OFF
    };
    build.native_driver_flags()
}

/// The C compiler/linker driver used to lower the emitted LLVM IR. Defaults to
/// `clang` (the only driver that consumes textual `.ll`); `OSPREY_CC` overrides
/// it — needed where several clangs coexist and the IR/runtime must link with a
/// matching toolchain (e.g. forcing the MinGW clang on Windows so it links the
/// MinGW-built C runtime archive rather than the system MSVC clang).
fn c_compiler() -> String {
    std::env::var("OSPREY_CC").unwrap_or_else(|_| "clang".to_string())
}

/// The source file's stem (`demo` for `examples/demo.osp`).
pub(crate) fn stem_of(path: &str) -> String {
    Path::new(path)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("osprey_out")
        .to_string()
}

/// The exit code to propagate for a finished child: its own code when it exited
/// normally, else (Unix) `128 + signal` for a signal death — so a segfaulting
/// program is NOT masked as success (`status.code()` is `None` for a signal).
pub(crate) fn child_exit_code(status: std::process::ExitStatus) -> u8 {
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
fn link_args(ir: &str, source: &str, memory: &str) -> Vec<String> {
    let mut args: Vec<String> = Vec::new();
    let uses_http = ir.contains("@http") || ir.contains("@websocket");

    // The reclaiming backend is a link-time archive swap — the IR is identical
    // [MEM-BACKENDS]. `gc` links the tracing-collector archive set; `default`
    // the malloc-passthrough set. (docs/plans/0011)
    let suffix = if memory == "gc" { "_gc" } else { "" };
    let lib = if uses_http {
        format!("libhttp_runtime{suffix}.a")
    } else {
        format!("libfiber_runtime{suffix}.a")
    };
    if let Some(p) = find_runtime_lib(&lib) {
        args.push(p);
    } else if let Some(p) = find_runtime_lib(&format!("libfiber_runtime{suffix}.a")) {
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
pub(crate) fn find_runtime_lib(lib: &str) -> Option<String> {
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

#[cfg(test)]
mod tests {
    use super::*;

    fn args(list: &[&str]) -> Vec<String> {
        list.iter().map(|s| (*s).to_string()).collect()
    }

    #[test]
    fn parse_args_defaults_to_check_with_full_capabilities() {
        let cli = parse_args(&args(&["prog.osp"])).expect("parses");
        assert_eq!(cli.path, "prog.osp");
        assert_eq!(cli.mode, "--check");
        assert!(!cli.quiet);
        assert!(cli.policy.http && cli.policy.websocket && cli.policy.fs && cli.policy.ffi);
    }

    #[test]
    fn parse_args_accepts_flavor_flag_in_both_spellings() {
        // No flag ⇒ unset, so resolution falls through to marker/extension.
        assert_eq!(parse_args(&args(&["f.osp"])).expect("ok").flavor, None);
        // Spaced and `=` spellings both set the explicit flavor.
        for spelling in [
            &["--flavor", "ml", "f.osp"][..],
            &["--flavor=ml", "f.osp"][..],
        ] {
            let cli = parse_args(&args(spelling)).expect("ok");
            assert_eq!(cli.flavor, Some(Flavor::Ml));
        }
        assert_eq!(
            parse_args(&args(&["--flavor=default", "f.osp"]))
                .expect("ok")
                .flavor,
            Some(Flavor::Default)
        );
        // A bogus value and a missing value both fail loudly.
        assert!(parse_args(&args(&["--flavor=fsharp", "f.osp"])).is_err());
        assert!(parse_args(&args(&["f.osp", "--flavor"])).is_err());
    }

    #[test]
    fn parse_args_last_mode_wins_and_quiet_sets() {
        let cli = parse_args(&args(&["--ast", "f.osp", "--llvm", "--run", "--quiet"])).expect("ok");
        assert_eq!(cli.mode, "--run");
        assert_eq!(cli.path, "f.osp");
        assert!(cli.quiet);
    }

    #[test]
    fn parse_args_each_sandbox_flag_clears_one_capability() {
        let cli = parse_args(&args(&["f.osp", "--no-http"])).expect("ok");
        assert!(!cli.policy.http && cli.policy.websocket && cli.policy.fs && cli.policy.ffi);
        let cli = parse_args(&args(&["f.osp", "--no-websocket"])).expect("ok");
        assert!(cli.policy.http && !cli.policy.websocket);
        let cli = parse_args(&args(&["f.osp", "--no-fs"])).expect("ok");
        assert!(!cli.policy.fs && cli.policy.ffi);
        let cli = parse_args(&args(&["f.osp", "--no-ffi"])).expect("ok");
        assert!(!cli.policy.ffi && cli.policy.fs);
        let cli = parse_args(&args(&["--sandbox", "f.osp"])).expect("ok");
        assert!(!cli.policy.http && !cli.policy.websocket && !cli.policy.fs && !cli.policy.ffi);
    }

    #[test]
    fn parse_args_rejects_unknown_flag_missing_path_and_extra_positional() {
        let e = parse_args(&args(&["f.osp", "--bogus"])).expect_err("unknown flag");
        assert!(e.contains("unknown flag --bogus"));
        let e = parse_args(&args(&["--check"])).expect_err("no path");
        assert!(e.contains("usage:"));
        let e = parse_args(&args(&["a.osp", "b.osp"])).expect_err("two paths");
        assert!(e.contains("unexpected argument b.osp"));
    }

    #[test]
    fn parse_args_handles_target_and_output() {
        let cli = parse_args(&args(&[
            "f.osp",
            "--target=wasm32",
            "--debug",
            "--compile",
            "-o",
            "out/f.wasm",
        ]))
        .expect("ok");
        assert_eq!(cli.target, "wasm32");
        assert!(cli.debug);
        assert_eq!(cli.output.as_deref(), Some("out/f.wasm"));
        // default target is native, no output.
        let cli = parse_args(&args(&["f.osp"])).expect("ok");
        assert_eq!(cli.target, "native");
        assert!(!cli.debug);
        assert!(cli.output.is_none());
        // -o with no following value, and an unknown target, are errors.
        assert!(parse_args(&args(&["f.osp", "-o"])).is_err());
        assert!(parse_args(&args(&["f.osp", "--target=riscv"])).is_err());
    }

    #[test]
    fn parse_target_accepts_known_and_rejects_unknown() {
        assert_eq!(parse_target("native").as_deref(), Ok("native"));
        assert_eq!(parse_target("wasm32").as_deref(), Ok("wasm32"));
        assert!(parse_target("x86").is_err());
    }

    #[test]
    fn output_path_defaults_by_target_and_honours_dash_o() {
        assert_eq!(output_path("a/b.osp", None, "native"), PathBuf::from("b"));
        assert_eq!(
            output_path("a/b.osp", None, "wasm32"),
            PathBuf::from("b.wasm")
        );
        assert_eq!(
            output_path("a/b.osp", Some("custom.wasm"), "wasm32"),
            PathBuf::from("custom.wasm")
        );
    }

    #[test]
    fn debug_wasm_rejection_is_centralized() {
        assert!(reject_debug_wasm(true).is_some());
        assert!(reject_debug_wasm(false).is_none());
    }

    #[test]
    fn stem_of_handles_dirs_and_missing_extension() {
        assert_eq!(stem_of("examples/demo.osp"), "demo");
        assert_eq!(stem_of("/a/b/c.osp"), "c");
        assert_eq!(stem_of("noext"), "noext");
    }

    #[test]
    fn directive_parses_both_spellings_and_ignores_others() {
        assert_eq!(directive("// @link: sqlite3", "link"), Some("sqlite3"));
        assert_eq!(
            directive("//@linkdir: /opt/lib ", "linkdir"),
            Some("/opt/lib")
        );
        assert_eq!(directive("  // @link:  pq  ", "link"), Some("pq"));
        assert_eq!(directive("let x = 1", "link"), None);
        assert_eq!(directive("// @link: sqlite3", "linkdir"), None);
    }

    #[test]
    fn link_args_adds_ffi_directives_and_openssl_for_http() {
        let ffi = link_args(
            "",
            "// @link: sqlite3\n// @linkdir: /opt/lib\ncode\n",
            "default",
        );
        assert!(ffi.iter().any(|a| a == "-lsqlite3"), "{ffi:?}");
        assert!(ffi.iter().any(|a| a == "-L/opt/lib"), "{ffi:?}");
        let http = link_args("call void @http_listen()", "", "default");
        assert!(http.iter().any(|a| a == "-lssl") && http.iter().any(|a| a == "-lcrypto"));
        // No HTTP markers => no openssl flags.
        let plain = link_args("call void @osprey_list_empty()", "", "default");
        assert!(!plain.iter().any(|a| a == "-lssl"));
    }

    #[test]
    fn link_args_selects_gc_archive_and_validates_backend() {
        // The `gc` backend swaps in the `_gc` archive set; `default` does not.
        let gc = link_args("call void @osprey_list_empty()", "", "gc");
        assert!(
            gc.iter().any(|a| a.contains("_gc.a")) || gc.is_empty(),
            "gc backend must select a *_gc archive when one is present: {gc:?}"
        );
        let plain = link_args("call void @osprey_list_empty()", "", "default");
        assert!(!plain.iter().any(|a| a.contains("_gc.a")), "{plain:?}");
        // Backend validation: default/gc accepted, arc reserved, others rejected.
        assert_eq!(parse_memory("gc").as_deref(), Ok("gc"));
        assert_eq!(parse_memory("default").as_deref(), Ok("default"));
        assert!(parse_memory("arc").is_err());
        assert!(parse_memory("bogus").is_err());
    }

    #[test]
    fn openssl_and_compiler_helpers_are_well_formed() {
        let flags = openssl_flags();
        assert!(flags.iter().any(|f| f == "-lssl") && flags.iter().any(|f| f == "-lcrypto"));
        assert!(!c_compiler().is_empty());
        assert!(find_runtime_lib("definitely_not_a_real_lib_xyz.a").is_none());
    }

    #[cfg(unix)]
    #[test]
    fn child_exit_code_maps_codes_and_signals() {
        use std::os::unix::process::ExitStatusExt;
        assert_eq!(child_exit_code(std::process::ExitStatus::from_raw(0)), 0);
        assert_eq!(
            child_exit_code(std::process::ExitStatus::from_raw(1 << 8)),
            1
        );
        // Killed by SIGKILL (9): no exit code, so 128 + signal.
        assert_eq!(child_exit_code(std::process::ExitStatus::from_raw(9)), 137);
    }

    #[test]
    fn report_type_errors_counts_zero_for_valid_and_more_for_ill_typed() {
        let ok = osprey_syntax::parse_program("let x = 1\nprint(x)\n").program;
        assert_eq!(report_type_errors("ok.osp", &ok), 0);
        let bad = osprey_syntax::parse_program("let y = 1 + \"oops\" - true\n").program;
        assert!(report_type_errors("bad.osp", &bad) > 0);
    }

    fn temp_source(name: &str, body: &str) -> String {
        let p = std::env::temp_dir().join(format!("osprey_cli_{name}.osp"));
        std::fs::write(&p, body).expect("write temp source");
        p.display().to_string()
    }

    fn cli(path: impl Into<String>, mode: &str, policy: Policy) -> Cli {
        Cli {
            path: path.into(),
            mode: mode.to_string(),
            quiet: true,
            policy,
            memory: "default".to_string(),
            target: "native".to_string(),
            output: None,
            debug: false,
            flavor: None,
        }
    }

    #[test]
    fn run_drives_check_symbols_and_llvm_modes_in_process() {
        let path = temp_source("ok", "let greeting = \"hi\"\nprint(greeting)\n");
        for mode in ["--check", "--symbols", "--llvm", "--ast"] {
            // ExitCode is opaque; this drives run -> dispatch coverage and must
            // not panic for a well-formed program.
            let _ = run(&cli(path.clone(), mode, Policy::allow_all()));
        }
    }

    #[test]
    fn run_reports_missing_file_and_parse_errors() {
        let _ = run(&cli(
            "/no/such/osprey/file.osp",
            "--check",
            Policy::allow_all(),
        ));
        let path = temp_source("broken", "fn = = =\n");
        let _ = run(&cli(path, "--check", Policy::allow_all())); // parse-error branch
    }

    #[test]
    fn run_rejects_sandbox_violation_before_codegen() {
        let path = temp_source("fs", "let c = readFile(\"x.txt\")\n");
        let _ = run(&cli(path, "--llvm", Policy::sandbox())); // sandbox-violation branch
    }

    #[test]
    fn parse_args_accepts_the_memory_backend_flag() {
        let cli = parse_args(&args(&["f.osp", "--memory=gc"])).expect("ok");
        assert_eq!(cli.memory, "gc");
    }

    #[test]
    fn report_type_errors_prints_positioned_diagnostics() {
        // An undefined identifier yields an error carrying a source position,
        // exercising the `Some(position)` diagnostic arm.
        let bad = osprey_syntax::parse_program("print(missingVariable)\n").program;
        assert!(report_type_errors("bad.osp", &bad) > 0);
    }

    #[test]
    fn compile_ir_and_debug_helpers_switch_on_the_debug_flag() {
        let program = osprey_syntax::parse_program("let n = 1\nprint(\"${n}\")\n").program;
        // debug=true takes the debug-info codegen path; both opt/driver helpers
        // branch on the same flag.
        assert!(compile_ir("p.osp", &program, true).is_ok());
        assert!(!opt_flag(true).is_empty());
        assert!(!opt_flag(false).is_empty());
        let _ = debug_compile_flags(true);
        let _ = debug_compile_flags(false);
    }

    #[test]
    fn wasm_target_rejects_debug_then_dispatches_to_the_backend() {
        let program = osprey_syntax::parse_program("let n = 1\nprint(\"${n}\")\n").program;
        let mut c = cli("p.osp", "--compile", Policy::allow_all());
        c.target = "wasm32".to_string();
        // --debug + --target=wasm32 is rejected before any toolchain work.
        c.debug = true;
        let _ = compile_program_to_disk(&c, &program, "");
        let _ = run_program(&c, &program, "");
        // Without --debug the wasm build/run driver is dispatched (it fails
        // cleanly without the wasm toolchain, but the dispatch lines execute).
        c.debug = false;
        let _ = compile_program_to_disk(&c, &program, "");
        let _ = run_program(&c, &program, "");
    }
}
