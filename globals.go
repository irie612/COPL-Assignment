// globals.go
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
)

//*************************************************************************

// Global Variables
var fstream *bufio.Reader
var lexeme [100]rune // the lexeme for each token
var nextChar rune    // the current char in the file
var lexLen int       // the current length of the lexeme
var nextToken int    // the current token
var charClass int    // classification of the current char
var rootExpressionNode *node
var rootTypeNode *node

//*************************************************************************

// Tokens
const (
	LEFT_P  = iota // left parenthesis
	RIGHT_P        // right parenthesis
	LAMBDA
	VARIABLE
	DOT
	APPLICATION
	COLON
	TYPE_ASS
	ARROW
	EOF = -2 // end of file
	EOL = -1 // end of line
)

//*************************************************************************

// Character Classes.
const (
	LETTER = iota + 10
	DIGIT
	UNKNOWN = 99
)

//*************************************************************************
