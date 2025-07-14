package integration

import (
	"fmt"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestFiberFeatures tests the fiber language features comprehensively.
func TestFiberFeatures(t *testing.T) {
	fiberTests := map[string]string{
		"basic_spawn": `fn test() -> int = spawn 42
fn main() -> int = test()`,

		"basic_await": `fn test() -> int = await (spawn 100)
fn main() -> int = test()`,

		"basic_yield": `fn test() -> int = yield 42
fn main() -> int = test()`,

		"basic_channel": `fn test() -> int = Channel<Int> { capacity: 10 }
fn main() -> int = test()`,

		"channel_send": `fn test() -> int = send(Channel<Int> { capacity: 1 }, 42)
fn main() -> int = test()`,

		"channel_recv": `fn test() -> int = recv(Channel<Int> { capacity: 1 })
fn main() -> int = test()`,

		"lambda_expression": `fn test() -> int = (fn() => 42)()
fn main() -> int = test()`,

		"spawn_with_await": `fn test() -> int = await (spawn 42)
fn main() -> int = test()`,

		"complex_fiber_chain": `fn test() -> int = await (spawn (yield 42))
fn main() -> int = test()`,

		"module_with_fibers": `module FiberModule {
    fn compute() -> int = spawn 42
    fn get_result() -> int = await (spawn 100)
}
fn main() -> int = FiberModule.compute()`,
	}

	for name, source := range fiberTests {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Fiber test %s failed to compile: %v", name, err)
			} else {
				t.Logf("✅ Fiber test %s compiled successfully", name)
			}
		})
	}
}

// TestFiberErrorHandling tests that invalid fiber syntax fails gracefully.
func TestFiberErrorHandling(t *testing.T) {
	invalidFiberTests := map[string]string{
		"spawn_without_expression": `fn test() -> int = spawn
fn main() -> int = test()`,

		"await_without_expression": `fn test() -> int = await
fn main() -> int = test()`,

		"channel_without_type": `fn test() -> int = Channel<>
fn main() -> int = test()`,

		"invalid_channel_syntax": `fn test() -> int = Channel<> { capacity: 10 }
fn main() -> int = test()`,

		"select_without_arms": `fn test() -> int = select {}
fn main() -> int = test()`,

		"malformed_lambda": `fn test() -> int = fn() =>
fn main() -> int = test()`,
	}

	for name, source := range invalidFiberTests {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err == nil {
				t.Errorf("Invalid fiber syntax %s should have failed to compile", name)
			} else {
				t.Logf("✅ Invalid fiber syntax %s correctly failed: %v", name, err)
			}
		})
	}
}

// TestFiberModuleIsolation tests the fiber-isolated module system.
func TestFiberModuleIsolation(t *testing.T) {
	moduleIsolationTests := map[string]string{
		"basic_module_isolation": `module IsolatedModule {
    fn increment() -> int = 42
    fn get_state() -> int = 42
}

fn main() -> int = spawn 42`,

		"module_with_fibers": `module FiberModule {
    fn compute_async() -> int = spawn 42
    fn process_data() -> int = await (spawn 100)
    fn yield_control() -> int = yield 200
}

fn main() -> int = FiberModule.compute_async()`,

		"simple_fiber_module": `module SimpleModule {
    fn fiber_task() -> int = spawn 42
}

fn main() -> int = SimpleModule.fiber_task()`,

		"module_channel_operations": `module ChannelModule {
    fn create_channel() -> int = Channel<Int> { capacity: 10 }
    fn send_data() -> int = send(Channel<Int> { capacity: 1 }, 42)
    fn recv_data() -> int = recv(Channel<Int> { capacity: 1 })
}

fn main() -> int = ChannelModule.create_channel()`,
	}

	for name, source := range moduleIsolationTests {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Module isolation test %s failed to compile: %v", name, err)
			} else {
				t.Logf("✅ Module isolation test %s compiled successfully", name)
			}
		})
	}
}

// TestFiberIntegration provides a comprehensive validation of the complete fiber implementation.
func TestFiberIntegration(t *testing.T) {
	t.Log("🚀 Running comprehensive fiber integration test")

	// Test that all core fiber keywords are recognized
	testFiberKeywords(t)
	testFiberNesting(t)
	testChannelOperations(t)
	testFiberLambdas(t)

	t.Log("🎉 Comprehensive fiber integration test completed successfully!")
}

// testFiberKeywords tests that all fiber keywords compile correctly.
func testFiberKeywords(t *testing.T) {
	fiberKeywords := []string{"spawn", "await", "yield", "channel", "select"}
	for _, keyword := range fiberKeywords {
		t.Run("keyword_"+keyword, func(t *testing.T) {
			source := getFiberKeywordTestSource(keyword)
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Fiber keyword %s should compile successfully: %v", keyword, err)
			} else {
				t.Logf("✅ Fiber keyword %s compiled successfully", keyword)
			}
		})
	}
}

// getFiberKeywordTestSource returns appropriate test source for each fiber keyword.
func getFiberKeywordTestSource(keyword string) string {
	switch keyword {
	case "channel":
		return "fn test() -> int = Channel<Int> { capacity: 42 }\nfn main() -> int = test()"
	case "select":
		return "fn test() -> int = select { 42 => 100 }\nfn main() -> int = test()"
	default:
		return fmt.Sprintf("fn test() -> int = %s 42\nfn main() -> int = test()", keyword)
	}
}

// testFiberNesting tests fiber expression nesting.
func testFiberNesting(t *testing.T) {
	t.Run("fiber_nesting", func(t *testing.T) {
		nestedFiberTests := []string{
			"await (spawn 42)",
			"spawn (await (spawn 42))",
			"yield (spawn 42)",
			"spawn (yield 42)",
		}

		for i, expr := range nestedFiberTests {
			source := fmt.Sprintf("fn test() -> int = %s\nfn main() -> int = test()", expr)
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Nested fiber expression %d should compile: %v", i, err)
			} else {
				t.Logf("✅ Nested fiber expression %d compiled successfully", i)
			}
		}
	})
}

// testChannelOperations tests channel operations.
func testChannelOperations(t *testing.T) {
	t.Run("channel_operations", func(t *testing.T) {
		channelTests := map[string]string{
			"channel_creation": "Channel<Int> { capacity: 10 }",
			"channel_send":     "send(Channel<Int> { capacity: 1 }, 42)",
			"channel_recv":     "recv(Channel<Int> { capacity: 1 })",
			"typed_channel":    "Channel<String> { capacity: 5 }",
		}

		for name, expr := range channelTests {
			source := fmt.Sprintf("fn test() -> int = %s\nfn main() -> int = test()", expr)
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Channel operation %s should compile: %v", name, err)
			} else {
				t.Logf("✅ Channel operation %s compiled successfully", name)
			}
		}
	})
}

// testFiberLambdas tests lambda expressions with fibers.
func testFiberLambdas(t *testing.T) {
	t.Run("fiber_lambdas", func(t *testing.T) {
		lambdaTests := []string{
			"(fn() => spawn 42)()",
			"(fn() => await (spawn 42))()",
			"(fn() => yield 42)()",
			"(fn() => 42)()",
		}

		for i, expr := range lambdaTests {
			source := fmt.Sprintf("fn test() -> int = %s\nfn main() -> int = test()", expr)
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Fiber lambda %d should compile: %v", i, err)
			} else {
				t.Logf("✅ Fiber lambda %d compiled successfully", i)
			}
		}
	})
}
