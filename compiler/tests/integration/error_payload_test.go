// Implements [ERR-PAYLOAD] Phase 6: per-builtin assertions that the message
// bound to `Error { message }` is the function-specific text wired in
// internal/codegen, not a generic placeholder.
//
// Every fallible builtin that returns Result<T, string> gets one test case
// here. The test fails if any future change drops a message back to null or
// rebinds it to a generic string.
//
// Spec: compiler/spec/0013-ErrorHandling.md [ERR-PAYLOAD].
// Plan:  docs/plans/error-payloads.md Phase 6.1.

package integration

import (
	"strings"
	"testing"
)

func TestErrorPayloadMessages(t *testing.T) {
	cases := []struct {
		name    string
		source  string
		message string
	}{
		{
			name:    "parseInt_nonnumeric",
			source:  `match parseInt("oops") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseInt: input is not a valid integer",
		},
		{
			name: "parseInt_overflow",
			source: `match parseInt("99999999999999999999") { ` +
				`Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseInt: input is not a valid integer",
		},
		{
			name:    "parseFloat_empty",
			source:  `match parseFloat("") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseFloat: input is not a valid number",
		},
		{
			name:    "parseFloat_trailing",
			source:  `match parseFloat("3.14x") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseFloat: input is not a valid number",
		},
		{
			name:    "substring_oob_high",
			source:  `match substring("abc", 5, 10) { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "substring: index out of range",
		},
		{
			name:    "substring_negative",
			source:  `match substring("abc", -1, 2) { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "substring: index out of range",
		},
		{
			name:    "substring_inverted",
			source:  `match substring("abc", 2, 0) { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "substring: index out of range",
		},
		{
			name:    "indexOf_not_found",
			source:  `match indexOf("hello", "xyz") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "indexOf: needle not found",
		},
		{
			name:    "replace_empty_needle",
			source:  `match replace("hello", "", "x") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "replace: needle must not be empty",
		},
		{
			name:    "repeat_negative",
			source:  `match repeat("ab", -1) { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "repeat: count must be non-negative",
		},
		{
			name:    "padStart_empty_fill",
			source:  `match padStart("hi", 5, "") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "padStart: fill must not be empty",
		},
		{
			name:    "padEnd_empty_fill",
			source:  `match padEnd("hi", 5, "") { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "padEnd: fill must not be empty",
		},
		{
			name: "user_constructed_error",
			source: `fn fail() -> Result<int, string> = Error { message: "boom" }` + "\n" +
				`match fail() { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "boom",
		},
		{
			name:    "arith_propagates_parseInt_message",
			source:  `match parseInt("oops") + 5 { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseInt: input is not a valid integer",
		},
		{
			name:    "arith_chain_propagates_message",
			source:  `match (parseInt("oops") * 3) - 1 { Success { value } => print("ok") Error { message } => print(message) }`,
			message: "parseInt: input is not a valid integer",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := captureJITOutput(tc.source)
			if err != nil {
				t.Fatalf("compile/run failed for %s: %v", tc.name, err)
			}
			trimmed := strings.TrimSpace(out)
			if trimmed != tc.message {
				t.Fatalf("%s [ERR-PAYLOAD] mismatch:\n  expected message: %q\n  got output:       %q",
					tc.name, tc.message, trimmed)
			}
		})
	}
}

// TestErrorPayloadNoGenericFallback proves [ERR-PAYLOAD] §"non-conforming":
// the legacy "Error occurred" global must not appear in any compiled IR.
// Catches accidental regressions that re-introduce the hardcoded placeholder.
func TestErrorPayloadNoGenericFallback(t *testing.T) {
	cases := []string{
		`match parseInt("oops") { Success { value } => print("ok") Error { message } => print(message) }`,
		`match substring("abc", 5, 10) { Success { value } => print("ok") Error { message } => print(message) }`,
		`match indexOf("a", "z") { Success { value } => print("ok") Error { message } => print(message) }`,
		`match replace("a", "", "b") { Success { value } => print("ok") Error { message } => print(message) }`,
		`fn fail() -> Result<int, string> = Error { message: "real msg" }` + "\n" +
			`match fail() { Success { value } => print("ok") Error { message } => print(message) }`,
	}

	for i, src := range cases {
		out, err := captureJITOutput(src)
		if err != nil {
			t.Fatalf("case %d compile/run failed: %v", i, err)
		}
		if strings.Contains(out, "Error occurred") {
			t.Fatalf("case %d emitted legacy placeholder 'Error occurred'; got: %q", i, out)
		}
	}
}
