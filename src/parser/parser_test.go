package parser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

var parseTestCases = []struct {
	name           string
	lexerOutput    []positionedToken
	expectedAST    ast.Prog
	expectedErrors []error
}{
	{
		name: "Declaration with assignment",
		lexerOutput: []positionedToken{
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "x"), token.Position{Line: 1, Column: 2}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.INTEGER, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.ASSIGN, ""), token.Position{Line: 1, Column: 5}},
			{token.New(token.INTEGER_LITERAL, "5"), token.Position{Line: 1, Column: 6}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 7}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
						Expression: ast.NullaryExpr{
							Operand: ast.NumberOpnd{Value: 5, Pos: token.Position{Line: 1, Column: 6}},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Declaration without assignment",
		lexerOutput: []positionedToken{
			// Int
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "x"), token.Position{Line: 1, Column: 2}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.INTEGER, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 5}},

			// String
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 6}},
			{token.New(token.IDENT, "y"), token.Position{Line: 1, Column: 7}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 8}},
			{token.New(token.STRING, ""), token.Position{Line: 1, Column: 9}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 10}},

			// Boolean
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 11}},
			{token.New(token.IDENT, "z"), token.Position{Line: 1, Column: 12}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 13}},
			{token.New(token.BOOLEAN, ""), token.Position{Line: 1, Column: 14}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 15}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "x"),
						VariableType: token.New(token.INTEGER, ""),
						Pos:          token.Position{Line: 1, Column: 1},
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "y"),
						VariableType: token.New(token.STRING, ""),
						Pos:          token.Position{Line: 1, Column: 6},
					},
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "z"),
						VariableType: token.New(token.BOOLEAN, ""),
						Pos:          token.Position{Line: 1, Column: 11},
					},
				},
			},
		},
	},
	{
		name: "Assignment",
		lexerOutput: []positionedToken{
			{token.New(token.IDENT, "foo"), token.Position{Line: 1, Column: 1}},
			{token.New(token.ASSIGN, ""), token.Position{Line: 1, Column: 2}},
			{token.New(token.STRING_LITERAL, "bar"), token.Position{Line: 1, Column: 3}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 4}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssignStmt{
						Identifier: ast.Ident{
							Id:  token.New(token.IDENT, "foo"),
							Pos: token.Position{Line: 1, Column: 1},
						},
						Expression: ast.NullaryExpr{
							Operand: ast.StringOpnd{
								Value: "bar",
								Pos:   token.Position{Line: 1, Column: 3},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "For statement with multiple inner statements",
		lexerOutput: []positionedToken{
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "i"), token.Position{Line: 1, Column: 2}},
			{token.New(token.IN, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.INTEGER_LITERAL, "3"), token.Position{Line: 1, Column: 4}},
			{token.New(token.PLUS, ""), token.Position{Line: 1, Column: 5}},
			{token.New(token.INTEGER_LITERAL, "2"), token.Position{Line: 1, Column: 6}},
			{token.New(token.DOTDOT, ""), token.Position{Line: 1, Column: 7}},
			{token.New(token.INTEGER_LITERAL, "25"), token.Position{Line: 1, Column: 8}},
			{token.New(token.MINUS, ""), token.Position{Line: 1, Column: 9}},
			{token.New(token.INTEGER_LITERAL, "1"), token.Position{Line: 1, Column: 10}},
			{token.New(token.DO, ""), token.Position{Line: 1, Column: 11}},
			{token.New(token.READ, ""), token.Position{Line: 1, Column: 12}},
			{token.New(token.IDENT, "x"), token.Position{Line: 1, Column: 13}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 14}},
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 15}},
			{token.New(token.IDENT, "x"), token.Position{Line: 1, Column: 16}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 17}},
			{token.New(token.END, ""), token.Position{Line: 1, Column: 18}},
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 19}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 20}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.ForStmt{
						Index: ast.Ident{
							Id:  token.New(token.IDENT, "i"),
							Pos: token.Position{Line: 1, Column: 2},
						},
						Low: ast.BinaryExpr{
							Left: ast.NumberOpnd{
								Value: 3,
								Pos:   token.Position{Line: 1, Column: 4},
							},
							Operator: token.New(token.PLUS, ""),
							Right: ast.NumberOpnd{
								Value: 2,
								Pos:   token.Position{Line: 1, Column: 6},
							},
						},
						High: ast.BinaryExpr{
							Left: ast.NumberOpnd{
								Value: 25,
								Pos:   token.Position{Line: 1, Column: 8},
							},
							Operator: token.New(token.MINUS, ""),
							Right: ast.NumberOpnd{
								Value: 1,
								Pos:   token.Position{Line: 1, Column: 10},
							},
						},
						Statements: ast.Stmts{
							Statements: []ast.Stmt{
								ast.ReadStmt{
									TargetIdentifier: ast.Ident{
										Id:  token.New(token.IDENT, "x"),
										Pos: token.Position{Line: 1, Column: 13},
									},
									Pos: token.Position{Line: 1, Column: 12},
								},
								ast.PrintStmt{
									Expression: ast.NullaryExpr{
										Operand: ast.Ident{
											Id:  token.New(token.IDENT, "x"),
											Pos: token.Position{Line: 1, Column: 16},
										},
									},
									Pos: token.Position{Line: 1, Column: 15},
								},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Assert statement",
		lexerOutput: []positionedToken{
			{token.New(token.ASSERT, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 2}},
			{token.New(token.STRING_LITERAL, "foo"), token.Position{Line: 1, Column: 3}},
			{token.New(token.EQ, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.STRING_LITERAL, "bar"), token.Position{Line: 1, Column: 5}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 6}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 7}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.AssertStmt{
						Expression: ast.BinaryExpr{
							Left: ast.StringOpnd{
								Value: "foo",
								Pos:   token.Position{Line: 1, Column: 3},
							},
							Operator: token.New(token.EQ, ""),
							Right: ast.StringOpnd{
								Value: "bar",
								Pos:   token.Position{Line: 1, Column: 5},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Parenthesised expressions",
		lexerOutput: []positionedToken{
			// print (1 * 2) / (4 - 3)
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 1}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 2}},
			{token.New(token.INTEGER_LITERAL, "1"), token.Position{Line: 1, Column: 3}},
			{token.New(token.MULTIPLY, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.INTEGER_LITERAL, "2"), token.Position{Line: 1, Column: 5}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 6}},

			{token.New(token.INTEGER_DIV, ""), token.Position{Line: 1, Column: 7}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 8}},
			{token.New(token.INTEGER_LITERAL, "4"), token.Position{Line: 1, Column: 9}},
			{token.New(token.MINUS, ""), token.Position{Line: 1, Column: 10}},
			{token.New(token.INTEGER_LITERAL, "3"), token.Position{Line: 1, Column: 11}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 12}},

			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 13}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.NumberOpnd{
									Value: 1,
									Pos:   token.Position{Line: 1, Column: 3},
								},
								Operator: token.New(token.MULTIPLY, ""),
								Right: ast.NumberOpnd{
									Value: 2,
									Pos:   token.Position{Line: 1, Column: 5},
								},
							},
							Operator: token.New(token.INTEGER_DIV, ""),
							Right: ast.BinaryExpr{
								Left: ast.NumberOpnd{
									Value: 4,
									Pos:   token.Position{Line: 1, Column: 9},
								},
								Operator: token.New(token.MINUS, ""),
								Right: ast.NumberOpnd{
									Value: 3,
									Pos:   token.Position{Line: 1, Column: 11},
								},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Less than comparison",
		lexerOutput: []positionedToken{
			// var foo : bool := 3 < 2
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "foo"), token.Position{Line: 1, Column: 2}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.BOOLEAN, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.ASSIGN, ""), token.Position{Line: 1, Column: 5}},
			{token.New(token.INTEGER_LITERAL, "3"), token.Position{Line: 1, Column: 6}},
			{token.New(token.LT, ""), token.Position{Line: 1, Column: 7}},
			{token.New(token.INTEGER_LITERAL, "2"), token.Position{Line: 1, Column: 8}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 9}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.BinaryExpr{
							Left: ast.NumberOpnd{
								Value: 3,
								Pos:   token.Position{Line: 1, Column: 6},
							},
							Operator: token.New(token.LT, ""),
							Right: ast.NumberOpnd{
								Value: 2,
								Pos:   token.Position{Line: 1, Column: 8},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Nested logical and operator with not",
		lexerOutput: []positionedToken{
			// var foo : bool := ((3 = 3) & (2 = 2)) & (1 = 0);
			{token.New(token.VAR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "foo"), token.Position{Line: 1, Column: 2}},
			{token.New(token.COLON, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.BOOLEAN, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.ASSIGN, ""), token.Position{Line: 1, Column: 5}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 6}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 7}},
			{token.New(token.INTEGER_LITERAL, "3"), token.Position{Line: 1, Column: 8}},
			{token.New(token.EQ, ""), token.Position{Line: 1, Column: 9}},
			{token.New(token.INTEGER_LITERAL, "3"), token.Position{Line: 1, Column: 10}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 11}},

			{token.New(token.AND, ""), token.Position{Line: 1, Column: 12}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 13}},
			{token.New(token.INTEGER_LITERAL, "2"), token.Position{Line: 1, Column: 14}},
			{token.New(token.EQ, ""), token.Position{Line: 1, Column: 15}},
			{token.New(token.INTEGER_LITERAL, "2"), token.Position{Line: 1, Column: 16}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 17}},

			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 18}},

			{token.New(token.AND, ""), token.Position{Line: 1, Column: 19}},

			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 20}},
			{token.New(token.INTEGER_LITERAL, "1"), token.Position{Line: 1, Column: 21}},
			{token.New(token.EQ, ""), token.Position{Line: 1, Column: 22}},
			{token.New(token.INTEGER_LITERAL, "0"), token.Position{Line: 1, Column: 23}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 24}},

			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 25}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.DeclStmt{
						Identifier:   token.New(token.IDENT, "foo"),
						VariableType: token.New(token.BOOLEAN, ""),
						Expression: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.BinaryExpr{
									Left: ast.NumberOpnd{
										Value: 3,
										Pos:   token.Position{Line: 1, Column: 8},
									},
									Operator: token.New(token.EQ, ""),
									Right: ast.NumberOpnd{
										Value: 3,
										Pos:   token.Position{Line: 1, Column: 10},
									},
								},
								Operator: token.New(token.AND, ""),
								Right: ast.BinaryExpr{
									Left: ast.NumberOpnd{
										Value: 2,
										Pos:   token.Position{Line: 1, Column: 14},
									},
									Operator: token.New(token.EQ, ""),
									Right: ast.NumberOpnd{
										Value: 2,
										Pos:   token.Position{Line: 1, Column: 16},
									},
								},
							},
							Operator: token.New(token.AND, ""),
							Right: ast.BinaryExpr{
								Left: ast.NumberOpnd{
									Value: 1,
									Pos:   token.Position{Line: 1, Column: 21},
								},
								Operator: token.New(token.EQ, ""),
								Right: ast.NumberOpnd{
									Value: 0,
									Pos:   token.Position{Line: 1, Column: 23},
								},
							},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	{
		name: "Logical not",
		lexerOutput: []positionedToken{
			// print !(notTrue);
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.NOT, ""), token.Position{Line: 1, Column: 2}},
			{token.New(token.LPAREN, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.IDENT, "notTrue"), token.Position{Line: 1, Column: 4}},
			{token.New(token.RPAREN, ""), token.Position{Line: 1, Column: 5}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 6}},
		},
		expectedAST: ast.Prog{
			Statements: ast.Stmts{
				Statements: []ast.Stmt{
					ast.PrintStmt{
						Expression: ast.UnaryExpr{
							Unary: token.New(token.NOT, ""),
							Operand: ast.NullaryExpr{
								Operand: ast.Ident{
									Id:  token.New(token.IDENT, "notTrue"),
									Pos: token.Position{Line: 1, Column: 4},
								},
							},
							Pos: token.Position{Line: 1, Column: 2},
						},
						Pos: token.Position{Line: 1, Column: 1},
					},
				},
			},
		},
	},
	// ERRORS
	{
		name: "Error if no EOF is returned by lexer when expected",
		lexerOutput: []positionedToken{
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.INTEGER_LITERAL, "22"), token.Position{Line: 1, Column: 2}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 4}},
		},
		expectedErrors: []error{
			errors.New("1:4: syntax error: expected EOF got SEMI"),
		},
	},
	{
		name: "Invalid for loop control block",
		lexerOutput: []positionedToken{
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "i"), token.Position{Line: 1, Column: 2}},
			{token.New(token.IN, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.DOTDOT, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.INTEGER_LITERAL, "25"), token.Position{Line: 1, Column: 5}},
			{token.New(token.DO, ""), token.Position{Line: 1, Column: 6}},
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 7}},
			{token.New(token.IDENT, "i"), token.Position{Line: 1, Column: 8}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 9}},
			{token.New(token.END, ""), token.Position{Line: 1, Column: 10}},
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 11}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 12}},
		},
		expectedErrors: []error{
			errors.New("1:4: syntax error: unexpected DOTDOT"),
		},
	},
	{
		name: "Invalid statements in for loop",
		lexerOutput: []positionedToken{
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 1}},
			{token.New(token.IDENT, "i"), token.Position{Line: 1, Column: 2}},
			{token.New(token.IN, ""), token.Position{Line: 1, Column: 3}},
			{token.New(token.INTEGER_LITERAL, "1"), token.Position{Line: 1, Column: 5}},
			{token.New(token.DOTDOT, ""), token.Position{Line: 1, Column: 4}},
			{token.New(token.INTEGER_LITERAL, "25"), token.Position{Line: 1, Column: 5}},
			{token.New(token.DO, ""), token.Position{Line: 1, Column: 6}},
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 7}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 9}},
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 22}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 23}},
			{token.New(token.PRINT, ""), token.Position{Line: 1, Column: 44}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 45}},
			{token.New(token.END, ""), token.Position{Line: 1, Column: 10}},
			{token.New(token.FOR, ""), token.Position{Line: 1, Column: 11}},
			{token.New(token.SEMI, ""), token.Position{Line: 1, Column: 12}},
		},
		expectedErrors: []error{
			errors.New("1:9: syntax error: unexpected SEMI"),
			errors.New("1:23: syntax error: unexpected SEMI"),
			errors.New("1:45: syntax error: unexpected SEMI"),
		},
	},
	{
		name:        "Error when no statements are present",
		lexerOutput: []positionedToken{},
		expectedErrors: []error{
			errors.New("99:99: syntax error: unexpected EOF"),
		},
	},
}

func TestParsePanicsWithNilLexer(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected a panic")
		}
	}()

	New(nil)
}

func TestParse(t *testing.T) {
	for _, testCase := range parseTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := New(newMockLexer(
				testCase.lexerOutput,
			))

			actual, errors := parser.Parse()
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
			} else if !reflect.DeepEqual(actual, testCase.expectedAST) {
				t.Errorf("Expected:\n%+#v\ngot:\n%+#v", testCase.expectedAST, actual)
			}

		})
	}
}

type mockLexer struct {
	tokens    []token.Token
	positions []token.Position
	pos       int
}

// Helper struct for simpler test case construction
type positionedToken struct {
	token token.Token
	pos   token.Position
}

func newMockLexer(output []positionedToken) *mockLexer {
	tokens := []token.Token{}
	positions := []token.Position{}

	for _, x := range output {
		tokens = append(tokens, x.token)
		positions = append(positions, x.pos)
	}

	return &mockLexer{tokens: tokens, positions: positions}
}

func (m *mockLexer) GetNextToken() (token.Token, token.Position) {
	if m.pos > len(m.tokens)-1 {
		return token.New(token.EOF, ""), token.Position{Line: 99, Column: 99}
	}

	token := m.tokens[m.pos]
	pos := m.positions[m.pos]
	m.pos++

	return token, pos
}
