package codegen

import (
	"fmt"
	"math/big"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) generateExpression(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral, *ast.StringLiteral, *ast.BooleanLiteral:

		return g.generateLiteralExpression(expr)
	case *ast.ListLiteral:

		return g.generateListLiteral(e)
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
	default:

		return g.generateFiberOrModuleExpression(expr)
	}
}

func (g *LLVMGenerator) getTypeOfExpression(expr ast.Expression) (string, error) {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return TypeString, nil
	case *ast.IntegerLiteral:
		return TypeInt, nil
	case *ast.BooleanLiteral:
		return TypeBool, nil
	case *ast.Identifier:
		if varType, ok := g.variableTypes[e.Name]; ok {
			return varType, nil
		}
		// Check for union type variants
		if _, ok := g.unionVariants[e.Name]; ok {
			return TypeInt, nil // Variants are i64 discriminants
		}
		return "", WrapUndefinedVariable(e.Name)
	case *ast.CallExpression:
		if fn, ok := e.Function.(*ast.Identifier); ok {
			if returnType, ok := g.functionReturnTypes[fn.Name]; ok {
				return returnType, nil
			}
			// It might be a user-defined function
			if _, ok := g.functions[fn.Name]; ok {
				if returnType, ok := g.functionReturnTypes[fn.Name]; ok {
					return returnType, nil
				}
				// If not in the map, it might not have been processed yet.
				// This is a limitation. For now, we'll assume int as a fallback
				// for user functions not yet in the map.
				return TypeInt, nil
			}
		}
		return "", WrapUnsupportedExpression(e)
	case *ast.BinaryExpression:
		// For simplicity, assuming numeric operations return Int
		// and comparisons return Bool.
		switch e.Operator {
		case "==", "!=", "<", "<=", ">", ">=":
			return TypeBool, nil
		default:
			return TypeInt, nil
		}
	case *ast.ResultExpression:
		// This is tricky. A result type is generic.
		// For now, let's say it's 'any'
		return TypeAny, nil
	case *ast.FieldAccessExpression:
		// This requires more sophisticated type tracking
		return TypeAny, nil
	default:
		return "", WrapUnsupportedExpression(e)
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
	if lit.Value {
		return constant.NewInt(types.I64, 1), nil
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
	// Check for regular variables first
	if val, exists := g.variables[ident.Name]; exists {
		// Check if this variable is of type 'any' - block direct access
		if varType, typeExists := g.variableTypes[ident.Name]; typeExists && varType == TypeAny {
			return nil, WrapAnyDirectVariableAccess(ident.Name)
		}

		return val, nil
	}

	// Check for union type variants (constants)
	if discriminant, exists := g.unionVariants[ident.Name]; exists {
		return constant.NewInt(types.I64, discriminant), nil
	}

	return nil, WrapUndefinedVariable(ident.Name)
}

func (g *LLVMGenerator) generateBinaryExpression(binExpr *ast.BinaryExpression) (value.Value, error) {
	// Validate that operands are not of type 'any' for arithmetic operations
	if err := g.validateNotAnyType(binExpr.Left, AnyOpArithmetic); err != nil {
		return nil, err
	}
	if err := g.validateNotAnyType(binExpr.Right, AnyOpArithmetic); err != nil {
		return nil, err
	}

	left, err := g.generateExpression(binExpr.Left)
	if err != nil {
		return nil, err
	}

	right, err := g.generateExpression(binExpr.Right)
	if err != nil {
		return nil, err
	}

	return g.generateBinaryOperation(binExpr.Operator, left, right)
}

// generateBinaryOperation generates the appropriate LLVM operation for the given operator.
func (g *LLVMGenerator) generateBinaryOperation(operator string, left, right value.Value) (value.Value, error) {
	switch operator {
	case "+", "-", "*", "/", "%":

		return g.generateArithmeticOperation(operator, left, right)
	case "==", "!=", "<", "<=", ">", ">=":

		return g.generateComparisonOperation(operator, left, right)
	default:

		return nil, WrapUnsupportedBinaryOp(operator)
	}
}

// generateArithmeticOperation generates LLVM arithmetic operations.
func (g *LLVMGenerator) generateArithmeticOperation(operator string, left, right value.Value) (value.Value, error) {
	switch operator {
	case "+":

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

		return nil, WrapUnsupportedBinaryOp(operator)
	}
}

// generateComparisonOperation generates LLVM comparison operations.
func (g *LLVMGenerator) generateComparisonOperation(operator string, left, right value.Value) (value.Value, error) {
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

		return nil, WrapUnsupportedBinaryOp(operator)
	}

	return g.builder.NewZExt(cmp, types.I64), nil
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

		return nil, WrapUnsupportedUnaryOp(unaryExpr.Operator)
	}
}

func (g *LLVMGenerator) generateResultExpression(resultExpr *ast.ResultExpression) (value.Value, error) {
	if resultExpr.IsSuccess {
		// Generate the actual value
		return g.generateExpression(resultExpr.Value)
	}
	// Return error sentinel value

	return constant.NewInt(types.I64, -1), nil
}

func (g *LLVMGenerator) generateFieldAccess(fieldAccess *ast.FieldAccessExpression) (value.Value, error) {
	// Validate that we're not trying to access fields on 'any' type
	if err := g.validateNotAnyType(fieldAccess.Object, AnyOpFieldAccess); err != nil {
		return nil, WrapAnyDirectFieldAccess(fieldAccess.FieldName)
	}

	// Check if this is field access on a constrained type constructor result
	if typeConstructor, isTypeConstructor := fieldAccess.Object.(*ast.TypeConstructorExpression); isTypeConstructor {
		// Check if this type has constraints
		if typeDecl, exists := g.typeDeclarations[typeConstructor.TypeName]; exists {
			// Check if any field has constraints
			if len(typeDecl.Variants) == 1 {
				variant := typeDecl.Variants[0]
				for _, field := range variant.Fields {
					if field.Constraint != nil {
						// This type has constraints, so constructor returns Result<T, E>
						// Field access on Result types should be a compilation error
						return nil, WrapFieldAccessOnResult(fieldAccess.FieldName, typeConstructor.TypeName)
					}
				}
			}
		}
	}

	// Check if this is field access on an identifier that might be a constrained type result
	if _, isIdent := fieldAccess.Object.(*ast.Identifier); isIdent {
		// Check if this identifier was declared as a constrained type constructor result
		// For now, we'll generate the object and then check if field access is valid
		obj, err := g.generateExpression(fieldAccess.Object)
		if err != nil {
			// If we can't generate the object, it's likely an undefined variable
			// Return the original error (which will be "undefined variable")
			return nil, err
		}

		// If we got here, the variable exists but field access may not be valid
		// For .value field access on result types, just return the object itself
		// since we're using simplified result types where the value IS the result
		if fieldAccess.FieldName == "value" {
			return obj, nil
		}

		// For other field access on identifiers, we need proper struct handling
		// This is where we would implement proper field access for non-Result types
		return nil, WrapFieldAccessNotImpl(fieldAccess.FieldName)
	}

	// Handle field access like a.value for other expression types
	obj, err := g.generateExpression(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

	// For .value field access on result types, just return the object itself
	// since we're using simplified result types where the value IS the result
	if fieldAccess.FieldName == "value" {
		return obj, nil
	}

	// For other field access, we'd need proper struct handling

	return nil, WrapFieldAccessNotImpl(fieldAccess.FieldName)
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
	// Look up the type declaration to get constraints
	typeDecl, exists := g.typeDeclarations[typeConstructor.TypeName]
	if !exists {
		return nil, WrapUndefinedType(typeConstructor.TypeName)
	}

	// For record types, validate constraints on each field
	if len(typeDecl.Variants) == 1 {
		variant := typeDecl.Variants[0]

		// Check each field for constraints
		for _, field := range variant.Fields {
			if field.Constraint != nil {
				// Get the field value from the constructor
				fieldValue, exists := typeConstructor.Fields[field.Name]
				if !exists {
					return nil, WrapMissingField(field.Name)
				}

				// Validate the constraint
				isValid, err := g.validateConstraint(field.Constraint, fieldValue)
				if err != nil {
					return nil, err
				}

				if !isValid {
					// Return constraint violation as Err variant
					// For now, we'll use a discriminant value that represents Err
					// TODO: Implement proper Result<T, E> type construction
					return constant.NewInt(types.I64, -1), nil
				}
			}
		}
	}

	// If all constraints pass, return success as Ok variant
	// For now, we'll use a discriminant value that represents Ok
	// TODO: Implement proper Result<T, E> type construction

	return constant.NewInt(types.I64, 1), nil
}

// validateConstraint validates a constraint function call against a field value.
func (g *LLVMGenerator) validateConstraint(
	constraint *ast.FunctionCallExpression,
	fieldValue ast.Expression,
) (bool, error) {
	// Handle compile-time evaluation for literal values
	switch fieldVal := fieldValue.(type) {
	case *ast.StringLiteral:

		return g.validateStringConstraint(constraint.Function, fieldVal.Value)
	case *ast.IntegerLiteral:

		return g.validateIntConstraint(constraint.Function, fieldVal.Value)
	case *ast.BooleanLiteral:

		return g.validateBoolConstraint(constraint.Function, fieldVal.Value)
	}

	// For non-literal values, we need runtime evaluation
	// Convert FunctionCallExpression to CallExpression for proper function calling
	callExpr := &ast.CallExpression{
		Function:  &ast.Identifier{Name: constraint.Function},
		Arguments: []ast.Expression{fieldValue},
	}

	// Generate the constraint function call
	result, err := g.generateCallExpression(callExpr)
	if err != nil {
		return false, err
	}

	// Check if the result is a compile-time constant
	if constResult, ok := result.(*constant.Int); ok {
		// Non-zero means constraint passed (check if X != 0)
		return constResult.X.Cmp(big.NewInt(0)) != 0, nil
	}

	// For runtime evaluation, we can't determine the result at compile time
	// So we'll assume it passes for now (this should be improved for runtime constraints)

	return true, nil
}

// validateStringConstraint evaluates string constraints at compile time.
func (g *LLVMGenerator) validateStringConstraint(constraintFunc, value string) (bool, error) {
	switch constraintFunc {
	case "notEmpty":

		return value != "", nil
	case "validEmail":

		return value != "", nil // Simple validation for now
	case "validLength":

		return value != "", nil
	default:

		return true, nil // Unknown constraint, assume valid
	}
}

// validateIntConstraint evaluates integer constraints at compile time.
func (g *LLVMGenerator) validateIntConstraint(constraintFunc string, value int64) (bool, error) {
	switch constraintFunc {
	case "isPositive":

		return value > 0, nil
	case "validAge":

		return value > 0, nil
	case "positive":

		return value > 0, nil
	default:

		return true, nil // Unknown constraint, assume valid
	}
}

// validateBoolConstraint evaluates boolean constraints at compile time.
func (g *LLVMGenerator) validateBoolConstraint(_ string, _ bool) (bool, error) {
	// Add boolean constraint logic here if needed

	return true, nil
}

// generateBlockExpression generates LLVM IR for block expressions.
func (g *LLVMGenerator) generateBlockExpression(blockExpr *ast.BlockExpression) (value.Value, error) {
	// Execute all statements in the block
	for _, stmt := range blockExpr.Statements {
		if err := g.generateStatement(stmt); err != nil {
			return nil, err
		}
	}

	// Return the final expression value, or a default value if no expression
	if blockExpr.Expression != nil {
		return g.generateExpression(blockExpr.Expression)
	}

	// If no return expression, return a default integer value

	return constant.NewInt(types.I64, 0), nil
}
