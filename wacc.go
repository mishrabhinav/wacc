package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		log.Print(err)
		os.Exit(100)
		os.Exit(100)
	}

	if len(os.Args) == 3 && os.Args[2] == "-t" {
		wacc.PrintSyntaxTree()
	}

	ast, err := ParseAST(wacc)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(ast.main)
	}
}
