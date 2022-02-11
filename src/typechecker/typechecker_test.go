package typechecker

import (
	"runtime/debug"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/symboltable"
	"github.com/mjjs/minipl-go/token"
)

var typeCheckerTestCases = []struct {
	name        string
	input       ast.Node
	symbols     *symboltable.SymbolTable
	shouldPanic bool
}{
	{
		name: "Assignment with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.STRING, ""),
					},
					ast.AssignStmt{
						Identifier: ast.Ident{Id: token.New(token.IDENT, "foo")},
						Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 67890}},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("foo", symboltable.STRING),
		shouldPanic: true,
	},
	// ASSERT
	{
		name: "Assert with non-boolean type",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssertStmt{
						Expression: ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Assert with boolean type",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssertStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	// NOT
	{
		name: "Not operator with non-boolean operand",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.UnaryExpr{
							Unary:   token.New(token.NOT, ""),
							Operand: ast.StringOpnd{Value: "12345"},
						},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("foo", symboltable.BOOLEAN),
		shouldPanic: true,
	},
	{
		name: "Not operator with boolean operand",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("foo", symboltable.BOOLEAN),
		shouldPanic: false,
	},
	// EQUALITY OPERATOR
	{
		name: "Equality with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Equality with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Equality with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.EQ, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Equality with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// PLUS OPERATOR
	{
		name: "Plus with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Plus with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Plus with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.PLUS, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Plus with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// MINUS OPERATOR
	{
		name: "Minus with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Minus with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Minus with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.MINUS, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Minus with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// MULTIPLY OPERATOR
	{
		name: "Multiply with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Multiply with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Multiply with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.MULTIPLY, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Multiply with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// symboltable.INTEGER DIVISION OPERATOR
	{
		name: "Integer division with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "Integer division with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Integer division with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "Integer division with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// AND OPERATOR
	{
		name: "AND with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "AND with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	{
		name: "AND with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.AND, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "AND with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// LESS THAN OPERATOR
	{
		name: "LESS THAN with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
							Operator: token.New(token.LT, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "LESS THAN with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
							Operator: token.New(token.LT, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "LESS THAN with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
							Operator: token.New(token.LT, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							},
						},
					},
				},
			},
		},
		shouldPanic: false,
	},
	{
		name: "LESS THAN with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
							Operator: token.New(token.LT, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		shouldPanic: true,
	},
	// FOR LOOP
	{
		name: "For statement with non-integer index",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.STRING, ""),
					},
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("i", symboltable.STRING),
		shouldPanic: true,
	},
	{
		name: "For statement with non-integer lower range expr",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:        ast.NullaryExpr{Operand: ast.StringOpnd{Value: "1"}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
		shouldPanic: true,
	},
	{
		name: "For statement with non-integer higher range expr",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						High:       ast.NullaryExpr{Operand: ast.StringOpnd{Value: "1"}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
		shouldPanic: true,
	},
	{
		name: "Valid for statement",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "i"),
						VariableType: token.New(token.INTEGER, ""),
					},
					ast.ForStmt{
						Index:      ast.Ident{Id: token.New(token.IDENT, "i")},
						Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 20}},
						Statements: ast.Stmts{},
					},
				},
			},
		},
		symbols:     symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
		shouldPanic: false,
	},
}

func TestCheckTypes(t *testing.T) {
	for _, testCase := range typeCheckerTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil && testCase.shouldPanic {
					t.Error("Expected a panic")
				} else if r != nil && !testCase.shouldPanic {
					t.Errorf("Did not expect a panic, got '%v'", r)
					debug.PrintStack()
				}
			}()

			typeChecker := New(testCase.symbols)
			typeChecker.CheckTypes(testCase.input)
		})
	}
}
