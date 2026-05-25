package codegen

// Registry entries for List<T> and Map<K, V> builtins.
//
// Spec: compiler/spec/0012-Built-InFunctions.md#collection-functions.
//
// Adds the full surface listed in the spec — length, contains, append,
// prepend, concat, reverse, set, remove, merge, etc. — to the global
// BuiltInFunctionRegistry. Wiring point: one call to
// `r.registerListMapBuiltins()` from NewBuiltInFunctionRegistry.

// registerListMapBuiltins adds every collection builtin specified in
// compiler/spec/0012-Built-InFunctions.md#collection-functions.
func (r *BuiltInFunctionRegistry) registerListMapBuiltins() {
	r.registerListBuiltins()
	r.registerMapBuiltins()
}

func (r *BuiltInFunctionRegistry) registerListBuiltins() {
	list := &ConcreteType{name: TypeList}
	intT := &ConcreteType{name: TypeInt}
	anyT := &ConcreteType{name: TypeAny}
	funT := &ConcreteType{name: "T -> Unit"}

	listParam := BuiltInParameter{Name: "list", Type: list, Description: "The list"}

	// NOTE: list and map operations use explicit `listX` / `mapX` prefixes
	// in this revision to avoid name collisions with the string builtins
	// (`length`, `contains`, …). When UFCS (currently in flight in
	// docs/plans/string-manipulation.md) lands, `xs.length()` will
	// dispatch on the receiver type, and these can collapse to plain
	// `length` etc.
	r.functions["listLength"] = &BuiltInFunction{
		Name:           "listLength",
		Signature:      "listLength(list: List<T>) -> int",
		Description:    "Returns the number of elements in a list. O(1).",
		ParameterTypes: []BuiltInParameter{listParam},
		ReturnType:     intT,
		Category:       CategoryFunctional,
		Generator:      (*LLVMGenerator).generateListLengthCall,
		Example:        `listLength([1, 2, 3])  // 3`,
	}
	r.functions["listGet"] = &BuiltInFunction{
		Name:        "listGet",
		Signature:   "listGet(list: List<T>, index: int) -> Result<T, string>",
		Description: "Returns Success(value) for an in-range index, Error otherwise. O(log32 n).",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "index", Type: intT, Description: "Zero-based element index"},
		},
		ReturnType: anyT, // Result<T, string>; element type via inference.
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListGetCall,
		Example:    `match listGet(xs, 0) { Success v => v; Error e => 0 }`,
	}
	r.functions["listSet"] = &BuiltInFunction{
		Name:        "listSet",
		Signature:   "listSet(list: List<T>, index: int, value: T) -> List<T>",
		Description: "Returns a new list with the element at index replaced. Out-of-range is a no-op. O(log32 n).",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "index", Type: intT, Description: "Zero-based element index"},
			{Name: "value", Type: anyT, Description: "Replacement value"},
		},
		ReturnType: list,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListSetCall,
		Example:    `listSet([1, 2, 3], 1, 99)  // [1, 99, 3]`,
	}
	r.functions["listDrop"] = &BuiltInFunction{
		Name:        "listDrop",
		Signature:   "listDrop(list: List<T>, n: int) -> List<T>",
		Description: "Returns a new list with the first n elements removed. O(log32 n).",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "n", Type: intT, Description: "Number of leading elements to drop"},
		},
		ReturnType: list,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListDropCall,
		Example:    `listDrop([1, 2, 3, 4], 2)  // [3, 4]`,
	}
	r.functions["listAppend"] = &BuiltInFunction{
		Name:        "listAppend",
		Signature:   "listAppend(list: List<T>, value: T) -> List<T>",
		Description: "Returns a new list with value at the end. O(log32 n) amortised.",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "value", Type: anyT, Description: "Value to append"},
		},
		ReturnType: list,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListAppendCall,
		Example:    `listAppend([1, 2], 3)  // [1, 2, 3]`,
	}
	r.functions["listPrepend"] = &BuiltInFunction{
		Name:        "listPrepend",
		Signature:   "listPrepend(list: List<T>, value: T) -> List<T>",
		Description: "Returns a new list with value at the front. O(n).",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "value", Type: anyT, Description: "Value to prepend"},
		},
		ReturnType: list,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListPrependCall,
		Example:    `listPrepend([2, 3], 1)  // [1, 2, 3]`,
	}
	r.functions["listConcat"] = &BuiltInFunction{
		Name:        "listConcat",
		Signature:   "listConcat(left: List<T>, right: List<T>) -> List<T>",
		Description: "Returns left ++ right. Same as left + right.",
		ParameterTypes: []BuiltInParameter{
			{Name: "left", Type: list, Description: "Left operand"},
			{Name: "right", Type: list, Description: "Right operand"},
		},
		ReturnType: list,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListConcatCall,
		Example:    `listConcat([1, 2], [3, 4])  // [1, 2, 3, 4]`,
	}
	r.functions["listReverse"] = &BuiltInFunction{
		Name:           "listReverse",
		Signature:      "listReverse(list: List<T>) -> List<T>",
		Description:    "Returns a new list in reverse order.",
		ParameterTypes: []BuiltInParameter{listParam},
		ReturnType:     list,
		Category:       CategoryFunctional,
		Generator:      (*LLVMGenerator).generateListReverseCall,
		Example:        `listReverse([1, 2, 3])  // [3, 2, 1]`,
	}
	r.functions["forEachList"] = &BuiltInFunction{
		Name:        "forEachList",
		Signature:   "forEachList(list: List<T>, function: fn(T) -> Unit) -> List<T>",
		Description: "Apply function to every element of list. Phase 7 of collections plan.",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "function", Type: funT, Description: "Function applied per element"},
		},
		ReturnType:  list,
		Category:    CategoryFunctional,
		IsProtected: true,
		Generator:   (*LLVMGenerator).generateForEachListCall,
		Example:     `forEachList(xs, print)`,
	}
	// `contains` on List is registered here; the Map version below would
	// otherwise clash on name. The codegen helper inspects the inferred
	// receiver type to dispatch.
	r.functions["listContains"] = &BuiltInFunction{
		Name:        "listContains",
		Signature:   "listContains(list: List<T>, value: T) -> bool",
		Description: "True iff some element equals value. O(n).",
		ParameterTypes: []BuiltInParameter{
			listParam,
			{Name: "value", Type: anyT, Description: "Value to find"},
		},
		ReturnType: &ConcreteType{name: TypeBool},
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateListContainsCall,
		Example:    `listContains([1, 2, 3], 2)  // true`,
	}
}

func (r *BuiltInFunctionRegistry) registerMapBuiltins() {
	mp := &ConcreteType{name: TypeMap}
	intT := &ConcreteType{name: TypeInt}
	boolT := &ConcreteType{name: TypeBool}
	anyT := &ConcreteType{name: TypeAny}

	mapParam := BuiltInParameter{Name: "map", Type: mp, Description: "The map"}

	r.functions["mapLength"] = &BuiltInFunction{
		Name:           "mapLength",
		Signature:      "mapLength(map: Map<K, V>) -> int",
		Description:    "Returns the number of entries in a map. O(1).",
		ParameterTypes: []BuiltInParameter{mapParam},
		ReturnType:     intT,
		Category:       CategoryFunctional,
		Generator:      (*LLVMGenerator).generateMapLengthCall,
		Example:        `mapLength({"a": 1, "b": 2})  // 2`,
	}
	r.functions["mapContains"] = &BuiltInFunction{
		Name:        "mapContains",
		Signature:   "mapContains(map: Map<K, V>, key: K) -> bool",
		Description: "True iff key is present in map.",
		ParameterTypes: []BuiltInParameter{
			mapParam,
			{Name: "key", Type: anyT, Description: "Key to find"},
		},
		ReturnType: boolT,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateMapContainsCall,
		Example:    `mapContains({"a": 1}, "a")  // true`,
	}
	r.functions["mapGet"] = &BuiltInFunction{
		Name:        "mapGet",
		Signature:   "mapGet(map: Map<K, V>, key: K) -> Result<V, string>",
		Description: "Returns Success(value) when key is present, Error(...) otherwise.",
		ParameterTypes: []BuiltInParameter{
			mapParam,
			{Name: "key", Type: anyT, Description: "Key to look up"},
		},
		ReturnType: anyT, // Result<V, string>; resolved via type inference.
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateMapGetCall,
		Example:    `match mapGet(m, "k") { Success v => print(v); Error e => print("missing") }`,
	}
	r.functions["mapSet"] = &BuiltInFunction{
		Name:        "mapSet",
		Signature:   "mapSet(map: Map<K, V>, key: K, value: V) -> Map<K, V>",
		Description: "Returns a new map with key bound to value (replaces prior binding).",
		ParameterTypes: []BuiltInParameter{
			mapParam,
			{Name: "key", Type: anyT, Description: "Key"},
			{Name: "value", Type: anyT, Description: "Value"},
		},
		ReturnType: mp,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateMapSetCall,
		Example:    `mapSet({"a": 1}, "b", 2)  // {"a": 1, "b": 2}`,
	}
	r.functions["mapRemove"] = &BuiltInFunction{
		Name:        "mapRemove",
		Signature:   "mapRemove(map: Map<K, V>, key: K) -> Map<K, V>",
		Description: "Returns a new map without key. No-op if key is absent.",
		ParameterTypes: []BuiltInParameter{
			mapParam,
			{Name: "key", Type: anyT, Description: "Key"},
		},
		ReturnType: mp,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateMapRemoveCall,
		Example:    `mapRemove({"a": 1, "b": 2}, "a")  // {"b": 2}`,
	}
	r.functions["mapMerge"] = &BuiltInFunction{
		Name:        "mapMerge",
		Signature:   "mapMerge(left: Map<K, V>, right: Map<K, V>) -> Map<K, V>",
		Description: "Right-biased union. Same as left + right.",
		ParameterTypes: []BuiltInParameter{
			{Name: "left", Type: mp, Description: "Left"},
			{Name: "right", Type: mp, Description: "Right"},
		},
		ReturnType: mp,
		Category:   CategoryFunctional,
		Generator:  (*LLVMGenerator).generateMapMergeCall,
		Example:    `mapMerge({"a": 1}, {"b": 2})  // {"a": 1, "b": 2}`,
	}
	list := &ConcreteType{name: TypeList}
	r.functions["mapKeys"] = &BuiltInFunction{
		Name:           "mapKeys",
		Signature:      "mapKeys(map: Map<K, V>) -> List<K>",
		Description:    "All keys of the map as a list. Order unspecified.",
		ParameterTypes: []BuiltInParameter{mapParam},
		ReturnType:     list,
		Category:       CategoryFunctional,
		Generator:      (*LLVMGenerator).generateMapKeysCall,
		Example:        `mapKeys(m)  // List<K>`,
	}
	r.functions["mapValues"] = &BuiltInFunction{
		Name:           "mapValues",
		Signature:      "mapValues(map: Map<K, V>) -> List<V>",
		Description:    "All values of the map as a list. Order matches mapKeys.",
		ParameterTypes: []BuiltInParameter{mapParam},
		ReturnType:     list,
		Category:       CategoryFunctional,
		Generator:      (*LLVMGenerator).generateMapValuesCall,
		Example:        `mapValues(m)  // List<V>`,
	}
}
