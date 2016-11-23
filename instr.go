package main

import (
	"fmt"
)

// Conditions

type Cond int

var condMap = map[int]string{
	0:  "EQ",
	1:  "NE",
	10: "GE",
	11: "LT",
	12: "GT",
	13: "LE",
	14: "AL",
}

const (
	EQ = 0
	NE = 1
	GE = 10
	LT = 11
	GT = 12
	LE = 13
)

func (m Cond) String() string {
	value, exists := condMap[int(m)]
	if !exists {
		return ""
	}

	return value
}

//------------------------------------------------------------------------------
//STORE AND LOAD
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

//StoreInstr struct
type StoreInstr struct {
	destination Reg
	value       Reg
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
	cond        Cond
	arg         Reg
	destination Reg
}

func (m *BaseUnaryInstr) String() string {
	return fmt.Sprintf("%s, %s",
		(m.destination).String(),
		(m.arg).String())
}

//NEGInstr struct
type NEGInstr struct {
	base BaseUnaryInstr
}

func (m *NEGInstr) String() string {
	return fmt.Sprintf("\tNEG %s", m.base.String())
}

//NOTInstr struct
type NOTInstr struct {
	base BaseUnaryInstr
}

func (m *NOTInstr) String() string {
	return fmt.Sprintf("\tNOT %s", m.base.String())
}

//------------------------------------------------------------------------------
// ARITHMETIC OPERATORS
//------------------------------------------------------------------------------

// lhs is destination REG

//BaseBinaryInstr struct
type BaseBinaryInstr struct {
	cond        Cond
	destination Reg
	lhs         Reg
	rhs         Reg
}

type Operand2 interface {
	String() string
}

type ImmediateOperand struct {
	n int
}

func (m ImmediateOperand) String() string {
	return fmt.Sprint("d", m.n)
}

func (m *BaseBinaryInstr) String() string {
	return fmt.Sprintf("%s, %s, %s",
		(m.destination).String(),
		(m.lhs).String(),
		(m.rhs).String())
}

//ADDInstr struct
type ADDInstr struct {
	base BaseBinaryInstr
}

func (m *ADDInstr) String() string {
	return fmt.Sprintf("\tADD%s %s",
		m.base.cond.String(),
		m.base.String())
}

//SUBInstr struct
type SUBInstr struct {
	base BaseBinaryInstr
}

func (m *SUBInstr) String() string {
	return fmt.Sprintf("\tSUB%s %s",
		m.base.cond.String(),
		m.base.String())
}

//RSBInstr struct
type RSBInstr struct {
	base BaseBinaryInstr
}

func (m *RSBInstr) String() string {
	return fmt.Sprintf("\tRSB%s %s",
		m.base.cond.String(),
		m.base.String())
}

//------------------------------------------------------------------------------
//COMPARISON OPERATORS
//------------------------------------------------------------------------------

type BaseComparisonInstr struct {
	cond Cond
	lhs  Reg
	rhs  int
}

func (m *BaseComparisonInstr) String() string {
	return fmt.Sprintf("%s, %d",
		(m.lhs).String(),
		m.rhs)
}

//CMPInstr struct
type CMPInstr struct {
	base BaseComparisonInstr
}

func (m *CMPInstr) String() string {
	return fmt.Sprintf("\tCMP%s %s",
		m.base.cond.String(),
		m.base.String())
}

//CMNInstr struct
type CMNInstr struct {
	base BaseComparisonInstr
}

func (m *CMNInstr) String() string {
	return fmt.Sprintf("\tCMN%s %s",
		m.base.cond.String(),
		m.base.String())
}

//TSTInstr struct
type TSTInstr struct {
	base BaseComparisonInstr
}

func (m *TSTInstr) String() string {
	return fmt.Sprintf("\tTST%s %s",
		m.base.cond.String(),
		m.base.String())
}

//TEQInstr struct
type TEQInstr struct {
	base BaseComparisonInstr
}

func (m *TEQInstr) String() string {
	return fmt.Sprintf("\tTEQ%s %s",
		m.base.cond.String(),
		m.base.String())
}

//------------------------------------------------------------------------------
//LOGICAL OPERATORS
//------------------------------------------------------------------------------

//ANDInstr struct
type ANDInstr struct {
	cond Cond
	base BaseBinaryInstr
}

func (m *ANDInstr) String() string {
	return fmt.Sprintf("\tAND %s", m.base.String())
}

//EORInstr struct
type EORInstr struct {
	base BaseBinaryInstr
}

func (m *EORInstr) String() string {
	return fmt.Sprintf("\tEOR %s", m.base.String())
}

//ORRInstr struct
type ORRInstr struct {
	base BaseBinaryInstr
}

func (m *ORRInstr) String() string {
	return fmt.Sprintf("\tORR %s", m.base.String())
}

//BICInstr struct
type BICInstr struct {
	base BaseBinaryInstr
}

func (m *BICInstr) String() string {
	return fmt.Sprintf("\tBIC %s", m.base.String())
}

//------------------------------------------------------------------------------
//DATA MOVEMENT
//------------------------------------------------------------------------------

//DataMovementInstr struct
type DataMovementInstr struct {
	destination Reg
	source      Operand2
}

func (m *DataMovementInstr) String() string {
	return fmt.Sprintf("\tMOV %s, %s", m.destination.String(), m.source.String())
}

//------------------------------------------------------------------------------
//MULTIPLICATION INSTRUCTION
//------------------------------------------------------------------------------

//MULInstr struct
type MULInstr struct {
	base BaseBinaryInstr
}

func (m *MULInstr) String() string {
	return fmt.Sprintf("\tMUL %s", m.base.String())
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
type STRInstr struct {
	source      Reg
	destination PreIndex
}

//LDRInstr struct
type LDRInstr struct {
	destination Reg
	source      PreIndex
}

func (m *STRInstr) String() string {
	return fmt.Sprintf("STR %s, [%s, %s, LSL #2]",
		m.source.String(),
		m.destination.Rn.String(),
		m.destination.Rm.String())
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

type DataWordInstr struct {
	n int
}

func (m *DataWordInstr) String() string {
	return fmt.Sprintf("\t.word %d", m.n)
}

type DataAsciiInstr struct {
	str string
}

func (m *DataAsciiInstr) String() string {
	return fmt.Sprintf("\t.ascii \"%s\"", m.str)
}
