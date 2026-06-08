#!/usr/bin/env zsh
# Differential harness: run every examples/tested/*.osp through osprey-rs --run,
# trim, and compare to the sibling .expectedoutput (mirrors the Go JIT test).
# Usage: diff_examples.sh [--verbose] [name-filter]
set -u
ROOT=/Users/christianfindlay/Documents/Code/osprey
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
  exp="$f.expectedoutput"
  if [[ ! -f "$exp" ]]; then
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
