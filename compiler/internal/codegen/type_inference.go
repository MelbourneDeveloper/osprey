package codegen

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/christianfindlay/osprey/internal/ast"
)

// Type represents a type in our type system
type Type interface {
	String() string
}

// TypeVar represents a type variable
type TypeVar struct {
	id int
}

func (tv *TypeVar) String() string {
	return fmt.Sprintf("t%d", tv.id)
}

// ConcreteType represents a concrete type like Int, String, etc.
type ConcreteType struct {
	name string
}

// NewConcreteType creates a new concrete type
func NewConcreteType(name string) *ConcreteType {
	return &ConcreteType{name: name}
}

func (ct *ConcreteType) String() string {
	return ct.name
}

// FunctionType represents a function type with parameters and return type
type FunctionType struct {
	paramTypes []Type
	returnType Type
}

// NewFunctionType creates a new function type
func NewFunctionType(paramTypes []Type, returnType Type) *FunctionType {
	return &FunctionType{
		paramTypes: paramTypes,
		returnType: returnType,
	}
}

func (ft *FunctionType) String() string {
	if len(ft.paramTypes) == 0 {
		return fmt.Sprintf("() -> %s", ft.returnType)
	}
	params := ""
	for i, p := range ft.paramTypes {
		if i > 0 {
			params += ", "
		}
		params += p.String()
	}
	return fmt.Sprintf("(%s) -> %s", params, ft.returnType)
}

// TypeScheme represents a polymorphic type with quantified variables
type TypeScheme struct {
	vars []int
	typ  Type
}

func (ts *TypeScheme) String() string {
	if len(ts.vars) == 0 {
		return ts.typ.String()
	}
	vars := ""
	for i, v := range ts.vars {
		if i > 0 {
			vars += ", "
		}
		vars += fmt.Sprintf("t%d", v)
	}
	return fmt.Sprintf("∀%s. %s", vars, ts.typ)
}

// TypeEnv represents the typing environment
type TypeEnv struct {
	vars map[string]Type
}

// NewTypeEnv creates a new type environment
func NewTypeEnv() *TypeEnv {
	return &TypeEnv{
		vars: make(map[string]Type),
	}
}

// Get retrieves a type from the environment by name
func (env *TypeEnv) Get(name string) (Type, bool) {
	t, exists := env.vars[name]
	return t, exists
}

// Set adds or updates a type in the environment
func (env *TypeEnv) Set(name string, t Type) {
	env.vars[name] = t
}

// GetAllVars returns a copy of all variables in the environment
func (env *TypeEnv) GetAllVars() map[string]Type {
	result := make(map[string]Type)
	for k, v := range env.vars {
		result[k] = v
	}
	return result
}

// Clone creates a deep copy of the type environment
func (env *TypeEnv) Clone() *TypeEnv {
	newEnv := NewTypeEnv()
	for k, v := range env.vars {
		newEnv.Set(k, v)
	}
	return newEnv
}

// Substitution represents a mapping of type variables to types
type Substitution map[int]Type

// TypeInferer handles type inference
type TypeInferer struct {
	nextID int
	env    *TypeEnv
	subst  Substitution
}

// NewTypeInferer creates a new type inferer with built-in functions initialized
func NewTypeInferer() *TypeInferer {
	ti := &TypeInferer{
		nextID: 0,
		env:    NewTypeEnv(),
		subst:  make(Substitution),
	}

	// Initialize built-in functions
	ti.initializeBuiltInFunctions()

	return ti
}

// Fresh creates a new type variable
func (ti *TypeInferer) Fresh() *TypeVar {
	tv := &TypeVar{id: ti.nextID}
	ti.nextID++
	return tv
}

// Instantiate creates a fresh instance of a type scheme
func (ti *TypeInferer) Instantiate(scheme *TypeScheme) Type {
	if len(scheme.vars) == 0 {
		return scheme.typ
	}
	subst := make(Substitution)
	for _, v := range scheme.vars {
		subst[v] = ti.Fresh()
	}
	return ti.applySubst(scheme.typ, subst)
}

// Generalize creates a type scheme by quantifying free type variables
func (ti *TypeInferer) Generalize(t Type) *TypeScheme {
	freeVars := ti.getFreeVars(t)

	// Remove vars that are already in the environment
	var schemeVars []int
	for _, v := range freeVars {
		inEnv := false
		for _, envType := range ti.env.GetAllVars() {
			envFreeVars := ti.getFreeVars(envType)
			for _, envVar := range envFreeVars {
				if envVar == v {
					inEnv = true
					break
				}
			}
			if inEnv {
				break
			}
		}
		if !inEnv {
			schemeVars = append(schemeVars, v)
		}
	}

	return &TypeScheme{vars: schemeVars, typ: t}
}

// Unify performs type unification between two types
func (ti *TypeInferer) Unify(t1, t2 Type) error {
	t1 = ti.prune(t1)
	t2 = ti.prune(t2)

	// Handle type variables
	if err := ti.unifyTypeVariables(t1, t2); !errors.Is(err, ErrNotTypeVariable) {
		return err
	}

	// Handle concrete types
	if err := ti.unifyConcreteTypes(t1, t2); !errors.Is(err, ErrNotConcreteType) {
		return err
	}

	// Handle function types
	if err := ti.unifyFunctionTypes(t1, t2); !errors.Is(err, ErrNotFunctionType) {
		return err
	}

	return fmt.Errorf("%w: cannot unify %s and %s", ErrTypeMismatch, t1.String(), t2.String())
}

// ResolveType resolves a type by following substitutions and resolving unbound type variables
func (ti *TypeInferer) ResolveType(t Type) Type {
	resolved := ti.prune(t)

	// If it's still a type variable after pruning, it's unbound
	if tv, ok := resolved.(*TypeVar); ok {
		// Try to infer a concrete type based on context
		concreteType := ti.resolveUnboundTypeVariable(tv)
		// Update the substitution for future references
		ti.subst[tv.id] = concreteType
		return concreteType
	}

	// For concrete types, return as-is (don't modify)
	if _, ok := resolved.(*ConcreteType); ok {
		return resolved
	}

	// For function types, only resolve if they contain unbound type variables
	if ft, ok := resolved.(*FunctionType); ok {
		hasUnboundVars := false

		// Check if any parameter types need resolution
		newParams := make([]Type, len(ft.paramTypes))
		for i, p := range ft.paramTypes {
			resolvedParam := ti.ResolveType(p)
			newParams[i] = resolvedParam
			if resolvedParam != p {
				hasUnboundVars = true
			}
		}

		// Check if return type needs resolution
		resolvedReturn := ti.ResolveType(ft.returnType)
		if resolvedReturn != ft.returnType {
			hasUnboundVars = true
		}

		// Only create a new function type if something actually changed
		if hasUnboundVars {
			return &FunctionType{
				paramTypes: newParams,
				returnType: resolvedReturn,
			}
		}
	}

	return resolved
}

// ResolveAllEnvironmentTypes resolves all types in the type environment
func (ti *TypeInferer) ResolveAllEnvironmentTypes() {
	for name, t := range ti.env.vars {
		// Skip function types to avoid recursive resolution issues
		if _, ok := t.(*FunctionType); ok {
			continue
		}
		ti.env.vars[name] = ti.ResolveType(t)
	}
}

// InferType performs type inference on an expression
func (ti *TypeInferer) InferType(expr ast.Expression) (Type, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral, *ast.StringLiteral, *ast.BooleanLiteral:
		return ti.inferLiteralType(e)
	case *ast.Identifier:
		return ti.inferIdentifierType(e)
	case *ast.LambdaExpression:
		return ti.inferLambdaExpression(e)
	case *ast.CallExpression:
		return ti.inferCallExpression(e)
	case *ast.InterpolatedStringLiteral:
		return &ConcreteType{name: TypeString}, nil
	case *ast.SpawnExpression, *ast.AwaitExpression, *ast.YieldExpression:
		return ti.inferConcurrencyExpression(e)
	case *ast.BinaryExpression:
		return ti.inferBinaryExpression(e)
	case *ast.FieldAccessExpression:
		return ti.inferFieldAccess(e)
	case *ast.TypeConstructorExpression:
		return ti.inferTypeConstructor(e)
	case *ast.MatchExpression:
		return ti.inferMatchExpression(e)
	case *ast.ResultExpression:
		return ti.inferResultExpression(e)
	case *ast.ListLiteral:
		return ti.inferListLiteral(e)
	case *ast.BlockExpression:
		return ti.inferBlockExpression(e)
	case *ast.UnaryExpression:
		return ti.inferUnaryExpression(e)
	case *ast.MethodCallExpression:
		return ti.inferMethodCall(e)
	case *ast.ListAccessExpression:
		return ti.inferListAccess(e)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedExpression, expr)
	}
}

// InferPattern performs type inference on a pattern
func (ti *TypeInferer) InferPattern(pattern ast.Pattern) (Type, error) {
	switch pattern.Constructor {
	case "_":
		// Wildcard pattern matches anything
		return ti.Fresh(), nil
	case "":
		// Variable pattern
		if pattern.Variable != "" {
			tv := ti.Fresh()
			ti.env.Set(pattern.Variable, tv)
			return tv, nil
		}
		return nil, ErrInvalidEmptyPattern
	case "true", "false":
		// Boolean literal pattern
		return &ConcreteType{name: TypeBool}, nil
	default:
		// Check if it's an integer literal pattern
		if _, err := strconv.ParseInt(pattern.Constructor, 10, 64); err == nil {
			return &ConcreteType{name: TypeInt}, nil
		}

		// Check if it's a string literal pattern (quoted)
		constructorLen := len(pattern.Constructor)
		if constructorLen >= 2 && pattern.Constructor[0] == '"' &&
			pattern.Constructor[constructorLen-1] == '"' {
			return &ConcreteType{name: TypeString}, nil
		}

		// Constructor pattern - look up in environment
		if t, ok := ti.env.Get(pattern.Constructor); ok {
			// Handle field bindings
			for _, field := range pattern.Fields {
				ti.env.Set(field, ti.Fresh())
			}
			return t, nil
		}
		return nil, fmt.Errorf("%w: %s", ErrUnknownConstructor, pattern.Constructor)
	}
}

// unifyTypeVariables handles unification when one or both types are type variables
func (ti *TypeInferer) unifyTypeVariables(t1, t2 Type) error {
	if t1v, ok := t1.(*TypeVar); ok {
		if !ti.occursCheck(t1v, t2) {
			ti.subst[t1v.id] = t2
			return nil
		}
		return ErrRecursiveType
	}

	if t2v, ok := t2.(*TypeVar); ok {
		if !ti.occursCheck(t2v, t1) {
			ti.subst[t2v.id] = t1
			return nil
		}
		return ErrRecursiveType
	}

	return ErrNotTypeVariable
}

// unifyConcreteTypes handles unification of concrete types
func (ti *TypeInferer) unifyConcreteTypes(t1, t2 Type) error {
	ct1, ok1 := t1.(*ConcreteType)
	if !ok1 {
		return ErrNotConcreteType
	}

	if ct2, ok2 := t2.(*ConcreteType); ok2 {
		if ct1.name == ct2.name || ct1.name == TypeAny || ct2.name == TypeAny {
			return nil
		}
		return fmt.Errorf("%w: %s != %s", ErrTypeMismatch, ct1.name, ct2.name)
	}

	// Handle any type compatibility with non-concrete types
	if ct1.name == TypeAny {
		return nil
	}
	if ct2, ok := t2.(*ConcreteType); ok && ct2.name == TypeAny {
		return nil
	}

	return fmt.Errorf("%w: %s != %s", ErrTypeMismatch, ct1.name, t2.String())
}

// unifyFunctionTypes handles unification of function types
func (ti *TypeInferer) unifyFunctionTypes(t1, t2 Type) error {
	ft1, ok1 := t1.(*FunctionType)
	if !ok1 {
		return ErrNotFunctionType
	}

	ft2, ok2 := t2.(*FunctionType)
	if !ok2 {
		return fmt.Errorf("%w: %s != %s", ErrTypeMismatch, ft1.String(), t2.String())
	}

	if len(ft1.paramTypes) != len(ft2.paramTypes) {
		return fmt.Errorf("%w: different arities", ErrTypeMismatch)
	}

	for i, p1 := range ft1.paramTypes {
		if err := ti.Unify(p1, ft2.paramTypes[i]); err != nil {
			return err
		}
	}

	return ti.Unify(ft1.returnType, ft2.returnType)
}

// Helper functions for operator checking
func isArithmeticOp(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" || op == "%"
}

func isComparisonOp(op string) bool {
	return op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">="
}

func isLogicalOp(op string) bool {
	return op == "&&" || op == "||"
}



func uniqueInts(ints []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, i := range ints {
		if !seen[i] {
			seen[i] = true
			result = append(result, i)
		}
	}
	return result
}

func (ti *TypeInferer) occursCheck(v *TypeVar, t Type) bool {
	t = ti.prune(t)
	switch t := t.(type) {
	case *TypeVar:
		return v.id == t.id
	case *FunctionType:
		for _, p := range t.paramTypes {
			if ti.occursCheck(v, p) {
				return true
			}
		}
		return ti.occursCheck(v, t.returnType)
	default:
		return false
	}
}

// resolveUnboundTypeVariable resolves an unbound type variable to a concrete type
func (ti *TypeInferer) resolveUnboundTypeVariable(_ *TypeVar) Type {
	// For now, use a simple heuristic:
	// Default to Int for most unbound variables since many operations produce integers
	// This could be enhanced with more sophisticated inference based on usage context
	return &ConcreteType{name: TypeInt}
}

// prune follows substitution chains to find the actual type
func (ti *TypeInferer) prune(t Type) Type {
	if tv, ok := t.(*TypeVar); ok {
		if subst, exists := ti.subst[tv.id]; exists {
			// Follow the substitution chain and update it
			pruned := ti.prune(subst)
			ti.subst[tv.id] = pruned
			return pruned
		}
	}
	return t
}

// applySubst applies a substitution to a type
func (ti *TypeInferer) applySubst(t Type, subst Substitution) Type {
	switch t := t.(type) {
	case *TypeVar:
		if newType, exists := subst[t.id]; exists {
			return ti.applySubst(newType, subst)
		}
		return t
	case *ConcreteType:
		return t
	case *FunctionType:
		newParams := make([]Type, len(t.paramTypes))
		for i, p := range t.paramTypes {
			newParams[i] = ti.applySubst(p, subst)
		}
		return &FunctionType{
			paramTypes: newParams,
			returnType: ti.applySubst(t.returnType, subst),
		}
	default:
		return t
	}
}

// getFreeVars returns the free type variables in a type
func (ti *TypeInferer) getFreeVars(t Type) []int {
	t = ti.prune(t)
	switch t := t.(type) {
	case *TypeVar:
		return []int{t.id}
	case *ConcreteType:
		return []int{}
	case *FunctionType:
		var vars []int
		for _, p := range t.paramTypes {
			vars = append(vars, ti.getFreeVars(p)...)
		}
		vars = append(vars, ti.getFreeVars(t.returnType)...)
		return uniqueInts(vars)
	default:
		return []int{}
	}
}

// inferLiteralType infers types for literal expressions
func (ti *TypeInferer) inferLiteralType(expr ast.Expression) (Type, error) {
	switch expr.(type) {
	case *ast.IntegerLiteral:
		return &ConcreteType{name: TypeInt}, nil
	case *ast.StringLiteral:
		return &ConcreteType{name: TypeString}, nil
	case *ast.BooleanLiteral:
		return &ConcreteType{name: TypeBool}, nil
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedExpression, expr)
	}
}

// inferIdentifierType infers types for identifier expressions
func (ti *TypeInferer) inferIdentifierType(e *ast.Identifier) (Type, error) {
	if t, ok := ti.env.Get(e.Name); ok {
		// If it's a type scheme, instantiate it
		if scheme, ok := t.(*TypeScheme); ok {
			return ti.Instantiate(scheme), nil
		}
		return t, nil
	}
	if e.Position != nil {
		return nil, fmt.Errorf("line %d:%d: %w: %s", e.Position.Line, e.Position.Column, ErrUndefinedVariable, e.Name)
	}
	return nil, fmt.Errorf("%w: %s", ErrUndefinedVariable, e.Name)
}

// inferLambdaExpression infers types for lambda expressions
func (ti *TypeInferer) inferLambdaExpression(e *ast.LambdaExpression) (Type, error) {
	// Create fresh type variables for parameters
	paramTypes := make([]Type, len(e.Parameters))
	oldEnv := ti.env.Clone()

	for i, param := range e.Parameters {
		paramType := ti.Fresh()
		paramTypes[i] = paramType
		ti.env.Set(param.Name, paramType)
	}

	// Infer body type
	bodyType, err := ti.InferType(e.Body)
	if err != nil {
		ti.env = oldEnv // Restore environment on error
		return nil, err
	}

	// Restore environment
	ti.env = oldEnv

	return &FunctionType{
		paramTypes: paramTypes,
		returnType: bodyType,
	}, nil
}

// inferCallExpression infers types for call expressions
func (ti *TypeInferer) inferCallExpression(e *ast.CallExpression) (Type, error) {
	// Infer function type
	funcType, err := ti.InferType(e.Function)
	if err != nil {
		return nil, err
	}

	// Handle both regular arguments and named arguments
	var argTypes []Type

	if len(e.NamedArguments) > 0 {
		// Process named arguments
		for _, namedArg := range e.NamedArguments {
			argType, err := ti.InferType(namedArg.Value)
			if err != nil {
				return nil, err
			}
			argTypes = append(argTypes, argType)
		}
	} else {
		// Process regular arguments
		for _, arg := range e.Arguments {
			argType, err := ti.InferType(arg)
			if err != nil {
				return nil, err
			}
			argTypes = append(argTypes, argType)
		}
	}

	// Create expected function type
	resultType := ti.Fresh()
	expectedFuncType := &FunctionType{
		paramTypes: argTypes,
		returnType: resultType,
	}

	// Unify with actual function type
	if err := ti.Unify(funcType, expectedFuncType); err != nil {
		return nil, fmt.Errorf("function call type mismatch: %w", err)
	}

	// CRITICAL FIX: Return the resolved/substituted result type, not the fresh variable
	return ti.prune(resultType), nil
}

// inferConcurrencyExpression infers types for concurrency expressions
func (ti *TypeInferer) inferConcurrencyExpression(expr ast.Expression) (Type, error) {
	switch e := expr.(type) {
	case *ast.SpawnExpression:
		// spawn returns a fiber handle
		return &ConcreteType{name: "Fiber"}, nil
	case *ast.YieldExpression:
		// yield returns Unit
		return &ConcreteType{name: TypeUnit}, nil
	case *ast.AwaitExpression:
		// await expression type depends on the awaited fiber
		return ti.InferType(e.Expression)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedExpression, expr)
	}
}

// inferBinaryExpression infers types for binary expressions
func (ti *TypeInferer) inferBinaryExpression(e *ast.BinaryExpression) (Type, error) {
	leftType, err := ti.InferType(e.Left)
	if err != nil {
		return nil, err
	}

	rightType, err := ti.InferType(e.Right)
	if err != nil {
		return nil, err
	}

	switch {
	case isArithmeticOp(e.Operator):
		// Arithmetic operations require Int operands and return Int
		intType := &ConcreteType{name: TypeInt}

		// Both operands must be Int
		if err := ti.Unify(leftType, intType); err != nil {
			return nil, fmt.Errorf("left operand of %s must be Int: %w", e.Operator, err)
		}

		if err := ti.Unify(rightType, intType); err != nil {
			return nil, fmt.Errorf("right operand of %s must be Int: %w", e.Operator, err)
		}

		return intType, nil

	case isComparisonOp(e.Operator):
		// Comparison operations require operands of same type and return Bool
		if err := ti.Unify(leftType, rightType); err != nil {
			return nil, fmt.Errorf("comparison operands must have same type: %w", err)
		}
		return &ConcreteType{name: TypeBool}, nil

	case isLogicalOp(e.Operator):
		// Logical operations require Bool operands and return Bool
		boolType := &ConcreteType{name: TypeBool}
		if err := ti.Unify(leftType, boolType); err != nil {
			return nil, fmt.Errorf("left operand of %s must be Bool: %w", e.Operator, err)
		}
		if err := ti.Unify(rightType, boolType); err != nil {
			return nil, fmt.Errorf("right operand of %s must be Bool: %w", e.Operator, err)
		}
		return boolType, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedBinaryOp, e.Operator)
	}
}

// inferFieldAccess infers types for field access expressions
func (ti *TypeInferer) inferFieldAccess(e *ast.FieldAccessExpression) (Type, error) {
	_, err := ti.InferType(e.Object)
	if err != nil {
		return nil, err
	}

	// For now, create a fresh type variable for field type
	// In a full implementation, this would look up the field type
	fieldType := ti.Fresh()

	return fieldType, nil
}

// inferTypeConstructor infers types for type constructor expressions
func (ti *TypeInferer) inferTypeConstructor(e *ast.TypeConstructorExpression) (Type, error) {
	// Look up constructor in environment
	if t, ok := ti.env.Get(e.TypeName); ok {
		return t, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownConstructor, e.TypeName)
}

// inferMatchExpression infers types for match expressions
func (ti *TypeInferer) inferMatchExpression(e *ast.MatchExpression) (Type, error) {
	// Infer discriminant type
	_, err := ti.InferType(e.Expression)
	if err != nil {
		return nil, err
	}

	if len(e.Arms) == 0 {
		return nil, ErrMatchNoArms
	}

	// Infer type of first arm
	firstArmType, err := ti.InferType(e.Arms[0].Expression)
	if err != nil {
		return nil, err
	}

	// All arms must have the same type
	for i := 1; i < len(e.Arms); i++ {
		armType, err := ti.InferType(e.Arms[i].Expression)
		if err != nil {
			return nil, err
		}
		if err := ti.Unify(firstArmType, armType); err != nil {
			return nil, fmt.Errorf("match arm %d type mismatch: %w", i, err)
		}
	}

	return firstArmType, nil
}

// inferResultExpression infers types for result expressions
func (ti *TypeInferer) inferResultExpression(e *ast.ResultExpression) (Type, error) {
	valueType, err := ti.InferType(e.Value)
	if err != nil {
		return nil, err
	}

	// CRITICAL FIX: If Success=true and ErrorType="", this is a transparent wrapper
	// that should just return the underlying type without wrapping in Result
	if e.Success && e.ErrorType == "" {
		return valueType, nil
	}

	// Otherwise, this is an explicit Result type construction
	if e.Success {
		return &ConcreteType{name: fmt.Sprintf("Result<%s, Error>", valueType.String())}, nil
	}
	return &ConcreteType{name: fmt.Sprintf("Result<T, %s>", valueType.String())}, nil
}

// inferListLiteral infers types for list literal expressions
func (ti *TypeInferer) inferListLiteral(e *ast.ListLiteral) (Type, error) {
	if len(e.Elements) == 0 {
		// Empty list - create fresh element type variable
		elementType := ti.Fresh()
		return &ConcreteType{name: fmt.Sprintf("List<%s>", elementType.String())}, nil
	}

	// Infer type of first element
	firstType, err := ti.InferType(e.Elements[0])
	if err != nil {
		return nil, err
	}

	// All elements must have the same type
	for i := 1; i < len(e.Elements); i++ {
		elemType, err := ti.InferType(e.Elements[i])
		if err != nil {
			return nil, err
		}
		if err := ti.Unify(firstType, elemType); err != nil {
			return nil, fmt.Errorf("list element %d type mismatch: %w", i, err)
		}
	}

	return &ConcreteType{name: fmt.Sprintf("List<%s>", firstType.String())}, nil
}

// inferBlockExpression infers types for block expressions
func (ti *TypeInferer) inferBlockExpression(e *ast.BlockExpression) (Type, error) {
	// Save the current environment so we can restore it after the block
	oldEnv := ti.env.Clone()

	// Process statements to handle let declarations and build up environment
	for _, stmt := range e.Statements {
		if letDecl, ok := stmt.(*ast.LetDeclaration); ok {
			// Infer the type of the let declaration's value
			valueType, err := ti.InferType(letDecl.Value)
			if err != nil {
				ti.env = oldEnv // Restore environment on error
				return nil, err
			}
			// Add the variable to the environment
			ti.env.Set(letDecl.Name, valueType)
		}
		// Other statement types can be handled here as needed
	}

	var blockType Type

	// The type of the block is the type of its final expression
	if e.Expression != nil {
		var err error
		blockType, err = ti.InferType(e.Expression)
		if err != nil {
			ti.env = oldEnv // Restore environment on error
			return nil, err
		}
	} else {
		// CRITICAL FIX: Check if the last statement is an expression statement
		// If so, the block's type is the type of that expression
		if len(e.Statements) > 0 {
			if exprStmt, ok := e.Statements[len(e.Statements)-1].(*ast.ExpressionStatement); ok {
				var err error
				blockType, err = ti.InferType(exprStmt.Expression)
				if err != nil {
					ti.env = oldEnv // Restore environment on error
					return nil, err
				}
			} else {
				blockType = &ConcreteType{name: TypeUnit}
			}
		} else {
			blockType = &ConcreteType{name: TypeUnit}
		}
	}

	// Restore the environment
	ti.env = oldEnv

	return blockType, nil
}

// inferUnaryExpression infers types for unary expressions
func (ti *TypeInferer) inferUnaryExpression(e *ast.UnaryExpression) (Type, error) {
	operandType, err := ti.InferType(e.Operand)
	if err != nil {
		return nil, err
	}

	switch e.Operator {
	case "+", "-":
		// Unary plus and minus require Int operand and return Int
		intType := &ConcreteType{name: TypeInt}
		if err := ti.Unify(operandType, intType); err != nil {
			return nil, fmt.Errorf("operand of %s must be Int: %w", e.Operator, err)
		}
		return intType, nil
	case "!":
		// Logical NOT requires Bool operand and returns Bool
		boolType := &ConcreteType{name: TypeBool}
		if err := ti.Unify(operandType, boolType); err != nil {
			return nil, fmt.Errorf("operand of ! must be Bool: %w", err)
		}
		return boolType, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedUnaryOperator, e.Operator)
	}
}

// inferMethodCall infers types for method call expressions
func (ti *TypeInferer) inferMethodCall(e *ast.MethodCallExpression) (Type, error) {
	// Infer object type
	_, err := ti.InferType(e.Object)
	if err != nil {
		return nil, err
	}

	// For now, create a fresh type variable for method result
	// In a full implementation, this would look up the method type
	resultType := ti.Fresh()

	return resultType, nil
}

// inferListAccess infers types for list access expressions
func (ti *TypeInferer) inferListAccess(e *ast.ListAccessExpression) (Type, error) {
	_, err := ti.InferType(e.List)
	if err != nil {
		return nil, err
	}

	indexType, err := ti.InferType(e.Index)
	if err != nil {
		return nil, err
	}

	// Index must be Int
	intType := &ConcreteType{name: TypeInt}
	if err := ti.Unify(indexType, intType); err != nil {
		return nil, fmt.Errorf("list index must be Int: %w", err)
	}

	// For now, create a fresh type variable for element type
	// In a full implementation, this would extract the element type from List<T>
	elementType := ti.Fresh()

	// List access returns Result<T, Error> for safety
	return &ConcreteType{name: fmt.Sprintf("Result<%s, Error>", elementType.String())}, nil
}

// initializeBuiltInFunctions adds built-in functions to the type environment
func (ti *TypeInferer) initializeBuiltInFunctions() {
	// Built-in function types
	anyType := &ConcreteType{name: TypeAny}
	intType := &ConcreteType{name: TypeInt}
	stringType := &ConcreteType{name: TypeString}

	// Type constructors
	ti.env.Set("Success", &ConcreteType{name: "Success"})
	ti.env.Set("Error", &ConcreteType{name: "Error"})
	ti.env.Set("HttpResponse", &ConcreteType{name: "HttpResponse"})
	ti.env.Set("HttpClient", &ConcreteType{name: "HttpClient"})
	ti.env.Set("HttpServer", &ConcreteType{name: "HttpServer"})
	ti.env.Set("ProcessHandle", &ConcreteType{name: "ProcessHandle"})
	ti.env.Set("Iterator", &ConcreteType{name: "Iterator"})

	// print(value: any) -> Int
	ti.env.Set("print", &FunctionType{
		paramTypes: []Type{anyType},
		returnType: intType,
	})

	// toString(value: any) -> String
	ti.env.Set("toString", &FunctionType{
		paramTypes: []Type{anyType},
		returnType: stringType,
	})

	// input() -> Result<Int, Error>
	ti.env.Set("input", &FunctionType{
		paramTypes: []Type{},
		returnType: &ConcreteType{name: "Result<Int, Error>"},
	})

	// length(s: String) -> Int (can't fail, returns plain Int)
	ti.env.Set("length", &FunctionType{
		paramTypes: []Type{stringType},
		returnType: intType,
	})

	// contains(haystack: String, needle: String) -> Result<Bool, Error>
	ti.env.Set("contains", &FunctionType{
		paramTypes: []Type{stringType, stringType},
		returnType: &ConcreteType{name: "Result<Bool, Error>"},
	})

	// substring(s: String, start: Int, length: Int) -> Result<String, Error>
	ti.env.Set("substring", &FunctionType{
		paramTypes: []Type{stringType, intType, intType},
		returnType: &ConcreteType{name: "Result<String, Error>"},
	})

	// File I/O functions
	// readFile(filename: String) -> Result<String, Error>
	ti.env.Set("readFile", &FunctionType{
		paramTypes: []Type{stringType},
		returnType: &ConcreteType{name: "Result<String, Error>"},
	})

	// writeFile(filename: String, content: String) -> Result<Unit, Error>
	ti.env.Set("writeFile", &FunctionType{
		paramTypes: []Type{stringType, stringType},
		returnType: &ConcreteType{name: "Result<Unit, Error>"},
	})

	// Process functions
	// spawnProcess(command: String, args: String) -> Result<ProcessHandle, Error>
	ti.env.Set("spawnProcess", &FunctionType{
		paramTypes: []Type{stringType, stringType},
		returnType: &ConcreteType{name: "Result<ProcessHandle, Error>"},
	})

	// awaitProcess(handle: ProcessHandle) -> Result<Int, Error>
	ti.env.Set("awaitProcess", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "ProcessHandle"}},
		returnType: &ConcreteType{name: "Result<Int, Error>"},
	})

	// cleanupProcess(handle: ProcessHandle) -> Result<Unit, Error>
	ti.env.Set("cleanupProcess", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "ProcessHandle"}},
		returnType: &ConcreteType{name: "Result<Unit, Error>"},
	})

	// HTTP functions
	// httpCreateClient(baseUrl: String, timeout: Int) -> Result<HttpClient, Error>
	ti.env.Set("httpCreateClient", &FunctionType{
		paramTypes: []Type{stringType, intType},
		returnType: &ConcreteType{name: "Result<HttpClient, Error>"},
	})

	// httpCreateServer(port: Int, address: String) -> Result<HttpServer, Error>
	ti.env.Set("httpCreateServer", &FunctionType{
		paramTypes: []Type{intType, stringType},
		returnType: &ConcreteType{name: "Result<HttpServer, Error>"},
	})

	// httpGet(client: HttpClient, path: String, headers: String) -> Result<HttpResponse, Error>
	ti.env.Set("httpGet", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "HttpClient"}, stringType, stringType},
		returnType: &ConcreteType{name: "Result<HttpResponse, Error>"},
	})

	// httpPost(client: HttpClient, path: String, headers: String, body: String) -> Result<HttpResponse, Error>
	ti.env.Set("httpPost", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "HttpClient"}, stringType, stringType, stringType},
		returnType: &ConcreteType{name: "Result<HttpResponse, Error>"},
	})

	// httpListen(server: Result<HttpServer, Error>, handler: Function) -> Result<Int, Error>
	ti.env.Set("httpListen", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "Result<HttpServer, Error>"}, anyType},
		returnType: &ConcreteType{name: "Result<Int, Error>"},
	})

	// Iterator functions
	// range(start: Int, end: Int) -> Iterator<Int>
	ti.env.Set("range", &FunctionType{
		paramTypes: []Type{intType, intType},
		returnType: &ConcreteType{name: "Iterator<Int>"},
	})

	// forEach(iter: Iterator<T>, fn: T -> Unit) -> Unit
	ti.env.Set("forEach", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "Iterator<T>"}, &ConcreteType{name: "T -> Unit"}},
		returnType: &ConcreteType{name: TypeUnit},
	})

	// map(iter: Iterator<T>, fn: T -> U) -> Iterator<U>
	ti.env.Set("map", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "Iterator<T>"}, &ConcreteType{name: "T -> U"}},
		returnType: &ConcreteType{name: "Iterator<U>"},
	})

	// filter(iter: Iterator<T>, predicate: T -> Bool) -> Iterator<T>
	ti.env.Set("filter", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "Iterator<T>"}, &ConcreteType{name: "T -> Bool"}},
		returnType: &ConcreteType{name: "Iterator<T>"},
	})

	// fold(iter: Iterator<T>, initial: U, fn: (U, T) -> U) -> U
	ti.env.Set("fold", &FunctionType{
		paramTypes: []Type{&ConcreteType{name: "Iterator<T>"}, &ConcreteType{name: "U"}, &ConcreteType{name: "(U, T) -> U"}},
		returnType: &ConcreteType{name: "U"},
	})

	// Add more built-in functions as needed
}
