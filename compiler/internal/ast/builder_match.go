package ast

import (
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

// buildPattern converts a pattern parse tree into the corresponding AST node.
//
// The function necessarily contains a lot of branching because it needs to
// inspect many possible grammar productions to decide which concrete pattern
// to build. Refactoring this would introduce significant duplication without
// improving readability. Therefore we opt to silence the gocognit linter for
// this specific case.
//
//nolint:gocognit
func (b *Builder) buildPattern(ctx parser.IPatternContext) Pattern {
	// Handle unary expressions (like negative numbers)
	if ctx.UnaryExpr() != nil {
		unaryCtx := ctx.UnaryExpr()

		// Check if it's a negative number
		if unaryCtx.MINUS() != nil && unaryCtx.PipeExpr() != nil {
			pipeCtx := unaryCtx.PipeExpr()
			if callCtx := pipeCtx.CallExpr(0); callCtx != nil {
				if primaryCtx := callCtx.Primary(); primaryCtx != nil {
					if literalCtx := primaryCtx.Literal(); literalCtx != nil {
						if literalCtx.INT() != nil {
							// This is a negative number like -1
							return Pattern{
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
		}

		// Handle positive numbers and other unary expressions
		if pipeCtx := unaryCtx.PipeExpr(); pipeCtx != nil {
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
						if literalCtx.INT() != nil {
							return Pattern{
								Constructor: literalCtx.INT().GetText(),
								Variable:    "",
								Fields:      nil,
								Nested:      nil,
								IsWildcard:  false,
							}
						} else if literalCtx.STRING() != nil {
							return Pattern{
								Constructor: literalCtx.STRING().GetText(),
								Variable:    "",
								Fields:      nil,
								Nested:      nil,
								IsWildcard:  false,
							}
						}
					}
				}
			}
		}
	}

	// Handle wildcard pattern
	if ctx.UNDERSCORE() != nil {
		return Pattern{
			Constructor: "_",
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  true,
		}
	}

	// Handle identifier patterns
	ids := ctx.AllID()
	if len(ids) == OneIdentifier {
		// Single identifier is a constructor pattern (like: Red => ...)
		return Pattern{
			Constructor: ids[0].GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	} else if len(ids) == TwoIdentifiers {
		// Two identifiers is a constructor with variable binding (like: Ok x => ...)
		return Pattern{
			Constructor: ids[0].GetText(),
			Variable:    ids[1].GetText(),
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}

	return Pattern{
		Constructor: "unknown",
		Variable:    "",
		Fields:      nil,
		Nested:      nil,
		IsWildcard:  false,
	}
}
