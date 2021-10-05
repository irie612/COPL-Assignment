package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

//*************************************************************************
/*
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
var outputString string

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

// Character Classes.
const (
	LETTER = iota + 10
	DIGIT
	UNKNOWN = 99
)

//*************************************************************************

// Function to read a given file byte for byte and set the appropriate charClass
// In the scenario that the EOF is reached
// return an error.
func getChar() error {
	var err error
	nextChar, err = fstream.ReadByte() //may be possible to implement support for UNICODE by using ReadRune
	if err != io.EOF {
		if unicode.IsLetter(rune(nextChar)) {
			charClass = LETTER
		} else if unicode.IsDigit(rune(nextChar)) {
			charClass = DIGIT
		} else {
			charClass = UNKNOWN
		}
		return err //we could probably get rid of it
	} else {
		charClass = UNKNOWN //Kinda ugly. Maybe could implement a new char class called EOF but it would be confusing
		return errors.New("EOF Reached")
	}
}

//*************************************************************************

// add char to the lexeme
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
	for unicode.IsSpace(rune(nextChar)) && nextChar != '\n' {
		getChar()
	}
}

//*************************************************************************

// Function that assigns nextToken and lexeme
func lex() int {
	clearLexeme()
	lexLen = 0
	getNonBlank()

	switch charClass {
	case LETTER: //if the lexeme starts with a letter nextToken is a variable
		addChar()
		getChar()
		for charClass == LETTER || charClass == DIGIT {
			addChar()
			getChar()
		}
		nextToken = VARIABLE
		break

	case DIGIT: //if the lexeme starts with a digit there's an error
		fmt.Fprintf(os.Stderr, "Variable starts with digit \n")
		os.Exit(1)
		break

	case UNKNOWN: // any other case
		lookup(nextChar)
		getChar()
		break

	}

	//fmt.Fprintf(os.Stdout, "Next token is: %d, next lexeme is %s \n", nextToken, lexeme)
	return nextToken
}

func clearLexeme() {
	for i := range lexeme {
		lexeme[i] = 0
	}
}

func addLexeme() {
	outputString += string(lexeme[:lexLen])
}

//*************************************************************************
//assigns nextToken BASED on the character.
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
	case '.':
		addChar()
		nextToken = DOT
	case '\n':
		lexeme[0] = 'E'
		lexeme[1] = 'O'
		lexeme[2] = 'L'
		lexeme[3] = 0
		nextToken = EOL
	default:
		lexeme[0] = 'E'
		lexeme[1] = 'O'
		lexeme[2] = 'F'
		lexeme[3] = 0
		nextToken = EOF
		break
	}
}

//*************************************************************************

func matchParenthesis(startPos int) {
	var noOpen = strings.Count(outputString[startPos:], "(")
	var noClose = strings.Count(outputString[startPos:], ")")
	for noClose > noOpen {
		outputString = outputString[:startPos] + "(" + outputString[startPos:]
		noOpen++
	}
}

func parse() {
	outputString = ""
	lex()
	if nextToken == EOF {
		return
	}
	expr()
	if nextToken != EOL && nextToken != EOF {
		fmt.Fprintf(os.Stderr, "INPUT STRING NOT FULLY PARSED\n")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "OUTPUT STRING IS: %s\n", outputString)
}

func expr() {
	//fmt.Fprintf(os.Stdout, "Enter <expr>\n")
	var exprStartPos = len(outputString)
	lexpr()
	expr_p()
	matchParenthesis(exprStartPos)
	//fmt.Fprintf(os.Stdout, "Exit <expr>\n")
}

func expr_p() {
	//fmt.Fprintf(os.Stdout, "Enter <expr_p>\n")

	if !(nextToken == EOF || nextToken == EOL || nextToken == RIGHT_P) {
		lexpr()
		expr_p()
	}
	//fmt.Fprintf(os.Stdout, "Exit <expr_p>\n")
}

func lexpr() {
	//fmt.Fprintf(os.Stdout, "Enter <lexpr>\n")
	if nextToken == LAMBDA { //check if we have a lambda abstraction
		addLexeme()
		lex()
		if nextToken == VARIABLE { //check if we have a variable after the lambda
			addLexeme()
			lex()
			if nextToken == VARIABLE {
				appendToOutputStr(" ")
			}
			if nextToken != EOL && nextToken != EOF {
				lexpr()
			} else {
				fmt.Fprintf(os.Stderr, "MISSING EXPRESSION AFTER LAMBDA ABSTRACTION\n")
				os.Exit(1)
			}
		} else { // nextToken != VARIABLE ERROR
			fmt.Fprintf(os.Stderr, "NO VARIABLE AFTER LAMBDA TOKEN\n")
			os.Exit(1)
		}
	} else {
		pexpr()
	}
	//fmt.Fprintf(os.Stdout, "Exit <lexpr>\n")
}

func appendToOutputStr(b string) {
	outputString += b
}

func pexpr() {
	//fmt.Fprintf(os.Stdout, "Enter <pexpr>\n")
	if nextToken == LEFT_P {
		addLexeme()
		lex()
		if nextToken == RIGHT_P {
			fmt.Fprintf(os.Stderr, "MISSING EXPRESSION AFTER OPENING PARENTHESIS\n")
			os.Exit(1)
		}
		expr()
		addLexeme()
		if nextToken != RIGHT_P {
			fmt.Fprintf(os.Stderr, "MISSING CLOSING PARENTHESIS\n")
			os.Exit(1)
		} else {
			lex()
		}
	} else { //var case
		addLexeme()
		lex()
		if nextToken == VARIABLE || nextToken == LEFT_P {
			appendToOutputStr(")")
		}
	}
	//fmt.Fprintf(os.Stdout, "Exit <pexpr>\n")
}

//*************************************************************************
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

	err = getChar() //read first character
	checkError(err)

	for err == nil && nextToken != EOF {
		parse()
		if nextToken == EOL {
			//fmt.Fprintf(os.Stdout, "END OF LINE\n")
		}
	}
}
