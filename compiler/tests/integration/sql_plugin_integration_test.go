package integration

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/plugins"
)

func TestSQLPluginFullIntegration(t *testing.T) {
	t.Logf("🚀 Starting SQL Plugin Integration Test!")

	// Get the compiler root directory
	compilerDir, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get compiler directory: %v", err)
	}

	// Create plugin system
	pluginSystem := plugins.NewPluginSystem(compilerDir)

	// Get test cases and run them
	testCases := getSQLIntegrationTestCases()
	runSQLIntegrationTestCases(t, pluginSystem, testCases)

	t.Logf("🎉 SQL Plugin Integration Test COMPLETED!")
}

// getSQLIntegrationTestCases returns the test cases for SQL plugin integration testing.
func getSQLIntegrationTestCases() []struct {
	name         string
	pluginFn     *ast.PluginFunctionDeclaration
	expectError  bool
	expectedType string
	description  string
} {
	return []struct {
		name         string
		pluginFn     *ast.PluginFunctionDeclaration
		expectError  bool
		expectedType string
		description  string
	}{
		{
			name: "SELECT Query",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:    "sql",
				FunctionName:  "getAllUsers",
				Parameters:    []ast.Parameter{},
				PluginContent: "SELECT id, name, email, created_at FROM users",
			},
			expectError:  false,
			expectedType: "Result",
			description:  "Simple SELECT query returning all users",
		},
		{
			name: "Parameterized SELECT",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:   "sql",
				FunctionName: "getUsersByStatus",
				Parameters: []ast.Parameter{
					{Name: "status", Type: &ast.TypeExpression{Name: "string"}},
					{Name: "limit", Type: &ast.TypeExpression{Name: "int"}},
				},
				PluginContent: "SELECT id, name, email FROM users WHERE status = $status LIMIT $limit",
			},
			expectError:  false,
			expectedType: "Result",
			description:  "Parameterized SELECT with multiple parameters",
		},
		{
			name: "INSERT with RETURNING",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:   "sql",
				FunctionName: "createUser",
				Parameters: []ast.Parameter{
					{Name: "name", Type: &ast.TypeExpression{Name: "string"}},
					{Name: "email", Type: &ast.TypeExpression{Name: "string"}},
				},
				PluginContent: "INSERT INTO users (name, email) VALUES ($name, $email) RETURNING id, created_at",
			},
			expectError:  false,
			expectedType: "Result",
			description:  "INSERT with RETURNING clause",
		},
		{
			name: "UPDATE Query",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:   "sql",
				FunctionName: "updateUserEmail",
				Parameters: []ast.Parameter{
					{Name: "user_id", Type: &ast.TypeExpression{Name: "int"}},
					{Name: "new_email", Type: &ast.TypeExpression{Name: "string"}},
				},
				PluginContent: "UPDATE users SET email = $new_email WHERE id = $user_id",
			},
			expectError:  false,
			expectedType: "Result",
			description:  "UPDATE query with parameters",
		},
		{
			name: "DELETE Query",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:   "sql",
				FunctionName: "deleteUser",
				Parameters: []ast.Parameter{
					{Name: "user_id", Type: &ast.TypeExpression{Name: "int"}},
				},
				PluginContent: "DELETE FROM users WHERE id = $user_id",
			},
			expectError:  false,
			expectedType: "Result",
			description:  "DELETE query with parameter",
		},
		{
			name: "Parameter Mismatch Error",
			pluginFn: &ast.PluginFunctionDeclaration{
				PluginName:   "sql",
				FunctionName: "badQuery",
				Parameters: []ast.Parameter{
					{Name: "wrong_param", Type: &ast.TypeExpression{Name: "string"}},
				},
				PluginContent: "SELECT * FROM users WHERE id = $missing_param",
			},
			expectError:  true,
			expectedType: "",
			description:  "Should fail due to parameter mismatch",
		},
	}
}

// runSQLIntegrationTestCases executes the SQL integration test cases.
func runSQLIntegrationTestCases(t *testing.T, pluginSystem *plugins.PluginSystem, testCases []struct {
	name         string
	pluginFn     *ast.PluginFunctionDeclaration
	expectError  bool
	expectedType string
	description  string
},
) {
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("📝 Testing: %s", tc.description)

			response, err := pluginSystem.ProcessPluginFunction(tc.pluginFn, "integration_test.osp", 1)

			validateSQLTestResponse(t, tc, response, err)
		})
	}
}

// validateSQLTestResponse validates the response from SQL plugin processing.
func validateSQLTestResponse(t *testing.T, tc struct {
	name         string
	pluginFn     *ast.PluginFunctionDeclaration
	expectError  bool
	expectedType string
	description  string
}, response *plugins.PluginResponse, err error,
) {
	if tc.expectError {
		handleExpectedError(t, err)
	} else {
		handleSuccessResponse(t, tc, response, err)
	}
}

// handleExpectedError handles test cases that expect errors.
func handleExpectedError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error but got none")
	} else {
		t.Logf("✅ Expected error caught: %s", err.Error())
	}
}

// handleSuccessResponse handles test cases that expect success.
func handleSuccessResponse(t *testing.T, tc struct {
	name         string
	pluginFn     *ast.PluginFunctionDeclaration
	expectError  bool
	expectedType string
	description  string
}, response *plugins.PluginResponse, err error,
) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !response.Success {
		t.Errorf("Plugin processing failed: %s", response.Error)
		return
	}

	validateSuccessResponseDetails(t, tc, response)
}

// validateSuccessResponseDetails validates the details of a successful response.
func validateSuccessResponseDetails(t *testing.T, tc struct {
	name         string
	pluginFn     *ast.PluginFunctionDeclaration
	expectError  bool
	expectedType string
	description  string
}, response *plugins.PluginResponse,
) {
	// Check return type
	if response.ReturnType.Type != tc.expectedType {
		t.Errorf("Expected return type '%s', got '%s'", tc.expectedType, response.ReturnType.Type)
	}

	// Verify generated code contains fiber implementation
	if !strings.Contains(response.GeneratedCode.FiberCode, "Fiber<") {
		t.Error("Generated code should contain Fiber type")
	}

	validateImports(t, response)

	t.Logf("✅ %s - SUCCESS", tc.name)
	t.Logf("   Return Type: %s", response.ReturnType.Type)
	t.Logf("   Generic Params: %d", len(response.ReturnType.GenericParams))
}

// validateImports verifies that expected imports are present.
func validateImports(t *testing.T, response *plugins.PluginResponse) {
	expectedImports := []string{"std.database", "std.sql", "std.fiber"}
	for _, expectedImport := range expectedImports {
		found := false
		for _, actualImport := range response.GeneratedCode.Imports {
			if actualImport == expectedImport {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Missing expected import: %s", expectedImport)
		}
	}
}

func TestSQLPluginTypeInference(t *testing.T) {
	t.Logf("🔍 Testing SQL Plugin Type Inference")

	compilerDir, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get compiler directory: %v", err)
	}

	pluginSystem := plugins.NewPluginSystem(compilerDir)

	// Test that the plugin correctly infers types based on column names
	pluginFn := &ast.PluginFunctionDeclaration{
		PluginName:    "sql",
		FunctionName:  "getUserStats",
		Parameters:    []ast.Parameter{},
		PluginContent: "SELECT id, name, email, age, active, price, created_at FROM user_stats",
	}

	response, err := pluginSystem.ProcessPluginFunction(pluginFn, "type_test.osp", 1)
	if err != nil {
		t.Fatalf("Plugin processing failed: %v", err)
	}

	validateTypeInferenceResponse(t, response)
}

// validateTypeInferenceResponse validates the type inference response from SQL plugin.
func validateTypeInferenceResponse(t *testing.T, response *plugins.PluginResponse) {
	// Check that we have a Result type with Array and Record
	if response.ReturnType.Type != "Result" {
		t.Errorf("Expected Result type, got %s", response.ReturnType.Type)
	}

	validateTypeParameters(t, response)
	recordParam := getRecordParameter(t, response)
	validateRecordFields(t, recordParam)

	t.Logf("✅ Type inference test passed!")
	t.Logf("📊 Inferred %d fields in record type", len(recordParam.Fields))
}

// validateTypeParameters validates the generic type parameters.
func validateTypeParameters(t *testing.T, response *plugins.PluginResponse) {
	if len(response.ReturnType.GenericParams) != 2 {
		t.Errorf("Expected 2 generic parameters, got %d", len(response.ReturnType.GenericParams))
	}

	arrayParam := response.ReturnType.GenericParams[0]
	if arrayParam.Type != "Array" {
		t.Errorf("Expected Array type, got %s", arrayParam.Type)
	}

	if len(arrayParam.GenericParams) != 1 {
		t.Errorf("Expected 1 array element type, got %d", len(arrayParam.GenericParams))
	}
}

// getRecordParameter extracts and validates the record parameter.
func getRecordParameter(t *testing.T, response *plugins.PluginResponse) *plugins.TypeInfo {
	arrayParam := response.ReturnType.GenericParams[0]
	recordParam := arrayParam.GenericParams[0]
	if recordParam.Type != "Record" {
		t.Errorf("Expected Record type, got %s", recordParam.Type)
	}
	return &recordParam
}

// validateRecordFields validates the fields in the record type.
func validateRecordFields(t *testing.T, recordParam *plugins.TypeInfo) {
	// Verify that common field types are inferred
	t.Logf("🔍 Actual fields in record: %d", len(recordParam.Fields))
	for i, field := range recordParam.Fields {
		t.Logf("   Field %d: %s (%s)", i, field.Name, field.Type)
	}

	expectedFields := []string{"ID", "NAME", "EMAIL", "CREATED_AT"}
	for _, expectedField := range expectedFields {
		found := false
		for _, field := range recordParam.Fields {
			if field.Name == expectedField {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found in record", expectedField)
		}
	}
}
