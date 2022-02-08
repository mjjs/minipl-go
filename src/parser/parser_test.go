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
			token.New(token.VAR, ""),
			token.New(token.IDENT, "x"),
			token.New(token.COLON, ""),
			token.New(token.INTEGER, ""),
			token.New(token.ASSIGN, ""),
			token.New(token.INTEGER_LITERAL, "5"),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "x"),
					VariableType: token.New(token.INTEGER, ""),
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
			token.New(token.VAR, ""),
			token.New(token.IDENT, "x"),
			token.New(token.COLON, ""),
			token.New(token.INTEGER, ""),
			token.New(token.SEMI, ""),

			// String
			token.New(token.VAR, ""),
			token.New(token.IDENT, "y"),
			token.New(token.COLON, ""),
			token.New(token.STRING, ""),
			token.New(token.SEMI, ""),

			// Boolean
			token.New(token.VAR, ""),
			token.New(token.IDENT, "z"),
			token.New(token.COLON, ""),
			token.New(token.BOOLEAN, ""),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "x"),
					VariableType: token.New(token.INTEGER, ""),
				},
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "y"),
					VariableType: token.New(token.STRING, ""),
				},
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "z"),
					VariableType: token.New(token.BOOLEAN, ""),
				},
			},
		},
		},
	},
	{
		name: "Assignment",
		lexerOutput: []token.Token{
			token.New(token.IDENT, "foo"),
			token.New(token.ASSIGN, ""),
			token.New(token.STRING_LITERAL, "bar"),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssignStmt{
					Identifier: ast.Ident{Id: token.New(token.IDENT, "foo")},
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
			token.New(token.FOR, ""),
			token.New(token.IDENT, "i"),
			token.New(token.IN, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.PLUS, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.DOTDOT, ""),
			token.New(token.INTEGER_LITERAL, "25"),
			token.New(token.MINUS, ""),
			token.New(token.INTEGER_LITERAL, "1"),
			token.New(token.DO, ""),
			token.New(token.READ, ""),
			token.New(token.IDENT, "x"),
			token.New(token.SEMI, ""),
			token.New(token.PRINT, ""),
			token.New(token.IDENT, "x"),
			token.New(token.SEMI, ""),
			token.New(token.END, ""),
			token.New(token.FOR, ""),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.ForStmt{
					Index: ast.Ident{Id: token.New(token.IDENT, "i")},
					Low: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: token.New(token.PLUS, ""),
						Right:    ast.NumberOpnd{Value: 2},
					},
					High: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 25},
						Operator: token.New(token.MINUS, ""),
						Right:    ast.NumberOpnd{Value: 1},
					},
					Statements: ast.Stmts{
						Statements: []ast.Stmt{
							ast.ReadStmt{TargetIdentifier: ast.Ident{Id: token.New(token.IDENT, "x")}},
							ast.PrintStmt{
								Expression: ast.NullaryExpr{
									Operand: ast.Ident{Id: token.New(token.IDENT, "x")},
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
			token.New(token.ASSERT, ""),
			token.New(token.LPAREN, ""),
			token.New(token.STRING_LITERAL, "foo"),
			token.New(token.EQ, ""),
			token.New(token.STRING_LITERAL, "bar"),
			token.New(token.RPAREN, ""),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.AssertStmt{
					Expression: ast.BinaryExpr{
						Left:     ast.StringOpnd{Value: "foo"},
						Operator: token.New(token.EQ, ""),
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
			token.New(token.PRINT, ""),
			token.New(token.INTEGER_LITERAL, "22"),
			token.New(token.SEMI, ""),
			token.New(token.SEMI, ""),
		},
		shouldError: true,
	},
	{
		name: "Parenthesised expressions",
		lexerOutput: []token.Token{
			// print (1 * 2) / (4 - 3)
			token.New(token.PRINT, ""),

			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "1"),
			token.New(token.MULTIPLY, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.RPAREN, ""),

			token.New(token.INTEGER_DIV, ""),

			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "4"),
			token.New(token.MINUS, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.RPAREN, ""),

			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: token.New(token.MULTIPLY, ""),
							Right:    ast.NumberOpnd{Value: 2},
						},
						Operator: token.New(token.INTEGER_DIV, ""),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 4},
							Operator: token.New(token.MINUS, ""),
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
			token.New(token.VAR, ""),
			token.New(token.IDENT, "foo"),
			token.New(token.COLON, ""),
			token.New(token.BOOLEAN, ""),
			token.New(token.ASSIGN, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.LT, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "foo"),
					VariableType: token.New(token.BOOLEAN, ""),
					Expression: ast.BinaryExpr{
						Left:     ast.NumberOpnd{Value: 3},
						Operator: token.New(token.LT, ""),
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
			token.New(token.VAR, ""),
			token.New(token.IDENT, "foo"),
			token.New(token.COLON, ""),
			token.New(token.BOOLEAN, ""),
			token.New(token.ASSIGN, ""),

			token.New(token.LPAREN, ""),

			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.EQ, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.RPAREN, ""),

			token.New(token.AND, ""),

			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.EQ, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.RPAREN, ""),

			token.New(token.RPAREN, ""),

			token.New(token.AND, ""),

			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "1"),
			token.New(token.EQ, ""),
			token.New(token.INTEGER_LITERAL, "0"),
			token.New(token.RPAREN, ""),

			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.DeclStmt{
					Identifier:   token.New(token.IDENT, "foo"),
					VariableType: token.New(token.BOOLEAN, ""),
					Expression: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 3},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NumberOpnd{Value: 3},
							},
							Operator: token.New(token.AND, ""),
							Right: ast.BinaryExpr{
								Left:     ast.NumberOpnd{Value: 2},
								Operator: token.New(token.EQ, ""),
								Right:    ast.NumberOpnd{Value: 2},
							},
						},
						Operator: token.New(token.AND, ""),
						Right: ast.BinaryExpr{
							Left:     ast.NumberOpnd{Value: 1},
							Operator: token.New(token.EQ, ""),
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
			token.New(token.PRINT, ""),
			token.New(token.NOT, ""),
			token.New(token.LPAREN, ""),
			token.New(token.IDENT, "notTrue"),
			token.New(token.RPAREN, ""),
			token.New(token.SEMI, ""),
		},
		expectedOutput: ast.Prog{Statements: ast.Stmts{
			Statements: []ast.Stmt{
				ast.PrintStmt{
					Expression: ast.UnaryExpr{
						Unary: token.New(token.NOT, ""),
						Operand: ast.NullaryExpr{
							Operand: ast.Ident{Id: token.New(token.IDENT, "notTrue")},
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
		return token.New(token.EOF, "")
	}

	token := m.tokens[m.pos]
	m.pos++
	return token
}
