---
layout: page
title: "Boolean Operations"
description: "Osprey Language Specification: Boolean Operations"
date: 2026-07-01
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0009-booleanoperations/"
---

# Boolean Operations

Osprey has no `if`/`else` statement. Conditional logic is written as a `match` on a boolean (which forces both arms to be considered) or as the ternary shorthand `cond ? then : else`, which desugars to the same `match`. The ternary is defined in [Pattern Matching](/spec/0007-patternmatching/#ternary-match-syntactic-sugar).

> **Flavor layer — shared core (AST and above).**  Boolean semantics are flavor-blind. `&&`, `||`, and the comparison operators lower to `Expr::Binary`, `!` to `Expr::Unary`, and every conditional to `Expr::Match` over the boolean — the canonical AST nodes the type checker, effect checker, and codegen consume without ever knowing which flavor produced them. Only the *spelling* differs: the Default flavor brackets `match` arms with braces, the ML flavor uses offside layout ([FLAVOR-ML-MATCH](/spec/0024-mlflavorsyntax/#match)). Both samples below lower to the same AST. See [Language Flavors](/spec/0023-languageflavors/).

Both flavors are shown side by side throughout this chapter — Default first, then its ML twin.

```osprey
let status = match isValid {
    true  => "Success"
    false => "Failure"
}

let max = match a > b {
    true  => a
    false => b
}
```

```osprey-ml
status =
    match isValid
        true => "Success"
        false => "Failure"

max =
    match a > b
        true => a
        false => b
```

Nested matches handle compound conditions:

```osprey
let category = match score >= 90 {
    true  => match score == 100 {
        true  => "Perfect"
        false => "Excellent"
    }
    false => match score >= 70 {
        true  => "Good"
        false => "Needs Improvement"
    }
}
```

```osprey-ml
category =
    match score >= 90
        true =>
            match score == 100
                true => "Perfect"
                false => "Excellent"
        false =>
            match score >= 70
                true => "Good"
                false => "Needs Improvement"
```

## Boolean Operators

`&&`, `||`, and `!` are short-circuiting; `==`, `!=`, `<`, `>`, `<=`, `>=` produce booleans. See [Lexical Structure](/spec/0002-lexicalstructure/) for the full operator list.

```osprey
let isAdult       = age >= 18
let hasPermission = isAdult && isAuthorized
let canAccess     = hasPermission || isAdmin
let isBlocked     = !isActive
let validUser     = !isBanned && (isVerified || hasInvite)
```

```osprey-ml
isAdult       = age >= 18
hasPermission = isAdult && isAuthorized
canAccess     = hasPermission || isAdmin
isBlocked     = !isActive
validUser     = !isBanned && (isVerified || hasInvite)
```