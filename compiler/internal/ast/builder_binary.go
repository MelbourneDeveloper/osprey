package ast

import (
	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildBinaryExpr(ctx parser.IBinaryExprContext) Expression {
	return b.buildTernaryExpr(ctx.TernaryExpr())
}

func (b *Builder) buildTernaryExpr(ctx parser.ITernaryExprContext) Expression {
	// Check if this is a type pattern ternary ({ type: pattern } expr ? then : else)
	if ctx.GetPat() != nil {
		conditionExpr := b.buildComparisonExpr(ctx.GetCond())
		thenExpr := b.buildTernaryExpr(ctx.GetThenExpr())
		elseExpr := b.buildTernaryExpr(ctx.GetElseExpr())
		
		// For structural patterns like v { value } ? value : "default"
		// Extract fields from the pattern and replace variable references with field access
		fieldPattern := b.buildFieldPattern("*", ctx.GetPat())
		
		if len(fieldPattern.Fields) == 1 {
			fieldName := fieldPattern.Fields[0]
			
			// If then expression is just the field name, return field access directly
			// Note: Result types will be caught at codegen and forced to use pattern matching
			if identExpr, ok := thenExpr.(*Identifier); ok && identExpr.Name == fieldName {
				return &FieldAccessExpression{
					Object:    conditionExpr,
					FieldName: fieldName,
					Position:  b.getPositionFromContext(ctx),
				}
			}
		}
		
		// Convert type pattern to structural pattern
		pattern := Pattern{
			Constructor: "",
			IsWildcard:  false,
			Fields:      fieldPattern.Fields,
		}
		
		thenArm := MatchArm{
			Pattern:    pattern,
			Expression: thenExpr,
		}

		elseArm := MatchArm{
			Pattern:    Pattern{Constructor: "_", IsWildcard: true},
			Expression: elseExpr,
		}

		return &MatchExpression{
			Expression: conditionExpr,
			Arms:       []MatchArm{thenArm, elseArm},
		}
	}

	// Check if this is a pattern-based ternary (expr { pattern } ? then : else)
	if ctx.LBRACE() != nil {
		conditionExpr := b.buildComparisonExpr(ctx.GetCond())
		thenExpr := b.buildTernaryExpr(ctx.GetThenExpr())
		elseExpr := b.buildTernaryExpr(ctx.GetElseExpr())
		
		// For structural patterns like v { value } ? value : "default"
		// Simplify to direct field access since pattern implies field exists
		fieldPattern := b.buildFieldPattern("*", ctx.GetPat())
		if len(fieldPattern.Fields) == 1 {
			fieldName := fieldPattern.Fields[0]
			
			// If then expression is just the field name, return field access directly
			// Note: Result types will be caught at codegen and forced to use pattern matching
			if identExpr, ok := thenExpr.(*Identifier); ok && identExpr.Name == fieldName {
				return &FieldAccessExpression{
					Object:    conditionExpr,
					FieldName: fieldName,
					Position:  b.getPositionFromContext(ctx),
				}
			}
		}
		
		// Fallback: create match expression (this shouldn't be reached for simple cases)
		pattern := b.buildFieldPattern("*", ctx.GetPat())
		
		thenArm := MatchArm{
			Pattern:    pattern,
			Expression: thenExpr,
		}

		elseArm := MatchArm{
			Pattern:    Pattern{Constructor: "_", IsWildcard: true},
			Expression: elseExpr,
		}

		return &MatchExpression{
			Expression: conditionExpr,
			Arms:       []MatchArm{thenArm, elseArm},
		}
	}
	
	// Check if this is a simple boolean ternary (expr ? then : else)
	if ctx.QUESTION() != nil {
		// First comparison expression is the condition 
		conditionExpr := b.buildComparisonExpr(ctx.ComparisonExpr())
		thenExpr := b.buildTernaryExpr(ctx.GetThenExpr())
		elseExpr := b.buildTernaryExpr(ctx.GetElseExpr())
		
		trueArm := MatchArm{
			Pattern:    Pattern{Constructor: "true"},
			Expression: thenExpr,
		}

		falseArm := MatchArm{
			Pattern:    Pattern{Constructor: "false"},
			Expression: elseExpr,
		}

		return &MatchExpression{
			Expression: conditionExpr,
			Arms:       []MatchArm{trueArm, falseArm},
		}
	}
	
	// Otherwise it's just a comparison expression
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
			Position: b.getPositionFromContext(ctx),
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
			Position: b.getPositionFromContext(ctx),
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

	// Get all operators in order they appear
	stars := ctx.AllSTAR()
	slashes := ctx.AllSLASH()
	mods := ctx.AllMOD_OP()
	
	// Determine operators by their positions in the input text
	operators := make([]string, len(unaryExprs)-1)
	starIdx, slashIdx, modIdx := 0, 0, 0
	
	for i := 0; i < len(operators); i++ {
		starPos := -1
		slashPos := -1
		modPos := -1
		
		if starIdx < len(stars) {
			starPos = stars[starIdx].GetSymbol().GetTokenIndex()
		}
		if slashIdx < len(slashes) {
			slashPos = slashes[slashIdx].GetSymbol().GetTokenIndex()
		}
		if modIdx < len(mods) {
			modPos = mods[modIdx].GetSymbol().GetTokenIndex()
		}
		
		// Find the operator with the smallest position (appears first)
		operators[i] = "*"
		earliestPos := starPos
		
		if slashPos != -1 && (earliestPos == -1 || slashPos < earliestPos) {
			earliestPos = slashPos
			operators[i] = "/"
		}
		if modPos != -1 && (earliestPos == -1 || modPos < earliestPos) {
			operators[i] = "%"
		}
		
		// Advance the index for the operator we just used
		switch operators[i] {
		case "*":
			starIdx++
		case "/":
			slashIdx++
		case "%":
			modIdx++
		}
	}

	for i := 1; i < len(unaryExprs); i++ {
		right := b.buildUnaryExpr(unaryExprs[i])

		// Wrap arithmetic operations in result types
		binExpr := &BinaryExpression{
			Left:     left,
			Operator: operators[i-1],
			Right:    right,
			Position: b.getPositionFromContext(ctx),
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
					Success:   false,
					Value:     &StringLiteral{Value: "Division by zero"},
					ErrorType: "DivisionByZero",
				}
			}
		}

		// Return success result with the arithmetic operation
		return &ResultExpression{
			Success:   true,
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
			Position: b.getPositionFromContext(ctx),
		}
	} else if ctx.MINUS() != nil {
		return &UnaryExpression{
			Operator: "-",
			Operand:  pipeExpr,
			Position: b.getPositionFromContext(ctx),
		}
	} else if ctx.NOT_OP() != nil {
		return &UnaryExpression{
			Operator: "!",
			Operand:  pipeExpr,
			Position: b.getPositionFromContext(ctx),
		}
	} else if ctx.AWAIT() != nil {
		// Handle await as unary operator: await expr
		return &AwaitExpression{
			Expression: pipeExpr,
			Position:   b.getPositionFromContext(ctx),
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
