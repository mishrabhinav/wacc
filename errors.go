package main

// WACC Group 34
//
// errors.go: Handles the different types of errors.
//
// File contains functions that return errors given a *token32 and supporting
// information

import (
	"fmt"
)

// WACCError is the base error type with filename, line and column number
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

// CreateWACCError pulls the position information from a *token32
func CreateWACCError(token *token32) WACCError {
	return WACCError{
		filename: token.filename,
		line:     token.line,
		column:   token.column,
	}
}

// SyntaxError is the base type for syntax errors
type SyntaxError struct {
	WACCError
}

// CreateSyntaxError initializes the WACCError position information
func CreateSyntaxError(token *token32) SyntaxError {
	return SyntaxError{
		WACCError: CreateWACCError(token),
	}
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf(
		"%s:syntax error",
		e.WACCError.Error(),
	)
}

// BigIntError is a syntax error when a number cannot fit the integer size
type BigIntError struct {
	SyntaxError
	number string
}

// CreateBigIntError creates an error from the token and number as a string
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

// MissingReturnError is a syntax error when a function does not return
type MissingReturnError struct {
	SyntaxError
	ident string
}

// CreateMissingReturnError creates an error from the token and the function
// identifier
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

// UnreachableStatementError is a syntax error when a statement is present after
// return
type UnreachableStatementError struct {
	SyntaxError
}

func (e *UnreachableStatementError) Error() string {
	return fmt.Sprintf(
		"%s: unreachable statement",
		e.SyntaxError.Error(),
	)
}

// CreateUnreachableStatementError creates an error from the token
func CreateUnreachableStatementError(token *token32) error {
	return &UnreachableStatementError{
		SyntaxError: CreateSyntaxError(token),
	}
}

// SemanticError is the base type for semantic errors
type SemanticError struct {
	WACCError
}

// CreateSemanticError initializes the WACCError position information
func CreateSemanticError(token *token32) SemanticError {
	return SemanticError{
		WACCError: CreateWACCError(token),
	}
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf(
		"%s:semantic error",
		e.WACCError.Error(),
	)
}

// VariableRedeclarationError is a semantic error when a variable is declared
// again within the same scope
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

// CreateVariableRedeclarationError creates an error from the token, variable
// identifier, previous and new type
func CreateVariableRedeclarationError(token *token32, ident string, oldt, newt Type) error {
	return &VariableRedeclarationError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
		prev:          oldt,
		new:           newt,
	}
}

// UndeclaredVariableError is a semantic error when trying to access an
// undeclared variable
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

// ClassRedeclarationError is a semantic error when a variable is declared
// again within the same scope
type ClassRedeclarationError struct {
	SemanticError
	ident string
}

func (e *ClassRedeclarationError) Error() string {
	return fmt.Sprintf(
		"%s: class '%s' already declared",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateClassRedeclarationError creates an error from the token, variable
// identifier, previous and new type
func CreateClassRedeclarationError(token *token32, ident string) error {
	return &ClassRedeclarationError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// CreateUndeclaredClassError creates an error from the token and variable
// identifier
func CreateUndeclaredClassError(token *token32, ident string) error {
	return &UndeclaredClassError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// UndeclaredClassError is a semantic error when trying to access an
// undeclared variable
type UndeclaredClassError struct {
	SemanticError
	ident string
}

func (e *UndeclaredClassError) Error() string {
	return fmt.Sprintf(
		"%s: class '%s' is undeclared",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateUndeclaredVariableError creates an error from the token and variable
// identifier
func CreateUndeclaredVariableError(token *token32, ident string) error {
	return &UndeclaredVariableError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// TypeMismatchError is a semantic error when trying to operate on incompatible
// types
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

// CreateTypeMismatchError creates an error from the token, the expected and the
// received type
func CreateTypeMismatchError(token *token32, expected, got Type) error {
	return &TypeMismatchError{
		SemanticError: CreateSemanticError(token),
		expected:      expected,
		got:           got,
	}
}

// CallingNonFunctionError is a semantic error trying to call an undeclared
// function
type CallingNonFunctionError struct {
	SemanticError
	ident string
}

func (e *CallingNonFunctionError) Error() string {
	return fmt.Sprintf(
		"%s: calling non function '%s'",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateCallingNonFunctionError creates an error from a token and a function
// identifier
func CreateCallingNonFunctionError(token *token32, ident string) error {
	return &CallingNonFunctionError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// FunctionCallWrongArityError is a semantic error trying to call a function
// with the wrong number of arguments
type FunctionCallWrongArityError struct {
	SemanticError
	ident    string
	expected int
	got      int
}

func (e *FunctionCallWrongArityError) Error() string {
	return fmt.Sprintf(
		"%s: '%s' called with '%d' arguments expected '%d'",
		e.SemanticError.Error(),
		e.ident,
		e.got,
		e.expected,
	)
}

// CreateFunctionCallWrongArityError creates and error from a token, function
// identifier, expected and received number of parameters
func CreateFunctionCallWrongArityError(
	token *token32,
	ident string,
	expected, got int) error {
	return &FunctionCallWrongArityError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
		expected:      expected,
		got:           got,
	}
}

// FunctionCallOnNonObjectError is a semantic error trying to call a function
// on a non class instance
type FunctionCallOnNonObjectError struct {
	SemanticError
	ident string
	wtype Type
}

func (e *FunctionCallOnNonObjectError) Error() string {
	return fmt.Sprintf(
		"%s: trying to call '%s' on non object of type '%s'",
		e.SemanticError.Error(),
		e.ident,
		e.wtype,
	)
}

// CreateFunctionCallOnNonObjectError creates and error from a token, function
// identifier, and type
func CreateFunctionCallOnNonObjectError(
	token *token32,
	ident string,
	wtype Type) error {
	return &FunctionCallOnNonObjectError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
		wtype:         wtype,
	}
}

// FunctionRedeclarationError is a semantic error when trying to declare a
// function again after it has been declared
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

// CreateFunctionRedelarationError creates an error from a token and a function
// identifier
func CreateFunctionRedelarationError(token *token32, ident string) error {
	return &FunctionRedeclarationError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// AmbigousFunctionCallError is a semantic error when multiple overloads of a
// function match the provided parameters
type AmbigousFunctionCallError struct {
	SemanticError
	ident string
}

func (e *AmbigousFunctionCallError) Error() string {
	return fmt.Sprintf(
		"%s: calling function '%s' where multiple overloads match",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateAmbigousFunctionCallError creates an error from a token and a function
// identifier
func CreateAmbigousFunctionCallError(token *token32, ident string) error {
	return &AmbigousFunctionCallError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// NoSuchOverloadError is a semantic error trying to call a non-existing variant
// of an overloaded function
type NoSuchOverloadError struct {
	SemanticError
	ident string
}

func (e *NoSuchOverloadError) Error() string {
	return fmt.Sprintf(
		"%s: calling function '%s' where no overload matches",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateNoSuchOverloadError creates an error from a token and a function
// identifier
func CreateNoSuchOverloadError(token *token32, ident string) error {
	return &NoSuchOverloadError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// InvalidVoidTypeError is a semantic error declaring a variable with void type
type InvalidVoidTypeError struct {
	SemanticError
	ident string
}

func (e *InvalidVoidTypeError) Error() string {
	return fmt.Sprintf(
		"%s: void type for variable '%s' is invalid",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateInvalidVoidTypeError creates an error from a token and an identifier
func CreateInvalidVoidTypeError(token *token32, ident string) error {
	return &InvalidVoidTypeError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}

// VoidAssignmentError is a semantic error when trying to assign a void return
// value
type VoidAssignmentError struct {
	SemanticError
	ident string
}

func (e *VoidAssignmentError) Error() string {
	return fmt.Sprintf(
		"%s: assigning result of void function '%s' is invalid",
		e.SemanticError.Error(),
		e.ident,
	)
}

// CreateVoidAssignmentError creates an error from a token and a function
// identifier
func CreateVoidAssignmentError(token *token32, ident string) error {
	return &VoidAssignmentError{
		SemanticError: CreateSemanticError(token),
		ident:         ident,
	}
}
