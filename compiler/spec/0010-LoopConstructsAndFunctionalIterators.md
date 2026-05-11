# Iterators and Iteration

Osprey has no `for`, `while`, or `loop` construct. Iteration is expressed as composition of `range`, `forEach`, `map`, `filter`, and `fold` using the pipe operator `|>`.

## Core Iterator Functions

All multi-argument iterator functions are intended to be used through the pipe; their full signatures are listed for reference. Direct calls require named arguments per [Function Calls](0005-FunctionCalls.md).

### `range(start: int, end: int) -> Iterator<int>`
Generates integers from `start` (inclusive) to `end` (exclusive).

```osprey
range(start: 1,  end: 5)    // 1, 2, 3, 4
range(start: 0,  end: 3)    // 0, 1, 2
range(start: 10, end: 13)   // 10, 11, 12
```

### `forEach(iterator: Iterator<T>, function: fn(T) -> U) -> unit`
Applies `function` to each element for its side effects.

```osprey
range(start: 1, end: 5) |> forEach(print)
```

### `map(iterator: Iterator<T>, function: fn(T) -> U) -> Iterator<U>`
Transforms each element.

```osprey
range(start: 1, end: 5) |> map(double)
```

### `filter(iterator: Iterator<T>, predicate: fn(T) -> bool) -> Iterator<T>`
Keeps elements that satisfy `predicate`.

```osprey
range(start: 1, end: 10) |> filter(isEven)
```

### `fold(iterator: Iterator<T>, initial: U, function: fn(U, T) -> U) -> U`
Reduces an iterator to a single value.

```osprey
range(start: 1, end: 5) |> fold(initial: 0, function: add)   // 0+1+2+3+4 = 10
```

## Pipe Operator

`|>` passes its left operand as the first argument to the function on its right.

```osprey
5 |> double |> print                                                 // print(double(5))
range(start: 1, end: 10) |> forEach(print)
range(start: 0, end: 20) |> filter(isEven) |> map(double) |> forEach(print)
```

## Stream Fusion

Chains of `map`, `filter`, `forEach`, and `fold` over an iterator are fused at compile time into a single loop with no intermediate collections. The chain

```osprey
range(start: 1, end: 5) |> map(double) |> filter(isEven) |> forEach(print)
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
range(start: 1, end: 20)
  |> map(square)
  |> filter(isEven)
  |> fold(initial: 0, function: add)
  |> print

// Pipeline of named stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```
