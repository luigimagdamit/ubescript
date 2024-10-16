package main

import "fmt"

func OBJ_TYPE(value Value) ObjType {
	return AS_OBJ(value).ObjType
}
func IS_STRING(value Value) bool {
	return isObjType(value, OBJ_STRING)
}

// should return objstring
// return the embedded OBJ String
func AS_STRING(value Value) *ObjString {
	return AS_OBJ(value).Value.(*ObjString)
}
func AS_CSTRING(value Value) string {
	return AS_OBJ(value).Value.(*ObjString).chars
}
func ALLOCATE_OBJ(objType ObjType) *Obj {
	var newString *ObjString = new(ObjString)
	return NEW_OBJ(OBJ_STRING, newString)
}
func copyString(chars string, length int) *Obj {
	return allocateString(chars, length)
}
func printObject(value Value) {
	switch OBJ_TYPE(value) {
	case OBJ_STRING:
		fmt.Printf("%s", AS_CSTRING(value))

	}
}
func allocateString(chars string, length int) *Obj {
	var str *Obj = ALLOCATE_OBJ(OBJ_STRING)
	objString := str.Value.(*ObjString)
	objString.chars = chars
	objString.length = length
	str.Value = objString
	return str
}
func takeString(chars string, length int) *Obj {
	return allocateString(chars, length)
}

//	func (o Obj) AS_STRING() ObjString {
//		return o.value
//	}
func NEW_OBJ(typ ObjType, value interface{}) *Obj {
	return &Obj{
		ObjType: typ,
		Value:   value,
	}
}

type ObjType int

const (
	OBJ_STRING ObjType = iota
)

type Obj struct {
	ObjType ObjType
	Value   interface{}
}

type ObjString struct {
	obj    Obj
	length int
	chars  string
}

func isObjType(value Value, objType ObjType) bool {
	return IS_OBJ(value) && AS_OBJ(value).ObjType == objType
}
