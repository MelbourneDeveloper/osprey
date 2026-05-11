# Boolean Operations

Osprey has no `if`/`else` statement. All conditional logic uses `match` on a boolean expression, which forces both arms to be considered.

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