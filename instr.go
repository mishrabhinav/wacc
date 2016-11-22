package main

import (
	"fmt"
)

//------------------------------------------------------------------------------
// UNARY OPERATORS
//------------------------------------------------------------------------------

//BaseUnaryInstr struct
type BaseUnaryInstr struct {
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
	lhs         Reg
	rhs         Reg
	destination Reg
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
	return fmt.Sprintf("\tADD %s", m.base.String())
}

//SUBInstr struct
type SUBInstr struct {
	base BaseBinaryInstr
}

func (m *SUBInstr) String() string {
	return fmt.Sprintf("\tSUB %s", m.base.String())
}

//RSBInstr struct
type RSBInstr struct {
	base BaseBinaryInstr
}

func (m *RSBInstr) String() string {
	return fmt.Sprintf("\tRSB %s", m.base.String())
}

//------------------------------------------------------------------------------
//COMPARISON OPERATORS
//------------------------------------------------------------------------------

//CMPInstr struct
type CMPInstr struct {
	base BaseBinaryInstr
}

func (m *CMPInstr) String() string {
	return fmt.Sprintf("\tCMP %s", m.base.String())
}

//CMNInstr struct
type CMNInstr struct {
	base BaseBinaryInstr
}

func (m *CMNInstr) String() string {
	return fmt.Sprintf("\tCMN %s", m.base.String())
}

//TSTInstr struct
type TSTInstr struct {
	base BaseBinaryInstr
}

func (m *TSTInstr) String() string {
	return fmt.Sprintf("\tTST %s", m.base.String())
}

//TEQInstr struct
type TEQInstr struct {
	base BaseBinaryInstr
}

func (m *TEQInstr) String() string {
	return fmt.Sprintf("\tTEQ %s", m.base.String())
}

//------------------------------------------------------------------------------
//LOGICAL OPERATORS
//------------------------------------------------------------------------------

//ANDInstr struct
type ANDInstr struct {
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
	base BaseBinaryInstr
}

func (m *DataMovementInstr) String() string {
	return fmt.Sprintf("\tMOV %s", m.base.String())
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
	Rn Reg
	Rm Reg
}

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
}

//------------------------------------------------------------------------------
// PUSH AND POP INSTRUCTIONS
//------------------------------------------------------------------------------

//BaseStackInstr struct
type BaseStackInstr struct {
	regs []Reg
}

//RegsToString .
func RegsToString(regs []Reg) string {
	//TODO Implement comment
	printedRegs := ""
	for i := 0; i < len(regs)-1; i++ {
		printedRegs += regs[i].String() + ", "
	}
	return "(" + printedRegs + ")"
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
