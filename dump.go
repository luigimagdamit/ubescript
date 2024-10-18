package main

import "os"

func write(filename string, contents string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	content := contents
	_, err = file.WriteString(content)

}
