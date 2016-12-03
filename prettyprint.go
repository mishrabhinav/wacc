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

// Util Functions:
//------------------------------------------------------------------------------

const startingIndent int = 1
const basicIndent string = "  "

// Given an indentationLevel (int),
// Returns a string corresponding to the given indent
func getIndentation(level int) string {
	return fmt.Sprint(strings.Repeat(basicIndent, level))
}

// Given an expr (Expresion) and operator (string),
// Returns a string with the given operator applied INLINE with the given expr.
// Format:
// [op][expr]
func generateUnaryOperator(expr Expression, op string) string {
	return fmt.Sprintf("%v%v", op, expr)
}

// Given an two expressions lhs and rhs (Expresion) and operator (string),
// Returns a string with the given operator applied INLINE with the given exprs.
// Format:
// [lhs][op][rhs]
func generateBinaryOperator(lhs Expression, rhs Expression, op string) string {
	return fmt.Sprintf("%v %v %v", lhs, op, rhs)
}

//------------------------------------------------------------------------------

// Prints invalid Types. Format:
//   "<invalid>"
func (m InvalidType) String() string {
	return "<invalid>"
}

// Prints unknown Types. Format:
//   "<unknown>"
func (m UnknownType) String() string {
	return "<unknown>"
}

// Prints integer Types. Format:
//   "int"
func (i IntType) String() string {
	return "int"
}

// Prints boolean Types. Format:
//   "bool"
func (b BoolType) String() string {
	return "bool"
}

// Prints char Types. Format:
//   "char"
func (c CharType) String() string {
	return "char"
}

// Prints pair Types. Format:
//   "pair([fst], [snd])"
// Recurses on fst and snd.
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

// Prints array Types. Format:
//   "[arr][]"
// Recurses on arr.
func (a ArrayType) String() string {
	return fmt.Sprintf("%v[]", a.base)
}

// Prints the file includes. Format:
//   "include <filename.wacc>"
func (incl *Include) String() string {
	return fmt.Sprintf("include \"%v\"", incl.file)
}

// Prints identifier Types. Format:
//   "[ident]"
// Recurses on ident.
func (ident *Ident) String() string {
	return ident.ident
}

// Prints integer Literals. Format:
//   "[int]"
// Recurses on int.
func (liter *IntLiteral) String() string {
	return fmt.Sprint(liter.value)
}

// Prints true bool Literal. Format:
//   "true"
func (liter *BoolLiteralTrue) String() string {
	return "true"
}

// Prints false bool Literal. Format:
//   "false"
func (liter *BoolLiteralFalse) String() string {
	return "false"
}

// Prints character Literals. Format:
//   "\'[char]\'"
// Recurses on char.
func (liter *CharLiteral) String() string {
	return fmt.Sprintf("'%v'", liter.char)
}

// Prints string Literals. Format:
//   "\"[str]\""
// Recurses on str.
func (liter *StringLiteral) String() string {
	return fmt.Sprintf("\"%v\"", liter.str)
}

// Prints pair Literals. Format:
//   "pair([fst], [snd])"
// Recurses on fst and snd.
func (liter *PairLiteral) String() string {
	return fmt.Sprintf("pair(%v, %v)", liter.fst, liter.snd)
}

// Prints a null Pair. Format:
//   "null"
func (liter *NullPair) String() string {
	return "null"
}

// Prints array Elements. Format:
//   "[arr][[elems]]"
// Recurses on arr and elems*.
func (elem *ArrayElem) String() string {
	var indexes string

	for _, index := range elem.indexes {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", elem.ident, indexes)
}

// Prints ! unaryOperator. Format:
//   "![expr]"
// Recurses on expr.
func (op *UnaryOperatorNot) String() string {
	return generateUnaryOperator(op.GetExpression(), "!")
}

// Prints - unaryOperator. Format:
//   "-[expr]"
// Recurses on expr.
func (op *UnaryOperatorNegate) String() string {
	return generateUnaryOperator(op.GetExpression(), "-")
}

// Prints len unaryOperator. Format:
//   "len [expr]"
// Recurses on expr.
func (op *UnaryOperatorLen) String() string {
	return generateUnaryOperator(op.GetExpression(), "len ")
}

// Prints ord unaryOperator. Format:
//   "ord [expr]"
// Recurses on expr.
func (op *UnaryOperatorOrd) String() string {
	return generateUnaryOperator(op.GetExpression(), "ord ")
}

// Prints chr unaryOperator. Format:
//   "chr [expr]"
// Recurses on expr.
func (op *UnaryOperatorChr) String() string {
	return generateUnaryOperator(op.GetExpression(), "chr ")
}

// Prints * unaryOperator. Format:
//   "*[expr]"
// Recurses on expr.
func (op *BinaryOperatorMult) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "*")
}

// Prints / unaryOperator. Format:
//   "/[expr]"
// Recurses on expr.
func (op *BinaryOperatorDiv) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "/")
}

// Prints % unaryOperator. Format:
//   "%[expr]"
// Recurses on expr.
func (op *BinaryOperatorMod) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "%")
}

// Prints + unaryOperator. Format:
//   "+[expr]"
// Recurses on expr.
func (op *BinaryOperatorAdd) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "+")
}

// Prints - unaryOperator. Format:
//   "-[expr]"
// Recurses on expr.
func (op *BinaryOperatorSub) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "-")
}

// Prints > unaryOperator. Format:
//   ">[expr]"
// Recurses on expr.
func (op *BinaryOperatorGreaterThan) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), ">")
}

// Prints >= unaryOperator. Format:
//   ">=[expr]"
// Recurses on expr.
func (op *BinaryOperatorGreaterEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), ">=")
}

// Prints < unaryOperator. Format:
//   "<[expr]"
// Recurses on expr.
func (op *BinaryOperatorLessThan) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "<")
}

// Prints <= unaryOperator. Format:
//   "<=[expr]"
// Recurses on expr.
func (op *BinaryOperatorLessEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "<=")
}

// Prints == unaryOperator. Format:
//   "==[expr]"
// Recurses on expr.
func (op *BinaryOperatorEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "==")
}

// Prints != unaryOperator. Format:
//   "!=[expr]"
// Recurses on expr.
func (op *BinaryOperatorNotEqual) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "!=")
}

// Prints && unaryOperator. Format:
//   "&&[expr]"
// Recurses on expr.
func (op *BinaryOperatorAnd) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "&&")
}

// Prints || unaryOperator. Format:
//   "||[expr]"
// Recurses on expr.
func (op *BinaryOperatorOr) String() string {
	return generateBinaryOperator(op.GetLHS(), op.GetRHS(), "||")
}

// Prints the lhs of a PairElem.
func (lhs *PairElemLHS) String() string {
	if lhs.snd {
		return fmt.Sprintf("snd %v", lhs.expr)
	}

	return fmt.Sprintf("fst %v", lhs.expr)
}

// Prints array Elements. Format:
//   "[arr][[elems]]"
// Recurses on arr and elems*.
func (lhs *ArrayLHS) String() string {
	var indexes string

	for _, index := range lhs.index {
		indexes = fmt.Sprintf("%v[%v]", indexes, index)
	}

	return fmt.Sprintf("%v%v", lhs.ident, indexes)
}

// Prints the lhs of a Variable. Format:
//   "[var]"
// Recurses on var.
func (lhs *VarLHS) String() string {
	return fmt.Sprintf(lhs.ident)
}

// Prints a new pairLiteral. Format:
//   "newpair([fst], [snd])"
// Recurses on fst and snd.
func (rhs *PairLiterRHS) String() string {
	return fmt.Sprintf("newpair(%v, %v)", rhs.fst, rhs.snd)
}

// Prints an array Elements. Format:
//   "[[elem1](, [elems])*]"
// Recurses on elem1 with optional elems.
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

// Prints the rhs of a PairElem.
func (rhs *PairElemRHS) String() string {
	if rhs.snd {
		return fmt.Sprintf("snd %v", rhs.expr)
	}

	return fmt.Sprintf("fst %v", rhs.expr)
}

// Prints a new functionCall. Format:
//   "call [fun]([args]*)"
// Recurses fun and optional args.
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

// Prints a new expression. Format:
//   "[expr]"
// Recurses on expr.
func (rhs *ExpressionRHS) String() string {
	return fmt.Sprintf("%v", rhs.expr)
}

// Prints a skip statement. Format:
//   "skip"
func (stmt *SkipStatement) istring(level int) string {
	return fmt.Sprintf("%vskip", getIndentation(level))
}

// Prints a useless block statement. Format:
//   ""
func (stmt *BlockStatement) istring(level int) string {
	return ""
}

// Prints a declaration assignment. Format:
//   "[type] [ident]=[rhs]"
// Recurses on type, ident and rhs.
func (stmt *DeclareAssignStatement) istring(level int) string {
	return fmt.Sprintf("%v%v %v = %v", getIndentation(level), stmt.waccType,
		stmt.ident, stmt.rhs)
}

// Prints an assignment. Format:
//   "[ident]=[rhs]"
// Recurses on ident and rhs.
func (stmt *AssignStatement) istring(level int) string {
	return fmt.Sprintf("%v%v = %v", getIndentation(level), stmt.target,
		stmt.rhs)
}

// Prints a read statement. Format:
//   "read"
func (stmt *ReadStatement) istring(level int) string {
	return fmt.Sprintf("%vread %v", getIndentation(level), stmt.target)
}

// Prints a free statement. Format:
//   "free"
func (stmt *FreeStatement) istring(level int) string {
	return fmt.Sprintf("%vfree %v", getIndentation(level), stmt.expr)
}

// Prints a return statement. Format:
//   "return"
func (ret *ReturnStatement) istring(level int) string {
	return fmt.Sprintf("%vreturn %v", getIndentation(level), ret.expr)
}

// Prints an exit statement. Format:
//   "exit"
func (stmt *ExitStatement) istring(level int) string {
	return fmt.Sprintf("%vexit %v", getIndentation(level), stmt.expr)
}

// Prints a println statement. Format:
//   "println"
func (stmt *PrintLnStatement) istring(level int) string {
	return fmt.Sprintf("%vprintln %v", getIndentation(level), stmt.expr)
}

// Prints a print statement. Format:
//   "print"
func (stmt *PrintStatement) istring(level int) string {
	return fmt.Sprintf("%vprint %v", getIndentation(level), stmt.expr)
}

// Prints an if statement. Format:
//   "if [cond]
//    then [trueStat]*
//    else [falseStat]*
//    fi"
// Recurses on cond, (multiple) trueStat and (multiple) falseStat.
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

// Prints a while loop. Format:
//   "while ([cond]) do
//    [body]*
//    done"
// Recurses on cond and (multiple) body.
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

// Prints a given function parameter. Format:
//   "[type] [name]"
// Recurses on type and name.
func (fp *FunctionParam) String() string {
	return fmt.Sprintf("%v %v", fp.waccType, fp.name)
}

// Prints a function definition. Format:
//   "[type] [name]([args]*) is
//    [body] (;\n [bodies])*
//    end"
// Recurses on type, name, (multiple) args, body and (multpiple/optional) bodies
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

// Prints the AST. Format:
//   "begin
//    ([functions])*
//    [body] (;\n [bodies])*
//    end"
// Recurses on (multpiple/optional) functions,
// body and (multpiple/optional) bodies.
func (ast *AST) String() string {
	var tree string

	tree = fmt.Sprintf("begin")

	for _, include := range ast.includes {
		tree = fmt.Sprintf("%v\n  %v\n", tree, include)
	}

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
