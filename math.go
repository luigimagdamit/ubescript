package main

func add(a Value, b Value) Value {
	return NUMBER_VAL(AS_NUMBER(a) + AS_NUMBER(b))
}
func sub(a Value, b Value) Value {
	return NUMBER_VAL(AS_NUMBER(a) - AS_NUMBER(b))
}
func div(a Value, b Value) Value {
	return NUMBER_VAL(AS_NUMBER(a) / AS_NUMBER(b))
}
func mul(a Value, b Value) Value {
	return NUMBER_VAL(AS_NUMBER(a) * AS_NUMBER(b))
}
func greater(a Value, b Value) Value {
	return BOOL_VAL(AS_NUMBER(a) > AS_NUMBER(b))
}
func less(a Value, b Value) Value {
	return BOOL_VAL(AS_NUMBER(a) < AS_NUMBER(b))
}
