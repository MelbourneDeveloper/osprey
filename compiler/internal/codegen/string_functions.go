// Implements [BUILTIN-STRING-*]
//
// Generators for the string builtins documented in
// compiler/spec/0012-Built-InFunctions.md, section "String Functions".
// Most generators are thin wrappers over C runtime helpers declared in
// compiler/runtime/string_runtime.h. Result-wrapping for fallible
// operations is built directly in LLVM IR using the same {value, i8 disc}
// layout that generateSubstringCall and generateParseIntCall use.

package codegen

import (
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// declareStringRuntime ensures a C runtime function is declared in the module.
// Memoised in g.functions like the other generators do.
func (g *LLVMGenerator) declareStringRuntime(name string, retType types.Type, params ...*ir.Param) *ir.Func {
	if fn, ok := g.functions[name]; ok {
		return fn
	}
	fn := g.module.NewFunc(name, retType, params...)
	g.functions[name] = fn
	return fn
}

// resultFromNullableString returns Result<string, StringError>:
// if ptr is NULL → Error; otherwise → Success { value: ptr }.
func (g *LLVMGenerator) resultFromNullableString(ptr value.Value) value.Value {
	resultType := g.getResultType(types.I8Ptr)
	result := g.builder.NewAlloca(resultType)

	isErr := g.builder.NewICmp(enum.IPredEQ, ptr, constant.NewNull(types.I8Ptr))
	// discriminant = 1 if err else 0
	disc := g.builder.NewSelect(isErr,
		constant.NewInt(types.I8, 1),
		constant.NewInt(types.I8, 0))
	// value = null if err else ptr
	val := g.builder.NewSelect(isErr, constant.NewNull(types.I8Ptr), ptr)

	vp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(val, vp)
	dp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(disc, dp)
	return result
}

// ===== INSPECTION (total) =====

// generateIsEmptyCall: isEmpty(s) -> bool. Wraps osp_string_is_empty.
func (g *LLVMGenerator) generateIsEmptyCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(IsEmptyFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_is_empty",
		types.I64, ir.NewParam("s", types.I8Ptr))
	raw := g.builder.NewCall(fn, s)
	// truncate i64 → i1
	return g.builder.NewICmp(enum.IPredNE, raw, constant.NewInt(types.I64, 0)), nil
}

// ===== SEARCH (total) =====

func (g *LLVMGenerator) generateStartsWithCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(StartsWithFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	prefix, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_starts_with", types.I64,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("prefix", types.I8Ptr))
	raw := g.builder.NewCall(fn, s, prefix)
	return g.builder.NewICmp(enum.IPredNE, raw, constant.NewInt(types.I64, 0)), nil
}

func (g *LLVMGenerator) generateEndsWithCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(EndsWithFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	suffix, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_ends_with", types.I64,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("suffix", types.I8Ptr))
	raw := g.builder.NewCall(fn, s, suffix)
	return g.builder.NewICmp(enum.IPredNE, raw, constant.NewInt(types.I64, 0)), nil
}

// ===== SEARCH (fallible) =====

// generateIndexOfCall: indexOf(s, needle) -> Result<int, StringError>.
// C returns -1 on not-found; we wrap accordingly.
func (g *LLVMGenerator) generateIndexOfCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(IndexOfFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	needle, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_index_of", types.I64,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("needle", types.I8Ptr))
	idx := g.builder.NewCall(fn, s, needle)

	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)
	isErr := g.builder.NewICmp(enum.IPredSLT, idx, constant.NewInt(types.I64, 0))
	disc := g.builder.NewSelect(isErr,
		constant.NewInt(types.I8, 1), constant.NewInt(types.I8, 0))
	val := g.builder.NewSelect(isErr, constant.NewInt(types.I64, 0), idx)
	vp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(val, vp)
	dp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(disc, dp)
	return result, nil
}

// ===== SUBSTRINGS (total) =====

func (g *LLVMGenerator) generateTakeCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateScalarStringIntCall(callExpr, TakeFunc, "osp_string_take")
}

func (g *LLVMGenerator) generateDropCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateScalarStringIntCall(callExpr, DropFunc, "osp_string_drop")
}

// shared: (s: string, n: int) -> string
func (g *LLVMGenerator) generateScalarStringIntCall(
	callExpr *ast.CallExpression, ospName, cName string,
) (value.Value, error) {
	err := validateBuiltInArgs(ospName, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	n, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime(cName, types.I8Ptr,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("n", types.I64))
	return g.builder.NewCall(fn, s, n), nil
}

// ===== TRANSFORM (total) =====

func (g *LLVMGenerator) generateToUpperCaseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, ToUpperCaseFunc, "osp_string_to_upper")
}

func (g *LLVMGenerator) generateToLowerCaseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, ToLowerCaseFunc, "osp_string_to_lower")
}

func (g *LLVMGenerator) generateTrimCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, TrimFunc, "osp_string_trim")
}

func (g *LLVMGenerator) generateTrimStartCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, TrimStartFunc, "osp_string_trim_start")
}

func (g *LLVMGenerator) generateTrimEndCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, TrimEndFunc, "osp_string_trim_end")
}

func (g *LLVMGenerator) generateReverseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateUnaryStringCall(callExpr, ReverseFunc, "osp_string_reverse")
}

// shared: (s: string) -> string
func (g *LLVMGenerator) generateUnaryStringCall(
	callExpr *ast.CallExpression, ospName, cName string,
) (value.Value, error) {
	err := validateBuiltInArgs(ospName, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime(cName, types.I8Ptr, ir.NewParam("s", types.I8Ptr))
	return g.builder.NewCall(fn, s), nil
}

// ===== TRANSFORM (fallible) =====

func (g *LLVMGenerator) generateReplaceCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(ReplaceFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	needle, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	replacement, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_replace", types.I8Ptr,
		ir.NewParam("s", types.I8Ptr),
		ir.NewParam("needle", types.I8Ptr),
		ir.NewParam("replacement", types.I8Ptr))
	ptr := g.builder.NewCall(fn, s, needle, replacement)
	return g.resultFromNullableString(ptr), nil
}

func (g *LLVMGenerator) generateRepeatCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(RepeatFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	n, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_repeat", types.I8Ptr,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("n", types.I64))
	ptr := g.builder.NewCall(fn, s, n)
	return g.resultFromNullableString(ptr), nil
}

func (g *LLVMGenerator) generatePadStartCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generatePadCall(callExpr, PadStartFunc, "osp_string_pad_start")
}

func (g *LLVMGenerator) generatePadEndCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generatePadCall(callExpr, PadEndFunc, "osp_string_pad_end")
}

// shared: (s: string, target_length: int, fill: string) -> Result<string, _>
func (g *LLVMGenerator) generatePadCall(
	callExpr *ast.CallExpression, ospName, cName string,
) (value.Value, error) {
	err := validateBuiltInArgs(ospName, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	target, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	fill, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime(cName, types.I8Ptr,
		ir.NewParam("s", types.I8Ptr),
		ir.NewParam("target_length", types.I64),
		ir.NewParam("fill", types.I8Ptr))
	ptr := g.builder.NewCall(fn, s, target, fill)
	return g.resultFromNullableString(ptr), nil
}

// ===== PARSING =====

// generateParseFloatCall: parseFloat(s) -> Result<float, StringError>.
// C returns 0 on success and writes through out-ptr, 1 on failure.
func (g *LLVMGenerator) generateParseFloatCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(ParseFloatFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_parse_float_strict", types.I64,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("out", types.NewPointer(types.Double)))
	outSlot := g.builder.NewAlloca(types.Double)
	g.builder.NewStore(constant.NewFloat(types.Double, 0.0), outSlot)
	rc := g.builder.NewCall(fn, s, outSlot)
	parsed := g.builder.NewLoad(types.Double, outSlot)

	resultType := g.getResultType(types.Double)
	result := g.builder.NewAlloca(resultType)
	isErr := g.builder.NewICmp(enum.IPredNE, rc, constant.NewInt(types.I64, 0))
	disc := g.builder.NewSelect(isErr,
		constant.NewInt(types.I8, 1), constant.NewInt(types.I8, 0))
	val := g.builder.NewSelect(isErr, constant.NewFloat(types.Double, 0.0), parsed)
	vp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(val, vp)
	dp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(disc, dp)
	return result, nil
}

// ===== LIST-RETURNING =====
//
// split / lines / words return a pointer to osp_string_list. Osprey's
// existing list literal layout is { i64 length, i8* data }. We expose
// the same shape so existing forEach/iter code can consume the result.

// listResultPtrType is the pointer type the C side returns and we hand back.
func listResultPtrType() *types.PointerType {
	return types.NewPointer(types.NewStruct(types.I64, types.I8Ptr))
}

func (g *LLVMGenerator) generateLinesCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(LinesFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_lines",
		listResultPtrType(), ir.NewParam("s", types.I8Ptr))
	return g.builder.NewCall(fn, s), nil
}

func (g *LLVMGenerator) generateWordsCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(WordsFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.declareStringRuntime("osp_string_words",
		listResultPtrType(), ir.NewParam("s", types.I8Ptr))
	return g.builder.NewCall(fn, s), nil
}

// generateSplitCall: split(s, sep) -> Result<List<string>, StringError>.
// C returns NULL when sep is empty.
func (g *LLVMGenerator) generateSplitCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(SplitFunc, callExpr)
	if err != nil {
		return nil, err
	}
	s, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	sep, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	listPtrTy := listResultPtrType()
	fn := g.declareStringRuntime("osp_string_split",
		listPtrTy,
		ir.NewParam("s", types.I8Ptr), ir.NewParam("sep", types.I8Ptr))
	listPtr := g.builder.NewCall(fn, s, sep)

	resultType := g.getResultType(listPtrTy)
	result := g.builder.NewAlloca(resultType)
	isErr := g.builder.NewICmp(enum.IPredEQ, listPtr, constant.NewNull(listPtrTy))
	disc := g.builder.NewSelect(isErr,
		constant.NewInt(types.I8, 1), constant.NewInt(types.I8, 0))
	val := g.builder.NewSelect(isErr, constant.NewNull(listPtrTy), listPtr)
	vp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(val, vp)
	dp := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(disc, dp)
	return result, nil
}
