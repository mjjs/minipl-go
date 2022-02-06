package interpreter

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/src/ast"
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
}

func TestInterpreter(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			interpreter := NewWithOutputWriter(w)
			interpreter.Run(testCase.input)

			if !reflect.DeepEqual(interpreter.variables, testCase.expectedVariables) {
				t.Errorf("VARIABLES DO NOT MATCH!!!!!!!!1")
			}

			if w.String() != testCase.expectedOutput.String() {
				t.Errorf("Expected %s, got %s", testCase.expectedOutput.String(), w.String())
			}
		})
	}
}
