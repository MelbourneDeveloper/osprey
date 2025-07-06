package ast

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/christianfindlay/osprey/parser"
)

const (
	// OneIdentifier represents a single identifier in a pattern.
	OneIdentifier = 1
	// TwoIdentifiers represents the number of identifiers in a two-part pattern.
	TwoIdentifiers = 2
	// InterpolationOffset represents the offset to skip ${ in string interpolation.
	InterpolationOffset = 2 // Skip ${ - offset of 2
)

// Builder builds an AST from the ANTLR parse tree.
type Builder struct {
	parser.BaseospreyListener
}

// NewBuilder creates a new AST builder instance.
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildProgram builds an AST from a parse tree.
func (b *Builder) BuildProgram(tree parser.IProgramContext) *Program {
	statements := make([]Statement, 0)

	for _, stmtCtx := range tree.AllStatement() {
		if stmt := b.buildStatement(stmtCtx); stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return &Program{Statements: statements}
}

// getPosition extracts position information from an ANTLR token.
func (b *Builder) getPosition(token antlr.Token) *Position {
	if token == nil {
		return &Position{Line: 1, Column: 0}
	}
	return &Position{
		Line:   token.GetLine(),
		Column: token.GetColumn(),
	}
}

// getPositionFromContext extracts position information from a parser context.
func (b *Builder) getPositionFromContext(ctx antlr.ParserRuleContext) *Position {
	if ctx == nil {
		return &Position{Line: 1, Column: 0}
	}
	start := ctx.GetStart()
	if start == nil {
		return &Position{Line: 1, Column: 0}
	}
	return &Position{
		Line:   start.GetLine(),
		Column: start.GetColumn(),
	}
}

func (b *Builder) buildStatement(ctx parser.IStatementContext) Statement {
	switch {
	case ctx.ImportStmt() != nil:
		return b.buildImport(ctx.ImportStmt())
	case ctx.LetDecl() != nil:
		return b.buildLetDecl(ctx.LetDecl())
	case ctx.AssignStmt() != nil:
		return b.buildAssignStmt(ctx.AssignStmt())
	case ctx.FnDecl() != nil:
		return b.buildFnDecl(ctx.FnDecl())
	case ctx.ExternDecl() != nil:
		return b.buildExternDecl(ctx.ExternDecl())
	case ctx.TypeDecl() != nil:
		return b.buildTypeDecl(ctx.TypeDecl())
	case ctx.EffectDecl() != nil:
		return b.buildEffectDecl(ctx.EffectDecl())
	case ctx.ModuleDecl() != nil:
		return b.buildModuleDecl(ctx.ModuleDecl())
	case ctx.ExprStmt() != nil:
		return b.buildExprStmt(ctx.ExprStmt())
	}

	return nil
}

func (b *Builder) buildExprStmt(ctx parser.IExprStmtContext) *ExpressionStatement {
	expr := b.buildExpression(ctx.Expr())

	return &ExpressionStatement{Expression: expr}
}

func (b *Builder) buildExpression(ctx parser.IExprContext) Expression {
	return b.buildMatchExpr(ctx.MatchExpr())
}
