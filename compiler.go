package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	Current   Token
	Previous  Token
	HadError  bool
	PanicMode bool
}

var parser Parser
var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}
func errorAt(token Token) {
	if parser.PanicMode {
		return
	}
	parser.PanicMode = true
	fmt.Printf("[Line %d] Error", token.Line)

	if token.Type == TOKEN_EOF {
		fmt.Printf(" at end")
	} else if token.Type == TOKEN_ERROR {

	} else {
		fmt.Printf(" at '%d'", token.Start)
	}
	fmt.Printf(": %s\n", token.Message)
	parser.HadError = true
}
func errorAtPrevious(token Token) {
	errorAt(parser.Previous)
}
func errorAtCurrent(token Token) {
	errorAt(parser.Current)
}

func printToken(token Token, name string) {
	word := string(scanner.Source[token.Start : token.Start+token.Length])
	fmt.Printf("[%40s] Type: %2d %20s| Length: %2d | Lexeme: %15s | Line: %d\n", name, token.Type, tokenName(token.Type), token.Length, word, token.Line)
}
func getLexeme(token Token) string {
	return string(scanner.Source[token.Start : token.Start+token.Length])
}
func parser_advance() {
	parser.Previous = parser.Current // consume the current token

	for {
		var curToken Token = scanToken()

		parser.Current = curToken // generate token from recent lexeme in source
		printToken(parser.Previous, "parser_advance(): parser.Previous")
		printToken(curToken, "parser_advance(): parser.Current")
		if parser.Current.Type != TOKEN_ERROR {
			break // if its valid, exit the loop
		}

		// otherwise lets generate the lexical error here
		errorAtCurrent(curToken)
	}
}

// Conditional advance that validates the the current token type
func consume(tokenType TokenType) {
	fmt.Println("[consume()]")
	if parser.Current.Type == tokenType {
		parser_advance()
		return
	}
	errorAtCurrent(parser.Current)
}

func emitByte(b uint8) {
	writeChunk(currentChunk(), b, parser.Previous.Line)
}
func emiteBytes(b1 uint8, b2 uint8) {
	emitByte(b1)
	emitByte(b2)
}
func emitReturn() {
	emitByte(OP_RETURN)
}
func emitConstant(val Value) {
	makeConstant(val)
}
func makeConstant(val Value) {
	writeConstant(compilingChunk, val, parser.Previous.Line)
}
func endCompiler() {
	emitReturn()
}

// parseNumber() will get the lexeme pointed at by parser.Previous
// meaning that the desired parssed number will need to have been consumed / advanced()
// wrapper for writeConstant() and turns lexeme into appropriate bytecode
func parseNumber() {
	var token Token = parser.Previous
	var lexeme string = getLexeme(token)
	fmt.Println("lex: ", lexeme, parser.Previous)
	value, err := strconv.Atoi(lexeme)
	if err != nil {
		fmt.Println("Error parsing number")
		return
	}
	emitConstant(float64(value))

}

// func compile(source *string, c *Chunk) {
// 	initScanner(source)
// 	line := -1
// 	for {
// 		var token Token = scanToken()

// 		if token.Line != line {
// 			if DEBUG_SCANNER_OUTPUT {
// 				fmt.Printf("\n==Line: %2d==\n", token.Line)
// 				fmt.Printf("   | ")
// 			}

// 			line = token.Line
// 		} else {
// 			if DEBUG_SCANNER_OUTPUT {
// 				fmt.Printf("   | ")
// 			}

// 		}
// 		if DEBUG_SCANNER_OUTPUT {
// 			word := string(scanner.Source[token.Start : token.Start+token.Length])
// 			fmt.Printf("Type: %2d %20s| Length: %2d | Lexeme: %15s | Line: %d\n", token.Type, tokenName(token.Type), token.Length, word, token.Line)
// 		}

// 		if token.Type == TOKEN_EOF || token.Type == TOKEN_ERROR {
// 			os.Exit(64)
// 		}

// 	}
// }

func compile(source *string, c *Chunk) bool {
	initScanner(source)
	compilingChunk = c
	parser.HadError = false
	parser.PanicMode = false

	parser_advance() // sets current = lexeme 1
	parser_advance() // consumes "20"

	// expression()
	//writeConstant(c, 1.0, 123)
	parseNumber()
	parser_advance() // consume ";"
	disassembleChunk(compilingChunk, "test")
	consume(TOKEN_EOF) // equality check for current
	endCompiler()
	return !parser.HadError
}
