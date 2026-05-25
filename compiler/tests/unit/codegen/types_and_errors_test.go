package codegen_test

// Coverage-focused tests for Type interface implementations and error
// wrapper functions in internal/codegen. These touch the marker methods
// (String, Category, Equals) that the rest of the codebase exercises only
// indirectly, and the error-wrapping helpers that fire on real failures.

import (
	"errors"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

var errTestSentinel = errors.New("boom")

// ---- Type interface marker methods --------------------------------------

func TestPrimitiveTypeMethods(t *testing.T) {
	pt := codegen.NewPrimitiveType("int")
	if pt.String() != "int" {
		t.Errorf("String mismatch: %q", pt.String())
	}
	if pt.Category() != codegen.PrimitiveTypeCategory {
		t.Errorf("Category mismatch")
	}
	if !pt.Equals(codegen.NewPrimitiveType("int")) {
		t.Error("equal primitives should compare equal")
	}
	if pt.Equals(codegen.NewPrimitiveType("bool")) {
		t.Error("different primitives must not equal")
	}
	if pt.Equals(codegen.NewGenericType("List", nil)) {
		t.Error("primitive must not equal non-primitive")
	}
}

func TestGenericTypeMethods(t *testing.T) {
	intT := codegen.NewPrimitiveType("int")
	g1 := codegen.NewGenericType("List", []codegen.Type{intT})
	g2 := codegen.NewGenericType("List", []codegen.Type{intT})
	g3 := codegen.NewGenericType("List", []codegen.Type{codegen.NewPrimitiveType("str")})
	g4 := codegen.NewGenericType("List", nil)

	if g1.String() != "List[int]" {
		t.Errorf("String: %q", g1.String())
	}
	if g4.String() != "List" {
		t.Errorf("empty arg list: %q", g4.String())
	}
	if g1.Category() != codegen.GenericTypeCategory {
		t.Error("category mismatch")
	}
	if !g1.Equals(g2) {
		t.Error("identical generics should equal")
	}
	if g1.Equals(g3) {
		t.Error("different arg generics must not equal")
	}
	if g1.Equals(codegen.NewGenericType("Map", []codegen.Type{intT})) {
		t.Error("different name should not equal")
	}
	if g1.Equals(intT) {
		t.Error("generic must not equal non-generic")
	}
}

func TestFunctionTypeMethods(t *testing.T) {
	intT := codegen.NewPrimitiveType("int")
	strT := codegen.NewPrimitiveType("str")
	f1 := codegen.NewFunctionType([]codegen.Type{intT, strT}, intT)
	f2 := codegen.NewFunctionType([]codegen.Type{intT, strT}, intT)
	f3 := codegen.NewFunctionType([]codegen.Type{intT}, intT)
	f4 := codegen.NewFunctionType([]codegen.Type{intT, intT}, intT)
	f5 := codegen.NewFunctionType([]codegen.Type{intT, strT}, strT)

	if f1.String() != "(int, str) -> int" {
		t.Errorf("String: %q", f1.String())
	}
	if f1.Category() != codegen.FunctionTypeCategory {
		t.Error("category mismatch")
	}
	if !f1.Equals(f2) {
		t.Error("identical fn types should equal")
	}
	if f1.Equals(f3) {
		t.Error("different arity should not equal")
	}
	if f1.Equals(f4) {
		t.Error("different param types should not equal")
	}
	if f1.Equals(f5) {
		t.Error("different return type should not equal")
	}
	if f1.Equals(intT) {
		t.Error("fn must not equal non-fn")
	}
}

func TestRecordTypeMethods(t *testing.T) {
	intT := codegen.NewPrimitiveType("int")
	strT := codegen.NewPrimitiveType("str")
	r1 := codegen.NewRecordType("P", map[string]codegen.Type{"x": intT, "y": intT})
	r2 := codegen.NewRecordType("P", map[string]codegen.Type{"x": intT, "y": intT})
	r3 := codegen.NewRecordType("P", map[string]codegen.Type{"x": intT, "z": intT})
	r4 := codegen.NewRecordType("Q", map[string]codegen.Type{"x": intT, "y": intT})
	r5 := codegen.NewRecordType("P", map[string]codegen.Type{"x": intT})
	r6 := codegen.NewOrderedRecordType("P", map[string]codegen.Type{"x": intT, "y": strT}, []string{"x", "y"})

	if !strings.HasPrefix(r1.String(), "{") || !strings.HasSuffix(r1.String(), "}") {
		t.Errorf("String shape: %q", r1.String())
	}
	if r1.Category() != codegen.RecordTypeCategory {
		t.Error("category mismatch")
	}
	if !r1.Equals(r2) {
		t.Error("identical records should equal")
	}
	if r1.Equals(r3) {
		t.Error("different fields should not equal")
	}
	if r1.Equals(r4) {
		t.Error("different name should not equal")
	}
	if r1.Equals(r5) {
		t.Error("different field count should not equal")
	}
	if r1.Equals(r6) {
		t.Error("different field types should not equal")
	}
	if r1.Equals(intT) {
		t.Error("record must not equal non-record")
	}
}

func TestUnionTypeMethods(t *testing.T) {
	intT := codegen.NewPrimitiveType("int")
	strT := codegen.NewPrimitiveType("str")
	u1 := codegen.NewUnionType("E", []codegen.Type{intT, strT})
	u2 := codegen.NewUnionType("E", []codegen.Type{intT, strT})
	u3 := codegen.NewUnionType("E", []codegen.Type{intT})
	u4 := codegen.NewUnionType("F", []codegen.Type{intT, strT})
	u5 := codegen.NewUnionType("E", []codegen.Type{strT, intT})

	if !strings.HasPrefix(u1.String(), "E(") {
		t.Errorf("String shape: %q", u1.String())
	}
	if u1.Category() != codegen.UnionTypeCategory {
		t.Error("category mismatch")
	}
	if !u1.Equals(u2) {
		t.Error("identical unions should equal")
	}
	if u1.Equals(u3) {
		t.Error("different variant count should not equal")
	}
	if u1.Equals(u4) {
		t.Error("different name should not equal")
	}
	if u1.Equals(u5) {
		t.Error("different variant order should not equal")
	}
	if u1.Equals(intT) {
		t.Error("union must not equal non-union")
	}
}

func TestTypeVarMethods(t *testing.T) {
	tv1 := codegen.NewTypeVar(1, "a")
	tv2 := codegen.NewTypeVar(1, "a")
	tv3 := codegen.NewTypeVar(2, "b")
	tvEmpty := codegen.NewTypeVar(7, "")

	if tv1.String() != "a" {
		t.Errorf("named String: %q", tv1.String())
	}
	if tvEmpty.String() != "t7" {
		t.Errorf("anonymous String: %q", tvEmpty.String())
	}
	if tv1.Category() != codegen.TypeVariableCategory {
		t.Error("category mismatch")
	}
	if !tv1.Equals(tv2) {
		t.Error("same-id type vars should equal")
	}
	if tv1.Equals(tv3) {
		t.Error("different-id type vars should not equal")
	}
	if tv1.Equals(codegen.NewPrimitiveType("a")) {
		t.Error("type var must not equal non-type-var")
	}
}

func TestConcreteTypeMethods(t *testing.T) {
	c1 := codegen.NewConcreteType("Foo")
	c2 := codegen.NewConcreteType("Foo")
	c3 := codegen.NewConcreteType("Bar")
	if c1.String() != "Foo" {
		t.Errorf("String: %q", c1.String())
	}
	if c1.Category() != codegen.PrimitiveTypeCategory {
		t.Error("category mismatch")
	}
	if !c1.Equals(c2) {
		t.Error("equal concretes should compare equal")
	}
	if c1.Equals(c3) {
		t.Error("different concretes should not equal")
	}
	if c1.Equals(codegen.NewPrimitiveType("Foo")) {
		t.Error("concrete must not equal non-concrete (PrimitiveType)")
	}
}

// ---- Error wrappers -----------------------------------------------------

func TestErrorWrapperHelpers(t *testing.T) {
	pos := &ast.Position{Line: 7, Column: 3}

	cases := []struct {
		name    string
		err     error
		mustHas string
	}{
		{"NoToStringImpl", codegen.WrapNoToStringImpl("Foo"), "Foo"},
		{"UnsupportedBinaryOpWithPos", codegen.WrapUnsupportedBinaryOpWithPos("@", pos), "@"},
		{"UnsupportedBinaryOpWithPosNil", codegen.WrapUnsupportedBinaryOpWithPos("@", nil), "@"},
		{"VoidArithmeticWithPos", codegen.WrapVoidArithmeticWithPos("+", pos), "+"},
		{"VoidArithmeticWithPosNil", codegen.WrapVoidArithmeticWithPos("+", nil), "+"},
		{"UnsupportedUnaryOpWithPos", codegen.WrapUnsupportedUnaryOpWithPos("?", pos), "?"},
		{"FieldAccessOnResult", codegen.WrapFieldAccessOnResult("x", "Result"), "Result"},
		{"ConstraintResultFieldAccessWithPos", codegen.WrapConstraintResultFieldAccessWithPos("x", pos), "x"},
		{"ConstraintResultFieldAccessWithPosNil", codegen.WrapConstraintResultFieldAccessWithPos("x", nil), "x"},
		{"FieldNotFoundInRecord", codegen.WrapFieldNotFoundInRecord("x", "Point"), "Point"},
		{"FieldAccessOnNonRecord", codegen.WrapFieldAccessOnNonRecord("x", "int"), "int"},
		{"FieldAccessOnLegacyRecord", codegen.WrapFieldAccessOnLegacyRecord("x", "P"), "P"},
		{"CannotUpdateNonRecord", codegen.WrapCannotUpdateNonRecord("int"), "int"},
		{"FieldNotInRecordType", codegen.WrapFieldNotInRecordType("x", "Q"), "Q"},
		{"UndefinedType", codegen.WrapUndefinedType("Foo"), "Foo"},
		{"MissingField", codegen.WrapMissingField("x"), "x"},
		{"HTTPStopServerUnknownNamedArg", codegen.WrapHTTPStopServerUnknownNamedArg("foo"), "foo"},
		{"HTTPFunctionMissingNamedArg", codegen.WrapHTTPFunctionMissingNamedArg("httpListen", "port"), "port"},
		{"UnsupportedCollectionType", codegen.WrapUnsupportedCollectionType("Bag"), "Bag"},
		{"MethodCallNotImplementedWithPos", codegen.WrapMethodCallNotImplementedWithPos("foo", pos), "foo"},
		{"ArithmeticTypeMismatch", codegen.WrapArithmeticTypeMismatch("+", "left", "int", "str"), "left"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatal("nil error")
			}
			if !strings.Contains(tc.err.Error(), tc.mustHas) {
				t.Errorf("missing %q in %q", tc.mustHas, tc.err.Error())
			}
		})
	}
}

// Additional wrappers (compile/link/match/argument families) covered in
// one shot. We only assert non-nil + substring; exact format is enforced
// by integration tests elsewhere.
func TestMoreErrorWrappers(t *testing.T) {
	pos := &ast.Position{Line: 11, Column: 4}
	sentinel := errTestSentinel

	cases := []struct {
		name string
		err  error
		want string
	}{
		{"ImmutableAssignmentErrorWithPos", codegen.WrapImmutableAssignmentErrorWithPos("v", pos), "v"},
		{"ImmutableAssignmentErrorWithPosNil", codegen.WrapImmutableAssignmentErrorWithPos("v", nil), "v"},
		{"BuiltInTwoArgs", codegen.WrapBuiltInTwoArgs("foo"), "foo"},
		{"FunctionNotDeclared", codegen.WrapFunctionNotDeclared("foo"), "foo"},
		{"UndefinedVariable", codegen.WrapUndefinedVariable("v"), "v"},
		{"MatchNotExhaustiveWithPos", codegen.WrapMatchNotExhaustiveWithPos([]string{"A", "B"}, pos), "A"},
		{"MatchNotExhaustiveWithPosNil", codegen.WrapMatchNotExhaustiveWithPos([]string{"A"}, nil), "A"},
		{"MatchTypeMismatchWithPos", codegen.WrapMatchTypeMismatchWithPos(2, "int", "str", pos), "int"},
		{"MatchTypeMismatchWithPosNil", codegen.WrapMatchTypeMismatchWithPos(2, "int", "str", nil), "int"},
		{"UnknownVariantWithPos", codegen.WrapUnknownVariantWithPos("Foo", "Bar", pos), "Foo"},
		{"UnknownVariantWithPosNil", codegen.WrapUnknownVariantWithPos("Foo", "Bar", nil), "Foo"},
		{"TypeMismatchWithPos", codegen.WrapTypeMismatchWithPos("int", "v", "str", pos), "int"},
		{"TypeMismatchWithPosNil", codegen.WrapTypeMismatchWithPos("int", "v", "str", nil), "int"},
		{"PrintConvertError", codegen.WrapPrintConvertError("X", "foo"), "X"},
		{"PrintDetermineError", codegen.WrapPrintDetermineError("foo"), "foo"},
		{"MissingArgument", codegen.WrapMissingArgument("a", "foo"), "a"},
		{"MissingArgumentWithPos", codegen.WrapMissingArgumentWithPos("a", "foo", pos), "a"},
		{"MissingArgumentWithPosNil", codegen.WrapMissingArgumentWithPos("a", "foo", nil), "a"},
		{"LLVMGenFailed", codegen.WrapLLVMGenFailed(sentinel), "boom"},
		{"CompileToObj", codegen.WrapCompileToObj(sentinel, "stderr"), "stderr"},
		{"LinkExecutable", codegen.WrapLinkExecutable("clang", sentinel, "stderr"), "clang"},
		{"BuiltInRedefine", codegen.WrapBuiltInRedefine("foo"), "foo"},
		{"WrongArgCount", codegen.WrapWrongArgCount("foo", 2, 3), "3"},
		{"BuiltInFunctionWrongArgs known", codegen.WrapBuiltInFunctionWrongArgs("listLength", 99), "99"},
		{"BuiltInFunctionWrongArgs unknown", codegen.WrapBuiltInFunctionWrongArgs("__nope__", 4), "4"},
		{"ToStringWrongArgs", codegen.WrapToStringWrongArgs(99), "99"},
		{"PrintWrongArgs", codegen.WrapPrintWrongArgs(99), "99"},
		{"UnsupportedCallExpressionSecurity", codegen.WrapUnsupportedCallExpressionSecurity("foo"), "foo"},
		{"MethodCallNotImplemented", codegen.WrapMethodCallNotImplemented("foo"), "foo"},
		{"FunctionRequiresNamedArgsWithPos", codegen.WrapFunctionRequiresNamedArgsWithPos("foo", []string{"a", "b"}, pos), "foo"},
		{"FunctionRequiresNamedArgsWithPosNil", codegen.WrapFunctionRequiresNamedArgsWithPos("foo", []string{"a", "b"}, nil), "foo"},
		{"UnsupportedBinaryOp", codegen.WrapUnsupportedBinaryOp("@"), "@"},
		{"UndefinedTypeWithPos", codegen.WrapUndefinedTypeWithPos("Foo", pos), "Foo"},
		{"UndefinedTypeWithPosNil", codegen.WrapUndefinedTypeWithPos("Foo", nil), "Foo"},
		{"UndefinedFunction", codegen.WrapUndefinedFunction("foo"), "foo"},
		{"FunctionArgsWithPosNil", codegen.WrapFunctionArgsWithPos("listLength", 1, 0, nil), "0"},
		{"FunctionArgsWithPos", codegen.WrapFunctionArgsWithPos("listLength", 1, 0, pos), "0"},
		{"UnsupportedExpression", codegen.WrapUnsupportedExpression(pos), "Position"},
		{"UndefinedVariableWithPos", codegen.WrapUndefinedVariableWithPos("v", pos), "v"},
		{"UndefinedVariableWithPosNil", codegen.WrapUndefinedVariableWithPos("v", nil), "v"},
		{"UnsupportedStatement", codegen.WrapUnsupportedStatement(pos), "Position"},
		{"ToolNotFound", codegen.WrapToolNotFound("clang"), "clang"},
		{"NoSuitableCompiler", codegen.WrapNoSuitableCompiler([]string{"clang", "gcc"}), "clang"},
		{"WriteIRFile", codegen.WrapWriteIRFile(sentinel), "boom"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatal("nil error")
			}
			if !strings.Contains(tc.err.Error(), tc.want) {
				t.Errorf("missing %q in %q", tc.want, tc.err.Error())
			}
		})
	}
}

// ---- collection / string error paths via compile errors ----------------

// TestCollectionBuiltinArgErrors covers the wrong-arg-count branches of the
// collection builtins by feeding the compiler invalid call sites and
// checking it surfaces a usable error.
func TestCollectionBuiltinArgErrors(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{"listLength wrong arity", "fn main() -> int { listLength(); 0 }"},
		{"listContains wrong arity", "fn main() -> int { listContains(); 0 }"},
		{"listAppend wrong arity", "fn main() -> int { listAppend(); 0 }"},
		{"listPrepend wrong arity", "fn main() -> int { listPrepend(); 0 }"},
		{"listConcat wrong arity", "fn main() -> int { listConcat(); 0 }"},
		{"listReverse wrong arity", "fn main() -> int { listReverse(); 0 }"},
		{"mapLength wrong arity", "fn main() -> int { mapLength(); 0 }"},
		{"mapContains wrong arity", "fn main() -> int { mapContains(); 0 }"},
		{"mapSet wrong arity", "fn main() -> int { mapSet(); 0 }"},
		{"mapRemove wrong arity", "fn main() -> int { mapRemove(); 0 }"},
		{"mapMerge wrong arity", "fn main() -> int { mapMerge(); 0 }"},
		{"mapKeys wrong arity", "fn main() -> int { mapKeys(); 0 }"},
		{"mapValues wrong arity", "fn main() -> int { mapValues(); 0 }"},
		{"forEachList wrong arity", "fn main() -> int { forEachList(); 0 }"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tc.source)
			if err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}

// TestStringBuiltinArgErrors covers wrong-arg-count branches across the
// string builtin generators.
func TestStringBuiltinArgErrors(t *testing.T) {
	cases := []string{
		`fn main() -> int { length(); 0 }`,
		`fn main() -> int { isEmpty(); 0 }`,
		`fn main() -> int { startsWith(); 0 }`,
		`fn main() -> int { endsWith(); 0 }`,
		`fn main() -> int { contains(); 0 }`,
		`fn main() -> int { indexOf(); 0 }`,
		`fn main() -> int { substring(); 0 }`,
		`fn main() -> int { take(); 0 }`,
		`fn main() -> int { drop(); 0 }`,
		`fn main() -> int { trim(); 0 }`,
		`fn main() -> int { trimStart(); 0 }`,
		`fn main() -> int { trimEnd(); 0 }`,
		`fn main() -> int { reverse(); 0 }`,
		`fn main() -> int { toUpperCase(); 0 }`,
		`fn main() -> int { toLowerCase(); 0 }`,
		`fn main() -> int { replace(); 0 }`,
		`fn main() -> int { repeat(); 0 }`,
		`fn main() -> int { padStart(); 0 }`,
		`fn main() -> int { padEnd(); 0 }`,
		`fn main() -> int { parseInt(); 0 }`,
		`fn main() -> int { parseFloat(); 0 }`,
		`fn main() -> int { lines(); 0 }`,
		`fn main() -> int { words(); 0 }`,
		`fn main() -> int { split(); 0 }`,
	}
	for _, src := range cases {
		_, err := codegen.CompileToLLVM(src)
		if err == nil {
			t.Errorf("expected error for source: %s", src)
		}
	}
}

// Sanity: package errors satisfy errors.Is wrapping.
func TestSentinelErrorsAreWrappable(t *testing.T) {
	wrapped := codegen.WrapUndefinedType("Foo")
	if !errors.Is(wrapped, codegen.ErrUndefinedType) {
		t.Error("wrapped error must satisfy errors.Is on the sentinel")
	}
}

// ---- Registry helpers ---------------------------------------------------

func TestRegistryHelpers(t *testing.T) {
	r := codegen.GlobalBuiltInRegistry

	// GetFunction happy path
	if _, ok := r.GetFunction("listLength"); !ok {
		t.Error("listLength should be registered")
	}
	if _, ok := r.GetFunction("__missing__"); ok {
		t.Error("missing fn must not resolve")
	}

	// GetAllFunctions returns a non-empty map containing listLength
	all := r.GetAllFunctions()
	if _, ok := all["listLength"]; !ok {
		t.Error("GetAllFunctions should include listLength")
	}

	// GetFunctionsByCategory: at least one functional category entry
	if got := r.GetFunctionsByCategory(codegen.CategoryFunctional); len(got) == 0 {
		t.Error("CategoryFunctional should return >0 entries")
	}

	// IsProtectedFunction: forEachList is protected; listLength is not
	if !r.IsProtectedFunction("forEachList") {
		t.Error("forEachList must be reported as protected")
	}
	if r.IsProtectedFunction("listLength") {
		t.Error("listLength must not be reported as protected")
	}
	if r.IsProtectedFunction("__missing__") {
		t.Error("missing function must not be reported as protected")
	}

	// RequiresPermission: PermissionNone matches any builtin that opts-in to it
	if !r.RequiresPermission("trimStart", codegen.PermissionNone) {
		t.Error("trimStart was registered with PermissionNone")
	}
	if r.RequiresPermission("__missing__", codegen.PermissionNone) {
		t.Error("missing function must not match any permission")
	}
}

// ---- Compile-only smoke tests covering more inference & codegen paths ---
//
// Each program compiles successfully and exercises a range of inference
// branches (records, unions, pattern matching, blocks, fibers, effects).
// They are coverage-focused: by getting CompileToLLVM through these
// constructs once each, we hit dozens of partial-coverage functions.

func TestCompileFeatureSmokeTests(t *testing.T) {
	programs := map[string]string{
		"record_basic": `
type Pair = { x: int, y: int }
fn main() -> int {
  let p = Pair { x: 1, y: 2 }
  p.x
  0
}`,
		"record_update": `
type Pair = { x: int, y: int }
fn main() -> int {
  let p = Pair { x: 1, y: 2 }
  let q = p { x: 99 }
  q.x
  0
}`,
		"union_match": `
type Shape = Circle { r: int } | Square { s: int }
fn area(sh: Shape) -> int = match sh {
  Circle { r } => r * r * 3
  Square { s } => s * s
}
fn main() -> int {
  let c = Circle { r: 5 }
  area(c)
  0
}`,
		"enum_union": `
type Color = Red | Green | Blue
fn main() -> int {
  let c = Red
  match c {
    Red => 1
    Green => 2
    Blue => 3
  }
  0
}`,
		"list_in_match": `
fn main() -> int {
  let xs = listAppend(listAppend(List(), 1), 2)
  match listLength(xs) {
    0 => print("empty")
    _ => print("nonempty")
  }
  0
}`,
		"map_literal_int_key": `
fn main() -> int {
  let m = {1: "a", 2: "b"}
  mapLength(m)
  0
}`,
		"map_literal_string_key": `
fn main() -> int {
  let m = {"a": 1, "b": 2}
  mapLength(m)
  0
}`,
		"nested_calls": `
fn main() -> int {
  let s = "  hello  ".trim().toUpperCase().take(3)
  length(s)
  0
}`,
		"chained_list_plus": `
fn main() -> int {
  let a = listAppend(List(), 1)
  let b = listAppend(List(), 2)
  let c = listAppend(List(), 3)
  listLength(a + b + c + a)
  0
}`,
		"chained_map_plus": `
fn main() -> int {
  let m1 = {"a": 1}
  let m2 = {"b": 2}
  let m3 = {"c": 3}
  mapLength(m1 + m2 + m3)
  0
}`,
		"string_interpolation": `
fn main() -> int {
  let name = "world"
  let n = 42
  print("hello ${name} ${n}")
  0
}`,
		"block_expression": `
fn main() -> int {
  let x = {
    let a = 10
    let b = 20
    a + b
  }
  match x {
    Success { value } => value
    Error { message } => 0
  }
  0
}`,
		"spawn_await": `
fn compute(n: int) -> int = n * n
fn main() -> int {
  let f = spawn compute(7)
  let r = await(f)
  r
  0
}`,
		"if_via_match": `
fn classify(n: int) -> int = match n {
  0 => 0
  1 => 1
  _ => 2
}
fn main() -> int {
  classify(5)
  0
}`,
		"list_index_access": `
fn main() -> int {
  let xs = [10, 20, 30]
  xs[1]
  0
}`,
		"nested_record_access": `
type Inner = { v: int }
type Outer = { inner: Inner }
fn main() -> int {
  let o = Outer { inner: Inner { v: 42 } }
  o.inner.v
  0
}`,
		"contains_polymorphic": `
fn main() -> int {
  let xs = listAppend(listAppend(List(), 1), 2)
  match listContains(xs, 2) {
    true => 1
    false => 0
  }
  0
}`,
	}

	for name, src := range programs {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(src)
			if err != nil {
				t.Fatalf("%s failed to compile: %v", name, err)
			}
		})
	}
}
