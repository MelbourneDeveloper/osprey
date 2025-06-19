package ast

import (
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

func (b *Builder) buildStatement(ctx parser.IStatementContext) Statement {
	switch {
	case ctx.ImportStmt() != nil:
		return b.buildImport(ctx.ImportStmt())
	case ctx.LetDecl() != nil:
		return b.buildLetDecl(ctx.LetDecl())
	case ctx.FnDecl() != nil:
		return b.buildFnDecl(ctx.FnDecl())
	case ctx.ExternDecl() != nil:
		return b.buildExternDecl(ctx.ExternDecl())
	case ctx.TypeDecl() != nil:
		return b.buildTypeDecl(ctx.TypeDecl())
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
	// Check for loop expression first
	if ctx.LoopExpr() != nil {
		return b.buildLoopExpr(ctx.LoopExpr())
	}

	return b.buildMatchExpr(ctx.MatchExpr())
}

// buildLoopExpr builds a LoopExpression from the parse tree.
func (b *Builder) buildLoopExpr(ctx parser.ILoopExprContext) *LoopExpression {
	var condition *Expression

	// Check if it's a conditional loop
	if ctx.Expr() != nil {
		expr := b.buildExpression(ctx.Expr())
		condition = &expr
	}

	// Build the loop body
	blockBody := b.buildBlockBody(ctx.BlockBody())
	body := BlockExpression{
		Statements: blockBody.Statements,
		Expression: blockBody.Expression,
	}

	return &LoopExpression{
		Condition: condition,
		Body:      body,
	}
}
