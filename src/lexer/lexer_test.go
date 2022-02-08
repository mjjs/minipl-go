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
		expectedTokens: []token.Token{token.New(token.EOF, "")},
	},
	{
		name:           "Integer",
		input:          "int",
		expectedTokens: []token.Token{token.New(token.INTEGER, "")},
	},
	{
		name:           "String type",
		input:          "string",
		expectedTokens: []token.Token{token.New(token.STRING, "")},
	},
	{
		name:           "Boolean type",
		input:          "bool",
		expectedTokens: []token.Token{token.New(token.BOOLEAN, "")},
	},
	{
		name:           "Integer literal",
		input:          "666",
		expectedTokens: []token.Token{token.New(token.INTEGER_LITERAL, "666")},
	},
	{
		name:        "Integer starting with zero",
		input:       "0123",
		shouldPanic: true,
	},
	{
		name:           "Integer with only zero",
		input:          "0",
		expectedTokens: []token.Token{token.New(token.INTEGER_LITERAL, "0")},
	},
	{
		name:           "String literal",
		input:          "\"C-beams\"",
		expectedTokens: []token.Token{token.New(token.STRING_LITERAL, "C-beams")},
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
	},
	{
		name:  "String literal with newlines",
		input: `"a\nb\r"`,
		expectedTokens: []token.Token{
			token.New(token.STRING_LITERAL, "a\nb\r"),
		},
	},
	{
		name:  "String literal with tab character",
		input: `"a:\tb"`,
		expectedTokens: []token.Token{
			token.New(token.STRING_LITERAL, "a:\tb"),
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
			token.New(token.IDENT, "myOwnVar"),
		},
	},
	{
		name:           "Plus operator",
		input:          "+",
		expectedTokens: []token.Token{token.New(token.PLUS, "")},
	},
	{
		name:           "Minus operator",
		input:          "-",
		expectedTokens: []token.Token{token.New(token.MINUS, "")},
	},
	{
		name:           "Multiply operator",
		input:          "*",
		expectedTokens: []token.Token{token.New(token.MULTIPLY, "")},
	},
	{
		name:           "Integer division operator",
		input:          "/",
		expectedTokens: []token.Token{token.New(token.INTEGER_DIV, "")},
	},
	{
		name:           "Less than operator",
		input:          "<",
		expectedTokens: []token.Token{token.New(token.LT, "")},
	},
	{
		name:           "Equality operator",
		input:          "=",
		expectedTokens: []token.Token{token.New(token.EQ, "")},
	},
	{
		name:           "Logical and operator",
		input:          "&",
		expectedTokens: []token.Token{token.New(token.AND, "")},
	},
	{
		name:           "Logical not operator",
		input:          "!",
		expectedTokens: []token.Token{token.New(token.NOT, "")},
	},
	{
		name:           "Assign operator",
		input:          ":=",
		expectedTokens: []token.Token{token.New(token.ASSIGN, "")},
	},
	{
		name:           "Left parenthesis",
		input:          "(",
		expectedTokens: []token.Token{token.New(token.LPAREN, "")},
	},
	{
		name:           "Right parenthesis",
		input:          ")",
		expectedTokens: []token.Token{token.New(token.RPAREN, "")},
	},
	{
		name:           "Semicolon",
		input:          ";",
		expectedTokens: []token.Token{token.New(token.SEMI, "")},
	},
	{
		name:           "Colon",
		input:          ":",
		expectedTokens: []token.Token{token.New(token.COLON, "")},
	},
	{
		name:           "For keyword",
		input:          "for",
		expectedTokens: []token.Token{token.New(token.FOR, "")},
	},
	{
		name:           "In keyword",
		input:          "in",
		expectedTokens: []token.Token{token.New(token.IN, "")},
	},
	{
		name:           "Do keyword",
		input:          "do",
		expectedTokens: []token.Token{token.New(token.DO, "")},
	},
	{
		name:           "End keyword",
		input:          "end",
		expectedTokens: []token.Token{token.New(token.END, "")},
	},
	{
		name:           "Double dot keyword",
		input:          "..",
		expectedTokens: []token.Token{token.New(token.DOTDOT, "")},
	},
	{
		name:           "Assert keyword",
		input:          "assert",
		expectedTokens: []token.Token{token.New(token.ASSERT, "")},
	},
	{
		name:           "Var keyword",
		input:          "var",
		expectedTokens: []token.Token{token.New(token.VAR, "")},
	},
	{
		name:           "Read keyword",
		input:          "read",
		expectedTokens: []token.Token{token.New(token.READ, "")},
	},
	{
		name:           "Print keyword",
		input:          "print",
		expectedTokens: []token.Token{token.New(token.PRINT, "")},
	},
	{
		name:  "Whitespace is skipped",
		input: "first   \n    second",
		expectedTokens: []token.Token{
			token.New(token.IDENT, "first"),
			token.New(token.IDENT, "second"),
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
	},
	{
		name:  "Parentheses tokenize correctly",
		input: "(52)",
		expectedTokens: []token.Token{
			token.New(token.LPAREN, ""),
			token.New(token.INTEGER_LITERAL, "52"),
			token.New(token.RPAREN, ""),
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
				for lexer.GetNextToken() != token.New(token.EOF, "") {
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
