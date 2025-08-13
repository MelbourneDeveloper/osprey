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
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s",
			position.Line, position.Column, ErrUndefinedVariable, varName)
	}
	return fmt.Errorf("%w: %s", ErrUndefinedVariable, varName)
}

// WrapUnsupportedBinaryOpWithPos wraps errors for unsupported binary operators
func WrapUnsupportedBinaryOpWithPos(op string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s", position.Line, position.Column, ErrUnsupportedBinaryOp, op)
	}
	return fmt.Errorf("%w: %s", ErrUnsupportedBinaryOp, op)
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
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s", position.Line, position.Column, ErrUnsupportedUnaryOperator, op)
	}
	return fmt.Errorf("%w: %s", ErrUnsupportedUnaryOperator, op)
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

// WrapMethodNotImpl wraps errors for method not implemented
func WrapMethodNotImpl(method string) error {
	return WrapMethodCallNotImplemented(method)
}

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
	return fmt.Errorf("%w: %s", ErrUndefinedType, typeName)
}

// WrapUndefinedFunction wraps errors for undefined functions
func WrapUndefinedFunction(funcName string) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotFound, funcName)
}

// WrapMissingField wraps errors for missing fields
func WrapMissingField(field string) error {
	return fmt.Errorf("%w: %s", ErrMissingField, field)
}

// WrapBuiltInRedefine wraps errors for redefining built-in functions
func WrapBuiltInRedefine(funcName string) error {
	return fmt.Errorf("%w: %s", ErrBuiltInRedefine, funcName)
}

// WrapHTTPCreateServerWrongArgs wraps errors for wrong number of HTTP server creation arguments
func WrapHTTPCreateServerWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPCreateServerOsprey)
	return WrapWrongArgCount(HTTPCreateServerOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPListenWrongArgs wraps errors for wrong number of HTTP listen arguments
func WrapHTTPListenWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPListenOsprey)
	return WrapWrongArgCount(HTTPListenOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPStopServerUnknownNamedArg wraps errors for unknown named arguments in HTTP stop server
func WrapHTTPStopServerUnknownNamedArg(argName string) error {
	return fmt.Errorf("%w: %s", ErrHTTPStopServerUnknownNamedArg, argName)
}

// WrapHTTPStopServerWrongArgCount wraps errors for wrong argument count in HTTP stop server
func WrapHTTPStopServerWrongArgCount(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPStopServerOsprey)
	return WrapWrongArgCount(HTTPStopServerOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPPostWrongArgs wraps errors for wrong number of HTTP POST arguments
func WrapHTTPPostWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPPostOsprey)
	return WrapWrongArgCount(HTTPPostOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPPutWrongArgs wraps errors for wrong number of HTTP PUT arguments
func WrapHTTPPutWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPPutOsprey)
	return WrapWrongArgCount(HTTPPutOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPDeleteWrongArgs wraps errors for wrong number of HTTP DELETE arguments
func WrapHTTPDeleteWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPDeleteOsprey)
	return WrapWrongArgCount(HTTPDeleteOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPRequestWrongArgs wraps errors for wrong number of HTTP request arguments
func WrapHTTPRequestWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPRequestOsprey)
	return WrapWrongArgCount(HTTPRequestOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPCloseClientWrongArgs wraps errors for wrong number of HTTP close client arguments
func WrapHTTPCloseClientWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(HTTPCloseClientOsprey)
	return WrapWrongArgCount(HTTPCloseClientOsprey, len(fn.ParameterTypes), argCount)
}

// WrapHTTPFunctionWrongArgs wraps errors for wrong number of HTTP function arguments
func WrapHTTPFunctionWrongArgs(functionName string, expected, actual int) error {
	return WrapWrongArgCount(functionName, expected, actual)
}

// WrapHTTPFunctionNotFound wraps errors for HTTP function not found
func WrapHTTPFunctionNotFound(functionName string) error {
	return WrapUndefinedFunction(functionName)
}

// WrapHTTPFunctionMissingNamedArg wraps errors for missing named arguments in HTTP functions
func WrapHTTPFunctionMissingNamedArg(functionName, argName string) error {
	return WrapMissingArgument(argName, functionName)
}

// WrapWebSocketConnectWrongArgs wraps errors for wrong number of WebSocket connect arguments
func WrapWebSocketConnectWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketConnectOsprey)
	return WrapWrongArgCount(WebSocketConnectOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketSendWrongArgs wraps errors for wrong number of WebSocket send arguments
func WrapWebSocketSendWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketSendOsprey)
	return WrapWrongArgCount(WebSocketSendOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketCloseWrongArgs wraps errors for wrong number of WebSocket close arguments
func WrapWebSocketCloseWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketCloseOsprey)
	return WrapWrongArgCount(WebSocketCloseOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketCreateServerWrongArgs wraps errors for wrong number of WebSocket server creation arguments
func WrapWebSocketCreateServerWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketCreateServerOsprey)
	return WrapWrongArgCount(WebSocketCreateServerOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketServerListenWrongArgs wraps errors for wrong number of WebSocket server listen arguments
func WrapWebSocketServerListenWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketServerListenOsprey)
	return WrapWrongArgCount(WebSocketServerListenOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketServerBroadcastWrongArgs wraps errors for wrong number of WebSocket server broadcast arguments
func WrapWebSocketServerBroadcastWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketServerBroadcastOsprey)
	return WrapWrongArgCount(WebSocketServerBroadcastOsprey, len(fn.ParameterTypes), argCount)
}

// WrapWebSocketStopServerWrongArgs wraps errors for wrong number of WebSocket stop server arguments
func WrapWebSocketStopServerWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(WebSocketStopServerOsprey)
	return WrapWrongArgCount(WebSocketStopServerOsprey, len(fn.ParameterTypes), argCount)
}

// WrapMapWrongArgs wraps errors for wrong number of map arguments
func WrapMapWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(MapFunc)
	return WrapWrongArgCount("map", len(fn.ParameterTypes), argCount)
}

// WrapFilterWrongArgs wraps errors for wrong number of filter arguments
func WrapFilterWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(FilterFunc)
	return WrapWrongArgCount("filter", len(fn.ParameterTypes), argCount)
}

// WrapFoldWrongArgs wraps errors for wrong number of fold arguments
func WrapFoldWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(FoldFunc)
	return WrapWrongArgCount("fold", len(fn.ParameterTypes), argCount)
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

// WrapFunctionNotFound wraps errors for function not found
func WrapFunctionNotFound(funcName string) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotFound, funcName)
}

// WrapBuiltInTwoArgs wraps errors for built-in functions expecting two arguments
func WrapBuiltInTwoArgs(funcName string) error {
	return fmt.Errorf("%w: %s", ErrBuiltInTwoArgs, funcName)
}

// WrapSleepWrongArgs wraps errors for wrong number of sleep arguments
func WrapSleepWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(SleepFunc)
	return WrapWrongArgCount("sleep", len(fn.ParameterTypes), argCount)
}

// Functions removed - definitions are above to match exact expected outputs

// WrapAwaitProcessWrongArgs wraps errors for wrong number of await process arguments
func WrapAwaitProcessWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(AwaitProcessFunc)
	return WrapWrongArgCount("awaitProcess", len(fn.ParameterTypes), argCount)
}

// WrapCleanupProcessWrongArgs wraps errors for wrong number of cleanup process arguments
func WrapCleanupProcessWrongArgs(argCount int) error {
	fn, _ := GlobalBuiltInRegistry.GetFunction(CleanupProcessFunc)
	return WrapWrongArgCount("cleanupProcess", len(fn.ParameterTypes), argCount)
}

// Additional wrapper functions expected by tests

// WrapFunctionNotDeclared wraps function not declared errors
func WrapFunctionNotDeclared(funcName string) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotDeclared, funcName)
}

// WrapUndefinedVariable wraps undefined variable errors
func WrapUndefinedVariable(varName string) error {
	return fmt.Errorf("%w: %s", ErrUndefinedVariable, varName)
}

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

// WrapMethodCallNotImplemented wraps errors for method calls not implemented
func WrapMethodCallNotImplemented(method string) error {
	return fmt.Errorf("%w: %s", ErrMethodCallNotImplemented, method)
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
