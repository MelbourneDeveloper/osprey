package ast

import (
	"strconv"
	"strings"

	"github.com/christianfindlay/osprey/parser"
)

func (b *Builder) buildPrimary(ctx parser.IPrimaryContext) Expression {
	if ctx == nil {
		return nil
	}

	// Handle fiber-related expressions
	if fiberExpr := b.buildFiberExpression(ctx); fiberExpr != nil {
		return fiberExpr
	}

	// Handle algebraic effects expressions
	if effectExpr := b.buildEffectExpression(ctx); effectExpr != nil {
		return effectExpr
	}

	// Handle other expressions
	switch {
	case ctx.BlockExpr() != nil:
		return b.buildBlockExpression(ctx.BlockExpr())
	case ctx.Literal() != nil:
		return b.buildLiteral(ctx.Literal())
	case ctx.LambdaExpr() != nil:
		return b.buildLambdaExpr(ctx.LambdaExpr())
	case ctx.UpdateExpr() != nil:
		return b.buildUpdateExpression(ctx.UpdateExpr())
	case ctx.ObjectLiteral() != nil:
		return b.buildObjectLiteral(ctx.ObjectLiteral())
	case ctx.ID(0) != nil:
		return &Identifier{
			Name:     ctx.ID(0).GetText(),
			Position: b.getPosition(ctx.ID(0).GetSymbol()),
		}
	case ctx.Expr(0) != nil:
		return b.buildExpression(ctx.Expr(0))
	}

	return nil
}

// buildFiberExpression handles fiber-related expressions.
func (b *Builder) buildFiberExpression(ctx parser.IPrimaryContext) Expression {
	switch {
	case ctx.SPAWN() != nil:
		return b.buildSpawnExpression(ctx)
	case ctx.AWAIT() != nil:
		return b.buildAwaitExpression(ctx)
	case ctx.SEND() != nil:
		return b.buildSendExpression(ctx)
	case ctx.RECV() != nil:
		return b.buildRecvExpression(ctx)
	case ctx.YIELD() != nil:
		return b.buildYieldExpression(ctx)
	case ctx.SELECT() != nil:
		return b.buildSelectExpression(ctx.SelectExpr().(*parser.SelectExprContext))
	case ctx.TypeConstructor() != nil:
		return b.buildTypeConstructor(ctx.TypeConstructor().(*parser.TypeConstructorContext))
	}

	return nil
}

// buildEffectExpression handles algebraic effects expressions.
func (b *Builder) buildEffectExpression(ctx parser.IPrimaryContext) Expression {
	switch {
	case ctx.PERFORM() != nil:
		return b.buildPerformExpression(ctx)
	case ctx.HandlerExpr() != nil:
		return b.buildHandlerExpression(ctx)
	}

	return nil
}

func (b *Builder) buildLiteral(ctx parser.ILiteralContext) Expression {
	switch {
	case ctx.FLOAT() != nil:
		text := ctx.FLOAT().GetText()
		value, _ := strconv.ParseFloat(text, 64)

		return &FloatLiteral{
			Value:    value,
			Position: b.getPosition(ctx.FLOAT().GetSymbol()),
		}
	case ctx.INT() != nil:
		text := ctx.INT().GetText()
		value, _ := strconv.ParseInt(text, 10, 64)

		return &IntegerLiteral{
			Value:    value,
			Position: b.getPosition(ctx.INT().GetSymbol()),
		}
	case ctx.STRING() != nil:
		text := ctx.STRING().GetText()
		// Remove the surrounding quotes and process escape sequences. Using
		// strings.Trim was a bug — `"Quote: \"hi\""` ended with `\"` followed
		// by the closing `"`, and Trim greedily stripped the `\"`'s `"` as
		// well, leaving a stray backslash and dropping the closing quote.
		value := strings.TrimSuffix(strings.TrimPrefix(text, "\""), "\"")
		value = b.processEscapeSequences(value)

		return &StringLiteral{
			Value:    value,
			Position: b.getPosition(ctx.STRING().GetSymbol()),
		}
	case ctx.INTERPOLATED_STRING() != nil:
		return b.buildInterpolatedString(ctx.INTERPOLATED_STRING().GetText())
	case ctx.TRUE() != nil:
		return &BooleanLiteral{
			Value:    true,
			Position: b.getPosition(ctx.TRUE().GetSymbol()),
		}
	case ctx.FALSE() != nil:
		return &BooleanLiteral{
			Value:    false,
			Position: b.getPosition(ctx.FALSE().GetSymbol()),
		}
	case ctx.ListLiteral() != nil:
		return b.buildListLiteral(ctx.ListLiteral())
	case ctx.MapLiteral() != nil:
		return b.buildMapLiteral(ctx.MapLiteral())
	}

	return nil
}

// buildInterpolatedString parses an interpolated string like "Hello ${name}!".
func (b *Builder) buildInterpolatedString(text string) Expression {
	// strings.Trim was greedy and ate trailing escape `\"` sequences along
	// with the closing quote — see the matching fix in the STRING branch.
	text = strings.TrimSuffix(strings.TrimPrefix(text, "\""), "\"")
	parts := b.parseInterpolatedParts(text)

	// Process escape sequences in text parts
	for i := range parts {
		if !parts[i].IsExpression {
			parts[i].Text = b.processEscapeSequences(parts[i].Text)
		}
	}

	return &InterpolatedStringLiteral{
		Parts:    parts,
		Position: nil, // Position will be set by caller if available
	}
}

// buildListLiteral builds a ListLiteral from a list literal context.
func (b *Builder) buildListLiteral(ctx parser.IListLiteralContext) Expression {
	if ctx == nil {
		return nil
	}

	elements := make([]Expression, 0)

	// Build each expression in the list
	for _, exprCtx := range ctx.AllExpr() {
		element := b.buildExpression(exprCtx)
		if element != nil {
			elements = append(elements, element)
		}
	}

	return &ListLiteral{
		Elements: elements,
		Position: b.getPositionFromContext(ctx),
	}
}

// buildMapLiteral builds a MapLiteral from a map literal context.
func (b *Builder) buildMapLiteral(ctx parser.IMapLiteralContext) Expression {
	if ctx == nil {
		return nil
	}

	entries := make([]MapEntry, 0)

	// Build each key-value pair in the map
	for _, entryCtx := range ctx.AllMapEntry() {
		key := b.buildExpression(entryCtx.AllExpr()[0])
		value := b.buildExpression(entryCtx.AllExpr()[1])

		if key != nil && value != nil {
			entries = append(entries, MapEntry{
				Key:   key,
				Value: value,
			})
		}
	}

	return &MapLiteral{
		Entries:  entries,
		Position: b.getPositionFromContext(ctx),
	}
}

// buildLambdaExpr builds a LambdaExpression from a lambda context.
func (b *Builder) buildLambdaExpr(ctx parser.ILambdaExprContext) Expression {
	if ctx == nil {
		return nil
	}

	params := make([]Parameter, 0)

	// Handle parameter list if present
	if ctx.ParamList() != nil {
		for _, paramCtx := range ctx.ParamList().AllParam() {
			param := Parameter{
				Name: paramCtx.ID().GetText(),
				Type: nil, // Parse type annotation if present
			}

			// Parse parameter type annotation if present
			if paramCtx.Type_() != nil {
				param.Type = b.buildTypeExpression(paramCtx.Type_())
			}

			params = append(params, param)
		}
	}

	// Parse return type annotation if present
	var returnType *TypeExpression
	if ctx.Type_() != nil {
		returnType = b.buildTypeExpression(ctx.Type_())
	}

	// Parse the lambda body
	body := b.buildExpression(ctx.Expr())

	return &LambdaExpression{
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}
}

const (
	// escByte is the ASCII ESC control character emitted by the "\e" escape; it
	// begins ANSI escape sequences so colored TUIs can be written in Osprey.
	escByte = 0x1b
	// hexLetterValue is the value of the first hex letter ('a'/'A' decode to 10).
	hexLetterValue = 10
	// hexEscapeDigits is the number of hex digits consumed by a "\xHH" escape.
	hexEscapeDigits = 2
)

// hexDigit decodes one ASCII hex digit, reporting whether it was valid.
func hexDigit(c byte) (int, bool) {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0'), true
	case c >= 'a' && c <= 'f':
		return int(c-'a') + hexLetterValue, true
	case c >= 'A' && c <= 'F':
		return int(c-'A') + hexLetterValue, true
	}

	return 0, false
}

// parseHexByte decodes the two hex digits at pos into a byte.
func parseHexByte(s string, pos int) (byte, bool) {
	if pos+1 >= len(s) {
		return 0, false
	}

	hi, okHi := hexDigit(s[pos])
	lo, okLo := hexDigit(s[pos+1])
	if !okHi || !okLo {
		return 0, false
	}

	return byte(hi<<4 | lo), true
}

// processEscapeSequences processes common escape sequences in string literals.
// Walks the input once, left to right, so that `\\n` is correctly read as
// literal backslash + n (the previous ReplaceAll cascade replaced the inner
// `\n` first and produced a real newline). Also supports `\e` (ESC) and
// `\xHH` (hex byte) so ANSI color codes can be written directly.
func (b *Builder) processEscapeSequences(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	for i := 0; i < len(input); i++ {
		c := input[i]
		if c != '\\' || i+1 >= len(input) {
			sb.WriteByte(c)
			continue
		}
		next := input[i+1]
		switch next {
		case 'n':
			sb.WriteByte('\n')
		case 't':
			sb.WriteByte('\t')
		case 'r':
			sb.WriteByte('\r')
		case '0':
			sb.WriteByte(0)
		case '\\':
			sb.WriteByte('\\')
		case '"':
			sb.WriteByte('"')
		case 'e':
			sb.WriteByte(escByte)
		case 'x':
			if hb, ok := parseHexByte(input, i+hexEscapeDigits); ok {
				sb.WriteByte(hb)
				i += hexEscapeDigits // skip the two hex digits; backslash+x consumed below
			} else {
				sb.WriteByte(c)
				sb.WriteByte(next)
			}
		default:
			// Unknown escape — preserve verbatim so users can spot the typo.
			sb.WriteByte(c)
			sb.WriteByte(next)
		}
		i++
	}

	return sb.String()
}

// buildTypeConstructor builds type constructor expressions like Fiber<Int> { computation: fn() => 42 }.
func (b *Builder) buildTypeConstructor(ctx *parser.TypeConstructorContext) Expression {
	typeName := ctx.ID().GetText()

	// Build field assignments
	fieldAssignments := make(map[string]Expression)

	if ctx.FieldAssignments() != nil {
		for _, fieldCtx := range ctx.FieldAssignments().AllFieldAssignment() {
			fieldName := fieldCtx.ID().GetText()
			fieldValue := b.buildExpression(fieldCtx.Expr())
			fieldAssignments[fieldName] = fieldValue
		}
	}

	// Handle specific fiber types
	switch typeName {
	case "Fiber":
		if computation, exists := fieldAssignments["computation"]; exists {
			return &SpawnExpression{
				Expression: computation,
			}
		}
	case "Channel":
		if capacity, exists := fieldAssignments["capacity"]; exists {
			return &ChannelCreateExpression{
				Capacity: capacity,
			}
		}
	}

	// For other types, return a generic type constructor
	return &TypeConstructorExpression{
		TypeName: typeName,
		Fields:   fieldAssignments,
		Position: b.getPositionFromContext(ctx),
	}
}

// buildBlockExpression builds a BlockExpression from a block context.
func (b *Builder) buildBlockExpression(ctx parser.IBlockExprContext) Expression {
	if ctx == nil {
		return nil
	}

	blockBody := ctx.BlockBody()
	if blockBody == nil {
		return &BlockExpression{
			Statements: []Statement{},
			Expression: nil,
		}
	}

	// Build all statements in the block
	statements := make([]Statement, 0)

	for _, stmtCtx := range blockBody.AllStatement() {
		stmt := b.buildStatement(stmtCtx)
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	// Build the final expression if present
	var finalExpr Expression
	if blockBody.Expr() != nil {
		finalExpr = b.buildExpression(blockBody.Expr())
	}

	return &BlockExpression{
		Statements: statements,
		Expression: finalExpr,
	}
}

// buildPerformExpression builds a PerformExpression from perform EffectName.operation(args) syntax.
func (b *Builder) buildPerformExpression(ctx parser.IPrimaryContext) *PerformExpression {
	// PERFORM ID DOT ID LPAREN argList? RPAREN
	effectName := ctx.ID(0).GetText()
	operationName := ctx.ID(1).GetText()

	var arguments []Expression

	if ctx.ArgList() != nil {
		args, _ := b.buildArguments(ctx.ArgList()) // Ignore named args for now
		arguments = args
	}

	return &PerformExpression{
		EffectName:    effectName,
		OperationName: operationName,
		Arguments:     arguments,
		Position:      b.getPositionFromContext(ctx),
	}
}

// buildObjectLiteral builds an ObjectLiteral from an object literal context.
func (b *Builder) buildObjectLiteral(ctx parser.IObjectLiteralContext) Expression {
	if ctx == nil {
		return nil
	}

	// Build field assignments
	fieldAssignments := make(map[string]Expression)

	if ctx.FieldAssignments() != nil {
		for _, fieldCtx := range ctx.FieldAssignments().AllFieldAssignment() {
			fieldName := fieldCtx.ID().GetText()
			fieldValue := b.buildExpression(fieldCtx.Expr())
			fieldAssignments[fieldName] = fieldValue
		}
	}

	return &ObjectLiteral{
		Fields:   fieldAssignments,
		Position: b.getPositionFromContext(ctx),
	}
}

// buildUpdateExpression builds an UpdateExpression for non-destructive record updates.
func (b *Builder) buildUpdateExpression(ctx parser.IUpdateExprContext) Expression {
	if ctx == nil {
		return nil
	}

	// Get the target identifier name
	targetName := ctx.ID().GetText()

	// Build field assignments for the update
	fieldAssignments := make(map[string]Expression)

	if ctx.FieldAssignments() != nil {
		for _, fieldCtx := range ctx.FieldAssignments().AllFieldAssignment() {
			fieldName := fieldCtx.ID().GetText()
			fieldValue := b.buildExpression(fieldCtx.Expr())
			fieldAssignments[fieldName] = fieldValue
		}
	}

	return &UpdateExpression{
		Target: &Identifier{
			Name:     targetName,
			Position: b.getPosition(ctx.ID().GetSymbol()),
		},
		Fields:   fieldAssignments,
		Position: b.getPositionFromContext(ctx),
	}
}

// buildHandlerExpression builds a HandlerExpression from handle...in syntax.
func (b *Builder) buildHandlerExpression(ctx parser.IPrimaryContext) *HandlerExpression {
	handlerCtx := ctx.HandlerExpr()
	if handlerCtx == nil {
		return nil
	}

	// Get effect name
	effectName := handlerCtx.ID().GetText()

	// Build handler arms
	handlers := make([]HandlerArm, 0)

	for _, armCtx := range handlerCtx.AllHandlerArm() {
		operationName := armCtx.ID().GetText()

		// Get parameters if present
		var parameters []string

		if armCtx.HandlerParams() != nil {
			for _, idCtx := range armCtx.HandlerParams().AllID() {
				parameters = append(parameters, idCtx.GetText())
			}
		}

		// Build handler body
		body := b.buildExpression(armCtx.Expr())

		handlers = append(handlers, HandlerArm{
			OperationName: operationName,
			Parameters:    parameters,
			Body:          body,
			Position:      b.getPositionFromContext(armCtx),
		})
	}

	// Build the expression that the handler wraps
	bodyExpr := b.buildExpression(handlerCtx.Expr())

	return &HandlerExpression{
		EffectName: effectName,
		Handlers:   handlers,
		Body:       bodyExpr,
		Position:   b.getPositionFromContext(ctx),
	}
}
