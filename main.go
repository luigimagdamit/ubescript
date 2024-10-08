package main

func main() {
	initVM()
	c := new(Chunk)
	initChunk(c)

	writeConstant(c, 23, 2)
	writeChunk(c, OP_NEGATE, 2)
	writeChunk(c, OP_RETURN, 2)

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
