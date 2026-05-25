package ast

import (
	"github.com/antlr4-go/antlr/v4"

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
//
// NOTE: `ctx.ArgList(i)` returns the i-th *present* ArgList — empty arg
// lists are skipped — so naively indexing by chain position lines up
// wrong as soon as one method in the chain has no args and another does.
// We walk the children in source order and find the ArgList that lives
// between LPAREN(index) and the matching RPAREN. Fixes UFCS chains like
// `s.trim().take(3)`.
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

	if argList := b.argListForChainElement(ctx, index); argList != nil {
		args, namedArgs = b.buildArguments(argList)
	}

	return &MethodCallExpression{
		Object:         object,
		MethodName:     methodName,
		Arguments:      args,
		NamedArguments: namedArgs,
		Position:       b.getPositionFromContext(ctx),
	}
}

// argListForChainElement returns the ArgList sitting inside the
// chainIndex-th `(...)` of the call expression, or nil if that call has
// no arguments. Walks children in source order to be robust to optional
// args being absent for some chain elements.
func (b *Builder) argListForChainElement(
	ctx parser.ICallExprContext,
	chainIndex int,
) parser.IArgListContext {
	lparenCount := -1
	insideTarget := false
	for _, child := range ctx.GetChildren() {
		if term, ok := child.(antlr.TerminalNode); ok {
			switch term.GetText() {
			case "(":
				lparenCount++
				if lparenCount == chainIndex {
					insideTarget = true
				}
			case ")":
				if insideTarget {
					return nil // matching `)` reached without an ArgList
				}
			}
			continue
		}
		if argList, ok := child.(parser.IArgListContext); ok && insideTarget {
			return argList
		}
	}
	return nil
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

	// Identifier callee (the common path) — wrap in CallExpression.
	if ident, ok := primary.(*Identifier); ok {
		return &CallExpression{
			Function:       ident,
			Arguments:      args,
			NamedArguments: namedArgs,
			Position:       b.getPositionFromContext(ctx),
		}
	}

	// Immediately-invoked lambdas and other parenthesized callee
	// expressions: `(fn(x: int) => x * 10)(5)`. Previously this dropped
	// the (5) args and returned just the lambda, so print(lambda)(5)
	// silently puts'd the lambda's function pointer. Wrap in
	// CallExpression — codegen's resolveFunctionValue handles the
	// non-identifier case via funcName=="" → generateExpression(Function).
	if _, isLambda := primary.(*LambdaExpression); isLambda {
		return &CallExpression{
			Function:       primary,
			Arguments:      args,
			NamedArguments: namedArgs,
			Position:       b.getPositionFromContext(ctx),
		}
	}

	return primary
}

// isModuleName checks if `name` was declared as `module Name { ... }`
// somewhere in this program. Populated by the pre-pass in BuildProgram.
//
// The previous "uppercase identifier == module" heuristic stole every
// `UppercaseVar.method()` call from the UFCS path (buildChainElement)
// and routed it to a placeholder that always returned 42. The actual
// module set is what we want — and effect-handler syntax
// `perform Logger.info(...)` takes a different grammar rule, so it is
// unaffected.
//
// Implements [BUILTIN-STRING-UFCS] disambiguation rule.
func (b *Builder) isModuleName(name string) bool {
	if b.moduleNames == nil {
		return false
	}
	return b.moduleNames[name]
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
