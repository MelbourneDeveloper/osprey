---
layout: page
title: "Iterators and Iteration"
description: "Osprey Language Specification: Iterators and Iteration"
date: 2026-06-30
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0010-loopconstructsandfunctionaliterators/"
---

# Iterators and Iteration

Osprey has no `for`, `while`, or `loop` construct. Iteration is expressed as composition of `range`, `forEach`, `map`, `filter`, and `fold` using the pipe operator `|>`.

## Core Iterator Functions

`Iterator<T>` is the type of a lazily-produced sequence. It exists during type-checking and is always fused away at compile time (see [Stream Fusion](#stream-fusion)); it is never a materialised runtime collection. Like all built-ins, iterator functions take positional arguments in subject-first order (rule 4 of [Function Calls](/spec/0005-functioncalls/#rules)); parameter names below are descriptive only.

### `range(start: int, end: int) -> Iterator<int>`
Generates integers from `start` (inclusive) to `end` (exclusive).

```osprey
range(1, 5)      // 1, 2, 3, 4
range(0, 3)      // 0, 1, 2
range(10, 13)    // 10, 11, 12
```

### `forEach(iterator: Iterator<T>, function: fn(T) -> U) -> unit`
Applies `function` to each element for its side effects.

```osprey
range(1, 5) |> forEach(print)
```

### `map(iterator: Iterator<T>, function: fn(T) -> U) -> Iterator<U>`
Transforms each element.

```osprey
range(1, 5) |> map(double)
```

### `filter(iterator: Iterator<T>, predicate: fn(T) -> bool) -> Iterator<T>`
Keeps elements that satisfy `predicate`.

```osprey
range(1, 10) |> filter(isEven)
```

### `fold(iterator: Iterator<T>, initial: U, function: fn(U, T) -> U) -> U`
Reduces an iterator to a single value.

```osprey
range(1, 5) |> fold(0, add)   // 0+1+2+3+4 = 10
```

## Pipe Operator

`|>` passes its left operand as the first argument to the function on its right.

```osprey
5 |> double |> print                                        // print(double(5))
range(1, 10) |> forEach(print)
range(0, 20) |> filter(isEven) |> map(double) |> forEach(print)
```

## Stream Fusion

Chains of `map`, `filter`, `forEach`, and `fold` over an iterator are fused at compile time into a single loop with no intermediate collections. The chain

```osprey
range(1, 5) |> map(double) |> filter(isEven) |> forEach(print)
```

compiles to one loop that applies `double`, the `isEven` test, and `print` per element — equivalent to:

```c
for (i = 1; i < 5; i++) {
    value = double(i);
    if (isEven(value)) print(value);
}
```

Fusion applies to any chain of `map` and `filter` terminated by `forEach` or `fold`.

## Patterns

```osprey
// Transform → filter → aggregate
range(1, 20)
  |> map(square)
  |> filter(isEven)
  |> fold(0, add)
  |> print

// Pipeline of named stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```