package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

//Global Variables
var fstream *bufio.Reader

//Tokens
// I'm not sure what we need to do with these tokens yet
const (
	LEFT_P = iota	// iota = enumerate, so everything below LEFT_P automatically gets assigned an INT of (LEFT_P + 1)
	RIGHT_P // = 2					// You can also do iota + 5 (any int), for when a range of numbers is already used
	LAMBDA // = 3
	VARIABLE // etc...
)

// Function to read a given file byte for byte.
// In the scenario that that the EOF is reached
// return an error.
func getChar() (byte, error) {
	if char, err := fstream.ReadByte(); err != io.EOF {
		return char, err
	} else {
			return 0, errors.New("EOF Reached")
	}
}


// Checks if the programs gives an error, if so
// quit the program with the return value 1.
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

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

	char, err := getChar()
	checkError(err) 

	for err == nil {
		fmt.Fprintf(os.Stdout, "Character: %s\n", string(char))
		char, err = getChar()	// Extra: when a function has two outputs, you can use _ for when
	}												// dont need a certain variable, so for example char, _ = getChar()

}