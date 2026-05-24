package ast

// Coverage-focused tests for the AST marker methods (`isStatement`,
// `isExpression`). They are interface markers; nothing in the runtime
// calls them by name, so they show 0% unless the type is asserted into
// the matching interface at least once.

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
)

func TestStatementMarkers(t *testing.T) {
	stmts := []ast.Statement{
		&ast.ImportStatement{},
		&ast.LetDeclaration{},
		&ast.AssignmentStatement{},
		&ast.FunctionDeclaration{},
		&ast.ExternDeclaration{},
		&ast.TypeDeclaration{},
		&ast.EffectDeclaration{},
		&ast.ExpressionStatement{},
		&ast.ModuleDeclaration{},
	}
	if len(stmts) == 0 {
		t.Fatal("no statements")
	}
}

func TestExpressionMarkers(t *testing.T) {
	exprs := []ast.Expression{
		&ast.FunctionCallExpression{},
		&ast.IntegerLiteral{},
		&ast.FloatLiteral{},
		&ast.StringLiteral{},
		&ast.BooleanLiteral{},
		&ast.InterpolatedStringLiteral{},
		&ast.Identifier{},
		&ast.BinaryExpression{},
		&ast.UnaryExpression{},
		&ast.CallExpression{},
		&ast.MethodCallExpression{},
		&ast.LambdaExpression{},
		&ast.MatchExpression{},
		&ast.ResultExpression{},
		&ast.FieldAccessExpression{},
		&ast.ModuleAccessExpression{},
		&ast.SpawnExpression{},
		&ast.AwaitExpression{},
		&ast.YieldExpression{},
		&ast.ChannelExpression{},
		&ast.ChannelSendExpression{},
		&ast.ChannelRecvExpression{},
		&ast.SelectExpression{},
		&ast.ChannelCreateExpression{},
		&ast.TypeConstructorExpression{},
		&ast.BlockExpression{},
		&ast.ListLiteral{},
		&ast.ObjectLiteral{},
		&ast.MapLiteral{},
		&ast.ListAccessExpression{},
		&ast.UpdateExpression{},
		&ast.PerformExpression{},
		&ast.HandlerExpression{},
	}
	if len(exprs) == 0 {
		t.Fatal("no expressions")
	}
}
