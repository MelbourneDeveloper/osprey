package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestRootLevelExamples tests ONLY the root-level examples in examples/tested/ (NOT subdirectories).
// Subdirectories like http/, websox/, fiber/ have their own separate test functions.
func TestRootLevelExamples(t *testing.T) {
	checkLLVMTools(t)

	examplesDir := "../../examples/tested"
	runTestExamples(t, examplesDir, getExpectedOutputs())
}

func runTestExamples(t *testing.T, examplesDir string, expectedOutputs map[string]string) {
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	processExampleFiles(t, examplesDir, entries, expectedOutputs)
}

// TestLanguageFeatures tests core language feature examples for better test explorer support.
func TestLanguageFeatures(t *testing.T) {
	checkLLVMTools(t)

	examplesDir := "../../examples/tested"
	expectedOutputs := getExpectedOutputs()

	languageExamples := []string{
		"hello.osp",
		"basic.osp",
		"simple.osp",
		"function.osp",
		"working_basics.osp",
		"simple_types.osp",
		"result_type_example.osp",
		"interpolation_math.osp",
		"interpolation_comprehensive.osp",
		"pattern_matching_basics.osp",
		"comparison_test.osp",
		"equality_test.osp",
		"comprehensive_bool_test.osp",
		"full_bool_test.osp",
		"modulo_test.osp",
		"block_statements_basic.osp",
		"block_statements_advanced.osp",
	}

	for _, exampleFile := range languageExamples {
		testName := strings.TrimSuffix(exampleFile, ".osp")
		t.Run(testName, func(t *testing.T) {
			filePath := filepath.Join(examplesDir, exampleFile)
			expectedOutput, exists := expectedOutputs[exampleFile]
			if !exists {
				t.Fatalf("❌ MISSING expected output for %s!", exampleFile)
			}
			testExampleFile(t, filePath, expectedOutput)
		})
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
			"Score 85: Very Good\nScore 50: Needs Improvement\n",
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
			"Function returning false:\nfalse\nBoolean literals:\nfalse\ntrue\nComparisons:\n1\n1\n1\n1\n1\n",
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
			"=== CONSTRAINT VALIDATION TEST ===\nTesting current constraint implementation:\n" +
			"✅ Valid Person (returns 1):\n1\n❌ Invalid Person - empty name (returns -1):\n-1\n" +
			"❌ Invalid Person - zero age (returns -1):\n-1\n✅ Valid Product (returns 1):\n1\n" +
			"❌ Invalid Product - empty name (returns -1):\n-1\n❌ Invalid Product - zero price (returns -1):\n-1\n" +
			"=== TYPE SAFETY ISSUES ===\nPROBLEM: These variables have type 'any' instead of Result<T, E>:\n" +
			"invalidPersonAge should be Result<Person, ConstraintError>\nBut we can treat it as an integer:\n-1\n" +
			"SOLUTION NEEDED: Proper Result<T, E> types\nThen we would need pattern matching:\n" +
			"match invalidPersonAge {\n  Ok { value } => use the person\n  Err { error } => handle constraint violation\n}\n" +
			"=== CONSTRAINT TESTS COMPLETE ===\n=== COMPREHENSIVE WHERE CONSTRAINT TESTS ===\n" +
			"PERSON CONSTRAINT TESTS:\n✅ Valid Person (should return 1):\n1\n" +
			"❌ Invalid Person - empty name (should return -1):\n-1\n❌ Invalid Person - zero age (should return -1):\n-1\n" +
			"USER CONSTRAINT TESTS:\n✅ Valid User (should return 1):\n1\n" +
			"❌ Invalid User - empty username (should return -1):\n-1\n❌ Invalid User - empty email (should return -1):\n-1\n" +
			"❌ Invalid User - zero userId (should return -1):\n-1\nPRODUCT CONSTRAINT TESTS:\n" +
			"✅ Valid Product (should return 1):\n1\n❌ Invalid Product - empty name (should return -1):\n-1\n" +
			"❌ Invalid Product - zero price (should return -1):\n-1\n=== WHERE CONSTRAINT VALIDATION COMPLETE ===\n",
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
		// TODO: REMOVED BULLSHIT "safe arithmetic" example - was just manual checking, not actually safe!
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
		"http_advanced_example.osp": "=== Advanced HTTP Test ===\n" +
			"Creating HTTP server on port 8080...\n" +
			"Server created with ID: 1\n" +
			"Starting server listener...\n" +
			"Server listening on http://127.0.0.1:8080\n" +
			"=== Creating Multiple Clients ===\n" +
			"Creating client 1...\n" +
			"Client 1 created with ID: 2\n" +
			"Creating client 2...\n" +
			"Client 2 created with ID: 3\n" +
			"Creating client 3...\n" +
			"Client 3 created with ID: 4\n" +
			"=== Concurrent Requests ===\n" +
			"Client 1: GET /api/users\n" +
			"Client 1 GET result: -7\n" +
			"Client 2: POST /api/posts\n" +
			"Client 2 POST result: -7\n" +
			"Client 3: GET /api/health\n" +
			"Client 3 health check: -7\n" +
			"=== API Versioning ===\n" +
			"Client 1: GET /v1/users\n" +
			"v1 API result: -7\n" +
			"Client 2: GET /v2/users\n" +
			"v2 API result: -7\n" +
			"=== Content Types ===\n" +
			"Client 1: POST /api/upload (XML)\n" +
			"XML POST result: -7\n" +
			"Client 2: PUT /api/config (YAML)\n" +
			"YAML PUT result: -7\n" +
			"Client 3: POST /api/data (Form)\n" +
			"Form POST result: -7\n" +
			"=== Authentication ===\n" +
			"Client 1: POST /auth/login\n" +
			"Login result: -5\n" +
			"Client 2: GET /protected (with token)\n" +
			"Protected GET result: -5\n" +
			"Client 3: DELETE /auth/logout\n" +
			"Logout result: -5\n" +
			"=== Error Scenarios ===\n" +
			"Client 1: GET /nonexistent\n" +
			"404 test result: -5\n" +
			"Client 2: POST /api/invalid (bad JSON)\n" +
			"Bad JSON result: -5\n" +
			"Stopping server...\n" +
			"Server stopped with result: 0\n" +
			"=== Advanced HTTP Test Complete ===\n",
	}
}

// processExampleFiles processes ONLY .osp files in the given directory (SKIPS subdirectories).
func processExampleFiles(t *testing.T, examplesDir string, entries []os.DirEntry, expectedOutputs map[string]string) {
	// Test each .osp file - SKIP ALL DIRECTORIES
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".osp") {
			// Create a meaningful test name by removing .osp extension
			testName := strings.TrimSuffix(entry.Name(), ".osp")

			t.Run(testName, func(t *testing.T) {
				filePath := filepath.Join(examplesDir, entry.Name())
				expectedOutput, exists := expectedOutputs[entry.Name()]
				if !exists {
					t.Fatalf("❌ MISSING expected output for %s!\n"+
						"🚨 ALL EXAMPLES MUST HAVE VERIFIED OUTPUT!\n"+
						"🚨 NO COMPILATION-ONLY TESTS ALLOWED!\n"+
						"🚨 RUN THE EXAMPLE AND ADD THE ACTUAL OUTPUT TO expectedOutputs MAP!\n"+
						"🚨 Use: ../../osprey %s --run\n"+
						"🚨 Then copy the output to the expectedOutputs map!", entry.Name(), entry.Name())
				}
				testExampleFile(t, filePath, expectedOutput)
			})
		}
	}
}

// testExampleFile tests a single example file.
func testExampleFile(t *testing.T, filePath, expectedOutput string) {
	t.Helper()

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", filePath, err)
	}

	source := string(content)

	// Try to compile first
	_, err = codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("Failed to compile %s: %v", filePath, err)
	}

	// ERROR: ALL EXAMPLES MUST HAVE VERIFIED OUTPUT!
	if expectedOutput == "" {
		t.Fatalf("❌ MISSING EXPECTED OUTPUT FOR %s!\n"+
			"🚨 ALL EXAMPLES MUST HAVE VERIFIED OUTPUT!\n"+
			"🚨 NO COMPILATION-ONLY TESTS ALLOWED!\n"+
			"🚨 RUN THE EXAMPLE AND ADD THE ACTUAL OUTPUT TO expectedOutputs MAP!\n"+
			"🚨 Use: ../../osprey %s --run\n"+
			"🚨 Then copy the output to the expectedOutputs map!",
			filepath.Base(filePath), filepath.Base(filePath))
	}

	// Execute and verify output
	output, err := captureJITOutput(source)
	if err != nil {
		// If JIT execution fails due to missing tools, fail the test
		if strings.Contains(err.Error(), "LLVM tools not found") ||
			strings.Contains(err.Error(), "no suitable compiler found") {
			t.Fatalf("❌ LLVM TOOLS NOT FOUND - TEST FAILED for %s: %v", filePath, err)
		}
		t.Fatalf("Failed to execute %s: %v", filePath, err)
	}

	if output != expectedOutput {
		t.Errorf("Output mismatch for %s:\nExpected: %q\nGot:      %q", filePath, expectedOutput, output)
	}

	t.Logf("✅ Example %s compiled and executed successfully", filepath.Base(filePath))
}

// Helper functions for expected outputs.
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
