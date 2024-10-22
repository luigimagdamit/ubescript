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
var current *Compiler
var compilingChunk *Chunk

func currentChunk() *Chunk {
	return compilingChunk
}
func errorAt(token Token) {
	if parser.PanicMode {
		return
	}
	parser.PanicMode = true
	fmt.Printf("[CompileError] Error on line %d", token.Line)

	if token.Type == TOKEN_EOF {
		fmt.Printf(" at end")
	} else if token.Type == TOKEN_ERROR {

	} else {
		fmt.Printf(" at %d:%d", token.Line, token.Start)
	}
	if token.Message != nil {
		fmt.Printf(": %s\n", *token.Message)
	}

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
		msg := "TOKEN_ERROR Found"
		curToken.Message = &msg
		errorAtCurrent(curToken)
	}
}

// Conditional advance that validates the the current token type
func consume(tokenType TokenType) {

	if parser.Current.Type == tokenType {
		parser_advance()
		return
	}
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("[consume()]", "Expected Type: ", tokenType)
		fmt.Println(parser.Current)
		printToken(parser.Current, "consume")
	}

	errorAtCurrent(parser.Current)
}
func check(tokenType TokenType) bool {
	return parser.Current.Type == tokenType
}
func parseMatch(tokenType TokenType) bool {
	if !check(tokenType) {
		return false
	}
	parser_advance()
	return true
}
func emitByte(b uint8) {
	writeChunk(currentChunk(), b, parser.Previous.Line)
}
func emitBytes(b1 uint8, b2 uint8) {
	emitByte(b1)
	emitByte(b2)
}
func emitLoop(loopStart int) {
	emitByte(OP_LOOP)

	offset := currentChunk().Count - loopStart + 2

	emitByte((uint8(offset >> 8)) & 0xff)
	emitByte(uint8(offset) & 0xff)
}
func emitJump(instruction uint8) int {
	emitByte(instruction)
	emitByte(0xff)
	emitByte(0xff)
	return currentChunk().Count - 2
}
func emitReturn() {
	emitByte(OP_RETURN)
}
func emitConstant(val Value) {
	writeChunk(currentChunk(), OP_CONSTANT_LONG, parser.Current.Line)
	res := makeConstant(val)
	emitBytes(res[0], res[1])
	emitBytes(res[2], res[3])

}
func patchJump(offset int) {
	jump := currentChunk().Count - offset - 1
	if jump > 65535 {
		panic("too much code to jump over")
	}

	currentChunk().Code[offset] = uint8((jump >> 8)) & 0xff
	currentChunk().Code[offset+1] = uint8((jump & 0xff))

}

func makeConstant(val Value) [4]uint8 {
	//fmt.Println("yij")
	return writeConstant(compilingChunk, val, parser.Previous.Line)
}
func initCompiler(compiler *Compiler) {
	compiler.Locals = [4096]Local{}
	compiler.LocalCount = 0
	compiler.ScopeDepth = 0
	current = compiler
}
func endCompiler() {
	if DEBUG_PRINT_CODE && !parser.HadError {
		disassembleChunk(currentChunk(), "code")
	}
	emitReturn()
}
func beginScope() {
	current.ScopeDepth++
}
func endScope() {
	current.ScopeDepth--
	for current.LocalCount > 0 && current.Locals[current.LocalCount-1].Depth >
		current.ScopeDepth {
		emitByte(OP_POP)
		current.LocalCount--
	}
}
func parseBinary(canAssign bool) {
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
func literal(canAssign bool) {

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
func grouping(canAssign bool) {
	expression() // Parse expression, parser.Current should land at )
	msg := "Expect ')' after expression but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	consume(TOKEN_RIGHT_PAREN)
}

// parseNumber() will get the lexeme pointed at by parser.Previous
// meaning that the desired parssed number will need to have been consumed / advanced()
// wrapper for writeConstant() and turns lexeme into appropriate bytecode
func parseNumber(canAssign bool) {
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
func or_(canAssign bool) {
	elseJump := emitJump(OP_JUMP_IF_FALSE)
	endJump := emitJump(OP_JUMP)

	patchJump(elseJump)
	emitByte(OP_POP)

	parsePrecedence(PREC_OR)
	patchJump(endJump)
}
func parseString(canAssign bool) {
	// emitConstant()
	c := (getLexeme(parser.Previous))
	c = c[1 : len(c)-1]
	var objString Obj = *copyString(c, len(c))
	emitConstant(OBJ_VAL(objString))
}
func namedVariable(name Token, canAssign bool) {
	var getOp uint8
	var setOp uint8
	var arg [4]uint8 = resolveLocal(current, &name)

	if combineUInt8Array(arg) != 4294967295 {

		getOp = OP_GET_LOCAL
		setOp = OP_SET_LOCAL
	} else {
		arg = identifierConstant(&name)
		getOp = OP_GET_GLOBAL
		setOp = OP_SET_GLOBAL
	}

	if canAssign && parseMatch(TOKEN_EQUAL) {

		expression()
		emitByte(setOp)

		for i := 0; i < 4; i++ {
			emitByte(arg[i])
		}
	} else {
		emitByte(getOp)
		for i := 0; i < 4; i++ {
			emitByte(arg[i])
		}
	}

}
func variable(canAssign bool) {
	namedVariable(parser.Previous, canAssign)
}

func unary(canAssign bool) {
	var operatorType TokenType = parser.Previous.Type
	// parse at higher level so it ignores binary operators
	parsePrecedence(PREC_UNARY)
	switch operatorType {
	case TOKEN_MINUS:
		emitByte(OP_NEGATE)
	case TOKEN_BANG:
		emitByte(OP_NOT)
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
	rules[TOKEN_IDENTIFIER] = ParseRule{variable, nil, PREC_NONE}
	rules[TOKEN_STRING] = ParseRule{parseString, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = ParseRule{parseNumber, nil, PREC_NONE}
	rules[TOKEN_AND] = ParseRule{nil, and_, PREC_AND}
	rules[TOKEN_CLASS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_FOR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_OR] = ParseRule{nil, or_, PREC_OR}
	rules[TOKEN_PRINT] = ParseRule{unary, nil, PREC_NONE}
	rules[TOKEN_RETURN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = ParseRule{literal, nil, PREC_NONE}
	rules[TOKEN_VAR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_WHILE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ERROR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TYPE] = ParseRule{nil, nil, PREC_NONE}
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
		msg := "Expect expression, but instead found " + getLexeme(parser.Current)
		prevTok.Message = &msg
		errorAtPrevious(prevTok)
		return
	}
	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	// get precendence rule for current token
	for precedence <= getRule(parser.Current.Type).Precedence {
		parser_advance() // consume it if the prec is higher / continue parsing the whole expr
		var infixRule ParseFn = getRule(parser.Previous.Type).Infix
		infixRule(canAssign)
	}
	if canAssign && parseMatch(TOKEN_EQUAL) {
		msg := "Invalid Assignment target"
		parser.Current.Message = &msg
		errorAtCurrent(parser.Current)
	}
}

func identifierConstant(name *Token) [4]uint8 {
	return makeConstant(OBJ_VAL(*copyString(getLexeme(*name), name.Length)))
}
func identifiersEqual(a *Token, b *Token) bool {
	if a.Length != b.Length {
		return false
	} else {
		return getLexeme(*a) == getLexeme(*b)
	}
}
func resolveLocal(compiler *Compiler, name *Token) [4]uint8 {
	for i := compiler.LocalCount - 1; i >= 0; i-- {

		local := &compiler.Locals[i]

		if identifiersEqual(name, &local.Name) {
			if local.Depth == -1 {
				msg := "Cant read local variable in its own initializer"
				parser.Current.Message = &msg
				errorAtCurrent(parser.Current)
			}
			return splitUInt32(uint32(i))
		}
	}

	return splitUInt32(4294967295)
}
func addLocal(name Token) {

	var newLocal *Local = new(Local)

	//local := &current.Locals[current.LocalCount]

	newLocal.Name = name
	// newLocal.Depth = current.ScopeDepth
	newLocal.Depth = -1

	current.Locals[current.LocalCount] = *newLocal
	current.LocalCount++

}
func declareVariable() {
	if current.ScopeDepth == 0 {
		return
	}

	var name *Token = &parser.Previous
	for i := current.LocalCount - 1; i >= 0; i-- {
		local := &current.Locals[i]
		if local.Depth != -1 && local.Depth < current.ScopeDepth {
			break
		}
		if identifiersEqual(name, &local.Name) {
			msg := "Error declaring locals"
			parser.Current.Message = &msg
			errorAtCurrent(parser.Current)
		}
	}
	addLocal(*name)
}
func parseVariable(errorMessage string) [4]uint8 {
	consume(TOKEN_IDENTIFIER)
	declareVariable()
	if current.ScopeDepth > 0 {
		return [4]uint8{0, 0, 0, 0}
	}
	return identifierConstant(&parser.Previous)
}
func defineVariable(global [4]uint8) {
	if current.ScopeDepth > 0 {
		markInitialized()
		return
	}
	emitByte(OP_DEFINE_GLOBAL)

	for i := 0; i < 4; i++ {
		emitByte(global[i])
	}
}
func and_(canAssign bool) {
	endJump := emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	parsePrecedence(PREC_AND)
	patchJump(endJump)
}

// returns rule at TokenType index
// called by parseBinary() to look up precedence of operator
func getRule(tokenType TokenType) *ParseRule {
	return &rules[tokenType]
}
func expression() {
	parsePrecedence(PREC_ASSIGNMENT) // the whole expression since it is the lowest precedence level
}
func block() {
	for !check(TOKEN_RIGHT_BRACE) && !check(TOKEN_EOF) {
		declaration()
	}
	msg := "Expect } after block but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	consume(TOKEN_RIGHT_BRACE)
}
func varDeclaration() {
	global := parseVariable("Expect variable name")

	msg := "Expected proper type annotation or =, but instead found '" + getLexeme(parser.Current) + "'"
	parser.Current.Message = &msg
	if parseMatch(TOKEN_COMMA) {

		var globals ([][4]uint8)
		globals = append(globals, global)

		for !check(TOKEN_EQUAL) && !check(TOKEN_SEMICOLON) {

			varName := parseVariable("expect variable name")
			globals = append(globals, varName)
			parseMatch(TOKEN_COMMA)
		}
		if parseMatch(TOKEN_EQUAL) {
			for i := 0; i < len(globals); i++ {
				expression() // expr
				// ,
				parseMatch(TOKEN_COMMA)
				defineVariable(globals[i])
			}
		} else {
			for i := 0; i < len(globals); i++ {
				emitByte(OP_NIL) // expr
				// ,
				defineVariable(globals[i])
			}
		}

		consume(TOKEN_SEMICOLON)
		return
	}
	if parseMatch(TOKEN_TYPE) {
		// idk do something i guess
	} else {

	}
	if parseMatch(TOKEN_EQUAL) {
		expression()
	} else {
		emitByte(OP_NIL)
	}
	msg = "Expected ; but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	consume(TOKEN_SEMICOLON)
	defineVariable(global)
}
func markInitialized() {
	current.Locals[current.LocalCount-1].Depth = current.ScopeDepth
}
func expressionStatement() {
	expression()

	consume(TOKEN_SEMICOLON)
	emitByte(OP_POP)
}
func ifStatement() {
	expression()
	var thenJump int = emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	statement()
	var elseJump = emitJump(OP_JUMP)
	patchJump(thenJump)
	emitByte(OP_POP)
	if parseMatch(TOKEN_ELSE) {
		statement()
	}
	patchJump(elseJump)
}
func printStatement() {
	keyword := (getLexeme(parser.Previous))

	expression()
	consume(TOKEN_SEMICOLON)
	if keyword == "println" {
		emitByte(OP_NEWLINE)
	} else {
		emitByte(OP_SHOW)
	}

}
func whileStatement() {
	loopStart := currentChunk().Count
	expression()
	exitJump := emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	statement()
	emitLoop(loopStart)

	patchJump(exitJump)
	emitByte(OP_POP)
}
func synchronize() {
	parser.PanicMode = false
	for parser.Current.Type != TOKEN_EOF {
		if parser.Previous.Type == TOKEN_SEMICOLON {
			return
		}
		switch parser.Current.Type {
		case TOKEN_CLASS:
		case TOKEN_FUN:
		case TOKEN_VAR:
		case TOKEN_IF:
		case TOKEN_WHILE:
		case TOKEN_PRINT:
		case TOKEN_RETURN:
			return
		default:

		}
		parser_advance()
	}
}
func declaration() {

	if parseMatch(TOKEN_VAR) {

		varDeclaration()
	} else {
		statement()
	}
	if parser.PanicMode {
		synchronize()
	}
}
func statement() {
	if parseMatch(TOKEN_PRINT) {

		printStatement()
	} else if parseMatch(TOKEN_IF) {
		ifStatement()
	} else if parseMatch(TOKEN_WHILE) {
		whileStatement()
	} else if parseMatch(TOKEN_LEFT_BRACE) {

		beginScope()

		block()
		endScope()
	} else {
		expressionStatement()
	}
}
func compile(source *string, c *Chunk) bool {
	initScanner(source)
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println(*source)
	}
	var compiler Compiler
	initCompiler(&compiler)

	compilingChunk = c
	parser.HadError = false
	parser.PanicMode = false

	parser_advance()
	// expression()
	// // NEED TO CHANGE BACK TO EOF PROBABLY
	// consume(TOKEN_SEMICOLON) // equality check for current
	for !parseMatch(TOKEN_EOF) {
		declaration()
	}
	endCompiler()
	return !parser.HadError
}
