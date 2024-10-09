package main

import "fmt"

func main() {
	initVM()
	c := new(Chunk)
	initChunk(c)

	// for i := 0; i < 1000; i++ {
	// 	writeConstant(c, float64(i), 2)
	// 	// writeChunk(c, OP_NEGATE, 2)
	// 	// writeChunk(c, OP_RETURN, 2)
	// }
	// for i := 0; i < 1000; i++ {
	// 	writeConstant(c, 2*float64(i), 2)
	// 	// writeChunk(c, OP_NEGATE, 2)
	// 	// writeChunk(c, OP_RETURN, 2)
	// }
	writeConstant(c, float64(3), 123)
	writeConstant(c, float64(1), 123)
	// writeConstant(c, float64(900), 123)
	writeChunk(c, OP_SUBTRACT, 123)
	writeConstant(c, 4.0, 123)
	writeChunk(c, OP_DIVIDE, 123)

	writeConstant(c, 5.0, 123)
	writeChunk(c, OP_MULTIPLY, 123)
	// writeConstant(c, float64(600), 123)
	// writeChunk(c, OP_SUBTRACT, 123)
	writeChunk(c, OP_RETURN, 123)

	// s := "}}}{{[."

	// fmt.Println(encodeRunLengthString(s))
	vm.chunk = c
	vm.ip = 0

	run()
	disassembleChunk(c, "genesis")
	c.LinesEncoded = encodeRunLengthString(c.LinesEncoded)
	fmt.Println(decodeRunLengthString(c.LinesEncoded))
	freeVM()
	res, err := preprocessFile("example.txt")
	if err != nil {
		fmt.Println("could not open file")
	}
	fmt.Println(res)
	// p := (encodeRunLengthString(c.LinesEncoded))
	// a := encodeRunLengthString(c.LinesEncoded)
	// decodeRunLengthString(a)
	// decodeRunLengthString(p)
	//freeChunk(c)

}
