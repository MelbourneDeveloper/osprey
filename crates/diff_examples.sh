#!/usr/bin/env zsh
# Differential harness: run every examples/tested/*.osp through osprey-rs --run,
# trim, and compare to the sibling .expectedoutput (mirrors the Go JIT test).
# Usage: diff_examples.sh [--verbose] [name-filter]
set -u
# Repo root: derived from this script's location (crates/diff_examples.sh) so the
# harness runs unchanged on a dev box and in CI; override with OSPREY_ROOT.
ROOT=${OSPREY_ROOT:-${0:A:h}/..}
ROOT=${ROOT:A}
BIN=$ROOT/target/release/osprey-rs
EXDIR=$ROOT/compiler/examples/tested
VERBOSE=0
FILTER=""
for a in "$@"; do
  case "$a" in
    --verbose) VERBOSE=1 ;;
    *) FILTER="$a" ;;
  esac
done

pass=0; fail=0; noexp=0; comperr=0
typeset -a FAILED
for f in $(find $EXDIR -name '*.osp' | sort); do
  rel=${f#$EXDIR/}
  [[ -n "$FILTER" && "$rel" != *"$FILTER"* ]] && continue
  # Expected-output precedence:
  #   1. .expectedoutput.rust  — a Rust-specific oracle for the rare case where
  #      osprey-rs is MORE correct than the Go binary that seeds the shared file
  #      (comprehensive_math: Rust computes the mathematically-correct 18 where a
  #      Go Result-chain bug yields 2). The Go test never reads this name, so its
  #      suite keeps asserting the Go value — both stay green on the same example.
  #   2. .expectedoutput        — the shared oracle (the Go binary's output).
  #   3. .expectedoutput.<uname> — OS-dependent output (callback_stdout_demo: its
  #      subprocess error text + exit code differ Darwin vs Linux; the Go test
  #      branches on runtime.GOOS for the same reason).
  if [[ -f "$f.expectedoutput.rust" ]]; then
    exp="$f.expectedoutput.rust"
  elif [[ -f "$f.expectedoutput" ]]; then
    exp="$f.expectedoutput"
  elif [[ -f "$f.expectedoutput.$(uname -s)" ]]; then
    exp="$f.expectedoutput.$(uname -s)"
  else
    noexp=$((noexp+1))
    [[ $VERBOSE -eq 1 ]] && echo "NOEXP  $rel"
    continue
  fi
  # Mirror the Go test exactly: compare strings.TrimSpace(actual) to
  # strings.TrimSpace(expected) — a single whole-string trim on each, never a
  # per-line strip (which would drop trailing whitespace the program emits).
  expected=$(cat "$exp")
  actual=$($BIN "$f" --run 2>/tmp/osprey_rs_err.txt)
  rc=$?
  expected_trim="${expected#"${expected%%[![:space:]]*}"}"; expected_trim="${expected_trim%"${expected_trim##*[![:space:]]}"}"
  actual_trim="${actual#"${actual%%[![:space:]]*}"}"; actual_trim="${actual_trim%"${actual_trim##*[![:space:]]}"}"
  if [[ "$actual_trim" == "$expected_trim" ]]; then
    pass=$((pass+1))
    [[ $VERBOSE -eq 1 ]] && echo "PASS   $rel"
  else
    fail=$((fail+1))
    FAILED+=("$rel")
    if [[ $VERBOSE -eq 1 ]]; then
      echo "FAIL   $rel (rc=$rc)"
      echo "  --- err ---"; head -3 /tmp/osprey_rs_err.txt | sed 's/^/  /'
    fi
  fi
done
echo "================================"
echo "PASS=$pass FAIL=$fail NOEXP=$noexp (of $((pass+fail+noexp)))"
echo "FAILED:"
for x in $FAILED; do echo "  $x"; done
