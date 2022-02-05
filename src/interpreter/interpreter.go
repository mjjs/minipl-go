package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/stack"
	"github.com/mjjs/minipl-go/src/token"
)

type Interpreter struct {
	stack *stack.Stack

	variables map[string]interface{}
}

func New() *Interpreter {
	return &Interpreter{
		stack:     stack.New(),
		variables: make(map[string]interface{}),
	}
}

func (i *Interpreter) Run(program ast.Prog) {
	program.Accept(i)
}

func (i *Interpreter) VisitProg(node ast.Prog) {
	node.Statements.Accept(i)
}

func (i *Interpreter) VisitStmts(node ast.Stmts) {
	for _, stmt := range node.Statements {
		stmt.Accept(i)
	}
}

func (i *Interpreter) VisitAssignStmt(node ast.AssignStmt) {
	varName := node.Identifier.Id.ValueString()
	node.Expression.Accept(i)
	value := i.stack.Pop()
	i.variables[varName] = value
}

func (i *Interpreter) VisitDeclStmt(node ast.DeclStmt) {
	varName := node.Identifier.ValueString()
	var value interface{}

	if node.Expression != nil {
		node.Expression.Accept(i)
		value = i.stack.Pop()
	} else {
		typ := node.VariableType.Type()
		if typ == token.INTEGER {
			value = 0
		} else if typ == token.STRING {
			value = ""
		} else {
			value = false
		}
	}

	i.variables[varName] = value
}

func (i *Interpreter) VisitForStmt(node ast.ForStmt) {
	idx := node.Index.Id.ValueString()

	node.Low.Accept(i)
	low := i.stack.Pop().(int)

	node.High.Accept(i)
	high := i.stack.Pop().(int)

	for j := low; j < high; j++ {
		i.variables[idx] = j
		node.Statements.Accept(i)
	}
}

func (i *Interpreter) VisitReadStmt(node ast.ReadStmt) {
	varName := node.TargetIdentifier.Id.ValueString()

	x := i.variables[varName]

	r := bufio.NewReader(os.Stdin)

	if _, ok := x.(int); ok {
		str, _ := r.ReadString('\n')
		x, _ = strconv.Atoi(str)
		i.variables[varName] = x
	} else if _, ok := x.(string); ok {
		x, _ = r.ReadString('\n')
		i.variables[varName] = x
	} else {
		panic("WTF NOT A PROPER TYPE (bool not supported yet!)")
	}
}

func (i *Interpreter) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(i)
	fmt.Print(i.stack.Pop())
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
	case token.PLUS:
		{
			l, leftOk := left.(int)
			r, rightOk := right.(int)

			if leftOk && rightOk {
				i.stack.Push(l + r)
				return
			}
		}

		{
			l, leftOk := left.(string)
			r, rightOk := right.(string)

			if leftOk && rightOk {
				i.stack.Push(l + r)
				return
			}
		}

	case token.MINUS:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l - r)
			return
		}

	case token.INTEGER_DIV:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l / r)
			return
		}

	case token.MULTIPLY:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l * r)
			return
		}

	case token.AND:
		l, leftOk := left.(bool)
		r, rightOk := right.(bool)

		if leftOk && rightOk {
			i.stack.Push(l && r)
			return
		}

	case token.LT:
		l, leftOk := left.(int)
		r, rightOk := right.(int)

		if leftOk && rightOk {
			i.stack.Push(l < r)
			return
		}

	case token.EQ:
		i.stack.Push(left == right)
		return

	default:
		panic(fmt.Sprintf("Unsupported operator %v", operator))
	}

	panic(fmt.Sprintf("Unsupported operands %v and %v for operator %v", left, right, operator))
}

func (i *Interpreter) VisitUnaryExpr(node ast.UnaryExpr) {
	node.Operand.Accept(i)
	val := i.stack.Pop()

	switch node.Unary.Type() {
	case token.NOT:
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

func (i *Interpreter) VisitIdent(node ast.Ident) {
	i.stack.Push(i.variables[node.Id.ValueString()])
}
