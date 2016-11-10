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
	switch t.(type) {
	case IntType:
		return true
	default:
		return false
	}
}

func (m BoolType) Match(t Type) bool {
	switch t.(type) {
	case BoolType:
		return true
	default:
		return false
	}
}

func (m CharType) Match(t Type) bool {
	switch t.(type) {
	case CharType:
		return true
	default:
		return false
	}
}

func (m PairType) Match(t Type) bool {
	switch o := t.(type) {
	case PairType:
		fst := m.first == nil ||
			o.first == nil ||
			m.first.Match(o.first)
		snd := m.second == nil ||
			o.second == nil ||
			m.second.Match(o.second)
		return fst && snd
	default:
		return false
	}
}

func (m ArrayType) Match(t Type) bool {
	switch o := t.(type) {
	case ArrayType:
		return m.base == nil ||
			o.base == nil ||
			m.base.Match(o.base)
	default:
		return false
	}
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
		errch <- &TypeMismatch{
			expected: m.waccType,
			got:      rhsT,
		}
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *AssignStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.target.TypeCheck(ts, errch)
	m.rhs.TypeCheck(ts, errch)

	lhsT := m.target.GetType(ts)
	rhsT := m.rhs.GetType(ts)

	if !lhsT.Match(rhsT) {
		errch <- &TypeMismatch{
			expected: lhsT,
			got:      rhsT,
		}
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ReadStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.target.TypeCheck(ts, errch)

	switch t := m.target.GetType(ts).(type) {
	case IntType:
	case CharType:
	default:
		errch <- &TypeMismatch{
			expected: IntType{},
			got:      t,
		}
		errch <- &TypeMismatch{
			expected: CharType{},
			got:      t,
		}
	}
}

func (m *FreeStatement) TypeCheck(ts *Scope, errch chan<- error) {

	m.expr.TypeCheck(ts, errch)
	freeT := m.expr.GetType(ts)

	switch t := freeT.(type) {
	case PairType:
	case ArrayType:
	default:
		errch <- &TypeMismatch{
			expected: PairType{},
			got:      t,
		}
		errch <- &TypeMismatch{
			expected: ArrayType{},
			got:      t,
		}
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ReturnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *ExitStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	exitT := m.expr.GetType(ts)

	if !(IntType{}.Match(exitT)) {
		errch <- &TypeMismatch{
			expected: IntType{},
			got:      exitT,
		}
	}

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PrintLnStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PrintStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *IfStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.cond.TypeCheck(ts, errch)
	boolT := m.cond.GetType(ts)

	if !(BoolType{}.Match(boolT)) {
		errch <- &TypeMismatch{
			expected: BoolType{},
			got:      boolT,
		}
	}

	m.trueStat.TypeCheck(ts, errch)
	m.falseStat.TypeCheck(ts, errch)

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *WhileStatement) TypeCheck(ts *Scope, errch chan<- error) {
	m.cond.TypeCheck(ts, errch)
	boolT := m.cond.GetType(ts)

	if !(BoolType{}.Match(boolT)) {
		errch <- &TypeMismatch{
			expected: BoolType{},
			got:      boolT,
		}
	}

	m.body.TypeCheck(ts, errch)

	m.BaseStatement.TypeCheck(ts, errch)
}

func (m *PairElemLHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch t := m.expr.GetType(ts).(type) {
	case PairType:
	default:
		errch <- &TypeMismatch{
			expected: PairType{},
			got:      t,
		}
	}
}

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

func (m *ArrayLHS) TypeCheck(ts *Scope, errch chan<- error) {
	t := ts.Lookup(m.ident)

	for _, i := range m.index {
		i.TypeCheck(ts, errch)

		if !(IntType{}).Match(i.GetType(ts)) {
			errch <- &TypeMismatch{
				expected: IntType{},
				got:      t,
			}
		}

		switch arr := t.(type) {
		case ArrayType:
			t = arr.base
		default:
			errch <- &TypeMismatch{
				expected: ArrayType{},
				got:      t,
			}
		}
	}
}

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

func (m *VarLHS) TypeCheck(ts *Scope, errch chan<- error) {
	switch ts.Lookup(m.ident).(type) {
	case InvalidType:
		errch <- &UndeclaredVariable{
			ident: m.ident,
		}
	}
}

func (m *VarLHS) GetType(ts *Scope) Type {
	return ts.Lookup(m.ident)
}

func (m *PairLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.fst.TypeCheck(ts, errch)
	m.snd.TypeCheck(ts, errch)
}

func (m *PairLiterRHS) GetType(ts *Scope) Type {
	fstT := m.fst.GetType(ts)
	sndT := m.snd.GetType(ts)

	if (InvalidType{}.Match(fstT) || InvalidType{}.Match(sndT)) {
		return InvalidType{}
	}

	return PairType{}
}

func (m *PairElemRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
	pairT := m.expr.GetType(ts)

	if !(PairType{}.Match(pairT)) {
		errch <- &TypeMismatch{
			expected: PairType{},
			got:      pairT,
		}
	}
}

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

func (m *ArrayLiterRHS) TypeCheck(ts *Scope, errch chan<- error) {
	if len(m.elements) == 0 {
		return
	}

	t := m.elements[0].GetType(ts)

	for _, elem := range m.elements {
		elem.TypeCheck(ts, errch)

		if !t.Match(elem.GetType(ts)) {
			errch <- &TypeMismatch{
				expected: t,
				got:      elem.GetType(ts),
			}
		}
	}
}

func (m *ArrayLiterRHS) GetType(ts *Scope) Type {
	if len(m.elements) == 0 {
		return ArrayType{}
	}

	t := m.elements[0].GetType(ts)

	for _, elem := range m.elements {
		if !t.Match(elem.GetType(ts)) {
			return InvalidType{}
		}
	}

	return t
}

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
			errch <- &TypeMismatch{
				expected: paramT,
				got:      argT,
			}
		}
	}
}

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

func (m *ExpressionRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
}

func (m *ExpressionRHS) GetType(ts *Scope) Type {
	return m.expr.GetType(ts)
}

func (m *Ident) TypeCheck(ts *Scope, errch chan<- error) {
	identT := m.GetType(ts)

	if (InvalidType{}.Match(identT)) {
		errch <- &UndeclaredVariable{
			ident: m.ident,
		}
	}
}

func (m *Ident) GetType(ts *Scope) Type {
	t := ts.Lookup(m.ident)
	if t == nil {
		return InvalidType{}
	} else {
		return t
	}
}

func (m *IntLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *IntLiteral) GetType(ts *Scope) Type {
	return IntType{}
}

func (m *BoolLiteralFalse) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BoolLiteralFalse) GetType(ts *Scope) Type {
	return BoolType{}
}

func (m *BoolLiteralTrue) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BoolLiteralTrue) GetType(ts *Scope) Type {
	return BoolType{}
}

func (m *CharLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *CharLiteral) GetType(ts *Scope) Type {
	return CharType{}
}

func (m *StringLiteral) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *StringLiteral) GetType(ts *Scope) Type {
	return ArrayType{CharType{}}
}

func (m *NullPair) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *NullPair) GetType(ts *Scope) Type {
	return PairType{}
}

func (m *ArrayElem) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *ArrayElem) GetType(ts *Scope) Type {
	array := ts.Lookup(m.ident)
	for i := 0; i < len(m.indexes); i++ {
		switch arrayT := array.(type) {
		case ArrayType:
			array = arrayT.base
		default:
			return InvalidType{}
		}
	}
	return array
}

func (m *UnaryOperatorNot) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case BoolType:
	default:
		errch <- &TypeMismatch{
			expected: BoolType{},
			got:      unopT,
		}
	}
}

func (m *UnaryOperatorNot) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case BoolType:
		return BoolType{}
	default:
		return InvalidType{}
	}
}

func (m *UnaryOperatorNegate) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case IntType:
	default:
		errch <- &TypeMismatch{
			expected: IntType{},
			got:      unopT,
		}
	}
}

func (m *UnaryOperatorNegate) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case IntType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

func (m *UnaryOperatorLen) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case ArrayType:
	default:
		errch <- &TypeMismatch{
			expected: ArrayType{},
			got:      unopT,
		}
	}
}

func (m *UnaryOperatorLen) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case ArrayType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

func (m *UnaryOperatorOrd) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case CharType:
	default:
		errch <- &TypeMismatch{
			expected: CharType{},
			got:      unopT,
		}
	}
}

func (m *UnaryOperatorOrd) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case CharType:
		return IntType{}
	default:
		return InvalidType{}
	}
}

func (m *UnaryOperatorChr) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)

	switch unopT := m.expr.GetType(ts).(type) {
	case IntType:
	default:
		errch <- &TypeMismatch{
			expected: IntType{},
			got:      unopT,
		}
	}
}

func (m *UnaryOperatorChr) GetType(ts *Scope) Type {
	switch m.GetExpression().GetType(ts).(type) {
	case IntType:
		return CharType{}
	default:
		return InvalidType{}
	}
}

func (m *BinaryOperatorMult) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorMult) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorDiv) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorDiv) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorMod) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorMod) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorAdd) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorAdd) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorSub) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorSub) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorGreaterThan) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorGreaterThan) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorGreaterEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorGreaterEqual) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorLessThan) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorLessThan) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorLessEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorLessEqual) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorEqual) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorNotEqual) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorNotEqual) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorAnd) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorAnd) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func (m *BinaryOperatorOr) TypeCheck(ts *Scope, errch chan<- error) {
}

func (m *BinaryOperatorOr) GetType(ts *Scope) Type {
	return GetTypeBinary(m, ts)
}

func GetTypeBinary(m BinaryOperator, ts *Scope) Type {
	lhs := m.GetLHS()
	rhs := m.GetRHS()
	if lhs != rhs {
		return InvalidType{}
	} else {
		switch m.(type) {
		case *BinaryOperatorMult,
			*BinaryOperatorDiv,
			*BinaryOperatorMod,
			*BinaryOperatorAdd,
			*BinaryOperatorSub:
			switch lhs.GetType(ts).(type) {
			case IntType:
				return IntType{}
			default:
				return InvalidType{}
			}
		case *BinaryOperatorGreaterThan,
			*BinaryOperatorGreaterEqual,
			*BinaryOperatorLessThan,
			*BinaryOperatorLessEqual:
			switch lhs.GetType(ts).(type) {
			case IntType,
				CharType:
				return BoolType{}
			default:
				return InvalidType{}
			}
		case *BinaryOperatorEqual,
			*BinaryOperatorNotEqual:
			return BoolType{}
		case *BinaryOperatorAnd,
			*BinaryOperatorOr:
			switch lhs.GetType(ts).(type) {
			case BoolType:
				return BoolType{}
			default:
				return InvalidType{}
			}
		default:
			return InvalidType{}
		}
	}
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
