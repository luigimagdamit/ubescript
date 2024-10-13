package main

type Precedence int

// Define the precedence constants
const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

// Given a token type, let us find
// Function to compile a prefix expression starting with this token type
// Function to compile an infix expression
// Precedence of the infix operator
// There is no need to keep track of prefix precedence since all are same precedence

type ParseRule struct {
	Prefix     ParseFn
	Infix      ParseFn
	Precedence Precedence
}

// just a function type
type ParseFn func()
