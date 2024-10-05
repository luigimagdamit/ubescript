package main

import (
	"fmt"
	"testing"
)

func TestRunLengthEncode(t *testing.T) {
	c := new(Chunk)
	initChunk(c)

	for i := 0; i < 300; i++ {
		c.LinesEncoded = appendSubstring(c.LinesEncoded, 42)
	}
	for i := 0; i < 200; i++ {
		c.LinesEncoded = appendSubstring(c.LinesEncoded, 10000)
	}
	for i := 0; i < 500; i++ {
		c.LinesEncoded = appendSubstring(c.LinesEncoded, 923)

	}
	actual := encodeRunLengthString(c.LinesEncoded)
	fmt.Println(actual)
	expected := "Ĭ*È✐ǴΛ"

	if actual != expected {
		t.Errorf("Constant Value Mismatch")
	}

}

func TestRunLengthDecode(t *testing.T) {
	s := ""
	for i := 0; i < 100; i++ {
		s = appendSubstring(s, 42)

	}
	for i := 0; i < 100; i++ {
		//writeChunk(c, OP_CONSTANT_LONG, 42)
		s = appendSubstring(s, 10000)

	}
	for i := 0; i < 101; i++ {
		s = appendSubstring(s, 923)

	}
	for i := 0; i < 102; i++ {
		s = appendSubstring(s, 9223)
	}
	es := encodeRunLengthString(s)
	fmt.Println(es)
	if es != "d*d✐eΛf␇" {
		t.Errorf("Constant Value Mismatch")
	}
	if decodeRunLengthString(es) != "(100 42)(100 10000)(101 923)(102 9223)" {
		t.Errorf("Decode not equal")
	}

}

func TestReEncodeMultipleRounds(t *testing.T) {
	s := "**********✐✐✐✐✐✐✐✐✐✐ΛΛΛΛΛ␇␇␇␇␇"
	s = encodeRunLengthString(s)
	s = decodeToOriginal(s)
	s = encodeRunLengthString(s)
	s = decodeToOriginal(s)
	s = encodeRunLengthString(s)
	s = decodeToOriginal(s)
	if s != "**********✐✐✐✐✐✐✐✐✐✐ΛΛΛΛΛ␇␇␇␇␇" {
		t.Errorf("Decode not equal")
	}
}
