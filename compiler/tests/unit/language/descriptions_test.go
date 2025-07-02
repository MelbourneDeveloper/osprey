package language_test

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/language/descriptions"
)

// TestGetOperatorDescription tests the individual operator lookup function.
func TestGetOperatorDescription(t *testing.T) {
	// Test valid operator
	desc := descriptions.GetOperatorDescription("+")
	if desc == nil {
		t.Fatal("Expected description for '+' operator, got nil")
	}
	if desc.Symbol != "+" {
		t.Errorf("Expected symbol '+', got %q", desc.Symbol)
	}
	if desc.Name != "Addition" {
		t.Errorf("Expected name 'Addition', got %q", desc.Name)
	}

	// Test another valid operator
	desc = descriptions.GetOperatorDescription("==")
	if desc == nil {
		t.Fatal("Expected description for '==' operator, got nil")
	}
	if desc.Symbol != "==" {
		t.Errorf("Expected symbol '==', got %q", desc.Symbol)
	}
	if desc.Name != "Equality" {
		t.Errorf("Expected name 'Equality', got %q", desc.Name)
	}

	// Test invalid operator
	desc = descriptions.GetOperatorDescription("invalid")
	if desc != nil {
		t.Errorf("Expected nil for invalid operator, got %+v", desc)
	}
}

// TestGetBuiltinTypeDescription tests the individual type lookup function.
func TestGetBuiltinTypeDescription(t *testing.T) {
	// Test valid type
	desc := descriptions.GetBuiltinTypeDescription("Int")
	if desc == nil {
		t.Fatal("Expected description for 'Int' type, got nil")
	}
	if desc.Name != "Int" {
		t.Errorf("Expected name 'Int', got %q", desc.Name)
	}

	// Test another valid type
	desc = descriptions.GetBuiltinTypeDescription("String")
	if desc == nil {
		t.Fatal("Expected description for 'String' type, got nil")
	}
	if desc.Name != "String" {
		t.Errorf("Expected name 'String', got %q", desc.Name)
	}

	// Test invalid type
	desc = descriptions.GetBuiltinTypeDescription("invalid")
	if desc != nil {
		t.Errorf("Expected nil for invalid type, got %+v", desc)
	}
}

// TestGetBuiltinFunctionDescription tests the individual function lookup function.
func TestGetBuiltinFunctionDescription(t *testing.T) {
	// Test valid function
	desc := descriptions.GetBuiltinFunctionDescription("print")
	if desc == nil {
		t.Fatal("Expected description for 'print' function, got nil")
	}
	if desc.Name != "print" {
		t.Errorf("Expected name 'print', got %q", desc.Name)
	}

	// Test another valid function
	desc = descriptions.GetBuiltinFunctionDescription("length")
	if desc == nil {
		t.Fatal("Expected description for 'length' function, got nil")
	}
	if desc.Name != "length" {
		t.Errorf("Expected name 'length', got %q", desc.Name)
	}

	// Test invalid function
	desc = descriptions.GetBuiltinFunctionDescription("invalid")
	if desc != nil {
		t.Errorf("Expected nil for invalid function, got %+v", desc)
	}
}

// TestGetKeywordDescription tests the individual keyword lookup function.
func TestGetKeywordDescription(t *testing.T) {
	// Test valid keyword
	desc := descriptions.GetKeywordDescription("fn")
	if desc == nil {
		t.Fatal("Expected description for 'fn' keyword, got nil")
	}
	if desc.Keyword != "fn" {
		t.Errorf("Expected keyword 'fn', got %q", desc.Keyword)
	}

	// Test another valid keyword
	desc = descriptions.GetKeywordDescription("let")
	if desc == nil {
		t.Fatal("Expected description for 'let' keyword, got nil")
	}
	if desc.Keyword != "let" {
		t.Errorf("Expected keyword 'let', got %q", desc.Keyword)
	}

	// Test invalid keyword
	desc = descriptions.GetKeywordDescription("invalid")
	if desc != nil {
		t.Errorf("Expected nil for invalid keyword, got %+v", desc)
	}
}
