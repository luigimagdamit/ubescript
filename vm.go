package main

import "fmt"

var STACK_MAX int = 256

// Interpret Result enum
type InterpretResult int

const (
	INTERPRET_OK            InterpretResult = 0
	INTERPRET_COMPILE_ERROR InterpretResult = 1
	INTERPRET_RUNTIME_ERROR InterpretResult = 2
)

type VM struct {
	chunk *Chunk
	ip    uint8 // will point to memory location of an OpCode, which is usually just a byte

	stack    [4096]Value
	stackTop int
}

// Use a single global VM variable since we only need one
var vm *VM = new(VM)

func resetStack() {
	vm.stack = [4096]Value{}
	vm.stackTop = 0
}
func initVM() {
	resetStack()
}
func freeVM() {

}

func push(val Value) {
	vm.stack[vm.stackTop] = val
	vm.stackTop++
}
func pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]

}

// replaces ip++
// we want to return the byte at offset n, then increment the byte to n + 1
func READ_BYTE() uint8 {
	vm.ip++
	tmp := vm.chunk.Code[vm.ip-1]

	return tmp
}

func BINARY_OP(op func(b float64, a float64) float64) {
	var b float64 = pop()
	var a float64 = pop()
	fmt.Println(a, b)
	push(op(b, a))
}

func READ_CONSTANT() Value {
	var constantIndex uint8 = READ_BYTE()
	return vm.chunk.Constants.Values[constantIndex]
}

func READ_LONG_CONSTANT() Value {
	var indexBytes [4]uint8
	for i := 0; i < 4; i++ {
		indexBytes[i] = READ_BYTE()
	}

	var longConstantIndex uint32 = combineUInt8Array(indexBytes)
	return vm.chunk.Constants.Values[longConstantIndex]
}

// interpret() takes a chunk pointer as an input, runs it in the VM and returns the output
func interpret(source *string) InterpretResult {
	// compile(source)
	// return INTERPRET_OK

	var c *Chunk = new(Chunk)
	initChunk(c)

	// Fill the new chunk with the bytecode from compile()
	// retuurns false if there is a compile error
	if !compile(source, c) {
		freeChunk(c)
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = c
	vm.ip = 0

	var result InterpretResult = run() // MAKE SURE TO UNDO THIS!!!!
	//result := INTERPRET_OK
	freeChunk(c)
	return result
}

func run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("          offset 0->")
			for i := 0; i < vm.stackTop; i++ {
				fmt.Printf("[")
				printValue(vm.stack[i])
				fmt.Printf("]")
			}
			fmt.Println("<-offset", vm.stackTop)
			disassembleInstruction(vm.chunk, int(vm.ip))
		}

		var instruction uint8
		instruction = READ_BYTE()
		switch instruction {

		// RET OpCode
		case OP_RETURN:
			fmt.Printf("Printed ")
			printValue(pop())
			fmt.Println()
			return INTERPRET_OK
		case OP_ADD:
			BINARY_OP(add)
		case OP_SUBTRACT:
			BINARY_OP(sub)
		case OP_MULTIPLY:
			BINARY_OP(mul)
		case OP_DIVIDE:
			BINARY_OP(div)

		case OP_NEGATE:
			fmt.Println(vm.stack, vm.stackTop)
			res := pop()

			fmt.Println(vm.stack, vm.stackTop, res)
			push(-res)
			fmt.Println(vm.stack, vm.stackTop, res)
			break
		case OP_CONSTANT:
			var constant Value = READ_CONSTANT()
			push(constant) // producing a vlue, we push it onto the stack to be used
			fmt.Printf("\n")
			break
		case OP_CONSTANT_LONG:
			var longConstant Value = READ_LONG_CONSTANT()
			push(longConstant)
			fmt.Println()
			break
		}

	}
}
