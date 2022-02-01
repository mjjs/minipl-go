package interpreter

import (
	"fmt"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/lexer"
	"github.com/mjjs/minipl-go/src/stack"
)

type Lexer interface {
	GetNextToken() lexer.Token
}

type Parser interface {
	Parse() (ast.Stmts, error)
}

type Interpreter struct {
	Lexer  Lexer
	Parser Parser

	stack *stack.Stack
}

func New(lexer Lexer, parser Parser) *Interpreter {
	return &Interpreter{
		Lexer:  lexer,
		Parser: parser,
		stack:  stack.New(),
	}
}

func (i *Interpreter) Run() {
	program, err := i.Parser.Parse()
	if err != nil {
		panic(err)
	}

	program.Accept(i)
}

func (i *Interpreter) VisitStmts(node ast.Stmts) {
	for _, stmt := range node.Statements {
		stmt.Accept(i)
	}
}

func (i *Interpreter) VisitAssignStmt(node ast.AssignStmt) {}
func (i *Interpreter) VisitDeclStmt(node ast.DeclStmt)     {}

func (i *Interpreter) VisitForStmt(node ast.ForStmt) {}

func (i *Interpreter) VisitReadStmt(node ast.ReadStmt) {}
func (i *Interpreter) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(i)
	fmt.Println(i.stack.Pop())
}
func (i *Interpreter) VisitAssertStmt(node ast.AssertStmt) {
	node.Expression.Accept(i)
	if !i.stack.Pop().(bool) {
		panic("ASSERT FAILED!")
	}
}

func (i *Interpreter) VisitBinaryExpr(node ast.BinaryExpr) {
	operator := node.Operator.Type()

	node.Left.Accept(i)
	left := i.stack.Pop()

	node.Right.Accept(i)
	right := i.stack.Pop()

	switch operator {
	case lexer.PLUS:
		{
			l, leftOk := left.(int)
			r, rightOk := right.(int)

			if leftOk && rightOk {
				i.stack.Push(l + r)
			}
		}

		{
			l, leftOk := left.(string)
			r, rightOk := right.(string)

			if leftOk && rightOk {
				i.stack.Push(l + r)
			}
		}

	case lexer.MINUS:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l - r)
		}

	case lexer.INTEGER_DIV:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l / r)
		}

	case lexer.MULTIPLY:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l * r)
		}

	case lexer.AND:
		l, leftOk := left.(bool)
		r, rightOk := right.(bool)

		if leftOk && rightOk {
			i.stack.Push(l && r)
		}

	case lexer.LT:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l < r)
		}

	case lexer.EQ:
		i.stack.Push(left == right)

	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	panic(fmt.Sprintf("Unsupported operands %v and %v for operator %v", left, right, operator))
}

func (i *Interpreter) VisitUnaryExpr(node ast.UnaryExpr) {
	node.Operand.Accept(i)
	val := i.stack.Pop()

	switch node.Unary.Type() {
	case lexer.NOT:
		i.stack.Push(!val.(bool))
	default:
		panic(fmt.Sprintf("Unsupported unary type %v", node.Unary.Type()))
	}
}

func (i *Interpreter) VisitNullaryExpr(node ast.NullaryExpr) {
	node.Operand.Accept(i)
}

func (i *Interpreter) VisitNumberOpnd(node ast.NumberOpnd) {
	i.stack.Push(node.Value)
}

func (i *Interpreter) VisitStringOpnd(node ast.StringOpnd) {
	i.stack.Push(node.Value)
}

func (i *Interpreter) VisitIdent(node ast.Ident) {}
