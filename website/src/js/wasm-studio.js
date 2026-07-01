import { runModule } from "/wasm/wasi-shim.mjs";

const $ = (id) => document.getElementById(id);
const esc = (s) => String(s).replace(/[&<>]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;" }[c]));
const SQLJS = "https://cdn.jsdelivr.net/npm/sql.js@1.11.0/dist/";

const OSP_RE = /(\/\/[^\n]*)|("(?:[^"\\]|\\.)*")|\b(fn|let|mut|type|match|effect|handle|in|perform|resume|import|return|true|false)\b|\b([A-Z][A-Za-z0-9_]*)\b|\b(\d+)\b|(:=|=>|->|\|>|>=|<=|==|!=|[=+\-*/%<>!])/g;
const SQL_RE = /(--[^\n]*)|('(?:[^'\\]|\\.)*'|"(?:[^"\\]|\\.)*")|\b(SELECT|FROM|WHERE|GROUP|BY|ORDER|DESC|ASC|LIMIT|AS|SUM|COUNT|AVG|MIN|MAX|ROUND|CASE|WHEN|THEN|ELSE|END|CREATE|TABLE|INSERT|INTO|VALUES|INTEGER|TEXT|PRIMARY|KEY|NOT|NULL|AND|OR|ON|JOIN|DISTINCT|HAVING)\b/gi;
const CLS = ["tok-c", "tok-s", "tok-k", "tok-t", "tok-n", "tok-o"];

const FLAVORS = {
  osp: { file: "/wasm/studio.osp", wasm: "/wasm/build/studio.osp.wasm", label: "Default .osp" },
  ospml: { file: "/wasm/studio.ospml", wasm: "/wasm/build/studio.ospml.wasm", label: "ML .ospml" },
};

const METRIC_CARDS = [
  { key: "total_revenue", lbl: "Total revenue", money: true, cls: "is-good", sql: "SELECT SUM(qty*price) v FROM sales" },
  { key: "total_units", lbl: "Units sold", sql: "SELECT SUM(qty) v FROM sales" },
  { key: "order_count", lbl: "Orders", sql: "SELECT COUNT(*) v FROM sales" },
  { key: "premium_orders", lbl: "Premium orders", cls: "is-good", sql: "SELECT COUNT(*) v FROM sales WHERE price >= 15" },
  { key: "revenue_tier", lbl: "Book tier", tier: true },
];

const state = { flavor: "osp", manifest: null, SQL: null, db: null, srcView: "osp", sources: {} };

function highlight(src, lang) {
  const re = lang === "sql" ? SQL_RE : OSP_RE;
  re.lastIndex = 0;
  let out = "";
  let last = 0;
  let m;
  while ((m = re.exec(src))) {
    out += esc(src.slice(last, m.index));
    let cls = "tok-k";
    for (let g = 1; g < m.length; g++) {
      if (m[g] !== undefined) {
        cls = CLS[g - 1] || "tok-k";
        break;
      }
    }
    out += `<span class="${cls}">${esc(m[0])}</span>`;
    last = re.lastIndex;
    if (m.index === re.lastIndex) re.lastIndex++;
  }
  return out + esc(src.slice(last));
}

function parseManifest(text) {
  const out = { ddl: [], seed: [], queries: [], metrics: {}, raw: text };
  let mode = null;
  let cur = null;
  for (const line of text.split("\n")) {
    if (line.startsWith("::")) {
      const [tag, a, b] = line.slice(2).split("|");
      if (tag === "DDL") mode = "ddl";
      else if (tag === "SEED") mode = "seed";
      else if (tag === "END-SEED") mode = null;
      else if (tag === "QUERY") {
        cur = { name: a, title: b, sql: [] };
        out.queries.push(cur);
        mode = "query";
      } else if (tag === "METRIC") {
        out.metrics[a] = b;
        mode = null;
      } else {
        mode = null;
      }
      continue;
    }
    if (mode === "ddl") out.ddl.push(line);
    else if (mode === "seed") out.seed.push(line);
    else if (mode === "query" && cur) cur.sql.push(line);
  }
  out.ddl = out.ddl.join("\n").trim();
  out.seed = out.seed.join("\n").trim();
  out.queries.forEach((q) => {
    q.sql = q.sql.join("\n").trim();
  });
  return out;
}

async function runOsprey(flavor) {
  const res = await fetch(FLAVORS[flavor].wasm);
  if (!res.ok) throw new Error(`fetch ${FLAVORS[flavor].wasm}: HTTP ${res.status}`);
  const bytes = await res.arrayBuffer();
  let text = "";
  await runModule(bytes, (t) => {
    text += t;
  });
  $("flavor-bytes").textContent = `${(bytes.byteLength / 1024).toFixed(1)} KB wasm -> ${text.length} B manifest`;
  return { manifest: parseManifest(text), size: bytes.byteLength };
}

function loadScript(src) {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[src="${src}"]`);
    if (existing) {
      resolve();
      return;
    }
    const script = document.createElement("script");
    script.src = src;
    script.onload = resolve;
    script.onerror = () => reject(new Error("could not load " + src));
    document.head.appendChild(script);
  });
}

async function initSqlite() {
  await loadScript(SQLJS + "sql-wasm.js");
  return window.initSqlJs({ locateFile: (file) => SQLJS + file });
}

function buildDb() {
  if (!state.SQL || !state.manifest) return;
  if (state.db) state.db.close();
  state.db = new state.SQL.Database();
  state.db.run(state.manifest.ddl);
  state.db.run(state.manifest.seed);
}

function query(sql) {
  const result = state.db.exec(sql);
  return result.length ? result[0] : { columns: [], values: [] };
}

function scalar(sql) {
  const result = query(sql);
  return result.values.length ? result.values[0][0] : null;
}

function fmtMoney(n) {
  return "$" + Number(n).toLocaleString();
}

function tableHTML(result) {
  if (!result.columns.length) return `<p class="muted">No rows.</p>`;
  const head = result.columns.map((c) => `<th>${esc(c)}</th>`).join("");
  const isNum = result.columns.map((_, i) => result.values.every((row) => typeof row[i] === "number"));
  const body = result.values
    .map(
      (row) =>
        "<tr>" +
        row.map((v, i) => `<td class="${isNum[i] ? "num" : ""}">${esc(v === null ? "null" : v)}</td>`).join("") +
        "</tr>"
    )
    .join("");
  return `<table class="data"><thead><tr>${head}</tr></thead><tbody>${body}</tbody></table>`;
}

function renderMetrics() {
  const m = state.manifest.metrics;
  $("metrics").innerHTML = METRIC_CARDS.map((card) => {
    if (card.tier) {
      return `<div class="metric"><span class="lbl">${esc(card.lbl)}</span><span class="badge-tier">${esc(
        m[card.key] || "-"
      )}</span><span class="check" id="chk-${card.key}">&nbsp;</span></div>`;
    }
    const value = m[card.key];
    const big = card.money ? fmtMoney(value) : Number(value).toLocaleString();
    return `<div class="metric ${card.cls || ""}"><span class="big">${esc(big)}</span><span class="lbl">${esc(
      card.lbl
    )}</span><span class="check" id="chk-${card.key}">checking</span></div>`;
  }).join("");
}

function reconcile() {
  const recEl = $("reconcile");
  if (!state.db) {
    recEl.className = "reconcile";
    recEl.innerHTML = `<span class="dot"></span><span>SQLite offline - showing Osprey answers only.</span>`;
    METRIC_CARDS.forEach((card) => {
      const el = $("chk-" + card.key);
      if (el) el.innerHTML = "";
    });
    $("pn-reconcile").textContent = "Osprey only";
    return;
  }

  let allOk = true;
  for (const card of METRIC_CARDS) {
    if (!card.sql) continue;
    const got = Number(scalar(card.sql));
    const want = Number(state.manifest.metrics[card.key]);
    const ok = got === want;
    allOk = allOk && ok;
    const el = $("chk-" + card.key);
    if (el) {
      el.className = "check " + (ok ? "ok" : "bad");
      el.textContent = ok ? `SQLite agrees (${got})` : `SQLite got ${got}`;
    }
  }
  recEl.className = "reconcile " + (allOk ? "ok" : "bad");
  recEl.innerHTML = `<span class="dot"></span><span>${
    allOk
      ? "Every Osprey metric matches SQLite."
      : "Mismatch detected between Osprey and SQLite."
  }</span>`;
  $("pn-reconcile").textContent = allOk ? "agree" : "differ";
}

function renderQueries() {
  $("queries").innerHTML = state.manifest.queries
    .map((q) => {
      let result = "";
      if (state.db) {
        try {
          result = tableHTML(query(q.sql));
        } catch (err) {
          result = `<p class="muted">SQLite error: ${esc(err.message || err)}</p>`;
        }
      } else {
        result = `<p class="muted">SQLite is offline. The query still ships in the Osprey manifest.</p>`;
      }
      return `<article class="wasm-query">
        <div class="qhead"><h3>${esc(q.title)}</h3><span class="qsql">${esc(q.name)}</span></div>
        <div class="tablewrap">${result}</div>
      </article>`;
    })
    .join("");
}

function renderManifest() {
  $("manifest").innerHTML = highlight(state.manifest.raw.replace(/\n+$/, ""), "sql");
}

function renderPresets() {
  const presets = [
    ...state.manifest.queries.map((q) => ({ label: q.name, sql: q.sql })),
    { label: "all rows", sql: "SELECT * FROM sales ORDER BY id;" },
    {
      label: "best margin",
      sql: "SELECT product, region, qty, price, qty*price AS revenue\nFROM sales ORDER BY revenue DESC LIMIT 3;",
    },
  ];
  $("presets").innerHTML = presets.map((p, i) => `<button class="preset" data-i="${i}">${esc(p.label)}</button>`).join("");
  $("presets")
    .querySelectorAll(".preset")
    .forEach((button) =>
      button.addEventListener("click", () => {
        $("sql").value = presets[Number(button.dataset.i)].sql;
        runConsole();
      })
    );
}

async function renderSource() {
  const view = state.srcView;
  if (!state.sources[view]) {
    try {
      state.sources[view] = await (await fetch(FLAVORS[view].file)).text();
    } catch {
      state.sources[view] = "// could not load source";
    }
  }
  $("src-code").innerHTML = highlight(state.sources[view], "osprey");
  $("src-name").textContent = FLAVORS[view].file.replace("/wasm/", "");
}

function runConsole() {
  const status = $("sql-status");
  const out = $("sql-result");
  if (!state.db) {
    status.className = "console-status err";
    status.textContent = "SQLite is offline";
    return;
  }
  const sql = $("sql").value.trim();
  if (!sql) return;
  try {
    const t0 = performance.now();
    const result = state.db.exec(sql);
    const rows = result.length ? result[0].values.length : 0;
    out.innerHTML = result.length ? result.map(tableHTML).join("") : `<p class="muted">Statement ran - no result set.</p>`;
    out.hidden = false;
    status.className = "console-status ok";
    status.textContent = `${rows} row${rows === 1 ? "" : "s"} in ${(performance.now() - t0).toFixed(1)} ms`;
  } catch (err) {
    status.className = "console-status err";
    status.textContent = String(err.message || err);
  }
}

function setBanner(kind, html) {
  const banner = $("banner");
  banner.className = "banner" + (kind ? " " + kind : "");
  banner.innerHTML = html;
}

async function loadFlavor(flavor) {
  state.flavor = flavor;
  $("flavor-osp").setAttribute("aria-pressed", String(flavor === "osp"));
  $("flavor-ospml").setAttribute("aria-pressed", String(flavor === "ospml"));
  $("flavor-name").textContent = FLAVORS[flavor].file.replace("/wasm/", "");
  $("pn-src").textContent = FLAVORS[flavor].file.replace("/wasm/", "");
  const { manifest } = await runOsprey(flavor);
  state.manifest = manifest;
  renderMetrics();
  renderManifest();
  renderPresets();
  if (state.SQL) buildDb();
  renderQueries();
  reconcile();
}

async function boot() {
  setBanner("", `<span class="spin"></span> Running the Osprey wasm module...`);
  try {
    await loadFlavor("osp");
  } catch (err) {
    setBanner(
      "warn",
      `Could not run the Osprey module: <code>${esc(err.message || err)}</code>. Build the assets with <code>make wasm-site</code>.`
    );
    return;
  }

  renderSource();
  setBanner("", `<span class="spin"></span> Loading SQLite from sql.js...`);
  try {
    state.SQL = await initSqlite();
    buildDb();
    renderQueries();
    reconcile();
    $("pn-sqlite").textContent = "executing SQL";
    setBanner("ok", `Both engines are live. Osprey emitted the database; SQLite is running it.`);
    runConsole();
  } catch (err) {
    setBanner(
      "warn",
      `Osprey ran, but SQLite could not load from the CDN: <code>${esc(
        err.message || err
      )}</code>. The Osprey metrics and manifest are still shown.`
    );
    $("pn-sqlite").textContent = "offline";
    reconcile();
  }
}

$("flavor-osp").addEventListener("click", () => {
  loadFlavor("osp").catch((err) => setBanner("warn", `Could not load Default flavor: <code>${esc(err.message || err)}</code>`));
});
$("flavor-ospml").addEventListener("click", () => {
  loadFlavor("ospml").catch((err) => setBanner("warn", `Could not load ML flavor: <code>${esc(err.message || err)}</code>`));
});
$("rerun").addEventListener("click", () => {
  loadFlavor(state.flavor).catch((err) => setBanner("warn", `Could not rerun pipeline: <code>${esc(err.message || err)}</code>`));
});
$("run-sql").addEventListener("click", runConsole);
$("reset-db").addEventListener("click", () => {
  buildDb();
  renderQueries();
  reconcile();
  const status = $("sql-status");
  status.className = "console-status ok";
  status.textContent = "database reset to Osprey seed";
});
$("sql").addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") runConsole();
});
$("src-tabs")
  .querySelectorAll("button")
  .forEach((button) =>
    button.addEventListener("click", () => {
      state.srcView = button.dataset.src;
      $("src-tabs")
        .querySelectorAll("button")
        .forEach((candidate) => candidate.setAttribute("aria-pressed", String(candidate === button)));
      renderSource();
    })
  );

window.ospreyWasmStudio = { state, loadFlavor, query, runOsprey };
boot().catch((err) => setBanner("warn", `Unexpected error: <code>${esc(err.message || err)}</code>`));
