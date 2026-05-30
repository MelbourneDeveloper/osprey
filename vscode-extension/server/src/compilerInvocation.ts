import { execFile } from 'child_process';

// Signature of Node's execFile that we actually rely on. Kept narrow so it can
// be substituted in tests.
export type ExecFileFn = (
  file: string,
  args: string[],
  callback: (error: Error | null, stdout: string, stderr: string) => void
) => void;

export interface CompilerResult {
  stdout: string;
  stderr: string;
  error?: Error;
}

// Runs the Osprey compiler against a file to gather diagnostics.
//
// The `compilerPath` MUST be honored: users point `osprey.server.compilerPath`
// at a freshly-built compiler so the extension does not analyze their code with
// a stale binary on PATH (which yields false "undefined variable" errors for
// builtins the old compiler predates). See issue #129.
export function runOspreyCompiler(
  filePath: string,
  compilerPath: string,
  exec: ExecFileFn = execFile as unknown as ExecFileFn
): Promise<CompilerResult> {
  return new Promise((resolve) => {
    exec(compilerPath, [filePath], (_error, stdout, stderr) => {
      // Don't treat non-zero exit codes as errors - they might just be syntax errors
      resolve({ stdout, stderr, error: undefined });
    });
  });
}
