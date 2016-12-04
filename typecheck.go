package main

// WACC Group 34
//
// typecheck.go: functions and auxiliary structures for typechecking
//
// Scope: structure for storing the functions, variables and their corresponging
//   types for use in the TypeCheck
// TypeCheck: recursively checks for type mismatches in AST, statements, and
//   expressions

// Scope stores the available variables, functions, and expected return type
// during lexical analysis
type Scope struct {
	parent     *Scope
	vars       map[string]Type
	funcs      map[string]*FunctionDef
	returnType Type
}

// CreateRootScope creates a global scope that has no parent
func CreateRootScope() *Scope {
	scope := &Scope{
		parent: nil,
		vars:   make(map[string]Type),
		funcs:  make(map[string]*FunctionDef),
	}

	return scope
}

// Child creates a child scope from a scope that inherits all properties but
// can declare variable independently
func (m *Scope) Child() *Scope {
	return &Scope{
		parent:     m,
		vars:       make(map[string]Type),
		funcs:      m.funcs,
		returnType: m.returnType,
	}
}

// Lookup tries to recusively search for the type of a given variable
// It returns InvalidType if not found
func (m *Scope) Lookup(ident string) Type {
	t, ok := m.vars[ident]

	if !ok {
		if m.parent != nil {
			t = m.parent.Lookup(ident)
		} else {
			t = InvalidType{}
		}
	}

	return t
}

// LookupFunction tries to return the function given it's identifier
// returns nil if not found.
func (m *Scope) LookupFunction(ident string) *FunctionDef {
	t, ok := m.funcs[ident]

	if !ok {
		return nil
	}

	return t
}

// Declare creates a new variable in the current scope returning the previous
// type in case of redeclaration, nil otherwise
func (m *Scope) Declare(ident string, t Type) Type {
	pt, ok := m.vars[ident]

	m.vars[ident] = t

	if ok {
		return pt
	}
	return nil
}

// DeclareFunction registers a new function in the scope returning the previous
// one in case of redeclaration, nil otherwise
func (m *Scope) DeclareFunction(ident string, f *FunctionDef) *FunctionDef {
	pf, ok := m.funcs[ident]

	m.funcs[ident] = f

	if ok {
		return pf
	}
	return nil
}

// Match checks whether a type is assignable to the current type
func (m InvalidType) Match(t Type) bool {
	return false
}

// Match checks whether a type is assignable to the current type
func (m UnknownType) Match(t Type) bool {
	return true
}

// Match checks whether a type is assignable to the current type
func (m IntType) Match(t Type) bool {
	switch t.(type) {
	case IntType:
		return true
	case UnknownType:
		return true
	default:
		return false
	}
}

// Match checks whether a type is assignable to the current type
func (m BoolType) Match(t Type) bool {
	switch t.(type) {
	case BoolType:
		return true
	case UnknownType:
		return true
	default:
		return false
	}
}

// Match checks whether a type is assignable to the current type
func (m CharType) Match(t Type) bool {
	switch t.(type) {
	case CharType:
		return true
	case UnknownType:
		return true
	default:
		return false
	}
}

// Match checks whether a type is assignable to the current type
func (m PairType) Match(t Type) bool {
	switch o := t.(type) {
	case PairType:
		fst := m.first.Match(o.first)
		snd := m.second.Match(o.second)
		return fst && snd
	case UnknownType:
		return true
	default:
		return false
	}
}

// Match checks whether a type is assignable to the current type
func (m ArrayType) Match(t Type) bool {
	switch o := t.(type) {
	case ArrayType:
		return m.base.Match(o.base)
	case UnknownType:
		return true
	default:
		return false
	}
}

// TypeCheck checks whether the AST has any type mismatches in expressions and
// assignments
func (m *AST) TypeCheck() []error {
	var errs []error

	errch := make(chan error)

	go func() {
		global := CreateRootScope()

		// add the functions to the scope
		for _, f := range m.functions {
			if pf := global.DeclareFunction(f.ident, f); pf != nil {
				errch <- CreateFunctionRedelarationError(
					f.Token(),
					f.ident,
				)
			}
		}

		// check the main program
		main := global.Child()
		main.returnType = InvalidType{}
		m.main.TypeCheck(main, errch)

		// check all the functions
		for _, f := range m.functions {
			fscope := global.Child()
			for _, arg := range f.params {
				pt := fscope.Declare(arg.name, arg.wtype)
				if pt != nil {
					errch <- CreateVariableRedeclarationError(
						arg.Token(),
						arg.name,
						pt,
						arg.wtype,
					)
				}
			}
			fscope.returnType = f.returnType
			f.body.TypeCheck(fscope, errch)
		}
		close(errch)
	}()

	for err := range errch {
		errs = append(errs, err)
	}

	return errs
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *BaseStatement) TypeCheck(ts *Scope, errch chan<- error) {
	if m.next == nil {
		return
	}
	m.next.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *BlockStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.body.TypeCheck(ts.Child(), errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *DeclareAssignStatement) TypeCheck(ts *Scope, errch chan<- error) {
	if pt := ts.Declare(m.ident, m.wtype); pt != nil {
		errch <- CreateVariableRedeclarationError(
			m.Token(),
			m.ident,
			pt,
			m.wtype,
		)
	}

	m.rhs.TypeCheck(ts, errch)
	if rhsT := m.rhs.Type(); !m.wtype.Match(rhsT) {
		errch <- CreateTypeMismatchError(
			m.rhs.Token(),
			m.wtype,
			rhsT,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *AssignStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.target.TypeCheck(ts, errch)
	m.rhs.TypeCheck(ts, errch)

	lhsT := m.target.Type()
	rhsT := m.rhs.Type()

	if !lhsT.Match(rhsT) {
		errch <- CreateTypeMismatchError(
			m.rhs.Token(),
			lhsT,
			rhsT,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *ReadStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.target.TypeCheck(ts, errch)

	switch t := m.target.Type().(type) {
	case IntType:
	case CharType:
	default:
		errch <- CreateTypeMismatchError(
			m.target.Token(),
			IntType{},
			t,
		)
		errch <- CreateTypeMismatchError(
			m.target.Token(),
			CharType{},
			t,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *FreeStatement) TypeCheck(ts *Scope, errch chan<- error) {

	m.expr.TypeCheck(ts, errch)
	freeT := m.expr.Type()

	switch t := freeT.(type) {
	case PairType:
	case ArrayType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			PairType{},
			t,
		)
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			ArrayType{},
			t,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *ReturnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	returnT := ts.returnType
	exprT := m.expr.Type()

	if !returnT.Match(exprT) {
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			returnT,
			exprT,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *ExitStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	exitT := m.expr.Type()

	if !(IntType{}.Match(exitT)) {
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			IntType{},
			exitT,
		)
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *PrintLnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *PrintStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *IfStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.cond.TypeCheck(ts, errch)
	boolT := m.cond.Type()

	if !(BoolType{}.Match(boolT)) {
		errch <- CreateTypeMismatchError(
			m.cond.Token(),
			BoolType{},
			boolT,
		)
	}

	m.trueStat.TypeCheck(ts.Child(), errch)
	m.falseStat.TypeCheck(ts.Child(), errch)

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *WhileStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.cond.TypeCheck(ts, errch)
	boolT := m.cond.Type()

	if !(BoolType{}.Match(boolT)) {
		errch <- CreateTypeMismatchError(
			m.cond.Token(),
			BoolType{},
			boolT,
		)
	}

	m.body.TypeCheck(ts.Child(), errch)

	m.BaseStatement.TypeCheck(ts, errch)
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *PairElemLHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch t := m.expr.Type().(type) {
	case PairType:
		if !m.snd {
			m.wtype = t.first
		} else {
			m.wtype = t.second
		}
	default:
		errch <- CreateTypeMismatchError(
			m.Token(),
			PairType{},
			t,
		)
		m.wtype = InvalidType{}
	}
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *ArrayLHS) TypeCheck(ts *Scope, errch chan<- error) {
	t := ts.Lookup(m.ident)

	for _, i := range m.index {
		i.TypeCheck(ts, errch)

		if !(IntType{}).Match(i.Type()) {
			errch <- CreateTypeMismatchError(
				i.Token(),
				IntType{},
				t,
			)
		}

		switch arr := t.(type) {
		case ArrayType:
			t = arr.base
		default:
			errch <- CreateTypeMismatchError(
				m.Token(),
				ArrayType{},
				t,
			)
		}
	}

	m.wtype = t
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *VarLHS) TypeCheck(ts *Scope, errch chan<- error) {
	t := ts.Lookup(m.ident)

	switch t.(type) {
	case InvalidType:
		errch <- CreateUndeclaredVariableError(
			m.Token(),
			m.ident,
		)
	}

	m.wtype = t
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *PairLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.fst.TypeCheck(ts, errch)
	m.snd.TypeCheck(ts, errch)
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *PairElemRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	pairT := m.expr.Type()

	switch pairT.(type) {
	case PairType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			PairType{},
			pairT,
		)
	}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *ArrayLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	if len(m.elements) == 0 {
		return
	}

	m.elements[0].TypeCheck(ts, errch)
	t := m.elements[0].Type()

	for i := 1; i < len(m.elements); i++ {
		elem := m.elements[i]
		elem.TypeCheck(ts, errch)

		if !t.Match(elem.Type()) {
			errch <- CreateTypeMismatchError(
				elem.Token(),
				t,
				elem.Type(),
			)
		}
	}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *FunctionCallRHS) TypeCheck(ts *Scope, errch chan<- error) {
	fun := ts.LookupFunction(m.ident)

	if fun == nil {
		errch <- CreateCallingNonFunctionError(
			m.Token(),
			m.ident,
		)
	}

	m.wtype = fun.returnType

	if len(fun.params) != len(m.args) {
		errch <- CreateFunctionCallWrongArityError(
			m.Token(),
			fun.ident,
			len(fun.params),
			len(m.args),
		)
	}

	for _, arg := range m.args {
		arg.TypeCheck(ts, errch)
	}

	for i := 0; i < len(fun.params) && i < len(m.args); i++ {
		paramT := fun.params[i].wtype
		argT := m.args[i].Type()
		if !paramT.Match(argT) {
			errch <- CreateTypeMismatchError(
				m.args[i].Token(),
				paramT,
				argT,
			)
		}
	}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *ExpressionRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *Ident) TypeCheck(ts *Scope, errch chan<- error) {
	t := ts.Lookup(m.ident)

	switch t.(type) {
	case InvalidType:
		errch <- CreateUndeclaredVariableError(
			m.Token(),
			m.ident,
		)
	}

	m.wtype = t
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *IntLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BoolLiteralFalse) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BoolLiteralTrue) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *CharLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *StringLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *NullPair) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *ArrayElem) TypeCheck(ts *Scope, errch chan<- error) {
	array := ts.Lookup(m.ident)

	for _, index := range m.indexes {
		index.TypeCheck(ts, errch)

		switch indexT := index.Type().(type) {
		case IntType:
		default:
			errch <- CreateTypeMismatchError(
				index.Token(),
				IntType{},
				indexT,
			)
		}

		switch arrayT := array.(type) {
		case ArrayType:
			array = arrayT.base
		default:
			errch <- CreateTypeMismatchError(
				m.Token(),
				ArrayType{},
				arrayT,
			)
		}
	}

	m.wtype = array
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorNot) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.Type().(type) {
	case BoolType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			BoolType{},
			unopT,
		)
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorNegate) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.Type().(type) {
	case IntType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			IntType{},
			unopT,
		)
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorLen) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.Type().(type) {
	case ArrayType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			ArrayType{},
			unopT,
		)
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorOrd) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.Type().(type) {
	case CharType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			CharType{},
			unopT,
		)
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorChr) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.Type().(type) {
	case IntType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			IntType{},
			unopT,
		)
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorMult) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorDiv) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorMod) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorAdd) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorSub) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorGreaterThan) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorGreaterEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorLessThan) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorLessEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckEquality(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorNotEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckEquality(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorAnd) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckBoolean(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorOr) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckBoolean(m, ts, errch)
}

func typeCheckArithmetic(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().Type()
	rhsT := m.GetRHS().Type()

	if !(lhsT.Match(rhsT)) {
		errch <- CreateTypeMismatchError(
			m.GetRHS().Token(),
			lhsT,
			rhsT,
		)
	}

	switch lhsT.(type) {
	case IntType:
	default:
		errch <- CreateTypeMismatchError(
			m.GetLHS().Token(),
			IntType{},
			lhsT,
		)
	}
}

func typeCheckComparator(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().Type()
	rhsT := m.GetRHS().Type()

	if !(lhsT.Match(rhsT)) {
		errch <- CreateTypeMismatchError(
			m.GetRHS().Token(),
			lhsT,
			rhsT,
		)
	}

	switch lhsT.(type) {
	case IntType:
	case CharType:
	default:
		errch <- CreateTypeMismatchError(
			m.GetLHS().Token(),
			IntType{},
			lhsT,
		)
	}
}

func typeCheckEquality(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().Type()
	rhsT := m.GetRHS().Type()

	if !(lhsT.Match(rhsT)) {
		errch <- CreateTypeMismatchError(
			m.GetRHS().Token(),
			lhsT,
			rhsT,
		)
	}
}

func typeCheckBoolean(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().Type()
	rhsT := m.GetLHS().Type()

	if !(lhsT.Match(rhsT)) {
		errch <- CreateTypeMismatchError(
			m.GetRHS().Token(),
			lhsT,
			rhsT,
		)
	}

	switch lhsT.(type) {
	case BoolType:
	default:
		errch <- CreateTypeMismatchError(
			m.GetLHS().Token(),
			BoolType{},
			lhsT,
		)
	}
}

// TypeCheck on ExpLPar to satisfy interface. Never called.
func (m *ExprParen) TypeCheck(ts *Scope, errch chan<- error) {
}
