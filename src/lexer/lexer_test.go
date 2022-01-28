package lexer

import "testing"

var getNextTokenTestCases = []struct {
	name           string
	input          string
	expectedTokens []Token
}{
	{
		name:           "Empty input returns EOF token",
		input:          "",
		expectedTokens: []Token{{tag: EOF}},
	},
	{
		name:           "Integer",
		input:          "int",
		expectedTokens: []Token{{tag: INTEGER}},
	},
	{
		name:           "String type",
		input:          "string",
		expectedTokens: []Token{{tag: STRING}},
	},
	{
		name:           "Boolean type",
		input:          "bool",
		expectedTokens: []Token{{tag: BOOLEAN}},
	},
	{
		name:           "Integer literal",
		input:          "666",
		expectedTokens: []Token{{tag: INTEGER_LITERAL, lexeme: 666}},
	},
	{
		name:           "String literal",
		input:          "\"C-beams\"",
		expectedTokens: []Token{{tag: STRING_LITERAL, lexeme: "C-beams"}},
	},
	{
		name:           "Single ID",
		input:          "myOwnVar",
		expectedTokens: []Token{{tag: IDENT, lexeme: "myOwnVar"}},
	},
	{
		name:           "Plus operator",
		input:          "+",
		expectedTokens: []Token{{tag: PLUS}},
	},
	{
		name:           "Minus operator",
		input:          "-",
		expectedTokens: []Token{{tag: MINUS}},
	},
	{
		name:           "Multiply operator",
		input:          "*",
		expectedTokens: []Token{{tag: MULTIPLY}},
	},
	{
		name:           "Integer division operator",
		input:          "/",
		expectedTokens: []Token{{tag: INTEGER_DIV}},
	},
	{
		name:           "Less than operator",
		input:          "<",
		expectedTokens: []Token{{tag: LT}},
	},
	{
		name:           "Equality operator",
		input:          "=",
		expectedTokens: []Token{{tag: EQ}},
	},
	{
		name:           "Logical and operator",
		input:          "&",
		expectedTokens: []Token{{tag: AND}},
	},
	{
		name:           "Logical not operator",
		input:          "!",
		expectedTokens: []Token{{tag: NOT}},
	},
	{
		name:           "Assign operator",
		input:          ":=",
		expectedTokens: []Token{{tag: ASSIGN}},
	},
	{
		name:           "Left parenthesis",
		input:          "(",
		expectedTokens: []Token{{tag: LPAREN}},
	},
	{
		name:           "Right parenthesis",
		input:          ")",
		expectedTokens: []Token{{tag: RPAREN}},
	},
	{
		name:           "Semicolon",
		input:          ";",
		expectedTokens: []Token{{tag: SEMI}},
	},
	{
		name:           "Colon",
		input:          ":",
		expectedTokens: []Token{{tag: COLON}},
	},
	{
		name:           "For keyword",
		input:          "for",
		expectedTokens: []Token{{tag: FOR}},
	},
	{
		name:           "In keyword",
		input:          "in",
		expectedTokens: []Token{{tag: IN}},
	},
	{
		name:           "Do keyword",
		input:          "do",
		expectedTokens: []Token{{tag: DO}},
	},
	{
		name:           "End keyword",
		input:          "end",
		expectedTokens: []Token{{tag: END}},
	},
	{
		name:           "Double dot keyword",
		input:          "..",
		expectedTokens: []Token{{tag: DOTDOT}},
	},
	{
		name:           "Assert keyword",
		input:          "assert",
		expectedTokens: []Token{{tag: ASSERT}},
	},
	{
		name:           "Var keyword",
		input:          "var",
		expectedTokens: []Token{{tag: VAR}},
	},
	{
		name:           "Read keyword",
		input:          "read",
		expectedTokens: []Token{{tag: READ}},
	},
	{
		name:           "Print keyword",
		input:          "print",
		expectedTokens: []Token{{tag: PRINT}},
	},
	{
		name:  "Whitespace is skipped",
		input: "first   \n    second",
		expectedTokens: []Token{
			{tag: IDENT, lexeme: "first"},
			{tag: IDENT, lexeme: "second"},
		},
	},
	{
		name:  "For loop",
		input: "for i in 1..n do\n\tv := v * i;\nend for;",
		expectedTokens: []Token{
			{tag: FOR},
			{tag: IDENT, lexeme: "i"},
			{tag: IN},
			{tag: INTEGER_LITERAL, lexeme: 1},
			{tag: DOTDOT},
			{tag: IDENT, lexeme: "n"},
			{tag: DO},
			{tag: IDENT, lexeme: "v"},
			{tag: ASSIGN},
			{tag: IDENT, lexeme: "v"},
			{tag: MULTIPLY},
			{tag: IDENT, lexeme: "i"},
			{tag: SEMI},
			{tag: END},
			{tag: FOR},
			{tag: SEMI},
		},
	},
	{
		name:  "Parentheses tokenize correctly",
		input: "(52)",
		expectedTokens: []Token{
			{tag: LPAREN},
			{tag: INTEGER_LITERAL, lexeme: 52},
			{tag: RPAREN},
		},
	},
	{
		name:  "Declaration with assignment",
		input: "var x : int := 4;",
		expectedTokens: []Token{
			{tag: VAR},
			{tag: IDENT, lexeme: "x"},
			{tag: COLON},
			{tag: INTEGER},
			{tag: ASSIGN},
			{tag: INTEGER_LITERAL, lexeme: 4},
			{tag: SEMI},
		},
	},
	{
		name:  "Declaration and assignment in different lines",
		input: "var firstLine : string;\nfirstLine := \"secondLine\";",
		expectedTokens: []Token{
			{tag: VAR},
			{tag: IDENT, lexeme: "firstLine"},
			{tag: COLON},
			{tag: STRING},
			{tag: SEMI},
			{tag: IDENT, lexeme: "firstLine"},
			{tag: ASSIGN},
			{tag: STRING_LITERAL, lexeme: "secondLine"},
			{tag: SEMI},
		},
	},
	{
		name:  "Assert with weird formatting",
		input: "            assert    (x     =         somethingElse)   ;",
		expectedTokens: []Token{
			{tag: ASSERT},
			{tag: LPAREN},
			{tag: IDENT, lexeme: "x"},
			{tag: EQ},
			{tag: IDENT, lexeme: "somethingElse"},
			{tag: RPAREN},
			{tag: SEMI},
		},
	},
	{
		name:  "Skips line comments",
		input: "3 + 3; // this is a comment\n4 + 4;",
		expectedTokens: []Token{
			{tag: INTEGER_LITERAL, lexeme: 3},
			{tag: PLUS},
			{tag: INTEGER_LITERAL, lexeme: 3},
			{tag: SEMI},
			{tag: INTEGER_LITERAL, lexeme: 4},
			{tag: PLUS},
			{tag: INTEGER_LITERAL, lexeme: 4},
			{tag: SEMI},
		},
	},
	{
		name:  "Skips block comments",
		input: "1 + 2; /*This is a comment*/ 3 + 4;",
		expectedTokens: []Token{
			{tag: INTEGER_LITERAL, lexeme: 1},
			{tag: PLUS},
			{tag: INTEGER_LITERAL, lexeme: 2},
			{tag: SEMI},
			{tag: INTEGER_LITERAL, lexeme: 3},
			{tag: PLUS},
			{tag: INTEGER_LITERAL, lexeme: 4},
			{tag: SEMI},
		},
	},
	{
		name:  "String with escaped characters",
		input: `"abc\\ def \""`,
		expectedTokens: []Token{
			{tag: STRING_LITERAL, lexeme: `abc\ def "`},
		},
	},
}

func TestGetNextToken(t *testing.T) {
	for _, testCase := range getNextTokenTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			lexer := New(testCase.input)

			for _, token := range testCase.expectedTokens {
				actual := lexer.GetNextToken()

				if actual != token {
					t.Errorf("Expected %v, got %v", token, actual)
				}
			}
		})
	}
}
