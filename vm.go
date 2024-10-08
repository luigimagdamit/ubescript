package main

import "fmt"

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
}

// Use a single global VM variable since we only need one
var vm *VM = new(VM)

func initVM() {

}
func freeVM() {

}

// replaces ip++
// we want to return the byte at offset n, then increment the byte to n + 1
func READ_BYTE() uint8 {
	tmp := vm.chunk.Code[vm.ip]
	vm.ip++
	return tmp
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
func interpret(c *Chunk) InterpretResult {
	vm.chunk = c
	vm.ip = 0 // should just hold the offset of the current byte
	return 0
}

func run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			disassembleInstruction(vm.chunk, int(vm.ip))
		}

		var instruction uint8
		instruction = READ_BYTE()
		switch instruction {

		// RET OpCode
		case OP_RETURN:
			fmt.Println("INTERPRET_OK")
			return INTERPRET_OK

		case OP_CONSTANT:
			var constant Value = READ_CONSTANT()
			printValue(constant)
			fmt.Printf("\n")
			break
		case OP_CONSTANT_LONG:
			var longConstant Value = READ_LONG_CONSTANT()
			printValue(longConstant)
			fmt.Println()
			break
		}

	}
}
