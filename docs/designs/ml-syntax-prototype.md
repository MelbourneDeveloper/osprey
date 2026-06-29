# Osprey ML Syntax Prototype

This is a syntax critique and prototype, not current compiler input.

The current API style is close, but it has two languages fighting each other:

- Normal functions, records, matches, and blocks look C/Rust-like.
- Effect handlers already look closer to ML: layout-oriented, clause-based, and expression-first.
- `{ ... }` is overloaded as block syntax, record construction, record patterns, match bodies, effect bodies, and type bodies.
- `handle Persist ... in handle Metrics ... in handle Db ...` is semantically fine but visually noisy. The repeated `in handle` chain is ceremony, not information.
- `Success { value } => length(snap)` exposes the representation of a one-field variant when the programmer only wants to bind the payload.
- `let` is useful in a C-ish grammar, but in a layout grammar it becomes low-value noise. The important distinction is not "let versus expression"; it is "binding versus mutation".

The proposed direction is to make Osprey consistently layout-based. This is not
"braces are optional" and it is not "the formatter can choose a nicer style".
It is a single syntax philosophy: the language has one canonical visual form,
and that form is layout.

Braces should not remain as an alternate spelling in the language design. A
one-time migration tool can translate old examples, and a transitional compiler
branch can exist while the change lands, but the target language should not have
two permanent block syntaxes. Two syntaxes means every example, formatter,
parser diagnostic, and style discussion has to pick sides forever.

- Indentation delimits blocks.
- `name = expr` creates an immutable binding.
- `mut name = expr` creates a mutable binding.
- `name := expr` mutates an existing mutable binding.
- `fn` is dropped from normal definitions.
- Function definitions and handler clauses use the same visual shape.
- `match` and `handle` use layout instead of braces.
- Single-field constructor patterns use `Constructor value`, not `Constructor { value }`.
- Record construction uses layout.
- Function calls may use ML-style whitespace application.
- A grouped `handle` uses one `do` to introduce the handled computation.

## Core Forms

### Values And Mutation

Dropping `let` only works if `=` stops meaning "assignment". In this design,
`=` is a binding operator. It introduces a name in the current layout block.
Mutation is deliberately spelled differently with `:=`, so a reader can tell
the difference without knowing whether a variable was declared earlier.

This gives the language less ceremony without making state invisible. Mutable
state still announces itself twice: once at declaration with `mut`, and again at
each write with `:=`.

```osp
answer = 42
label = "tasks"

mut requests = 0
requests := requests + 1
```

`=` binds. `:=` mutates. This lets `let` disappear without making reassignment look identical to local binding.

Same-scope rebinding should be rejected. Shadowing in a nested pattern or nested
block is still useful, but accidentally writing a second `name = ...` in the
same block should get a diagnostic that says to use `:=` if mutation was meant.

### Functions

Function definitions keep the current ML-like shape. Dropping `let` does not
remove `inc x = x + 1`; it makes that form the normal function-binding syntax.
The left side of a statement-level `=` is a binding head: a name plus zero or
more parameter patterns.

This is the rule that prevents confusion between bindings and returned values:
a binding line has a top-level `=`, while a returned expression does not. The
final expression in a layout block is still the value of the block.

```osp
esc : string -> string -> string
esc code s =
    "\e[${code}m${s}\e[0m"

inc : int -> int
inc x = x + 1
```

The signature syntax is ML-like, and this prototype adopts the ML meaning:
`a -> b -> c` is right-associative and means `a -> (b -> c)`. A function takes
one argument and returns either the final result or another function waiting for
the next argument.

This sounds strange from C-family languages, but it makes partial application
fall out naturally. If a function normally needs two arguments, supplying one
argument gives back a smaller function.

```osp
main : Unit -> int
main () =
    x = 41
    inc x
```

Here `x = 41` is a local binding. `inc x` has no top-level `=`, so it is the
final expression returned from `main`.

The parser can keep this simple:

```osp
name = expr             // value binding
name arg1 arg2 = expr   // function binding
name := expr            // mutation
expr                    // expression; if final in block, returned
```

### Currying

Currying is settled for this syntax. Osprey should not borrow the ML arrow shape
and then give it non-ML semantics. If the type says `int -> int -> int`, the
function takes an `int` and returns a function from `int` to `int`.

```osp
add : int -> int -> int
add x y = x + y

answer = add 1 2
```

Function application is left-associative, so this:

```osp
add 1 2
```

means this:

```osp
(add 1) 2
```

The payoff is that reusable adapters do not need lambda noise:

```osp
addOne : int -> int
addOne = add 1

okText : string -> HttpResponse
okText = textResp 200

writeTasks : string -> Result<Unit, Error>
writeTasks = writeFile "/tmp/osprey_tasks.db"
```

This should shape API design. Put stable, reusable, configuration-like arguments
first, and the final data argument last. That makes partial application useful:

```osp
replace : string -> string -> string -> string
replace old new text = ...

removeSpaces = replace " " ""
slug = removeSpaces "hello world"
```

Uncurried tuple functions still exist, but they are not the normal API style.
They mean "this function takes one tuple value":

```osp
distance : (int, int) -> int
distance point = ...

distance (3, 4)
```

### Effects

Effect declarations are named operation signatures. They look a bit like
Haskell type classes or Rust traits because they are a named bundle of
operations, but the semantics are different. A trait says "this type implements
these methods". An algebraic effect says "this computation may request these
operations".

That distinction matters because effects are handled by computation scope, not
implemented by a receiver type. A function that performs `Db.add` does not know
which database implementation it is using. The surrounding handler decides what
that operation means for this run.

In a layout language, the braces do not add information: the indented operation
signatures already say exactly what belongs to the effect.

Moving effects to layout also makes them visually match handlers. The effect
declares operations with names and types; a handler later interprets those same
operation names with clauses.

```osp
effect Db
    add : string -> int
    list : Unit -> string
    count : Unit -> int

effect Log
    info : string -> Unit
```

No effect-body braces. Operations read like a small scoped signature.

### Match

`match` should use the same clause style as handlers. The discriminant follows
`match`, and each indented arm is a pattern plus `=>`. This removes the match
body braces and makes nested matches easier to scan because indentation carries
the tree shape.

The common `Result` case should not expose one-field record destructuring.
`Success { value }` is precise, but noisy. If a constructor has one payload,
`Success value` should bind that payload directly.

```osp
diskBytes =
    match saved
        Success value => length snap
        Error message => -1
```

Single-field variants bind directly. Multi-field constructors can still have a named form later, but the common `Result` case should not need record-pattern braces.

### Records

Record construction is one of the worst brace clusters in the current syntax,
especially in API code where response objects are common. A layout record keeps
the constructor name as the header and places fields under it.

The field form deliberately reuses `field = value`. Inside a record literal the
left side is a field name, not a new local binding. That context is unambiguous
because the fields are indented under a constructor.

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

Record construction is a layout block headed by the constructor name. This removes the heaviest brace cluster in API code.

### Blocks

Blocks become ordinary layout regions. A function body, match arm, handler arm,
or `do` body can contain bindings, mutations, effect performs, and a final
expression. There is no separate `{ ... }` expression form in the target syntax.

This is the biggest aesthetic change, but it is what makes the language feel
cohesive. If every multi-line expression follows the same indentation rule,
there is less punctuation to remember and fewer local style choices.

A layout block returns its final expression:

```osp
onPost body =
    id = perform Db.add body
    snap = perform Db.list
    written = perform Persist.flush snap
    perform Log.info "created"
    textResp 201 "created\n"
```

This keeps Osprey's existing block-expression behavior, but makes the block the default shape rather than a `{ ... }` island.

### Function Calls

Calls may use whitespace application:

```osp
length snap
textResp 201 "created\n"
httpCreateServer 8080 "127.0.0.1"
```

Parentheses remain for grouping and precedence, not as mandatory call
punctuation. That keeps complex expressions readable without making every
simple call look like C.

```osp
gauge id
c256 "213" (blocks 0 (mn n 28))
```

Named arguments can still exist where the language needs explicit labels, but
the surface style of normal Osprey should bias toward the lighter ML call form.

### Handlers

Handlers are the syntax that already points in the right direction. The current
problem is not the handler arms themselves; it is the repeated nesting:
`handle Persist ... in handle Metrics ... in handle Db ... in ...`. The
important thing is the set of handlers installed for a computation, not the
mechanical nesting needed by the current AST.

The core split should be:

```osp
effect Db       // declares operations
handler Db      // creates a handler value that interprets those operations
handle db do    // installs handler values around a computation
```

Handlers need to be first-class values. That lets code build them, configure
them, pass them into functions, return them from functions, and swap them for
tests. This is the effect-system version of dependency injection, without
turning effects into traits or interfaces.

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

The mutable cells belong to the handler value. Creating a new handler creates
fresh state; passing the same handler around shares that interpreter and its
state.

```osp
db = memoryDb ()

handle db
do
    handleReq method path headers body
```

Parameterized handlers compose naturally with currying:

```osp
filePersist : string -> Handler Persist
filePersist path =
    writeSnapshot = writeFile path

    handler Persist
        flush snap =>
            saved = writeSnapshot snap
            match saved
                Success _ => length snap
                Error _ => -1

persist = filePersist "/tmp/osprey_tasks.db"
```

A grouped `handle` installs several handler values for one body. `do` marks the
start of that body. It reads better than `in` for non-ML users: "handle these
effects, then do this computation." The word is short, familiar, and still
creates an explicit boundary between handler installation and handled code.

```osp
persist = filePersist "/tmp/osprey_tasks.db"
metrics = metricsCounter ()
db = memoryDb ()
log = fileLog "/tmp/osprey_server.log"

handle persist metrics db log
do
    serverId = httpCreateServer 8080 "127.0.0.1"
    listening = httpListen serverId handleReq
    serveForever ()
```

Inline handler syntax can still exist for local one-off handlers:

```osp
handle
    handler Metrics
        hit p => requests := requests + 1
        served => requests
do
    serveForever ()
```

That is shorthand for creating an anonymous `Handler Metrics` value and
installing it immediately. The underlying concept is still handler values.

One `handle` installs a group of handler values for one body. The compiler can
lower this to nested handler machinery internally. Repeated `in handle`
disappears, but the single `do` still earns its keep as the boundary between
handler installation and the handled computation.

## Prototype: `examples/statefulhttp/server.osp`

```osp
effect Db
    add : string -> int
    list : Unit -> string
    count : Unit -> int

effect Metrics
    hit : string -> Unit
    served : Unit -> int

effect Persist
    flush : string -> int

effect Log
    info : string -> Unit


esc : string -> string -> string
esc code s =
    "\e[${code}m${s}\e[0m"

c256 : string -> string -> string
c256 n s =
    "\e[38;5;${n}m${s}\e[0m"

rose : string -> string
rose = c256 "213"

teal : string -> string
teal = c256 "159"

lime : string -> string
lime = c256 "46"

muted : string -> string
muted = c256 "239"

bold : string -> string
bold s = esc "1" s

dim : string -> string
dim s = esc "2" s

cyan : string -> string
cyan s = esc "36" s

green : string -> string
green s = esc "92" s

yellow : string -> string
yellow s = esc "93" s

magenta : string -> string
magenta s = esc "95" s

inc : int -> int
inc x = x + 1

mn : int -> int -> int
mn a b =
    match a < b
        true => a
        false => b


badge : string -> string
badge method =
    match method
        "POST" => "\e[1;30;48;5;213m POST \e[0m"
        "GET" => "\e[1;30;48;5;117m GET  \e[0m"
        _ => "\e[1;30;48;5;250m ${method} \e[0m"


blocks : int -> int -> string
blocks i n =
    match i < n
        true => "▰" + blocks (inc i) n
        false => ""

gauge : int -> string
gauge n =
    rose (blocks 0 (mn n 28))

divider : string
divider =
    muted "  ──────────────────────────────────────────────────────"

art : string
art =
    c256 "51"  "   ██████╗ ███████╗██████╗ ██████╗ ███████╗██╗   ██╗" + "\n" +
    c256 "45"  "  ██╔═══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝╚██╗ ██╔╝" + "\n" +
    c256 "39"  "  ██║   ██║███████╗██████╔╝██████╔╝█████╗   ╚████╔╝ " + "\n" +
    c256 "99"  "  ██║   ██║╚════██║██╔═══╝ ██╔══██╗██╔══╝    ╚██╔╝  " + "\n" +
    c256 "135" "  ╚██████╔╝███████║██║     ██║  ██║███████╗   ██║   " + "\n" +
    c256 "171" "   ╚═════╝ ╚══════╝╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   "

banner : string
banner =
    "\n" + art + "\n\n  " +
    lime "●" + dim " live  " + bold (teal "http://127.0.0.1:8080") +
    dim "   ·   state lives in algebraic-effect handlers" + "\n" +
    divider + "\n  " +
    badge "GET" + "  " + dim "/" + "         list tasks\n  " +
    badge "GET" + "  " + dim "/stats" + "    server metrics\n  " +
    badge "POST" + " " + dim "/tasks" + "    add a task " + dim "(body = task text)" + "\n" +
    divider


textResp : int -> string -> HttpResponse
textResp status bodyText =
    HttpResponse
        status = status
        headers = "Content-Type: text/plain"
        contentType = "text/plain"
        streamFd = -1
        isComplete = true
        partialBody = bodyText


filePersist : string -> Handler Persist
filePersist path =
    mut diskBytes = 0
    writeTasks = writeFile path

    handler Persist
        flush snap =>
            saved = writeTasks snap
            diskBytes := match saved
                Success value => length snap
                Error message => -1
            diskBytes


metricsCounter : Unit -> Handler Metrics
metricsCounter () =
    mut requests = 0

    handler Metrics
        hit p =>
            requests := requests + 1

        served =>
            requests


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


fileLog : string -> Handler Log
fileLog path =
    mut logBuf = ""
    writeServerLog = writeFile path

    handler Log
        info m =>
            logBuf := "${logBuf}${m}\n"
            _ = writeServerLog logBuf
            print "  ${m}"


onPost : string -> HttpResponse
onPost body =
    createdText = textResp 201
    id = perform Db.add body
    snap = perform Db.list
    written = perform Persist.flush snap
    perform Log.info (
        badge "POST" +
        "  " +
        green "✚" +
        " " +
        bold (yellow "\"${body}\"") +
        "  " +
        cyan "#${toString id}" +
        "  " +
        gauge id +
        dim "  ${toString written}B"
    )
    createdText "created task #${toString id}\n"


onGet : string -> HttpResponse
onGet path =
    okText = textResp 200
    perform Log.info (
        badge "GET" +
        " " +
        dim path +
        dim "  ·  ${toString (perform Metrics.served)} served"
    )

    match path
        "/stats" =>
            okText "requests=${toString (perform Metrics.served)} tasks=${toString (perform Db.count)}\n"

        _ =>
            okText (perform Db.list)


handleReq : string -> string -> string -> string -> HttpResponse
handleReq method path headers body =
    perform Metrics.hit path

    match method
        "POST" => onPost body
        _ => onGet path


serveForever : Unit -> Unit
serveForever () =
    sleep 3600000
    serveForever ()


main : Unit -> int
main () =
    persist = filePersist "/tmp/osprey_tasks.db"
    metrics = metricsCounter ()
    db = memoryDb ()
    log = fileLog "/tmp/osprey_server.log"

    handle persist metrics db log
    do
        createLocalServer = httpCreateServer 8080
        serverId = createLocalServer "127.0.0.1"
        listening = httpListen serverId handleReq
        print banner
        serveForever ()
        0
```

## Settled Choices For This Prototype

- Grouped handlers use `do`, not `in`. `do` is easier to read as "install these handlers, then run this computation".
- Old block braces are not accepted in the target syntax. Osprey should have one canonical syntax, not permanent brace and layout dialects.
- String interpolation keeps `${...}` for now. It is inside string literal syntax rather than block syntax, so it can stay unless it creates a concrete parser or readability problem.
- Function arrows are curried. `a -> b -> c` means `a -> (b -> c)`, and partial application is part of normal Osprey style.
- Handlers are first-class values. `handler Db ...` creates a `Handler Db`, and `handle db do ...` installs handler values around a computation.

## Remaining Syntax Questions

- Should zero-argument functions be called as `banner`, `banner ()`, or both? The prototype uses values for pure constants and `()` only where recursion or effects make the call boundary useful.
