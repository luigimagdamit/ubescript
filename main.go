package main

func main() {
	initVM()
	c := new(Chunk)
	initChunk(c)

	writeConstant(c, float64(1000), 123)
	constant := addContant(c, 1.2)
	writeChunk(c, OP_CONSTANT, 14)
	writeChunk(c, uint8(constant), 127)

	writeChunk(c, OP_RETURN, 444)
	writeChunk(c, OP_RETURN, 444)
	writeChunk(c, OP_RETURN, 444)
	// s := "}}}{{[."

	// fmt.Println(encodeRunLengthString(s))
	vm.chunk = c
	vm.ip = 0

	run()
	//disassembleChunk(c, "genesis")

	freeVM()
	// p := (encodeRunLengthString(c.LinesEncoded))
	// a := encodeRunLengthString(c.LinesEncoded)
	// decodeRunLengthString(a)
	// decodeRunLengthString(p)
	//freeChunk(c)

}
