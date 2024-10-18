package main

type Table struct {
	Count    int
	Capacity int
	Entries  map[string]Value
}

func initTable(table *Table) {
	table.Count = 0
	table.Capacity = 0

	table.Entries = make(map[string]Value)
}

func freeTable(table *Table) {
	table.Entries = nil
	table = nil
}

func tableSet(table *Table, key *ObjString, value Value) bool {
	var keyStr string = key.chars

	if table == nil {

		initTable(table)

	} else {

		table.Entries[keyStr] = value
		table.Count++
	}

	return true
}
func tableGet(table *Table, key *ObjString, value Value) bool {
	var keyStr string = key.chars
	retrieved, exists := table.Entries[keyStr]

	if exists {
		//fmt.Println("retrieved ", retrieved)
		value = retrieved
		return true
	}

	return false
}
func tableDelete(table *Table, key *ObjString, value Value) bool {
	delete(table.Entries, key.chars)
	table.Count--
	return true
}
func tableFindString(table *Table, key string) *Obj {
	valRes, exists := table.Entries[key]

	if exists {
		return &valRes.as.obj
	}
	return nil
}
