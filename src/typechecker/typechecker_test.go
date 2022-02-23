package typechecker

import (
	"fmt"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/symboltable"
	"github.com/mjjs/minipl-go/token"
)

var typeCheckerTestCases = []struct {
	name           string
	input          ast.Node
	symbols        *symboltable.SymbolTable
	expectedErrors []error
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
						Pos:        token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("foo", symboltable.STRING),
		expectedErrors: []error{
			fmt.Errorf(
				"1:1: cannot assign type %s to variable foo of type %s",
				symboltable.INTEGER, symboltable.STRING,
			),
		},
	},
	{
		name: "Declare without assignment",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.STRING, ""),
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("foo", symboltable.STRING),
	},
	// ASSERT
	{
		name: "Assert with non-boolean type",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssertStmt{
						Expression: ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						Pos:        token.Position{Line: 22, Column: 1},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"22:1: assert statement is only defined for type %s, not %s",
				symboltable.BOOLEAN, symboltable.STRING,
			),
		},
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
							Pos:     token.Position{Line: 13, Column: 1},
						},
						Pos: token.Position{Line: 12, Column: 1},
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("foo", symboltable.BOOLEAN),
		expectedErrors: []error{
			fmt.Errorf(
				"13:1: unary operator %s not defined for type %s",
				token.NOT, symboltable.STRING,
			),
			fmt.Errorf(
				"12:1: cannot assign type %s to variable foo of type %s",
				symboltable.STRING, symboltable.BOOLEAN,
			),
		},
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
		symbols: symboltable.NewSymbolTable().Insert("foo", symboltable.BOOLEAN),
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
	},
	{
		name: "Equality with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.EQ, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.EQ,
			),
		},
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
	},
	{
		name: "Plus with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.NullaryExpr{
									Operand: ast.StringOpnd{
										Value: "foo",
										Pos:   token.Position{Line: 5, Column: 1},
									},
								},
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
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.PLUS, symboltable.BOOLEAN,
			),
		},
	},
	{
		name: "Plus with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.PLUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.PLUS,
			),
		},
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
	},
	{
		name: "Minus with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "a",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
						Pos: token.Position{Line: 5, Column: 1},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MINUS, symboltable.STRING,
			),
		},
	},
	{
		name: "Minus with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.NullaryExpr{
									Operand: ast.StringOpnd{
										Value: "foo",
										Pos:   token.Position{Line: 5, Column: 1},
									},
								},
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
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MINUS, symboltable.BOOLEAN,
			),
		},
	},
	{
		name: "Minus with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.MINUS, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.MINUS,
			),
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MINUS, symboltable.STRING,
			),
		},
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
	},
	{
		name: "Multiply with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "a",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MULTIPLY, symboltable.STRING,
			),
		},
	},
	{
		name: "Multiply with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.NullaryExpr{
									Operand: ast.StringOpnd{
										Value: "foo",
										Pos:   token.Position{Line: 5, Column: 1},
									},
								},
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
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MULTIPLY, symboltable.BOOLEAN,
			),
		},
	},
	{
		name: "Multiply with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.MULTIPLY,
			),
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.MULTIPLY, symboltable.STRING,
			),
		},
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
	},
	{
		name: "Integer division with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "a",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.INTEGER_DIV, symboltable.STRING,
			),
		},
	},
	{
		name: "Integer division with booleans",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.NullaryExpr{
									Operand: ast.StringOpnd{
										Value: "foo",
										Pos:   token.Position{Line: 5, Column: 1},
									},
								},
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
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.INTEGER_DIV, symboltable.BOOLEAN,
			),
		},
	},
	{
		name: "Integer division with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.INTEGER_DIV,
			),
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.INTEGER_DIV, symboltable.STRING,
			),
		},
	},
	// AND OPERATOR
	{
		name: "AND with ints",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.NumberOpnd{
									Value: 1,
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.AND, symboltable.INTEGER,
			),
		},
	},
	{
		name: "AND with strings",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "a",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.AND, symboltable.STRING,
			),
		},
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
	},
	{
		name: "AND with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 5, Column: 1},
								},
							},
							Operator: token.New(token.AND, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"5:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.AND,
			),
			fmt.Errorf(
				"5:1: operator %s not defined for type %s",
				token.AND, symboltable.STRING,
			),
		},
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
	},
	{
		name: "LESS THAN with unmatched types",
		input: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.NullaryExpr{
								Operand: ast.StringOpnd{
									Value: "6",
									Pos:   token.Position{Line: 666, Column: 1},
								},
							},
							Operator: token.New(token.LT, ""),
							Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
						},
					},
				},
			},
		},
		expectedErrors: []error{
			fmt.Errorf(
				"666:1: unmatched types %s and %s for binary expression %s",
				symboltable.STRING, symboltable.INTEGER, token.LT,
			),
		},
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
						Pos:        token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("i", symboltable.STRING),
		expectedErrors: []error{
			fmt.Errorf(
				"1:1: loop index must be %s, not %s",
				symboltable.INTEGER, symboltable.STRING,
			),
		},
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
						Pos:        token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
		expectedErrors: []error{
			fmt.Errorf(
				"1:1: for loop range lower bound must be %s, not %s",
				symboltable.INTEGER, symboltable.STRING,
			),
		},
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
						Pos:        token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
		symbols: symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
		expectedErrors: []error{
			fmt.Errorf(
				"1:1: for loop range upper bound must be %s, not %s",
				symboltable.INTEGER, symboltable.STRING,
			),
		},
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
		symbols: symboltable.NewSymbolTable().Insert("i", symboltable.INTEGER),
	},
}

func TestCheckTypes(t *testing.T) {
	for _, testCase := range typeCheckerTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			typeChecker := New(testCase.symbols)
			errors := typeChecker.CheckTypes(testCase.input)

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
			}
		})
	}
}
