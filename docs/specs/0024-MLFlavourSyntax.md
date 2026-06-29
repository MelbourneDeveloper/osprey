# ML Flavour Syntax

The **ML flavour** is a layout-based source surface for Osprey: indentation
delimits blocks, functions curry by default, and effect handlers are first-class
values. It is one of Osprey's [language flavours](0023-LanguageFlavours.md) — a
parsing-and-lowering profile, not a separate language. Every construct here
lowers to the same `osprey_ast::Program` the Default (brace) flavour produces,
and from there shares one type checker, effect checker, and backend.

This chapter is the **surface reference**. The boundary rules, the lowering
contract, currying canonicalisation, and the shared-core handler-value feature
are normative in [Language Flavours](0023-LanguageFlavours.md); this chapter is
subordinate to that contract. Implementation is tracked in
[plan 0013](../plans/0013-ml-flavour-frontend.md).

- [Status](#status)
- [Layout Model](#layout-model)
- [Bindings and Mutation](#bindings-and-mutation)
- [Functions and Currying](#functions-and-currying)
- [Function Calls](#function-calls)
- [Effects](#effects)
- [Handlers](#handlers)
- [Match](#match)
- [Records](#records)
- [Blocks](#blocks)
- [Canonical Lowering Table](#canonical-lowering-table)
- [Worked Example](#worked-example)
- [Resolved Syntax Questions](#resolved-syntax-questions)

## Status

Not implemented. The Default flavour (specs `0001`–`0022`) is the only frontend
today. This chapter specifies the ML surface that
[plan 0013](../plans/0013-ml-flavour-frontend.md) will build as a second
frontend. Select it with `--flavour ml`, the `.ospml` extension, or a
`// osprey: flavour=ml` marker (see
[Flavour Selection](0023-LanguageFlavours.md#flavour-selection)).

## Layout Model

`[FLAVOUR-ML-LAYOUT]` The ML flavour uses the **offside rule**. A block is
introduced by a header line and continued by the lines indented under it; a line
indented less than the block's column closes it. Blocks nest by indentation.

```ebnf
INDENT  ::= (* start of a more-indented region *)
DEDENT  ::= (* return to a less-indented region *)
NEWLINE ::= (* significant end-of-line within a layout region *)
```

These tokens are produced by a stateful **external scanner** (the tree-sitter
brace grammar has none today — `tree-sitter-osprey/` ships no `scanner.c`). The
scanner tracks an indentation stack, emits `INDENT`/`DEDENT`/`NEWLINE`, and
ignores blank lines and comment-only lines. Source positions (row/column) are
preserved on every token so diagnostics and the LSP keep working
([FLAVOUR-LOWER-CONTRACT](0023-LanguageFlavours.md#the-lowering-contract)).

String interpolation keeps `${…}`. Parentheses remain available for grouping and
precedence; they are not mandatory call punctuation.

## Bindings and Mutation

`[FLAVOUR-ML-BIND]` `=` **binds**, `:=` **mutates**. There is no `let`: a bare
`name = expr` introduces an immutable binding in the current layout block. `mut`
marks a mutable binding, and every write to it uses `:=`, so mutation is visible
without scanning back to the declaration.

```ebnf
binding   ::= "mut"? bindingHead "=" expr
bindingHead ::= ID paramPattern*          (* zero patterns ⇒ value; one+ ⇒ function *)
mutation  ::= ID ":=" expr
```

```osp
answer = 42
mut requests = 0
requests := requests + 1
```

Same-scope rebinding with `=` is rejected; the diagnostic suggests `:=` if
mutation was meant. Shadowing in a nested block or pattern is allowed.

Lowering: `name = e` → `Stmt::Let { mutable: false }`; `mut name = e` →
`Stmt::Let { mutable: true }`; `name := e` → `Stmt::Assignment`. These are the
*same* canonical nodes the Default flavour emits for `let`, `mut`, and `=`
reassignment respectively — only the spelling differs.

## Functions and Currying

`[FLAVOUR-ML-FN]` A function definition is a binding whose head has one or more
parameter patterns. The optional signature line above it uses ML arrows.

```ebnf
signature ::= ID ":" type
funDef    ::= ID paramPattern+ "=" blockOrExpr
type      ::= type "->" type            (* right-associative: a -> b -> c = a -> (b -> c) *)
            | "(" type ("," type)* ")" "->" type   (* one tupled argument *)
            | typeAtom
```

```osp
inc : int -> int
inc x = x + 1

add : int -> int -> int
add x y = x + y
```

`[FLAVOUR-ML-CURRY]` Arrows are **right-associative** and application is
**left-associative**, so `int -> int -> int` is `int -> (int -> int)` and
`add 1 2` is `(add 1) 2`. Multi-argument function syntax therefore *reads as
curried*, and partial application falls out for free:

```osp
addOne : int -> int
addOne = add 1

answer = addOne 41
```

A tupled function takes **one** tuple value and is not the normal API style:

```osp
distance : (int, int) -> int
distance point = ...

d = distance (3, 4)
```

Lowering (normative in
[FLAVOUR-CURRY](0023-LanguageFlavours.md#currying-canonicalisation)): `add x y =
body` lowers to a one-parameter binding returning a one-parameter `Expr::Lambda`
— **identical** canonical AST to the Default flavour's explicit
`fn add(x) -> (int) -> int = fn(y) => body`. It is deliberately **not** the same
AST as Default `fn add(x, y) = body` (one two-parameter `Stmt::Function`).
`add 1 2` lowers to nested single-argument `Expr::Call`s; no partial-application
support is added to the type checker, because every curried function and every
application is one-argument and thus always saturated.

API guidance: put stable, configuration-like arguments first and the data
argument last, so partial application is useful (`replace " " ""` ⇒ a
space-remover).

## Function Calls

`[FLAVOUR-ML-CALL]` Calls use whitespace application; parentheses group.

```ebnf
application ::= app atom
             | atom
atom        ::= ID | literal | "(" expr ")"
```

```osp
length snap
textResp 201 "created\n"
c256 "213" (blocks 0 (mn n 28))
```

Lowering: `f a b` → nested `Expr::Call`, one argument each
(`Call(Call(f,[a]),[b])`). A parenthesised group lowers to its inner expression.

## Effects

`[FLAVOUR-ML-EFFECT]` An effect declaration is a layout block of operation
signatures. Operations use `=>` so that `->` keeps its one meaning — function
and currying type. An operation is a request with a **payload** and a **result**,
not a curried function.

```ebnf
effectDecl ::= "effect" ID INDENT opSig+ DEDENT
opSig      ::= ID ":" type "=>" type
```

```osp
effect Db
    add : string => int
    list : Unit => string
    count : Unit => int

effect Log
    info : string => Unit
```

Zero-payload operations take `Unit`. Multi-field requests use a record payload,
not a fake multi-argument operation:

```osp
type AddTask =
    body : string
    priority : int

effect Db
    add : AddTask => int
```

Lowering: `effect E` + arms → `Stmt::Effect { operations }`, where each
`op : P => R` becomes `EffectOperation { name, parameters: [P], return_type: R }`
— the same canonical node the Default `op : fn(P) -> R` produces.
`perform E.op a` → `Expr::Perform`.

> `->` belongs to functions and currying. `=>` belongs to clauses and requests
> that yield a result: it appears in `effect` operations, `handler` arms, and
> `match` arms, always meaning "the left yields the right."

## Handlers

`[FLAVOUR-ML-HANDLER]` Handlers are **first-class values**. `handler E` followed
by indented arms evaluates to a value of type `Handler E`. `handle` installs one
or more such values around a computation, with `do` marking the handled body.

```ebnf
handlerValue ::= "handler" ID INDENT handlerArm+ DEDENT
handlerArm   ::= ID param* "=>" blockOrExpr
install      ::= "handle" expr+ "do" blockOrExpr
```

```osp
memoryDb : Unit -> Handler Db
memoryDb () =
    mut tasks = ""
    mut taskCount = 0

    handler Db
        add t =>
            taskCount := taskCount + 1
            tasks := "${tasks}#${toString taskCount} ${t}\n"
            taskCount

        list =>
            tasks

        count =>
            taskCount
```

Installing several at once replaces the Default flavour's repeated nesting:

```osp
db = memoryDb ()
log = silentLog ()

handle db log
do
    createTask "buy milk"
```

The mutable cells belong to the handler value: a fresh `handler` makes fresh
state; passing the same value around shares it. Parameterised handlers compose
with currying (`filePersist path = … handler Persist …`).

First-class handler values, the `Handler E` type, and multi-install are a
**shared-core feature**, not ML-only sugar — see
[FLAVOUR-HANDLER-VALUE](0023-LanguageFlavours.md#shared-core-additions). They
lower to `Expr::HandlerValue` and `Expr::Install`; `handle a b c do body`
desugars to nested installs. The Default flavour gains the same feature in brace
spelling.

## Match

`[FLAVOUR-ML-MATCH]` `match` uses the same clause style as handlers: the
scrutinee follows `match`, and each indented arm is `Pattern => body`. A
one-payload constructor binds its payload directly — `Success value`, not
`Success { value }`.

```ebnf
matchExpr ::= "match" expr INDENT matchArm+ DEDENT
matchArm  ::= pattern "=>" blockOrExpr
```

```osp
diskBytes =
    match saved
        Success value => length snap
        Error message => -1
```

Lowering: `Expr::Match` + `MatchArm`; `Success value` →
`Pattern::Constructor { name: "Success", fields: ["value"] }` — the same node the
Default `Success { value }` produces. Wildcard `_` → `Pattern::Wildcard`.

## Records

`[FLAVOUR-ML-RECORD]` Record construction is a layout block headed by the
constructor name, with `field = value` lines. Inside a record literal the left
of `=` is a field name, not a new binding; the indentation under a constructor
makes that unambiguous.

```ebnf
recordExpr ::= ID INDENT fieldInit+ DEDENT
fieldInit  ::= ID "=" expr
```

```osp
textResp status bodyText =
    HttpResponse
        status = status
        headers = "Content-Type: text/plain"
        contentType = "text/plain"
        streamFd = -1
        isComplete = true
        partialBody = bodyText
```

Lowering: `Expr::TypeConstructor { name, fields }`; record update lowers to
`Expr::Update`.

## Blocks

`[FLAVOUR-ML-BLOCK]` A function body, match arm, handler arm, or `do` body is an
ordinary layout region containing bindings, mutations, performs, and a final
expression. The final expression is the block's value. There is no separate
`{ … }` expression form in this flavour.

```osp
onPost body =
    id = perform Db.add body
    snap = perform Db.list
    written = perform Persist.flush snap
    perform Log.info "created"
    textResp 201 "created\n"
```

Lowering: `Expr::Block { statements, value }`, where `value` is the trailing
expression — the same node the Default `{ … }` block produces.

## Canonical Lowering Table

Every ML form on the left lowers to the canonical node on the right
(`crates/osprey-ast/src/lib.rs`). The Default-flavour spelling of the same node
is in [FLAVOUR-LAYER](0023-LanguageFlavours.md#flavour-concern-vs-shared-core-concern).

| ML surface | Canonical AST node |
| --- | --- |
| `x = e` | `Stmt::Let { mutable: false }` |
| `mut x = e` | `Stmt::Let { mutable: true }` |
| `x := e` | `Stmt::Assignment` |
| `f x y = e` (curried) | one-param `Stmt::Function` → nested `Expr::Lambda` |
| `\y => e` | `Expr::Lambda` |
| `f a b` | nested one-arg `Expr::Call` |
| layout block | `Expr::Block` |
| `match v` + arms | `Expr::Match` + `MatchArm` |
| `Success value` | `Pattern::Constructor { fields: ["value"] }` |
| `T` + `f = v` lines | `Expr::TypeConstructor` |
| `effect E` + `op : P => R` | `Stmt::Effect` + `EffectOperation` |
| `perform E.op a` | `Expr::Perform` |
| `handler E` + arms | `Expr::HandlerValue` *(shared-core addition)* |
| `handle a b do body` | `Expr::Install` *(shared-core addition)* |

## Worked Example

The same program a Default-flavour author would write with braces, `fn`, named
arguments, and nested `handle … in`. It exercises curried definitions, partial
application (`textResp 201`, `c256 "213"`), `=>` effect operations, first-class
handler values with owned `mut` state, and one grouped `handle … do`.

```osp
effect Db
    add : string => int
    list : Unit => string
    count : Unit => int

effect Log
    info : string => Unit

c256 : string -> string -> string
c256 n s =
    "\e[38;5;${n}m${s}\e[0m"

rose : string -> string
rose = c256 "213"

textResp : int -> string -> HttpResponse
textResp status bodyText =
    HttpResponse
        status = status
        headers = "Content-Type: text/plain"
        contentType = "text/plain"
        streamFd = -1
        isComplete = true
        partialBody = bodyText

memoryDb : Unit -> Handler Db
memoryDb () =
    mut tasks = ""
    mut taskCount = 0

    handler Db
        add t =>
            taskCount := taskCount + 1
            tasks := "${tasks}#${toString taskCount} ${t}\n"
            taskCount

        list => tasks
        count => taskCount

silentLog : Unit -> Handler Log
silentLog () =
    handler Log
        info m => ()

createTask : string -> HttpResponse
createTask body =
    id = perform Db.add body
    snap = perform Db.list
    perform Log.info "created #${toString id} ${snap}"
    textResp 201 "created task #${toString id}\n"

main : Unit -> int
main () =
    db = memoryDb ()
    log = silentLog ()

    handle db log
    do
        response = createTask "buy milk"
        print (httpResponseBody response)
        0
```

The first-class handlers make test doubles trivial — a test installs spy or
stub handlers that close over the test's own `mut` cells around the call under
test, with no `Db`/`Log` parameters polluting the production signature:

```osp
test "createTask stores the task and logs" =
    mut stored = ""
    mut logLine = ""

    db =
        handler Db
            add task =>
                stored := task
                1
            list => "#1 ${stored}\n"
            count => 1

    log =
        handler Log
            info message => logLine := message

    response =
        handle db log
        do
            createTask "buy milk"

    expectEqual 201 (httpResponseStatus response)
    expectEqual "buy milk" stored
    expectEqual "created #1 #1 buy milk\n" logLine
```

## Resolved Syntax Questions

- **Zero-argument functions:** a parameterless `name = expr` is a value binding;
  a `name () = expr` is a `Unit -> T` function. Pure constants are values
  (`banner`); `()` is used where recursion or effects make the call boundary
  meaningful (`serveForever ()`).
- **Lambdas:** anonymous functions are written `\param* => body` (lowering to
  `Expr::Lambda`), keeping `=>` as the clause/yield arrow and `->` as the type
  arrow.
- **Effect annotations on signatures:** the effect row follows the result type,
  as in the Default flavour (`saveTask : string -> int ![Store, Log]`).

## Cross-references

- [Language Flavours](0023-LanguageFlavours.md) — the normative boundary,
  contract, currying canonicalisation, and shared-core handler-value feature.
- [Algebraic Effects](0017-AlgebraicEffects.md) — effect semantics shared by both
  flavours.
- [Plan 0013 — ML Flavour Frontend](../plans/0013-ml-flavour-frontend.md).
