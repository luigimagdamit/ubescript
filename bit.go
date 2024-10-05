package main

import (
	"encoding/binary"
)

func splitUInt32(num uint32) [4]uint8 {
	var arr [4]uint8
	binary.BigEndian.PutUint32(arr[0:4], uint32(num))

	return arr
}

func combineUInt8Array(arr [4]uint8) uint32 {
	var res uint32 = 0

	for i, x := range arr {
		res |= uint32(x) << (24 - (8 * i))
	}
	return res
}
