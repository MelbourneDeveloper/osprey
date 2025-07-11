---
layout: page
title: "Loop Constructs and Functional Iterators"
description: "Osprey Language Specification: Loop Constructs and Functional Iterators"
date: 2025-07-11
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0011-loopconstructsandfunctionaliterators/"
---

11. [Loop Constructs and Functional Iterators](0011-LoopConstructsAndFunctionalIterators.md)
    - [Functional Iteration Philosophy](#111-functional-iteration-philosophy)
    - [Core Iterator Functions](#112-core-iterator-functions)
        - [`range(start: int, end: int) -> Iterator<int>`](#rangestart-int-end-int---iteratorint)
        - [`forEach(iterator: Iterator<T>, function: T -> U) -> T`](#foreachiterator-iteratort-function-t---u---t)
        - [`map(iterator: Iterator<T>, function: T -> U) -> U`](#mapiterator-iteratort-function-t---u---u)
        - [`filter(iterator: Iterator<T>, predicate: T -> bool) -> T`](#filteriterator-iteratort-predicate-t---bool---t)
        - [`fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`](#folditerator-iteratort-initial-u-function-u-t---u---u)
    - [Pipe Operator](#113-pipe-operator)
        - [`|>` - Pipe Operator](#---pipe-operator)
    - [Functional Programming Patterns](#114-functional-programming-patterns)
        - [Chaining Pattern](#chaining-pattern)
        - [Side Effect Pattern](#side-effect-pattern)
        - [Data Transformation Pattern](#data-transformation-pattern)
    - [Why No Imperative Loops?](#115-why-no-imperative-loops)

# 11. Loop Constructs and Functional Iterators

🚧 **PARTIAL IMPLEMENTATION**: Basic iterator functions (`range`, `forEach`, `map`, `filter`, `fold`) are implemented and working. The pipe operator (`|>`) is implemented.

## 11.1 Functional Iteration Philosophy

**Osprey is a functional language and does NOT support imperative loop constructs.** Instead, Osprey provides powerful functional iteration patterns that are:

1. **Composable** - Functional iterators can be chained with `|>`
2. **Safe** - No mutable state, no infinite loop bugs
3. **Concurrent** - Fibers provide better parallelism than loops
4. **Testable** - Pure functions are easier to test than stateful loops

## 11.2 Core Iterator Functions

### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator that generates integers from start (inclusive) to end (exclusive).

```osprey
range(1, 5)    // generates 1, 2, 3, 4
range(0, 3)    // generates 0, 1, 2
range(10, 13)  // generates 10, 11, 12
```

### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element in an iterator for side effects.

```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
forEach(range(0, 3), double)           // calls double(0), double(1), double(2)
```

### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element in an iterator by applying a function.

```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
map(range(0, 3), square)      // applies square to 0, 1, 2
```

### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements from an iterator based on a predicate function.

```osprey
range(1, 10) |> filter(isEven)
filter(range(0, 20), isPositive)
```

### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value by repeatedly applying a function.

```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
fold(range(1, 6), 1, multiply)       // product: 1*1*2*3*4*5 = 120
```

## 11.3 Pipe Operator

### `|>` - Pipe Operator
The pipe operator creates clean, readable function composition by allowing you to chain operations from left to right.

```osprey
// Basic piping
5 |> double |> print                 // Equivalent to: print(double(5))

// Iterator chaining
range(1, 10) |> forEach(print)
range(1, 5) |> map(square) |> fold(0, add)

// Complex chains
range(0, 20) |> filter(isEven) |> map(double) |> forEach(print)
```

## 11.4 Functional Programming Patterns

### Chaining Pattern
```osprey
// Transform -> Filter -> Aggregate
range(1, 20)
  |> map(square)           // Square each number
  |> filter(isEven)        // Keep only even results
  |> fold(0, add)          // Sum them up
  |> print                 // Print the result
```

### Side Effect Pattern
```osprey
// Process each element for side effects
range(1, 100)
  |> filter(isPrime)
  |> forEach(print)        // Print each prime
```

### Data Transformation Pattern
```osprey
// Transform data through multiple stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```

## 11.5 Why No Imperative Loops?

**Anti-Pattern:**
```osprey
// ❌ WRONG - Imperative loops (NOT SUPPORTED)
loop {
    let request = getRequest()
    processRequest(request)
}
```

**Functional Pattern:**
```osprey
// ✅ CORRECT - Functional approach
fn serverHandler() -> Unit = {
    requestStream() 
    |> map(processRequest)
    |> forEach(sendResponse)
}
```

This functional approach provides better maintainability, testability, and performance than traditional imperative loops.