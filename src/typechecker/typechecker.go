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
}

func New(symbols *symboltable.SymbolTable) *TypeChecker {
	return &TypeChecker{
		stack:   stack.New(),
		symbols: symbols,
	}
}

func (tc *TypeChecker) CheckTypes(root ast.Node) {
	root.Accept(tc)
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

	if variableType != tc.stack.Pop().(symboltable.SymbolType) {
		panic("NO MATCH")
	}
}

func (tc *TypeChecker) VisitAssignStmt(node ast.AssignStmt) {
	node.Identifier.Accept(tc)
	idType := tc.stack.Pop().(symboltable.SymbolType)

	node.Expression.Accept(tc)
	exprType := tc.stack.Pop().(symboltable.SymbolType)

	if exprType != idType {
		panic("Types do not match!")
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
		panic("INDEX MUST BE INT")
	}
	if lowType != symboltable.INTEGER {
		panic("LOW MUST BE INT")
	}
	if highType != symboltable.INTEGER {
		panic("HIGH MUST BE INT")
	}
}

func (tc *TypeChecker) VisitReadStmt(node ast.ReadStmt) {
}

func (tc *TypeChecker) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(tc)
}

func (tc *TypeChecker) VisitAssertStmt(node ast.AssertStmt) {
	node.Expression.Accept(tc)

	if tc.stack.Pop().(symboltable.SymbolType) != symboltable.BOOLEAN {
		panic("Not a bool")
	}
}

func (tc *TypeChecker) VisitBinaryExpr(node ast.BinaryExpr) {
	node.Left.Accept(tc)
	left := tc.stack.Pop().(symboltable.SymbolType)

	node.Right.Accept(tc)
	right := tc.stack.Pop().(symboltable.SymbolType)

	if left != right {
		panic("NO MATCH!")
	}

	switch node.Operator.Type() {
	case token.PLUS:
		if left != symboltable.INTEGER && left != symboltable.STRING {
			panic(fmt.Errorf("Operator + not defined for type %s", left))
		}
		tc.stack.Push(left)

	case token.MINUS:
		if left != symboltable.INTEGER {
			panic(fmt.Errorf("Operator - not defined for type %s", left))
		}
		tc.stack.Push(left)

	case token.MULTIPLY:
		if left != symboltable.INTEGER {
			panic(fmt.Errorf("Operator * not defined for type %s", left))
		}
		tc.stack.Push(left)

	case token.INTEGER_DIV:
		if left != symboltable.INTEGER {
			panic(fmt.Errorf("Operator / not defined for type %s", left))
		}
		tc.stack.Push(left)

	case token.AND:
		if left != symboltable.BOOLEAN {
			panic(fmt.Errorf("Operator & not defined for type %s", left))
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
		panic("NOT A symboltable.BOOLEAN")
	}
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
