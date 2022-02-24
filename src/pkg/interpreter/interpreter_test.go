package interpreter

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/pkg/ast"
	"github.com/mjjs/minipl-go/pkg/token"
)

var testCases = []struct {
	name              string
	input             ast.Prog
	expectedVariables map[string]interface{}
	expectedOutput    *bytes.Buffer
}{
	{
		name: "Print can print strings and ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.NullaryExpr{
							Operand: ast.StringOpnd{Value: "foo\n"},
						},
					},
					ast.PrintStmt{
						Expression: ast.NullaryExpr{
							Operand: ast.NumberOpnd{Value: 666},
						},
					},
				},
			},
		},
		expectedVariables: make(map[string]interface{}),
		expectedOutput:    bytes.NewBufferString("foo\n666"),
	},
	{
		name: "For loop runs correct number of times",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.ForStmt{
						Index: ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:   ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 0}},
						High:  ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						Statements: ast.Stmts{
							Statements: []ast.Stmt{
								ast.PrintStmt{
									Expression: ast.NullaryExpr{
										Operand: ast.Ident{Id: token.New(token.IDENT, "i")},
									},
								},
								ast.PrintStmt{
									Expression: ast.NullaryExpr{
										Operand: ast.StringOpnd{Value: "\n"},
									},
								},
							},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{"i": 4},
		expectedOutput:    bytes.NewBufferString("0\n1\n2\n3\n4\n"),
	},
	{
		name: "Plus operation works properly",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
						Expression:   ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "s"),
						VariableType: token.New(token.STRING, ""),
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "x")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "s")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "Hello "}},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "World\n"}},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{"x": 10, "s": "Hello World\n"},
		expectedOutput:    bytes.NewBufferString(""),
	},
	{
		name: "Minus operation works properly",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
						Expression:   ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "x")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{"x": 0},
		expectedOutput:    bytes.NewBufferString(""),
	},
	{
		name: "Integer division works properly",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
						Expression:   ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "x")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{"x": 1},
		expectedOutput:    bytes.NewBufferString(""),
	},
	{
		name: "Multiplication works properly",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "x")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{"x": 25},
		expectedOutput:    bytes.NewBufferString(""),
	},
	{
		name: "Logical AND works properly",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "true"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "false"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 0}},
						},
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "isTrue"),
						VariableType: token.New(token.BOOLEAN, ""),
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "isFalse"),
						VariableType: token.New(token.BOOLEAN, ""),
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "isTrue")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.Ident{Id: token.New(token.IDENT, "true")}},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.Ident{Id: token.New(token.IDENT, "true")}},
						},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "isFalse")},
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.Ident{Id: token.New(token.IDENT, "true")}},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.Ident{Id: token.New(token.IDENT, "false")}},
						},
					},
				},
			},
		},
		expectedVariables: map[string]interface{}{
			"true":    true,
			"false":   false,
			"isTrue":  true,
			"isFalse": false,
		},
		expectedOutput: bytes.NewBufferString(""),
	},
}

func TestInterpreter(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			interpreter := NewWithOutputWriter(w)
			interpreter.Run(testCase.input)

			if !reflect.DeepEqual(interpreter.variables, testCase.expectedVariables) {
				t.Errorf("Expected variables to be in state %v, got %v", testCase.expectedVariables, interpreter.variables)
			}

			if w.String() != testCase.expectedOutput.String() {
				t.Errorf("Expected %s, got %s", testCase.expectedOutput.String(), w.String())
			}
		})
	}
}
