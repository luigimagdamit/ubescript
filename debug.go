package main

import "fmt"

func disassembleChunk(c *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	fmt.Println("OFFS  LINE\tOPCCODE")
	for offset := 0; offset < c.Count; {
		offset = disassembleInstruction(c, offset)
	}
}
func disassembleInstruction(c *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Println("   | ") // this is done for any other instructions that come from the same source line
	} else {
		fmt.Printf("%4d\t", c.Lines[offset])
	}

	var inst uint8 = c.Code[offset]
	switch inst {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_CONSTANT_LONG:
		return constantInstructionLong("OP_CONSTANT_LONG", c, offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown OpCode %d at offset %04d\n", inst, offset)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, c *Chunk, offset int) int {
	var constant uint8 = c.Code[offset+1] // obtain the operand
	fmt.Printf("%-16s C.Index %4d ' Value: ", name, constant)

	printValue(c.Constants.Values[constant]) // print the actual value from within the constant pool
	fmt.Printf("\n")
	return offset + 2
}

func constantInstructionLong(name string, c *Chunk, offset int) int {
	//var constant uint8 = c.Code[offset+1] // obtain the operand

	var arr [4]uint8 = [4]uint8{}
	for i := 0; i < 4; i++ {
		arr[i] = c.Code[1+offset+i]
	}
	long_con := combineUInt8Array(arr)
	fmt.Printf("%s C.Index: %4d ' Value: ", name, long_con)

	printValue(c.Constants.Values[long_con]) // print the actual value from within the constant pool
	fmt.Printf("\n")
	return offset + 5
}
