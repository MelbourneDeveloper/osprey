package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestCompilationFailures tests that examples in failscompilation directory fail compilation with expected errors.
func TestCompilationFailures(t *testing.T) {
	failsDir := "../../examples/failscompilation"
	entries, err := os.ReadDir(failsDir)
	if err != nil {
		t.Fatalf("Failed to read failscompilation directory: %v", err)
	}

	// YOU ARE NOT ALLOWED TO LET RUNTIME ERRORS APPEAR HERE!!!

	// Define expected error patterns for each test case
	expectedErrors := map[string]string{
		"match_not_exhaustive.osp":   "match expression not exhaustive: missing patterns: [Blue]",
		"match_not_exhaustive.ospo":  "match expression not exhaustive: missing patterns: [Blue]",
		"match_unknown_variant.osp":  "unknown variant 'Maybe' in match expression for type 'Color'",
		"match_unknown_variant.ospo": "unknown variant 'Maybe' in match expression for type 'Color'",
		"built_in_function_redefinition.osp": "Parameter 'x' in function 'toString' " +
			"requires explicit type annotation - type cannot be inferred from usage",
		"built_in_function_redefinition.ospo": "Parameter 'x' in function 'toString' " +
			"requires explicit type annotation - type cannot be inferred from usage",
		"named_args_violation.osp": "function requires named arguments 'add' has 2 parameters " +
			"and requires named arguments. Use: add(x: value, y: value)",
		"named_args_violation.ospo": "function requires named arguments 'add' has 2 parameters " +
			"and requires named arguments. Use: add(x: value, y: value)",
		"type_inference_ambiguity.osp": "Function 'identity' requires explicit return type " +
			"annotation - type cannot be inferred from body",
		"type_inference_ambiguity.ospo": "Function 'identity' requires explicit return type " +
			"annotation - type cannot be inferred from body",
		"undefined_variable.osp":                 "undefined variable 'unknownVar': undefined variable",
		"undefined_variable.ospo":                "undefined variable 'unknownVar': undefined variable",
		"constraint_field_access_violation.osp":  "field access not implemented for field 'name'",
		"constraint_field_access_violation.ospo": "field access not implemented for field 'name'",
		"simple_field_access.osp":                "field access not implemented for field 'name'",
		"simple_field_access.ospo":               "field access not implemented for field 'name'",
		"debug_interpolation.osp":                "field access not implemented for field 'name'",
		"debug_interpolation.ospo":               "field access not implemented for field 'name'",
		"any_function_arg.osp":                   "cannot pass 'any' type to function expecting specific type",
		"any_direct_arithmetic.ospo":             "cannot pass 'any' type to function expecting specific type",
		"any_direct_variable_access.osp":         "direct access to 'any' type variable",
	}

	// Test each .osp and .ospo file (both extensions should be tested)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".osp") && !strings.HasSuffix(entry.Name(), ".ospo") {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			filePath := filepath.Join(failsDir, entry.Name())

			// Read the file
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", filePath, err)
			}

			source := string(content)

			// Attempt compilation - this should fail
			_, err = codegen.CompileToLLVM(source)

			if err == nil {
				t.Errorf("File %s should have failed compilation but succeeded", entry.Name())
				return //nolint:nlreturn
			}

			// Check if error message contains expected pattern
			if expectedPattern, ok := expectedErrors[entry.Name()]; ok {
				if !strings.Contains(err.Error(), expectedPattern) {
					t.Errorf("File %s failed compilation but with unexpected error.\nExpected to contain: %s\nActual error: %s",
						entry.Name(), expectedPattern, err.Error())
				} else {
					t.Logf("✓ File %s correctly failed with expected error: %s", entry.Name(), err.Error())
				}
			} else {
				// No specific pattern defined, just verify it failed
				t.Logf("✓ File %s correctly failed compilation: %s", entry.Name(), err.Error())
			}
		})
	}
}
