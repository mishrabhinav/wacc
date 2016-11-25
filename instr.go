package main

// WACC Group 34
//
// instr.go: Contains interfaces and structs to print instructions as strings.
//
// The File contains interfaces and structs for String() functions.

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
	8: "CS",
	9: "VS",
}

var oppCondMap = map[int]int{
	1: condNE,
	2: condEQ,
	3: condLT,
	4: condGE,
	5: condLE,
	6: condGT,
	7: condAL,
	8: condCS,
	9: condVS,
}

const (
	condEQ = 1
	condNE = 2
	condGE = 3
	condLT = 4
	condGT = 5
	condLE = 6
	condAL = 7
	condCS = 8
	condVS = 9
)

// Returns String representation Cond given,
// uses map
func (m Cond) String() string {
	value, exists := condMap[int(m)]
	if !exists {
		return ""
	}

	return value
}

// Returns String representation of opposite of Cond,
// uses map
func (m Cond) getOpposite() Cond {
	value, exists := oppCondMap[int(m)]
	if !exists {
		return Cond(condAL)
	}
	return Cond(value)
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

// Returns strinf representation of Shift value given
// uses map
func (m Shift) String() string {
	value, exists := shiftMap[int(m)]
	if !exists {
		return ""
	}

	return value
}

//------------------------------------------------------------------------------
// UNARY OPERATORS
//------------------------------------------------------------------------------

//BaseUnaryInstr struct
// --> (UNARYINSTRUCTION)(COND) dest, arg
type BaseUnaryInstr struct {
	cond Cond
	arg  Reg
	dest Reg
}

//NEGInstr struct
// --> NEGS dest, arg
type NEGInstr struct {
	BaseUnaryInstr
}

//NOTInstr struct
// --> EOR dest, arg
type NOTInstr struct {
	BaseUnaryInstr
}

// Returns string representation of NEGString struct
// --> NEG dest, arg
func (m *NEGInstr) String() string {
	return fmt.Sprintf("\tNEGS %v, %v", m.dest, m.arg)
}

// Returns string representation of NEGString struct
// --> EOR dest, arg
func (m *NOTInstr) String() string {
	return fmt.Sprintf("\tEOR %v, %v, #1", m.dest, m.arg)
}

//------------------------------------------------------------------------------
// ARITHMETIC OPERATORS
//------------------------------------------------------------------------------

// lhs is destination REG

//BaseBinaryInstr struct
// (BINARYINSTR)(COND) dest, lhs, rhs
type BaseBinaryInstr struct {
	cond Cond
	dest Reg
	lhs  Reg
	rhs  Operand2
}

//Operand2 interface
type Operand2 interface {
	String() string
}

//ImmediateOperand struct
type ImmediateOperand struct {
	n int
}

//RegisterOperand struct
//-->reg, shift, #amount
type RegisterOperand struct {
	reg    Reg
	shift  Shift
	amount int
}

//CharOperand struct
type CharOperand struct {
	char string
}

//ADDInstr struct
//--> ADD(COND) dest, lhs, rhs
type ADDInstr struct {
	BaseBinaryInstr
}

//SUBInstr struct
//-->SUB(COND) dest, lhs, rhs
type SUBInstr struct {
	BaseBinaryInstr
}

//RSBInstr struct
//-->RSB(COND) dest, lhs, rhs
type RSBInstr struct {
	BaseBinaryInstr
}

// Returns String representation of ImmediateOperand
//--> #n
func (m ImmediateOperand) String() string {
	return fmt.Sprintf("#%d", m.n)
}

// Returns String representation of RegisterOperand
// --> reg, shift, #amount
func (m RegisterOperand) String() string {
	if m.shift > 0 {
		return fmt.Sprintf("%v, %v #%d", m.reg, m.shift, m.amount)
	}

	return fmt.Sprintf("%v", m.reg)
}

// Returns the String representation of CharOperand
//--> #char
func (m CharOperand) String() string {
	return fmt.Sprintf("#'%s'", m.char)
}

// Returns the String representation of the ADD Instruction
//--> ADD(COND) dest, lhs, rhs
func (m *ADDInstr) String() string {
	return fmt.Sprintf("\tADDS%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

// Returns the String representation of the SUB Instruction
//--> SUB(COND) dest, lhs, rhs
func (m *SUBInstr) String() string {
	return fmt.Sprintf("\tSUBS%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

// Returns the String representation of the RSB Instruction
//--> RSB(COND) dest, lhs, rhs
func (m *RSBInstr) String() string {
	return fmt.Sprintf("\tRSBS%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//COMPARISON OPERATORS
//------------------------------------------------------------------------------

//BaseComparisonInstr struct
// -->CMP(COND) lhs, rhs
type BaseComparisonInstr struct {
	cond Cond
	lhs  Reg
	rhs  Operand2
}

//CMPInstr struct
//-->CMP(COND) lhs, rhs
type CMPInstr struct {
	BaseComparisonInstr
}

//CMNInstr struct
//-->CMN(COND) lhs, rhs
type CMNInstr struct {
	BaseComparisonInstr
}

//TSTInstr struct
//-->TST(COND) lhs, rhs
type TSTInstr struct {
	BaseComparisonInstr
}

//TEQInstr struct
//-->TEQ(COND) lhs, rhs
type TEQInstr struct {
	BaseComparisonInstr
}

//Returns String representation of CMPInstr given
//-->CMP(COND) lhs, rhs
func (m *CMPInstr) String() string {
	return fmt.Sprintf("\tCMP%v %v, %s", m.cond, m.lhs, m.rhs)
}

//Returns String representation of CMPInstr given
//-->CMP(COND) lhs, rhs
func (m *CMNInstr) String() string {
	return fmt.Sprintf("\tCMN%v %v, %s", m.cond, m.lhs, m.rhs)
}

//Returns String representation of CMPInstr given
//-->CMP(COND) lhs, rhs
func (m *TSTInstr) String() string {
	return fmt.Sprintf("\tTST%v %v, %s", m.cond, m.lhs, m.rhs)
}

//Returns String representation of CMPInstr given
//-->CMP(COND) lhs, rhs
func (m *TEQInstr) String() string {
	return fmt.Sprintf("\tTEQ%v %v, %s", m.cond, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//LOGICAL OPERATORS
//------------------------------------------------------------------------------

//ANDInstr struct
//--> AND(COND) dest, lhs, rhs
type ANDInstr struct {
	BaseBinaryInstr
}

//EORInstr struct
//--> EOR(COND) dest, lhs, rhs
type EORInstr struct {
	BaseBinaryInstr
}

//ORRInstr struct
//--> OOR(COND) dest, lhs, rhs
type ORRInstr struct {
	BaseBinaryInstr
}

//BICInstr struct
//--> BIC(COND) dest, lhs, rhs
type BICInstr struct {
	BaseBinaryInstr
}

//Returns string representation of ANDInstr given
//--> AND(COND) dest, lhs, rhs
func (m *ANDInstr) String() string {
	return fmt.Sprintf("\tAND%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//Returns string representation of EORInstr given
//--> EOR(COND) dest, lhs, rhs
func (m *EORInstr) String() string {
	return fmt.Sprintf("\tEOR%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//Returns string representation of ORRInstr given
//--> ORR(COND) dest, lhs, rhs
func (m *ORRInstr) String() string {
	return fmt.Sprintf("\tORR%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//Returns string representation of BICInstr given
//--> BIC(COND) dest, lhs, rhs
func (m *BICInstr) String() string {
	return fmt.Sprintf("\tBIC%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//------------------------------------------------------------------------------
//DATA MOVEMENT
//------------------------------------------------------------------------------

//MOVInstr struct
//--> MOV(COND) dest, source
type MOVInstr struct {
	cond   Cond
	dest   Reg
	source Operand2
}

//Returns string representation of MOVInstr given
//--> MOV(COND) dest, source
func (m *MOVInstr) String() string {
	return fmt.Sprintf("\tMOV%v %v, %v", m.cond, m.dest, m.source)
}

//------------------------------------------------------------------------------
//MULTIPLICATION INSTRUCTION
//------------------------------------------------------------------------------

//MULInstr struct
//MUL(COND) dest, lhs, rhs
type MULInstr struct {
	BaseBinaryInstr
}

//Returns String representation of MULInstr given
//MUL(COND) dest, lhs, rhs
func (m *MULInstr) String() string {
	return fmt.Sprintf("\tMUL%v %v, %v, %v", m.cond, m.dest, m.lhs, m.rhs)
}

//SMULLInstr struct
// SMULL(COND) RdLo, RdHi, Rm, Rs
type SMULLInstr struct {
	cond Cond
	RdLo Reg
	RdHi Reg
	Rm   Reg
	Rs   Reg
}

//Returns the string representation of the SMULLInstr given
// SMULL(COND) RdLo, RdHi, Rm, Rs
func (m *SMULLInstr) String() string {
	return fmt.Sprintf("\tSMULL%v %v, %v, %v, %v", m.cond, m.RdLo, m.RdHi, m.Rm, m.Rs)
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

//LoadOperand interface
type LoadOperand interface {
	String() string
}

//BasicLoadOperand struct
//--> =value
//where value is a string
type BasicLoadOperand struct {
	value string
}

//Returns the string representation of BasicLoadOperand given
// --> =value
// where value is a string
func (m *BasicLoadOperand) String() string {
	return fmt.Sprintf("=%s", m.value)
}

//ConstLoadOperand struct
// --> =value
// where value is an int
type ConstLoadOperand struct {
	value int
}

//Returns the string representation of ConstLoadOperand given
// --> =value
// where value is an int
func (m *ConstLoadOperand) String() string {
	return fmt.Sprintf("=%d", m.value)
}

// RegisterLoadOperand struct
type RegisterLoadOperand struct {
	value int
	reg   Reg
}

//Returns the string representation of RegisterLoadOperand given
//--> [reg]
func (m *RegisterLoadOperand) String() string {
	if m.value == 0 {
		return fmt.Sprintf("[%v]", m.reg)
	}
	return fmt.Sprintf("[%v, #%d]", m.reg, m.value)
}

//LoadInstr struct
type LoadInstr struct {
	reg   Reg
	value LoadOperand
	cond  Cond
}

//LDRInstr struct
type LDRInstr struct {
	LoadInstr
}

// Returns the string representation of the LDRInstr given
// --> LDR(COND) reg, value
func (m *LDRInstr) String() string {
	return fmt.Sprintf("\tLDR%v %v, %v", m.cond, m.reg, m.value)
}

//StoreOperand interface
type StoreOperand interface {
	String() string
}

//MemoryStoreOperand struct
type MemoryStoreOperand struct {
	value int
}

//Returns string representation of MemoryStoreOperand given
// --> [sp, #value]
func (m *MemoryStoreOperand) String() string {
	if m.value == 0 {
		return fmt.Sprintf("[sp]")
	}
	return fmt.Sprintf("[sp, #%d]", m.value)
}

//RegStoreOperand struct
type RegStoreOperand struct {
	reg Reg
}

//Returns string representation of RegStoreOperand given
//-->[reg]
func (m *RegStoreOperand) String() string {
	return fmt.Sprintf("[%v]", m.reg)
}

//RegStoreOffsetOperand struct
type RegStoreOffsetOperand struct {
	reg    Reg
	offset int
}

// Returns string representation of RegStoreOffsetOperand given
// --> [reg, #offset]
func (m *RegStoreOffsetOperand) String() string {
	return fmt.Sprintf("[%v, #%d]", m.reg, m.offset)
}

//StoreInstr struct
type StoreInstr struct {
	reg   Reg
	value StoreOperand
}

// Returns the string representation of the StoreInstr given
// --> reg, value
func (m *StoreInstr) String() string {
	return fmt.Sprintf("%s, %s",
		(m.reg).String(),
		(m.value).String())
}

//STRInstr struct
type STRInstr struct {
	base StoreInstr
}

// Returns the string representation of the STRInstr given
// --> STR, base
func (m *STRInstr) String() string {
	return fmt.Sprintf("\tSTR %s", m.base.String())
}

//------------------------------------------------------------------------------
// PUSH AND POP INSTRUCTIONS
//------------------------------------------------------------------------------

//BaseStackInstr struct
//-->PUSH/POP(COND) {regs}
type BaseStackInstr struct {
	cond Cond
	regs []Reg
}

//RegsToString returns the string representation
// of a list of registers
func RegsToString(regs []Reg) string {
	var printedRegs string
	if len(regs) == 1 {
		printedRegs = regs[0].String()
	} else {
		for i := 0; i < len(regs)-1; i++ {
			printedRegs += regs[i].String() + ", "
		}
		printedRegs += regs[len(regs)-1].String()
	}
	return "{" + printedRegs + "}"
}

//PUSHInstr struct
//--> PUSH(COND) regs
type PUSHInstr struct {
	BaseStackInstr
}

//POPInstr struct
//-->POP(COND) {regs}
type POPInstr struct {
	BaseStackInstr
}

//Returns string representation of PUSHInstr given
//--> PUSH(COND) regs
func (m *PUSHInstr) String() string {
	return fmt.Sprintf("\tPUSH %s", RegsToString(m.regs))
}

//Returns string representation of POPInstr given
//--> POP(COND) {regs}
func (m *POPInstr) String() string {
	return fmt.Sprintf("\tPOP %s", RegsToString(m.regs))
}

//------------------------------------------------------------------------------
//LABELS
//------------------------------------------------------------------------------

//LABELInstr struct
//--> label:
type LABELInstr struct {
	ident string
}

func (m *LABELInstr) String() string {
	return fmt.Sprintf("%s:", m.ident)
}

//------------------------------------------------------------------------------
//SHIFTED OPERANDS
//------------------------------------------------------------------------------

//LSLRegOperand struct
type LSLRegOperand struct {
	reg    Reg
	offset int
}

// Returns the string representation of LSLRegOperand
// --> reg, LSL #offset
func (m *LSLRegOperand) String() string {
	return fmt.Sprintf("%v, LSL #%d", m.reg, m.offset)
}

//------------------------------------------------------------------------------
//BRANCH
//------------------------------------------------------------------------------

//BInstr struct
//--> B(COND) label
type BInstr struct {
	label string
	cond  Cond
}

//BLInstr struct
//--> BL(COND) label
type BLInstr struct {
	BInstr
}

//Returns the string representation of BInstr given
// --> B(COND) label
func (m *BInstr) String() string {
	return fmt.Sprintf("\tB%s %s", m.cond.String(), m.label)
}

//Returns the string representation of BLInstr given
//--> BL(COND) label
func (m *BLInstr) String() string {
	return fmt.Sprintf("\tBL%s %s", m.cond.String(), m.label)
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
type LTORGInstr struct{}

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

//DataASCIIInstr type
type DataASCIIInstr struct {
	str string
}

func (m *DataASCIIInstr) String() string {
	return fmt.Sprintf("\t.ascii \"%s\"", m.str)
}
