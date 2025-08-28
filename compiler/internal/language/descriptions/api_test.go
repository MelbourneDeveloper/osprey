package descriptions

import (
	"strings"
	"testing"
)

func TestGetHoverDocumentation_NonexistentElement(t *testing.T) {
	result := GetHoverDocumentation("nonexistent")
	if result != "" {
		t.Errorf("Expected empty string for nonexistent element, got: %s", result)
	}
}

func TestGetHoverDocumentation_Function(t *testing.T) {
	// Test with the 'print' function that should exist
	result := GetHoverDocumentation("print")
	
	if result == "" {
		t.Error("Expected non-empty result for 'print' function")
	}
	
	if !strings.Contains(result, "```osprey") {
		t.Error("Expected code block for function signature")
	}
	
	if !strings.Contains(result, "print") {
		t.Error("Expected function name in result")
	}
}

func TestGetHoverDocumentation_Type(t *testing.T) {
	// Test with a type that should exist
	result := GetHoverDocumentation("String")
	
	if result == "" {
		t.Error("Expected non-empty result for 'String' type")
	}
	
	if !strings.Contains(result, "```osprey") {
		t.Error("Expected code block for type")
	}
}

func TestGetHoverDocumentation_Operator(t *testing.T) {
	// Test with an operator that should exist
	result := GetHoverDocumentation("+")
	
	if result == "" {
		t.Error("Expected non-empty result for '+' operator")
	}
	
	if !strings.Contains(result, "**Operator:**") {
		t.Error("Expected operator label in result")
	}
}

func TestGetHoverDocumentation_Keyword(t *testing.T) {
	// Test with a keyword that should exist
	result := GetHoverDocumentation("let")
	
	if result == "" {
		t.Error("Expected non-empty result for 'let' keyword")
	}
	
	if !strings.Contains(result, "**Keyword:**") {
		t.Error("Expected keyword label in result")
	}
}

func TestGetLanguageElementDescription_ExistingElement(t *testing.T) {
	// Test with a function that should exist
	result := GetLanguageElementDescription("print")
	
	if result == nil {
		t.Error("Expected non-nil result for existing 'print' function")
	} else {
		if result.Type != ElementTypeFunction {
			t.Errorf("Expected type 'function', got: %s", result.Type)
		}
		if result.Name != "print" {
			t.Errorf("Expected name 'print', got: %s", result.Name)
		}
	}
}

func TestGetLanguageElementDescription_NonexistentElement(t *testing.T) {
	result := GetLanguageElementDescription("nonexistent")
	
	if result != nil {
		t.Errorf("Expected nil result for nonexistent element, got: %v", result)
	}
}

func TestGetAllLanguageElements_ReturnsElements(t *testing.T) {
	elements := GetAllLanguageElements()
	
	if len(elements) == 0 {
		t.Error("Expected non-empty elements map")
	}
	
	// Check that we have some expected elements
	if _, exists := elements["print"]; !exists {
		t.Error("Expected 'print' function to exist in elements")
	}
}

func TestGetHoverDocumentation_FunctionWithExample(t *testing.T) {
	// Test that functions with examples include them
	result := GetHoverDocumentation("print")
	
	if result == "" {
		t.Error("Expected non-empty result for 'print' function")
	}
	
	// Look for example section
	if strings.Contains(result, "**Example:**") {
		if !strings.Contains(result, "```osprey") {
			t.Error("Expected code block for example")
		}
	}
}

func TestGetHoverDocumentation_TypeWithDescription(t *testing.T) {
	// Test that types include descriptions
	result := GetHoverDocumentation("Int")
	
	if result == "" {
		t.Error("Expected non-empty result for 'Int' type")
	}
	
	if !strings.Contains(result, "type Int") {
		t.Error("Expected type declaration in result")
	}
}

func TestGetHoverDocumentation_OperatorWithDescription(t *testing.T) {
	// Test that operators include descriptions
	result := GetHoverDocumentation("-")
	
	if result == "" {
		t.Error("Expected non-empty result for '-' operator")
	}
	
	if !strings.Contains(result, "`-`") {
		t.Error("Expected operator symbol in backticks")
	}
}

func TestGetHoverDocumentation_KeywordWithDescription(t *testing.T) {
	// Test that keywords include descriptions
	result := GetHoverDocumentation("fn")
	
	if result == "" {
		t.Error("Expected non-empty result for 'fn' keyword")
	}
	
	if !strings.Contains(result, "`fn`") {
		t.Error("Expected keyword in backticks")
	}
}

func TestGetHoverDocumentation_EmptyExample(t *testing.T) {
	// This tests the path where Example is empty
	// We can't easily mock this, but we can test that the function handles all cases
	elements := GetAllLanguageElements()
	
	for name, element := range elements {
		result := GetHoverDocumentation(name)
		
		// Each element should produce some result
		if result == "" && element != nil {
			t.Errorf("Element %s produced empty result", name)
		}
		
		// Test that the result format is correct based on type
		switch element.Type {
		case ElementTypeFunction:
			if !strings.Contains(result, "```osprey") && element.Signature != "" {
				t.Errorf("Function %s missing code block", name)
			}
		case ElementTypeType:
			if !strings.Contains(result, "type "+element.Name) {
				t.Errorf("Type %s missing type declaration", name)
			}
		case ElementTypeOperator:
			if !strings.Contains(result, "**Operator:**") {
				t.Errorf("Operator %s missing operator label", name)
			}
		case ElementTypeKeyword:
			if !strings.Contains(result, "**Keyword:**") {
				t.Errorf("Keyword %s missing keyword label", name)
			}
		}
	}
}

func TestGetHoverDocumentation_StringFormatting(t *testing.T) {
	// Test that strings are properly joined
	result := GetHoverDocumentation("print")
	
	if result == "" {
		t.Error("Expected non-empty result")
	}
	
	// Should not have triple newlines (indicating improper joining)
	if strings.Contains(result, "\n\n\n") {
		t.Error("Result contains triple newlines, indicating formatting issue")
	}
}