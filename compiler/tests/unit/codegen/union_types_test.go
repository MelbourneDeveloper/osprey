package codegen

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestUnionTypeDeclarationCompiles(t *testing.T) {
	source := `type Color = Red | Green | Blue`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Expected union type declaration to compile successfully, got: %v", err)
	}
}

func TestUnionTypeVariantUsage(t *testing.T) {
	source := `
		type Color = Red | Green | Blue
		let currentColor = Red
	`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Expected union type variant usage to compile successfully, got: %v", err)
	}
}

func TestUnionTypeInMatchExpression(t *testing.T) {
	source := `
		type Grade = A | B | C | D | F
		fn gradeMessage(grade: Grade) -> string = match grade {
			A => "Excellent!"
			B => "Good"
			C => "Average"
			D => "Poor"
			F => "Fail"
		}
	`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Expected union type in match expression to compile successfully, got: %v", err)
	}
}

func TestUndefinedVariantShouldFail(t *testing.T) {
	source := `
		type Color = Red | Green | Blue
		let invalidColor = Purple
	`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Errorf("Expected undefined variant 'Purple' to cause compilation failure")
	}

	if !strings.Contains(err.Error(), "undefined variable") {
		t.Errorf("Expected 'undefined variable' error, got: %v", err)
	}
}

func TestMatchExhaustivenessShouldPass(t *testing.T) {
	// This should pass - all variants are covered
	source := `
		type Color = Red | Green | Blue
		let color = Red
		let description = match color {
			Red => "red"
			Green => "green"
			Blue => "blue"
		}
	`

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Errorf("Expected exhaustive match to compile successfully, got: %v", err)
	}
}

func TestUnknownVariantInMatchShouldFail(t *testing.T) {
	// This should fail with unknown variant error
	source := `
		type Color = Red | Green | Blue
		let color = Red
		let description = match color {
			Red => "red"
			Green => "green"
			Blue => "blue"
			Purple => "invalid"
		}
	`

	_, err := codegen.CompileToLLVM(source)
	if err == nil {
		t.Errorf("Expected unknown variant 'Purple' to cause compilation failure")
	}

	if !strings.Contains(err.Error(), "Purple") || !strings.Contains(err.Error(), "not defined in type") {
		t.Errorf("Expected 'unknown variant' error for Purple, got: %v", err)
	}
}
