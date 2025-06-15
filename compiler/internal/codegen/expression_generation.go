package codegen

import (
	"math/big"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) generateExpression(expr ast.Expression) (value.Value, error) {
	return g.generateExpressionByType(expr)
}

func (g *LLVMGenerator) generateExpressionByType(expr ast.Expression) (value.Value, error) {
	result, handled, err := g.generateBasicExpressions(expr)
	if err != nil {
		return nil, err
	}
	if handled {
		return result, nil
	}
	return g.generateAdvancedExpressions(expr)
}

func (g *LLVMGenerator) generateBasicExpressions(expr ast.Expression) (value.Value, bool, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral, *ast.StringLiteral, *ast.BooleanLiteral:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateLiteralExpression(expr)
		})
	case *ast.InterpolatedStringLiteral:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateInterpolatedString(e)
		})
	case *ast.Identifier:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateIdentifier(e)
		})
	case *ast.BinaryExpression:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateBinaryExpression(e)
		})
	case *ast.FieldAccessExpression:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateFieldAccess(e)
		})
	case *ast.UnaryExpression:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateUnaryExpression(e)
		})
	case *ast.ResultExpression:
		return g.handleExpressionWithError(func() (value.Value, error) {
			return g.generateResultExpression(e)
		})
	}
	return nil, false, nil
}

// handleExpressionWithError is a helper function to reduce cyclomatic complexity
// by extracting the common error handling pattern.
func (g *LLVMGenerator) handleExpressionWithError(fn func() (value.Value, error)) (value.Value, bool, error) {
	result, err := fn()
	if err != nil {
		return nil, false, err
	}
	return result, true, nil
}

func (g *LLVMGenerator) generateAdvancedExpressions(expr ast.Expression) (value.Value, error) {
	switch e := expr.(type) {
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

// generateFiberOrModuleExpression handles fiber and module-related expressions.
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
	if err := g.validateNotAnyType(fieldAccess.Object, AnyOpFieldAccess); err != nil {
		return nil, WrapAnyDirectFieldAccess(fieldAccess.FieldName)
	}

	// Check if this is field access on a constrained type constructor result
	if typeConstructor, isTypeConstructor := fieldAccess.Object.(*ast.TypeConstructorExpression); isTypeConstructor {
		return g.handleTypeConstructorFieldAccess(fieldAccess, typeConstructor)
	}

	// Check if this is field access on an identifier
	if _, isIdent := fieldAccess.Object.(*ast.Identifier); isIdent {
		return g.handleIdentifierFieldAccess(fieldAccess)
	}

	return g.handleGeneralFieldAccess(fieldAccess)
}

func (g *LLVMGenerator) handleTypeConstructorFieldAccess(
	fieldAccess *ast.FieldAccessExpression,
	typeConstructor *ast.TypeConstructorExpression,
) (value.Value, error) {
	if typeDecl, exists := g.typeDeclarations[typeConstructor.TypeName]; exists {
		if len(typeDecl.Variants) == 1 {
			variant := typeDecl.Variants[0]
			for _, field := range variant.Fields {
				if field.Constraint != nil {
					return nil, WrapFieldAccessOnResult(fieldAccess.FieldName, typeConstructor.TypeName)
				}
			}
		}
	}
	return g.handleGeneralFieldAccess(fieldAccess)
}

func (g *LLVMGenerator) handleIdentifierFieldAccess(fieldAccess *ast.FieldAccessExpression) (value.Value, error) {
	obj, err := g.generateExpression(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

	if fieldAccess.FieldName == "value" {
		return obj, nil
	}

	return nil, WrapFieldAccessNotImpl(fieldAccess.FieldName)
}

func (g *LLVMGenerator) handleGeneralFieldAccess(fieldAccess *ast.FieldAccessExpression) (value.Value, error) {
	obj, err := g.generateExpression(fieldAccess.Object)
	if err != nil {
		return nil, err
	}

	if fieldAccess.FieldName == "value" {
		return obj, nil
	}

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
