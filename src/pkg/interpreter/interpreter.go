package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mjjs/minipl-go/pkg/ast"
	"github.com/mjjs/minipl-go/pkg/stack"
	"github.com/mjjs/minipl-go/pkg/token"
)

type Interpreter struct {
	stack *stack.Stack

	variables    map[string]interface{}
	outputWriter io.Writer
	inputReader  io.Reader
}

func New(outputWriter io.Writer, inputReader io.Reader) *Interpreter {
	return &Interpreter{
		stack:        stack.New(),
		variables:    make(map[string]interface{}),
		outputWriter: outputWriter,
		inputReader:  inputReader,
	}
}

func NewWithOutputWriter(output io.Writer) *Interpreter {
	return &Interpreter{
		stack:        stack.New(),
		variables:    make(map[string]interface{}),
		outputWriter: output,
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
	varName := node.Identifier.Id.Value()
	node.Expression.Accept(i)
	value := i.stack.Pop()
	i.variables[varName] = value
}

func (i *Interpreter) VisitDeclStmt(node ast.DeclStmt) {
	varName := node.Identifier.Value()
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
	idx := node.Index.Id.Value()

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
	varName := node.TargetIdentifier.Id.Value()

	x := i.variables[varName]

	r := bufio.NewReader(i.inputReader)

	if _, ok := x.(int); ok {
		str, _ := r.ReadString('\n')
		x, err := strconv.Atoi(strings.Trim(str, "\n"))
		if err != nil {
			i.terminate(fmt.Sprintf("%s: runtime error: failed to parse integer", node.Position()))
		}
		i.variables[varName] = x
	} else if _, ok := x.(string); ok {
		x, err := r.ReadString('\n')
		if err != nil {
			i.terminate(fmt.Sprintf("%s: runtime error: failed to parse string", node.Position()))
		}
		i.variables[varName] = x
	} else {
		i.terminate(fmt.Sprintf("%s: runtime error: could not read user input", node.Position()))
	}
}

func (i *Interpreter) VisitPrintStmt(node ast.PrintStmt) {
	node.Expression.Accept(i)
	fmt.Fprint(i.outputWriter, i.stack.Pop())
}

func (i *Interpreter) VisitAssertStmt(node ast.AssertStmt) {
	node.Expression.Accept(i)
	if !i.stack.Pop().(bool) {
		i.terminate(fmt.Sprintf("%s: runtime error: assert failed", node.Position()))
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
		panic(fmt.Sprintf("Encountered an unsupported operator %v", operator))
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
	i.stack.Push(i.variables[node.Id.Value()])
}

func (i *Interpreter) terminate(message string) {
	fmt.Fprintln(i.outputWriter, message)
	os.Exit(1)
}
