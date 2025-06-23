package ast

import (
	"strconv"
	"strings"

	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildPrimary(ctx parser.IPrimaryContext) Expression {
	if ctx == nil {
		return nil
	}

	// Handle fiber-related expressions
	if fiberExpr := b.buildFiberExpression(ctx); fiberExpr != nil {
		return fiberExpr
	}

	// Handle other expressions
	switch {
	case ctx.BlockExpr() != nil:
		return b.buildBlockExpression(ctx.BlockExpr())
	case ctx.Literal() != nil:
		return b.buildLiteral(ctx.Literal())
	case ctx.LambdaExpr() != nil:
		return b.buildLambdaExpr(ctx.LambdaExpr())
	case ctx.ID() != nil && ctx.LSQUARE() != nil && ctx.INT() != nil && ctx.RSQUARE() != nil:
		// Array indexing: ID[INT]
		return b.buildListAccess(ctx)
	case ctx.ID() != nil:
		return &Identifier{Name: ctx.ID().GetText()}
	case ctx.Expr(0) != nil:
		return b.buildExpression(ctx.Expr(0))
	}

	return nil
}

// buildFiberExpression handles fiber-related expressions.
func (b *Builder) buildFiberExpression(ctx parser.IPrimaryContext) Expression {
	switch {
	case ctx.SPAWN() != nil:
		return b.buildSpawnExpression(ctx)
	case ctx.AWAIT() != nil:
		return b.buildAwaitExpression(ctx)
	case ctx.SEND() != nil:
		return b.buildSendExpression(ctx)
	case ctx.RECV() != nil:
		return b.buildRecvExpression(ctx)
	case ctx.YIELD() != nil:
		return b.buildYieldExpression(ctx)
	case ctx.SELECT() != nil:
		return b.buildSelectExpression(ctx.SelectExpr().(*parser.SelectExprContext))
	case ctx.TypeConstructor() != nil:
		return b.buildTypeConstructor(ctx.TypeConstructor().(*parser.TypeConstructorContext))
	}

	return nil
}

func (b *Builder) buildLiteral(ctx parser.ILiteralContext) Expression {
	switch {
	case ctx.INT() != nil:
		text := ctx.INT().GetText()
		value, _ := strconv.ParseInt(text, 10, 64)

		return &IntegerLiteral{Value: value}
	case ctx.STRING() != nil:
		text := ctx.STRING().GetText()
		// Remove quotes and process escape sequences
		value := strings.Trim(text, "\"")
		value = b.processEscapeSequences(value)

		return &StringLiteral{Value: value}
	case ctx.INTERPOLATED_STRING() != nil:
		return b.buildInterpolatedString(ctx.INTERPOLATED_STRING().GetText())
	case ctx.TRUE() != nil:
		return &BooleanLiteral{Value: true}
	case ctx.FALSE() != nil:
		return &BooleanLiteral{Value: false}
	case ctx.ListLiteral() != nil:
		return b.buildListLiteral(ctx.ListLiteral())
	}

	return nil
}

// buildInterpolatedString parses an interpolated string like "Hello ${name}!".
func (b *Builder) buildInterpolatedString(text string) Expression {
	text = strings.Trim(text, "\"")
	parts := b.parseInterpolatedParts(text)

	// Process escape sequences in text parts
	for i := range parts {
		if !parts[i].IsExpression {
			parts[i].Text = b.processEscapeSequences(parts[i].Text)
		}
	}

	return &InterpolatedStringLiteral{Parts: parts}
}

// buildListLiteral builds a ListLiteral from a list literal context.
func (b *Builder) buildListLiteral(ctx parser.IListLiteralContext) Expression {
	if ctx == nil {
		return nil
	}

	elements := make([]Expression, 0)

	// Build each expression in the list
	for _, exprCtx := range ctx.AllExpr() {
		element := b.buildExpression(exprCtx)
		if element != nil {
			elements = append(elements, element)
		}
	}

	return &ListLiteral{Elements: elements}
}

// buildListAccess builds a ListAccessExpression from array indexing syntax.
func (b *Builder) buildListAccess(ctx parser.IPrimaryContext) Expression {
	if ctx == nil || ctx.ID() == nil || ctx.INT() == nil {
		return nil
	}

	// Get the list identifier
	listExpr := &Identifier{Name: ctx.ID().GetText()}

	// Parse the index
	indexText := ctx.INT().GetText()
	indexValue, _ := strconv.ParseInt(indexText, 10, 64)
	indexExpr := &IntegerLiteral{Value: indexValue}

	return &ListAccessExpression{
		List:  listExpr,
		Index: indexExpr,
	}
}

// buildLambdaExpr builds a LambdaExpression from a lambda context.
func (b *Builder) buildLambdaExpr(ctx parser.ILambdaExprContext) Expression {
	if ctx == nil {
		return nil
	}

	params := make([]Parameter, 0)

	// Handle parameter list if present
	if ctx.ParamList() != nil {
		for _, paramCtx := range ctx.ParamList().AllParam() {
			param := Parameter{
				Name: paramCtx.ID().GetText(),
				Type: nil, // Parse type annotation if present
			}

			// Parse parameter type annotation if present
			if paramCtx.Type_() != nil {
				param.Type = b.buildTypeExpression(paramCtx.Type_())
			}

			params = append(params, param)
		}
	}

	// Parse return type annotation if present
	var returnType *TypeExpression
	if ctx.Type_() != nil {
		returnType = b.buildTypeExpression(ctx.Type_())
	}

	// Parse the lambda body
	body := b.buildExpression(ctx.Expr())

	return &LambdaExpression{
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}
}

// processEscapeSequences processes common escape sequences in string literals.
func (b *Builder) processEscapeSequences(input string) string {
	result := strings.ReplaceAll(input, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\t", "\t")
	result = strings.ReplaceAll(result, "\\r", "\r")
	result = strings.ReplaceAll(result, "\\\\", "\\")
	result = strings.ReplaceAll(result, "\\\"", "\"")

	return result
}

// buildTypeConstructor builds type constructor expressions like Fiber<Int> { computation: fn() => 42 }.
func (b *Builder) buildTypeConstructor(ctx *parser.TypeConstructorContext) Expression {
	typeName := ctx.ID().GetText()

	// Build field assignments
	fieldAssignments := make(map[string]Expression)
	if ctx.FieldAssignments() != nil {
		for _, fieldCtx := range ctx.FieldAssignments().AllFieldAssignment() {
			fieldName := fieldCtx.ID().GetText()
			fieldValue := b.buildExpression(fieldCtx.Expr())
			fieldAssignments[fieldName] = fieldValue
		}
	}

	// Handle specific fiber types
	switch typeName {
	case "Fiber":
		if computation, exists := fieldAssignments["computation"]; exists {
			return &SpawnExpression{
				Expression: computation,
			}
		}
	case "Channel":
		if capacity, exists := fieldAssignments["capacity"]; exists {
			return &ChannelCreateExpression{
				Capacity: capacity,
			}
		}
	}

	// For other types, return a generic type constructor
	return &TypeConstructorExpression{
		TypeName: typeName,
		Fields:   fieldAssignments,
	}
}

// buildBlockExpression builds a BlockExpression from a block context.
func (b *Builder) buildBlockExpression(ctx parser.IBlockExprContext) Expression {
	if ctx == nil {
		return nil
	}

	blockBody := ctx.BlockBody()
	if blockBody == nil {
		return &BlockExpression{
			Statements: []Statement{},
			Expression: nil,
		}
	}

	// Build all statements in the block
	statements := make([]Statement, 0)
	for _, stmtCtx := range blockBody.AllStatement() {
		stmt := b.buildStatement(stmtCtx)
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	// Build the final expression if present
	var finalExpr Expression
	if blockBody.Expr() != nil {
		finalExpr = b.buildExpression(blockBody.Expr())
	}

	return &BlockExpression{
		Statements: statements,
		Expression: finalExpr,
	}
}
