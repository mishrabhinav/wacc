package main

import (
	"io/ioutil"
	"log"
)

func main() {
	buffer, err := ioutil.ReadFile("try.wacc")
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
