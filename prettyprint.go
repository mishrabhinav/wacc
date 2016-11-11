package main

// WACC Group 34
//
// prettyprint.go: Pretty print the AST.
//
// File contains functions that return indented strings and produces a pretty
// string for the abstract syntax tree.

import (
	"fmt"
	"strings"
)

const startingIndent int = 1
const basicIndent string = "  "

func getIndentation(level int) string {
	return fmt.Sprint(strings.Repeat(basicIndent, level))
}

func generateUnaryOperator(expr Expression, op string) string {
	return fmt.Sprintf("%v%v", op, expr)
}

func generateBinaryOperator(lhs Expression, rhs Expression, op string) string {
	return fmt.Sprintf("%v %v %v", lhs, op, rhs)
}

func (m InvalidType) String() string {
	return "<invalid>"
}

func (m UnknownType) String() string {
	return "<unknown>"
}

func (i IntType) String() string {
	return "int"
}

func (b BoolType) String() string {
	return "bool"
}

func (c CharType) String() string {
	return "char"
}

func (p PairType) String() string {
	var first = fmt.Sprintf("%v", p.first)
	var second = fmt.Sprintf("%v", p.second)

	if p.first == nil {
		first = "pair"
	}
	if p.second == nil {
		second = "pair"
	}
	return fmt.Sprintf("pair(%v, %v)", first, second)
}

func (a ArrayType) String() string {
	return fmt.Sprintf("%v[]", a.base)
}

func (ident *Ident) String() string {
	return ident.ident
}

func (liter *IntLiteral) String() string {
	return fmt.Sprint(liter.value)
}

func (liter *BoolLiteralTrue) String() string {
	return "true"
}

func (liter *BoolLiteralFalse) String() string {
	return "false"
}

func (liter *CharLiteral) String() string {
	return fmt.Sprintf("'%v'", liter.char)
}

func (liter *StringLiteral) String() string {
	return fmt.Sprintf("\"%v\"", liter.str)
}

func (liter *PairLiteral) String() string {
	return fmt.Sprintf("pair(%v, %v)", liter.fst, liter.snd)
}

func (liter *NullPair) String() string {
	return "null"
}

func (elem *ArrayElem) String() string {
	var indexes string

	for _, index := range elem.indexes {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", elem.ident, indexes)
}

func (op *UnaryOperatorNot) String() string {
	return generateUnaryOperator(op.GetExpression(), "!")
}

func (op *UnaryOperatorNegate) String() string {
	return generateUnaryOperator(op.GetExpression(), "-")
}

func (op *UnaryOperatorLen) String() string {
	return generateUnaryOperator(op.GetExpression(), "len ")
}

func (op *UnaryOperatorOrd) String() string {
	return generateUnaryOperator(op.GetExpression(), "ord ")
}

func (op *UnaryOperatorChr) String() string {
	return generateUnaryOperator(op.GetExpression(), "chr ")
}

func (op *BinaryOperatorMult) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "*")
}

func (op *BinaryOperatorDiv) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "/")
}

func (op *BinaryOperatorMod) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "%")
}

func (op *BinaryOperatorAdd) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "+")
}

func (op *BinaryOperatorSub) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "-")
}

func (op *BinaryOperatorGreaterThan) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), ">")
}

func (op *BinaryOperatorGreaterEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), ">=")
}

func (op *BinaryOperatorLessThan) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "<")
}

func (op *BinaryOperatorLessEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "<=")
}

func (op *BinaryOperatorEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "==")
}

func (op *BinaryOperatorNotEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "!=")
}

func (op *BinaryOperatorAnd) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "&&")
}

func (op *BinaryOperatorOr) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "||")
}

func (lhs *PairElemLHS) String() string {
	if lhs.snd {
		return fmt.Sprintf("snd %v", lhs.expr)
	}

	return fmt.Sprintf("fst %v", lhs.expr)
}

func (lhs *ArrayLHS) String() string {
	var indexes string

	for _, index := range lhs.index {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", lhs.ident, indexes)
}

func (lhs *VarLHS) String() string {
	return fmt.Sprintf(lhs.ident)
}

func (rhs *PairLiterRHS) String() string {
	return fmt.Sprintf("newpair(%v, %v)", rhs.fst, rhs.snd)
}

func (rhs *ArrayLiterRHS) String() string {
	var elements string

	if len(rhs.elements) > 0 {
		elements = fmt.Sprintf("%v", rhs.elements[0])

		for _, element := range rhs.elements[1:] {
			elements = fmt.Sprintf("%v, %v", elements, element)
		}
	}

	return fmt.Sprintf("[%v]", elements)
}

func (rhs *PairElemRHS) String() string {
	if rhs.snd {
		return fmt.Sprintf("snd %v", rhs.expr)
	}

	return fmt.Sprintf("fst %v", rhs.expr)
}

func (rhs *FunctionCallRHS) String() string {
	var params string

	if len(rhs.args) > 0 {
		params = fmt.Sprintf("%v", rhs.args[0])

		for _, param := range rhs.args[1:] {
			params = fmt.Sprintf("%v, %v", params, param)
		}
	}

	return fmt.Sprintf("call %v(%v)", rhs.ident, params)
}

func (rhs *ExpressionRHS) String() string {
	return fmt.Sprintf("%v", rhs.expr)
}

func (stmt *SkipStatement) istring(level int) string {
	return fmt.Sprintf("%vskip", getIndentation(level))
}

func (stmt *BlockStatement) istring(level int) string {
	return ""
}

func (stmt *DeclareAssignStatement) istring(level int) string {
	return fmt.Sprintf("%v%v %v = %v", getIndentation(level), stmt.waccType,
		stmt.ident, stmt.rhs)
}
func (stmt *AssignStatement) istring(level int) string {
	return fmt.Sprintf("%v%v = %v", getIndentation(level), stmt.target,
		stmt.rhs)
}

func (stmt *ReadStatement) istring(level int) string {
	return fmt.Sprintf("%vread %v", getIndentation(level), stmt.target)
}

func (stmt *FreeStatement) istring(level int) string {
	return fmt.Sprintf("%vfree %v", getIndentation(level), stmt.expr)
}

func (ret *ReturnStatement) istring(level int) string {
	return fmt.Sprintf("%vreturn %v", getIndentation(level), ret.expr)
}

func (stmt *ExitStatement) istring(level int) string {
	return fmt.Sprintf("%vexit %v", getIndentation(level), stmt.expr)
}

func (stmt *PrintLnStatement) istring(level int) string {
	return fmt.Sprintf("%vprintln %v", getIndentation(level), stmt.expr)
}

func (stmt *PrintStatement) istring(level int) string {
	return fmt.Sprintf("%vprint %v", getIndentation(level), stmt.expr)
}

func (stmt *IfStatement) istring(level int) string {
	var trueStats string
	var falseStats string

	var indent = getIndentation(level)

	st := stmt.trueStat
	for st.GetNext() != nil {
		trueStats = fmt.Sprintf("%v\n%v ;", trueStats,
			st.istring(level+1))
		st = st.GetNext()
	}

	trueStats = fmt.Sprintf("%v\n%v", trueStats, st.istring(level+1))

	st = stmt.falseStat
	for st.GetNext() != nil {
		falseStats = fmt.Sprintf("%v\n%v ;", falseStats,
			st.istring(level+1))
		st = st.GetNext()
	}

	falseStats = fmt.Sprintf("%v\n%v", falseStats, st.istring(level+1))

	return fmt.Sprintf("%vif %v\n%vthen %v\n%velse %v\n%vfi", indent,
		stmt.cond, indent, trueStats, indent, falseStats, indent)
}

func (stmt *WhileStatement) istring(level int) string {
	var body string
	var indent = getIndentation(level)

	st := stmt.body
	for st.GetNext() != nil {
		body = fmt.Sprintf("%v\n%v ;", body, st.istring(level+1))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v\n%v", body, st.istring(level+1))

	return fmt.Sprintf("%vwhile (%v) do%v\n%vdone", indent, stmt.cond, body,
		indent)
}

func (fp *FunctionParam) String() string {
	return fmt.Sprintf("%v %v", fp.waccType, fp.name)
}

func (fd *FunctionDef) istring(level int) string {
	var params string
	var body string

	indent := getIndentation(level)

	if len(fd.params) > 0 {
		params = fmt.Sprintf("%v", fd.params[0])

		for _, param := range fd.params[1:] {
			params = fmt.Sprintf("%v, %v", params, param)
		}
	}

	declaration := fmt.Sprintf("%v%v %v(%v) is", indent, fd.returnType,
		fd.ident, params)

	st := fd.body
	for st.GetNext() != nil {
		body = fmt.Sprintf("%v\n%v ;", body, st.istring(level+1))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v\n%v", body, st.istring(level+1))

	return fmt.Sprintf("%v %v\n%vend", declaration, body, indent)
}

func (ast *AST) String() string {
	var tree string

	tree = fmt.Sprintf("begin")

	for _, function := range ast.functions {
		tree = fmt.Sprintf("%v\n%v\n", tree,
			function.istring(startingIndent))
	}

	stmt := ast.main
	for stmt.GetNext() != nil {
		tree = fmt.Sprintf("%v\n%v ;", tree,
			stmt.istring(startingIndent))
		stmt = stmt.GetNext()
	}

	tree = fmt.Sprintf("%v\n%v", tree, stmt.istring(startingIndent))

	return fmt.Sprintf("%v\nend", tree)
}
