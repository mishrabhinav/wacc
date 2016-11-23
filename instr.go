package main

import (
	"fmt"
)

// Conditions

//Cond type
type Cond int

var condMap = map[int]string{
	1: "EQ",
	2: "NE",
	3: "GE",
	4: "LT",
	5: "GT",
	6: "LE",
	7: "AL",
}

const (
	condEQ = 1
	condNE = 2
	condGE = 3
	condLT = 4
	condGT = 5
	condLE = 6
	condAL = 7
)

func (m Cond) String() string {
	value, exists := condMap[int(m)]
	if !exists {
		return ""
	}

	return value
}

//Shift type
type Shift int

var shiftMap = map[int]string{
	1: "LSL",
	2: "LSR",
	3: "ASR",
	4: "ROR",
}

const (
	shiftLSL = 1
	shiftLSR = 2
	shiftASR = 3
	shiftROR = 4
)

func (m Shift) String() string {
	value, exists := shiftMap[int(m)]
	if !exists {
		return ""
	}

	return value
}

//------------------------------------------------------------------------------
//LOAD / STORE INSTRUCTIONS
//------------------------------------------------------------------------------

//LoadOperand interface
type LoadOperand interface {
	String() string
}

//BasicLoadOperand struct
type BasicLoadOperand struct {
	value string
}

func (m *BasicLoadOperand) String() string {
	return fmt.Sprintf("=%s", m.value)
}

//MemoryLoadOperand struct
type MemoryLoadOperand struct {
	value int
}

func (m *MemoryLoadOperand) String() string {
	return fmt.Sprintf("[sp, #%d]", m.value)
}

//LoadInstr struct
type LoadInstr struct {
	destination Reg
	value       LoadOperand
}

func (m *LoadInstr) String() string {
	return fmt.Sprintf("%s, %s",
		(m.destination).String(),
		(m.value).String())
}

//LDRInstr struct
type LDRInstr struct {
	base LoadInstr
}

func (m *LDRInstr) String() string {
	return fmt.Sprintf("\tLDR %s", m.base.String())
}

//StoreOperand interface
type StoreOperand interface {
	String() string
}

//MemoryStoreOperand struct
type MemoryStoreOperand struct {
	value int
}

func (m *MemoryStoreOperand) String() string {
	return fmt.Sprintf("[sp, #%d]", m.value)
}

//StoreInstr struct
type StoreInstr struct {
	destination Reg
	value       StoreOperand
}

func (m *StoreInstr) String() string {
	return fmt.Sprintf("%s, %s",
		(m.destination).String(),
		(m.value).String())
}

//STRInstr struct
type STRInstr struct {
	base StoreInstr
}

func (m *STRInstr) String() string {
	return fmt.Sprintf("\tSTR %s", m.base.String())
}

//------------------------------------------------------------------------------
// UNARY OPERATORS
//------------------------------------------------------------------------------

//BaseUnaryInstr struct
type BaseUnaryInstr struct {
	cond Cond
	arg  Reg
	dest Reg
}

//NEGInstr struct
type NEGInstr struct {
	BaseUnaryInstr
}

//NOTInstr struct
type NOTInstr struct {
	BaseUnaryInstr
}

func (m *NEGInstr) String() string {
	return fmt.Sprintf("\tNEG %v, %v", m.dest, m.arg)
}

func (m *NOTInstr) String() string {
	return fmt.Sprintf("\tNOT %v, %v", m.dest, m.arg)
}

//------------------------------------------------------------------------------
// ARITHMETIC OPERATORS
//------------------------------------------------------------------------------

// lhs is destination REG

//BaseBinaryInstr struct
type BaseBinaryInstr struct {
	cond Cond
	dest Reg
	lhs  Reg
	rhs  Reg
}

//Operand2 interface
type Operand2 interface {
	String() string
}

//ImmediateOperand struct
type ImmediateOperand struct {
	n int
}

//ADDInstr struct
type ADDInstr struct {
	BaseBinaryInstr
}

//SUBInstr struct
type SUBInstr struct {
	BaseBinaryInstr
}

//RSBInstr struct
type RSBInstr struct {
	BaseBinaryInstr
}

func (m ImmediateOperand) String() string {
	return fmt.Sprint("%d", m.n)
}

func (m *ADDInstr) String() string {
	return fmt.Sprintf("\tADD%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

func (m *SUBInstr) String() string {
	return fmt.Sprintf("\tSUB%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

func (m *RSBInstr) String() string {
	return fmt.Sprintf("\tRSB%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//COMPARISON OPERATORS
//------------------------------------------------------------------------------

//BaseComparisonInstr struct
type BaseComparisonInstr struct {
	cond Cond
	lhs  Reg
	rhs  int
}

//CMPInstr struct
type CMPInstr struct {
	BaseComparisonInstr
}

//CMNInstr struct
type CMNInstr struct {
	BaseComparisonInstr
}

//TSTInstr struct
type TSTInstr struct {
	BaseComparisonInstr
}

//TEQInstr struct
type TEQInstr struct {
	BaseComparisonInstr
}

func (m *CMPInstr) String() string {
	return fmt.Sprintf("\tCMP%v %v, %d", m.cond, m.lhs, m.rhs)
}

func (m *CMNInstr) String() string {
	return fmt.Sprintf("\tCMN%v %v, %d", m.cond, m.lhs, m.rhs)
}

func (m *TSTInstr) String() string {
	return fmt.Sprintf("\tTST%v %v, %d", m.cond, m.lhs, m.rhs)
}

func (m *TEQInstr) String() string {
	return fmt.Sprintf("\tTEQ%v %v, %d", m.cond, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//LOGICAL OPERATORS
//------------------------------------------------------------------------------

//ANDInstr struct
type ANDInstr struct {
	BaseBinaryInstr
}

//EORInstr struct
type EORInstr struct {
	BaseBinaryInstr
}

//ORRInstr struct
type ORRInstr struct {
	BaseBinaryInstr
}

//BICInstr struct
type BICInstr struct {
	BaseBinaryInstr
}

func (m *ANDInstr) String() string {
	return fmt.Sprintf("\tAND%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

func (m *EORInstr) String() string {
	return fmt.Sprintf("\tEOR%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

func (m *ORRInstr) String() string {
	return fmt.Sprintf("\tORR%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

func (m *BICInstr) String() string {
	return fmt.Sprintf("\tBIC%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//DATA MOVEMENT
//------------------------------------------------------------------------------

//DataMovementInstr struct
type DataMovementInstr struct {
	dest   Reg
	source Operand2
}

func (m *DataMovementInstr) String() string {
	return fmt.Sprintf("\tMOV %v, %v", m.dest, m.source)
}

//------------------------------------------------------------------------------
//MULTIPLICATION INSTRUCTION
//------------------------------------------------------------------------------

//MULInstr struct
type MULInstr struct {
	BaseBinaryInstr
}

func (m *MULInstr) String() string {
	return fmt.Sprintf("\tMUL%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
// LOAD / STORE INSTRUCTIONS
//------------------------------------------------------------------------------

//PreIndex struct
type PreIndex struct {
	cond Cond
	Rn   Reg
	Rm   Reg
}

/*
TODO:
Check above declaration.
DEPRECATED CODE.

//STRInstr struct
type STRPreIndexInstr struct {
	source Reg
	PreIndex
}

//LDRInstr struct
type LDRPreIndexInstr struct {
	dest Reg
	PreIndex
}

func (m *STRPreIndexInstr) String() string {
	return fmt.Sprintf("STR %s, [%s, %s, LSL #2]", m.source.String(),
		m.Rn.String(), m.Rm.String())
}

func (m *LDRInstr) String() string {
	return fmt.Sprintf("STR %s, [%s, %s, LSL #2]",
		m.destination.String(),
		m.source.Rn.String(),
		m.source.Rm.String())
} */

//------------------------------------------------------------------------------
// PUSH AND POP INSTRUCTIONS
//------------------------------------------------------------------------------

//BaseStackInstr struct
type BaseStackInstr struct {
	cond Cond
	regs []Reg
}

//RegsToString .
func RegsToString(regs []Reg) string {
	//TODO Implement comment
	printedRegs := ""
	for i := 0; i < len(regs)-1; i++ {
		printedRegs += regs[i].String() + ", "
	}
	return "{" + printedRegs + "}"
}

//PUSHInstr struct
type PUSHInstr struct {
	BaseStackInstr
}

//POPInstr struct
type POPInstr struct {
	BaseStackInstr
}

func (m *PUSHInstr) String() string {
	return fmt.Sprintf("\tPUS %s", RegsToString(m.regs))
}

func (m *POPInstr) String() string {
	return fmt.Sprintf("\tPOP %s", RegsToString(m.regs))
}

//------------------------------------------------------------------------------
//LABELS
//------------------------------------------------------------------------------

//LABELInstr struct
type LABELInstr struct {
	ident string
}

func (m *LABELInstr) String() string {
	return fmt.Sprintf("%s:", m.ident)
}

//------------------------------------------------------------------------------
//BRANCH
//------------------------------------------------------------------------------

//BInstr struct
type BInstr struct {
	label string
	cond  Cond
}

func (m *BInstr) String() string {
	return fmt.Sprintf("B%s %s", m.cond.String(), m.label)
}

//------------------------------------------------------------------------------
//SEGMENTS
//------------------------------------------------------------------------------

// DataSegInstr signals the beginning of the data segment
type DataSegInstr struct{}

func (m *DataSegInstr) String() string {
	return ".data"
}

// TextSegInstr signals the beginning of the text segment
type TextSegInstr struct{}

func (m *TextSegInstr) String() string {
	return ".text"
}

// GlobalInstr exposes the argument to the linker
type GlobalInstr struct {
	label string
}

func (m *GlobalInstr) String() string {
	return fmt.Sprintf(".global %s", m.label)
}

// LTORGInstr ensures subroutines are within range of literal pools
type LTORGInstr struct {
	label string
}

func (m *LTORGInstr) String() string {
	return "\t.ltorg"
}

//------------------------------------------------------------------------------
// DATA
//------------------------------------------------------------------------------

//DataWordInstr struct
type DataWordInstr struct {
	n int
}

func (m *DataWordInstr) String() string {
	return fmt.Sprintf("\t.word %d", m.n)
}

//DataAsciiInstr type
type DataAsciiInstr struct {
	str string
}

func (m *DataAsciiInstr) String() string {
	return fmt.Sprintf("\t.ascii \"%s\"", m.str)
}
