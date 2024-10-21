# `ubescript is a small, portable, dynamically typed scripting language.`

 **Note: This project is still currently under construction**

## Intro
---
ubescript is a gradually typed, high level, compiled programming language. It focuses on maintaining a balance between readability and expressiveness. It currently compiles to a bytecode format and is run from a custom made virtual machine architecture.

ubescript pulls inspiration from the simplicity of Go, while developing a strong typing system found in Rust, and stealing from the convienience of Python semantics 

## Planned Features
---
- Target Compile to C down to machine code
- Gradual typing
- Compile to register based IR > LLVM
- Record Types
- Loop interchange optimization
- Some() Construct
- Further OpCode optimizations

## Examples

### Hello World.ube
```rust
println("hello world")

let breakfast str = "beignets"
let beverage str = "cafe au lait"

println(breakfast)

// "beignets"
breakfast = breakfast + " with " + beverage

println(breakfast)
// "beignets with au lait"
```

```
hello world
beignets
beignets with cafe au lait
```

### Strings
```rust
println(len("hello"))
// prints 5

println("ube" + "script")
//print 'ubescript'

println("abc" == "abc")
// print 'true'
```

### Numerical Operations
```rust
println(1 + 2)
// prints 3
println(1 * 2 * 3)
// prints 6
1..400
// places 1-400 on the VM stack
```


### Variables
```rust
let greeting str = "hello world"
let n int = len(greeting)

let equals12 = n == 12

println(greeting + " has the length: ") 
println(n) 

println("and its equal to 12: ")
println(equals12)
```


### Fibonacci

```rust
let n1, n2 = 0, 1
let tmp int

let i, n = 0, 10

while i < n  {
    println(n1)
    tmp = n2
    n2 = n1 + n2
    n1 = tmp
    i = i + 1
}
```

### Scoping
```rust
let x int = 24 / 2
let y int = 4

println(x / y)

{
    let y int = 2
    println(x / y)
}

println(x / y)
```