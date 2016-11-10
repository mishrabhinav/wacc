package main

import (
	"fmt"
)

func (i IntType) ASTString(indent string) string {
	return fmt.Sprintf("int")
}

func (b BoolType) ASTString(indent string) string {
	return fmt.Sprintf("bool")
}

func (c CharType) ASTString(indent string) string {
	return fmt.Sprintf("char")
}

func (p PairType) ASTString(indent string) string {
	var first string = fmt.Sprintf("%v", p.first)
	var second string = fmt.Sprintf("%v", p.second)

	if p.first == nil {
		first = "pair"
	}
	if p.second == nil {
		second = "pair"
	}
	return fmt.Sprintf("pair(%v, %v)", first, second)
}

func (a ArrayType) ASTString(indent string) string {
	return fmt.Sprintf("%v[]", a.base)
}

func (stmt SkipStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vskip", indent)
}

func (stmt BlockStatement) ASTString(indent string) string {
	return fmt.Sprintf("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXx")
}

func (stmt DeclareAssignStatement) ASTString(indent string) string {
	return fmt.Sprintf("%v%v %v = %v", indent, stmt.waccType, stmt.ident, stmt.rhs)
}

func (lhs PairElemLHS) ASTString(indent string) string {
	if lhs.snd {
		return fmt.Sprintf("snd %v", lhs.expr)
	} else {
		return fmt.Sprintf("fst %v", lhs.expr)
	}
}

func (lhs ArrayLHS) ASTString(indent string) string {
	var indexes string

	for _, index := range lhs.index {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", lhs.ident, indexes)
}

func (lhs VarLHS) ASTString(indent string) string {
	return fmt.Sprintf(lhs.ident)
}

func (rhs PairLiterRHS) ASTString(indent string) string {
	return fmt.Sprintf("newpair(%v, %v)", rhs.fst, rhs.snd)
}

func (rhs ArrayLiterRHS) ASTString(indent string) string {
	var elements string

	if len(rhs.elements) > 0 {
		elements = fmt.Sprintf("%v", rhs.elements[0])

		for _, element := range rhs.elements[1:] {
			elements = fmt.Sprintf("%v, %v", elements, element)
		}
	}

	return fmt.Sprintf("[%v]", elements)
}

func (rhs PairElemRHS) ASTString(indent string) string {
	if rhs.snd {
		return fmt.Sprintf("snd %v", rhs.expr)
	} else {
		return fmt.Sprintf("fst %v", rhs.expr)
	}
}

func (rhs FunctionCallRHS) ASTString(indent string) string {
	var params string

	if len(rhs.args) > 0 {
		params = fmt.Sprintf("%v", rhs.args[0])

		for _, param := range rhs.args[1:] {
			params = fmt.Sprintf("%v, %v", params, param)
		}
	}

	return fmt.Sprintf("call %v(%v)", rhs.ident, params)
}

func (lpar ExprLPar) ASTString(indent string) string {
	return ""
}

func (rpar ExprRPar) ASTString(indent string) string {
	return ""
}

func (rhs ExpressionRHS) ASTString(indent string) string {
	return fmt.Sprintf("%v", rhs.expr)
}

func (stmt AssignStatement) ASTString(indent string) string {
	return fmt.Sprintf("%v%v = %v", indent, stmt.target, stmt.rhs)
}

func (stmt ReadStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vread %v", indent, stmt.target)
}

func (stmt FreeStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vfree %v", indent, stmt.expr)
}

func (ret ReturnStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vreturn %v", indent, ret.expr)
}

func (stmt ExitStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vexit %v", indent, stmt.expr)
}

func (stmt PrintLnStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vprintln %v", indent, stmt.expr)
}

func (stmt PrintStatement) ASTString(indent string) string {
	return fmt.Sprintf("%vprint %v", indent, stmt.expr)
}

func (stmt IfStatement) ASTString(indent string) string {
	var trueStats string
	var falseStats string
	var innerIndent string = fmt.Sprintf("%v  ", indent)

	st := stmt.trueStat
	for st.GetNext() != nil {
		trueStats = fmt.Sprintf("%v\n%v ;", trueStats, st.ASTString(innerIndent))
		st = st.GetNext()
	}

	trueStats = fmt.Sprintf("%v\n%v", trueStats, st.ASTString(innerIndent))

	st = stmt.falseStat
	for st.GetNext() != nil {
		falseStats = fmt.Sprintf("%v\n%v ;", falseStats, st.ASTString(innerIndent))
		st = st.GetNext()
	}

	falseStats = fmt.Sprintf("%v\n%v", falseStats, st.ASTString(innerIndent))

	return fmt.Sprintf("%vif %v then %v\n%velse %v\n%vfi", indent, stmt.cond, trueStats, indent, falseStats, indent)
}

func (stmt WhileStatement) ASTString(indent string) string {
	var body string
	var innerIndent string = fmt.Sprintf("%v  ", indent)

	st := stmt.body
	for st.GetNext() != nil {
		body = fmt.Sprintf("%v\n%v ;", body, st.ASTString(innerIndent))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v\n%v", body, st.ASTString(innerIndent))

	return fmt.Sprintf("%vwhile %v do%v\n%vdone", indent, stmt.cond, body, indent)
}

func (fp FunctionParam) ASTString(indent string) string {
	return fmt.Sprintf("%v %v", fp.waccType, fp.name)
}

func (fd FunctionDef) ASTString(indent string) string {

	var params string
	var body string

	innerIndent := fmt.Sprintf("%v  ", indent)

	if len(fd.params) > 0 {
		params = fmt.Sprintf("%v", fd.params[0])

		for _, param := range fd.params[1:] {
			params = fmt.Sprintf("%v, %v", params, param)
		}
	}

	declaration := fmt.Sprintf("%v%v %v(%v) is", indent, fd.returnType, fd.ident, params)

	st := fd.body
	for st.GetNext() != nil {
		body = fmt.Sprintf("%v\n%v ;", body, st.ASTString(innerIndent))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v\n%v", body, st.ASTString(innerIndent))

	return fmt.Sprintf("%v %v\n%vend", declaration, body, indent)
}

func (ident Ident) ASTString(indent string) string {
	return fmt.Sprintf(ident.ident)
}

func (liter IntLiteral) ASTString(indent string) string {
	return fmt.Sprintf("%v", liter.value)
}

func (liter BoolLiteralTrue) ASTString(indent string) string {
	return fmt.Sprintf("true")
}

func (liter BoolLiteralFalse) ASTString(indent string) string {
	return fmt.Sprintf("false")
}

func (liter CharLiteral) ASTString(indent string) string {
	return fmt.Sprintf("'%v'", liter.char)
}

func (liter StringLiteral) ASTString(indent string) string {
	return fmt.Sprintf("\"%v\"", liter.str)
}

func (liter PairLiteral) ASTString(indent string) string {
	return fmt.Sprintf("pair(%v, %v)", liter.fst, liter.snd)
}

func (liter NullPair) ASTString(indent string) string {
	return fmt.Sprintf("null")
}

func (elem ArrayElem) ASTString(indent string) string {
	var indexes string

	for _, index := range elem.indexes {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", elem.ident, indexes)
}

func (op UnaryOperatorNot) ASTString(indent string) string {
	return fmt.Sprintf("!%v", op.GetExpression())
}

func (op UnaryOperatorNegate) ASTString(indent string) string {
	return fmt.Sprintf("-%v", op.GetExpression())
}

func (op UnaryOperatorLen) ASTString(indent string) string {
	return fmt.Sprintf("len %v", op.GetExpression())
}

func (op UnaryOperatorOrd) ASTString(indent string) string {
	return fmt.Sprintf("ord %v", op.GetExpression())
}

func (op UnaryOperatorChr) ASTString(indent string) string {
	return fmt.Sprintf("chr %v", op.GetExpression())
}

func (op BinaryOperatorMult) ASTString(indent string) string {
	return fmt.Sprintf("%v * %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorDiv) ASTString(indent string) string {
	return fmt.Sprintf("%v / %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorMod) ASTString(indent string) string {
	return fmt.Sprintf("%v % %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorAdd) ASTString(indent string) string {
	return fmt.Sprintf("%v + %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorSub) ASTString(indent string) string {
	return fmt.Sprintf("%v - %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorGreaterThan) ASTString(indent string) string {
	return fmt.Sprintf("%v > %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorGreaterEqual) ASTString(indent string) string {
	return fmt.Sprintf("%v >= %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorLessThan) ASTString(indent string) string {
	return fmt.Sprintf("%v < %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorLessEqual) ASTString(indent string) string {
	return fmt.Sprintf("%v <= %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorEqual) ASTString(indent string) string {
	return fmt.Sprintf("%v == %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorNotEqual) ASTString(indent string) string {
	return fmt.Sprintf("%v != %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorAnd) ASTString(indent string) string {
	return fmt.Sprintf("%v && %v", op.GetLHS(), op.GetRHS())
}

func (op BinaryOperatorOr) ASTString(indent string) string {
	return fmt.Sprintf("%v || %v", op.GetLHS(), op.GetRHS())
}

func (ast AST) ASTString() string {
	var tree string

	tree = fmt.Sprintf("begin")

	for _, function := range ast.functions {
		tree = fmt.Sprintf("%v\n%v", tree, function.ASTString(basicIndent))
	}

	stmt := ast.main
	for stmt.GetNext() != nil {
		tree = fmt.Sprintf("%v\n%v ;", tree, stmt.ASTString(basicIndent))
		stmt = stmt.GetNext()
	}

	tree = fmt.Sprintf("%v\n%v", tree, stmt.ASTString(basicIndent))

	return fmt.Sprintf("%v\nend", tree)
}
