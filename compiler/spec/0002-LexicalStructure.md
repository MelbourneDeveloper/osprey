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

**Examples:**
```osprey
let numbers = [1, 2, 3, 4]  // Fixed size: 4 elements
let names = ["Alice", "Bob", "Charlie"]  // Fixed size: 3 elements
let pair = [x, y]  // Fixed size: 2 elements
```

## Operators

### Arithmetic Operators

All arithmetic operators return `Result` types to handle overflow, underflow, and division by zero.

**Integer Operations:**
- `+` Addition: `(int, int) -> Result<int, MathError>`
- `-` Subtraction: `(int, int) -> Result<int, MathError>`
- `*` Multiplication: `(int, int) -> Result<int, MathError>`
- `/` Division: `(int, int) -> Result<float, MathError>` — always returns float
- `%` Modulo: `(int, int) -> Result<int, MathError>`

**Floating-Point Operations:**
- `+`, `-`, `*`, `/`, `%`: `(float, float) -> Result<float, MathError>`

**Type Safety:**
- No automatic type promotion between int and float
- Use `toFloat(int)` and `toInt(float)` for explicit conversion
- Division `/` always returns float, even for integer operands

**Examples:**
```osprey
let sum = 5 + 3           // Result<int, MathError>
let quotient = 10 / 3     // Result<float, MathError> - returns 3.333...
let remainder = 10 % 3    // Result<int, MathError> - returns 1

let precise = 10.0 / 3.0  // Result<float, MathError>
let divZero = 10 / 0      // Result<float, MathError> - Error(DivisionByZero)
```

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