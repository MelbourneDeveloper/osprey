---
layout: page.njk
title: "Building Type-Safe Web APIs with Osprey's Pattern Matching"
excerpt: "Discover how Osprey's exhaustive pattern matching and algebraic data types eliminate entire classes of runtime errors in web API development, making your services more reliable and maintainable."
date: 2025-01-10
tags: ["blog", "web-development", "pattern-matching", "type-safety", "apis"]
author: "Christian Findlay"
readingTime: 6
---

Web APIs are the backbone of modern applications, but they're also where things go wrong most often. Null reference exceptions, unhandled error cases, and missing validation logic plague even the most carefully written services. **What if your programming language could eliminate these problems entirely?**

Osprey's **exhaustive pattern matching** and **algebraic data types** provide exactly that guarantee. Let's explore how to build rock-solid web APIs that handle every possible case—and how the compiler ensures you never forget an edge case.

## The Problem with Traditional Error Handling

Most web frameworks handle errors through exceptions or error codes. Here's typical Node.js/Express code:

```javascript
app.post('/users', async (req, res) => {
  try {
    const user = await createUser(req.body);
    if (user) {
      res.json({ success: true, data: user });
    } else {
      res.status(400).json({ error: "Failed to create user" });
    }
  } catch (error) {
    if (error.code === 'DUPLICATE_EMAIL') {
      res.status(409).json({ error: "Email already exists" });
    } else {
      res.status(500).json({ error: "Internal server error" });
    }
  }
});
```

**What's wrong here?** The compiler can't verify you've handled every error type. You might forget to handle a database timeout, miss a validation error, or incorrectly assume what `createUser` returns.

## Modeling API Responses with Algebraic Data Types

Osprey takes a different approach. We model all possible outcomes as **algebraic data types**:

```osprey
type CreateUserResult = 
  | Success(user: User, id: UserId)
  | ValidationError(fields: List<String>, messages: List<String>)
  | DuplicateEmail(email: String)
  | DatabaseTimeout(retryAfter: Int)
  | DatabaseError(message: String)
  | InternalError(context: String)

type ApiResponse<T> = 
  | Ok(data: T, status: Int)
  | Error(message: String, code: String, status: Int)
```

Now our API handler **must** handle every case:

```osprey
fn createUserHandler(request: HttpRequest) -> HttpResponse {
  let userData = parseUserData(request.body)
  
  match createUser(userData) {
    Success(user, id) -> 
      Ok({ user: user, id: id }, 201),
      
    ValidationError(fields, messages) ->
      Error(
        "Validation failed: " + messages.join(", "),
        "VALIDATION_ERROR",
        400
      ),
      
    DuplicateEmail(email) ->
      Error(
        "Email " + email + " is already registered",
        "DUPLICATE_EMAIL", 
        409
      ),
      
    DatabaseTimeout(retryAfter) ->
      Error(
        "Database temporarily unavailable",
        "DATABASE_TIMEOUT",
        503
      ).withHeader("Retry-After", retryAfter.toString()),
      
    DatabaseError(message) ->
      Error(
        "Database error occurred",
        "DATABASE_ERROR",
        502
      ).withLogging("Database error: " + message),
      
    InternalError(context) ->
      Error(
        "Internal server error",
        "INTERNAL_ERROR",
        500
      ).withLogging("Internal error in " + context)
  }
}
```

**The compiler guarantees** you handle every case. If you add a new error type, every pattern match that handles `CreateUserResult` will fail to compile until you add the new case.

## Building a Complete REST API

Let's build a complete user management API to see how this scales:

```osprey
module UserApi {
  type User = {
    id: UserId,
    email: String,
    name: String,
    createdAt: DateTime,
    isActive: Boolean
  }
  
  type UserFilter = 
    | All
    | Active
    | Inactive
    | CreatedAfter(date: DateTime)
    | EmailDomain(domain: String)
  
  type UserOperation<T> = 
    | Success(data: T)
    | NotFound(id: UserId)
    | ValidationError(field: String, message: String)
    | PermissionDenied(action: String)
    | DatabaseError(details: String)
    | RateLimited(resetTime: DateTime)
  
  // GET /users
  fn listUsers(filter: UserFilter, auth: AuthToken) -> UserOperation<List<User>> {
    if (!auth.hasPermission("users:read")) {
      return PermissionDenied("list users")
    }
    
    match Database.queryUsers(filter) {
      Ok(users) -> Success(users),
      Err(dbError) -> DatabaseError(dbError.toString())
    }
  }
  
  // GET /users/:id
  fn getUser(id: UserId, auth: AuthToken) -> UserOperation<User> {
    if (!auth.hasPermission("users:read")) {
      return PermissionDenied("read user")
    }
    
    match Database.findUser(id) {
      Some(user) -> Success(user),
      None -> NotFound(id)
    }
  }
  
  // PUT /users/:id
  fn updateUser(id: UserId, updates: UserUpdates, auth: AuthToken) -> UserOperation<User> {
    if (!auth.hasPermission("users:write")) {
      return PermissionDenied("update user")
    }
    
    match validateUserUpdates(updates) {
      Invalid(field, message) -> ValidationError(field, message),
      Valid(validUpdates) -> 
        match Database.updateUser(id, validUpdates) {
          Some(updatedUser) -> Success(updatedUser),
          None -> NotFound(id)
        }
    }
  }
}
```

## HTTP Response Conversion

Converting our domain types to HTTP responses becomes a pure mapping function:

```osprey
fn toHttpResponse<T>(operation: UserOperation<T>, encoder: T -> Json) -> HttpResponse {
  match operation {
    Success(data) -> 
      HttpResponse(200, encoder(data)),
      
    NotFound(id) ->
      HttpResponse(404, {
        error: "User not found",
        code: "USER_NOT_FOUND", 
        userId: id.toString()
      }),
      
    ValidationError(field, message) ->
      HttpResponse(400, {
        error: "Validation failed",
        code: "VALIDATION_ERROR",
        field: field,
        message: message
      }),
      
    PermissionDenied(action) ->
      HttpResponse(403, {
        error: "Permission denied",
        code: "PERMISSION_DENIED",
        action: action
      }),
      
    DatabaseError(details) ->
      HttpResponse(500, {
        error: "Internal server error",
        code: "DATABASE_ERROR"
      }).withLogging("DB Error: " + details),
      
    RateLimited(resetTime) ->
      HttpResponse(429, {
        error: "Rate limit exceeded", 
        code: "RATE_LIMITED"
      }).withHeader("X-Rate-Limit-Reset", resetTime.toIsoString())
  }
}
```

## Request Routing with Pattern Matching

Osprey's pattern matching shines for request routing too:

```osprey
fn routeRequest(request: HttpRequest) -> HttpResponse {
  match (request.method, request.path, request.auth) {
    (GET, "/users", Some(auth)) ->
      let filter = parseUserFilter(request.query)
      UserApi.listUsers(filter, auth) |> toHttpResponse(encodeUserList),
      
    (GET, "/users/" + userId, Some(auth)) ->
      match UserId.parse(userId) {
        Some(id) -> UserApi.getUser(id, auth) |> toHttpResponse(encodeUser),
        None -> HttpResponse(400, { error: "Invalid user ID" })
      },
      
    (PUT, "/users/" + userId, Some(auth)) ->
      match (UserId.parse(userId), parseUserUpdates(request.body)) {
        (Some(id), Some(updates)) -> 
          UserApi.updateUser(id, updates, auth) |> toHttpResponse(encodeUser),
        (None, _) -> 
          HttpResponse(400, { error: "Invalid user ID" }),
        (_, None) -> 
          HttpResponse(400, { error: "Invalid request body" })
      },
      
    (_, _, None) ->
      HttpResponse(401, { error: "Authentication required" }),
      
    (method, path, _) ->
      HttpResponse(404, { 
        error: "Endpoint not found",
        method: method.toString(),
        path: path 
      })
  }
}
```

## Middleware as Function Composition

Middleware becomes simple function composition:

```osprey
fn withRateLimit<T>(handler: HttpRequest -> T, limit: RateLimit) -> HttpRequest -> UserOperation<T> {
  fn(request) {
    match RateLimit.check(request.clientIp, limit) {
      Allowed -> Success(handler(request)),
      Limited(resetTime) -> RateLimited(resetTime)
    }
  }
}

fn withAuth<T>(handler: (HttpRequest, AuthToken) -> T) -> HttpRequest -> UserOperation<T> {
  fn(request) {
    match Auth.validateToken(request.headers.authorization) {
      Valid(token) -> Success(handler(request, token)),
      Invalid -> PermissionDenied("invalid token"),
      Missing -> PermissionDenied("authentication required")
    }
  }
}

// Compose middleware naturally
let protectedHandler = listUsersHandler
  |> withAuth
  |> withRateLimit(RateLimit.perMinute(100))
```

## Testing Becomes Trivial

Since everything is a pure function returning data, testing is incredibly straightforward:

```osprey
test "user creation handles all error cases" {
  // Test successful creation
  assert UserApi.createUser(validUserData) == Success(expectedUser, expectedId)
  
  // Test validation errors
  assert UserApi.createUser(invalidEmail) == ValidationError(["email"], ["Invalid format"])
  
  // Test duplicate email
  assert UserApi.createUser(duplicateEmail) == DuplicateEmail("test@example.com")
  
  // The compiler ensures we test every possible return value
}
```

## The Reliability Advantage

This approach eliminates entire categories of production bugs:

- **No null reference exceptions** - Options and Results make nullability explicit
- **No unhandled error cases** - Pattern matching forces you to handle every scenario  
- **No silent failures** - Every operation's outcome is explicitly modeled
- **No incorrect status codes** - HTTP responses are generated deterministically
- **No missing validation** - Validation is built into the type system

## Performance Benefits

Despite the high-level abstractions, Osprey compiles to efficient code:

```osprey
// This high-level code...
let result = users
  |> filter(u -> u.isActive)
  |> map(u -> { ...u, lastSeen: now() })
  |> take(10)

// ...compiles to an efficient loop with no intermediate allocations
```

The pattern matching compiles to efficient jump tables, and the functional pipelines are optimized away entirely.

## Conclusion

Building web APIs with Osprey's pattern matching transforms error-prone, defensive programming into **compiler-verified correctness**. You can't ship a handler that forgets to handle an error case because **it won't compile**.

This isn't just academic elegance—it's practical reliability. When your API serves millions of requests per day, the difference between "it should work" and "it can't fail" is the difference between 3 AM pages and peaceful sleep.

The functional programming revolution isn't coming—**it's here**. And for web development, the benefits are too compelling to ignore.

---

*Ready to try building your own type-safe APIs? Check out our [web development guide](/docs/web-development/) or experiment in the [playground](/playground/).* 