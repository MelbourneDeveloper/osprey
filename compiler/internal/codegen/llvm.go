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

// Constants for type sizes and arrays (shared with expression_generation.go)
const (
	LargeArraySizeForCasting = 1000 // Large array size for type casting operations
)

func (g *LLVMGenerator) generateCallExpression(callExpr *ast.CallExpression) (value.Value, error) {
	// Get function name if it's an identifier
	var funcName string
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		funcName = ident.Name

		// First check if this is a built-in function
		if result, err := g.handleBuiltInFunction(funcName, callExpr); result != nil || err != nil {
			return result, err
		}
	}

	// Validate any type usage in function calls
	if funcName != "" {
		if err := g.validateFunctionCallArguments(funcName, callExpr); err != nil {
			return nil, err
		}
		
		// Validate named arguments requirement for multi-parameter functions
		if err := g.validateNamedArguments(funcName, callExpr); err != nil {
			return nil, err
		}
	}

	// Not a built-in function, generate the function value
	funcValue, err := g.generateExpression(callExpr.Function)
	if err != nil {
		return nil, err
	}

	// Generate arguments
	var args []value.Value
	if len(callExpr.NamedArguments) > 0 {
		// Handle named arguments
		reorderedExprs, err := g.reorderNamedArguments(funcName, callExpr.NamedArguments)
		if err != nil {
			return nil, err
		}
		// Convert expressions to values
		args = make([]value.Value, len(reorderedExprs))
		for i, expr := range reorderedExprs {
			val, err := g.generateExpression(expr)
			if err != nil {
				return nil, err
			}
			args[i] = val
		}
	} else {
		// Handle regular arguments
		args = make([]value.Value, len(callExpr.Arguments))
		for i, arg := range callExpr.Arguments {
			val, err := g.generateExpression(arg)
			if err != nil {
				return nil, err
			}
			args[i] = val
		}
	}

	// Create the function call
	return g.builder.NewCall(funcValue, args...), nil
}

// validateFunctionCallArguments validates that any types are not passed to functions expecting specific types
func (g *LLVMGenerator) validateFunctionCallArguments(funcName string, callExpr *ast.CallExpression) error {
	// Check if we have the function type in the type environment
	funcType, exists := g.typeInferer.env.Get(funcName)
	if !exists {
		return nil // Built-in functions or unknown functions are handled elsewhere
	}

	// Check if it's a function type
	if fnType, ok := funcType.(*FunctionType); ok {
		// Get the arguments to check
		var args []ast.Expression
		if len(callExpr.NamedArguments) > 0 {
			// Reorder named arguments to match parameter order
			reorderedExprs, err := g.reorderNamedArguments(funcName, callExpr.NamedArguments)
			if err != nil {
				return err
			}
			args = reorderedExprs
		} else {
			args = callExpr.Arguments
		}

		// Check each argument against its corresponding parameter type
		for i, arg := range args {
			if i >= len(fnType.paramTypes) {
				continue // Too many arguments, will be caught elsewhere
			}

			paramType := fnType.paramTypes[i]

			// Check for any type mismatch, passing function name and parameter index for better error messages
			if err := g.checkAnyTypeMismatchParam(arg, paramType, callExpr.Position, funcName, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkAnyTypeMismatchParam checks if an any type is being passed to a function expecting a specific type
func (g *LLVMGenerator) checkAnyTypeMismatchParam(
	arg ast.Expression,
	paramType Type,
	pos *ast.Position,
	funcName string,
	paramIndex int,
) error {
	// Check if argument is an identifier with any type
	if ident, ok := arg.(*ast.Identifier); ok {
		// Look up the variable type
		if varType, exists := g.typeInferer.env.Get(ident.Name); exists {
			resolvedType := g.typeInferer.ResolveType(varType)
			resolvedParamType := g.typeInferer.ResolveType(paramType)

			// Check if it's an any type being passed to a specific type
			if resolvedType.String() == TypeAny && resolvedParamType.String() != TypeAny {
				// Get the parameter name from the function declaration
				paramName := g.getParameterName(funcName, paramIndex)
				if paramName == "" {
					paramName = resolvedParamType.String() // fallback to type name
				}

				// Use the identifier's position for accuracy, fallback to provided position
				identPos := ident.Position
				if identPos == nil {
					identPos = pos
				}

				if identPos != nil {
					return fmt.Errorf("line %d:%d: %w - pattern matching required: function '%s' expecting '%s'",
						identPos.Line, identPos.Column, ErrAnyTypeMismatch, funcName, paramName)
				}
				return fmt.Errorf("%w - pattern matching required: function '%s' expecting '%s'",
					ErrAnyTypeMismatch, funcName, paramName)
			}
		}
	}

	return nil
}

// getParameterName gets the parameter name for a function at the given index
func (g *LLVMGenerator) getParameterName(funcName string, paramIndex int) string {
	// Check if we have parameter names stored for this function
	if paramNames, exists := g.functionParameters[funcName]; exists {
		if paramIndex >= 0 && paramIndex < len(paramNames) {
			return paramNames[paramIndex]
		}
	}
	return ""
}

// handleBuiltInFunction handles all built-in function calls using the unified registry.
func (g *LLVMGenerator) handleBuiltInFunction(name string, callExpr *ast.CallExpression) (value.Value, error) {
	// Check if this is a built-in function using the registry
	fn, exists := GlobalBuiltInRegistry.GetFunction(name)
	if !exists {
		return nil, nil // Not a built-in function
	}

	// Check security permissions
	if !g.checkSecurityPermission(fn.SecurityFlag) {
		return nil, WrapUnsupportedCallExpressionSecurity(name)
	}

	// Call the function's generator
	return fn.Generator(g, callExpr)
}

// checkSecurityPermission checks if the generator's security config allows the given permission
func (g *LLVMGenerator) checkSecurityPermission(permission SecurityPermission) bool {
	switch permission {
	case PermissionNone:
		return true
	case PermissionHTTP:
		return g.security.AllowHTTP
	case PermissionWebSocket:
		return g.security.AllowWebSocket
	case PermissionFileRead:
		return g.security.AllowFileRead
	case PermissionFileWrite:
		return g.security.AllowFileWrite
	case PermissionProcess:
		return g.security.AllowProcessExecution
	case PermissionFFI:
		return g.security.AllowFFI
	default:
		return false
	}
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

	// True case - return "1"
	g.builder = trueBlock
	trueStr := constant.NewCharArrayFromString("1\x00")
	trueGlobal := g.module.NewGlobalDef("", trueStr)
	truePtr := trueBlock.NewGetElementPtr(trueStr.Typ, trueGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	trueBlock.NewBr(endBlock)

	// False case - return "0"
	g.builder = falseBlock
	falseStr := constant.NewCharArrayFromString("0\x00")
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

	for i, arm := range arms {
		// Set builder to the arm block at the start of each iteration
		g.builder = armBlocks[i]

		// Handle variable binding in patterns
		armValue, err := g.processMatchArm(arm, discriminant)
		if err != nil {
			return nil, nil, err
		}

		armValues = append(armValues, armValue)

		// Handle branch termination
		currentBuilderBlock := g.addBranchTermination(endBlock)
		predecessorBlocks = append(predecessorBlocks, currentBuilderBlock)
	}

	return armValues, predecessorBlocks, nil
}

// processMatchArm handles the processing of a single match arm
func (g *LLVMGenerator) processMatchArm(arm ast.MatchArm, discriminant value.Value) (value.Value, error) {
	if arm.Pattern.Variable != "" || len(arm.Pattern.Fields) > 0 {
		return g.processMatchArmWithBinding(arm, discriminant)
	}
	return g.processMatchArmWithoutBinding(arm)
}

// processMatchArmWithBinding handles match arms that have variable binding
func (g *LLVMGenerator) processMatchArmWithBinding(arm ast.MatchArm, discriminant value.Value) (value.Value, error) {
	// Save the current variable scope
	oldVariables := g.saveVariableScope()

	// Bind the pattern variable to the discriminant value
	if arm.Pattern.Variable != "" {
		g.variables[arm.Pattern.Variable] = discriminant
	}

	// Handle field extraction for structural matching and discriminated unions
	if len(arm.Pattern.Fields) > 0 {
		g.extractPatternFields(arm.Pattern, discriminant)
	}

	// Generate the arm expression
	armValue, err := g.generateExpression(arm.Expression)
	if err != nil {
		g.restoreVariableScope(oldVariables)
		return nil, err
	}

	// Restore the previous variable scope
	g.restoreVariableScope(oldVariables)

	return g.normalizeArmValue(armValue), nil
}

// processMatchArmWithoutBinding handles match arms without variable binding
func (g *LLVMGenerator) processMatchArmWithoutBinding(arm ast.MatchArm) (value.Value, error) {
	armValue, err := g.generateExpression(arm.Expression)
	if err != nil {
		return nil, err
	}

	return g.normalizeArmValue(armValue), nil
}

// saveVariableScope saves the current variable scope
func (g *LLVMGenerator) saveVariableScope() map[string]value.Value {
	oldVariables := make(map[string]value.Value)
	for k, v := range g.variables {
		oldVariables[k] = v
	}
	return oldVariables
}

// restoreVariableScope restores a saved variable scope
func (g *LLVMGenerator) restoreVariableScope(oldVariables map[string]value.Value) {
	g.variables = oldVariables
}

// extractPatternFields handles field extraction for patterns
func (g *LLVMGenerator) extractPatternFields(pattern ast.Pattern, discriminant value.Value) {
	if pattern.Constructor == "*" {
		// For structural matching, extract fields from the object
		g.extractStructuralFields(pattern.Fields, discriminant)
	} else if _, exists := g.unionVariants[pattern.Constructor]; exists {
		// FIXED: Handle discriminated union field extraction
		err := g.extractDiscriminatedUnionFields(discriminant, pattern, g.variables)
		if err != nil {
			// If field extraction fails, bind fields to zero values
			g.bindFieldsToZeroValues(pattern.Fields)
		}
	}
}

// extractStructuralFields extracts fields for structural matching
func (g *LLVMGenerator) extractStructuralFields(fields []string, discriminant value.Value) {
	for _, fieldName := range fields {
		// Extract the actual field value from the object
		fieldValue, err := g.extractFieldFromObject(discriminant, fieldName)
		if err != nil {
			// If field extraction fails, bind to null/zero value
			fieldValue = constant.NewNull(types.I8Ptr)
		}
		g.variables[fieldName] = fieldValue
	}
}

// bindFieldsToZeroValues binds pattern fields to zero values when extraction fails
func (g *LLVMGenerator) bindFieldsToZeroValues(fields []string) {
	for _, fieldName := range fields {
		g.variables[fieldName] = constant.NewInt(types.I64, 0)
	}
}

// normalizeArmValue handles Unit expressions in match arms
func (g *LLVMGenerator) normalizeArmValue(armValue value.Value) value.Value {
	// CRITICAL FIX: Handle Unit expressions in match arms
	// If the expression returns Unit, use void type
	if armValue == nil || armValue.Type() == types.Void {
		armValue = constant.NewUndef(types.Void)
	}
	return armValue
}

// addBranchTermination ensures proper branch termination for match arms
func (g *LLVMGenerator) addBranchTermination(endBlock *ir.Block) *ir.Block {
	// TODO: FIX THIS! DON'T IGNORE IT!!
	// CRITICAL FIX: After generating the expression (which might be a nested match),
	// the builder might be pointing to a different block. We need to ensure the
	// branch comes from the current builder block (where the expression ended),
	// but ONLY if that block doesn't already have a terminator.
	currentBuilderBlock := g.builder

	// Check if the current block already has a terminator instruction
	if currentBuilderBlock.Term == nil {
		currentBuilderBlock.NewBr(endBlock)
	}
	return currentBuilderBlock
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

	// Handle structural matching (any constructor with field extraction)
	if pattern.Constructor == "*" {
		// For structural matching, we assume the pattern always matches for now
		// Real field checking would require runtime type information
		// But we need to extract and bind the field values when this pattern is used
		return constant.NewBool(true)
	}

	// Handle variable binding patterns (empty constructor means variable binding)
	if pattern.Constructor == "" && pattern.Variable != "" {
		return constant.NewBool(true)
	}

	// Handle boolean patterns for ternary expressions
	if pattern.Constructor == "true" {
		// Check if discriminant is true (non-zero)
		zero := constant.NewInt(types.I64, 0)
		return currentBlock.NewICmp(enum.IPredNE, discriminant, zero)
	}

	if pattern.Constructor == "false" {
		// Check if discriminant is false (zero)
		zero := constant.NewInt(types.I64, 0)
		return currentBlock.NewICmp(enum.IPredEQ, discriminant, zero)
	}

	// Check if it's a union type variant
	if discriminantValue, exists := g.unionVariants[pattern.Constructor]; exists {
		// FIXED: Handle both simple enums and discriminated unions
		discriminantType := discriminant.Type()

		// Check if this is a tagged union (pointer to struct with tag + data)
		if ptrType, ok := discriminantType.(*types.PointerType); ok {
			if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == 2 {
				// This is a discriminated union - extract tag field (index 0)
				tagPtr := currentBlock.NewGetElementPtr(structType, discriminant,
					constant.NewInt(types.I32, 0), // struct index
					constant.NewInt(types.I32, 0)) // tag field index
				tagValue := currentBlock.NewLoad(types.I8, tagPtr)

				// Convert discriminant value to i8 for comparison
				patternConst := constant.NewInt(types.I8, discriminantValue)
				return currentBlock.NewICmp(enum.IPredEQ, tagValue, patternConst)
			}
		}

		// Fallback: simple enum discriminant (i64)
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

// extractFieldFromObject extracts a field value from an object literal
func (g *LLVMGenerator) extractFieldFromObject(objectValue value.Value, fieldName string) (value.Value, error) {
	// Check if this is a pointer to a struct
	objectType := objectValue.Type()
	var structType *types.StructType
	var isPointer bool

	if ptrType, ok := objectType.(*types.PointerType); ok {
		if st, ok := ptrType.ElemType.(*types.StructType); ok {
			structType = st
			isPointer = true
		}
	} else if st, ok := objectType.(*types.StructType); ok {
		structType = st
		isPointer = false
	}

	if structType == nil {
		// If not a struct, we can't extract fields
		return nil, fmt.Errorf("%w: %s", ErrCannotExtractField, fieldName)
	}

	// For now, assume the first field contains the value we want
	// In a real implementation, we'd need to map field names to indices
	fieldIndex := 0

	// Get pointer to the field
	var fieldPtr value.Value
	if isPointer {
		// Object is already a pointer to the struct
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			objectValue,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	} else {
		// Object is a struct value, need to get its address first
		structAddr := g.builder.NewAlloca(structType)
		g.builder.NewStore(objectValue, structAddr)
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			structAddr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	}

	// Load the field value
	fieldType := structType.Fields[fieldIndex]
	fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

	return fieldValue, nil
}

// extractDiscriminatedUnionFields extracts field values from a discriminated union variant
func (g *LLVMGenerator) extractDiscriminatedUnionFields(
	discriminant value.Value,
	pattern ast.Pattern,
	variables map[string]value.Value,
) error {
	// Check if this is a tagged union (pointer to struct with tag + data)
	discriminantType := discriminant.Type()
	ptrType, ok := discriminantType.(*types.PointerType)
	if !ok {
		return ErrDiscriminantNotPointer
	}

	structType, ok := ptrType.ElemType.(*types.StructType)
	if !ok || len(structType.Fields) != 2 {
		return ErrDiscriminantNotTaggedUnion
	}

	// Get pointer to the data area (second field in the tagged union)
	dataPtr := g.builder.NewGetElementPtr(structType, discriminant,
		constant.NewInt(types.I32, 0), // struct index
		constant.NewInt(types.I32, 1)) // data field index

	// Find the union type and variant information
	var variant *ast.TypeVariant

	// Find which union type this discriminant belongs to
	for _, typeDecl := range g.typeDeclarations {
		if len(typeDecl.Variants) > 1 {
			// Check if this pattern constructor matches any variant
			for _, v := range typeDecl.Variants {
				if v.Name == pattern.Constructor {
					variant = &v
					break
				}
			}
			if variant != nil {
				break
			}
		}
	}

	if variant == nil {
		return fmt.Errorf("%w: %s", ErrVariantNotFound, pattern.Constructor)
	}

	// Extract each field from the data area
	offset := int64(0)
	for _, field := range variant.Fields {
		// Check if this field is requested in the pattern
		fieldRequested := false
		for _, patternField := range pattern.Fields {
			if patternField == field.Name {
				fieldRequested = true
				break
			}
		}

		if fieldRequested {
			// Get the field type and calculate its size
			fieldType := g.getFieldType(field.Type)

			// Cast data array to appropriate pointer type for this field
			fieldPtr := g.builder.NewBitCast(
				g.builder.NewGetElementPtr(
					types.NewArray(uint64(LargeArraySizeForCasting), types.I8), // Use large array for casting
					dataPtr,
					constant.NewInt(types.I32, 0),      // array index
					constant.NewInt(types.I32, offset), // byte offset
				),
				types.NewPointer(fieldType),
			)

			// Load the field value
			fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

			// Bind the field value to the pattern variable
			variables[field.Name] = fieldValue
			// Type tracking is now handled by Hindley-Milner inference
		}

		// Move to next field offset regardless of whether it was requested
		fieldSize := g.getTypeSize(g.getFieldType(field.Type))
		offset += fieldSize
	}

	return nil
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
		// For single arm, we still need to set the builder but don't need PHI
		return armValues[0], nil
	}

	// CRITICAL FIX: Check if all arm values are void type
	// PHI nodes cannot be created with void values
	allVoid := true
	for _, val := range armValues {
		if val.Type() != types.Void {
			allVoid = false
			break
		}
	}

	// If all arms return void, return void directly without PHI
	if allVoid {
		return constant.NewUndef(types.Void), nil
	}

	// Check if we need type coercion
	coercedValues, err := g.coerceArmValuesToCommonType(armValues)
	if err != nil {
		return nil, err
	}

	// CRITICAL FIX: Only include predecessors that actually have terminators
	var validIncomings []*ir.Incoming
	for i, val := range coercedValues {
		// Skip void values in PHI nodes
		if val.Type() != types.Void {
			// Check if the predecessor block has a terminator (meaning it actually branches)
			if i < len(predecessorBlocks) && predecessorBlocks[i].Term != nil {
				validIncomings = append(validIncomings, ir.NewIncoming(val, predecessorBlocks[i]))
			}
		}
	}

	// If no valid (non-void) values for PHI, return void
	if len(validIncomings) == 0 {
		return constant.NewUndef(types.Void), nil
	}

	// If only one valid incoming, don't create PHI
	if len(validIncomings) == 1 {
		return validIncomings[0].X, nil
	}

	phi := endBlock.NewPhi(validIncomings...)

	// The end block now has a PHI node and the builder is set to this block.
	// The calling function (like generateStandardMatchExpression) should handle
	// adding any necessary terminator when the match is used in a larger context.

	return phi, nil
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

	// Track which blocks actually branch to the end
	var actualSuccessBlock *ir.Block
	var actualErrorBlock *ir.Block

	// Generate success block and track the actual ending block
	g.builder = blocks.Success
	successValue, err := g.generateSuccessBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}
	// The builder is now pointing to the block that will branch to the end
	actualSuccessBlock = g.builder

	// Generate error block and track the actual ending block
	g.builder = blocks.Error
	errorValue, err := g.generateErrorBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}
	// The builder is now pointing to the block that will branch to the end
	actualErrorBlock = g.builder

	// Create PHI with the actual predecessor blocks
	return g.createResultMatchPhiWithActualBlocks(successValue, errorValue, actualSuccessBlock, actualErrorBlock, blocks)
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
	if successExpr != nil {
		// Generate the expression (which might be a nested match)
		val, err := g.generateExpression(successExpr)
		if err != nil {
			return nil, err
		}
		successValue = val

		// CRITICAL FIX: After generating a nested expression, the builder might have changed
		// We need to ensure the branch to the end block comes from the correct block
		// But only add the branch if the current block doesn't already have a terminator
		if g.builder.Term == nil {
			g.builder.NewBr(blocks.End)
		}

		// For PHI node creation, we need to track which block actually branches to the end
		// This might be different from the original success block if we had nested expressions
		return successValue, nil
	}

	// Fallback: use the bound variable from pattern matching
	if successArm := g.findSuccessArm(matchExpr); successArm != nil && len(successArm.Pattern.Fields) > 0 {
		fieldName := successArm.Pattern.Fields[0]
		if extractedValue, exists := g.variables[fieldName]; exists {
			successValue = extractedValue
		} else {
			successValue = constant.NewInt(types.I64, ArrayIndexZero)
		}
	} else {
		successValue = constant.NewInt(types.I64, ArrayIndexZero)
	}

	// Only add branch if the current block doesn't already have a terminator
	if g.builder.Term == nil {
		g.builder.NewBr(blocks.End)
	}

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
	if errorExpr != nil {
		// Generate the expression (which might be a nested match)
		val, err := g.generateExpression(errorExpr)
		if err != nil {
			return nil, err
		}
		errorValue = val

		// CRITICAL FIX: After generating a nested expression, the builder might have changed
		// We need to ensure the branch to the end block comes from the correct block
		// But only add the branch if the current block doesn't already have a terminator
		if g.builder.Term == nil {
			g.builder.NewBr(blocks.End)
		}

		return errorValue, nil
	}

	// Fallback: use the bound variable from pattern matching
	if errorArm := g.findErrorArm(matchExpr); errorArm != nil && len(errorArm.Pattern.Fields) > 0 {
		fieldName := errorArm.Pattern.Fields[0]
		if extractedError, exists := g.variables[fieldName]; exists {
			errorValue = extractedError
		} else {
			errorValue = constant.NewInt(types.I64, ArrayIndexZero)
		}
	} else {
		errorValue = constant.NewInt(types.I64, ArrayIndexZero)
	}

	// Only add branch if the current block doesn't already have a terminator
	if g.builder.Term == nil {
		g.builder.NewBr(blocks.End)
	}

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

// createResultMatchPhiWithActualBlocks creates the PHI node for result matching with actual predecessor blocks.
func (g *LLVMGenerator) createResultMatchPhiWithActualBlocks(
	successValue, errorValue value.Value,
	actualSuccessBlock, actualErrorBlock *ir.Block,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.End

	// CRITICAL FIX: Use the actual blocks that branch to the end
	var validPredecessors []*ir.Incoming

	// Check if the actual success block has a terminator and branches to end
	if actualSuccessBlock != nil && actualSuccessBlock.Term != nil {
		validPredecessors = append(validPredecessors, ir.NewIncoming(successValue, actualSuccessBlock))
	}

	// Check if the actual error block has a terminator and branches to end
	if actualErrorBlock != nil && actualErrorBlock.Term != nil {
		validPredecessors = append(validPredecessors, ir.NewIncoming(errorValue, actualErrorBlock))
	}

	// If we don't have valid predecessors, return a default value
	if len(validPredecessors) == 0 {
		// Return the success value as a fallback
		return successValue, nil
	}

	// If we only have one valid predecessor, don't create a PHI
	if len(validPredecessors) == 1 {
		return validPredecessors[0].X, nil
	}

	// Create PHI node with valid predecessors
	phi := blocks.End.NewPhi(validPredecessors...)

	return phi, nil
}
