package ast

// EffectDeclaration represents effect declarations like 'effect ProcessOutput { captureStdout: fn(string) -> Unit }'
type EffectDeclaration struct {
	Name       string
	Operations []EffectOperation
}

func (ed *EffectDeclaration) isStatement() {}

// EffectOperation represents individual operations in an effect declaration
type EffectOperation struct {
	Name       string
	Type       string
	Parameters []Parameter
	ReturnType string
}
