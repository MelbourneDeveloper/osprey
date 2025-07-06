package ast

// PerformExpression represents effect operations like 'perform Effect.operation(args)'
type PerformExpression struct {
	EffectName    string
	OperationName string
	Arguments     []Expression
	Position      *Position
}

func (pe *PerformExpression) isExpression() {}

// HandlerExpression represents effect handlers that catch and handle effects
type HandlerExpression struct {
	EffectName string
	Handlers   []HandlerArm
	Body       Expression
	Position   *Position
}

func (he *HandlerExpression) isExpression() {}

// HandlerArm represents individual effect operation handlers
type HandlerArm struct {
	OperationName string
	Parameters    []string
	Body          Expression
	Position      *Position
}
