---
layout: page
title: "match (Keyword)"
description: "Pattern matching expression. Used for destructuring values and control flow based on patterns."
---

**Description:** Pattern matching expression. Used for destructuring values and control flow based on patterns.

## Example

```osprey
match value {
    Some(x) -> x
    None -> 0
}

match status {
    Active -> "User is active"
    Inactive -> "User is inactive"
}
```
