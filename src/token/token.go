package token

import "strconv"

type TokenTag string

const (
	// Types
	INTEGER TokenTag = "INTEGER"
	STRING           = "STRING"
	BOOLEAN          = "BOOLEAN"

	// Constants
	INTEGER_LITERAL = "INTEGER_LITERAL"
	STRING_LITERAL  = "STRING_LITERAL"
	BOOLEAN_LITERAL = "BOOLEAN_LITERAL"

	// Identifier
	IDENT = "IDENT"

	// Operators
	PLUS        = "PLUS"        // +
	MINUS       = "MINUS"       // -
	MULTIPLY    = "MULTIPLY"    // *
	INTEGER_DIV = "INTEGER_DIV" // /
	LT          = "LT"          // <
	EQ          = "EQ"          // =
	AND         = "AND"         // &
	NOT         = "NOT"         // !
	ASSIGN      = "ASSIGN"      // :=

	LPAREN = "LPAREN" // (
	RPAREN = "RPAREN" // )
	SEMI   = "SEMI"   // ;
	COLON  = "COLON"  // :

	// For loop
	FOR    = "FOR"
	IN     = "IN"
	DO     = "DO"
	END    = "END"
	DOTDOT = "DOTDOT"

	ASSERT = "ASSERT"
	VAR    = "VAR"
	READ   = "READ"
	PRINT  = "PRINT"

	EOF = "EOF"

	ERROR = "ERROR" // For reporting invalid tokens
)

// Token represents a single token found in the program by scanning it.
// It consists of a tag and an optional lexeme.
type Token struct {
	tag    TokenTag
	lexeme string
}

// New constructs a Token with the given tag and lexeme.
// An empty lexeme can be passed in if the token is not expecting a lexeme.
func New(tag TokenTag, lexeme string) Token { return Token{tag, lexeme} }

// ValueInt returns the lexeme of the token as an integer or panics if the lexeme is not an integer.
func (t Token) ValueInt() int {
	if t.lexeme == "" {
		panic("Attempting to take value of an empty lexeme")
	}

	num, err := strconv.Atoi(t.lexeme)
	if err != nil {
		panic(err)
	}
	return num
}

// ValueBool is like ValueInt but returns a boolean.
func (t Token) ValueBool() bool {
	if t.lexeme == "" {
		panic("Attempting to take value of an empty lexeme")
	}

	x, err := strconv.ParseBool(t.lexeme)
	if err != nil {
		panic(err)
	}
	return x
}

// Value returns the value of the lexeme.
func (t Token) Value() string {
	if t.lexeme == "" {
		panic("Attempting to take value of an empty lexeme")
	}

	return t.lexeme
}

// Type returns the tag of the token.
func (t Token) Type() TokenTag { return t.tag }
