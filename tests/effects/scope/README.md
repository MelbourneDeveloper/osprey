# Handler lifetime and scope tests

These paired Default/ML suites pin the parts of the [effects contract](../../../docs/specs/language/0017-AlgebraicEffects.md) that decide *which* activation answers and *when* an activation ends:

- `finally` cleanup runs exactly once per activation, after the return clause, inner to outer, and when a control arm abandons the remainder (`[EFFECTS-FINALIZATION]`, `[MULTI-COST-ABORT]`);
- `mask E` skips one enclosing activation per level and keeps handlers installed inside its body (`[EFFECTS-MASK]`);
- `instance E` creates a capability that only its own handlers answer, distinct from the default identity (`[EFFECTS-INSTANCE]`).

[Plan 0016](../../../docs/plans/0016-algebraic-effects-and-handlers.md) lists finalizers, masking and named instances as unfinished, so these suites stay red until they ship. The goldens record the behaviour the specification requires, not what the compiler does today.

```sh
target/release/osprey test tests/effects/scope
```
