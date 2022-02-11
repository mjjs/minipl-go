package symboltable

import (
	"fmt"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

type SymbolTableCreator struct {
	symbols       *SymbolTable
	lockedSymbols map[string]struct{}
}

func (stc *SymbolTableCreator) Create(root ast.Node) *SymbolTable {
	stc.symbols = NewSymbolTable()
	stc.lockedSymbols = make(map[string]struct{})

	root.Accept(stc)

	return stc.symbols
}

func (stc *SymbolTableCreator) VisitProg(node ast.Prog) {
	node.Statements.Accept(stc)
}

func (stc *SymbolTableCreator) VisitStmts(node ast.Stmts) {
	for _, stmt := range node.Statements {
		stmt.Accept(stc)
	}
}

func (stc *SymbolTableCreator) VisitDeclStmt(node ast.DeclStmt) {
	name := node.Identifier.Value()
	_, exists := stc.symbols.Get(name)
	if exists {
		panic(fmt.Errorf("Variable %s has already been declared", name))
	}

	switch node.VariableType.Type() {
	case token.INTEGER:
		stc.symbols.Insert(name, INTEGER)
	case token.STRING:
		stc.symbols.Insert(name, STRING)
	case token.BOOLEAN:
		stc.symbols.Insert(name, BOOLEAN)
	}
}

func (stc *SymbolTableCreator) VisitAssignStmt(node ast.AssignStmt) {
	node.Identifier.Accept(stc)

	_, locked := stc.lockedSymbols[node.Identifier.Id.Value()]

	if locked {
		panic(fmt.Errorf(
			"Cannot assign to a locked variable %s",
			node.Identifier.Id.Value(),
		))
	}
}

func (stc *SymbolTableCreator) VisitForStmt(node ast.ForStmt) {
	node.Index.Accept(stc)
	stc.lockedSymbols[node.Index.Id.Value()] = struct{}{}

	node.Low.Accept(stc)
	node.High.Accept(stc)

	node.Statements.Accept(stc)

	delete(stc.lockedSymbols, node.Index.Id.Value())
}

func (stc *SymbolTableCreator) VisitReadStmt(node ast.ReadStmt) {
	node.TargetIdentifier.Accept(stc)
}

func (stc *SymbolTableCreator) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(stc)
}

func (stc *SymbolTableCreator) VisitAssertStmt(node ast.AssertStmt) {
	node.Expression.Accept(stc)
}

func (stc *SymbolTableCreator) VisitBinaryExpr(node ast.BinaryExpr) {
	node.Left.Accept(stc)
	node.Right.Accept(stc)
}

func (stc *SymbolTableCreator) VisitUnaryExpr(node ast.UnaryExpr) {
	node.Operand.Accept(stc)
}

func (stc *SymbolTableCreator) VisitNullaryExpr(node ast.NullaryExpr) {
	node.Operand.Accept(stc)
}

func (stc *SymbolTableCreator) VisitNumberOpnd(node ast.NumberOpnd) {
	// Nothing to do
}

func (stc *SymbolTableCreator) VisitStringOpnd(node ast.StringOpnd) {
	// Nothing to do
}

func (stc *SymbolTableCreator) VisitIdent(node ast.Ident) {
	name := node.Id.Value()
	_, exists := stc.symbols.Get(name)
	if !exists {
		panic(fmt.Errorf("Variable %s has not been declared", name))
	}
}
