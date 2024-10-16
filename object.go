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

	res := NEW_OBJ(OBJ_STRING, newString)
	res.next = vm.objects
	vm.objects = res
	return res
}
func copyString(chars string, length int) *Obj {
	fmt.Println(vm.strings.Entries)
	interned := tableFindString(&vm.strings, chars)
	if interned != nil {
		fmt.Println("COPY FOUND")
		return interned
	}
	fmt.Println("allocating new string")
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
	fmt.Println("table set")
	tableSet(&vm.strings, objString, OBJ_VAL(*str))
	fmt.Println(vm.strings.Entries)

	return str
}
func takeString(chars string, length int) *Obj {
	interned := tableFindString(&vm.strings, chars)
	if interned != nil {
		chars = ""
		return interned
	}
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
	next    *Obj
}

type ObjString struct {
	obj    Obj
	length int
	chars  string
}

func isObjType(value Value, objType ObjType) bool {
	return IS_OBJ(value) && AS_OBJ(value).ObjType == objType
}
