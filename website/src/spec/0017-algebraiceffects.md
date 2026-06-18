---
layout: page
title: "Algebraic Effects"
description: "Osprey Language Specification: Algebraic Effects"
date: 2026-06-18
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0017-algebraiceffects/"
---

# Algebraic Effects

Osprey treats effects as first-class language features. An effect declares a set of operations; functions list the effects they may perform; handlers give meaning to operations. The compiler rejects any program that performs an unhandled effect.

## Status

Effect declarations, `perform` expressions, effect annotations on function types, handler parsing, and full compile-time unhandled-effect checking are implemented. Continuation/`resume` semantics inside handlers are not yet implemented; current handlers act as value substitutions, which is sufficient for many uses but does not yet model the full algebraic-effects calculus.

## Keywords

```
effect perform handle in
```

## Effect Declarations

```ebnf
effectDecl ::= docComment? "effect" IDENT "{" opDecl* "}"
opDecl     ::= IDENT ":" fnType
```

```osprey
effect State {
    get : fn() -> int
    set : fn(int) -> unit
}
```

## Effectful Function Types

A function declares the effects it may perform with `!E` after its return type. `E` is either a single effect or a bracketed set.

```osprey
fn read() -> string !IO = perform IO.readLine()
fn fetch(url: string) -> string ![IO, Net] = ...
```

A function with no `!E` is pure; calling an effectful function from a pure context is a compilation error.

## Performing Operations

```ebnf
performExpr ::= "perform" IDENT "." IDENT "(" args? ")"
```

```osprey
fn incrementTwice() -> int !State = {
    let current = perform State.get()
    perform State.set(current + 1)
    perform State.get()
}
```

If no enclosing handler covers an effect, the program does not compile.

## Handlers

```ebnf
handlerExpr ::= "handle" IDENT handlerArm+ "in" expr
handlerArm  ::= IDENT paramList? "=>" expr
```

```osprey
handle State
    get        => 42
    set newVal => print("set to " + toString(newVal))
in
    incrementTwice()
```

The innermost matching handler wins for each effect. Handlers may be nested freely:

```osprey
handle Logger
    log msg => print("[OUTER] " + msg)
in
    handle Logger
        log msg => print("[INNER] " + msg)
    in
        perform Logger.log("test")    // prints "[INNER] test"
```

## Effect Inference

The compiler infers the minimal effect set of every expression. Functions either declare their effects or are required to be pure. A function may be polymorphic over an effect set:

```osprey
fn loggedCalculation<E>(x: int) -> int !E = {
    perform Logger.log("calculating")     // E must include Logger
    x * 2
}
```

## Static Safety Checks

The compiler enforces three static checks on effect programs. Each failure is a compile-time error, not a runtime fault.

| Check                              | Failure mode in other languages |
| ---------------------------------- | ------------------------------- |
| Every `perform` has a handler      | Runtime crash / unhandled exn   |
| No circular effect dependency      | Stack overflow                  |
| No handler that performs the same effect it handles | Infinite loop |

### Circular Dependency Example

```osprey
effect StateA { getFromB: fn() -> int }
effect StateB { getFromA: fn() -> int }

fn circularA() -> int !StateA = perform StateA.getFromB()
fn circularB() -> int !StateB = perform StateB.getFromA()

handle StateA
    getFromB => circularB()       // ❌ circular dependency
in
    handle StateB
        getFromA => circularA()   // ❌ circular dependency
    in
        circularA()
```

### Handler-Self-Recursion Example

```osprey
effect Counter { increment: fn(int) -> int }

fn performIncrement(n: int) -> int !Counter = perform Counter.increment(n)

handle Counter
    increment n => performIncrement(n + 1)   // ❌ handler performs the effect it handles
in
    performIncrement(5)
```

## Worked Example

`x * 2` returns `Result<int, MathError>`; the function below performs `Exception` on overflow and `State` to record the success.

```osprey
effect Exception { raise: fn(string) -> unit }
effect State     { get: fn() -> int, set: fn(int) -> unit }

fn doubleAndStore(x: int) -> int ![Exception, State] = match x * 2 {
    Success { value }   => {
        perform State.set(value)
        value
    }
    Error   { message } => {
        perform Exception.raise(message)
        0
    }
}

handle Exception
    raise msg => { print("error: " + msg); -1 }
in
    handle State
        get        => 0
        set newVal => print("state: " + toString(newVal))
    in
        let result = doubleAndStore(21)
        print("result: " + toString(result))
```