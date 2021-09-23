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
	LEFT_P = 1
	RIGHT_P = 2
	LAMBDA = 3
	MULT_OP = 4
  VARIABLE = 5
)

//*************************************************************************

// Character Classes
const (
  LETTER = 6
  DIGIT = 7
  UNKNOWN = 99
  EOF = -1
)


//*************************************************************************

// Function to read a given file byte for byte.
// In the scenario that that the EOF is reached
// return an error.
func getChar() (error) {
  var err error
	nextChar, err = fstream.ReadByte()
	if err != io.EOF {
		if unicode.IsLetter(rune(nextChar)) || unicode.IsDigit(rune(nextChar)) {
      charClass = ALPHA_NUM
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
  if (lexLen < 99) {
    lexeme[lexLen++] = nextChar
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

func getNonBlank () {
  for unicode.IsSpace(rune(nextChar)){
    getChar()
  }
}

//*************************************************************************

func lex() {
  lexLen = 0
  getNonBlank()

  switch charClass {
    case VAR:
      addChar()
      getChar()
      for charClass == VAR {
        addChar()
        getChar()
      }
      nextToken = 
      break
    }

    case UNKNOWN {
      lookup(nextChar)
      getChar()
      break
    }
      

  }
}

//*************************************************************************

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "No arguments given. \n")
		os.Exit(1)
	}	// Given a program that has exactly 1 argument
	// then the length of the "argument" should be 2.
	// As the first argument is the command program itself, and the
	// second the actual argument.

	f, err := os.Open(os.Args[1])	// Opens file
	checkError(err)

	fmt.Fprint(os.Stdout, "test 1 \n")

	fstream = bufio.NewReader(f) // Buffer for the reader

	err = getChar()
	checkError(err) 

	for err == nil {
		fmt.Fprintf(os.Stdout, "Character: %s\n", string(nextChar))
		err = getChar()
	}

}