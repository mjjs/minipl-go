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
