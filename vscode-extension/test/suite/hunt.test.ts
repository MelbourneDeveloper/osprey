// Tests for turning an `osprey hunt --json` report into editor diagnostics
// ([SYMEX-CLI], [SYMEX-FINDINGS]). The engine is KLEE; the editor only reads
// the report.

import * as assert from "assert";
import * as path from "path";
import { DiagnosticSeverity } from "vscode";
import {
  HUNT_SOURCE,
  HuntReport,
  findingMessage,
  findingPath,
  findingSeverity,
  huntDiagnostics,
  huntSummary,
  parseHuntReport,
} from "../../client/src/hunt";

const program = path.resolve("/work/app/main.osp");

function report(overrides: Partial<HuntReport> = {}): HuntReport {
  return {
    program,
    engine: "KLEE",
    budgetSeconds: 60,
    explored: { paths: 12, completed: 10, instructions: 5000, elapsed: "00:00:03" },
    findings: [
      {
        class: "memory",
        message: "nullptr-with-offset",
        file: "main.osp",
        line: 7,
        runtimeFrame: "runtime/json_runtime.c:237 in close_container",
        stdin: "9\\n",
        replay: "reproduced",
        replayDetail: "SIGSEGV",
      },
      {
        class: "unmodelled",
        message: "unmodelled: pipe",
        file: null,
        line: null,
        runtimeFrame: null,
        stdin: "",
        replay: "skipped",
        replayDetail: "an unmodelled call stops the path under KLEE only",
      },
    ],
    incomplete: [],
    ...overrides,
  };
}

suite("hunt report", () => {
  test("parses the JSON report and rejects anything else", () => {
    assert.strictEqual(parseHuntReport(JSON.stringify(report()))?.engine, "KLEE");
    assert.strictEqual(parseHuntReport("not json"), undefined);
    assert.strictEqual(parseHuntReport("{}"), undefined);
  });

  test("summarises findings, replay and exploration in one line", () => {
    assert.strictEqual(
      huntSummary(report()),
      "Hunt: 2 findings (1 reproduced natively), 12 paths (10 complete) explored by KLEE.",
    );
    assert.strictEqual(
      huntSummary(report({ findings: [] })),
      "Hunt found no bugs: 12 paths (10 complete) explored by KLEE.",
    );
  });

  test("a reproduced finding is an error, klee-only a warning, unmodelled information", () => {
    const [memory, unmodelled] = report().findings;
    assert.strictEqual(findingSeverity(memory), DiagnosticSeverity.Error);
    assert.strictEqual(findingSeverity({ ...memory, replay: "klee-only" }), DiagnosticSeverity.Warning);
    assert.strictEqual(findingSeverity(unmodelled), DiagnosticSeverity.Information);
  });

  test("a relative file resolves beside the program, an absolute one stands", () => {
    const [memory] = report().findings;
    assert.strictEqual(findingPath(memory, program), path.resolve("/work/app/main.osp"));
    assert.strictEqual(findingPath({ ...memory, file: "/abs/other.osp" }, program), "/abs/other.osp");
    assert.strictEqual(findingPath({ ...memory, file: null }, program), program);
  });

  test("diagnostics sit on the finding's line and name the engine", () => {
    const byFile = huntDiagnostics(report(), program);
    const diagnostics = byFile.get(program) ?? [];
    assert.strictEqual(diagnostics.length, 2);
    assert.strictEqual(diagnostics[0].range.start.line, 6);
    assert.strictEqual(diagnostics[0].source, HUNT_SOURCE);
    assert.strictEqual(diagnostics[0].code, "memory");
    assert.strictEqual(diagnostics[1].range.start.line, 0, "a finding without a line lands on line 1");
    assert.ok(findingMessage(report().findings[0]).includes("replay reproduced (SIGSEGV)"));
    assert.ok(findingMessage(report().findings[0]).includes('stdin "9\\n"'));
  });
});
