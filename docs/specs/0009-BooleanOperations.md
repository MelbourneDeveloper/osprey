# Boolean Operations

Osprey has no `if`/`else` statement. Conditional logic is written as a `match` on a boolean (which forces both arms to be considered) or as the ternary shorthand `cond ? then : else`, which desugars to the same `match`. The ternary is defined in [Pattern Matching](0007-PatternMatching.md#ternary-match-syntactic-sugar).

> **Flavor layer — shared core (AST and above).**  Boolean semantics are flavor-blind. `&&`, `||`, and the comparison operators lower to `Expr::Binary`, `!` to `Expr::Unary`, and every conditional to `Expr::Match` over the boolean — the canonical AST nodes the type checker, effect checker, and codegen consume without ever knowing which flavor produced them. The `match` *spelling* differs between the Default surface shown here (braces) and the ML offside form in [ML Flavor Syntax](0024-MLFlavorSyntax.md), but the desugaring and short-circuit semantics described in this chapter are identical across both. See [Language Flavors](0023-LanguageFlavors.md).

```osprey
let status = match isValid {
    true  => "Success"
    false => "Failure"
}

let max = match a > b {
    true  => a
    false => b
}
```

Nested matches handle compound conditions:

```osprey
let category = match score >= 90 {
    true  => match score == 100 {
        true  => "Perfect"
        false => "Excellent"
    }
    false => match score >= 70 {
        true  => "Good"
        false => "Needs Improvement"
    }
}
```

## Boolean Operators

`&&`, `||`, and `!` are short-circuiting; `==`, `!=`, `<`, `>`, `<=`, `>=` produce booleans. See [Lexical Structure](0002-LexicalStructure.md) for the full operator list.

```osprey
let isAdult       = age >= 18
let hasPermission = isAdult && isAuthorized
let canAccess     = hasPermission || isAdmin
let isBlocked     = !isActive
let validUser     = !isBanned && (isVerified || hasInvite)
```