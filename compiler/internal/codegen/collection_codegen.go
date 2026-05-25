package codegen

// LLVM IR generation for List<T> and Map<K, V> backed by the C runtime in
// compiler/runtime/list_runtime.c and compiler/runtime/map_runtime.c.
//
// Implements [TYPE-LIST-OPS], [TYPE-MAP-OPS] from
// compiler/spec/0004-TypeSystem.md and the builtins in
// compiler/spec/0012-Built-InFunctions.md#collection-functions.
//
// All collection handles cross the LLVM/C boundary as opaque i8*; the C
// runtime owns layout. Element values cross as i64; pointers (strings,
// nested collections, records) are bit-cast to i64 at the call site.

import (
	"errors"
	"fmt"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Key-type tag matching enum OspreyKeyType in collection_runtime.h. Only
// the string tag is currently referenced from codegen; the int/bool tags
// will be reintroduced when typed-map literals land.
const collectionKeyString int64 = 1

// Argument-count constants for collection builtins. Named per builtin so
// the call site reads as a domain assertion rather than a magic literal.
const (
	collectionArgsOne   = 1
	collectionArgsTwo   = 2
	collectionMapSetArg = 3
)

// Static errors satisfying err113 — wrap with fmt.Errorf("…: %w", err) when
// the caller needs to attach a name.
var (
	errCollectionArgCount   = errors.New("wrong argument count")
	errCollectionExternMiss = errors.New("collection extern not declared")
	errForEachListSecondArg = errors.New("forEachList: second argument must be a function name")
)

// ensureCollectionExtern declares an osprey_list_* / osprey_map_* function
// once and caches it. The codegen helpers below all funnel through here.
func (g *LLVMGenerator) ensureCollectionExtern(name string, retType types.Type, paramTypes ...types.Type) {
	if _, ok := g.functions[name]; ok {
		return
	}
	params := make([]*ir.Param, 0, len(paramTypes))
	for i, pt := range paramTypes {
		params = append(params, ir.NewParam(fmt.Sprintf("p%d", i), pt))
	}
	g.functions[name] = g.module.NewFunc(name, retType, params...)
}

// boxToI64 widens a non-i64 value (i8*, i32, i1, etc.) to i64 for the
// uniform collection element ABI. i64 passes through.
func (g *LLVMGenerator) boxToI64(v value.Value) value.Value {
	t := v.Type()
	if t == types.I64 {
		return v
	}
	if pt, ok := t.(*types.PointerType); ok {
		_ = pt
		return g.builder.NewPtrToInt(v, types.I64)
	}
	if t == types.I1 || t == types.I8 || t == types.I32 {
		return g.builder.NewZExt(v, types.I64)
	}
	if t == types.Double {
		return g.builder.NewBitCast(v, types.I64)
	}
	// Fallback: assume it's already pointer-sized or codegen elsewhere
	// has converted it.
	return v
}

// ============================================================================
//                              List codegen
// ============================================================================

// declareListExterns ensures every osprey_list_* function is declared in
// the module. Idempotent; safe to call from any codegen entry point.
func (g *LLVMGenerator) declareListExterns() {
	g.ensureCollectionExtern("osprey_list_empty", types.I8Ptr)
	g.ensureCollectionExtern("osprey_list_length", types.I64, types.I8Ptr)
	g.ensureCollectionExtern("osprey_list_in_bounds", types.I32, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_get", types.I64, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_set", types.I8Ptr, types.I8Ptr, types.I64, types.I64)
	g.ensureCollectionExtern("osprey_list_append", types.I8Ptr, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_prepend", types.I8Ptr, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_concat", types.I8Ptr, types.I8Ptr, types.I8Ptr)
	g.ensureCollectionExtern("osprey_list_drop", types.I8Ptr, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_reverse", types.I8Ptr, types.I8Ptr)
	g.ensureCollectionExtern("osprey_list_builder_new", types.I8Ptr)
	g.ensureCollectionExtern("osprey_list_builder_push", types.Void, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_list_builder_seal", types.I8Ptr, types.I8Ptr)
}

// ============================================================================
//                              Map codegen
// ============================================================================

func (g *LLVMGenerator) declareMapExterns() {
	g.ensureCollectionExtern("osprey_map_empty", types.I8Ptr, types.I32)
	g.ensureCollectionExtern("osprey_map_length", types.I64, types.I8Ptr)
	g.ensureCollectionExtern("osprey_map_contains", types.I32, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_map_get", types.I64, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_map_set", types.I8Ptr, types.I8Ptr, types.I64, types.I64)
	g.ensureCollectionExtern("osprey_map_remove", types.I8Ptr, types.I8Ptr, types.I64)
	g.ensureCollectionExtern("osprey_map_merge", types.I8Ptr, types.I8Ptr, types.I8Ptr)
	g.ensureCollectionExtern("osprey_map_builder_new", types.I8Ptr, types.I32)
	g.ensureCollectionExtern("osprey_map_builder_put", types.Void, types.I8Ptr, types.I64, types.I64)
	g.ensureCollectionExtern("osprey_map_builder_seal", types.I8Ptr, types.I8Ptr)
}

// ============================================================================
//                         Builtin generators (Phase 3)
// ============================================================================
//
// Each function below is a generator suitable for plugging into the
// BuiltInFunctionRegistry as a *BuiltInFunction.Generator field. They share
// a small evaluation shape: evaluate args, declare the runtime extern,
// emit the call.

func (g *LLVMGenerator) callOneArg(name string, ret types.Type, callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsOne {
		return nil, fmt.Errorf("%s expects %d argument: %w", name, collectionArgsOne, errCollectionArgCount)
	}
	g.declareListExterns()
	g.declareMapExterns()
	arg, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.functions[name]
	if fn == nil {
		return nil, fmt.Errorf("%s: %w", name, errCollectionExternMiss)
	}
	arg = g.coerceToI8Ptr(arg)
	res := g.builder.NewCall(fn, arg)
	if ret == types.I32 {
		return g.builder.NewICmp(enum.IPredNE, res, constant.NewInt(types.I32, 0)), nil
	}
	return res, nil
}

// coerceToI8Ptr bitcasts a pointer value to i8* so it can be passed to a
// collection extern that takes an opaque handle. No-op if already i8* or
// not a pointer. Needed when an expression returns a typed pointer
// (e.g. osp_string_list*) that the runtime treats as an opaque list handle.
//
// Result-wrapped lists (`split(...)` returns Result<osp_string_list*>)
// are auto-unwrapped first per spec 0004-TypeSystem.md — without this,
// `listLength(splitResult)` bit-cast the *Result struct* pointer
// straight to i8* and the runtime then read garbage from the
// {ptr, disc} layout, returning a huge wrong length.
func (g *LLVMGenerator) coerceToI8Ptr(v value.Value) value.Value {
	if v.Type() == types.I8Ptr {
		return v
	}
	const errorFlagBitSize = 8
	if ptrTy, isPtr := v.Type().(*types.PointerType); isPtr {
		if structTy, isStruct := ptrTy.ElemType.(*types.StructType); isStruct {
			if len(structTy.Fields) == ResultFieldCount {
				if intTy, ok := structTy.Fields[1].(*types.IntType); ok && intTy.BitSize == errorFlagBitSize {
					vp := g.builder.NewGetElementPtr(structTy, v,
						constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
					v = g.builder.NewLoad(structTy.Fields[0], vp)
				}
			}
		}
	}
	if v.Type() == types.I8Ptr {
		return v
	}
	if _, isPtr := v.Type().(*types.PointerType); isPtr {
		return g.builder.NewBitCast(v, types.I8Ptr)
	}
	return v
}

// callTwoArgs evaluates two arguments, optionally boxes the second to i64
// (the first is always a collection handle, already i8*), and calls the
// named extern. Returns its result.
func (g *LLVMGenerator) callTwoArgs(name string, callExpr *ast.CallExpression, boxB bool) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsTwo {
		return nil, fmt.Errorf("%s expects %d arguments: %w", name, collectionArgsTwo, errCollectionArgCount)
	}
	g.declareListExterns()
	g.declareMapExterns()
	a, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	b, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	if boxB {
		b = g.boxToI64(b)
	}
	a = g.coerceToI8Ptr(a)
	fn := g.functions[name]
	if fn == nil {
		return nil, fmt.Errorf("%s: %w", name, errCollectionExternMiss)
	}
	return g.builder.NewCall(fn, a, b), nil
}

// Public generators wired into the registry.

func (g *LLVMGenerator) generateListLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_list_length", types.I64, callExpr)
}

func (g *LLVMGenerator) generateListContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	// `contains(list, value)` is O(n) linear scan via iter.
	if len(callExpr.Arguments) != collectionArgsTwo {
		return nil, fmt.Errorf("listContains expects %d arguments: %w", collectionArgsTwo, errCollectionArgCount)
	}
	g.declareListExterns()
	list, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	needle, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	// Lower to a runtime loop: for i in 0..length: if get(i) == needle: return 1
	lenFn := g.functions["osprey_list_length"]
	getFn := g.functions["osprey_list_get"]
	listLen := g.builder.NewCall(lenFn, list)
	idxPtr := g.builder.NewAlloca(types.I64)
	g.builder.NewStore(constant.NewInt(types.I64, 0), idxPtr)
	resPtr := g.builder.NewAlloca(types.I1)
	g.builder.NewStore(constant.NewInt(types.I1, 0), resPtr)

	loop := g.function.NewBlock(fmt.Sprintf("contains_loop_%p", callExpr))
	body := g.function.NewBlock(fmt.Sprintf("contains_body_%p", callExpr))
	done := g.function.NewBlock(fmt.Sprintf("contains_done_%p", callExpr))
	g.builder.NewBr(loop)

	g.builder = loop
	idxLoad := g.builder.NewLoad(types.I64, idxPtr)
	cont := g.builder.NewICmp(enum.IPredSLT, idxLoad, listLen)
	g.builder.NewCondBr(cont, body, done)

	g.builder = body
	elem := g.builder.NewCall(getFn, list, idxLoad)
	eq := g.builder.NewICmp(enum.IPredEQ, elem, g.boxToI64(needle))
	foundBlock := g.function.NewBlock(fmt.Sprintf("contains_found_%p", callExpr))
	contBlock := g.function.NewBlock(fmt.Sprintf("contains_cont_%p", callExpr))
	g.builder.NewCondBr(eq, foundBlock, contBlock)

	g.builder = foundBlock
	g.builder.NewStore(constant.NewInt(types.I1, 1), resPtr)
	g.builder.NewBr(done)

	g.builder = contBlock
	next := g.builder.NewAdd(idxLoad, constant.NewInt(types.I64, 1))
	g.builder.NewStore(next, idxPtr)
	g.builder.NewBr(loop)

	g.builder = done
	return g.builder.NewLoad(types.I1, resPtr), nil
}

func (g *LLVMGenerator) generateListAppendCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_list_append", callExpr, true)
}

func (g *LLVMGenerator) generateListPrependCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_list_prepend", callExpr, true)
}

func (g *LLVMGenerator) generateListConcatCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_list_concat", callExpr, false)
}

func (g *LLVMGenerator) generateListReverseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_list_reverse", types.I8Ptr, callExpr)
}

func (g *LLVMGenerator) generateMapLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_map_length", types.I64, callExpr)
}

func (g *LLVMGenerator) generateMapContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsTwo {
		return nil, fmt.Errorf("mapContains expects %d arguments: %w", collectionArgsTwo, errCollectionArgCount)
	}
	g.declareMapExterns()
	m, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	k, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	res := g.builder.NewCall(g.functions["osprey_map_contains"], m, g.boxToI64(k))
	return g.builder.NewICmp(enum.IPredNE, res, constant.NewInt(types.I32, 0)), nil
}

func (g *LLVMGenerator) generateMapSetCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != collectionMapSetArg {
		return nil, fmt.Errorf("mapSet expects %d arguments: %w", collectionMapSetArg, errCollectionArgCount)
	}
	g.declareMapExterns()
	m, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	k, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	v, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}
	return g.builder.NewCall(g.functions["osprey_map_set"], m, g.boxToI64(k), g.boxToI64(v)), nil
}

func (g *LLVMGenerator) generateMapRemoveCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsTwo {
		return nil, fmt.Errorf("mapRemove expects %d arguments: %w", collectionArgsTwo, errCollectionArgCount)
	}
	g.declareMapExterns()
	m, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	k, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}
	return g.builder.NewCall(g.functions["osprey_map_remove"], m, g.boxToI64(k)), nil
}

func (g *LLVMGenerator) generateMapMergeCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_map_merge", callExpr, false)
}

// generateMapKeysCall returns a List<K> containing every key in the map.
// Phase 3.4 — implements [TYPE-MAP-CONV].
func (g *LLVMGenerator) generateMapKeysCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateMapToListCall(callExpr, true /* take key */)
}

// generateMapValuesCall returns a List<V> containing every value.
func (g *LLVMGenerator) generateMapValuesCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateMapToListCall(callExpr, false /* take value */)
}

// generateMapToListCall is the shared body for keys(map)/values(map). It
// iterates via osprey_map_iter_next and builds a List via the transient
// builder.
func (g *LLVMGenerator) generateMapToListCall(callExpr *ast.CallExpression, takeKey bool) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsOne {
		return nil, fmt.Errorf("keys/values expect %d argument: %w", collectionArgsOne, errCollectionArgCount)
	}
	g.declareMapExterns()
	g.declareListExterns()
	g.ensureCollectionExtern("osprey_map_iter_new", types.I8Ptr, types.I8Ptr)
	g.ensureCollectionExtern("osprey_map_iter_next", types.I32,
		types.I8Ptr, types.NewPointer(types.I64), types.NewPointer(types.I64))

	m, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	bld := g.builder.NewCall(g.functions["osprey_list_builder_new"])
	iter := g.builder.NewCall(g.functions["osprey_map_iter_new"], m)
	keyPtr := g.builder.NewAlloca(types.I64)
	valPtr := g.builder.NewAlloca(types.I64)

	loopCond := g.function.NewBlock(fmt.Sprintf("map_to_list_cond_%p", callExpr))
	loopBody := g.function.NewBlock(fmt.Sprintf("map_to_list_body_%p", callExpr))
	loopEnd := g.function.NewBlock(fmt.Sprintf("map_to_list_end_%p", callExpr))
	g.builder.NewBr(loopCond)

	g.builder = loopCond
	hasNext := g.builder.NewCall(g.functions["osprey_map_iter_next"], iter, keyPtr, valPtr)
	cond := g.builder.NewICmp(enum.IPredNE, hasNext, constant.NewInt(types.I32, 0))
	g.builder.NewCondBr(cond, loopBody, loopEnd)

	g.builder = loopBody
	var elem value.Value
	if takeKey {
		elem = g.builder.NewLoad(types.I64, keyPtr)
	} else {
		elem = g.builder.NewLoad(types.I64, valPtr)
	}
	g.builder.NewCall(g.functions["osprey_list_builder_push"], bld, elem)
	g.builder.NewBr(loopCond)

	g.builder = loopEnd
	return g.builder.NewCall(g.functions["osprey_list_builder_seal"], bld), nil
}

// generateForEachListCall iterates a List<T> via the runtime iterator and
// calls a user function for each element. Phase 7 — implements [TYPE-LIST]
// integration with the iterator family.
func (g *LLVMGenerator) generateForEachListCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != collectionArgsTwo {
		return nil, fmt.Errorf("forEachList expects %d arguments: %w", collectionArgsTwo, errCollectionArgCount)
	}
	g.declareListExterns()
	list, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	// List literals (`[1, 2, 3]`) lower to an inline `{i64, i8*}` struct
	// on the stack, but forEachList expects an opaque `i8*` handle backed
	// by the C runtime's OspreyList. Without this guard the runtime
	// segfaults at osprey_list_length. Use listAppend(List(), …) instead
	// for now until [TYPE-LIST-LITERAL] is rewired to the C runtime.
	if list.Type() != types.I8Ptr {
		return nil, fmt.Errorf("%w: forEachList currently requires a "+
			"runtime-allocated list (built with List() / listAppend / split / ...). "+
			"List literals like [1, 2, 3] aren't wired through the C list runtime yet",
			errCollectionExternMiss)
	}
	// Accept either a named-fn identifier or an inline lambda (matches the
	// resolveCallbackIdent helper used by forEach / map / filter / fold).
	funcArg := callExpr.Arguments[1]
	funcIdent, err := g.resolveCallbackIdent(funcArg, errForEachListSecondArg)
	if err != nil {
		return nil, err
	}

	// Use osprey_list_length + osprey_list_get for a simple counted loop.
	// (osprey_list_iter is also available; counted-loop is friendlier to
	// fusion because we already do that for ranges.)
	lenFn := g.functions["osprey_list_length"]
	getFn := g.functions["osprey_list_get"]
	listLen := g.builder.NewCall(lenFn, list)
	idxPtr := g.builder.NewAlloca(types.I64)
	g.builder.NewStore(constant.NewInt(types.I64, 0), idxPtr)

	loopCond := g.function.NewBlock(fmt.Sprintf("foreach_list_cond_%p", callExpr))
	loopBody := g.function.NewBlock(fmt.Sprintf("foreach_list_body_%p", callExpr))
	loopEnd := g.function.NewBlock(fmt.Sprintf("foreach_list_end_%p", callExpr))
	g.builder.NewBr(loopCond)

	g.builder = loopCond
	idx := g.builder.NewLoad(types.I64, idxPtr)
	cond := g.builder.NewICmp(enum.IPredSLT, idx, listLen)
	g.builder.NewCondBr(cond, loopBody, loopEnd)

	g.builder = loopBody
	elem := g.builder.NewCall(getFn, list, idx)
	_, err = g.callFunctionWithValue(funcIdent, elem)
	if err != nil {
		return nil, err
	}
	next := g.builder.NewAdd(idx, constant.NewInt(types.I64, 1))
	g.builder.NewStore(next, idxPtr)
	g.builder.NewBr(loopCond)

	g.builder = loopEnd
	return list, nil
}
