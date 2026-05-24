package codegen_test

// Wave 2 coverage tests — compile-only smoke tests that exercise a wide
// range of partial-coverage codegen paths. Each program is the smallest
// piece of osprey that drives the targeted function down a specific arm.

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestWave2CodegenSmokes(t *testing.T) {
	programs := map[string]string{
		// generateInputCall — was 6%.
		"input_builtin": `
fn main() -> int {
  let line = input()
  match line {
    Success { value }   => print(value)
    Error   { message } => print(message)
  }
  0
}`,

		// generateAssignmentStatement — was 17.6%.
		"mut_reassign": `
fn main() -> int {
  mut x = 1
  x = 99
  x
  0
}`,

		// typeExpressionToLLVMType — was 17.6%. Hit Unit, function,
		// generic, and user-defined branches.
		"type_expressions_extern": `
extern fn ext_void() -> Unit
extern fn ext_int() -> int
extern fn ext_str() -> string
extern fn ext_fn(cb: fn(int) -> int) -> int
fn main() -> int {
  0
}`,

		// generateErrorConstructor (was 25%) + Result patterns.
		"error_constructor_in_match": `
fn classify(n: int) = match n {
  0 => Error { message: "zero" }
  _ => Success { value: n }
}
fn main() -> int {
  match classify(7) {
    Success { value }   => value
    Error   { message } => 0
  }
  0
}`,

		// inferTypeConstructor (was 37.5%) — discriminated union ctor.
		"discriminated_union": `
type Pri = Low | Med | High
type Box = OK { v: int } | Bad { msg: string }
fn main() -> int {
  let p = High
  let b = OK { v: 42 }
  match p {
    Low  => 0
    Med  => 1
    High => 2
  }
  match b {
    OK  { v }   => v
    Bad { msg } => 0
  }
  0
}`,

		// generateBlockExpression (was 64.7%) deep blocks.
		"nested_blocks": `
fn main() -> int {
  let x = {
    let a = 10
    let b = {
      let c = 5
      a + c
    }
    b * 2
  }
  match x {
    Success { value }   => value
    Error   { message } => 0
  }
  0
}`,

		// generateListAccess (was 81.3%) — list indexing returning Result.
		"list_index_match": `
fn main() -> int {
  let xs = [10, 20, 30]
  match xs[1] {
    Success { value }   => value
    Error   { message } => 0
  }
  0
}`,

		// Pipe operator chains for fusion bookkeeping.
		"pipe_chain": `
fn dbl(x: int) -> int = x + x
fn inc(x: int) -> int = x + 1
fn main() -> int {
  3 |> dbl |> inc |> print
  0
}`,

		// Range iterator forEach (range takes named args).
		"range_foreach": `
fn main() -> int {
  forEach(range(0, 5), print)
  0
}`,

		// Module declaration + access (covers ModuleAccessExpression).
		"module_access": `
module Math {
  fn pi() = 3
  fn e()  = 2
}
fn main() -> int {
  Math.pi() + Math.e()
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
