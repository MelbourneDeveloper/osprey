import * as assert from 'assert';
import { runOspreyCompiler, ExecFileFn } from '../../server/src/compilerInvocation';

// Regression test for issue #129: the extension surfaced a false
// "undefined variable 'listAppend'" diagnostic because diagnostics always ran a
// stale `osprey` binary from PATH instead of the compiler the user configured
// via `osprey.server.compilerPath`. runOspreyCompiler must invoke the configured
// path so users can point diagnostics at a freshly-built compiler.
suite('runOspreyCompiler compiler-path resolution', () => {
  test('invokes the configured compilerPath, not a hardcoded binary name', async () => {
    let invokedBinary: string | undefined;
    const fakeExec: ExecFileFn = (file, _args, callback) => {
      invokedBinary = file;
      callback(null, '', '');
    };

    await runOspreyCompiler('/tmp/example.osp', '/custom/bin/osprey', fakeExec);

    assert.strictEqual(
      invokedBinary,
      '/custom/bin/osprey',
      'diagnostics must run the configured compiler; otherwise a stale PATH binary produces false errors (issue #129)'
    );
  });
});
