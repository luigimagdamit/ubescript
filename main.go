package main

func main() {
	c := new(Chunk)
	initChunk(c)
	constant := addContant(c, 1.2)
	writeChunk(c, OP_CONSTANT, 123)
	writeChunk(c, uint8(constant), 123) // write the actual offset of the constant in the constant pool into the chunk
	// POSSIBLE ISSUE: we may run into constrains with how many constants we have access to
	writeChunk(c, OP_RETURN, 123)

	for i := 0; i < 10000; i++ {
		//constant := addContant(c, float64(i))

		// basically constant pool overwrites itself after 255, then the constant instruction debug tries to go betond the limit, but valuearr
		// doeesn't go beyond that so
		//fmt.Println(uint8(constant))
		writeConstant(c, float64(i), 123)

	}
	disassembleChunk(c, "genesis")
	//freeChunk(c)

}
