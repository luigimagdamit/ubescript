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

### Intro
```rust
// basics
let name str = "ube"
let age int = 400

println("My name is " + name + "!")

let favorite = "saffron and rose"

if name == "ube" and favorite == "saffron and rose" {
    println("persian icecream")
} else {
    println(":(")
}

if age == 400 println("i belong in a museum")
else println("the museum belongs to me")

// variables
let cat str = "meow"
let dog str = "woof"
let isDog = true

if isDog println(dog); else println(cat)

let mouse1, mouse2, mouse3 (str) = "squeek", "squeek", "squeek"

let pets int = 3
println("Cat says " + cat)
println("Dog says " + dog)

mouse1 = "deadmau5"
println("Mouse says " + mouse1)

// logic
let name = "beach"
let sleepy = false

if !sleepy println("Let's drive!")
else println("Let's go home...")

// control flow
for 0..3 {
    println(dog)
}
let dogsLoose = 0

while dogsLoose < 10 {
    dogsLoose = dogsLoose + 1
    println(dog)
}
if dogsLoose == 10 {
    println("Who let the dogs out?")
}
```
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
let n1, n2, tmp (int) = 0, 1, 0
let runs int = 10

for 0..runs {
    println(n1)
    tmp = n2
    n2 = n1 + n2
    n1 = tmp
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