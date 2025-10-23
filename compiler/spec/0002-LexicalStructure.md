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

All arithmetic operators are type-preserving and return `Result` types to handle errors (overflow, underflow, division by zero).

**Integer Arithmetic:**
- `+` Addition: `(int, int) -> Result<int, MathError>`
- `-` Subtraction: `(int, int) -> Result<int, MathError>`
- `*` Multiplication: `(int, int) -> Result<int, MathError>`
- `/` Division: `(int, int) -> Result<int, MathError>` - Returns quotient only (10 / 3 = 3)
- `%` Modulo: `(int, int) -> Result<int, MathError>` - Returns remainder (10 % 3 = 1)

**Floating-Point Arithmetic:**
- `+` Addition: `(float, float) -> Result<float, MathError>`
- `-` Subtraction: `(float, float) -> Result<float, MathError>`
- `*` Multiplication: `(float, float) -> Result<float, MathError>`
- `/` Division: `(float, float) -> Result<float, MathError>` - IEEE 754 division (10.0 / 3.0 = 3.333...)
- `%` Modulo: `(float, float) -> Result<float, MathError>` - IEEE 754 remainder

**Type Safety:**
- No automatic type promotion: cannot mix int and float in operations
- Use `toFloat(int)` to convert int to float: `toFloat(10) / 3.0`
- Use `toInt(float)` to truncate float to int: `toInt(3.7) = 3`

**Examples:**
```osprey
// Integer arithmetic
let sum = 5 + 3           // Result<int, MathError> - Success(8)
let quotient = 10 / 3     // Result<int, MathError> - Success(3)
let remainder = 10 % 3    // Result<int, MathError> - Success(1)

// Floating-point arithmetic
let precise = 10.0 / 3.0  // Result<float, MathError> - Success(3.333...)
let area = 3.14 * 2.5     // Result<float, MathError> - Success(7.85)

// Mixed requires explicit conversion
let mixed = toFloat(10) / 3.0  // Result<float, MathError> - Success(3.333...)

// Error cases
let divZero = 10 / 0      // Result<int, MathError> - Error(DivisionByZero)
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