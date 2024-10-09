package main

import (
	"fmt"
	"os"
)

func compile(source *string) {
	initScanner(source)
	line := -1
	for {
		var token Token = scanToken()

		if token.Line != line {
			fmt.Printf("%4d", token.Line)
			line = token.Line
		} else {
			fmt.Printf("   | ")
		}
		word := scanner.Source[scanner.Current : scanner.Current+1]
		fmt.Printf("%d %d %s\n", token.Type, token.Length, word)
		if token.Type == TOKEN_EOF || token.Type == TOKEN_ERROR {
			os.Exit(64)
		}

	}
}
