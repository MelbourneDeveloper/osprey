# Built-in Functions

Reference for built-in functions available in every Osprey program. Operations that can fail return `Result`; see [Error Handling](0013-ErrorHandling.md).

## Basic I/O Functions

```osprey
print(value: int | string | bool) -> int
```
Prints values to standard output with automatic type conversion.

```osprey
print("Hello World")
print(42)
print(true)
```

### `input() -> int`
Reads an integer from stdin.

```osprey
let x = input()
```

### `toString(value: int | string | bool) -> string`
Converts any value to its string representation.

### String Functions

#### `length(s: string) -> Result<int, StringError>`
Returns string length. Requires pattern matching for safety.

```osprey
match length("hello") {
    Success { value } => print("Length: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `contains(haystack: string, needle: string) -> Result<bool, StringError>`
Checks if a string contains a substring.

```osprey
match contains(haystack: "hello", needle: "ell") {
    Success { value }   => print("Found: ${value}")
    Error   { message } => print("Error: ${message}")
}
```

#### `substring(s: string, start: int, end: int) -> Result<string, StringError>`
Extracts a substring from `start` (inclusive) to `end` (exclusive).

## File System Functions

### `writeFile(path: string, content: string) -> Result<Success, string>`
Writes content to a file.

### `readFile(path: string) -> Result<string, string>`
Reads file content as string.

### `deleteFile(path: string) -> Result<Success, string>`
Deletes a file.

### `createDirectory(path: string) -> Result<Success, string>`
Creates a directory.

### `fileExists(path: string) -> bool`
Checks if file exists.

## Process Operations

### `spawnProcess(command: string, callback: fn(int, int, string) -> unit) -> Result<ProcessResult, string>`
Spawns an external process. The callback is invoked for each stdout/stderr line and on exit.

```osprey
fn processEventHandler(processID: int, eventType: int, data: string) -> unit = match eventType {
    1 => print("[STDOUT] ${data}")
    2 => print("[STDERR] ${data}")
    3 => print("[EXIT] Code: ${data}")
    _ => print("[UNKNOWN] ${data}")
}

let result = spawnProcess(command: "echo 'Hello'", callback: processEventHandler)
```

### `awaitProcess(processId: int) -> int`
Waits for process completion and returns the exit code.

### `cleanupProcess(processId: int) -> unit`
Releases process resources.

## Iterators and Pipe

`range`, `forEach`, `map`, `filter`, `fold`, and `|>` are documented in [Iterators and Iteration](0010-LoopConstructsAndFunctionalIterators.md).

## HTTP

See [HTTP](0014-HTTP.md).

## WebSockets

See [WebSockets](0015-WebSockets.md).

## Fibers and Channels

`spawn`, `await`, `send`, `recv`, `yield`, `Fiber<T>`, `Channel<T>` are documented in [Fibers and Concurrency](0011-LightweightFibersAndConcurrency.md).