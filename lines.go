package main

import (
	"strconv"
)

func appendSubstring(s string, i int) string {
	res := s + string(rune(i))

	return res
}

func encodeRunLengthString(s string) string {
	var count int = 1
	var res string = ""
	rs := []rune(s)

	if len(rs) == 1 {
		res += string(rune(count)) + string((rs[0]))
	}
	for i := 1; i < len(rs); i++ {
		r := rs[i]
		num := int(r)

		if num != int(rs[i-1]) {

			//fmt.Println(res)
			res += string(rune(count)) + string((rs[i-1]))
			count = 1
		} else {
			count++
		}
	}
	// if string(rs[len(rs) - 1]) {

	// }
	res += string(rune(count)) + string((rs[len(rs)-1]))
	decodeRunLengthString(res)
	return res

}

func decodeRunLengthString(s string) string {
	rs := []rune(s)
	//fmt.Println(len(rs))
	res := ""
	original := ""
	for i := 0; i < len(rs); i += 2 {
		count := int(rs[i])
		line := int(rs[i+1])
		for j := 0; j < count; j++ {
			original += string(rune(line))
		}
		//fmt.Printf("(%d %d)", count, line)

		//fmt.Println(count, line)
		res += "(" + strconv.Itoa(count) + " " + strconv.Itoa(line) + ")"

	}

	return res
}

func decodeToOriginal(s string) string {
	rs := []rune(s)
	//fmt.Println(len(rs))
	res := ""
	original := ""
	for i := 0; i < len(rs); i += 2 {
		count := int(rs[i])
		line := int(rs[i+1])
		for j := 0; j < count; j++ {
			original += string(rune(line))
		}
		//fmt.Printf("(%d %d)", count, line)

		//fmt.Println(count, line)
		res += "(" + strconv.Itoa(count) + " " + strconv.Itoa(line) + ")"

	}
	//fmt.Println(original)
	return original
}
