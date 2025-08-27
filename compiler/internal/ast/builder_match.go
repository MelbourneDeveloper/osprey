package ast

import (
	"github.com/christianfindlay/osprey/parser"
)

const (
	unknownPatternConstructor = "unknown"
)

func (b *Builder) buildMatchExpr(ctx parser.IMatchExprContext) Expression {
	if ctx.MATCH() != nil {
		// This is a match expression
		expr := b.buildBinaryExpr(ctx.BinaryExpr())

		arms := make([]MatchArm, 0)

		for _, armCtx := range ctx.AllMatchArm() {
			arm := b.buildMatchArm(armCtx)
			arms = append(arms, arm)
		}

		return &MatchExpression{
			Expression: expr,
			Arms:       arms,
			Position:   b.getPositionFromContext(ctx),
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
func (b *Builder) buildPattern(ctx parser.IPatternContext) Pattern {
	// Handle unary expressions (like negative numbers)
	if ctx.UnaryExpr() != nil {
		return b.buildUnaryExprPattern(ctx.UnaryExpr())
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
	return b.buildIdentifierPattern(ctx)
}

// buildUnaryExprPattern handles unary expression patterns like negative numbers.
func (b *Builder) buildUnaryExprPattern(unaryCtx parser.IUnaryExprContext) Pattern {
	// Check if it's a negative number
	if unaryCtx.MINUS() != nil && unaryCtx.PipeExpr() != nil {
		if pattern := b.tryBuildNegativeNumberPattern(unaryCtx.PipeExpr()); pattern.Constructor != unknownPatternConstructor {
			return pattern
		}
	}

	// Handle positive numbers and other unary expressions
	if unaryCtx.PipeExpr() != nil {
		return b.buildPipeExprPattern(unaryCtx.PipeExpr())
	}

	return Pattern{
		Constructor: unknownPatternConstructor,
		Variable:    "",
		Fields:      nil,
		Nested:      nil,
		IsWildcard:  false,
	}
}

// tryBuildNegativeNumberPattern attempts to build a negative number pattern.
func (b *Builder) tryBuildNegativeNumberPattern(pipeCtx parser.IPipeExprContext) Pattern {
	if callCtx := pipeCtx.CallExpr(0); callCtx != nil {
		if primaryCtx := callCtx.Primary(); primaryCtx != nil {
			if literalCtx := primaryCtx.Literal(); literalCtx != nil {
				if literalCtx.INT() != nil {
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

	return Pattern{Constructor: unknownPatternConstructor}
}

// buildPipeExprPattern builds patterns from pipe expressions.
func (b *Builder) buildPipeExprPattern(pipeCtx parser.IPipeExprContext) Pattern {
	if callCtx := pipeCtx.CallExpr(0); callCtx != nil {
		if primaryCtx := callCtx.Primary(); primaryCtx != nil {
			// Check for simple identifier (like Red, Green, Blue)
			if primaryCtx.ID(0) != nil {
				return Pattern{
					Constructor: primaryCtx.ID(0).GetText(),
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

	return Pattern{Constructor: unknownPatternConstructor}
}

// buildLiteralPattern builds patterns from literal contexts.
func (b *Builder) buildLiteralPattern(literalCtx parser.ILiteralContext) Pattern {
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
	} else if literalCtx.TRUE() != nil {
		return Pattern{
			Constructor: literalCtx.TRUE().GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	} else if literalCtx.FALSE() != nil {
		return Pattern{
			Constructor: literalCtx.FALSE().GetText(),
			Variable:    "",
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}

	return Pattern{Constructor: unknownPatternConstructor}
}

// buildIdentifierPattern handles identifier patterns including field patterns.
func (b *Builder) buildIdentifierPattern(ctx parser.IPatternContext) Pattern {
	ids := ctx.AllID()

	// Handle constructor with field pattern (like: Success { value } => ...)
	if len(ids) == OneIdentifier && ctx.LBRACE() != nil && ctx.FieldPattern() != nil {
		return b.buildFieldPattern(ids[0].GetText(), ctx.FieldPattern())
	}

	// Handle type annotation pattern (like: p: PersonData => ...)  
	if len(ids) == OneIdentifier && ctx.COLON() != nil && ctx.Type_() != nil {
		typeName := b.buildTypeFromContext(ctx.Type_())
		return Pattern{
			Constructor: typeName,
			Variable:    ids[0].GetText(),
			Fields:      nil,
			Nested:      nil,
			IsWildcard:  false,
		}
	}

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
		Constructor: unknownPatternConstructor,
		Variable:    "",
		Fields:      nil,
		Nested:      nil,
		IsWildcard:  false,
	}
}

// buildFieldPattern builds a pattern with field bindings.
func (b *Builder) buildFieldPattern(constructor string, fieldPattern parser.IFieldPatternContext) Pattern {
	fields := make([]string, 0)

	// Extract all field names from the field pattern
	for _, fieldID := range fieldPattern.AllID() {
		fields = append(fields, fieldID.GetText())
	}

	return Pattern{
		Constructor: constructor,
		Variable:    "",
		Fields:      fields,
		Nested:      nil,
		IsWildcard:  false,
	}
}

// buildTypeFromContext extracts type name from a type context
func (b *Builder) buildTypeFromContext(ctx parser.ITypeContext) string {
	if ctx == nil {
		return unknownPatternConstructor
	}
	// For simple type ID patterns like "PersonData", just return the ID text
	if ctx.ID() != nil {
		return ctx.ID().GetText()
	}
	return unknownPatternConstructor
}
