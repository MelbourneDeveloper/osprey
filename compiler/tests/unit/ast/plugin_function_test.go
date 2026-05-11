package ast_test

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
)

func TestPluginFunctionParsingPreservesInlineSQL(t *testing.T) {
	source := "fn sql getUsers(status: String, limit: Int) -> UserRows = " +
		"SELECT id, name FROM users WHERE status = $status LIMIT $limit;"

	program := parseToAST(t, source)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	pluginFn, ok := program.Statements[0].(*ast.PluginFunctionDeclaration)
	if !ok {
		t.Fatalf("expected plugin function declaration, got %T", program.Statements[0])
	}

	if pluginFn.PluginName != "sql" {
		t.Fatalf("expected plugin name sql, got %q", pluginFn.PluginName)
	}
	if pluginFn.FunctionName != "getUsers" {
		t.Fatalf("expected function name getUsers, got %q", pluginFn.FunctionName)
	}
	if got := len(pluginFn.Parameters); got != 2 {
		t.Fatalf("expected 2 parameters, got %d", got)
	}
	if pluginFn.ReturnType == nil || pluginFn.ReturnType.Name != "UserRows" {
		t.Fatalf("expected UserRows return type, got %#v", pluginFn.ReturnType)
	}

	expectedSQL := "SELECT id, name FROM users WHERE status = $status LIMIT $limit"
	if pluginFn.PluginContent != expectedSQL {
		t.Fatalf("expected SQL %q, got %q", expectedSQL, pluginFn.PluginContent)
	}
}
