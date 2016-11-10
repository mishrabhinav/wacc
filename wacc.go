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

	if len(os.Args) == 3 && os.Args[2] == "-a" {
		fmt.Println(ast.ASTString())
	}

	if err != nil {
		fmt.Println(err.Error())
		switch err.(type) {
		case *SyntaxError:
			os.Exit(100)
		default:
			os.Exit(1)
		}
	} else if len(os.Args) == 3 && os.Args[2] == "-s" {
		fmt.Println(ast)
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
}
