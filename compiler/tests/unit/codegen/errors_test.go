package codegen

import (
	"errors"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestStaticErrors(t *testing.T) {
	// Test that static errors are defined
	staticErrors := []error{
		codegen.ErrToStringReserved,
		codegen.ErrUnsupportedStatement,
		codegen.ErrFunctionNotDeclared,
		codegen.ErrUndefinedVariable,
		codegen.ErrUnsupportedExpression,
		codegen.ErrUnsupportedBinaryOp,
		codegen.ErrFieldAccessNotImpl,
		codegen.ErrToStringWrongArgs,
		codegen.ErrPrintWrongArgs,
		codegen.ErrInputWrongArgs,
		codegen.ErrUnsupportedCall,
		codegen.ErrMethodNotImpl,
		codegen.ErrNoToStringImpl,
		codegen.ErrNoToStringForFunc,
		codegen.ErrPrintCannotConvert,
		codegen.ErrPrintComplexExpr,
		codegen.ErrPrintUnknownFunc,
		codegen.ErrFunctionRequiresNamed,
		codegen.ErrWrongArgCount,
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
		codegen.ErrRangeWrongArgs,
		codegen.ErrForEachWrongArgs,
		codegen.ErrForEachNotFunction,
		codegen.ErrMapWrongArgs,
		codegen.ErrMapNotFunction,
		codegen.ErrFilterWrongArgs,
		codegen.ErrFilterNotFunction,
		codegen.ErrFoldWrongArgs,
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

func TestWrapFieldAccessNotImpl(t *testing.T) {
	err := codegen.WrapFieldAccessNotImpl("testField")

	if err == nil {
		t.Error("WrapFieldAccessNotImpl should return error")
	}
	if !errors.Is(err, codegen.ErrFieldAccessNotImpl) {
		t.Error("Should wrap ErrFieldAccessNotImpl")
	}
	if !strings.Contains(err.Error(), "testField") {
		t.Error("Should contain field name")
	}
}

func TestWrapToStringWrongArgs(t *testing.T) {
	err := codegen.WrapToStringWrongArgs(3)

	if err == nil {
		t.Error("WrapToStringWrongArgs should return error")
	}
	if !errors.Is(err, codegen.ErrToStringWrongArgs) {
		t.Error("Should wrap ErrToStringWrongArgs")
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
	if !errors.Is(err, codegen.ErrPrintWrongArgs) {
		t.Error("Should wrap ErrPrintWrongArgs")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Error("Should contain argument count")
	}
}

func TestWrapInputWrongArgs(t *testing.T) {
	err := codegen.WrapInputWrongArgs(1)

	if err == nil {
		t.Error("WrapInputWrongArgs should return error")
	}
	if !errors.Is(err, codegen.ErrInputWrongArgs) {
		t.Error("Should wrap ErrInputWrongArgs")
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
	if !errors.Is(err, codegen.ErrParseErrors) {
		t.Error("Should wrap ErrParseErrors")
	}
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
		{"WrapRangeWrongArgs", codegen.WrapRangeWrongArgs, codegen.ErrRangeWrongArgs},
		{"WrapForEachWrongArgs", codegen.WrapForEachWrongArgs, codegen.ErrForEachWrongArgs},
		{"WrapMapWrongArgs", codegen.WrapMapWrongArgs, codegen.ErrMapWrongArgs},
		{"WrapFilterWrongArgs", codegen.WrapFilterWrongArgs, codegen.ErrFilterWrongArgs},
		{"WrapFoldWrongArgs", codegen.WrapFoldWrongArgs, codegen.ErrFoldWrongArgs},
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

func TestErrorChaining(t *testing.T) {
	// Test that error unwrapping works correctly
	baseErr := codegen.ErrLLVMGenFailed
	wrappedErr := codegen.WrapLLVMGenFailed(baseErr)

	if !errors.Is(wrappedErr, baseErr) {
		t.Error("Should be able to unwrap to inner error")
	}
	if !errors.Is(wrappedErr, codegen.ErrLLVMGenFailed) {
		t.Error("Should be able to unwrap to wrapper error")
	}
}
