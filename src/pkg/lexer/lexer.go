package lexer

import (
	"fmt"
	"unicode"

	"github.com/mjjs/minipl-go/pkg/token"
)

// reservedKeywords maps the reserved keywords of MiniPL into the tokens for the keywords.
var reservedKeywords map[string]token.Token = map[string]token.Token{
	"var":    token.New(token.VAR, ""),
	"for":    token.New(token.FOR, ""),
	"end":    token.New(token.END, ""),
	"in":     token.New(token.IN, ""),
	"do":     token.New(token.DO, ""),
	"read":   token.New(token.READ, ""),
	"print":  token.New(token.PRINT, ""),
	"int":    token.New(token.INTEGER, ""),
	"string": token.New(token.STRING, ""),
	"bool":   token.New(token.BOOLEAN, ""),
	"assert": token.New(token.ASSERT, ""),
}

// Lexer is the main structure of the lexer package. It takes in the source code
// of the MiniPL application and turns it into tokens.
type Lexer struct {
	sourceCode  []rune
	currentChar rune
	pos         int
	eof         bool

	tokenPos token.Position
}

// New returns a properly initialized pointer to a new Lexer instance using
// sourceCode as the input program.
func New(sourceCode string) *Lexer {
	lexer := &Lexer{
		sourceCode: []rune(sourceCode),
		tokenPos:   token.Position{Line: 1, Column: 1},
	}

	if len(sourceCode) == 0 {
		lexer.eof = true
	} else {
		lexer.currentChar = lexer.sourceCode[lexer.pos]
	}

	return lexer
}

// GetNextToken returns the next token that the Lexer can parse from the
// sourceCode given during initialization.
func (l *Lexer) GetNextToken() (token.Token, token.Position) {
	for !l.eof {
		pos := l.tokenPos

		if unicode.IsSpace(l.currentChar) {
			l.skipWhitespace()
			continue
		}

		if unicode.IsLetter(l.currentChar) {
			return l.ident()
		}

		if unicode.IsNumber(l.currentChar) {
			return l.number()
		}

		if l.currentChar == '/' {
			next, eof := l.peek()
			if !eof && next == '/' {
				l.skipLineComment()
				continue
			}

			if !eof && next == '*' {
				tok, pos := l.skipBlockComment()
				if tok != nil {
					return *tok, pos
				}

				continue
			}

			l.advance()
			return token.New(token.INTEGER_DIV, ""), pos
		}

		if l.currentChar == '"' {
			return l.string()
		}

		if l.currentChar == ':' {
			next, eof := l.peek()
			if !eof && next == '=' {
				l.advance()
				l.advance()
				return token.New(token.ASSIGN, ""), pos
			}

			l.advance()
			return token.New(token.COLON, ""), pos
		}

		if l.currentChar == '.' {
			next, eof := l.peek()
			if !eof && next == '.' {
				l.advance()
				l.advance()
				return token.New(token.DOTDOT, ""), pos
			}
		}

		if l.currentChar == ';' {
			l.advance()
			return token.New(token.SEMI, ""), pos
		}

		if l.currentChar == '!' {
			l.advance()
			return token.New(token.NOT, ""), pos
		}

		if l.currentChar == '+' {
			l.advance()
			return token.New(token.PLUS, ""), pos
		}

		if l.currentChar == '-' {
			l.advance()
			return token.New(token.MINUS, ""), pos
		}

		if l.currentChar == '*' {
			l.advance()
			return token.New(token.MULTIPLY, ""), pos
		}

		if l.currentChar == '<' {
			l.advance()
			return token.New(token.LT, ""), pos
		}

		if l.currentChar == '=' {
			l.advance()
			return token.New(token.EQ, ""), pos
		}

		if l.currentChar == '&' {
			l.advance()
			return token.New(token.AND, ""), pos
		}

		if l.currentChar == '(' {
			l.advance()
			return token.New(token.LPAREN, ""), pos
		}

		if l.currentChar == ')' {
			l.advance()
			return token.New(token.RPAREN, ""), pos
		}

		l.advance()

		errorToken := token.New(token.ERROR,
			fmt.Sprintf("unrecognized character '%c'", l.currentChar))

		return errorToken, pos
	}

	return token.New(token.EOF, ""), l.tokenPos
}

// advance moves the position of the lexer forward one character and sets the
// EOF flag to true if we have reached the end of the input program.
func (l *Lexer) advance() {
	l.pos++

	if l.pos > len(l.sourceCode)-1 {
		l.eof = true
	} else {
		if l.currentChar == '\n' {
			l.tokenPos.Line++
			l.tokenPos.Column = 1
		} else {
			l.tokenPos.Column++
		}

		l.currentChar = rune(l.sourceCode[l.pos])
	}
}

// peek returns the next rune of the source code without advancing the position
// of the lexer. The returned boolean indicates whether we have reached EOF
// or not. If EOF is reached, the returned rune should be discarded.
func (l *Lexer) peek() (rune, bool) {
	if l.pos+1 > len(l.sourceCode)-1 {
		return 0, true
	}

	return l.sourceCode[l.pos+1], false
}

// skipWhitespace advances the lexer until the next non-whitespace character.
func (l *Lexer) skipWhitespace() {
	for !l.eof && unicode.IsSpace(l.currentChar) {
		l.advance()
	}
}

// skipLineComment advances the lexer until the end of a line.
func (l *Lexer) skipLineComment() {
	l.advance()
	l.advance()

	for !l.eof && l.currentChar != '\n' {
		l.advance()
	}

	l.advance()
}

// skipBlockComment advances the lexer until it has skipped all the characters
// inside a block comment.
// Returns a non-nil Error token if the block comment is unterminated.
func (l *Lexer) skipBlockComment() (*token.Token, token.Position) {
	pos := l.tokenPos
	l.advance()
	l.advance()

	for !l.eof {
		if l.currentChar != '*' {
			l.advance()
			continue
		}

		next, eof := l.peek()

		if !eof && next == '/' {
			l.advance()
			l.advance()
			return nil, pos
		} else if !eof {
			l.advance()
		}
	}

	tok := token.New(token.ERROR, "Unterminated block comment")

	return &tok, pos
}

// ident reads a A-Za-z0-9_ string from the input program and returns an IDENT
// token with the read string as a lexeme.
func (l *Lexer) ident() (token.Token, token.Position) {
	pos := l.tokenPos

	id := ""

	// TODO: Check if unicode.X is allowed
	for !l.eof && unicode.In(l.currentChar, unicode.Number, unicode.Letter) || l.currentChar == '_' {
		id += string(l.currentChar)
		l.advance()
	}

	t, ok := reservedKeywords[id]
	if ok {
		return t, pos
	}

	return token.New(token.IDENT, id), pos
}

// number reads a number from the input program and returns an INTEGER_LITERAL
// token with the number as a lexeme. MiniPL only supports integers, so we
// do not consider floating point numbers.
func (l *Lexer) number() (token.Token, token.Position) {
	pos := l.tokenPos
	numString := ""

	for !l.eof && unicode.IsNumber(l.currentChar) {
		numString += string(l.currentChar)
		l.advance()
	}

	return token.New(token.INTEGER_LITERAL, numString), pos
}

// string reads a string from the input program and returns a STRING_LITERAL
// token containing the read string as a lexeme.
func (l *Lexer) string() (token.Token, token.Position) {
	pos := l.tokenPos
	str := ""

	l.advance()

	for !l.eof {
		if l.currentChar == '\\' {
			l.advance()

			switch l.currentChar {
			case 'n':
				str += "\n"
			case 't':
				str += "\t"
			case 'r':
				str += "\r"
			default:
				str += string(l.currentChar)
			}

			l.advance()
			continue
		}

		if l.currentChar == '"' {
			l.advance()
			return token.New(token.STRING_LITERAL, str), pos
		}

		if l.currentChar == '\n' || l.currentChar == '\r' {
			break
		}

		str += string(l.currentChar)
		l.advance()
	}

	tok := token.New(
		token.ERROR,
		fmt.Sprintf("unterminated string literal %s", str),
	)

	return tok, pos

}
