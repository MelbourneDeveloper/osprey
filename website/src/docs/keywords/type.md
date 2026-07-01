---
layout: page
title: "type (Keyword)"
description: "Type declaration keyword. Used to define custom types and type aliases."
---

**Description:** Type declaration keyword. Used to define custom types and type aliases.

## Example

```osprey
type UserId = int
type Status = Active | Inactive
type User = { name: string, age: int }
```

```osprey-ml
type UserId = int
type Status = Active | Inactive
type User =
    name : string
    age : int
```
