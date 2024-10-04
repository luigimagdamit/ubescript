package main

import "fmt"

// opcode enum

const (
	OP_RETURN = iota
)

type Chunk struct {
	Count    int
	Capacity int

	Code []uint8
}

// tbh this is isn't really needed i'll just  do it since we likely will need it later
func initChunk(c *Chunk) {
	c.Capacity = 0
	c.Count = 0
	c.Code = []uint8{} // initialize an empty slice
}

func writeChunk(c *Chunk, inst uint8) {
	c.Code = append(c.Code, inst)
	c.Capacity++
	c.Count++
}

func freeChunk(c *Chunk) {
	c.Code = []uint8{}
	initChunk(c)

}

func disassembleChunk(c *Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for offset := 0; offset < c.Count; {
		offset = disassembleInstruction(c, offset)
	}
}
func disassembleInstruction(c *Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	var inst uint8 = c.Code[offset]
	switch inst {
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
