package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	inputFile, inputError := os.Open("../README.md")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	buf := make([]byte, 1024)
	for {
		n, err := inputReader.Read(buf)
		fmt.Println(n,err)
		if  n==0 || err == io.EOF {
			fmt.Println("Read finished")
			break
		}
	}
}
