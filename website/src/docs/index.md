---
layout: page
title: "API Reference - Osprey Programming Language"
description: "Complete reference documentation for the Osprey programming language"
---

## Quick Navigation

- [Functions](functions/) - Built-in functions for I/O, iteration, and data transformation
- [Types](types/) - Built-in data types (Int, String, Bool, Any)
- [Operators](operators/) - Arithmetic, comparison, and logical operators
- [Keywords](keywords/) - Language keywords (fn, let, type, match, import)

## Function Reference

| Function | Description |
|----------|-------------|
| [filter](functions/filter/) | Filters elements in an iterator based on a predicate function. |
| [fold](functions/fold/) | Reduces an iterator to a single value using an accumulator function. |
| [forEach](functions/foreach/) | Applies a function to each element in an iterator. |
| [input](functions/input/) | Reads an integer from the user's input. |
| [map](functions/map/) | Transforms each element in an iterator using a function, returning a new iterator. |
| [print](functions/print/) | Prints a value to the console. Automatically converts the value to a string representation. |
| [range](functions/range/) | Creates an iterator that generates numbers from start to end (exclusive). |
| [toString](functions/tostring/) | Converts a value to its string representation. |

## Type Reference

| Type | Description |
|------|-------------|
| [Any](types/any/) | A type that can represent any value. Useful for generic programming but should be used carefully as it bypasses type checking. |
| [Bool](types/bool/) | A boolean type that can be either true or false. Used for logical operations and conditionals. |
| [Int](types/int/) | A 64-bit signed integer type. Can represent whole numbers from -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807. |
| [String](types/string/) | A sequence of characters representing text. Supports string interpolation and escape sequences. |

## Operator Reference

| Operator | Name | Description |
|----------|------|-------------|
| [!=](operators/not-equal/) | Inequality | Compares two values for inequality. |
| [%](operators/modulo/) | Modulo | Returns the remainder of dividing the first number by the second. |
| [*](operators/multiply/) | Multiplication | Multiplies two numbers. |
| [+](operators/plus/) | Addition | Adds two numbers together. |
| [-](operators/minus/) | Subtraction | Subtracts the second number from the first. |
| [/](operators/divide/) | Division | Divides the first number by the second. |
| [<](operators/less-than/) | Less Than | Checks if the first value is less than the second. |
| [<=](operators/less-equal/) | Less Than or Equal | Checks if the first value is less than or equal to the second. |
| [==](operators/equal/) | Equality | Compares two values for equality. |
| [>](operators/greater-than/) | Greater Than | Checks if the first value is greater than the second. |
| [>=](operators/greater-equal/) | Greater Than or Equal | Checks if the first value is greater than or equal to the second. |
| [|>](operators/pipe-operator/) | Pipe Operator | Takes the result of the left expression and passes it as the first argument to the right function. Enables functional programming and method chaining. |

## Keyword Reference

| Keyword | Description |
|---------|-------------|
| [false](keywords/false/) | Boolean literal representing the logical value false. |
| [fn](keywords/fn/) | Function declaration keyword. Used to define functions with parameters and return types. |
| [import](keywords/import/) | Import declaration keyword. Used to bring modules and their exports into the current scope. |
| [let](keywords/let/) | Variable declaration keyword. Used to bind values to identifiers. Variables are immutable by default in Osprey. |
| [match](keywords/match/) | Pattern matching expression. Used for destructuring values and control flow based on patterns. |
| [true](keywords/true/) | Boolean literal representing the logical value true. |
| [type](keywords/type/) | Type declaration keyword. Used to define custom types and type aliases. |

