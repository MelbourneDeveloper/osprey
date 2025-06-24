package codegen

import (
	"errors"
	"fmt"
	"strings"
)

// Static error definitions to replace dynamic errors.
var (
	ErrToStringReserved      = errors.New("cannot redefine built-in function 'toString'. Built-in functions are reserved")
	ErrUnsupportedStatement  = errors.New("unsupported statement type")
	ErrFunctionNotDeclared   = errors.New("function not declared")
	ErrUndefinedVariable     = errors.New("undefined variable")
	ErrUnsupportedExpression = errors.New("unsupported expression type")
	ErrUnsupportedBinaryOp   = errors.New("unsupported binary operator")
	ErrUnsupportedUnaryOp    = errors.New("unsupported unary operator")
	ErrFieldAccessNotImpl    = errors.New("field access not implemented for field")
	ErrToStringWrongArgs     = errors.New("toString expects exactly 1 argument")
	ErrPrintWrongArgs        = errors.New("print expects exactly 1 argument")
	ErrInputWrongArgs        = errors.New("input expects exactly 0 arguments")
	ErrUnsupportedCall       = errors.New("unsupported call expression")
	ErrMethodNotImpl         = errors.New("method call not implemented")
	ErrNoToStringImpl        = errors.New("no toString implementation for type")
	ErrNoToStringForFunc     = errors.New("no toString implementation for function return type")
	ErrPrintCannotConvert    = errors.New("print() cannot convert to string")
	ErrPrintComplexExpr      = errors.New("print() cannot handle complex expression. Use toString() to convert")
	ErrPrintUnknownFunc      = errors.New("print() cannot determine return type of function")
	ErrFunctionRequiresNamed = errors.New("function requires named arguments")
	ErrWrongArgCount         = errors.New("function expects different number of arguments")
	ErrMissingArgument       = errors.New("missing argument for parameter")
	ErrParseErrors           = errors.New("parse errors")
	ErrParseTreeNil          = errors.New("failed to parse source code: parse tree is nil")
	ErrASTBuildFailed        = errors.New("failed to build AST: program is nil")
	ErrLLVMGenFailed         = errors.New("failed to generate LLVM IR")
	ErrWriteIRFile           = errors.New("failed to write IR file")
	ErrCompileToObj          = errors.New("failed to compile IR to object file")
	ErrLinkExecutable        = errors.New("failed to link executable")
	ErrToolNotFound          = errors.New("tool not found in PATH or common locations")
	ErrNoSuitableCompiler    = errors.New("no suitable compiler found")
	ErrPrintComplexCall      = errors.New("print() cannot handle complex call expression. Use toString() to convert")
	ErrPrintConvertError     = errors.New("print() cannot convert function return type to string")
	ErrPrintDetermineError   = errors.New("print() cannot determine return type of function")
	ErrRangeWrongArgs        = errors.New("range expects exactly 2 arguments (start, end)")
	ErrForEachWrongArgs      = errors.New("forEach expects exactly 2 arguments (iterator, function)")
	ErrForEachNotFunction    = errors.New("forEach second argument must be a function")
	ErrMapWrongArgs          = errors.New("map expects exactly 2 arguments (iterator, function)")
	ErrMapNotFunction        = errors.New("map second argument must be a function")
	ErrFilterWrongArgs       = errors.New("filter expects exactly 2 arguments (iterator, predicate)")
	ErrFilterNotFunction     = errors.New("filter second argument must be a predicate function")
	ErrFoldWrongArgs         = errors.New("fold expects exactly 3 arguments (iterator, initial, function)")
	ErrFoldNotFunction       = errors.New("fold third argument must be a function")
	ErrInputNoArgs           = errors.New("input function takes no arguments")
	ErrBuiltInTwoArgs        = errors.New("built-in function does not accept two arguments")
	ErrBuiltInRedefine       = errors.New("cannot redefine built-in function. Built-in functions are reserved")
	ErrFunctionNotFound      = errors.New("function not found")
	ErrImplicitAnyReturn     = errors.New("function cannot implicitly return 'any' type")

	// Match expression type safety errors.
	ErrMatchTypeMismatch            = errors.New("match expression type mismatch")
	ErrMatchNotExhaustive           = errors.New("match expression not exhaustive")
	ErrMatchInvalidPattern          = errors.New("invalid pattern in match expression")
	ErrMatchUnknownVariant          = errors.New("unknown variant in match expression")
	ErrMatchMixedPatterns           = errors.New("match arms contain mixed pattern types")
	ErrMatchUnreachableArm          = errors.New("unreachable match arm")
	ErrMatchWildcardNotLast         = errors.New("wildcard pattern must be the last arm")
	ErrMatchDuplicateArm            = errors.New("duplicate match arm")
	ErrMatchExpressionNotExhaustive = errors.New("match expression is not exhaustive")

	// Type constructor errors.
	ErrUndefinedType       = errors.New("undefined type")
	ErrMissingField        = errors.New("missing field in type constructor")
	ErrConstraintViolation = errors.New("constraint violation")

	// Result type field access errors.
	ErrFieldAccessOnResult         = errors.New("cannot access field on Result type - pattern matching required")
	ErrConstraintResultFieldAccess = errors.New("field access requires pattern matching on Result type")

	// Any type access violation errors.
	ErrAnyDirectArithmetic = errors.New(
		"cannot use 'any' type directly in arithmetic operation - pattern matching required")
	ErrAnyDirectFieldAccess = errors.New("cannot access field on 'any' type without pattern matching")
	ErrAnyDirectAssignment  = errors.New("cannot assign 'any' to typed variable without pattern matching")
	ErrAnyDirectFunctionArg = errors.New(
		"cannot pass 'any' type to function expecting specific type - pattern matching required")
	ErrAnyImplicitConversion = errors.New(
		"cannot implicitly convert 'any' to specific type - use pattern matching to extract specific type")
	ErrAnyDirectVariableAccess = errors.New(
		"cannot access variable of type 'any' directly - pattern matching required")
	ErrAnyPatternNotExhaustive = errors.New(
		"pattern matching on 'any' type must handle all possible types or include wildcard")
	ErrAnyPatternImpossible  = errors.New("pattern type is not a possible type for expression")
	ErrAnyPatternUnreachable = errors.New("unreachable pattern: type cannot occur based on context analysis")

	// HTTP Server function errors.
	ErrHTTPCreateServerWrongArgs = errors.New("httpCreateServer expects exactly 2 arguments (port, address)")
	ErrHTTPListenWrongArgs       = errors.New("httpListen expects exactly 2 arguments (server_id, handler)")
	ErrHTTPStopServerWrongArgs   = errors.New("httpStopServer expects exactly 1 argument (server_id)")

	// HTTP Client function errors.
	ErrHTTPCreateClientWrongArgs = errors.New("httpCreateClient expects exactly 2 arguments (base_url, timeout)")
	ErrHTTPGetWrongArgs          = errors.New("httpGet expects exactly 3 arguments (client_id, path, headers)")
	ErrHTTPPostWrongArgs         = errors.New("httpPost expects exactly 4 arguments (client_id, path, body, headers)")
	ErrHTTPPutWrongArgs          = errors.New("httpPut expects exactly 4 arguments (client_id, path, body, headers)")
	ErrHTTPDeleteWrongArgs       = errors.New("httpDelete expects exactly 3 arguments (client_id, path, headers)")
	ErrHTTPRequestWrongArgs      = errors.New(
		"httpRequest expects exactly 5 arguments (client_id, method, path, headers, body)")
	ErrHTTPCloseClientWrongArgs = errors.New("httpCloseClient expects exactly 1 argument (client_id)")

	// WebSocket function errors.
	ErrWebSocketConnectWrongArgs = errors.New("websocketConnect expects exactly 2 arguments (url, message_handler)")
	ErrWebSocketSendWrongArgs    = errors.New("websocketSend expects exactly 2 arguments (ws_id, message)")
	ErrWebSocketCloseWrongArgs   = errors.New("websocketClose expects exactly 1 argument (ws_id)")

	// WebSocket Server function errors.
	ErrWebSocketCreateServerWrongArgs = errors.New(
		"websocketCreateServer expects exactly 3 arguments (port, address, path)")
	ErrWebSocketServerListenWrongArgs = errors.New(
		"websocketServerListen expects exactly 1 argument (server_id)")
	ErrWebSocketServerBroadcastWrongArgs = errors.New(
		"websocketServerBroadcast expects exactly 2 arguments (server_id, message)")
	ErrWebSocketStopServerWrongArgs = errors.New(
		"websocketStopServer expects exactly 1 argument (server_id)")
	ErrWebSocketKeepAliveWrongArgs = errors.New("websocketKeepAlive requires no arguments")

	// HTTP Server Named Arguments errors.
	ErrHTTPStopServerUnknownNamedArg = errors.New("httpStopServer: unknown named argument")
	ErrHTTPStopServerWrongArgCount   = errors.New("httpStopServer expects exactly 1 argument (serverID)")

	// Match expression validation errors.
	ErrMatchUnknownVariantType = errors.New("unknown variant in match expression")

	// String utility function errors
	ErrLengthWrongArgs    = errors.New("length expects exactly 1 argument")
	ErrContainsWrongArgs  = errors.New("contains expects exactly 2 arguments")
	ErrSubstringWrongArgs = errors.New("substring expects exactly 3 arguments")

	// Process spawning errors.
	ErrSpawnProcessWrongArgs = errors.New("spawnProcess expects exactly 2 arguments (command, callback)")

	ErrAwaitProcessWrongArgs   = errors.New("awaitProcess expects exactly 1 argument (processId)")
	ErrCleanupProcessWrongArgs = errors.New("cleanupProcess expects exactly 1 argument (processId)")
	ErrSleepWrongArgs          = errors.New("sleep expects exactly 1 argument (milliseconds)")
	ErrWriteFileWrongArgs      = errors.New("writeFile expects exactly 2 arguments (filename, content)")
	ErrReadFileWrongArgs       = errors.New("readFile expects exactly 1 argument (filename)")
	ErrParseJSONWrongArgs      = errors.New("parseJSON expects exactly 1 argument (json_string)")
	ErrExtractCodeWrongArgs    = errors.New("extractCode expects exactly 1 argument (json_string)")
)

// Helper functions to wrap static errors with context

// WrapUnsupportedStatement wraps unsupported statement errors with type information.
func WrapUnsupportedStatement(t interface{}) error {
	return fmt.Errorf("%w: %T", ErrUnsupportedStatement, t)
}

// WrapFunctionNotDeclared wraps function not declared errors with function name.
func WrapFunctionNotDeclared(name string) error {
	return fmt.Errorf("function '%s' not declared: %w", name, ErrFunctionNotDeclared)
}

// WrapUndefinedVariable wraps undefined variable errors with variable name.
func WrapUndefinedVariable(name string) error {
	return fmt.Errorf("undefined variable '%s': %w", name, ErrUndefinedVariable)
}

// WrapUnsupportedExpression wraps unsupported expression errors with type information.
func WrapUnsupportedExpression(t interface{}) error {
	return fmt.Errorf("%w: %T", ErrUnsupportedExpression, t)
}

// WrapUnsupportedBinaryOp wraps unsupported binary operator errors with operator.
func WrapUnsupportedBinaryOp(op string) error {
	return fmt.Errorf("unsupported binary operator '%s': %w", op, ErrUnsupportedBinaryOp)
}

// WrapUnsupportedUnaryOp wraps unsupported unary operator errors with operator.
func WrapUnsupportedUnaryOp(op string) error {
	return fmt.Errorf("unsupported unary operator '%s': %w", op, ErrUnsupportedUnaryOp)
}

// WrapFieldAccessNotImpl wraps field access not implemented errors with field name.
func WrapFieldAccessNotImpl(field string) error {
	return fmt.Errorf("field access not implemented for field '%s': %w", field, ErrFieldAccessNotImpl)
}

// WrapToStringWrongArgs wraps toString wrong arguments errors with argument count.
func WrapToStringWrongArgs(got int) error {
	return fmt.Errorf("toString expects exactly 1 argument, got %d: %w", got, ErrToStringWrongArgs)
}

// WrapPrintWrongArgs wraps print wrong arguments errors with argument count.
func WrapPrintWrongArgs(got int) error {
	return fmt.Errorf("print expects exactly 1 argument, got %d: %w", got, ErrPrintWrongArgs)
}

// WrapInputWrongArgs wraps input wrong arguments errors with argument count.
func WrapInputWrongArgs(got int) error {
	return fmt.Errorf("input expects exactly 0 arguments, got %d: %w", got, ErrInputWrongArgs)
}

// WrapMethodNotImpl wraps method not implemented errors with method name.
func WrapMethodNotImpl(method string) error {
	return fmt.Errorf("method call not implemented for method '%s': %w", method, ErrMethodNotImpl)
}

// WrapNoToStringImpl wraps no toString implementation errors with type name.
func WrapNoToStringImpl(typeName string) error {
	return fmt.Errorf("%w: %s", ErrNoToStringImpl, typeName)
}

// WrapNoToStringForFunc wraps no toString for function errors with type name.
func WrapNoToStringForFunc(typeName string) error {
	return fmt.Errorf("%w: %s", ErrNoToStringForFunc, typeName)
}

// WrapPrintCannotConvert wraps print cannot convert errors with variable and type info.
func WrapPrintCannotConvert(varName, typeName string) error {
	return fmt.Errorf("%w variable '%s' of type '%s'. Use toString(%s) explicitly",
		ErrPrintCannotConvert, varName, typeName, varName)
}

// WrapPrintUnknownFunc wraps print unknown function errors with function name.
func WrapPrintUnknownFunc(funcName string) error {
	return fmt.Errorf("%w '%s'. Use toString(%s) to convert", ErrPrintUnknownFunc, funcName, funcName)
}

// WrapPrintCannotConvertFunc wraps print cannot convert function errors with type and function info.
func WrapPrintCannotConvertFunc(returnType, funcName string) error {
	return fmt.Errorf("%w function return type '%s'. Use toString(%s) explicitly",
		ErrPrintCannotConvert, returnType, funcName)
}

// WrapPrintComplexExpr wraps print complex expression errors with expression type.
func WrapPrintComplexExpr(expr interface{}) error {
	return fmt.Errorf("%w of type %T", ErrPrintComplexExpr, expr)
}

// WrapFunctionRequiresNamed wraps function requires named arguments errors with details.
func WrapFunctionRequiresNamed(funcName string, paramCount int, example string) error {
	return fmt.Errorf("%w '%s' has %d parameters and requires named arguments. Use: %s(%s)",
		ErrFunctionRequiresNamed, funcName, paramCount, funcName, example)
}

// WrapWrongArgCount wraps wrong argument count errors with function and count details.
func WrapWrongArgCount(funcName string, expected, got int) error {
	return fmt.Errorf("%w %s expects %d arguments, got %d", ErrWrongArgCount, funcName, expected, got)
}

// WrapMissingArgument wraps missing argument errors with parameter and function info.
func WrapMissingArgument(paramName, funcName string) error {
	return fmt.Errorf("%w %s in function %s", ErrMissingArgument, paramName, funcName)
}

// WrapParseErrors wraps parse errors with error details.
func WrapParseErrors(errors string) error {
	return fmt.Errorf("%w:\n%s", ErrParseErrors, errors)
}

// WrapLLVMGenFailed wraps LLVM generation failed errors with underlying error.
func WrapLLVMGenFailed(err error) error {
	return fmt.Errorf("%w: %w", ErrLLVMGenFailed, err)
}

// WrapWriteIRFile wraps write IR file errors with underlying error.
func WrapWriteIRFile(err error) error {
	return fmt.Errorf("%w: %w", ErrWriteIRFile, err)
}

// WrapCompileToObj wraps compile to object errors with underlying error and output.
func WrapCompileToObj(err error, output string) error {
	return fmt.Errorf("%w: %w\nllc output: %s", ErrCompileToObj, err, output)
}

// WrapLinkExecutable wraps link executable errors with compiler, error and output info.
func WrapLinkExecutable(compiler string, err error, output string) error {
	return fmt.Errorf("%w with %s: %w\nOutput: %s", ErrLinkExecutable, compiler, err, output)
}

// WrapToolNotFound wraps tool not found errors with tool name.
func WrapToolNotFound(toolName string) error {
	return fmt.Errorf("%w: %s", ErrToolNotFound, toolName)
}

// WrapNoSuitableCompiler wraps no suitable compiler errors with tried compilers list.
func WrapNoSuitableCompiler(compilers []string) error {
	return fmt.Errorf("%w (tried: %v)", ErrNoSuitableCompiler, compilers)
}

// WrapPrintConvertError wraps print convert errors with return type and function name.
func WrapPrintConvertError(returnType, funcName string) error {
	return fmt.Errorf("%w: function return type '%s'. Use toString(%s) explicitly",
		ErrPrintConvertError, returnType, funcName)
}

// WrapPrintDetermineError wraps print determine errors with function name.
func WrapPrintDetermineError(funcName string) error {
	return fmt.Errorf("%w: '%s'. Use toString(%s) to convert", ErrPrintDetermineError, funcName, funcName)
}

// WrapBuiltInRedefine wraps built-in function redefinition errors.
func WrapBuiltInRedefine(fnName string) error {
	return fmt.Errorf("%w: '%s'", ErrBuiltInRedefine, fnName)
}

// WrapRangeWrongArgs wraps range wrong arguments error.
func WrapRangeWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrRangeWrongArgs, got)
}

// WrapForEachWrongArgs wraps forEach wrong arguments error.
func WrapForEachWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrForEachWrongArgs, got)
}

// WrapMapWrongArgs wraps map wrong arguments error.
func WrapMapWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrMapWrongArgs, got)
}

// WrapFilterWrongArgs wraps filter wrong arguments error.
func WrapFilterWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrFilterWrongArgs, got)
}

// WrapFoldWrongArgs wraps fold wrong arguments error.
func WrapFoldWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrFoldWrongArgs, got)
}

// WrapBuiltInTwoArgs wraps built-in two arguments error.
func WrapBuiltInTwoArgs(fnName string) error {
	return fmt.Errorf("%w: %s", ErrBuiltInTwoArgs, fnName)
}

// WrapFunctionNotFound wraps function not found errors with function name.
func WrapFunctionNotFound(fnName string) error {
	return fmt.Errorf("%w: %s", ErrFunctionNotFound, fnName)
}

// WrapImplicitAnyReturn wraps implicit any return errors with function name.
func WrapImplicitAnyReturn(fnName string) error {
	return fmt.Errorf("%w: function '%s' - if 'any' return type is intended, declare it explicitly with '-> any'",
		ErrImplicitAnyReturn, fnName)
}

// Match expression error wrappers

// WrapMatchTypeMismatch wraps match type mismatch errors with detailed type information.
func WrapMatchTypeMismatch(exprType, patternType string) error {
	return fmt.Errorf("%w: cannot match expression of type '%s' against pattern of type '%s'",
		ErrMatchTypeMismatch, exprType, patternType)
}

// WrapMatchNotExhaustive wraps non-exhaustive match errors with missing patterns.
func WrapMatchNotExhaustive(missingPatterns []string) error {
	return fmt.Errorf("%w: missing patterns: %v", ErrMatchNotExhaustive, missingPatterns)
}

// WrapMatchInvalidPattern wraps invalid pattern errors with pattern details.
func WrapMatchInvalidPattern(pattern, reason string) error {
	return fmt.Errorf("%w: pattern '%s' is invalid: %s", ErrMatchInvalidPattern, pattern, reason)
}

// WrapMatchUnknownVariant wraps unknown variant errors with variant and type info.
func WrapMatchUnknownVariant(variant, typeName string) error {
	return fmt.Errorf("%w: variant '%s' is not defined in type '%s'", ErrMatchUnknownVariant, variant, typeName)
}

// WrapMatchMixedPatterns wraps mixed pattern type errors with pattern details.
func WrapMatchMixedPatterns(patterns []string) error {
	return fmt.Errorf("%w: found patterns of different types: %v", ErrMatchMixedPatterns, patterns)
}

// WrapMatchUnreachableArm wraps unreachable arm errors with arm details.
func WrapMatchUnreachableArm(armPattern string) error {
	return fmt.Errorf("%w: pattern '%s' will never be matched", ErrMatchUnreachableArm, armPattern)
}

// WrapMatchWildcardNotLast wraps wildcard position errors.
func WrapMatchWildcardNotLast() error {
	return fmt.Errorf("%w: wildcard '_' pattern can only appear as the final match arm", ErrMatchWildcardNotLast)
}

// WrapMatchDuplicateArm wraps duplicate arm errors with pattern info.
func WrapMatchDuplicateArm(pattern string) error {
	return fmt.Errorf("%w: pattern '%s' appears multiple times", ErrMatchDuplicateArm, pattern)
}

// Type constructor error wrappers

// WrapUndefinedType wraps undefined type errors with type name.
func WrapUndefinedType(typeName string) error {
	return fmt.Errorf("%w: '%s'", ErrUndefinedType, typeName)
}

// WrapMissingField wraps missing field errors with field name.
func WrapMissingField(fieldName string) error {
	return fmt.Errorf("%w: '%s'", ErrMissingField, fieldName)
}

// WrapConstraintViolation wraps constraint violation errors with field and constraint info.
func WrapConstraintViolation(fieldName, constraintName string) error {
	return fmt.Errorf("%w: field '%s' failed constraint '%s'", ErrConstraintViolation, fieldName, constraintName)
}

// WrapFieldAccessOnResult wraps field access on Result type errors with field and type info.
func WrapFieldAccessOnResult(fieldName, typeName string) error {
	return fmt.Errorf("cannot access field '%s' on Result<%s, ConstructionError> type - pattern matching required: %w",
		fieldName, typeName, ErrFieldAccessOnResult)
}

// WrapConstraintResultFieldAccess wraps constraint result field access errors.
func WrapConstraintResultFieldAccess() error {
	return ErrConstraintResultFieldAccess
}

// WrapHTTPCreateServerWrongArgs wraps httpCreateServer wrong arguments errors with argument count.
func WrapHTTPCreateServerWrongArgs(got int) error {
	return fmt.Errorf("httpCreateServer expects exactly 2 arguments (port, address), got %d: %w",
		got, ErrHTTPCreateServerWrongArgs)
}

// WrapHTTPListenWrongArgs wraps httpListen wrong arguments errors with argument count.
func WrapHTTPListenWrongArgs(got int) error {
	return fmt.Errorf("httpListen expects exactly 2 arguments (server_id, handler), got %d: %w",
		got, ErrHTTPListenWrongArgs)
}

// WrapHTTPStopServerWrongArgs wraps httpStopServer wrong arguments errors with argument count.
func WrapHTTPStopServerWrongArgs(got int) error {
	return fmt.Errorf("httpStopServer expects exactly 1 argument (server_id), got %d: %w",
		got, ErrHTTPStopServerWrongArgs)
}

// WrapHTTPCreateClientWrongArgs wraps httpCreateClient wrong arguments errors with argument count.
func WrapHTTPCreateClientWrongArgs(got int) error {
	return fmt.Errorf("httpCreateClient expects exactly 2 arguments (base_url, timeout), got %d: %w",
		got, ErrHTTPCreateClientWrongArgs)
}

// WrapHTTPGetWrongArgs wraps httpGet wrong arguments errors with argument count.
func WrapHTTPGetWrongArgs(got int) error {
	return fmt.Errorf("httpGet expects exactly 3 arguments (client_id, path, headers), got %d: %w",
		got, ErrHTTPGetWrongArgs)
}

// WrapHTTPPostWrongArgs wraps httpPost wrong arguments errors with argument count.
func WrapHTTPPostWrongArgs(got int) error {
	return fmt.Errorf("httpPost expects exactly 4 arguments (client_id, path, body, headers), got %d: %w",
		got, ErrHTTPPostWrongArgs)
}

// WrapHTTPPutWrongArgs wraps httpPut wrong arguments errors with argument count.
func WrapHTTPPutWrongArgs(got int) error {
	return fmt.Errorf("httpPut expects exactly 4 arguments (client_id, path, body, headers), got %d: %w",
		got, ErrHTTPPutWrongArgs)
}

// WrapHTTPDeleteWrongArgs wraps httpDelete wrong arguments errors with argument count.
func WrapHTTPDeleteWrongArgs(got int) error {
	return fmt.Errorf("httpDelete expects exactly 3 arguments (client_id, path, headers), got %d: %w",
		got, ErrHTTPDeleteWrongArgs)
}

// WrapHTTPRequestWrongArgs wraps httpRequest wrong arguments errors with argument count.
func WrapHTTPRequestWrongArgs(got int) error {
	return fmt.Errorf("httpRequest expects exactly 5 arguments (client_id, method, path, headers, body), got %d: %w",
		got, ErrHTTPRequestWrongArgs)
}

// WrapHTTPCloseClientWrongArgs wraps httpCloseClient wrong arguments errors with argument count.
func WrapHTTPCloseClientWrongArgs(got int) error {
	return fmt.Errorf("httpCloseClient expects exactly 1 argument (client_id), got %d: %w",
		got, ErrHTTPCloseClientWrongArgs)
}

// WrapWebSocketConnectWrongArgs wraps websocketConnect wrong arguments errors with argument count.
func WrapWebSocketConnectWrongArgs(got int) error {
	return fmt.Errorf("websocketConnect expects exactly 2 arguments (url, message_handler), got %d: %w",
		got, ErrWebSocketConnectWrongArgs)
}

// WrapWebSocketSendWrongArgs wraps websocketSend wrong arguments errors with argument count.
func WrapWebSocketSendWrongArgs(got int) error {
	return fmt.Errorf("websocketSend expects exactly 2 arguments (ws_id, message), got %d: %w",
		got, ErrWebSocketSendWrongArgs)
}

// WrapWebSocketCloseWrongArgs wraps websocketClose wrong arguments errors with argument count.
func WrapWebSocketCloseWrongArgs(got int) error {
	return fmt.Errorf("websocketClose expects exactly 1 argument, got %d: %w", got, ErrWebSocketCloseWrongArgs)
}

// WrapAnyDirectArithmetic wraps any direct arithmetic operation errors.
func WrapAnyDirectArithmetic(operation string) error {
	return fmt.Errorf("%w: operation '%s'", ErrAnyDirectArithmetic, operation)
}

// WrapAnyDirectFieldAccess wraps any direct field access errors.
func WrapAnyDirectFieldAccess(fieldName string) error {
	return fmt.Errorf("%w: field '%s'", ErrAnyDirectFieldAccess, fieldName)
}

// WrapAnyDirectAssignment wraps any direct assignment errors.
func WrapAnyDirectAssignment(varName, targetType string) error {
	return fmt.Errorf("%w: variable '%s' to type '%s'", ErrAnyDirectAssignment, varName, targetType)
}

// WrapAnyDirectFunctionArg wraps any direct function argument errors.
func WrapAnyDirectFunctionArg(funcName, expectedType string) error {
	return fmt.Errorf("%w: function '%s' expecting '%s'", ErrAnyDirectFunctionArg, funcName, expectedType)
}

// WrapAnyImplicitConversion wraps any implicit conversion errors.
func WrapAnyImplicitConversion(varName, targetType string) error {
	return fmt.Errorf("%w: variable '%s' to '%s'", ErrAnyImplicitConversion, varName, targetType)
}

// WrapAnyDirectVariableAccess wraps any direct variable access errors.
func WrapAnyDirectVariableAccess(varName string) error {
	return fmt.Errorf("%w: variable '%s'", ErrAnyDirectVariableAccess, varName)
}

// WrapAnyPatternNotExhaustive wraps any pattern not exhaustive errors.
func WrapAnyPatternNotExhaustive(missingTypes []string) error {
	return fmt.Errorf("%w: missing patterns for types [%s]", ErrAnyPatternNotExhaustive, strings.Join(missingTypes, ", "))
}

// WrapAnyPatternImpossible wraps any pattern impossible errors.
func WrapAnyPatternImpossible(patternType string, possibleTypes []string) error {
	return fmt.Errorf("%w: pattern '%s' not in documented types [%s]",
		ErrAnyPatternImpossible, patternType, strings.Join(possibleTypes, ", "))
}

// WrapAnyPatternUnreachable wraps any pattern unreachable errors.
func WrapAnyPatternUnreachable(patternType string) error {
	return fmt.Errorf("%w: pattern '%s'", ErrAnyPatternUnreachable, patternType)
}

// WrapWebSocketCreateServerWrongArgs wraps websocketCreateServer wrong arguments errors with argument count.
func WrapWebSocketCreateServerWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrWebSocketCreateServerWrongArgs, got)
}

// WrapWebSocketServerListenWrongArgs wraps websocketServerListen wrong arguments errors with argument count.
func WrapWebSocketServerListenWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrWebSocketServerListenWrongArgs, got)
}

// WrapWebSocketServerBroadcastWrongArgs wraps websocketServerBroadcast wrong arguments errors with argument count.
func WrapWebSocketServerBroadcastWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrWebSocketServerBroadcastWrongArgs, got)
}

// WrapWebSocketStopServerWrongArgs wraps websocketStopServer wrong arguments errors with argument count.
func WrapWebSocketStopServerWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrWebSocketStopServerWrongArgs, got)
}

// WrapHTTPStopServerUnknownNamedArg wraps httpStopServer unknown named argument errors.
func WrapHTTPStopServerUnknownNamedArg(namedArg string) error {
	return fmt.Errorf("%w '%s'", ErrHTTPStopServerUnknownNamedArg, namedArg)
}

// WrapHTTPStopServerWrongArgCount wraps httpStopServer wrong argument count errors.
func WrapHTTPStopServerWrongArgCount(got int) error {
	return fmt.Errorf("%w, got %d", ErrHTTPStopServerWrongArgCount, got)
}

// WrapMatchUnknownVariantType wraps match unknown variant type errors.
func WrapMatchUnknownVariantType(variant, typeName string) error {
	return fmt.Errorf("unknown variant '%s' in match expression for type '%s': %w",
		variant, typeName, ErrMatchUnknownVariantType)
}

// WrapLengthWrongArgs wraps length wrong arguments errors with argument count.
func WrapLengthWrongArgs(got int) error {
	return fmt.Errorf("length expects exactly 1 argument, got %d: %w", got, ErrLengthWrongArgs)
}

// WrapContainsWrongArgs wraps contains wrong arguments errors with argument count.
func WrapContainsWrongArgs(got int) error {
	return fmt.Errorf("contains expects exactly 2 arguments, got %d: %w", got, ErrContainsWrongArgs)
}

// WrapSubstringWrongArgs wraps substring wrong arguments error.
func WrapSubstringWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrSubstringWrongArgs, got)
}

// WrapSpawnProcessWrongArgs wraps spawnProcess wrong arguments error.
func WrapSpawnProcessWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrSpawnProcessWrongArgs, got)
}

// WrapSpawnProcessWithCallbackWrongArgs wraps spawnProcessWithCallback wrong arguments error.

// WrapAwaitProcessWrongArgs wraps awaitProcess wrong arguments error.
func WrapAwaitProcessWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrAwaitProcessWrongArgs, got)
}

// WrapCleanupProcessWrongArgs wraps cleanupProcess wrong arguments error.
func WrapCleanupProcessWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrCleanupProcessWrongArgs, got)
}

// WrapSleepWrongArgs wraps sleep wrong arguments error.
func WrapSleepWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrSleepWrongArgs, got)
}

// WrapWriteFileWrongArgs wraps writeFile wrong arguments error.
func WrapWriteFileWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrWriteFileWrongArgs, got)
}

// WrapReadFileWrongArgs wraps readFile wrong arguments error.
func WrapReadFileWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrReadFileWrongArgs, got)
}

// WrapParseJSONWrongArgs wraps parseJSON wrong arguments error.
func WrapParseJSONWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrParseJSONWrongArgs, got)
}

// WrapExtractCodeWrongArgs wraps extractCode wrong arguments error.
func WrapExtractCodeWrongArgs(got int) error {
	return fmt.Errorf("%w, got %d", ErrExtractCodeWrongArgs, got)
}
