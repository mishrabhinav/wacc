package main

// WACC Group 34
//
// wacc.go: Main program
//
// Handles the input file(s)
// Handles flag(s) parsing
// Handles exit codes in case Syntax/Semantic errors are encountered

import (
	"bufio"
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

	//Start Code Generation
	armFile := bufio.NewWriter(os.Stdout)

	if !flags.printAssembly {
		armFileHandle, err := os.Create(flags.assemblyfile)
		if err != nil {
			log.Fatal(err)
		}
		armFile = bufio.NewWriter(armFileHandle)
	}

	for instr := range ast.CodeGen() {
		fInstr := fmt.Sprintf("%v\n", instr)
		fmt.Fprint(armFile, fInstr)
	}

	armFile.Flush()

	flags.Finish()
}
