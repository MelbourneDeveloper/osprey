package ast

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/christianfindlay/osprey/parser"
)

// parseInterpolatedParts parses the parts of an interpolated string.
func (b *Builder) parseInterpolatedParts(text string) []InterpolatedPart {
	var parts []InterpolatedPart

	i := 0

	for i < len(text) {
		nextPart, newIndex := b.parseNextInterpolatedPart(text, i)
		if nextPart != nil {
			parts = append(parts, *nextPart)
		}

		if newIndex <= i {
			break
		}

		i = newIndex
	}

	return parts
}

// parseNextInterpolatedPart parses the next part of an interpolated string.
func (b *Builder) parseNextInterpolatedPart(text string, startIndex int) (*InterpolatedPart, int) {
	start := strings.Index(text[startIndex:], "${")
	if start == -1 {
		return b.createTextPart(text, startIndex), len(text)
	}

	if start > 0 {
		return b.createTextPartWithRange(text, startIndex, startIndex+start), startIndex + start
	}

	return b.parseExpressionPart(text, startIndex)
}

// createTextPart creates a text part from the remaining text.
func (b *Builder) createTextPart(text string, startIndex int) *InterpolatedPart {
	if startIndex >= len(text) {
		return nil
	}

	return &InterpolatedPart{
		IsExpression: false,
		Text:         text[startIndex:],
		Expression:   nil,
	}
}

// createTextPartWithRange creates a text part with a specific range.
func (b *Builder) createTextPartWithRange(text string, start, end int) *InterpolatedPart {
	return &InterpolatedPart{
		IsExpression: false,
		Text:         text[start:end],
		Expression:   nil,
	}
}

// parseExpressionPart parses an expression part ${...}.
func (b *Builder) parseExpressionPart(text string, startIndex int) (*InterpolatedPart, int) {
	exprStart := startIndex + InterpolationOffset
	depth := 1
	j := exprStart

	for j < len(text) && depth > 0 {
		switch text[j] {
		case '{':
			depth++
		case '}':
			depth--
		}

		j++
	}

	if depth == 0 {
		exprText := text[exprStart : j-1]
		expr := b.parseExpressionFromString(exprText)

		return &InterpolatedPart{
			IsExpression: true,
			Text:         "",
			Expression:   expr,
		}, j
	}

	return &InterpolatedPart{
		IsExpression: false,
		Text:         text[startIndex:],
	}, len(text)
}

// parseExpressionFromString parses a string into an expression using the actual ANTLR grammar.
func (b *Builder) parseExpressionFromString(exprText string) Expression {
	input := antlr.NewInputStream(exprText)
	lexer := parser.NewospreyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewospreyParser(stream)

	// Parse expression
	tree := p.Expr()

	return b.buildExpression(tree)
}
