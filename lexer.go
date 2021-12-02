// lexer.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type Checking
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 30th November, 2021.
//

//*************************************************************************

package main

import (
	"fmt"
	"os"
	"unicode"
)

//*************************************************************************

// Add char to the lexeme.
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

// Gets the first nonblank character, but can also be a newline char.
func getNonBlank() {
	for unicode.IsSpace(nextChar) && nextChar != '\n' {
		err := getChar()
		checkError(err)
	}
}

//*************************************************************************

// Function that assigns nextToken and lexeme
func lex() {
	clearLexeme()
	lexLen = 0
	getNonBlank()

	switch charClass {
	case LETTER: //if lexeme starts with a letter nextToken is a variable
		addChar()
		getChar()
		for charClass == LETTER || charClass == DIGIT {
			addChar()
			checkError(getChar())
		}
		nextToken = VARIABLE
	case DIGIT: //if the lexeme starts with a digit there's an error
		fmt.Fprintf(os.Stderr, "Variable starts with digit \n")
		os.Exit(1)
	case UNKNOWN: //any other case
		lookup(nextChar)
		getChar()
	}
}

//*************************************************************************

// Assigns nextToken BASED on the character.
func lookup(char rune) {
	switch char {
	case '(':
		nextToken = LEFT_P
	case ')':
		nextToken = RIGHT_P
	case 'λ':
		fallthrough
	case '\\':
		addChar()
		nextToken = LAMBDA
	case ':':
		nextToken = COLON
	case '-':
		checkError(getChar())
		nextToken = ARROW
	case '^':
		nextToken = TYPE_ASS
	case '.':
		nextToken = DOT
	case '\n':
		nextToken = EOL
	default:
		nextToken = EOF
	}
}

//*************************************************************************

// Clears lexeme rune array.
func clearLexeme() {
	for i := range lexeme {
		lexeme[i] = 0
	}
}

//*************************************************************************