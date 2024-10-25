package main

// OpCode Enums will be the bytecode instructions that the virtual machine writes
// OpCodes will be treated as a singular byte (0x00) uint8 instructions
// Operands such as constant pool indexes should be multiples of 2, or combinations of multiple bytes

const (
	OP_CONSTANT = iota
	OP_CONSTANT_LONG
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_POP
	OP_GET_LOCAL
	OP_SET_LOCAL
	OP_GET_GLOBAL
	OP_DEFINE_GLOBAL
	OP_SET_GLOBAL
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_DOTDOT
	OP_MOD
	OP_EMIT_BREAK
	OP_PLUS_PLUS
	OP_LEN
	OP_SHOW
	OP_NEWLINE
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_LOOP
	OP_NOT
	OP_RETURN
)

// Chunk struct will hold all the information for an array of OpCOdes
// It also maintains the count, capacity of how many we can carry. These parameters are not necessary since we allow
// Go to handle the dynamic array
// LinesEncoded will hold an encoded version of line information of every instruction within the chunk
// This is done in the format of Run Length Encoding
// ValueArray holds all of the constant values that are used by the operations
type Chunk struct {
	Count    int
	Capacity int

	Code         []uint8
	Lines        []int
	LinesEncoded string

	Constants ValueArray
}

// initChunk initializes the chunk fields as well as the init function for the ValueArray type
func initChunk(c *Chunk) {
	c.Capacity = 0
	c.Count = 0

	c.Code = []uint8{} // initialize an empty slice
	c.Lines = []int{}
	c.LinesEncoded = ""
	initValueArray(&c.Constants) // initialize constant pool
}

// writeChunk will dynamically append a 1 byte instruction to the `Code` slice
// Decodes the encoded line information, appends the new instructions line info, then re encodes it
func writeChunk(c *Chunk, inst uint8, line int) {
	c.Code = append(c.Code, inst)
	c.Lines = append(c.Lines, line)

	//tmp := decodeToOriginal(c.LinesEncoded)
	tmp := appendSubstring(c.LinesEncoded, line)
	//tmp = encodeRunLengthString(tmp)
	c.LinesEncoded = tmp

	c.Capacity++
	c.Count++
}

// writeConstants is an alternate function that is used for writing constants to the constant pool and the index to the chunk
// The index is written from a uint32 version of an index, relative to the size of the ValueArray(constant pool) - necessary to encode more than 255 constants - Refer to Ch 1 of CI
// The 32-bit int is split into an array of four 8 bit instructions, and fed into the chunk as the OpCodes are
func writeConstant(c *Chunk, val Value, line int) [4]uint8 {
	var index uint32 = uint32(addContant(c, val))
	var arr [4]uint8 = splitUInt32(index)

	if combineUInt8Array(arr) != index {
		panic("conversion failed")
	}
	// may have to remove the writeChunk so that the OP_CONSTANT_LONG emit can be handled separately
	//writeChunk(c, OP_CONSTANT_LONG, line)
	// This just the index but split up
	for i := 0; i < 4; i++ {
		//writeChunk(c, arr[i], line)
	}
	//fmt.Println("OFFSET ", combineUInt8Array(arr), arr)
	return arr

}

// addConstant (old) writes the Value to the chunk's ValueArray
// Returns an integer that is the index of the Value in the array when its appended, since the writeValueArray function appends the count field in the chunk.ValueArray
func addContant(c *Chunk, val Value) int {
	writeValueArray(&c.Constants, val)
	return c.Constants.Count - 1 // accessing the count field in the ValueArray
}

// freeChunk resets all the fields to their default state
func freeChunk(c *Chunk) {
	c.Code = []uint8{}
	initChunk(c)
	freeValueArray(&c.Constants)

}
