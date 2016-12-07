package main

// WACC Group 34
//
// semantic.go: TODO
//
// TODO

import (
	"sync"
)

// mergeErrors merges all the errors in the channel
func mergeErrors(cs []<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done
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
	// done.  This must start after the wg.Add call
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// hasReturn returns true if the statement is/has a return value
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

// checkFunctionReturns checks if the function has a return type. If it is
//  missing it creates an error
func checkFunctionReturns(f *FunctionDef) <-chan error {
	out := make(chan error)

	go func() {
		switch f.returnType.(type) {
		case VoidType:
		default:
			returns := hasReturn(f.body)
			if !returns {
				out <- CreateMissingReturnError(f.token, f.ident)
			}
		}
		close(out)
	}()

	return out
}

// checkJunkStatement checks if there are statments after the function returns
// or exits and produces an error accordingly
func checkJunkStatement(stm Statement) <-chan error {
	out := make(chan error)

	go func() {
		switch t := stm.(type) {
		case *BlockStatement:
			for err := range checkJunkStatement(t.body) {
				out <- err
			}
		case *IfStatement:
			for err := range checkJunkStatement(t.trueStat) {
				out <- err
			}
			for err := range checkJunkStatement(t.falseStat) {
				out <- err
			}
		case *WhileStatement:
			for err := range checkJunkStatement(t.body) {
				out <- err
			}
		case *ForStatement:
			for err := range checkJunkStatement(t.init) {
				out <- err
			}
			for err := range checkJunkStatement(t.after) {
				out <- err
			}
			for err := range checkJunkStatement(t.body) {
				out <- err
			}
		case *ReturnStatement:
			if n := t.next; n != nil {
				out <- CreateUnreachableStatementError(
					stm.Token(),
				)
			}
		}

		if n := stm.GetNext(); n != nil {
			for err := range checkJunkStatement(n) {
				out <- err
			}
		}

		close(out)
	}()

	return out
}

// CheckFunctionCodePaths checks for a return or exit statment in the function
// It also removes the deadcode from the function
func (m *AST) CheckFunctionCodePaths() (errs []error) {
	var errorChannels []<-chan error

	for _, f := range m.functions {
		errorChannels = append(
			errorChannels,
			checkFunctionReturns(f),
		)
		errorChannels = append(
			errorChannels,
			checkJunkStatement(f.body),
		)
	}

	for err := range mergeErrors(errorChannels) {
		errs = append(errs, err)
	}

	return
}
