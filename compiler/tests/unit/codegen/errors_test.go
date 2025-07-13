package codegen

import (
	"errors"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestErrorDefinition(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrToStringWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrPrintWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrInputWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrMethodNotImpl", codegen.ErrMethodNotImpl, "method not implemented"},
		{"ErrNoToStringForFunc", codegen.ErrNoToStringForFunc, "no toString implementation for function"},
		{"ErrPrintComplexExpr", codegen.ErrPrintComplexExpr, "print complex expression error"},
		{"ErrPrintUnknownFunc", codegen.ErrPrintUnknownFunc, "print unknown function error"},
		{"ErrFunctionRequiresNamed", codegen.ErrFunctionRequiresNamed, "function requires named arguments"},
		{"ErrWrongArgCount", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrMissingArgument", codegen.ErrMissingArgument, "missing argument"},
		{"ErrParseErrors", codegen.ErrParseErrors, "parse errors"},
		{"ErrLLVMGenFailed", codegen.ErrLLVMGenFailed, "LLVM generation failed"},
		{"ErrWriteIRFile", codegen.ErrWriteIRFile, "write IR file failed"},
		{"ErrCompileToObj", codegen.ErrCompileToObj, "compile to object failed"},
		{"ErrLinkExecutable", codegen.ErrLinkExecutable, "link executable failed"},
		{"ErrToolNotFound", codegen.ErrToolNotFound, "tool not found"},
		{"ErrNoSuitableCompiler", codegen.ErrNoSuitableCompiler, "no suitable compiler found"},
		{"ErrPrintComplexCall", codegen.ErrPrintComplexCall, "print complex call error"},
		{"ErrPrintConvertError", codegen.ErrPrintConvertError, "print convert error"},
		{"ErrPrintDetermineError", codegen.ErrPrintDetermineError, "print determine error"},
		{"ErrUnsupportedStatement", codegen.ErrUnsupportedStatement, "unsupported statement"},
		{"ErrUndefinedVariable", codegen.ErrUndefinedVariable, "undefined variable"},
		{"ErrUnsupportedBinaryOp", codegen.ErrUnsupportedBinaryOp, "unsupported binary operator"},
		{"ErrRangeWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrForEachWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrBuiltInTwoArgs", codegen.ErrBuiltInTwoArgs, "built-in function expects two arguments"},
		{"ErrMapWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrBuiltInRedefine", codegen.ErrBuiltInRedefine, "cannot redefine built-in function"},
		{"ErrFilterWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
		{"ErrFunctionNotFound", codegen.ErrFunctionNotFound, "function not found"},
		{"ErrFoldWrongArgs", codegen.ErrWrongArgCount, "wrong argument count"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.err.Error() != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, test.err.Error())
			}
		})
	}
}

func TestStaticErrors(t *testing.T) {
	// Test that static errors are defined
	staticErrors := []error{
		codegen.ErrToStringReserved,
		codegen.ErrUnsupportedStatement,
		codegen.ErrFunctionNotDeclared,
		codegen.ErrUndefinedVariable,
		codegen.ErrUnsupportedExpression,
		codegen.ErrUnsupportedBinaryOp,
		codegen.ErrWrongArgCount, // consolidated all WrongArgs errors into this one
		codegen.ErrUnsupportedCall,
		codegen.ErrMethodNotImpl,
		codegen.ErrNoToStringImpl,
		codegen.ErrNoToStringForFunc,
		codegen.ErrPrintCannotConvert,
		codegen.ErrPrintComplexExpr,
		codegen.ErrPrintUnknownFunc,
		codegen.ErrFunctionRequiresNamed,
		codegen.ErrMissingArgument,
		codegen.ErrParseErrors,
		codegen.ErrParseTreeNil,
		codegen.ErrASTBuildFailed,
		codegen.ErrLLVMGenFailed,
		codegen.ErrWriteIRFile,
		codegen.ErrCompileToObj,
		codegen.ErrLinkExecutable,
		codegen.ErrToolNotFound,
		codegen.ErrNoSuitableCompiler,
		codegen.ErrPrintComplexCall,
		codegen.ErrPrintConvertError,
		codegen.ErrPrintDetermineError,
		codegen.ErrForEachNotFunction,
		codegen.ErrMapNotFunction,
		codegen.ErrFilterNotFunction,
		codegen.ErrFoldNotFunction,
		codegen.ErrInputNoArgs,
		codegen.ErrBuiltInTwoArgs,
		codegen.ErrBuiltInRedefine,
		codegen.ErrFunctionNotFound,
	}

	for _, err := range staticErrors {
		if err == nil {
			t.Error("Static error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Static error should have non-empty message")
		}
	}
}

func TestWrapUnsupportedStatement(t *testing.T) {
	type testStruct struct{}
	err := codegen.WrapUnsupportedStatement(&testStruct{})

	if err == nil {
		t.Error("WrapUnsupportedStatement should return error")
	}
	if !errors.Is(err, codegen.ErrUnsupportedStatement) {
		t.Error("Should wrap ErrUnsupportedStatement")
	}
	if !strings.Contains(err.Error(), "testStruct") {
		t.Error("Should contain type information")
	}
}

func TestWrapFunctionNotDeclared(t *testing.T) {
	err := codegen.WrapFunctionNotDeclared("testFunc")

	if err == nil {
		t.Error("WrapFunctionNotDeclared should return error")
	}
	if !errors.Is(err, codegen.ErrFunctionNotDeclared) {
		t.Error("Should wrap ErrFunctionNotDeclared")
	}
	if !strings.Contains(err.Error(), "testFunc") {
		t.Error("Should contain function name")
	}
}

func TestWrapUndefinedVariable(t *testing.T) {
	err := codegen.WrapUndefinedVariable("testVar")

	if err == nil {
		t.Error("WrapUndefinedVariable should return error")
	}
	if !errors.Is(err, codegen.ErrUndefinedVariable) {
		t.Error("Should wrap ErrUndefinedVariable")
	}
	if !strings.Contains(err.Error(), "testVar") {
		t.Error("Should contain variable name")
	}
}

func TestWrapUnsupportedExpression(t *testing.T) {
	type testExpr struct{}
	err := codegen.WrapUnsupportedExpression(&testExpr{})

	if err == nil {
		t.Error("WrapUnsupportedExpression should return error")
	}
	if !errors.Is(err, codegen.ErrUnsupportedExpression) {
		t.Error("Should wrap ErrUnsupportedExpression")
	}
	if !strings.Contains(err.Error(), "testExpr") {
		t.Error("Should contain type information")
	}
}

func TestWrapUnsupportedBinaryOp(t *testing.T) {
	err := codegen.WrapUnsupportedBinaryOp("@@")

	if err == nil {
		t.Error("WrapUnsupportedBinaryOp should return error")
	}
	if !errors.Is(err, codegen.ErrUnsupportedBinaryOp) {
		t.Error("Should wrap ErrUnsupportedBinaryOp")
	}
	if !strings.Contains(err.Error(), "@@") {
		t.Error("Should contain operator")
	}
}

func TestWrapToStringWrongArgs(t *testing.T) {
	err := codegen.WrapToStringWrongArgs(3)

	if err == nil {
		t.Error("WrapToStringWrongArgs should return error")
	}
	if !errors.Is(err, codegen.ErrWrongArgCount) {
		t.Error("Should wrap ErrWrongArgCount")
	}
	if !strings.Contains(err.Error(), "3") {
		t.Error("Should contain argument count")
	}
}

func TestWrapPrintWrongArgs(t *testing.T) {
	err := codegen.WrapPrintWrongArgs(2)

	if err == nil {
		t.Error("WrapPrintWrongArgs should return error")
	}
	if !errors.Is(err, codegen.ErrWrongArgCount) {
		t.Error("Should wrap ErrWrongArgCount")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Error("Should contain argument count")
	}
}

func TestWrapInputWrongArgs(t *testing.T) {
	err := codegen.WrapWrongArgCount("input", 0, 1)

	if err == nil {
		t.Error("WrapWrongArgCount should return error")
	}
	if !errors.Is(err, codegen.ErrWrongArgCount) {
		t.Error("Should wrap ErrWrongArgCount")
	}
	if !strings.Contains(err.Error(), "1") {
		t.Error("Should contain argument count")
	}
}

func TestWrapWrongArgCount(t *testing.T) {
	err := codegen.WrapWrongArgCount("testFunc", 2, 3)

	if err == nil {
		t.Error("WrapWrongArgCount should return error")
	}
	if !errors.Is(err, codegen.ErrWrongArgCount) {
		t.Error("Should wrap ErrWrongArgCount")
	}
	if !strings.Contains(err.Error(), "testFunc") {
		t.Error("Should contain function name")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Error("Should contain expected count")
	}
	if !strings.Contains(err.Error(), "3") {
		t.Error("Should contain actual count")
	}
}

func TestWrapMissingArgument(t *testing.T) {
	err := codegen.WrapMissingArgument("param1", "testFunc")

	if err == nil {
		t.Error("WrapMissingArgument should return error")
	}
	if !errors.Is(err, codegen.ErrMissingArgument) {
		t.Error("Should wrap ErrMissingArgument")
	}
	if !strings.Contains(err.Error(), "param1") {
		t.Error("Should contain parameter name")
	}
	if !strings.Contains(err.Error(), "testFunc") {
		t.Error("Should contain function name")
	}
}

func TestWrapParseErrors(t *testing.T) {
	err := codegen.WrapParseErrors([]string{"line 1:5: syntax error at 'unexpected token'"})

	if err == nil {
		t.Error("WrapParseErrors should return error")
	}
	// ParseError is a custom error type, not a wrapped ErrParseErrors
	if !strings.Contains(err.Error(), "line 1:5: syntax error") {
		t.Error("Should contain error details")
	}
}

func TestWrapLLVMGenFailed(t *testing.T) {
	baseErr := codegen.ErrLLVMGenFailed
	err := codegen.WrapLLVMGenFailed(baseErr)

	if err == nil {
		t.Error("WrapLLVMGenFailed should return error")
	}
	if !errors.Is(err, codegen.ErrLLVMGenFailed) {
		t.Error("Should wrap ErrLLVMGenFailed")
	}
	if !errors.Is(err, baseErr) {
		t.Error("Should preserve inner error")
	}
}

func TestWrapWriteIRFile(t *testing.T) {
	baseErr := codegen.ErrWriteIRFile
	err := codegen.WrapWriteIRFile(baseErr)

	if err == nil {
		t.Error("WrapWriteIRFile should return error")
	}
	if !errors.Is(err, codegen.ErrWriteIRFile) {
		t.Error("Should wrap ErrWriteIRFile")
	}
	if !errors.Is(err, baseErr) {
		t.Error("Should preserve inner error")
	}
}

func TestWrapCompileToObj(t *testing.T) {
	baseErr := codegen.ErrCompileToObj
	err := codegen.WrapCompileToObj(baseErr, "llc output here")

	if err == nil {
		t.Error("WrapCompileToObj should return error")
	}
	if !errors.Is(err, codegen.ErrCompileToObj) {
		t.Error("Should wrap ErrCompileToObj")
	}
	if !strings.Contains(err.Error(), "llc output here") {
		t.Error("Should contain llc output")
	}
}

func TestWrapLinkExecutable(t *testing.T) {
	baseErr := codegen.ErrLinkExecutable
	err := codegen.WrapLinkExecutable("gcc", baseErr, "linker output")

	if err == nil {
		t.Error("WrapLinkExecutable should return error")
	}
	if !errors.Is(err, codegen.ErrLinkExecutable) {
		t.Error("Should wrap ErrLinkExecutable")
	}
	if !strings.Contains(err.Error(), "gcc") {
		t.Error("Should contain compiler name")
	}
	if !strings.Contains(err.Error(), "linker output") {
		t.Error("Should contain linker output")
	}
}

func TestWrapToolNotFound(t *testing.T) {
	err := codegen.WrapToolNotFound("llc")

	if err == nil {
		t.Error("WrapToolNotFound should return error")
	}
	if !errors.Is(err, codegen.ErrToolNotFound) {
		t.Error("Should wrap ErrToolNotFound")
	}
	if !strings.Contains(err.Error(), "llc") {
		t.Error("Should contain tool name")
	}
}

func TestWrapNoSuitableCompiler(t *testing.T) {
	compilers := []string{"gcc", "clang", "cc"}
	err := codegen.WrapNoSuitableCompiler(compilers)

	if err == nil {
		t.Error("WrapNoSuitableCompiler should return error")
	}
	if !errors.Is(err, codegen.ErrNoSuitableCompiler) {
		t.Error("Should wrap ErrNoSuitableCompiler")
	}
	if !strings.Contains(err.Error(), "gcc") {
		t.Error("Should contain compiler list")
	}
}

func TestWrapPrintConvertError(t *testing.T) {
	err := codegen.WrapPrintConvertError("CustomType", "myFunc")

	if err == nil {
		t.Error("WrapPrintConvertError should return error")
	}
	if !errors.Is(err, codegen.ErrPrintConvertError) {
		t.Error("Should wrap ErrPrintConvertError")
	}
	if !strings.Contains(err.Error(), "CustomType") {
		t.Error("Should contain return type")
	}
	if !strings.Contains(err.Error(), "myFunc") {
		t.Error("Should contain function name")
	}
}

func TestWrapPrintDetermineError(t *testing.T) {
	err := codegen.WrapPrintDetermineError("unknownFunc")

	if err == nil {
		t.Error("WrapPrintDetermineError should return error")
	}
	if !errors.Is(err, codegen.ErrPrintDetermineError) {
		t.Error("Should wrap ErrPrintDetermineError")
	}
	if !strings.Contains(err.Error(), "unknownFunc") {
		t.Error("Should contain function name")
	}
}

func TestWrapBuiltInRedefine(t *testing.T) {
	err := codegen.WrapBuiltInRedefine("print")

	if err == nil {
		t.Error("WrapBuiltInRedefine should return error")
	}
	if !errors.Is(err, codegen.ErrBuiltInRedefine) {
		t.Error("Should wrap ErrBuiltInRedefine")
	}
	if !strings.Contains(err.Error(), "print") {
		t.Error("Should contain function name")
	}
}

func TestWrapIteratorErrors(t *testing.T) {
	tests := []struct {
		name     string
		wrapFunc func(int) error
		baseErr  error
	}{
		{"WrapMapWrongArgs", codegen.WrapMapWrongArgs, codegen.ErrWrongArgCount},
		{"WrapFilterWrongArgs", codegen.WrapFilterWrongArgs, codegen.ErrWrongArgCount},
		{"WrapFoldWrongArgs", codegen.WrapFoldWrongArgs, codegen.ErrWrongArgCount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.wrapFunc(5)

			if err == nil {
				t.Errorf("%s should return error", test.name)
			}
			if !errors.Is(err, test.baseErr) {
				t.Errorf("%s should wrap %v", test.name, test.baseErr)
			}
			if !strings.Contains(err.Error(), "5") {
				t.Errorf("%s should contain argument count", test.name)
			}
		})
	}
}

func TestWrapBuiltInTwoArgs(t *testing.T) {
	err := codegen.WrapBuiltInTwoArgs("input")

	if err == nil {
		t.Error("WrapBuiltInTwoArgs should return error")
	}
	if !errors.Is(err, codegen.ErrBuiltInTwoArgs) {
		t.Error("Should wrap ErrBuiltInTwoArgs")
	}
	if !strings.Contains(err.Error(), "input") {
		t.Error("Should contain function name")
	}
}

func TestWrapFunctionNotFound(t *testing.T) {
	err := codegen.WrapFunctionNotFound("missingFunc")

	if err == nil {
		t.Error("WrapFunctionNotFound should return error")
	}
	if !errors.Is(err, codegen.ErrFunctionNotFound) {
		t.Error("Should wrap ErrFunctionNotFound")
	}
	if !strings.Contains(err.Error(), "missingFunc") {
		t.Error("Should contain function name")
	}
}

func TestErrorWrappingWithPosition(t *testing.T) {
	position := &ast.Position{Line: 10, Column: 5}

	tests := []struct {
		name        string
		wrapperFunc func() error
		expectedMsg string
	}{
		{
			"WrapUndefinedVariableWithPos",
			func() error { return codegen.WrapUndefinedVariableWithPos("myVar", position) },
			"line 10:5: undefined variable: myVar",
		},
		{
			"WrapUnsupportedBinaryOpWithPos",
			func() error { return codegen.WrapUnsupportedBinaryOpWithPos("***", position) },
			"line 10:5: unsupported binary operator: ***",
		},
		{
			"WrapVoidArithmeticWithPos",
			func() error { return codegen.WrapVoidArithmeticWithPos("+", position) },
			"line 10:5: arithmetic operation on void type +",
		},
		{
			"WrapUnsupportedUnaryOpWithPos",
			func() error { return codegen.WrapUnsupportedUnaryOpWithPos("***", position) },
			"line 10:5: unsupported unary operator: ***",
		},
		{
			"WrapConstraintResultFieldAccessWithPos",
			func() error { return codegen.WrapConstraintResultFieldAccessWithPos("myField", position) },
			"line 10:5: constraint result field access not allowed: myField",
		},
		{
			"WrapFunctionArgsWithPos print",
			func() error { return codegen.WrapFunctionArgsWithPos("print", 1, 2, position) },
			"line 10:5: print expects exactly 1 argument(s), got 2: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos input",
			func() error { return codegen.WrapFunctionArgsWithPos("input", 0, 1, position) },
			"line 10:5: input expects exactly 0 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos length",
			func() error { return codegen.WrapFunctionArgsWithPos("length", 1, 2, position) },
			"line 10:5: length expects exactly 1 argument(s), got 2: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos contains",
			func() error { return codegen.WrapFunctionArgsWithPos("contains", 2, 1, position) },
			"line 10:5: contains expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos substring",
			func() error { return codegen.WrapFunctionArgsWithPos("substring", 3, 2, position) },
			"line 10:5: substring expects exactly 3 argument(s), got 2: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos range",
			func() error { return codegen.WrapFunctionArgsWithPos("range", 2, 1, position) },
			"line 10:5: range expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos forEach",
			func() error { return codegen.WrapFunctionArgsWithPos("forEach", 2, 1, position) },
			"line 10:5: forEach expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos httpCreateClient",
			func() error { return codegen.WrapFunctionArgsWithPos("httpCreateClient", 2, 1, position) },
			"line 10:5: httpCreateClient expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos httpGet",
			func() error { return codegen.WrapFunctionArgsWithPos("httpGet", 3, 2, position) },
			"line 10:5: httpGet expects exactly 3 argument(s), got 2: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos spawnProcess",
			func() error { return codegen.WrapFunctionArgsWithPos("spawnProcess", 2, 1, position) },
			"line 10:5: spawnProcess expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos writeFile",
			func() error { return codegen.WrapFunctionArgsWithPos("writeFile", 2, 1, position) },
			"line 10:5: writeFile expects exactly 2 argument(s), got 1: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos readFile",
			func() error { return codegen.WrapFunctionArgsWithPos("readFile", 1, 2, position) },
			"line 10:5: readFile expects exactly 1 argument(s), got 2: wrong argument count",
		},
		{
			"WrapImmutableAssignmentErrorWithPos",
			func() error { return codegen.WrapImmutableAssignmentErrorWithPos("myVar", position) },
			"line 10:5: cannot assign to immutable variable: myVar",
		},
		{
			"WrapMissingArgumentWithPos",
			func() error { return codegen.WrapMissingArgumentWithPos("argName", "funcName", position) },
			"line 10:5: missing argument: argName for function funcName",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.wrapperFunc()
			if err.Error() != test.expectedMsg {
				t.Errorf("Expected %q, got %q", test.expectedMsg, err.Error())
			}
		})
	}
}

// Test error wrapping without position
func TestErrorWrappingWithoutPosition(t *testing.T) {
	// Use existing static error for testing
	errTest := codegen.ErrLLVMGenFailed

	tests := []struct {
		name        string
		wrapperFunc func() error
		expectedMsg string
	}{
		{
			"WrapUndefinedVariableWithPos nil position",
			func() error { return codegen.WrapUndefinedVariableWithPos("myVar", nil) },
			"undefined variable: myVar",
		},
		{
			"WrapFunctionArgsWithPos nil position for print",
			func() error { return codegen.WrapFunctionArgsWithPos("print", 1, 2, nil) },
			"function print expects 1 arguments, got 2: wrong argument count",
		},
		{
			"WrapFunctionArgsWithPos nil position for length",
			func() error { return codegen.WrapFunctionArgsWithPos("length", 1, 2, nil) },
			"function length expects 1 arguments, got 2: wrong argument count",
		},
		{
			"WrapToStringWrongArgs",
			func() error { return codegen.WrapToStringWrongArgs(2) },
			"function toString expects 1 arguments, got 2: wrong argument count",
		},
		{
			"WrapMissingArgument",
			func() error { return codegen.WrapMissingArgumentWithPos("argName", "funcName", nil) },
			"missing argument: argName for function funcName",
		},
		{
			"WrapLLVMGenFailed",
			func() error { return codegen.WrapLLVMGenFailed(errTest) },
			"LLVM generation failed: LLVM generation failed",
		},
		{
			"WrapCompileToObj",
			func() error { return codegen.WrapCompileToObj(errTest, "output") },
			"compile to object failed: LLVM generation failed - output",
		},
		{
			"WrapLinkExecutable",
			func() error { return codegen.WrapLinkExecutable("gcc", errTest, "output") },
			"link executable failed: gcc failed - LLVM generation failed - output",
		},
		{
			"WrapPrintConvertError",
			func() error { return codegen.WrapPrintConvertError("int", "myFunc") },
			"print convert error: cannot convert return type int for function myFunc to string",
		},
		{
			"WrapPrintDetermineError",
			func() error { return codegen.WrapPrintDetermineError("myFunc") },
			"print determine error: cannot determine return type for function myFunc",
		},
		{
			"WrapUnsupportedBinaryOp string",
			func() error { return codegen.WrapUnsupportedBinaryOp("***") },
			"unsupported binary operator: ***",
		},
		{
			"WrapUnsupportedCallExpressionSecurity",
			func() error { return codegen.WrapUnsupportedCallExpressionSecurity("myFunc") },
			"unsupported call expression in current security mode: myFunc",
		},
		{
			"WrapMethodCallNotImplemented",
			func() error { return codegen.WrapMethodCallNotImplemented("myMethod") },
			"method call not implemented: myMethod",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.wrapperFunc()
			if err.Error() != test.expectedMsg {
				t.Errorf("Expected %q, got %q", test.expectedMsg, err.Error())
			}
		})
	}
}

// Test if error is wrapped by checking if it's an instance of the base error
func TestErrorWrapping(t *testing.T) {
	baseErr := codegen.ErrWrongArgCount
	wrappedErr := codegen.WrapToStringWrongArgs(2)

	if !errors.Is(wrappedErr, baseErr) {
		t.Errorf("Expected wrapped error to be instance of %v", baseErr)
	}
}
