package symboltable

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

var symbolCreatorTestCases = []struct {
	name           string
	input          ast.Node
	expectedOutput *SymbolTable
	expectedErrors []error
}{
	// ASSIGNMENTS
	{
		name: "Assignment before declaration",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssignStmt{
						Identifier: ast.Ident{
							Id:  token.New(token.IDENT, "x"),
							Pos: token.Position{Line: 5, Column: 1},
						},
						Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		expectedErrors: []error{
			errors.New("5:1: variable x used before declaration"),
		},
	},
	{
		name: "Assignment after declaration",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.STRING, ""),
						Expression:   ast.NullaryExpr{Operand: ast.StringOpnd{Value: "12345"}},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "foo")},
						Expression: ast.NullaryExpr{Operand: ast.StringOpnd{Value: "67890"}},
					},
				},
			},
		},
		expectedOutput: &SymbolTable{
			symbols: map[string]Symbol{
				"foo": {STRING},
			},
		},
	},
	// DECLARATION
	{
		name: "Duplicate declaration",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.INTEGER, ""),
						Pos:          token.Position{Line: 13, Column: 37},
					},
				},
			},
		},
		expectedErrors: []error{
			errors.New("13:37: redeclaration of variable foo"),
		},
	},
	// READ
	{
		name: "Read statement to undeclared variable",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.ReadStmt{
						TargetIdentifier: ast.Ident{
							Id:  token.New(token.IDENT, "asd"),
							Pos: token.Position{Line: 666, Column: 666},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			errors.New("666:666: variable asd used before declaration"),
		},
	},
	// FOR LOOP
	{
		name: "For statement with undeclared index",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.ForStmt{
						Index: ast.Ident{
							Id:  token.New(token.IDENT, "foo"),
							Pos: token.Position{Line: 1, Column: 2},
						},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		expectedErrors: []error{
			errors.New("1:2: variable foo used before declaration"),
		},
	},
	{
		name: "For loop index modified inside loop",
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
								ast.AssignStmt{
									Identifier: ast.Ident{Id: token.New(token.IDENT, "i")},
									Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 15}},
									Pos:        token.Position{Line: 7, Column: 2},
								},
							},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			errors.New("7:2: cannot modify loop index i during loop"),
		},
	},
	{
		name: "For loop index modified after loop",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 0}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
						Statements: ast.Stmts{},
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "i")},
						Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 15}},
					},
				},
			},
		},
		expectedOutput: &SymbolTable{
			symbols: map[string]Symbol{
				"i": {INTEGER},
			},
		},
	},
}

func TestCreateSymbolTable(t *testing.T) {
	for _, testCase := range symbolCreatorTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			stc := &SymbolTableCreator{}
			symbols, errors := stc.Create(testCase.input)

			if len(testCase.expectedErrors) > 0 {
				if len(testCase.expectedErrors) != len(errors) {
					t.Errorf(
						"\nExpected %d error(s) (%s),\ngot %d (%s)",
						len(testCase.expectedErrors),
						testCase.expectedErrors,
						len(errors),
						errors,
					)
				}

				for i, err := range testCase.expectedErrors {
					actual := errors[i]
					if actual.Error() != err.Error() {
						t.Errorf("Expected:\n%s\ngot:\n%s", err.Error(), actual.Error())
					}
				}
			} else if len(testCase.expectedErrors) == 0 && len(errors) > 0 {
				t.Errorf(
					"Expected no errors, got %d (%s)",
					len(errors), errors,
				)
			} else if !reflect.DeepEqual(symbols, testCase.expectedOutput) {
				t.Errorf("Expected:\n%+v\ngot:\n%+v", testCase.expectedOutput, symbols)
			}
		})
	}
}
