package symboltable

import (
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

var symbolCreatorTestCases = []struct {
	name           string
	input          ast.Node
	expectedOutput *SymbolTable
	shouldPanic    bool
}{
	// ASSIGNMENTS
	{
		name: "Assignment before declaration",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "x")},
						Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldPanic: true,
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
		shouldPanic: false,
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
					},
				},
			},
		},
		shouldPanic: true,
	},
	// READ
	{
		name: "Read statement to undeclared variable",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.ReadStmt{
						TargetIdentifier: ast.Ident{Id: token.New(token.IDENT, "asd")},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// FOR LOOP
	{
		name: "For statement with undeclared index",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "foo")},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		shouldPanic: true,
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
								},
							},
						},
					},
				},
			},
		},
		shouldPanic: true,
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
		shouldPanic: false,
	},
}

func TestCreateSymbolTable(t *testing.T) {
	for _, testCase := range symbolCreatorTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil && testCase.shouldPanic {
					t.Error("Expected a panic")
				} else if r != nil && !testCase.shouldPanic {
					t.Errorf("Did not expect a panic, got '%v'", r)
					debug.PrintStack()
				}
			}()

			stc := &SymbolTableCreator{}
			symbols := stc.Create(testCase.input)

			if !testCase.shouldPanic && !reflect.DeepEqual(symbols, testCase.expectedOutput) {
				t.Errorf("Expected:\n%+v\ngot:\n%+v", testCase.expectedOutput, symbols)
			}
		})
	}
}
