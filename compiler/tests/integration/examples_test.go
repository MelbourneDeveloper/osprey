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
			testName = strings.ReplaceAll(testName, string(filepath.Separator), "_")

			t.Run(testName, func(t *testing.T) {
				// Try to read from .expectedoutput file first
				expectedOutputPath := path + ".expectedoutput"
				if expectedContent, err := os.ReadFile(expectedOutputPath); err == nil {
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

// getExpectedOutputs returns the map of expected outputs for each test file.
func getExpectedOutputs() map[string]string {
	return map[string]string{
		"hello.osp": "Hello, World!\nHello from function!\n",
		"interpolation_math.osp": "Next year you'll be 26\nLast year you were 24\n" +
			"Double your age: 50\nHalf your age: 12\n",
		"interpolation_comprehensive.osp": "Hello Alice!\nYou are 25 years old\n" +
			"Your score is 95 points\nNext year you'll be 26\n" +
			"Double your score: 190\nAlice (25) scored 95/100\n",
		"working_basics.osp": "x = 42\nname = Alice\ndouble(21) = 42\n" +
			"greeting = Hello\n10 + 5 = 15\n6 * 7 = 42\nmatch 42 = 1\n",
		"simple_types.osp":        "Type definitions compiled successfully\nred\nworking\n",
		"result_type_example.osp": "Result type defined successfully\n42\n",
		"simple_input.osp": "Greeting code: 1\nNumber result: 999\n" +
			"Unknown code: 0\nSmall number: 100\n",
		"pattern_matching_basics.osp": "Number analysis:\n0 is Zero\n" +
			"42 is The answer to everything!\n7 is Some other number\n" +
			"\nEven number check:\n42 is even: 0\n7 is even: 0\n2 is even: 1\n" +
			"\nScore categories:\nScore 100: Perfect!\n" +
			"Score 85: Very Good\nScore 50: Needs Improvement\n" +
			"Nested: Both zero\n",
		"safe_arithmetic_demo.osp": "=== Type-Safe Arithmetic Demo ===\n" +
			"Future: All operators return Result<T, Error>\n\n10 / 2 = 5\n" +
			"Error: Cannot divide 15 by 0!\n20 / 4 = 5\n\n" +
			"✅ No panics! All division operations handled safely\n" +
			"🔮 Future: Built-in Result<T, E> types for all fallible operations\n",
		"script_style_working.osp": "Script starting...\nFactorial computed!\n",
		"calculator_fixed.osp": "=== Osprey Interactive Calculator ===\n" +
			"Enter a number:\nComputing operations...\nMany!\nAll computations complete!\n",
		"math_calculator_fixed.osp": "=== Advanced Math Calculator ===\n" +
			"Enter base number:\nEnter multiplier:\nComputing advanced operations...\n" +
			"=== Results ===\nBase cubed:\n125\nFactorial approximation:\n15\n" +
			"Fibonacci approximation:\n15\nComplex formula result:\n2\n" +
			"=== Calculator Complete ===\n",
		"space_trader.osp":   getSpaceTraderExpectedOutput(),
		"adventure_game.osp": getAdventureGameExpectedOutput(),
		"basic_iterator_test.osp": "=== Basic Iterator Test ===\n" +
			"Test 1: Simple pipe with double\n10\n\n" +
			"Test 2: Range 1 to 5 with double function\n\n" +
			"Test 3: Range 1 to 5 with print\n1\n2\n3\n4\n5\n\n" +
			"Test 4: Range 1 to 4 with square function\n\n" +
			"Test 5: Chained pipe operations\n400\n\n" +
			"Test 6: 3 -> addFive -> double -> print\n16\n\n" +
			"Test 7: Range 0 to 3 with addFive\n\n" +
			"Test 8: Multiple small ranges\n1\n2\n10\n11\n=== Test Complete ===\n",
		"comprehensive_iterators.osp": "=== Comprehensive Iterator Test ===\n" +
			"Test 1: Count 1 to 10\n1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n\n" +
			"Test 2: Count 10 to 15\n10\n11\n12\n13\n14\n15\n\n" +
			"Test 3: Count 0 to 5\n0\n1\n2\n3\n4\n5\n\n" +
			"Test 4: Count -3 to 3\n-3\n-2\n-1\n0\n1\n2\n3\n\n" +
			"Test 5: Count 100 to 105\n100\n101\n102\n103\n104\n105\n\n" +
			"Test 6: Single value 42\n42\n\nTest 7: Empty range (5 to 5)\n" +
			"=== All Tests Complete ===\n",
		"functional_showcase.osp": "=== Functional Programming Showcase ===\n" +
			"Example 1: Basic range iteration\n1\n2\n3\n4\n5\n" +
			"Example 2: Single value pipe operations\n18\n" +
			"Example 3: Business logic pipeline\n88\n" +
			"Example 4: Range forEach\n42\n43\n44\n" +
			"Example 5: Small range\n10\n11\n12\n" +
			"Example 6: Range 0 to 4\n0\n1\n2\n3\n4\n" +
			"Example 7: Fold operations\n15\n42\n" +
			"Example 8: Chained single value operations\n21\n" +
			"Example 9: Conditional operations\n1\n0\n=== Showcase Complete ===\n",
		"explicit_any_allowed.osp": "Explicit any return type works\n" +
			"getDynamicValue() = 42\n" +
			"processAnyValue(5) = 15\n",
		"explicit_any_simple.osp": "Explicit any return type works\n",
		"functional_iterators.osp": "=== Functional Iterator Examples ===\n" +
			"1. Basic forEach:\n1\n2\n3\n4\n" +
			"2. Single value transformations:\n10\n9\n" +
			"3. Different ranges:\n10\n11\n12\n0\n1\n2\n" +
			"4. Fold operations:\n15\n125\n" +
			"5. Chained single value operations:\n16\n" +
			"=== Examples Complete ===\n",
		"documentation_test.osp": "Testing documentation\n1\n2\n3\n4\n",
		// Boolean examples that work with current parser
		"comparison_test.osp": "1\n",    // Prints result of 5 > 3
		"equality_test.osp":   "true\n", // Prints result of isEqual(5, 5)
		"comprehensive_bool_test.osp": "=== Boolean Test ===\nFunction returning true:\ntrue\n" +
			"Function returning false:\nfalse\nBoolean literals:\nfalse\ntrue\nComparisons:\ntrue\ntrue\ntrue\ntrue\ntrue\n",
		"full_bool_test.osp": "=== Boolean Test Results ===\n5 > 3:\ntrue\n" +
			"10 == 10:\ntrue\ntrue literal:\ntrue\nfalse literal:\nfalse\n",
		"modulo_test.osp": "true\nfalse\n",
		// Compilation-only tests (no output expected)
		"basic.osp": "Basic test results:\nx = 42\ntestGood(10) = 10\n" +
			"getIntResult() = 42\ngetStringResult() = asd\naddOne(5) = 6\n",
		"comprehensive.osp": "=== Comprehensive Osprey Demo ===\n" +
			"Student Alice scored 95 points\n" +
			"Doubled score: 190\n" +
			"Excellent!\n" +
			"Status: System operational\n" +
			"Double of 42: 84\n" +
			"Student Bob scored 92 points\n" +
			"=== Demo Complete ===\n",
		"debug_module.osp": "Debug module test:\nsimple() = 42\n",
		"function.osp":     "Function test:\nadd(3, 7) = 10\nadd(10, 20) = 30\n",
		"function_composition_test.osp": "=== Function Composition Test ===\n" +
			"Testing function composition...\n" +
			"Starting value: 10\n" +
			"After double: 20\n" +
			"After triple: 30\n" +
			"After add5: 15\n" +
			"square(5) = 25\n" +
			"Function composition working correctly!\n" +
			"=== Function Composition Test Complete ===\n",
		"minimal_test.osp": "Minimal test:\nx = 5\n",
		"simple.osp":       "Simple test:\nx = 42\ngreeting = hello\n",
		// Constraint validation test files
		"constraint_validation_test.osp": "=== CONSTRAINT VALIDATION WITH FAILURE DETECTION ===\n" +
			"Test 1: Valid Person construction\nResult: 1\nSuccess: 1\nFailure: 0\n\n" +
			"Test 2: Invalid Person - empty name constraint violation\nResult: -1\nSuccess: 0\nFailure: 1\n" +
			"Expected: Failure = 1 (constraint violation)\n\n" +
			"Test 3: Invalid Person - zero age constraint violation\nResult: -1\nSuccess: 0\nFailure: 1\n" +
			"Expected: Failure = 1 (constraint violation)\n\n" +
			"Test 4: Valid Product construction\nResult: 1\nSuccess: 1\nFailure: 0\n\n" +
			"Test 5: Invalid Product - zero price constraint violation\nResult: -1\nSuccess: 0\nFailure: 1\n" +
			"Expected: Failure = 1 (constraint violation)\n\n" +
			"Test 6: Multiple constraint violations\nResult: -1\nSuccess: 0\nFailure: 1\n" +
			"Expected: Failure = 1 (multiple constraint violations)\n\n" +
			"=== CONSTRAINT VALIDATION TESTS COMPLETE ===\n" +
			"This test demonstrates that WHERE constraints work correctly:\n" +
			"✅ Valid constructions return 1 (success)\n" +
			"❌ Invalid constructions return -1 (constraint violation)\n" +
			"✅ notEmpty constraint rejects empty strings\n" +
			"✅ validAge constraint rejects zero age\n" +
			"✅ positive constraint rejects zero prices\n" +
			"✅ Multiple violations are properly detected\n\n" +
			"FUTURE: Should return Result<T, ConstraintError> types for type safety.\n",
		"working_constraint_test.osp": "=== CONSTRAINT FUNCTION VERIFICATION ===\n" +
			"Testing notEmpty function:\nnotEmpty(\"\") should be false:\nfalse\n" +
			"notEmpty(\"alice\") should be true:\ntrue\nTesting isPositive function:\n" +
			"isPositive(0) should be false:\nfalse\nisPositive(100) should be true:\ntrue\n" +
			"Testing validAge function:\nvalidAge(0) should be false:\nfalse\n" +
			"validAge(25) should be true:\ntrue\nTesting validEmail function:\n" +
			"validEmail(\"\") should be false:\nfalse\nvalidEmail(\"test@email.com\") should be true:\ntrue\n" +
			"=== BASIC TYPE CONSTRUCTION TEST ===\n" +
			"✅ Creating Person:\nPerson created successfully\n" +
			"✅ Creating User:\nUser created successfully\n" +
			"✅ Creating Product:\nProduct created successfully\n" +
			"=== CONSTRAINT FUNCTIONS AND TYPE CONSTRUCTION COMPLETE ===\n" +
			"📝 Note: Type-level validation with Result types to be implemented later\n",
		"proper_validation_test.osp": "Testing validation functions:\nfalse\ntrue\nfalse\ntrue\ntrue\nfalse\n",
		"match_type_mismatch.osp":    "none\n",
		// Website examples
		"website_hero_example.osp": "x = 42\nname = Alice\nResult: The answer!\n",
		"website_type_safe_example.osp": "Testing functions:\nZero\nThe answer!\nSomething else\n" +
			"5 doubled is 10\n10 squared is 100\n",
		"website_string_interpolation_example.osp": "Hello Alice!\nNext year you'll be 26\n" +
			"Double score: 190\nAlice (25) scored 95/100\n",
		"website_pattern_matching_grade_example.osp": "Grade for 100: Perfect!\nGrade for 95: Excellent\n" +
			"Grade for 85: Very Good\nGrade for 75: Good\nGrade for 50: Needs Improvement\n",
		"website_functional_programming_example.osp": "100\nRange operations:\n1\n2\n3\n4\n5\n6\n7\n8\n9\n",
		"website_fiber_isolation_example.osp": "Fiber 1 result: 1\nFiber 2 result: 2\n" +
			"Processed 10 items\n=== Fiber Example Complete ===\n",
		// Block statement examples
		"block_statements_basic.osp": "=== Basic Block Statements Test ===\n" +
			"Test 1 - Simple block: 0\n" +
			"Test 2 - Block computation: 0\n" +
			"Test 3 - Multiple statements: 0\n" +
			"=== Basic Block Statements Complete ===\n",
		"block_statements_advanced.osp": "=== Advanced Block Statements Test ===\n" +
			"Test 1 - Function block: 0\n" +
			"Test 2 - Nested with shadowing: 0\n" +
			"Test 3 - Block with match: 0\n" +
			"Test 4 - Complex function: 0\n" +
			"=== Advanced Block Statements Complete ===\n",
		"process_spawn_basic.osp": "Hello World\n" +
			"=== Basic Process Spawning Test ===\n" +
			"Process result: 0\n" +
			"=== Test Complete ===\n",
		"process_spawn_fiber.osp": "=== Process Spawning in Fibers ===\n" +
			"Process result: 0\n" +
			"=== Fiber Test Complete ===\n",
		"simple_process_test.osp": "Testing simple process spawn...\n" +
			"Process spawned successfully\n" +
			"[STDOUT] Process 1: hello\n\n" +
			"[EXIT] Process 1 exited with code: 0\n" +
			"Process finished\n" +
			"Test complete\n",
		"async_process_management.osp": "=== Async Process Management Demo ===\n" +
			"--- Test 1: Basic Process Spawning ---\n" +
			"✓ Process spawned successfully\n" +
			"[STDOUT] Process 1: Hello from async process!\n\n" +
			"[EXIT] Process 1 exited with code: 0\n" +
			"✓ Process completed successfully\n" +
			"✓ Process resources cleaned up\n" +
			"--- Test 2: Another Process ---\n" +
			"Process 2 spawned successfully\n" +
			"[STDOUT] Process 2: Process 2 output\n\n" +
			"[EXIT] Process 2 exited with code: 0\n" +
			"Process 2 finished\n" +
			"--- Test 3: Error Handling ---\n" +
			"[EXIT] Process 3 exited with code: 1\n" +
			"Error process returned non-zero exit code\n" +
			"=== Async Process Management Demo Complete ===\n" +
			"Note: Process output appears via C runtime callbacks during execution\n",
		"callback_stdout_demo.osp": getCallbackStdoutDemoExpectedOutput(),
		"process_spawn_workflow.osp": "Step 1\n" +
			"Step 2\n" +
			"=== Process Spawning Workflow ===\n" +
			"Step 1 result: 0\n" +
			"Step 2 result: 0\n" +
			"Fiber result: 1\n" +
			"=== Workflow Complete ===\n",
		"process_spawn_simple.osp": "Hello from process!\n" +
			"2025\n" +
			"=== Simple Process Spawning ===\n" +
			"Result 1: 0\n" +
			"Result 2: 0\n" +
			"Fiber ID: 1\n" +
			"=== Test Complete ===\n",
		"result_type_workflow.osp": "=== Result Type Workflow Test ===\n\n" +
			"Length: \n5\n\n\n" +
			"Contains 'ell': \n1\n\n\n" +
			"Contains 'xyz': \n0\n\n\n",
		"file_io_json_workflow.osp": "=== File I/O Workflow Test ===\n" +
			"-- Step 1: Writing to file --\n" +
			"File written successfully!\n" +
			"-- Step 2: Reading from file --\n" +
			"Read successful!\n" +
			"=== Test Complete ===\n",
		"string_utils_combined.osp": "=== String Utils Test ===\n\nOriginal: \\\nhello world\n\"\n\n" +
			"Length: \n11\n\n\nContains 'world': \n1\n\n\nContains 'galaxy': \n0\n\n\n" +
			"Substring(6, 11): \\\nworld\n\"\n\nSubstring(0, 20): \\\nhello world\n\"\n\n" +
			"=== Test Complete ===\n\n",
		"list_and_process.osp": "=== Array Access Test ===\n\nCreated array with 3 commands\n\n" +
			"Testing array access with pattern matching:\n\n✅ commands[0] = \\\necho hello\n\"\n\n" +
			"✅ commands[1] = \\\necho world\n\"\n\n✅ commands[2] = \\\necho test\n\"\n\n" +
			"Testing out-of-bounds access:\n\n✅ Correctly caught out-of-bounds: commands[5] -> Error\n\n" +
			"=== Array Test Complete ===\n\n",
		// Effects examples
		"algebraic_effects.osp": "Pure function result: 42\n🎉 BASIC TEST COMPLETE! 🎉",
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
		t.Errorf("Output mismatch for %s:\nExpected: %q\nGot:      %q", filePath, expectedOutput, actualOutput)
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
			"[CALLBACK] Process 2 STDOUT: Line 1\\\nLine 2\\\nLine 3\\\n\n" +
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
		"[CALLBACK] Process 2 STDOUT: Line 1\\\nLine 2\\\nLine 3\\\n\n" +
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

func getSpaceTraderExpectedOutput() string {
	return "🌌 Welcome to the Galactic Trade Network! 🌌\n" +
		"You are Captain Alex, commander of the starship Osprey-7\n" +
		"Your mission: Build a trading empire across the galaxy!\n\n" +
		"🛸 MISSION BRIEFING 🛸\n" +
		"Ship: Osprey-7 Starfreighter\n" +
		"Fuel: 100% ⛽\n" +
		"Credits: 1000 💰\n" +
		"Cargo Space: 0/50 📦\n" +
		"Reputation: Unknown Trader\n\n" +
		"🌍 GALACTIC TRADING SIMULATION 🌍\n\n" +
		"📍 Arriving at Nebula Prime\n" +
		"This planet specializes in: Quantum Crystals\n" +
		"Market price: 50 credits per unit\n" +
		"Purchasing 10 units of Quantum Crystals\n" +
		"Total cost: 500 credits\n" +
		"Remaining credits: 500 💰\n" +
		"Cargo: 10/50 📦\n\n" +
		"🚀 Traveling to Crystal Moon...\n" +
		"Fuel consumed: 20%\n" +
		"Current fuel: 80% ⛽\n\n" +
		"📍 Arrived at Crystal Moon\n" +
		"Local specialty: Space Metal\n" +
		"Market price: 25 credits per unit\n" +
		"Selling 10 units of Quantum Crystals\n" +
		"Sale price: 75 credits per unit\n" +
		"Revenue: 750 credits 💰\n" +
		"New balance: 1250 credits\n" +
		"Cargo space freed: 0/50 📦\n\n" +
		"Purchasing 15 units of Space Metal\n" +
		"Cost: 375 credits\n" +
		"Remaining credits: 875 💰\n\n" +
		"🚀 Long-range jump to Trade Station Alpha\n" +
		"Fuel consumed: 30%\n" +
		"Current fuel: 50% ⛽\n\n" +
		"📍 Docking at Trade Station Alpha\n" +
		"This is the galaxy's premier trading hub!\n" +
		"Selling 15 units of Space Metal\n" +
		"Hub premium price: 55 credits per unit\n" +
		"Major revenue: 825 credits! 💰\n" +
		"New balance: 1700 credits\n\n" +
		"📈 TRADING RESULTS 📈\n" +
		"Starting credits: 1000\n" +
		"Final credits: 1700\n" +
		"Total profit: 700 credits! 💰\n" +
		"Planets visited: 3\n" +
		"New reputation: Novice Merchant\n\n" +
		"🛸 SHIP STATUS REPORT 🛸\n" +
		"Fuel level: 50% (Fair)\n" +
		"Cargo bay: 0/50 units\n" +
		"Ship condition: Operational\n\n" +
		"📊 ADVANCED ANALYTICS 📊\n" +
		"Fuel efficiency: 16% per planet\n" +
		"Profit per planet: 233 credits\n" +
		"Projected wealth (if doubled): 3400 credits\n\n" +
		"🏆 MISSION COMPLETE! 🏆\n" +
		"Congratulations, Captain Novice Merchant!\n" +
		"You have successfully established trade routes across the galaxy!\n\n" +
		"Next objectives:\n" +
		"  ⭐ Explore more distant sectors\n" +
		"  ⭐ Upgrade ship cargo capacity\n" +
		"  ⭐ Establish permanent trade agreements\n" +
		"  ⭐ Recruit specialized crew members\n\n" +
		"🌟 Your trading empire awaits! 🌟\n" +
		"End of Galactic Trade Simulation\n" +
		"Thank you for playing Osprey Space Trader!\n"
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
		"You need 2 successful attacks to defeat the Ancient Dragon!\n\n" +
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
