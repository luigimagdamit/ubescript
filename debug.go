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

	// if offset > 0 {
	// 	curLine := c.Lines[offset]
	// 	prevLine := c.Lines[offset-1]
	// 	fmt.Println(curLine == prevLine, []int{curLine, prevLine}, []int{offset, offset - 1})

	// }

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		//fmt.Println("   | ") // this is done for any other instructions that come from the same source line
	} else {
		fmt.Printf("%4d\t", c.Lines[offset])
	}

	var inst uint8 = c.Code[offset]
	switch inst {
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_CONSTANT_LONG:
		return constantInstructionLong("OP_CONSTANT_LONG", c, offset)
	case OP_NIL:
		return simpleInstruction("OP_NIL", offset)
	case OP_TRUE:
		return simpleInstruction("OP_TRUE", offset)
	case OP_FALSE:
		return simpleInstruction("OP_FALSE", offset)
	case OP_POP:
		return simpleInstruction("OP_POP", offset)
	case OP_SET_LOCAL:
		return simpleInstruction("OP_SET_LOCAL", offset)
	case OP_GET_LOCAL:
		return simpleInstruction("OP_GET_LOCAL", offset)
	case OP_DEFINE_GLOBAL:
		return simpleInstruction("OP_DEFINE_GLOBAL", offset)
	case OP_SET_GLOBAL:
		return simpleInstruction("OP_SET_GLOBAL", offset)
	case OP_GET_GLOBAL:
		return simpleInstruction("OP_GET_GLOBAL", offset)
	case OP_EQUAL:
		return simpleInstruction("OP_EQUAL", offset)
	case OP_GREATER:
		return simpleInstruction("OP_GREATER", offset)
	case OP_LESS:
		return simpleInstruction("OP_LESS", offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	case OP_DOTDOT:
		return simpleInstruction("OP_DOTDOT", offset)
	case OP_LEN:
		return simpleInstruction("OP_LEN", offset)
	case OP_SHOW:
		return simpleInstruction("OP_SHOW", offset)
	case OP_NOT:
		return simpleInstruction("OP_NOT", offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_JUMP_IF_FALSE:
		return jumpInstruction("OP_JUMP_IF_FALSE", 1, c, offset)
	case OP_LOOP:
		return jumpInstruction("OP_LOOP", -1, c, offset)
	case OP_JUMP:
		return jumpInstruction("OP_JUMP", 1, c, offset)
	default:
		fmt.Printf("Unknown OpCode %d at offset %04d\n", inst, offset)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1 // Thiinggss  liikeeee  RET
}

func constantInstruction(name string, c *Chunk, offset int) int {
	var constant uint8 = c.Code[offset+1] // obtain the operand
	fmt.Printf("%-16s (C.Index %4d ' Value: ", name, constant)

	printValue(c.Constants.Values[constant]) // print the actual value from within the constant pool
	fmt.Printf(")\n")
	return offset + 2 // since the original constant implementation is [01] [02] = [OP_CONSTANT][CONST_INDEX]
}

func constantInstructionLong(name string, c *Chunk, offset int) int {
	//var constant uint8 = c.Code[offset+1] // obtain the operand

	var arr [4]uint8 = [4]uint8{}
	for i := 0; i < 4; i++ {
		arr[i] = c.Code[1+offset+i]
	}
	long_con := combineUInt8Array(arr)
	fmt.Printf("%s (C.Index: %4d ' Value: ", name, long_con)

	printValue(c.Constants.Values[long_con]) // print the actual value from within the constant pool
	fmt.Printf(")\n")
	return offset + 5
}
func jumpInstruction(name string, sign int, chunk *Chunk, offset int) int {
	jump := uint16(uint16(chunk.Code[offset+1]) << 8)
	jump |= uint16(uint16(chunk.Code[offset+2]))

	return offset + 3
}
