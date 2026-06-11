//! The capability sandbox: `--sandbox` / `--no-http` / `--no-websocket` /
//! `--no-fs` / `--no-ffi` gate which built-in capabilities a program may use.
//! Enforcement is a pre-codegen pass over the parsed program: a gated builtin
//! referenced anywhere (or an `extern` declaration under `--no-ffi`) is a
//! compile error, so untrusted code is rejected before any IR exists.

use osprey_ast::{Program, Stmt};

/// Which capability groups the invocation allows. Everything defaults to on;
/// `--sandbox` turns every group off at once.
#[expect(
    clippy::struct_excessive_bools,
    reason = "a capability sandbox is by nature a set of independent on/off switches"
)]
#[derive(Clone, Copy)]
pub(crate) struct Policy {
    pub http: bool,
    pub websocket: bool,
    pub fs: bool,
    pub ffi: bool,
    pub process: bool,
}

impl Policy {
    pub(crate) fn allow_all() -> Policy {
        Policy {
            http: true,
            websocket: true,
            fs: true,
            ffi: true,
            process: true,
        }
    }

    /// `--sandbox`: every risky capability off.
    pub(crate) fn sandbox() -> Policy {
        Policy {
            http: false,
            websocket: false,
            fs: false,
            ffi: false,
            process: false,
        }
    }
}

const HTTP_FNS: &[&str] = &[
    "httpCreateServer",
    "httpListen",
    "httpStopServer",
    "httpCreateClient",
    "httpGet",
    "httpPost",
    "httpPut",
    "httpDelete",
    "httpCloseClient",
    "httpGetResponse",
    "httpResponseStatus",
    "httpResponseBody",
    "httpResponseHeader",
    "httpResponseFree",
];
const WEBSOCKET_FNS: &[&str] = &[
    "websocketConnect",
    "websocketSend",
    "websocketKeepAlive",
    "websocketClose",
    "websocketCreateServer",
    "websocketListen",
    "websocketServerBroadcast",
    "websocketStopServer",
];
const FS_FNS: &[&str] = &["readFile", "writeFile"];
const PROCESS_FNS: &[&str] = &["spawnProcess", "awaitProcess", "cleanupProcess"];

/// Every policy violation in `program`, as ready-to-print messages. Empty means
/// the program is allowed to compile under `policy`.
pub(crate) fn violations(program: &Program, policy: Policy) -> Vec<String> {
    let idents = osprey_codegen::referenced_idents(program);
    let mut out = Vec::new();
    let gated: &[(bool, &[&str], &str)] = &[
        (policy.http, HTTP_FNS, "--no-http"),
        (policy.websocket, WEBSOCKET_FNS, "--no-websocket"),
        (policy.fs, FS_FNS, "--no-fs"),
        (policy.process, PROCESS_FNS, "--sandbox"),
    ];
    for (allowed, fns, flag) in gated {
        if *allowed {
            continue;
        }
        for f in fns.iter().filter(|f| idents.contains(**f)) {
            out.push(format!("security: `{f}` is disabled by {flag}"));
        }
    }
    if !policy.ffi {
        extern_violations(&program.statements, &mut out);
    }
    out
}

/// `--no-ffi`: any `extern` declaration (including inside modules) is rejected.
fn extern_violations(statements: &[Stmt], out: &mut Vec<String>) {
    for s in statements {
        match s {
            Stmt::Extern { name, .. } => out.push(format!(
                "security: extern function `{name}` is disabled by --no-ffi"
            )),
            Stmt::Module { body, .. } => extern_violations(body, out),
            _ => {}
        }
    }
}
