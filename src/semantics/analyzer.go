package semantics

import (
	"fmt"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/lexer"
)

type Analyzer struct {
	symbols *SymbolTable

	lastType SymbolType
	err      error
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{symbols: NewSymbolTable()}
}

func (a *Analyzer) Analyze(statements ast.Stmts) error {
	statements.Accept(a)
	return a.err
}

func (i *Analyzer) VisitStmts(node ast.Stmts) {
	if i.err != nil {
		return
	}

	for _, stmt := range node.Statements {
		stmt.Accept(i)
	}
}

func (i *Analyzer) VisitDeclStmt(node ast.DeclStmt) {
	if i.err != nil {
		return
	}

	name := node.Identifier.ValueString()
	_, exists := i.symbols.Get(name)
	if exists {
		i.err = fmt.Errorf("Variable %s has already been declared", name)
		return
	}

	var variableType SymbolType

	switch node.VariableType.Type() {
	case lexer.INTEGER:
		variableType = INTEGER
	case lexer.STRING:
		variableType = STRING
	case lexer.BOOLEAN:
		variableType = BOOLEAN
	}

	i.symbols.Insert(name, variableType)

	if node.Expression != nil {
		node.Expression.Accept(i)

		if i.lastType != variableType {
			i.err = fmt.Errorf("Types %s and %s do not match", variableType, i.lastType)
			return
		}
	}
}

func (i *Analyzer) VisitAssignStmt(node ast.AssignStmt) {
	if i.err != nil {
		return
	}

	node.Identifier.Accept(i)
	lsType := i.lastType

	node.Expression.Accept(i)
	rsType := i.lastType

	if lsType != rsType {
		i.err = fmt.Errorf("Types %s %s do not match", lsType, rsType)
		return
	}
}

func (i *Analyzer) VisitForStmt(node ast.ForStmt) {
	if i.err != nil {
		return
	}

	node.Index.Accept(i)
	indexType := i.lastType

	node.Low.Accept(i)
	lowType := i.lastType

	node.High.Accept(i)
	highType := i.lastType

	if indexType != INTEGER {
		i.err = fmt.Errorf("Unsupported type %s as loop index", indexType)
		return
	}

	if lowType != INTEGER {
		i.err = fmt.Errorf("Unsupported type %s as lower range bound", lowType)
		return
	}

	if highType != INTEGER {
		i.err = fmt.Errorf("Unsupported type %s as lower range bound", highType)
		return
	}

	node.Statements.Accept(i)
}

func (i *Analyzer) VisitReadStmt(node ast.ReadStmt) {
	if i.err != nil {
		return
	}

	node.TargetIdentifier.Accept(i)
}

func (i *Analyzer) VisitPrintStmt(node ast.PrintStmt) {
	if i.err != nil {
		return
	}

	node.Expression.Accept(i)
}

func (i *Analyzer) VisitAssertStmt(node ast.AssertStmt) {
	if i.err != nil {
		return
	}

	node.Expression.Accept(i)
	if i.lastType != BOOLEAN {
		i.err = fmt.Errorf("Invalid type assert(%s)", i.lastType)
		return
	}
}

func (i *Analyzer) VisitBinaryExpr(node ast.BinaryExpr) {
	if i.err != nil {
		return
	}

	node.Left.Accept(i)
	leftType := i.lastType

	node.Right.Accept(i)
	rightType := i.lastType

	if leftType != rightType {
		i.err = fmt.Errorf("Mismatched types %s %v %s", leftType, node.Operator.Type(), rightType)
		return
	}

	switch node.Operator.Type() {
	case lexer.PLUS:
		if leftType != INTEGER && leftType != STRING {
			i.err = fmt.Errorf("Operator + not defined for type %s", leftType)
			return
		}
		i.lastType = leftType

	case lexer.MINUS:
		if leftType != INTEGER {
			i.err = fmt.Errorf("Operator - not defined for type %s", leftType)
			return
		}
		i.lastType = INTEGER

	case lexer.MULTIPLY:
		if leftType != INTEGER {
			i.err = fmt.Errorf("Operator * not defined for type %s", leftType)
			return
		}
		i.lastType = INTEGER

	case lexer.INTEGER_DIV:
		if leftType != INTEGER {
			i.err = fmt.Errorf("Operator / not defined for type %s", leftType)
			return
		}
		i.lastType = INTEGER

	case lexer.AND:
		if leftType != BOOLEAN {
			i.err = fmt.Errorf("Operator & not defined for type %s", leftType)
			return
		}

		i.lastType = BOOLEAN

	case lexer.LT:
		i.lastType = BOOLEAN

	case lexer.EQ:
		i.lastType = BOOLEAN
	}
}

func (i *Analyzer) VisitUnaryExpr(node ast.UnaryExpr) {
	if i.err != nil {
		return
	}

	node.Operand.Accept(i)
	if node.Unary.Type() == lexer.NOT && i.lastType != BOOLEAN {
		i.err = fmt.Errorf("Operator ! not defined for type %s", i.lastType)
		return
	}
}

func (i *Analyzer) VisitNullaryExpr(node ast.NullaryExpr) {
	if i.err != nil {
		return
	}

	node.Operand.Accept(i)
}

func (i *Analyzer) VisitNumberOpnd(node ast.NumberOpnd) {
	if i.err != nil {
		return
	}

	i.lastType = INTEGER
}

func (i *Analyzer) VisitStringOpnd(node ast.StringOpnd) {
	if i.err != nil {
		return
	}

	i.lastType = STRING
}

func (i *Analyzer) VisitIdent(node ast.Ident) {
	if i.err != nil {
		return
	}

	name := node.Id.ValueString()
	symbolType, exists := i.symbols.Get(name)
	if !exists {
		i.err = fmt.Errorf("Variable %s has not been declared", name)
		return
	}

	i.lastType = symbolType
}
