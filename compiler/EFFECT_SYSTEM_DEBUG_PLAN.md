# ALGEBRAIC EFFECTS DEBUGGING PLAN 🔥

## STATUS: 10/12 Effects Tests Passing ✅

**FAILING TESTS:**
1. `comprehensive_bugs_demo` - **SYNTAX ERROR** (parser issue)
2. `nested_handler_scoping_bug` - **LEXICAL SCOPING BUG** (handler resolution issue)

---

## 🚨 CRITICAL BUG: Lexical Scoping (`nested_handler_scoping_bug`)

### THE PROBLEM
Functions with declared effects (e.g., `fn innerOperation() !Logger`) are **resolving handlers at function definition time instead of call time**. This creates hardcoded calls to specific handlers instead of dynamic resolution.

**Expected Output:**
```
[OUTER-LOG] Before inner handler
[INNER-LOG] Inner operation executing    ← Should use INNER handler
[INNER-DEBUG] Debug info from inner     ← Should use INNER handler  
[INNER-LOG] Inner result: 42            ← Should use INNER handler
[OUTER-LOG] After inner handler         ← Should use OUTER handler
```

**Actual Output:**
```
[OUTER-LOG] Before inner handler
[OUTER-LOG] Inner operation executing    ← WRONG! Using OUTER handler
[OUTER-DEBUG] Debug info from inner     ← WRONG! Using OUTER handler
[INNER-LOG] Inner result: 42            ← Correct
[INNER-LOG] After inner handler         ← WRONG! Using INNER handler
```

### ROOT CAUSE
The `generateDeclaredEffectCall()` function needs **TRUE RUNTIME HANDLER LOOKUP** that searches the active handler stack when the function is executed, not when it's defined.

### SOLUTION APPROACH
1. **Functions with declared effects should use CURRENT ACTIVE HANDLERS**
2. **Dynamic resolution at call time, not definition time**
3. **Evidence passing system needs to be operation-specific**

---

## 🛠️ GOOD CODE ALREADY APPLIED FROM `dd4e0ed1e9cc078ae07b04a2be591329855088b5`

### ✅ APPLIED FIXES:
- **Lexical depth tracking** (`LexicalDepth` field in `HandlerFrame`)
- **Improved handler resolution** (innermost handler wins by depth)
- **Better type inference** (operation-specific parameter types)
- **Evidence passing constants** (`OpLog`, `OpGet`, `OpSet`, etc.)
- **Safety checks** (null handler validation)
- **Proper handler stack restoration** (both `currentHandlers` and `handlerStack`)

### 🔧 KEY IMPROVEMENTS FROM THAT COMMIT:
```go
// GOOD: Lexical depth tracking
type HandlerFrame struct {
    EffectName   string
    Operations   map[string]*ir.Func
    Continuation *ir.Func
    LexicalDepth int // ← This is critical
}

// GOOD: Evidence passing fields
type EffectCodegen struct {
    // ... existing fields ...
    currentEvidenceParams map[string]*ir.Param  // ← For operation-specific evidence
    currentLexicalDepth   int                   // ← Track nesting depth
}

// GOOD: Handler resolution prioritizes handlers over declared effects
if len(ec.currentHandlers) > 0 || len(ec.handlerStack) > 0 {
    if result, err := ec.tryCurrentScopeHandlers(perform); err != nil || result != nil {
        return result, err
    }
    if result, err := ec.tryStackHandlers(perform); err != nil || result != nil {
        return result, err
    }
}
```

---

## 🎯 NEXT STEPS TO FIX LEXICAL SCOPING

### STEP 1: Debug Handler Resolution
Add **temporary diagnostic prints** (remove after debugging):
```go
// In generateDeclaredEffectCall()
fmt.Printf("DEBUG: Function %s.%s - currentHandlers=%d, stack=%d\n", 
    perform.EffectName, perform.OperationName, len(ec.currentHandlers), len(ec.handlerStack))

// In findHandlerByEffectName()  
fmt.Printf("DEBUG: Searching for %s.%s, found handler at depth %d\n",
    effectName, perform.OperationName, frame.LexicalDepth)
```

### STEP 2: Fix Evidence Parameter Resolution
The issue might be that **operation-specific evidence parameters** aren't being created/used correctly:

```go
// In generateDeclaredEffectCall() - try this approach:
evidenceParamName := fmt.Sprintf("%s_%s", perform.EffectName, perform.OperationName)
if evidenceParam, exists := ec.currentEvidenceParams[evidenceParamName]; exists {
    // Use evidence parameter directly - NO LOOKUP!
    return ec.generator.builder.NewCall(evidenceParam, args...), nil
}
```

### STEP 3: Implement True Runtime Handler Lookup
If evidence passing doesn't work, implement **dynamic handler resolution**:
```go
// Functions with declared effects should ALWAYS check current handlers first:
// PRIORITY 1: Current lexical scope
// PRIORITY 2: Handler stack  
// PRIORITY 3: Runtime lookup (fallback)
```

---

## 🔧 SYNTAX ERROR: `comprehensive_bugs_demo` 

### THE PROBLEM
Parser error on line 82 - malformed handler syntax:
```
line 82:22 no viable alternative at input 'withhandlerLoggerlog(msg)=>print...'
```

### SOLUTION
**Replace the broken file with the working version from `dd4e0ed1e9cc078ae07b04a2be591329855088b5`:**

```bash
git show dd4e0ed1e9cc078ae07b04a2be591329855088b5:compiler/examples/tested/effects/comprehensive_bugs_demo.osp > examples/tested/effects/comprehensive_bugs_demo.osp
```

The working version has proper syntax with correct spacing and formatting.

---

## 🧪 DEBUGGING STRATEGY

### Create Isolated Unit Test
```go
func TestLexicalHandlerResolution(t *testing.T) {
    // Test that functions with declared effects use current active handlers
    // This will help isolate the problem without running full integration tests
}
```

### Diagnostic Questions to Answer:
1. **Are handlers being created with correct lexical depth?**
2. **Are functions with declared effects checking current handlers?**
3. **Is the evidence parameter system working correctly?**
4. **Are handlers being restored properly after scope ends?**

---

## 📚 REFERENCES

- **Memory ID 1275188797194468042**: Root cause is static handler resolution during function definition
- **Memory ID 2134669719688179820**: Evidence passing needs to be operation-specific  
- **Commit dd4e0ed1e9cc078ae07b04a2be591329855088b5**: Contains working syntax and better handler resolution logic

---

## 🎯 SUCCESS CRITERIA

**Test should pass when:**
1. `innerOperation()` calls use **INNER** handler (depth 2) 
2. `outerOperation()` calls use **OUTER** handler (depth 1)
3. Handler scopes are properly restored after each block
4. All syntax errors are resolved

**Key insight:** The bug is that **declared effect functions are compiled with hardcoded handler references instead of dynamic lookups**. 