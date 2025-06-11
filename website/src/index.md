---
layout: base.njk
title: "Osprey Programming Language"
description: "A modern functional programming oriented language designed for elegance, safety, and performance."
---

<section class="hero">
  <div class="container">
    <div class="hero-content">
      <h1 class="hero-title animate-on-scroll">
        Build the Future with
        Osprey
      </h1>
      <p class="hero-subtitle animate-on-scroll">
        A modern functional programming oriented language designed for elegance, safety, and performance.
      </p>
      <div class="hero-actions animate-on-scroll">
        <a href="/playground/" class="btn btn-primary btn-lg">
          Try Osprey Online
        </a>
        <a href="/docs/" class="btn btn-secondary btn-lg">
          View Documentation
        </a>
      </div>
    </div>
    
    <div class="hero-code animate-on-scroll parallax typewriter-enabled" data-speed="0.3">
      <div class="code-example">
        <pre class="language-osprey"><code class="language-osprey">// Simple, clean function definitions
fn double(n: int) -> int = n * 2
fn greet(name: string) -> string = "Hello " + name

// String interpolation that works
let x = 42
let name = "Alice"
print("x = ${x}")
print("Greeting: ${greet(name)}")

// Pattern matching on values
let result = match x {
  42 => "The answer!"
  0 => "Zero"
  _ => "Something else"
}
print("Result: ${result}")</code></pre>
      </div>
    </div>
  </div>
</section>

<section class="features">
  <div class="container">
    <h2 class="section-title animate-on-scroll">Why Choose Osprey?</h2>
    <div class="features-grid">
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">🎯</div>
        <h3>Clear Function Definitions</h3>
        <p>Explicit type annotations and expression-bodied functions create self-documenting code.</p>
      </div>
      
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">🔀</div>
        <h3>Pattern Matching</h3>
        <p>Elegant pattern matching on values with exhaustiveness checking for safe code.</p>
      </div>
      
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">✨</div>
        <h3>String Interpolation</h3>
        <p>Built-in string interpolation with expression support for readable output formatting.</p>
      </div>
      
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">🔗</div>
        <h3>Functional Programming</h3>
        <p>Pipe operators and functional iterators for elegant data processing pipelines.</p>
      </div>
      
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">🛡️</div>
        <h3>Type Safety</h3>
        <p>Strong static typing prevents runtime errors and catches issues at compile time.</p>
      </div>
      
      <div class="feature-card animate-on-scroll">
        <div class="feature-icon">⚡</div>
        <h3>Fast Compilation</h3>
        <p>Quick compilation cycles for rapid development and testing of your programs.</p>
      </div>
    </div>
  </div>
</section>

<section class="code-showcase typewriter-enabled">
  <div class="container">
    <h2 class="section-title animate-on-scroll">See Osprey in Action</h2>
    
    <div class="showcase-grid">
      <div class="showcase-item animate-on-scroll">
        <h3>Clean Functions</h3>
        <div class="code-example">
          <pre class="language-osprey"><code class="language-osprey">// Expression-bodied functions
fn analyzeNumber(n: int) -> string = match n {
  0 => "Zero"
  1 => "One"
  42 => "The answer to everything!"
  _ => "Some other number"
}

// Simple calculations
fn double(x: int) -> int = x * 2
fn square(x: int) -> int = x * x

let result = double(21)
print("double(21) = ${result}")

let analysis = analyzeNumber(42)
print("Analysis: ${analysis}")</code></pre>
        </div>
      </div>
      
      <div class="showcase-item animate-on-scroll">
        <h3>String Interpolation</h3>
        <div class="code-example">
          <pre class="language-osprey"><code class="language-osprey">// Variables and expressions in strings
let name = "Alice"
let age = 25
let score = 95

print("Hello ${name}!")
print("You are ${age} years old")
print("Next year you'll be ${age + 1}")

// Calculate and interpolate
let doubled = score * 2
print("Double your score: ${doubled}")

// Multiple expressions
print("${name} (${age}) scored ${score}/100")</code></pre>
        </div>
      </div>
      
      <div class="showcase-item animate-on-scroll">
        <h3>Functional Pipelines</h3>
        <div class="code-example">
          <pre class="language-osprey"><code class="language-osprey">// Pipe operator for clean data flow
fn double(x: int) -> int = x * 2
fn square(x: int) -> int = x * x

// Single value transformations
5 |> double |> print
3 |> square |> print

// Range operations
range(1, 5) |> forEach(print)

// Chained operations
let result = 2 |> double |> square
print("Result: ${result}")</code></pre>
        </div>
      </div>
      
      <div class="showcase-item animate-on-scroll">
        <h3>Pattern Matching</h3>
        <div class="code-example">
          <pre class="language-osprey"><code class="language-osprey">// Match on different values
fn getCategory(score: int) -> string = match score {
  100 => "Perfect!"
  95 => "Excellent"
  85 => "Very Good"
  75 => "Good"
  _ => "Needs Improvement"
}

// Simple boolean logic through matching
fn isEven(n: int) -> int = match n {
  0 => 1
  2 => 1
  4 => 1
  _ => 0
}

print("Score 95: ${getCategory(95)}")
print("4 is even: ${isEven(4)}")</code></pre>
        </div>
      </div>
    </div>
  </div>
</section>

<section class="philosophy">
  <div class="container">
    <h2 class="section-title animate-on-scroll">Core Principles</h2>
    <div class="philosophy-grid">
      <div class="philosophy-card animate-on-scroll">
        <h3>Simple & Elegant</h3>
        <p>Clean syntax that reads naturally. Expression-bodied functions and minimal ceremony.</p>
      </div>
      <div class="philosophy-card animate-on-scroll">
        <h3>Type Safe</h3>
        <p>Strong static typing catches errors at compile time. Explicit types prevent confusion.</p>
      </div>
      <div class="philosophy-card animate-on-scroll">
        <h3>Functional First</h3>
        <p>Pattern matching, pipe operators, and functional iterators for elegant data processing.</p>
      </div>
      <div class="philosophy-card animate-on-scroll">
        <h3>Fast Development</h3>
        <p>Quick compilation and immediate feedback. Get from idea to working code fast.</p>
      </div>
    </div>
  </div>
</section>

<section class="getting-started">
  <div class="container">
    <div class="cta-content animate-on-scroll">
      <h2>Ready to Get Started?</h2>
      <p>Experience clean, functional programming with strong typing.</p>
      <div class="cta-actions">
        <a href="/playground/" class="btn btn-primary btn-lg">
          Try Osprey Now
        </a>
        <a href="/docs/" class="btn btn-outline btn-lg">
          Read the Docs
        </a>
      </div>
    </div>
  </div>
</section> 