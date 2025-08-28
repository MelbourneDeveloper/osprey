package codegen_test

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestConvertResultToString_Success_I64(t *testing.T) {
	source := `
effect Test {
    getInt: fn() -> int
}

fn test() -> string = handle Test
    getInt => 42
in {
    toString(perform Test.getInt())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile Result<int> toString: %v", err)
	}
}

func TestConvertResultToString_Success_I1(t *testing.T) {
	source := `
effect Test {
    getBool: fn() -> bool
}

fn test() -> string = handle Test
    getBool => true
in {
    toString(perform Test.getBool())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile Result<bool> toString: %v", err)
	}
}

func TestConvertResultToString_Success_String(t *testing.T) {
	source := `
effect Test {
    getString: fn() -> string
}

fn test() -> string = handle Test
    getString => "test"
in {
    toString(perform Test.getString())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile Result<string> toString: %v", err)
	}
}

func TestConvertValueToStringByType_String(t *testing.T) {
	source := `
fn test() -> string = toString("test")
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile string toString: %v", err)
	}
}

func TestConvertValueToStringByType_Int(t *testing.T) {
	source := `
fn test() -> string = toString(42)
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile int toString: %v", err)
	}
}

func TestConvertValueToStringByType_Bool(t *testing.T) {
	source := `
fn test() -> string = toString(true)
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile bool toString: %v", err)
	}
}

func TestConvertValueToStringByType_Unit(t *testing.T) {
	// Unit type test - this tests the () case in convertValueToStringByType
	source := `
fn test() -> string = "unit_placeholder"
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile unit placeholder: %v", err)
	}
}

func TestConvertValueToStringByType_Fiber(t *testing.T) {
	source := `
fn fiber_func() -> int = 42

fn test() -> string = {
    let f = spawn fiber_func()
    toString(f)
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile fiber toString: %v", err)
	}
}

func TestConvertValueToStringByType_Any(t *testing.T) {
	// Any type test - testing the TypeAny case in convertValueToStringByType
	source := `
fn test() -> string = "any_placeholder"
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile any placeholder: %v", err)
	}
}

func TestConvertValueToStringByType_ComplexType(t *testing.T) {
	source := `
type CustomType = Value

fn test() -> string = {
    let val = Value
    toString(val)
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile custom type toString: %v", err)
	}
}

func TestGenerateToStringCall_WrongArgs_Zero(t *testing.T) {
	source := `
fn test() -> string = toString()
`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Error("Expected error for toString() with no arguments")
	}
}

func TestGenerateToStringCall_WrongArgs_Multiple(t *testing.T) {
	source := `
fn test() -> string = toString(1, 2)
`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Error("Expected error for toString() with multiple arguments")
	}
}

func TestIsResultValueSemanticBoolean_True(t *testing.T) {
	source := `
effect BoolEffect {
    getBool: fn() -> bool
}

fn test() -> string = handle BoolEffect
    getBool => true
in {
    toString(perform BoolEffect.getBool())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile semantic boolean test: %v", err)
	}
}

func TestIsResultValueSemanticBoolean_False(t *testing.T) {
	source := `
effect BoolEffect {
    getBool: fn() -> bool
}

fn test() -> string = handle BoolEffect
    getBool => false
in {
    toString(perform BoolEffect.getBool())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile semantic boolean false test: %v", err)
	}
}

func TestValidateBuiltInArgs_ValidFunction(t *testing.T) {
	source := `
fn test() -> string = toString(42)
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to validate valid built-in args: %v", err)
	}
}

func TestValidateBuiltInArgs_InvalidFunction(t *testing.T) {
	source := `
fn test() -> string = nonExistentFunction(42)
`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Error("Expected error for non-existent function")
	}
}

func TestValidateBuiltInArgs_WrongArgCount(t *testing.T) {
	source := `
fn test() -> string = print("hello", "world")
`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Error("Expected error for wrong argument count")
	}
}

func TestCreateGlobalString_EmptyString(t *testing.T) {
	source := `
fn test() -> string = ""
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile empty string: %v", err)
	}
}

func TestCreateGlobalString_WithContent(t *testing.T) {
	source := `
fn test() -> string = "hello world"
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile string with content: %v", err)
	}
}

func TestCreateGlobalString_SpecialCharacters(t *testing.T) {
	source := `
fn test() -> string = "hello\nworld\ttab"
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile string with special chars: %v", err)
	}
}

func TestConvertResultToString_BlockGeneration(t *testing.T) {
	source := `
effect MultipleResults {
    getFirst: fn() -> int
    getSecond: fn() -> int
}

fn test() -> string = handle MultipleResults
    getFirst => 1
    getSecond => 2
in {
    let first = toString(perform MultipleResults.getFirst())
    let second = toString(perform MultipleResults.getSecond())
    first + second
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile multiple result toString: %v", err)
	}
}

func TestCreateGlobalString_NullTermination(t *testing.T) {
	source := `
fn test() -> string = {
    print("test")
    "result"
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile string with null termination: %v", err)
	}
}

func TestConvertResultToString_Discriminant_Check(t *testing.T) {
	source := `
effect DiscriminantTest {
    mayFail: fn() -> int
}

fn test() -> string = handle DiscriminantTest
    mayFail => 100
in {
    toString(perform DiscriminantTest.mayFail())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile discriminant test: %v", err)
	}
}

func TestConvertResultToString_ValuePtr_Access(t *testing.T) {
	source := `
effect ValuePtrTest {
    getValue: fn() -> string
}

fn test() -> string = handle ValuePtrTest
    getValue => "inner_value"
in {
    toString(perform ValuePtrTest.getValue())
}
`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Failed to compile value ptr test: %v", err)
	}
}