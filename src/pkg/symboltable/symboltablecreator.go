package symboltable

import (
	"fmt"

	"github.com/mjjs/minipl-go/pkg/ast"
	"github.com/mjjs/minipl-go/pkg/token"
)

type SymbolTableCreator struct {
	symbols       *SymbolTable
	lockedSymbols map[string]struct{}

	errors []error
}

func (stc *SymbolTableCreator) Create(root ast.Node) (*SymbolTable, []error) {
	stc.symbols = NewSymbolTable()
	stc.lockedSymbols = make(map[string]struct{})

	root.Accept(stc)

	return stc.symbols, stc.errors
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
		err := fmt.Errorf("%s: redeclaration of variable %s", node.Position(), name)
		stc.errors = append(stc.errors, err)
		return
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
		err := fmt.Errorf(
			"%s: cannot modify loop index %s during loop",
			node.Position(),
			node.Identifier.Id.Value(),
		)

		stc.errors = append(stc.errors, err)
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
		err := fmt.Errorf("%s: variable %s used before declaration", node.Position(), name)
		stc.errors = append(stc.errors, err)
	}
}
