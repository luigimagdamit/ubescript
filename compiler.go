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
	if DEBUG_COMPILER_OUTPUT {
		fmt.Printf("[%40s] Type: %2d %20s| Length: %2d | Lexeme: %15s | Line: %d\n", name, token.Type, tokenName(token.Type), token.Length, word, token.Line)
	}

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
		curToken.Message = "TOKEN_ERROR Found"
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
	fmt.Println(parser.Current)
	printToken(parser.Current, "consume")
	parser.Current.Message = "Not correct token type: " + string(parser.Current.Type)
	errorAtCurrent(parser.Current)
}

func emitByte(b uint8) {
	writeChunk(currentChunk(), b, parser.Previous.Line)
}
func emitBytes(b1 uint8, b2 uint8) {
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
	if DEBUG_PRINT_CODE && !parser.HadError {
		disassembleChunk(currentChunk(), "code")
	}
	emitReturn()
}

func parseBinary() {
	// left hand operand consumed
	// we have consumed the operator
	var operatorType TokenType = parser.Previous.Type
	var rule *ParseRule = getRule(operatorType)
	// why pass it in with 1?
	parsePrecedence(rule.Precedence + 1)
	switch operatorType {
	case TOKEN_BANG_EQUAL:
		emitBytes(OP_EQUAL, OP_NOT)
		break
	case TOKEN_EQUAL_EQUAL:
		emitByte(OP_EQUAL)
		break
	case TOKEN_GREATER:
		emitByte(OP_GREATER)
		break
	case TOKEN_GREATER_EQUAL:
		emitBytes(OP_LESS, OP_NOT)
		break
	case TOKEN_LESS:
		emitByte(OP_LESS)
		break
	case TOKEN_LESS_EQUAL:
		emitBytes(OP_GREATER, OP_NOT)
		break
	case TOKEN_PLUS:
		emitByte(OP_ADD)
		break
	case TOKEN_MINUS:
		emitByte(OP_SUBTRACT)
		break
	case TOKEN_STAR:
		emitByte(OP_MULTIPLY)
		break
	case TOKEN_SLASH:
		emitByte(OP_DIVIDE)
		break
	case TOKEN_DOTDOT:
		emitByte(OP_DOTDOT)
	default:
		return
	}

}
func literal() {
	switch parser.Previous.Type {
	case TOKEN_FALSE:
		emitByte(OP_FALSE)
		break
	case TOKEN_NIL:
		emitByte(OP_NIL)
		break
	case TOKEN_TRUE:
		emitByte(OP_TRUE)
		break
	default:
		return
	}

}

// Parsing functions for () grouping. We assume ( as already been consumed
// parser.Previous == "(" Token
func grouping() {
	expression() // Parse expression, parser.Current should land at )
	parser.Current.Message = "Expect ')' after expression."
	consume(TOKEN_RIGHT_PAREN)
}

// parseNumber() will get the lexeme pointed at by parser.Previous
// meaning that the desired parssed number will need to have been consumed / advanced()
// wrapper for writeConstant() and turns lexeme into appropriate bytecode
func parseNumber() {
	var token Token = parser.Previous
	var lexeme string = getLexeme(token)
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("lex: ", lexeme, parser.Previous)
	}

	value, err := strconv.Atoi(lexeme)
	if err != nil {
		fmt.Println("Error parsing number")
		return
	}
	emitConstant(NUMBER_VAL(float64(value)))

}
func parseString() {
	// emitConstant()
	c := (getLexeme(parser.Previous))
	c = c[1 : len(c)-1]
	var objString Obj = *copyString(c, len(c))
	emitConstant(OBJ_VAL(objString))
}
func unary() {
	var operatorType TokenType = parser.Previous.Type
	// parse at higher level so it ignores binary operators
	parsePrecedence(PREC_UNARY)
	switch operatorType {
	case TOKEN_MINUS:
		emitByte(OP_NEGATE)
		break
	case TOKEN_BANG:
		emitByte(OP_NOT)
		break
	case TOKEN_LEN:
		emitByte(OP_LEN)
	case TOKEN_PRINT:

		emitByte(OP_SHOW)
	default:
		return
	}
}

var rules = make([]ParseRule, TOKEN_EOF+1)

func init() {
	rules[TOKEN_LEFT_PAREN] = ParseRule{grouping, nil, PREC_NONE}
	rules[TOKEN_RIGHT_PAREN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACKET] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACKET] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_COMMA] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOT] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_MINUS] = ParseRule{unary, parseBinary, PREC_TERM}
	rules[TOKEN_PLUS] = ParseRule{nil, parseBinary, PREC_TERM}
	rules[TOKEN_SEMICOLON] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SLASH] = ParseRule{nil, parseBinary, PREC_FACTOR}
	rules[TOKEN_STAR] = ParseRule{nil, parseBinary, PREC_FACTOR}
	rules[TOKEN_BANG] = ParseRule{unary, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = ParseRule{nil, parseBinary, PREC_EQUALITY}
	rules[TOKEN_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = ParseRule{nil, parseBinary, PREC_EQUALITY}
	rules[TOKEN_GREATER] = ParseRule{nil, parseBinary, PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = ParseRule{nil, parseBinary, PREC_COMPARISON}
	rules[TOKEN_LESS] = ParseRule{nil, parseBinary, PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = ParseRule{nil, parseBinary, PREC_COMPARISON}
	rules[TOKEN_IDENTIFIER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_STRING] = ParseRule{parseString, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = ParseRule{parseNumber, nil, PREC_NONE}
	rules[TOKEN_AND] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_CLASS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_FOR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_OR] = ParseRule{unary, nil, PREC_NONE}
	rules[TOKEN_PRINT] = ParseRule{unary, nil, PREC_NONE}
	rules[TOKEN_RETURN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_VAR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_WHILE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ERROR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOTDOT] = ParseRule{nil, parseBinary, PREC_TERM}
	rules[TOKEN_LEN] = ParseRule{unary, nil, PREC_NONE}
	rules[TOKEN_EOF] = ParseRule{nil, nil, PREC_NONE}

}

func parsePrecedence(precedence Precedence) {
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("[parsePrecedence()]======")
	}

	parser_advance()
	var prevTok Token = parser.Previous
	printToken(prevTok, "[parsePrecedence()] Consumed this token.")
	// access prefix func from rule
	var prefixRule ParseFn = getRule(prevTok.Type).Prefix
	if prefixRule == nil {
		//fmt.Println("eee")
		prevTok.Message = "Expect expression"
		errorAtPrevious(prevTok)
		return
	}
	prefixRule()

	// get precendence rule for current token
	for precedence <= getRule(parser.Current.Type).Precedence {
		parser_advance() // consume it if the prec is higher / continue parsing the whole expr
		var infixRule ParseFn = getRule(parser.Previous.Type).Infix
		infixRule()
	}
}

// returns rule at TokenType index
// called by parseBinary() to look up precedence of operator
func getRule(tokenType TokenType) *ParseRule {
	return &rules[tokenType]
}
func expression() {
	parsePrecedence(PREC_ASSIGNMENT) // the whole expression since it is the lowest precedence level
}

func compile(source *string, c *Chunk) bool {
	initScanner(source)
	fmt.Println(*source)
	compilingChunk = c
	parser.HadError = false
	parser.PanicMode = false

	parser_advance()
	expression()
	// NEED TO CHANGE BACK TO EOF PROBABLY
	consume(TOKEN_SEMICOLON) // equality check for current
	endCompiler()
	return !parser.HadError
}
