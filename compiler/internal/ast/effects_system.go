package ast

import "fmt"

// EffectType represents an effect in the type system
type EffectType struct {
	Name       string
	Operations map[string]*EffectOpDef
}

// EffectOpDef represents a single operation definition within an effect
type EffectOpDef struct {
	Name       string
	Parameters []Parameter
	ReturnType *TypeExpression
}

// EffectSet represents a set of effects that a function can perform
type EffectSet struct {
	Effects []*EffectType
}

// ContinuationType represents the type of a continuation
type ContinuationType struct {
	ParamType  *TypeExpression
	ReturnType *TypeExpression
	EffectSet  *EffectSet
}

// EffectSignature represents the full signature of an effectful function
type EffectSignature struct {
	Parameters []Parameter
	ReturnType *TypeExpression
	Effects    *EffectSet
}

// ResumeExpression represents a resume operation in a handler
type ResumeExpression struct {
	Value Expression
}

// HandlerFrame represents a handler frame on the handler stack
type HandlerFrame struct {
	EffectName string
	Handlers   map[string]*HandlerClause
	Parent     *HandlerFrame
}

// HandlerClause represents a single handler clause
type HandlerClause struct {
	OperationName string
	Parameters    []string
	Continuation  string // Name of the continuation parameter
	Body          Expression
}

// ContinuationExpression represents a continuation value
type ContinuationExpression struct {
	Parameters  []Parameter
	Body        Expression
	CapturedEnv map[string]*Variable // Captured environment
}

// EffectContext represents the current effect handling context
type EffectContext struct {
	HandlerStack []*HandlerFrame
	CurrentFrame *HandlerFrame
}

// Methods for effect type checking and inference

// (e *EffectSet) Contains checks if an effect is in the set
func (e *EffectSet) Contains(effectName string) bool {
	for _, effect := range e.Effects {
		if effect.Name == effectName {
			return true
		}
	}
	return false
}

// (e *EffectSet) Union creates a union of two effect sets
func (e *EffectSet) Union(other *EffectSet) *EffectSet {
	result := &EffectSet{Effects: make([]*EffectType, 0)}

	// Add all effects from this set
	for _, effect := range e.Effects {
		result.Effects = append(result.Effects, effect)
	}

	// Add effects from other set that aren't already present
	for _, effect := range other.Effects {
		if !result.Contains(effect.Name) {
			result.Effects = append(result.Effects, effect)
		}
	}

	return result
}

// (e *EffectSet) String provides string representation
func (e *EffectSet) String() string {
	if len(e.Effects) == 0 {
		return "Pure"
	}

	names := make([]string, len(e.Effects))
	for i, effect := range e.Effects {
		names[i] = effect.Name
	}

	if len(names) == 1 {
		return fmt.Sprintf("!%s", names[0])
	}

	return fmt.Sprintf("![%s]", joinStrings(names, ", "))
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// EffectInference performs effect inference on expressions
type EffectInference struct {
	EffectRegistry  map[string]*EffectType
	FunctionEffects map[string]*EffectSet
}

// NewEffectInference creates a new effect inference engine
func NewEffectInference() *EffectInference {
	return &EffectInference{
		EffectRegistry:  make(map[string]*EffectType),
		FunctionEffects: make(map[string]*EffectSet),
	}
}

// RegisterEffect registers an effect in the inference engine
func (ei *EffectInference) RegisterEffect(effect *EffectType) {
	ei.EffectRegistry[effect.Name] = effect
}

// InferEffects infers the effect set of an expression
func (ei *EffectInference) InferEffects(expr Expression) *EffectSet {
	switch e := expr.(type) {
	case *PerformExpression:
		// Performing an effect adds that effect to the set
		if effect, exists := ei.EffectRegistry[e.EffectName]; exists {
			return &EffectSet{Effects: []*EffectType{effect}}
		}
		return &EffectSet{Effects: []*EffectType{}}

	case *FunctionCallExpression:
		// Function calls inherit the effects of the called function
		if effects, exists := ei.FunctionEffects[e.FunctionName]; exists {
			return effects
		}
		return &EffectSet{Effects: []*EffectType{}}

	case *BlockExpression:
		// Block expressions have the union of all statement effects
		resultSet := &EffectSet{Effects: []*EffectType{}}
		for _, stmt := range e.Statements {
			if exprStmt, ok := stmt.(*ExpressionStatement); ok {
				stmtEffects := ei.InferEffects(exprStmt.Expression)
				resultSet = resultSet.Union(stmtEffects)
			}
		}
		if e.Expression != nil {
			exprEffects := ei.InferEffects(e.Expression)
			resultSet = resultSet.Union(exprEffects)
		}
		return resultSet

	case *HandlerExpression:
		// Handler expressions have the effects of their body minus handled effects
		bodyEffects := ei.InferEffects(e.Body)
		// TODO: Subtract handled effects
		return bodyEffects

	default:
		// Pure expressions have no effects
		return &EffectSet{Effects: []*EffectType{}}
	}
}

// CPSTransformer transforms effectful code into continuation-passing style
type CPSTransformer struct {
	HandlerStack        []*HandlerFrame
	ContinuationCounter int
}

// NewCPSTransformer creates a new CPS transformer
func NewCPSTransformer() *CPSTransformer {
	return &CPSTransformer{
		HandlerStack:        make([]*HandlerFrame, 0),
		ContinuationCounter: 0,
	}
}

// TransformExpression transforms an expression to CPS form
func (cps *CPSTransformer) TransformExpression(expr Expression) Expression {
	switch e := expr.(type) {
	case *PerformExpression:
		return cps.transformPerformExpression(e)
	case *HandlerExpression:
		return cps.transformHandlerExpression(e)
	case *ResumeExpression:
		return cps.transformResumeExpression(e)
	default:
		return expr // Pure expressions don't need CPS transformation
	}
}

// transformPerformExpression transforms a perform expression
func (cps *CPSTransformer) transformPerformExpression(perform *PerformExpression) Expression {
	// Find the handler for this effect on the stack
	for i := len(cps.HandlerStack) - 1; i >= 0; i-- {
		frame := cps.HandlerStack[i]
		if frame.EffectName == perform.EffectName {
			if handler, exists := frame.Handlers[perform.OperationName]; exists {
				// Create a continuation for the rest of the computation
				contName := fmt.Sprintf("__cont_%d", cps.ContinuationCounter)
				cps.ContinuationCounter++

				// Transform the handler body with the continuation
				return &FunctionCallExpression{
					FunctionName: fmt.Sprintf("__handler_%s_%s", perform.EffectName, perform.OperationName),
					Arguments:    append(perform.Arguments, &Identifier{Name: contName}),
				}
			}
		}
	}

	// No handler found - this should be a compile error in a real system
	return perform
}

// transformHandlerExpression transforms a handler expression
func (cps *CPSTransformer) transformHandlerExpression(handler *HandlerExpression) Expression {
	// Create a new handler frame
	frame := &HandlerFrame{
		EffectName: handler.EffectName,
		Handlers:   make(map[string]*HandlerClause),
	}

	// Convert handler arms to clauses
	for _, arm := range handler.Handlers {
		clause := &HandlerClause{
			OperationName: arm.OperationName,
			Parameters:    arm.Parameters,
			Continuation:  "__k", // Standard continuation parameter name
			Body:          arm.Body,
		}
		frame.Handlers[arm.OperationName] = clause
	}

	// Push handler frame onto stack
	cps.HandlerStack = append(cps.HandlerStack, frame)

	// Transform the body with the new handler
	transformedBody := cps.TransformExpression(handler.Body)

	// Pop handler frame
	cps.HandlerStack = cps.HandlerStack[:len(cps.HandlerStack)-1]

	return transformedBody
}

// transformResumeExpression transforms a resume expression
func (cps *CPSTransformer) transformResumeExpression(resume *ResumeExpression) Expression {
	// Resume calls the continuation with the given value
	return &FunctionCallExpression{
		FunctionName: "__k", // The continuation parameter
		Arguments:    []Expression{resume.Value},
	}
}
