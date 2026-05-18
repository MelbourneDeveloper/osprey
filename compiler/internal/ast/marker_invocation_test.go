package ast

// Invokes every isStatement() / isExpression() marker so coverage records
// them. The methods have empty bodies — type-asserting alone doesn't
// execute them. They must be called by name from inside this package.

import "testing"

func TestStatementMarkerInvocations(t *testing.T) {
	(&ImportStatement{}).isStatement()
	(&LetDeclaration{}).isStatement()
	(&AssignmentStatement{}).isStatement()
	(&FunctionDeclaration{}).isStatement()
	(&ExternDeclaration{}).isStatement()
	(&TypeDeclaration{}).isStatement()
	(&EffectDeclaration{}).isStatement()
	(&ExpressionStatement{}).isStatement()
	(&ModuleDeclaration{}).isStatement()
}

func TestExpressionMarkerInvocations(t *testing.T) {
	(&FunctionCallExpression{}).isExpression()
	(&IntegerLiteral{}).isExpression()
	(&FloatLiteral{}).isExpression()
	(&StringLiteral{}).isExpression()
	(&BooleanLiteral{}).isExpression()
	(&InterpolatedStringLiteral{}).isExpression()
	(&Identifier{}).isExpression()
	(&BinaryExpression{}).isExpression()
	(&UnaryExpression{}).isExpression()
	(&CallExpression{}).isExpression()
	(&MethodCallExpression{}).isExpression()
	(&LambdaExpression{}).isExpression()
	(&MatchExpression{}).isExpression()
	(&ResultExpression{}).isExpression()
	(&FieldAccessExpression{}).isExpression()
	(&ModuleAccessExpression{}).isExpression()
	(&SpawnExpression{}).isExpression()
	(&AwaitExpression{}).isExpression()
	(&YieldExpression{}).isExpression()
	(&ChannelExpression{}).isExpression()
	(&ChannelSendExpression{}).isExpression()
	(&ChannelRecvExpression{}).isExpression()
	(&SelectExpression{}).isExpression()
	(&ChannelCreateExpression{}).isExpression()
	(&TypeConstructorExpression{}).isExpression()
	(&BlockExpression{}).isExpression()
	(&ListLiteral{}).isExpression()
	(&ObjectLiteral{}).isExpression()
	(&MapLiteral{}).isExpression()
	(&ListAccessExpression{}).isExpression()
	(&UpdateExpression{}).isExpression()
	(&PerformExpression{}).isExpression()
	(&HandlerExpression{}).isExpression()
}
