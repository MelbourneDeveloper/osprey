// Package descriptions provides comprehensive documentation for built-in functions.
package descriptions

import (
	"github.com/christianfindlay/osprey/internal/codegen"
)

// BuiltinFunctionDesc represents documentation for a built-in function.
type BuiltinFunctionDesc struct {
	Name        string
	Signature   string
	Description string
	Parameters  []ParameterDesc
	ReturnType  string
	Example     string
}

// ParameterDesc represents documentation for a function parameter.
type ParameterDesc struct {
	Name        string
	Type        string
	Description string
}

// GetBuiltinFunctionDescriptions returns all built-in function descriptions from the unified registry.
func GetBuiltinFunctionDescriptions() map[string]*BuiltinFunctionDesc {
	// Convert from the unified registry format to the legacy format
	descriptions := make(map[string]*BuiltinFunctionDesc)
	
	for _, fn := range codegen.GlobalBuiltInRegistry.GetAllFunctions() {
		// Convert parameters
		params := make([]ParameterDesc, len(fn.ParameterTypes))
		for i, param := range fn.ParameterTypes {
			params[i] = ParameterDesc{
				Name:        param.Name,
				Type:        param.Type.String(),
				Description: param.Description,
			}
		}
		
		descriptions[fn.Name] = &BuiltinFunctionDesc{
			Name:        fn.Name,
			Signature:   fn.Signature,
			Description: fn.Description,
			Parameters:  params,
			ReturnType:  fn.ReturnType.String(),
			Example:     fn.Example,
		}
	}
	
	return descriptions
}

// GetBuiltinFunctionDescription returns description for a single built-in function.
func GetBuiltinFunctionDescription(name string) *BuiltinFunctionDesc {
	descriptions := GetBuiltinFunctionDescriptions()
	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return nil
}

// ValidateAllBuiltinFunctionsDocumented checks that all built-in functions are documented.
// This function should be called during build/test to ensure documentation completeness.
func ValidateAllBuiltinFunctionsDocumented() []string {
	// Get the authoritative list of built-in functions from the compiler's constants
	builtinFunctions := GetCompilerBuiltinFunctionNames()

	descriptions := GetBuiltinFunctionDescriptions()
	var missing []string

	for _, funcName := range builtinFunctions {
		if _, exists := descriptions[funcName]; !exists {
			missing = append(missing, funcName)
		}
	}

	return missing
}

// GetCompilerBuiltinFunctionNames returns all built-in function names from the unified registry.
// This is the authoritative source - it reads directly from the unified built-in function registry.
func GetCompilerBuiltinFunctionNames() []string {
	return codegen.GlobalBuiltInRegistry.GetFunctionNames()
}

// GetAllBuiltinFunctionNames returns all documented built-in function names.
// This can be used to cross-check against the actual compiler implementation.
func GetAllBuiltinFunctionNames() []string {
	descriptions := GetBuiltinFunctionDescriptions()
	var names []string
	for name := range descriptions {
		names = append(names, name)
	}
	return names
}