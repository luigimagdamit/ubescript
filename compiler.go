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
	OpMode    bool
	LH        string
	RH        string
	ValCount  int
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
func Advance() {
	parser.Previous = parser.Current // Consume the current token

	for {
		var curToken Token = scanToken()

		parser.Current = curToken // generate token from recent lexeme in source
		printToken(parser.Previous, "Advance(): parser.Previous")
		printToken(curToken, "Advance(): parser.Current")
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
func Consume(tokenType TokenType) {

	if parser.Current.Type == tokenType {
		Advance()
		return
	}
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("[Consume()]", "Expected Type: ", tokenType)
		fmt.Println(parser.Current)
		printToken(parser.Current, "Consume")
	}

	errorAtCurrent(parser.Current)
}
func check(tokenType TokenType) bool {
	return parser.Current.Type == tokenType
}
func Match(tokenType TokenType) bool {
	if !check(tokenType) {
		return false
	}
	Advance()
	return true
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
func Binary(canAssign bool) {
	parser.OpMode = true
	// left hand operand consumed
	// we have consumed the operator
	var operatorType TokenType = parser.Previous.Type
	var rule *ParseRule = getRule(operatorType)

	// why pass it in with 1?
	parsePrecedence(rule.Precedence + 1)
	switch operatorType {
	case TOKEN_BANG_EQUAL:
		emitBytes(OP_EQUAL, OP_NOT)
	case TOKEN_EQUAL_EQUAL:
		emitByte(OP_EQUAL)
	case TOKEN_GREATER:
		emitByte(OP_GREATER)
	case TOKEN_GREATER_EQUAL:
		emitBytes(OP_LESS, OP_NOT)
	case TOKEN_LESS:
		emitByte(OP_LESS)
	case TOKEN_LESS_EQUAL:
		emitBytes(OP_GREATER, OP_NOT)
	case TOKEN_PLUS:
		emitByte(OP_ADD)
		if COMPILER_MODE {
			fmt.Print("%")
			fmt.Printf("%d", parser.ValCount)
			fmt.Println(" = add i32 " + parser.LH + ", " + parser.RH)
			parser.LH = fmt.Sprintf("%%%d", parser.ValCount)
			parser.ValCount++
		}

	case TOKEN_MINUS:
		emitByte(OP_SUBTRACT)

		if COMPILER_MODE {
			fmt.Print("%")
			fmt.Printf("%d", parser.ValCount)
			fmt.Println(" = sub i32 " + parser.LH + ", " + parser.RH)
			parser.LH = fmt.Sprintf("%%%d", parser.ValCount)
			parser.ValCount++
		}
	case TOKEN_STAR:
		emitByte(OP_MULTIPLY)
	case TOKEN_SLASH:
		emitByte(OP_DIVIDE)
	case TOKEN_DOTDOT:
		emitByte(OP_DOTDOT)
	case TOKEN_PERCENT:
		emitByte(OP_MOD)
	default:
		return
	}

}
func Literal(canAssign bool) {
	switch parser.Previous.Type {
	case TOKEN_FALSE:
		emitByte(OP_FALSE)
	case TOKEN_NIL:
		emitByte(OP_NIL)
	case TOKEN_TRUE:
		emitByte(OP_TRUE)
	default:
		return
	}

}

// Parsing functions for () Grouping. We assume ( as already been consumed
// parser.Previous == "(" Token
func Grouping(canAssign bool) {
	Expr() // Parse Expr, parser.Current should land at )
	msg := "Expect ')' after Expr but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	Consume(TOKEN_RIGHT_PAREN)
}

// Number() will get the lexeme pointed at by parser.Previous
// meaning that the desired parssed number will need to have been consumed / advanced()
// wrapper for writeConstant() and turns lexeme into appropriate bytecode
func Number(canAssign bool) {
	var token Token = parser.Previous
	var lexeme string = getLexeme(token)
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("lex: ", lexeme, parser.Previous)
	}
	//fmt.Print(" i32\n")

	//fmt.Printf("i32 " + getLexeme(token) + "\n")

	value, err := strconv.Atoi(lexeme)
	if err != nil {
		fmt.Println("Error parsing number")
		return
	}
	//fmt.Print(getLexeme(token))
	emitConstant(NUMBER_VAL(float64(value)))

	//
	if !parser.OpMode {
		parser.LH = getLexeme(token)
	} else {
		parser.RH = getLexeme(token)
		parser.OpMode = false
	}

}
func Or(canAssign bool) {
	elseJump := emitJump(OP_JUMP_IF_FALSE)
	endJump := emitJump(OP_JUMP)

	patchJump(elseJump)
	emitByte(OP_POP)

	parsePrecedence(PREC_OR)
	patchJump(endJump)
}
func createString(c string) Obj {
	return *copyString(c, len(c))
}
func trimQuotes(s string) string {
	return s[1 : len(s)-1]
}
func String(canAssign bool) {
	str := (getLexeme(parser.Previous))
	str = trimQuotes(str)
	newString := createString(str)
	emitConstant(OBJ_VAL(newString))
}
func NameVar(name Token, canAssign bool) {
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

	if canAssign && Match(TOKEN_EQUAL) {

		Expr()
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
func Variable(canAssign bool) {
	NameVar(parser.Previous, canAssign)
}

func Unary(canAssign bool) {
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
	case TOKEN_PLUSPLUS:
		emitByte(OP_PLUS_PLUS)
	default:
		return
	}
}

var rules = make([]ParseRule, TOKEN_EOF+1)

func init() {
	rules[TOKEN_LEFT_PAREN] = ParseRule{Grouping, nil, PREC_NONE}
	rules[TOKEN_RIGHT_PAREN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_LEFT_BRACKET] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_RIGHT_BRACKET] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_COMMA] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOT] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_MINUS] = ParseRule{Unary, Binary, PREC_TERM}
	rules[TOKEN_PLUS] = ParseRule{nil, Binary, PREC_TERM}
	rules[TOKEN_SEMICOLON] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SLASH] = ParseRule{nil, Binary, PREC_FACTOR}
	rules[TOKEN_STAR] = ParseRule{nil, Binary, PREC_FACTOR}
	rules[TOKEN_BANG] = ParseRule{Unary, nil, PREC_NONE}
	rules[TOKEN_BANG_EQUAL] = ParseRule{nil, Binary, PREC_EQUALITY}
	rules[TOKEN_EQUAL] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_EQUAL_EQUAL] = ParseRule{nil, Binary, PREC_EQUALITY}
	rules[TOKEN_GREATER] = ParseRule{nil, Binary, PREC_COMPARISON}
	rules[TOKEN_GREATER_EQUAL] = ParseRule{nil, Binary, PREC_COMPARISON}
	rules[TOKEN_LESS] = ParseRule{nil, Binary, PREC_COMPARISON}
	rules[TOKEN_LESS_EQUAL] = ParseRule{nil, Binary, PREC_COMPARISON}
	rules[TOKEN_IDENTIFIER] = ParseRule{Variable, nil, PREC_NONE}
	rules[TOKEN_STRING] = ParseRule{String, nil, PREC_NONE}
	rules[TOKEN_NUMBER] = ParseRule{Number, nil, PREC_NONE}
	rules[TOKEN_AND] = ParseRule{nil, And, PREC_AND}
	rules[TOKEN_CLASS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ELSE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FALSE] = ParseRule{Literal, nil, PREC_NONE}
	rules[TOKEN_FOR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_FUN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_IF] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_NIL] = ParseRule{Literal, nil, PREC_NONE}
	rules[TOKEN_OR] = ParseRule{nil, Or, PREC_OR}
	rules[TOKEN_PRINT] = ParseRule{Unary, nil, PREC_NONE}
	rules[TOKEN_RETURN] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_SUPER] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_THIS] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TRUE] = ParseRule{Literal, nil, PREC_NONE}
	rules[TOKEN_VAR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_WHILE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_ERROR] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_TYPE] = ParseRule{nil, nil, PREC_NONE}
	rules[TOKEN_DOTDOT] = ParseRule{nil, Binary, PREC_TERM}
	rules[TOKEN_LEN] = ParseRule{Unary, nil, PREC_NONE}
	rules[TOKEN_PLUSPLUS] = ParseRule{Unary, nil, PREC_NONE}
	rules[TOKEN_PERCENT] = ParseRule{nil, Binary, PREC_FACTOR}
	rules[TOKEN_EOF] = ParseRule{nil, nil, PREC_NONE}

}

func parsePrecedence(precedence Precedence) {
	if DEBUG_COMPILER_OUTPUT {
		fmt.Println("[parsePrecedence()]======")
	}

	Advance()
	var prevTok Token = parser.Previous
	printToken(prevTok, "[parsePrecedence()] Consumed this token.")
	// access prefix func from rule
	var prefixRule ParseFn = getRule(prevTok.Type).Prefix
	if prefixRule == nil {
		//fmt.Println("eee")
		msg := "Expect Expr, but instead found " + getLexeme(parser.Current)
		prevTok.Message = &msg
		errorAtPrevious(prevTok)
		return
	}
	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	// get precendence rule for current token
	for precedence <= getRule(parser.Current.Type).Precedence {
		Advance() // Consume it if the prec is higher / continue parsing the whole expr
		var infixRule ParseFn = getRule(parser.Previous.Type).Infix
		infixRule(canAssign)
	}
	if canAssign && Match(TOKEN_EQUAL) {
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
				msg := "Cant read local Variable in its own initializer"
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
	if COMPILER_MODE {
		fmt.Printf(getLexeme(parser.Current))
	}

	Consume(TOKEN_IDENTIFIER)
	declareVariable()
	if current.ScopeDepth > 0 {
		return [4]uint8{0, 0, 0, 0}
	}
	return identifierConstant(&parser.Previous)
}
func parseVariableString(errorMessage string) string {
	fmt.Printf(getLexeme(parser.Current))
	Consume(TOKEN_IDENTIFIER)
	declareVariable()

	return getLexeme(parser.Current)

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
func And(canAssign bool) {
	endJump := emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	parsePrecedence(PREC_AND)
	patchJump(endJump)
}

// returns rule at TokenType index
// called by Binary() to look up precedence of operator
func getRule(tokenType TokenType) *ParseRule {
	return &rules[tokenType]
}
func Expr() {
	parsePrecedence(PREC_ASSIGNMENT) // the whole Expr since it is the lowest precedence level
}
func Block() {
	for !check(TOKEN_RIGHT_BRACE) && !check(TOKEN_EOF) {
		Declaration()
	}
	msg := "Expect } after Block but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	Consume(TOKEN_RIGHT_BRACE)
}
func SomeDeclaration() [][4]uint8 {
	var someGlobals ([][4]uint8)
	for !check(TOKEN_EQUAL) && !check(TOKEN_SEMICOLON) && !check(TOKEN_LEFT_PAREN) {

		varName := parseVariable("expect Variable name")
		someGlobals = append(someGlobals, varName)
		Match(TOKEN_COMMA)
	}
	return someGlobals
}
func varDeclaration() {
	if COMPILER_MODE {
		fmt.Print("%")
	}

	global := parseVariable("Expect Variable name")

	msg := "Expected proper type annotation or =, but instead found '" + getLexeme(parser.Current) + "'"
	parser.Current.Message = &msg
	if Match(TOKEN_COMMA) {

		var globals ([][4]uint8)
		globals = append(globals, global)
		globals = append(globals, SomeDeclaration()...)
		if Match(TOKEN_LEFT_PAREN) {
			Advance()
			Consume(TOKEN_RIGHT_PAREN)
		}
		if Match(TOKEN_EQUAL) {
			for i := 0; i < len(globals); i++ {
				Expr() // expr
				// ,
				Match(TOKEN_COMMA)
				defineVariable(globals[i])
			}
		} else {
			for i := 0; i < len(globals); i++ {
				emitByte(OP_NIL) // expr
				// ,
				defineVariable(globals[i])
			}
		}

		Consume(TOKEN_SEMICOLON)
		return
	}
	if Match(TOKEN_TYPE) {
		// idk do something i guess

	} else {

	}
	if Match(TOKEN_EQUAL) {
		if COMPILER_MODE {
			fmt.Printf(" = alloca")
		}

		Expr()
		index := combineUInt8Array(global)
		if COMPILER_MODE {
			fmt.Print(" %", AS_STRING(currentChunk().Constants.Values[index]).chars)
		}

	} else {
		emitByte(OP_NIL)
	}
	msg = "Expected ; but found " + getLexeme(parser.Current)
	parser.Current.Message = &msg
	Consume(TOKEN_SEMICOLON)

	defineVariable(global)
}
func markInitialized() {
	current.Locals[current.LocalCount-1].Depth = current.ScopeDepth
}
func ExprStatement() {
	Expr()

	Consume(TOKEN_SEMICOLON)
	emitByte(OP_POP)
}
func For() {
	if check(TOKEN_NUMBER) {
		Expr()

		loopStart := currentChunk().Count

		exitJump := emitJump(OP_JUMP_IF_FALSE)
		emitByte(OP_EMIT_BREAK)
		emitByte(OP_POP)
		Statement()

		emitLoop(loopStart)

		patchJump(exitJump)
		emitByte(OP_POP)
	}

}
func If() {
	Expr()
	// then statement body
	var thenJump int = emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	Statement()

	// else statement body
	var elseJump = emitJump(OP_JUMP)
	patchJump(thenJump)
	emitByte(OP_POP)
	if Match(TOKEN_ELSE) {
		Statement()
	}
	patchJump(elseJump)
}
func Print() {
	keyword := (getLexeme(parser.Previous))

	Expr()
	if COMPILER_MODE {
		fmt.Print("%")
		fmt.Printf("%d = getelementptr [4 x i8], [4 x i8]* @.str, i32 0, i32 0", parser.ValCount)

		fmt.Printf("\ncall i32 (i8*, ...) @printf(i8* %%%d, i32 %%%d)", parser.ValCount, parser.ValCount-1)
	}
	Consume(TOKEN_SEMICOLON)
	if keyword == "println" {
		if !COMPILER_MODE {
			emitByte(OP_NEWLINE)
		}
	} else {
		emitByte(OP_SHOW)
	}
}
func While() {
	loopStart := currentChunk().Count
	Expr()
	exitJump := emitJump(OP_JUMP_IF_FALSE)
	emitByte(OP_POP)
	Statement()
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
		Advance()
	}
}
func Declaration() {

	if Match(TOKEN_VAR) {

		varDeclaration()
	} else {

		Statement()
	}
	if parser.PanicMode {
		synchronize()
	}
}
func Scope() {
	beginScope()
	Block()
	endScope()
}
func Statement() {

	if Match(TOKEN_PRINT) {

		Print()
	} else if Match(TOKEN_FOR) {
		For()
	} else if Match(TOKEN_IF) {
		If()
	} else if Match(TOKEN_WHILE) {
		While()
	} else if Match(TOKEN_LEFT_BRACE) {
		Scope()
	} else {
		ExprStatement()
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

	Advance()
	// Expr()
	// // NEED TO CHANGE BACK TO EOF PROBABLY
	// Consume(TOKEN_SEMICOLON) // equality check for current
	if COMPILER_MODE {
		fmt.Printf("@.str = private unnamed_addr constant [4 x i8] c\"%%d\\0A\\00\", align 1\ndeclare i32 @printf(i8*, ...)\ndefine i32 @main() {\nentry:\n")
	}

	for !Match(TOKEN_EOF) {
		Declaration()
	}

	if COMPILER_MODE {
		fmt.Print("\nret i32 0\n}\n")
	}

	endCompiler()
	return !parser.HadError
}
