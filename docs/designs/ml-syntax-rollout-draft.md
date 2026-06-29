# ML Syntax Rollout Draft

This is a non-authoritative design rollout draft. It is not a language spec and
not an implementation plan in `docs/plans/`. It exists to preserve the design
decisions and work breakdown so a later pass can promote the pieces into formal
spec changes and implementation plans deliberately.

Source design: [ML syntax prototype](ml-syntax-prototype.md).

## Target Language Shape

Osprey moves to one canonical layout syntax. Braces are not an alternate style
and should not remain as a permanent compatibility dialect. A migration tool can
exist, and a temporary compiler branch can accept old syntax during rollout, but
the target language has one form.

The core decisions:

- Layout/offside rule delimits blocks.
- `let` is removed from ordinary bindings.
- `name = expr` creates an immutable binding.
- `mut name = expr` creates a mutable binding.
- `name := expr` mutates an existing mutable binding.
- Same-scope rebinding with `=` is rejected.
- Function definitions use ML-style binding heads: `inc x = x + 1`.
- Function arrows are curried: `a -> b -> c` means `a -> (b -> c)`.
- Function application is whitespace application and left-associative.
- Tupled functions exist as one-argument functions over tuple values, but are not normal API style.
- Effect declarations are operation signatures, not trait-like function lists.
- Effect operations use `op : Payload => Result`, reserving `->` for functions.
- Handlers are first-class values: `handler Db ...` has type `Handler Db`.
- `handle db log persist do ...` installs handler values around a computation.
- Inline handler syntax is shorthand for anonymous handler values.
- Test doubles are ordinary handler values installed around the call under test.
- `do`, not `in`, separates handler installation from the handled computation.
- String interpolation keeps `${...}` for now.

## Syntax Inventory

### Bindings

```osp
answer = 42
name = "Osprey"

mut requests = 0
requests := requests + 1
```

TODO:

- [ ] Specify `=` as binding only, never mutation.
- [ ] Specify `:=` as mutation only, requiring a visible mutable binding.
- [ ] Specify same-scope rebinding diagnostics.
- [ ] Specify allowed shadowing in nested blocks and patterns.
- [ ] Specify top-level binding order and forward-reference behavior.
- [ ] Specify whether `_ = expr` is the canonical ignored binding.

### Functions And Currying

```osp
add : int -> int -> int
add x y = x + y

addOne : int -> int
addOne = add 1

answer = addOne 41
```

TODO:

- [ ] Specify `->` right-associativity.
- [ ] Specify function application left-associativity.
- [ ] Specify partial application typing.
- [ ] Specify fully saturated call optimization as an implementation permission.
- [ ] Specify tuple function types: `(a, b) -> c` means one tuple argument.
- [ ] Specify record-argument API style for named/config-style arguments.
- [ ] Decide explicit syntax for lambdas in the new layout syntax.
- [ ] Decide whether explicit effect annotations attach after the result type or elsewhere.
- [ ] Decide whether zero-argument functions are values, `Unit -> T`, or both at call sites.

### Effects

Effect declarations are named operation signatures:

```osp
effect Db
    add : string => int
    list : Unit => string
    count : Unit => int

effect Persist
    flush : string => int

effect Log
    info : string => Unit
```

This deliberately avoids:

```osp
add : string -> int
```

because `->` now means function/currying. An effect operation is a request with a
payload and result, not a method in a trait.

TODO:

- [ ] Replace effect operation `fn(...) -> ...` grammar with `name : Payload => Result`.
- [ ] Specify `Unit` payload operations.
- [ ] Specify that multi-field operation inputs should use records or named payload types.
- [ ] Specify operation result typing.
- [ ] Specify `perform Effect.operation payload` typing.
- [ ] Specify zero-payload perform style, likely `perform Db.list` or `perform Db.list ()`.
- [ ] Specify effect inference with the new operation syntax.
- [ ] Specify diagnostics for unknown effects and unknown operations.
- [ ] Specify diagnostics for operation payload/result mismatch.

### Handlers

Handlers are first-class values:

```osp
memoryDb : Unit -> Handler Db
memoryDb () =
    mut tasks = ""
    mut taskCount = 0

    handler Db
        add task =>
            taskCount := taskCount + 1
            tasks := "${tasks}#${toString taskCount} ${task}\n"
            taskCount

        list =>
            tasks

        count =>
            taskCount
```

Install them with `handle`:

```osp
db = memoryDb ()
log = silentLog ()

handle db log
do
    createTask "buy milk"
```

TODO:

- [ ] Add `Handler EffectName` type constructor.
- [ ] Specify `handler EffectName ...` expression typing.
- [ ] Specify handler operation coverage: all operations required, or partial handlers allowed only with explicit syntax.
- [ ] Specify handler arm parameter binding rules.
- [ ] Specify handler arm result checking against operation result type.
- [ ] Specify captured mutable state inside handler values.
- [ ] Specify handler identity and state sharing when a handler value is reused.
- [ ] Specify installing multiple handler values with one `handle`.
- [ ] Specify duplicate handlers for the same effect: reject by default unless explicit shadowing syntax is introduced.
- [ ] Specify handler lookup order for nested `handle` blocks.
- [ ] Specify handler interaction with C callback boundaries.
- [ ] Specify handler interaction with fibers.
- [ ] Specify handler interaction with future `resume`.
- [ ] Specify whether inline handler syntax remains:

```osp
handle
    handler Log
        info m => ()
do
    work ()
```

### Match And Patterns

```osp
match saved
    Success value => length snap
    Error message => -1
```

TODO:

- [ ] Specify layout match grammar.
- [ ] Specify single-field constructor patterns: `Success value`.
- [ ] Specify zero-field constructor patterns.
- [ ] Specify multi-field constructor patterns and record patterns.
- [ ] Specify wildcard `_`.
- [ ] Specify typed patterns in layout syntax.
- [ ] Specify exhaustiveness diagnostics after syntax migration.
- [ ] Specify unreachable/duplicate arm diagnostics after syntax migration.

### Records

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

TODO:

- [ ] Specify record construction as a layout block headed by constructor name.
- [ ] Specify record field assignment syntax.
- [ ] Specify record update syntax in layout form.
- [ ] Specify record pattern syntax in layout form.
- [ ] Specify record type declaration layout syntax.
- [ ] Specify inline record types if retained.

### Blocks

```osp
onPost body =
    id = perform Db.add body
    snap = perform Db.list
    written = perform Persist.flush snap
    textResp 201 "created\n"
```

TODO:

- [ ] Specify all layout regions that create lexical scopes.
- [ ] Specify final-expression return behavior.
- [ ] Specify empty block behavior.
- [ ] Specify statement separators via newline/layout.
- [ ] Specify expression continuation lines.
- [ ] Specify parenthesized grouping in layout syntax.

## Test Doubles As A Core Use Case

The language should make this test shape effortless:

```osp
response =
    handle fakeDb fakePersist fakeLog
    do
        createTask "buy milk"
```

The function under test performs operations:

```osp
createTask body =
    id = perform Db.add body
    snap = perform Db.list
    written = perform Persist.flush snap
    perform Log.info "created #${toString id}"
    textResp 201 "created task #${toString id}\n"
```

The test supplies test doubles:

```osp
test "createTask flushes and logs" =
    mut taskText = ""
    mut flushedSnapshot = ""
    mut logLine = ""

    db =
        handler Db
            add task =>
                taskText := task
                1

            list =>
                "#1 ${taskText}\n"

            count => 1

    persist =
        handler Persist
            flush snap =>
                flushedSnapshot := snap
                length snap

    log =
        handler Log
            info message =>
                logLine := message

    response =
        handle db persist log
        do
            createTask "buy milk"

    expectEqual 201 (httpResponseStatus response)
    expectEqual "#1 buy milk\n" flushedSnapshot
    expectEqual "created #1" logLine
```

TODO:

- [ ] Add examples showing inline spy handlers.
- [ ] Add examples showing reusable fake handler factories.
- [ ] Add examples showing failure-path handlers.
- [ ] Add examples showing handler values passed into helper functions.
- [ ] Add examples showing handler state isolation between tests.
- [ ] Add examples showing shared handler state when the same handler value is reused.
- [ ] Decide whether Osprey has a built-in `test` form or tests remain external examples.
- [ ] If a built-in `test` form is added, specify effect scoping per test.
- [ ] If a built-in `test` form is added, specify assertion functions and expected-output integration.

## Spec Promotion Checklist

When this draft becomes real specs, update these areas.

TODO:

- [ ] `docs/specs/0002-LexicalStructure.md`: layout/offside tokens, newline significance, indentation diagnostics.
- [ ] `docs/specs/0003-Syntax.md`: bindings, functions, calls, records, matches, handlers, effect declarations.
- [ ] `docs/specs/0004-TypeSystem.md`: curried functions, partial application, `Handler E`, operation payload/result typing.
- [ ] `docs/specs/0005-FunctionCalls.md`: whitespace application and currying.
- [ ] `docs/specs/0007-PatternMatching.md`: layout match and constructor payload patterns.
- [ ] `docs/specs/0008-BlockExpressions.md`: layout blocks and final-expression returns.
- [ ] `docs/specs/0017-AlgebraicEffects.md`: operation signatures, first-class handlers, handler installation, testing model.
- [ ] `docs/specs/0013-ErrorHandling.md`: Result pattern examples in new syntax.
- [ ] `docs/specs/0014-HTTP.md`: update examples to curried/layout syntax.
- [ ] Website generated spec mirrors after authoritative specs are updated.
- [ ] Examples README and stateful HTTP docs after compiler support lands.

## Implementation Breakdown

This is not a real implementation plan yet. These TODOs are the raw material for
later plans.

### Parser And AST

TODO:

- [ ] Define layout token model: `NEWLINE`, `INDENT`, `DEDENT`, or equivalent parser strategy.
- [ ] Parse binding heads: `name patterns = expr`.
- [ ] Parse mutation: `name := expr`.
- [ ] Parse layout function bodies.
- [ ] Parse whitespace application.
- [ ] Parse layout match.
- [ ] Parse layout record construction.
- [ ] Parse effect operation signatures with `=>`.
- [ ] Parse `handler Effect ...` expressions.
- [ ] Parse `handle handlerValue+ do expr`.
- [ ] Parse inline anonymous handler syntax if retained.
- [ ] Preserve source spans across layout tokens for diagnostics and LSP.
- [ ] Update AST types for `Handler` expressions and values.

### Tree-Sitter, VS Code, And LSP

TODO:

- [ ] Update `tree-sitter-osprey/grammar.js`.
- [ ] Regenerate parser artifacts.
- [ ] Update tree-sitter corpus tests.
- [ ] Update VS Code TextMate grammar.
- [ ] Update language configuration for indentation.
- [ ] Update folding queries.
- [ ] Update highlight queries for `handler`, `handle`, `do`, `:=`, `=>`.
- [ ] Update LSP parser integration.
- [ ] Update completion around effect operations and handler arms.
- [ ] Update hover/signature help for curried functions and handlers.

### Type Checker

TODO:

- [ ] Represent curried function types internally.
- [ ] Type-check partial application.
- [ ] Type-check saturated calls.
- [ ] Type-check tuple function arguments as ordinary tuple values.
- [ ] Reject same-scope rebinding.
- [ ] Type-check mutation only against `mut` bindings.
- [ ] Type-check effect operation signatures with payload/result.
- [ ] Type-check `perform` payloads and result types.
- [ ] Add `Handler E` type.
- [ ] Type-check handler operation coverage.
- [ ] Type-check handler arm parameters and result values.
- [ ] Type-check `handle` installation lists.
- [ ] Detect duplicate installed handlers.
- [ ] Preserve unhandled-effect checking with first-class handlers.
- [ ] Update diagnostics for old syntax if temporary migration parser exists.

### Codegen And Runtime

TODO:

- [ ] Lower curried functions to closures.
- [ ] Optimize fully saturated known-function calls.
- [ ] Support partial application values.
- [ ] Ensure captured values and mutable cells survive closure lowering.
- [ ] Represent `Handler E` values at runtime.
- [ ] Lower handler expressions to handler-value construction.
- [ ] Lower `handle values do body` to handler stack installation.
- [ ] Preserve handler-owned mutable state.
- [ ] Preserve handler behavior across C HTTP callback boundaries.
- [ ] Preserve handler behavior across fibers.
- [ ] Verify interaction with current effect runtime.
- [ ] Leave `resume` compatibility hooks for future continuation semantics.
- [ ] Update runtime tests for handler state isolation and sharing.

### Builtins, FFI, And Standard Library

TODO:

- [ ] Decide wrapper strategy for builtins under curried call syntax.
- [ ] Decide FFI strategy for curried Osprey functions over uncurried C ABI calls.
- [ ] Update builtin docs to new call style.
- [ ] Update HTTP examples to curried wrappers.
- [ ] Update file I/O examples to partial application where useful.
- [ ] Update collection functions to curried style friendly to `map`, `filter`, and `fold`.
- [ ] Audit argument order for partial-application ergonomics.

### Migration

TODO:

- [ ] Build a syntax migration tool or formatter mode from current syntax to layout syntax.
- [ ] Migrate `examples/tested`.
- [ ] Migrate `examples/statefulhttp`.
- [ ] Migrate `examples/websocketserver`.
- [ ] Migrate benchmark cases.
- [ ] Migrate website code examples.
- [ ] Migrate failure examples and expected diagnostics.
- [ ] Remove old syntax once migrated; do not keep permanent dual syntax.
- [ ] Add clear error messages for old braces/`let`/old effect operation syntax.

### Tests And Acceptance

TODO:

- [ ] Parser golden tests for layout indentation.
- [ ] Parser golden tests for error indentation.
- [ ] Type tests for currying and partial application.
- [ ] Type tests for same-scope rebinding rejection.
- [ ] Type tests for mutation diagnostics.
- [ ] Effect declaration tests for `Payload => Result`.
- [ ] Handler value tests.
- [ ] Handler state isolation tests.
- [ ] Handler state sharing tests.
- [ ] Test-double examples that compile and run.
- [ ] Stateful HTTP prototype migrated and runnable.
- [ ] LSP smoke tests on new syntax.
- [ ] Tree-sitter parse-all over migrated examples.
- [ ] Website docs build after examples migrate.

## Suggested Phasing

### Phase 0: Design Freeze

TODO:

- [ ] Confirm syntax decisions in `docs/designs/ml-syntax-prototype.md`.
- [ ] Confirm unresolved zero-argument call rule.
- [ ] Confirm lambda syntax.
- [ ] Confirm effect annotation syntax for function signatures.
- [ ] Promote design decisions into formal spec-change tasks.

### Phase 1: Parser Prototype

TODO:

- [ ] Implement layout parser in an experimental branch.
- [ ] Parse a small curated corpus before touching all examples.
- [ ] Parse the stateful HTTP prototype.
- [ ] Parse test-double examples.
- [ ] Update tree-sitter in parallel or explicitly defer editor support.

### Phase 2: Type System

TODO:

- [ ] Implement curried function types.
- [ ] Implement partial application.
- [ ] Implement handler value types.
- [ ] Implement new effect operation signature typing.
- [ ] Preserve unhandled-effect diagnostics.

### Phase 3: Runtime Semantics

TODO:

- [ ] Lower curried calls and partial applications.
- [ ] Lower handler values.
- [ ] Install multiple handler values with `handle ... do`.
- [ ] Verify handler state across HTTP callbacks and fibers.

### Phase 4: Migration And Examples

TODO:

- [ ] Migrate examples.
- [ ] Migrate docs.
- [ ] Migrate benchmarks.
- [ ] Add old-syntax diagnostics.
- [ ] Remove temporary old syntax support.

### Phase 5: Formalization

TODO:

- [ ] Split this draft into actual `docs/specs` edits.
- [ ] Split implementation work into actual `docs/plans` entries.
- [ ] Add links from the plans README after formal plans exist.
- [ ] Keep this draft as historical design context or delete it after promotion.
