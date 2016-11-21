package main

import (
	"fmt"
)

// Instr is the interface for the ARM assembly instructions
type Instr interface {
	String() string
}

// Location is either a register or memory address
type Location interface {
	String() string
}

// Reg represents a register in ARM
type Reg interface {
	Location
	Reg() int
}

// ARMReg is a specific ARM register that also tracks usage information
type ARMReg struct {
	r    int
	used int
}

func (m *ARMReg) String() string {
	return fmt.Sprintf("r%d", m.r)
}

// Reg returns the register number
func (m *ARMReg) Reg() int {
	return m.r
}

// RegAllocator tracks register usage
type RegAllocator struct {
	regs []*ARMReg
}

// GetReg returns a register that is free and ready for use
func (m *RegAllocator) GetReg(insch chan<- Instr) Reg {
	r := m.regs[0]

	if r.used > 0 {
		// TODO push register
	}

	r.used++

	m.regs = append(m.regs[1:], r)

	return r
}

// FreeReg frees a register loading back the previous value if necessary
func (m *RegAllocator) FreeReg(re Reg, insch chan<- Instr) {
	if re.Reg() != m.regs[len(m.regs)-1].Reg() {
		panic("Register free order mismatch")
	}

	r := re.(*ARMReg)

	if r.used > 1 {
		// TODO pop register
	}

	r.used--

	m.regs = append([]*ARMReg{r}, m.regs[:len(m.regs)-1]...)
}

// DeclareVar registers a new variable for use
func (m *RegAllocator) DeclareVar(ident string) {
}

// ResolveVar returns the location of a variable
func (m *RegAllocator) ResolveVar(ident string) Location {
	// TODO
	return nil
}

// StartScope starts a new scope with new variable mappings possible
func (m *RegAllocator) StartScope() {
	// TODO
}

// CleanupScope starts a new scope with new variable mappings possible
func (m *RegAllocator) CleanupScope() {
	// TODO
}

// CodeGen for skip statements
func (m *SkipStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *BlockStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *DeclareAssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *AssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *ReadStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *FreeStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *ReturnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *ExitStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *PrintLnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *PrintStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *IfStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *WhileStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

func (m *PairElemLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *ArrayLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *VarLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *PairLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *ArrayLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *PairElemRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *FunctionCallRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *ExpressionRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *Ident) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *IntLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BoolLiteralTrue) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BoolLiteralFalse) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *CharLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *StringLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *PairLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *NullPair) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *ArrayElem) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *UnaryOperatorNot) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *UnaryOperatorNegate) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *UnaryOperatorLen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *UnaryOperatorOrd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *UnaryOperatorChr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorMult) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorDiv) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorMod) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorAdd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorSub) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorGreaterThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorGreaterEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorLessThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorLessEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorNotEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorAnd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *BinaryOperatorOr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *ExprParen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

func (m *Ident) Weight() int {
	return -1
}

func (m *IntLiteral) Weight() int {
	return -1
}

func (m *BoolLiteralTrue) Weight() int {
	return -1
}

func (m *BoolLiteralFalse) Weight() int {
	return -1
}

func (m *CharLiteral) Weight() int {
	return -1
}

func (m *StringLiteral) Weight() int {
	return -1
}

func (m *PairLiteral) Weight() int {
	return -1
}

func (m *NullPair) Weight() int {
	return -1
}

func (m *ArrayElem) Weight() int {
	return -1
}

func (m *UnaryOperatorNot) Weight() int {
	return -1
}

func (m *UnaryOperatorNegate) Weight() int {
	return -1
}

func (m *UnaryOperatorLen) Weight() int {
	return -1
}

func (m *UnaryOperatorOrd) Weight() int {
	return -1
}

func (m *UnaryOperatorChr) Weight() int {
	return -1
}

func (m *BinaryOperatorMult) Weight() int {
	return -1
}

func (m *BinaryOperatorDiv) Weight() int {
	return -1
}

func (m *BinaryOperatorMod) Weight() int {
	return -1
}

func (m *BinaryOperatorAdd) Weight() int {
	return -1
}

func (m *BinaryOperatorSub) Weight() int {
	return -1
}

func (m *BinaryOperatorGreaterThan) Weight() int {
	return -1
}

func (m *BinaryOperatorGreaterEqual) Weight() int {
	return -1
}

func (m *BinaryOperatorLessThan) Weight() int {
	return -1
}

func (m *BinaryOperatorLessEqual) Weight() int {
	return -1
}

func (m *BinaryOperatorEqual) Weight() int {
	return -1
}

func (m *BinaryOperatorNotEqual) Weight() int {
	return -1
}

func (m *BinaryOperatorAnd) Weight() int {
	return -1
}

func (m *BinaryOperatorOr) Weight() int {
	return -1
}

func (m *ExprParen) Weight() int {
	return -1
}

func (m *AST) FunctionDef() <-chan Instr {
	ch := make(chan Instr)

	return ch
}

func (m *AST) CodeGen() <-chan Instr {
	ch := make(chan Instr)

	return ch
}
