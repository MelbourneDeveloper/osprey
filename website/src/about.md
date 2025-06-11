---
layout: base.njk
title: "About Osprey"
description: "Osprey's design principles and language philosophy"
---

<section class="page-hero">
  <div class="container">
    <h1>About Osprey</h1>
    <p class="lead">
      A functional programming language designed for safety, elegance, and performance.
    </p>
  </div>
</section>

<section class="about-story">
  <div class="container">
    <div class="two-column">
      <div class="content">
        <h2>The Vision</h2>
        <p>
          Osprey was born from the belief that elegance <strong>is</strong> the best defence against bugs. Too many hours are wasted debugging null pointer exceptions, handling unexpected panics, and tracking down race conditions. Osprey is different.
        </p>
        <p>
          Osprey's type system aims to prevent entire classes of bugs at compile time, make concurrency safe by default, and keep your code maintainable for years to come.
        </p>
      </div>
      <div class="visual typewriter-enabled">
        <div class="code-example">
          <pre class="language-osprey"><code class="language-osprey">
type User = {
  email: String where isValidEmail(email),
  age: Int where between(age, 0, 150)
}

// Result type ensures error handling
let user = User {
  email: "alice@example.com",
  age: 25
}

match user {
  Ok { value } => welcomeUser(value),
  Err { error } => showValidationError(error)
}</code></pre>
        </div>
      </div>
    </div>
  </div>
</section>

<section class="design-principles">
  <div class="container">
    <h2 class="section-title">Design Principles</h2>
    <div class="principles-detailed">
      <div class="principle-item">
        <div class="principle-header">
          <h3>🎯 Simplicity First</h3>
        </div>
        <div class="principle-content">
          <p>One way to accomplish any task. No multiple syntax variations, minimal ceremony.</p>
          <ul>
            <li>Single approach for each language feature</li>
            <li>Named arguments for clarity</li>
            <li>Self-documenting code</li>
          </ul>
        </div>
      </div>

      <div class="principle-item">
        <div class="principle-header">
          <h3>🛡️ Safety by Design</h3>
        </div>
        <div class="principle-content">
          <p>Type system prevents bugs at compile time. No null pointers, buffer overflows, or data races.</p>
          <ul>
            <li>Result types for all error conditions</li>
            <li>Memory safety without garbage collection</li>
            <li>Compile-time race condition prevention</li>
          </ul>
        </div>
      </div>

      <div class="principle-item">
        <div class="principle-header">
          <h3>⚡ Fiber Concurrency</h3>
        </div>
        <div class="principle-content">
          <p>Isolated module instances per fiber eliminate data races. Scale to millions of concurrent tasks.</p>
          <ul>
            <li>Zero-cost fiber creation</li>
            <li>No shared mutable state</li>
            <li>Typed channels for communication</li>
          </ul>
        </div>
      </div>

      <div class="principle-item">
        <div class="principle-header">
          <h3>🔗 Rust Interoperability</h3>
        </div>
        <div class="principle-content">
          <p>Seamless integration with Rust for maximum performance. Type-safe FFI.</p>
          <ul>
            <li>Direct Rust async/await integration</li>
            <li>Zero-overhead interop</li>
            <li>Gradual migration support</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</section>

<section class="philosophy">
  <div class="container">
    <h2 class="section-title">Core Philosophy</h2>
    <div class="philosophy-grid">
      <div class="philosophy-card">
        <h3>Referential Transparency</h3>
        <p>Functions return the same output for the same input. Side effects are explicit.</p>
      </div>
      
      <div class="philosophy-card">
        <h3>Immutability by Default</h3>
        <p>Data immutable unless explicitly marked mutable. Safe sharing across fibers.</p>
      </div>
      
      <div class="philosophy-card">
        <h3>Explicit Error Handling</h3>
        <p>No exceptions or panics. All failures return Result types.</p>
      </div>
      
      <div class="philosophy-card">
        <h3>Zero-Cost Abstractions</h3>
        <p>High-level features compile to optimal machine code.</p>
      </div>
    </div>
  </div>
</section>

<section class="comparison">
  <div class="container">
    <h2 class="section-title">Key Differences</h2>
    <div class="comparison-table">
      <table>
        <thead>
          <tr>
            <th>Feature</th>
            <th>Traditional</th>
            <th>Osprey</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Error Handling</td>
            <td>Exceptions, panics</td>
            <td>Result types</td>
          </tr>
          <tr>
            <td>Null Safety</td>
            <td>Null pointer exceptions</td>
            <td>Option types only</td>
          </tr>
          <tr>
            <td>Concurrency</td>
            <td>Shared mutable state</td>
            <td>Fiber-isolated modules</td>
          </tr>
          <tr>
            <td>Type Safety</td>
            <td>Runtime type errors possible</td>
            <td>Compile-time prevention of type errors</td>
          </tr>
          <tr>
            <td>Memory Management</td>
            <td>Manual memory or garbage collection</td>
            <td>Memory safety with automatic reference counting</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</section>

<section class="innovation">
  <div class="container">
    <h2 class="section-title">Key Innovations</h2>
    <div class="innovations-grid">
      <div class="innovation-item">
        <h3>Fiber-Isolated Modules</h3>
        <p>
          Revolutionary approach to concurrency where each fiber gets its own isolated instances of modules. Eliminates data races while maintaining clean encapsulation.
        </p>
        <div class="innovation-code typewriter-enabled">
          <pre class="language-osprey"><code class="language-osprey">module Counter {
  let mut count = 0
  fn increment() = { count = count + 1; count }
}

// Each fiber has its own Counter instance
let fiber1 = Fiber<Int> { 
  computation: fn() => Counter.increment() 
}
let fiber2 = Fiber<Int> { 
  computation: fn() => Counter.increment() 
}

// Both return 1, not 1 and 2
await(fiber1)  // 1
await(fiber2)  // 1 (separate instance)</code></pre>
        </div>
      </div>

      <div class="innovation-item">
        <h3>Constraint-Based Types</h3>
        <p>
          Type constructors with WHERE constraints ensure data validity at creation time. Invalid data simply cannot exist in your program.
        </p>
        <div class="innovation-code typewriter-enabled">
          <pre class="language-osprey"><code class="language-osprey">type Person = {
  name: String where notEmpty(name),
  age: Int where between(age, 0, 150),
  email: String where validateEmail(email)
}

// Returns Result<Person, ValidationError>
let person = Person { 
  name: "Alice", 
  age: 25, 
  email: "alice@example.com" 
}</code></pre>
        </div>
      </div>

      <div class="innovation-item">
        <h3>Safe Arithmetic by Default</h3>
        <p>
          All arithmetic operations return Result types to handle overflow, underflow, and division by zero. Math errors are impossible to ignore.
        </p>
        <div class="innovation-code typewriter-enabled">
          <pre class="language-osprey"><code class="language-osprey">// Safe division that cannot panic
fn divide(a: Int, b: Int) -> Result<Int, MathError> = 
  match b {
    0 => Err { error: DivisionByZero }
    _ => Ok { value: a / b }
  }

// Must handle the result
match divide(a: 10, b: 0) {
  Ok { value } => print("Result: ${value}")
  Err { error } => print("Cannot divide by zero")
}</code></pre>
        </div>
      </div>
    </div>
  </div>
</section>

<section class="future">
  <div class="container">
    <h2 class="section-title">The Future</h2>
    <p style="text-align: center; font-size: var(--font-size-xl); margin-bottom: var(--space-12); color: var(--color-text-secondary); max-width: 800px; margin-left: auto; margin-right: auto;">
      We are just getting started. We're building a complete ecosystem for safe, concurrent programming.
    </p>
    
    <div class="features-grid">
      <div class="feature-card">
        <div class="feature-icon">🔗</div>
        <h3>Haskell Integration</h3>
        <p>Future interoperability with Haskell for formal verification and mathematical proofs of program correctness.</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">📦</div>
        <h3>Package Ecosystem</h3>
        <p>Growing library ecosystem with built-in safety guarantees and comprehensive documentation.</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">🛠️</div>
        <h3>Development Tooling</h3>
        <p>Advanced IDE support, intelligent debuggers, and powerful development tools for productive programming.</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">🎓</div>
        <h3>Education & Learning</h3>
        <p>Comprehensive resources to teach safe programming principles and functional programming concepts.</p>
      </div>
    </div>

    <div class="cta-section">
      <h3>Join the Revolution</h3>
      <p>Help us build the future of programming. Safe, fast, and elegant code for everyone.</p>
      <div class="cta-buttons">
        <a href="/playground/" class="btn btn-primary">Try Osprey Now</a>
        <a href="/docs/" class="btn btn-secondary">Read the Docs</a>
      </div>
    </div>
  </div>
</section> 