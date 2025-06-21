package ast

import (
	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildBinaryExpr(ctx parser.IBinaryExprContext) Expression {
	return b.buildComparisonExpr(ctx.ComparisonExpr())
}

func (b *Builder) buildComparisonExpr(ctx parser.IComparisonExprContext) Expression {
	addExprs := ctx.AllAddExpr()
	if len(addExprs) == 1 {
		return b.buildAddExpr(addExprs[0])
	}

	// Build left-associative comparison expression
	left := b.buildAddExpr(addExprs[0])

	for i := 1; i < len(addExprs); i++ {
		right := b.buildAddExpr(addExprs[i])

		// Determine comparison operator
		operator := "=="
		if ctx.NE_OP(i-1) != nil {
			operator = "!="
		} else if ctx.LT(i-1) != nil {
			operator = "<"
		} else if ctx.GT(i-1) != nil {
			operator = ">"
		} else if ctx.LE_OP(i-1) != nil {
			operator = "<="
		} else if ctx.GE_OP(i-1) != nil {
			operator = ">="
		}

		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (b *Builder) buildAddExpr(ctx parser.IAddExprContext) Expression {
	mulExprs := ctx.AllMulExpr()
	if len(mulExprs) == 1 {
		return b.buildMulExpr(mulExprs[0])
	}

	// Build left-associative addition/subtraction expression
	left := b.buildMulExpr(mulExprs[0])

	for i := 1; i < len(mulExprs); i++ {
		right := b.buildMulExpr(mulExprs[i])

		// Determine operator (+ or -)
		operator := "+"
		if ctx.MINUS(i-1) != nil {
			operator = "-"
		}

		// Wrap arithmetic operations in result types
		binExpr := &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}

		left = b.wrapInResultType(binExpr)
	}

	return left
}

func (b *Builder) buildMulExpr(ctx parser.IMulExprContext) Expression {
	unaryExprs := ctx.AllUnaryExpr()
	if len(unaryExprs) == 1 {
		return b.buildUnaryExpr(unaryExprs[0])
	}

	// Build left-associative multiplication/division expression
	left := b.buildUnaryExpr(unaryExprs[0])

	for i := 1; i < len(unaryExprs); i++ {
		right := b.buildUnaryExpr(unaryExprs[i])

		// Determine operator (*, /, or %)
		operator := "*"
		if ctx.SLASH(i-1) != nil {
			operator = "/"
		} else if ctx.MOD_OP(i-1) != nil {
			operator = "%"
		}

		// Wrap arithmetic operations in result types
		binExpr := &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}

		left = b.wrapInResultType(binExpr)
	}

	return left
}

// Add this helper function to wrap expressions in result types.
func (b *Builder) wrapInResultType(expr Expression) Expression {
	if binExpr, ok := expr.(*BinaryExpression); ok && b.isArithmeticOperation(binExpr.Operator) {
		// Check for potential division by zero
		if binExpr.Operator == "/" {
			if intLit, ok := binExpr.Right.(*IntegerLiteral); ok && intLit.Value == 0 {
				// Division by zero - return error result
				return &ResultExpression{
					Success: false,
					Value:     &StringLiteral{Value: "Division by zero"},
					ErrorType: "DivisionByZero",
				}
			}
		}

		// Return success result with the arithmetic operation
		return &ResultExpression{
			Success: true,
			Value:     binExpr,
			ErrorType: "",
		}
	}

	return expr
}

func (b *Builder) isArithmeticOperation(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" || op == "%"
}

func (b *Builder) buildUnaryExpr(ctx parser.IUnaryExprContext) Expression {
	pipeExpr := b.buildPipeExpr(ctx.PipeExpr())

	// Handle unary operators
	if ctx.PLUS() != nil {
		return &UnaryExpression{
			Operator: "+",
			Operand:  pipeExpr,
		}
	} else if ctx.MINUS() != nil {
		return &UnaryExpression{
			Operator: "-",
			Operand:  pipeExpr,
		}
	} else if ctx.NOT_OP() != nil {
		return &UnaryExpression{
			Operator: "!",
			Operand:  pipeExpr,
		}
	} else if ctx.AWAIT() != nil {
		// Handle await as unary operator: await expr
		return &AwaitExpression{
			Expression: pipeExpr,
		}
	}

	return pipeExpr
}

func (b *Builder) buildPipeExpr(ctx parser.IPipeExprContext) Expression {
	callExprs := ctx.AllCallExpr()
	if len(callExprs) == 1 {
		return b.buildCallExpr(callExprs[0])
	}

	// Build pipe chain
	result := b.buildCallExpr(callExprs[0])

	for i := 1; i < len(callExprs); i++ {
		right := b.buildCallExpr(callExprs[i])

		// Handle both CallExpression and Identifier cases
		switch rightExpr := right.(type) {
		case *CallExpression:
			// Create a call where the left expression is the first argument
			args := make([]Expression, 0, len(rightExpr.Arguments)+1)
			args = append(args, result)
			args = append(args, rightExpr.Arguments...)

			result = &CallExpression{
				Function:       rightExpr.Function,
				Arguments:      args,
				NamedArguments: nil,
			}
		case *Identifier:
			// Convert identifier to call expression with piped argument
			result = &CallExpression{
				Function:       rightExpr,
				Arguments:      []Expression{result},
				NamedArguments: nil,
			}
		}
	}

	return result
}
