---
layout: page
title: "13. Built-in Functions"
description: "Osprey Language Specification: 13. Built-in Functions"
date: 2025-06-26
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0013-built-infunctions/"
---

## 13.1 Basic I/O Functions

### `print(value: int | string | bool) -> int`
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
match contains("hello", "ell") {
    Success { value } => print("Found: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `substring(s: string, start: int, end: int) -> Result<string, StringError>`
Extracts a substring from start to end.

## 13.2 File System Functions

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

## 13.3 Process Operations

### `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessResult, string>`
Spawns external process with asynchronous stdout/stderr collection via callbacks.

```osprey
fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("[STDOUT] ${data}")
        2 => print("[STDERR] ${data}")
        3 => print("[EXIT] Code: ${data}")
        _ => print("[UNKNOWN] ${data}")
    }
}

let result = spawnProcess("echo 'Hello'", processEventHandler)
```

### `awaitProcess(processId: int) -> int`
Waits for process completion and returns exit code.

### `cleanupProcess(processId: int) -> void`
Cleans up process resources after completion.

## 13.4 Functional Programming

### Iterator Functions

#### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator from start (inclusive) to end (exclusive).

```osprey
range(1, 5)    // generates 1, 2, 3, 4
```

#### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element for side effects.

```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
```

#### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element by applying a function.

```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
```

#### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements based on a predicate function.

```osprey
range(1, 10) |> filter(isEven)
```

#### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value.

```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
```

### Pipe Operator `|>`

The pipe operator passes the left expression as the first argument to the right function.

```osprey
5 |> double |> print                 // Equivalent to: print(double(5))
range(1, 10) |> map(square) |> filter(isEven) |> forEach(print)
```

## 13.5 HTTP Functions

HTTP functions for server and client operations are documented in [Chapter 15 - HTTP](0015-HTTP.md).

## 13.6 WebSocket Functions

WebSocket functions for real-time bidirectional communication are documented in [Chapter 16 - WebSockets](0016-WebSockets.md).

## 13.7 Fiber and Concurrency Functions

Osprey provides lightweight concurrency through fibers.

### Fiber Types

```osprey
// Create a fiber
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(1000) 
}

// Spawn syntax sugar
let result = spawn 42

// Channels for communication
let ch = Channel<String> { capacity: 10 }
```

### Fiber Operations

#### `await(fiber: Fiber<T>) -> T`
Wait for fiber completion and get result.

#### `send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`
Send value to channel.

#### `recv(channel: Channel<T>) -> Result<T, ChannelError>`
Receive value from channel.

#### `yield() -> Unit`
Voluntarily yield control to scheduler.

### Example Usage

```osprey
// Producer-consumer pattern
let ch = Channel<Int> { capacity: 3 }

let producer = spawn {
    send(ch, 1)
    send(ch, 2)
    send(ch, 3)
}

let consumer = spawn {
    let value1 = recv(ch)
    let value2 = recv(ch)
    let value3 = recv(ch)
    print("Received values")
}

await(producer)
await(consumer)
```

## 13.8 Functional Programming Examples

Combining functional programming capabilities for data processing:

```osprey
fn main() -> Int = {
    // Calculate sum of squares of even numbers from 1 to 10
    let evenSquareSum = range(1, 11)
        |> filter(isEven)
        |> map(square)
        |> fold(0, add)
    
    print("Sum of squares of even numbers: ${toString(evenSquareSum)}")
    
    // Process user data with functional pipeline
    print("Processing user data:")
    range(1, 6)
        |> map(createUserData)
        |> forEach(print)
    
    // Concurrent processing with fibers
    let ch = Channel<String> { capacity: 3 }
    
    let producer = spawn {
        range(1, 4) |> forEach(fn(i) => send(ch, "Message ${toString(i)}"))
    }
    
    let consumer = spawn {
        range(1, 4) |> forEach(fn(_) => {
            match recv(ch) {
                Success { value } => print("Received: ${value}")
                Error { message } => print("Error: ${message}")
            }
        })
    }
    
    await(producer)
    await(consumer)
    
    0
}

fn isEven(x: Int) -> Bool = x % 2 == 0
fn square(x: Int) -> Int = x * x
fn add(a: Int, b: Int) -> Int = a + b

fn createUserData(id: Int) -> String = 
    "{\"id\": ${toString(id)}, \"name\": \"User${toString(id)}\"}"