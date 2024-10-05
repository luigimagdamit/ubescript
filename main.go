package main

import "fmt"

func main() {
	splitUInt32(0xdeadbeef)
	c := new(Chunk)
	initChunk(c)

	fmt.Println(c.Constants.Values)
	constant := addContant(c, 1.2)
	writeChunk(c, OP_CONSTANT, 14)
	writeChunk(c, OP_CONSTANT, 14)
	writeChunk(c, uint8(constant), 127)
	writeChunk(c, uint8(constant), 127)
	writeChunk(c, OP_RETURN, 444)
	writeChunk(c, OP_RETURN, 444)
	// s := "}}}{{[."

	// fmt.Println(encodeRunLengthString(s))
	disassembleChunk(c, "genesis")

	fmt.Println(c.LinesEncoded)
	decodeRunLengthString(c.LinesEncoded)
	// p := (encodeRunLengthString(c.LinesEncoded))
	// a := encodeRunLengthString(c.LinesEncoded)
	// decodeRunLengthString(a)
	// decodeRunLengthString(p)
	//freeChunk(c)

}
