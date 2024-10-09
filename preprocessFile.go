package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func preprocessFile(filename string) (string, error) {
	// Open a file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", errors.New("could not open file")
	}
	defer file.Close()

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	// Read file line by line
	var res string = ""
	for scanner.Scan() {
		line := scanner.Text()
		for i := 0; i < len(line); i++ {
			fmt.Println(string(line[i]))
			res += string(line[i])
		}
		fmt.Println("Line:", line)
		res += ";"
		res += "\n"

		// Check if the line ends with a newline (this is automatically handled by Scanner)
		// No explicit check is needed since Scanner reads lines one by one
	}
	res += "\x01"
	fmt.Println(res)
	for i := 0; i < len(res); i++ {
		fmt.Println(string(res[i]) == "\n")
	}
	// Check for errors while reading
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return res, nil
}
