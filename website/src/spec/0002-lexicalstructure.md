---
layout: page
title: "Lexical Structure"
description: "Osprey Language Specification: Lexical Structure"
date: 2026-07-01
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0002-lexicalstructure/"
---

# Lexical Structure

- [Identifiers](#identifiers)
- [Keywords](#keywords)
- [Literals](#literals)
- [Operators](#operators)
- [Delimiters](#delimiters)

> **Flavor layer — surface (CST).**  This chapter is the Default flavor (`.osp`) lexical grammar: a flavor-internal, below-the-AST concern owned by `crates/osprey-syntax/src/default/`. Tokens are a CST artifact — they never reach the shared core, which sees only the canonical `osprey_ast::Program` after lowering. The ML flavor (`.ospml`) has its OWN offside-rule layout lexer (`crates/osprey-syntax/src/ml/lexer.rs`) that derives `INDENT`/`DEDENT`/`NEWLINE` from an explicit indent stack ([FLAVOR-ML-LAYOUT] in [ML Flavor Syntax](/spec/0024-mlflavorsyntax/)); the rules below describe only this surface. Lexical structure differs per flavor; both feed lowering. See [Language Flavors](/spec/0023-languageflavors/).

## Identifiers

Start with letter or underscore, followed by letters, digits, or underscores.
```
ID := [a-zA-Z_][a-zA-Z0-9_]*
```

## Keywords

```
fn let mut type match extern import module where
effect perform handle in
spawn await yield select
```

## Literals

### Integer Literals
```
INTEGER := [0-9]+
```

**Examples:**
```osprey
let count = 42
let negative = -17
let zero = 0
```

### Float Literals
```
FLOAT := [0-9]+ '.' [0-9]+ ([eE] [+-]? [0-9]+)?
       | [0-9]+ [eE] [+-]? [0-9]+
```

**Examples:**
```osprey
let pi = 3.14159
let temperature = -273.15
let scientific = 6.022e23
let small = 1.5e-10
```

**Type Inference:**
- Integer literals without decimal point infer to `int`
- Literals with decimal point or scientific notation infer to `float`

### String Literals
```
STRING := '"' (CHAR | ESCAPE_SEQUENCE)* '"'
ESCAPE_SEQUENCE := '\n' | '\t' | '\r' | '\\' | '\"'
```

### Interpolated String Literals
```
INTERPOLATED_STRING := '"' (CHAR | INTERPOLATION)* '"'
INTERPOLATION := '${' EXPRESSION '}'
```

### List Literals
```
LIST := '[' (expression (',' expression)*)? ']'
```

```osprey
let numbers = [1, 2, 3, 4]
let names   = ["Alice", "Bob", "Charlie"]
let pair    = [x, y]
```

## Operators

### Arithmetic Operators

`+`, `-`, `*`, `/`, `%`. All arithmetic returns `Result`; full signatures and semantics are in [Error Handling](/spec/0013-errorhandling/).

### Comparison Operators
- `==` Equality
- `!=` Inequality  
- `<` Less than
- `>` Greater than
- `<=` Less than or equal
- `>=` Greater than or equal

### Logical Operators
- `&&` Logical AND (short-circuit evaluation)
- `||` Logical OR (short-circuit evaluation)
- `!` Logical NOT

### Assignment Operator
- `=` Assignment

### Other Operators
- `->` Function return type
- `=>` Lambda body and match arm
- `|` Union type separator
- `|>` Pipe
- `!` Effect-set marker on a function type

## Delimiters

- `(` `)` Parentheses
- `{` `}` Braces
- `[` `]` Brackets
- `,` Comma
- `:` Colon
- `;` Semicolon
- `.` Dot