## 12. Lightweight Fibers and Concurrency

🚧 **IMPLEMENTATION STATUS**: Fiber syntax is partially implemented. Basic fiber operations (`spawn`, `await`, `yield`) are in the grammar but runtime support is limited.

❌ **NOT IMPLEMENTED**: The fiber-isolated module system is a design goal but not yet implemented. Current module support is basic.

### 12.1 Fiber Types and Concurrency

Osprey provides lightweight concurrency through fiber types. Unlike traditional function-based approaches, fibers are proper type instances constructed using Osprey's standard type construction syntax.

#### Core Fiber Types

**`Fiber<T>`** - A lightweight concurrent computation that produces a value of type T
**`Channel<T>`** - A communication channel for passing values of type T between fibers

#### Fiber Construction

Fibers are created using standard type construction syntax:

```osprey
// Create a fiber that computes a value
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(n: 1000) 
}

// Create a fiber with more complex computation
let worker = Fiber<String> { 
    computation: fn() => {
        processData()
        "completed"
    }
}

// Create a parameterized fiber
let calculator = Fiber<Int> { 
    computation: fn() => multiply(x: 10, y: 20) 
}
```

#### Spawn Syntax Sugar

For convenience, Osprey provides `spawn` as syntax sugar for creating and immediately starting a fiber:

```osprey
// Using spawn (syntax sugar)
let result = spawn 42

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => 42 }
let result = fiber

// More complex spawn
let computation = spawn (x * 2 + y)

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => x * 2 + y }
let computation = fiber
```

The `spawn` keyword immediately evaluates the expression in a new fiber context, making it convenient for quick concurrent computations without the full type construction syntax.

#### Channel Construction

Channels are created using type construction syntax:

```osprey
// Unbuffered (synchronous) channel
let sync_channel = Channel<Int> { capacity: 0 }

// Buffered (asynchronous) channel  
let async_channel = Channel<String> { capacity: 10 }

// Large buffer channel
let buffer_channel = Channel<Int> { capacity: 100 }
```

#### Fiber Operations

Once created, fibers and channels are manipulated using functional operations:

**`await(fiber: Fiber<T>) -> T`** - Wait for fiber completion and get result
**`send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`** - Send value to channel
**`recv(channel: Channel<T>) -> Result<T, ChannelError>`** - Receive value from channel
**`yield() -> Unit`** - Voluntarily yield control to scheduler

```osprey
// Create and await a fiber
let task = Fiber<Int> { computation: fn() => heavyComputation() }
let result = await(task)

// Channel communication
let ch = Channel<String> { capacity: 5 }
send(ch, "hello")
let message = recv(ch)

// Yielding control
yield()
```

#### Complete Fiber Example

```osprey
// Producer fiber
let producer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        send(ch, 1)
        send(ch, 2) 
        send(ch, 3)
    }
}

// Consumer fiber
let consumer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        let value1 = recv(ch)
        let value2 = recv(ch)
        let value3 = recv(ch)
        print("Received: ${value1}, ${value2}, ${value3}")
    }
}

// Start both fibers
await(producer)
await(consumer)
```

#### Select Expression for Channel Multiplexing

The `select` expression allows waiting on multiple channel operations:

```osprey
let ch1 = Channel<String> { capacity: 1 }
let ch2 = Channel<Int> { capacity: 1 }

let result = select {
    msg => recv(ch1) => process_string(msg)
    num => recv(ch2) => process_number(num)
    _ => timeout_handler()
}
```

#### Rust Interoperability

Osprey fibers are designed to interoperate with Rust's async/await system:

```osprey
// Osprey fiber that calls Rust async function
extern fn rust_async_task() -> Future<Int>

let osprey_task = Fiber<Int> {
    computation: fn() => await(rust_async_task())
}

let result = await(osprey_task)
```

## 12.2 Fiber-Isolated Module System

❌ **NOT IMPLEMENTED**: The fiber-isolated module system is a design goal but not yet implemented. Current module support is basic.

### Module Isolation Principles

The fiber-isolated module system eliminates data races by design through:

1. **Fiber-Local State**: Each fiber gets its own isolated copy of module state
2. **No Shared Mutable State**: Modules cannot share mutable data between fibers
3. **Immutable Sharing**: Only immutable data can be shared between fibers
4. **Automatic Isolation**: Module isolation happens automatically without explicit synchronization

## 13.2 Module Declaration Syntax

```osprey
module ModuleName {
    // Module declarations
    let value = 42
    mut counter = 0
    
    fn increment() -> Int = {
        counter = counter + 1
        counter
    }
    
    fn getValue() -> Int = value
}
```

### Fiber Isolation Behavior

When a fiber accesses a module, it gets its own isolated instance:

```osprey
module Counter {
    mut count = 0
    
    fn increment() -> Int = {
        count = count + 1
        count
    }
    
    fn get() -> Int = count
}

// Each fiber gets its own Counter instance
let fiber1 = spawn Counter.increment()  // Returns 1
let fiber2 = spawn Counter.increment()  // Also returns 1 (separate instance)

let result1 = await(fiber1)  // 1
let result2 = await(fiber2)  // 1 (not 2!)
```

### Memory and Performance Characteristics

- **Copy-on-First-Access**: Module instances are copied when first accessed by a fiber
- **Memory Isolation**: Each fiber's module state is completely isolated
- **No Synchronization Overhead**: No locks, atomics, or other synchronization primitives needed
- **Deterministic Behavior**: Same input always produces same output within a fiber

## 13.5 Inter-Fiber Communication

Since modules are isolated, inter-fiber communication must use explicit channels:

```osprey
module Database {
    mut connections = []
    
    fn connect() -> Connection = {
        // This connection is fiber-local
        let conn = createConnection()
        connections = conn :: connections
        conn
    }
}

// Fibers communicate via channels, not shared module state
let resultChannel = Channel<String> { capacity: 10 }

let worker1 = spawn {
    let conn = Database.connect()  // Fiber-local connection
    let result = query(conn, "SELECT * FROM users")
    send(resultChannel, result)
}

let worker2 = spawn {
    let conn = Database.connect()  // Different fiber-local connection  
    let result = query(conn, "SELECT * FROM products")
    send(resultChannel, result)
}
```

This design ensures that concurrent access to modules is always safe without requiring explicit synchronization. 