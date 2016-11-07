package main

import (
	"errors"
	"fmt"
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
	fst Expression
	snd Expression
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

type AST struct {
	main Statement
}

func nextNode(node *node32, rule pegRule) *node32 {
	for ; node != nil; node = node.next {
		if node.pegRule == rule {
			return node
		}
	}

	return nil
}

func parseExpr(node *node32) (Expression, error) {
	// TODO implement
	return nil, fmt.Errorf("Not implemented")
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
		return nil, fmt.Errorf("Unexpected %s", node.match)
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

		for argNode := nextNode(arglistNode.up, ruleEXPR); argNode != nil; argNode = nextNode(argNode, ruleEXPR) {
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
		return nil, fmt.Errorf("Unexpected rule %s", node.match)
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
		return nil, fmt.Errorf("unexpected %s", node.match)
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

func parseWACC(node *node32) (*AST, error) {
	var err error
	ast := &AST{}

	for node != nil {
		switch node.pegRule {
		case ruleBEGIN:
			ast.main, err = parseStatement(node.next.up)
			if err != nil {
				return nil, err
			}
			return ast, nil
		default:
			return nil, fmt.Errorf("Unexpected %s", node.match)
		}
	}

	return nil, errors.New("expected ruleBEGIN")
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
