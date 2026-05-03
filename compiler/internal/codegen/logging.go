package codegen

import (
	"context"
	"log/slog"

	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/logging"
)

const nilLogValue = "<nil>"

func (g *LLVMGenerator) codegenLogger() *slog.Logger {
	if g.logger != nil {
		return g.logger
	}

	return logging.Logger("codegen")
}

func (g *LLVMGenerator) logTrace(message string, args ...any) {
	ctx := context.Background()
	codegenLogger := g.codegenLogger()
	if codegenLogger.Enabled(ctx, logging.LevelTrace) {
		codegenLogger.Log(ctx, logging.LevelTrace, message, args...)
	}
}

func (g *LLVMGenerator) logDebug(message string, args ...any) {
	ctx := context.Background()
	codegenLogger := g.codegenLogger()
	if codegenLogger.Enabled(ctx, slog.LevelDebug) {
		codegenLogger.DebugContext(ctx, message, args...)
	}
}

func (g *LLVMGenerator) typeNames(types []Type) []string {
	names := make([]string, len(types))
	for i, typ := range types {
		names[i] = typeName(typ)
	}

	return names
}

func typeName(typ Type) string {
	if typ == nil {
		return nilLogValue
	}

	return typ.String()
}

func typeKind(typ Type) string {
	switch typ.(type) {
	case *TypeScheme:
		return "type_scheme"
	case *FunctionType:
		return "function"
	case *TypeVar:
		return "type_var"
	case *ConcreteType:
		return "concrete"
	case *PrimitiveType:
		return "primitive"
	case *GenericType:
		return "generic"
	default:
		return "unknown"
	}
}

func llvmValueName(val value.Value) string {
	if val == nil {
		return nilLogValue
	}

	return val.Ident()
}

func llvmValueType(val value.Value) string {
	if val == nil || val.Type() == nil {
		return nilLogValue
	}

	return val.Type().String()
}
