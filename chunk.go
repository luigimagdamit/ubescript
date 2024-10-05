package main

import (
	"encoding/binary"
)

// opcode enum

const (
	OP_CONSTANT      = iota
	OP_CONSTANT_LONG = iota
	OP_RETURN        = iota
)

type Chunk struct {
	Count    int
	Capacity int

	Code      []uint8
	Lines     []int
	Constants ValueArray
}

// tbh this is isn't really needed i'll just  do it since we likely will need it later
func initChunk(c *Chunk) {
	c.Capacity = 0
	c.Count = 0

	c.Code = []uint8{} // initialize an empty slice
	c.Lines = []int{}
	initValueArray(&c.Constants) // initialize constant pool
}

func writeChunk(c *Chunk, inst uint8, line int) {
	c.Code = append(c.Code, inst)
	c.Lines = append(c.Lines, line)
	c.Capacity++
	c.Count++
}

// to handle long operands
func writeConstant(c *Chunk, val Value, line int) {
	//b1 := uint8()
	var index uint32 = uint32(addContant(c, val))

	var arr [4]byte
	binary.BigEndian.PutUint32(arr[0:4], uint32(index))

	writeChunk(c, OP_CONSTANT_LONG, line)
	writeChunk(c, arr[0], 123)
	writeChunk(c, arr[1], 123)
	writeChunk(c, arr[2], 123)
	writeChunk(c, arr[3], 123)

	a := (uint32(arr[2]) << 8)
	b := (uint32(arr[3]) << 0)

	if a|b != index {
		panic("conversion failed")
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
