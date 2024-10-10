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
func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}
func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") ||
		(c >= "A" && c <= "Z") ||
		(c == "_")
}
func scanToken() Token {
	skipWhitespace()
	scanner.Start = scanner.Current // points to current character since we scan one token at a time

	if isAtEnd() {
		return makeToken(TOKEN_EOF)
	}

	var c string = advance()
	if isDigit(c) {
		return number()
	}
	switch c {

	case "(":
		return makeToken(TOKEN_LEFT_PAREN)
	case ")":
		return makeToken(TOKEN_RIGHT_PAREN)
	case "{":
		return makeToken(TOKEN_LEFT_BRACE)
	case "}":
		return makeToken(TOKEN_RIGHT_BRACE)
	case ";":
		return makeToken(TOKEN_SEMICOLON)
	case ",":
		return makeToken(TOKEN_COMMA)
	case ".":
		return makeToken(TOKEN_DOT)
	case "-":
		return makeToken(TOKEN_MINUS)
	case "+":
		return makeToken(TOKEN_PLUS)
	case "/":
		return makeToken(TOKEN_SLASH)
	case "*":
		return makeToken(TOKEN_STAR)
	case "!":
		return makeToken(compare("=", TOKEN_BANG_EQUAL, TOKEN_BANG))
	case "<":
		return makeToken(compare("=", TOKEN_LESS_EQUAL, TOKEN_LESS))
	case ">":
		return makeToken(compare("=", TOKEN_GREATER_EQUAL, TOKEN_GREATER))
	case "\"":
		return str()
	}

	return errorToken("Unexpected character")
}

func isAtEnd() bool {
	return string(scanner.Source[scanner.Current]) == "\x01"
}

func advance() string {
	scanner.Current++
	return string(scanner.Source[scanner.Current-1])
}

func peek() string {
	return string(scanner.Source[scanner.Current])
}
func peekNext() string {
	if isAtEnd() {
		return "\x01"
	}
	return string(scanner.Source[scanner.Current+1])
}
func match(expected string) bool {
	if isAtEnd() {
		return false
	}
	if string(scanner.Source[scanner.Current]) != expected {
		return false
	}
	scanner.Current++
	return true
}
func compare(expected string, truthToken TokenType, falseToken TokenType) TokenType {
	if match(expected) {
		return truthToken
	}
	return falseToken
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
	fmt.Printf("===\n" + "===\n")
	token.Line = scanner.Line
	return *token
}
func skipWhitespace() {
	for {
		var c string = peek()
		switch c {
		case " ":
			advance()
			break
		case "\r":
			advance()
			break
		case "\t":
			advance()
			break
		case "\n":
			scanner.Line++
			advance()
			break
		case "/":
			if peekNext() == "/" {
				for peek() != "\n" && !isAtEnd() { // while there;es no newline or end of file byte
					advance()
				}
			} else {
				return
			}

		default:
			return
		}

	}
}

func str() Token {
	for peek() != "\"" && !isAtEnd() { // if the character to be consumed is not a quotation mark
		if peek() == "\n" {
			scanner.Line++
		}
		advance()
	}
	if isAtEnd() {
		return errorToken("unterminated string")
	}
	advance()
	return makeToken(TOKEN_STRING)
}

func number() Token {
	for isDigit(peek()) {
		advance()
	}

	if peek() == "." && isDigit(peekNext()) {
		advance() // consume the decimal .
		for isDigit(peek()) {
			advance()
		}
	}
	return makeToken(TOKEN_NUMBER)
}
