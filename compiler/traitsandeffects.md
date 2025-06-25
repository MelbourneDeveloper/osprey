## 19. Traits and Type Constraints ([ospreylang.dev][1], [ospreylang.dev][1])

OSPREY HAS A FIRST-CLASS EFFECTS SYSTEM

https://www.inner-product.com/posts/direct-style-effects/

https://www.eff-lang.org/handlers-tutorial.pdf

https://en.wikipedia.org/wiki/Effect_system

https://dl.acm.org/doi/pdf/10.1145/3290319

### 19.0 IMPLEMENTATION STATUS

NOT IMPLEMENTED. You can write the syntax today, but the compiler ignores `trait` and `impl` blocks. The team will land parsing and name-resolution support in **0.2.0-alpha**.

### 19.1 New Keywords

Add these to the **reserved keyword** list:

```
trait impl where
```

### 19.2 Syntax

```ebnf
traitDecl  ::= docComment? "trait" IDENT genericParams? "{" traitMember* "}"
traitMember::= fnSig defaultBody?
defaultBody::= "=" expr

implDecl   ::= docComment? "impl" traitRef "for" typeRef "{" implMember* "}"
implMember ::= fnDecl        // full body required unless default provided

genericParams ::= "<" IDENT ("," IDENT)* ">"
traitRef      ::= IDENT genericArgs?
typeRef       ::= IDENT genericArgs?
genericArgs   ::= "<" type ("," type)* ">"
```

### 19.3 Rules

| Rule                 | Description                                                                                           |
| -------------------- | ----------------------------------------------------------------------------------------------------- |
| **Trait definition** | You group related capabilities. A trait owns **no data**.                                             |
| **Implementation**   | You write `impl` blocks to supply code for a specific `type`.                                         |
| **Coherence**        | One of *trait* or *type* must live in your module. The compiler rejects code that violates this rule. |
| **Default methods**  | Provide defaults with `=` expressions. Implementations may override.                                  |
| **Super-traits**     | Declare with `trait Ord : Eq { ... }`. Any `impl Ord for T` must also `impl Eq for T`.                |
| **Trait bounds**     | Use either inline `fn foo<T: Read + Write>(x)` or a trailing `where` clause.                          |
| **Static dispatch**  | Generic calls monomorphise at compile time. No vtables.                                               |
| **Trait objects**    | Not in 0.1.0. We will add `dyn` syntax after the static path is rock-solid.                           |

### 19.4 Examples

```osprey
// 1. Declare a trait
trait Addable {
  add: fn(self, other: Self) -> Self
  zero: fn() -> Self = self.add(self, self)  // default uses 'add'
}

// 2. Implement it for Int
impl Addable for Int {
  fn add(self, other: Int) -> Int = self + other
  fn zero() -> Int = 0
}

// 3. Use a trait bound
fn sum<T: Addable>(left: T, right: T) -> T = left.add(right)
```

---

## 20. Algebraic Effects ([ospreylang.dev][1])

### 20.0 IMPLEMENTATION STATUS

PARSER ONLY. Type-checking and lowering to LLVM drop **0.2.0-alpha**.

### 20.1 New Keywords

```
effect perform handler with
```

### 20.2 Effect Declarations

```ebnf
effectDecl   ::= docComment? "effect" IDENT "{" opDecl* "}"
opDecl       ::= IDENT ":" fnType
```

Example:

```osprey
effect IO {
  readLine : fn() -> String
  writeLine: fn(String) -> Unit
}
```

### 20.3 Effectful Function Types

Append `!EffectSet` to a function type. A set literal uses `[...]`.

```osprey
fn read() -> String !IO = perform IO.readLine()

fn fetch(url: String) -> String ![IO, Net] = ...
```

### 20.4 Performing an Operation

```
perform EffectName.operation(args...)
```

**CRITICAL COMPILE-TIME SAFETY**: If no handler intercepts the call, the compiler produces a **compilation error**. Unhandled effects are **NEVER** permitted at runtime - this ensures complete effect safety and prevents runtime crashes.

### 20.5 Handlers

```ebnf
handlerExpr ::= "handler" IDENT "{" handlerArm+ "}"
handlerArm  ::= IDENT pattern? "=>" expr

withExpr    ::= "with" handlerExpr "do" block
```

Example:

```osprey
with handler IO {
  readLine()      => "mock"
  writeLine(msg)  => log(msg)
} do {
  runApp()
}
```

The handler captures `perform` calls inside `runApp`. You may nest handlers; the innermost wins.

### 20.6 Effects in Traits

A trait method may list effects:

```osprey
trait Parsable {
  parse: fn(self) -> Int !IO
}

impl Parsable for String {
  fn parse(self) -> Int !IO = perform IO.readLine() |> toInt
}
```

When you call `parse`, its effect set bubbles to the caller. The compiler unions all effect sets on the call path and forces you to handle or forward them.

### 20.7 Inference and Checking

* The compiler infers the minimal effect set for every `expr`.
* If inference fails, you add an explicit `!EffectSet`.
* A function that *declares* no effects **must not** perform any.
* **COMPILE-TIME EFFECT SAFETY**: All effects must be explicitly handled or forwarded in function signatures. The compiler **rejects** any code with unhandled effects - **NO RUNTIME EFFECT ERRORS ALLOWED**.

### 20.8 Compilation Model (Design)

1. Front-end tags each `perform` with its effect ID.
2. Middle-end rewrites effectful code into CPS with explicit handler frames.
3. Back-end lowers the CPS form to LLVM using jump tables. Zero overhead when you run without handlers.

---

### 20.9 Examples

```osprey
// Pure function
fn double(x) = x * 2

// Effectful function that writes a line
fn greet(name: String) !IO = perform IO.writeLine("Hello " + name)

// Handle IO in main
with handler IO {
  writeLine(msg) => print(msg)
  readLine()     => ""       // never called
} do {
  greet("Osprey")
}
```

---

## 21. Summary

Traits give you zero-cost static polymorphism. Algebraic effects give you explicit, type-checked side-effects. Both features fit Osprey’s goals: **clarity, safety, and performance**.

[1]: https://www.ospreylang.dev/spec/ "Osprey Language Specification - Osprey Programming Language"
