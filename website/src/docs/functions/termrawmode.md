---
layout: page
title: "termRawMode (Function)"
description: "Enables (1) or disables (0) raw terminal input mode, so keypresses arrive unbuffered."
---

**Signature:** `termRawMode(enabled: int) -> Unit`

**Description:** Enables (1) or disables (0) raw terminal input mode, so keypresses arrive unbuffered.

## Parameters

- **enabled** (int): 1 to enable raw mode, 0 to restore cooked mode

**Returns:** Unit

## Example

```osprey
termRawMode(1)
```
