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

  ```lento
  42
  3.14
  -17.5
  ```

- **Boolean**: Logical true/false values

  ```lento
  true
  false
  ```

- **String**: Text enclosed in quotes (single, double, or backticks for multiline)

  ```lento
  "Hello, World!"
  'Single quotes work too'
  `Multiline
  strings
  supported`
  ```

- Arrays: 1D list of datas (can be string, number, boolean, and even another array)

```lento
  [0, "Hello, World!", true, [x, [y]]]
```

- Object: List of key-value pair

```lento
  {
    firstName: "Foo",
    lastName: "Bar",
    age: 22,
  };
```

### Variables

#### Declaration

Declare variables using `var` for mutable values or `const` for immutable constants:

```lento
var foo = "bar";
const baz = -10;
var uninitialized;  // defaults to nil
```

### Printing values

Values can be printed in two ways (or more)

#### Arrays

You can print array values by either accessing the array index or printing the array itself

```lento
  print([]) // Prints an empty array

  var fruits = ["Apple", "Orange", "Banana"]

  print(fruits[0])        // Prints "Apple"
  print(fruits)           // Prints the whole array
  print([true, false][0]) // Prints `true` (not recommended to print arrays this way)
```

#### Objects

Same as the array, Objects can also be printed by either printing it directly or accessing the key of the object

```lento
  var object = {
    name: "Bob",
    age: 25,
    location: {
      street: "FooBar street",
      continent: "Asia",
    }
  }
  print({}) // Prints an empty object

  print(object[name])             // Prints "Bob"
  print(object[location])         // Prints the `location` object
  print(object[location][street]) // Prints "FooBar street"
```

### Reassignment

#### Variables

Non-constant variables can be reassigned:

```lento
  var foo = "Bar";
  print(foo); // Outputs "Bar"
  foo = "Hello, World!";
  print(foo); // Outputs "Hello, World!"
```

#### Arrays

```lento
  var bar = ["Foo", "Baz"];
  bar[1] = "Bar";

  print(bar[1]) // Outputs "Bar"
```

### Operators

#### Arithmetic Operators

```lento
  +      // Addition
  -      // Subtraction
  *      // Multiplication
  /      // Division
  %      // Modulo
```

#### Comparison Operators

```lento
  <      // Less than
  <=     // Less than or equal to
  >      // Greater than
  >=     // Greater than or equal to
  ==     // Equal to
  !=     // Not equal to
```

#### Logical Operators

```lento
  and    // Logical AND
  or     // Logical OR
  not    // Logical NOT
```

#### Assignment Operators

```lento
  =      // Assignment
  +=     // Add and assign
  -=     // Subtract and assign
  *=     // Multiply and assign
  /=     // Divide and assign
  %=     // Modulo and assign
```

### Control Flow

#### If Statements

```lento
  // Simple if
  if (x > 10) {
      x = x + 1;  // x is 11 now
  }

  // If-else
  if (x > 10) {
      x = x * 2;  // x is doubled
  } else {
      x = 5;      // x is 5 now
  }

  // Else-if chains
  if (x > 100) {
      x = 100;    // cap at 100
  } else if (x > 10) {
      x = x + 5;  // add 5
  } else {
      x = 0;      // reset to 0
  }

  // Single-line syntax (parentheses optional with braces)
  if true { x = 42; }
  if (x > 0) x = x - 1;
```

### Functions

#### Function Declaration

Define functions using the `fn` keyword:

```lento
  fn recursivePrint(y) { // NOTE: `y` is a pass-by-value parameter
    var z = y+1;         // Modifying it won't affect the argument
    print(z);  b         // of the caller
    recursivePrint(z);
  }
```

Calling a function

```lento
  print("Hello, World!"); // Calls the `print` native function
  add(x, y);
```

## Examples

```lento
  // Variable declaration and manipulation
  var counter = 0;
  const maxValue = 100;

  // Arithmetic operations
  var result = (10 + 5) * 2;  // result is 30
  var remainder = 17 % 5;      // remainder is 2

  // Logical operations
  var isValid = true and (counter < maxValue);

  // Comparison operations
  var isComplete = counter >= maxValue;

  // String manipulation
  var greeting = "Hello";
  greeting = greeting + ", World!";  // greeting is "Hello, World!" now

  // Control flow
  if (counter < maxValue) {
      counter = counter + 1;  // counter is 1 now
  }

  // Functions
  fn printHelloWorld() { // Why not?
    print("Hello, World!");
  }

  // Function with closures
  var x = 10;
  fn makeAdder() {
      x;  // Captures x from outer scope
  }
```

## REPL

Lento includes an interactive REPL (Read-Eval-Print Loop) for quick experimentation:

```bash
$ lento
>> var x = 42;
42
>> x * 2;
84
>> fn greet(name) name = "Hello, " + name;
[ greet function ]
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests on the [Lento repository](https://github.com/caelondev/lento).

## License

[MIT](./LICENSE)

## Credits

Created and maintained by [caelondev](https://github.com/caelondev)

---

**Note**: Lento is under active development. Features and syntax may evolve as the project matures.
