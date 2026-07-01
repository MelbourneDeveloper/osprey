---
layout: page
title: "Pattern Matching"
description: "Osprey Language Specification: Pattern Matching"
date: 2026-07-01
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0007-patternmatching/"
---

# Pattern Matching

`match` is the only branching construct in Osprey. Record patterns are matched structurally by field name, not by field order. See [Type System](/spec/0004-typesystem/) for type unification rules.

> **Flavor layer — mixed.**  A `match` lowers to `Expr::Match` over `MatchArm`s, each carrying a `Pattern` (`Wildcard`, `Literal`, `Constructor { name, fields, sub_patterns }`, `TypeAnnotated`, `Structural`, `List`, `Binding`). Only the *spelling* of these patterns is a surface (CST) concern. This chapter documents both flavors: samples below appear in both surfaces — Default (`.osp`) then its ML (`.ospml`) twin, each tagged with a flavor badge. A one-field variant is `Success { value }` in Default, where the ML flavor writes `Success value` ([`[FLAVOR-ML-MATCH]`](/spec/0024-mlflavorsyntax/#match)) — but both flavors lower to the **same** `Pattern::Constructor { name, fields }`. Everything else here — exhaustiveness checking, `any`/union narrowing, and arm semantics — is shared-core: it runs on the canonical AST and is flavor-blind ([`[FLAVOR-BOUNDARY]`](/spec/0023-languageflavors/#the-one-law)). See [Language Flavors](/spec/0023-languageflavors/) and [ML Flavor Syntax](/spec/0024-mlflavorsyntax/).

## Basic Patterns

```osprey
let result = match value {
    0 => "zero"
    1 => "one"
    n => "other: " + toString(n)
}
```

The ML flavor drops the braces — the scrutinee sits on the `match` line and each `pattern => body` arm is an indented (offside) region:

```osprey-ml
result =
    match value
        0 => "zero"
        1 => "one"
        n => "other: " + toString n
```

## Union Type Patterns

A union pattern names the variant. Variants with fields are destructured using `{ field, ... }`; variants without fields are matched by name alone. Both forms lower to `Pattern::Constructor`; the brace destructuring shown here is the Default surface, spelled `Success value` in the ML flavor ([`[FLAVOR-ML-MATCH]`](/spec/0024-mlflavorsyntax/#match)).

```osprey
type Option = Some { value: int } | None

let message = match option {
    Some { value } => "Value: " + toString(value)
    None           => "No value"
}
```

In the ML flavor the one-field variant pattern is written `Some value` — the payload binds positionally, with no braces — and the arms are laid out by indentation:

```osprey-ml
type Option = Some { value: int } | None

message =
    match option
        Some value => "Value: " + toString value
        None       => "No value"
```

## Wildcard Patterns

The underscore `_` matches any value:

```osprey
let category = match score {
    100 => "perfect"
    90 => "excellent"
    _ => "good"
}
```

## Type Annotation Patterns

A pattern of the form `name: type` matches when the value has the named type and binds it. This is the required form for narrowing an `any` value. The grammar for all pattern forms is in [Syntax](/spec/0003-syntax/#match-expressions).

```osprey
// Narrowing an any value
match anyValue {
    n: int    => n + 1
    s: string => length(s)
    b: bool   => match b {
        true  => 1
        false => 0
    }
    _ => 0
}

// Structural matching: any type with these field names
match anyValue {
    { name, age }       => print("${name}: ${age}")
    p: { name, age }    => print("person ${p.name}: ${p.age}")   // bind whole + destructure
    u: User { id }      => print("user ${id}")                   // typed structural
    _                   => print("unknown")
}

// Type-narrowed structural fields
match anyValue {
    { x, y }                       => print("point: (${x}, ${y})")
    p: { name }                    => print("named: ${p.name}")
    { id, email, active: bool }    => print("active user: ${id}")
    _                              => print("no match")
}

// Type pattern with destructuring of a known constructor
match result {
    success: Success { value, timestamp } => processSuccess(value: value, timestamp: timestamp)
    error:   Error   { code, message }    => handleError(code: code, message: message)
    _                                     => defaultHandler()
}
```

## Result Patterns

`Result<T, E>` is matched the same way as any other union. See [Error Handling](/spec/0013-errorhandling/) for the type and arithmetic semantics.

```osprey
let calculation = 1 + 3 + (300 / 5)  // Result<int, MathError>

match calculation {
    Success { value }   => print("Result: ${value}")
    Error   { message } => print("Math error: ${message}")
}
```

The same `Result` match in the ML flavor binds each payload positionally (`Success value`, `Error message`) with offside arms:

```osprey-ml
calculation = 1 + 3 + (300 / 5)  // Result<int, MathError>

match calculation
    Success value   => print "Result: ${value}"
    Error message   => print "Math error: ${message}"
```

Compound arithmetic expressions yield a single `Result`, not nested `Result`s; the compiler unwraps intermediate values inside the chain. Only the final value needs to be matched.

## Ternary Match (Syntactic Sugar)

A two-arm match has a shorthand. Two equivalent forms exist:

```ebnf
ternary ::= expr "{" pattern "}" "?" expr ":" expr   (* structural form *)
          | expr "?:" expr                            (* Result default form *)
```

Structural form — pick out a field, fall back if the pattern fails:

```osprey
let calculation = 10 + 5
let value = calculation { value } ? value : -1   // 15
```

Desugars to:

```osprey
match calculation {
    { value } => value
    _         => -1
}
```

Result-default form — extract `Success { value }` or use the default on `Error`:

```osprey
let safeValue = divide(a: 10, b: 2) ?: -1   // 5
let errorVal  = divide(a: 10, b: 0) ?: -1   // -1
```

A boolean expression with `?:` works because `true`/`false` desugar to the same match:

```osprey
let status = isActive ? "Active" : "Inactive"
```