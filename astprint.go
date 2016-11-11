package main

// WACC Group 34
//
// astprint.go: Print the AST in the reference compiler's format.
//
// File contains functions that return indented strings and produces an AST
// indented with the same format as the reference compiler

import (
	"fmt"
	"strconv"
)

//Util Functions:
//------------------------------------------------------------------------------

// Adds a dash + space at the end of the current indent
// Is invoked every time before printing
func addMinToIndent(indent string) string {
	return (indent + "- ")
}

// Adds the current string at the next indent level
func addAtGreaterIndent(indent, value string) string {
	return addIndAndNewLine(getGreaterIndent(indent), value)
}

// Adds the current string at the given indent value
func addIndAndNewLine(indent, value string) string {
	return fmt.Sprintf("%v%v\n", addMinToIndent(indent), value)
}

// Adds two strings s1 and s2, indenting s1 only
func addIndentForFirst(indent, a1, a2 string) string {
	return fmt.Sprintf("%v%v", addIndAndNewLine(indent, a1), a2)
}

// Adds two strings s1 and s2 indenting both s1 and s2
func addDoubleIndent(indent, a1, a2 string) string {
	return addTripleIndent(indent, a1, a2, "")
}

// Adds a string s1 and a []string t.
// Only s1 is indented
func addArrayIndent(indent, a1 string, arr []string) string {
	var innerStats string

	typeStats := addIndAndNewLine(indent, a1)
	for _, element := range arr {
		innerStats = fmt.Sprintf("%v%v", innerStats, element)
	}

	return fmt.Sprintf("%v%v", typeStats, innerStats)
}

// Returns the next indenting level given by basicIndent
func getGreaterIndent(indent string) string {
	return fmt.Sprintf("%v%v", indent, basicIndent)
}

// Adds soubleIndent with caption "TYPE"
func addType(indent, argument string) string {
	return addDoubleIndent(indent, "TYPE", argument)
}

// Adds three strings s1, s2 and s3, indenting s1 only
func addTripleIndentOnlyFst(indent, a1, a2, a3 string) string {
	innerStats2 := a3

	typeStats := addIndAndNewLine(indent, a1)
	innerStats := a2
	if a3 != "" {
		innerStats2 = a3
	}

	return fmt.Sprintf("%v%v%v", typeStats, innerStats, innerStats2)
}

// Adds three strings s1, s2 and s3, indenting all strings
func addTripleIndent(indent, a1, a2, a3 string) string {
	innerStats2 := a3

	innerIndent := getGreaterIndent(indent)

	typeStats := addIndAndNewLine(indent, a1)
	innerStats := addIndAndNewLine(innerIndent, a2)
	if a3 != "" {
		innerStats2 = addIndAndNewLine(innerIndent, a3)
	}

	return fmt.Sprintf("%v%v%v", typeStats, innerStats, innerStats2)
}

//------------------------------------------------------------------------------

func (p PairType) ASTString(indent string) string {
	var first string
	var second string

	if p.first != nil {
		first = fmt.Sprintf("%v", p.first.ASTString(indent))
	}
	if p.second != nil {
		second = fmt.Sprintf("%v", p.second.ASTString(indent))
	}
	return fmt.Sprintf("%v%v", first, second)
}

func (stmt DeclareAssignStatement) ASTString(indent string) string {

	declareStats := fmt.Sprintf("%vDECLARE\n", addMinToIndent(indent))
	innerIndent := getGreaterIndent(indent)
	lhsIndent := addDoubleIndent(innerIndent, "LHS", stmt.ident)
	rhsIndent := addIndentForFirst(
		innerIndent,
		"RHS",
		stmt.rhs.ASTString(getGreaterIndent(innerIndent)),
	)

	return fmt.Sprintf(
		"%v%v%v%v",
		declareStats,
		stmt.waccType.ASTString(innerIndent),
		lhsIndent,
		rhsIndent,
	)
}

func (lhs PairElemLHS) ASTString(indent string) string {
	if lhs.snd {
		return addIndentForFirst(
			indent,
			"SND",
			lhs.expr.ASTString(getGreaterIndent(indent)),
		)
	}
	return addIndentForFirst(
		indent,
		"FST",
		lhs.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (lhs ArrayLHS) ASTString(indent string) string {
	var indexes string
	var tmpIndex string

	nextIndent := getGreaterIndent(indent)

	for _, index := range lhs.index {
		tmpIndex = addIndentForFirst(
			nextIndent,
			"[]",
			index.ASTString(getGreaterIndent(nextIndent)),
		)
		indexes = fmt.Sprintf("%v%v", indexes, tmpIndex)
	}

	return addIndentForFirst(indent, lhs.ident, indexes)
}

func (lhs VarLHS) ASTString(indent string) string {
	return addIndAndNewLine(indent, lhs.ident)
}

func (rhs PairLiterRHS) ASTString(indent string) string {
	nextIndent := getGreaterIndent(indent)
	fstStats := addIndentForFirst(
		nextIndent,
		"FST",
		rhs.fst.ASTString(getGreaterIndent(nextIndent)),
	)
	sndStats := addIndentForFirst(
		nextIndent,
		"SND",
		rhs.snd.ASTString(getGreaterIndent(nextIndent)),
	)
	return addTripleIndentOnlyFst(indent, "NEWPAIR", fstStats, sndStats)
}

func (rhs ArrayLiterRHS) ASTString(indent string) string {
	elemArr := []string{}

	for _, element := range rhs.elements {
		elemArr = append(elemArr, element.ASTString(getGreaterIndent((indent))))
	}

	return addArrayIndent(indent, "ARRAY LITERAL", elemArr)
}

func (rhs PairElemRHS) ASTString(indent string) string {
	if rhs.snd {
		return addIndentForFirst(
			indent,
			"SND",
			rhs.expr.ASTString(getGreaterIndent(indent)),
		)
	}
	return addIndentForFirst(indent,
		"FST",
		rhs.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (rhs FunctionCallRHS) ASTString(indent string) string {
	var innerStats string
	nameStats := addIndAndNewLine(indent, rhs.ident)

	for _, param := range rhs.args {
		innerStats = fmt.Sprintf(
			"%v%v",
			innerStats,
			param.ASTString(indent),
		)
	}

	return fmt.Sprintf("%v%v", nameStats, innerStats)
}

func (stmt AssignStatement) ASTString(indent string) string {
	declareStats := fmt.Sprintf("%vASSIGNMENT\n", addMinToIndent(indent))
	innerIndent := getGreaterIndent(indent)
	lhsIndent := addIndentForFirst(
		innerIndent,
		"LHS",
		stmt.target.ASTString(getGreaterIndent(innerIndent)),
	)
	rhsIndent := addIndentForFirst(
		innerIndent,
		"RHS",
		stmt.rhs.ASTString(getGreaterIndent(innerIndent)),
	)

	return fmt.Sprintf("%v%v%v", declareStats, lhsIndent, rhsIndent)
}

func (stmt IfStatement) ASTString(indent string) string {
	var stmtStats string
	var trueStats string
	var falseStats string
	var trueTmp string
	var falseTmp string
	innerIndent := getGreaterIndent(indent)
	doubleInnerIndent := getGreaterIndent(innerIndent)

	ifStats := fmt.Sprintf("%vIF\n", addMinToIndent(indent))
	condStats := fmt.Sprintf("%vCONDITION\n", addMinToIndent(innerIndent))
	thenStats := fmt.Sprintf("%vTHEN\n", addMinToIndent(innerIndent))
	elseStats := fmt.Sprintf("%vELSE\n", addMinToIndent(innerIndent))

	stmtStats = stmt.cond.ASTString(doubleInnerIndent)

	st := stmt.trueStat
	for st.GetNext() != nil {
		trueTmp = st.ASTString(doubleInnerIndent)
		trueStats = fmt.Sprintf("%v%v", trueStats, trueTmp)
		st = st.GetNext()
	}

	trueTmp = st.ASTString(doubleInnerIndent)
	trueStats = fmt.Sprintf("%v%v", trueStats, trueTmp)

	st = stmt.falseStat
	for st.GetNext() != nil {
		falseTmp = st.ASTString(doubleInnerIndent)
		falseStats = fmt.Sprintf("%v%v", falseStats, falseTmp)
		st = st.GetNext()
	}

	falseTmp = st.ASTString(doubleInnerIndent)
	falseStats = fmt.Sprintf("%v%v", falseStats, falseTmp)

	return fmt.Sprintf(
		"%v%v%v%v%v%v%v",
		ifStats,
		condStats,
		stmtStats,
		thenStats,
		trueStats,
		elseStats,
		falseStats,
	)
}

func (stmt WhileStatement) ASTString(indent string) string {
	var body string
	var doStats string
	innerIndent := getGreaterIndent(indent)

	doStats = addIndAndNewLine(innerIndent, "DO")

	condStats := addIndentForFirst(
		innerIndent,
		"CONDITION",
		stmt.cond.ASTString(getGreaterIndent(innerIndent)),
	)

	st := stmt.body
	for st.GetNext() != nil {
		body = st.ASTString(getGreaterIndent(innerIndent))
		doStats = fmt.Sprintf("%v%v", doStats, body)
		st = st.GetNext()
	}
	body = st.ASTString(getGreaterIndent(innerIndent))
	doStats = fmt.Sprintf("%v%v", doStats, body)

	loopStats := addIndAndNewLine(indent, "LOOP")

	return fmt.Sprintf("%v%v%v", loopStats, condStats, doStats)
}

func (fp FunctionParam) ASTString(indent string) string {
	return fmt.Sprintf("%v %v", fp.waccType, fp.name)
}

func (fd FunctionDef) ASTString(indent string) string {

	var params string
	var body string

	innerIndent := getGreaterIndent(indent)

	if len(fd.params) > 0 {
		params = fmt.Sprintf("%v", fd.params[0])

		for _, param := range fd.params[1:] {
			params = fmt.Sprintf("%v, %v", params, param)
		}
	}

	declaration :=
		addIndAndNewLine(indent,
			fmt.Sprintf(
				"%v %v(%v)",
				fd.returnType,
				fd.ident,
				params))

	st := fd.body
	for st.GetNext() != nil {
		body = fmt.Sprintf("%v%v", body, st.ASTString(innerIndent))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v%v", body, st.ASTString(innerIndent))

	return fmt.Sprintf("%v%v", declaration, body)
}

func (elem ArrayElem) ASTString(indent string) string {
	var indexes string
	var tmpIndex string

	nextIndent := getGreaterIndent(indent)

	for _, index := range elem.indexes {
		tmpIndex =
			addIndentForFirst(nextIndent,
				"[]",
				index.ASTString(getGreaterIndent(nextIndent)))
		indexes = fmt.Sprintf("%v%v", indexes, tmpIndex)
	}

	return addIndentForFirst(indent, elem.ident, indexes)
}

func (a ArrayType) ASTString(indent string) string {
	typeStats := fmt.Sprintf("%v[]", a.base)
	return addType(indent, typeStats)
}

func (i InvalidType) ASTString(indent string) string {
	return addType(indent, "<invalid>")
}

func (i UnknownType) ASTString(indent string) string {
	return addType(indent, "<unknown>")
}

func (i IntType) ASTString(indent string) string {
	return addType(indent, "int")
}

func (b BoolType) ASTString(indent string) string {
	return addType(indent, "bool")
}

func (c CharType) ASTString(indent string) string {
	return addType(indent, "char")
}

func (stmt SkipStatement) ASTString(indent string) string {
	return addIndAndNewLine(indent, "SKIP")
}

func (stmt BlockStatement) ASTString(indent string) string {
	return fmt.Sprintf("")
}

func (lpar ExprParen) ASTString(indent string) string {
	return ""
}

func (rhs ExpressionRHS) ASTString(indent string) string {
	return rhs.expr.ASTString(indent)
}

func (stmt ReadStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"READ",
		stmt.target.ASTString(getGreaterIndent(indent)),
	)
}

func (stmt FreeStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"FREE",
		stmt.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (ret ReturnStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"RETURN",
		ret.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (stmt ExitStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"EXIT",
		stmt.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (stmt PrintLnStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"PRINTLN",
		stmt.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (stmt PrintStatement) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"PRINT",
		stmt.expr.ASTString(getGreaterIndent(indent)),
	)
}

func (ident Ident) ASTString(indent string) string {
	return addIndAndNewLine(indent, ident.ident)
}

func (liter IntLiteral) ASTString(indent string) string {
	return addIndAndNewLine(indent, strconv.Itoa(liter.value))
}

func (liter BoolLiteralTrue) ASTString(indent string) string {
	return addIndAndNewLine(indent, "true")
}

func (liter BoolLiteralFalse) ASTString(indent string) string {
	return addIndAndNewLine(indent, "false")
}

func (liter CharLiteral) ASTString(indent string) string {
	tmpStats := fmt.Sprintf("'%v'", liter.char)
	return addIndAndNewLine(indent, tmpStats)
}

func (liter StringLiteral) ASTString(indent string) string {
	tmp := fmt.Sprintf("\"%v\"", liter.str)
	return addIndAndNewLine(indent, tmp)
}

func (liter PairLiteral) ASTString(indent string) string {
	return fmt.Sprintf("pair(%v, %v)", liter.fst, liter.snd)
}

func (liter NullPair) ASTString(indent string) string {
	return addIndAndNewLine(indent, "null")
}

func (op UnaryOperatorNot) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"!",
		op.GetExpression().ASTString(getGreaterIndent(indent)),
	)
}

func (op UnaryOperatorNegate) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"-",
		op.GetExpression().ASTString(getGreaterIndent(indent)),
	)
}

func (op UnaryOperatorLen) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"len",
		op.GetExpression().ASTString(getGreaterIndent(indent)),
	)
}

func (op UnaryOperatorOrd) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"ord",
		op.GetExpression().ASTString(getGreaterIndent(indent)),
	)
}

func (op UnaryOperatorChr) ASTString(indent string) string {
	return addIndentForFirst(
		indent,
		"chr",
		op.GetExpression().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorMult) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"*",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorDiv) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"/",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorMod) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"%",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorAdd) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"+",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorSub) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"-",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorGreaterThan) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		">",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorGreaterEqual) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		">=",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorLessThan) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"<",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorLessEqual) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"<=",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorEqual) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"==",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorNotEqual) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"!=",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorAnd) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"&&",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

func (op BinaryOperatorOr) ASTString(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"||",
		op.GetLHS().ASTString(getGreaterIndent(indent)),
		op.GetRHS().ASTString(getGreaterIndent(indent)),
	)
}

//------------------------------------------------------------------------------

func (ast AST) ASTString() string {
	var tree string
	var tmpIndent string

	tree = addIndAndNewLine("", "Program")

	for _, function := range ast.functions {
		tree = fmt.Sprintf(
			"%v%v",
			tree,
			function.ASTString(basicIndent),
		)
	}

	tmpIndent = getGreaterIndent(basicIndent)

	tree = fmt.Sprintf("%v%v",
		tree,
		addIndAndNewLine(basicIndent, "int main()"),
	)

	stmt := ast.main
	for stmt.GetNext() != nil {
		tree = fmt.Sprintf("%v%v", tree, stmt.ASTString(tmpIndent))
		stmt = stmt.GetNext()
	}
	tree = fmt.Sprintf("%v%v", tree, stmt.ASTString(tmpIndent))

	return tree
}
