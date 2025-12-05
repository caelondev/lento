# Lento

A lightweight, tree-walk interpreter built with Go, designed for simplicity and ease of use.

## Overview

Lento is a modern scripting language with clean syntax, perfect for automation tasks and learning language implementation. Built from scratch in Go, it offers fast execution with a straightforward learning curve.

## Getting Started

### Prerequisites

- Go 1.25.4 or higher

### Installation

```bash
go install github.com/caelondev/lento@latest
```

## Language Guide

### Data Types

Lento supports five core data types:

**Number** - Integers and floating-point values
```lento
42
3.14
-17.5
```

**Boolean** - Logical true/false values
```lento
true
false
```

**String** - Text in quotes (single, double, or backticks for multiline)
```lento
"Hello, World!"
'Single quotes work too'
`Multiline
strings
supported`
```

**Array** - Ordered lists that can hold any type, including nested arrays
```lento
[1, 2, 3]
["Apple", "Orange", "Banana"]
[0, "Hello, World!", true, [x, [y]]]
```

**Object** - Key-value pairs for structured data
```lento
{
  firstName: "Foo",
  lastName: "Bar",
  age: 22
}
```

### Variables

Declare variables using `var` for mutable values or `const` for immutable constants:

```lento
var foo = "bar";
const baz = -10;
var uninitialized;  // defaults to nil
```

Non-constant variables can be reassigned:

```lento
var foo = "Bar";
print(foo);  // Outputs "Bar"

foo = "Hello, World!";
print(foo);  // Outputs "Hello, World!"
```

### Working with Arrays

You can print array values by either accessing a specific index or printing the entire array:

```lento
print([])  // Prints an empty array

var fruits = ["Apple", "Orange", "Banana"]

print(fruits[0])         // Prints "Apple"
print(fruits)            // Prints the whole array
print([true, false][0])  // Prints true (not recommended)
```

Arrays can be modified by reassigning values at specific indices:

```lento
var bar = ["Foo", "Baz"]
bar[1] = "Bar"
print(bar[1])  // Outputs "Bar"
```

### Working with Objects

Objects can be printed directly or you can access specific properties using bracket notation or dot notation:

```lento
var person = { 
  name: "Bob",
  age: 25,
  location: {
    street: "FooBar street",
    continent: "Asia"
  }
}

print({})  // Prints an empty object

// Bracket notation
print(person[name])              // Prints "Bob"
print(person[location])          // Prints the location object
print(person[location][street])  // Prints "FooBar street"

// Dot notation
print(person.name)               // Prints "Bob"
print(person.location.continent) // Prints "Asia"
```

Objects can be modified using either bracket or dot notation:

```lento
// Bracket notation
person[location][continent] = "Europe"

// Dot notation
person.location.continent = "Europe"

print(person.location.continent)  // Prints "Europe"
```

### Operators

**Arithmetic**
```lento
+      // Addition
-      // Subtraction
*      // Multiplication
/      // Division
%      // Modulo
```

**Comparison**
```lento
<      // Less than
<=     // Less than or equal to
>      // Greater than
>=     // Greater than or equal to
==     // Equal to
!=     // Not equal to
```

**Logical**
```lento
and    // Logical AND
or     // Logical OR
not    // Logical NOT
```

**Assignment**
```lento
=      // Assignment
+=     // Add and assign
-=     // Subtract and assign
*=     // Multiply and assign
/=     // Divide and assign
%=     // Modulo and assign
```

### Control Flow

Lento supports standard if-else statements with flexible syntax:

```lento
// Simple if
if (x > 10) {
  x = x + 1
}

// If-else
if (x > 10) {
  x = x * 2
} else {
  x = 5
}

// Else-if chains
if (x > 100) {
  x = 100
} else if (x > 10) {
  x = x + 5
} else {
  x = 0
}

// Single-line syntax (parentheses optional with braces)
if true { x = 42; }
if (x > 0) x = x - 1;
```

### Functions

Define functions using the `fn` keyword:

```lento
fn greet(name) {
  print("Hello, " + name + "!")
}

fn add(a, b) {
  a + b  // Last expression is returned
}
```

**Note**: Parameters are pass-by-value, so modifying them won't affect the original arguments.

Functions support closures and capture their surrounding environment:

```lento
var x = 10

fn makeAdder() {
  x  // Captures x from outer scope
}
```

Call functions like you'd expect:

```lento
print("Hello, World!")  // Built-in function
greet("Alice")          // User-defined function
var sum = add(5, 3)
```

### Loops

#### While loops

While loops can be written with this syntax

```lento
var x = 0;

while (x < 100) { // This loops over-and-over again until
  x = x+10;       // x is >= 100
  print(x) // Outputs (x + 10)
}

print(x) // Outputs 100
```

## Interactive REPL

Lento includes an interactive REPL for quick experimentation:

```bash
$ lento
>> var x = 42
42
>> x * 2
84
>> fn greet(name) "Hello, " + name
[ greet function ]
>> greet("World")
Hello, World!
```

Exit the REPL with `*exit` or Ctrl+C.

## Performance

Lento is designed for speed. Even as a tree-walk interpreter, it executes scripts in microseconds thanks to Go's efficient runtime.

## Contributing

Contributions are welcome! Feel free to:
- Report bugs or request features via [issues](https://github.com/caelondev/lento/issues)
- Submit pull requests with improvements
- Share your Lento scripts and projects

## License

[MIT](./LICENSE)

## Credits

Created and maintained by [caelondev](https://github.com/caelondev)

---

**Note**: Lento is under active development. Features and syntax may evolve as the project matures.
