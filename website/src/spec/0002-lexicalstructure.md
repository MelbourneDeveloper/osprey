---
layout: page
title: "Lexical Structure"
description: "Osprey Language Specification: Lexical Structure"
date: 2025-07-05
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0002-lexicalstructure/"
---

2. [Lexical Structure](0002-LexicalStructure)
   - [Identifiers](#21-identifiers)
   - [Keywords](#22-keywords)
   - [Literals](#23-literals)
   - [Operators](#24-operators)
   - [Delimiters](#25-delimiters)

## 2. Lexical Structure

### 2.1 Identifiers

Identifiers must start with a letter or underscore, followed by letters, digits, or underscores:

```
ID := [a-zA-Z_][a-zA-Z0-9_]*
```

### 2.2 Keywords

Reserved keywords:
```
fn let mut type import match if then else case of extern
```

### 2.3 Literals

#### Integer Literals
```
INTEGER := [0-9]+
```

#### String Literals
```
STRING := '"' (CHAR | ESCAPE_SEQUENCE)* '"'
ESCAPE_SEQUENCE := '\n' | '\t' | '\r' | '\\' | '\"'
```

#### Interpolated String Literals
```
INTERPOLATED_STRING := '"' (CHAR | INTERPOLATION)* '"'
INTERPOLATION := '${' EXPRESSION '}'
```

#### Immutable List Literals
```
LIST := '[' (expression (',' expression)*)? ']'
```

**Examples:**
```osprey
let numbers = [1, 2, 3, 4]  // Fixed size: 4 elements
let names = ["Alice", "Bob", "Charlie"]  // Fixed size: 3 elements
let pair = [x, y]  // Fixed size: 2 elements
```

### 2.4 Operators

#### Arithmetic Operators (All Safe by Default)
- `+` Addition: `(Int, Int) -> Result<Int, MathError>`
- `-` Subtraction: `(Int, Int) -> Result<Int, MathError>`
- `*` Multiplication: `(Int, Int) -> Result<Int, MathError>`
- `/` Division: `(Int, Int) -> Result<Int, MathError>`

**CRITICAL**: All arithmetic operators return Result types to handle overflow, underflow, and division by zero.

#### Assignment Operator
- `=` Assignment

#### Other Operators
- `=>` Lambda/Match arm arrow
- `|` Union type separator
- `::` Type annotation

### 2.5 Delimiters

- `(` `)` Parentheses
- `{` `}` Braces
- `[` `]` Brackets
- `,` Comma
- `:` Colon
- `;` Semicolon
- `.` Dot