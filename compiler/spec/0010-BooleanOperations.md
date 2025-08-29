10. [Boolean Operations](0010-BooleanOperations.md)
    - [Boolean Pattern Matching](#101-boolean-pattern-matching)
    - [Boolean Operators](#102-boolean-operators)

# 10. Boolean Operations

Use pattern matching for conditional logic:

**Examples:**
```osprey
let result = match x > 0 {
    true => "positive"
    false => "zero or negative"
}

let max = match a > b {
    true => a
    false => b
}
```

## 10.1 Boolean Pattern Matching

Osprey uses pattern matching instead of traditional if-else statements for boolean operations. This ensures exhaustive handling of both true and false cases.

**Basic Boolean Matching:**
```osprey
let status = match isValid {
    true => "Success"
    false => "Failure"
}
```

**Complex Boolean Logic:**
```osprey
let category = match score >= 90 {
    true => match score == 100 {
        true => "Perfect"
        false => "Excellent"
    }
    false => match score >= 70 {
        true => "Good"
        false => "Needs Improvement"
    }
}
```

## 10.2 Boolean Operators

- `&&` - Logical AND
- `||` - Logical OR  
- `!` - Logical NOT
- `==` - Equality
- `!=` - Inequality
- `>`, `<`, `>=`, `<=` - Comparison operators

**Operator Examples:**
```osprey
let isAdult = age >= 18
let hasPermission = isAdult && isAuthorized
let canAccess = hasPermission || isAdmin
let isBlocked = !isActive

// Complex logical expressions with parentheses
let complexLogic = (score >= 90) && (attendance > 0.8)
let shouldNotify = (status == "urgent") || (priority > 5)
let validUser = !isBanned && (isVerified || hasInvite)
```

**Short-Circuit Evaluation:**

Logical operators use short-circuit evaluation for performance:
- `&&` (AND): If left operand is false, right operand is not evaluated
- `||` (OR): If left operand is true, right operand is not evaluated