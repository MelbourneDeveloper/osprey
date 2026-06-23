# Plan 0003 — Match Exhaustiveness With Catch-Alls

**Subsystem:** `crates/osprey-types`
**Status:** Partially implemented
**Spec:** [0007-PatternMatching.md](../specs/0007-PatternMatching.md), [0004-TypeSystem.md](../specs/0004-TypeSystem.md) (`any` exhaustiveness)

## Summary

Exhaustiveness checking exists and is correct for the no-catch-all case: a `bool`
match needs both `true` and `false`; a match on a known union needs every
variant. Two gaps weaken it: (1) the presence of **any** catch-all arm disables
the check entirely, and (2) `is_catch_all` is too broad, so an ordinary binding
name can be misread as a catch-all.

## What works today

- `check_exhaustive` reports missing union variants and missing bool arms —
  [crates/osprey-types/src/pattern.rs:289](../../crates/osprey-types/src/pattern.rs#L289).
- `check_bool_exhaustive` requires both boolean constructors —
  [crates/osprey-types/src/pattern.rs](../../crates/osprey-types/src/pattern.rs).

## Where it falls short

```rust
// crates/osprey-types/src/pattern.rs:289 (check_exhaustive)
if arms.iter().any(|a| self.is_catch_all(&a.pattern)) {
    return;            // <-- whole check skipped if ANY arm is a catch-all
}
```

```rust
// is_catch_all (crates/osprey-types/src/pattern.rs:335)
Pattern::Binding(name) => self.ctors.get(name).is_none_or(|i| !i.fields.is_empty()),
```

Effects:

1. `match x { A => …, _ => … }` on a 3-variant union is accepted even though the
   author may have meant to cover each variant — the catch-all silently absorbs
   the rest. (This is intentionally legal, but we lose *redundancy/anomaly*
   feedback and there is no "useless arm after catch-all" check.)
2. A bare binding that is not a known nullary constructor is treated as a
   catch-all. A typo of a constructor name (`Sucess` instead of `Success`) becomes
   a "binding", silently turns the arm into a catch-all, and suppresses the
   missing-variant error.

## Implementation plan

1. **Tighten `is_catch_all`.** Only `Pattern::Wildcard` and
   `Pattern::TypeAnnotated` (a typed binding) are true catch-alls. A bare
   `Pattern::Binding(name)` should be a catch-all **only** when `name` is not a
   known constructor *and* the discriminant type is not a union with a
   same-spelled variant — otherwise treat it as an attempted (mis-spelled)
   constructor and let the missing-variant path report it.
2. **Keep catch-all legal but add redundancy detection.** When a catch-all is
   present, do not error on missing variants, but emit a warning-style diagnostic
   for any arm that is unreachable because an earlier arm (or the catch-all)
   already covers it.
3. **Validate constructor spelling.** When an arm's head is an identifier that
   looks like a constructor (capitalized / matches a variant namespace) but is
   unknown, report "unknown constructor `X`" rather than silently binding.
4. **(Optional) nested exhaustiveness.** For `Ctor { a, b }` patterns, verify all
   declared fields are accounted for; defer if it grows the scope.

## Testing

- Add `failscompilation` cases: (a) misspelled constructor that today slips
  through as a binding; (b) genuinely non-exhaustive union match with no
  catch-all (already covered — keep).
- Extend a `tested/basics` match example with a legitimate catch-all to prove it
  still compiles and runs.

## Risks / considerations

- Distinguishing "intended constructor, misspelled" from "intended binding"
  needs the discriminant type in scope — `check_exhaustive` already receives
  `disc`, so thread it into the spelling check.
- Avoid false positives on legitimately generic bindings (e.g. matching an `any`
  with a single binding arm).

## TODO

- [ ] Narrow `is_catch_all` so a bare binding is a catch-all only when it is not a
      (possibly misspelled) constructor for the discriminant type.
- [ ] Report "unknown constructor" for capitalized/variant-shaped unknown heads.
- [ ] Add redundant/unreachable-arm diagnostics when a catch-all is present.
- [ ] (Optional) nested-field exhaustiveness for record-constructor patterns.
- [ ] Add `failscompilation` cases for misspelled constructor + non-exhaustive
      union.
- [ ] Confirm existing catch-all `tested` examples still pass; `make ci` green.
