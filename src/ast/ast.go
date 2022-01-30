package ast

import "github.com/mjjs/minipl-go/src/lexer"

type Visitor interface {
	VisitStmts(Stmts)

	VisitAssignStmt(AssignStmt)
	VisitDeclStmt(DeclStmt)
	VisitForStmt(ForStmt)
	VisitReadStmt(ReadStmt)
	VisitPrintStmt(PrintStmt)
	VisitAssertStmt(AssertStmt)

	VisitBinaryExpr(BinaryExpr)
	VisitUnaryExpr(UnaryExpr)
	VisitNullaryExpr(NullaryExpr)

	VisitNumberOpnd(NumberOpnd)
	VisitStringOpnd(StringOpnd)

	VisitIdent(Ident)
}

type Node interface{ Accept(Visitor) }
type Expr interface {
	Node
	exprNode()
}
type Stmt interface {
	Node
	stmtNode()
}

type Stmts struct{ Statements []Stmt }
type ReadStmt struct{ TargetIdentifier lexer.Token }
type PrintStmt struct{ Expression Expr }
type AssertStmt struct{ Expression Expr }
type ForStmt struct {
	Index      lexer.Token
	Low        Expr
	High       Expr
	Statements Stmts
}
type AssignStmt struct {
	Identifier lexer.Token
	Expression Expr
}
type DeclStmt struct {
	Identifier   lexer.Token
	VariableType lexer.Token
	// nil when no value is assigned to the variable during declaration
	Expression Expr
}

type BinaryExpr struct {
	Left     Node
	Operator lexer.Token
	Right    Node
}
type UnaryExpr struct {
	Unary   lexer.Token
	Operand Node
}
type NullaryExpr struct {
	Operand Node
}

type NumberOpnd struct{ Value int }
type StringOpnd struct{ Value string }

type Ident struct{ Id lexer.Token }

func (n Stmts) Accept(v Visitor)       { v.VisitStmts(n) }
func (n ForStmt) Accept(v Visitor)     { v.VisitForStmt(n) }
func (n NumberOpnd) Accept(v Visitor)  { v.VisitNumberOpnd(n) }
func (n StringOpnd) Accept(v Visitor)  { v.VisitStringOpnd(n) }
func (n Ident) Accept(v Visitor)       { v.VisitIdent(n) }
func (n BinaryExpr) Accept(v Visitor)  { v.VisitBinaryExpr(n) }
func (n UnaryExpr) Accept(v Visitor)   { v.VisitUnaryExpr(n) }
func (n NullaryExpr) Accept(v Visitor) { v.VisitNullaryExpr(n) }
func (n AssignStmt) Accept(v Visitor)  { v.VisitAssignStmt(n) }
func (n ReadStmt) Accept(v Visitor)    { v.VisitReadStmt(n) }
func (n PrintStmt) Accept(v Visitor)   { v.VisitPrintStmt(n) }
func (n AssertStmt) Accept(v Visitor)  { v.VisitAssertStmt(n) }
func (n DeclStmt) Accept(v Visitor)    { v.VisitDeclStmt(n) }

func (n BinaryExpr) exprNode()  {}
func (n UnaryExpr) exprNode()   {}
func (n NullaryExpr) exprNode() {}
func (n NumberOpnd) exprNode()  {}
func (n StringOpnd) exprNode()  {}

func (n ForStmt) stmtNode()    {}
func (n PrintStmt) stmtNode()  {}
func (n ReadStmt) stmtNode()   {}
func (n AssertStmt) stmtNode() {}
func (n AssignStmt) stmtNode() {}
func (n DeclStmt) stmtNode()   {}
