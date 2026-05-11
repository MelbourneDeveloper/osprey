# Iterators and Iteration

Osprey has no `for`, `while`, or `loop` construct. Iteration is expressed as composition of `range`, `forEach`, `map`, `filter`, and `fold` using the pipe operator `|>`.

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

`|>` passes its left operand as the first argument to the function on its right.

```osprey
5 |> double |> print                                       // print(double(5))
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
