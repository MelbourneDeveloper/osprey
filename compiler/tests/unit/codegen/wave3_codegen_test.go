package codegen_test

// Wave 3 coverage — programs that drive every comparison/logical/unary
// branch, every field-access shape, every map-key type, and every error
// path that bubbles out of inference or codegen. Each test is a tiny
// compile-only smoke; the value is in the set of programs being broad
// rather than any single program being deep.

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestWave3HappySmokes(t *testing.T) {
	programs := map[string]string{
		// Every comparison operator on every primitive.
		"comparison_ops_int": `
fn main() -> int {
  let a = 5
  let b = 7
  print(a == b)
  print(a != b)
  print(a < b)
  print(a <= b)
  print(a > b)
  print(a >= b)
  0
}`,
		"comparison_ops_bool": `
fn main() -> int {
  let a = true
  let b = false
  print(a == b)
  print(a != b)
  0
}`,
		"comparison_ops_string": `
fn main() -> int {
  let a = "hello"
  let b = "world"
  print(a == b)
  print(a != b)
  0
}`,
		// Every logical operator.
		"logical_ops": `
fn main() -> int {
  let a = true
  let b = false
  print(a && b)
  print(a || b)
  print(!a)
  0
}`,
		// Float arithmetic + comparison.
		"float_arith": `
fn main() -> int {
  let a = 1.5
  let b = 2.25
  match a + b {
    Success { value }   => print(value)
    Error   { message } => print(message)
  }
  print(a < b)
  print(a >= b)
  0
}`,
		// Bool map keys (collectionKeyBool tag).
		"map_bool_keys": `
fn main() -> int {
  let m = {true: 1, false: 2}
  mapLength(m)
  0
}`,
		// Result chaining via match.
		"result_chain_match": `
fn parse(s: string) -> int = match parseInt(s) {
  Success { value }   => value
  Error   { message } => 0
}
fn main() -> int {
  parse("42") + parse("oops")
  0
}`,
		// Nested if-via-match returning different branches.
		"nested_match_branches": `
type Tag = Red | Green | Blue
fn pick(t: Tag) -> int = match t {
  Red   => match 10 {
    0 => 0
    _ => 100
  }
  Green => 200
  Blue  => 300
}
fn main() -> int {
  pick(Red)
}`,
		// Big record with mixed field types.
		"record_mixed_fields": `
type Big = { name: string, age: int, active: bool, score: float }
fn main() -> int {
  let b = Big { name: "x", age: 1, active: true, score: 9.5 }
  print(b.name)
  print(b.age)
  print(b.active)
  print(b.score)
  0
}`,
		// Recursive function via match (covers function-signature generation).
		"recursive_fn": `
fn fact(n: int) -> int = match n {
  0 => 1
  _ => match n * fact(n - 1) {
    Success { value }   => value
    Error   { message } => 0
  }
}
fn main() -> int {
  fact(5)
  0
}`,
		// Concat empty + non-empty list (boundary).
		"list_concat_empty": `
fn main() -> int {
  let a = listAppend(listAppend(List(), 1), 2)
  let b = List()
  listLength(a + b) + listLength(b + a)
  0
}`,
		// Negative integers and unary minus.
		"unary_minus_chain": `
fn main() -> int {
  let a = -5
  let b = -a
  print(a)
  print(b)
  0
}`,
		// String + string (infallible concat).
		"string_concat_op": `
fn main() -> int {
  let s = "hello" + ", " + "world"
  print(s)
  print(length(s))
  0
}`,
		// fold over range.
		"range_fold": `
fn add(acc: int, x: int) -> int = acc + x
fn main() -> int {
  match fold(range(1, 5), 0, add) {
    Success { value }   => print(value)
    Error   { message } => print(message)
  }
  0
}`,
		// filter + map + forEach.
		"pipe_filter_map_foreach": `
fn isEven(x: int) -> bool = (x % 2) == 0
fn sq(x: int) -> int = x * x
fn main() -> int {
  range(1, 6) |> filter(isEven) |> map(sq) |> forEach(print)
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

// Compile-error smokes — exercise the various Wrap*Error code paths.
func TestWave3CompileErrors(t *testing.T) {
	cases := map[string]struct {
		src  string
		want string // substring expected in error
	}{
		"undefined_variable": {
			src:  `fn main() -> int { x }`,
			want: "x",
		},
		"undefined_function_call": {
			src:  `fn main() -> int { unknownFn(1) }`,
			want: "unknownFn",
		},
		"unsupported_binary_op": {
			src:  `fn main() -> int { 1 @ 2 }`,
			want: "@", // parse error before codegen
		},
		"immutable_assignment": {
			src: `
fn main() -> int {
  let x = 1
  x = 2
  0
}`,
			want: "x",
		},
		"field_not_in_record": {
			src: `
type P = { x: int }
fn main() -> int {
  let p = P { x: 1 }
  p.y
  0
}`,
			want: "",
		},
		"undefined_type_in_let": {
			src: `
fn main() -> int {
  let x: NoSuchType = 1
  0
}`,
			want: "",
		},
		"wrong_arg_count_known_builtin": {
			src: `
fn main() -> int {
  length("a", "b")
  0
}`,
			want: "",
		},
		"match_non_exhaustive_union": {
			src: `
type T = A | B | C
fn main() -> int {
  let t = A
  match t {
    A => 1
    B => 2
  }
  0
}`,
			want: "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tc.src)
			if err == nil {
				t.Fatalf("expected compile error for %s", name)
			}
			if tc.want != "" && !strings.Contains(err.Error(), tc.want) {
				t.Errorf("%s: error %q missing %q", name, err.Error(), tc.want)
			}
		})
	}
}
