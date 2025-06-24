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

### 17.3 Process Operations

#### `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessResult, string>`
Spawns external process with asynchronous stdout/stderr collection via callbacks.

**Key Features:**
- **Event-driven output collection** - Uses callbacks from the runtime for things like new stdout
- **Non-blocking execution** - Runs in background threads via fiber runtime
- **Real-time stdout/stderr** - Process output is sent to Osprey via callbacks as it's generated
- **Thread-safe** - Safe for concurrent use with multiple processes

**Architecture:**
The process runtime follows the same callback pattern as the HTTP server:
1. Osprey calls `spawnProcess(command)` 
2. C runtime spawns process and monitoring thread
3. C runtime calls back to Osprey with stdout/stderr events as they occur
4. Process completion triggers exit event callback

**ProcessResult Type:**
```osprey
type ProcessResult = {
    processId: int
}
```

**Event Types (handled by C runtime callbacks):**
- `PROCESS_STDOUT_DATA` - New stdout data available
- `PROCESS_STDERR_DATA` - New stderr data available  
- `PROCESS_EXIT` - Process completed with exit code

**Stdout Callback Handling:**

The process runtime uses callback functions to deliver stdout/stderr data as it's generated. The C runtime calls back into Osprey functions for real-time process events.

🚨 **CALLBACK IS MANDATORY!** The callback parameter is **REQUIRED** for `spawnProcess` - it cannot be omitted.

```osprey
// Define the callback function that C runtime will call for ALL process events
fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("[STDOUT] Process ${toString(processID)}: ${data}")
        2 => print("[STDERR] Process ${toString(processID)}: ${data}")
        3 => print("[EXIT] Process ${toString(processID)} exited with code: ${data}")
        _ => print("[UNKNOWN] Process ${toString(processID)} event ${toString(eventType)}: ${data}")
    }
}

// Process with callback-based stdout collection (CALLBACK IS MANDATORY!)
let result = spawnProcess("echo 'Hello from callback!'", processEventHandler)
match result {
    Success { value } => {
        print("Process spawned with ID: ${toString(value)}")
        let exitCode = awaitProcess(value)
        print("Process finished with exit code: ${toString(exitCode)}")
        cleanupProcess(value)
        print("Process cleaned up")
    }
    Error { message } => print("Failed to spawn process")
}
```

**Advanced Usage:**
```osprey
// Wait for process completion and get exit code
let result = spawnProcess("gcc myprogram.c -o myprogram")
match result {
    Success { value } => {
        let exitCode = awaitProcess(value.processId)
        cleanupProcess(value.processId)
        print("Compilation finished with exit code: ${toString(exitCode)}")
    }
    Error { message } => print("Compilation failed: ${message}")
}
```

#### `awaitProcess(processId: int) -> int`
Waits for process completion and returns exit code.

**Parameters:**
- `processId: int` - Process ID returned by spawnProcess

**Returns:** `int` - Process exit code (0 = success, non-zero = error)

#### `cleanupProcess(processId: int) -> void`
Cleans up process resources after completion.

**Parameters:**  
- `processId: int` - Process ID to clean up

**Note:** Always call this after `awaitProcess` to prevent memory leaks.

### 17.2 Functional Iterator Functions

#### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator that generates integers from start (inclusive) to end (exclusive). Used with functional iterator functions like forEach, map, filter, and fold.

**Parameters:**
- `start: int` - Starting value (inclusive)
- `end: int` - Ending value (exclusive)

**Returns:** `Iterator<int>` - Iterator struct containing start and end values

**Examples:**
```osprey
range(1, 5)    // generates 1, 2, 3, 4
range(0, 3)    // generates 0, 1, 2
range(10, 13)  // generates 10, 11, 12
range(1, 10) |> forEach(print)
```

#### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element in an iterator for side effects. This is the primary way to iterate through ranges and apply operations to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to traverse (usually from range())
- `function: T -> U` - Function to apply to each element

**Returns:** `T` - Final counter value after iteration

**Examples:**
```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
forEach(range(0, 3), double)           // calls double(0), double(1), double(2)
range(1, 10) |> forEach(square)
forEach(range(-2, 3), print)          // prints -2, -1, 0, 1, 2
```

#### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element in an iterator by applying a function. Returns the result of the transformation function applied to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to transform (usually from range())
- `function: T -> U` - Transformation function to apply

**Returns:** `U` - Result of applying function to each element

**Examples:**
```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
map(range(0, 3), square)      // applies square to 0, 1, 2
range(1, 6) |> map(addFive)
```

#### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements from an iterator based on a predicate function. Only elements where the predicate returns true are processed.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to filter (usually from range())
- `predicate: T -> bool` - Function that returns true for elements to keep

**Returns:** `T` - Filtered results

**Examples:**
```osprey
range(1, 10) |> filter(isEven)
filter(range(0, 20), isPositive)
range(-5, 6) |> filter(isGreaterThanZero)
```

#### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value by repeatedly applying a function. Also known as reduce or accumulate in other languages.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to reduce (usually from range())
- `initial: U` - Initial accumulator value
- `function: (U, T) -> U` - Function that combines accumulator with each element

**Returns:** `U` - Final accumulated value

**Examples:**
```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
fold(range(1, 6), 1, multiply)       // product: 1*1*2*3*4*5 = 120
range(0, 10) |> fold(0, max)
```

### 17.3 Pipe Operator

#### `|>` - Pipe Operator
The pipe operator takes the result of the left expression and passes it as the first argument to the right function. This enables elegant functional programming and method chaining.

**Syntax:** `expression |> function`

**Type:** `T |> (T -> U) -> U`

**Description:**
The pipe operator creates clean, readable function composition by allowing you to chain operations from left to right, making the data flow explicit and natural to read.

**Rules:**
- Left side can be any expression
- Right side must be a function or function call
- Creates clean, readable function composition
- Enables Haskell/Rust-style functional programming

**Examples:**
```osprey
// Basic piping
5 |> double |> print                 // Equivalent to: print(double(5))

// Iterator chaining
range(1, 10) |> forEach(print)
range(1, 5) |> map(square) |> fold(0, add)

// Complex chains
range(0, 20) |> filter(isEven) |> map(double) |> forEach(print)

// Multiple operations
let result = input() |> double |> square |> toString

// Nested operations
range(1, 10) 
  |> map(square) 
  |> filter(isEven) 
  |> fold(0, add) 
  |> print
```

### 17.4 Functional Programming Patterns

The combination of iterator functions and the pipe operator enables powerful functional programming patterns:

#### Chaining Pattern
```osprey
// Transform -> Filter -> Aggregate
range(1, 20)
  |> map(square)           // Square each number
  |> filter(isEven)        // Keep only even results
  |> fold(0, add)          // Sum them up
  |> print                 // Print the result
```

#### Side Effect Pattern
```osprey
// Process each element for side effects
range(1, 100)
  |> filter(isPrime)
  |> forEach(print)        // Print each prime
```

#### Data Transformation Pattern
```osprey
// Transform data through multiple stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```

### 17.5 Fiber Types and Concurrency

Osprey provides lightweight concurrency through fiber types. Unlike traditional function-based approaches, fibers are proper type instances constructed using Osprey's standard type construction syntax.

#### Core Fiber Types

**`Fiber<T>`** - A lightweight concurrent computation that produces a value of type T
**`Channel<T>`** - A communication channel for passing values of type T between fibers

#### Fiber Construction

Fibers are created using standard type construction syntax:

```osprey
// Create a fiber that computes a value
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(n: 1000) 
}

// Create a fiber with more complex computation
let worker = Fiber<String> { 
    computation: fn() => {
        processData()
        "completed"
    }
}

// Create a parameterized fiber
let calculator = Fiber<Int> { 
    computation: fn() => multiply(x: 10, y: 20) 
}
```

#### Spawn Syntax Sugar

For convenience, Osprey provides `spawn` as syntax sugar for creating and immediately starting a fiber:

```osprey
// Using spawn (syntax sugar)
let result = spawn 42

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => 42 }
let result = fiber

// More complex spawn
let computation = spawn (x * 2 + y)

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => x * 2 + y }
let computation = fiber
```

The `spawn` keyword immediately evaluates the expression in a new fiber context, making it convenient for quick concurrent computations without the full type construction syntax.

#### Channel Construction

Channels are created using type construction syntax:

```osprey
// Unbuffered (synchronous) channel
let sync_channel = Channel<Int> { capacity: 0 }

// Buffered (asynchronous) channel  
let async_channel = Channel<String> { capacity: 10 }

// Large buffer channel
let buffer_channel = Channel<Int> { capacity: 100 }
```

#### Fiber Operations

Once created, fibers and channels are manipulated using functional operations:

**`await(fiber: Fiber<T>) -> T`** - Wait for fiber completion and get result
**`send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`** - Send value to channel
**`recv(channel: Channel<T>) -> Result<T, ChannelError>`** - Receive value from channel
**`yield() -> Unit`** - Voluntarily yield control to scheduler

```osprey
// Create and await a fiber
let task = Fiber<Int> { computation: fn() => heavyComputation() }
let result = await(task)

// Channel communication
let ch = Channel<String> { capacity: 5 }
send(ch, "hello")
let message = recv(ch)

// Yielding control
yield()
```

#### Complete Fiber Example

```osprey
// Producer fiber
let producer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        send(ch, 1)
        send(ch, 2) 
        send(ch, 3)
    }
}

// Consumer fiber
let consumer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        let value1 = recv(ch)
        let value2 = recv(ch)
        let value3 = recv(ch)
        print("Received: ${value1}, ${value2}, ${value3}")
    }
}

// Start both fibers
await(producer)
await(consumer)
```

#### Select Expression for Channel Multiplexing

The `select` expression allows waiting on multiple channel operations:

```osprey
let ch1 = Channel<String> { capacity: 1 }
let ch2 = Channel<Int> { capacity: 1 }

let result = select {
    msg => recv(ch1) => process_string(msg)
    num => recv(ch2) => process_number(num)
    _ => timeout_handler()
}
```

#### Rust Interoperability

Osprey fibers are designed to interoperate with Rust's async/await system:

```osprey
// Osprey fiber that calls Rust async function
extern fn rust_async_task() -> Future<Int>

let osprey_task = Fiber<Int> {
    computation: fn() => await(rust_async_task())
}

let result = await(osprey_task)
```

#### Fiber-Isolated Modules

Each fiber gets its own isolated instance of modules, preventing data races:

```osprey
module Counter {
    let mut count = 0
    fn increment() -> Int = { count = count + 1; count }
    fn get() -> Int = count
}

// Each fiber has its own Counter instance
let fiber1 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}
let fiber2 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}

// These will both return 1, not 1 and 2
let result1 = await(fiber1)  // 1
let result2 = await(fiber2)  // 1 (separate instance)
```