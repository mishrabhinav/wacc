package main

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
	return (m.destination).String() +
		", " + (m.lhs).String() +
		", " + (m.rhs).String()
}

func (m *AddInstr) String() string {
	return "ADD " + m.base.String()
}

type SubInstr struct {
	base BaseBinaryInstr
}

func (m *SubInstr) String() string {
	return "SUB " + m.base.String()
}

type RSBInstr struct {
	base BaseBinaryInstr
}

func (m *RSBInstr) String() string {
	return "RSB " + m.base.String()
}

//COMPARISON OPERATORS
type CMPInstr struct {
	base BaseBinaryInstr
}

func (m *CMPInstr) String() string {
	return "CMP " + m.base.String()
}

type CMNInstr struct {
	base BaseBinaryInstr
}

func (m *CMNInstr) String() string {
	return "CMN " + m.base.String()
}

type TSTInstr struct {
	base BaseBinaryInstr
}

func (m *TSTInstr) String() string {
	return "TST " + m.base.String()
}

type TEQInstr struct {
	base BaseBinaryInstr
}

func (m *TEQInstr) String() string {
	return "TEQ " + m.base.String()
}

//LOGICAL OPERATORS
type ANDInstr struct {
	base BaseBinaryInstr
}

func (m *ANDInstr) String() string {
	return "AND " + m.base.String()
}

type EORInstr struct {
	base BaseBinaryInstr
}

func (m *EORInstr) String() string {
	return "EOR " + m.base.String()
}

type ORRInstr struct {
	base BaseBinaryInstr
}

func (m *ORRInstr) String() string {
	return "ORR " + m.base.String()
}

type BICInstr struct {
	base BaseBinaryInstr
}

func (m *BICInstr) String() string {
	return "BIC " + m.base.String()
}

//DATA MOVEMENT
type DataMovementInstr struct {
	base BaseBinaryInstr
}

func (m *DataMovementInstr) String() string {
	return "MOV " + m.base.String()
}

//MULTIPLICATION INSTRUCTION
type MULInstr struct {
	base BaseBinaryInstr
}

func (m *MULInstr) String() string {
	return "MUL " + m.base.String()
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
	return "STR " + m.source.String() +
		", [" + m.destination.Rn.String() +
		", " + m.destination.Rm.String() +
		", LSL #2]"
}

func (m *LDRInstr) String() string {
	return "LDR " + m.destination.String() +
		", [" + m.source.Rn.String() +
		", " + m.source.Rm.String() +
		", LSL #2]"
}

// PUSH AND POP INSTRUCTIONS

type BaseStackInstr struct {
	regs []Reg
}

func String(regs []Reg) string {
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
	return "PUSH, " + String(m.regs)
}

func (m *PopInstr) String() string {
	return "POP, " + String(m.regs)
}
