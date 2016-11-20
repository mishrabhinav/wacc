package main

type Instr interface {
	String() string
}

type Reg interface {
	String() string
	Reg() int
}

type RegAllocator struct {
}

func (m *RegAllocator) GetReg(insch chan<- Instr) Reg {
	return nil
}

func (m *RegAllocator) FreeReg(re Reg, insch chan<- Instr) {
}

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
