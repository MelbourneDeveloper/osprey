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

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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

// extractResultDiscriminant reads the discriminant slot (index 1) from a
// Result value (struct value or pointer-to-struct).
func (g *LLVMGenerator) extractResultDiscriminant(result value.Value) value.Value {
	if ptrType, ok := result.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok {
			discPtr := g.builder.NewGetElementPtr(structType, result,
				constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultDiscFieldIdx))
			return g.builder.NewLoad(structType.Fields[resultDiscFieldIdx], discPtr)
		}
	}
	return g.builder.NewExtractValue(result, resultDiscFieldIdx)
}

// extractResultValueSlot reads the value slot (index 0) from a Result value.
// Mirrors extractResultDiscriminant for struct-value vs pointer-to-struct.
func (g *LLVMGenerator) extractResultValueSlot(result value.Value) value.Value {
	if ptrType, ok := result.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok {
			vPtr := g.builder.NewGetElementPtr(structType, result,
				constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultValueFieldIdx))
			return g.builder.NewLoad(structType.Fields[resultValueFieldIdx], vPtr)
		}
	}
	return g.builder.NewExtractValue(result, resultValueFieldIdx)
}

// withResultErrorPropagation runs `compute` only if neither operand is an
// Error Result. If either operand is an Error, the surrounding arithmetic
// short-circuits to that Error — preserving the err_msg slot — and `compute`
// is not invoked. `compute` receives the operands with any Result wrapping
// already stripped, and must return a value of type
// Result<resultElemType, string>.
//
// Implements the auto-unwrap propagation requirement of spec [ERR-PAYLOAD]
// for binary arithmetic chains (see docs/plans/error-payloads.md Phase 4).
func (g *LLVMGenerator) withResultErrorPropagation(
	left, right value.Value,
	resultElemType types.Type,
	compute func(l, r value.Value) (value.Value, error),
) (value.Value, error) {
	leftIsRes := g.isResultType(left)
	rightIsRes := g.isResultType(right)
	if !leftIsRes && !rightIsRes {
		return compute(left, right)
	}

	blockID := len(g.function.Blocks)
	propBlock := g.function.NewBlock(fmt.Sprintf("arith_prop_%d", blockID))
	okBlock := g.function.NewBlock(fmt.Sprintf("arith_ok_%d", blockID))
	endBlock := g.function.NewBlock(fmt.Sprintf("arith_end_%d", blockID))

	errDisc := constant.NewInt(types.I8, 1)
	falseVal := constant.NewBool(false)
	var leftDiscIsErr, rightDiscIsErr value.Value = falseVal, falseVal
	if leftIsRes {
		leftDiscIsErr = g.builder.NewICmp(enum.IPredEQ,
			g.extractResultDiscriminant(left), errDisc)
	}
	if rightIsRes {
		rightDiscIsErr = g.builder.NewICmp(enum.IPredEQ,
			g.extractResultDiscriminant(right), errDisc)
	}
	anyErr := g.builder.NewOr(leftDiscIsErr, rightDiscIsErr)
	g.builder.NewCondBr(anyErr, propBlock, okBlock)

	g.builder = propBlock
	propResult := g.buildPropagatedErrorResult(
		left, right, leftIsRes, rightIsRes, leftDiscIsErr, resultElemType,
	)
	propLastBlock := g.builder
	g.builder.NewBr(endBlock)

	g.builder = okBlock
	leftExtracted := left
	if leftIsRes {
		leftExtracted = g.extractResultValueSlot(left)
	}
	rightExtracted := right
	if rightIsRes {
		rightExtracted = g.extractResultValueSlot(right)
	}
	okResult, err := compute(leftExtracted, rightExtracted)
	if err != nil {
		return nil, err
	}
	okLastBlock := g.builder
	if g.builder.Term == nil {
		g.builder.NewBr(endBlock)
	}

	g.builder = endBlock
	phi := g.builder.NewPhi(
		ir.NewIncoming(propResult, propLastBlock),
		ir.NewIncoming(okResult, okLastBlock),
	)
	return phi, nil
}

// buildPropagatedErrorResult constructs the Result that withResultErrorPropagation
// emits in its propagation block. The err_msg comes from the leftmost Error
// operand (or the only Result operand if just one is a Result). The value
// slot is a typed zero of resultElemType so the propagation block's Result
// type unifies with the success block's at the trailing PHI.
func (g *LLVMGenerator) buildPropagatedErrorResult(
	left, right value.Value,
	leftIsRes, rightIsRes bool,
	leftDiscIsErr value.Value,
	resultElemType types.Type,
) value.Value {
	var propMsg value.Value
	switch {
	case leftIsRes && rightIsRes:
		leftMsg := g.extractErrMsgSlot(left)
		rightMsg := g.extractErrMsgSlot(right)
		propMsg = g.builder.NewSelect(leftDiscIsErr, leftMsg, rightMsg)
	case leftIsRes:
		propMsg = g.extractErrMsgSlot(left)
	default:
		propMsg = g.extractErrMsgSlot(right)
	}
	defaultVal := zeroValueForType(resultElemType)
	return g.makeResultValue(defaultVal, 1, propMsg)
}

// extractErrMsgSlot reads the err_msg slot (index 2) without going through
// the generic loadResultErrorMessage path, which is wired for read-back at
// match arms; here we're reading directly off a Result-typed operand to
// propagate it through arithmetic.
func (g *LLVMGenerator) extractErrMsgSlot(result value.Value) value.Value {
	if ptrType, ok := result.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok {
			msgPtr := g.builder.NewGetElementPtr(structType, result,
				constant.NewInt(types.I32, 0), constant.NewInt(types.I32, resultErrMsgFieldIdx))
			return g.builder.NewLoad(types.I8Ptr, msgPtr)
		}
	}
	return g.builder.NewExtractValue(result, resultErrMsgFieldIdx)
}
