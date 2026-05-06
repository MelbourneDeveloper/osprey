const { defineConfig } = require('@vscode/test-cli');

module.exports = defineConfig({
  files: 'out/test/suite/**/*.test.js',
  srcDir: 'client/src',
  version: 'stable',
  mocha: {
    ui: 'tdd',
    timeout: 10000,
    color: true
  },
  launchArgs: [
    '--disable-extensions',
    '--disable-workspace-trust'
  ],
  coverage: {
    reporter: ['text-summary', 'json-summary', 'html'],
    include: ['out/client/src/**/*.js', 'out/server/src/**/*.js'],
    exclude: ['out/test/**', '**/node_modules/**'],
    includeAll: true
  }
});
