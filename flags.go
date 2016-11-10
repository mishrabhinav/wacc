package main

import (
	"flag"
)

type Flags struct {
	filename     string
	printPEGTree bool
	printPretty  bool
	printAST     bool
}

func (f *Flags) Parse() {
	flag.StringVar(&f.filename, "file", "", "Input File")

	flag.BoolVar(&f.printPEGTree, "peg", false, "Print PEG tree for the supplied file")
	flag.BoolVar(&f.printPretty, "pretty", false, "Pretty print the supplied file")
	flag.BoolVar(&f.printAST, "ast", false, "Print AST for the supplied file")

	flag.Parse()
}
