//! Live-attach CPU profiling of a real Osprey program through lspkit's shared
//! profiler [PROF-ATTACH-ENDPOINT] [PROF-ATTACH-PROTOCOL] [PROF-LSP-RUNTIME].
//!
//! No mocks: this build's `osprey` compiles the fixture with debug info, the
//! program runs as its own process with its own endpoint, and the tests drive
//! `ProfilerService` over its wire methods exactly as the language server
//! does — the same contract lspkit's .NET suite pins for its runtime.
// The endpoint and its Unix socket are POSIX-only [PROF-ATTACH-ENDPOINT].
#![cfg(unix)]
#![expect(
    clippy::expect_used,
    clippy::panic,
    reason = "a failed setup step must fail the end-to-end test that needed it"
)]

use std::path::Path;
use std::sync::Arc;
use std::time::Duration;

use async_trait::async_trait;
use lsp_profiler::wire::methods;
use lsp_profiler::{ProfilerConfig, ProfilerEvents, ProfilerService};
use lsp_profiler_osprey::{OspreyConfig, OspreyRuntime};
use serde_json::{json, Value};

#[path = "cases/live_fixture.rs"]
mod live_fixture;

use live_fixture::{Fixture, HOT_FUNCTION};

/// `spin` in the fixture: its signature through its last arm.
const SPIN_LINES: std::ops::RangeInclusive<i64> = 9..=13;

#[derive(Debug)]
struct Silent;

#[async_trait]
impl ProfilerEvents for Silent {
    async fn notify(&self, _method: &'static str, _params: Value) {}
}

/// The shared service with only the Osprey runtime, reading `fixture`'s
/// endpoint directory and writing exports under its scratch directory.
fn service(fixture: &Fixture) -> ProfilerService {
    let runtime = OspreyRuntime::new(OspreyConfig {
        socket_dir: Some(fixture.socket_dir.clone()),
    });
    let config = ProfilerConfig {
        output_dir: fixture.tmp.join("out"),
        ..ProfilerConfig::default()
    };
    ProfilerService::new(vec![Arc::new(runtime)], config, Arc::new(Silent))
}

async fn call(service: &ProfilerService, method: &str, params: Value) -> Value {
    service
        .handle(method, params)
        .await
        .unwrap_or_else(|error| panic!("{method} failed: {error}"))
}

async fn error_code(service: &ProfilerService, method: &str, params: Value) -> i64 {
    match service.handle(method, params).await {
        Ok(value) => panic!("{method} must be refused, got {value}"),
        Err(error) => error.code(),
    }
}

/// Record `secs` of the fixture, then stop and return the report.
async fn profile(service: &ProfilerService, fixture: &Fixture, params: Value, secs: u64) -> Value {
    let mut start = json!({ "pid": fixture.pid, "sample_rate": 997 });
    if let (Some(start), Some(extra)) = (start.as_object_mut(), params.as_object()) {
        start.extend(extra.clone());
    }
    let started = call(service, methods::CPU_START, start).await;
    tokio::time::sleep(Duration::from_secs(secs)).await;
    call(
        service,
        methods::CPU_STOP,
        json!({ "session_id": started.get("session_id") }),
    )
    .await
}

fn hot_functions(report: &Value) -> &[Value] {
    report
        .get("hot_functions")
        .and_then(Value::as_array)
        .map_or(&[], Vec::as_slice)
}

fn samples(report: &Value) -> i64 {
    report
        .get("total_samples")
        .and_then(Value::as_i64)
        .unwrap_or(0)
}

#[tokio::test(flavor = "multi_thread")]
async fn the_program_is_listed_as_an_attachable_osprey_process() {
    let fixture = Fixture::start("lp-list");
    let service = service(&fixture);
    let list = call(
        &service,
        methods::PROCESSES,
        json!({ "workspace_roots": [] }),
    )
    .await;
    let entry = list
        .get("processes")
        .and_then(Value::as_array)
        .and_then(|all| {
            all.iter()
                .find(|p| p.get("pid") == Some(&json!(fixture.pid)))
        })
        .unwrap_or_else(|| panic!("PID {} missing from {list}", fixture.pid));
    assert_eq!(entry.get("runtime"), Some(&json!("osprey")), "{entry}");
    assert_eq!(entry.get("attachable"), Some(&json!(true)), "{entry}");
    assert_eq!(entry.get("entry"), Some(&json!("live_attach")), "{entry}");
    let version = entry
        .get("runtime_version")
        .and_then(Value::as_str)
        .unwrap_or_default();
    assert!(!version.is_empty(), "{entry}");
    let started = call(&service, methods::CPU_START, json!({ "pid": fixture.pid })).await;
    assert_eq!(
        started.get("runtime_version").and_then(Value::as_str),
        Some(version)
    );
    drop(
        call(
            &service,
            methods::CPU_STOP,
            json!({ "session_id": started.get("session_id") }),
        )
        .await,
    );
    service.shutdown().await;
}

#[tokio::test(flavor = "multi_thread")]
async fn cpu_profiling_finds_the_hot_function_and_its_lines() {
    let fixture = Fixture::start("lp-cpu");
    let service = service(&fixture);
    let report = profile(&service, &fixture, json!({}), 3).await;
    assert!(samples(&report) > 100, "{report}");
    let spin = hot_functions(&report)
        .iter()
        .find(|f| f.get("name") == Some(&json!(HOT_FUNCTION)))
        .unwrap_or_else(|| panic!("{HOT_FUNCTION} absent: {report}"));
    let file = spin.get("file").and_then(Value::as_str).unwrap_or_default();
    assert!(file.ends_with("live_attach.osp"), "{spin}");
    let share = spin
        .get("percentage")
        .and_then(Value::as_f64)
        .unwrap_or(0.0);
    assert!(
        share > 50.0,
        "{HOT_FUNCTION} owns the busy main thread, {share}%: {report}"
    );
    let lines = report
        .get("hot_lines")
        .and_then(Value::as_array)
        .map_or(&[][..], Vec::as_slice);
    assert!(
        lines
            .iter()
            .any(|l| l.get("file").and_then(Value::as_str) == Some(file)
                && SPIN_LINES.contains(&l.get("line").and_then(Value::as_i64).unwrap_or(0))),
        "a hot line inside {HOT_FUNCTION}: {report}"
    );
    for key in ["speedscope_path", "flamegraph_path", "cpuprofile_path"] {
        let path = report
            .pointer(&format!("/artifacts/{key}"))
            .and_then(Value::as_str);
        assert!(
            path.is_some_and(|p| Path::new(p).is_file()),
            "{key} not written: {report}"
        );
    }
    service.shutdown().await;
}

#[tokio::test(flavor = "multi_thread")]
async fn native_frames_appear_only_on_request() {
    let fixture = Fixture::start("lp-native");
    let service = service(&fixture);
    let native = |report: &Value| {
        hot_functions(report)
            .iter()
            .any(|f| f.get("file").is_none())
    };
    let source_only = profile(&service, &fixture, json!({}), 2).await;
    assert!(
        !native(&source_only),
        "without native frames every hotspot is source: {source_only}"
    );
    let with_native = profile(&service, &fixture, json!({ "include_native": true }), 2).await;
    assert!(
        native(&with_native),
        "the runtime's own frames appear on request: {with_native}"
    );
    service.shutdown().await;
}

/// [PROF-ERRORS]: a pid that names no process is `PROCESS_NOT_FOUND`; a live
/// process without an endpoint (this test binary) is `UNSUPPORTED_RUNTIME`.
#[tokio::test(flavor = "multi_thread")]
async fn a_pid_without_an_osprey_endpoint_is_classified() {
    let fixture = Fixture::start("lp-errors");
    let service = service(&fixture);
    let missing = json!({ "pid": 999_999 });
    assert_eq!(
        error_code(&service, methods::CPU_START, missing).await,
        -32001
    );
    let own = json!({ "pid": std::process::id() });
    assert_eq!(error_code(&service, methods::CPU_START, own).await, -32002);
    service.shutdown().await;
}

/// [PROF-CPU-START]: a pid with a session is `ALREADY_PROFILING`.
#[tokio::test(flavor = "multi_thread")]
async fn a_second_session_on_the_same_program_is_refused() {
    let fixture = Fixture::start("lp-busy");
    let service = service(&fixture);
    let started = call(&service, methods::CPU_START, json!({ "pid": fixture.pid })).await;
    assert_eq!(
        error_code(&service, methods::CPU_START, json!({ "pid": fixture.pid })).await,
        -32004
    );
    drop(
        call(
            &service,
            methods::CPU_STOP,
            json!({ "session_id": started.get("session_id") }),
        )
        .await,
    );
    service.shutdown().await;
}

/// [PROF-CPU-SESSIONS]: a duration ends the session by itself.
#[tokio::test(flavor = "multi_thread")]
async fn a_duration_caps_the_recorded_time() {
    let fixture = Fixture::start("lp-duration");
    let service = service(&fixture);
    let report = profile(&service, &fixture, json!({ "duration_secs": 1.0 }), 3).await;
    let secs = report
        .get("duration_secs")
        .and_then(Value::as_f64)
        .unwrap_or(0.0);
    assert!(
        (0.8..=1.5).contains(&secs),
        "recorded about the 1 s requested: {report}"
    );
    assert!(samples(&report) > 0, "{report}");
    service.shutdown().await;
}

/// [PROF-CPU-SNAPSHOT]: samples stream while recording, so a snapshot has them.
#[tokio::test(flavor = "multi_thread")]
async fn a_mid_session_snapshot_reports_the_samples_so_far() {
    let fixture = Fixture::start("lp-snapshot");
    let service = service(&fixture);
    let started = call(&service, methods::CPU_START, json!({ "pid": fixture.pid })).await;
    let session = json!({ "session_id": started.get("session_id") });
    tokio::time::sleep(Duration::from_secs(2)).await;
    let snapshot = call(&service, methods::CPU_SNAPSHOT, session.clone()).await;
    drop(call(&service, methods::CPU_STOP, session).await);
    assert!(
        samples(&snapshot) > 0,
        "two seconds in, the snapshot has samples: {snapshot}"
    );
    service.shutdown().await;
}

/// [PROF-CPU-STOP]: a program killed mid-session keeps the samples taken while
/// it lived.
#[tokio::test(flavor = "multi_thread")]
async fn a_program_killed_mid_session_keeps_its_samples() {
    let mut fixture = Fixture::start("lp-killed");
    let service = service(&fixture);
    let started = call(&service, methods::CPU_START, json!({ "pid": fixture.pid })).await;
    tokio::time::sleep(Duration::from_secs(2)).await;
    fixture.kill();
    tokio::time::sleep(Duration::from_millis(500)).await;
    let report = call(
        &service,
        methods::CPU_STOP,
        json!({ "session_id": started.get("session_id") }),
    )
    .await;
    assert!(
        samples(&report) > 0,
        "samples taken before the kill are reported: {report}"
    );
    service.shutdown().await;
}
