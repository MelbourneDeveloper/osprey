// The singleton flame-graph WebviewPanel ([PROF-VSCODE-FLAME]): thin vscode
// shell around flame-html. Click-to-source messages from the canvas/table are
// parsed by the pure `parseSelectMessage` and revealed via `revealSource`.

import { bindPanelPalette } from "@nimblesite/lsp-profiler-ui/out/src/palette-vscode";
import { bindDebugSelection } from "@nimblesite/lsp-profiler-ui/out/src/debug-selection";
import { handleSourceNavigation, revealProfileSource as revealSource } from "@nimblesite/lsp-profiler-ui/out/src/profiler-webview";
import * as path from "path";
import * as vscode from "vscode";
import { buildFlamegraphHtml, openNativeTraceViewer, openSpeedscopeImport } from "@nimblesite/lsp-profiler-ui/out/src/profiler-flamegraph-html";
import { installIds } from "@nimblesite/lsp-profiler-ui/out/src/ids";
import { configurePalette, resolvePalette } from "@nimblesite/lsp-profiler-ui/out/src/palette";
import { nativeProfileResult } from "./native-profile";

import { OSPREY_PALETTE } from "./palette";
installIds({ idPrefix: "osprey", brand: "Osprey", configSection: "osprey", debugType: "osprey" });
configurePalette(() => resolvePalette(OSPREY_PALETTE,
  vscode.workspace.getConfiguration("osprey").get("profiler.palette")));
import type { FlameModel } from "./flame-model";
import type { ProfileSummary } from "./summary";

export interface SelectMessage {
  type: "select";
  file: string;
  line: number;
}

/** Validate a webview message; only well-formed selects navigate. Pure. */
export function parseSelectMessage(
  message: unknown,
): SelectMessage | undefined {
  const select = message as SelectMessage;
  const valid =
    typeof select === "object" &&
    select !== null &&
    select.type === "select" &&
    typeof select.file === "string" &&
    select.file !== "" &&
    typeof select.line === "number";
  return valid ? select : undefined;
}

export { revealSource };

/** The webview message handler (exported so it is directly testable). */
export async function handlePanelMessage(message: unknown): Promise<void> {
  if (typeof message === "object" && message !== null && handleSourceNavigation(message as { type: string })) return;
  const candidate = message as { type?: string; file?: string } | undefined;
  if (candidate?.type === "openCpuProfile" && candidate.file) {
    await openNativeTraceViewer(candidate.file, candidate.file); return;
  }
  if (candidate?.type === "openSpeedscope" && candidate.file) {
    await openSpeedscopeImport(candidate.file); return;
  }
  const select = parseSelectMessage(candidate?.type === "navigateToSource"
    ? { ...candidate, type: "select" } : message);
  if (select === undefined) {
    return;
  }
  try {
    await revealSource(select);
  } catch {
    void vscode.window.showWarningMessage(`Cannot open ${select.file}`);
  }
}

let panel: vscode.WebviewPanel | undefined;

function createPanel(title: string): vscode.WebviewPanel {
  const created = vscode.window.createWebviewPanel(
    "ospreyProfiler",
    title,
    vscode.ViewColumn.Beside,
    { enableScripts: true },
  );
  bindDebugSelection(created);
  bindPanelPalette(created);
  created.onDidDispose(() => {
    panel = undefined;
  });
  created.webview.onDidReceiveMessage(
    (message) => void handlePanelMessage(message),
  );
  return created;
}

/** Show (or refresh) the singleton profiler panel for a new run. */
export function showFlamePanel(
  model: FlameModel,
  summary: ProfileSummary,
  sourcePath: string,
  artifacts?: { speedscope: string; cpuprofile: string },
): vscode.WebviewPanel {
  const title = `Profile: ${path.basename(sourcePath)}`;
  if (panel === undefined) {
    panel = createPanel(title);
  } else {
    panel.title = title;
    panel.reveal(vscode.ViewColumn.Beside, true);
  }
  panel.webview.html = buildFlamegraphHtml(nativeProfileResult(model, summary, artifacts));
  return panel;
}
