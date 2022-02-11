package lexer

import (
	"testing"

	"github.com/mjjs/minipl-go/token"
)

var getNextTokenTestCases = []struct {
	name              string
	input             string
	expectedTokens    []token.Token
	expectedPositions []token.Position
	shouldPanic       bool
}{
	{
		name:              "Empty input returns EOF token",
		input:             "",
		expectedTokens:    []token.Token{token.New(token.EOF, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Integer",
		input:             "int",
		expectedTokens:    []token.Token{token.New(token.INTEGER, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "String type",
		input:             "string",
		expectedTokens:    []token.Token{token.New(token.STRING, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Boolean type",
		input:             "bool",
		expectedTokens:    []token.Token{token.New(token.BOOLEAN, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Integer literal",
		input:             "666",
		expectedTokens:    []token.Token{token.New(token.INTEGER_LITERAL, "666")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:        "Integer starting with zero",
		input:       "0123",
		shouldPanic: true,
	},
	{
		name:              "Integer with only zero",
		input:             "0",
		expectedTokens:    []token.Token{token.New(token.INTEGER_LITERAL, "0")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "String literal",
		input:             "\"C-beams\"",
		expectedTokens:    []token.Token{token.New(token.STRING_LITERAL, "C-beams")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:        "Unterminated string literal 1",
		input:       "\"Tännhauser gate\n\"",
		shouldPanic: true,
	},
	{
		name:        "Unterminated string literal 2",
		input:       "\"Tännhauser gate\r\"",
		shouldPanic: true,
	},
	{
		name:  "String literal with escaped characters",
		input: `"abc\\ def \""`,
		expectedTokens: []token.Token{
			token.New(token.STRING_LITERAL, "abc\\ def \""),
		},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:  "String literal with newlines",
		input: `"a\nb\r"`,
		expectedTokens: []token.Token{
			token.New(token.STRING_LITERAL, "a\nb\r"),
		},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:  "String literal with tab character",
		input: `"a:\tb"`,
		expectedTokens: []token.Token{
			token.New(token.STRING_LITERAL, "a:\tb"),
		},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:        "Unsupported character",
		input:       `🦊`,
		shouldPanic: true,
	},
	{
		name:  "Single ID",
		input: "myOwnVar",
		expectedTokens: []token.Token{
			token.New(token.IDENT, "myOwnVar"),
		},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Plus operator",
		input:             "+",
		expectedTokens:    []token.Token{token.New(token.PLUS, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Minus operator",
		input:             "-",
		expectedTokens:    []token.Token{token.New(token.MINUS, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Multiply operator",
		input:             "*",
		expectedTokens:    []token.Token{token.New(token.MULTIPLY, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Integer division operator",
		input:             "/",
		expectedTokens:    []token.Token{token.New(token.INTEGER_DIV, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Less than operator",
		input:             "<",
		expectedTokens:    []token.Token{token.New(token.LT, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Equality operator",
		input:             "=",
		expectedTokens:    []token.Token{token.New(token.EQ, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Logical and operator",
		input:             "&",
		expectedTokens:    []token.Token{token.New(token.AND, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Logical not operator",
		input:             "!",
		expectedTokens:    []token.Token{token.New(token.NOT, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Assign operator",
		input:             ":=",
		expectedTokens:    []token.Token{token.New(token.ASSIGN, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Left parenthesis",
		input:             "(",
		expectedTokens:    []token.Token{token.New(token.LPAREN, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Right parenthesis",
		input:             ")",
		expectedTokens:    []token.Token{token.New(token.RPAREN, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Semicolon",
		input:             ";",
		expectedTokens:    []token.Token{token.New(token.SEMI, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Colon",
		input:             ":",
		expectedTokens:    []token.Token{token.New(token.COLON, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "For keyword",
		input:             "for",
		expectedTokens:    []token.Token{token.New(token.FOR, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "In keyword",
		input:             "in",
		expectedTokens:    []token.Token{token.New(token.IN, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Do keyword",
		input:             "do",
		expectedTokens:    []token.Token{token.New(token.DO, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "End keyword",
		input:             "end",
		expectedTokens:    []token.Token{token.New(token.END, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Double dot keyword",
		input:             "..",
		expectedTokens:    []token.Token{token.New(token.DOTDOT, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Assert keyword",
		input:             "assert",
		expectedTokens:    []token.Token{token.New(token.ASSERT, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Var keyword",
		input:             "var",
		expectedTokens:    []token.Token{token.New(token.VAR, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Read keyword",
		input:             "read",
		expectedTokens:    []token.Token{token.New(token.READ, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:              "Print keyword",
		input:             "print",
		expectedTokens:    []token.Token{token.New(token.PRINT, "")},
		expectedPositions: []token.Position{{Line: 1, Column: 1}},
	},
	{
		name:  "Whitespace is skipped",
		input: "first   \n    second",
		expectedTokens: []token.Token{
			token.New(token.IDENT, "first"),
			token.New(token.IDENT, "second"),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 2, Column: 5},
		},
	},
	{
		name:  "For loop",
		input: "for i in 1..n do\n\tv := v * i;\nend for;",
		expectedTokens: []token.Token{
			token.New(token.FOR, ""),
			token.New(token.IDENT, "i"),
			token.New(token.IN, ""),
			token.New(token.INTEGER_LITERAL, "1"),
			token.New(token.DOTDOT, ""),
			token.New(token.IDENT, "n"),
			token.New(token.DO, ""),

			token.New(token.IDENT, "v"),
			token.New(token.ASSIGN, ""),
			token.New(token.IDENT, "v"),
			token.New(token.MULTIPLY, ""),
			token.New(token.IDENT, "i"),
			token.New(token.SEMI, ""),

			token.New(token.END, ""),
			token.New(token.FOR, ""),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 5},
			{Line: 1, Column: 7},
			{Line: 1, Column: 10},
			{Line: 1, Column: 11},
			{Line: 1, Column: 13},
			{Line: 1, Column: 15},

			{Line: 2, Column: 2},
			{Line: 2, Column: 4},
			{Line: 2, Column: 7},
			{Line: 2, Column: 9},
			{Line: 2, Column: 11},
			{Line: 2, Column: 12},

			{Line: 3, Column: 1},
			{Line: 3, Column: 5},
			{Line: 3, Column: 8},
		},
	},
	{
		name:  "Parentheses tokenize correctly",
		input: "(52)",
		expectedTokens: []token.Token{
			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "52"),
			token.New(token.RPAREN, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 2},
			{Line: 1, Column: 4},
		},
	},
	{
		name:  "Declaration with assignment",
		input: "var x : int := 4;",
		expectedTokens: []token.Token{
			token.New(token.VAR, ""),
			token.New(token.IDENT, "x"),
			token.New(token.COLON, ""),
			token.New(token.INTEGER, ""),
			token.New(token.ASSIGN, ""),
			token.New(token.INTEGER_LITERAL, "4"),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 5},
			{Line: 1, Column: 7},
			{Line: 1, Column: 9},
			{Line: 1, Column: 13},
			{Line: 1, Column: 16},
			{Line: 1, Column: 17},
		},
	},
	{
		name:  "Declaration and assignment in different lines",
		input: "var firstLine : string;\nfirstLine := \"secondLine\";",
		expectedTokens: []token.Token{
			token.New(token.VAR, ""),
			token.New(token.IDENT, "firstLine"),
			token.New(token.COLON, ""),
			token.New(token.STRING, ""),
			token.New(token.SEMI, ""),

			token.New(token.IDENT, "firstLine"),
			token.New(token.ASSIGN, ""),
			token.New(token.STRING_LITERAL, "secondLine"),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 5},
			{Line: 1, Column: 15},
			{Line: 1, Column: 17},
			{Line: 1, Column: 23},

			{Line: 2, Column: 1},
			{Line: 2, Column: 11},
			{Line: 2, Column: 14},
			{Line: 2, Column: 26},
		},
	},
	{
		name:  "Assert with weird formatting",
		input: "            assert    (x     =         somethingElse)   ;",
		expectedTokens: []token.Token{
			token.New(token.ASSERT, ""),
			token.New(token.LPAREN, ""),
			token.New(token.IDENT, "x"),
			token.New(token.EQ, ""),
			token.New(token.IDENT, "somethingElse"),
			token.New(token.RPAREN, ""),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 13},
			{Line: 1, Column: 23},
			{Line: 1, Column: 24},
			{Line: 1, Column: 30},
			{Line: 1, Column: 40},
			{Line: 1, Column: 53},
			{Line: 1, Column: 57},
		},
	},
	{
		name:  "Skips line comments",
		input: "3 + 3; // this is a comment\n4 + 4;",
		expectedTokens: []token.Token{
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.PLUS, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.SEMI, ""),

			token.New(token.INTEGER_LITERAL, "4"),
			token.New(token.PLUS, ""),
			token.New(token.INTEGER_LITERAL, "4"),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 3},
			{Line: 1, Column: 5},
			{Line: 1, Column: 6},

			{Line: 2, Column: 1},
			{Line: 2, Column: 3},
			{Line: 2, Column: 5},
			{Line: 2, Column: 6},
		},
	},
	{
		name:  "Skips block comments",
		input: "1 + 2; /*This is a comment*/ 3 + 4;",
		expectedTokens: []token.Token{
			token.New(token.INTEGER_LITERAL, "1"),
			token.New(token.PLUS, ""),
			token.New(token.INTEGER_LITERAL, "2"),
			token.New(token.SEMI, ""),
			token.New(token.INTEGER_LITERAL, "3"),
			token.New(token.PLUS, ""),
			token.New(token.INTEGER_LITERAL, "4"),
			token.New(token.SEMI, ""),
		},
		expectedPositions: []token.Position{
			{Line: 1, Column: 1},
			{Line: 1, Column: 3},
			{Line: 1, Column: 5},
			{Line: 1, Column: 6},
			{Line: 1, Column: 30},
			{Line: 1, Column: 32},
			{Line: 1, Column: 34},
			{Line: 1, Column: 35},
		},
	},
}

func TestGetNextToken(t *testing.T) {
	for _, testCase := range getNextTokenTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil && testCase.shouldPanic {
					t.Error("Expected a panic")
				} else if r != nil && !testCase.shouldPanic {
					t.Errorf("Did not expect a panic, got '%v'", r)
				}
			}()

			lexer := New(testCase.input)

			if len(testCase.expectedTokens) == 0 {
				for tok, _ := lexer.GetNextToken(); tok != token.New(token.EOF, ""); {
				}
			}

			for i, expectedToken := range testCase.expectedTokens {
				actualToken, actualPos := lexer.GetNextToken()

				if actualToken != expectedToken {
					t.Errorf("Expected %v, got %v", expectedToken, actualToken)
				}

				if expectedPos := testCase.expectedPositions[i]; actualPos != expectedPos {
					t.Errorf("Expected %v, got %v", expectedPos, actualPos)
				}
			}
		})
	}
}
