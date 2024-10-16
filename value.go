package main

import "fmt"

type ValueType int

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

type as struct {
	boolean bool
	number  float64
	obj     Obj
}
type Value struct {
	valueType ValueType
	as        as
}

func IS_BOOL(val Value) bool {
	return val.valueType == VAL_BOOL
}
func IS_NIL(val Value) bool {
	return val.valueType == VAL_NIL
}
func IS_NUMBER(val Value) bool {
	return val.valueType == VAL_NUMBER
}
func IS_OBJ(val Value) bool {
	return val.valueType == VAL_OBJ
}
func AS_BOOL(val Value) bool {
	return val.as.boolean
}
func AS_NUMBER(val Value) float64 {
	return val.as.number
}
func AS_OBJ(val Value) Obj {
	return val.as.obj
}

// Raw GO Primitive to Ube Value Primitive
func BOOL_VAL(value bool) Value {
	var newVal *Value = new(Value)
	newVal.valueType = VAL_BOOL
	newVal.as.boolean = value
	return *newVal
}
func NIL_VAL(nilval float64) Value {
	var newVal *Value = new(Value)
	newVal.valueType = VAL_NIL
	newVal.as.number = 0
	return *newVal
}
func NUMBER_VAL(number float64) Value {
	var newVal *Value = new(Value)
	newVal.valueType = VAL_NUMBER
	newVal.as.number = number
	return *newVal
}

func OBJ_VAL(object Obj) Value {
	var newVal *Value = new(Value)
	newVal.valueType = VAL_OBJ
	newVal.as.obj = object
	return *newVal
}

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
	switch val.valueType {
	case VAL_BOOL:
		if AS_BOOL(val) {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		break
	case VAL_NIL:
		fmt.Println("nil")
		break
	case VAL_NUMBER:
		fmt.Printf("%g", AS_NUMBER(val))
	case VAL_OBJ:
		printObject(val)
	}

}

func valuesEqual(b Value, a Value) bool {
	if b.valueType != a.valueType {
		return false
	}
	switch a.valueType {
	case VAL_BOOL:
		return AS_BOOL(a) == AS_BOOL(b)
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return AS_NUMBER(a) == AS_NUMBER(b)
	case VAL_OBJ:
		return AS_STRING(a).chars == AS_STRING(b).chars
	default:
		return false
	}
}
