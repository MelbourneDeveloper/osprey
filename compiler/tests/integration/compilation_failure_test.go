package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
	"github.com/stretchr/testify/require"
)

// TestCompilationFailures tests that examples in failscompilation directory fail compilation with expected errors.
func TestCompilationFailures(t *testing.T) {
	failsDir := "../../examples/failscompilation"
	entries, err := os.ReadDir(failsDir)
	require.NoError(t, err, "Failed to read failscompilation directory: %v", err)

	// YOU ARE NOT ALLOWED TO LET RUNTIME ERRORS APPEAR HERE!!!
	// EVERY .ospo FILE MUST HAVE A CORRESPONDING .expectedoutput FILE!

	// Test each .osp and .ospo file (both extensions should be tested)
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".osp") && !strings.HasSuffix(entry.Name(), ".ospo") {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			filePath := filepath.Join(failsDir, entry.Name())
			expectedOutputPath := filePath + ".expectedoutput"

			// 🚨 SCREAM FAILURE IF NO .expectedoutput FILE EXISTS! 🚨
			expectedContent, err := os.ReadFile(expectedOutputPath)
			require.NoError(t, err,
				"🚨🚨🚨 MISSING .expectedoutput FILE FOR %s! 🚨🚨🚨\n"+
					"EVERY FAILURE TEST MUST HAVE EXACT EXPECTED ERROR!\n"+
					"Create: %s\n"+
					"Run the test manually and copy the EXACT error message!",
				entry.Name(), expectedOutputPath)

			expectedError := strings.TrimSpace(string(expectedContent))
			require.NotEmpty(t, expectedError,
				"🚨🚨🚨 EMPTY .expectedoutput FILE FOR %s! 🚨🚨🚨\n"+
					"The .expectedoutput file must contain the EXACT error message!",
				entry.Name())

			// Read the source file
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Failed to read %s: %v", filePath, err)

			source := string(content)

			// Attempt compilation - this should fail
			_, err = codegen.CompileToLLVM(source)

			require.Error(t, err, "File %s should have failed compilation but succeeded", entry.Name())

			// 🎯 EXACT ERROR MESSAGE MATCH REQUIRED! 🎯
			actualError := strings.TrimSpace(err.Error())
			require.Equal(t, expectedError, actualError,
				"🚨 EXACT ERROR MISMATCH FOR %s! 🚨\n"+
					"Expected EXACT error: %q\n"+
					"Actual error:         %q\n"+
					"Update the .expectedoutput file with the correct error message!",
				entry.Name(), expectedError, actualError)

			t.Logf("✅ File %s correctly failed with EXACT expected error: %s", entry.Name(), actualError)
		})
	}
}
