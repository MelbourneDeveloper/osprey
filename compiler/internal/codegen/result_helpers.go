// Implements [ERR-PAYLOAD]
//
// Construction helpers for Result<T, string>. Every Result built by codegen
// flows through these helpers so the err_msg slot (index 2) is uniformly
// initialised: null for Success, an interned static C string for Error.
//
// Reader side lives in generateErrorBlock (llvm.go); intern table backs onto
// the pre-existing stringConstants map on LLVMGenerator.

package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// nullErrorMessage returns the i8* null constant used to populate the err_msg
// slot of a Success Result.
func (g *LLVMGenerator) nullErrorMessage() value.Value {
	return constant.NewNull(types.I8Ptr)
}

// internErrorMessage returns an i8* pointing to a deduplicated, null-terminated
// global string in .rodata. Repeated calls with the same text return the same
// pointer. Lifetime is the lifetime of the program image.
func (g *LLVMGenerator) internErrorMessage(msg string) value.Value {
	if existing, ok := g.stringConstants[msg]; ok {
		return existing
	}
	bytes := msg + "\x00"
	globalIdx := len(g.stringConstants)
	name := errMsgGlobalName(globalIdx)
	global := g.module.NewGlobalDef(name, constant.NewCharArrayFromString(bytes))
	global.Immutable = true
	ptr := constant.NewGetElementPtr(global.ContentType, global,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.stringConstants[msg] = ptr
	return ptr
}

// errMsgGlobalName produces a stable, unique global name without leaking
// raw user message text into the symbol table.
func errMsgGlobalName(idx int) string {
	return fmt.Sprintf("osp_err_msg_%d", idx)
}

// Result struct field indices. Implements [ERR-PAYLOAD]:
// slot 2 is the new err_msg pointer added so error paths can carry a real
// message instead of aliasing one hardcoded global.
const (
	resultValueFieldIdx = 0
	resultDiscFieldIdx  = 1
	resultErrMsgFieldIdx = 2
)

// makeResultValue builds a Result struct value (not a pointer) carrying
// {value, discriminant, err_msg}. Used by every codegen site that returns a
// Result by value (InsertValue chains).
func (g *LLVMGenerator) makeResultValue(val value.Value, discriminant int64, errMsg value.Value) value.Value {
	resultType := g.getResultType(val.Type())
	undef := constant.NewUndef(resultType)
	s1 := g.builder.NewInsertValue(undef, val, resultValueFieldIdx)
	s2 := g.builder.NewInsertValue(s1, constant.NewInt(types.I8, discriminant), resultDiscFieldIdx)
	s3 := g.builder.NewInsertValue(s2, errMsg, resultErrMsgFieldIdx)
	return s3
}

// makeSuccessValue: Result { value, discriminant=0, err_msg=null }.
func (g *LLVMGenerator) makeSuccessValue(val value.Value) value.Value {
	return g.makeResultValue(val, 0, g.nullErrorMessage())
}

// makeErrorValueWithMessage: Result { defaultVal (typed slot, ignored), discriminant=1, err_msg=interned(msg) }.
// defaultVal supplies the type-correct value slot — its actual content is meaningless when discriminant==1.
func (g *LLVMGenerator) makeErrorValueWithMessage(defaultVal value.Value, msg string) value.Value {
	return g.makeResultValue(defaultVal, 1, g.internErrorMessage(msg))
}

// loadResultErrorMessage reads the err_msg slot (index 2) from a Result.
// The Result may be supplied either as a pointer to a struct (the common
// alloca form) or as a struct value (the InsertValue form used by Result
// constructors and the wrap helpers in function_signatures.go). Implements
// the read side of [ERR-PAYLOAD] for generateErrorBlock.
//
// If the input is somehow not a Result we fall back to a null i8* so callers
// don't crash; in practice the type-checker rules that case out, but
// generateErrorBlock runs over the matched discriminant whose type is fully
// known by the time we get here.
func (g *LLVMGenerator) loadResultErrorMessage(result value.Value) value.Value {
	if result == nil {
		return g.nullErrorMessage()
	}
	if ptrType, ok := result.Type().(*types.PointerType); ok {
		if _, ok := ptrType.ElemType.(*types.StructType); ok {
			errMsgPtr := g.builder.NewGetElementPtr(ptrType.ElemType, result,
				constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultErrMsgFieldIdx))
			return g.builder.NewLoad(types.I8Ptr, errMsgPtr)
		}
	}
	if _, ok := result.Type().(*types.StructType); ok {
		return g.builder.NewExtractValue(result, resultErrMsgFieldIdx)
	}
	return g.nullErrorMessage()
}

// storeResultFields writes the three Result fields into an already-allocated
// pointer (alloca). Used by codegen sites that build Results by pointer.
func (g *LLVMGenerator) storeResultFields(
	resultPtr, val, discriminant, errMsg value.Value,
) {
	resultType := resultPtr.Type().(*types.PointerType).ElemType
	valuePtr := g.builder.NewGetElementPtr(resultType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultValueFieldIdx))
	g.builder.NewStore(val, valuePtr)
	discPtr := g.builder.NewGetElementPtr(resultType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultDiscFieldIdx))
	g.builder.NewStore(discriminant, discPtr)
	errMsgPtr := g.builder.NewGetElementPtr(resultType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultErrMsgFieldIdx))
	g.builder.NewStore(errMsg, errMsgPtr)
}
