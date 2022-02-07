package lexer

import (
	"testing"

	"github.com/mjjs/minipl-go/token"
)

var getNextTokenTestCases = []struct {
	name           string
	input          string
	expectedTokens []token.Token
	shouldPanic    bool
}{
	{
		name:           "Empty input returns EOF token",
		input:          "",
		expectedTokens: []token.Token{token.NewToken(token.EOF, nil)},
	},
	{
		name:           "Integer",
		input:          "int",
		expectedTokens: []token.Token{token.NewToken(token.INTEGER, nil)},
	},
	{
		name:           "String type",
		input:          "string",
		expectedTokens: []token.Token{token.NewToken(token.STRING, nil)},
	},
	{
		name:           "Boolean type",
		input:          "bool",
		expectedTokens: []token.Token{token.NewToken(token.BOOLEAN, nil)},
	},
	{
		name:           "Integer literal",
		input:          "666",
		expectedTokens: []token.Token{token.NewToken(token.INTEGER_LITERAL, 666)},
	},
	{
		name:        "Integer starting with zero",
		input:       "0123",
		shouldPanic: true,
	},
	{
		name:           "Integer with only zero",
		input:          "0",
		expectedTokens: []token.Token{token.NewToken(token.INTEGER_LITERAL, 0)},
	},
	{
		name:           "String literal",
		input:          "\"C-beams\"",
		expectedTokens: []token.Token{token.NewToken(token.STRING_LITERAL, "C-beams")},
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
			token.NewToken(token.STRING_LITERAL, "abc\\ def \""),
		},
	},
	{
		name:  "String literal with newlines",
		input: `"a\nb\r"`,
		expectedTokens: []token.Token{
			token.NewToken(token.STRING_LITERAL, "a\nb\r"),
		},
	},
	{
		name:  "String literal with tab character",
		input: `"a:\tb"`,
		expectedTokens: []token.Token{
			token.NewToken(token.STRING_LITERAL, "a:\tb"),
		},
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
			token.NewToken(token.IDENT, "myOwnVar"),
		},
	},
	{
		name:           "Plus operator",
		input:          "+",
		expectedTokens: []token.Token{token.NewToken(token.PLUS, nil)},
	},
	{
		name:           "Minus operator",
		input:          "-",
		expectedTokens: []token.Token{token.NewToken(token.MINUS, nil)},
	},
	{
		name:           "Multiply operator",
		input:          "*",
		expectedTokens: []token.Token{token.NewToken(token.MULTIPLY, nil)},
	},
	{
		name:           "Integer division operator",
		input:          "/",
		expectedTokens: []token.Token{token.NewToken(token.INTEGER_DIV, nil)},
	},
	{
		name:           "Less than operator",
		input:          "<",
		expectedTokens: []token.Token{token.NewToken(token.LT, nil)},
	},
	{
		name:           "Equality operator",
		input:          "=",
		expectedTokens: []token.Token{token.NewToken(token.EQ, nil)},
	},
	{
		name:           "Logical and operator",
		input:          "&",
		expectedTokens: []token.Token{token.NewToken(token.AND, nil)},
	},
	{
		name:           "Logical not operator",
		input:          "!",
		expectedTokens: []token.Token{token.NewToken(token.NOT, nil)},
	},
	{
		name:           "Assign operator",
		input:          ":=",
		expectedTokens: []token.Token{token.NewToken(token.ASSIGN, nil)},
	},
	{
		name:           "Left parenthesis",
		input:          "(",
		expectedTokens: []token.Token{token.NewToken(token.LPAREN, nil)},
	},
	{
		name:           "Right parenthesis",
		input:          ")",
		expectedTokens: []token.Token{token.NewToken(token.RPAREN, nil)},
	},
	{
		name:           "Semicolon",
		input:          ";",
		expectedTokens: []token.Token{token.NewToken(token.SEMI, nil)},
	},
	{
		name:           "Colon",
		input:          ":",
		expectedTokens: []token.Token{token.NewToken(token.COLON, nil)},
	},
	{
		name:           "For keyword",
		input:          "for",
		expectedTokens: []token.Token{token.NewToken(token.FOR, nil)},
	},
	{
		name:           "In keyword",
		input:          "in",
		expectedTokens: []token.Token{token.NewToken(token.IN, nil)},
	},
	{
		name:           "Do keyword",
		input:          "do",
		expectedTokens: []token.Token{token.NewToken(token.DO, nil)},
	},
	{
		name:           "End keyword",
		input:          "end",
		expectedTokens: []token.Token{token.NewToken(token.END, nil)},
	},
	{
		name:           "Double dot keyword",
		input:          "..",
		expectedTokens: []token.Token{token.NewToken(token.DOTDOT, nil)},
	},
	{
		name:           "Assert keyword",
		input:          "assert",
		expectedTokens: []token.Token{token.NewToken(token.ASSERT, nil)},
	},
	{
		name:           "Var keyword",
		input:          "var",
		expectedTokens: []token.Token{token.NewToken(token.VAR, nil)},
	},
	{
		name:           "Read keyword",
		input:          "read",
		expectedTokens: []token.Token{token.NewToken(token.READ, nil)},
	},
	{
		name:           "Print keyword",
		input:          "print",
		expectedTokens: []token.Token{token.NewToken(token.PRINT, nil)},
	},
	{
		name:  "Whitespace is skipped",
		input: "first   \n    second",
		expectedTokens: []token.Token{
			token.NewToken(token.IDENT, "first"),
			token.NewToken(token.IDENT, "second"),
		},
	},
	{
		name:  "For loop",
		input: "for i in 1..n do\n\tv := v * i;\nend for;",
		expectedTokens: []token.Token{
			token.NewToken(token.FOR, nil),
			token.NewToken(token.IDENT, "i"),
			token.NewToken(token.IN, nil),
			token.NewToken(token.INTEGER_LITERAL, 1),
			token.NewToken(token.DOTDOT, nil),
			token.NewToken(token.IDENT, "n"),
			token.NewToken(token.DO, nil),
			token.NewToken(token.IDENT, "v"),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.IDENT, "v"),
			token.NewToken(token.MULTIPLY, nil),
			token.NewToken(token.IDENT, "i"),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.END, nil),
			token.NewToken(token.FOR, nil),
			token.NewToken(token.SEMI, nil),
		},
	},
	{
		name:  "Parentheses tokenize correctly",
		input: "(52)",
		expectedTokens: []token.Token{
			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.INTEGER_LITERAL, 52),
			token.NewToken(token.RPAREN, nil),
		},
	},
	{
		name:  "Declaration with assignment",
		input: "var x : int := 4;",
		expectedTokens: []token.Token{
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.INTEGER, nil),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.INTEGER_LITERAL, 4),
			token.NewToken(token.SEMI, nil),
		},
	},
	{
		name:  "Declaration and assignment in different lines",
		input: "var firstLine : string;\nfirstLine := \"secondLine\";",
		expectedTokens: []token.Token{
			token.NewToken(token.VAR, nil),
			token.NewToken(token.IDENT, "firstLine"),
			token.NewToken(token.COLON, nil),
			token.NewToken(token.STRING, nil),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.IDENT, "firstLine"),
			token.NewToken(token.ASSIGN, nil),
			token.NewToken(token.STRING_LITERAL, "secondLine"),
			token.NewToken(token.SEMI, nil),
		},
	},
	{
		name:  "Assert with weird formatting",
		input: "            assert    (x     =         somethingElse)   ;",
		expectedTokens: []token.Token{
			token.NewToken(token.ASSERT, nil),
			token.NewToken(token.LPAREN, nil),
			token.NewToken(token.IDENT, "x"),
			token.NewToken(token.EQ, nil),
			token.NewToken(token.IDENT, "somethingElse"),
			token.NewToken(token.RPAREN, nil),
			token.NewToken(token.SEMI, nil),
		},
	},
	{
		name:  "Skips line comments",
		input: "3 + 3; // this is a comment\n4 + 4;",
		expectedTokens: []token.Token{
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.PLUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.INTEGER_LITERAL, 4),
			token.NewToken(token.PLUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 4),
			token.NewToken(token.SEMI, nil),
		},
	},
	{
		name:  "Skips block comments",
		input: "1 + 2; /*This is a comment*/ 3 + 4;",
		expectedTokens: []token.Token{
			token.NewToken(token.INTEGER_LITERAL, 1),
			token.NewToken(token.PLUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 2),
			token.NewToken(token.SEMI, nil),
			token.NewToken(token.INTEGER_LITERAL, 3),
			token.NewToken(token.PLUS, nil),
			token.NewToken(token.INTEGER_LITERAL, 4),
			token.NewToken(token.SEMI, nil),
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
				for lexer.GetNextToken() != token.NewToken(token.EOF, nil) {
				}
			}

			for _, expectedToken := range testCase.expectedTokens {
				actual := lexer.GetNextToken()

				if actual != expectedToken {
					t.Errorf("Expected %v, got %v", expectedToken, actual)
				}
			}
		})
	}
}
