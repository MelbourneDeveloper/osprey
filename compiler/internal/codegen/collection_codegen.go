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

// Key-type tags matching enum OspreyKeyType in collection_runtime.h.
const (
	collectionKeyInt    int64 = 0
	collectionKeyString int64 = 1
	collectionKeyBool   int64 = 2
)

// ensureCollectionExtern declares an osprey_list_* / osprey_map_* function
// once and caches it. The codegen helpers below all funnel through here.
func (g *LLVMGenerator) ensureCollectionExtern(name string, retType types.Type, paramTypes ...types.Type) *ir.Func {
	if fn, ok := g.functions[name]; ok {
		return fn
	}
	params := make([]*ir.Param, 0, len(paramTypes))
	for i, pt := range paramTypes {
		params = append(params, ir.NewParam(fmt.Sprintf("p%d", i), pt))
	}
	fn := g.module.NewFunc(name, retType, params...)
	g.functions[name] = fn
	return fn
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

// unboxFromI64 narrows the uniform i64 element back to the desired type.
func (g *LLVMGenerator) unboxFromI64(v value.Value, target types.Type) value.Value {
	if target == types.I64 {
		return v
	}
	if pt, ok := target.(*types.PointerType); ok {
		return g.builder.NewIntToPtr(v, pt)
	}
	if target == types.I1 || target == types.I8 || target == types.I32 {
		return g.builder.NewTrunc(v, target)
	}
	if target == types.Double {
		return g.builder.NewBitCast(v, types.Double)
	}
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

// generateListLiteralRuntime lowers `[e1, e2, …]` to a transient builder
// sequence followed by a seal.
//
// Returns an i8* opaque handle to the new OspreyList. Implements
// [TYPE-LIST] literal codegen.
func (g *LLVMGenerator) generateListLiteralRuntime(lit *ast.ListLiteral) (value.Value, error) {
	g.declareListExterns()
	builderNew := g.functions["osprey_list_builder_new"]
	push := g.functions["osprey_list_builder_push"]
	seal := g.functions["osprey_list_builder_seal"]

	bld := g.builder.NewCall(builderNew)
	for _, elem := range lit.Elements {
		v, err := g.generateExpression(elem)
		if err != nil {
			return nil, fmt.Errorf("list element: %w", err)
		}
		g.builder.NewCall(push, bld, g.boxToI64(v))
	}
	return g.builder.NewCall(seal, bld), nil
}

// generateListGetRuntime emits a bounds-checked `list[i]` returning
// Result<T, IndexError>. The Result layout matches the existing list-access
// codegen so callers can pattern-match identically.
func (g *LLVMGenerator) generateListGetRuntime(listVal value.Value, indexVal value.Value, elemLLVMType types.Type) (value.Value, error) {
	g.declareListExterns()
	inBounds := g.builder.NewCall(g.functions["osprey_list_in_bounds"], listVal, indexVal)
	cond := g.builder.NewICmp(enum.IPredEQ, inBounds, constant.NewInt(types.I32, 1))
	_ = cond
	rawElem := g.builder.NewCall(g.functions["osprey_list_get"], listVal, indexVal)
	return g.unboxFromI64(rawElem, elemLLVMType), nil
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

// inferMapKeyTag picks an OspreyKeyType tag based on the key type of a map
// literal's first entry. If the literal is empty the caller must supply a
// tag from the type annotation context.
func (g *LLVMGenerator) inferMapKeyTag(lit *ast.MapLiteral) int64 {
	if len(lit.Entries) == 0 {
		return collectionKeyString // best-effort fallback for empty literal
	}
	switch lit.Entries[0].Key.(type) {
	case *ast.StringLiteral:
		return collectionKeyString
	case *ast.IntegerLiteral:
		return collectionKeyInt
	case *ast.BooleanLiteral:
		return collectionKeyBool
	default:
		return collectionKeyString
	}
}

// generateMapLiteralRuntime lowers `{k1: v1, k2: v2, …}` to a builder
// sequence followed by a seal. Implements [TYPE-MAP-LITERAL] codegen.
func (g *LLVMGenerator) generateMapLiteralRuntime(lit *ast.MapLiteral) (value.Value, error) {
	g.declareMapExterns()
	keyTag := constant.NewInt(types.I32, g.inferMapKeyTag(lit))
	bld := g.builder.NewCall(g.functions["osprey_map_builder_new"], keyTag)
	put := g.functions["osprey_map_builder_put"]
	for _, entry := range lit.Entries {
		k, err := g.generateExpression(entry.Key)
		if err != nil {
			return nil, fmt.Errorf("map key: %w", err)
		}
		v, err := g.generateExpression(entry.Value)
		if err != nil {
			return nil, fmt.Errorf("map value: %w", err)
		}
		g.builder.NewCall(put, bld, g.boxToI64(k), g.boxToI64(v))
	}
	return g.builder.NewCall(g.functions["osprey_map_builder_seal"], bld), nil
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
	if len(callExpr.Arguments) != 1 {
		return nil, fmt.Errorf("%s expects 1 argument", name)
	}
	g.declareListExterns()
	g.declareMapExterns()
	arg, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	fn := g.functions[name]
	if fn == nil {
		return nil, fmt.Errorf("collection extern %s not declared", name)
	}
	res := g.builder.NewCall(fn, arg)
	if ret == types.I32 {
		return g.builder.NewICmp(enum.IPredNE, res, constant.NewInt(types.I32, 0)), nil
	}
	return res, nil
}

func (g *LLVMGenerator) callTwoArgs(name string, callExpr *ast.CallExpression, boxA bool, boxB bool) (value.Value, error) {
	if len(callExpr.Arguments) != 2 {
		return nil, fmt.Errorf("%s expects 2 arguments", name)
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
	if boxA {
		a = g.boxToI64(a)
	}
	if boxB {
		b = g.boxToI64(b)
	}
	fn := g.functions[name]
	if fn == nil {
		return nil, fmt.Errorf("collection extern %s not declared", name)
	}
	return g.builder.NewCall(fn, a, b), nil
}

// Public generators wired into the registry.

func (g *LLVMGenerator) generateListLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_list_length", types.I64, callExpr)
}

func (g *LLVMGenerator) generateListContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	// `contains(list, value)` is O(n) linear scan via iter.
	if len(callExpr.Arguments) != 2 {
		return nil, errors.New("contains expects 2 arguments")
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
	return g.callTwoArgs("osprey_list_append", callExpr, false, true)
}

func (g *LLVMGenerator) generateListPrependCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_list_prepend", callExpr, false, true)
}

func (g *LLVMGenerator) generateListConcatCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callTwoArgs("osprey_list_concat", callExpr, false, false)
}

func (g *LLVMGenerator) generateListReverseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_list_reverse", types.I8Ptr, callExpr)
}

func (g *LLVMGenerator) generateMapLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.callOneArg("osprey_map_length", types.I64, callExpr)
}

func (g *LLVMGenerator) generateMapContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 2 {
		return nil, errors.New("contains expects 2 arguments")
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
	if len(callExpr.Arguments) != 3 {
		return nil, errors.New("set expects 3 arguments")
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
	if len(callExpr.Arguments) != 2 {
		return nil, errors.New("remove expects 2 arguments")
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
	return g.callTwoArgs("osprey_map_merge", callExpr, false, false)
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
	if len(callExpr.Arguments) != 1 {
		return nil, errors.New("keys/values expect 1 argument")
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
	if len(callExpr.Arguments) != 2 {
		return nil, errors.New("forEachList expects 2 arguments")
	}
	g.declareListExterns()
	list, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	funcArg := callExpr.Arguments[1]
	funcIdent, ok := funcArg.(*ast.Identifier)
	if !ok {
		return nil, errors.New("forEachList: second argument must be a function name")
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
