package semantics

import (
	"testing"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/lexer"
)

var analyzerTestCases = []struct {
	name        string
	input       ast.Stmts
	shouldError bool
}{
	// ASSIGNMENTS
	{
		name: "Assignment before declaration",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssignStmt{
					Identifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "x")},
					Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Assignment after declaration",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.STRING, nil),
					Expression:   ast.NullaryExpr{Operand: ast.StringOpnd{Value: "12345"}},
				},
				ast.AssignStmt{
					Identifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "foo")},
					Expression: ast.NullaryExpr{Operand: ast.StringOpnd{Value: "67890"}},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Assignment with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.STRING, nil),
				},
				ast.AssignStmt{
					Identifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "foo")},
					Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 67890}},
				},
			},
		},
		shouldError: true,
	},
	// DECLARATION
	{
		name: "Duplicate declaration",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
			},
		},
		shouldError: true,
	},
	// ASSERT
	{
		name: "Assert with non-boolean type",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssertStmt{
					Expression: ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Assert with boolean type",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssertStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
					},
				},
			},
		},
		shouldError: false,
	},
	// NOT
	{
		name: "Not operator with non-boolean operand",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.BOOLEAN, nil),
					Expression: ast.UnaryExpr{
						Unary:   lexer.NewToken(lexer.NOT, nil),
						Operand: ast.StringOpnd{Value: "12345"},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Not operator with boolean operand",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.BOOLEAN, nil),
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
					},
				},
			},
		},
		shouldError: false,
	},
	// EQUALITY OPERATOR
	{
		name: "Equality with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Equality with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Equality with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Equality with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// PLUS OPERATOR
	{
		name: "Plus with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.PLUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Plus with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.PLUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Plus with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.PLUS, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Plus with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.PLUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// MINUS OPERATOR
	{
		name: "Minus with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.MINUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Minus with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.MINUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Minus with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.MINUS, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Minus with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.MINUS, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// MULTIPLY OPERATOR
	{
		name: "Multiply with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.MULTIPLY, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Multiply with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.MULTIPLY, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Multiply with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.MULTIPLY, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Multiply with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.MULTIPLY, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// INTEGER DIVISION OPERATOR
	{
		name: "Integer division with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.INTEGER_DIV, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Integer division with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.INTEGER_DIV, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Integer division with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.INTEGER_DIV, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "Integer division with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.INTEGER_DIV, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// AND OPERATOR
	{
		name: "AND with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.AND, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "AND with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.AND, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "AND with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.AND, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "AND with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.AND, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// LESS THAN OPERATOR
	{
		name: "LESS THAN with ints",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
						Operator: lexer.NewToken(lexer.LT, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "LESS THAN with strings",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "a"}},
						Operator: lexer.NewToken(lexer.LT, nil),
						Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "b"}},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "LESS THAN with booleans",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
						Operator: lexer.NewToken(lexer.LT, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NullaryExpr{Operand: ast.StringOpnd{Value: "foo"}},
						},
					},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "LESS THAN with unmatched types",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.NullaryExpr{Operand: ast.StringOpnd{Value: "6"}},
						Operator: lexer.NewToken(lexer.LT, nil),
						Right:    ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 6}},
					},
				},
			},
		},
		shouldError: true,
	},
	// READ
	{
		name: "Read statement to undeclared variable",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.ReadStmt{
					TargetIdentifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "asd")},
				},
			},
		},
		shouldError: true,
	},
	// FOR LOOP
	{
		name: "For statement with undeclared index",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "foo")},
					Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					Statements: ast.Stmts{},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "For statement with non-integer index",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.STRING, nil),
				},
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					Statements: ast.Stmts{},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "For statement with non-integer lower range expr",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:        ast.NullaryExpr{Operand: ast.StringOpnd{Value: "1"}},
					High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					Statements: ast.Stmts{},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "For statement with non-integer higher range expr",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					High:       ast.NullaryExpr{Operand: ast.StringOpnd{Value: "1"}},
					Statements: ast.Stmts{},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "For loop index modified inside loop",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.ForStmt{
					Index: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:   ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 0}},
					High:  ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
					Statements: ast.Stmts{
						Statements: []ast.Stmt{
							ast.AssignStmt{
								Identifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
								Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 15}},
							},
						},
					},
				},
			},
		},
		shouldError: true,
	},
	{
		name: "For loop index modified after loop",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 0}},
					High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 5}},
					Statements: ast.Stmts{},
				},
				ast.AssignStmt{
					Identifier: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Expression: ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 15}},
				},
			},
		},
		shouldError: false,
	},
	{
		name: "Valid for statement",
		input: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "i"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.ForStmt{
					Index:      ast.Ident{Id: lexer.NewToken(lexer.IDENT, "i")},
					Low:        ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 1}},
					High:       ast.NullaryExpr{Operand: ast.NumberOpnd{Value: 20}},
					Statements: ast.Stmts{},
				},
			},
		},
		shouldError: false,
	},
}

func TestAnalyze(t *testing.T) {
	for _, testCase := range analyzerTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			analyzer := NewAnalyzer()
			err := analyzer.Analyze(testCase.input)

			if testCase.shouldError && err == nil {
				t.Error("Expected an error, got nil")
			} else if !testCase.shouldError && err != nil {
				t.Errorf("Expected nil error, got %v", err)
			}
		})
	}
}
