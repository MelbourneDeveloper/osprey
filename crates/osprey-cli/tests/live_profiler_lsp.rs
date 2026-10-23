//! The profiler through the real language server [PROF-LSP-HOST]: `osprey lsp`
//! is spawned over stdio, as the editor spawns it, and a running Osprey program
//! is listed, attached to and profiled with `osprey.lspProfiler.*` commands.
// The endpoint and its Unix socket are POSIX-only [PROF-ATTACH-ENDPOINT].
#![cfg(unix)]
#![expect(
    clippy::expect_used,
    clippy::panic,
    reason = "a failed setup step must fail the end-to-end test that needed it"
)]

use std::io::{BufRead, BufReader, Read, Write};
use std::process::{Child, ChildStdin, ChildStdout, Command, Stdio};
use std::time::Duration;

use serde_json::{json, Value};

#[path = "cases/live_fixture.rs"]
mod live_fixture;

use live_fixture::{Fixture, HOT_FUNCTION};

/// `osprey lsp` with its stdio, killed on drop.
struct Server {
    child: Child,
    stdin: ChildStdin,
    stdout: BufReader<ChildStdout>,
}

impl Server {
    /// The server sees the fixture's `TMPDIR`, so it finds its endpoint where
    /// every program's endpoint lives: `$TMPDIR/osprey-diagnostics-<uid>`.
    fn start(fixture: &Fixture) -> Self {
        let mut child = Command::new(env!("CARGO_BIN_EXE_osprey"))
            .arg("lsp")
            .env("TMPDIR", &fixture.tmp)
            .stdin(Stdio::piped())
            .stdout(Stdio::piped())
            .stderr(Stdio::null())
            .spawn()
            .expect("osprey lsp starts");
        let stdin = child.stdin.take().expect("stdin");
        let stdout = BufReader::new(child.stdout.take().expect("stdout"));
        Self {
            child,
            stdin,
            stdout,
        }
    }

    fn send(&mut self, message: &Value) {
        let body = message.to_string();
        write!(self.stdin, "Content-Length: {}\r\n\r\n{body}", body.len()).expect("write");
        self.stdin.flush().expect("flush");
    }

    fn read(&mut self) -> Value {
        let mut length = 0;
        loop {
            let mut header = String::new();
            assert!(
                self.stdout.read_line(&mut header).expect("header") > 0,
                "the server closed its stdout"
            );
            match header.trim_end().strip_prefix("Content-Length: ") {
                Some(n) => length = n.parse().expect("length"),
                None if header.trim_end().is_empty() => break,
                None => {}
            }
        }
        let mut body = vec![0; length];
        self.stdout.read_exact(&mut body).expect("body");
        serde_json::from_slice(&body).expect("json")
    }

    /// Send a request and return its response, skipping notifications.
    fn request(&mut self, id: i64, method: &str, params: &Value) -> Value {
        self.send(&json!({ "jsonrpc": "2.0", "id": id, "method": method, "params": params }));
        loop {
            let message = self.read();
            if message.get("id") == Some(&json!(id)) {
                return message;
            }
        }
    }

    fn command(&mut self, id: i64, name: &str, params: &Value) -> Value {
        let command =
            json!({ "command": format!("osprey.lspProfiler.{name}"), "arguments": [params] });
        let reply = self.request(id, "workspace/executeCommand", &command);
        reply
            .get("result")
            .cloned()
            .unwrap_or_else(|| panic!("{name} failed: {reply}"))
    }
}

impl Drop for Server {
    fn drop(&mut self) {
        drop(self.child.kill());
        drop(self.child.wait());
    }
}

#[test]
fn the_language_server_lists_and_profiles_a_running_program() {
    let fixture = Fixture::start("lsp-cpu");
    let socket = fixture.socket_dir.join(format!("{}.sock", fixture.pid));
    assert!(socket.exists(), "the program serves {}", socket.display());
    let mut server = Server::start(&fixture);
    let init = server.request(1, "initialize", &json!({ "capabilities": {} }));
    assert!(
        init.pointer("/result/capabilities/experimental/lspProfiler")
            .is_some(),
        "{init}"
    );
    let list = server.command(2, "processes", &json!({ "workspace_roots": [] }));
    let listed = list
        .get("processes")
        .and_then(Value::as_array)
        .and_then(|all| {
            all.iter()
                .find(|p| p.get("pid") == Some(&json!(fixture.pid)))
        })
        .unwrap_or_else(|| panic!("PID {} missing from {list}", fixture.pid));
    assert_eq!(listed.get("runtime"), Some(&json!("osprey")), "{listed}");
    let started = server.command(
        3,
        "cpu.start",
        &json!({ "pid": fixture.pid, "sample_rate": 997 }),
    );
    std::thread::sleep(Duration::from_secs(2));
    let report = server.command(
        4,
        "cpu.stop",
        &json!({ "session_id": started.get("session_id") }),
    );
    let hot = report
        .get("hot_functions")
        .and_then(Value::as_array)
        .map_or(&[][..], Vec::as_slice);
    assert!(
        hot.iter()
            .any(|f| f.get("name") == Some(&json!(HOT_FUNCTION))),
        "{report}"
    );
    let shutdown = server.request(5, "shutdown", &Value::Null);
    assert!(shutdown.get("error").is_none(), "{shutdown}");
    server.send(&json!({ "jsonrpc": "2.0", "method": "exit" }));
}
