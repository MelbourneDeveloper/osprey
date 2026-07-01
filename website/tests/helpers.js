// Shared fixtures/data for the Osprey website e2e tests.

// Every distinct page template the site renders.
const PAGES = [
  { name: "home", path: "/", kind: "home" },
  { name: "spec-index", path: "/spec/", kind: "prose" },
  { name: "spec-prose", path: "/spec/0001-introduction/", kind: "prose" },
  { name: "docs-index", path: "/docs/", kind: "prose" },
  { name: "docs-function", path: "/docs/functions/map/", kind: "prose" },
  { name: "docs-type", path: "/docs/types/string/", kind: "prose" },
  { name: "blog-index", path: "/blog/", kind: "listing" },
  { name: "blog-post", path: "/blog/2026-05-17-persistent-collections/", kind: "prose" },
  { name: "status", path: "/status/", kind: "prose" },
  { name: "wasm", path: "/wasm/", kind: "wasm" },
  { name: "playground", path: "/playground/", kind: "app" },
];

const VIEWPORTS = [
  { name: "desktop", width: 1440, height: 900 },
  { name: "tablet", width: 768, height: 1024 },
  { name: "mobile", width: 390, height: 844 },
];

// Per-page error suppressions. Empty: the playground's "totalResults is not
// defined" crash (an unescaped ${} in the embedded sample) is fixed at the
// source (website/scripts/update-playground.js now escapes ${), so the suite
// asserts every page — playground included — is genuinely error-free.
const KNOWN_ERRORS = [];

// Third-party/CDN request failures we don't control (playground's Monaco editor).
const IGNORED_REQUEST_HOSTS = /monaco|cdnjs|jsdelivr|unpkg|googleapis|gstatic/;

// Attaches console-error / pageerror / failed-request collectors to a page.
function collectProblems(page, pageName) {
  const consoleErrors = [];
  const pageErrors = [];
  const failedRequests = [];
  page.on("console", (m) => {
    if (m.type() === "error" && !isKnown(pageName, m.text())) consoleErrors.push(m.text());
  });
  page.on("pageerror", (e) => {
    if (!isKnown(pageName, String(e))) pageErrors.push(String(e));
  });
  page.on("requestfailed", (r) => {
    if (!IGNORED_REQUEST_HOSTS.test(r.url())) failedRequests.push(`${r.url()} (${r.failure()?.errorText})`);
  });
  page.on("response", (r) => {
    if (r.status() >= 400 && !IGNORED_REQUEST_HOSTS.test(r.url())) failedRequests.push(`${r.status()} ${r.url()}`);
  });
  return { consoleErrors, pageErrors, failedRequests };
}

function isKnown(pageName, text) {
  return KNOWN_ERRORS.some((k) => k.page === pageName && k.pattern.test(text));
}

module.exports = { PAGES, VIEWPORTS, collectProblems };
