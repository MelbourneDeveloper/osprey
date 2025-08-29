package descriptions

import (
	"testing"
)

// TestAllBuiltinFunctionsDocumented ensures all built-in functions have documentation.
// This test will fail the build if any built-in function is missing documentation.
func TestAllBuiltinFunctionsDocumented(t *testing.T) {
	missing := ValidateAllBuiltinFunctionsDocumented()

	if len(missing) > 0 {
		t.Errorf("The following built-in functions are missing documentation: %v", missing)
		t.Error("All built-in functions must be documented in GetBuiltinFunctionDescriptions()")
		t.Error("Add the missing functions to the descriptions map in functions.go")
	}
}

// TestDocumentationCompleteness validates that all documented functions have complete information.
func TestDocumentationCompleteness(t *testing.T) {
	descriptions := GetBuiltinFunctionDescriptions()

	for name, desc := range descriptions {
		if desc.Name == "" {
			t.Errorf("Function %s is missing Name field", name)
		}

		if desc.Signature == "" {
			t.Errorf("Function %s is missing Signature field", name)
		}

		if desc.Description == "" {
			t.Errorf("Function %s is missing Description field", name)
		}

		if desc.ReturnType == "" {
			t.Errorf("Function %s is missing ReturnType field", name)
		}

		if desc.Example == "" {
			t.Errorf("Function %s is missing Example field", name)
		}

		// Check that the name in the description matches the key
		if desc.Name != name {
			t.Errorf("Function %s has mismatched name in description: %s", name, desc.Name)
		}
	}
}

// TestFunctionSignatureConsistency validates that function signatures follow consistent patterns.
func TestFunctionSignatureConsistency(t *testing.T) {
	descriptions := GetBuiltinFunctionDescriptions()

	for name, desc := range descriptions {
		// Check that signature starts with function name
		if len(desc.Signature) < len(name) || desc.Signature[:len(name)] != name {
			t.Errorf("Function %s signature does not start with function name: %s", name, desc.Signature)
		}

		// Check that signature contains -> for return type
		if desc.Signature[len(name):len(name)+1] != "(" {
			t.Errorf("Function %s signature does not have opening parenthesis after name: %s", name, desc.Signature)
		}
	}
}
