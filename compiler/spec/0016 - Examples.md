 ## 16. Examples

The `examples/` directory contains a variety of sample programs demonstrating Osprey's features. These examples are tested as part of the standard build process to ensure they remain up-to-date and functional.

## 17. Built-in Functions Reference

### 17.1 Basic I/O Functions

#### `print(value: int | string | bool) -> int`
Prints the given value to standard output with automatic type conversion.

**Parameters:**
- `value: int | string | bool` - The value to print (int, bool, string, or expression)

**Returns:** `int` - Exit code from puts function

**Examples:**
```osprey
print("Hello World")
print(42)
print(true)
print(x + y)
```

#### `input() -> int`
Reads an integer from stdin. Blocks until user enters a number.

**Parameters:** None

**Returns:** `int` - The number entered by user

**Examples:**
```osprey
let x = input()
let age = input()
```

#### `toString(value: int | string | bool) -> string`
Converts any value to its string representation.

#### `length(s: string) -> Result<int, StringError>`
🚨 **CRITICAL**: Returns the length of a string wrapped in a Result type for safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match length("hello") {
    Success { value } => print("Length: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `contains(haystack: string, needle: string) -> Result<bool, StringError>`
🚨 **CRITICAL**: Checks if a string contains a substring, returns Result for safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match contains("hello", "ell") {
    Success { value } => print("Found: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `substring(s: string, start: int, end: int) -> Result<string, StringError>`
🚨 **CRITICAL**: Extracts a substring from start to end, returns Result for bounds safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match substring("hello", 1, 3) {
    Success { value } => print("Substring: ${value}")
    Error { message } => print("Error: ${message}")
}
```

**FUNDAMENTAL PRINCIPLE**: All string operations that could conceptually fail MUST return Result types. This enforces explicit error handling and prevents runtime panics.

### 17.2 File System Functions

#### `writeFile(path: string, content: string) -> Result<Success, string>`
Writes content to a file.

#### `readFile(path: string) -> Result<string, string>`
Reads file content as string.

#### `deleteFile(path: string) -> Result<Success, string>`
Deletes a file.

#### `createDirectory(path: string) -> Result<Success, string>`
Creates a directory.

#### `fileExists(path: string) -> bool`
Checks if file exists.