package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) generateCallExpression(callExpr *ast.CallExpression) (value.Value, error) {
	// Handle function calls
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		// Handle built-in functions
		if result, err := g.handleBuiltInFunction(ident.Name, callExpr); result != nil || err != nil {
			return result, err
		}

		// Handle user-defined functions
		if fn, exists := g.functions[ident.Name]; exists {
			return g.generateUserFunctionCall(ident.Name, fn, callExpr)
		}
	}

	return nil, ErrUnsupportedCall
}

// handleBuiltInFunction handles all built-in function calls.
func (g *LLVMGenerator) handleBuiltInFunction(name string, callExpr *ast.CallExpression) (value.Value, error) {
	// Try core functions first (always available)
	if result, err := g.handleCoreFunctions(name, callExpr); result != nil || err != nil {
		return result, err
	}

	// Try HTTP functions only if allowed by security policy
	if g.security.AllowHTTP {
		if result, err := g.handleHTTPFunctions(name, callExpr); result != nil || err != nil {
			return result, err
		}
	}

	// Try WebSocket functions only if allowed by security policy
	if g.security.AllowWebSocket {
		if result, err := g.handleWebSocketFunctions(name, callExpr); result != nil || err != nil {
			return result, err
		}
	}

	// Not a built-in function or not allowed by security policy
	return nil, nil
}

// handleCoreFunctions handles core built-in functions like print, toString, etc.
func (g *LLVMGenerator) handleCoreFunctions(name string, callExpr *ast.CallExpression) (value.Value, error) {
	switch name {
	case ToStringFunc:
		return g.generateToStringCall(callExpr)
	case PrintFunc:
		return g.generatePrintCall(callExpr)
	case InputFunc:
		return g.generateInputCall(callExpr)
	case RangeFunc:
		return g.generateRangeCall(callExpr)
	case ForEachFunc:
		return g.generateForEachCall(callExpr)
	case MapFunc:
		return g.generateMapCall(callExpr)
	case FilterFunc:
		return g.generateFilterCall(callExpr)
	case FoldFunc:
		return g.generateFoldCall(callExpr)
	case LengthFunc:
		return g.generateLengthCall(callExpr)
	case ContainsFunc:
		return g.generateContainsCall(callExpr)
	case SubstringFunc:
		return g.generateSubstringCall(callExpr)
	case SpawnProcessFunc:
		return g.generateSpawnProcessCall(callExpr)
	case WriteFileFunc:
		return g.generateWriteFileCall(callExpr)
	case ReadFileFunc:
		return g.generateReadFileCall(callExpr)
	case ParseJSONFunc:
		return g.generateParseJSONCall(callExpr)
	case ExtractCodeFunc:
		return g.generateExtractCodeCall(callExpr)
	default:
		return nil, nil
	}
}

// handleHTTPFunctions handles HTTP-related built-in functions.
func (g *LLVMGenerator) handleHTTPFunctions(name string, callExpr *ast.CallExpression) (value.Value, error) {
	switch name {
	case HTTPCreateServerFunc:
		return g.generateHTTPCreateServerCall(callExpr)
	case HTTPListenFunc:
		return g.generateHTTPListenCall(callExpr)
	case HTTPStopServerFunc:
		return g.generateHTTPStopServerCall(callExpr)
	case HTTPCreateClientFunc:
		return g.generateHTTPCreateClientCall(callExpr)
	case HTTPGetFunc:
		return g.generateHTTPGetCall(callExpr)
	case HTTPPostFunc:
		return g.generateHTTPPostCall(callExpr)
	case HTTPPutFunc:
		return g.generateHTTPPutCall(callExpr)
	case HTTPDeleteFunc:
		return g.generateHTTPDeleteCall(callExpr)
	case HTTPRequestFunc:
		return g.generateHTTPRequestCall(callExpr)
	case HTTPCloseClientFunc:
		return g.generateHTTPCloseClientCall(callExpr)
	default:
		return nil, nil
	}
}

// handleWebSocketFunctions handles WebSocket-related built-in functions.
func (g *LLVMGenerator) handleWebSocketFunctions(name string, callExpr *ast.CallExpression) (value.Value, error) {
	switch name {
	case WebSocketConnectFunc:
		return g.generateWebSocketConnectCall(callExpr)
	case WebSocketSendFunc:
		return g.generateWebSocketSendCall(callExpr)
	case WebSocketCloseFunc:
		return g.generateWebSocketCloseCall(callExpr)
	case WebSocketCreateServerFunc:
		return g.generateWebSocketCreateServerCall(callExpr)
	case WebSocketServerListenFunc:
		return g.generateWebSocketServerListenCall(callExpr)
	case WebSocketServerBroadcastFunc:
		return g.generateWebSocketServerBroadcastCall(callExpr)
	case WebSocketStopServerFunc:
		return g.generateWebSocketStopServerCall(callExpr)
	case WebSocketKeepAlive:
		return g.generateWebSocketKeepAliveCall(callExpr)
	default:
		return nil, nil
	}
}

// generateUserFunctionCall handles user-defined function calls.
func (g *LLVMGenerator) generateUserFunctionCall(
	funcName string,
	fn *ir.Func,
	callExpr *ast.CallExpression,
) (value.Value, error) {
	// Check if function has multiple parameters and enforce named arguments
	params, paramExists := g.functionParameters[funcName]
	if paramExists && len(params) > 1 {
		// Multi-parameter function - MUST use named arguments
		if len(callExpr.NamedArguments) == 0 {
			example := g.buildNamedArgumentsExample(params)

			return nil, WrapFunctionRequiresNamed(funcName, len(params), example)
		}

		if len(callExpr.Arguments) > 0 {
			example := g.buildNamedArgumentsExample(params)

			return nil, WrapFunctionRequiresNamed(funcName, len(params), example)
		}
	}

	// Handle named arguments vs positional arguments
	if len(callExpr.NamedArguments) > 0 {
		// Validate named arguments are not of type 'any'
		for _, namedArg := range callExpr.NamedArguments {
			if err := g.validateNotAnyType(namedArg.Value, AnyOpFunctionArgument); err != nil {
				return nil, WrapAnyDirectFunctionArg(funcName, "unknown")
			}
		}

		// Named arguments - need to reorder them to match function signature
		args, err := g.reorderNamedArguments(funcName, callExpr.NamedArguments)
		if err != nil {
			return nil, err
		}

		return g.builder.NewCall(fn, args...), nil
	}

	// Positional arguments (traditional)
	args := make([]value.Value, len(callExpr.Arguments))

	for i, arg := range callExpr.Arguments {
		// Validate that argument is not of type 'any'
		if err := g.validateNotAnyType(arg, AnyOpFunctionArgument); err != nil {
			return nil, WrapAnyDirectFunctionArg(funcName, "unknown")
		}

		val, err := g.generateExpression(arg)
		if err != nil {
			return nil, err
		}

		args[i] = val
	}

	return g.builder.NewCall(fn, args...), nil
}

// generateInterpolatedString generates LLVM IR for interpolated strings by concatenating parts.
func (g *LLVMGenerator) generateInterpolatedString(interpStr *ast.InterpolatedStringLiteral) (value.Value, error) {
	// For now, we'll use a simple approach: build the string by calling printf multiple times
	// into a buffer. A more sophisticated implementation would use string concatenation.
	// If there's only one part and it's text, treat it as a regular string
	if len(interpStr.Parts) == 1 && !interpStr.Parts[0].IsExpression {
		str := constant.NewCharArrayFromString(interpStr.Parts[0].Text + StringTerminator)
		global := g.module.NewGlobalDef("", str)

		return g.builder.NewGetElementPtr(str.Typ, global,
			constant.NewInt(types.I32, ArrayIndexZero),
			constant.NewInt(types.I32, ArrayIndexZero)), nil
	}

	// For multiple parts or expressions, we'll use sprintf to build the string
	// First, we need sprintf declared
	sprintf := g.ensureSprintfDeclaration()

	// Build format string and collect arguments
	var formatParts []string

	var args []value.Value

	for _, part := range interpStr.Parts {
		if part.IsExpression {
			// Auto-call toString() on all expressions in string interpolation
			toStringCall := &ast.CallExpression{
				Function:  &ast.Identifier{Name: "toString"},
				Arguments: []ast.Expression{part.Expression},
			}

			// Generate the toString call which will return a string
			val, err := g.generateExpression(toStringCall)
			if err != nil {
				return nil, err
			}

			args = append(args, val)

			// All expressions become %s since toString() always returns string
			formatParts = append(formatParts, "%s")
		} else {
			// Escape % characters in literal text by replacing % with %%
			escapedText := strings.ReplaceAll(part.Text, "%", "%%")
			formatParts = append(formatParts, escapedText)
		}
	}

	// Create the format string
	formatString := strings.Join(formatParts, "") + StringTerminator
	formatStr := constant.NewCharArrayFromString(formatString)
	formatGlobal := g.module.NewGlobalDef("", formatStr)
	formatPtr := g.builder.NewGetElementPtr(formatStr.Typ, formatGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Allocate buffer for result string (simplified - use fixed size)
	bufferType := types.NewArray(BufferSize1KB, types.I8) // 1KB buffer
	buffer := g.builder.NewAlloca(bufferType)
	bufferPtr := g.builder.NewGetElementPtr(bufferType, buffer,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Call sprintf(buffer, format, args...)
	sprintfArgs := make([]value.Value, 0, len(args)+TwoArgs)
	sprintfArgs = append(sprintfArgs, bufferPtr, formatPtr)
	sprintfArgs = append(sprintfArgs, args...)

	g.builder.NewCall(sprintf, sprintfArgs...)

	return bufferPtr, nil
}

// ensureSprintfDeclaration declares sprintf if not already declared.
func (g *LLVMGenerator) ensureSprintfDeclaration() *ir.Func {
	if sprintf, exists := g.functions["sprintf"]; exists {
		return sprintf
	}

	// Declare sprintf: i32 @sprintf(i8* %str, i8* %format, ...)
	sprintf := g.module.NewFunc("sprintf", types.I32,
		ir.NewParam("str", types.I8Ptr),
		ir.NewParam("format", types.I8Ptr))
	sprintf.Sig.Variadic = true
	g.functions["sprintf"] = sprintf

	return sprintf
}

func (g *LLVMGenerator) generateIntToString(arg value.Value) (value.Value, error) {
	// Ensure sprintf is declared
	sprintf := g.ensureSprintfDeclaration()

	// Create format string for integer conversion
	formatStr := constant.NewCharArrayFromString("%ld\x00")
	formatGlobal := g.module.NewGlobalDef("", formatStr)
	formatPtr := g.builder.NewGetElementPtr(formatStr.Typ, formatGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Allocate buffer for result string (64 bytes should be enough for any 64-bit integer)
	bufferType := types.NewArray(BufferSize64Bytes, types.I8)
	buffer := g.builder.NewAlloca(bufferType)
	bufferPtr := g.builder.NewGetElementPtr(bufferType, buffer,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Call sprintf(buffer, "%ld", arg)
	g.builder.NewCall(sprintf, bufferPtr, formatPtr, arg)

	return bufferPtr, nil
}

func (g *LLVMGenerator) generateBoolToString(arg value.Value) (value.Value, error) {
	// Create blocks for true/false cases
	blockSuffix := fmt.Sprintf("_%p", arg) // Use pointer address for uniqueness
	currentBlock := g.builder

	trueBlock := g.function.NewBlock("bool_true" + blockSuffix)
	falseBlock := g.function.NewBlock("bool_false" + blockSuffix)
	endBlock := g.function.NewBlock("bool_end" + blockSuffix)

	// Check if arg == 1 (true) or 0 (false)
	zero := constant.NewInt(types.I64, ArrayIndexZero)
	isTrue := currentBlock.NewICmp(enum.IPredNE, arg, zero)
	currentBlock.NewCondBr(isTrue, trueBlock, falseBlock)

	// True case - return "true"
	g.builder = trueBlock
	trueStr := constant.NewCharArrayFromString(TrueString)
	trueGlobal := g.module.NewGlobalDef("", trueStr)
	truePtr := trueBlock.NewGetElementPtr(trueStr.Typ, trueGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	trueBlock.NewBr(endBlock)

	// False case - return "false"
	g.builder = falseBlock
	falseStr := constant.NewCharArrayFromString(FalseString)
	falseGlobal := g.module.NewGlobalDef("", falseStr)
	falsePtr := falseBlock.NewGetElementPtr(falseStr.Typ, falseGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	falseBlock.NewBr(endBlock)

	// Create PHI node in end block
	g.builder = endBlock
	phi := endBlock.NewPhi(ir.NewIncoming(truePtr, trueBlock), ir.NewIncoming(falsePtr, falseBlock))

	return phi, nil
}

func (g *LLVMGenerator) generateMatchExpression(matchExpr *ast.MatchExpression) (value.Value, error) {
	// Validate match expression for exhaustiveness and unknown variants
	if err := g.validateMatchExpression(matchExpr); err != nil {
		return nil, err
	}

	discriminant, err := g.generateExpression(matchExpr.Expression)
	if err != nil {
		return nil, err
	}

	return g.generateMatchExpressionWithDiscriminant(matchExpr, discriminant)
}

// generateMatchExpressionWithDiscriminant generates match expression with pre-computed discriminant.
func (g *LLVMGenerator) generateMatchExpressionWithDiscriminant(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	if g.hasResultPatterns(matchExpr.Arms) {
		return g.generateResultMatchExpression(matchExpr, discriminant)
	}

	if len(matchExpr.Arms) == 0 {
		return constant.NewInt(types.I64, ArrayIndexZero), nil
	}

	return g.generateStandardMatchExpression(matchExpr, discriminant)
}

// hasResultPatterns checks if the match expression has Success/Err patterns.
func (g *LLVMGenerator) hasResultPatterns(arms []ast.MatchArm) bool {
	for _, arm := range arms {
		if arm.Pattern.Constructor == SuccessPattern || arm.Pattern.Constructor == ErrorPattern {
			return true
		}
	}

	return false
}

// generateStandardMatchExpression generates a standard (non-result) match expression.
func (g *LLVMGenerator) generateStandardMatchExpression(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	// Create unique block names for this match expression
	blockSuffix := fmt.Sprintf("_%p", matchExpr)
	endBlock := g.function.NewBlock("match_end" + blockSuffix)
	armBlocks := g.createMatchArmBlocks(matchExpr.Arms, blockSuffix)

	// Generate conditions first (creates branching logic)
	g.generateMatchConditions(matchExpr.Arms, armBlocks, discriminant, blockSuffix)

	// Then generate arm values (adds terminators to arm blocks)
	armValues, predecessorBlocks, err := g.generateMatchArmValues(matchExpr.Arms, armBlocks, endBlock, discriminant)
	if err != nil {
		return nil, err
	}

	return g.createMatchResult(armValues, predecessorBlocks, endBlock)
}

// createMatchArmBlocks creates LLVM blocks for each match arm.
func (g *LLVMGenerator) createMatchArmBlocks(arms []ast.MatchArm, blockSuffix string) []*ir.Block {
	var armBlocks []*ir.Block
	for i := range arms {
		armBlocks = append(armBlocks, g.function.NewBlock(fmt.Sprintf("match_arm_%d%s", i, blockSuffix)))
	}

	return armBlocks
}

// generateMatchArmValues generates values for each match arm.
func (g *LLVMGenerator) generateMatchArmValues(
	arms []ast.MatchArm,
	armBlocks []*ir.Block,
	endBlock *ir.Block,
	discriminant value.Value,
) ([]value.Value, []*ir.Block, error) {
	var armValues []value.Value
	var predecessorBlocks []*ir.Block
	oldBuilder := g.builder

	for i, arm := range arms {
		g.builder = armBlocks[i]

		// Handle variable binding in patterns
		if arm.Pattern.Variable != "" {
			// Save the current variable scope
			oldVariables := make(map[string]value.Value)
			for k, v := range g.variables {
				oldVariables[k] = v
			}

			// Bind the pattern variable to the discriminant value
			g.variables[arm.Pattern.Variable] = discriminant

			// Generate the arm expression
			armValue, err := g.generateExpression(arm.Expression)
			if err != nil {
				return nil, nil, err
			}

			// Restore the previous variable scope
			g.variables = oldVariables

			armValues = append(armValues, armValue)
		} else {
			// No variable binding, generate normally
			armValue, err := g.generateExpression(arm.Expression)
			if err != nil {
				return nil, nil, err
			}

			armValues = append(armValues, armValue)
		}

		armBlocks[i].NewBr(endBlock)
		predecessorBlocks = append(predecessorBlocks, armBlocks[i])
	}

	g.builder = oldBuilder

	return armValues, predecessorBlocks, nil
}

// generateMatchConditions generates the conditional branches for pattern matching.
func (g *LLVMGenerator) generateMatchConditions(
	arms []ast.MatchArm,
	armBlocks []*ir.Block,
	discriminant value.Value,
	blockSuffix string,
) {
	currentBlock := g.builder

	for i, arm := range arms {
		// Ensure the builder is set to the current block
		g.builder = currentBlock
		condition := g.createPatternCondition(arm.Pattern, discriminant, currentBlock)

		if i == len(arms)-1 {
			currentBlock.NewBr(armBlocks[i])
		} else {
			nextCheckBlock := g.function.NewBlock(fmt.Sprintf("match_check_%d%s", i+1, blockSuffix))
			currentBlock.NewCondBr(condition, armBlocks[i], nextCheckBlock)
			currentBlock = nextCheckBlock
		}
	}
}

// createPatternCondition creates a condition for pattern matching.
func (g *LLVMGenerator) createPatternCondition(
	pattern ast.Pattern,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	if pattern.Constructor == "_" || pattern.Constructor == UnknownPattern {
		return constant.NewBool(true)
	}

	// Handle variable binding patterns (empty constructor means variable binding)
	if pattern.Constructor == "" && pattern.Variable != "" {
		return constant.NewBool(true)
	}

	// Check if it's a union type variant
	if discriminantValue, exists := g.unionVariants[pattern.Constructor]; exists {
		patternConst := constant.NewInt(types.I64, discriminantValue)

		return currentBlock.NewICmp(enum.IPredEQ, discriminant, patternConst)
	}

	// Check if it's a numeric literal
	if patternValue, err := strconv.ParseInt(pattern.Constructor, 10, 64); err == nil {
		patternConst := constant.NewInt(types.I64, patternValue)

		return currentBlock.NewICmp(enum.IPredEQ, discriminant, patternConst)
	}

	return g.createStringPatternCondition(pattern.Constructor, discriminant, currentBlock)
}

// createStringPatternCondition creates a condition for string pattern matching.
func (g *LLVMGenerator) createStringPatternCondition(
	constructor string,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	if !strings.HasPrefix(constructor, "\"") || !strings.HasSuffix(constructor, "\"") {
		return constant.NewBool(false)
	}

	patternStr := constructor[1 : len(constructor)-1]
	patternConstStr := constant.NewCharArrayFromString(patternStr + StringTerminator)
	patternGlobal := g.module.NewGlobalDef("", patternConstStr)
	patternPtr := currentBlock.NewGetElementPtr(patternConstStr.Typ, patternGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	strcmp := g.functions["strcmp"]
	compareResult := currentBlock.NewCall(strcmp, discriminant, patternPtr)
	zero := constant.NewInt(types.I32, ArrayIndexZero)

	return currentBlock.NewICmp(enum.IPredEQ, compareResult, zero)
}

// createMatchResult creates the final result value using PHI nodes if needed.
func (g *LLVMGenerator) createMatchResult(
	armValues []value.Value,
	predecessorBlocks []*ir.Block,
	endBlock *ir.Block,
) (value.Value, error) {
	g.builder = endBlock

	if len(armValues) == 1 {
		return armValues[0], nil
	}

	// Check if we need type coercion
	coercedValues, err := g.coerceArmValuesToCommonType(armValues)
	if err != nil {
		return nil, err
	}

	var incomings []*ir.Incoming
	for i, val := range coercedValues {
		incomings = append(incomings, ir.NewIncoming(val, predecessorBlocks[i]))
	}

	return endBlock.NewPhi(incomings...), nil
}

// coerceArmValuesToCommonType ensures all arm values have compatible types.
func (g *LLVMGenerator) coerceArmValuesToCommonType(armValues []value.Value) ([]value.Value, error) {
	expectedType := armValues[0].Type()

	if !g.needsTypeCoercion(armValues, expectedType) {
		return armValues, nil
	}

	return g.performTypeCoercion(armValues, expectedType)
}

// needsTypeCoercion checks if type coercion is needed for arm values.
func (g *LLVMGenerator) needsTypeCoercion(armValues []value.Value, expectedType types.Type) bool {
	for i := 1; i < len(armValues); i++ {
		if armValues[i].Type() != expectedType {
			return true
		}
	}

	return false
}

// performTypeCoercion converts all values to the expected type.
func (g *LLVMGenerator) performTypeCoercion(armValues []value.Value, expectedType types.Type) ([]value.Value, error) {
	coercedValues := make([]value.Value, len(armValues))

	for i, val := range armValues {
		if val.Type() == expectedType {
			coercedValues[i] = val
		} else {
			coercedVal, err := g.coerceValueToType(val, expectedType)
			if err != nil {
				return nil, err
			}
			coercedValues[i] = coercedVal
		}
	}

	return coercedValues, nil
}

// coerceValueToType converts a single value to the target type.
func (g *LLVMGenerator) coerceValueToType(val value.Value, targetType types.Type) (value.Value, error) {
	switch targetType {
	case types.I8Ptr:

		return g.convertToString(val)
	case types.I64:

		return constant.NewInt(types.I64, 0), nil
	default:

		return val, nil
	}
}

// convertToString converts a value to string type.
func (g *LLVMGenerator) convertToString(val value.Value) (value.Value, error) {
	if val.Type() == types.I64 {
		return g.generateIntToString(val)
	}

	return val, nil
}

func (g *LLVMGenerator) generateResultMatchExpression(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	blocks := g.createResultMatchBlocks(matchExpr)

	// Store the Result value for pattern binding
	g.currentResultValue = discriminant

	g.generateResultMatchCondition(discriminant, blocks)

	successValue, err := g.generateSuccessBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}

	errorValue, err := g.generateErrorBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}

	return g.createResultMatchPhi(successValue, errorValue, blocks)
}

// ResultMatchBlocks holds the blocks for result match expressions.
type ResultMatchBlocks struct {
	Success *ir.Block
	Error   *ir.Block
	End     *ir.Block
}

// createResultMatchBlocks creates blocks for result match expressions.
func (g *LLVMGenerator) createResultMatchBlocks(matchExpr *ast.MatchExpression) *ResultMatchBlocks {
	blockSuffix := fmt.Sprintf("_%p", matchExpr)

	return &ResultMatchBlocks{
		Success: g.function.NewBlock("success_case" + blockSuffix),
		Error:   g.function.NewBlock("error_case" + blockSuffix),
		End:     g.function.NewBlock("match_end" + blockSuffix),
	}
}

// generateResultMatchCondition generates the condition for result matching.
func (g *LLVMGenerator) generateResultMatchCondition(discriminant value.Value, blocks *ResultMatchBlocks) {
	// Check if the discriminant is a pointer to a struct (Result type) or just an integer
	if ptrType, ok := discriminant.Type().(*types.PointerType); ok {
		// Extract the discriminant field from the Result struct
		// Result struct: [value, discriminant] where discriminant is at index 1
		resultType := ptrType.ElemType
		discriminantPtr := g.builder.NewGetElementPtr(resultType, discriminant,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
		discriminantValue := g.builder.NewLoad(discriminantPtr.Type().(*types.PointerType).ElemType, discriminantPtr)

		// 0 = Success, 1 = Error
		zero := constant.NewInt(types.I8, 0)
		isSuccess := g.builder.NewICmp(enum.IPredEQ, discriminantValue, zero)
		g.builder.NewCondBr(isSuccess, blocks.Success, blocks.Error)
	} else {
		// Fallback: treat as integer discriminant
		zero := constant.NewInt(types.I64, 0)
		isSuccess := g.builder.NewICmp(enum.IPredSGE, discriminant, zero)
		g.builder.NewCondBr(isSuccess, blocks.Success, blocks.Error)
	}
}

// generateSuccessBlock generates the success block for result matching.
func (g *LLVMGenerator) generateSuccessBlock(
	matchExpr *ast.MatchExpression,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.Success

	// Find the success arm and bind pattern variables
	successArm := g.findSuccessArm(matchExpr)
	if successArm != nil && len(successArm.Pattern.Fields) > 0 {
		// Bind the Result value to the pattern variable
		fieldName := successArm.Pattern.Fields[0] // First field is the value

		// Get the Result value from the matched expression
		// The Result struct has: [value, discriminant]
		// We need to extract the value field (index 0)
		if g.currentResultValue != nil {
			if ptrType, ok := g.currentResultValue.Type().(*types.PointerType); ok {
				resultType := ptrType.ElemType
				valuePtr := g.builder.NewGetElementPtr(resultType, g.currentResultValue,
					constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				extractedValue := g.builder.NewLoad(valuePtr.Type().(*types.PointerType).ElemType, valuePtr)
				g.variables[fieldName] = extractedValue
			} else {
				// Fallback: use the discriminant value directly
				g.variables[fieldName] = g.currentResultValue
			}
		}
	}

	successExpr := g.findSuccessValue(matchExpr)
	var successValue value.Value
	if successExpr == nil {
		successValue = constant.NewInt(types.I64, ArrayIndexZero)
	} else {
		val, err := g.generateExpression(successExpr)
		if err != nil {
			return nil, err
		}
		successValue = val
	}

	blocks.Success.NewBr(blocks.End)

	return successValue, nil
}

// generateErrorBlock generates the error block for result matching.
func (g *LLVMGenerator) generateErrorBlock(
	matchExpr *ast.MatchExpression,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.Error

	// Find the Error arm and bind pattern variables
	errorArm := g.findErrorArm(matchExpr)
	if errorArm != nil && len(errorArm.Pattern.Fields) > 0 {
		// Bind the Result error message to the pattern variable
		fieldName := errorArm.Pattern.Fields[0] // First field is the message
		// Create a unique global string for the error message
		blockSuffix := fmt.Sprintf("_%p", matchExpr)
		errorStr := g.module.NewGlobalDef("error_msg"+blockSuffix, constant.NewCharArrayFromString("Error occurred\\x00"))
		errorPtr := g.builder.NewGetElementPtr(errorStr.ContentType, errorStr,
			constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
		g.variables[fieldName] = errorPtr
	}

	errorExpr := g.findErrorValue(matchExpr)
	var errorValue value.Value
	if errorExpr == nil {
		errorValue = constant.NewInt(types.I64, ArrayIndexZero)
	} else {
		val, err := g.generateExpression(errorExpr)
		if err != nil {
			return nil, err
		}
		errorValue = val
	}

	blocks.Error.NewBr(blocks.End)

	return errorValue, nil
}

// findSuccessArm finds the success match arm.
func (g *LLVMGenerator) findSuccessArm(matchExpr *ast.MatchExpression) *ast.MatchArm {
	for i, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == SuccessPattern {
			return &matchExpr.Arms[i]
		}
	}

	return nil
}

// findErrorArm finds the error match arm.
func (g *LLVMGenerator) findErrorArm(matchExpr *ast.MatchExpression) *ast.MatchArm {
	for i, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == ErrorPattern {
			return &matchExpr.Arms[i]
		}
	}

	return nil
}

// findSuccessValue finds the success expression in match arms.
func (g *LLVMGenerator) findSuccessValue(matchExpr *ast.MatchExpression) ast.Expression {
	for _, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == SuccessPattern {
			return arm.Expression
		}
	}

	return nil
}

// findErrorValue finds the error expression in match arms.
func (g *LLVMGenerator) findErrorValue(matchExpr *ast.MatchExpression) ast.Expression {
	for _, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == ErrorPattern {
			return arm.Expression
		}
	}

	return nil
}

// createResultMatchPhi creates the PHI node for result matching.
func (g *LLVMGenerator) createResultMatchPhi(
	successValue, errorValue value.Value,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.End

	// Ensure both values have the same type for PHI node
	if successValue.Type() != errorValue.Type() {
		// Determine which type to use - prefer string if either is string
		if successValue.Type() == types.I8Ptr {
			// Convert errorValue to string if it's not already
			if errorValue.Type() == types.I64 {
				convertedError, err := g.generateIntToString(errorValue)
				if err != nil {
					return nil, err
				}
				errorValue = convertedError
			}
		} else if errorValue.Type() == types.I8Ptr {
			// Convert successValue to string if it's not already
			if successValue.Type() == types.I64 {
				convertedSuccess, err := g.generateIntToString(successValue)
				if err != nil {
					return nil, err
				}
				successValue = convertedSuccess
			}
		}
		// If both are different non-string types, convert both to int
		if successValue.Type() != errorValue.Type() {
			successValue = constant.NewInt(types.I64, 0)
			errorValue = constant.NewInt(types.I64, 0)
		}
	}

	phi := blocks.End.NewPhi(
		ir.NewIncoming(successValue, blocks.Success),
		ir.NewIncoming(errorValue, blocks.Error),
	)

	return phi, nil
}
