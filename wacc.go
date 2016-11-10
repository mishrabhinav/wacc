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

	if flags.verbose {
		fmt.Println("-- Compiling...")
	}

	wacc := &WACC{Buffer: string(buffer)}
	wacc.Init()

	if err := wacc.Parse(); err != nil {
		log.Print(err)
		os.Exit(100)
		os.Exit(100)
	}

	if flags.printPEGTree {
		fmt.Println("-- Printing PEG Tree")
		wacc.PrintSyntaxTree()
	}

	ast, err := ParseAST(wacc)

	if err != nil {
		fmt.Println(err.Error())
		switch err.(type) {
		case *SyntaxError:
			os.Exit(100)
		default:
			os.Exit(1)
		}
	}

	if flags.printPretty {
		fmt.Println("-- Printing Pretty Code")
		fmt.Println(ast)
	}

	if flags.printAST {
		fmt.Println("-- Printing AST")
		fmt.Println(ast.ASTString())
	}

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

	if flags.verbose {
		fmt.Println("-- Finished")
	}
}
