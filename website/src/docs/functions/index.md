---
layout: page
title: "Built-in Functions"
description: "Complete reference for all built-in functions in Osprey"
---

All built-in functions available in Osprey.

## [contains](contains/)

**Signature:** `contains(haystack: string, needle: string) -> bool`

Checks if a string contains a substring.

## [filter](filter/)

**Signature:** `filter(iterator: iterator, predicate: function) -> iterator`

Filters elements in an iterator based on a predicate function.

## [fold](fold/)

**Signature:** `fold(iterator: iterator, initial: any, fn: function) -> any`

Reduces an iterator to a single value using an accumulator function.

## [forEach](foreach/)

**Signature:** `forEach(iterator: iterator, fn: function) -> int`

Applies a function to each element in an iterator.

## [input](input/)

**Signature:** `input() -> int`

Reads an integer from the user's input.

## [length](length/)

**Signature:** `length(s: string) -> int`

Returns the length of a string.

## [map](map/)

**Signature:** `map(iterator: iterator, fn: function) -> iterator`

Transforms each element in an iterator using a function, returning a new iterator.

## [print](print/)

**Signature:** `print(value: any) -> int`

Prints a value to the console. Automatically converts the value to a string representation.

## [range](range/)

**Signature:** `range(start: int, end: int) -> iterator`

Creates an iterator that generates numbers from start to end (exclusive).

## [substring](substring/)

**Signature:** `substring(s: string, start: int, end: int) -> string`

Extracts a substring from start to end index.

## [toString](tostring/)

**Signature:** `toString(value: any) -> string`

Converts a value to its string representation.

