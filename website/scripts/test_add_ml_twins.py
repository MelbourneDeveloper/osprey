#!/usr/bin/env python3
# Tests for add-ml-twins.py translation rules.
# Run: python3 scripts/test_add_ml_twins.py   (from website/)
import importlib.util
import pathlib

_HERE = pathlib.Path(__file__).resolve().parent
_spec = importlib.util.spec_from_file_location("add_ml_twins", _HERE / "add-ml-twins.py")
add_ml_twins = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(add_ml_twins)


def test_match_arm_brace_pattern_binds_by_juxtaposition():
    """A `match` arm whose pattern carries a brace payload — `Success { value }`
    — must lower to ML juxtaposition `Success value`. Braces in a match pattern
    are invalid ML (spec 0024 [FLAVOR-ML-MATCH]; crates/osprey-syntax/src/ml/cst.rs).
    """
    default = (
        "match byteAt(\"hi\", 0) {\n"
        "    Success { value } => print(\"byte: ${value}\")\n"
        "    Error { message } => print(message)\n"
        "}"
    )
    ml = add_ml_twins.translate_block(default)
    assert ml is not None, "expected a translated ML twin"
    assert "Success value =>" in ml, f"pattern not juxtaposed:\n{ml}"
    assert "Error message =>" in ml, f"pattern not juxtaposed:\n{ml}"
    # No brace payload may survive in an arm *pattern* (before the `=>`).
    # String-interpolation braces `${...}` in an arm *body* are valid ML and OK.
    for arm in ml.splitlines():
        pat = arm.split("=>", 1)[0]
        assert "{" not in pat and "}" not in pat, f"braces leaked into ML pattern:\n{arm}"


if __name__ == "__main__":
    test_match_arm_brace_pattern_binds_by_juxtaposition()
    print("PASS: match-arm brace pattern binds by juxtaposition")
