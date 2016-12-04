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
	"path/filepath"
)

// parseInput checks the syntax of the input file and exits if there
// are any errors
func parseInput(filename string) *WACC {
	// Open the input wacc file and read the code
	file, err := os.Open(filename)

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Initialise the Lexer and Parser
	wacc := &WACC{Buffer: string(buffer), File: filename}
	wacc.Init()

	// Parse the supplied code and check for errors
	err = wacc.Parse()
	if err != nil {
		log.Print(err)
		os.Exit(100)
	}

	return wacc
}

// generateASTFromWACC takes the wacc file and the included files and generates
// the AST
func generateASTFromWACC(wacc *WACC, ifm *IncludeFiles) *AST {
	// Parse the library generated tree and return the sanitized AST
	ast, err := ParseAST(wacc, ifm)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(100)
	}

	// Look for further syntax errors, which were missed by the grammar
	if retErrs := ast.CheckFunctionCodePaths(); len(retErrs) > 0 {
		for _, err := range retErrs {
			fmt.Println(err.Error())
		}
		os.Exit(100)
	}

	return ast
}

// semanticAnalysis checks the semantics of the imput file and exits if there
// any errors
func semanticAnalysis(ast *AST) {
	// Check the semantics of the syntactically correct program
	if typeErrs := ast.TypeCheck(); len(typeErrs) > 0 {
		for _, err := range typeErrs {
			fmt.Println(err.Error())
		}
		os.Exit(200)
	}
}

// codeGeneration generates the assembly code for the input file and puts it in
// a `.s` file
func codeGeneration(ast *AST, flags *Flags) {
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
}

// main is the starting point of the compiler
func main() {

	// Parse all the supplied arguments
	var flags = &Flags{}
	flags.Parse()

	// Prints compiler stage, if verbose flag is supplied
	flags.Start()

	// Get the directory of the base file
	dir := filepath.Dir(flags.filename)

	// Create a new instace of the IncludeFiles struct
	ifm := &IncludeFiles{dir: dir}
	ifm.Include(flags.filename)

	// Initial syntax analysis by the lexer/parser library
	wacc := parseInput(flags.filename)

	// Generate AST from the WACC struct produced by the peg library
	ast := generateASTFromWACC(wacc, ifm)

	// Prints the PEG Tree structure, if peg flag is supplied
	if flags.printPEGTree {
		fmt.Println("-- Printing PEG Tree")
		wacc.PrintSyntaxTree()
	}

	// Prints the AST in pretty format, if appropriate flag supplied
	flags.PrintPrettyAST(ast)

	// Perform semantic analysis on the AST
	semanticAnalysis(ast)

	// Generate assembly code for the input wacc file
	codeGeneration(ast, flags)

	// Prints compiler stage, if verbose flag is supplied
	flags.Finish()
}
