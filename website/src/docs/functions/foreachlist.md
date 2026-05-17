---
layout: page
title: "forEachList (Function)"
description: "Apply function to every element of list. Phase 7 of collections plan."
---

**Signature:** `forEachList(list: List<T>, function: fn(T) -> Unit) -> List<T>`

**Description:** Apply function to every element of list. Phase 7 of collections plan.

## Parameters

- **list** (List): The list
- **function** (T -> Unit): Function applied per element

**Returns:** List

## Example

```osprey
forEachList(xs, print)
```
