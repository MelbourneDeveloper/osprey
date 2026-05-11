# Syntax

This chapter defines the syntactic forms that make up an Osprey program. Semantics for individual constructs are in their dedicated chapters; cross-references are noted inline.

- [Program Structure](#program-structure)
- [Imports](#imports)
- [Let Declarations](#let-declarations)
- [Function Declarations](#function-declarations)
- [Extern Declarations](#extern-declarations)
- [Type Declarations](#type-declarations)
- [Records](#records)
- [Expressions](#expressions)
- [Field Access](#field-access)
- [Match Expressions](#match-expressions)
- [Variable Binding](#variable-binding)

## Program Structure

```ebnf
program   ::= statement* EOF
statement ::= importStmt
            | letDecl
            | fnDecl
            | externDecl
            | typeDecl
            | moduleDecl
            | exprStmt
```

## Imports

```ebnf
importStmt ::= "import" ID ("." ID)*
```

```osprey
import std
import std.io
import graphics.canvas
```

## Let Declarations

```ebnf
letDecl ::= ("let" | "mut") ID (":" type)? "=" expr
```

```osprey
let x       = 42
let name    = "Alice"
mut counter = 0
let result  = calculateValue(input: data)
```

`let` binds immutably; `mut` binds mutably. Type annotations are optional.

## Function Declarations

```ebnf
fnDecl    ::= docComment? "fn" ID "(" paramList? ")" ("->" type)? effectSet?
              ("=" expr | "{" blockBody "}")
paramList ::= param ("," param)*
param     ::= ID (":" type)?
```

```osprey
fn double(x)   = x * 2
fn add(x, y)   = x + y
fn greet(name) = "Hello " + name
fn getValue()  = 42
```

Effect sets (`!E`) are described in [Algebraic Effects](0017-AlgebraicEffects.md). Functions of two or more parameters require named arguments at call sites; see [Function Calls](0005-FunctionCalls.md). The plugin-function variant `fn <plugin> <name>(...) = ... ;` is described in [Language Plugins](0018-LanguagePlugins.md).

## Extern Declarations

`extern` declares an interface to a foreign function (Rust, C, or any C-ABI library). It has no body.

```ebnf
externDecl      ::= docComment? "extern" "fn" ID "(" externParamList? ")" ("->" type)?
externParamList ::= externParam ("," externParam)*
externParam     ::= ID ":" type
```

Parameter types are required. Calls use named arguments (single-parameter functions may use positional).

```osprey
extern fn rust_add(a: int, b: int) -> int
extern fn rust_is_prime(n: int) -> bool

let sum     = rust_add(a: 15, b: 25)
let isPrime = rust_is_prime(17)
```

ABI mapping:

| Osprey   | Rust                | C            |
| -------- | ------------------- | ------------ |
| `int`    | `i64`               | `int64_t`    |
| `bool`   | `bool`              | `bool`       |
| `string` | `*const c_char`     | `char *`     |

The foreign function must use the C ABI (`extern "C"` and `#[no_mangle]` in Rust) and be linked at compile time.

## Type Declarations

```ebnf
typeDecl          ::= docComment? "type" ID ("<" typeParamList ">")? "=" (unionType | recordType)
unionType         ::= variant ("|" variant)*
recordType        ::= "{" fieldDeclarations "}"
variant           ::= ID ("{" fieldDeclarations "}")?
fieldDeclarations ::= fieldDeclaration ("," fieldDeclaration)*
fieldDeclaration  ::= ID ":" type constraint?
constraint        ::= "where" function_name
```

```osprey
type Color = Red | Green | Blue

type Shape = Circle    { radius: int }
           | Rectangle { width: int, height: int }
```

## Records

A record type names a fixed set of fields. Construction uses `TypeName { field: value, ... }`; field order at the call site is irrelevant.

```osprey
type Point  = { x: int, y: int }
type Person = { name: string, age: int } where validatePerson

let point  = Point { x: 10, y: 20 }
let person = Person { name: "Alice", age: 25 }
```

Validation, non-destructive update (`record { field: value }`), and full field-access semantics are in [Type System](0004-TypeSystem.md).

## Expressions

```ebnf
expression          ::= logicalOrExpression
logicalOrExpression ::= logicalAndExpression ("||" logicalAndExpression)*
logicalAndExpression::= comparisonExpression ("&&" comparisonExpression)*
comparisonExpression::= additiveExpression (("==" | "!=" | "<" | ">" | "<=" | ">=") additiveExpression)*
additiveExpression  ::= multiplicativeExpression (("+" | "-") multiplicativeExpression)*
multiplicativeExpression ::= unaryExpression (("*" | "/") unaryExpression)*
unaryExpression     ::= ("+" | "-" | "!")? pipeExpression
pipeExpression      ::= callExpression ("|>" callExpression)*
callExpression      ::= primaryExpression (
                          "." ID "(" argumentList? ")"
                        | "(" argumentList? ")"
                        | "[" expression "]"
                        | "." ID
                      )*
primaryExpression   ::= literal | ID | "(" expression ")"
                      | lambdaExpression | blockExpression | matchExpression

argumentList        ::= namedArgument ("," namedArgument)+
                      | expression ("," expression)*
namedArgument       ::= ID ":" expression
```

Precedence, highest to lowest:

1. Unary `!`, `-`
2. Multiplicative `*`, `/`, `%`
3. Additive `+`, `-`
4. Comparison `==`, `!=`, `<`, `>`, `<=`, `>=`
5. Logical AND `&&`
6. Logical OR `||`

Block expressions and their scoping are defined in [Block Expressions](0008-BlockExpressions.md). Pattern-matching for booleans (the only conditional construct) is in [Boolean Operations](0009-BooleanOperations.md).

## List Access

```ebnf
listAccess ::= expression "[" expression "]"
```

Indexing returns `Result<T, IndexError>`:

```osprey
let numbers = [1, 2, 3, 4]

match numbers[0] {
    Success { value }   => print("first: ${value}")
    Error   { message } => print("index error: ${message}")
}
```

## Field Access

```ebnf
fieldAccess ::= expression "." ID
```

Fields are accessible directly only on record values:

```osprey
type User = { id: int, name: string }
let user  = User { id: 1, name: "Alice" }
let n     = user.name
```

Field access on `any`, `Result`, or any union type requires a `match` to narrow the value first. See [Type System](0004-TypeSystem.md) for the full rules.

Records are immutable. Use the non-destructive update form to produce a modified copy:

```osprey
let p2 = point { x: 15 }   // y carried over
```

## Match Expressions

```ebnf
matchExpr   ::= "match" expr "{" matchArm+ "}"
matchArm    ::= pattern "=>" expr
pattern     ::= unaryExpr                              (* literals incl. -1, +42 *)
              | ID ("{" fieldPattern "}")?             (* constructor / destructure *)
              | ID "(" pattern ("," pattern)* ")"      (* positional constructor *)
              | ID ":" type                            (* type annotation *)
              | ID ":" "{" fieldPattern "}"            (* named structural *)
              | "{" fieldPattern "}"                   (* anonymous structural *)
              | "_"                                    (* wildcard *)
fieldPattern::= ID ("," ID)*
```

```osprey
let result = match status {
    Success    => "OK"
    Error msg  => "Failed: " + msg
    _          => "Unknown"
}
```

Pattern semantics, exhaustiveness, and the ternary shorthand are in [Pattern Matching](0007-PatternMatching.md).

## Variable Binding

- `let` creates an immutable binding; `mut` creates a mutable one.
- Every binding is initialised at declaration.
- Inner scopes may shadow outer bindings.
- Function arguments evaluate left to right before the call. `&&` and `||` short-circuit.
