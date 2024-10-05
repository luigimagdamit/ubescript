package main

import (
	"fmt"
	"testing"
)

// expected result should be (0, 0, [])
func TestInitChunk(t *testing.T) {
	var c *Chunk = new(Chunk)
	initChunk(c)

	if c.Capacity != 0 {
		t.Errorf("Capacity should be 0")
	}
	if c.Count != 0 {
		t.Errorf("Count should be 0")
	}
	if len(c.Code) != 0 {
		t.Errorf("Code [] size should be 0")
	}
}

func TestWriteChunk(t *testing.T) {
	c := new(Chunk)
	initChunk(c)
	writeChunk(c, OP_RETURN, 123)

	if c.Capacity != 1 {
		t.Errorf("Capacity should be 0")
	}
	if c.Count != 1 {
		t.Errorf("Count should be 0")
	}
	if len(c.Code) != 1 {
		t.Errorf("Code [] size should be 0")
	}

}

func TestMultipleWriteChunk(t *testing.T) {
	c := new(Chunk)
	initChunk(c)

	for i := 0; i < 100; i++ {
		writeChunk(c, OP_RETURN, 123)
	}
	if c.Capacity != 100 {
		t.Errorf("Capacity should be 0")
	}
	if c.Count != 100 {
		t.Errorf("Count should be 0")
	}
	if len(c.Code) != 100 {
		t.Errorf("Code [] size should be 0")
	}
	if ((c).Code)[99] != OP_RETURN {
		t.Errorf("OpCode Should be OP_RETURN (0)")
	}

}

func TestWriteConstant(t *testing.T) {
	c := new(Chunk)
	initChunk(c)
	for i := 0; i < 10000; i++ {
		writeConstant(c, float64(i), 123)
	}

	for i := 0; i < 10000; i++ {
		if c.Constants.Values[i] != float64(i) {
			fmt.Println(c, c.Constants.Values[i])
			t.Errorf("Constant Value Mismatch")
		}

	}
}
