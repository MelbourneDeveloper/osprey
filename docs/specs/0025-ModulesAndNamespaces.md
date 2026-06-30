# Modules and Namespaces

Osprey multi-file programs are built from **logical namespaces** and
**explicit modules**, not from file paths. A source file's path decides whether
it belongs to a project; it does **not** decide the names it exports.

> **Flavor layer - shared core (AST and above).** Namespace/import resolution,
> module signatures, exports, state ownership, separate compilation, and project
> assembly are shared-core semantics. The Default flavor and ML flavor may spell
> declarations differently, but both lower to the same canonical project model:
> `NamespaceDecl`, `ModuleDecl`, `SignatureDecl`, `Import`, and qualified names.
> No type checker, effect checker, code generator, runtime, or LSP feature may
> infer semantics from `.osp` vs `.ospml` once lowering has happened. See
> [Language Flavors](0023-LanguageFlavors.md).

## Status

`import` and `module` syntax are parsed today, and module bodies are checked in a
child scope. Cross-file resolution, open namespaces, explicit exports,
signatures, module-owned state rules, and project assembly are planned. This
chapter is the normative contract for those features and supersedes the
fiber-isolated module sketch in [Fibers and Concurrency](0011-LightweightFibersAndConcurrency.md#fiber-isolated-modules-planned).

## Research Basis

`[MODULES-RESEARCH]` The design combines .NET-style logical namespaces with
ML-style abstraction boundaries and Osprey's algebraic effects.

- Parnas set the bar for modularity: "The effectiveness of a \"modularization\" is
  dependent upon the criteria used" ([Parnas 1972](https://wstomv.win.tue.nl/edu/2ip30/references/criteria_for_modularization.pdf)).
- The C# specification gives the path-independent naming model: "Namespaces are
  open-ended" ([Microsoft C# spec, namespaces](https://learn.microsoft.com/en-us/dotnet/csharp/language-reference/language-specification/namespaces)).
- Harper and Lillibridge identify the core problem as "the management of the
  flow of information between program units" ([POPL 1994](https://www.cs.cmu.edu/~rwh/papers/sharing/popl94.pdf)).
- Rossberg, Russo, and Dreyer summarize the ML lesson: "ML modules are a
  powerful language mechanism for decomposing programs" ([F-ing Modules](https://people.mpi-sws.org/~dreyer/courses/modules/f-ing.pdf)).
- Leroy's manifest-types work requires a "strict distinction between abstract
  types and manifest types" ([POPL 1994](https://caml.inria.fr/pub/papers/xleroy-manifest_types-popl94.pdf)).
- Backpack states the separate-compilation target: "explicit interfaces express
  assumptions about dependencies" ([Kilpatrick, Dreyer, Peyton Jones, Marlow 2014](https://plv.mpi-sws.org/backpack/)).
- Launchbury and Peyton Jones justify encapsulated mutable state: "Some
  algorithms make critical internal use of updatable state" ([Lazy Functional State Threads](https://www.microsoft.com/en-us/research/publication/lazy-functional-state-threads/)).
- Plotkin and Pretnar make state a handled effect: effects include "state, time,
  and their combinations" ([Handlers of Algebraic Effects](https://homepages.inf.ed.ac.uk/gdp/publications/Effect_Handlers.pdf)).
- Moseley and Marks give the architectural rule: "Separate" essential state from
  essential logic and accidental state/control ([Out of the Tar Pit](https://curtclifton.net/papers/MoseleyMarks06a.pdf)).
- Linear Haskell points at the resource-state horizon: "typestates ... are
  actually enforced by the type system" ([Bernardy et al. 2018](https://arxiv.org/pdf/1710.09756)).
- Modern lexical effect handlers aim at "local-reasoning principles"
  ([Ma, Ge, Lee, Zhang 2024](https://cs.uwaterloo.ca/~yizhou/papers/lexa-oopsla2024.pdf)).
- Redux captures the state-management operational rule: "single source of truth"
  ([Redux Three Principles](https://redux.js.org/understanding/thinking-in-redux/three-principles)).

These are not ornamental citations. They drive the rules below: names are
logical, interfaces are explicit, abstract state does not leak, and mutable state
has one owner.

## Design Goals

`[MODULES-GOALS]` The module system must make the good structure the easy
structure:

- **Path-independent names.** A namespace name comes from source text, not from
  `src/foo/bar.osp`.
- **Open namespaces, closed modules.** Namespaces organize; modules encapsulate.
- **Explicit imports and exports.** Wildcard visibility is the escape hatch, not
  the default.
- **Separate compilation by interface.** A file can be checked against imported
  signatures without loading every implementation detail.
- **State has a declared owner.** Top-level mutable state is forbidden outside a
  state module or handler-owned state region.
- **Pure logic stays pure.** Modules expose state through effect-typed operations
  or pure query/update functions, not exported cells.
- **Cross-flavor interop.** A `.osp` module and `.ospml` module import each other
  through canonical signatures.

## Namespaces

`[MODULES-NAMESPACE]` A `namespace` declaration contributes declarations to an
open logical namespace. Multiple files may contribute to the same namespace.

```ebnf
namespaceDecl ::= "namespace" qualifiedName ("{" statement* "}" | ";")
qualifiedName ::= IDENT ("." IDENT)*
```

```osprey
namespace Acme.Billing {
    type Money = { cents: int, currency: string }
}

namespace Acme.Billing {
    fn zero(currency: string) -> Money = Money { cents: 0, currency: currency }
}
```

The two declarations above define one namespace, `Acme.Billing`. The compiler
merges namespace bodies before semantic analysis. Duplicate exported names in the
same namespace are compile-time errors unless they are overloads explicitly
allowed by a later overload spec.

`[MODULES-FILE-SCOPED-NAMESPACE]` A file-scoped namespace declaration applies to
all declarations after it in the file:

```osprey
namespace Acme.Billing;

type Invoice = { id: string, total: int }
fn emptyInvoice(id: string) = Invoice { id: id, total: 0 }
```

A file may contain either one file-scoped namespace declaration or any number of
block-scoped namespace declarations, not both.

`[MODULES-PATH-INDEPENDENCE]` The physical file path is never part of the
namespace identity. A file `src/weird/place/x.osp` may declare
`namespace Acme.Billing;`. The compiler may warn when path and namespace drift
from project convention, but it must not change symbol identity or import
resolution.

## Modules

`[MODULES-MODULE]` A `module` is a closed implementation boundary inside a
namespace. It may contain values, functions, types, effects, nested modules, and
private mutable state. It exports only declarations marked `export` or listed by
its signature.

```ebnf
moduleDecl ::= plainModuleDecl | stateModuleDecl
plainModuleDecl ::= "module" qualifiedName signatureAscription? "{" moduleItem* "}"
stateModuleDecl ::= "state" "module" qualifiedName signatureAscription? "{" moduleItem* "}"
signatureAscription ::= ":" qualifiedName
moduleItem ::= exportDecl | statement
exportDecl ::= "export" statement
```

```osprey
namespace Acme.Billing;

module Tax {
    let defaultRate = 10

    export fn addTax(cents: int) -> int =
        cents + cents * defaultRate / 100
}
```

`Tax.defaultRate` is private. `Tax.addTax` is exported.

`[MODULES-NAMESPACE-VS-MODULE]` Namespaces are open and stateless. Modules are
closed and may own private implementation details. A namespace cannot be used as
a runtime value; a module can be referenced as a qualified declaration space and,
when it is a `state module`, has a runtime state owner.

## Imports

`[MODULES-IMPORT]` Imports name namespaces or modules, not files.

```ebnf
importStmt ::= "import" qualifiedName importTail?
importTail ::= "as" IDENT
             | "." "{" importMember ("," importMember)* "}"
             | "." "*"
importMember ::= IDENT ("as" IDENT)?
```

```osprey
import Acme.Billing.Tax
import Acme.Billing.Tax.{addTax}
import Acme.Billing.Tax as Tax
```

Resolution rules:

- Fully qualified names always work: `Acme.Billing.Tax.addTax(100)`.
- `import A.B.C` brings the exported module or namespace `C` into the local
  scope as `C`.
- `import A.B.C.{x, y}` brings only listed exported members into local scope.
- `import A.B.C as Alias` brings `Alias` into local scope.
- `import A.B.C.*` is allowed only in examples, scripts, and tests unless the
  project enables `allow_wildcard_imports = true`; it is forbidden for state
  modules.
- Ambiguous unqualified names are compile-time errors. The diagnostic must show
  every imported candidate and suggest qualification or aliasing.

Imports do not execute code, allocate module state, or load files by relative
path.

## Exports And Visibility

`[MODULES-EXPORTS]` Declarations are private by default inside modules and
public by default inside namespaces. A module controls its public surface through
`export` or a signature.

```osprey
module Parser {
    type Token = { text: string }       // private
    export type Ast = Expr | Stmt
    export fn parse(source: string) -> Result<Ast, Error> = ...
}
```

`[MODULES-OPAQUE-TYPES]` A module may export an opaque type, hiding its
representation:

```osprey
module UserIds {
    export opaque type UserId = int

    export fn parseUserId(raw: string) -> Result<UserId, Error> = ...
    export fn showUserId(id: UserId) -> string = ...
}
```

Outside `UserIds`, `UserId` is distinct from `int`. Inside `UserIds`, the
manifest representation is available. This is the Osprey form of ML abstract
types and Leroy-style manifest types.

## Signatures

`[MODULES-SIGNATURE]` A `signature` is an explicit interface for a module. It
lists the names, types, effects, and opacity visible to clients.

```ebnf
signatureDecl ::= "signature" IDENT "{" signatureItem* "}"
signatureItem ::= typeSpec | effectSpec | fnSpec | moduleSpec
```

```osprey
signature StoreSig {
    opaque type Store
    effect StoreFx {
        load : fn() -> Store
        save : fn(Store) -> Unit
    }
    fn empty() -> Store
}

module MemoryStore : StoreSig {
    export opaque type Store = { values: [string] }
    export effect StoreFx {
        load : fn() -> Store
        save : fn(Store) -> Unit
    }
    export fn empty() = Store { values: [] }
}
```

Signature conformance is checked structurally:

- Every signature item must have a matching exported declaration.
- Types must match after applying opacity rules.
- Effect operations must match names, parameter types, return types, and effect
  rows.
- Extra private declarations are allowed.
- Extra exported declarations are rejected unless the ascription is marked
  `: StoreSig + extra`.

`[MODULES-SEPARATE-CHECKING]` A compiler may type-check an importing file using
only the imported module's signature. The implementation body is needed only
when compiling that module or linking the final project.

## Parameterised Modules

`[MODULES-FUNCTOR]` A parameterised module is a module-level function from
signatures to modules. This is planned after basic signatures.

```osprey
module MakeRepo(Db: DatabaseSig, Clock: ClockSig) : RepoSig {
    export fn save(item: Item) -> Unit !Db.Database =
        Db.insert(table: "items", value: encode(item, Clock.now()))
}
```

Parameterised modules are the dependency-injection mechanism for reusable
libraries. They are preferred over ambient globals.

## State Ownership

`[MODULES-STATE]` Mutable state may appear only in three places:

- inside a function or block as an ordinary local `mut`;
- inside an algebraic-effect handler's owned state region
  ([EFFECTS-HANDLER-STATE](0017-AlgebraicEffects.md#handler-owned-state));
- inside a `state module`.

Namespace-level `mut` is a compile-time error.

```osprey
namespace Acme.Bad;

mut count = 0
// error [MODULES-STATE-TOPLEVEL]:
// mutable state must live in a function, handler, or state module
```

`[MODULES-STATE-MODULE]` A `state module` is the declared owner of durable
module state. All state cells are private, and no `mut` cell may be exported.

```osprey
state module Counter {
    mut count = 0

    export effect CounterFx {
        next : fn() -> int
        read : fn() -> int
    }

    export let counterHandler = handler CounterFx {
        next => {
            count = count + 1
            count
        }
        read => count
    }
}
```

Clients perform the effect; the module owns the cell:

```osprey
import Counter.{CounterFx, counterHandler}

fn allocate() -> int !CounterFx =
    perform CounterFx.next()

handle counterHandler in {
    print(toString(allocate()))
}
```

Rules:

- `state module` cells are private by construction.
- Exporting a `mut`, a pointer to a `mut`, or a closure that directly exposes
  assignment is a compile-time error.
- A `state module` must export at least one handler, effect, or function that is
  the declared access path.
- A namespace may contain at most one unannotated `state module`. Additional
  state owners require `@state_boundary("reason")` and are reported by LSP and
  docs tooling as architecture-visible state boundaries.
- Derived state should be expressed as pure functions over owner state. Cached
  derived state is forbidden in Phase 1; a later `cache mut` feature must name
  the owner state it derives from, so invalidation can be checked.

`[MODULES-STATE-SOURCE-OF-TRUTH]` The compiler and tooling treat each state
module as a **single source of truth** for the state it owns. Cross-module writes
are impossible. Cross-module reads happen through exported pure queries or
effect operations. This is the language-level answer to scattered app state.

## Effects And Capabilities

`[MODULES-EFFECTS]` Modules do not hide effects. Exported functions and handlers
carry ordinary Osprey effect rows. Importing a module never grants ambient
permission; a caller must still handle or forward every effect.

State modules are encouraged to expose capabilities as algebraic effects:

```osprey
signature LedgerSig {
    effect Ledger {
        post : fn(int) -> int
        balance : fn() -> int
    }
}
```

This keeps application logic pure except for explicit `!Ledger`, while the
module decides whether state is in memory, SQLite, HTTP, or a test fake.

## Initialisation

`[MODULES-INIT]` Imports have no runtime effect. Module initialization is explicit.

- Pure `let` declarations may be evaluated at compile time or lowered as
  constants.
- Effectful setup must live in an exported `init` function or handler factory.
- `state module` initial state is allocated only when its handler or instance is
  explicitly constructed.
- Cyclic initialization between state modules is a compile-time error.

```osprey
state module DbStore {
    mut conn = None

    export fn init(path: string) -> Unit !Database =
        conn = Some(perform Database.open(path))
}
```

## Project Assembly

`[MODULES-PROJECT]` A project compile scans configured source roots, parses every
`.osp` and `.ospml` file, resolves each file's flavor, and builds one project
namespace graph.

```toml
[project]
name = "acme.billing"
source_roots = ["src", "tests"]
root_namespace = "Acme"

[modules]
allow_wildcard_imports = false
```

Single-file mode remains valid for scripts and examples. Project mode adds:

- all source files in the project graph;
- namespace merge;
- import resolution;
- signature checking;
- duplicate-name and ambiguity diagnostics;
- one entry point.

`[MODULES-ENTRYPOINT]` In project mode, executable top-level statements are
allowed only in the designated entry file or in `fn main()`. Library files must
contain declarations only. This avoids hidden initialization order and makes
multi-file apps deterministic.

## Cycles

`[MODULES-CYCLES]` Namespace declarations may be mutually visible after merging,
but module implementation cycles are restricted.

- Pure type/function cycles are allowed only when ordinary Osprey recursion rules
  allow them.
- Signature cycles are allowed only through explicit opaque types.
- `state module` cycles are rejected.
- Parameterised modules may depend on signatures, not implementation bodies, to
  preserve separate compilation.

Recursive modules are a later feature and must require explicit signatures, as in
the ML literature.

## Name Mangling And ABI

`[MODULES-ABI]` Canonical names are fully qualified:

```text
Acme.Billing.Tax.addTax
```

Codegen must mangle qualified names deterministically and collision-free. The
mangled form is an implementation detail; diagnostics, docs, LSP, debugger, and
stack traces use source-level qualified names.

Cross-flavor exports use the same ABI rules as
[Cross-Flavor Interop](0023-LanguageFlavors.md#cross-flavor-interop).

## Diagnostics

`[MODULES-DIAG]` Module diagnostics must be architecture-facing:

- unknown import: show candidate namespaces from the project graph;
- ambiguous import: show all providers and suggest aliases;
- exported private dependency: show the hidden type/value in the public signature;
- state scatter: show every state module in the namespace and require
  `@state_boundary`;
- top-level mutable state: suggest `state module` or handler-owned state;
- path drift: warn, never change semantics.

## Examples

### Multi-file, Path-Independent Namespace

`src/a.osp`:

```osprey
namespace Acme.App;

fn hello(name: string) = "Hello ${name}"
```

`src/deeply/nested/b.ospml`:

```osprey
namespace Acme.App

greet name = hello name
```

Both files contribute to `Acme.App`. The path `deeply/nested` is irrelevant.

### Centralised State

```osprey
namespace Acme.App;

state module SessionStore {
    mut sessions = []

    export effect Sessions {
        add : fn(string) -> Unit
        count : fn() -> int
    }

    export let liveSessions = handler Sessions {
        add id => { sessions = listAppend(sessions, id) }
        count => listLength(sessions)
    }
}

fn login(id: string) -> Unit !Sessions =
    perform Sessions.add(id)
```

Application code cannot mutate `sessions`. It can only perform `Sessions`.

## References

- David L. Parnas. "On the Criteria To Be Used in Decomposing Systems into
  Modules." Communications of the ACM, 1972.
- David MacQueen. "Modules for Standard ML." LFP, 1984.
- John C. Mitchell and Gordon D. Plotkin. "Abstract Types Have Existential
  Type." POPL 1985 / TOPLAS 1988.
- Robert Harper and Mark Lillibridge. "A Type-Theoretic Approach to Higher-Order
  Modules with Sharing." POPL 1994.
- Xavier Leroy. "Manifest Types, Modules, and Separate Compilation." POPL 1994.
- Xavier Leroy. "Applicative Functors and Fully Transparent Higher-Order
  Modules." POPL 1995.
- Xavier Leroy. "A Modular Module System." JFP, 2000.
- Karl Crary, Robert Harper, and Sidd Puri. "What is a Recursive Module?" PLDI
  1999.
- Keiko Nakata and Jacques Garrigue. "Recursive Modules for Programming." ICFP
  2006.
- Andreas Rossberg, Claudio V. Russo, and Derek Dreyer. "F-ing Modules." TLDI
  2010 / JFP.
- Andreas Rossberg. "1ML - Core and Modules United." ICFP 2015 / JFP.
- Scott Kilpatrick, Derek Dreyer, Simon Peyton Jones, and Simon Marlow.
  "Backpack: Retrofitting Haskell with Interfaces." POPL 2014.
- Gordon Plotkin and Matija Pretnar. "Handlers of Algebraic Effects." ESOP 2009.
- John Launchbury and Simon L. Peyton Jones. "Lazy Functional State Threads."
  PLDI 1994.
- Simon Peyton Jones and Philip Wadler. "Imperative Functional Programming."
  POPL 1993.
- Ben Moseley and Peter Marks. "Out of the Tar Pit." 2006.
- Jean-Philippe Bernardy, Mathieu Boespflug, Ryan R. Newton, Simon Peyton Jones,
  and Arnaud Spiwack. "Linear Haskell: Practical Linearity in a Higher-Order
  Polymorphic Language." POPL 2018.
- Cong Ma, Zhaoyi Ge, Edward Lee, and Yizhou Zhang. "Lexical Effect Handlers,
  Directly." OOPSLA 2024.
- Microsoft. "C# Language Specification - Namespaces."
- Redux. "Three Principles."
