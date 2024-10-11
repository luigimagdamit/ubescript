package main

import (
	"fmt"
	"os"
)

func compile(source *string) {
	initScanner(source)
	for i := 0; i < len(*source); i++ {
		fmt.Println((*source)[i])
	}
	line := -1
	for {
		var token Token = scanToken()

		if token.Line != line {
			line = token.Line
		} else {
			fmt.Printf("   | ")
		}
		word := string(scanner.Source[token.Start : token.Start+token.Length])
		fmt.Println(token.Type == TOKEN_PLUS)
		fmt.Printf("Type: %2d %s| Length: %d | Lexeme: %s\n", token.Type, tokenName(token.Type), token.Length, word)
		if token.Type == TOKEN_EOF || token.Type == TOKEN_ERROR {
			os.Exit(64)
		}

	}
}
