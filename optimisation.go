package main

// WACC Group 34
//
// optimisation.go: Contains functions to optimise a given AST
//
// The File contains functions to optimise a given AST.

// OptimisationContext holds information that can be useful during optimisation
type OptimisationContext struct {
	conditional []bool
	literals    []map[string]Expression
}

// StartScope starts a new scope where variables can be declared
func (m *OptimisationContext) StartScope() {
	m.conditional = append([]bool{false}, m.conditional...)
	m.literals = append([]map[string]Expression{nil}, m.literals...)
}

// StartCondScope starts a new scope where variables can be assigned but
// execution is conditional. All assignments escaping it taint the variables
func (m *OptimisationContext) StartCondScope() {
	m.conditional = append([]bool{true}, m.conditional...)
	m.literals = append([]map[string]Expression{nil}, m.literals...)
}

// EndScope discards the last opened scope
func (m *OptimisationContext) EndScope() {
	m.conditional = m.conditional[1:]
	m.literals = m.literals[1:]
}

// DeclareLiteral assigns a literal expression to an identifier in the scope
func (m *OptimisationContext) DeclareLiteral(ident string, expr Expression) {
	if m.literals[0] == nil {
		m.literals[0] = make(map[string]Expression)
	}
	m.literals[0][ident] = expr
}

// AssignLiteral assigns a literal expression to an identifier in some scope
func (m *OptimisationContext) AssignLiteral(ident string, expr Expression) {
	for i, lmap := range m.literals {
		if m.conditional[i] {
			expr = nil
		}
		if _, ok := lmap[ident]; ok {
			lmap[ident] = expr
		}
	}
}

// LookupLiteral looks for a literal for the identifier in the scope
// ok is true if found and has value
func (m *OptimisationContext) LookupLiteral(ident string) (Expression, bool) {
	for i, lmap := range m.literals {
		if m.conditional[i] {
			return nil, false
		}
		if expr, ok := lmap[ident]; ok {
			if expr != nil {
				return expr, ok
			}
			return nil, false
		}
	}
	return nil, false
}

//------------------------------------------------------------------------------
// Optimise
//------------------------------------------------------------------------------

// Optimise for skip statements
func (m *SkipStatement) Optimise(context *OptimisationContext) Statement {
	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m.next
}

//Optimise for block statements
func (m *BlockStatement) Optimise(context *OptimisationContext) Statement {
	context.StartScope()
	m.body = m.body.Optimise(context)
	context.EndScope()

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	switch m.body {
	case nil:
		return m.next
	default:
		return m
	}
}

//Optimise optimises for DeclareAssignStatement
func (m *DeclareAssignStatement) Optimise(context *OptimisationContext) Statement {
	m.rhs = m.rhs.Optimise(context)

	switch eRHS := m.rhs.(type) {
	case *ExpressionRHS:
		switch e := eRHS.expr.(type) {
		case *IntLiteral,
			*CharLiteral,
			*BoolLiteralTrue,
			*BoolLiteralFalse:
			context.DeclareLiteral(m.ident, e)
		default:
			context.DeclareLiteral(m.ident, nil)
		}
	default:
		context.DeclareLiteral(m.ident, nil)
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for AssignStatement
func (m *AssignStatement) Optimise(context *OptimisationContext) Statement {
	m.target = m.target.Optimise(context)
	m.rhs = m.rhs.Optimise(context)

	if eLHS, lok := m.target.(*VarLHS); lok {
		switch eRHS := m.rhs.(type) {
		case *ExpressionRHS:
			switch e := eRHS.expr.(type) {
			case *IntLiteral,
				*CharLiteral,
				*BoolLiteralTrue,
				*BoolLiteralFalse:
				context.AssignLiteral(eLHS.ident, e)
			default:
				context.AssignLiteral(eLHS.ident, nil)
			}
		default:
			context.AssignLiteral(eLHS.ident, nil)
		}
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for ReadStatement
func (m *ReadStatement) Optimise(context *OptimisationContext) Statement {
	m.target = m.target.Optimise(context)

	if eLHS, lok := m.target.(*VarLHS); lok {
		context.AssignLiteral(eLHS.ident, nil)
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for FreeStatement
func (m *FreeStatement) Optimise(context *OptimisationContext) Statement {
	m.expr = m.expr.Optimise(context)

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for ReturnStatement
func (m *ReturnStatement) Optimise(context *OptimisationContext) Statement {
	m.expr = m.expr.Optimise(context)

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for ExitStatement
func (m *ExitStatement) Optimise(context *OptimisationContext) Statement {
	m.expr = m.expr.Optimise(context)

	m.next = nil
	return m
}

//Optimise optimises for PrintLnStatement
func (m *PrintLnStatement) Optimise(context *OptimisationContext) Statement {
	m.expr = m.expr.Optimise(context)

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for PrintStatement
func (m *PrintStatement) Optimise(context *OptimisationContext) Statement {
	m.expr = m.expr.Optimise(context)

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for FunctionCallStat
func (m *FunctionCallStat) Optimise(context *OptimisationContext) Statement {
	for i, arg := range m.args {
		m.args[i] = arg.Optimise(context)
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for IfStatement
func (m *IfStatement) Optimise(context *OptimisationContext) Statement {
	m.cond = m.cond.Optimise(context)

	switch m.cond.(type) {
	case *BoolLiteralTrue:
		block := &BlockStatement{}
		block.body = m.trueStat
		block.SetNext(m.GetNext())
		return block.Optimise(context)
	case *BoolLiteralFalse:
		block := &BlockStatement{}
		block.body = m.falseStat
		block.SetNext(m.GetNext())
		return block.Optimise(context)
	}

	context.StartCondScope()
	m.trueStat = m.trueStat.Optimise(context)
	context.EndScope()
	if m.falseStat != nil {
		context.StartCondScope()
		m.falseStat = m.falseStat.Optimise(context)
		context.EndScope()
	}

	if m.trueStat != nil {
		context.StartScope()
		m.trueStat = m.trueStat.Optimise(context)
		context.EndScope()
	}
	if m.falseStat != nil {
		context.StartScope()
		m.falseStat = m.falseStat.Optimise(context)
		context.EndScope()
	}

	if m.trueStat == nil {
		m.trueStat = m.falseStat
		m.falseStat = nil
		m.cond = &UnaryOperatorNot{
			UnaryOperatorBase{expr: m.cond},
		}
		m.cond = m.cond.Optimise(context)
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	if m.trueStat == nil {
		return m.next
	}

	return m
}

//Optimise optimises for WhileStatement
func (m *WhileStatement) Optimise(context *OptimisationContext) Statement {
	context.StartCondScope()
	m.body = m.body.Optimise(context)
	m.cond = m.cond.Optimise(context)
	context.EndScope()

	if m.body != nil {
		context.StartScope()
		m.cond = m.cond.Optimise(context)
		m.body = m.body.Optimise(context)
		context.EndScope()
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	if m.body == nil {
		return m.next
	}

	switch m.cond.(type) {
	case *BoolLiteralFalse:
		return m.next
	}

	return m
}

//Optimise optimises for SwitchStatement
func (m *SwitchStatement) Optimise(context *OptimisationContext) Statement {
	m.cond = m.cond.Optimise(context)

	for i, cs := range m.cases {
		m.cases[i] = cs.Optimise(context)
	}

	for i, stm := range m.bodies {
		context.StartCondScope()
		m.bodies[i] = stm.Optimise(context)
		context.EndScope()
		if m.bodies[i] != nil {
			context.StartScope()
			m.bodies[i] = m.bodies[i].Optimise(context)
			context.EndScope()
		}
	}

	for i := 0; i < len(m.bodies); i++ {
		if m.bodies[i] == nil {
			m.bodies = append(m.bodies[:i], m.bodies[i+1:]...)
			m.cases = append(m.cases[:i], m.cases[i+1:]...)
		}
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	if len(m.bodies) == 0 {
		return m.next
	}

	return m
}

//Optimise optimises for DoWhileStatement
func (m *DoWhileStatement) Optimise(context *OptimisationContext) Statement {
	context.StartCondScope()
	m.body = m.body.Optimise(context)

	m.cond = m.cond.Optimise(context)
	context.EndScope()

	if m.body != nil {
		context.StartScope()
		m.cond = m.cond.Optimise(context)
		m.body = m.body.Optimise(context)
		context.EndScope()
	}

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	if m.body == nil {
		return m.next
	}

	switch m.cond.(type) {
	case *BoolLiteralFalse:
		return m.next
	}

	return m
}

//Optimise optimises for BreakStatement
func (m *BreakStatement) Optimise(context *OptimisationContext) Statement {
	m.next = nil

	return m
}

//Optimise optimises for ContinueStatement
func (m *ContinueStatement) Optimise(context *OptimisationContext) Statement {
	m.next = nil

	return m
}

//Optimise optimises for ForStatement
func (m *ForStatement) Optimise(context *OptimisationContext) Statement {
	context.StartScope()
	m.init = m.init.Optimise(context)
	context.StartCondScope()
	m.cond = m.cond.Optimise(context)
	m.after = m.after.Optimise(context)
	m.body = m.body.Optimise(context)
	context.EndScope()
	context.EndScope()

	context.StartScope()
	m.cond = m.cond.Optimise(context)
	m.body = m.body.Optimise(context)
	m.after = m.after.Optimise(context)
	context.EndScope()

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for PairElemLHS
func (m *PairElemLHS) Optimise(context *OptimisationContext) LHS {
	m.expr = m.expr.Optimise(context)
	return m
}

//Optimise optimises for ArrayLHS
func (m *ArrayLHS) Optimise(context *OptimisationContext) LHS {
	for i, index := range m.index {
		m.index[i] = index.Optimise(context)
	}
	return m
}

//Optimise optimises for VarLHS
func (m *VarLHS) Optimise(context *OptimisationContext) LHS {
	return m
}

//Optimise optimises for PairLiterRHS
func (m *PairLiterRHS) Optimise(context *OptimisationContext) RHS {
	m.fst = m.fst.Optimise(context)
	m.snd = m.snd.Optimise(context)
	return m
}

//Optimise optimises for ArrayLiterRHS
func (m *ArrayLiterRHS) Optimise(context *OptimisationContext) RHS {
	for i, elem := range m.elements {
		m.elements[i] = elem.Optimise(context)
	}
	return m
}

//Optimise optimises for PairElemRHS
func (m *PairElemRHS) Optimise(context *OptimisationContext) RHS {
	m.expr = m.expr.Optimise(context)
	return m
}

//Optimise optimises for FunctionCallRHS
func (m *FunctionCallRHS) Optimise(context *OptimisationContext) RHS {
	for i, arg := range m.args {
		m.args[i] = arg.Optimise(context)
	}
	return m
}

//Optimise optimises for ExpressionRHS
func (m *ExpressionRHS) Optimise(context *OptimisationContext) RHS {
	m.expr = m.expr.Optimise(context)
	return m
}

//Optimise optimises for NewInstanceRHS
func (m *NewInstanceRHS) Optimise(context *OptimisationContext) RHS {
	for i, arg := range m.args {
		m.args[i] = arg.Optimise(context)
	}
	return m
}

//------------------------------------------------------------------------------
// LITERALS AND ELEMENTS OPTIMISATION
//------------------------------------------------------------------------------

//Optimise optimises for Ident
func (m *Ident) Optimise(context *OptimisationContext) Expression {
	if liter, ok := context.LookupLiteral(m.ident); ok {
		if !m.Type().Match(liter.Type()) {
			panic("replacing wrong type")
		}
		return liter
	}
	return m
}

//Optimise optimises for IntLiteral
func (m *IntLiteral) Optimise(context *OptimisationContext) Expression {
	return m
}

func (m *EnumLiteral) Optimise(context *OptimisationContext) Expression {
	return &IntLiteral{value: m.value}
}

//Optimise optimises for BoolLiteralTrue
func (m *BoolLiteralTrue) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BoolLiteralFalse
func (m *BoolLiteralFalse) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for CharLiteral
func (m *CharLiteral) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for StringLiteral
func (m *StringLiteral) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for NullPair
func (m *NullPair) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for ArrayElem
func (m *ArrayElem) Optimise(context *OptimisationContext) Expression {
	for i, index := range m.indexes {
		m.indexes[i] = index.Optimise(context)
	}
	return m
}

//------------------------------------------------------------------------------
// UNARY OPERATOR OPTIMISATION
//------------------------------------------------------------------------------

func (m *UnaryOperatorBase) optimiseUnary(context *OptimisationContext) {
	m.expr = m.expr.Optimise(context)
}

func getIntLiter(unaryExpr UnaryOperator) (value int, ok bool) {
	val, ok := unaryExpr.GetExpression().(*IntLiteral)

	if ok {
		value = val.value
	}

	return
}

func getBoolLiter(unaryExpr UnaryOperator) (value bool, ok bool) {
	ok = true

	switch unaryExpr.GetExpression().(type) {
	case *BoolLiteralFalse:
		value = false
	case *BoolLiteralTrue:
		value = true
	default:
		ok = false
	}

	return
}

func getCharLiter(unaryExpr UnaryOperator) (value byte, ok bool) {
	val, ok := unaryExpr.GetExpression().(*CharLiteral)

	if ok {
		value = val.char[0]
	}

	return
}

//Optimise optimises for UnaryOperatorNot
func (m *UnaryOperatorNot) Optimise(context *OptimisationContext) Expression {
	m.UnaryOperatorBase.optimiseUnary(context)
	if val, ok := getBoolLiter(m); ok {
		return toWACCBool(!val)
	}
	return m
}

//Optimise optimises for UnaryOperatorNegate
func (m *UnaryOperatorNegate) Optimise(context *OptimisationContext) Expression {
	m.UnaryOperatorBase.optimiseUnary(context)
	if val, ok := getIntLiter(m); ok && in32(-val) {
		return &IntLiteral{value: -val}
	}
	return m
}

//Optimise optimises for UnaryOperatorLen
func (m *UnaryOperatorLen) Optimise(context *OptimisationContext) Expression {
	m.UnaryOperatorBase.optimiseUnary(context)
	return m
}

//Optimise optimises for UnaryOperatorOrd
func (m *UnaryOperatorOrd) Optimise(context *OptimisationContext) Expression {
	m.UnaryOperatorBase.optimiseUnary(context)
	if val, ok := getCharLiter(m); ok {
		return &IntLiteral{value: int(val)}
	}
	return m
}

//Optimise optimises for UnaryOperatorChr
func (m *UnaryOperatorChr) Optimise(context *OptimisationContext) Expression {
	m.UnaryOperatorBase.optimiseUnary(context)
	if val, ok := getCharLiter(m); ok {
		return &CharLiteral{char: string([]byte{byte(val)})}
	}
	return m
}

//------------------------------------------------------------------------------
// BINARY OPERATOR OPTIMISATION
//------------------------------------------------------------------------------

func (m *BinaryOperatorBase) optimiseBinary(context *OptimisationContext) {
	m.lhs = m.lhs.Optimise(context)
	m.rhs = m.rhs.Optimise(context)
}

func in32(value int) bool {
	return -2147483648 <= value && value <= 2147483647
}

func getIntLiters(binExpr BinaryOperator) (lhsv, rhsv int, ok bool) {
	lhs, ok1 := binExpr.GetLHS().(*IntLiteral)
	rhs, ok2 := binExpr.GetRHS().(*IntLiteral)

	ok = ok1 && ok2
	if ok {
		lhsv = lhs.value
		rhsv = rhs.value
	}

	return
}

func getBoolLiters(binExpr BinaryOperator) (lhsv, rhsv bool, ok bool) {
	ok = true

	switch binExpr.GetLHS().(type) {
	case *BoolLiteralFalse:
		lhsv = false
	case *BoolLiteralTrue:
		lhsv = true
	default:
		ok = false
	}

	switch binExpr.GetRHS().(type) {
	case *BoolLiteralFalse:
		rhsv = false
	case *BoolLiteralTrue:
		rhsv = true
	default:
		ok = false
	}

	return
}

func getCharLiters(binExpr BinaryOperator) (lhsv, rhsv byte, ok bool) {
	lhs, ok1 := binExpr.GetLHS().(*CharLiteral)
	rhs, ok2 := binExpr.GetRHS().(*CharLiteral)

	ok = ok1 && ok2
	if ok {
		lhsv = lhs.char[0]
		rhsv = rhs.char[0]
	}

	return
}

func toWACCBool(b bool) Expression {
	if b {
		return &BoolLiteralTrue{}
	}
	return &BoolLiteralFalse{}
}

//Optimise optimises for BinaryOperatorMult
func (m *BinaryOperatorMult) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok && in32(lhsv*rhsv) {
		return &IntLiteral{value: lhsv * rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorDiv
func (m *BinaryOperatorDiv) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok && rhsv != 0 {
		return &IntLiteral{value: lhsv / rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorMod
func (m *BinaryOperatorMod) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok && rhsv != 0 {
		return &IntLiteral{value: lhsv % rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorAdd
func (m *BinaryOperatorAdd) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok && in32(lhsv+rhsv) {
		return &IntLiteral{value: lhsv + rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorSub
func (m *BinaryOperatorSub) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok && in32(lhsv-rhsv) {
		return &IntLiteral{value: lhsv - rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorGreaterThan
func (m *BinaryOperatorGreaterThan) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv > rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv > rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorGreaterEqual
func (m *BinaryOperatorGreaterEqual) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv >= rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv >= rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorLessThan
func (m *BinaryOperatorLessThan) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv < rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv < rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorLessEqual
func (m *BinaryOperatorLessEqual) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv <= rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv <= rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorEqual
func (m *BinaryOperatorEqual) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv == rhsv)
	}
	if lhsv, rhsv, ok := getBoolLiters(m); ok {
		return toWACCBool(lhsv == rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv == rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorNotEqual
func (m *BinaryOperatorNotEqual) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return toWACCBool(lhsv != rhsv)
	}
	if lhsv, rhsv, ok := getBoolLiters(m); ok {
		return toWACCBool(lhsv != rhsv)
	}
	if lhsv, rhsv, ok := getCharLiters(m); ok {
		return toWACCBool(lhsv != rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorAnd
func (m *BinaryOperatorAnd) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getBoolLiters(m); ok {
		return toWACCBool(lhsv && rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorBitAnd
func (m *BinaryOperatorBitAnd) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return &IntLiteral{value: lhsv & rhsv}
	}
	return m
}

//Optimise optimises for BinaryOperatorOr
func (m *BinaryOperatorOr) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getBoolLiters(m); ok {
		return toWACCBool(lhsv || rhsv)
	}
	return m
}

//Optimise optimises for BinaryOperatorBitOr
func (m *BinaryOperatorBitOr) Optimise(context *OptimisationContext) Expression {
	m.BinaryOperatorBase.optimiseBinary(context)
	if lhsv, rhsv, ok := getIntLiters(m); ok {
		return &IntLiteral{value: lhsv | rhsv}
	}
	return m
}

//Optimise optimises for VoidExpr
func (m *VoidExpr) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for ExprParen
func (m *ExprParen) Optimise(context *OptimisationContext) Expression {
	return m
}

//------------------------------------------------------------------------------
// GENERAL OPTIMISATION UTILITY
//------------------------------------------------------------------------------

// Optimise generates instructions for functions
func (m *FunctionDef) Optimise() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		ctx := &OptimisationContext{}
		ctx.StartScope()
		m.body = m.body.Optimise(ctx)
		ctx.EndScope()
		close(ch)
	}()
	return ch
}

// Optimise generates instructions for the whole program
func (m *AST) Optimise() {
	var chs []<-chan interface{}
	for _, c := range m.classes {
		for _, m := range c.methods {
			chs = append(chs, m.Optimise())
		}
	}
	for _, f := range m.functions {
		chs = append(chs, f.Optimise())
	}
	{
		ctx := &OptimisationContext{}
		ctx.StartScope()
		m.main = m.main.Optimise(ctx)
		ctx.EndScope()
	}

	for _, ch := range chs {
		<-ch
	}
}
