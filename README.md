# `ubescript`

`ubescript` is a small, gradually typed programming language built around the semantics and ease of readability of Python and the simplicity of Go. It is yet another C-inspired language, with a focus on conciseness, convenience, and expressivity.

 **Note: This project is still currently under construction**

Many thanks to Bob Nystrom and his `lox` language, serving as the foundational basis for this project.
## Intro
---

```rust
let cat str = "meow" // optional type annotations
let dog str = "woof"
let isDog = true

if isDog println(dog); else println(cat) // python style one line if-else

let mouse1, mouse2, mouse3 (str) = "squeek", "squeek", "squeek" // multiple assignment 

let miceCount int = 3
println("Cat says " + cat) // "meow"
println("Dog says " + dog) // "woof"
for 0..miceCount println("Mouse says " + mouse3) // "squeek", "squeek", "squeek"

mouse1 = "deadmau5"
println("Mouse also says " + mouse1) // "deadmau5"
```
## Planned Features
---
- Target Compile to C down to machine code
- Record Types
- Some() Construct
- Many more...

## Examples

### Intro
```rust
// basics

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

echo "hello world!" // prints 'hello world' without \n
```

### Numerical Operations
```rust
println(1 + 2)
// prints 3
println(1 * 2 * 3)
// prints 6
```


### Variables
```rust
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

```

### Fizzbuzz
```rust
let i int = 0

for 0..30 {
    if (i % 5 == 0 and i % 3 == 0) println("fizzbuzz")
    else if (i % 3 == 0) println("fizz")
    else if (i % 5 == 0) println("buzz")
    else println(i)
    i = i + 1
}
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