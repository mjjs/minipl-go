package lexer

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
)

// reservedKeywords maps the reserved keywords of MiniPL into the tokens for the keywords.
var reservedKeywords map[string]Token = map[string]Token{
	"var":    {tag: VAR},
	"for":    {tag: FOR},
	"end":    {tag: END},
	"in":     {tag: IN},
	"do":     {tag: DO},
	"read":   {tag: READ},
	"print":  {tag: PRINT},
	"int":    {tag: INTEGER},
	"string": {tag: STRING},
	"bool":   {tag: BOOLEAN},
	"assert": {tag: ASSERT},
}

// Token represents a single token found in the program by scanning it.
// It consists of a tag and an optional lexeme.
type Token struct {
	tag    TokenTag
	lexeme interface{}
}

// NewToken constructs a Token with the given tag and lexeme.
// A nil lexeme can be passed in if the token is not expecting a lexeme.
func NewToken(tag TokenTag, lexeme interface{}) Token { return Token{tag, lexeme} }

// ValueInt returns the lexeme of the token as an integer or panics if the lexeme is not an integer.
func (t Token) ValueInt() int { return t.lexeme.(int) }

// ValueBool is like ValueInt but returns a boolean.
func (t Token) ValueBool() bool { return t.lexeme.(bool) }

// ValueString is like ValueInt but returns a String.
func (t Token) ValueString() string { return t.lexeme.(string) }

// Type returns the tag of the token.
func (t Token) Type() TokenTag { return t.tag }
