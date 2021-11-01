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
var lexeme [100]rune    //the lexeme for each token
var nextChar rune       //the current char in the file
var lexLen int          //the current length of the lexeme
var nextToken int       //the current token
var charClass int       //classification of the current char
var outputString string //the final output for the parsing
var rootNode *node

//*************************************************************************

//Tokens
const (
	LEFT_P  = iota //left parenthesis
	RIGHT_P        //right parenthesis
	LAMBDA
	VARIABLE
	DOT
	APPLICATION
	EOF = -2 //end of file
	EOL = -1 //end of line
)

//*************************************************************************

//Character Classes.
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
		return err
	} else {
		charClass = UNKNOWN
		return errors.New("EOF Reached")
	}
}

//*************************************************************************

//Add char to the lexeme.
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

//Gets the first nonblank character, but can also be a newline char.
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
		outputString = outputString[:startPos] + "(" +
			outputString[startPos:]
		noOpen++
	}
}

//*************************************************************************

func expr() *node {
	lexprNode := lexpr()
	exprPNodes := expr_p()
	if len(exprPNodes) == 0 {
		return lexprNode
	} else {
		return appTreeCreate(append([]*node{lexprNode}, exprPNodes...))
	}
}

//*************************************************************************

//Finds valid expr_p expressions. May be "empty"
func expr_p() []*node {
	if !(nextToken == EOF || nextToken == EOL || nextToken == RIGHT_P) {
		lexprNode := lexpr()
		exprPNodes := expr_p()
		exprPNodes = append([]*node{lexprNode}, exprPNodes...)
		return exprPNodes
	}
	return []*node{}
}

//*************************************************************************

//Finds valid lambda abstractions. If there's no lambda abstractions,
//continue to pexpr.
func lexpr() *node {
	if nextToken == LAMBDA { //check if we have a lambda abstraction
		lex()
		if nextToken == VARIABLE { //check if we have a variable
			lambdaNode := newNode(string(lexeme[:lexLen]), LAMBDA)
			lex()
			if nextToken != EOL && nextToken != EOF {
				lambdaNode.linkNodes(lexpr())
				return lambdaNode
			} else {
				fmt.Fprintf(os.Stderr,
					"MISSING EXPRESSION AFTER LAMBDA ABSTRACTION\n")
				os.Exit(1)
			}
		} else { // nextToken != VARIABLE ERROR
			fmt.Fprintf(os.Stderr, "NO VARIABLE AFTER LAMBDA TOKEN\n")
			os.Exit(1)
		}
		return nil
	} else {
		return pexpr()
	}
}

//*************************************************************************

//Looks for a valid pexpr expression.
func pexpr() *node {
	if nextToken == LEFT_P {
		lex()
		if nextToken == RIGHT_P {
			fmt.Fprintf(os.Stderr,
				"MISSING EXPRESSION AFTER OPENING PARENTHESIS\n")
			os.Exit(1)
		}
		exprNode := expr()
		if nextToken != RIGHT_P {
			fmt.Fprintf(os.Stderr, "MISSING CLOSING PARENTHESIS\n")
			os.Exit(1)
		} else {
			lex()
		}
		return exprNode
	} else { //var case
		varNode := newNode(string(lexeme[:lexLen]), VARIABLE)
		lex()
		/*if nextToken == VARIABLE || nextToken == LEFT_P ||
			nextToken == LAMBDA {
			appendToOutputStr(")")
		}*/
		return varNode
	}
}

//*************************************************************************

//Parses each line in the text file, and outputs the parsed string
//given that no errors has been encountered.
func parse() {
	lex()
	if nextToken == EOF {
		return
	}
	rootNode = expr()
	if nextToken != EOL && nextToken != EOF {
		fmt.Fprintf(os.Stderr, "INPUT STRING NOT FULLY PARSED\n")
		os.Exit(1)
	}
	print("BEFORE REDUCTION:\n")
	printTree(rootNode)
	/*DEBUG
	fmt.Fprintf(os.Stdout, "rootNode.token = %d\n", rootNode.token)
	fmt.Fprintf(os.Stdout, "rootNode.value = %s\n", rootNode.value)
	*/
	betaDriver(rootNode)

	print("AFTER REDUCTION\n")
	printTree(rootNode)
}

//*************************************************************************

//Main driver of the parsing
func main() {
	//fmt.Fprintf(os.Stdout, "LEFT_P = %s\n RIGHT_P = %s\n LAMBDA = %s\n VARIABLE = %s\n DOT = %s\n APPLICATION = %s\n LETTER = %s\n DIGIT = %s\n UNKNOWN = %s\n ", strconv.Itoa(LEFT_P), strconv.Itoa(RIGHT_P), strconv.Itoa(LAMBDA), strconv.Itoa(VARIABLE), strconv.Itoa(DOT), strconv.Itoa(APPLICATION), strconv.Itoa(LETTER), strconv.Itoa(DIGIT), strconv.Itoa(UNKNOWN))
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
	} //parses each line until EOF
	os.Exit(0) //exits the program with status 0 when everything is
} //parsed correctly.

//*************************************************************************
