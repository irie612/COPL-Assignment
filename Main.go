package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readFile() {
	file, err := os.Open(os.Args[1]) //open file specified in argument for reading
	if err != nil {
		log.Fatal(err) //throw error if file opened incorrectly
	}
	defer file.Close() //make sure to close file when done with it

	scanner := bufio.NewScanner(file) //create a scanner
	for scanner.Scan() {              //read in each line
		lexicallyAnalyze(scanner.Text())
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil { //check for scanner error
		log.Fatal(err)
	}
}

func lexicallyAnalyze(expression string) {

}

func main() {
	if len(os.Args) < 2 { //check if arguments parsed
		log.Fatal("Error: file name not specified") //error if not parsed
	} else {
		readFile()
	}

}
