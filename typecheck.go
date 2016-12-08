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
	classes    map[string]*ClassType
	members    map[string]Type
	funcs      map[string]map[string]map[string]*FunctionDef
	class      *ClassType
	returnType Type
}

// CreateRootScope creates a global scope that has no parent
func CreateRootScope() *Scope {
	scope := &Scope{
		parent:  nil,
		vars:    make(map[string]Type),
		classes: make(map[string]*ClassType),
		members: make(map[string]Type),
		funcs:   make(map[string]map[string]map[string]*FunctionDef),
	}

	return scope
}

// Child creates a child scope from a scope that inherits all properties but
// can declare variable independently
func (m *Scope) Child() *Scope {
	return &Scope{
		parent:     m,
		vars:       make(map[string]Type),
		classes:    m.classes,
		members:    m.members,
		funcs:      m.funcs,
		class:      m.class,
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

// LookupMember tries to search for the type of a given member
// It returns InvalidType if not found
func (m *Scope) LookupMember(ident string) Type {
	t, ok := m.members[ident]

	if !ok {
		return InvalidType{}
	}

	return t
}

// LookupFunction tries to return the function given it's identifier
// returns nil if not found.
func (m *Scope) LookupFunction(ident string) map[string]*FunctionDef {
	t, ok := m.funcs[""][ident]

	if !ok {
		return nil
	}

	return t
}

// LookupMethod tries to return the function given it's identifier and the name
// of the class it is on
// returns nil if not found.
func (m *Scope) LookupMethod(class, ident string) map[string]*FunctionDef {
	t, ok := m.funcs[class][ident]

	if !ok {
		return nil
	}

	return t
}

// LookupClass tries to return the class given it's identifier
// returns nil if not found.
func (m *Scope) LookupClass(ident string) *ClassType {
	t, ok := m.classes[ident]

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

// DeclareMember creates a new variable in the current scope returning the previous
// type in case of redeclaration, nil otherwise
func (m *Scope) DeclareMember(ident string, t Type) Type {
	pt, ok := m.members[ident]

	m.members[ident] = t

	if ok {
		return pt
	}
	return nil
}

// DeclareFunction registers a new function in the scope returning the previous
// one in case of redeclaration, nil otherwise
func (m *Scope) DeclareFunction(ident, symbol string, f *FunctionDef) *FunctionDef {
	if m.funcs[""] == nil {
		m.funcs[""] = make(map[string]map[string]*FunctionDef)
	}

	if m.funcs[""][ident] == nil {
		m.funcs[""][ident] = make(map[string]*FunctionDef)
	}

	pf, ok := m.funcs[""][ident][symbol]

	m.funcs[""][ident][symbol] = f

	if ok {
		return pf
	}

	return nil
}

// DeclareMethod registers a new function in the scope returning the previous
// one in case of redeclaration, nil otherwise
func (m *Scope) DeclareMethod(class, ident, symbol string, f *FunctionDef) *FunctionDef {
	if m.funcs[class] == nil {
		m.funcs[class] = make(map[string]map[string]*FunctionDef)
	}

	if m.funcs[class][ident] == nil {
		m.funcs[class][ident] = make(map[string]*FunctionDef)
	}

	pf, ok := m.funcs[class][ident][symbol]

	m.funcs[class][ident][symbol] = f

	if ok {
		return pf
	}

	return nil
}

// DeclareClass registers a new class in the scope returning the previous
// one in case of redeclaration, nil otherwise
func (m *Scope) DeclareClass(ident string, c *ClassType) *ClassType {
	if m.classes == nil {
		m.classes = make(map[string]*ClassType)
	}

	pf, ok := m.classes[ident]

	m.classes[ident] = c

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
func (m VoidType) Match(t Type) bool {
	return true
}

// Match checks whether a type is assignable to the current type
func (m IntType) Match(t Type) bool {
	switch t.(type) {
	case IntType:
		return true
	case VoidType:
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
	case VoidType:
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
	case VoidType:
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
	case VoidType:
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
	case VoidType:
		return true
	default:
		return false
	}
}

// Match checks whether a type is assignable to the current type
func (m *ClassType) Match(t Type) bool {
	switch o := t.(type) {
	case *ClassType:
		return m.name == o.name
	case VoidType:
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

		// add the classes and methods to the scope
		for _, c := range m.classes {
			if pc := global.DeclareClass(c.name, c); pc != nil {
				errch <- CreateClassRedeclarationError(
					c.Token(),
					c.name,
				)
			}
			for _, m := range c.methods {
				m.class = c
				if pm := global.DeclareMethod(c.name, m.ident, m.Symbol(), m); pm != nil {
					errch <- CreateFunctionRedelarationError(
						m.Token(),
						m.ident,
					)
				}
			}
		}

		// add the functions to the scope
		for _, f := range m.functions {
			if pf := global.DeclareFunction(f.ident, f.Symbol(), f); pf != nil {
				errch <- CreateFunctionRedelarationError(
					f.Token(),
					f.ident,
				)
			}
		}

		// check class methods
		for _, c := range m.classes {
			cs := global.Child()
			cs.class = c
			cs.members = make(map[string]Type)
			// add the members
			for _, m := range c.members {
				switch m.wtype.(type) {
				case VoidType:
					errch <- CreateInvalidVoidTypeError(
						m.Token(),
						m.ident,
					)
				}
				if pt := cs.DeclareMember(m.ident, m.wtype); pt != nil {
					errch <- CreateVariableRedeclarationError(
						m.Token(),
						m.ident,
						pt,
						m.wtype,
					)
				}
			}
			// typecheck methods
			for _, m := range c.methods {
				mscope := cs.Child()
				for _, arg := range m.params {
					switch arg.wtype.(type) {
					case VoidType:
						errch <- CreateInvalidVoidTypeError(
							arg.Token(),
							arg.name,
						)
					}
					pt := mscope.Declare(arg.name, arg.wtype)
					if pt != nil {
						errch <- CreateVariableRedeclarationError(
							arg.Token(),
							arg.name,
							pt,
							arg.wtype,
						)
					}
				}
				mscope.returnType = m.returnType
				m.body.TypeCheck(mscope, errch)
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
				switch arg.wtype.(type) {
				case VoidType:
					errch <- CreateInvalidVoidTypeError(
						arg.Token(),
						arg.name,
					)
				}
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
	switch m.wtype.(type) {
	case VoidType:
		errch <- CreateInvalidVoidTypeError(
			m.Token(),
			m.ident,
		)
	}

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
	case *ClassType:
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
		errch <- CreateTypeMismatchError(
			m.expr.Token(),
			&ClassType{name: "any class"},
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

	switch returnT.(type) {
	case VoidType:
		switch exprT.(type) {
		case VoidType:
		default:
			errch <- CreateTypeMismatchError(
				m.expr.Token(),
				returnT,
				exprT,
			)
		}
	}

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

// TypeCheck checks whether the call is valid
// The check is propagated recursively.
func (m *FunctionCallStat) TypeCheck(ts *Scope, errch chan<- error) {
	for _, arg := range m.args {
		arg.TypeCheck(ts, errch)
	}

	var classname string
	if len(m.obj) > 0 {
		var recvT Type
		switch m.obj[0] {
		case '@':
			recvT = ts.LookupMember(m.obj[1:])
		default:
			recvT = ts.Lookup(m.obj)
		}

		switch t := recvT.(type) {
		case *ClassType:
			classname = t.name
		default:
			errch <- CreateFunctionCallOnNonObjectError(
				m.Token(),
				m.ident,
				t,
			)
		}
	}

	var overloads map[string]*FunctionDef
	if len(classname) > 0 {
		overloads = ts.LookupMethod(classname, m.ident)
	} else {
		overloads = ts.LookupFunction(m.ident)
	}

	if overloads == nil {
		errch <- CreateCallingNonFunctionError(
			m.Token(),
			m.ident,
		)
	}

	found := false
	mangledIdent := ""

	m.wtype = InvalidType{}

	for symbol, fun := range overloads {
		if len(fun.params) != len(m.args) {
			continue
		}

		match := true

		for i := 0; i < len(fun.params) && i < len(m.args) && match; i++ {
			paramT := fun.params[i].wtype
			argT := m.args[i].Type()
			if !argT.Match(paramT) {
				match = false
			}
		}

		if match && found {
			errch <- CreateAmbigousFunctionCallError(
				m.Token(),
				m.ident,
			)
		}

		if match {
			found = true
			mangledIdent = symbol
			m.wtype = fun.returnType
		}
	}

	if !found {
		errch <- CreateNoSuchOverloadError(m.Token(), m.ident)
	}

	m.mangledIdent = mangledIdent

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
	if m.falseStat != nil {
		m.falseStat.TypeCheck(ts.Child(), errch)
	}

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

// TypeCheck checks whether the statement has any type mismatches in expressions
// and assignments. The check is propagated recursively
func (m *SwitchStatement) TypeCheck(ts *Scope, errch chan<- error) {

	m.cond.TypeCheck(ts, errch)
	condT := m.cond.Type()

	if !(BoolType{}.Match(condT) || IntType{}.Match(condT) || CharType{}.Match(condT)) {
		errch <- CreateTypeMismatchError(
			m.cond.Token(),
			IntType{},
			condT,
		)
	}

	for index := 0; index < len(m.cases); index++ {
		m.cases[index].TypeCheck(ts.Child(), errch)

		exprT := m.cases[index].Type()

		if !condT.Match(exprT) {
			errch <- CreateTypeMismatchError(
				m.cond.Token(),
				condT,
				exprT,
			)
		}

		m.bodies[index].TypeCheck(ts.Child(), errch)
	}

	if m.defaultCase != nil {
		m.defaultCase.TypeCheck(ts.Child(), errch)
	}
}

// TypeCheck checks whether the left hand is a valid assignment target.
// The check propagated recursively.
func (m *DoWhileStatement) TypeCheck(ts *Scope, errch chan<- error) {
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

func (m *ForStatement) TypeCheck(ts *Scope, errch chan<- error) {
	child := ts.Child()

	switch t := m.init.(type) {
	case *DeclareAssignStatement:
		t.TypeCheck(child, errch)
	default:
		errch <- CreateForLoopError(m.init.Token(),
			"Declare Assign statement expected")
	}

	m.cond.TypeCheck(child, errch)
	boolT := m.cond.Type()

	if !(BoolType{}.Match(boolT)) {
		errch <- CreateTypeMismatchError(
			m.cond.Token(),
			BoolType{},
			boolT,
		)
	}

	switch t := m.after.(type) {
	case *AssignStatement:
		t.TypeCheck(child, errch)
	default:
		errch <- CreateForLoopError(m.after.Token(),
			"Assign statement expected")
	}

	m.body.TypeCheck(child, errch)

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
	var t Type
	switch m.ident[0] {
	case '@':
		t = ts.LookupMember(m.ident[1:])
	default:
		t = ts.Lookup(m.ident)
	}

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
	var t Type
	switch m.ident[0] {
	case '@':
		t = ts.LookupMember(m.ident[1:])
	default:
		t = ts.Lookup(m.ident)
	}

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
	for _, arg := range m.args {
		arg.TypeCheck(ts, errch)
	}

	var classname string
	if len(m.obj) > 0 {
		var recvT Type
		switch m.obj[0] {
		case '@':
			recvT = ts.LookupMember(m.obj[1:])
		default:
			recvT = ts.Lookup(m.obj)
		}

		switch t := recvT.(type) {
		case *ClassType:
			classname = t.name
		default:
			errch <- CreateFunctionCallOnNonObjectError(
				m.Token(),
				m.ident,
				t,
			)
		}
	}

	var overloads map[string]*FunctionDef
	if len(classname) > 0 {
		overloads = ts.LookupMethod(classname, m.ident)
	} else {
		overloads = ts.LookupFunction(m.ident)
	}

	if overloads == nil {
		errch <- CreateCallingNonFunctionError(
			m.Token(),
			m.ident,
		)
	}

	found := false
	mangledIdent := ""

	m.wtype = InvalidType{}

	for symbol, fun := range overloads {
		if len(fun.params) != len(m.args) {
			continue
		}

		match := true

		for i := 0; i < len(fun.params) && i < len(m.args) && match; i++ {
			paramT := fun.params[i].wtype
			argT := m.args[i].Type()
			if !argT.Match(paramT) {
				match = false
			}
		}

		if match && found {
			errch <- CreateAmbigousFunctionCallError(
				m.Token(),
				m.ident,
			)
		}

		if match {
			found = true
			mangledIdent = symbol
			m.wtype = fun.returnType
		}
	}

	if !found {
		errch <- CreateNoSuchOverloadError(m.Token(), m.ident)
	}

	m.mangledIdent = mangledIdent

	switch m.wtype.(type) {
	case VoidType:
		errch <- CreateVoidAssignmentError(
			m.Token(),
			m.ident,
		)
	}
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *ExpressionRHS) TypeCheck(ts *Scope, errch chan<- error) {
	m.expr.TypeCheck(ts, errch)
}

// TypeCheck checks whether the right hand side is valid and assignable
// The check is propagated recursively.
func (m *NewInstanceRHS) TypeCheck(ts *Scope, errch chan<- error) {
	var ct *ClassType
	switch t := m.wtype.(type) {
	case *ClassType:
		ct = t
	default:
		errch <- CreateTypeMismatchError(
			m.Token(),
			&ClassType{name: "any class"},
			t,
		)
	}

	c := ts.LookupClass(ct.name)

	if c == nil {
		errch <- CreateUndeclaredClassError(
			m.Token(),
			ct.name,
		)
		m.wtype = InvalidType{}
	} else {
		m.wtype = c
	}

	for _, arg := range m.args {
		arg.TypeCheck(ts, errch)
	}

	var classname = ct.name

	var overloads map[string]*FunctionDef
	if len(classname) > 0 {
		overloads = ts.LookupMethod(classname, "init")
	}

	if overloads == nil {
		errch <- CreateCallingNonFunctionError(
			m.Token(),
			"init",
		)
	}

	found := false
	constr := ""

	for symbol, fun := range overloads {
		if len(fun.params) != len(m.args) {
			continue
		}

		match := true

		for i := 0; i < len(fun.params) && i < len(m.args) && match; i++ {
			paramT := fun.params[i].wtype
			argT := m.args[i].Type()
			if !argT.Match(paramT) {
				match = false
			}
		}

		if match && found {
			errch <- CreateAmbigousFunctionCallError(
				m.Token(),
				"ident",
			)
		}

		if match {
			found = true
			constr = symbol
		}
	}

	if !found {
		errch <- CreateNoSuchOverloadError(m.Token(), "init")
	}

	m.constr = constr
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *Ident) TypeCheck(ts *Scope, errch chan<- error) {
	var t Type
	switch m.ident[0] {
	case '@':
		t = ts.LookupMember(m.ident[1:])
	default:
		t = ts.Lookup(m.ident)
	}

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

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorBitAnd) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
}

// TypeCheck checks expression whether all operators get the type they can
// operate on, all variables are declared, arrays are indexed properly.
// The check is propagated recursively.
func (m *BinaryOperatorBitOr) TypeCheck(ts *Scope, errch chan<- error) {
	typeCheckArithmetic(m, ts, errch)
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

// TypeCheck on VoidExpr
func (m *VoidExpr) TypeCheck(ts *Scope, errch chan<- error) {
}

// TypeCheck on ExpLPar to satisfy interface. Never called.
func (m *ExprParen) TypeCheck(ts *Scope, errch chan<- error) {
}
