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
			var source string
			if tt.operator == "!" {
				source = "fn test() -> bool = " + tt.operator + tt.operand
			} else {
				source = "fn test() -> int = " + tt.operator + tt.operand
			}

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

func TestFieldAccessExpression(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		wantErr  bool
		contains string
	}{
		{
			name: "valid field access on struct",
			source: `
				type Point = { x: int, y: int }
				let point = Point { x: 10, y: 20 }
				let result = point.x
			`,
			wantErr: false,
		},
		{
			name: "invalid field access on integer",
			source: `
				let x = 42
				let result = x.value
			`,
			wantErr:  true,
			contains: "cannot access field 'value' on non-struct type",
		},
		{
			name: "invalid field access on string",
			source: `
				let s = "hello"
				let result = s.length
			`,
			wantErr:  true,
			contains: "cannot access field 'length' on non-struct type",
		},
		{
			name: "valid field access in expression",
			source: `
				type Rectangle = { width: int, height: int }
				let rect = Rectangle { width: 10, height: 5 }
				let result = rect.width + rect.height
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
					t.Errorf("Expected error to contain '%s', but got: %v", tt.contains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
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
			errMsg:  "print expects exactly 1 arguments (value), got 0",
		},
		{
			name:    "print too many args",
			source:  `print(1, 2)`,
			wantErr: true,
			errMsg:  "print expects exactly 1 arguments (value), got 2",
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
			errMsg:  "input expects exactly 0 arguments, got 1",
		},
		{
			name:    "input too many args",
			source:  `input("prompt", "extra")`,
			wantErr: true,
			errMsg:  "input expects exactly 0 arguments, got 2",
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
			source:  `fn test() -> bool = 1 < 2`,
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
