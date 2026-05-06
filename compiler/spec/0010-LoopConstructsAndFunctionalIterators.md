11. [Loop Constructs and Functional Iterators](0011-LoopConstructsAndFunctionalIterators.md)
    - [Functional Iteration Philosophy](#functional-iteration-philosophy)
    - [Core Iterator Functions](#core-iterator-functions)
        - [`range(start: int, end: int) -> Iterator<int>`](#rangestart-int-end-int---iteratorint)
        - [`forEach(iterator: Iterator<T>, function: T -> U) -> T`](#foreachiterator-iteratort-function-t---u---t)
        - [`map(iterator: Iterator<T>, function: T -> U) -> U`](#mapiterator-iteratort-function-t---u---u)
        - [`filter(iterator: Iterator<T>, predicate: T -> bool) -> T`](#filteriterator-iteratort-predicate-t---bool---t)
        - [`fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`](#folditerator-iteratort-initial-u-function-u-t---u---u)
    - [Pipe Operator](#pipe-operator)
        - [`|>` - Pipe Operator](#---pipe-operator)
    - [Stream Fusion Optimization](#stream-fusion-optimization)
        - [How Stream Fusion Works](#how-stream-fusion-works)
        - [Performance Benefits](#performance-benefits)
        - [Supported Fusion Chains](#supported-fusion-chains)
    - [Functional Programming Patterns](#functional-programming-patterns)
        - [Chaining Pattern](#chaining-pattern)
        - [Side Effect Pattern](#side-effect-pattern)
        - [Data Transformation Pattern](#data-transformation-pattern)
    - [Why No Imperative Loops?](#why-no-imperative-loops)

# Loop Constructs and Functional Iterators

Osprey is a functional language without imperative loop constructs. Iteration uses functional patterns with the core functions `range`, `forEach`, `map`, `filter`, and `fold`, composed with the pipe operator `|>`.

## Functional Iteration

Functional iteration provides composable, safe patterns without mutable state. Fibers handle concurrent iteration, and pure functions simplify testing.

## Core Iterator Functions

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

## Pipe Operator

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

## Stream Fusion Optimization

Osprey implements **stream fusion** - a compile-time optimization that eliminates intermediate data structures when chaining iterator operations. This provides zero-cost abstractions: you write elegant functional code that compiles to the same performance as hand-optimized loops.

### How Stream Fusion Works

When you write:
```osprey
range(1, 5) |> map(double) |> filter(isEven) |> forEach(print)
```

**Without stream fusion** (naive approach):
1. `range(1, 5)` creates array `[1, 2, 3, 4]`
2. `map(double)` creates new array `[2, 4, 6, 8]`
3. `filter(isEven)` creates new array `[2, 4, 6, 8]`
4. `forEach(print)` iterates and prints

**With stream fusion** (Osprey's approach):
- Compiler detects the chain at compile time
- `map()` stores the transformation function, returns range unchanged
- `filter()` stores the predicate function, returns range unchanged
- `forEach()` generates a single optimized loop that applies all transformations inline

The generated LLVM IR is equivalent to:
```c
// Hand-optimized loop - what Osprey generates
for (i = 1; i < 5; i++) {
    value = double(i);           // map applied inline
    if (isEven(value)) {         // filter applied inline
        print(value);            // forEach applied inline
    }
}
```

### Performance Benefits

**Zero-cost abstractions:**
- ✅ No intermediate arrays or memory allocations
- ✅ Single pass through data instead of multiple iterations
- ✅ Better CPU cache utilization
- ✅ Same performance as hand-written optimized loops

**Example:**
```osprey
// Elegant functional code
range(1, 1000000)
  |> map(square)
  |> filter(isEven)
  |> fold(0, add)

// Compiles to single optimized loop with:
// - Zero memory allocations
// - Zero intermediate arrays
// - Optimal CPU cache usage
```

### Supported Fusion Chains

Stream fusion works with any combination of:
- `map()` - Transforms are fused inline
- `filter()` - Predicates are fused as conditional branches
- `forEach()` - Terminal operation that consumes the fused chain
- `fold()` - Terminal operation that consumes and reduces

Multiple transformations and filters can be chained together and will all be fused into a single optimized loop.

## Functional Programming Patterns

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

## Alternative to Imperative Loops

Use functional patterns instead of imperative loops:

```osprey
// Functional approach
fn serverHandler() -> unit = {
    requestStream()
    |> map(processRequest)
    |> forEach(sendResponse)
}
``` 