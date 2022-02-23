package ast

import "github.com/mjjs/minipl-go/pkg/token"

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
type Node interface {
	Accept(Visitor)
	Position() token.Position
}

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
type Prog struct {
	Statements Stmts
}

func (p Prog) Position() token.Position {
	return p.Statements.Position()
}

// Stmts is an abstract syntax tree node containing all the statements of the program.
type Stmts struct{ Statements []Stmt }

func (s Stmts) Position() token.Position {
	return s.Statements[0].Position()
}

// ReadStmt is an abstract syntax tree node which defines the read statement.
// The destination of the read is defined in TargetIdentifier.
type ReadStmt struct {
	TargetIdentifier Ident
	Pos              token.Position
}

func (r ReadStmt) Position() token.Position { return r.Pos }

// PrintStmt is an abstract syntax tree node which defines the print statement.
type PrintStmt struct {
	Expression Expr
	Pos        token.Position
}

func (p PrintStmt) Position() token.Position { return p.Pos }

// AssertStmt is an abstract syntax tree node which defines the assert statement.
type AssertStmt struct {
	Expression Expr
	Pos        token.Position
}

func (a AssertStmt) Position() token.Position { return a.Pos }

// ForStmt defines a for loop statement over a range of values ranging from
// Low to High. Index is the control variable used in the statement.
type ForStmt struct {
	Index      Ident
	Low        Expr
	High       Expr
	Statements Stmts
	Pos        token.Position
}

func (f ForStmt) Position() token.Position { return f.Pos }

// AssignStmt defines a statement node.
type AssignStmt struct {
	Identifier Ident
	Expression Expr
	Pos        token.Position
}

func (a AssignStmt) Position() token.Position { return a.Pos }

// DeclStmt defines a declaration of a new variable.
type DeclStmt struct {
	Identifier   token.Token
	VariableType token.Token
	// nil when no value is assigned to the variable during declaration
	Expression Expr
	Pos        token.Position
}

func (d DeclStmt) Position() token.Position { return d.Pos }

// BinaryExpr is an expression with two operands and and operator.
type BinaryExpr struct {
	Left     Node
	Operator token.Token
	Right    Node
}

func (b BinaryExpr) Position() token.Position { return b.Left.Position() }

// UnaryExpr is an expression with one operand preceded by an unary operator.
type UnaryExpr struct {
	Unary   token.Token
	Operand Node
	Pos     token.Position
}

func (u UnaryExpr) Position() token.Position { return u.Pos }

// NullaryExpr is an expression with a single operand and no operators.
type NullaryExpr struct {
	Operand Node
}

func (n NullaryExpr) Position() token.Position { return n.Operand.Position() }

// NumberOpnd is an integer operand.
type NumberOpnd struct {
	Value int
	Pos   token.Position
}

func (n NumberOpnd) Position() token.Position { return n.Pos }

// StringOpnd is a string operand.
type StringOpnd struct {
	Value string
	Pos   token.Position
}

func (s StringOpnd) Position() token.Position { return s.Pos }

// Ident is an identifier node.
type Ident struct {
	Id  token.Token
	Pos token.Position
}

func (i Ident) Position() token.Position { return i.Pos }

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
