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

var ReservedKeywords map[string]Token = map[string]Token{
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

type Token struct {
	tag    TokenTag
	lexeme interface{}
}

func NewToken(tag TokenTag, lexeme interface{}) Token { return Token{tag, lexeme} }

func (t Token) ValueInt() int       { return t.lexeme.(int) }
func (t Token) ValueBool() bool     { return t.lexeme.(bool) }
func (t Token) ValueString() string { return t.lexeme.(string) }
func (t Token) Type() TokenTag      { return t.tag }
