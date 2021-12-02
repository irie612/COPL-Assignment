// main.go
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

// Checks if the programs gives an error, if so
// quit the program with the return value 1.
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

//*************************************************************************

// Main driver of the parsing.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments given. \n")
		os.Exit(1)
	} // Check whether an argument (text file) is given

	f, err := os.Open(os.Args[1]) // Opens file
	checkError(err)               // Checks whether file is valid
	fstream = bufio.NewReader(f)  // Buffer for the reader

	err = getChar() // Read first character
	checkError(err)

	for nextToken != EOF {
		parse()
		err = typeCheck(rootExpressionNode, rootTypeNode)
		if err == nil {
			fmt.Fprintf(os.Stdout, rootExpressionNode.toString()+
				":"+rootTypeNode.toString()+"\n")
		} else {
			fmt.Fprintf(os.Stdout, rootExpressionNode.toString()+
				" : ")
			checkError(err)
		}
	}
	
	os.Exit(0) //exits the program with status 0 when everything is
	//parsed correctly.
}

//*************************************************************************