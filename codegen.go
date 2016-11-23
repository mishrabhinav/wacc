package main

import (
	"bytes"
	"fmt"
	"sync"
)

//------------------------------------------------------------------------------
// RUNTIME ERRORS
//------------------------------------------------------------------------------

const (
	mPrintString      = "%.*s\\0"
	mPrintInt         = "%d\\0"
	mPrintReference   = "%p\\0"
	mNullChar         = "\\0"
	mPutChar          = "putchar"
	mPuts             = "puts"
	mTrue             = "true\\0"
	mFalse            = "false\\0"
	mFFlush           = "fflush"
	mPrintf           = "printf"
	mFreeLabel        = "free"
	mPrintStringLabel = "p_print_string"
	mExitLabel        = "exit"
	mMalloc           = "malloc"
	mThrowRuntimeErr  = "p_throw_runtime_error"
	mDivideByZeroLbl  = "p_check_divide_by_zero"
	mNullReferenceLbl = "pi_check_null_pointer"
	mOverflowLbl      = "p_throw_overflow_error"
	mArrayBoundLbl    = "p_check_array_bounds"
	mDivideByZeroErr  = "DivideByZeroError: divide or modulo by zero\\n\\0"
	mNullReferenceErr = "NullReferenceError: dereference a null reference" +
		"\\n\\0"
	mArrayNegIndexErr = "ArrayIndexOutOfBoundsError: negative index\\n\\0"
	mArrayLrgIndexErr = "ArrayIndexOutOfBoundsError: index too large\\n\\0"
	mOverflowErr      = "OverflowError: the result is too small/large to" +
		"store in a 4-byte signed-integer.\\n\\0"
)

//------------------------------------------------------------------------------
// INTERFACES
//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------
// REG ALLOCATOR
//------------------------------------------------------------------------------

// ARMGenReg is a general purpose ARM register
type ARMGenReg struct {
	r    int
	used int
}

func (m *ARMGenReg) String() string {
	return fmt.Sprintf("r%d", m.r)
}

// Reg returns the register number
func (m *ARMGenReg) Reg() int {
	return m.r
}

// ARMNamedReg is an ARM register with a specific purpose
type ARMNamedReg struct {
	r    int
	name string
}

func (m *ARMNamedReg) String() string {
	return m.name
}

// Reg returns the register number
func (m *ARMNamedReg) Reg() int {
	return m.r
}

// registers that can be used
var r0 = &ARMGenReg{r: 0}
var r1 = &ARMGenReg{r: 1}
var r2 = &ARMGenReg{r: 2}
var r3 = &ARMGenReg{r: 3}
var r4 = &ARMGenReg{r: 4}
var r5 = &ARMGenReg{r: 5}
var r6 = &ARMGenReg{r: 6}
var r7 = &ARMGenReg{r: 7}
var r8 = &ARMGenReg{r: 8}
var r9 = &ARMGenReg{r: 9}
var r10 = &ARMGenReg{r: 10}
var r11 = &ARMGenReg{r: 11}
var ip = &ARMNamedReg{name: "ip", r: 12}
var sp = &ARMNamedReg{name: "sp", r: 13}
var lr = &ARMNamedReg{name: "lr", r: 14}
var pc = &ARMNamedReg{name: "pc", r: 15}

var argRegs = []Reg{r0, r1, r2, r3}
var resReg = r0

// RegAllocator tracks register usage
type RegAllocator struct {
	stringPool   *StringPool
	fname        string
	labelCounter int
	regs         []*ARMGenReg
	stackSize    int
	stack        []map[string]int
}

// CreateRegAllocator returns an allocator initialized with all the general
// purpose registers
func CreateRegAllocator() *RegAllocator {
	return &RegAllocator{
		regs: []*ARMGenReg{
			r4, r5, r6, r7, r8, r9, r10, r11,
		},
	}
}

// GetReg returns a register that is free and ready for use
func (m *RegAllocator) GetReg(insch chan<- Instr) Reg {
	r := m.regs[0]

	if r.used > 0 {
		insch <- &PUSHInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{r},
			},
		}
		m.PushStack(4)
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

	r := re.(*ARMGenReg)

	if r.used > 1 {
		insch <- &POPInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{r},
			},
		}
		m.PopStack(4)
	}

	r.used--

	m.regs = append([]*ARMGenReg{r}, m.regs[:len(m.regs)-1]...)
}

// GetUniqueLabelSuffix returns a new unique label suffix
func (m *RegAllocator) GetUniqueLabelSuffix() string {
	defer func() {
		m.labelCounter++
	}()
	return fmt.Sprintf("_%s_%d", m.fname, m.labelCounter)
}

// PushStack increases the stack size by size
func (m *RegAllocator) PushStack(size int) {
	m.stackSize += size
}

// PopStack decreases the stack size by size
func (m *RegAllocator) PopStack(size int) {
	m.stackSize -= size
}

// DeclareVar registers a new variable for use
func (m *RegAllocator) DeclareVar(ident string, insch chan<- Instr) {
	m.PushStack(4)
	m.stack[0][ident] = m.stackSize
	insch <- &SUBInstr{
		BaseBinaryInstr: BaseBinaryInstr{
			dest: sp,
			lhs:  sp,
			rhs:  ImmediateOperand{4},
		},
	}
}

// ResolveVar returns the location of a variable
func (m *RegAllocator) ResolveVar(ident string) int {
	for _, scope := range m.stack {
		if v, ok := scope[ident]; ok {
			return (m.stackSize - v)
		}
	}

	panic(fmt.Sprintf("var %s not found in scope", ident))
}

// StartScope starts a new scope with new variable mappings possible
func (m *RegAllocator) StartScope(insch chan<- Instr) {
	m.stack = append([]map[string]int{make(map[string]int)}, m.stack...)
}

// CleanupScope starts a new scope with new variable mappings possible
func (m *RegAllocator) CleanupScope(insch chan<- Instr) {
	sl := len(m.stack[0])
	insch <- &ADDInstr{
		BaseBinaryInstr: BaseBinaryInstr{
			dest: sp,
			lhs:  sp,
			rhs:  ImmediateOperand{sl * 4},
		},
	}
	m.PopStack(sl)
	m.stack = m.stack[1:]
}

//------------------------------------------------------------------------------
// GLOBAL STRING STORAGE
//------------------------------------------------------------------------------

type DataString struct {
	len int
	str string
}

// StringPool holds the string literals that have been declared in the program
type StringPool struct {
	sync.RWMutex
	pool map[int]*DataString
}

// Lookup8 returns the msg label of a string literal
func (m *StringPool) Lookup8(msg string) string {
	m.Lock()
	defer m.Unlock()

	if m.pool == nil {
		m.pool = make(map[int]*DataString)
	}

	// TODO deduplicate strings

	l := len(m.pool)

	m.pool[l] = &DataString{len: len(msg), str: msg}

	return fmt.Sprintf("msg_%d", l)
}

// Lookup32 returns the msg label of a string literal, converted to 32 bit chars
func (m *StringPool) Lookup32(msg string) string {
	m.Lock()
	defer m.Unlock()

	if m.pool == nil {
		m.pool = make(map[int]*DataString)
	}

	// TODO deduplicate strings

	l := len(m.pool)

	var buffer bytes.Buffer

	for i := 0; i < len(msg); i++ {
		buffer.WriteString(fmt.Sprintf("%c\\0\\0\\0", msg[i]))
	}

	m.pool[l] = &DataString{len: len(msg), str: buffer.String()}

	return fmt.Sprintf("msg_%d", l)
}

//------------------------------------------------------------------------------
// CODEGEN
//------------------------------------------------------------------------------

// CodeGen base for next instruction
func (m *BaseStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	if m.next != nil {
		m.next.CodeGen(alloc, insch)
	}
}

// CodeGen for skip statements
func (m *SkipStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen for block statements
func (m *BlockStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	suffix := alloc.GetUniqueLabelSuffix()

	insch <- &LABELInstr{fmt.Sprintf("block%s", suffix)}
	alloc.StartScope(insch)

	m.body.CodeGen(alloc, insch)

	alloc.CleanupScope(insch)
	insch <- &LABELInstr{fmt.Sprintf("block_end%s", suffix)}

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for DeclareAssignStatement
func (m *DeclareAssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	// TODO: waccType
	// m.waccType <- Not caring now
	lhs := m.ident
	alloc.DeclareVar(lhs, insch)

	rhs := m.rhs

	baseReg := alloc.GetReg(insch)
	rhs.CodeGen(alloc, baseReg, insch)

	storeValue := &MemoryStoreOperand{alloc.ResolveVar(lhs)}
	insch <- &STRInstr{StoreInstr{dest: baseReg, value: storeValue}}

	alloc.FreeReg(baseReg, insch)

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for AssignStatement
func (m *AssignStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	//TODO: Implement LHS

	lhs := m.target

	rhs := m.rhs

	lhsReg := alloc.GetReg(insch)
	lhs.CodeGen(alloc, lhsReg, insch)

	rhsReg := alloc.GetReg(insch)
	rhs.CodeGen(alloc, rhsReg, insch)

	storeValue := &RegStoreOperand{lhsReg.String()}
	insch <- &STRInstr{StoreInstr{dest: rhsReg, value: storeValue}}

	alloc.FreeReg(rhsReg, insch)
	alloc.FreeReg(lhsReg, insch)

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for ReadStatement
func (m *ReadStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	//TODO

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for FreeStatement
func (m *FreeStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mNullReferenceErr)

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	insch <- &CMPInstr{BaseComparisonInstr{lhs: r0,
		rhs: &ImmediateOperand{n: 0}}}

	insch <- &LDRInstr{LoadInstr{dest: r0, cond: condEQ,
		value: &BasicLoadOperand{value: msg}}}

	insch <- &BLInstr{BInstr{cond: condEQ, label: mThrowRuntimeErr}}

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{r0}}}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &RegisterLoadOperand{reg: r0}}}

	insch <- &BLInstr{BInstr{label: mFreeLabel}}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &RegisterLoadOperand{reg: sp}}}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &RegisterLoadOperand{reg: r0, value: 4}}}

	insch <- &BLInstr{BInstr{label: mFreeLabel}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{r0}}}

	insch <- &BLInstr{BInstr{label: mFreeLabel}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{pc}}}

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for ReturnStatement
func (m *ReturnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	reg := alloc.GetReg(insch)

	m.expr.CodeGen(alloc, reg, insch)
	if r0.Reg() != reg.Reg() {
		insch <- &MOVInstr{
			dest:   r0,
			source: reg,
		}
	}

	insch <- &ADDInstr{
		BaseBinaryInstr: BaseBinaryInstr{
			dest: sp,
			lhs:  sp,
			rhs:  ImmediateOperand{alloc.stackSize},
		},
	}

	insch <- &BInstr{
		label: fmt.Sprintf("%s_return", alloc.fname),
	}

	alloc.FreeReg(reg, insch)

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for ExitStatement
func (m *ExitStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	reg := alloc.GetReg(insch)

	m.expr.CodeGen(alloc, reg, insch)

	insch <- &MOVInstr{
		dest:   r0,
		source: reg,
	}

	insch <- &BLInstr{BInstr: BInstr{label: "exit"}}

	alloc.FreeReg(reg, insch)

	m.BaseStatement.CodeGen(alloc, insch)
}

func print(m Expression, alloc *RegAllocator, insch chan<- Instr) {
	m.CodeGen(alloc, r0, insch)
	switch t := m.Type().(type) {
	case IntType:
		PrintInt(alloc, insch)
	case BoolType:
		PrintBool(alloc, insch)
	case CharType:
		PrintChar(alloc, insch)
	case PairType:
		PrintReference(alloc, insch)
	case ArrayType:
		switch t.base.(type) {
		case CharType:
			PrintString(alloc, insch)
		default:
			PrintReference(alloc, insch)
		}
	}
}

//CodeGen generates code for PrintLnStatement
func (m *PrintLnStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	print(m.expr, alloc, insch)
	msg := alloc.stringPool.Lookup32(mNullChar)
	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	insch <- &LDRInstr{LoadInstr{dest: r0, value: &BasicLoadOperand{msg}}}

	insch <- &ADDInstr{BaseBinaryInstr{dest: r0, lhs: r0, rhs: ImmediateOperand{4}}}

	insch <- &BLInstr{BInstr{label: mPuts}}

	insch <- &MOVInstr{dest: r0, source: ImmediateOperand{0}}

	insch <- &BLInstr{BInstr{label: mFFlush}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{pc}}}

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for PrintStatement
func (m *PrintStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	print(m.expr, alloc, insch)
	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for IfStatement
func (m *IfStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	suffix := alloc.GetUniqueLabelSuffix()

	labelIf := fmt.Sprintf("if%s", suffix)
	labelThen := fmt.Sprintf("then%s", suffix)
	labelElse := fmt.Sprintf("else%s", suffix)
	labelEnd := fmt.Sprintf("end%s", suffix)

	// Condition
	insch <- &LABELInstr{ident: labelIf}
	target := alloc.GetReg(insch)

	m.cond.CodeGen(alloc, target, insch)

	// CMP Check
	TruthValue := &ImmediateOperand{0}
	insch <- &CMPInstr{BaseComparisonInstr{lhs: target, rhs: TruthValue}}

	alloc.FreeReg(target, insch)

	insch <- &BInstr{label: labelElse, cond: condEQ}

	//TruthCases
	insch <- &LABELInstr{ident: labelThen}
	alloc.StartScope(insch)

	m.trueStat.CodeGen(alloc, insch)

	alloc.CleanupScope(insch)
	insch <- &BInstr{label: labelEnd}

	//FalseCases
	insch <- &LABELInstr{ident: labelElse}
	alloc.StartScope(insch)

	m.falseStat.CodeGen(alloc, insch)

	alloc.CleanupScope(insch)
	insch <- &LABELInstr{ident: labelEnd}

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for WhileStatement
func (m *WhileStatement) CodeGen(alloc *RegAllocator, insch chan<- Instr) {
	suffix := alloc.GetUniqueLabelSuffix()

	labelWhile := fmt.Sprintf("while%s", suffix)
	labelCond := fmt.Sprintf("cond%s", suffix)
	labelDo := fmt.Sprintf("do%s", suffix)
	labelEnd := fmt.Sprintf("end%s", suffix)

	// CMP Check

	insch <- &LABELInstr{ident: labelWhile}
	insch <- &BInstr{label: labelCond}

	//Body
	insch <- &LABELInstr{ident: labelDo}
	alloc.StartScope(insch)

	m.body.CodeGen(alloc, insch)

	alloc.CleanupScope(insch)

	// Condition
	insch <- &LABELInstr{ident: labelCond}
	target := alloc.GetReg(insch)

	m.cond.CodeGen(alloc, target, insch)

	alloc.FreeReg(target, insch)

	insch <- &CMPInstr{BaseComparisonInstr{lhs: target, rhs: &ImmediateOperand{1}}}

	insch <- &BInstr{label: labelDo, cond: condEQ}

	insch <- &LABELInstr{ident: labelEnd}

	m.BaseStatement.CodeGen(alloc, insch)
}

//CodeGen generates code for PairElemLHS
func (m *PairElemLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	insch <- &LDRInstr{LoadInstr{dest: r0, value: &ConstLoadOperand{8}}}
	insch <- &BLInstr{BInstr{label: mMalloc}}
}

//CodeGen generates code for ArrayLHS
func (m *ArrayLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {

	//Load array Address
	rhsVal := &ImmediateOperand{alloc.ResolveVar(m.ident)}
	insch <- &ADDInstr{BaseBinaryInstr{dest: target, lhs: sp, rhs: rhsVal}}

	//Retrieve content of Array Address
	//rhsValue := alloc.ResolveVar(m.ident)
	//insch <- &LDRInstr{LoadInstr{dest: target, value: &RegisterLoadOperand{reg: sp, value: rhsValue}}}

	//Place index in new Register
	indexReg := alloc.GetReg(insch)
	for index := 0; index < len(m.index); index++ {

		//Retrieve content of Array Address
		insch <- &LDRInstr{LoadInstr{dest: target, value: &RegisterLoadOperand{reg: target}}}

		m.index[index].CodeGen(alloc, indexReg, insch)

		//Check array Bounds
		insch <- &MOVInstr{dest: r0, source: indexReg}
		insch <- &MOVInstr{dest: r1, source: target}
		insch <- &BLInstr{BInstr{label: "p_check_array_bounds"}}

		//Target now points to the first element
		rhsVal := &ImmediateOperand{4}
		insch <- &ADDInstr{BaseBinaryInstr{dest: target, lhs: target, rhs: rhsVal}}

		//Target now points to the index element
		OpTwoRegLSL := &LSLRegOperand{reg: indexReg, offset: 2}
		insch <- &ADDInstr{BaseBinaryInstr{dest: target, lhs: target, rhs: OpTwoRegLSL}}
	}

	alloc.FreeReg(indexReg, insch)

}

//CodeGen generates code for VarLHS
func (m *VarLHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	insch <- &MOVInstr{dest: target, source: sp}
	rhsVal := &ImmediateOperand{alloc.ResolveVar(m.ident)}
	insch <- &ADDInstr{BaseBinaryInstr{dest: target, lhs: target, rhs: rhsVal}}
}

//CodeGen generates code for PairLiterRHS
func (m *PairLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	//TODO
}

//CodeGen generates code for ArrayLiterRHS
func (m *ArrayLiterRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {

	//Call Malloc
	leng := &ConstLoadOperand{len(m.elements)*4 + 4}
	insch <- &LDRInstr{LoadInstr{dest: r0, value: leng}}

	insch <- &BLInstr{BInstr{label: mMalloc}}

	insch <- &MOVInstr{dest: target, source: r0}

	//Array Pos Reg
	arrayReg := alloc.GetReg(insch)

	//Populate Heap at array indexes
	for pos := 1; pos <= len(m.elements); pos++ {
		element := m.elements[pos-1]
		element.CodeGen(alloc, arrayReg, insch)

		regOffset := &RegStoreOffsetOperand{reg: target, offset: (pos * 4)}
		insch <- &STRInstr{StoreInstr{dest: arrayReg, value: regOffset}}
	}

	alloc.FreeReg(arrayReg, insch)

	//Mov length into position 0
	lenInt := &ConstLoadOperand{len(m.elements)}
	insch <- &LDRInstr{LoadInstr{dest: arrayReg, value: lenInt}}

	insch <- &STRInstr{StoreInstr{dest: arrayReg, value: &RegStoreOperand{target.String()}}}
}

//CodeGen generates code for PairElemRHS
func (m *PairElemRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	//TODO
}

//CodeGen generates code for FunctionCallRHS
func (m *FunctionCallRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	for i := len(m.args) - 1; i >= 0; i-- {
		reg := alloc.GetReg(insch)
		m.args[i].CodeGen(alloc, reg, insch)
		insch <- &PUSHInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{reg},
			},
		}
		alloc.PushStack(4)
		alloc.FreeReg(reg, insch)
	}

	for i := 0; i < 4 && i < len(m.args); i++ {
		insch <- &POPInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{argRegs[i]},
			},
		}
	}

	insch <- &BLInstr{BInstr: BInstr{label: m.ident}}

	insch <- &MOVInstr{
		dest:   target,
		source: r0,
	}

	if pl := len(m.args); pl > 4 {
		insch <- &ADDInstr{
			BaseBinaryInstr: BaseBinaryInstr{
				dest: sp,
				lhs:  sp,
				rhs:  ImmediateOperand{(pl - 4) * 4},
			},
		}
	}

	alloc.PopStack(len(m.args) * 4)
}

//CodeGen generates code for ExpressionRHS
func (m *ExpressionRHS) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	m.expr.CodeGen(alloc, target, insch)
}

//------------------------------------------------------------------------------
// LITERALS AND ELEMENTS CODEGEN
//------------------------------------------------------------------------------

//CodeGen generates code for Ident
func (m *Ident) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	loadValue := &RegisterLoadOperand{reg: sp, value: alloc.ResolveVar(m.ident)}
	LoadInstruction := &LDRInstr{LoadInstr{dest: target, value: loadValue}}
	insch <- LoadInstruction
}

//CodeGen generates code for IntLiteral
func (m *IntLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	loadValue := &ConstLoadOperand{m.value}
	LoadInstruction := &LDRInstr{LoadInstr{dest: target, value: loadValue}}
	insch <- LoadInstruction
}

//CodeGen generates code for BoolLiteralTrue
func (m *BoolLiteralTrue) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	insch <- &MOVInstr{dest: target, source: &ImmediateOperand{1}}
}

//CodeGen generates code for BoolLiteralFalse
func (m *BoolLiteralFalse) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	insch <- &MOVInstr{dest: target, source: &ImmediateOperand{0}}
}

//CodeGen generates code for CharLiteral
func (m *CharLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	MOVInstr := &MOVInstr{dest: target, source: CharOperand{m.char}}
	insch <- MOVInstr
}

//CodeGen generates code for StringLiteral
func (m *StringLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup32(m.str)

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest:  target,
			value: &BasicLoadOperand{msg},
		},
	}
}

//CodeGen generates code for PairLiteral
func (m *PairLiteral) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	//TODO
}

//CodeGen generates code for NullPair
func (m *NullPair) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	insch <- &MOVInstr{
		dest:   target,
		source: ImmediateOperand{0},
	}
}

//CodeGen generates code for ArrayElem
func (m *ArrayElem) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	//TODO
}

//------------------------------------------------------------------------------
// UNARY OPERATOR CODEGEN
//------------------------------------------------------------------------------

//CodeGen generates code for UnaryOperatorNot
func (m *UnaryOperatorNot) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	expr := m.GetExpression()
	expr.CodeGen(alloc, target, insch)
	UnaryInstrNot := &NOTInstr{BaseUnaryInstr{arg: target, dest: target}}
	insch <- UnaryInstrNot
}

//CodeGen generates code for UnaryOperatorNegate
func (m *UnaryOperatorNegate) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	expr := m.GetExpression()
	expr.CodeGen(alloc, target, insch)
	UnaryInstrNeg := &NEGInstr{BaseUnaryInstr{arg: target, dest: target}}
	insch <- UnaryInstrNeg
}

//CodeGen generates code for UnaryOperatorLen
func (m *UnaryOperatorLen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	//TODO Implement
}

//CodeGen generates code for UnaryOperatorOrd
func (m *UnaryOperatorOrd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	m.expr.CodeGen(alloc, target, insch)
}

//CodeGen generates code for UnaryOperatorChr
func (m *UnaryOperatorChr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	m.expr.CodeGen(alloc, target, insch)
}

//------------------------------------------------------------------------------
// BINARY OPERATOR CODEGEN
//------------------------------------------------------------------------------

//CodeGen generates code for BinaryOperatorMult
// If LHS.Weight > RHS.Weight LHS is executed first
// otherwise RHS is executed first
func (m *BinaryOperatorMult) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
	}
	binaryInstrMul := &MULInstr{BaseBinaryInstr{dest: target, lhs: target,
		rhs: target2}}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrMul
}

//CodeGen generates code for BinaryOperatorDiv
func (m *BinaryOperatorDiv) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	lhsResult := target
	rhsResult := target2
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
		lhsResult = target2
		rhsResult = target
	}
	insch <- &MOVInstr{dest: r0, source: lhsResult}
	insch <- &MOVInstr{dest: r1, source: rhsResult}
	insch <- &BLInstr{BInstr: BInstr{label: mDivideByZeroLbl}}
	insch <- &BLInstr{BInstr: BInstr{label: "__aeabi_idiv"}}
	insch <- &MOVInstr{dest: target, source: r0}
	alloc.FreeReg(target2, insch)

}

//CodeGen generates code for BinaryOperatorMod
func (m *BinaryOperatorMod) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	lhsResult := target
	rhsResult := target2
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
		lhsResult = target2
		rhsResult = target
	}
	insch <- &MOVInstr{dest: r0, source: lhsResult}
	insch <- &MOVInstr{dest: r1, source: rhsResult}
	insch <- &BLInstr{BInstr: BInstr{label: mDivideByZeroLbl}}
	insch <- &BLInstr{BInstr: BInstr{label: "__aeabi_idivmod"}}
	insch <- &MOVInstr{dest: target, source: r1}
	alloc.FreeReg(target2, insch)
}

//CodeGen generates code for BinaryOperatorAdd
func (m *BinaryOperatorAdd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
	}
	binaryInstrAdd := &ADDInstr{BaseBinaryInstr{dest: target, lhs: target2,
		rhs: target}}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrAdd
}

//CodeGen generates code for BinaryOperatorSub
func (m *BinaryOperatorSub) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	var binaryInstrSub *SUBInstr
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
		binaryInstrSub = &SUBInstr{BaseBinaryInstr{dest: target,
			lhs: target, rhs: target2}}
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
		binaryInstrSub = &SUBInstr{BaseBinaryInstr{dest: target,
			lhs: target2, rhs: target}}
	}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrSub
}

func codeGenComparators(m BinaryOperator, alloc *RegAllocator, target Reg, insch chan<- Instr, condCode int) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	var binaryInstrCMP *CMPInstr
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
		binaryInstrCMP = &CMPInstr{BaseComparisonInstr{lhs: target2,
			rhs: target}}
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
		binaryInstrCMP = &CMPInstr{BaseComparisonInstr{lhs: target2,
			rhs: target}}
	}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrCMP
	insch <- &MOVInstr{cond: Cond(condCode),
		dest:   target,
		source: ImmediateOperand{1}}
	insch <- &MOVInstr{cond: Cond(condCode).getOpposite(),
		dest:   target,
		source: ImmediateOperand{0}}
}

//CodeGen generates code for BinaryOperatorGreaterThan
func (m *BinaryOperatorGreaterThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condGT)
}

//CodeGen generates code for BinaryOperatorGreaterEqual
func (m *BinaryOperatorGreaterEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condGE)
}

//CodeGen generates code for BinaryOperatorLessThan
func (m *BinaryOperatorLessThan) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condLT)
}

//CodeGen generates code for BinaryOperatorLessEqual
func (m *BinaryOperatorLessEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condLE)
}

//CodeGen generates code for BinaryOperatorEqual
func (m *BinaryOperatorEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condEQ)
}

//CodeGen generates code for BinaryOperatorNotEqual
func (m *BinaryOperatorNotEqual) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	codeGenComparators(m, alloc, target, insch, condNE)
}

//CodeGen generates code for BinaryOperatorAnd
func (m *BinaryOperatorAnd) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
	}
	binaryInstrAnd := &ANDInstr{BaseBinaryInstr{dest: target, lhs: target2,
		rhs: target}}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrAnd
}

//CodeGen generates code for BinaryOperatorOr
func (m *BinaryOperatorOr) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	var target2 Reg
	if lhs.Weight() > rhs.Weight() {
		lhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		rhs.CodeGen(alloc, target2, insch)
	} else {
		rhs.CodeGen(alloc, target, insch)
		target2 = alloc.GetReg(insch)
		lhs.CodeGen(alloc, target2, insch)
	}
	binaryInstrOrr := &ORRInstr{BaseBinaryInstr{dest: target, lhs: target2,
		rhs: target}}
	alloc.FreeReg(target2, insch)
	insch <- binaryInstrOrr
}

//CodeGen generates code for ExprParen
func (m *ExprParen) CodeGen(alloc *RegAllocator, target Reg, insch chan<- Instr) {
}

//------------------------------------------------------------------------------
// WEIGHT
//------------------------------------------------------------------------------

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
	return 1
}

//Weight returns weight of StringLiteral
func (m *StringLiteral) Weight() int {
	//TODO ??
	return 2
}

//Weight returns weight of PairLiteral
func (m *PairLiteral) Weight() int {
	return maxWeight(m.fst.Weight(), m.snd.Weight()) + 1
}

//Weight returns weight of NullPair
func (m *NullPair) Weight() int {
	return 1
}

//Weight returns weight of ArrayElem
func (m *ArrayElem) Weight() int {
	return len(m.indexes) + 1
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

func maxWeight(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func minWeight(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func binaryWeight(e1, e2 Expression) int {
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

func PrintString(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mPrintString)

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	insch <- &LDRInstr{LoadInstr{dest: r1,
		value: &RegisterLoadOperand{reg: r0}}}

	insch <- &ADDInstr{BaseBinaryInstr: BaseBinaryInstr{dest: r2, lhs: r0,
		rhs: ImmediateOperand{n: 4}}}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &BasicLoadOperand{value: msg}}}

	insch <- &ADDInstr{BaseBinaryInstr: BaseBinaryInstr{dest: r0, lhs: r0,
		rhs: ImmediateOperand{n: 4}}}

	insch <- &BLInstr{BInstr{label: mPrintf}}

	insch <- &MOVInstr{dest: r0, source: &ImmediateOperand{n: 0}}

	insch <- &BLInstr{BInstr{label: mFFlush}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{pc}}}
}

func PrintInt(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mPrintInt)

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	insch <- &MOVInstr{dest: r1, source: r0}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &BasicLoadOperand{value: msg}}}

	insch <- &ADDInstr{BaseBinaryInstr: BaseBinaryInstr{dest: r0, lhs: r0,
		rhs: ImmediateOperand{n: 4}}}

	insch <- &BLInstr{BInstr{label: mPrintf}}

	insch <- &MOVInstr{dest: r0, source: &ImmediateOperand{n: 0}}

	insch <- &BLInstr{BInstr{label: mFFlush}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{pc}}}
}

func PrintChar(alloc *RegAllocator, insch chan<- Instr) {
	insch <- &BLInstr{BInstr{label: mPutChar}}
}

func PrintBool(alloc *RegAllocator, insch chan<- Instr) {
	msg0 := alloc.stringPool.Lookup8(mTrue)
	msg1 := alloc.stringPool.Lookup8(mFalse)

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	insch <- &CMPInstr{BaseComparisonInstr{lhs: r0,
		rhs: &ImmediateOperand{n: 0}}}

	insch <- &LDRInstr{LoadInstr{dest: r0, cond: condNE,
		value: &BasicLoadOperand{value: msg0}}}

	insch <- &LDRInstr{LoadInstr{dest: r0, cond: condEQ,
		value: &BasicLoadOperand{value: msg1}}}

	insch <- &ADDInstr{BaseBinaryInstr: BaseBinaryInstr{dest: r0, lhs: r0,
		rhs: ImmediateOperand{n: 4}}}

	insch <- &BLInstr{BInstr{label: mPrintf}}

	insch <- &MOVInstr{dest: r0, source: &ImmediateOperand{n: 0}}

	insch <- &BLInstr{BInstr{label: mFFlush}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{pc}}}
}

func PrintReference(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mPrintReference)

	insch <- &PUSHInstr{BaseStackInstr{regs: []Reg{lr}}}

	nReg := alloc.GetReg(insch)

	insch <- &MOVInstr{dest: nReg, source: r0}

	insch <- &LDRInstr{LoadInstr{dest: r0,
		value: &BasicLoadOperand{value: msg}}}

	insch <- &ADDInstr{BaseBinaryInstr{
		dest: r0,
		lhs:  r0,
		rhs:  &ImmediateOperand{4}}}

	insch <- &BLInstr{BInstr{label: mPrintf}}

	insch <- &MOVInstr{dest: r0, source: ImmediateOperand{0}}

	insch <- &BLInstr{BInstr{label: mFFlush}}

	insch <- &POPInstr{BaseStackInstr{regs: []Reg{lr}}}
}

//CheckDivideByZero function
func CheckDivideByZero(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mDivideByZeroErr)

	insch <- &LABELInstr{
		ident: mDivideByZeroLbl,
	}

	insch <- &PUSHInstr{
		BaseStackInstr: BaseStackInstr{regs: []Reg{lr}},
	}

	insch <- &CMPInstr{
		BaseComparisonInstr: BaseComparisonInstr{
			lhs: r1,
			rhs: &ImmediateOperand{
				n: 0,
			},
		},
	}

	insch <- &LDRInstr{
		LoadInstr{
			dest: r0,
			cond: condEQ,
			value: &BasicLoadOperand{
				value: msg,
			},
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			cond:  condEQ,
			label: mThrowRuntimeErr,
		},
	}

	insch <- &POPInstr{
		BaseStackInstr: BaseStackInstr{
			regs: []Reg{pc},
		},
	}

}

func checkNullPointer(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mNullReferenceErr)

	insch <- &LABELInstr{
		ident: mNullReferenceLbl,
	}

	insch <- &PUSHInstr{
		BaseStackInstr: BaseStackInstr{
			regs: []Reg{lr},
		},
	}

	insch <- &CMPInstr{
		BaseComparisonInstr: BaseComparisonInstr{
			lhs: r0,
			rhs: &ImmediateOperand{
				n: 0,
			},
		},
	}

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest: r0,
			cond: condEQ,
			value: &BasicLoadOperand{
				value: msg,
			},
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			cond:  condEQ,
			label: mThrowRuntimeErr,
		},
	}

	insch <- &POPInstr{
		BaseStackInstr: BaseStackInstr{
			regs: []Reg{pc},
		},
	}
}

func checkArrayBounds(alloc *RegAllocator, insch chan<- Instr) {
	msg0 := alloc.stringPool.Lookup8(mArrayNegIndexErr)
	msg1 := alloc.stringPool.Lookup8(mArrayLrgIndexErr)

	insch <- &LABELInstr{
		ident: mArrayBoundLbl,
	}

	insch <- &PUSHInstr{
		BaseStackInstr: BaseStackInstr{
			regs: []Reg{lr},
		},
	}

	insch <- &CMPInstr{
		BaseComparisonInstr: BaseComparisonInstr{
			lhs: r0,
			rhs: &ImmediateOperand{
				n: 0,
			},
		},
	}

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest: r0,
			cond: condLT,
			value: &BasicLoadOperand{
				value: msg0,
			},
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			cond:  condLT,
			label: mThrowRuntimeErr,
		},
	}

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest: r1,
			value: &RegisterLoadOperand{
				reg: r1,
			},
		},
	}

	insch <- &CMPInstr{
		BaseComparisonInstr: BaseComparisonInstr{
			lhs: r0,
			rhs: r1,
		},
	}

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest: r0,
			cond: condCS,
			value: &BasicLoadOperand{
				value: msg1,
			},
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			cond:  condCS,
			label: mThrowRuntimeErr,
		},
	}

	insch <- &POPInstr{
		BaseStackInstr: BaseStackInstr{
			regs: []Reg{pc},
		},
	}
}

func checkOverflowUnderflow(alloc *RegAllocator, insch chan<- Instr) {
	msg := alloc.stringPool.Lookup8(mOverflowErr)

	insch <- &LABELInstr{
		ident: mOverflowLbl,
	}

	insch <- &LDRInstr{
		LoadInstr: LoadInstr{
			dest: r0,
			value: &BasicLoadOperand{
				value: msg,
			},
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			label: mThrowRuntimeErr,
		},
	}
}

func throwRuntimeError(alloc *RegAllocator, insch chan<- Instr) {
	insch <- &LABELInstr{
		ident: mThrowRuntimeErr,
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			cond:  condEQ,
			label: mPrintStringLabel,
		},
	}

	insch <- &MOVInstr{
		dest: r0,
		source: &ImmediateOperand{
			n: -1,
		},
	}

	insch <- &BLInstr{
		BInstr: BInstr{
			label: mExitLabel,
		},
	}
}

// CodeGen generates instructions for functions
func (m *FunctionDef) CodeGen(strPool *StringPool) <-chan Instr {
	ch := make(chan Instr)

	go func() {
		alloc := CreateRegAllocator()
		alloc.stringPool = strPool
		alloc.fname = m.ident

		ch <- &LABELInstr{m.ident}

		alloc.StartScope(ch)

		ch <- &PUSHInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{lr},
			},
		}

		ch <- &PUSHInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{ip},
			},
		}

		ch <- &PUSHInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{r4, r5, r6, r7, r8, r9, r10, r11},
			},
		}

		pl := len(m.params)
		switch {
		case pl >= 4:
			ch <- &PUSHInstr{
				BaseStackInstr: BaseStackInstr{
					regs: []Reg{r3},
				},
			}
			fallthrough
		case pl == 3:
			ch <- &PUSHInstr{
				BaseStackInstr: BaseStackInstr{
					regs: []Reg{r2},
				},
			}
			fallthrough
		case pl == 2:
			ch <- &PUSHInstr{
				BaseStackInstr: BaseStackInstr{
					regs: []Reg{r1},
				},
			}
			fallthrough
		case pl == 1:
			ch <- &PUSHInstr{
				BaseStackInstr: BaseStackInstr{
					regs: []Reg{r0},
				},
			}
		}

		for i := 0; i < 4 && i < len(m.params); i++ {
			p := m.params[i]
			alloc.stack[0][p.name] = i * -4
		}

		for i := 4; i < len(m.params); i++ {
			p := m.params[i]
			alloc.stack[0][p.name] = -4 + -4 + i*-4 + 8*-4
		}

		alloc.StartScope(ch)

		m.body.CodeGen(alloc, ch)

		alloc.CleanupScope(ch)

		ch <- &LABELInstr{fmt.Sprintf("%s_return", m.ident)}

		if pl := len(m.params); pl > 0 {
			ppregs := pl * 4
			if ppregs > 16 {
				ppregs = 16
			}
			ch <- &ADDInstr{
				BaseBinaryInstr: BaseBinaryInstr{
					dest: sp,
					lhs:  sp,
					rhs:  ImmediateOperand{ppregs},
				},
			}
		}

		ch <- &POPInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{r4, r5, r6, r7, r8, r9, r10, r11},
			},
		}

		ch <- &POPInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{ip},
			},
		}

		ch <- &POPInstr{
			BaseStackInstr: BaseStackInstr{
				regs: []Reg{pc},
			},
		}

		ch <- &LTORGInstr{}

		close(ch)
	}()

	return ch
}

// CodeGen generates instructions for the whole program
func (m *AST) CodeGen() <-chan Instr {
	ch := make(chan Instr)
	var charr []<-chan Instr

	strPool := &StringPool{}

	for _, f := range m.functions {
		charr = append(charr, f.CodeGen(strPool))
	}
	mainF := &FunctionDef{
		ident:      "main",
		returnType: InvalidType{},
		body:       m.main,
	}
	charr = append(charr, mainF.CodeGen(strPool))

	go func() {
		ch <- &DataSegInstr{}

		var txtInstr []Instr
		txtInstr = append(txtInstr, &TextSegInstr{})
		txtInstr = append(txtInstr, &GlobalInstr{"main"})

		for _, fch := range charr {
			for instr := range fch {
				txtInstr = append(txtInstr, instr)
			}
		}

		for i := 0; i < len(strPool.pool); i++ {
			v := strPool.pool[i]
			ch <- &LABELInstr{fmt.Sprintf("msg_%d", i)}
			ch <- &DataWordInstr{v.len}
			ch <- &DataASCIIInstr{v.str}
		}

		for _, tin := range txtInstr {
			ch <- tin
		}

		close(ch)
	}()

	return ch
}
