package ast

import "github.com/mjjs/minipl-go/src/lexer"

// Visitor defines an interface for a visitor to the abstract syntax tree.
// This interface is used for defining separate compiler passes that need to
// go through the abstract syntax tree using the visitor pattern.
type Visitor interface {
	VisitProg(Prog)
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

// Node is a basic node for the abstract syntax tree.
type Node interface{ Accept(Visitor) }

// Expr defines all the expression nodes.
type Expr interface {
	Node
	exprNode()
}

// Stmt defines all the expression nodes.
type Stmt interface {
	Node
	stmtNode()
}

// Prog is the root node of the abstract syntax tree.
type Prog struct{ Statements Stmts }

// Stmts is an abstract syntax tree node containing all the statements of the program.
type Stmts struct{ Statements []Stmt }

// ReadStmt is an abstract syntax tree node which defines the read statement.
// The destination of the read is defined in TargetIdentifier.
type ReadStmt struct{ TargetIdentifier Ident }

// PrintStmt is an abstract syntax tree node which defines the print statement.
type PrintStmt struct{ Expression Expr }

// AssertStmt is an abstract syntax tree node which defines the assert statement.
type AssertStmt struct{ Expression Expr }

// ForStmt defines a for loop statement over a range of values ranging from
// Low to High. Index is the control variable used in the statement.
type ForStmt struct {
	Index      Ident
	Low        Expr
	High       Expr
	Statements Stmts
}

// AssignStmt defines a statement node.
type AssignStmt struct {
	Identifier Ident
	Expression Expr
}

// DeclStmt defines a declaration of a new variable.
type DeclStmt struct {
	Identifier   lexer.Token
	VariableType lexer.Token
	// nil when no value is assigned to the variable during declaration
	Expression Expr
}

// BinaryExpr is an expression with two operands and and operator.
type BinaryExpr struct {
	Left     Node
	Operator lexer.Token
	Right    Node
}

// UnaryExpr is an expression with one operand preceded by an unary operator.
type UnaryExpr struct {
	Unary   lexer.Token
	Operand Node
}

// NullaryExpr is an expression with a single operand and no operators.
type NullaryExpr struct {
	Operand Node
}

// NumberOpnd is an integer operand.
type NumberOpnd struct{ Value int }

// StringOpnd is a string operand.
type StringOpnd struct{ Value string }

// Ident is an identifier node.
type Ident struct{ Id lexer.Token }

func (n Prog) Accept(v Visitor)        { v.VisitProg(n) }
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
