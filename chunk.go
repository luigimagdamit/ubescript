package main

// opcode enum

const (
	OP_CONSTANT      = iota
	OP_CONSTANT_LONG = iota
	OP_RETURN        = iota
)

type Chunk struct {
	Count    int
	Capacity int

	Code         []uint8
	Lines        []int
	LinesEncoded string

	Constants ValueArray
}

// tbh this is isn't really needed i'll just  do it since we likely will need it later
func initChunk(c *Chunk) {
	c.Capacity = 0
	c.Count = 0

	c.Code = []uint8{} // initialize an empty slice
	c.Lines = []int{}
	c.LinesEncoded = ""
	initValueArray(&c.Constants) // initialize constant pool
}

func writeChunk(c *Chunk, inst uint8, line int) {
	c.Code = append(c.Code, inst)
	c.Lines = append(c.Lines, line)

	tmp := decodeToOriginal(c.LinesEncoded)
	tmp = appendSubstring(tmp, line)
	tmp = encodeRunLengthString(tmp)
	c.LinesEncoded = tmp
	// a := encodeRunLengthString(c.LinesEncoded)
	// fmt.Println("cle", c.LinesEncoded)
	// fmt.Println(decodeToOriginal(a))
	c.Capacity++
	c.Count++
}

// to handle long operands
func writeConstant(c *Chunk, val Value, line int) {
	var index uint32 = uint32(addContant(c, val))
	var arr [4]uint8 = splitUInt32(index)

	if combineUInt8Array(arr) != index {
		panic("conversion failed")
	}
	writeChunk(c, OP_CONSTANT_LONG, line)
	// This just the index but split up
	for i := 0; i < 4; i++ {
		writeChunk(c, arr[i], line)
	}

}

func addContant(c *Chunk, val Value) int {
	writeValueArray(&c.Constants, val)
	return c.Constants.Count - 1 // accessing the count field in the ValueArray
}

func freeChunk(c *Chunk) {
	c.Code = []uint8{}
	initChunk(c)
	freeValueArray(&c.Constants)

}
