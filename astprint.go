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

// Print the type of a given Pair(fst T1, snd T2), given that T1/T2 are not null
func (p PairType) aststring(indent string) string {
	var first string
	var second string

	if p.first != nil {
		first = fmt.Sprintf("%v", p.first.aststring(indent))
	}
	if p.second != nil {
		second = fmt.Sprintf("%v", p.second.aststring(indent))
	}
	return fmt.Sprintf("%v%v", first, second)
}

// Prints a DECLARE statement. Format:
// - DECLARE
//   - LHS
//     - [lhsEXPR]
//   - RHS
//     - [rhsEXPR]
// REcurses on lhsEXPR and rhsEXPR.
func (stmt DeclareAssignStatement) aststring(indent string) string {
	declareStats := fmt.Sprintf("%vDECLARE\n", addMinToIndent(indent))
	innerIndent := getGreaterIndent(indent)
	lhsIndent := addDoubleIndent(innerIndent, "LHS", stmt.ident)
	rhsIndent := addIndentForFirst(
		innerIndent,
		"RHS",
		stmt.rhs.aststring(getGreaterIndent(innerIndent)),
	)

	return fmt.Sprintf(
		"%v%v%v%v",
		declareStats,
		stmt.waccType.aststring(innerIndent),
		lhsIndent,
		rhsIndent,
	)
}

// Prints the LHS of a PairElem
func (lhs PairElemLHS) aststring(indent string) string {
	if lhs.snd {
		return addIndentForFirst(
			indent,
			"SND",
			lhs.expr.aststring(getGreaterIndent(indent)),
		)
	}
	return addIndentForFirst(
		indent,
		"FST",
		lhs.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints (listing), the array's identifier (lhs.ident), "[]", and its elems
func (lhs ArrayLHS) aststring(indent string) string {
	var indexes string
	var tmpIndex string

	nextIndent := getGreaterIndent(indent)

	for _, index := range lhs.index {
		tmpIndex = addIndentForFirst(
			nextIndent,
			"[]",
			index.aststring(getGreaterIndent(nextIndent)),
		)
		indexes = fmt.Sprintf("%v%v", indexes, tmpIndex)
	}

	return addIndentForFirst(indent, lhs.ident, indexes)
}

// Prints the LHS of a Variable
func (lhs VarLHS) aststring(indent string) string {
	return addIndAndNewLine(indent, lhs.ident)
}

// Prints the RHS of a PairLiteral
func (rhs PairLiterRHS) aststring(indent string) string {
	nextIndent := getGreaterIndent(indent)
	fstStats := addIndentForFirst(
		nextIndent,
		"FST",
		rhs.fst.aststring(getGreaterIndent(nextIndent)),
	)
	sndStats := addIndentForFirst(
		nextIndent,
		"SND",
		rhs.snd.aststring(getGreaterIndent(nextIndent)),
	)
	return addTripleIndentOnlyFst(indent, "NEWPAIR", fstStats, sndStats)
}

// Introduces an Array Literal. Recurses on rhs.elements and prints.
func (rhs ArrayLiterRHS) aststring(indent string) string {
	elemArr := []string{}

	for _, element := range rhs.elements {
		elemArr = append(elemArr, element.aststring(getGreaterIndent((indent))))
	}

	return addArrayIndent(indent, "ARRAY LITERAL", elemArr)
}

// Prints the RHS of a PairElem
func (rhs PairElemRHS) aststring(indent string) string {
	if rhs.snd {
		return addIndentForFirst(
			indent,
			"SND",
			rhs.expr.aststring(getGreaterIndent(indent)),
		)
	}
	return addIndentForFirst(indent,
		"FST",
		rhs.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints the RHS of a function call
// Lists parameters (rhs.args) at a greater indent
func (rhs FunctionCallRHS) aststring(indent string) string {
	var innerStats string
	nameStats := addIndAndNewLine(indent, rhs.ident)

	for _, param := range rhs.args {
		innerStats = fmt.Sprintf(
			"%v%v",
			innerStats,
			param.aststring(indent),
		)
	}

	return fmt.Sprintf("%v%v", nameStats, innerStats)
}

// Prints the LHS and RHS of an assignment. Format:
// - ASSIGNMENT
//   - LHS
//     - [lhsEXPR]
//   - RHS
//     - [rhsEXPR]
// lhsEXPR and rhsEXPR are recursed upon
func (stmt AssignStatement) aststring(indent string) string {
	declareStats := fmt.Sprintf("%vASSIGNMENT\n", addMinToIndent(indent))
	innerIndent := getGreaterIndent(indent)
	lhsIndent := addIndentForFirst(
		innerIndent,
		"LHS",
		stmt.target.aststring(getGreaterIndent(innerIndent)),
	)
	rhsIndent := addIndentForFirst(
		innerIndent,
		"RHS",
		stmt.rhs.aststring(getGreaterIndent(innerIndent)),
	)

	return fmt.Sprintf("%v%v%v", declareStats, lhsIndent, rhsIndent)
}

// Prints an IfStatement. Format:
// - IF
//   - CONDITION
//     - [bool]
//   - THEN
//     - [thenEXPR]
//   - ELSE
//     - [elseEXPR]
// thenEXPR and elseEXPR are recursed upon
func (stmt IfStatement) aststring(indent string) string {
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

	stmtStats = stmt.cond.aststring(doubleInnerIndent)

	st := stmt.trueStat
	for st.GetNext() != nil {
		trueTmp = st.aststring(doubleInnerIndent)
		trueStats = fmt.Sprintf("%v%v", trueStats, trueTmp)
		st = st.GetNext()
	}

	trueTmp = st.aststring(doubleInnerIndent)
	trueStats = fmt.Sprintf("%v%v", trueStats, trueTmp)

	st = stmt.falseStat
	for st.GetNext() != nil {
		falseTmp = st.aststring(doubleInnerIndent)
		falseStats = fmt.Sprintf("%v%v", falseStats, falseTmp)
		st = st.GetNext()
	}

	falseTmp = st.aststring(doubleInnerIndent)
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

// Prints a WhileLoop. Format:
// - LOOP
//   - CONDITION
//     - [bool]
//   - DO
//     - [doEXPR]
// doEXPR is recursed upon
func (stmt WhileStatement) aststring(indent string) string {
	var body string
	var doStats string
	innerIndent := getGreaterIndent(indent)

	doStats = addIndAndNewLine(innerIndent, "DO")

	condStats := addIndentForFirst(
		innerIndent,
		"CONDITION",
		stmt.cond.aststring(getGreaterIndent(innerIndent)),
	)

	st := stmt.body
	for st.GetNext() != nil {
		body = st.aststring(getGreaterIndent(innerIndent))
		doStats = fmt.Sprintf("%v%v", doStats, body)
		st = st.GetNext()
	}
	body = st.aststring(getGreaterIndent(innerIndent))
	doStats = fmt.Sprintf("%v%v", doStats, body)

	loopStats := addIndAndNewLine(indent, "LOOP")

	return fmt.Sprintf("%v%v%v", loopStats, condStats, doStats)
}

// Prints FunctionParameters in function declaration.
func (fp FunctionParam) aststring(indent string) string {
	return fmt.Sprintf("%v %v", fp.waccType, fp.name)
}

// Prints a FunctionDefinition. Format:
// - [type] [ident]([params])
func (fd FunctionDef) aststring(indent string) string {

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
		body = fmt.Sprintf("%v%v", body, st.aststring(innerIndent))
		st = st.GetNext()
	}

	body = fmt.Sprintf("%v%v", body, st.aststring(innerIndent))

	return fmt.Sprintf("%v%v", declaration, body)
}

// Prints all ArrayElements. Format:
// - []
//   -[elems]
// elems is recursed upon.
func (elem ArrayElem) aststring(indent string) string {
	var indexes string
	var tmpIndex string

	nextIndent := getGreaterIndent(indent)

	for _, index := range elem.indexes {
		tmpIndex =
			addIndentForFirst(nextIndent,
				"[]",
				index.aststring(getGreaterIndent(nextIndent)))
		indexes = fmt.Sprintf("%v%v", indexes, tmpIndex)
	}

	return addIndentForFirst(indent, elem.ident, indexes)
}

// Prints "[]" in of ArrayType
func (a ArrayType) aststring(indent string) string {
	typeStats := fmt.Sprintf("%v[]", a.base)
	return addType(indent, typeStats)
}

// Prints and invalid Type. Format:
// - TYPE
//   - <invalid>
func (i InvalidType) aststring(indent string) string {
	return addType(indent, "<invalid>")
}

// Prints and unknown Type. Format:
// - TYPE
//   - <unknown>
func (i UnknownType) aststring(indent string) string {
	return addType(indent, "<unknown>")
}

// Prints and int Type. Format:
// - TYPE
//   - int
func (i IntType) aststring(indent string) string {
	return addType(indent, "int")
}

// Prints and bool Type. Format:
// - TYPE
//   - bool
func (b BoolType) aststring(indent string) string {
	return addType(indent, "bool")
}

// Prints and char Type. Format:
// - TYPE
//   - char
func (c CharType) aststring(indent string) string {
	return addType(indent, "char")
}

// Prints a SKIP statement. Format:
// - SKIP
func (stmt SkipStatement) aststring(indent string) string {
	return addIndAndNewLine(indent, "SKIP")
}

// Prints a useless BlockStatement.
func (stmt BlockStatement) aststring(indent string) string {
	return fmt.Sprintf("")
}

// Prints a useless parenthesis.
func (par ExprParen) aststring(indent string) string {
	return ""
}

// Recurses oh the RHS of an Expression.
func (rhs ExpressionRHS) aststring(indent string) string {
	return rhs.expr.aststring(indent)
}

// Prints a READ statement. Format:
// - READ
//   - [args]
// Recurses on args.
func (stmt ReadStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"READ",
		stmt.target.aststring(getGreaterIndent(indent)),
	)
}

// Prints a FREE statement. Format:
// - FREE
//   - [args]
// Recurses on args.
func (stmt FreeStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"FREE",
		stmt.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints a RETURN statement. Format:
// - RETURN
//   - [args]
// Recurses on args.
func (ret ReturnStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"RETURN",
		ret.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints a EXIT statement. Format:
// - EXIT
//   - [args]
// Recurses on args.
func (stmt ExitStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"EXIT",
		stmt.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints a PRINTLN statement. Format:
// - PRINTLN
//   - [args]
// Recurses on args.
func (stmt PrintLnStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"PRINTLN",
		stmt.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints a PRINT statement. Format:
// - PRINT
//   - [args]
// Recurses on args.
func (stmt PrintStatement) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"PRINT",
		stmt.expr.aststring(getGreaterIndent(indent)),
	)
}

// Prints an identifier on a new line.
func (ident Ident) aststring(indent string) string {
	return addIndAndNewLine(indent, ident.ident)
}

// Prints a literal on a new line.
func (liter IntLiteral) aststring(indent string) string {
	return addIndAndNewLine(indent, strconv.Itoa(liter.value))
}

// Prints bool literal "true" on a new line.
func (liter BoolLiteralTrue) aststring(indent string) string {
	return addIndAndNewLine(indent, "true")
}

// Prints bool literal "false" on a new line.
func (liter BoolLiteralFalse) aststring(indent string) string {
	return addIndAndNewLine(indent, "false")
}

// Prints a char literal on a new line.
func (liter CharLiteral) aststring(indent string) string {
	tmpStats := fmt.Sprintf("'%v'", liter.char)
	return addIndAndNewLine(indent, tmpStats)
}

// Prints a string literal on a new line.
func (liter StringLiteral) aststring(indent string) string {
	tmp := fmt.Sprintf("\"%v\"", liter.str)
	return addIndAndNewLine(indent, tmp)
}

// Prints a pair on a new line.
func (liter PairLiteral) aststring(indent string) string {
	return fmt.Sprintf("pair(%v, %v)", liter.fst, liter.snd)
}

// Prints a null Pair on a new line.
func (liter NullPair) aststring(indent string) string {
	return addIndAndNewLine(indent, "null")
}

// Prints a ! unaryOperator. Format:
// - !
//   - [args]
// Recurses on args.
func (op UnaryOperatorNot) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"!",
		op.GetExpression().aststring(getGreaterIndent(indent)),
	)
}

// Prints a - unaryOperator. Format:
// - -
//   - [args]
// Recurses on args.
func (op UnaryOperatorNegate) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"-",
		op.GetExpression().aststring(getGreaterIndent(indent)),
	)
}

// Prints a len unaryOperator. Format:
// - len
//   - [args]
// Recurses on args.
func (op UnaryOperatorLen) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"len",
		op.GetExpression().aststring(getGreaterIndent(indent)),
	)
}

// Prints a ord unaryOperator. Format:
// - ord
//   - [args]
// Recurses on args.
func (op UnaryOperatorOrd) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"ord",
		op.GetExpression().aststring(getGreaterIndent(indent)),
	)
}

// Prints a chr unaryOperator. Format:
// - chr
//   - [args]
// Recurses on args.
func (op UnaryOperatorChr) aststring(indent string) string {
	return addIndentForFirst(
		indent,
		"chr",
		op.GetExpression().aststring(getGreaterIndent(indent)),
	)
}

// Prints a * binaryOperator. Format:
// - *
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorMult) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"*",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a / binaryOperator. Format:
// - /
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorDiv) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"/",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a % binaryOperator. Format:
// - %
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorMod) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"%",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a + binaryOperator. Format:
// - +
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorAdd) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"+",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a - binaryOperator. Format:
// - -
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorSub) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"-",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a > binaryOperator. Format:
// - >
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorGreaterThan) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		">",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a >= binaryOperator. Format:
// - >=
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorGreaterEqual) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		">=",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a < binaryOperator. Format:
// - <
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorLessThan) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"<",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a <= binaryOperator. Format:
// - <=
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorLessEqual) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"<=",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a == binaryOperator. Format:
// - ==
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorEqual) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"==",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a != binaryOperator. Format:
// - !=
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorNotEqual) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"!=",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a && binaryOperator. Format:
// - &&
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorAnd) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"&&",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

// Prints a || binaryOperator. Format:
// - ||
//   - [arg1]
//   - [arg2]
// Recurses on arg1 and arg2.
func (op BinaryOperatorOr) aststring(indent string) string {
	return addTripleIndentOnlyFst(
		indent,
		"||",
		op.GetLHS().aststring(getGreaterIndent(indent)),
		op.GetRHS().aststring(getGreaterIndent(indent)),
	)
}

//------------------------------------------------------------------------------

// Main method. Format:
// - [functions]
// - int main()
//   - [main]
// Recurses on functions and main
func (ast AST) aststring() string {
	var tree string
	var tmpIndent string

	tree = addIndAndNewLine("", "Program")

	for _, function := range ast.functions {
		tree = fmt.Sprintf(
			"%v%v",
			tree,
			function.aststring(basicIndent),
		)
	}

	tmpIndent = getGreaterIndent(basicIndent)

	tree = fmt.Sprintf("%v%v",
		tree,
		addIndAndNewLine(basicIndent, "int main()"),
	)

	stmt := ast.main
	for stmt.GetNext() != nil {
		tree = fmt.Sprintf("%v%v", tree, stmt.aststring(tmpIndent))
		stmt = stmt.GetNext()
	}
	tree = fmt.Sprintf("%v%v", tree, stmt.aststring(tmpIndent))

	return tree
}
