package main

import (
	"fmt"
)

var scanner *Scanner = new(Scanner)

func initScanner(source *string) {
	fmt.Println(*source)
	scanner.Source = *source
	scanner.Start = 0
	scanner.Current = 0
	scanner.Line = 1
}

func scanToken() Token {
	scanner.Start = scanner.Current // points to current character since we scan one token at a time

	if isAtEnd() {
		return makeToken(TOKEN_EOF)
	}
	return errorToken("Unexpected character")
}

func isAtEnd() bool {
	return string(scanner.Source[scanner.Current]) == "\x01"
}

func makeToken(tokenType TokenType) Token {
	var token *Token = new(Token)
	token.Type = tokenType
	token.Start = scanner.Start
	token.Length = scanner.Current - scanner.Start
	token.Line = scanner.Line
	return *token
}
func errorToken(message string) Token {
	var token *Token = new(Token)
	token.Type = TOKEN_ERROR

	token.Start = scanner.Start
	token.Length = scanner.Current - scanner.Start
	token.Line = scanner.Line
	return *token
}
