package lexer

import (
	"fmt"
	"unicode"

	"github.com/mjjs/minipl-go/token"
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

	currentLine int
	currentCol  int
}

// New returns a properly initialized pointer to a new Lexer instance using
// sourceCode as the input program.
func New(sourceCode string) *Lexer {
	lexer := &Lexer{
		sourceCode:  []rune(sourceCode),
		currentLine: 1,
		currentCol:  1,
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
		if unicode.IsSpace(l.currentChar) {
			l.skipWhitespace()
			continue
		}

		if unicode.IsLetter(l.currentChar) {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			return l.ident(), pos
		}

		if unicode.IsNumber(l.currentChar) {
			if l.currentChar == '0' {
				if next, eof := l.peek(); !eof && unicode.IsNumber(next) {
					panic("Number beginning with 0")
				}
			}

			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			return l.number(), pos
		}

		if l.currentChar == '/' {
			next, eof := l.peek()
			if !eof && next == '/' {
				l.advance()
				l.advance()
				l.skipLineComment()
				continue
			}

			if !eof && next == '*' {
				l.advance()
				l.advance()
				l.skipBlockComment()
				continue
			}

			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.INTEGER_DIV, ""), pos
		}

		if l.currentChar == '"' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return l.string(), pos
		}

		if l.currentChar == ':' {
			next, eof := l.peek()
			if !eof && next == '=' {
				pos := token.Position{Line: l.currentLine, Column: l.currentCol}
				l.advance()
				l.advance()
				return token.New(token.ASSIGN, ""), pos
			}

			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.COLON, ""), pos
		}

		if l.currentChar == '.' {
			next, eof := l.peek()
			if !eof && next == '.' {
				pos := token.Position{Line: l.currentLine, Column: l.currentCol}
				l.advance()
				l.advance()
				return token.New(token.DOTDOT, ""), pos
			}
		}

		if l.currentChar == ';' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.SEMI, ""), pos
		}

		if l.currentChar == '!' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.NOT, ""), pos
		}

		if l.currentChar == '+' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.PLUS, ""), pos
		}

		if l.currentChar == '-' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.MINUS, ""), pos
		}

		if l.currentChar == '*' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.MULTIPLY, ""), pos
		}

		if l.currentChar == '<' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.LT, ""), pos
		}

		if l.currentChar == '=' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.EQ, ""), pos
		}

		if l.currentChar == '&' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.AND, ""), pos
		}

		if l.currentChar == '(' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.LPAREN, ""), pos
		}

		if l.currentChar == ')' {
			pos := token.Position{Line: l.currentLine, Column: l.currentCol}
			l.advance()
			return token.New(token.RPAREN, ""), pos
		}

		panic(fmt.Sprintf("Could not tokenize character '%c'", l.currentChar))
	}

	pos := token.Position{Line: l.currentLine, Column: l.currentCol}
	return token.New(token.EOF, ""), pos
}

// advance moves the position of the lexer forward one character and sets the
// EOF flag to true if we have reached the end of the input program.
func (l *Lexer) advance() {
	l.pos++

	if l.pos > len(l.sourceCode)-1 {
		l.eof = true
	} else {
		if l.currentChar == '\n' {
			l.currentLine++
			l.currentCol = 1
		} else {
			l.currentCol++
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
	for !l.eof && l.currentChar != '\n' {
		l.advance()
	}

	l.advance()
}

// skipBlockComment advances the lexer until it has skipped all the characters
// inside a block comment.
func (l *Lexer) skipBlockComment() {
	for !l.eof {
		if l.currentChar != '*' {
			l.advance()
			continue
		}

		next, eof := l.peek()

		if !eof && next == '/' {
			l.advance()
			l.advance()
			return
		}
	}
}

// ident reads a A-Za-z0-9_ string from the input program and returns an IDENT
// token with the read string as a lexeme.
func (l *Lexer) ident() token.Token {
	id := ""

	// TODO: Check if unicode.X is allowed
	for !l.eof && unicode.In(l.currentChar, unicode.Number, unicode.Letter) || l.currentChar == '_' {
		id += string(l.currentChar)
		l.advance()
	}

	t, ok := reservedKeywords[id]
	if ok {
		return t
	}

	return token.New(token.IDENT, id)
}

// number reads a number from the input program and returns an INTEGER_LITERAL
// token with the number as a lexeme. MiniPL only supports integers, so we
// do not consider floating point numbers.
func (l *Lexer) number() token.Token {
	numString := ""

	for !l.eof && unicode.IsNumber(l.currentChar) {
		numString += string(l.currentChar)
		l.advance()
	}

	return token.New(token.INTEGER_LITERAL, numString)
}

// string reads a string from the input program and returns a STRING_LITERAL
// token containing the read string as a lexeme.
func (l *Lexer) string() token.Token {
	str := ""

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
			break
		}

		if l.currentChar == '\n' || l.currentChar == '\r' {
			panic("Unterminated string literal")
		}

		str += string(l.currentChar)
		l.advance()
	}

	return token.New(token.STRING_LITERAL, str)
}
