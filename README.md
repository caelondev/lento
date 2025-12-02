# Lento

A lightweight, tree-walk interpreter built with Go, designed for simplicity and ease of use.

## Overview

Lento is a modern interpreter that provides a clean syntax for scripting and automation tasks. Built from the ground up in Go, it offers fast execution and a straightforward learning curve.

## Getting Started

### Prerequisites

- Go version 1.25.4 or higher

### Installation

Install Lento using Go's package manager:

```bash
go install github.com/caelondev/lento@latest
```

## Language Reference

### Data Types

Lento supports three core data types:

- **Number**: Integer and floating-point values
  ```
  42
  3.14
  -17.5
  ```

- **Boolean**: Logical true/false values
  ```
  true
  false
  ```

- **String**: Text enclosed in double quotes
  ```
  "Hello, World!"
  "Lento"
  ```

### Operators

#### Logical Operators

```
and    // Logical AND
or     // Logical OR
not    // Logical NOT
```

#### Comparison Operators

```
<      // Less than
<=     // Less than or equal to
>      // Greater than
>=     // Greater than or equal to
==     // Equal to
!=     // Not equal to
```

### Variables

#### Declaration

Declare variables using `var` for mutable values or `const` for immutable constants:

```
var foo = "bar";
const baz = -10;
```

#### Reassignment

Mutable variables can be reassigned:

```
foo = "Hello, World!";
```

Note: Constants cannot be reassigned after declaration.

## Examples

```
// Variable declaration and manipulation
var counter = 0;
const maxValue = 100;

// Logical operations
var isValid = true and (counter < maxValue);

// Comparison operations
var isComplete = counter >= maxValue;

// String manipulation
var greeting = "Hello";
greeting = greeting + ", World!";
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests on the [Lento repository](https://github.com/caelondev/lento).

## License

[MIT](./LICENSE)

## Credits

Created and maintained by [caelondev](https://github.com/caelondev)

---

**Note**: Lento is under active development. Features and syntax may evolve as the project matures.
