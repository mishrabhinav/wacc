package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Type interface{}

type IntType struct{}

type BoolType struct{}

type CharType struct{}

type PairType struct {
	first  Type
	second Type
}

type ArrayType struct {
	base Type
}

type Expression interface{}

type Statement interface {
	GetNext() Statement
	SetNext(Statement)
}

type BaseStatement struct {
	next Statement
}

func (m *BaseStatement) GetNext() Statement {
	return m.next
}

func (m *BaseStatement) SetNext(next Statement) {
	m.next = next
}

type SkipStatement struct {
	BaseStatement
}

type BlockStatement struct {
	BaseStatement
	body Statement
}

type DeclareAssignStatement struct {
	BaseStatement
	waccType Type
	ident    string
	rhs      RHS
}

type LHS interface{}

type PairElemLHS struct {
	snd  bool
	expr Expression
}

type ArrayLHS struct {
	ident string
	index []Expression
}

type VarLHS struct {
	ident string
}

type RHS interface{}

type PairLiterRHS struct {
	PairLiteral
}

type ArrayLiterRHS struct {
	elements []Expression
}

type PairElemRHS struct {
	snd  bool
	expr Expression
}

type FunctionCallRHS struct {
	ident string
	args  []Expression
}

type ExpressionRHS struct {
	expr Expression
}

type AssignStatement struct {
	BaseStatement
	target LHS
	rhs    RHS
}

type ReadStatement struct {
	BaseStatement
	target LHS
}

type FreeStatement struct {
	BaseStatement
	expr Expression
}

type ReturnStatement struct {
	BaseStatement
	expr Expression
}

type ExitStatement struct {
	BaseStatement
	expr Expression
}

type PrintLnStatement struct {
	BaseStatement
	expr Expression
}

type PrintStatement struct {
	BaseStatement
	expr Expression
}

type IfStatement struct {
	BaseStatement
	cond      Expression
	trueStat  Statement
	falseStat Statement
}

type WhileStatement struct {
	BaseStatement
	cond Expression
	body Statement
}

type FunctionParam struct {
	name     string
	waccType Type
}

type FunctionDef struct {
	ident      string
	returnType Type
	params     []*FunctionParam
	body       Statement
}

type AST struct {
	main      Statement
	functions []*FunctionDef
}

func nodeRange(node *node32) <-chan *node32 {
	out := make(chan *node32)
	go func() {
		for ; node != nil; node = node.next {
			out <- node
		}
		close(out)
	}()
	return out
}

func nextNode(node *node32, rule pegRule) *node32 {
	for cnode := range nodeRange(node) {
		if cnode.pegRule == rule {
			return cnode
		}
	}

	return nil
}

func parseArrayElem(node *node32) (Expression, error) {
	arrElem := &ArrayElem{}

	arrElem.ident = node.match

	for enode := nextNode(node, ruleEXPR); enode != nil; enode = nextNode(enode.next, ruleEXPR) {
		var exp Expression
		var err error
		if exp, err = parseExpr(enode.up); err != nil {
			return nil, err
		}
		arrElem.indexes = append(arrElem.indexes, exp)
	}

	return arrElem, nil
}

type Ident struct {
	ident string
}

type IntLiteral struct {
	value int
}

type BoolLiteralTrue struct{}

type BoolLiteralFalse struct{}

type CharLiteral struct {
	char string
}

type StringLiteral struct {
	str string
}

type PairLiteral struct {
	fst Expression
	snd Expression
}

type NullPair struct{}

type ArrayElem struct {
	ident   string
	indexes []Expression
}

type UnaryOperator interface {
	GetExpression() Expression
	SetExpression(Expression)
}

type UnaryOperatorBase struct {
	expr Expression
}

func (m *UnaryOperatorBase) GetExpression() Expression {
	return m.expr
}

func (m *UnaryOperatorBase) SetExpression(exp Expression) {
	m.expr = exp
}

type UnaryOperatorNot struct {
	UnaryOperatorBase
}

type UnaryOperatorNegate struct {
	UnaryOperatorBase
}

type UnaryOperatorLen struct {
	UnaryOperatorBase
}

type UnaryOperatorOrd struct {
	UnaryOperatorBase
}

type UnaryOperatorChr struct {
	UnaryOperatorBase
}

type BinaryOperator interface {
	GetRHS() Expression
	SetRHS(Expression)
	GetLHS() Expression
	SetLHS(Expression)
}

type BinaryOperatorBase struct {
	lhs Expression
	rhs Expression
}

func (m *BinaryOperatorBase) GetLHS() Expression {
	return m.lhs
}

func (m *BinaryOperatorBase) SetLHS(exp Expression) {
	m.lhs = exp
}

func (m *BinaryOperatorBase) GetRHS() Expression {
	return m.rhs
}

func (m *BinaryOperatorBase) SetRHS(exp Expression) {
	m.rhs = exp
}

type BinaryOperatorMult struct {
	BinaryOperatorBase
}

type BinaryOperatorDiv struct {
	BinaryOperatorBase
}

type BinaryOperatorMod struct {
	BinaryOperatorBase
}

type BinaryOperatorAdd struct {
	BinaryOperatorBase
}

type BinaryOperatorSub struct {
	BinaryOperatorBase
}

type BinaryOperatorGreaterThan struct {
	BinaryOperatorBase
}

type BinaryOperatorGreaterEqual struct {
	BinaryOperatorBase
}

type BinaryOperatorLessThan struct {
	BinaryOperatorBase
}

type BinaryOperatorLessEqual struct {
	BinaryOperatorBase
}

type BinaryOperatorEqual struct {
	BinaryOperatorBase
}

type BinaryOperatorNotEqual struct {
	BinaryOperatorBase
}

type BinaryOperatorAnd struct {
	BinaryOperatorBase
}

type BinaryOperatorOr struct {
	BinaryOperatorBase
}

type ExprLPar struct{}

type ExprRPar struct{}

func exprStream(node *node32) <-chan *node32 {
	out := make(chan *node32)
	go func() {
		for ; node != nil; node = node.next {
			switch node.pegRule {
			case ruleSPACE:
			case ruleBOOLLITER:
				out <- node.up
			case ruleEXPR:
				for inode := range exprStream(node.up) {
					out <- inode
				}
			default:
				out <- node
			}
		}
		close(out)
	}()
	return out
}

func parseExpr(node *node32) (Expression, error) {
	var stack []Expression
	var opstack []Expression

	push := func(e Expression) {
		stack = append(stack, e)
	}

	pop := func() (ret Expression) {
		ret, stack = stack[len(stack)-1], stack[:len(stack)-1]
		return
	}

	pushop := func(e Expression) {
		opstack = append(opstack, e)
	}

	peekop := func() Expression {
		if len(opstack) == 0 {
			return nil
		}
		return opstack[len(opstack)-1]
	}

	popop := func() {
		var exp Expression

		exp, opstack = opstack[len(opstack)-1], opstack[:len(opstack)-1]

		switch t := exp.(type) {
		case UnaryOperator:
			t.SetExpression(pop())
		case BinaryOperator:
			t.SetRHS(pop())
			t.SetLHS(pop())
		case *ExprLPar, ExprRPar:
			exp = nil
		}

		if exp != nil {
			push(exp)
		}
	}

	prio := func(exp Expression) int {
		switch exp.(type) {
		case *UnaryOperatorNot:
			return 2
		case *UnaryOperatorNegate:
			return 2
		case *UnaryOperatorLen:
			return 2
		case *UnaryOperatorOrd:
			return 2
		case *UnaryOperatorChr:
			return 2
		case *BinaryOperatorMult:
			return 3
		case *BinaryOperatorDiv:
			return 3
		case *BinaryOperatorMod:
			return 3
		case *BinaryOperatorAdd:
			return 4
		case *BinaryOperatorSub:
			return 4
		case *BinaryOperatorGreaterThan:
			return 6
		case *BinaryOperatorGreaterEqual:
			return 6
		case *BinaryOperatorLessThan:
			return 6
		case *BinaryOperatorLessEqual:
			return 6
		case *BinaryOperatorEqual:
			return 7
		case *BinaryOperatorNotEqual:
			return 7
		case *BinaryOperatorAnd:
			return 11
		case *BinaryOperatorOr:
			return 12
		case *ExprLPar:
			return 13
		default:
			return 42
		}
	}

	rightAssoc := func(exp Expression) bool {
		switch exp.(type) {
		case *UnaryOperatorNot:
			return true
		case *UnaryOperatorNegate:
			return true
		case *UnaryOperatorLen:
			return true
		case *UnaryOperatorOrd:
			return true
		case *UnaryOperatorChr:
			return true
		default:
			return false
		}
	}

	ruleToOp := func(outer, inner pegRule) Expression {
		switch outer {
		case ruleUNARYOPER:
			switch inner {
			case ruleBANG:
				return &UnaryOperatorNot{}
			case ruleMINUS:
				return &UnaryOperatorNegate{}
			case ruleLEN:
				return &UnaryOperatorLen{}
			case ruleORD:
				return &UnaryOperatorOrd{}
			case ruleCHR:
				return &UnaryOperatorChr{}
			}
		case ruleBINARYOPER:
			switch inner {
			case ruleSTAR:
				return &BinaryOperatorMult{}
			case ruleDIV:
				return &BinaryOperatorDiv{}
			case ruleMOD:
				return &BinaryOperatorMod{}
			case rulePLUS:
				return &BinaryOperatorAdd{}
			case ruleMINUS:
				return &BinaryOperatorSub{}
			case ruleGT:
				return &BinaryOperatorGreaterThan{}
			case ruleGE:
				return &BinaryOperatorGreaterEqual{}
			case ruleLT:
				return &BinaryOperatorLessThan{}
			case ruleLE:
				return &BinaryOperatorLessEqual{}
			case ruleEQUEQU:
				return &BinaryOperatorEqual{}
			case ruleBANGEQU:
				return &BinaryOperatorNotEqual{}
			case ruleANDAND:
				return &BinaryOperatorAnd{}
			case ruleOROR:
				return &BinaryOperatorOr{}
			}
		}

		return nil
	}

	for enode := range exprStream(node) {
		switch enode.pegRule {
		case ruleINTLITER:
			num, err := strconv.ParseInt(enode.match, 10, 32)
			if err != nil {
				return nil, err
			}
			push(&IntLiteral{int(num)})
		case ruleFALSE:
			push(&BoolLiteralFalse{})
		case ruleTRUE:
			push(&BoolLiteralTrue{})
		case ruleCHARLITER:
			push(&CharLiteral{enode.up.next.match})
		case ruleSTRLITER:
			push(&StringLiteral{enode.up.next.match})
		case rulePAIRLITER:
			push(&NullPair{})
		case ruleIDENT:
			push(&Ident{enode.match})
		case ruleARRAYELEM:
			arrElem, err := parseArrayElem(enode.up)
			if err != nil {
				return nil, err
			}
			push(arrElem)
		case ruleUNARYOPER, ruleBINARYOPER:
			op1 := ruleToOp(enode.pegRule, enode.up.pegRule)
		op2l:
			for op2 := peekop(); op2 != nil; op2 = peekop() {
				if op2 == nil {
					break
				}

				switch {
				case !rightAssoc(op1) && prio(op1) >= prio(op2),
					rightAssoc(op1) && prio(op1) > prio(op2):
					popop()
				default:
					break op2l
				}
			}
			pushop(op1)
		case ruleLPAR:
			pushop(&ExprLPar{})
		case ruleRPAR:
		parloop:
			for {
				switch peekop().(type) {
				case *ExprLPar:
					popop()
					break parloop
				default:
					popop()
				}
			}
		}
	}

	for peekop() != nil {
		popop()
	}

	return pop(), nil
}

func parseLHS(node *node32) (LHS, error) {
	switch node.pegRule {
	case rulePAIRELEM:
		target := new(PairElemLHS)

		fstNode := nextNode(node.up, ruleFST)
		target.snd = fstNode == nil

		exprNode := nextNode(node.up, ruleEXPR)
		var err error
		if target.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		return target, nil
	case ruleARRAYELEM:
		target := new(ArrayLHS)

		identNode := nextNode(node.up, ruleIDENT)
		target.ident = identNode.match

		for exprNode := nextNode(node.up, ruleEXPR); exprNode != nil; exprNode = nextNode(exprNode.next, ruleEXPR) {
			var expr Expression
			var err error
			if expr, err = parseExpr(exprNode.up); err != nil {
				return nil, err
			}
			target.index = append(target.index, expr)
		}

		return target, nil
	case ruleIDENT:
		return &VarLHS{ident: node.match}, nil
	default:
		return nil, fmt.Errorf("Unexpected %s %s", node.String(), node.match)
	}
}

func parseRHS(node *node32) (RHS, error) {
	switch node.pegRule {
	case ruleNEWPAIR:
		var err error
		pair := new(PairLiterRHS)

		fstNode := nextNode(node, ruleEXPR)
		if pair.fst, err = parseExpr(fstNode.up); err != nil {
			return nil, err
		}

		sndNode := nextNode(fstNode.next, ruleEXPR)
		if pair.snd, err = parseExpr(sndNode.up); err != nil {
			return nil, err
		}

		return pair, nil
	case ruleARRAYLITER:
		node = node.up

		arr := new(ArrayLiterRHS)

		for node = nextNode(node, ruleEXPR); node != nil; node = nextNode(node.next, ruleEXPR) {
			var err error
			var expr Expression

			if expr, err = parseExpr(node.up); err != nil {
				return nil, err
			}
			arr.elements = append(arr.elements, expr)
		}

		return arr, nil
	case rulePAIRELEM:
		target := new(PairElemRHS)

		fstNode := nextNode(node.up, ruleFST)
		target.snd = fstNode == nil

		exprNode := nextNode(node.up, ruleEXPR)
		var err error
		if target.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		return target, nil
	case ruleCALL:
		call := new(FunctionCallRHS)

		identNode := nextNode(node, ruleIDENT)
		call.ident = identNode.match

		arglistNode := nextNode(node, ruleARGLIST)
		if arglistNode == nil {
			return call, nil
		}

		for argNode := nextNode(arglistNode.up, ruleEXPR); argNode != nil; argNode = nextNode(argNode.next, ruleEXPR) {
			var err error
			var expr Expression

			if expr, err = parseExpr(argNode.up); err != nil {
				return nil, err
			}

			call.args = append(call.args, expr)
		}

		return call, nil
	case ruleEXPR:
		exprRHS := new(ExpressionRHS)

		var err error
		var expr Expression
		if expr, err = parseExpr(node.up); err != nil {
			return nil, err
		}

		exprRHS.expr = expr

		return exprRHS, nil
	default:
		return nil, fmt.Errorf("Unexpected rule %s %s", node.String(), node.match)
	}
}

func parseBaseType(node *node32) (Type, error) {
	switch node.pegRule {
	case ruleINT:
		return &IntType{}, nil
	case ruleBOOL:
		return &BoolType{}, nil
	case ruleCHAR:
		return &CharType{}, nil
	case ruleSTRING:
		return &ArrayType{base: &CharType{}}, nil
	default:
		return nil, fmt.Errorf("Unknown type: %s", node.up.match)
	}
}

func parsePairType(node *node32) (Type, error) {
	var err error

	pairType := &PairType{}

	first := nextNode(node, rulePAIRELEMTYPE)

	if first == nil {
		return pairType, nil
	}

	second := nextNode(first.next, rulePAIRELEMTYPE)

	if pairType.first, err = parseType(first.up); err != nil {
		return nil, err
	}
	if pairType.second, err = parseType(second.up); err != nil {
		return nil, err
	}

	return pairType, nil

}

func parseType(node *node32) (Type, error) {
	var err error
	var waccType Type

	switch node.pegRule {
	case ruleBASETYPE:
		if waccType, err = parseBaseType(node.up); err != nil {
			return nil, err
		}
	case rulePAIRTYPE:
		if waccType, err = parsePairType(node.up); err != nil {
			return nil, err
		}
	}

	for node = nextNode(node.next, ruleARRAYTYPE); node != nil; node = nextNode(node.next, ruleARRAYTYPE) {
		waccType = &ArrayType{base: waccType}
	}

	return waccType, nil
}

func parseStatement(node *node32) (Statement, error) {
	var stm Statement
	var err error

	switch node.pegRule {
	case ruleSKIP:
		stm = &SkipStatement{}
	case ruleBEGIN:
		block := new(BlockStatement)

		bodyNode := nextNode(node, ruleSTAT)
		if block.body, err = parseStatement(bodyNode.up); err != nil {
			return nil, err
		}

		stm = block
	case ruleTYPE:
		decl := new(DeclareAssignStatement)

		typeNode := nextNode(node, ruleTYPE)
		if decl.waccType, err = parseType(typeNode.up); err != nil {
			return nil, err
		}

		identNode := nextNode(node, ruleIDENT)
		decl.ident = identNode.match

		rhsNode := nextNode(node, ruleASSIGNRHS)
		if decl.rhs, err = parseRHS(rhsNode.up); err != nil {
			return nil, err
		}

		stm = decl
	case ruleASSIGNLHS:
		assign := new(AssignStatement)

		lhsNode := nextNode(node, ruleASSIGNLHS)
		if assign.target, err = parseLHS(lhsNode.up); err != nil {
			return nil, err
		}

		rhsNode := nextNode(node, ruleASSIGNRHS)
		if assign.rhs, err = parseRHS(rhsNode.up); err != nil {
			return nil, err
		}

		stm = assign
	case ruleREAD:
		read := new(ReadStatement)

		lhsNode := nextNode(node, ruleASSIGNLHS)
		if read.target, err = parseLHS(lhsNode.up); err != nil {
			return nil, err
		}

		stm = read
	case ruleFREE:
		free := new(FreeStatement)

		exprNode := nextNode(node, ruleEXPR)
		if free.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		stm = free
	case ruleRETURN:
		retur := new(ReturnStatement)

		exprNode := nextNode(node, ruleEXPR)
		if retur.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		stm = retur
	case ruleEXIT:
		exit := new(ExitStatement)

		exprNode := nextNode(node, ruleEXPR)
		if exit.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		stm = exit
	case rulePRINTLN:
		println := new(PrintLnStatement)

		exprNode := nextNode(node, ruleEXPR)
		if println.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		stm = println
	case rulePRINT:
		print := new(PrintStatement)

		exprNode := nextNode(node, ruleEXPR)
		if print.expr, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		stm = print
	case ruleIF:
		ifs := new(IfStatement)

		exprNode := nextNode(node, ruleEXPR)
		if ifs.cond, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		bodyNode := nextNode(node, ruleSTAT)
		if ifs.trueStat, err = parseStatement(bodyNode.up); err != nil {
			return nil, err
		}

		elseNode := nextNode(bodyNode.next, ruleSTAT)
		if ifs.falseStat, err = parseStatement(elseNode.up); err != nil {
			return nil, err
		}

		stm = ifs
	case ruleWHILE:
		whiles := new(WhileStatement)

		exprNode := nextNode(node, ruleEXPR)
		if whiles.cond, err = parseExpr(exprNode.up); err != nil {
			return nil, err
		}

		bodyNode := nextNode(node, ruleSTAT)
		if whiles.body, err = parseStatement(bodyNode.up); err != nil {
			return nil, err
		}

		stm = whiles
	default:
		return nil, fmt.Errorf("unexpected %s %s", node.String(), node.match)
	}

	if semi := nextNode(node, ruleSEMI); semi != nil {
		var next Statement
		if next, err = parseStatement(semi.next.up); err != nil {
			return nil, err
		}
		stm.SetNext(next)
	}

	return stm, nil
}

func parseParam(node *node32) (*FunctionParam, error) {
	var err error

	param := &FunctionParam{}

	param.waccType, err = parseType(nextNode(node, ruleTYPE).up)
	if err != nil {
		return nil, err
	}

	param.name = nextNode(node, ruleIDENT).match

	return param, nil
}

func parseFunction(node *node32) (*FunctionDef, error) {
	var err error
	function := &FunctionDef{}

	function.returnType, err = parseType(nextNode(node, ruleTYPE).up)
	if err != nil {
		return nil, err
	}

	function.ident = nextNode(node, ruleIDENT).match

	paramListNode := nextNode(node, rulePARAMLIST)
	if paramListNode != nil {
		for pnode := range nodeRange(paramListNode.up) {
			if pnode.pegRule == rulePARAM {
				var param *FunctionParam
				param, err = parseParam(pnode.up)
				if err != nil {
					return nil, err
				}
				function.params = append(function.params, param)
			}
		}
	}

	function.body, err = parseStatement(nextNode(node, ruleSTAT).up)
	if err != nil {
		return nil, err
	}

	return function, nil
}

func parseWACC(node *node32) (*AST, error) {
	var err error
	ast := &AST{}

	for node := range nodeRange(node) {
		switch node.pegRule {
		case ruleBEGIN:
		case ruleEND:
		case ruleSPACE:
		case ruleFUNC:
			f, err := parseFunction(node.up)
			ast.functions = append(ast.functions, f)
			if err != nil {
				return nil, err
			}
		case ruleSTAT:
			ast.main, err = parseStatement(node.up)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("Unexpected %s %s", node.String(), node.match)
		}
	}

	return ast, nil
}

func ParseAST(wacc *WACC) (*AST, error) {
	node := wacc.AST()
	switch node.pegRule {
	case ruleWACC:
		return parseWACC(node.up)
	default:
		return nil, errors.New("expected ruleWACC")
	}
}
