package codegen_test

// Wave 4 — targeted unit tests for specific exported helpers that lack
// coverage, plus compile-only programs that drive the remaining
// 0%/low-coverage codegen branches identified by `go tool cover -func`.

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

// CheckProtectedFunction returns an error when a user function redeclares
// a builtin name. Was 0% — exported but never called from a test.
func TestCheckProtectedFunction(t *testing.T) {
	// Redefining a builtin — must error.
	err := codegen.CheckProtectedFunction(&ast.FunctionDeclaration{Name: "length"})
	if err == nil {
		t.Fatal("expected redefining `length` builtin to error")
	}
	if !strings.Contains(err.Error(), "length") {
		t.Errorf("error should mention `length`, got %q", err.Error())
	}

	// Non-builtin name — must succeed (return nil).
	err = codegen.CheckProtectedFunction(&ast.FunctionDeclaration{Name: "myUniqueName"})
	if err != nil {
		t.Errorf("non-builtin name must not error, got %v", err)
	}
}

// CompileAndCapture / CompileAndCaptureWithSecurity — 0% covered.
// These compile + link + run, capturing stdout. Smallest hello-world.
func TestCompileAndCapture(t *testing.T) {
	src := `fn main() -> int { print("xyz") 0 }`
	out, err := codegen.CompileAndCapture(src)
	if err != nil {
		t.Fatalf("compile-and-capture: %v", err)
	}
	if !strings.Contains(out, "xyz") {
		t.Errorf("expected output to contain xyz, got %q", out)
	}
}

func TestCompileAndCaptureWithSecurity(t *testing.T) {
	src := `fn main() -> int { print("zzz") 0 }`
	out, err := codegen.CompileAndCaptureWithSecurity(src, codegen.SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
	if err != nil {
		t.Fatalf("compile-and-capture-secured: %v", err)
	}
	if !strings.Contains(out, "zzz") {
		t.Errorf("expected output to contain zzz, got %q", out)
	}
}

// CompileToExecutable / CompileAndRun were also untouched.
func TestCompileToExecutable(t *testing.T) {
	src := `fn main() -> int { 0 }`
	tmp := t.TempDir() + "/out"
	err := codegen.CompileToExecutable(src, tmp)
	if err != nil {
		t.Fatalf("compile-to-executable: %v", err)
	}
}

func TestCompileAndRun(t *testing.T) {
	src := `fn main() -> int { 0 }`
	err := codegen.CompileAndRun(src)
	if err != nil {
		t.Fatalf("compile-and-run: %v", err)
	}
}

// Programs that exercise extra-language features whose generators are
// still partial. These are all compile-only — no execution.
func TestWave4ExtraSmokes(t *testing.T) {
	programs := map[string]string{
		// generateJoinCall — was 0%. List literals lower to a struct
		// value that doesn't satisfy the runtime-allocated list
		// contract — use listAppend to get an OspreyList handle.
		"join_call": `
fn main() -> int {
  let parts = listAppend(listAppend(listAppend(List(), "a"), "b"), "c")
  let s = join(parts, ",")
  print(s)
  0
}`,
		// Record update expression — was 72%.
		"record_update_chain": `
type P = { x: int, y: int, z: int }
fn main() -> int {
  let p0 = P { x: 0, y: 0, z: 0 }
  let p1 = p0 { x: 1 }
  let p2 = p1 { y: 2 }
  let p3 = p2 { z: 3 }
  p3.x + p3.y + p3.z
  0
}`,
		// File read/write builtins (compile-only).
		"file_io_compile": `
fn main() -> int {
  match writeFile("/tmp/xx", "data") {
    Success { value }   => print("wrote")
    Error   { message } => print(message)
  }
  match readFile("/tmp/xx") {
    Success { value }   => print(value)
    Error   { message } => print(message)
  }
  0
}`,
		// Toll-free Math chain.
		"math_chain": `
fn main() -> int {
  match 10 * 20 / 5 - 3 + 4 {
    Success { value }   => print(value)
    Error   { message } => print(message)
  }
  0
}`,
		// String formatting + interpolation with multiple types.
		"interpolation_mixed": `
fn main() -> int {
  let name = "x"
  let n = 42
  let b = true
  print("name=${name} n=${n} b=${b}")
  0
}`,
		// mapList / filterList / foldList builtins — covers
		// generateListBuilderLoop (shared by map+filter) and
		// generateFoldListCall counted-loop.
		"map_list": `
fn main() -> int {
  let xs = listAppend(listAppend(List(), 10), 20)
  let doubled = mapList(xs, fn(x: int) => x * 2)
  forEachList(doubled, print)
  0
}`,
		"filter_list": `
fn main() -> int {
  let xs = listAppend(listAppend(listAppend(List(), 1), 2), 3)
  let kept = filterList(xs, fn(x: int) => x > 1)
  forEachList(kept, print)
  0
}`,
		"fold_list": `
fn main() -> int {
  let xs = listAppend(listAppend(listAppend(List(), 1), 2), 3)
  let total = foldList(xs, 0, fn(acc: int, x: int) => match acc + x {
    Success v => v
    Error e => acc
  })
  print(total)
  0
}`,
	}

	for name, src := range programs {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(src)
			if err != nil {
				t.Fatalf("%s failed: %v", name, err)
			}
		})
	}
}
