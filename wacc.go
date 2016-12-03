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

// main is the starting point of the compiler
func main() {

	// Parse all the supplied arguments
	var flags Flags
	flags.Parse()

	// Open the input wacc file and read the code
	file, err := os.Open(flags.filename)
	addFileToInclude(flags.filename)

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Prints compiler stage, if verbose flag is supplied
	flags.Start()

	// Initialise the Lexer and Parser
	wacc := &WACC{Buffer: string(buffer), File: flags.filename}
	wacc.Init()

	// Parse the supplied code and check for errors
	err = wacc.Parse()
	if err != nil {
		log.Print(err)
		os.Exit(100)
	}

	// Prints the PEG Tree structure, if peg flag is supplied
	if flags.printPEGTree {
		fmt.Println("-- Printing PEG Tree")
		wacc.PrintSyntaxTree()
	}

	// Parse the library generated tree and return the sanitized AST
	ast, err := ParseAST(wacc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(100)
	}

	// Prints the AST in pretty format, if appropriate flag supplied
	flags.PrintPrettyAST(ast)

	// Look for further syntax errors, which were missed by the grammar
	if retErrs := ast.CheckFunctionCodePaths(); len(retErrs) > 0 {
		for _, err := range retErrs {
			fmt.Println(err.Error())
		}
		os.Exit(100)
	}

	// Check the semantics of the syntactically correct program
	if typeErrs := ast.TypeCheck(); len(typeErrs) > 0 {
		for _, err := range typeErrs {
			fmt.Println(err.Error())
		}
		os.Exit(200)
	}

	// Initialise Code Generation
	armFile := bufio.NewWriter(os.Stdout)

	// Put the assembly code in a file, if assembly flag missing
	if !flags.printAssembly {
		armFileHandle, err := os.Create(flags.assemblyfile)
		if err != nil {
			log.Fatal(err)
		}

		defer armFileHandle.Close()

		armFile = bufio.NewWriter(armFileHandle)
	}

	// Take all the instructions in the channel and push them to the defined
	// IO Writer
	for instr := range ast.CodeGen() {
		fInstr := fmt.Sprintf("%v\n", instr)
		fmt.Fprint(armFile, fInstr)
	}

	// Flush the buffered writer
	armFile.Flush()

	// Prints compiler stage, if verbose flag is supplied
	flags.Finish()
}
