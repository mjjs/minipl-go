package typechecker

import (
	"fmt"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/stack"
	"github.com/mjjs/minipl-go/symboltable"
	"github.com/mjjs/minipl-go/token"
)

type TypeChecker struct {
	stack   *stack.Stack
	symbols *symboltable.SymbolTable

	errors []error
}

func New(symbols *symboltable.SymbolTable) *TypeChecker {
	return &TypeChecker{
		stack:   stack.New(),
		symbols: symbols,
	}
}

func (tc *TypeChecker) CheckTypes(root ast.Node) []error {
	root.Accept(tc)
	return tc.errors
}

func (tc *TypeChecker) VisitProg(node ast.Prog) {
	node.Statements.Accept(tc)
}

func (tc *TypeChecker) VisitStmts(node ast.Stmts) {
	for _, stmt := range node.Statements {
		stmt.Accept(tc)
	}
}

func (tc *TypeChecker) VisitDeclStmt(node ast.DeclStmt) {
	if node.Expression == nil {
		return
	}

	var variableType symboltable.SymbolType

	switch node.VariableType.Type() {
	case token.INTEGER:
		variableType = symboltable.INTEGER
	case token.STRING:
		variableType = symboltable.STRING
	case token.BOOLEAN:
		variableType = symboltable.BOOLEAN
	}

	node.Expression.Accept(tc)

	rhsType := tc.stack.Pop().(symboltable.SymbolType)
	if variableType != rhsType {
		err := fmt.Errorf(
			"%s: cannot assign type %s to variable %s of type %s",
			node.Position(), rhsType, node.Identifier.Value(), variableType,
		)

		tc.errors = append(tc.errors, err)
	}
}

func (tc *TypeChecker) VisitAssignStmt(node ast.AssignStmt) {
	node.Identifier.Accept(tc)
	idType := tc.stack.Pop().(symboltable.SymbolType)

	node.Expression.Accept(tc)
	exprType := tc.stack.Pop().(symboltable.SymbolType)

	if exprType != idType {
		err := fmt.Errorf(
			"%s: cannot assign type %s to variable %s of type %s",
			node.Position(), exprType, node.Identifier.Id.Value(), idType,
		)

		tc.errors = append(tc.errors, err)
	}
}

func (tc *TypeChecker) VisitForStmt(node ast.ForStmt) {
	node.Index.Accept(tc)
	indexType := tc.stack.Pop().(symboltable.SymbolType)

	node.Low.Accept(tc)
	lowType := tc.stack.Pop().(symboltable.SymbolType)

	node.High.Accept(tc)
	highType := tc.stack.Pop().(symboltable.SymbolType)

	if indexType != symboltable.INTEGER {
		err := fmt.Errorf(
			"%s: loop index must be %s, not %s",
			node.Position(), symboltable.INTEGER, indexType,
		)

		tc.errors = append(tc.errors, err)
	}

	if lowType != symboltable.INTEGER {
		err := fmt.Errorf(
			"%s: for loop range lower bound must be %s, not %s",
			node.Position(), symboltable.INTEGER, lowType,
		)

		tc.errors = append(tc.errors, err)
	}

	if highType != symboltable.INTEGER {
		err := fmt.Errorf(
			"%s: for loop range upper bound must be %s, not %s",
			node.Position(), symboltable.INTEGER, highType,
		)

		tc.errors = append(tc.errors, err)
	}
}

func (tc *TypeChecker) VisitReadStmt(node ast.ReadStmt) {
}

func (tc *TypeChecker) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(tc)
}

func (tc *TypeChecker) VisitAssertStmt(node ast.AssertStmt) {
	node.Expression.Accept(tc)

	exprType := tc.stack.Pop().(symboltable.SymbolType)
	if exprType != symboltable.BOOLEAN {
		err := fmt.Errorf(
			"%s: assert statement is only defined for type %s, not %s",
			node.Position(), symboltable.BOOLEAN, exprType,
		)

		tc.errors = append(tc.errors, err)
	}
}

func (tc *TypeChecker) VisitBinaryExpr(node ast.BinaryExpr) {
	node.Left.Accept(tc)
	left := tc.stack.Pop().(symboltable.SymbolType)

	node.Right.Accept(tc)
	right := tc.stack.Pop().(symboltable.SymbolType)

	if left != right {
		err := fmt.Errorf(
			"%s: unmatched types %s and %s for binary expression %s",
			node.Position(), left, right, node.Operator.Type(),
		)

		tc.errors = append(tc.errors, err)
	}

	switch node.Operator.Type() {
	case token.PLUS:
		if left != symboltable.INTEGER && left != symboltable.STRING {
			err := fmt.Errorf(
				"%s: operator %s not defined for type %s",
				node.Position(), token.PLUS, left,
			)

			tc.errors = append(tc.errors, err)
		}
		tc.stack.Push(left)

	case token.MINUS:
		if left != symboltable.INTEGER {
			err := fmt.Errorf(
				"%s: operator %s not defined for type %s",
				node.Position(), token.MINUS, left,
			)

			tc.errors = append(tc.errors, err)
		}
		tc.stack.Push(left)

	case token.MULTIPLY:
		if left != symboltable.INTEGER {
			err := fmt.Errorf(
				"%s: operator %s not defined for type %s",
				node.Position(), token.MULTIPLY, left,
			)

			tc.errors = append(tc.errors, err)
		}
		tc.stack.Push(left)

	case token.INTEGER_DIV:
		if left != symboltable.INTEGER {
			err := fmt.Errorf(
				"%s: operator %s not defined for type %s",
				node.Position(), token.INTEGER_DIV, left,
			)

			tc.errors = append(tc.errors, err)
		}
		tc.stack.Push(left)

	case token.AND:
		if left != symboltable.BOOLEAN {
			err := fmt.Errorf(
				"%s: operator %s not defined for type %s",
				node.Position(), token.AND, left,
			)

			tc.errors = append(tc.errors, err)
		}

		tc.stack.Push(left)

	case token.LT:
		tc.stack.Push(symboltable.BOOLEAN)

	case token.EQ:
		tc.stack.Push(symboltable.BOOLEAN)
	}
}

func (tc *TypeChecker) VisitUnaryExpr(node ast.UnaryExpr) {
	node.Operand.Accept(tc)
	t := tc.stack.Pop().(symboltable.SymbolType)

	if node.Unary.Type() == token.NOT && t != symboltable.BOOLEAN {
		err := fmt.Errorf(
			"%s: unary operator %s not defined for type %s",
			node.Position(), token.NOT, t,
		)

		tc.errors = append(tc.errors, err)
	}

	tc.stack.Push(t)
}

func (tc *TypeChecker) VisitNullaryExpr(node ast.NullaryExpr) {
	node.Operand.Accept(tc)
}

func (tc *TypeChecker) VisitNumberOpnd(node ast.NumberOpnd) {
	tc.stack.Push(symboltable.INTEGER)
}

func (tc *TypeChecker) VisitStringOpnd(node ast.StringOpnd) {
	tc.stack.Push(symboltable.STRING)
}

func (tc *TypeChecker) VisitIdent(node ast.Ident) {
	symbol, ok := tc.symbols.Get(node.Id.Value())
	if !ok {
		panic("Type checker came across a symbol not in the symbol table")
	}

	tc.stack.Push(symbol.Type())
}
