package main

// WACC Group 34
//
// optimisation.go: Contains functions to optimise a given AST
//
// The File contains functions to optimise a given AST.

// OptimisationContext holds information that can be useful during optimisation
type OptimisationContext struct {
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
	m.body = m.body.Optimise(context)

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

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for AssignStatement
func (m *AssignStatement) Optimise(context *OptimisationContext) Statement {
	m.target = m.target.Optimise(context)
	m.rhs = m.rhs.Optimise(context)

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for ReadStatement
func (m *ReadStatement) Optimise(context *OptimisationContext) Statement {
	m.target = m.target.Optimise(context)

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

	m.trueStat = m.trueStat.Optimise(context)
	if m.falseStat != nil {
		m.falseStat = m.falseStat.Optimise(context)
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
	m.cond = m.cond.Optimise(context)

	m.body = m.body.Optimise(context)

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

	if m.next != nil {
		m.SetNext(m.next.Optimise(context))
	}

	return m
}

//Optimise optimises for DoWhileStatement
func (m *DoWhileStatement) Optimise(context *OptimisationContext) Statement {
	m.cond = m.cond.Optimise(context)

	m.body = m.body.Optimise(context)

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
	return m
}

//------------------------------------------------------------------------------
// LITERALS AND ELEMENTS OPTIMISATION
//------------------------------------------------------------------------------

//Optimise optimises for Ident
func (m *Ident) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for IntLiteral
func (m *IntLiteral) Optimise(context *OptimisationContext) Expression {
	return m
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
	return m
}

//------------------------------------------------------------------------------
// UNARY OPERATOR OPTIMISATION
//------------------------------------------------------------------------------

//Optimise optimises for UnaryOperatorNot
func (m *UnaryOperatorNot) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for UnaryOperatorNegate
func (m *UnaryOperatorNegate) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for UnaryOperatorLen
func (m *UnaryOperatorLen) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for UnaryOperatorOrd
func (m *UnaryOperatorOrd) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for UnaryOperatorChr
func (m *UnaryOperatorChr) Optimise(context *OptimisationContext) Expression {
	return m
}

//------------------------------------------------------------------------------
// BINARY OPERATOR OPTIMISATION
//------------------------------------------------------------------------------

//Optimise optimises for BinaryOperatorMult
func (m *BinaryOperatorMult) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorDiv
func (m *BinaryOperatorDiv) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorMod
func (m *BinaryOperatorMod) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorAdd
func (m *BinaryOperatorAdd) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorSub
func (m *BinaryOperatorSub) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorGreaterThan
func (m *BinaryOperatorGreaterThan) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorGreaterEqual
func (m *BinaryOperatorGreaterEqual) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorLessThan
func (m *BinaryOperatorLessThan) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorLessEqual
func (m *BinaryOperatorLessEqual) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorEqual
func (m *BinaryOperatorEqual) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorNotEqual
func (m *BinaryOperatorNotEqual) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorAnd
func (m *BinaryOperatorAnd) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorBitAnd
func (m *BinaryOperatorBitAnd) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorOr
func (m *BinaryOperatorOr) Optimise(context *OptimisationContext) Expression {
	return m
}

//Optimise optimises for BinaryOperatorBitOr
func (m *BinaryOperatorBitOr) Optimise(context *OptimisationContext) Expression {
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
		m.body = m.body.Optimise(&OptimisationContext{})
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
	m.main = m.main.Optimise(&OptimisationContext{})

	for _, ch := range chs {
		<-ch
	}
}
