package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"
)

//*************************************************************************
/*
TODO: try lex in main
TODO: store string of tokens
TODO: figuring out merge conlifcts
*/
//*************************************************************************

//Global Variables
var fstream *bufio.Reader
var lexeme [100]byte
var nextChar byte
var lexLen int
var token int
var nextToken int
var charClass int

//*************************************************************************

//Tokens
const (
	EOF    = -2
	EOL    = -1
	LEFT_P = iota
	RIGHT_P
	LAMBDA
	VARIABLE
	DOT
)

//*************************************************************************

// Character Classes
const (
	LETTER = iota + 10
	DIGIT
	UNKNOWN = 99
)

//*************************************************************************

// Function to read a given file byte for byte.
// In the scenario that that the EOF is reached
// return an error.
func getChar() error {
	var err error
	nextChar, err = fstream.ReadByte()
	if err != io.EOF {
		if unicode.IsLetter(rune(nextChar)) {
			charClass = LETTER
		} else if unicode.IsDigit(rune(nextChar)) {
			charClass = DIGIT
		} else {
			charClass = UNKNOWN
		}
		return err
	} else {
		return errors.New("EOF Reached")
	}
}

//*************************************************************************

// addChar
func addChar() {
	if lexLen < 99 {
		lexeme[lexLen] = nextChar
		lexLen++
		lexeme[lexLen] = 0
	} else {
		fmt.Fprintf(os.Stderr, "Error - lexeme is too long \n")
		os.Exit(1)
	}
}

//*************************************************************************

// Checks if the programs gives an error, if so
// quit the program with the return value 1.
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//*************************************************************************

func getNonBlank() {
	for unicode.IsSpace(rune(nextChar)) {
		getChar()
	}
}

//*************************************************************************

func lex() int {
	lexLen = 0
	getNonBlank()

	switch charClass {
	case LETTER:
		addChar()
		getChar()
		for charClass == LETTER || charClass == DIGIT {
			addChar()
			getChar()
		}
		nextToken = VARIABLE
		break

	case DIGIT:
		fmt.Fprintf(os.Stderr, "Variable starts with digit \n")
		os.Exit(1)
		break

	case UNKNOWN:
		lookup(nextChar)
		getChar()
		break

	case EOF:
		lexeme[0] = 'E'
		lexeme[1] = 'O'
		lexeme[2] = 'F'
		lexeme[3] = 0
		break
	}

	fmt.Fprintf(os.Stdout, "Next token is: %d, next lexeme is %s \n", nextToken, lexeme)
	return nextToken
}

//*************************************************************************
func lookup(char byte) {
	switch char {
	case '(':
		addChar()
		nextToken = LEFT_P
		break
	case ')':
		addChar()
		nextToken = RIGHT_P
		break
	case '\\':
		addChar()
		nextToken = LAMBDA
		break
	case '\n':
		addChar()
		nextToken = EOL
	case '.':
		addChar()
		nextToken = DOT
	default:
		addChar()
		nextToken = EOF
		break
	}
}

//*************************************************************************

func clearLexeme() {
	for i := range lexeme {
		lexeme[i] = 0
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments given. \n")
		os.Exit(1)
	} // Given a program that has exactly 1 argument
	// then the length of the "argument" should be 2.
	// As the first argument is the command program itself, and the
	// second the actual argument.

	f, err := os.Open(os.Args[1]) // Opens file
	checkError(err)

	fstream = bufio.NewReader(f) // Buffer for the reader

	err = getChar()
	checkError(err)

	for err == nil && nextToken != EOF {
		lex()
		clearLexeme()

		if nextToken == EOL {
			fmt.Fprintf(os.Stdout, "EOL.")
		}
	}
}
