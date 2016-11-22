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

//CodeGen for block statements
func (m *BlockStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	alloc.StartScope()
	m.body.CodeGen(alloc, insch)
	alloc.CleanupScope()
}

//CodeGen generates code for DeclareAssignStatement
func (m *DeclareAssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for AssignStatement
func (m *AssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for ReadStatement
func (m *ReadStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for FreeStatement
func (m *FreeStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for ReturnStatement
func (m *ReturnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for ExitStatement
func (m *ExitStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for PrintLnStatement
func (m *PrintLnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for PrintStatement
func (m *PrintStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for IfStatement
func (m *IfStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for WhileStatement
func (m *WhileStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
}

//CodeGen generates code for PairElemLHS
func (m *PairElemLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for ArrayLHS
func (m *ArrayLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for VarLHS
func (m *VarLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for PairLiterRHS
func (m *PairLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for ArrayLiterRHS
func (m *ArrayLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for PairElemRHS
func (m *PairElemRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for FunctionCallRHS
func (m *FunctionCallRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for ExpressionRHS
func (m *ExpressionRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for Ident
func (m *Ident) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for IntLiteral
func (m *IntLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BoolLiteralTrue
func (m *BoolLiteralTrue) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BoolLiteralFalse
func (m *BoolLiteralFalse) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for CharLiteral
func (m *CharLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for StringLiteral
func (m *StringLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for PairLiteral
func (m *PairLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for NullPair
func (m *NullPair) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for ArrayElem
func (m *ArrayElem) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for UnaryOperatorNot
func (m *UnaryOperatorNot) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for UnaryOperatorNegate
func (m *UnaryOperatorNegate) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for UnaryOperatorLen
func (m *UnaryOperatorLen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for UnaryOperatorOrd
func (m *UnaryOperatorOrd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for UnaryOperatorChr
func (m *UnaryOperatorChr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorMult
func (m *BinaryOperatorMult) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorDiv
func (m *BinaryOperatorDiv) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorMod
func (m *BinaryOperatorMod) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorAdd
func (m *BinaryOperatorAdd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorSub
func (m *BinaryOperatorSub) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorGreaterThan
func (m *BinaryOperatorGreaterThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorGreaterEqual
func (m *BinaryOperatorGreaterEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorLessThan
func (m *BinaryOperatorLessThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorLessEqual
func (m *BinaryOperatorLessEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorEqual
func (m *BinaryOperatorEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorNotEqual
func (m *BinaryOperatorNotEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorAnd
func (m *BinaryOperatorAnd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for BinaryOperatorOr
func (m *BinaryOperatorOr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//CodeGen generates code for ExprParen
func (m *ExprParen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//Weight returns weight of Ident
func (m *Ident) Weight() int {
	return 1
}

//Weight returns weight of IntLiteral
func (m *IntLiteral) Weight() int {
	return 1
}

//Weight returns weight of BoolLiteralTrue
func (m *BoolLiteralTrue) Weight() int {
	return 1
}

//Weight returns weight of BoolLiteralFalse
func (m *BoolLiteralFalse) Weight() int {
	return 1
}

//Weight returns weight of CharLiteral
func (m *CharLiteral) Weight() int {
	//TODO
	return -1
}

//Weight returns weight of StringLiteral
func (m *StringLiteral) Weight() int {
	//TODO
	return -1
}

//Weight returns weight of PairLiteral
func (m *PairLiteral) Weight() int {
	//TODO
	return -1
}

//Weight returns weight of NullPair
func (m *NullPair) Weight() int {
	//TODO
	return -1
}

//Weight returns weight of ArrayElem
func (m *ArrayElem) Weight() int {
	//TODO
	return -1
}

//Weight returns weight of UnaryOperatorNot
func (m *UnaryOperatorNot) Weight() int {
	return m.GetExpression().Weight()
}

//Weight returns weight of UnaryOperatorNegate
func (m *UnaryOperatorNegate) Weight() int {
	return m.GetExpression().Weight()
}

//Weight returns weight of UnaryOperatorLen
func (m *UnaryOperatorLen) Weight() int {
	return m.GetExpression().Weight()
}

//Weight returns weight of UnaryOperatorOrd
func (m *UnaryOperatorOrd) Weight() int {
	return m.GetExpression().Weight()
}

//Weight returns weight of UnaryOperatorChr
func (m *UnaryOperatorChr) Weight() int {
	return m.GetExpression().Weight()
}

func maxWeight(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func minWeight(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func binaryWeight(e1 Expression, e2 Expression) int {
	cost1 := maxWeight(e1.Weight(), e2.Weight()+1)
	cost2 := maxWeight(e1.Weight()+1, e2.Weight())
	return minWeight(cost1, cost2)
}

//Weight returns weight of BinaryOperatorMult
func (m *BinaryOperatorMult) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorDiv
func (m *BinaryOperatorDiv) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorMod
func (m *BinaryOperatorMod) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorAdd
func (m *BinaryOperatorAdd) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorSub
func (m *BinaryOperatorSub) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorGreaterThan
func (m *BinaryOperatorGreaterThan) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorGreaterEqual
func (m *BinaryOperatorGreaterEqual) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorLessThan
func (m *BinaryOperatorLessThan) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorLessEqual
func (m *BinaryOperatorLessEqual) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorEqual
func (m *BinaryOperatorEqual) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorNotEqual
func (m *BinaryOperatorNotEqual) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorAnd
func (m *BinaryOperatorAnd) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of BinaryOperatorOr
func (m *BinaryOperatorOr) Weight() int {
	return binaryWeight(m.GetLHS(), m.GetRHS())
}

//Weight returns weight of ExprParen
func (m *ExprParen) Weight() int {
	//TODO
	return -1
}

//FunctionDef .
func (m *AST) FunctionDef() <-chan Instr {
	ch := make(chan Instr)

	return ch
}

//CodeGen .
func (m *AST) CodeGen() <-chan Instr {
	ch := make(chan Instr)

	return ch
}
