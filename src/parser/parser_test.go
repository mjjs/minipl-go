package parser

import (
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/lexer"
)

var parseTestCases = []struct {
	name           string
	lexerOutput    []lexer.Token
	expectedOutput ast.Stmts
	shouldError    bool
}{
	{
		name: "Declaration with assignment",
		lexerOutput: []lexer.Token{
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "x"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.INTEGER, nil),
			lexer.NewToken(lexer.ASSIGN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 5),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "x"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
					Expression: ast.NullaryExpr{
						Operand: ast.NumberOpnd{Value: 5},
					},
				},
			},
		},
	},
	{
		name: "Declaration without assignment",
		lexerOutput: []lexer.Token{
			// Int
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "x"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.INTEGER, nil),
			lexer.NewToken(lexer.SEMI, nil),

			// String
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "y"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.STRING, nil),
			lexer.NewToken(lexer.SEMI, nil),

			// Boolean
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "z"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.BOOLEAN, nil),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "x"),
					VariableType: lexer.NewToken(lexer.INTEGER, nil),
				},
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "y"),
					VariableType: lexer.NewToken(lexer.STRING, nil),
				},
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "z"),
					VariableType: lexer.NewToken(lexer.BOOLEAN, nil),
				},
			},
		},
	},
	{
		name: "Assignment",
		lexerOutput: []lexer.Token{
			lexer.NewToken(lexer.IDENT, "foo"),
			lexer.NewToken(lexer.ASSIGN, nil),
			lexer.NewToken(lexer.STRING_LITERAL, "bar"),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssignStmt{
					Identifier: lexer.NewToken(lexer.IDENT, "foo"),
					Expression: ast.NullaryExpr{
						Operand: ast.StringOpnd{Value: "bar"},
					},
				},
			},
		},
	},
	{
		name: "For statement with multiple inner statements",
		lexerOutput: []lexer.Token{
			lexer.NewToken(lexer.FOR, nil),
			lexer.NewToken(lexer.IDENT, "i"),
			lexer.NewToken(lexer.IN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 3),
			lexer.NewToken(lexer.PLUS, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 2),
			lexer.NewToken(lexer.DOTDOT, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 25),
			lexer.NewToken(lexer.MINUS, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 1),
			lexer.NewToken(lexer.DO, nil),
			lexer.NewToken(lexer.READ, nil),
			lexer.NewToken(lexer.IDENT, "x"),
			lexer.NewToken(lexer.SEMI, nil),
			lexer.NewToken(lexer.PRINT, nil),
			lexer.NewToken(lexer.IDENT, "x"),
			lexer.NewToken(lexer.SEMI, nil),
			lexer.NewToken(lexer.END, nil),
			lexer.NewToken(lexer.FOR, nil),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.ForStmt{
					Index: lexer.NewToken(lexer.IDENT, "i"),
					Low: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: lexer.NewToken(lexer.PLUS, nil),
						Right:    ast.NumberOpnd{Value: 2},
					},
					High: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 25},
						Operator: lexer.NewToken(lexer.MINUS, nil),
						Right:    ast.NumberOpnd{Value: 1},
					},
					Statements: ast.Stmts{
						Statements: []ast.Stmt{
							ast.ReadStmt{TargetIdentifier: lexer.NewToken(lexer.IDENT, "x")},
							ast.PrintStmt{
								Expression: ast.NullaryExpr{
									Operand: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "x")},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "Assert statement",
		lexerOutput: []lexer.Token{
			lexer.NewToken(lexer.ASSERT, nil),
			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.STRING_LITERAL, "foo"),
			lexer.NewToken(lexer.EQ, nil),
			lexer.NewToken(lexer.STRING_LITERAL, "bar"),
			lexer.NewToken(lexer.RPAREN, nil),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssertStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.StringOpnd{Value: "foo"},
						Operator: lexer.NewToken(lexer.EQ, nil),
						Right:    ast.StringOpnd{Value: "bar"},
					},
				},
			},
		},
	},
	{
		name: "Error if no EOF is returned by lexer when expected",
		lexerOutput: []lexer.Token{
			lexer.NewToken(lexer.PRINT, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 22),
			lexer.NewToken(lexer.SEMI, nil),
			lexer.NewToken(lexer.SEMI, nil),
		},
		shouldError: true,
	},
	{
		name: "Parenthesised expressions",
		lexerOutput: []lexer.Token{
			// print (1 * 2) / (4 - 3)
			lexer.NewToken(lexer.PRINT, nil),

			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 1),
			lexer.NewToken(lexer.MULTIPLY, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 2),
			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.INTEGER_DIV, nil),

			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 4),
			lexer.NewToken(lexer.MINUS, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 3),
			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: lexer.NewToken(lexer.MULTIPLY, nil),
							Right:    ast.NumberOpnd{Value: 2},
						},
						Operator: lexer.NewToken(lexer.INTEGER_DIV, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 4},
							Operator: lexer.NewToken(lexer.MINUS, nil),
							Right:    ast.NumberOpnd{Value: 3},
						},
					},
				},
			},
		},
	},
	{
		name: "Less than comparison",
		lexerOutput: []lexer.Token{
			// var foo : bool := 3 < 2
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "foo"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.BOOLEAN, nil),
			lexer.NewToken(lexer.ASSIGN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 3),
			lexer.NewToken(lexer.LT, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 2),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.BOOLEAN, nil),
					Expression: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: lexer.NewToken(lexer.LT, nil),
						Right:    ast.NumberOpnd{Value: 2},
					},
				},
			},
		},
	},
	{
		name: "Nested logical and operator with not",
		lexerOutput: []lexer.Token{
			// var foo : bool := ((3 = 3) & (2 = 2)) & (1 = 0);
			lexer.NewToken(lexer.VAR, nil),
			lexer.NewToken(lexer.IDENT, "foo"),
			lexer.NewToken(lexer.COLON, nil),
			lexer.NewToken(lexer.BOOLEAN, nil),
			lexer.NewToken(lexer.ASSIGN, nil),

			lexer.NewToken(lexer.LPAREN, nil),

			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 3),
			lexer.NewToken(lexer.EQ, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 3),
			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.AND, nil),

			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 2),
			lexer.NewToken(lexer.EQ, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 2),
			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.AND, nil),

			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 1),
			lexer.NewToken(lexer.EQ, nil),
			lexer.NewToken(lexer.INTEGER_LITERAL, 0),
			lexer.NewToken(lexer.RPAREN, nil),

			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   lexer.NewToken(lexer.IDENT, "foo"),
					VariableType: lexer.NewToken(lexer.BOOLEAN, nil),
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 3},
								Operator: lexer.NewToken(lexer.EQ, nil),
								Right:    ast.NumberOpnd{Value: 3},
							},
							Operator: lexer.NewToken(lexer.AND, nil),
							Right: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 2},
								Operator: lexer.NewToken(lexer.EQ, nil),
								Right:    ast.NumberOpnd{Value: 2},
							},
						},
						Operator: lexer.NewToken(lexer.AND, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: lexer.NewToken(lexer.EQ, nil),
							Right:    ast.NumberOpnd{Value: 0},
						},
					},
				},
			},
		},
	},
	{
		name: "Logical not",
		lexerOutput: []lexer.Token{
			// print !(notTrue);
			lexer.NewToken(lexer.PRINT, nil),
			lexer.NewToken(lexer.NOT, nil),
			lexer.NewToken(lexer.LPAREN, nil),
			lexer.NewToken(lexer.IDENT, "notTrue"),
			lexer.NewToken(lexer.RPAREN, nil),
			lexer.NewToken(lexer.SEMI, nil),
		},
		expectedOutput: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.UnaryExpr{
						Unary: lexer.NewToken(lexer.NOT, nil),
						Operand: ast.NullaryExpr{
							Operand: ast.Ident{Id: lexer.NewToken(lexer.IDENT, "notTrue")},
						},
					},
				},
			},
		},
	},
}

func TestParse(t *testing.T) {
	for _, testCase := range parseTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := New(newMockLexer(
				testCase.lexerOutput...,
			))

			actual, err := parser.Parse()
			if testCase.shouldError && err == nil {
				t.Error("Expected an error, got nil")
			} else if !testCase.shouldError && err != nil {
				t.Errorf("Expected nil error, got %v", err)
			}

			if !reflect.DeepEqual(actual, testCase.expectedOutput) {
				t.Errorf("Expected:\n%+#v\ngot:\n%+#v", testCase.expectedOutput, actual)
			}
		})
	}
}

type mockLexer struct {
	tokens []lexer.Token
	pos    int
}

func newMockLexer(tokens ...lexer.Token) *mockLexer {
	return &mockLexer{tokens: tokens}
}

func (m *mockLexer) GetNextToken() lexer.Token {
	if m.pos > len(m.tokens)-1 {
		return lexer.NewToken(lexer.EOF, nil)
	}

	token := m.tokens[m.pos]
	m.pos++
	return token
}
