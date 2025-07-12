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

	// Parse errors
	ErrParseTreeNil   = errors.New("parse tree is nil")
	ErrASTBuildFailed = errors.New("AST build failed")

	// Core function errors
	ErrPrintCannotConvert = errors.New("cannot convert value for printing")
	ErrUnsupportedCall    = errors.New("unsupported function call")
	ErrToStringReserved   = errors.New("toString is a reserved function name")

	// WebSocket errors
	ErrWebSocketKeepAliveWrongArgs = errors.New("webSocketKeepAlive function has wrong number of arguments")

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
	ErrFieldAccessNotImpl   = errors.New("field access not implemented")

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
)

// Static error definitions for match expressions
var (
	ErrMatchNotExhaustive = errors.New("match expression not exhaustive")
	ErrMatchTypeMismatch  = errors.New("match expression type mismatch")
	ErrUnknownVariant     = errors.New("unknown variant in match expression")
)

// Function argument count constants
const (
	PrintExpectedArgs                    = 1
	InputExpectedArgs                    = 0
	LengthExpectedArgs                   = 1
	ContainsExpectedArgs                 = 2
	SubstringExpectedArgs                = 3
	HTTPCreateServerExpectedArgs         = 2
	HTTPListenExpectedArgs               = 2
	HTTPStopServerExpectedArgs           = 1
	HTTPCreateClientExpectedArgs         = 2
	HTTPGetExpectedArgs                  = 3
	HTTPPostExpectedArgs                 = 4
	HTTPPutExpectedArgs                  = 4
	HTTPDeleteExpectedArgs               = 3
	HTTPRequestExpectedArgs              = 5
	HTTPCloseClientExpectedArgs          = 1
	WebSocketConnectExpectedArgs         = 1
	WebSocketSendExpectedArgs            = 2
	WebSocketCloseExpectedArgs           = 1
	WebSocketCreateServerExpectedArgs    = 2
	WebSocketServerListenExpectedArgs    = 2
	WebSocketServerBroadcastExpectedArgs = 2
	WebSocketStopServerExpectedArgs      = 1
	RangeExpectedArgs                    = 2
	ForEachExpectedArgs                  = 2
	MapExpectedArgs                      = 2
	FilterExpectedArgs                   = 2
	FoldExpectedArgs                     = 3
	SpawnProcessExpectedArgs             = 2
	SleepExpectedArgs                    = 1
	WriteFileExpectedArgs                = 2
	ReadFileExpectedArgs                 = 1
	AwaitProcessExpectedArgs             = 1
	CleanupProcessExpectedArgs           = 1
	ToStringExpectedArgs                 = 1
)

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
	return WrapWrongArgCount(ToStringFunc, ToStringExpectedArgs, argCount)
}

// WrapNoToStringImpl wraps errors when no toString implementation is found
func WrapNoToStringImpl(typeName string) error {
	return fmt.Errorf("%w for type %s", ErrNoToStringImpl, typeName)
}

// WrapFunctionArgsWithPos is a consolidated function for all "wrong arguments" errors
func WrapFunctionArgsWithPos(funcName string, expected int, actual int, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %s expects exactly %d argument(s), got %d: %w",
			position.Line, position.Column, funcName, expected, actual, ErrWrongArgCount)
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
	return WrapWrongArgCount("httpCreateServer", HTTPCreateServerExpectedArgs, argCount)
}

// WrapHTTPListenWrongArgs wraps errors for wrong number of HTTP listen arguments
func WrapHTTPListenWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpListen", HTTPListenExpectedArgs, argCount)
}

// WrapHTTPStopServerUnknownNamedArg wraps errors for unknown named arguments in HTTP stop server
func WrapHTTPStopServerUnknownNamedArg(argName string) error {
	return fmt.Errorf("%w: %s", ErrHTTPStopServerUnknownNamedArg, argName)
}

// WrapHTTPStopServerWrongArgCount wraps errors for wrong argument count in HTTP stop server
func WrapHTTPStopServerWrongArgCount(argCount int) error {
	return WrapWrongArgCount("httpStopServer", HTTPStopServerExpectedArgs, argCount)
}

// WrapHTTPPostWrongArgs wraps errors for wrong number of HTTP POST arguments
func WrapHTTPPostWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpPost", HTTPPostExpectedArgs, argCount)
}

// WrapHTTPPutWrongArgs wraps errors for wrong number of HTTP PUT arguments
func WrapHTTPPutWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpPut", HTTPPutExpectedArgs, argCount)
}

// WrapHTTPDeleteWrongArgs wraps errors for wrong number of HTTP DELETE arguments
func WrapHTTPDeleteWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpDelete", HTTPDeleteExpectedArgs, argCount)
}

// WrapHTTPRequestWrongArgs wraps errors for wrong number of HTTP request arguments
func WrapHTTPRequestWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpRequest", HTTPRequestExpectedArgs, argCount)
}

// WrapHTTPCloseClientWrongArgs wraps errors for wrong number of HTTP close client arguments
func WrapHTTPCloseClientWrongArgs(argCount int) error {
	return WrapWrongArgCount("httpCloseClient", HTTPCloseClientExpectedArgs, argCount)
}

// WrapWebSocketConnectWrongArgs wraps errors for wrong number of WebSocket connect arguments
func WrapWebSocketConnectWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketConnect", WebSocketConnectExpectedArgs, argCount)
}

// WrapWebSocketSendWrongArgs wraps errors for wrong number of WebSocket send arguments
func WrapWebSocketSendWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketSend", WebSocketSendExpectedArgs, argCount)
}

// WrapWebSocketCloseWrongArgs wraps errors for wrong number of WebSocket close arguments
func WrapWebSocketCloseWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketClose", WebSocketCloseExpectedArgs, argCount)
}

// WrapWebSocketCreateServerWrongArgs wraps errors for wrong number of WebSocket server creation arguments
func WrapWebSocketCreateServerWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketCreateServer", WebSocketCreateServerExpectedArgs, argCount)
}

// WrapWebSocketServerListenWrongArgs wraps errors for wrong number of WebSocket server listen arguments
func WrapWebSocketServerListenWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketServerListen", WebSocketServerListenExpectedArgs, argCount)
}

// WrapWebSocketServerBroadcastWrongArgs wraps errors for wrong number of WebSocket server broadcast arguments
func WrapWebSocketServerBroadcastWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketServerBroadcast", WebSocketServerBroadcastExpectedArgs, argCount)
}

// WrapWebSocketStopServerWrongArgs wraps errors for wrong number of WebSocket stop server arguments
func WrapWebSocketStopServerWrongArgs(argCount int) error {
	return WrapWrongArgCount("websocketStopServer", WebSocketStopServerExpectedArgs, argCount)
}

// WrapMapWrongArgs wraps errors for wrong number of map arguments
func WrapMapWrongArgs(argCount int) error {
	return WrapWrongArgCount("map", MapExpectedArgs, argCount)
}

// WrapFilterWrongArgs wraps errors for wrong number of filter arguments
func WrapFilterWrongArgs(argCount int) error {
	return WrapWrongArgCount("filter", FilterExpectedArgs, argCount)
}

// WrapFoldWrongArgs wraps errors for wrong number of fold arguments
func WrapFoldWrongArgs(argCount int) error {
	return WrapWrongArgCount("fold", FoldExpectedArgs, argCount)
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
	return WrapWrongArgCount("sleep", SleepExpectedArgs, argCount)
}

// Functions removed - definitions are above to match exact expected outputs

// WrapAwaitProcessWrongArgs wraps errors for wrong number of await process arguments
func WrapAwaitProcessWrongArgs(argCount int) error {
	return WrapWrongArgCount("awaitProcess", AwaitProcessExpectedArgs, argCount)
}

// WrapCleanupProcessWrongArgs wraps errors for wrong number of cleanup process arguments
func WrapCleanupProcessWrongArgs(argCount int) error {
	return WrapWrongArgCount("cleanupProcess", CleanupProcessExpectedArgs, argCount)
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

// WrapFieldAccessNotImpl wraps field access not implemented errors
func WrapFieldAccessNotImpl(fieldName string) error {
	return fmt.Errorf("%w: %s", ErrFieldAccessNotImpl, fieldName)
}

// WrapFieldAccessNotImplWithPos wraps field access errors with position information
func WrapFieldAccessNotImplWithPos(fieldName string, pos interface{}) error {
	if position, ok := pos.(*ast.Position); ok && position != nil {
		return fmt.Errorf("line %d:%d: %w: %s",
			position.Line, position.Column, ErrFieldAccessNotImpl, fieldName)
	}
	return fmt.Errorf("%w: %s", ErrFieldAccessNotImpl, fieldName)
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
	return WrapWrongArgCount("print", PrintExpectedArgs, argCount)
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
		return fmt.Errorf("line %d:%d: %w: %s has %d parameters and requires named arguments. Use: %s(x: value, y: value)",
			position.Line, position.Column, ErrFunctionRequiresNamed, funcName, paramCount, funcName)
	}
	return fmt.Errorf("%w: %s has %d parameters and requires named arguments",
		ErrFunctionRequiresNamed, funcName, paramCount)
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
