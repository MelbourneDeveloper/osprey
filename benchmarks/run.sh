#!/usr/bin/env zsh
# Cross-language benchmark harness for Osprey.
#
# For every case under benchmarks/cases/<name>/, compiles each language's
# implementation to a NATIVE binary ONCE, verifies its output byte-for-byte
# against expected.txt (a broken program is never timed), then measures:
#   * CPU time   — hyperfine (warmup + many runs, statistical mean/stddev)
#   * Peak memory — /usr/bin/time RSS (max over a few runs)
# Absent toolchains (rustc/ghc/ocamlopt/cc) are skipped and reported, so the
# suite runs today with whatever is installed and lights up the rest later.
#
# Mirrors the conventions of crates/diff_examples.sh (zsh, set -u, ROOT from
# the script path). Results land in benchmarks/results/ (gitignored).
#
# Usage: run.sh [name-filter]
#   BENCH_WARMUP   (default 3)   warmup runs per command
#   BENCH_MINRUNS  (default 10)  minimum timed runs per command
#   BENCH_MEMRUNS  (default 3)   memory-sampling runs per command
set -u

BENCHDIR=${0:A:h}
ROOT=${BENCHDIR}/..
ROOT=${ROOT:A}
OSP=$ROOT/target/release/osprey
CASEDIR=$BENCHDIR/cases
OUT=$BENCHDIR/results
TMP=$OUT/tmp
BINDIR=$OUT/bin
HFDIR=$OUT/hf
RAW=$OUT/raw.jsonl

FILTER=${1:-}
WARMUP=${BENCH_WARMUP:-3}
MINRUNS=${BENCH_MINRUNS:-10}
MEMRUNS=${BENCH_MEMRUNS:-3}

# Language order is the report's column order. "Speed of light" baselines (C,
# Rust) first after Osprey so the gap to Osprey reads left-to-right. `osprey-gc`
# is the SAME .osp compiled with the tracing GC backend (--memory=gc), so the
# allocation cases (binarytrees) show reclamation next to the default backend.
LANGS=(osprey osprey-gc rust c ocaml haskell)
typeset -A EXT
EXT=(osprey osp  osprey-gc osp  rust rs  c c  ocaml ml  haskell hs)

have() { command -v "$1" >/dev/null 2>&1 }

# toolchain_ok <lang> — is the compiler for <lang> installed?
toolchain_ok() {
  case "$1" in
    osprey|osprey-gc) [[ -x "$OSP" ]] ;;
    rust)    have rustc ;;
    c)       have cc ;;
    ocaml)   have ocamlopt ;;
    haskell) have ghc ;;
  esac
}

# build <lang> <casedir> <name> <out-binary> — compile to a native binary.
build() {
  local lang=$1 dir=$2 name=$3 out=$4
  case "$lang" in
    osprey)    ( cd "$dir" && "$OSP" "$name.osp" --compile >/dev/null 2>&1 ) && mv -f "$dir/$name" "$out" ;;
    osprey-gc) ( cd "$dir" && "$OSP" "$name.osp" --memory=gc --compile >/dev/null 2>&1 ) && mv -f "$dir/$name" "$out" ;;
    rust)    rustc -C opt-level=3 -C overflow-checks=off -o "$out" "$dir/$name.rs" 2>/dev/null ;;
    c)       cc -O2 -o "$out" "$dir/$name.c" 2>/dev/null ;;
    ocaml)   cp "$dir/$name.ml" "$TMP/$name.ml" && \
             ( cd "$TMP" && ocamlopt -O3 -unsafe -o "$out" "$name.ml" >/dev/null 2>&1 ) ;;  # compile a copy: ocamlopt litters .cmi/.cmx/.o beside the source
    haskell) ghc -O2 -outputdir "$TMP/hs_$name" -o "$out" "$dir/$name.hs" >/dev/null 2>&1 ;;
  esac
}

# peak_rss <binary> — max resident set size in bytes over MEMRUNS runs.
peak_rss() {
  local bin=$1 best=0 v
  for _ in $(seq 1 $MEMRUNS); do
    if [[ "$(uname)" == Darwin ]]; then
      /usr/bin/time -l "$bin" <"$MODE0" >/dev/null 2>"$TMP/mem.txt"
      v=$(awk '/maximum resident set size/ {print $1}' "$TMP/mem.txt")
    else
      /usr/bin/time -v "$bin" <"$MODE0" >/dev/null 2>"$TMP/mem.txt"
      v=$(awk -F: '/Maximum resident set size/ {gsub(/ /,"",$2); print $2*1024}' "$TMP/mem.txt")
    fi
    [[ -n "$v" && "$v" -gt "$best" ]] && best=$v
  done
  print -r -- "$best"
}

# json_row — append one {case,lang,...} record to the raw results log.
json_row() {
  printf '{"case":"%s","lang":"%s","status":"%s","output":"%s","expected":"%s","rss":%s}\n' \
    "$1" "$2" "$3" "$4" "$5" "${6:-0}" >> "$RAW"
}

rm -rf "$OUT"; mkdir -p "$TMP" "$BINDIR" "$HFDIR"; : > "$RAW"

# Benchmark input modes. A case may read its first stdin line to pick a token
# seed: "0" => a fixed seed (fully deterministic — what we time and oracle), "1"
# => a fresh cryptographically-secure seed (same workload, unpredictable data).
# We always feed MODE0 so the suite is reproducible AND so a case that calls
# input() never blocks on a tty. `BENCH_RANDOM=1` adds a randomized demo pass.
MODE0="$TMP/mode_const"; print -- "0" > "$MODE0"
MODE1="$TMP/mode_rand";  print -- "1" > "$MODE1"

if [[ ! -x "$OSP" ]]; then
  echo "FATAL: osprey binary not found at $OSP — run 'make build' first." >&2
  exit 1
fi

echo "==> toolchains:"
for l in $LANGS; do printf '    %-8s %s\n' "$l" "$(toolchain_ok $l && echo present || echo ABSENT)"; done
echo "==> warmup=$WARMUP min-runs=$MINRUNS mem-runs=$MEMRUNS"

fail=0
for dir in $CASEDIR/*/(/); do
  name=${${dir%/}:t}
  [[ -n "$FILTER" && "$name" != *"$FILTER"* ]] && continue
  [[ -f "$dir/expected.txt" ]] || { echo "SKIP $name (no expected.txt)"; continue }
  expected=$(<"$dir/expected.txt"); expected=${expected//[[:space:]]/}
  echo "\n==> $name  (expected: $expected)"

  typeset -a hf_args; hf_args=()
  typeset -a ok_pairs; ok_pairs=()
  for lang in $LANGS; do
    src="$dir/$name.${EXT[$lang]}"
    toolchain_ok "$lang" || continue
    [[ -f "$src" ]] || { echo "    $lang: no source"; continue }
    bin="$BINDIR/${name}__${lang}"
    if ! build "$lang" "$dir" "$name" "$bin"; then
      echo "    $lang: BUILD FAILED"; json_row "$name" "$lang" "build_failed" "" "$expected"; fail=1; continue
    fi
    actual=$("$bin" <"$MODE0" 2>/dev/null); actual=${actual//[[:space:]]/}
    if [[ "$actual" != "$expected" ]]; then
      echo "    $lang: WRONG OUTPUT ($actual != $expected) — excluded from timing"
      json_row "$name" "$lang" "wrong_output" "$actual" "$expected"; fail=1; continue
    fi
    rss=$(peak_rss "$bin")
    json_row "$name" "$lang" "ok" "$actual" "$expected" "$rss"
    echo "    $lang: ok  (rss $(( rss / 1024 )) KiB)"
    hf_args+=(-n "$lang" "$bin")
    ok_pairs+=("$lang=$bin")
  done

  if (( ${#hf_args} > 0 )); then
    # --input feeds MODE0 to every timed run: the constant seed keeps the
    # measurement reproducible (and matches the oracle above).
    hyperfine -N --input "$MODE0" --warmup "$WARMUP" --min-runs "$MINRUNS" \
      --export-json "$HFDIR/$name.json" $hf_args >/dev/null 2>&1 \
      || echo "    (hyperfine failed for $name)"
  fi

  # Randomized demo pass: prove each case also runs on unpredictable input. Feeds
  # MODE1 (cryptographically-secure seed) — outputs vary run-to-run, so this is
  # never oracle-checked or charted, just shown. Opt-in: BENCH_RANDOM=1.
  if [[ -n "${BENCH_RANDOM:-}" && ${#ok_pairs} -gt 0 ]]; then
    echo "    [random] same workload, cryptographically-secure seed:"
    for pair in $ok_pairs; do
      rl=${pair%%=*}; rb=${pair#*=}
      echo "      $rl: $("$rb" <"$MODE1" 2>/dev/null) / $("$rb" <"$MODE1" 2>/dev/null) (constant=$expected)"
    done
  fi
done

echo "\n==> rendering report"
python3 "$BENCHDIR/report.py" "$OUT" || { echo "report failed" >&2; exit 1; }
rm -rf "$TMP"
echo "==> done. Open $OUT/results.html"
exit $fail
