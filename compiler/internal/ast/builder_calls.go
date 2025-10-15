package ast

import (
	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildCallExpr(ctx parser.ICallExprContext) Expression {
	if !b.isValidCallContext(ctx) {
		return nil
	}

	primary := b.buildPrimary(ctx.Primary())
	if primary == nil {
		return nil
	}

	return b.buildCallFromPrimary(ctx, primary)
}

// isValidCallContext checks if the call context is valid.
func (b *Builder) isValidCallContext(ctx parser.ICallExprContext) bool {
	return ctx != nil && ctx.Primary() != nil
}

// buildCallFromPrimary builds a call expression from a primary expression.
func (b *Builder) buildCallFromPrimary(ctx parser.ICallExprContext, primary Expression) Expression {
	if len(ctx.AllLSQUARE()) > 0 {
		return b.buildArrayAccess(ctx, primary)
	}

	if len(ctx.AllDOT()) > 0 {
		return b.buildChainedCall(ctx, primary)
	}

	return b.buildSimpleCall(ctx, primary)
}

// buildChainedCall handles method chaining and field access.
func (b *Builder) buildChainedCall(ctx parser.ICallExprContext, primary Expression) Expression {
	result := primary

	for i := range ctx.AllDOT() {
		// Check if the ID at this index exists
		if i >= len(ctx.AllID()) || ctx.ID(i) == nil {
			continue
		}

		fieldName := ctx.ID(i).GetText()
		result = b.buildChainElement(ctx, result, fieldName, i)
	}

	return result
}

// buildChainElement builds a single element in a method/field chain.
func (b *Builder) buildChainElement(
	ctx parser.ICallExprContext,
	object Expression,
	fieldName string,
	index int,
) Expression {
	// Check if this is module access (Module.function) - only for identifiers that are module names
	if ident, ok := object.(*Identifier); ok && b.isModuleName(ident.Name) {
		return b.buildModuleAccess(ctx, ident, fieldName, index)
	}

	if b.isMethodCallAtIndex(ctx, index) {
		return b.buildMethodCallAtIndex(ctx, object, fieldName, index)
	}

	// This should handle sum.value, x.value, etc.
	return &FieldAccessExpression{
		Object:    object,
		FieldName: fieldName,
		Position:  b.getPositionFromContext(ctx),
	}
}

// buildModuleAccess builds module access expressions.
func (b *Builder) buildModuleAccess(
	ctx parser.ICallExprContext,
	ident *Identifier,
	fieldName string,
	index int,
) Expression {
	if b.isMethodCallAtIndex(ctx, index) {
		// Get arguments if this is a function call
		var (
			args      []Expression
			namedArgs []NamedArgument
		)

		if index < len(ctx.AllArgList()) && ctx.ArgList(index) != nil {
			args, namedArgs = b.buildArguments(ctx.ArgList(index))
		}

		return &ModuleAccessExpression{
			ModuleName:     ident.Name,
			MemberName:     fieldName,
			Arguments:      args,
			NamedArguments: namedArgs,
			Position:       b.getPositionFromContext(ctx),
		}
	}

	// Simple module member access (Module.constant)
	return &ModuleAccessExpression{
		ModuleName: ident.Name,
		MemberName: fieldName,
		Position:   b.getPositionFromContext(ctx),
	}
}

// isMethodCallAtIndex checks if the element at the given index is a method call.
func (b *Builder) isMethodCallAtIndex(ctx parser.ICallExprContext, index int) bool {
	return index < len(ctx.AllLPAREN()) && ctx.LPAREN(index) != nil
}

// buildMethodCallAtIndex builds a method call at the given index.
func (b *Builder) buildMethodCallAtIndex(
	ctx parser.ICallExprContext,
	object Expression,
	methodName string,
	index int,
) Expression {
	var (
		args      []Expression
		namedArgs []NamedArgument
	)

	if index < len(ctx.AllArgList()) && ctx.ArgList(index) != nil {
		args, namedArgs = b.buildArguments(ctx.ArgList(index))
	}

	return &MethodCallExpression{
		Object:         object,
		MethodName:     methodName,
		Arguments:      args,
		NamedArguments: namedArgs,
		Position:       b.getPositionFromContext(ctx),
	}
}

// buildSimpleCall handles simple function calls.
func (b *Builder) buildSimpleCall(ctx parser.ICallExprContext, primary Expression) Expression {
	if ctx.LPAREN(0) == nil {
		return primary
	}

	var (
		args      []Expression
		namedArgs []NamedArgument
	)

	if ctx.ArgList(0) != nil {
		args, namedArgs = b.buildArguments(ctx.ArgList(0))
	}

	if ident, ok := primary.(*Identifier); ok {
		return &CallExpression{
			Function:       ident,
			Arguments:      args,
			NamedArguments: namedArgs,
			Position:       b.getPositionFromContext(ctx),
		}
	}

	return primary
}

// isModuleName checks if the given name is a known module name.
// For now, we'll use a simple heuristic: module names start with uppercase
// In a full implementation, this would check against a registry of declared modules.
func (b *Builder) isModuleName(name string) bool {
	if len(name) == 0 {
		return false
	}

	return name[0] >= 'A' && name[0] <= 'Z'
}

// buildArrayAccess handles array/map access expressions like expr[index].
func (b *Builder) buildArrayAccess(ctx parser.ICallExprContext, primary Expression) Expression {
	result := primary

	// Handle multiple array accesses: arr[0][1][2]
	for i := range ctx.AllLSQUARE() {
		if i >= len(ctx.AllExpr()) || ctx.Expr(i) == nil {
			continue
		}

		indexExpr := b.buildExpression(ctx.Expr(i))
		if indexExpr == nil {
			continue
		}

		result = &ListAccessExpression{
			List:     result,
			Index:    indexExpr,
			Position: b.getPositionFromContext(ctx),
		}
	}

	return result
}
