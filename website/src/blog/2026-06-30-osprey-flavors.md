---
layout: page.njk
title: "Osprey Flavors: One Core, Two Surfaces, Zero Compromise"
excerpt: "Braces or layout — pick your tribe and go all in. The ML flavor isn't braces-optional and the Default flavor isn't deprecated. It's the same language underneath."
date: 2026-06-30
tags: ["blog", "language-design", "flavors", "functional-programming", "ml-syntax"]
author: "Christian Findlay"
readingTime: 7
image: /assets/images/og-image.png
---

Every language picks a side, and every side loses someone. Curly braces or significant whitespace. `fn add(x, y)` or `add x y`. The systems programmer who wants explicit blocks and named arguments, versus the FP devotee who wants layout and curry-by-default. Pick braces and the Haskell crowd wrinkles their nose; pick layout and the C crowd walks away. The syntax wars are real, and they force you into a tribe before you've written a line.

Osprey's answer is to stop pretending there's one right answer. **One core. Two surfaces. Zero compromise.**

## The problem: syntax forces a tribe

The FP-snob-versus-systems-programmer divide is mostly about spelling. The ideas — algebraic data types, exhaustive matching, immutability, effects — are not in dispute. What's in dispute is whether `do`-blocks should have braces, whether application should need parentheses, whether a function with two arguments is one value or two. These are aesthetic and ergonomic preferences, and they are *strong* preferences. Telling someone their preferred surface is wrong is how you lose them.

Most languages resolve this by declaring a winner and grudgingly bolting on the loser as an afterthought — a "lite" mode, an optional layout extension, a deprecated legacy syntax kept alive for migration. The afterthought is always second-class, and everybody can tell.

## Osprey's answer: flavors

Osprey ships **two first-class, permanent syntaxes** called flavors. Neither is the watered-down one.

- **Default flavor (`.osp`)** — C-style braces, `fn`, `f(x: a, y: b)` calls with named arguments. Explicit, familiar, block-structured. This is the surface a systems programmer reaches for. Fully implemented today.
- **ML flavor (`.ospml`)** — offside-rule layout (indentation, no braces), curry-by-default, whitespace application `f a b`, `\x => e` lambdas, `:=` mutation, `->` for types and `=>` for clauses. Terse, expression-first, ML/Haskell-shaped. This is the surface an FP devotee reaches for. In active development.

The point is **no compromise**. The ML flavor is not "braces optional." The Default flavor is not deprecated or transitional. Each surface goes all the way in its own direction. Systems programmers get real braces and real named arguments; FP folks get real layout and real currying. Nobody is asked to swallow the other camp's spelling. The language belongs to *your* tribe — pick your flavor and go all in.

Here's the ML flavor saying hello — this runs today:

```osprey
greeting = "Hello from the ML flavor"
print greeting
print "2 + 3 = ${2 + 3}"
```

No `fn`, no braces, no parentheses around the print argument. Layout and whitespace application all the way down.

## How it actually works

A flavor is not a preprocessor or a transpiler bolted onto a host language. Each flavor is a **parser plus a lowerer** that converge on **one canonical AST** — `osprey_ast::Program`. After lowering, there is exactly one type checker, one effect system, one optimiser, one LLVM/wasm backend. None of them know which flavor you wrote. The flavor is gone by the time any analysis runs.

That's what makes the "no compromise" claim more than a slogan: both surfaces meet at the same tree, so both get the same Hindley-Milner inference, the same compile-time effect safety, the same performance. There is no second-class path.

The one honest difference between the surfaces is currying, and it's machine-checked. In ML, every function is curried by default:

```osprey
inc : int -> int
inc x = x + 1

add : int -> int -> int
add x y = x + y

// partial application falls straight out of currying:
addTen = add 10
answer = addTen 32        // 42
```

That ML `add x y` lowers to *exactly* the same canonical AST as this Default-flavor explicit-curry definition:

```osprey
// Default flavor (.osp):
fn add(x) = fn(y) => x + y
// ML flavor (.ospml) — identical canonical AST:
add x y = x + y
```

We have a test that asserts the two produce the same tree. Note the precision here: ML `add x y` equals the *explicit-curry* Default form, not the multi-parameter `fn add(x, y)`. The latter is deliberately a different value — a single two-argument function, not a chain. The flavors converge where they should and stay distinct where the semantics genuinely differ.

The ML surface carries its FP shape all the way through. Layout-driven `match`:

```osprey
classify n =
    match n
        0 => "zero"
        1 => "one"
        _ => "many"
```

Higher-order functions and `Result` payload matching (integer division and mod return `Result<int, MathError>`, so you match the payload):

```osprey
twice : (int -> int) -> int -> int
twice f x = f (f x)

bump x = x + 10

safeMod a b =
    match a % b
        Success value => value
        Error e => -1
```

Bindings and mutation, with `=` to bind and `:=` to mutate:

```osprey
mut counter = 0
counter := counter + 1      // := mutates; = binds
print "counter = ${counter}"
```

## Same folder, compiled together

Because every flavor lowers to the same canonical AST *before* any type checking, the flavor is a **per-file** choice — not a per-project one. A `.osp` file and a `.ospml` file can sit in the same folder and compile into one program:

```osprey
// One project folder, two flavors, one compiled program:
//   project/
//     math.ospml     # ML flavor — curry-by-default module
//     app.osp        # Default flavor — braces; imports math
// Each file is wholly one flavor (chosen by extension/marker/--flavor). Both lower to
// the SAME canonical AST, so they share one type checker and one binary. Exports are
// canonical signatures, so a Default module and an ML module import each other normally.
```

Exports are canonical signatures with stable names and ordering, so a Default module and an ML module reference each other with no glue layer. The team is never forced to pick one tribe; each developer picks the flavor for the file they're writing.

To be precise about what ships today: **per-file flavor selection is implemented and green.** You select the ML surface with the `.ospml` extension, the `--flavor ml` CLI flag, or a leading `// osprey: flavor=ml` marker (precedence: flag > marker > extension > Default). That mechanism is exercised by tested examples right now. The multi-file *cross-flavor import* — a Default module pulling in an ML module in the same build — is the design direction the canonical-AST architecture is built for, but it is not yet covered by a tested example, so we're showing you the folder model and the per-file selection that is green, not a runnable cross-flavor import program.

## Effects: shown in the Default flavor

Osprey's headline feature is compile-time-safe algebraic effects. Here's the effect demo in the **Default flavor**, where it's fully implemented:

```osprey
effect Logger {
    log: fn(string) -> Unit
}

fn greet(name: string) -> Unit !Logger =
  perform Logger.log("Hello, ${name}!")

// Production: write to stdout
handle Logger
  log msg => print(msg)
in greet("Alice")

// Test: stay silent — same code, new handler
handle Logger
  log msg => 0
in greet("Bob")
```

The `!Logger` row says `greet` performs a `Logger` effect; an unhandled effect is a compile error. Swap the handler and the same code logs to stdout or stays silent — no global mutable wiring, just a different `handle` block.

## Status, honestly

The **Default flavor is fully implemented** — specs 0001 through 0022, the complete effect system, the persistent collections, the lot.

The **ML flavor is in active development**, with runnable proof you can read and run: the [tested ML examples](https://github.com/Nimblesite/osprey/tree/main/examples/tested/ml) cover hello-world, curry-by-default with partial application, higher-order functions, `Result` matching, layout `match`, and mutation — each one runs through the compiler and its `stdout` is byte-compared against a checked-in `.expectedoutput`.

The honest caveat: **ML effect syntax does not run yet.** The ML surface for `effect` / `handle … do` is *designed* but errors loudly today — it's the deferred Phase 0 shared-core feature. The next thing landing is **first-class handler values plus ML effects**, the piece that lets the effect system reach the ML surface the same way everything else already has. Until then, every effect demo you see in our docs is in the Default flavor, because that's what works.

## Pick your flavor

If you live in braces and named arguments, write `.osp` and never think about layout again. If you live in layout and currying, write `.ospml` and never type a brace. Either way you get the same Hindley-Milner type checker, the same compile-time effect safety, the same backend, the same standard library — because after lowering, nothing downstream can even tell which flavor you wrote.

Pick your flavor. Go all in. It's the same Osprey.
