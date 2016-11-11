package main

// WACC Group 34
//
// typecheck.go: functions and auxiliary structures for typechecking
//
// Scope: structure for storing the functions, variables and their corresponging
//   types for use in the TypeCheck and GetType
// TypeCheck: recursively checks for type mismatches in AST, statements, and
//   expressions
// GetType: walks the expression tree to deduce the type of a particular
//   expression

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
				pt := fscope.Declare(arg.name, arg.waccType)
				if pt != nil {
					errch <- CreateVariableRedeclarationError(
						arg.Token(),
						arg.name,
						pt,
						arg.waccType,
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
	if pt := ts.Declare(m.ident, m.waccType); pt != nil {
		errch <- CreateVariableRedeclarationError(
			m.Token(),
			m.ident,
			pt,
			m.waccType,
		)
	}

	m.rhs.TypeCheck(ts, errch)
	if rhsT := m.rhs.GetType(ts); !m.waccType.Match(rhsT) {
		errch <- CreateTypeMismatchError(
			m.rhs.Token(),
			m.waccType,
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

	lhsT := m.target.GetType(ts)
	rhsT := m.rhs.GetType(ts)

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

	switch t := m.target.GetType(ts).(type) {
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
}

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *FreeStatement) TypeCheck(ts *Scope, errch chan<- error) {

	m.expr.TypeCheck(ts, errch)
	freeT := m.expr.GetType(ts)

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
	exprT := m.expr.GetType(ts)

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
	exitT := m.expr.GetType(ts)

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
	boolT := m.cond.GetType(ts)

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
	boolT := m.cond.GetType(ts)

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

	switch t := m.expr.GetType(ts).(type) {
	case PairType:
	default:
		errch <- CreateTypeMismatchError(
			m.Token(),
			PairType{},
			t,
		)
	}
}

// GetType returns the deduced type of the left hand side assignment target
// InvalidType is returned in case of mismatch
func (m *PairElemLHS) GetType(ts *Scope) Type {
	switch t := m.expr.GetType(ts).(type) {
	case PairType:
		if !m.snd {
			return t.first
		}
		return t.second
	default:
		return InvalidType{}
	}
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *ArrayLHS) TypeCheck(ts *Scope, errch chan<- error) {
	t := ts.Lookup(m.ident)

	for _, i := range m.index {
		i.TypeCheck(ts, errch)

		if !(IntType{}).Match(i.GetType(ts)) {
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
}

// GetType returns the deduced type of the left hand side assignment target
// InvalidType is returned in case of mismatch
func (m *ArrayLHS) GetType(ts *Scope) Type {
	t := ts.Lookup(m.ident)

	for _, i := range m.index {
		switch i.GetType(ts).(type) {
		case IntType:
		default:
			return InvalidType{}
		}

		switch arr := t.(type) {
		case ArrayType:
			t = arr.base
		default:
			return InvalidType{}
		}
	}

	return t
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *VarLHS) TypeCheck(ts *Scope, errch chan<- error) {
	switch ts.Lookup(m.ident).(type) {
	case InvalidType:
		errch <- CreateUndelaredVariableError(
			m.Token(),
			m.ident,
		)
	}
}

// GetType returns the deduced type of the left hand side assignment target.
// InvalidType is returned in case of mismatch.
func (m *VarLHS) GetType(ts *Scope) Type {
	return ts.Lookup(m.ident)
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *PairLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.fst.TypeCheck(ts, errch)
	m.snd.TypeCheck(ts, errch)
}

// GetType returns the deduced type of the right hand side assignment source.
// InvalidType is returned in case of mismatch.
func (m *PairLiterRHS) GetType(ts *Scope) Type {
	fstT := m.fst.GetType(ts)
	sndT := m.snd.GetType(ts)

	switch fstT.(type) {
	case InvalidType:
		return InvalidType{}
	}

	switch sndT.(type) {
	case InvalidType:
		return InvalidType{}
	}

	return PairType{first: fstT, second: sndT}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *PairElemRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	pairT := m.expr.GetType(ts)

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

// GetType returns the deduced type of the right hand side assignment source.
// InvalidType is returned in case of mismatch.
func (m *PairElemRHS) GetType(ts *Scope) Type {
	switch t := m.expr.GetType(ts).(type) {
	case PairType:
		if !m.snd {
			return t.first
		}
		return t.second
	default:
		return InvalidType{}
	}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *ArrayLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	if len(m.elements) == 0 {
		return
	}

	t := m.elements[0].GetType(ts)

	for _, elem := range m.elements {
		elem.TypeCheck(ts, errch)

		if !t.Match(elem.GetType(ts)) {
			errch <- CreateTypeMismatchError(
				elem.Token(),
				t,
				elem.GetType(ts),
			)
		}
	}
}

// GetType returns the deduced type of the right hand side assignment source.
// InvalidType is returned in case of mismatch.
func (m *ArrayLiterRHS) GetType(ts *Scope) Type {
	if len(m.elements) == 0 {
		return ArrayType{base: UnknownType{}}
	}

	t := m.elements[0].GetType(ts)

	for _, elem := range m.elements {
		if !t.Match(elem.GetType(ts)) {
			return InvalidType{}
		}
	}

	return ArrayType{t}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *FunctionCallRHS) TypeCheck(ts *Scope, errch chan<- error) {
	fun := ts.LookupFunction(m.ident)

	if fun == nil {
		errch <- &CallingNonFunction{
			ident: m.ident,
		}
	}

	if len(fun.params) != len(m.args) {
		errch <- &FunctionCallWrongArity{
			ident:    fun.ident,
			expected: len(fun.params),
			got:      len(m.args),
		}
	}

	for _, arg := range m.args {
		arg.TypeCheck(ts, errch)
	}

	for i := 0; i < len(fun.params) && i < len(m.args); i++ {
		paramT := fun.params[i].waccType
		argT := m.args[i].GetType(ts)
		if !paramT.Match(argT) {
			errch <- CreateTypeMismatchError(
				m.args[i].Token(),
				paramT,
				argT,
			)
		}
	}
}

// GetType returns the deduced type of the right hand side assignment source.
// InvalidType is returned in case of mismatch.
func (m *FunctionCallRHS) GetType(ts *Scope) Type {
	fun := ts.LookupFunction(m.ident)

	if fun == nil {
		return InvalidType{}
	}

	if len(fun.params) != len(m.args) {
		return InvalidType{}
	}

	for i := 0; i < len(fun.params) && i < len(m.args); i++ {
		paramT := fun.params[i].waccType
		argT := m.args[i].GetType(ts)
		if !paramT.Match(argT) {
			return InvalidType{}
		}
	}

	return fun.returnType
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *ExpressionRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
}

// GetType returns the deduced type of the right hand side assignment source.
// InvalidType is returned in case of mismatch.
func (m *ExpressionRHS) GetType(ts *Scope) Type {
	return m.expr.GetType(ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *Ident) TypeCheck(ts *Scope, errch chan<- error) {
	identT := m.GetType(ts)

	switch identT.(type) {
	case InvalidType:
		errch <- CreateUndelaredVariableError(
			m.Token(),
			m.ident,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *Ident) GetType(ts *Scope) Type {
	t := ts.Lookup(m.ident)
	if t == nil {
		return InvalidType{}
	}
	return t
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *IntLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *IntLiteral) GetType(ts *Scope) Type {
	return IntType{}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BoolLiteralFalse) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BoolLiteralFalse) GetType(ts *Scope) Type {
	return BoolType{}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BoolLiteralTrue) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BoolLiteralTrue) GetType(ts *Scope) Type {
	return BoolType{}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *CharLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *CharLiteral) GetType(ts *Scope) Type {
	return CharType{}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *StringLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *StringLiteral) GetType(ts *Scope) Type {
	return ArrayType{CharType{}}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *NullPair) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *NullPair) GetType(ts *Scope) Type {
	return PairType{first: UnknownType{}, second: UnknownType{}}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *ArrayElem) TypeCheck(ts *Scope, errch chan<- error) {
	array := ts.Lookup(m.ident)

	for _, index := range m.indexes {
		index.TypeCheck(ts, errch)

		switch indexT := index.GetType(ts).(type) {
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
		default:
			errch <- CreateTypeMismatchError(
				m.Token(),
				ArrayType{},
				arrayT,
			)
		}
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *ArrayElem) GetType(ts *Scope) Type {
	array := ts.Lookup(m.ident)
	for _, i := range m.indexes {
		switch i.GetType(ts).(type) {
		case IntType:
		default:
			return InvalidType{}
		}

		switch arrayT := array.(type) {
		case ArrayType:
			array = arrayT.base
		default:
			return InvalidType{}
		}
	}
	return array
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorNot) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case BoolType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			BoolType{},
			unopT,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *UnaryOperatorNot) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case BoolType:
		return BoolType{}
	default:
		return InvalidType{}
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorNegate) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case IntType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			IntType{},
			unopT,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *UnaryOperatorNegate) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case IntType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorLen) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case ArrayType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			ArrayType{},
			unopT,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *UnaryOperatorLen) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case ArrayType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorOrd) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case CharType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			CharType{},
			unopT,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *UnaryOperatorOrd) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case CharType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *UnaryOperatorChr) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case IntType:
	default:
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			IntType{},
			unopT,
		)
	}
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *UnaryOperatorChr) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case IntType:
		return CharType{}
	default:
		return InvalidType{}
	}
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorMult) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorMult) GetType(ts *Scope) Type {
	return getTypeBinaryArithmetic(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorDiv) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorDiv) GetType(ts *Scope) Type {
	return getTypeBinaryArithmetic(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorMod) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorMod) GetType(ts *Scope) Type {
	return getTypeBinaryArithmetic(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorAdd) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorAdd) GetType(ts *Scope) Type {
	return getTypeBinaryArithmetic(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorSub) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorSub) GetType(ts *Scope) Type {
	return getTypeBinaryArithmetic(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorGreaterThan) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorGreaterThan) GetType(ts *Scope) Type {
	return getTypeBinaryComparator(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorGreaterEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorGreaterEqual) GetType(ts *Scope) Type {
	return getTypeBinaryComparator(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorLessThan) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorLessThan) GetType(ts *Scope) Type {
	return getTypeBinaryComparator(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorLessEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckComparator(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorLessEqual) GetType(ts *Scope) Type {
	return getTypeBinaryComparator(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckEquality(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorEqual) GetType(ts *Scope) Type {
	return getTypeBinaryEquality(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorNotEqual) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckEquality(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorNotEqual) GetType(ts *Scope) Type {
	return getTypeBinaryEquality(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorAnd) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckBoolean(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorAnd) GetType(ts *Scope) Type {
	return getTypeBinaryBoolean(m, ts)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorOr) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckBoolean(m, ts, errch)
}

// GetType returns the deduced type of the expression.
// InvalidType is returned in case of error.
func (m *BinaryOperatorOr) GetType(ts *Scope) Type {
	return getTypeBinaryBoolean(m, ts)
}

func typeCheckArithmetic(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().GetType(ts)
	rhsT := m.GetRHS().GetType(ts)

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

func getTypeBinaryArithmetic(m BinaryOperator, ts *Scope) Type {
	lhs := m.GetLHS()
	rhs := m.GetRHS()

	if !(lhs.GetType(ts).Match(rhs.GetType(ts))) {
		return InvalidType{}
	}

	if (IntType{}.Match(lhs.GetType(ts))) {
		return IntType{}
	}

	return InvalidType{}
}

func typeCheckComparator(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().GetType(ts)
	rhsT := m.GetRHS().GetType(ts)

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

func getTypeBinaryComparator(m BinaryOperator, ts *Scope) Type {
	lhs := m.GetLHS()
	rhs := m.GetRHS()

	if !(lhs.GetType(ts).Match(rhs.GetType(ts))) {
		return InvalidType{}
	}

	switch lhs.GetType(ts).(type) {
	case IntType, CharType:
		return BoolType{}
	default:
		return InvalidType{}
	}
}

func typeCheckEquality(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().GetType(ts)
	rhsT := m.GetRHS().GetType(ts)

	if !(lhsT.Match(rhsT)) {
		errch <- CreateTypeMismatchError(
			m.GetRHS().Token(),
			lhsT,
			rhsT,
		)
	}
}

func getTypeBinaryEquality(m BinaryOperator, ts *Scope) Type {
	lhs := m.GetLHS()
	rhs := m.GetRHS()

	if !(lhs.GetType(ts).Match(rhs.GetType(ts))) {
		return InvalidType{}
	}

	return BoolType{}
}

func typeCheckBoolean(m BinaryOperator, ts *Scope, errch chan<- error) {
	m.GetLHS().TypeCheck(ts, errch)
	m.GetRHS().TypeCheck(ts, errch)

	lhsT := m.GetLHS().GetType(ts)
	rhsT := m.GetLHS().GetType(ts)

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

func getTypeBinaryBoolean(m BinaryOperator, ts *Scope) Type {
	lhs := m.GetLHS()
	rhs := m.GetRHS()

	if !(lhs.GetType(ts).Match(rhs.GetType(ts))) {
		return InvalidType{}
	}

	if (BoolType{}.Match(lhs.GetType(ts))) {
		return BoolType{}
	}

	return InvalidType{}
}

// TypeCheck on ExpLPar to satisfy interface. Never called.
func (m *ExprParen) TypeCheck(ts *Scope, errch chan<- error) {
}

// GetType on ExpLPar to satisfy interface. Never called.
func (m *ExprParen) GetType(ts *Scope) Type {
	return InvalidType{}
}
