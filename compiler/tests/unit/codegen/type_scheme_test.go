package codegen

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestTypeSchemeString(t *testing.T) {
	// Create a TypeInferer to use Generalize method
	ti := codegen.NewTypeInferer()

	// Test TypeScheme with no free vars (should not be generalized)
	concreteType := codegen.NewConcreteType("int")
	scheme1 := ti.Generalize(concreteType)
	result1 := scheme1.String()
	expected1 := "int"
	if result1 != expected1 {
		t.Errorf("Expected %s, got %s", expected1, result1)
	}

	// Test TypeScheme with type variables
	typeVar := ti.Fresh()
	scheme2 := ti.Generalize(typeVar)
	result2 := scheme2.String()
	// Should contain "forall" for type variables
	if !strings.Contains(result2, "forall") {
		t.Errorf("Expected scheme with type variables to contain 'forall', got %s", result2)
	}
}

func TestTypeSchemeCategory(t *testing.T) {
	ti := codegen.NewTypeInferer()

	// Test that Category delegates to underlying type
	intType := codegen.NewConcreteType("int")
	scheme := ti.Generalize(intType)

	result := scheme.Category()
	expected := intType.Category()
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTypeSchemeEquals(t *testing.T) {
	ti := codegen.NewTypeInferer()

	// Test equal TypeSchemes
	intType := codegen.NewConcreteType("int")
	scheme1 := ti.Generalize(intType)
	scheme2 := ti.Generalize(intType)
	if !scheme1.Equals(scheme2) {
		t.Error("Expected equal TypeSchemes to be equal")
	}

	// Test different underlying types
	stringType := codegen.NewConcreteType("string")
	scheme3 := ti.Generalize(stringType)
	if scheme1.Equals(scheme3) {
		t.Error("Expected TypeSchemes with different types to be not equal")
	}

	// Test non-TypeScheme type
	if scheme1.Equals(intType) {
		t.Error("Expected TypeScheme to not equal non-TypeScheme")
	}
}

func TestTypeSchemeInstantiate(t *testing.T) {
	ti := codegen.NewTypeInferer()

	// Test Instantiate method
	typeVar := ti.Fresh()
	scheme := ti.Generalize(typeVar)

	// Instantiate should create a fresh instance
	instance1 := ti.Instantiate(scheme)
	instance2 := ti.Instantiate(scheme)

	// Both instances should be valid types
	if instance1 == nil || instance2 == nil {
		t.Error("Instantiate should return valid types")
	}
}