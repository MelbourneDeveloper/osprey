package ast

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildMatchExpr(ctx parser.IMatchExprContext) Expression {
	if ctx.MATCH() != nil {
		// This is a match expression
		expr := b.buildExpression(ctx.Expr())

		arms := make([]MatchArm, 0)

		for _, armCtx := range ctx.AllMatchArm() {
			arm := b.buildMatchArm(armCtx)
			arms = append(arms, arm)
		}

		return &MatchExpression{
			Expression: expr,
			Arms:       arms,
		}
	}

	if ctx.SelectExpr() != nil {
		// This is a select expression
		return b.buildSelectExpression(ctx.SelectExpr().(*parser.SelectExprContext))
	}

	// Otherwise it's a binary expression
	return b.buildBinaryExpr(ctx.BinaryExpr())
}

func (b *Builder) buildMatchArm(ctx parser.IMatchArmContext) MatchArm {
	pattern := b.buildPattern(ctx.Pattern())
	expr := b.buildExpression(ctx.Expr())

	return MatchArm{
		Pattern:    pattern,
		Expression: expr,
	}
}

func (b *Builder) buildPattern(ctx parser.IPatternContext) Pattern {
	// Handle unary expressions (like negative numbers)
	if ctx.UnaryExpr() != nil {
		return b.buildUnaryPattern(ctx.UnaryExpr())
	}

	// Handle wildcard pattern
	if ctx.UNDERSCORE() != nil {
		return b.buildWildcardPattern()
	}

	// Handle identifier patterns
	return b.buildIdentifierPattern(ctx.AllID())
}

// buildUnaryPattern handles unary expression patterns
func (b *Builder) buildUnaryPattern(unaryCtx parser.IUnaryExprContext) Pattern {
	// Check if it's a negative number
	if unaryCtx.MINUS() != nil && unaryCtx.PipeExpr() != nil {
		if pattern := b.buildNegativeNumberPattern(unaryCtx.PipeExpr()); pattern != nil {
			return *pattern
		}
	}

	// Handle positive numbers and other unary expressions
	if unaryCtx.PipeExpr() != nil {
		return b.buildPipeExprPattern(unaryCtx.PipeExpr())
	}

	return b.buildUnknownPattern()
}

// buildNegativeNumberPattern handles negative number patterns like -1
func (b *Builder) buildNegativeNumberPattern(pipeCtx parser.IPipeExprContext) *Pattern {
	if callCtx := pipeCtx.CallExpr(0); callCtx != nil {
		if primaryCtx := callCtx.Primary(); primaryCtx != nil {
			if literalCtx := primaryCtx.Literal(); literalCtx != nil {
				if literalCtx.INT() != nil {
					return &Pattern{
						Constructor: "-" + literalCtx.INT().GetText(),
						Variable:    "",
						Fields:      nil,
						Nested:      nil,
						IsWildcard:  false,
					}
				}
			}
		}
	}
	return nil
}

// buildPipeExprPattern handles pipe expression patterns
func (b *Builder) buildPipeExprPattern(pipeCtx parser.IPipeExprContext) Pattern {
	if callCtx := pipeCtx.CallExpr(0); callCtx != nil {
		if primaryCtx := callCtx.Primary(); primaryCtx != nil {
			// Check for simple identifier (like Red, Green, Blue)
			if primaryCtx.ID() != nil {
				return Pattern{
					Constructor: primaryCtx.ID().GetText(),
					Variable:    "",
					Fields:      nil,
					Nested:      nil,
					IsWildcard:  false,
				}
			}

			// Check for literals
			if literalCtx := primaryCtx.Literal(); literalCtx != nil {
				return b.buildLiteralPattern(literalCtx)
			}
		}
	}
	return b.buildUnknownPattern()
}

// buildLiteralPattern handles literal patterns
func (b *Builder) buildLiteralPattern(literalCtx parser.ILiteralContext) Pattern {
	if literalCtx.INT() != nil {
		return Pattern{
			Constructor: literalCtx.INT().GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}
	if literalCtx.STRING() != nil {
		return Pattern{
			Constructor: literalCtx.STRING().GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}
	return b.buildUnknownPattern()
}

// buildWildcardPattern builds a wildcard pattern
func (b *Builder) buildWildcardPattern() Pattern {
	return Pattern{
		Constructor: "_",
		Variable:    "",
		Fields:      nil,
		Nested:      nil,
		IsWildcard:  true,
	}
}

// buildIdentifierPattern handles identifier patterns
func (b *Builder) buildIdentifierPattern(ids []antlr.TerminalNode) Pattern {
	if len(ids) == OneIdentifier {
		// Single identifier is a constructor pattern (like: Red => ...)
		return Pattern{
			Constructor: ids[0].GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}
	if len(ids) == TwoIdentifiers {
		// Two identifiers is a constructor with variable binding (like: Ok x => ...)
		return Pattern{
			Constructor: ids[0].GetText(),
			Variable:    ids[1].GetText(),
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}
	return b.buildUnknownPattern()
}

// buildUnknownPattern builds a fallback unknown pattern
func (b *Builder) buildUnknownPattern() Pattern {
	return Pattern{
		Constructor: "unknown",
		Variable:    "",
		Fields:      nil,
		Nested:      nil,
		IsWildcard:  false,
	}
}
