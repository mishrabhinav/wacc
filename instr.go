package main

import "fmt"

// ARITHMETIC OPERATORS
// lhs is destination REG
type BaseBinaryInstr struct {
	lhs         Reg
	rhs         Reg
	destination Reg
}

type AddInstr struct {
	base BaseBinaryInstr
}

func (m *BaseBinaryInstr) String() string {
	return fmt.Sprintf("%s, %s, %s",
		(m.destination).String(),
		(m.lhs).String(),
		(m.rhs).String())
}

func (m *AddInstr) String() string {
	return fmt.Sprintf("\tADD %s", m.base.String())
}

type SubInstr struct {
	base BaseBinaryInstr
}

func (m *SubInstr) String() string {
	return fmt.Sprintf("\tSUB %s", m.base.String())
}

type RSBInstr struct {
	base BaseBinaryInstr
}

func (m *RSBInstr) String() string {
	return fmt.Sprintf("\tRSB %s", m.base.String())
}

//COMPARISON OPERATORS
type CMPInstr struct {
	base BaseBinaryInstr
}

func (m *CMPInstr) String() string {
	return fmt.Sprintf("\tCMP %s", m.base.String())
}

type CMNInstr struct {
	base BaseBinaryInstr
}

func (m *CMNInstr) String() string {
	return fmt.Sprintf("\tCMN %s", m.base.String())
}

type TSTInstr struct {
	base BaseBinaryInstr
}

func (m *TSTInstr) String() string {
	return fmt.Sprintf("\tTST %s", m.base.String())
}

type TEQInstr struct {
	base BaseBinaryInstr
}

func (m *TEQInstr) String() string {
	return fmt.Sprintf("\tTEQ %s", m.base.String())
}

//LOGICAL OPERATORS
type ANDInstr struct {
	base BaseBinaryInstr
}

func (m *ANDInstr) String() string {
	return fmt.Sprintf("\tAND %s", m.base.String())
}

type EORInstr struct {
	base BaseBinaryInstr
}

func (m *EORInstr) String() string {
	return fmt.Sprintf("\tEOR %s", m.base.String())
}

type ORRInstr struct {
	base BaseBinaryInstr
}

func (m *ORRInstr) String() string {
	return fmt.Sprintf("\tORR %s", m.base.String())
}

type BICInstr struct {
	base BaseBinaryInstr
}

func (m *BICInstr) String() string {
	return fmt.Sprintf("\tBIC %s", m.base.String())
}

//DATA MOVEMENT
type DataMovementInstr struct {
	base BaseBinaryInstr
}

func (m *DataMovementInstr) String() string {
	return fmt.Sprintf("\tMOV %s", m.base.String())
}

//MULTIPLICATION INSTRUCTION
type MULInstr struct {
	base BaseBinaryInstr
}

func (m *MULInstr) String() string {
	return fmt.Sprintf("\tMUL %s", m.base.String())
}

// LOAD / STORE INSTRUCTIONS
type PreIndex struct {
	Rn Reg
	Rm Reg
}

type STRInstr struct {
	source      Reg
	destination PreIndex
}

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

// PUSH AND POP INSTRUCTIONS

type BaseStackInstr struct {
	regs []Reg
}

func RegsToString(regs []Reg) string {
	printedRegs := ""
	for i := 0; i < len(regs)-1; i++ {
		printedRegs += regs[i].String() + ", "
	}
	return "(" + printedRegs + ")"
}

type PushInstr struct {
	BaseStackInstr
}

type PopInstr struct {
	BaseStackInstr
}

func (m *PushInstr) String() string {
	return fmt.Sprintf("\tPUS %s", RegsToString(m.regs))
}

func (m *PopInstr) String() string {
	return fmt.Sprintf("\tPOP %s", RegsToString(m.regs))
}

//LABELS

type label struct {
	ident string
}

func (m *label) String() string {
	return fmt.Sprintf("%s:", m.ident)
}
