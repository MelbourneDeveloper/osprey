package codegen

import (
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
	var elementType types.Type
	var elementSize int64

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
	var fieldTypes []types.Type
	var fieldValues []value.Value

	// Sort field names to ensure consistent ordering
	fieldNames := make([]string, 0, len(lit.Fields))
	for fieldName := range lit.Fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

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
		// Boolean NOT: convert to 0/1 and XOR with 1
		cmp := g.builder.NewICmp(enum.IPredEQ, operand, constant.NewInt(types.I64, 0))

		return g.builder.NewZExt(cmp, types.I64), nil
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
			// Look for Result< pattern in the type
			if strings.Contains(varType.String(), "Result<") {
				// This is field access on a Result type - requires pattern matching
				return nil, WrapConstraintResultFieldAccessWithPos(fieldAccess.FieldName, fieldAccess.Position)
			}
		}
	}

	// Generate the object value
	objectValue, err := g.generateExpression(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

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

	// If we can't find a record type, this is an error
	return nil, fmt.Errorf("line %d:%d: cannot access field '%s' on non-struct type", //nolint:err113
		fieldAccess.Position.Line, fieldAccess.Position.Column, fieldAccess.FieldName)
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
	var structType *types.StructType
	var isPointer bool

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

	// Find the field index by iterating through the record type fields in sorted order
	// The order should match the ObjectLiteral field iteration order
	fieldNames := make([]string, 0, len(recordType.fields))
	for fieldName := range recordType.fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	fieldIndex := -1
	for i, fieldName := range fieldNames {
		if fieldName == fieldAccess.FieldName {
			fieldIndex = i
			break
		}
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
	// Check field-level constraints first
	if len(typeDecl.Variants) > 0 {
		variant := &typeDecl.Variants[0] // Record types have one variant
		for _, field := range variant.Fields {
			if field.Constraint != nil {
				// Get the field value from the constructor
				fieldValue, exists := typeConstructor.Fields[field.Name]
				if !exists {
					return nil, WrapMissingField(field.Name)
				}
				
				// Generate the field value
				fieldLLVMValue, err := g.generateExpression(fieldValue)
				if err != nil {
					return nil, err
				}
				
				// Call the constraint function
				constraintFunc, exists := g.functions[field.Constraint.Function]
				if !exists {
					return nil, WrapUndefinedFunction(field.Constraint.Function)
				}
				
				// Call constraint function with field value
				result := g.builder.NewCall(constraintFunc, fieldLLVMValue)
				
				// Convert boolean result to integer (1 = true, 0 = false)
				resultAsInt := g.builder.NewZExt(result, types.I64)
				
				// Check if constraint failed (result == 0)
				zero := constant.NewInt(types.I64, 0)
				constraintPassed := g.builder.NewICmp(enum.IPredNE, resultAsInt, zero)
				
				// If constraint failed, return -1
				failureBlock := g.function.NewBlock("constraint_failure")
				successBlock := g.function.NewBlock("constraint_success")
				
				g.builder.NewCondBr(constraintPassed, successBlock, failureBlock)
				
				// Failure case: return -1
				g.builder = failureBlock
				g.builder.NewRet(constant.NewInt(types.I32, -1))
				
				// Success case: continue to next constraint
				g.builder = successBlock
			}
		}
	}
	
	// If all constraints passed or no constraints, return 1 for success
	return constant.NewInt(types.I64, 1), nil
}

// generateUnconstrainedRecordConstructor creates actual struct instances for unconstrained record types

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
			if err := g.generateStatement(stmt); err != nil {
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
		if err := g.generateStatement(lastStmt); err != nil {
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
	return fn
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
