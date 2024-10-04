package main

import "fmt"

func main() {
	c := new(Chunk)
	initChunk(c)
	writeChunk(c, OP_RETURN)
	writeChunk(c, OP_RETURN)
	writeChunk(c, 0x21)
	disassembleChunk(c, "genesis")

	freeChunk(c)

	fmt.Println(c)
}
