package codegen

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestTypeInference(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expression
		expected string
	}{
		{
			name: "identity function",
			expr: &ast.LambdaExpression{
				Parameters: []ast.Parameter{{Name: "x"}},
				Body:       &ast.Identifier{Name: "x"},
			},
			expected: "(t0) -> t0",
		},
		{
			name: "constant function",
			expr: &ast.LambdaExpression{
				Parameters: []ast.Parameter{{Name: "x"}},
				Body:       &ast.IntegerLiteral{Value: 42},
			},
			expected: "(t0) -> " + codegen.TypeInt,
		},
		{
			name: "function application",
			expr: &ast.CallExpression{
				Function: &ast.LambdaExpression{
					Parameters: []ast.Parameter{{Name: "x"}},
					Body:       &ast.Identifier{Name: "x"},
				},
				Arguments: []ast.Expression{
					&ast.IntegerLiteral{Value: 42},
				},
			},
			expected: codegen.TypeInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inferer := codegen.NewTypeInferer()

			typ, err := inferer.InferType(tt.expr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if typ.String() != tt.expected {
				t.Errorf("expected type %q, got %q", tt.expected, typ.String())
			}
		})
	}
}

func TestUnification(t *testing.T) {
	tests := []struct {
		name          string
		t1            codegen.Type
		t2            codegen.Type
		shouldSucceed bool
	}{
		{
			name:          "same concrete types",
			t1:            codegen.NewConcreteType(codegen.TypeInt),
			t2:            codegen.NewConcreteType(codegen.TypeInt),
			shouldSucceed: true,
		},
		{
			name:          "different concrete types",
			t1:            codegen.NewConcreteType(codegen.TypeInt),
			t2:            codegen.NewConcreteType(codegen.TypeString),
			shouldSucceed: false,
		},
		{
			name: "function types same arity",
			t1: codegen.NewFunctionType(
				[]codegen.Type{codegen.NewConcreteType(codegen.TypeInt)},
				codegen.NewConcreteType(codegen.TypeInt),
			),
			t2: codegen.NewFunctionType(
				[]codegen.Type{codegen.NewConcreteType(codegen.TypeInt)},
				codegen.NewConcreteType(codegen.TypeInt),
			),
			shouldSucceed: true,
		},
		{
			name: "function types different arity",
			t1: codegen.NewFunctionType(
				[]codegen.Type{
					codegen.NewConcreteType(codegen.TypeInt),
					codegen.NewConcreteType(codegen.TypeInt),
				},
				codegen.NewConcreteType(codegen.TypeInt),
			),
			t2: codegen.NewFunctionType(
				[]codegen.Type{
					codegen.NewConcreteType(codegen.TypeInt),
				},
				codegen.NewConcreteType(codegen.TypeInt),
			),
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inferer := codegen.NewTypeInferer()

			err := inferer.Unify(tt.t1, tt.t2)
			if tt.shouldSucceed && err != nil {
				t.Errorf("expected unification to succeed, got error: %v", err)
			}

			if !tt.shouldSucceed && err == nil {
				t.Error("expected unification to fail, but it succeeded")
			}
		})
	}
}

func TestRecursiveTypes(t *testing.T) {
	inferer := codegen.NewTypeInferer()
	tv := inferer.Fresh()
	ft := codegen.NewFunctionType([]codegen.Type{tv}, tv)

	err := inferer.Unify(tv, ft)
	if err == nil {
		t.Error("expected recursive type error, got nil")
	}
}
