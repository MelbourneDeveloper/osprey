#!/usr/bin/env python3
"""Aggregate harness output into results.md + results.json.

Reads <out>/raw.jsonl (per case/lang: status + peak RSS) and the hyperfine
exports in <out>/hf/<case>.json (per case/lang: timing), and renders a CPU
table, a peak-memory table, and Osprey-relative slowdown factors.
"""
import json
import math
import sys
from pathlib import Path
from typing import Callable, cast

ORDER: list[str] = ["osprey", "rust", "c", "ocaml", "haskell"]
LABEL: dict[str, str] = {"osprey": "Osprey", "rust": "Rust", "c": "C", "ocaml": "OCaml", "haskell": "Haskell"}

Cell = dict[str, object]
Data = dict[str, dict[str, Cell]]


def load(out: Path) -> Data:
    """Merge raw status/RSS records with hyperfine timings, keyed by case/lang."""
    data: Data = cast(Data, {})
    for line in (out / "raw.jsonl").read_text().splitlines():
        r = json.loads(line)
        data.setdefault(r["case"], {})[r["lang"]] = {"status": r["status"], "rss": r["rss"]}
    for hf in (out / "hf").glob("*.json"):
        case = hf.stem
        for res in json.loads(hf.read_text())["results"]:
            cell = data.setdefault(case, {}).setdefault(res["command"], {})
            cell.update(mean=res["mean"], stddev=res["stddev"], min=res["min"], max=res["max"])
    return data


def present_langs(data: Data) -> list[str]:
    seen = {lang for case in data.values() for lang, c in case.items() if c.get("status") == "ok"}
    return [l for l in ORDER if l in seen]


def fmt_time(sec: float) -> str:
    ms = sec * 1000.0
    return f"{ms:.1f} ms" if ms < 1000 else f"{sec:.3f} s"


def fmt_mem(b: float) -> str:
    mib = b / (1024 * 1024)
    return f"{b / 1024:.0f} KiB" if mib < 1 else f"{mib:.1f} MiB"


def cell_time(c: Cell) -> str:
    if c.get("status") != "ok":
        return {"build_failed": "build✗", "wrong_output": "wrong✗"}.get(str(c.get("status")), "—")
    if "mean" not in c:
        return "—"
    return f"{fmt_time(float(c['mean']))} ± {float(c['stddev']) * 1000:.1f}"


def cell_mem(c: Cell) -> str:
    return fmt_mem(float(c["rss"])) if c.get("status") == "ok" and c.get("rss") else "—"


def md_table(data: Data, langs: list[str], cases: list[str], cellfn: Callable[[Cell], str]) -> str:
    head = "| Benchmark | " + " | ".join(LABEL[l] for l in langs) + " |\n"
    sep = "|" + "---|" * (len(langs) + 1) + "\n"
    rows = "".join(
        "| " + name + " | " + " | ".join(cellfn(data[name].get(l, {})) for l in langs) + " |\n"
        for name in cases
    )
    return head + sep + rows


def fastest(case: dict[str, Cell]) -> float:
    means = [float(c["mean"]) for c in case.values() if c.get("status") == "ok" and "mean" in c]
    return min(means) if means else float("nan")


def geomean(xs: list[float]) -> float:
    ys = [x for x in xs if x and x > 0]
    return math.exp(sum(math.log(x) for x in ys) / len(ys)) if ys else float("nan")


def ratio_table(data: Data, langs: list[str], cases: list[str]) -> str:
    """Each cell = mean / fastest-in-row (1.0× == fastest language for that benchmark)."""
    def cell(case: str, lang: str) -> str:
        c = data[case].get(lang, {})
        if c.get("status") != "ok" or "mean" not in c:
            return "—"
        f = fastest(data[case])
        return f"{float(c['mean']) / f:.1f}×" if f == f else "—"
    head = "| Benchmark | " + " | ".join(LABEL[l] for l in langs) + " |\n"
    sep = "|" + "---|" * (len(langs) + 1) + "\n"
    rows = "".join(
        "| " + n + " | " + " | ".join(cell(n, l) for l in langs) + " |\n" for n in cases
    )
    return head + sep + rows


def _other_min(case: dict[str, Cell], key: str) -> float:
    """Best (lowest) value of `key` among the non-Osprey languages — the bar Osprey is measured against."""
    vs = [float(c[key]) for lang, c in case.items()
          if lang != "osprey" and c.get("status") == "ok" and c.get(key)]
    return min(vs) if vs else float("nan")


def tuning_table(data: Data, cases: list[str]) -> str:
    """Rank benchmarks by Osprey's CPU slowdown vs the fastest OTHER language — worst gap first.

    This is the work list: the top rows are where Osprey's runtime/codegen needs
    attention most. Peak-memory overhead (Osprey RSS vs the lightest other
    language) is shown alongside — it is the dominant signal in this suite.
    """
    rows: list[tuple[float, str, float, float, float]] = []
    for n in cases:
        o = data[n].get("osprey", {})
        if o.get("status") != "ok" or "mean" not in o:
            continue
        cpu_ref = _other_min(data[n], "mean")
        rss_ref = _other_min(data[n], "rss")
        rss = float(o.get("rss") or 0)
        slow = float(o["mean"]) / cpu_ref if cpu_ref == cpu_ref else float("nan")
        memx = rss / rss_ref if (rss_ref == rss_ref and rss_ref) else float("nan")
        rows.append((slow, n, float(o["mean"]), rss, memx))
    rows.sort(reverse=True)
    head = "| Rank | Benchmark | Osprey CPU | × slower vs fastest | Osprey peak RSS | × heavier |\n"
    sep = "|---|---|---|---|---|---|\n"
    body = "".join(
        f"| {i + 1} | {n} | {fmt_time(mean)} | {slow:.0f}× | {fmt_mem(rss)} | {memx:.0f}× |\n"
        for i, (slow, n, mean, rss, memx) in enumerate(rows)
    )
    return head + sep + body


def _ratios(data: Data, cases: list[str], lang: str, key: str) -> list[float]:
    rs: list[float] = []
    for n in cases:
        o, x = data[n].get("osprey", {}), data[n].get(lang, {})
        if o.get("status") == "ok" and x.get("status") == "ok" and o.get(key) and x.get(key):
            rs.append(float(o[key]) / float(x[key]))
    return rs


def osprey_summary(data: Data, langs: list[str], cases: list[str]) -> str:
    """Geomean of Osprey's slowdown (CPU) and memory overhead vs each language."""
    out: list[str] = ["**CPU** (geometric mean across benchmarks):"]
    for lang in langs:
        if lang != "osprey" and (rs := _ratios(data, cases, lang, "mean")):
            out.append(f"- {geomean(rs):.1f}× slower than {LABEL[lang]}")
    out.append("\n**Peak memory** (geometric mean across benchmarks):")
    for lang in langs:
        if lang != "osprey" and (rs := _ratios(data, cases, lang, "rss")):
            out.append(f"- {geomean(rs):.0f}× more than {LABEL[lang]}")
    return "\n".join(out)


def render(out: Path) -> None:
    data = load(out)
    langs = present_langs(data)
    cases = sorted(data)
    md = ["# Osprey cross-language benchmarks\n",
          f"Languages compared: {', '.join(LABEL[l] for l in langs)}.\n",
          "Lower is better. CPU = hyperfine mean ± stddev; memory = peak resident set size.\n",
          "\n## Tuning priorities — where Osprey needs work (worst gap first)\n",
          "Ranked by Osprey's CPU slowdown against the *fastest other* language on each "
          "benchmark. Start at the top. Memory overhead (vs the lightest other language) is "
          "shown alongside — in this suite it is the bigger problem and tracks operation count, "
          "not data size.\n\n",
          tuning_table(data, cases),
          "\n## CPU time\n", md_table(data, langs, cases, cell_time),
          "\n## Relative speed (× the fastest language per benchmark)\n", ratio_table(data, langs, cases),
          "\n## Peak memory\n", md_table(data, langs, cases, cell_mem),
          "\n## Where Osprey sits\n", osprey_summary(data, langs, cases), "\n"]
    (out / "results.md").write_text("".join(md))
    (out / "results.json").write_text(json.dumps({"languages": langs, "cases": data}, indent=2))
    print("".join(md))


if __name__ == "__main__":
    render(Path(sys.argv[1]))
