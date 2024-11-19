package main

func emitByte(b uint8) {
	writeChunk(currentChunk(), b, parser.Previous.Line)
}
func emitBytes(b1 uint8, b2 uint8) {
	emitByte(b1)
	emitByte(b2)
}
func emitLoop(loopStart int) {
	emitByte(OP_LOOP)

	offset := currentChunk().Count - loopStart + 2

	emitByte((uint8(offset >> 8)) & 0xff)
	emitByte(uint8(offset) & 0xff)
}
func emitJump(instruction uint8) int {
	emitByte(instruction)
	emitByte(0xff)
	emitByte(0xff)
	return currentChunk().Count - 2
}
func emitReturn() {
	emitByte(OP_RETURN)
}
func emitConstant(val Value) int {
	writeChunk(currentChunk(), OP_CONSTANT_LONG, parser.Current.Line)
	res := makeConstant(val)
	emitBytes(res[0], res[1])
	emitBytes(res[2], res[3])
	return int(combineUInt8Array(res))
}
