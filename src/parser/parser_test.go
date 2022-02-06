package parser

import (
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

var parseTestCases = []struct {
	name           string
	lexerOutput    []token.Token
	expectedOutput ast.Prog
	shouldError    bool
}{
	{
		name: "Declaration with assignment",
		lexerOutput: []token.Token{
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.INTEGER, nil),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.INTEGER_LITERAL, 5),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "x"),
					VariableType: token.NewToken(token.INTEGER, nil),
					Expression: ast.NullaryExpr{
						Operand: ast.NumberOpnd{Value: 5},
					},
				},
			},
		},
		},
	},
	{
		name: "Declaration without assignment",
		lexerOutput: []token.Token{
			// Int
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.INTEGER, nil),
			token.NewToken(token.SEMI, nil),

			// String
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "y"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.STRING, nil),
			token.NewToken(token.SEMI, nil),

			// Boolean
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "z"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.BOOLEAN, nil),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "x"),
					VariableType: token.NewToken(token.INTEGER, nil),
				},
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "y"),
					VariableType: token.NewToken(token.STRING, nil),
				},
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "z"),
					VariableType: token.NewToken(token.BOOLEAN, nil),
				},
			},
		},
		},
	},
	{
		name: "Assignment",
		lexerOutput: []token.Token{
			token.NewToken(token.IDENT, "foo"),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.STRING_LITERAL, "bar"),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssignStmt{
					Identifier: ast.Ident{Id: token.NewToken(token.IDENT, "foo")},
					Expression: ast.NullaryExpr{
						Operand: ast.StringOpnd{Value: "bar"},
					},
				},
			},
		},
		},
	},
	{
		name: "For statement with multiple inner statements",
		lexerOutput: []token.Token{
			token.NewToken(token.FOR, nil),
			token.NewToken(token.IDENT, "i"),
			token.NewToken(token.IN, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.PLUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.DOTDOT, nil),
			token.NewToken(token.INTEGER_LITERAL, 25),
			token.NewToken(token.MINUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 1),
			token.NewToken(token.DO, nil),
			token.NewToken(token.READ, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.PRINT, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.END, nil),
			token.NewToken(token.FOR, nil),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.ForStmt{
					Index: ast.Ident{Id: token.NewToken(token.IDENT, "i")},
					Low: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: token.NewToken(token.PLUS, nil),
						Right:    ast.NumberOpnd{Value: 2},
					},
					High: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 25},
						Operator: token.NewToken(token.MINUS, nil),
						Right:    ast.NumberOpnd{Value: 1},
					},
					Statements: ast.Stmts{
						Statements: []ast.Stmt{
							ast.ReadStmt{TargetIdentifier: ast.Ident{Id: token.NewToken(token.IDENT, "x")}},
							ast.PrintStmt{
								Expression: ast.NullaryExpr{
									Operand: ast.Ident{Id: token.NewToken(token.IDENT, "x")},
								},
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
		lexerOutput: []token.Token{
			token.NewToken(token.ASSERT, nil),
			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.STRING_LITERAL, "foo"),
			token.NewToken(token.EQ, nil),
			token.NewToken(token.STRING_LITERAL, "bar"),
			token.NewToken(token.RPAREN, nil),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssertStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.StringOpnd{Value: "foo"},
						Operator: token.NewToken(token.EQ, nil),
						Right:    ast.StringOpnd{Value: "bar"},
					},
				},
			},
		},
		},
	},
	{
		name: "Error if no EOF is returned by lexer when expected",
		lexerOutput: []token.Token{
			token.NewToken(token.PRINT, nil),
			token.NewToken(token.INTEGER_LITERAL, 22),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.SEMI, nil),
		},
		shouldError: true,
	},
	{
		name: "Parenthesised expressions",
		lexerOutput: []token.Token{
			// print (1 * 2) / (4 - 3)
			token.NewToken(token.PRINT, nil),

			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 1),
			token.NewToken(token.MULTIPLY, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.INTEGER_DIV, nil),

			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 4),
			token.NewToken(token.MINUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: token.NewToken(token.MULTIPLY, nil),
							Right:    ast.NumberOpnd{Value: 2},
						},
						Operator: token.NewToken(token.INTEGER_DIV, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 4},
							Operator: token.NewToken(token.MINUS, nil),
							Right:    ast.NumberOpnd{Value: 3},
						},
					},
				},
			},
		},
		},
	},
	{
		name: "Less than comparison",
		lexerOutput: []token.Token{
			// var foo : bool := 3 < 2
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "foo"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.BOOLEAN, nil),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.LT, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "foo"),
					VariableType: token.NewToken(token.BOOLEAN, nil),
					Expression: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: token.NewToken(token.LT, nil),
						Right:    ast.NumberOpnd{Value: 2},
					},
				},
			},
		},
		},
	},
	{
		name: "Nested logical and operator with not",
		lexerOutput: []token.Token{
			// var foo : bool := ((3 = 3) & (2 = 2)) & (1 = 0);
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "foo"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.BOOLEAN, nil),
			token.NewToken(token.ASSIGN, nil),

			token.NewToken(token.LPAREN, nil),

			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.EQ, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.AND, nil),

			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.EQ, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.AND, nil),

			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 1),
			token.NewToken(token.EQ, nil),
			token.NewToken(token.INTEGER_LITERAL, 0),
			token.NewToken(token.RPAREN, nil),

			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.NewToken(token.IDENT, "foo"),
					VariableType: token.NewToken(token.BOOLEAN, nil),
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 3},
								Operator: token.NewToken(token.EQ, nil),
								Right:    ast.NumberOpnd{Value: 3},
							},
							Operator: token.NewToken(token.AND, nil),
							Right: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 2},
								Operator: token.NewToken(token.EQ, nil),
								Right:    ast.NumberOpnd{Value: 2},
							},
						},
						Operator: token.NewToken(token.AND, nil),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: token.NewToken(token.EQ, nil),
							Right:    ast.NumberOpnd{Value: 0},
						},
					},
				},
			},
		},
		},
	},
	{
		name: "Logical not",
		lexerOutput: []token.Token{
			// print !(notTrue);
			token.NewToken(token.PRINT, nil),
			token.NewToken(token.NOT, nil),
			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.IDENT, "notTrue"),
			token.NewToken(token.RPAREN, nil),
			token.NewToken(token.SEMI, nil),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.UnaryExpr{
						Unary: token.NewToken(token.NOT, nil),
						Operand: ast.NullaryExpr{
							Operand: ast.Ident{Id: token.NewToken(token.IDENT, "notTrue")},
						},
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
	tokens []token.Token
	pos    int
}

func newMockLexer(tokens ...token.Token) *mockLexer {
	return &mockLexer{tokens: tokens}
}

func (m *mockLexer) GetNextToken() token.Token {
	if m.pos > len(m.tokens)-1 {
		return token.NewToken(token.EOF, nil)
	}

	token := m.tokens[m.pos]
	m.pos++
	return token
}
