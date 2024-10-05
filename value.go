package main

import "fmt"

type Value = float64

type ValueArray struct {
	Capacity int
	Count    int
	Values   []Value
}

func initValueArray(arr *ValueArray) {
	arr.Values = []Value{}
	arr.Capacity = 0
	arr.Count = 0
}

func writeValueArray(arr *ValueArray, val Value) {
	arr.Values = append(arr.Values, val)
	arr.Capacity++
	arr.Count++
}
func freeValueArray(arr *ValueArray) {
	arr.Values = []Value{}
	initValueArray(arr)
}
func printValue(val Value) {
	fmt.Printf("%g", val)
}
