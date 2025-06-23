---
layout: page
title: "|> (Pipe Operator Operator)"
description: "Takes the result of the left expression and passes it as the first argument to the right function. Enables functional programming and method chaining."
---

**Description:** Takes the result of the left expression and passes it as the first argument to the right function. Enables functional programming and method chaining.

## Example

```osprey
5 |> double |> print  // (5) -> double -> print
range(1, 10) |> forEach(print)
```
