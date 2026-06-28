# Osprey Examples

This directory contains comprehensive examples demonstrating all features of the Osprey programming language.

## Directory Structure

### `tested/` - Working Examples
Complete, validated examples that compile and run successfully. These demonstrate core language features and best practices.

### `broken/` - Incomplete Examples  
Examples that don't fully work yet or are under development. Used for testing compiler error handling.

### `failscompilation/` - Compilation Error Tests
Examples that should fail compilation with specific error messages. Used for testing compiler validation.

### `rust_integration/` - Rust Interop Examples
Examples demonstrating Rust-Osprey interoperability and extern function declarations.

## Example Categories

### Basic Language Features

#### `hello.osp`
**Hello World Example**
- Basic print functionality
- Simple function definitions
- Program structure

#### `basic.osp`
**Basic Syntax**
- Variable declarations
- Function definitions
- Basic operations

#### `function.osp`
**Function Definitions**
- Simple function declarations
- Function calls
- Return values

#### `simple.osp`
**Simple Variable Operations**
- Variable assignment
- Basic string operations

### Pattern Matching

#### `pattern_matching_basics.osp`
**Basic Pattern Matching**
- Match expressions with integers
- Wildcard patterns
- Function return types
- Pattern-based conditional logic

### Type System

#### `result_type_example.osp`
**Result Type for Error Handling**
- Safe division with error handling
- Result type pattern (Ok/Error)
- Pattern matching for error propagation

#### `simple_types.osp`
**Simple Type Operations**
- Basic type definitions
- Function parameter types

### Functional Programming

#### `functional_iterators.osp`
**Functional Iterator Operations**
- Range generation
- Pipe operator `|>`
- forEach, map, filter operations
- Higher-order functions

#### `comprehensive_iterators.osp`
**Advanced Iterator Operations**
- Fold operations
- Complex function chaining
- Functional composition

#### `functional_showcase.osp`
**Comprehensive Functional Features**
- Multiple functional programming patterns
- Complex data transformations

### String Interpolation

#### `interpolation_comprehensive.osp`
**String Interpolation Features**
- Variable interpolation
- Expression interpolation
- Multiple data types in strings

#### `interpolation_math.osp`
**Mathematical String Interpolation**
- Arithmetic expressions in strings
- Real-time calculation display

### Boolean Operations

#### `comprehensive_bool_test.osp`
**Boolean Logic Operations**
- AND, OR, NOT operations
- Boolean function definitions
- Logical expression evaluation

#### `full_bool_test.osp`
**Complete Boolean Testing**
- Extended boolean operations
- Complex boolean expressions

### Input/Output

#### `simple_input.osp`
**User Input Example**
- Reading user input with `input()`
- Processing user data
- Interactive programs

### Constraint Validation

#### `constraint_validation_test.osp`
**Type Constraints with WHERE**
- Record types with constraints
- Constraint validation
- Error handling for invalid data

#### `working_constraint_test.osp`
**Working Constraint Examples**
- Practical constraint usage
- Real-world validation scenarios

#### `proper_validation_test.osp`
**Proper Validation Patterns**
- Validation function patterns
- Error reporting
- Type safety

### Complex Applications

#### `adventure_game.osp`
**Text Adventure Game**
- Complete game implementation
- State management
- Interactive gameplay
- Complex pattern matching

#### `space_trader.osp`
**Space Trading Simulation**
- Economic simulation
- Resource management
- Complex game mechanics

#### `calculator_fixed.osp`
**Mathematical Calculator**
- Safe arithmetic operations
- Error handling
- User interface

### Testing and Validation

#### `comprehensive.osp`
**Comprehensive Language Demo**
- Multiple language features
- Integration testing
- Feature showcase

#### `working_basics.osp`
**Basic Language Validation**
- Core feature testing
- Syntax validation

#### `script_style_working.osp`
**Script-Style Programming**
- No main function required
- Direct execution
- Simple scripting patterns

## Test Categories

### Compilation Tests
- `comparison_test.osp` - Comparison operations
- `equality_test.osp` - Equality testing
- `modulo_test.osp` - Modulo operations
- `documentation_test.osp` - Documentation examples

### Iterator Tests
- `basic_iterator_test.osp` - Basic iterator functionality
- Various iterator operation tests

### Minimal Tests
- `minimal_test.osp` - Minimal working example

## Running Examples

All examples can be executed using the Osprey compiler:

```bash
# Run an example
osprey examples/tested/hello.osp --run

# View AST
osprey examples/tested/basic.osp --ast

# Generate LLVM IR
osprey examples/tested/function.osp --llvm

# Compile only
osprey examples/tested/pattern_matching_basics.osp --compile
```

## Key Features Demonstrated

### Language Features
- **Type Safety**: Strong typing with inference
- **Pattern Matching**: Exhaustive pattern matching
- **Functional Programming**: Higher-order functions, immutability
- **String Interpolation**: Expression embedding in strings
- **Error Handling**: Result types instead of exceptions

### Programming Patterns
- **Script-style execution**: No main function required
- **Functional composition**: Pipe operator and function chaining
- **Safe arithmetic**: Result types for all operations
- **Constraint validation**: WHERE clauses for type safety
- **Interactive programming**: Input/output operations

These examples demonstrate Osprey's core philosophy of safety, expressiveness, and functional programming while maintaining simplicity and clarity. 