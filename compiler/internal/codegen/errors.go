package codegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/christianfindlay/osprey/internal/ast"
)

// Static error definitions for the codegen package
var (
	// Variable-related errors
	ErrUndefinedIdentifier  = errors.New("undefined identifier")
	ErrVariableNotInRuntime = errors.New("variable exists in type environment but not in runtime")
	ErrUnknownParameterName = errors.New("unknown parameter name")

	// Type-related errors
	ErrTypeMismatch              = errors.New("type mismatch")
	ErrAnyTypeMismatch           = errors.New("cannot pass 'any' type to function expecting specific type")
	ErrRecursiveType             = errors.New("recursive type detected")
	ErrNotTypeVariable           = errors.New("not a type variable")
	ErrNotConcreteType           = errors.New("not a concrete type")
	ErrNotFunctionType           = errors.New("not a function type")
	ErrFunctionArityMismatch     = errors.New("function arity mismatch")
	ErrUnsupportedType           = errors.New("unsupported type")
	ErrUnsupportedExpression     = errors.New("unsupported expression")
	ErrUnsupportedBinaryOperator = errors.New("unsupported binary operator")
	ErrUnsupportedUnaryOperator  = errors.New("unsupported unary operator")
	ErrMethodCallsNotImplemented = errors.New("method calls are not implemented")

	// Result constructor errors
	ErrSuccessConstructorMissingValue = errors.New("success constructor requires 'value' field")
	ErrErrorConstructorMissingMessage = errors.New("error constructor requires 'message' field")

	// Function-related errors
	ErrFunctionNotDeclared = errors.New("function not declared")
	ErrNotAFunction        = errors.New("type is not a function")

	// Union/variant-related errors
	ErrNoVariantFound     = errors.New("no variant found matching field structure")
	ErrNoVariantsFound    = errors.New("no variants found for type")
	ErrVariantNotFound    = errors.New("variant not found")
	ErrUnknownConstructor = errors.New("unknown constructor")

	// Pattern matching errors
	ErrInvalidEmptyPattern = errors.New("invalid empty pattern")
	ErrMatchNoArms         = errors.New("match expression must have at least one arm")

	// Field access errors
	ErrCannotExtractField         = errors.New("cannot extract field from non-struct type")
	ErrDiscriminantNotPointer     = errors.New("discriminant is not a pointer type")
	ErrDiscriminantNotTaggedUnion = errors.New("discriminant is not a tagged union")
	ErrFieldAccessOnInteger       = errors.New("cannot access field on integer type")
	ErrFieldAccessOnString        = errors.New("cannot access field on string type")
	ErrFieldAccessOnBoolean       = errors.New("cannot access field on boolean type")
	ErrFieldAccessOnNonStruct     = errors.New("cannot access field on non-struct type")
	ErrFieldAccessOnUnknownStruct = errors.New("cannot access field on unknown struct type")
	ErrFieldNotFoundInRecord      = errors.New("field not found in record type")
	ErrFieldAccessOnNonRecord     = errors.New("cannot access field on non-record type")
	ErrFieldAccessOnLegacyRecord  = errors.New("field access on legacy record type not supported")
	ErrFieldNotFoundOnType        = errors.New("field not found on type")

	// Parse errors
	ErrParseTreeNil   = errors.New("parse tree is nil")
	ErrASTBuildFailed = errors.New("AST build failed")

	// Core function errors
	ErrPrintCannotConvert = errors.New("cannot convert value for printing")
	ErrUnsupportedCall    = errors.New("unsupported function call")
	ErrToStringReserved   = errors.New("toString is a reserved function name")

	// WebSocket errors
	ErrWebSocketKeepAliveWrongArgs = errors.New("websocketKeepAlive function has wrong number of arguments")

	// Iterator errors
	ErrForEachNotFunction = errors.New("forEach callback is not a function")
	ErrMapNotFunction     = errors.New("map callback is not a function")
	ErrFilterNotFunction  = errors.New("filter callback is not a function")
	ErrFoldNotFunction    = errors.New("fold callback is not a function")

	// System errors
	ErrInputNoArgs    = errors.New("input function takes no arguments")
	ErrNoToStringImpl = errors.New("no toString implementation found")

	// Collection access errors
	ErrUnsupportedCollectionType = errors.New("unsupported collection type for access")

	// Additional error constants expected by tests
	ErrUnsupportedStatement = errors.New("unsupported statement")
	ErrUndefinedVariable    = errors.New("undefined variable")
	ErrUnsupportedBinaryOp  = errors.New("unsupported binary operator") // Alias for ErrUnsupportedBinaryOperator

	ErrMethodNotImpl         = errors.New("method not implemented")
	ErrNoToStringForFunc     = errors.New("no toString implementation for function")
	ErrPrintComplexExpr      = errors.New("print complex expression error")
	ErrPrintUnknownFunc      = errors.New("print unknown function error")
	ErrFunctionRequiresNamed = errors.New("function requires named arguments")
	ErrWrongArgCount         = errors.New("wrong argument count")
	ErrMissingArgument       = errors.New("missing argument")
	ErrParseErrors           = errors.New("parse errors")
	ErrLLVMGenFailed         = errors.New("LLVM generation failed")
	ErrWriteIRFile           = errors.New("write IR file failed")
	ErrCompileToObj          = errors.New("compile to object failed")
	ErrLinkExecutable        = errors.New("link executable failed")
	ErrToolNotFound          = errors.New("tool not found")
	ErrNoSuitableCompiler    = errors.New("no suitable compiler found")
	ErrPrintComplexCall      = errors.New("print complex call error")
	ErrPrintConvertError     = errors.New("print convert error")
	ErrPrintDetermineError   = errors.New("print determine error")

	ErrBuiltInTwoArgs   = errors.New("built-in function expects two arguments")
	ErrBuiltInRedefine  = errors.New("cannot redefine built-in function")
	ErrFunctionNotFound = errors.New("function not found")

	// Consolidated error constants
	ErrVoidArithmetic                = errors.New("arithmetic operation on void type")
	ErrFieldAccessOnResult           = errors.New("field access on result type")
	ErrConstraintResultFieldAccess   = errors.New("constraint result field access not allowed")
	ErrUndefinedType                 = errors.New("undefined type")
	ErrMissingField                  = errors.New("missing field")
	ErrHTTPStopServerUnknownNamedArg = errors.New("unknown named argument for httpStopServer")
	ErrImmutableAssignmentError      = errors.New("cannot assign to immutable variable")
	ErrArithmeticTypeMismatch        = errors.New("arithmetic operand type mismatch")

	// Security-related errors
	ErrUnsupportedCallExpressionSecurity = errors.New("unsupported call expression in current security mode")
	ErrMethodCallNotImplemented          = errors.New("method call not implemented")

	// Fiber-related errors
	ErrAwaitTypeMismatch = errors.New("await can only be used on Fiber types")

	// Channel-related errors
	ErrChannelSendFunctionNotFound  = errors.New("channel send function not found")
	ErrChannelRecvFunctionNotFound  = errors.New("channel recv function not found")
	ErrChannelCreateInvalidArgCount = errors.New("channel create invalid argument count")

	// Record type errors
	ErrRecordFieldTypeMismatch = errors.New("record field type mismatch")
)

// Static error definitions for match expressions
var (
	ErrMatchNotExhaustive = errors.New("match expression not exhaustive")
	ErrMatchTypeMismatch  = errors.New("match expression type mismatch")
	ErrUnknownVariant     = errors.New("unknown variant in match expression")
	ErrNilOperand         = errors.New("arithmetic operation has nil operand")
)

// NOTE: Function argument count constants have been moved to the unified built-in function registry
// (builtin_registry.go). These constants are no longer needed as the registry contains all
// built-in function metadata in one place. Use len(GlobalBuiltInRegistry.GetFunction(name).ParameterTypes) instead.

// Error wrapper functions

// ParseError represents multiple parse errors
type ParseError struct {
	Errors []string
}

// Error implements the error interface
func (pe *ParseError) Error() string {
	return strings.Join(pe.Errors, "\n")
}

// WrapParseErrors wraps parse errors from the parser
func WrapParseErrors(parseErrors []string) error {
	if len(parseErrors) == 0 {
		return fmt.Errorf("%w: unknown parse error", ErrParseErrors)
	}
	// Return a custom error type to avoid dynamic error creation
	return &ParseError{Errors: parseErrors}
}

// WrapWriteIRFile wraps errors when writing IR files
func WrapWriteIRFile(err error) error {
	return fmt.Errorf("failed to write IR file: %w", err)
}

// WrapToStringWrongArgs wraps errors for wrong number of toString arguments
func WrapToStringWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(ToStringFunc)
	return WrapWrongArgCount(ToStringFunc, len(fn.ParameterTypes), argCount)
}

// WrapNoToStringImpl wraps errors when no toString implementation is found
func WrapNoToStringImpl(typeName string) error {
	return fmt.Errorf("%w for type %s", ErrNoToStringImpl, typeName)
}

// WrapFunctionArgsWithPos is a consolidated function for all "wrong arguments" errors
func WrapFunctionArgsWithPos(funcName string, expected int, actual int, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		// Try to get parameter names from the built-in function registry
		paramNames := ""

		if fn, exists := GlobalBuiltInRegistry.GetFunction(funcName); exists {
			var names []string
			for _, param := range fn.ParameterTypes {
				names = append(names, param.Name)
			}

			if len(names) > 0 {
				paramNames = fmt.Sprintf(" (%s)", strings.Join(names, ", "))
			}
		}

		//nolint:err113 // Dynamic error needed for exact test format matching
		return fmt.Errorf("line %d:%d: %s expects exactly %d arguments%s, got %d",
			position.Line, position.Column, funcName, expected, paramNames, actual)
	}

	return WrapWrongArgCount(funcName, expected, actual)
}

// WrapUnsupportedExpression wraps errors for unsupported expressions
func WrapUnsupportedExpression(expr interface{}) error {
	return fmt.Errorf("%w: %T", ErrUnsupportedExpression, expr)
}

// WrapUndefinedVariableWithPos wraps errors for undefined variables
func WrapUndefinedVariableWithPos(varName string, pos interface{}) error {
	return WrapSimpleErrorWithPos(ErrUndefinedVariable, varName, pos)
}

// WrapUnsupportedBinaryOpWithPos wraps errors for unsupported binary operators
func WrapUnsupportedBinaryOpWithPos(op string, pos interface{}) error {
	return WrapSimpleErrorWithPos(ErrUnsupportedBinaryOp, op, pos)
}

// WrapVoidArithmeticWithPos wraps errors for arithmetic on void types
func WrapVoidArithmeticWithPos(op string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w %s", position.Line, position.Column, ErrVoidArithmetic, op)
	}

	return fmt.Errorf("%w %s", ErrVoidArithmetic, op)
}

// WrapUnsupportedUnaryOpWithPos wraps errors for unsupported unary operators
func WrapUnsupportedUnaryOpWithPos(op string, pos interface{}) error {
	return WrapSimpleErrorWithPos(ErrUnsupportedUnaryOperator, op, pos)
}

// WrapFieldAccessOnResult wraps errors for field access on result types
func WrapFieldAccessOnResult(field string, resultType string) error {
	return fmt.Errorf("%w %s on result type %s", ErrFieldAccessOnResult, field, resultType)
}

// WrapConstraintResultFieldAccessWithPos wraps errors for constraint result field access
func WrapConstraintResultFieldAccessWithPos(field string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s", position.Line, position.Column, ErrConstraintResultFieldAccess, field)
	}

	return fmt.Errorf("%w: %s", ErrConstraintResultFieldAccess, field)
}

// Removed - Use WrapMethodCallNotImplemented directly

// WrapFieldNotFoundInRecord wraps the field not found error with additional context
func WrapFieldNotFoundInRecord(field, recordType string) error {
	return fmt.Errorf("%w '%s' in record type %s", ErrFieldNotFoundInRecord, field, recordType)
}

// WrapFieldAccessOnNonRecord wraps the non-record field access error with context
func WrapFieldAccessOnNonRecord(field, typeStr string) error {
	return fmt.Errorf("%w '%s' on non-record type %s", ErrFieldAccessOnNonRecord, field, typeStr)
}

// WrapFieldAccessOnLegacyRecord wraps the legacy record field access error
func WrapFieldAccessOnLegacyRecord(field, typeName string) error {
	return fmt.Errorf("%w, field '%s' on type %s", ErrFieldAccessOnLegacyRecord, field, typeName)
}

// WrapUndefinedType wraps errors for undefined types
func WrapUndefinedType(typeName string) error {
	return WrapSimpleError(ErrUndefinedType, typeName)
}

// WrapUndefinedTypeWithPos wraps errors for undefined types with position
func WrapUndefinedTypeWithPos(typeName string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		//nolint:err113 // Dynamic error needed for exact test format matching
		return fmt.Errorf("line %d:%d: undefined type: %s",
			position.Line, position.Column, typeName)
	}
	//nolint:err113 // Dynamic error needed for exact test format matching
	return fmt.Errorf("undefined type: %s", typeName)
}

// WrapUndefinedFunction wraps errors for undefined functions
func WrapUndefinedFunction(funcName string) error {
	return WrapSimpleError(ErrFunctionNotFound, funcName)
}

// WrapMissingField wraps errors for missing fields
func WrapMissingField(field string) error {
	return WrapSimpleError(ErrMissingField, field)
}

// WrapBuiltInRedefine wraps errors for redefining built-in functions
func WrapBuiltInRedefine(funcName string) error {
	return WrapSimpleError(ErrBuiltInRedefine, funcName)
}

// WrapBuiltInFunctionWrongArgs wraps errors for wrong number of arguments to any built-in function
func WrapBuiltInFunctionWrongArgs(functionName string, argCount int) error {
	fn, exists := GlobalBuiltInRegistry.GetFunction(functionName)
	if !exists {
		return WrapWrongArgCount(functionName, 0, argCount) // Fallback for unknown functions
	}

	return WrapWrongArgCount(functionName, len(fn.ParameterTypes), argCount)
}

// Generic wrapper functions to reduce duplication

// WrapSimpleError wraps a base error with a single string argument
func WrapSimpleError(baseErr error, arg string) error {
	return fmt.Errorf("%w: %s", baseErr, arg)
}

// WrapSimpleErrorWithPos wraps a base error with position and single string argument
func WrapSimpleErrorWithPos(baseErr error, arg string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s", position.Line, position.Column, baseErr, arg)
	}

	return fmt.Errorf("%w: %s", baseErr, arg)
}

// Remaining wrapper functions that are NOT just one-line calls to WrapBuiltInFunctionWrongArgs

// WrapHTTPStopServerUnknownNamedArg wraps errors for unknown named arguments in HTTP stop server
func WrapHTTPStopServerUnknownNamedArg(argName string) error {
	return fmt.Errorf("%w: %s", ErrHTTPStopServerUnknownNamedArg, argName)
}

// WrapToolNotFound wraps errors for tool not found
func WrapToolNotFound(tool string) error {
	return fmt.Errorf("%w: %s", ErrToolNotFound, tool)
}

// WrapNoSuitableCompiler wraps errors for no suitable compiler
func WrapNoSuitableCompiler(compilers []string) error {
	return fmt.Errorf("%w: %v", ErrNoSuitableCompiler, compilers)
}

// WrapUnsupportedStatement wraps errors for unsupported statements
func WrapUnsupportedStatement(stmt interface{}) error {
	return fmt.Errorf("%w: %T", ErrUnsupportedStatement, stmt)
}

// WrapImmutableAssignmentErrorWithPos wraps errors for immutable assignment
func WrapImmutableAssignmentErrorWithPos(varName string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s", position.Line, position.Column, ErrImmutableAssignmentError, varName)
	}

	return fmt.Errorf("%w: %s", ErrImmutableAssignmentError, varName)
}

// Removed - Use WrapUndefinedFunction or WrapSimpleError(ErrFunctionNotFound, funcName) directly

// WrapBuiltInTwoArgs wraps errors for built-in functions expecting two arguments
func WrapBuiltInTwoArgs(funcName string) error {
	return WrapSimpleError(ErrBuiltInTwoArgs, funcName)
}

// Functions removed - duplicate wrapper functions consolidated to WrapBuiltInFunctionWrongArgs

// Additional wrapper functions expected by tests

// WrapFunctionNotDeclared wraps function not declared errors
func WrapFunctionNotDeclared(funcName string) error {
	return WrapSimpleError(ErrFunctionNotDeclared, funcName)
}

// WrapUndefinedVariable wraps undefined variable errors
func WrapUndefinedVariable(varName string) error {
	return WrapSimpleError(ErrUndefinedVariable, varName)
}

// Removed - Use WrapUndefinedFunction or WrapSimpleError(ErrFunctionNotFound, funcName) directly

// WrapHTTPFunctionMissingNamedArg wraps HTTP function missing named argument errors
func WrapHTTPFunctionMissingNamedArg(funcName, argName string) error {
	return fmt.Errorf("%w: %s missing argument %s", ErrMissingArgument, funcName, argName)
}

// Map/Filter/Fold wrapper functions removed - use WrapBuiltInFunctionWrongArgs directly

// WrapWrongArgCount wraps wrong argument count errors
func WrapWrongArgCount(funcName string, expected, actual int) error {
	return fmt.Errorf("function %s expects %d arguments, got %d: %w", funcName, expected, actual, ErrWrongArgCount)
}

// WrapMissingArgument wraps missing argument errors
func WrapMissingArgument(argName string, funcName string) error {
	return fmt.Errorf("%w: %s for function %s", ErrMissingArgument, argName, funcName)
}

// WrapMissingArgumentWithPos wraps missing argument errors with position
func WrapMissingArgumentWithPos(argName string, funcName string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s for function %s",
			position.Line, position.Column, ErrMissingArgument, argName, funcName)
	}

	return fmt.Errorf("%w: %s for function %s", ErrMissingArgument, argName, funcName)
}

// WrapLLVMGenFailed wraps LLVM generation failed errors
func WrapLLVMGenFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrLLVMGenFailed, err)
}

// WrapCompileToObj wraps compile to object errors
func WrapCompileToObj(err error, output string) error {
	return fmt.Errorf("%w: %w - %s", ErrCompileToObj, err, output)
}

// WrapLinkExecutable wraps link executable errors
func WrapLinkExecutable(compiler string, err error, output string) error {
	return fmt.Errorf("%w: %s failed - %w - %s", ErrLinkExecutable, compiler, err, output)
}

// WrapPrintConvertError wraps print convert errors
func WrapPrintConvertError(returnType string, funcName string) error {
	return fmt.Errorf("%w: cannot convert return type %s for function %s to string",
		ErrPrintConvertError, returnType, funcName)
}

// WrapPrintDetermineError wraps print determine errors
func WrapPrintDetermineError(funcName string) error {
	return fmt.Errorf("%w: cannot determine return type for function %s", ErrPrintDetermineError, funcName)
}

// Additional wrapper functions expected by tests

// WrapUnsupportedBinaryOp wraps errors for unsupported binary operators (simplified version)
func WrapUnsupportedBinaryOp(op string) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedBinaryOp, op)
}

// WrapPrintWrongArgs wraps errors for wrong number of print arguments (simplified version)
func WrapPrintWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(PrintFunc)
	return WrapWrongArgCount("print", len(fn.ParameterTypes), argCount)
}

// WrapUnsupportedCallExpressionSecurity wraps errors for security violations
func WrapUnsupportedCallExpressionSecurity(funcName string) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedCallExpressionSecurity, funcName)
}

// WrapUnsupportedCollectionType wraps errors for unsupported collection types
func WrapUnsupportedCollectionType(collectionType string) error {
	return fmt.Errorf("%w: %s", ErrUnsupportedCollectionType, collectionType)
}

// WrapMethodCallNotImplemented wraps errors for method calls not implemented
func WrapMethodCallNotImplemented(method string) error {
	return WrapSimpleError(ErrMethodCallNotImplemented, method)
}

// WrapMethodCallNotImplementedWithPos wraps errors for method calls not implemented with position
func WrapMethodCallNotImplementedWithPos(method string, pos interface{}) error {
	return WrapSimpleErrorWithPos(ErrMethodCallNotImplemented, method, pos)
}

// WrapFunctionRequiresNamedArgsWithPos wraps errors for functions requiring named arguments
func WrapFunctionRequiresNamedArgsWithPos(funcName string, paramCount int, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		//nolint:err113 // Dynamic error needed for exact test format matching
		return fmt.Errorf("line %d:%d: function requires named arguments '%s' has %d parameters "+
			"and requires named arguments. Use: %s(x: value, y: value)",
			position.Line, position.Column, funcName, paramCount, funcName)
	}
	//nolint:err113 // Dynamic error needed for exact test format matching
	return fmt.Errorf("function requires named arguments '%s' has %d parameters and requires named arguments",
		funcName, paramCount)
}

// WrapMatchNotExhaustiveWithPos wraps errors for non-exhaustive match expressions
func WrapMatchNotExhaustiveWithPos(missingPatterns []string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: missing patterns: %v",
			position.Line, position.Column, ErrMatchNotExhaustive, missingPatterns)
	}

	return fmt.Errorf("%w: missing patterns: %v", ErrMatchNotExhaustive, missingPatterns)
}

// WrapMatchTypeMismatchWithPos wraps errors for match expression type mismatches
func WrapMatchTypeMismatchWithPos(armIndex int, returnType, expectedType string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: arm %d returns '%s' but expected '%s'",
			position.Line, position.Column, ErrMatchTypeMismatch, armIndex, returnType, expectedType)
	}

	return fmt.Errorf("%w: arm %d returns '%s' but expected '%s'",
		ErrMatchTypeMismatch, armIndex, returnType, expectedType)
}

// WrapUnknownVariantWithPos wraps errors for unknown variants in match expressions
func WrapUnknownVariantWithPos(variantName, typeName string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: variant '%s' is not defined in type '%s'",
			position.Line, position.Column, ErrUnknownVariant, variantName, typeName)
	}

	return fmt.Errorf("%w: variant '%s' is not defined in type '%s'",
		ErrUnknownVariant, variantName, typeName)
}

// WrapTypeMismatchWithPos wraps type mismatch errors with position and detailed context
func WrapTypeMismatchWithPos(valueType, varName, annotatedType string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		//nolint:err113 // Dynamic error needed for exact test format matching
		return fmt.Errorf("line %d:%d: type mismatch: cannot assign %s to variable '%s' of type %s",
			position.Line, position.Column, valueType, varName, annotatedType)
	}
	//nolint:err113 // Dynamic error needed for exact test format matching
	return fmt.Errorf("type mismatch: cannot assign %s to variable '%s' of type %s",
		valueType, varName, annotatedType)
}

// WrapArithmeticTypeMismatch wraps arithmetic type mismatch errors
func WrapArithmeticTypeMismatch(operator, operand, expectedType, actualType string) error {
	return fmt.Errorf("%w: %s operand of %s must be %s: got %s",
		ErrArithmeticTypeMismatch, operand, operator, expectedType, actualType)
}
