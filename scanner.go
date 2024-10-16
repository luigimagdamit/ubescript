package main

import (
	"fmt"
)

var scanner *Scanner = new(Scanner)

func initScanner(source *string) {
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
	if isAlpha(c) {
		return identifier()
	}
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
	case "[":
		return makeToken(TOKEN_LEFT_BRACKET)
	case "]":
		return makeToken(TOKEN_RIGHT_BRACKET)
	case ";":
		return makeToken(TOKEN_SEMICOLON)
	case ",":
		return makeToken(TOKEN_COMMA)
	case ".":
		return makeToken(compare(".", TOKEN_DOTDOT, TOKEN_DOT))
	case "-":
		return makeToken(TOKEN_MINUS)
	case "+":
		return makeToken(TOKEN_PLUS)
	case "/":
		return makeToken(TOKEN_SLASH)
	case "*":
		return makeToken(TOKEN_STAR)
	case "=":
		return makeToken(compare("=", TOKEN_EQUAL_EQUAL, TOKEN_EQUAL))
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
	token.Message = message
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

func checkKeyword(start int, length int, rest string, tokenType TokenType) TokenType {
	source := scanner.Source
	var sameLength bool = scanner.Current-scanner.Start == start+length
	var remainder string = string(source[scanner.Start+start : scanner.Start+start+length]) // compare "nd" with the rest of the current lexeme
	//fmt.Println(rest)
	if sameLength && remainder == rest {
		return tokenType
	}
	return TOKEN_IDENTIFIER
}

func identifierType() TokenType {
	switch string(scanner.Source[scanner.Start]) {
	case "a":
		return checkKeyword(1, 2, "nd", TOKEN_AND)
	case "c":
		return checkKeyword(1, 4, "lass", TOKEN_CLASS)
	case "e":
		return checkKeyword(1, 3, "lse", TOKEN_ELSE)
	case "f":
		if scanner.Current-scanner.Start > 1 {
			switch string(scanner.Source[scanner.Start+1]) {
			case "a":
				return checkKeyword(2, 3, "lse", TOKEN_FALSE)
			case "o":
				return checkKeyword(2, 1, "r", TOKEN_FOR)
			case "n":
				return checkKeyword(2, 0, "", TOKEN_FUN)
			}
		}
	case "i":
		if scanner.Current-scanner.Start > 1 {
			switch string(scanner.Source[scanner.Start+1]) {
			case "f":
				return checkKeyword(2, 0, "", TOKEN_IF)
			case "n":
				return checkKeyword(2, 1, "t", TOKEN_INT_TAG)
			}
		}

	case "n":
		return checkKeyword(1, 3, "one", TOKEN_NIL)
	case "p":
		return checkKeyword(1, 4, "rint", TOKEN_PRINT)
	case "r":
		return checkKeyword(1, 5, "eturn", TOKEN_RETURN)
	case "s":
		if scanner.Current-scanner.Start > 1 {
			switch string(scanner.Source[scanner.Start+1]) {
			case "u":
				return checkKeyword(2, 3, "per", TOKEN_SUPER)
			case "h":
				return checkKeyword(2, 2, "ow", TOKEN_PRINT)
			case "t":
				if scanner.Current-scanner.Start > 2 {
					switch string(scanner.Source[scanner.Start+3]) {
					case "u":
						return checkKeyword(4, 2, "ct", TOKEN_STRUCT)
					case "i":
						return checkKeyword(4, 2, "ng", TOKEN_STRING_TAG)
					}
				}

			}
		}
	case "t":
		if scanner.Current-scanner.Start > 1 {
			switch string(scanner.Source[scanner.Start+1]) {
			case "h":
				return checkKeyword(2, 2, "is", TOKEN_THIS)
			case "r":
				return checkKeyword(2, 2, "ue", TOKEN_TRUE)
			}
		}
	case "l":

		if scanner.Current-scanner.Start > 1 {
			switch string(scanner.Source[scanner.Start+1]) {
			case "e":
				if scanner.Current-scanner.Start > 2 {
					switch string(scanner.Source[scanner.Start+2]) {
					case "t":
						return checkKeyword(3, 0, "", TOKEN_VAR)
					case "n":
						return checkKeyword(3, 0, "", TOKEN_LEN)
					}
				}
			}
		}
	case "w":
		return checkKeyword(1, 4, "hile", TOKEN_WHILE)
	}
	return TOKEN_IDENTIFIER
}

func identifier() Token {
	for isAlpha(peek()) || isDigit(peek()) {
		advance()
	}
	return makeToken(identifierType())
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
