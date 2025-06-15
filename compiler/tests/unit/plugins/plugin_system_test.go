package plugins

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/plugins"
)

func TestSQLPluginBasicFunctionality(t *testing.T) {
	// Get the current working directory (should be the compiler root)
	compilerDir, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get compiler directory: %v", err)
	}

	// Create plugin system
	pluginSystem := plugins.NewPluginSystem(compilerDir)

	// Create a test plugin function declaration
	pluginFn := &ast.PluginFunctionDeclaration{
		PluginName:   "sql",
		FunctionName: "getUsers",
		Parameters: []ast.Parameter{
			{Name: "type", Type: &ast.TypeExpression{Name: "int"}},
		},
		PluginContent: "SELECT id, name, email FROM users WHERE type = $type",
	}

	// Process the plugin function
	response, err := pluginSystem.ProcessPluginFunction(pluginFn, "test.osp", 1)
	if err != nil {
		t.Fatalf("Failed to process plugin function: %v", err)
	}

	// Verify the response
	if !response.Success {
		t.Fatalf("Plugin processing failed: %s", response.Error)
	}

	validateBasicFunctionalityResponse(t, response)

	t.Logf("✅ SQL Plugin test passed!")
	t.Logf("📦 Return Type: %s", response.ReturnType.Type)
	t.Logf("🔧 Generated Fiber Code: %s", response.GeneratedCode.FiberCode[:100]+"...")
}

// validateBasicFunctionalityResponse validates the response from basic functionality test.
func validateBasicFunctionalityResponse(t *testing.T, response *plugins.PluginResponse) {
	// Check return type
	if response.ReturnType.Type != "Result" {
		t.Errorf("Expected return type 'Result', got '%s'", response.ReturnType.Type)
	}

	validateGenericParameters(t, response)
	validateGeneratedCode(t, response)
	validateImportsAndErrorTypes(t, response)
}

// validateGenericParameters validates the generic parameters of the return type.
func validateGenericParameters(t *testing.T, response *plugins.PluginResponse) {
	// Check that it has generic parameters for Result<Array<Record>, DatabaseError>
	if len(response.ReturnType.GenericParams) != 2 {
		t.Errorf("Expected 2 generic parameters, got %d", len(response.ReturnType.GenericParams))
	}

	// Check first generic parameter is Array
	if response.ReturnType.GenericParams[0].Type != "Array" {
		t.Errorf("Expected first generic param to be 'Array', got '%s'", response.ReturnType.GenericParams[0].Type)
	}

	// Check second generic parameter is DatabaseError
	if response.ReturnType.GenericParams[1].Type != "DatabaseError" {
		t.Errorf("Expected second generic param to be 'DatabaseError', got '%s'", response.ReturnType.GenericParams[1].Type)
	}
}

// validateGeneratedCode validates the generated code content.
func validateGeneratedCode(t *testing.T, response *plugins.PluginResponse) {
	// Check generated code
	if response.GeneratedCode.FiberCode == "" {
		t.Error("Expected generated fiber code, got empty string")
	}
}

// validateImportsAndErrorTypes validates imports and error types.
func validateImportsAndErrorTypes(t *testing.T, response *plugins.PluginResponse) {
	// Check imports
	expectedImports := []string{"std.database", "std.sql", "std.fiber"}
	if len(response.GeneratedCode.Imports) != len(expectedImports) {
		t.Errorf("Expected %d imports, got %d", len(expectedImports), len(response.GeneratedCode.Imports))
	}

	for i, expectedImport := range expectedImports {
		if i >= len(response.GeneratedCode.Imports) || response.GeneratedCode.Imports[i] != expectedImport {
			t.Errorf("Expected import '%s', got '%s'", expectedImport, response.GeneratedCode.Imports[i])
		}
	}

	// Check error types
	expectedErrorTypes := []string{"ConnectionError", "QueryError", "TimeoutError"}
	if len(response.GeneratedCode.ErrorTypes) != len(expectedErrorTypes) {
		t.Errorf("Expected %d error types, got %d", len(expectedErrorTypes), len(response.GeneratedCode.ErrorTypes))
	}
}

func TestSQLPluginParameterValidation(t *testing.T) {
	compilerDir, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get compiler directory: %v", err)
	}

	pluginSystem := plugins.NewPluginSystem(compilerDir)

	// Test with missing parameter
	pluginFn := &ast.PluginFunctionDeclaration{
		PluginName:   "sql",
		FunctionName: "getUsers",
		Parameters: []ast.Parameter{
			{Name: "status", Type: &ast.TypeExpression{Name: "string"}},
		},
		PluginContent: "SELECT id, name, email FROM users WHERE type = $type", // uses $type but parameter is $status
	}

	_, err = pluginSystem.ProcessPluginFunction(pluginFn, "test.osp", 1)

	// The plugin should return an error for parameter mismatch
	if err == nil {
		t.Error("Expected plugin processing to fail due to parameter mismatch, but it succeeded")
	}

	// Check that the error message mentions the parameter mismatch
	expectedErrorSubstring := "query parameter $type not found in function parameters"
	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErrorSubstring, err.Error())
	}

	t.Logf("✅ Parameter validation test passed!")
	t.Logf("❌ Error (expected): %s", err.Error())
}

func TestSQLPluginInsertQuery(t *testing.T) {
	compilerDir, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get compiler directory: %v", err)
	}

	pluginSystem := plugins.NewPluginSystem(compilerDir)

	// Test INSERT query
	pluginFn := &ast.PluginFunctionDeclaration{
		PluginName:   "sql",
		FunctionName: "createUser",
		Parameters: []ast.Parameter{
			{Name: "name", Type: &ast.TypeExpression{Name: "string"}},
			{Name: "email", Type: &ast.TypeExpression{Name: "string"}},
		},
		PluginContent: "INSERT INTO users (name, email) VALUES ($name, $email) RETURNING id",
	}

	response, err := pluginSystem.ProcessPluginFunction(pluginFn, "test.osp", 1)
	if err != nil {
		t.Fatalf("Failed to process INSERT plugin function: %v", err)
	}

	if !response.Success {
		t.Fatalf("Plugin processing failed: %s", response.Error)
	}

	// For INSERT with RETURNING, should still return Result type
	if response.ReturnType.Type != "Result" {
		t.Errorf("Expected return type 'Result', got '%s'", response.ReturnType.Type)
	}

	t.Logf("✅ INSERT query test passed!")
	t.Logf("📝 INSERT Return Type: %s", response.ReturnType.Type)
}
