package main

func freeObjects() {
	obj := vm.objects
	for obj != nil {
		next := obj.next
		freeObject(obj)
		obj = next
	}
}

func freeObject(obj *Obj) {
	switch obj.ObjType {
	case OBJ_STRING:
		obj.Value = nil
		obj = nil
	}
}
