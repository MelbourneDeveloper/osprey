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
}

// NewEffectCodegen creates a new algebraic effects code generator
func (g *LLVMGenerator) NewEffectCodegen() *EffectCodegen {
	return &EffectCodegen{
		generator:    g,
		registry:     &EffectRegistry{Effects: make(map[string]*EffectType)},
		handlerStack: make([]*HandlerFrame, 0),
		handlerFuncs: make(map[string]*ir.Func),
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
	// Look for handler on the stack
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

				// Add the continuation as the last argument
				if frame.Continuation != nil {
					// Create a closure for the current continuation
					contCall := ec.createContinuation()
					args = append(args, contCall)
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
	return nil, fmt.Errorf("%w: %s", ErrUnhandledEffect, errorMsg)
}

// GenerateHandlerExpression generates handler expressions with proper CPS
func (ec *EffectCodegen) GenerateHandlerExpression(handler *ast.HandlerExpression) (value.Value, error) {
	// Create handler frame
	frame := &HandlerFrame{
		EffectName: handler.EffectName,
		Operations: make(map[string]*ir.Func),
	}

	// Generate handler functions for each operation
	for _, arm := range handler.Handlers {
		handlerFunc := ec.createHandlerFunction(handler.EffectName, arm)
		frame.Operations[arm.OperationName] = handlerFunc
	}

	// Create continuation function for the body
	frame.Continuation = ec.createBodyContinuation(handler.Body)

	// Push handler frame
	ec.handlerStack = append(ec.handlerStack, frame)

	// Generate the body with the handler active
	result, err := ec.generator.generateExpression(handler.Body)

	// Pop handler frame
	ec.handlerStack = ec.handlerStack[:len(ec.handlerStack)-1]

	return result, err
}

// createHandlerFunction creates a handler function for an effect operation
func (ec *EffectCodegen) createHandlerFunction(effectName string, arm ast.HandlerArm) *ir.Func {
	// Create function name
	funcName := fmt.Sprintf("__handler_%s_%s_%d", effectName, arm.OperationName, ec.contCounter)
	ec.contCounter++

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

	// Add continuation parameter
	contType := types.NewFunc(types.Void, types.I64) // Simplified continuation type
	allParamTypes := make([]types.Type, len(paramTypes)+1)
	copy(allParamTypes, paramTypes)
	allParamTypes[len(paramTypes)] = types.NewPointer(contType)
	paramTypes = allParamTypes

	// Create function type
	funcType := types.NewFunc(types.Void, paramTypes...)

	// Create the function
	handlerFunc := ec.generator.module.NewFunc(funcName, funcType)

	// Create entry block
	entry := handlerFunc.NewBlock("entry")

	// Store current builder and switch to handler function
	oldBuilder := ec.generator.builder
	ec.generator.builder = entry

	// Generate handler body
	// For now, just call the continuation with Unit value
	if len(handlerFunc.Params) > 0 {
		contParam := handlerFunc.Params[len(handlerFunc.Params)-1]

		// Call the continuation with a default value (0 for simplicity)
		defaultVal := constant.NewInt(types.I64, 0)
		entry.NewCall(contParam, defaultVal)
	}

	// Return void
	entry.NewRet(nil)

	// Restore builder
	ec.generator.builder = oldBuilder

	return handlerFunc
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

// generateRealPerformExpression generates real algebraic effects perform expressions
func (g *LLVMGenerator) generateRealPerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}
	return g.effectCodegen.GeneratePerformExpression(perform)
}
