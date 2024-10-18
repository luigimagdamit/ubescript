# `ubescript is a small, portable, dynamically typed scripting language.`

 **Note: This project is still currently under construction**

## Intro
---
ubescript is a gradually typed, high level, compiled programming language. It focuses on maintaining a balance between readability and expressiveness. It currently compiles to a bytecode format and is run from a custom made virtual machine architecture.

ubescript pulls inspiration from the simplicity of Go, while developing a strong typing system found in Rust.

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
---
### Hello World.ube
```rust
print("hello world")
```

```
hello world
```

### Strings
```rust
print(len("hello"))
// prints 5

print("ube" + "script")
//print 'ubescript'

print("abc" == "abc")
// print 'true'
```

### Numerical Operations
```rust
print(1 + 2)
// prints 3
print(1 * 2 * 3)
// prints 6
1..400
// places 1-400 on the VM stack
```


### Variables
```rust
let greeting str= "hello world"
let length int= len(greeting)

print("The length of ")
print(greeting + " is")
print(length)
```