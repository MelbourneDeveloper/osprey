---
layout: page
title: "10. Boolean Operations"
description: "Osprey Language Specification: 10. Boolean Operations"
date: 2025-06-26
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0010-booleanoperations/"
---

## 10.1 Boolean Pattern Matching

Osprey uses pattern matching instead of traditional if-else statements for boolean operations. This ensures exhaustive handling of both true and false cases.

**Basic Boolean Matching:**
```osprey
let status = match isValid {
    true => "Success"
    false => "Failure"
}
```

**Complex Boolean Logic:**
```osprey
let category = match score >= 90 {
    true => match score == 100 {
        true => "Perfect"
        false => "Excellent"
    }
    false => match score >= 70 {
        true => "Good"
        false => "Needs Improvement"
    }
}
```

## 10.2 Boolean Operators

- `&&` - Logical AND
- `||` - Logical OR  
- `!` - Logical NOT
- `==` - Equality
- `!=` - Inequality
- `>`, `<`, `>=`, `<=` - Comparison operators

**Operator Examples:**
```osprey
let isAdult = age >= 18
let hasPermission = isAdult && isAuthorized
let canAccess = hasPermission || isAdmin
let isBlocked = !isActive
```