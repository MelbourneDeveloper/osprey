---
layout: page
title: "Fibers and Concurrency"
description: "Osprey Language Specification: Fibers and Concurrency"
date: 2026-06-18
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0011-lightweightfibersandconcurrency/"
---

# Fibers and Concurrency

Fibers are lightweight concurrent computations. They are constructed as values of `Fiber<T>` and communicate through `Channel<T>`. There are no OS threads exposed to user code; the runtime schedules fibers cooperatively. Values cross fiber boundaries — `spawn` captures and channel `send` — by move or copy, never by sharing ([MEM-FIBER-ISOLATION] in [Memory Management](/spec/0018-memorymanagement/)).

## Status

`spawn`, `await`, `yield`, and basic channel operations are implemented. The `select` expression and the fiber-isolated module system below are planned and not yet wired through code generation.

## Core Types

```osprey
Fiber<T>     // a concurrent computation that produces a value of type T
Channel<T>   // an in-process communication channel carrying values of type T
```

## Constructing Fibers

Fibers are constructed using ordinary record-construction syntax:

```osprey
let task = Fiber<int> {
    computation: fn() => calculatePrimes(n: 1000)
}
```

`spawn <expr>` is sugar for the equivalent `Fiber` construction:

```osprey
let result = spawn 42
// equivalent to:
let result = Fiber<int> { computation: fn() => 42 }
```

## Constructing Channels

```osprey
let sync   = Channel<int>    { capacity: 0  }   // unbuffered (rendezvous)
let buf    = Channel<string> { capacity: 10 }   // buffered
```

## Operations

| Operation                                       | Signature                                            |
| ----------------------------------------------- | ---------------------------------------------------- |
| Wait for a fiber to produce its value           | `await(fiber: Fiber<T>) -> T`                        |
| Send a value to a channel                       | `send(channel: Channel<T>, value: T) -> Result<unit, ChannelError>` |
| Receive a value from a channel                  | `recv(channel: Channel<T>) -> Result<T, ChannelError>` |
| Yield to the scheduler                          | `yield() -> unit`                                    |

## Producer / Consumer Example

```osprey
let ch = Channel<int> { capacity: 3 }

let producer = spawn {
    range(1, 4) |> forEach(fn(i) => send(ch, i))
}

let consumer = spawn {
    range(1, 4) |> forEach(fn(i) => match recv(ch) {
        Success { value }   => print("got ${value}")
        Error   { message } => print("recv error: ${message}")
    })
}

await(producer)
await(consumer)
```

## select (planned)

`select` waits on multiple channel operations and runs the arm whose operation completes first:

```ebnf
selectExpr ::= "select" "{" selectArm+ "}"
selectArm  ::= IDENT "=>" channelOp "=>" expr
             | "_"   "=>" expr   (* timeout / default *)
```

```osprey
let ch1 = Channel<string> { capacity: 1 }
let ch2 = Channel<int>    { capacity: 1 }

select {
    msg => recv(ch1) => processString(msg)
    num => recv(ch2) => processNumber(num)
    _   => timeoutHandler()
}
```

## Fiber-Isolated Modules (planned)

Each fiber that touches a `module` receives its own private instance. There is no shared mutable state across fibers; communication is via channels.

```osprey
module Counter {
    mut count = 0
    fn increment() -> int = { count = count + 1; count }
    fn get()       -> int = count
}

let f1 = spawn Counter.increment()   // 1
let f2 = spawn Counter.increment()   // 1, not 2 — separate instance

await(f1)
await(f2)
```

A fiber's module instance is initialised on first access (copy-on-first-access) and is destroyed with the fiber.