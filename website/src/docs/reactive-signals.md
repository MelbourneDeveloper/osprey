---
mlTwins: manual
layout: page
title: Reactive signals
description: Derive signal dependencies, cache immutable results, and rebuild affected computations.
permalink: /docs/reactive-signals/
---

# Rebuild only the views that depend on a change

Declare signal reads as static effects. A handler supplies their current values; the compiler derives each computation's possible signal reads through helpers and callbacks. You do not maintain a dependency array.

```osprey
static effect Signal<T> { read: fn() -> T }
fn count() = perform Signal<int>.read()
fn label() = "count=${count()}"
let first = {
    handle static Signal<int> { read => 1 }
    signalSnapshot(label)
}
let updated = {
    handle static Signal<int> { read => 2 }
    signalRefresh(first, signalDependencies(count), label)
}
print("${first.value} / ${updated.value}")
```

```ospml
static effect Signal T
    read : Unit => T
count () = perform Signal<int>.read ()
label () = "count=${count ()}"
first =
    handle static Signal<int>
        read => 1
    signalSnapshot label
updated =
    handle static Signal<int>
        read => 2
    signalRefresh first (signalDependencies count) label
print "${first.value} / ${updated.value}"
```

This prints `count=1 / count=2`. `first` retains its original value. Each snapshot contains `value`, `dependencies`, and `revision`; the revision starts at 1 and increases when the computation rebuilds.

`signalDependencies(work)` returns sorted, distinct identities such as `Signal<int>.read`. It evaluates the callback argument once without invoking the callback. `signalSnapshot(work)` invokes it once and caches its result. `signalRefresh(cached, changed, work)` invokes it only when a dependency is in `changed`, the computation's dependency set changes, or either dependency information or changes include `<unknown>`. Otherwise it returns the existing snapshot. Arguments evaluate once, left to right, including when the cached value is reused. Cached callable values keep their own requirements. A callable that escapes a statically selected handler while still requiring its operations is rejected before erasure; the compiler must not bind its future calls to the construction handler.

Supply every changed signal identity together in one refresh to batch an update. Signal identity includes the effect, operation and resolved generic arguments. Two independently changeable inputs need different identities; using the same payload type twice does not create two signals. The caller supplies the current values and the changed identities; refresh runs synchronously, and a later update builds on the returned snapshot.

Dependencies are conservative **may-read** sets. Both arms of a runtime choice contribute possible reads, including an arm not taken on this execution. A resolved callback keeps its actual requirements. An unresolved callback, generic instance or open remainder contributes `<unknown>`, which forces refresh rather than risk retaining a stale value. Builtin I/O is still checked as an effect obligation; metadata queries do not make their argument expressions pure.

Run the [complete dashboard examples](https://github.com/Nimblesite/osprey/tree/main/examples/handlers), in Default and ML, with `python3 examples/handlers/run.py all --demo staged-rows --check`. They update counter and cursor views independently, rebuild a combined status view for either change, and retain all three on an idle update. Native programs can also rebuild computations that use fibers and deep one-shot continuations.

The [effects specification](/spec/0017-algebraiceffects/#reactive-signals--stage-signals) defines the contract.
