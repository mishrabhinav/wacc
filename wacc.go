package main

import (
	"io/ioutil"
	"log"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v FILE\n", os.Args[0])
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	wacc := &WACC{Buffer: string(buffer)}
	wacc.Init()

	if err := wacc.Parse(); err != nil {
		log.Fatal(err)
	}
	wacc.PrintSyntaxTree()
}
