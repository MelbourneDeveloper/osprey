---
layout: page
title: "Language Flavors — Default vs ML"
description: "Osprey has one core and two surface syntaxes: the Default (.osp) C-style braces flavor and the ML (.ospml) offside-layout, curry-by-default flavor. Both lower to the exact same AST."
---

## One core, two flavors

Osprey is **one language** with **two surface syntaxes**. Every program — whichever
flavor you write it in — lowers to the *same* canonical AST, is checked by the same
Hindley-Milner type system, and runs identically. Nothing about effects, fibers,
types, or semantics changes between flavors; **only the way you spell it does.**

| | **Default** (`.osp`) | **ML** (`.ospml`) |
|---|---|---|
| Style | C-style braces | Offside (indentation) layout |
| Blocks | `{ … }` | Indentation |
| Function def | `fn f(a, b) = …` | `f (a, b) = …` (no `fn`) |
| Binding | `let x = v` | `x = v` (no `let`) |
| Reassignment | `x = x + 1` | `x := x + 1` (walrus) |
| Application | `f(x)` | `f x` (juxtaposition) |
| Lambda | `fn(x) => e` | `\x => e` |
| Match arm | `pat => body` in `{ }` | `pat => body`, indented |

Pick per file by extension. `.osp` is the Default flavor; `.ospml` is the ML flavor.
There is no runtime or performance difference — the choice is purely stylistic.

## Side by side

The same program in both flavors. They compile to identical output.

```osprey
type Shape = Circle | Square

fn area(s, size) = match s {
    Circle => size * size * 3
    Square => size * size
}

fn main() = {
    let total = area(Circle, 4) + area(Square, 2)
    print("total: ${total}")
}
```

```osprey-ml
type Shape = Circle | Square

area (s, size) =
    match s
        Circle => size * size * 3
        Square => size * size

main () =
    total = area (Circle, 4) + area (Square, 2)
    print "total: ${total}"
```

## The Default flavor (`.osp`)

Braces delimit blocks, `fn` introduces functions, `let` binds, and calls use
parentheses — familiar to anyone coming from C, Rust, Swift, or TypeScript.

```osprey
fn classify(n) = match n {
    0 => "zero"
    _ => "many"
}

fn main() = {
    let xs = [1, 2, 3]
    let doubled = xs |> map(fn(x) => x * 2)
    print("${classify(0)}, ${classify(9)}")
}
```

## The ML flavor (`.ospml`)

Indentation delimits blocks (the *offside rule*), functions are **curried by
default**, application is whitespace (`f x` instead of `f(x)`), lambdas are `\x => e`,
and mutation uses the walrus `:=`. Type annotations, when genuinely needed, sit on
their own line above the definition. It reads like OCaml, F#, or Haskell.

```osprey-ml
classify n =
    match n
        0 => "zero"
        _ => "many"

main () =
    xs = [1, 2, 3]
    doubled = xs |> map \x => x * 2
    print "${classify 0}, ${classify 9}"
```

## Things that never change

Both flavors are **fully identical** in everything that matters:

- **The type system** — Hindley-Milner inference, so you rarely annotate types in
  either flavor.
- **Algebraic effects** — declared with `effect` and discharged with `handle … in`;
  unhandled effects are a *compile error* in both flavors.
- **Pattern matching** — exhaustive `match`; a missing arm is rejected by the compiler.
- **Fibers, channels, and immutable persistent collections** — same runtime.
- **The pipe operator** `|>` — spelled identically in both.

Because both surfaces lower to one canonical AST, a Default file and an ML file can
describe the same program and produce byte-for-byte identical results. The `/wasm/`
demo ships exactly such a pair (`studio.osp` and `studio.ospml`).

For the full specification, see [Language Flavors](/spec/0023-languageflavors/) and
[ML Flavor Syntax](/spec/0024-mlflavorsyntax/).
