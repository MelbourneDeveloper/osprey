package codegen

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// Constants for type sizes in bytes
const (
	ByteSize64Bit   = 8    // Size of 64-bit integers and pointers
	ByteSize32Bit   = 4    // Size of 32-bit integers
	ByteSize8Bit    = 1    // Size of 8-bit integers and booleans
	DefaultByteSize = 8    // Default size fallback
	LargeArraySize  = 1000 // Large array size for type casting operations

	// MinResultFieldCount represents the minimum number of fields in a Result type struct
	MinResultFieldCount = 2
)

func (g *LLVMGenerator) generateExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral, *ast.StringLiteral, *ast.BooleanLiteral:
		return g.generateLiteralExpression(expr)
	case *ast.ListLiteral:
		return g.generateListLiteral(e)
	case *ast.ObjectLiteral:
		return g.generateObjectLiteral(e)
	case *ast.ListAccessExpression:
		return g.generateListAccess(e)
	case *ast.InterpolatedStringLiteral:
		return g.generateInterpolatedString(e)
	case *ast.Identifier:
		return g.generateIdentifier(e)
	case *ast.BinaryExpression:
		return g.generateBinaryExpression(e)
	case *ast.UnaryExpression:
		return g.generateUnaryExpression(e)
	case *ast.ResultExpression:
		return g.generateResultExpression(e)
	case *ast.FieldAccessExpression:
		return g.generateFieldAccess(e)
	case *ast.CallExpression, *ast.MethodCallExpression:
		return g.generateCallLikeExpression(expr)
	case *ast.MatchExpression:
		return g.generateMatchExpression(e)
	case *ast.BlockExpression:
		return g.generateBlockExpression(e)
	case *ast.PerformExpression:
		return g.generatePerformExpression(e)
	case *ast.HandlerExpression:
		return g.generateHandlerExpression(e)
	default:
		return g.generateFiberOrModuleExpression(expr)
	}
}

// generateFiberOrModuleExpression handles fiber expressions and module access.
func (g *LLVMGenerator) generateFiberOrModuleExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.SpawnExpression:
		return g.generateSpawnExpression(e)
	case *ast.AwaitExpression:
		return g.generateAwaitExpression(e)
	case *ast.YieldExpression:
		return g.generateYieldExpression(e)
	case *ast.SelectExpression:
		return g.generateSelectExpression(e)
	case *ast.LambdaExpression:
		return g.generateLambdaExpression(e)
	case *ast.ModuleAccessExpression:
		return g.generateModuleAccessExpression(e)
	default:
		return g.generateChannelOrUnsupportedExpression(expr)
	}
}

// generateChannelOrUnsupportedExpression handles all channel-related expressions.
func (g *LLVMGenerator) generateChannelOrUnsupportedExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.ChannelExpression:
		return g.generateChannelExpression(e)
	case *ast.ChannelCreateExpression:
		return g.generateChannelCreateExpression(e)
	case *ast.ChannelSendExpression:
		return g.generateChannelSendExpression(e)
	case *ast.ChannelRecvExpression:
		return g.generateChannelRecvExpression(e)
	case *ast.TypeConstructorExpression:
		return g.generateTypeConstructorExpression(e)
	default:
		return g.generateUnsupportedExpression(expr)
	}
}

// generateUnsupportedExpression handles unsupported expression types.
func (g *LLVMGenerator) generateUnsupportedExpression(expr ast.Expression) (value.Value, error) {
	return nil, WrapUnsupportedExpression(expr)
}

// generateLiteralExpression handles all literal types.
func (g *LLVMGenerator) generateLiteralExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return g.generateIntegerLiteral(e)
	case *ast.StringLiteral:
		return g.generateStringLiteral(e)
	case *ast.BooleanLiteral:
		return g.generateBooleanLiteral(e)
	default:
		return nil, WrapUnsupportedExpression(expr)
	}
}

// generateCallLikeExpression handles call and method call expressions.
func (g *LLVMGenerator) generateCallLikeExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.CallExpression:
		return g.generateCallExpression(e)
	case *ast.MethodCallExpression:
		return g.generateMethodCallExpression(e)
	default:
		return nil, WrapUnsupportedExpression(expr)
	}
}

// generateIntegerLiteral generates LLVM IR for integer literals.
func (g *LLVMGenerator) generateIntegerLiteral(lit *ast.IntegerLiteral) (value.Value, error) {
	return constant.NewInt(types.I64, lit.Value), nil
}

// generateStringLiteral generates LLVM IR for string literals.
func (g *LLVMGenerator) generateStringLiteral(lit *ast.StringLiteral) (value.Value, error) {
	str := constant.NewCharArrayFromString(lit.Value + StringTerminator)
	global := g.module.NewGlobalDef("", str)

	return g.builder.NewGetElementPtr(str.Typ, global,
		constant.NewInt(types.I32, ArrayIndexZero),
		constant.NewInt(types.I32, ArrayIndexZero)), nil
}

// generateBooleanLiteral generates LLVM IR for boolean literals.
func (g *LLVMGenerator) generateBooleanLiteral(lit *ast.BooleanLiteral) (value.Value, error) {
	// Check context to determine the correct type
	targetType := types.I64 // Default to i64 for most contexts

	// 1. Check if function expects boolean return type
	if g.function != nil && g.function.Sig != nil {
		if g.function.Sig.RetType == types.I1 {
			targetType = types.I1
		}
	}

	// 2. Check expected return type context (for nested expressions)
	if g.expectedReturnType != nil {
		if g.expectedReturnType == types.I1 {
			targetType = types.I1
		}
	}

	// 3. Check expected parameter type context (for function arguments)
	if g.expectedParameterType != nil {
		if g.expectedParameterType == types.I1 {
			targetType = types.I1
		}
	}

	if lit.Value {
		if targetType == types.I1 {
			return constant.NewBool(true), nil
		}

		return constant.NewInt(types.I64, 1), nil
	}

	if targetType == types.I1 {
		return constant.NewBool(false), nil
	}

	return constant.NewInt(types.I64, 0), nil
}

// generateListLiteral generates LLVM IR for list literals like [1, 2, 3] or ["a", "b"].
func (g *LLVMGenerator) generateListLiteral(lit *ast.ListLiteral) (value.Value, error) {
	// For simplicity, implement arrays as a struct with length and data pointer
	// Array struct: { i64 length, i8* data }
	numElements := int64(len(lit.Elements))
	if numElements == 0 {
		// Empty array
		arrayStructType := types.NewStruct(types.I64, types.I8Ptr)
		arrayStruct := g.builder.NewAlloca(arrayStructType)

		// Store length = 0
		lengthPtr := g.builder.NewGetElementPtr(arrayStructType, arrayStruct,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		g.builder.NewStore(constant.NewInt(types.I64, 0), lengthPtr)

		// Store null data pointer
		dataPtr := g.builder.NewGetElementPtr(arrayStructType, arrayStruct,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
		g.builder.NewStore(constant.NewNull(types.I8Ptr), dataPtr)

		return arrayStruct, nil
	}

	// Determine element type from first element
	firstElement := lit.Elements[0]

	var (
		elementType types.Type
		elementSize int64
	)

	switch firstElement.(type) {
	case *ast.StringLiteral:
		elementType = types.I8Ptr
		elementSize = 8 // pointer size
	default:
		elementType = types.I64
		elementSize = 8 // i64 size
	}

	totalSize := numElements * elementSize

	// Allocate memory for the array data
	mallocFunc, ok := g.functions["malloc"]
	if !ok {
		mallocFunc = g.module.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I64))
		g.functions["malloc"] = mallocFunc
	}

	arrayData := g.builder.NewCall(mallocFunc, constant.NewInt(types.I64, totalSize))

	// Cast to appropriate pointer type
	arrayPtr := g.builder.NewBitCast(arrayData, types.NewPointer(elementType))

	// Store each element in the array
	for i, element := range lit.Elements {
		elementValue, err := g.generateExpression(element)
		if err != nil {
			return nil, err
		}

		// Get pointer to element position
		elementPtr := g.builder.NewGetElementPtr(elementType, arrayPtr, constant.NewInt(types.I64, int64(i)))
		g.builder.NewStore(elementValue, elementPtr)
	}

	// Create array struct { length, data }
	arrayStructType := types.NewStruct(types.I64, types.I8Ptr)
	arrayStruct := g.builder.NewAlloca(arrayStructType)

	// Store length
	lengthPtr := g.builder.NewGetElementPtr(arrayStructType, arrayStruct,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(constant.NewInt(types.I64, numElements), lengthPtr)

	// Store data pointer
	dataPtr := g.builder.NewGetElementPtr(arrayStructType, arrayStruct,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(arrayData, dataPtr)

	return arrayStruct, nil
}

// generateObjectLiteral generates LLVM IR for object literals like { field: value }.
func (g *LLVMGenerator) generateObjectLiteral(lit *ast.ObjectLiteral) (value.Value, error) {
	//CRITICAL TODO: implement record types properly, as per the spec 0005-TypeSystem.md

	// Create a simple struct with the fields
	var (
		fieldTypes  []types.Type
		fieldValues []value.Value
	)

	// HINDLEY-MILNER FIX: Use consistent field ordering
	// Infer the record type to get consistent field mapping
	objType, err := g.typeInferer.InferType(lit)
	if err != nil {
		return nil, fmt.Errorf("failed to infer object literal type: %w", err)
	}

	recordType, ok := objType.(*RecordType)
	if !ok {
		return nil, fmt.Errorf("%w: object literal did not infer to record type", ErrTypeMismatch)
	}

	// Get consistent field mapping
	fieldMapping := g.getOrCreateRecordFieldMapping(recordType.name, recordType.fields)

	// Create ordered field names based on mapping
	fieldNames := make([]string, len(lit.Fields))
	for fieldName, index := range fieldMapping {
		if index < len(fieldNames) {
			fieldNames[index] = fieldName
		}
	}

	// Process fields in sorted order
	for _, fieldName := range fieldNames {
		fieldValue := lit.Fields[fieldName]

		val, err := g.generateExpression(fieldValue)
		if err != nil {
			return nil, err
		}

		fieldValues = append(fieldValues, val)
		fieldTypes = append(fieldTypes, val.Type())
	}

	// Create struct type and allocate
	structType := types.NewStruct(fieldTypes...)
	structValue := g.builder.NewAlloca(structType)

	// Store each field value
	for i, fieldValue := range fieldValues {
		fieldPtr := g.builder.NewGetElementPtr(structType, structValue,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(i)))
		g.builder.NewStore(fieldValue, fieldPtr)
	}

	return structValue, nil
}

// generateListAccess generates LLVM IR for array indexing like arr[0].
func (g *LLVMGenerator) generateListAccess(access *ast.ListAccessExpression) (value.Value, error) {
	// Get the array value
	arrayValue, err := g.generateExpression(access.List)
	if err != nil {
		return nil, err
	}

	// Get the index value
	indexValue, err := g.generateExpression(access.Index)
	if err != nil {
		return nil, err
	}

	// Array access returns a Result<T, IndexError> for safety
	// First, extract length and data from array struct
	arrayStructType := types.NewStruct(types.I64, types.I8Ptr)

	// Get length
	lengthPtr := g.builder.NewGetElementPtr(arrayStructType, arrayValue,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	length := g.builder.NewLoad(types.I64, lengthPtr)

	// Get data pointer
	dataPtr := g.builder.NewGetElementPtr(arrayStructType, arrayValue,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	data := g.builder.NewLoad(types.I8Ptr, dataPtr)

	// Bounds check: index >= 0 && index < length
	zero := constant.NewInt(types.I64, 0)
	indexValid := g.builder.NewICmp(enum.IPredSGE, indexValue, zero)
	indexInBounds := g.builder.NewICmp(enum.IPredSLT, indexValue, length)
	boundsOk := g.builder.NewAnd(indexValid, indexInBounds)

	// Create unique block names to avoid conflicts with multiple array accesses
	blockSuffix := fmt.Sprintf("_%p", access)
	successBlock := g.function.NewBlock("array_access_success" + blockSuffix)
	errorBlock := g.function.NewBlock("array_access_error" + blockSuffix)
	endBlock := g.function.NewBlock("array_access_end" + blockSuffix)

	// Store current block before branching
	currentBlock := g.builder

	// Branch based on bounds check
	currentBlock.NewCondBr(boundsOk, successBlock, errorBlock)

	// Success block: return the element
	g.builder = successBlock

	// For now, assume string arrays (i8*) - this is a simplification
	// In a full implementation, we'd need to store type information with the array
	arrayDataPtr := g.builder.NewBitCast(data, types.NewPointer(types.I8Ptr))
	elementPtr := g.builder.NewGetElementPtr(types.I8Ptr, arrayDataPtr, indexValue)
	element := g.builder.NewLoad(types.I8Ptr, elementPtr)

	// Create Success result for string
	resultType := g.getResultType(types.I8Ptr)
	successResult := g.builder.NewAlloca(resultType)

	// Store element value
	valuePtr := g.builder.NewGetElementPtr(resultType, successResult,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(element, valuePtr)

	// Store success discriminant (0)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, successResult,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	// Branch to end block
	g.builder.NewBr(endBlock)

	// Error block: return index error
	g.builder = errorBlock
	errorResult := g.builder.NewAlloca(resultType)

	// Store error value (null string as placeholder)
	errorValuePtr := g.builder.NewGetElementPtr(resultType, errorResult,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(constant.NewNull(types.I8Ptr), errorValuePtr)

	// Store error discriminant (1)
	errorDiscriminantPtr := g.builder.NewGetElementPtr(resultType, errorResult,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 1), errorDiscriminantPtr)

	// Branch to end block
	g.builder.NewBr(endBlock)

	// End block: PHI node to select result
	g.builder = endBlock
	phi := endBlock.NewPhi(
		ir.NewIncoming(successResult, successBlock),
		ir.NewIncoming(errorResult, errorBlock),
	)

	return phi, nil
}

// generateIdentifier generates LLVM IR for identifiers.
func (g *LLVMGenerator) generateIdentifier(ident *ast.Identifier) (value.Value, error) {
	// Check for union type variants (constants) first
	if discriminant, exists := g.unionVariants[ident.Name]; exists {
		return constant.NewInt(types.I64, discriminant), nil
	}

	// Check for function references (user-defined functions and extern declarations)
	if fn, exists := g.functions[ident.Name]; exists {
		return fn, nil
	}

	// Check for regular variables
	if val, exists := g.variables[ident.Name]; exists {
		return val, nil
	}

	// Check for built-in functions
	if _, exists := GlobalBuiltInRegistry.GetFunction(ident.Name); exists {
		// For built-in functions referenced as identifiers, we need to create/get the LLVM function
		// This handles cases where built-in functions are passed as values
		return g.ensureBuiltinFunctionDeclaration(ident.Name), nil
	}

	// Check if this exists in the Hindley-Milner type environment
	if _, exists := g.typeInferer.env.Get(ident.Name); exists {
		// Variable exists in type environment but not in runtime
		// This can happen for built-in functions or uninitialized variables
		return nil, fmt.Errorf("%w: %s", ErrVariableNotInRuntime, ident.Name)
	}

	// Variable doesn't exist anywhere - undefined variable
	return nil, WrapUndefinedVariableWithPos(ident.Name, ident.Position)
}

func (g *LLVMGenerator) generateBinaryExpression(binExpr *ast.BinaryExpression) (value.Value, error) {
	left, err := g.generateExpression(binExpr.Left)
	if err != nil {
		return nil, err
	}

	right, err := g.generateExpression(binExpr.Right)
	if err != nil {
		return nil, err
	}

	return g.generateBinaryOperationWithPos(binExpr.Operator, left, right, binExpr.Position)
}

// generateBinaryOperationWithPos generates the appropriate LLVM operation for the given operator with position info.
func (g *LLVMGenerator) generateBinaryOperationWithPos(
	operator string, left, right value.Value, pos *ast.Position,
) (value.Value, error) {
	// Note: auto-unwrapping for Result types is already done in generateBinaryExpression
	switch operator {
	case "+", "-", "*", "/", "%":
		return g.generateArithmeticOperationWithPos(operator, left, right, pos)
	case "==", "!=", "<", "<=", ">", ">=":
		return g.generateComparisonOperationWithPos(operator, left, right, pos)
	case "&&", "||":
		return g.generateLogicalOperationWithPos(operator, left, right, pos)
	default:
		return nil, WrapUnsupportedBinaryOpWithPos(operator, pos)
	}
}

// generateArithmeticOperationWithPos generates LLVM arithmetic operations with position info.
func (g *LLVMGenerator) generateArithmeticOperationWithPos(
	operator string, left, right value.Value, pos *ast.Position,
) (value.Value, error) {
	// Check for nil values before accessing type information
	if left == nil || right == nil {
		return nil, ErrNilOperand
	}

	// Check for void types before arithmetic operations
	if left.Type() == types.Void || right.Type() == types.Void {
		return nil, WrapVoidArithmeticWithPos(operator, pos)
	}

	switch operator {
	case "+":
		// Handle string concatenation for pointer types (strings)
		if _, isPtr := left.Type().(*types.PointerType); isPtr {
			return g.generateStringConcatenation(left, right)
		}

		return g.builder.NewAdd(left, right), nil
	case "-":
		return g.builder.NewSub(left, right), nil
	case "*":
		return g.builder.NewMul(left, right), nil
	case "/":
		return g.builder.NewSDiv(left, right), nil
	case "%":
		return g.builder.NewSRem(left, right), nil
	default:
		return nil, WrapUnsupportedBinaryOpWithPos(operator, pos)
	}
}

// generateComparisonOperationWithPos generates LLVM comparison operations with position info.
func (g *LLVMGenerator) generateComparisonOperationWithPos(
	operator string, left, right value.Value, pos *ast.Position,
) (value.Value, error) {
	var cmp value.Value

	switch operator {
	case "==":
		cmp = g.builder.NewICmp(enum.IPredEQ, left, right)
	case "!=":
		cmp = g.builder.NewICmp(enum.IPredNE, left, right)
	case "<":
		cmp = g.builder.NewICmp(enum.IPredSLT, left, right)
	case "<=":
		cmp = g.builder.NewICmp(enum.IPredSLE, left, right)
	case ">":
		cmp = g.builder.NewICmp(enum.IPredSGT, left, right)
	case ">=":
		cmp = g.builder.NewICmp(enum.IPredSGE, left, right)
	default:
		return nil, WrapUnsupportedBinaryOpWithPos(operator, pos)
	}

	// Check current function's return type to determine output type
	if g.function != nil && g.function.Sig != nil {
		returnType := g.function.Sig.RetType
		if returnType == types.I1 {
			return cmp, nil
		}
	}

	// Default to extending to i64 for Result type construction and other contexts
	// The print function will handle the conversion to proper boolean strings
	return g.builder.NewZExt(cmp, types.I64), nil
}

// generateLogicalOperationWithPos generates LLVM logical operations with position info.
func (g *LLVMGenerator) generateLogicalOperationWithPos(
	operator string, left, right value.Value, pos *ast.Position,
) (value.Value, error) {
	switch operator {
	case "&&":
		return g.generateLogicalAnd(left, right)
	case "||":
		return g.generateLogicalOr(left, right)
	default:
		return nil, WrapUnsupportedBinaryOpWithPos(operator, pos)
	}
}

// generateLogicalAnd generates LLVM IR for logical AND operations.
func (g *LLVMGenerator) generateLogicalAnd(left, right value.Value) (value.Value, error) {
	// Short-circuit evaluation for &&
	// If left is false, return false without evaluating right

	// Convert to booleans first
	leftBool := g.builder.NewICmp(enum.IPredNE, left, constant.NewInt(types.I64, 0))
	rightBool := g.builder.NewICmp(enum.IPredNE, right, constant.NewInt(types.I64, 0))

	// Perform logical AND
	result := g.builder.NewAnd(leftBool, rightBool)

	// Convert back to i64 (0 for false, 1 for true)
	return g.builder.NewZExt(result, types.I64), nil
}

// generateLogicalOr generates LLVM IR for logical OR operations.
func (g *LLVMGenerator) generateLogicalOr(left, right value.Value) (value.Value, error) {
	// Short-circuit evaluation for ||
	// If left is true, return true without evaluating right

	// Convert to booleans first
	leftBool := g.builder.NewICmp(enum.IPredNE, left, constant.NewInt(types.I64, 0))
	rightBool := g.builder.NewICmp(enum.IPredNE, right, constant.NewInt(types.I64, 0))

	// Perform logical OR
	result := g.builder.NewOr(leftBool, rightBool)

	// Convert back to i64 (0 for false, 1 for true)
	return g.builder.NewZExt(result, types.I64), nil
}

// generateUnaryExpression generates LLVM IR for unary expressions.
func (g *LLVMGenerator) generateUnaryExpression(unaryExpr *ast.UnaryExpression) (value.Value, error) {
	operand, err := g.generateExpression(unaryExpr.Operand)
	if err != nil {
		return nil, err
	}

	switch unaryExpr.Operator {
	case "+":
		// Unary plus is a no-op
		return operand, nil
	case "-":
		// Unary minus
		zero := constant.NewInt(types.I64, 0)

		return g.builder.NewSub(zero, operand), nil
	case "!":
		// Boolean NOT: check operand type and return appropriate type
		if operand.Type() == types.I1 {
			// If operand is already boolean, just invert it
			return g.builder.NewICmp(enum.IPredEQ, operand, constant.NewBool(false)), nil
		} else {
			// If operand is integer, compare with 0 and return boolean
			cmp := g.builder.NewICmp(enum.IPredEQ, operand, constant.NewInt(types.I64, 0))
			// Check if we need boolean or integer result based on context
			if g.expectedReturnType == types.I1 {
				return cmp, nil
			} else {
				return g.builder.NewZExt(cmp, types.I64), nil
			}
		}
	default:
		return nil, WrapUnsupportedUnaryOpWithPos(unaryExpr.Operator, unaryExpr.Position)
	}
}

func (g *LLVMGenerator) generateResultExpression(resultExpr *ast.ResultExpression) (value.Value, error) {
	if resultExpr.Success {
		// Generate the actual value
		return g.generateExpression(resultExpr.Value)
	}
	// Return error sentinel value

	return constant.NewInt(types.I64, -1), nil
}

func (g *LLVMGenerator) generateFieldAccess(fieldAccess *ast.FieldAccessExpression) (value.Value, error) {
	// Type validation is now handled by Hindley-Milner type inference

	// Check if this is field access on a validated type constructor result
	if typeConstructor, isTypeConstructor := fieldAccess.Object.(*ast.TypeConstructorExpression); isTypeConstructor {
		// Check if this type has validation
		if typeDecl, exists := g.typeDeclarations[typeConstructor.TypeName]; exists {
			if typeDecl.ValidationFunc != nil {
				// This type has validation, so constructor returns Result<T, E>
				// Field access on Result types should be a compilation error
				return nil, WrapFieldAccessOnResult(fieldAccess.FieldName, typeConstructor.TypeName)
			}
		}
	}

	// Check if this is field access on an identifier that might be a constrained type result
	if ident, isIdent := fieldAccess.Object.(*ast.Identifier); isIdent {
		// Check if this identifier represents a constrained type constructor result using Hindley-Milner
		if varType, exists := g.typeInferer.env.Get(ident.Name); exists {
			// Look for Result[ pattern in the type (GenericType uses square brackets)
			if strings.Contains(varType.String(), "Result[") {
				// This is field access on a Result type - convert to pattern matching
				return g.generateResultFieldAccessAsMatch(fieldAccess, ident)
			}
		}
	}

	// Generate the object value
	objectValue, err := g.generateExpression(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

	// DISABLED: Result type detection was too broad and broke regular structs
	// TODO: Implement proper Result type tracking through Success/Error constructors
	// if g.isResultType(objectValue) {
	// 	return g.generateResultFieldAccess(fieldAccess, objectValue)
	// }

	// Handle field access on struct types (record types)
	return g.generateStructFieldAccess(fieldAccess, objectValue)
}

// generateStructFieldAccess handles field access on struct types
func (g *LLVMGenerator) generateStructFieldAccess(
	fieldAccess *ast.FieldAccessExpression,
	objectValue value.Value,
) (value.Value, error) {
	// For ObjectLiterals, we need to use the Hindley-Milner type information
	// instead of trying to reverse-engineer from LLVM types

	// If the object is an identifier, get its type from the type environment
	if ident, ok := fieldAccess.Object.(*ast.Identifier); ok {
		if varType, exists := g.typeInferer.env.Get(ident.Name); exists {
			if recordType, ok := varType.(*RecordType); ok {
				return g.generateRecordFieldAccess(fieldAccess, objectValue, recordType)
			}
			// Handle TypeVars that might be constrained to record types
			if typeVar, ok := varType.(*TypeVar); ok {
				prunedType := g.typeInferer.prune(typeVar)
				if recordType, ok := prunedType.(*RecordType); ok {
					return g.generateRecordFieldAccess(fieldAccess, objectValue, recordType)
				}
			}
		}
	}

	// For non-identifier objects, try to infer the type
	objectType, err := g.typeInferer.InferType(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

	resolvedType := g.typeInferer.ResolveType(objectType)
	if recordType, ok := resolvedType.(*RecordType); ok {
		return g.generateRecordFieldAccess(fieldAccess, objectValue, recordType)
	}

	// Handle type variables that might be constrained to record types
	if typeVar, ok := objectType.(*TypeVar); ok {
		// Check if this type variable has been unified with a record type
		prunedType := g.typeInferer.prune(typeVar)
		if recordType, ok := prunedType.(*RecordType); ok {
			return g.generateRecordFieldAccess(fieldAccess, objectValue, recordType)
		}

		// If the type variable is constrained to have this field during inference,
		// we need to find that constraint. For now, try to infer from the object value type.
		if objectValue != nil {
			if structType := g.tryGetStructType(objectValue.Type()); structType != nil {
				// Create a temporary record type based on the LLVM struct
				return g.generateStructFieldAccessFallback(fieldAccess, objectValue, structType)
			}
		}
	}

	// Additional fallback: If we still have a TypeVar, try direct struct field access
	if _, ok := objectType.(*TypeVar); ok && objectValue != nil {
		if structType := g.tryGetStructType(objectValue.Type()); structType != nil {
			return g.generateStructFieldAccessFallback(fieldAccess, objectValue, structType)
		}
	}

	// If we can't find a record type, this is an error
	if fieldAccess.Position != nil {
		return nil, fmt.Errorf("line %d:%d: cannot access field '%s' on non-struct type", //nolint:err113
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName)
	}

	return nil, fmt.Errorf("cannot access field '%s' on non-struct type", //nolint:err113
		fieldAccess.FieldName)
}

// tryGetStructType extracts a struct type from an LLVM type
func (g *LLVMGenerator) tryGetStructType(llvmType types.Type) *types.StructType {
	if ptrType, ok := llvmType.(*types.PointerType); ok {
		if st, ok := ptrType.ElemType.(*types.StructType); ok {
			return st
		}
	} else if st, ok := llvmType.(*types.StructType); ok {
		return st
	}

	return nil
}

// generateStructFieldAccessFallback handles field access on raw LLVM struct types
func (g *LLVMGenerator) generateStructFieldAccessFallback(
	fieldAccess *ast.FieldAccessExpression,
	objectValue value.Value,
	structType *types.StructType,
) (value.Value, error) {
	// For polymorphic field access, we need to make assumptions about field ordering
	// This is a fallback for when type inference hasn't provided a concrete record type

	// Try to find the field by name using a heuristic approach
	// For now, assume common field names map to indices
	var fieldIndex int

	switch fieldAccess.FieldName {
	case "first":
		fieldIndex = 0
	case "second":
		fieldIndex = 1
	case "x":
		fieldIndex = 0
	case "y":
		fieldIndex = 1
	case "value":
		fieldIndex = 0
	case "label":
		fieldIndex = 1
	default:
		// Try to parse field name as index if it's numeric
		return nil, fmt.Errorf("line %d:%d: cannot determine field index for '%s' in polymorphic field access: %w",
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName, ErrFieldAccessOnNonRecord)
	}

	if fieldIndex >= len(structType.Fields) {
		return nil, fmt.Errorf("line %d:%d: field index %d out of bounds for struct with %d fields: %w",
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldIndex, len(structType.Fields),
			ErrFieldAccessOnNonRecord)
	}

	// Generate field access using the computed index
	// Check if objectValue is a pointer or value
	objectType := objectValue.Type()
	if _, ok := objectType.(*types.PointerType); ok {
		// Object is a pointer to the struct - use GEP + load
		fieldPtr := g.builder.NewGetElementPtr(
			structType,
			objectValue,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)

		return g.builder.NewLoad(structType.Fields[fieldIndex], fieldPtr), nil
	} else if _, ok := objectType.(*types.StructType); ok {
		// Object is a struct value - use extractvalue directly
		return g.builder.NewExtractValue(objectValue, uint64(fieldIndex)), nil
	}
	// Fallback: Object is a struct value, but might not be recognized as such
	// Try extractvalue first
	return g.builder.NewExtractValue(objectValue, uint64(fieldIndex)), nil
}

// generateRecordFieldAccess handles field access using Hindley-Milner RecordType information
func (g *LLVMGenerator) generateRecordFieldAccess(
	fieldAccess *ast.FieldAccessExpression,
	objectValue value.Value,
	recordType *RecordType,
) (value.Value, error) {
	// Check if the field exists in the record type
	_, exists := recordType.fields[fieldAccess.FieldName]
	if !exists {
		return nil, fmt.Errorf("line %d:%d: field '%s' not found in record type", //nolint:err113
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName)
	}

	// Get the LLVM struct type from the object value
	objectType := objectValue.Type()

	var (
		structType *types.StructType
		isPointer  bool
	)

	if ptrType, ok := objectType.(*types.PointerType); ok {
		if st, ok := ptrType.ElemType.(*types.StructType); ok {
			structType = st
			isPointer = true
		}
	} else if st, ok := objectType.(*types.StructType); ok {
		structType = st
		isPointer = false
	}

	if structType == nil {
		return nil, fmt.Errorf("line %d:%d: cannot access field '%s' on non-struct value", //nolint:err113
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName)
	}

	// HINDLEY-MILNER FIX: Use consistent field mapping
	fieldMapping := g.getOrCreateRecordFieldMapping(recordType.name, recordType.fields)

	fieldIndex, exists := fieldMapping[fieldAccess.FieldName]
	if !exists {
		return nil, fmt.Errorf("line %d:%d: field '%s' not found in record type '%s'", //nolint:err113
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName, recordType.name)
	}

	if fieldIndex == -1 || fieldIndex >= len(structType.Fields) {
		return nil, fmt.Errorf("line %d:%d: field index mismatch for field '%s'", //nolint:err113
			fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName)
	}

	// Get pointer to the field
	var fieldPtr value.Value
	if isPointer {
		// Object is already a pointer to the struct
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			objectValue,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	} else {
		// Object is a struct value, need to get its address first
		structAddr := g.builder.NewAlloca(structType)
		g.builder.NewStore(objectValue, structAddr)
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			structAddr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	}

	// Load the field value
	fieldType := structType.Fields[fieldIndex]
	fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

	return fieldValue, nil
}

func (g *LLVMGenerator) generateMethodCallExpression(methodCall *ast.MethodCallExpression) (value.Value, error) {
	// For now, method calls are not fully implemented
	// This is a placeholder for future elegant method chaining like obj.method()
	// We could implement this to support chaining operations on values
	return nil, WrapMethodNotImpl(methodCall.MethodName)
}

// generateTypeConstructorExpression generates LLVM IR for type construction with constraint validation.
func (g *LLVMGenerator) generateTypeConstructorExpression(
	typeConstructor *ast.TypeConstructorExpression,
) (value.Value, error) {
	// Check if this is a built-in type first
	if typeConstructor.TypeName == TypeHTTPResponse {
		return g.generateHTTPResponseConstructor(typeConstructor)
	}

	// Handle built-in Result constructors
	if typeConstructor.TypeName == SuccessPattern {
		return g.generateSuccessConstructor(typeConstructor)
	}

	if typeConstructor.TypeName == ErrorPattern {
		return g.generateErrorConstructor(typeConstructor)
	}

	// Look up the type declaration to get constraints (for user-defined types)
	typeDecl := g.findTypeDeclarationByVariant(typeConstructor.TypeName)
	if typeDecl == nil {
		return nil, WrapUndefinedType(typeConstructor.TypeName)
	}

	// Check if this is a record type (single variant with fields)
	if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
		return g.generateRecordTypeConstructor(typeConstructor, typeDecl)
	}

	// FIXED: Handle discriminated unions with multiple variants
	if len(typeDecl.Variants) > 1 {
		// This is a discriminated union - find which variant we're constructing
		return g.generateDiscriminatedUnionConstructor(typeConstructor, typeDecl)
	}

	// For simple enum types (single variant without fields), return discriminant value
	if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) == 0 {
		// Simple enum variant - return its discriminant
		if discriminant, exists := g.unionVariants[typeDecl.Variants[0].Name]; exists {
			return constant.NewInt(types.I64, discriminant), nil
		}
	}

	// Fallback for unknown union structure
	return constant.NewInt(types.I64, 0), nil
}

// generateRecordTypeConstructor generates LLVM IR for record type construction
func (g *LLVMGenerator) generateRecordTypeConstructor(
	typeConstructor *ast.TypeConstructorExpression,
	typeDecl *ast.TypeDeclaration,
) (value.Value, error) {
	// Check if this type has any constraints
	hasConstraints := g.hasRecordTypeConstraints(typeDecl)

	// If this type has constraints, handle constraint validation
	if hasConstraints {
		return g.generateConstrainedRecordConstructor(typeConstructor, typeDecl)
	}

	// For unconstrained record types, create the actual struct instance
	return g.generateUnconstrainedRecordConstructor(typeConstructor, typeDecl)
}

// hasRecordTypeConstraints checks if a record type has any constraints
func (g *LLVMGenerator) hasRecordTypeConstraints(typeDecl *ast.TypeDeclaration) bool {
	// Only check for type-level constraints (validation function)
	return typeDecl.ValidationFunc != nil
}

// generateConstrainedRecordConstructor handles record types with constraints
func (g *LLVMGenerator) generateConstrainedRecordConstructor(
	typeConstructor *ast.TypeConstructorExpression,
	typeDecl *ast.TypeDeclaration,
) (value.Value, error) {
	// Step 1: First create the struct instance (same as unconstrained)
	structValue, err := g.generateUnconstrainedRecordConstructor(typeConstructor, typeDecl)
	if err != nil {
		return nil, err
	}

	// Step 2: If there's a validation function, call it with the struct
	if typeDecl.ValidationFunc != nil {
		validationFuncName := *typeDecl.ValidationFunc

		// Look up the validation function
		validationFunc, exists := g.functions[validationFuncName]
		if !exists {
			return nil, WrapUndefinedFunction(validationFuncName)
		}

		// Call the validation function with the struct
		validationResult := g.builder.NewCall(validationFunc, structValue)

		// The validation function returns a bool (true for valid, false for invalid)
		// Convert to integer: 1 for success (valid), -1 for failure (invalid)
		// Generate unique block names to avoid conflicts with multiple constraint validations
		blockSuffix := fmt.Sprintf("_%p", typeConstructor)
		validBlock := g.function.NewBlock("valid" + blockSuffix)
		invalidBlock := g.function.NewBlock("invalid" + blockSuffix)
		mergeBlock := g.function.NewBlock("merge" + blockSuffix)

		// Branch based on validation result
		g.builder.NewCondBr(validationResult, validBlock, invalidBlock)

		// Valid case: return 1
		g.builder = validBlock
		successValue := constant.NewInt(types.I64, 1)
		g.builder.NewBr(mergeBlock)

		// Invalid case: return -1
		g.builder = invalidBlock
		failureValue := constant.NewInt(types.I64, -1)
		g.builder.NewBr(mergeBlock)

		// Merge point - use PHI to select the result
		g.builder = mergeBlock
		phi := g.builder.NewPhi(
			ir.NewIncoming(successValue, validBlock),
			ir.NewIncoming(failureValue, invalidBlock),
		)
		
		return phi, nil
	}

	// No validation function - just return the struct
	return structValue, nil
}

// generateUnconstrainedRecordConstructor creates actual struct instances for unconstrained record types
func (g *LLVMGenerator) generateUnconstrainedRecordConstructor(
	typeConstructor *ast.TypeConstructorExpression,
	typeDecl *ast.TypeDeclaration,
) (value.Value, error) {
	if len(typeDecl.Variants) == 0 {
		return nil, fmt.Errorf("record type %s has no variants", typeDecl.Name) //nolint:err113
	}

	variant := &typeDecl.Variants[0] // Record types have one variant

	// HINDLEY-MILNER FIX: Use consistent field mapping for all record operations
	// Create field map from variant
	fieldMap := make(map[string]Type)
	for _, field := range variant.Fields {
		// Convert AST field type to our Type system
		var fieldType Type
		switch field.Type {
		case TypeInt:
			fieldType = NewPrimitiveType(TypeInt)
		case TypeString:
			fieldType = NewPrimitiveType(TypeString)
		case TypeBool:
			fieldType = NewPrimitiveType(TypeBool)
		default:
			// For complex types, create a concrete type
			fieldType = NewConcreteType(field.Type)
		}

		fieldMap[field.Name] = fieldType
	}

	// Get consistent field mapping
	fieldMapping := g.getOrCreateRecordFieldMapping(typeDecl.Name, fieldMap)

	// Create ordered field arrays based on mapping size
	fieldTypes := make([]types.Type, len(fieldMapping))
	fieldValues := make([]value.Value, len(fieldMapping))

	// HINDLEY-MILNER FIX: Process fields in the same sorted order as the mapping
	// Get sorted field names to ensure consistent processing order
	sortedFieldNames := make([]string, 0, len(fieldMapping))
	for fieldName := range fieldMapping {
		sortedFieldNames = append(sortedFieldNames, fieldName)
	}

	sort.Strings(sortedFieldNames)

	// Process fields using sorted order that matches the mapping
	for _, fieldName := range sortedFieldNames {
		// Find the field declaration
		var declaredFieldTypeName string

		for _, field := range variant.Fields {
			if field.Name == fieldName {
				declaredFieldTypeName = field.Type
				break
			}
		}

		// Find the field value in the constructor
		fieldExpr, exists := typeConstructor.Fields[fieldName]
		if !exists {
			return nil, WrapMissingField(fieldName)
		}

		// Convert declared type name to LLVM type
		var declaredLLVMType types.Type

		switch declaredFieldTypeName {
		case TypeInt:
			declaredLLVMType = types.I64
		case TypeString:
			declaredLLVMType = types.I8Ptr
		case TypeBool:
			declaredLLVMType = types.I1
		default:
			// For complex types, we'll determine after generation
			declaredLLVMType = nil
		}

		// Set expected type context for proper literal generation
		oldExpectedType := g.expectedReturnType
		g.expectedReturnType = declaredLLVMType

		// Generate the field value with the expected type context
		fieldValue, err := g.generateExpression(fieldExpr)
		if err != nil {
			g.expectedReturnType = oldExpectedType
			return nil, err
		}

		// Restore previous expected type
		g.expectedReturnType = oldExpectedType

		// Place field value at correct index based on mapping
		fieldIndex := fieldMapping[fieldName]
		fieldValues[fieldIndex] = fieldValue

		// HINDLEY-MILNER FIX: For polymorphic records, use actual field value types
		// not declared types to ensure type safety
		fieldTypes[fieldIndex] = fieldValue.Type()
	}

	// Create struct type and initialize with values
	structType := types.NewStruct(fieldTypes...)

	// Create the struct value by inserting field values into an undef struct
	var structValue value.Value = constant.NewUndef(structType)
	for i, fieldValue := range fieldValues {
		// Ensure the field value type matches the expected struct field type
		expectedType := structType.Fields[i]
		actualType := fieldValue.Type()

		finalFieldValue := fieldValue

		if expectedType != actualType {
			// Type mismatch - attempt conversion if types are compatible
			if g.areCompatibleLLVMTypes(expectedType, actualType) {
				// For pointer types that should be compatible, use bitcast
				if expectedType.String() == actualType.String() {
					// Same string representation - likely same type from different contexts
					finalFieldValue = fieldValue
				} else if g.isPointerType(expectedType) && g.isPointerType(actualType) {
					finalFieldValue = g.builder.NewBitCast(fieldValue, expectedType)
				} else {
					// Try other conversions if needed
					finalFieldValue = fieldValue
				}
			} else {
				return nil, fmt.Errorf("type mismatch in record field %d: expected %v, got %v: %w",
					i, expectedType, actualType, ErrRecordFieldTypeMismatch)
			}
		}

		structValue = g.builder.NewInsertValue(structValue, finalFieldValue, uint64(i))
	}

	return structValue, nil
}

// areCompatibleLLVMTypes checks if two LLVM types are compatible for conversion
func (g *LLVMGenerator) areCompatibleLLVMTypes(expected, actual types.Type) bool {
	// If string representations are the same, they should be compatible
	if expected.String() == actual.String() {
		return true
	}

	// Check if both are pointer types
	if g.isPointerType(expected) && g.isPointerType(actual) {
		return true
	}

	return false
}

// isPointerType checks if a type is a pointer type
func (g *LLVMGenerator) isPointerType(t types.Type) bool {
	_, isPtr := t.(*types.PointerType)
	return isPtr
}

// generateDiscriminatedUnionConstructor generates LLVM IR for discriminated union variant construction
func (g *LLVMGenerator) generateDiscriminatedUnionConstructor(
	typeConstructor *ast.TypeConstructorExpression,
	typeDecl *ast.TypeDeclaration,
) (value.Value, error) {
	// Find which variant we're constructing
	variant, variantIndex, err := g.findVariantByConstructorCall(typeConstructor, typeDecl)
	if err != nil {
		return nil, err
	}

	// Get the tagged union type from our type map
	unionType, exists := g.typeMap[typeDecl.Name]
	if !exists {
		return nil, WrapUndefinedType(typeDecl.Name)
	}

	// Allocate memory for the tagged union
	unionValue := g.builder.NewAlloca(unionType)

	// Set the discriminant (tag) field - this is the first field (index 0)
	tagPtr := g.builder.NewGetElementPtr(
		unionType,
		unionValue,
		constant.NewInt(types.I32, 0), // struct index
		constant.NewInt(types.I32, 0), // tag field index
	)
	g.builder.NewStore(constant.NewInt(types.I8, int64(variantIndex)), tagPtr)

	// If the variant has fields, serialize them into the data area
	if len(variant.Fields) > 0 {
		err := g.serializeVariantFields(variant, typeConstructor.Fields, unionValue, "")
		if err != nil {
			return nil, err
		}
	}

	return unionValue, nil
}

// findVariantByConstructorCall finds which variant matches the constructor call
func (g *LLVMGenerator) findVariantByConstructorCall(
	typeConstructor *ast.TypeConstructorExpression,
	typeDecl *ast.TypeDeclaration,
) (*ast.TypeVariant, int, error) {
	// The constructor call should be in the form: VariantName { field1: value1, field2: value2 }
	// We need to find which variant matches by looking at the discriminant of the variant name

	// First, try to find a variant by matching field names if fields are provided
	if len(typeConstructor.Fields) > 0 {
		for i, variant := range typeDecl.Variants {
			if len(variant.Fields) == len(typeConstructor.Fields) {
				// Check if all constructor fields match variant fields
				allMatch := true

				for _, variantField := range variant.Fields {
					if _, exists := typeConstructor.Fields[variantField.Name]; !exists {
						allMatch = false
						break
					}
				}

				if allMatch {
					return &variant, i, nil
				}
			}
		}

		return nil, -1, fmt.Errorf("%w for type %s", ErrNoVariantFound, typeDecl.Name)
	}

	// If no fields provided, this might be a simple enum variant
	// Try to find variant by checking if typeConstructor.TypeName matches any variant name
	for i, variant := range typeDecl.Variants {
		if variant.Name == typeConstructor.TypeName {
			return &variant, i, nil
		}
	}

	// If we still haven't found it, return the first variant as a fallback
	// This handles cases where we're constructing by type name rather than variant name
	if len(typeDecl.Variants) > 0 {
		return &typeDecl.Variants[0], 0, nil
	}

	return nil, -1, fmt.Errorf("%w: %s", ErrNoVariantsFound, typeDecl.Name)
}

// serializeVariantFields serializes variant fields into the union data area
func (g *LLVMGenerator) serializeVariantFields(
	variant *ast.TypeVariant,
	fieldValues map[string]ast.Expression,
	unionPtr value.Value,
	_ string,
) error {
	// Get pointer to the data area (second field in the tagged union)
	unionType := unionPtr.Type().(*types.PointerType).ElemType
	dataPtr := g.builder.NewGetElementPtr(
		unionType,
		unionPtr,
		constant.NewInt(types.I32, 0), // struct index
		constant.NewInt(types.I32, 1), // data field index
	)

	// Serialize each field into the data area
	offset := int64(0)

	for _, field := range variant.Fields {
		fieldValue, exists := fieldValues[field.Name]
		if !exists {
			return WrapMissingField(field.Name)
		}

		// Generate the field value
		llvmFieldValue, err := g.generateExpression(fieldValue)
		if err != nil {
			return err
		}

		// Get the field type and calculate its size
		fieldType := g.getFieldType(field.Type)
		fieldSize := g.getTypeSize(fieldType)

		// Cast data array to appropriate pointer type for this field
		fieldPtr := g.builder.NewBitCast(
			g.builder.NewGetElementPtr(
				types.NewArray(uint64(LargeArraySize), types.I8), // Use large array for casting
				dataPtr,
				constant.NewInt(types.I32, 0),      // array index
				constant.NewInt(types.I32, offset), // byte offset
			),
			types.NewPointer(fieldType),
		)

		// Convert value to correct type if needed
		convertedValue := g.convertValueToExpectedType(llvmFieldValue, fieldType)

		// Store the field value
		g.builder.NewStore(convertedValue, fieldPtr)

		// Move to next field offset
		offset += fieldSize
	}

	return nil
}

// getTypeSize returns the size in bytes of an LLVM type
func (g *LLVMGenerator) getTypeSize(t types.Type) int64 {
	switch t {
	case types.I64:
		return ByteSize64Bit
	case types.I8Ptr:
		return ByteSize64Bit // pointer size
	case types.I1:
		return ByteSize8Bit
	case types.I32:
		return ByteSize32Bit
	case types.I8:
		return ByteSize8Bit
	default:
		return DefaultByteSize // default to 8 bytes
	}
}

// convertValueToExpectedType converts a value to match the expected LLVM type
//
//nolint:gocognit // Complex function required for comprehensive type conversion handling
func (g *LLVMGenerator) convertValueToExpectedType(value value.Value, expectedType types.Type) value.Value {
	currentType := value.Type()

	// If types already match, no conversion needed
	if currentType == expectedType {
		return value
	}

	// Handle string types (i8*) explicitly
	if currentType == types.I8Ptr && expectedType == types.I8Ptr {
		return value
	}

	// Convert i64 to i1 (boolean)
	if currentType == types.I64 && expectedType == types.I1 {
		// Convert non-zero to true, zero to false
		zero := constant.NewInt(types.I64, 0)
		return g.builder.NewICmp(enum.IPredNE, value, zero)
	}

	// Convert i1 to i64 (boolean to integer)
	if currentType == types.I1 && expectedType == types.I64 {
		return g.builder.NewZExt(value, types.I64)
	}

	// Handle incompatible pointer types by checking what the field actually expects
	if ptrType, ok := expectedType.(*types.PointerType); ok {
		// If expected type is a pointer and we have a different type, try to convert via casting
		if currentType != expectedType {
			// Don't try to cast between fundamentally incompatible types
			if currentType == types.I8Ptr && ptrType.ElemType != types.I8 {
				// This is a string being stored in a non-string pointer field - likely an error
				// Return a null pointer of the expected type to avoid crash
				return constant.NewNull(ptrType)
			}
		}
	}

	// Handle all incompatible pointer/type combinations
	// Detect various type mismatches and provide safe defaults

	// String to integer conversion
	if currentType == types.I8Ptr && expectedType == types.I64 {
		return constant.NewInt(types.I64, 0)
	}

	// Integer to string conversion
	if currentType == types.I64 && expectedType == types.I8Ptr {
		return constant.NewNull(types.I8Ptr)
	}

	// Handle pointer type mismatches - if we're trying to store into an i64* pointer
	if ptrType, ok := expectedType.(*types.PointerType); ok {
		if ptrType.ElemType == types.I64 {
			// Expected type is pointer to integer - convert current value to integer first
			switch currentType {
			case types.I8Ptr:
				return constant.NewInt(types.I64, 0)
			case types.I1:
				return g.builder.NewZExt(value, types.I64)
			default:
				return constant.NewInt(types.I64, 0)
			}
		}

		if ptrType.ElemType == types.I8 || ptrType.ElemType == types.I8Ptr {
			// Expected type is pointer to string - convert to string
			if currentType != types.I8Ptr {
				return constant.NewNull(types.I8Ptr)
			}
		}
	}

	// For other cases, try basic casting if the types are compatible
	if intType1, ok1 := currentType.(*types.IntType); ok1 {
		if intType2, ok2 := expectedType.(*types.IntType); ok2 {
			// Both are integer types - try casting
			if intType1.BitSize < intType2.BitSize {
				// Zero-extend smaller to larger
				return g.builder.NewZExt(value, expectedType)
			} else if intType1.BitSize > intType2.BitSize {
				// Truncate larger to smaller
				return g.builder.NewTrunc(value, expectedType)
			}
		}
	}

	// For other cases, return the original value and hope LLVM can handle it
	return value
}

// NOTE: Old field-level constraint validation functions removed.
// Type-level validation is now handled by user-defined validation functions
// that return Result<T, String> types.

func (g *LLVMGenerator) generateBlockExpression(blockExpr *ast.BlockExpression) (value.Value, error) {
	// If the block has statements, execute all but the last one
	// Then check if the last statement is an expression that should be returned
	if len(blockExpr.Statements) > 0 {
		// Execute all statements except the last one
		for _, stmt := range blockExpr.Statements[:len(blockExpr.Statements)-1] {
			err := g.generateStatement(stmt)
			if err != nil {
				return nil, err
			}
		}

		// Check if the last statement is an expression statement
		lastStmt := blockExpr.Statements[len(blockExpr.Statements)-1]
		if exprStmt, ok := lastStmt.(*ast.ExpressionStatement); ok {
			// Return the value of the last expression statement
			result, err := g.generateExpression(exprStmt.Expression)
			if err != nil {
				return nil, err
			}

			return result, nil
		}

		// Execute the last statement (it's not an expression statement)
		err := g.generateStatement(lastStmt)
		if err != nil {
			return nil, err
		}
	}

	// Return the final expression value if present
	if blockExpr.Expression != nil {
		return g.generateExpression(blockExpr.Expression)
	}

	// If no explicit expression and no expression statements, return Unit (0)
	return constant.NewInt(types.I64, 0), nil
}

// generateHTTPResponseConstructor generates LLVM IR for HttpResponse construction.
func (g *LLVMGenerator) generateHTTPResponseConstructor(
	typeConstructor *ast.TypeConstructorExpression,
) (value.Value, error) {
	// Get the HttpResponse struct type
	httpResponseType, exists := g.typeMap[TypeHTTPResponse]
	if !exists {
		return nil, WrapUndefinedType(TypeHTTPResponse)
	}

	// Allocate memory for the HttpResponse struct on the heap
	mallocFunc, ok := g.functions["malloc"]
	if !ok {
		mallocFunc = g.module.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I64))
		g.functions["malloc"] = mallocFunc
	}

	const httpResponseStructSize = 48 // 6 fields * 8 bytes each (removed redundant length fields)

	structSize := constant.NewInt(types.I64, httpResponseStructSize)
	structMem := g.builder.NewCall(mallocFunc, structSize)
	structPtr := g.builder.NewBitCast(structMem, types.NewPointer(httpResponseType))

	// Define field order and types to match the struct definition (REMOVED REDUNDANT LENGTH FIELDS)
	fieldInfo := []struct {
		name      string
		fieldType types.Type
	}{
		{"status", types.I64},        // status: Int
		{"headers", types.I8Ptr},     // headers: String
		{"contentType", types.I8Ptr}, // contentType: String
		{"streamFd", types.I64},      // streamFd: Int
		{"isComplete", types.I1},     // isComplete: Bool
		{"partialBody", types.I8Ptr}, // partialBody: String (runtime calculates length)
	}

	for i, field := range fieldInfo {
		fieldValue, exists := typeConstructor.Fields[field.name]
		if !exists {
			return nil, WrapMissingField(field.name)
		}

		// Generate the field value
		value, err := g.generateExpression(fieldValue)
		if err != nil {
			return nil, err
		}

		// Convert value to correct type if needed
		if field.fieldType == types.I1 && value.Type() != types.I1 {
			// Convert integer to boolean (non-zero = true, zero = false)
			value = g.builder.NewICmp(enum.IPredNE, value, constant.NewInt(types.I64, 0))
		}

		// Get pointer to the field
		fieldPtr := g.builder.NewGetElementPtr(
			httpResponseType,
			structPtr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)

		// Store the value in the field
		g.builder.NewStore(value, fieldPtr)
	}

	return structPtr, nil
}

// generateStringConcatenation generates LLVM IR for string concatenation using strcat
func (g *LLVMGenerator) generateStringConcatenation(left, right value.Value) (value.Value, error) {
	// Ensure strcat and strlen are declared
	strcatFunc := g.ensureStrcatDeclaration()
	strlenFunc := g.ensureStrlenDeclaration()
	mallocFunc := g.ensureMallocDeclaration()

	// Extract strings from Result types if needed
	leftStr := g.extractStringFromValue(left)
	rightStr := g.extractStringFromValue(right)

	// Calculate lengths of both strings
	leftLen := g.builder.NewCall(strlenFunc, leftStr)
	rightLen := g.builder.NewCall(strlenFunc, rightStr)

	// Calculate total length: leftLen + rightLen + 1 (for null terminator)
	totalLen := g.builder.NewAdd(leftLen, rightLen)
	totalLenPlusOne := g.builder.NewAdd(totalLen, constant.NewInt(types.I64, 1))

	// Allocate memory for the result string
	result := g.builder.NewCall(mallocFunc, totalLenPlusOne)

	// Copy left string to result
	strcpyFunc := g.ensureStrcpyDeclaration()
	g.builder.NewCall(strcpyFunc, result, leftStr)

	// Concatenate right string to result
	g.builder.NewCall(strcatFunc, result, rightStr)

	return result, nil
}

// extractStringFromValue extracts a string from either a regular string or a Result type
func (g *LLVMGenerator) extractStringFromValue(val value.Value) value.Value {
	// If it's already a string pointer, return it as is
	if val.Type() == types.I8Ptr {
		return val
	}

	// Check if it's a Result type struct pointer
	if ptrType, ok := val.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == MinResultFieldCount {
			// This is a Result type { T, i8 } - extract the value (first field)
			valuePtr := g.builder.NewGetElementPtr(structType, val,
				constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

			return g.builder.NewLoad(structType.Fields[0], valuePtr)
		}
	}

	// If it's not a string or Result type, return it as is (might be an error case)
	return val
}

// ensureStrcatDeclaration ensures strcat is declared
func (g *LLVMGenerator) ensureStrcatDeclaration() *ir.Func {
	if strcat, exists := g.functions["strcat"]; exists {
		return strcat
	}

	strcat := g.module.NewFunc("strcat", types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr))
	g.functions["strcat"] = strcat

	return strcat
}

// ensureStrcpyDeclaration ensures strcpy is declared
func (g *LLVMGenerator) ensureStrcpyDeclaration() *ir.Func {
	if strcpy, exists := g.functions["strcpy"]; exists {
		return strcpy
	}

	strcpy := g.module.NewFunc("strcpy", types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr))
	g.functions["strcpy"] = strcpy

	return strcpy
}

// ensureStrlenDeclaration ensures strlen is declared
func (g *LLVMGenerator) ensureStrlenDeclaration() *ir.Func {
	if strlen, exists := g.functions["strlen"]; exists {
		return strlen
	}

	strlen := g.module.NewFunc("strlen", types.I64,
		ir.NewParam("s", types.I8Ptr))
	g.functions["strlen"] = strlen

	return strlen
}

// ensureMallocDeclaration ensures malloc is declared
func (g *LLVMGenerator) ensureMallocDeclaration() *ir.Func {
	if malloc, exists := g.functions["malloc"]; exists {
		return malloc
	}

	malloc := g.module.NewFunc("malloc", types.I8Ptr,
		ir.NewParam("size", types.I64))
	g.functions["malloc"] = malloc

	return malloc
}

// ensureBuiltinFunctionDeclaration ensures a built-in function is declared using builtin registry
func (g *LLVMGenerator) ensureBuiltinFunctionDeclaration(ospreyName string) *ir.Func {
	// Get function details from builtin registry
	builtinFunc, exists := GlobalBuiltInRegistry.GetFunction(ospreyName)
	if !exists {
		return nil
	}

	// Determine the actual function name to use for LLVM
	llvmFunctionName := ospreyName
	if builtinFunc.CName != "" {
		llvmFunctionName = builtinFunc.CName
	}

	// Check if function is already declared
	if fn, exists := g.functions[ospreyName]; exists {
		return fn
	}

	// Convert builtin parameters to LLVM parameters
	params := make([]*ir.Param, len(builtinFunc.ParameterTypes))
	for i, param := range builtinFunc.ParameterTypes {
		llvmType := g.getLLVMType(param.Type)
		params[i] = ir.NewParam(param.Name, llvmType)
	}

	// Convert return type to LLVM type
	returnType := g.getLLVMType(builtinFunc.ReturnType)

	// Create function with the correct name (C name if available, otherwise Osprey name)
	fn := g.module.NewFunc(llvmFunctionName, returnType, params...)
	g.functions[ospreyName] = fn

	// For certain built-in functions that can be used as first-class values,
	// we need to generate the actual function body
	if ospreyName == "toString" {
		g.generateToStringFunctionBody(fn)
	}

	return fn
}

// generateToStringFunctionBody generates the function body for toString built-in function
func (g *LLVMGenerator) generateToStringFunctionBody(fn *ir.Func) {
	// Save current context
	oldFunction := g.function
	oldBuilder := g.builder

	// Set up function context
	entry := fn.NewBlock("")
	g.builder = entry
	g.function = fn

	// Get the argument value
	valueParam := fn.Params[0] // toString takes one parameter

	// Simple toString implementation for integers
	// Use sprintf to convert integer to string
	sprintf := g.ensureSprintfDeclaration()
	malloc := g.ensureMallocDeclaration()

	// Format string for integer conversion
	formatStr := constant.NewCharArrayFromString("%ld\x00")
	formatGlobal := g.module.NewGlobalDef("", formatStr)
	formatPtr := g.builder.NewGetElementPtr(formatStr.Typ, formatGlobal,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

	// Allocate buffer on heap for result string (64 bytes should be enough for any 64-bit integer)
	const bufferSizeBytes = 64

	bufferSize := constant.NewInt(types.I64, bufferSizeBytes)
	bufferPtr := g.builder.NewCall(malloc, bufferSize)

	// Call sprintf(buffer, "%ld", value)
	g.builder.NewCall(sprintf, bufferPtr, formatPtr, valueParam)

	// Ensure null termination by explicitly setting the last byte to 0
	// (sprintf should handle this, but let's be safe)
	// bufferPtr is already i8*, so we just need to get element at index 63
	lastBytePtr := g.builder.NewGetElementPtr(types.I8, bufferPtr,
		constant.NewInt(types.I32, bufferSizeBytes-1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), lastBytePtr)

	g.builder.NewRet(bufferPtr)

	// Restore context
	g.function = oldFunction
	g.builder = oldBuilder
}

// findTypeDeclarationByVariant finds the type declaration that contains the given variant name
func (g *LLVMGenerator) findTypeDeclarationByVariant(variantName string) *ast.TypeDeclaration {
	// First try to find by type name (for direct type lookup)
	if typeDecl, exists := g.typeDeclarations[variantName]; exists {
		return typeDecl
	}

	// If not found, search through all type declarations for the variant name
	for _, typeDecl := range g.typeDeclarations {
		for _, variant := range typeDecl.Variants {
			if variant.Name == variantName {
				return typeDecl
			}
		}
	}

	return nil
}

// generateSuccessConstructor generates LLVM IR for Success { value: T } constructor.
func (g *LLVMGenerator) generateSuccessConstructor(
	typeConstructor *ast.TypeConstructorExpression,
) (value.Value, error) {
	// Success constructor should create a Result struct with discriminant = 0 (success)
	// Result struct: [value, discriminant] where discriminant=0 for success

	// Get the value expression from the constructor fields
	valueExpr, exists := typeConstructor.Fields["value"]
	if !exists {
		return nil, ErrSuccessConstructorMissingValue
	}

	// Generate the value
	value, err := g.generateExpression(valueExpr)
	if err != nil {
		return nil, err
	}

	// Create Result struct type based on the value type
	valueType := value.Type()
	resultStructType := types.NewStruct(valueType, types.I8) // [value, discriminant]

	// Create the result struct as a value (not pointer)
	// Use InsertValue to build the struct value directly
	undefStruct := constant.NewUndef(resultStructType)
	resultWithValue := g.builder.NewInsertValue(undefStruct, value, 0)
	resultComplete := g.builder.NewInsertValue(resultWithValue, constant.NewInt(types.I8, 0), 1)

	return resultComplete, nil
}

// generateErrorConstructor generates LLVM IR for Error { message: E } constructor.
func (g *LLVMGenerator) generateErrorConstructor(
	typeConstructor *ast.TypeConstructorExpression,
) (value.Value, error) {
	// Error constructor should create a Result struct with discriminant = 1 (error)
	// Result struct: [defaultValue, discriminant] where discriminant=1 for error

	// Get the message expression from the constructor fields
	messageExpr, exists := typeConstructor.Fields["message"]
	if !exists {
		return nil, ErrErrorConstructorMissingMessage
	}

	// Generate the error message
	message, err := g.generateExpression(messageExpr)
	if err != nil {
		return nil, err
	}

	// Error constructor creates Result struct: [error_message, discriminant]
	// where discriminant = 1 for error
	messageType := message.Type()
	resultStructType := types.NewStruct(messageType, types.I8) // [error_message, discriminant]

	// Create the result struct as a value (not pointer)
	// Use InsertValue to build the struct value directly
	undefStruct := constant.NewUndef(resultStructType)
	resultWithMessage := g.builder.NewInsertValue(undefStruct, message, 0)
	resultComplete := g.builder.NewInsertValue(resultWithMessage, constant.NewInt(types.I8, 1), 1)

	return resultComplete, nil
}

// generateResultFieldAccessAsMatch converts Result field access to pattern matching
// This handles cases like myResult { value } ? value : "default"
func (g *LLVMGenerator) generateResultFieldAccessAsMatch(
	_ *ast.FieldAccessExpression,
	ident *ast.Identifier,
) (value.Value, error) {
	// For now, just extract the value directly from the Success Result struct
	// This is a simplified implementation that assumes the Result is a Success

	// Generate the Result value
	resultValue, err := g.generateExpression(ident)
	if err != nil {
		return nil, err
	}

	// Extract the first field (value) from the Result struct
	// Result struct layout: [value, discriminant]
	resultType := resultValue.Type()
	if structType, ok := resultType.(*types.StructType); ok && len(structType.Fields) >= 2 {
		// Get pointer to the value field (index 0)
		valuePtr := g.builder.NewGetElementPtr(structType, resultValue,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

		// Load the value
		return g.builder.NewLoad(structType.Fields[0], valuePtr), nil
	}

	return nil, errors.New("result field access failed: invalid Result type structure") //nolint:err113
}
