package main

import (
	"fmt"
)

type WACCError struct {
	file         string
	line, column int
}

func (e *WACCError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d",
		e.file,
		e.line,
		e.column,
	)
}

type SyntaxError struct {
	file         string
	line, column int
	msg          string
}

func (m *SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d:error: %s", m.file, m.line, m.column, m.msg)
}

type SemanticError struct {
	WACCError
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf(
		"%s:semantic error",
		e.WACCError.Error(),
	)
}

type VariableRedeclaration struct {
	SemanticError
	ident string
	prev  Type
	new   Type
}

func (e *VariableRedeclaration) Error() string {
	return fmt.Sprintf(
		"%s: '%s' redaclared from '%s' to '%s'",
		e.SemanticError.Error(),
		e.ident,
		e.prev.String(),
		e.new.String(),
	)
}

type UndeclaredVariable struct {
	SemanticError
	ident string
}

func (e *UndeclaredVariable) Error() string {
	return fmt.Sprintf(
		"%s: '%s' is undeclared",
		e.SemanticError.Error(),
		e.ident,
	)
}

type TypeMismatch struct {
	SemanticError
	expected Type
	got      Type
}

func (e *TypeMismatch) Error() string {
	return fmt.Sprintf(
		"%s: type mismatch expected '%s' got '%s'",
		e.SemanticError.Error(),
		e.expected.String(),
		e.got.String(),
	)
}

type CallingNonFunction struct {
	SemanticError
	ident string
}

func (e *CallingNonFunction) Error() string {
	return fmt.Sprintf(
		"%s: calling non function '%s'",
		e.SemanticError.Error(),
		e.ident,
	)
}

type FunctionCallWrongArity struct {
	SemanticError
	ident    string
	expected int
	got      int
}

func (e *FunctionCallWrongArity) Error() string {
	return fmt.Sprintf(
		"%s: '%s' called with '%d' arguments expected '%d'",
		e.SemanticError.Error(),
		e.ident,
		e.got,
		e.expected,
	)
}
