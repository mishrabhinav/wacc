package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var flags Flags
	flags.Parse()

	file, err := os.Open(flags.filename)

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	flags.Start()

	wacc := &WACC{Buffer: string(buffer), File: flags.filename}
	wacc.Init()

	err = wacc.Parse()
	if err != nil {
		log.Print(err)
		os.Exit(100)
	}

	if flags.printPEGTree {
		fmt.Println("-- Printing PEG Tree")
		wacc.PrintSyntaxTree()
	}

	ast, err := ParseAST(wacc)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(100)
	}

	flags.PrintPrettyAST(ast)

	if retErrs := ast.CheckFunctionCodePaths(); len(retErrs) > 0 {
		for _, err := range retErrs {
			fmt.Println(err.Error())
		}
		os.Exit(100)
	}

	if typeErrs := ast.TypeCheck(); len(typeErrs) > 0 {
		for _, err := range typeErrs {
			fmt.Println(err.Error())
		}
		os.Exit(200)
	}

	flags.Finish()
}
