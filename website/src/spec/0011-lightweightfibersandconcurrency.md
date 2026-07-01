---
layout: page
title: "Fibers and Concurrency"
description: "Osprey Language Specification: Fibers and Concurrency"
date: 2026-07-01
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0011-lightweightfibersandconcurrency/"
---

# Fibers and Concurrency

Fibers are lightweight concurrent computations. They are constructed as values of `Fiber<T>` and communicate through `Channel<T>`. There are no OS threads exposed to user code; the runtime schedules fibers cooperatively. Values cross fiber boundaries — `spawn` captures and channel `send` — by move or copy, never by sharing ([MEM-FIBER-ISOLATION] in [Memory Management](/spec/0018-memorymanagement/)).

> **Flavor layer — shared core (AST and above).**  Concurrency is a shared-core concern. The constructs here lower to canonical `osprey_ast` nodes — `Expr::Spawn`, `Expr::Yield`, `Expr::Await`, `Expr::Send`, `Expr::Recv`, and `Expr::Select` — and the runtime scheduler operates on those nodes alone ([FLAVOR-BOUNDARY] in [Language Flavors](/spec/0023-languageflavors/)). The semantics, the cooperative scheduling, and the channel runtime are one across every flavor; only the surface spelling differs. The Default (`.osp`) spelling is shown below; the ML (`.ospml`) counterpart is described in [ML Flavor Syntax](/spec/0024-mlflavorsyntax/). No phase below the AST can tell which flavor produced a fiber, send, or select.

## Status

`spawn`, `await`, `yield`, and basic channel operations are implemented. `yield`
performs a real cooperative hand-off: in concurrent (thread-backed) execution it
donates the CPU to the scheduler and resumes when next scheduled, forwarding its
operand unchanged. Under the deterministic execution mode used by the test
harness, fibers run sequentially to completion, so `yield` forwards its value
without re-ordering; true cross-fiber interleaving under deterministic mode would
require stackful context switching and is not yet implemented. The `select`
expression and the fiber-isolated module system below are planned and not yet
wired through code generation.

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

```osprey-ml
task = Fiber<int>
    computation = \() => calculatePrimes (n: 1000)
```

`spawn <expr>` is sugar for the equivalent `Fiber` construction:

```osprey
let result = spawn 42
// equivalent to:
let result = Fiber<int> { computation: fn() => 42 }
```

```osprey-ml
result = spawn 42
// equivalent to:
result = Fiber<int>
    computation = \() => 42
```

## Constructing Channels

```osprey
let sync   = Channel<int>    { capacity: 0  }   // unbuffered (rendezvous)
let buf    = Channel<string> { capacity: 10 }   // buffered
```

```osprey-ml
sync = Channel<int>    { capacity = 0  }   // unbuffered (rendezvous)
buf  = Channel<string> { capacity = 10 }   // buffered
```

## Operations

| Operation                                       | Signature                                            |
| ----------------------------------------------- | ---------------------------------------------------- |
| Wait for a fiber to produce its value           | `await(fiber: Fiber<T>) -> T`                        |
| Send a value to a channel                       | `send(channel: Channel<T>, value: T) -> Result<unit, ChannelError>` |
| Receive a value from a channel                  | `recv(channel: Channel<T>) -> Result<T, ChannelError>` |
| Yield to the scheduler, forwarding the value    | `yield(value: T) -> T`                               |

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

```osprey-ml
ch = Channel<int> { capacity = 3 }

producer = spawn
    range (1, 4) |> forEach (\i => send (ch, i))

consumer = spawn
    range (1, 4) |> forEach (\i =>
        match recv ch
            Success value   => print "got ${value}"
            Error   message => print "recv error: ${message}")

await producer
await consumer
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

```osprey-ml
ch1 = Channel<string> { capacity = 1 }
ch2 = Channel<int>    { capacity = 1 }

select
    msg => recv ch1 => processString msg
    num => recv ch2 => processNumber num
    _   => timeoutHandler ()
```

## Fiber-Isolated Modules (planned)

> **Superseded design note.** This section records the older sketch that existed
> before the multi-file module design. The normative module/state model is now
> [Modules and Namespaces](/spec/0025-modulesandnamespaces/), especially
> `[MODULES-STATE]` and `[MODULES-STATE-MODULE]`.

Each fiber that touches a `module` receives its own private instance. There is no shared mutable state across fibers; communication is via channels.

```osprey
module Counter {
    mut count = 0
    fn increment() = { count = count + 1; count }
    fn get()       = count
}

let f1 = spawn Counter.increment()   // 1
let f2 = spawn Counter.increment()   // 1, not 2 — separate instance

await(f1)
await(f2)
```

```osprey-ml
module Counter
    mut count = 0
    increment () =
        count := count + 1
        count
    get () = count

f1 = spawn Counter.increment ()   // 1
f2 = spawn Counter.increment ()   // 1, not 2 — separate instance

await f1
await f2
```

A fiber's module instance is initialised on first access (copy-on-first-access) and is destroyed with the fiber.