// main.go
// Programming Language: GoLang
//
// Course: Concepts of Programming Language
// Assignment 3: Type checker
// Class 2, Group 11
// Author(s) :	Emanuele Greco (s3375951),
//				Irie Railton (s3292037),
//				Kah ming Wong (s2641976).
//
// Date: 13th November, 2021.
//

//*************************************************************************

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

// Function to read a given file byte for byte and set the appropriate
// charClass. In the scenario that the EOF is reached return an error.
func getChar() error {
	var err error
	nextChar, _, err = fstream.ReadRune()

	if err != io.EOF {
		if unicode.IsLetter(nextChar) && nextChar != 'λ' {
			charClass = LETTER
		} else if unicode.IsDigit(nextChar) {
			charClass = DIGIT
		} else {
			charClass = UNKNOWN
		}
		return err
	} else {
		charClass = UNKNOWN
		return errors.New("EOF Reached")
	}
}

//*************************************************************************

// Add char to the lexeme.
func addChar() {
	if lexLen < 99 {
		lexeme[lexLen] = nextChar
		lexLen++
		lexeme[lexLen] = 0
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Error - lexeme is too long \n")
		os.Exit(1)
	}
}

//*************************************************************************

// Checks if the programs gives an error, if so
// quit the program with the return value 1.
func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
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
	case LETTER: //if the lexeme starts with a letter nextToken is a variable
		addChar()
		_ = getChar()
		for charClass == LETTER || charClass == DIGIT {
			addChar()
			checkError(getChar())
		}
		nextToken = VARIABLE
	case DIGIT: //if the lexeme starts with a digit there's an error
		_, _ = fmt.Fprintf(os.Stderr, "Variable starts with digit \n")
		os.Exit(1)
	case UNKNOWN: //any other case
		lookup(nextChar)
		_ = getChar()
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

// Main driver of the parsing.
func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "No arguments given. \n")
		os.Exit(1)
	} // Check whether an argument (text file) is given

	f, err := os.Open(os.Args[1]) // Opens file
	checkError(err)               // Checks whether file is valid
	fstream = bufio.NewReader(f)  // Buffer for the reader

	err = getChar() // Read first character
	checkError(err)
	for nextToken != EOF {
		context := contextStack{nil} //initialize context
		parse()                      // Parses once
		theJudgment := typeCheck(context, rootExpressionNode, rootTypeNode)
		if theJudgment {
			fmt.Fprintf(os.Stdout, rootExpressionNode.toString()+":"+
				rootTypeNode.toString())
		} else {
			fmt.Fprintf(os.Stderr, rootExpressionNode.toString()+
				" : "+"Cannot type check")
			os.Exit(1)
		}
		println()
	}
	checkError(err)

	os.Exit(0) //exits the program with status 0 when everything is
	//parsed correctly.
}

//*************************************************************************
