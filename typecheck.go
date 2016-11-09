package main

type Scope struct {
	parent *Scope
	vars   map[string]Type
	funcs  map[string]*FunctionDef
}

func CreateRootScope(ast *AST) *Scope {
	scope := &Scope{
		parent: nil,
		vars:   make(map[string]Type),
		funcs:  make(map[string]*FunctionDef),
	}

	for _, f := range ast.functions {
		scope.funcs[f.ident] = f
	}

	return scope
}

func (m *Scope) Child() *Scope {
	return &Scope{
		parent: m,
		vars:   make(map[string]Type),
		funcs:  m.funcs,
	}
}

func (m *Scope) Lookup(ident string) Type {
	t, ok := m.vars[ident]

	if !ok {
		if m.parent != nil {
			t = m.parent.Lookup(ident)
		} else {
			t = nil
		}
	}

	return t
}

func (m *Scope) LookupFunction(ident string) *FunctionDef {
	t, ok := m.funcs[ident]

	if !ok {
		return nil
	}

	return t
}

func (m *Scope) Declare(ident string, t Type) Type {
	pt, ok := m.vars[ident]

	m.vars[ident] = t

	if ok {
		return pt
	}
	return nil
}

func (m InvalidType) Match(t Type) bool {
	return false
}

func (m IntType) Match(t Type) bool {
	return false
}

func (m BoolType) Match(t Type) bool {
	return false
}

func (m CharType) Match(t Type) bool {
	return false
}

func (m PairType) Match(t Type) bool {
	return false
}

func (m ArrayType) Match(t Type) bool {
	return false
}

func (m InvalidType) String() string {
	return ""
}

func (m *AST) TypeCheck() []error {
	var errs []error

	errch := make(chan error)

	go func() {
		m.main.TypeCheck(CreateRootScope(m), errch)
		close(errch)
	}()

	for err := range errch {
		errs = append(errs, err)
	}

	return errs
}

func (m *BaseStatement) TypeCheck(ts *Scope, errch chan<- error) {
	if m.next == nil {
		return
	}
	m.next.TypeCheck(ts, errch)
}

func (m *BlockStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.body.TypeCheck(ts.Child(), errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *DeclareAssignStatement) TypeCheck(ts *Scope, errch chan<- error) {
	if pt := ts.Declare(m.ident, m.waccType); pt != nil {
		errch <- &VariableRedeclaration{
			ident: m.ident,
			prev:  pt,
			new:   m.waccType,
		}
	}

	m.rhs.TypeCheck(ts, errch)
	if rhsT := m.rhs.GetType(ts); !m.waccType.Match(rhsT) {
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *AssignStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ReadStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *FreeStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ReturnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ExitStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PrintLnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PrintStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *IfStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *WhileStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PairElemLHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *PairElemLHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ArrayLHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ArrayLHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *VarLHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *VarLHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *PairLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *PairLiterRHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *PairElemRHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *PairElemRHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ArrayLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ArrayLiterRHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *FunctionCallRHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *FunctionCallRHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ExpressionRHS) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ExpressionRHS) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *Ident) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *Ident) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *IntLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *IntLiteral) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BoolLiteralFalse) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BoolLiteralFalse) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BoolLiteralTrue) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BoolLiteralTrue) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *CharLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *CharLiteral) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *StringLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *StringLiteral) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *NullPair) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *NullPair) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ArrayElem) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ArrayElem) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *UnaryOperatorNot) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *UnaryOperatorNot) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *UnaryOperatorNegate) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *UnaryOperatorNegate) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *UnaryOperatorLen) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *UnaryOperatorLen) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *UnaryOperatorOrd) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *UnaryOperatorOrd) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *UnaryOperatorChr) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *UnaryOperatorChr) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorMult) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorMult) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorDiv) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorDiv) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorMod) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorMod) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorAdd) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorAdd) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorSub) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorSub) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorGreaterThan) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorGreaterThan) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorGreaterEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorGreaterEqual) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorLessThan) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorLessThan) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorLessEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorLessEqual) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorEqual) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorNotEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorNotEqual) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorAnd) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorAnd) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *BinaryOperatorOr) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorOr) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ExprLPar) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ExprLPar) GetType(ts *Scope) Type {
	return InvalidType{}
}

func (m *ExprRPar) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ExprRPar) GetType(ts *Scope) Type {
	return InvalidType{}
}
