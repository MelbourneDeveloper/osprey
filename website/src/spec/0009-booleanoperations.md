---
layout: page
title: "Boolean Operations"
description: "Osprey Language Specification: Boolean Operations"
date: 2026-05-17
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0009-booleanoperations/"
---

# Boolean Operations

Osprey has no `if`/`else` statement. Conditional logic is written as a `match` on a boolean (which forces both arms to be considered) or as the ternary shorthand `cond ? then : else`, which desugars to the same `match`. The ternary is defined in [Pattern Matching](0007-PatternMatching.md#ternary-match-syntactic-sugar).

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

## Boolean Operators

`&&`, `||`, and `!` are short-circuiting; `==`, `!=`, `<`, `>`, `<=`, `>=` produce booleans. See [Lexical Structure](0002-LexicalStructure.md) for the full operator list.

```osprey
let isAdult       = age >= 18
let hasPermission = isAdult && isAuthorized
let canAccess     = hasPermission || isAdmin
let isBlocked     = !isActive
let validUser     = !isBanned && (isVerified || hasInvite)
```