package main

import "fmt"

func main() {
	splitUInt32(0xdeadbeef)
	c := new(Chunk)
	initChunk(c)

	for i := 0; i < 10000; i++ {
		//writeChunk(c, OP_CONSTANT_LONG, 42)
		writeConstant(c, float64(2*i), 42)

	}

	fmt.Println(c.Constants.Values)
	// constant := addContant(c, 1.2)
	// writeChunk(c, OP_CONSTANT, 123)
	// writeChunk(c, uint8(constant), 123)
	// writeChunk(c, OP_RETURN, 123)

	disassembleChunk(c, "genesis")
	//freeChunk(c)

}
