package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestBasicsExamples tests the basic language feature examples.
func TestBasicsExamples(t *testing.T) {
	checkLLVMTools(t)

	runTestExamplesRecursive(t, "../../examples/tested/basics", getExpectedOutputs())
}

// TestEffectsExamples tests the algebraic effects examples.
func TestEffectsExamples(t *testing.T) {
	checkLLVMTools(t)

	examplesDir := "../../examples/tested/effects"
	// Effects examples use same expected outputs map as basics, with .expectedoutput file fallback
	runTestExamplesRecursive(t, examplesDir, getExpectedOutputs())
}

// TestRustIntegrationExamples tests the Rust interop examples.
func TestRustIntegrationExamples(t *testing.T) {
	checkLLVMTools(t)

	// Check if Rust tools are available before running the test
	_, _, err := findRustTools()
	if err != nil {
		t.Fail()
		return
	}

	examplesDir := "../../examples/rust_integration"
	// Rust integration examples use same expected outputs map with .expectedoutput file fallback
	runTestExamplesRecursive(t, examplesDir, getExpectedOutputs())
}

// TestDatabaseExamples tests the SQLite-over-FFI database examples
// (examples/tested/db). They link libsqlite3 via the `// @link: sqlite3`
// directive — CI installs libsqlite3-dev; macOS ships it in the SDK.
func TestDatabaseExamples(t *testing.T) {
	checkLLVMTools(t)

	runTestExamplesRecursive(t, "../../examples/tested/db", getExpectedOutputs())
}

func runTestExamplesRecursive(t *testing.T, examplesDir string, expectedOutputs map[string]string) {
	err := filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .osp files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".osp") {
			// Create test name from the file path relative to examples/tested
			relPath, _ := filepath.Rel(examplesDir, path)
			testName := strings.TrimSuffix(relPath, ".osp")
			testName = strings.ReplaceAll(testName, string(filepath.Separator), "/")

			t.Run(testName, func(t *testing.T) {
				// Try to read from .expectedoutput file first
				expectedOutputPath := path + ".expectedoutput"
				expectedContent, err := os.ReadFile(expectedOutputPath)
				if err == nil {
					// Use .expectedoutput file content, trimmed to match captureJITOutput behavior
					expectedOutput := strings.TrimSpace(string(expectedContent))
					testExampleFileWithTrimming(t, path, expectedOutput, true)

					return
				}

				// Fallback to hardcoded expected outputs
				expectedOutput, exists := expectedOutputs[info.Name()]
				if !exists {
					t.Fatalf("❌ MISSING expected output for %s!\n"+
						"🚨 CREATE: %s\n"+
						"🚨 OR ADD TO expectedOutputs MAP!\n"+
						"🚨 RUN: ../../osprey %s --run\n"+
						"🚨 Then copy the output to create the .expectedoutput file!",
						info.Name(), expectedOutputPath, info.Name())
				}

				testExampleFileWithTrimming(t, path, expectedOutput, false)
			})
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk examples directory: %v", err)
	}
}

// getExpectedOutputs returns expected outputs for the handful of .osp files
// that intentionally don't carry a sibling .expectedoutput on disk. New
// tests should add an .expectedoutput file alongside the .osp rather than
// extend this map — the file-on-disk path is preferred and platform-aware
// outputs are kept here only when they need to branch on runtime.GOOS.
func getExpectedOutputs() map[string]string {
	return map[string]string{
		// Adventure game keeps its expected output inline because it's the
		// reference exercise of nested matches across rooms.
		"adventure_game.osp": getAdventureGameExpectedOutput(),
		// Lambdas, string compare, mut auto-unwrap, generic union payload,
		// escape-sequence handling, unary-minus on floats — pinned literal.
		"function_composition_test.osp": "=== Function Composition Test ===\n" +
			"Testing function composition...\n" +
			"Starting value: 10\n" +
			"After double: 20\n" +
			"After triple: 30\n" +
			"After add5: 15\n" +
			"square(5) = 25\n" +
			"Function composition working correctly!\n" +
			"hi!\n" +
			"11\n" +
			"110\n" +
			"strEq=hi\n" +
			"strNe=ok\n" +
			"strLt=ok\n" +
			"strLt2=ok\n" +
			"3\n" +
			"7\n" +
			"Quote: \"hi\"\n" +
			"Hi \"Bob\"\n" +
			"lit: ab\\nc\n" +
			"-2.5\n" +
			"=== Function Composition Test Complete ===\n",
		// Every [BUILTIN-STRING-*] error path and boundary in one pinned blob.
		"string_edge_cases.osp": "empty-emp=true\n\nempty-x=false\n\n" +
			"starts-empty=true\n\nends-empty=true\n\n" +
			"starts-long=false\n\nends-long=false\n\n" +
			"take-neg=\"\"\n\ntake-big=\"hi\"\n\n" +
			"drop-neg=\"hi\"\n\ndrop-big=\"\"\n\n" +
			"trim-allws=\"\"\n\ntrim-empty=\"\"\n\n" +
			"sub-empty=\"\"\n\nsub-full=\"hi\"\n\n" +
			"sub-neg rejected\n\nsub-over rejected\n\nsub-inverted rejected\n\n" +
			"idx-empty=0\n\nidx-missing rejected\n\n" +
			"rep-empty rejected\n\nrep-shrink=\"heo\"\n\n" +
			"rep0=\"\"\n\nrep-neg rejected\n\n" +
			"padS-empty rejected\n\npadE-empty rejected\n\n" +
			"padS-noop=\"hi\"\n\npadS-multi=\"aba42\"\n\n" +
			"pf-empty rejected\n\npf-trailing rejected\n\n" +
			"pi-overflow rejected\n\npi-min=-9223372036854775808\n\n" +
			"pi-leading-space rejected\n\n" +
			"=== Edge Cases Complete ===\n\n",
		// Platform-specific: ls error message + exit code differ across OSes.
		"callback_stdout_demo.osp": getCallbackStdoutDemoExpectedOutput(),
		// Rust interop demo (separate TestRustIntegrationExamples runner).
		"demo.osp": "Rust add(15, 25) = 40\n" +
			"Rust multiply(6, 7) = 42\n" +
			"Rust factorial(5) = 120\n" +
			"Rust fibonacci(10) = 55\n" +
			"Rust is_prime(17) = 1\n" +
			"✅ Rust-Osprey integration demo completed successfully!\n",
	}
}


// testExampleFileWithTrimming tests a single example file with optional output trimming.
func testExampleFileWithTrimming(t *testing.T, filePath, expectedOutput string, trimActualOutput bool) {
	t.Helper()

	// Read the file content - needed for captureJITOutput
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", filePath, err)
	}

	source := string(content)

	// ERROR: ALL EXAMPLES MUST HAVE VERIFIED OUTPUT!
	if expectedOutput == "" {
		expectedOutputPath := filePath + ".expectedoutput"
		t.Fatalf("❌ MISSING EXPECTED OUTPUT FOR %s!\n"+
			"🚨 ALL EXAMPLES MUST HAVE VERIFIED OUTPUT!\n"+
			"🚨 NO COMPILATION-ONLY TESTS ALLOWED!\n"+
			"🚨 CREATE: %s\n"+
			"🚨 RUN: ../../osprey %s --run\n"+
			"🚨 Then copy the output to create the .expectedoutput file!",
			filepath.Base(filePath), expectedOutputPath, filepath.Base(filePath))
	}

	// Execute via CLI interface AND capture output - we need both coverage AND verification!
	// This exercises the runRunProgram function while capturing the actual output
	output, err := captureJITOutput(source)
	if err != nil {
		// If JIT execution fails due to missing tools, fail the test
		if strings.Contains(err.Error(), "LLVM tools not found") ||
			strings.Contains(err.Error(), "no suitable compiler found") {
			t.Fatalf("❌ LLVM TOOLS NOT FOUND - TEST FAILED for %s: %v", filePath, err)
		}

		t.Fatalf("Failed to execute %s: %v", filePath, err)
	}

	// THE MOST CRITICAL PART: Verify output matches expected!
	// If expected output came from a file (and was trimmed), trim actual output too
	actualOutput := output
	if trimActualOutput {
		actualOutput = strings.TrimSpace(output)
	}

	if actualOutput != expectedOutput {
		t.Fatalf("Output mismatch for %s:\nExpected: %q\nGot:      %q", filePath, expectedOutput, actualOutput)
	}

	t.Logf("✅ Example %s executed and output verified", filepath.Base(filePath))
}

// Helper functions for expected outputs.
func getCallbackStdoutDemoExpectedOutput() string {
	// Platform-specific expected output due to different ls behavior
	if runtime.GOOS == "darwin" {
		// macOS behavior: exit code 1, different error message format
		return "=== CALLBACK-BASED STDOUT COLLECTION DEMO ===\n" +
			"--- Test 1: Basic Stdout Callback ---\n" +
			"✓ Process spawned with ID: 1\n" +
			"[CALLBACK] Process 1 STDOUT: Hello from callback!\n\n" +
			"[CALLBACK] Process 1 EXIT: 0\n" +
			"✓ Process finished with exit code: 0\n" +
			"✓ Process cleaned up\n" +
			"--- Test 2: Multiple Lines Callback ---\n" +
			"✓ Multi-line process spawned with ID: 2\n" +
			"[CALLBACK] Process 2 STDOUT: Line 1\nLine 2\nLine 3\n\n" +
			"[CALLBACK] Process 2 EXIT: 0\n" +
			"✓ Multi-line process finished\n" +
			"--- Test 3: Error Process Callback ---\n" +
			"✓ Error process spawned with ID: 3\n" +
			"[CALLBACK] Process 3 STDERR: ls: /nonexistent/directory: No such file or directory\n\n" +
			"[CALLBACK] Process 3 EXIT: 1\n" +
			"✓ Error process finished with exit code: 1\n" +
			"=== CALLBACK DEMO COMPLETE ===\n" +
			"The [CALLBACK] lines above show C runtime calling into Osprey!\n"
	}

	// Linux behavior: exit code 2, different error message format
	return "=== CALLBACK-BASED STDOUT COLLECTION DEMO ===\n" +
		"--- Test 1: Basic Stdout Callback ---\n" +
		"✓ Process spawned with ID: 1\n" +
		"[CALLBACK] Process 1 STDOUT: Hello from callback!\n\n" +
		"[CALLBACK] Process 1 EXIT: 0\n" +
		"✓ Process finished with exit code: 0\n" +
		"✓ Process cleaned up\n" +
		"--- Test 2: Multiple Lines Callback ---\n" +
		"✓ Multi-line process spawned with ID: 2\n" +
		"[CALLBACK] Process 2 STDOUT: Line 1\nLine 2\nLine 3\n\n" +
		"[CALLBACK] Process 2 EXIT: 0\n" +
		"✓ Multi-line process finished\n" +
		"--- Test 3: Error Process Callback ---\n" +
		"✓ Error process spawned with ID: 3\n" +
		"[CALLBACK] Process 3 STDERR: ls: cannot access '/nonexistent/directory': No such file or directory\n\n" +
		"[CALLBACK] Process 3 EXIT: 2\n" +
		"✓ Error process finished with exit code: 2\n" +
		"=== CALLBACK DEMO COMPLETE ===\n" +
		"The [CALLBACK] lines above show C runtime calling into Osprey!\n"
}

func getAdventureGameExpectedOutput() string {
	return "🏰 Welcome to the Mystical Castle Adventure! 🏰\n" +
		"You stand before an ancient castle shrouded in mystery...\n\n" +
		"⚔️  Your Quest Begins! ⚔️\n\n" +
		"You are Novice Adventurer (Level 1)\n" +
		"Health: 100 ❤️  | Gold: 50 💰\n\n" +
		"🚪 Room 1: You enter the Grand Entrance Hall with marble columns\n" +
		"You find 10 gold coins! Total: 60 💰\n\n" +
		"📚 Room 2: You discover a dusty Library filled with ancient tomes\n" +
		"You find 25 gold coins and acquire a mysterious key! 🗝️\n" +
		"Total gold: 85 💰\n\n" +
		"⚔️  Room 3: You enter the Armory containing gleaming weapons\n" +
		"You acquire a gleaming sword! ⚔️\n" +
		"Your combat prowess has increased dramatically!\n\n" +
		"🐉 BOSS BATTLE: Ancient Dragon Appears! 🐉\n" +
		"The ground trembles as a massive dragon blocks your path!\n\n" +
		"Enemy: Ancient Dragon\n" +
		"Enemy Health: 120 ❤️\n" +
		"Your attack power: 60 ⚔️\n\n" +
		"⚡ BATTLE COMMENCES! ⚡\n" +
		"You need 2.0 successful attacks to defeat the Ancient Dragon!\n\n" +
		"🥊 Round 1: You strike for 60 damage!\n" +
		"Dragon health remaining: 60\n\n" +
		"🥊 Round 2: Another powerful blow for 60 damage!\n" +
		"Dragon health remaining: 0\n\n" +
		"🥊 FINAL ROUND: You deliver the finishing blow!\n" +
		"Critical hit for 0 damage!\n\n" +
		"🎉 VICTORY! 🎉\n" +
		"The Ancient Dragon has been defeated!\n" +
		"You gain 200 gold coins as reward!\n" +
		"Total gold: 285 💰\n\n" +
		"📈 LEVEL UP! 📈\n" +
		"Previous: Novice Adventurer (Level 1)\n" +
		"New: Brave Explorer (Level 2)\n\n" +
		"🏆 Room 4: You enter the Treasure Chamber sparkling with gold\n" +
		"You discover the legendary treasure chest!\n" +
		"Inside: 100 gold coins! 💎\n" +
		"Your final wealth: 385 💰\n\n" +
		"🎭 QUEST COMPLETE! 🎭\n" +
		"Congratulations, Brave Explorer!\n" +
		"You have conquered the Mystical Castle!\n" +
		"Final Stats:\n" +
		"  - Level: 2\n" +
		"  - Monsters Defeated: 1\n" +
		"  - Gold Collected: 385 💰\n" +
		"  - Artifacts: Sword ⚔️ & Key 🗝️\n\n" +
		"🌟 Your legend will be remembered forever! 🌟\n" +
		"Thanks for playing the Osprey Adventure Game!\n"
}
