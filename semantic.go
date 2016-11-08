package main

import (
	"fmt"
	"sync"
)

func mergeErrors(cs []<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func hasReturn(stm Statement) bool {
	if stm == nil {
		return false
	}

	switch t := stm.(type) {
	case *BlockStatement:
		return hasReturn(t.body) || hasReturn(t.next)
	case *ReturnStatement:
		return true
	case *ExitStatement:
		return true
	case *IfStatement:
		return (hasReturn(t.trueStat) && hasReturn(t.falseStat)) ||
			hasReturn(t.next)
	default:
		return hasReturn(t.GetNext())
	}
}

func checkFunctionReturns(f *FunctionDef) <-chan error {
	out := make(chan error)

	go func() {
		returns := hasReturn(f.body)
		if !returns {
			out <- &SyntaxError{
				file:   "",
				line:   0,
				column: 0,
				msg:    fmt.Sprintf("Expection function '%s' to return", f.ident),
			}
		}
		close(out)
	}()

	return out
}

func (m *AST) CheckFunctionCodePaths() (errs []error) {
	var errorChannels []<-chan error

	for _, f := range m.functions {
		errorChannels = append(errorChannels, checkFunctionReturns(f))
	}

	for err := range mergeErrors(errorChannels) {
		errs = append(errs, err)
	}

	return
}
