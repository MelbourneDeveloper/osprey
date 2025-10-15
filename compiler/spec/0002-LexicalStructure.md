# Lexical Structure

- [Identifiers](#identifiers)
- [Keywords](#keywords)
- [Literals](#literals)
- [Operators](#operators)
- [Delimiters](#delimiters)

## Identifiers

Start with letter or underscore, followed by letters, digits, or underscores.
```
ID := [a-zA-Z_][a-zA-Z0-9_]*
```

## Keywords

```
fn let mut type import match if then else case of extern
```

## Literals

### Integer Literals
```
INTEGER := [0-9]+
```

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

**Examples:**
```osprey
let numbers = [1, 2, 3, 4]  // Fixed size: 4 elements
let names = ["Alice", "Bob", "Charlie"]  // Fixed size: 3 elements
let pair = [x, y]  // Fixed size: 2 elements
```

## Operators

### Arithmetic Operators
- `+` Addition: `(Int, Int) -> Result<Int, MathError>`
- `-` Subtraction: `(Int, Int) -> Result<Int, MathError>`
- `*` Multiplication: `(Int, Int) -> Result<Int, MathError>`
- `/` Division: `(Int, Int) -> Result<Int, MathError>`

All arithmetic operators return Result types to handle overflow, underflow, and division by zero.

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
- `=>` Lambda/Match arm arrow
- `|` Union type separator
- `::` Type annotation

## Delimiters

- `(` `)` Parentheses
- `{` `}` Braces
- `[` `]` Brackets
- `,` Comma
- `:` Colon
- `;` Semicolon
- `.` Dot