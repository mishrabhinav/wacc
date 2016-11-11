package main

// WACC Group 34
//
// errors.go: Handles the different types of errors.
//
// File contains functions that return errors given a *token32

import (
	"fmt"
)

type WACCError struct {
	filename     string
	line, column int
}

func (e *WACCError) Error() string {
	return fmt.Sprintf(
		"%s:%d:%d",
		e.filename,
		e.line,
		e.column,
	)
}

func CreateWaccError(token *token32) WACCError {
	return WACCError{
		filename: token.filename,
		line:     token.line,
		column:   token.column,
	}
}

type SyntaxError struct {
	WACCError
}

func CreateSyntaxError(token *token32) SyntaxError {
	return SyntaxError{
		WACCError: CreateWaccError(token),
	}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf(
		"%s:syntax error",
		e.WACCError.Error(),
	)
}

type BigIntError struct {
	SyntaxError
	number string
}

func CreateBigIntError(token *token32, number string) error {
	return &BigIntError{
		SyntaxError: CreateSyntaxError(token),
		number:      number,
	}
}

func (m *BigIntError) Error() string {
	return fmt.Sprintf(
		"%s: number '%s' does not fit in integer",
		m.SyntaxError.Error(),
		m.number,
	)
}

type MissingReturnError struct {
	SyntaxError
	ident string
}

func CreateMissingReturnError(token *token32, ident string) error {
	return &MissingReturnError{
		SyntaxError: CreateSyntaxError(token),
		ident:       ident,
	}
}

func (e *MissingReturnError) Error() string {
	return fmt.Sprintf(
		"%s: expected function '%s' to return",
		e.SyntaxError.Error(),
		e.ident,
	)
}

type UnreachableStatementError struct {
	SyntaxError
}

func (e *UnreachableStatementError) Error() string {
	return fmt.Sprintf(
		"%s: unreachable statement",
		e.SyntaxError.Error(),
	)
}

func CreateUnreachableStatementError(token *token32) error {
	return &UnreachableStatementError{
		SyntaxError: CreateSyntaxError(token),
	}
}

type SemanticError struct {
	WACCError
}

func CreateSemanticError(token *token32) SemanticError {
	return SemanticError{
		WACCError: CreateWaccError(token),
	}
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf(
		"%s:semantic error",
		e.WACCError.Error(),
	)
}

type VariableRedeclarationError struct {
	SemanticError
	ident string
	prev  Type
	new   Type
}

func (e *VariableRedeclarationError) Error() string {
	return fmt.Sprintf(
		"%s: '%s' redaclared from '%s' to '%s'",
		e.SemanticError.Error(),
		e.ident,
		e.prev.String(),
		e.new.String(),
	)
}

func CreateVariableRedeclarationError(token *token32, ident string, oldt, newt Type) error {
	return &VariableRedeclarationError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
		prev:          oldt,
		new:           newt,
	}
}

type UndeclaredVariableError struct {
	SemanticError
	ident string
}

func (e *UndeclaredVariableError) Error() string {
	return fmt.Sprintf(
		"%s: '%s' is undeclared",
		e.SemanticError.Error(),
		e.ident,
	)
}

func CreateUndelaredVariableError(token *token32, ident string) error {
	return &UndeclaredVariableError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

type TypeMismatchError struct {
	SemanticError
	expected Type
	got      Type
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf(
		"%s: type mismatch expected '%s' got '%s'",
		e.SemanticError.Error(),
		e.expected.String(),
		e.got.String(),
	)
}

func CreateTypeMismatchError(token *token32, expected, got Type) error {
	return &TypeMismatchError{
		SemanticError: CreateSemanticError(token),
		expected:      expected,
		got:           got,
	}
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

type FunctionRedeclarationError struct {
	SemanticError
	ident string
}

func (e *FunctionRedeclarationError) Error() string {
	return fmt.Sprintf(
		"%s: function '%s' already declared",
		e.SemanticError.Error(),
		e.ident,
	)
}

func CreateFunctionRedelarationError(token *token32, ident string) error {
	return &FunctionRedeclarationError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}
