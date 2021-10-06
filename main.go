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

//Global Variables
var fstream *bufio.Reader
var lexeme [100]rune
var nextChar rune
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

//Function to read a given file byte for byte and set the appropriate
//charClass. In the scenario that the EOF is reached return an error.
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
		return err //we could probably get rid of it
	} else {
		charClass = UNKNOWN
		return errors.New("EOF Reached")
	}
}

//*************************************************************************

//add char to the lexeme
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

//Checks if the programs gives an error, if so
//quit the program with the return value 1.
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//*************************************************************************

func getNonBlank() {
	for unicode.IsSpace(nextChar) && nextChar != '\n' {
		getChar()
	}
}

//*************************************************************************

//Function that assigns nextToken and lexeme
func lex() {
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
	case UNKNOWN: //any other case
		lookup(nextChar)
		getChar()
		break
	}
}

//*************************************************************************

//Assigns nextToken BASED on the character.
func lookup(char rune) {
	switch char {
	case '(':
		addChar()
		nextToken = LEFT_P
		break
	case ')':
		addChar()
		nextToken = RIGHT_P
		break
	case 'λ':
		fallthrough
	case '\\':
		addChar()
		nextToken = LAMBDA
		break
	case '.':
		addChar()
		nextToken = DOT
		break
	case '\n':
		lexeme[0] = 'E'
		lexeme[1] = 'O'
		lexeme[2] = 'L'
		lexeme[3] = 0
		nextToken = EOL
		break
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

//Clears lexeme rune array.
func clearLexeme() {
	for i := range lexeme {
		lexeme[i] = 0
	}
}

//*************************************************************************

//Adds lexeme to outputString
func addLexeme() {
	outputString += string(lexeme[:lexLen])
}

//*************************************************************************

//Appends string <b> to outputString
func appendToOutputStr(b string) {
	outputString += b
}

//*************************************************************************

//Resolves matching left parentheses in instances where there are more
//right parentheses.
func matchParenthesis(startPos int) {
	var noOpen = strings.Count(outputString[startPos:], "(")
	var noClose = strings.Count(outputString[startPos:], ")")
	for noClose > noOpen {
		outputString = outputString[:startPos] + "(" + outputString[startPos:]
		noOpen++
	}
}

//*************************************************************************

//Finds valid expressions
func expr() {
	var exprStartPos = len(outputString)
	lexpr()
	expr_p()
	matchParenthesis(exprStartPos)
}

//*************************************************************************

//Finds valid expr_p expressions. May be "empty"
func expr_p() {
	if !(nextToken == EOF || nextToken == EOL || nextToken == RIGHT_P) {
		lexpr()
		expr_p()
	}
}

//*************************************************************************

//Finds valid lambda abstractions. If there's no lambda abstractions,
//continue to pexpr.
func lexpr() {
	if nextToken == LAMBDA { //check if we have a lambda abstraction
		addLexeme()
		lex()
		if nextToken == VARIABLE { //check if we have a variable
			addLexeme() //after the lambda
			lex()
			if nextToken == VARIABLE {
				appendToOutputStr(" ")
			}
			if nextToken != EOL && nextToken != EOF {
				lexpr()
			} else {
				fmt.Fprintf(os.Stderr,
					"MISSING EXPRESSION AFTER LAMBDA ABSTRACTION\n")
				os.Exit(1)
			}
		} else { // nextToken != VARIABLE ERROR
			fmt.Fprintf(os.Stderr, "NO VARIABLE AFTER LAMBDA TOKEN\n")
			os.Exit(1)
		}
	} else {
		pexpr()
	}
}

//*************************************************************************

//Looks for a valid pexpr expression.
func pexpr() {
	if nextToken == LEFT_P {
		addLexeme()
		lex()
		if nextToken == RIGHT_P {
			fmt.Fprintf(os.Stderr,
				"MISSING EXPRESSION AFTER OPENING PARENTHESIS\n")
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
		if nextToken == VARIABLE || nextToken == LEFT_P ||
			nextToken == LAMBDA {
			appendToOutputStr(")")
		}
	}
}

//*************************************************************************

//Parses each line in the text file, and outputs the parsed string
//given that no errors has been encountered.
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

//*************************************************************************

// Main driver of the parsing
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments given. \n")
		os.Exit(1)
	} // Check whether an argument (text file) is given

	f, err := os.Open(os.Args[1]) // Opens file
	checkError(err)               // Checks whether file is valid
	fstream = bufio.NewReader(f)  // Buffer for the reader

	err = getChar() //read first character
	checkError(err)

	for err == nil && nextToken != EOF {
		parse()
	}
}
