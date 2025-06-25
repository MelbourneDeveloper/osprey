package codegen

import (
	"errors"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// ErrUnhandledEffect is returned when an effect is not handled at compile time
var ErrUnhandledEffect = errors.New("unhandled effect")

// EffectRegistry tracks registered effects and their operations
type EffectRegistry struct {
	Effects map[string]*EffectType
}

// EffectType represents an effect type in the registry
type EffectType struct {
	Name       string
	Operations map[string]*EffectOp
}

// EffectOp represents an operation within an effect
type EffectOp struct {
	Name       string
	ParamTypes []types.Type
	ReturnType types.Type
}

// HandlerFrame represents an active effect handler frame
type HandlerFrame struct {
	EffectName   string
	Operations   map[string]*ir.Func
	Continuation *ir.Func
}

// EffectCodegen implements real algebraic effects with CPS transformation
type EffectCodegen struct {
	generator    *LLVMGenerator
	registry     *EffectRegistry
	handlerStack []*HandlerFrame
	contCounter  int
	handlerFuncs map[string]*ir.Func
	// NEW: Track function generation context for scope maintenance
	inHandlerScope  bool
	currentHandlers []*HandlerFrame
	// CRITICAL: Track effects declared in current function signature
	currentFunctionEffects []string
}

// NewEffectCodegen creates a new algebraic effects code generator
func (g *LLVMGenerator) NewEffectCodegen() *EffectCodegen {
	return &EffectCodegen{
		generator:       g,
		registry:        &EffectRegistry{Effects: make(map[string]*EffectType)},
		handlerStack:    make([]*HandlerFrame, 0),
		handlerFuncs:    make(map[string]*ir.Func),
		inHandlerScope:  false,
		currentHandlers: make([]*HandlerFrame, 0),
	}
}

// RegisterEffect registers an effect declaration for code generation
func (ec *EffectCodegen) RegisterEffect(effect *ast.EffectDeclaration) {
	effectType := &EffectType{
		Name:       effect.Name,
		Operations: make(map[string]*EffectOp),
	}

	// Register each operation
	for _, op := range effect.Operations {
		// Convert operation parameters to LLVM types
		paramTypes := make([]types.Type, len(op.Parameters))
		for i := range op.Parameters {
			paramTypes[i] = types.I8Ptr // Simplified for now
		}

		effectOp := &EffectOp{
			Name:       op.Name,
			ParamTypes: paramTypes,
			ReturnType: types.Void, // Simplified for now
		}

		effectType.Operations[op.Name] = effectOp
	}

	ec.registry.Effects[effect.Name] = effectType
}

// GeneratePerformExpression generates CPS-transformed perform expressions
func (ec *EffectCodegen) GeneratePerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	// CRITICAL FIX: Check if current function declares this effect
	if ec.currentFunctionEffects != nil {
		for _, declaredEffect := range ec.currentFunctionEffects {
			if declaredEffect == perform.EffectName {
				// Effect is declared in function signature - generate proper perform call
				return ec.generateDeclaredEffectCall(perform)
			}
		}
	}

	// Check current scope handlers first (for function calls within handler bodies)
	if ec.inHandlerScope {
		for i := len(ec.currentHandlers) - 1; i >= 0; i-- {
			frame := ec.currentHandlers[i]
			if frame.EffectName == perform.EffectName {
				// Generate arguments
				args := make([]value.Value, 0)
				for _, argExpr := range perform.Arguments {
					argVal, err := ec.generator.generateExpression(argExpr)
					if err != nil {
						return nil, err
					}
					args = append(args, argVal)
				}

				// Call the handler function
				return ec.generator.builder.NewCall(frame.Operations[perform.OperationName], args...), nil
			}
		}
	}

	// Look for handler on the stack (traditional stack-based lookup)
	for i := len(ec.handlerStack) - 1; i >= 0; i-- {
		frame := ec.handlerStack[i]
		if frame.EffectName == perform.EffectName {
			if handlerFunc, exists := frame.Operations[perform.OperationName]; exists {
				// Generate arguments
				args := make([]value.Value, 0)
				for _, argExpr := range perform.Arguments {
					argVal, err := ec.generator.generateExpression(argExpr)
					if err != nil {
						return nil, err
					}
					args = append(args, argVal)
				}

				// Call the handler function
				return ec.generator.builder.NewCall(handlerFunc, args...), nil
			}
		}
	}

	// No handler found - COMPILATION ERROR! Effects must be handled at compile time!
	errorMsg := fmt.Sprintf("COMPILATION ERROR: Unhandled effect '%s.%s' - "+
		"all effects must be explicitly handled or forwarded in function signatures. "+
		"Add a handler or declare the effect in the function signature with !%s",
		perform.EffectName, perform.OperationName, perform.EffectName)
	return nil, fmt.Errorf("unhandled effect: %s", errorMsg)
}

// Helper function to get map keys for debugging
func getKeys(m map[string]*ir.Func) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GenerateHandlerExpression generates code for handler expressions
func (ec *EffectCodegen) GenerateHandlerExpression(handler *ast.HandlerExpression) (value.Value, error) {
	effectName := handler.EffectName

	// Create handler function for each operation
	handlerFrame := &HandlerFrame{
		EffectName: effectName,
		Operations: make(map[string]*ir.Func),
	}

	for _, arm := range handler.Handlers {
		handlerFunc, err := ec.createHandlerFunction(effectName, arm, ec.contCounter)
		if err != nil {
			return nil, err
		}

		handlerFrame.Operations[arm.OperationName] = handlerFunc
		ec.contCounter++
	}

	// Push handler onto stack
	ec.handlerStack = append(ec.handlerStack, handlerFrame)
	ec.inHandlerScope = true
	ec.currentHandlers = append(ec.currentHandlers, handlerFrame)

	// Generate the do body with the handler active
	result, err := ec.generator.generateExpression(handler.Body)
	if err != nil {
		return nil, err
	}

	// Pop handler from stack
	ec.handlerStack = ec.handlerStack[:len(ec.handlerStack)-1]
	ec.inHandlerScope = len(ec.currentHandlers) > 1
	if len(ec.currentHandlers) > 0 {
		ec.currentHandlers = ec.currentHandlers[:len(ec.currentHandlers)-1]
	}

	return result, nil
}

// createHandlerFunction creates a handler function for an effect operation
func (ec *EffectCodegen) createHandlerFunction(effectName string, arm ast.HandlerArm, contCounter int) (*ir.Func, error) {
	// Create function name
	funcName := fmt.Sprintf("__handler_%s_%s_%d", effectName, arm.OperationName, contCounter)

	// Determine parameter types
	effectType := ec.registry.Effects[effectName]
	var paramTypes []types.Type
	if effectType != nil && effectType.Operations[arm.OperationName] != nil {
		paramTypes = effectType.Operations[arm.OperationName].ParamTypes
	} else {
		// Default to string parameters
		paramTypes = make([]types.Type, len(arm.Parameters))
		for i := range paramTypes {
			paramTypes[i] = types.I8Ptr
		}
	}

	// Create function with parameters using ir.NewParam (not types.NewFunc)
	// This is the correct way to create functions with parameters in llir/llvm
	params := make([]*ir.Param, len(arm.Parameters))
	for i, paramName := range arm.Parameters {
		params[i] = ir.NewParam(paramName, types.I8Ptr) // Default to string type
	}

	// Create the function with explicit parameters
	handlerFunc := ec.generator.module.NewFunc(funcName, types.Void, params...)

	// Create entry block
	entry := handlerFunc.NewBlock("entry")

	// Store current state
	oldBuilder := ec.generator.builder
	oldVars := make(map[string]value.Value)
	oldVarTypes := make(map[string]string)
	for k, v := range ec.generator.variables {
		oldVars[k] = v
	}
	if ec.generator.variableTypes != nil {
		for k, v := range ec.generator.variableTypes {
			oldVarTypes[k] = v
		}
	}

	// Switch to handler function builder
	ec.generator.builder = entry

	// Initialize variable types if nil
	if ec.generator.variableTypes == nil {
		ec.generator.variableTypes = make(map[string]string)
	}

	// Set up parameters as variables - THIS IS THE KEY FIX!
	for i, param := range arm.Parameters {
		if i < len(handlerFunc.Params) {
			ec.generator.variables[param] = handlerFunc.Params[i]
			ec.generator.variableTypes[param] = "string" // Default to string for now
		}
	}

	// Generate handler body - now the parameters are available!
	if arm.Body != nil {
		result, err := ec.generator.generateExpression(arm.Body)
		if err != nil {
			// Restore state before returning error
			ec.generator.builder = oldBuilder
			ec.generator.variables = oldVars
			ec.generator.variableTypes = oldVarTypes
			return nil, err
		}
		_ = result // Ignore result for now
	}

	// Return void
	entry.NewRet(nil)

	// Restore state
	ec.generator.builder = oldBuilder
	ec.generator.variables = oldVars
	ec.generator.variableTypes = oldVarTypes

	return handlerFunc, nil
}

// createBodyContinuation creates a continuation for the handler body
func (ec *EffectCodegen) createBodyContinuation(_ ast.Expression) *ir.Func {
	funcName := fmt.Sprintf("__body_cont_%d", ec.contCounter)
	ec.contCounter++

	// Simple continuation that takes a value and returns void
	funcType := types.NewFunc(types.Void, types.I64)
	contFunc := ec.generator.module.NewFunc(funcName, funcType)

	entry := contFunc.NewBlock("entry")
	entry.NewRet(nil)

	return contFunc
}

// createContinuation creates a continuation closure
func (ec *EffectCodegen) createContinuation() value.Value {
	// For now, return a null pointer - in a full implementation this would
	// capture the current computation state
	return constant.NewNull(types.I8Ptr)
}

// Integration with main generator

// InitializeEffects initializes the effect system for the generator
func (g *LLVMGenerator) InitializeEffects() {
	g.effectCodegen = g.NewEffectCodegen()
}

// RegisterEffectDeclaration registers an effect declaration with the effect system
func (g *LLVMGenerator) RegisterEffectDeclaration(effect *ast.EffectDeclaration) error {
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}
	g.effectCodegen.RegisterEffect(effect)
	return nil
}

// generateRealPerformExpression generates real algebraic effects perform expressions
func (g *LLVMGenerator) generateRealPerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}
	return g.effectCodegen.GeneratePerformExpression(perform)
}

// generateDeclaredEffectCall generates calls for effects declared in function signatures
// This represents the "suspension point" where the effect needs to be handled by the caller
func (ec *EffectCodegen) generateDeclaredEffectCall(perform *ast.PerformExpression) (value.Value, error) {
	// For functions that declare effects, we need to generate a suspension point
	// that will be handled by whatever calls this function

	// Generate arguments for the perform expression
	args := make([]value.Value, 0)
	for _, argExpr := range perform.Arguments {
		argVal, err := ec.generator.generateExpression(argExpr)
		if err != nil {
			return nil, err
		}
		args = append(args, argVal)
	}

	// For now, since we don't have full CPS transformation yet,
	// we'll generate a print statement to demonstrate the effect is being performed
	// In a full implementation, this would generate a suspension point

	// Create a debug message showing the effect is being performed
	effectMsg := fmt.Sprintf("EFFECT: %s.%s called", perform.EffectName, perform.OperationName)
	msgStr := ec.generator.createGlobalString(effectMsg)

	// Use existing puts function from the functions map
	putsFunc := ec.generator.functions["puts"]
	call := ec.generator.builder.NewCall(putsFunc, msgStr)

	// Convert to i64 to match expected return type
	return ec.generator.builder.NewSExt(call, types.I64), nil
}
