package codegen_test

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestUnaryExpressionGeneration(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		operand  string
		wantErr  bool
	}{
		{"unary plus", "+", "42", false},
		{"unary minus", "-", "42", false},
		{"boolean not", "!", "true", false},
		{"unsupported operator", "@", "42", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := "fn test() -> int = " + tt.operator + tt.operand
			_, err := codegen.CompileToLLVM(source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error for unsupported unary operator")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMethodCallExpression(t *testing.T) {
	// Test method call expressions (should fail with WrapMethodNotImpl)
	source := `
		let obj = 42
		obj.toString()
	`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Error("Expected error for method call")
	}

	if !strings.Contains(err.Error(), "method call not implemented") {
		t.Error("Expected method call error message")
	}
}

func TestFieldAccessExpression(_ *testing.T) {

	TODO: put a test here that verifies that field access is actually working

}

func TestToStringConversions(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "string to string",
			source:  `print(toString("hello"))`,
			wantErr: false,
		},
		{
			name:    "int to string",
			source:  `print(toString(42))`,
			wantErr: false,
		},
		{
			name:    "bool to string",
			source:  `print(toString(true))`,
			wantErr: false,
		},
		{
			name:    "wrong arg count",
			source:  `toString()`,
			wantErr: true,
			errMsg:  "function toString expects 1 arguments, got 0",
		},
		{
			name:    "too many args",
			source:  `toString(1, 2)`,
			wantErr: true,
			errMsg:  "function toString expects 1 arguments, got 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPrintExpressionTypes(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "print string literal",
			source:  `print("hello")`,
			wantErr: false,
		},
		{
			name:    "print integer",
			source:  `print(42)`,
			wantErr: false,
		},
		{
			name:    "print boolean",
			source:  `print(true)`,
			wantErr: false,
		},
		{
			name:    "print binary expression",
			source:  `print(1 + 2)`,
			wantErr: false,
		},
		{
			name: "print identifier",
			source: `let x = 42
print(x)`,
			wantErr: false,
		},
		{
			name:    "print wrong args",
			source:  `print()`,
			wantErr: true,
			errMsg:  "print expects exactly 1 argument(s), got 0",
		},
		{
			name:    "print too many args",
			source:  `print(1, 2)`,
			wantErr: true,
			errMsg:  "print expects exactly 1 argument(s), got 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestInputFunction(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "input with wrong args",
			source:  `input(42)`,
			wantErr: true,
			errMsg:  "input expects exactly 0 argument(s), got 1",
		},
		{
			name:    "input too many args",
			source:  `input("prompt", "extra")`,
			wantErr: true,
			errMsg:  "input expects exactly 0 argument(s), got 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUnsupportedExpressions(t *testing.T) {
	// Test unsupported expression types that hit WrapUnsupportedExpression
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "unsupported expression in print",
			source: `print(someUnsupportedExpr)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			// These should generate errors due to undefined variables or unsupported expressions
			if err == nil {
				t.Error("Expected error for unsupported expression")
			}
		})
	}
}

func TestResultExpressions(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name:    "successful result expression",
			source:  `let x = (1 + 2)`,
			wantErr: false,
		},
		{
			name:    "result in print",
			source:  `print((1 + 2))`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr && err == nil {
				t.Error("Expected error")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestBinaryOperatorErrors(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid arithmetic",
			source:  `fn test() -> int = 1 + 2`,
			wantErr: false,
		},
		{
			name:    "valid comparison",
			source:  `fn test() -> int = 1 < 2`,
			wantErr: false,
		},
		{
			name:    "division",
			source:  `fn test() -> int = 10 / 2`,
			wantErr: false,
		},
		{
			name:    "modulo",
			source:  `fn test() -> int = 10 % 3`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
