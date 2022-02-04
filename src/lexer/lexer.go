package lexer

import (
	"fmt"
	"strconv"
	"unicode"
)

// Lexer is the main structure of the lexer package. It takes in the source code
// of the MiniPL application and turns it into tokens.
type Lexer struct {
	sourceCode  []rune
	currentChar rune
	pos         int
	eof         bool
}

// New returns a properly initialized pointer to a new Lexer instance using
// sourceCode as the input program.
func New(sourceCode string) *Lexer {
	lexer := &Lexer{
		sourceCode: []rune(sourceCode),
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
func (l *Lexer) GetNextToken() Token {
	for !l.eof {
		if unicode.IsSpace(l.currentChar) {
			l.skipWhitespace()
			continue
		}

		if unicode.IsLetter(l.currentChar) {
			return l.ident()
		}

		if unicode.IsNumber(l.currentChar) {
			if l.currentChar == '0' {
				panic("Number beginning with 0")
			}
			return l.number()
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

			l.advance()
			return Token{tag: INTEGER_DIV}
		}

		if l.currentChar == '"' {
			l.advance()
			return l.string()
		}

		if l.currentChar == ':' {
			next, eof := l.peek()
			if !eof && next == '=' {
				l.advance()
				l.advance()
				return Token{tag: ASSIGN}
			}

			l.advance()
			return Token{tag: COLON}
		}

		if l.currentChar == '.' {
			next, eof := l.peek()
			if !eof && next == '.' {
				l.advance()
				l.advance()
				return Token{tag: DOTDOT}
			}
		}

		if l.currentChar == ';' {
			l.advance()
			return Token{tag: SEMI}
		}

		if l.currentChar == '!' {
			l.advance()
			return Token{tag: NOT}
		}

		if l.currentChar == '+' {
			l.advance()
			return Token{tag: PLUS}
		}

		if l.currentChar == '-' {
			l.advance()
			return Token{tag: MINUS}
		}

		if l.currentChar == '*' {
			l.advance()
			return Token{tag: MULTIPLY}
		}

		if l.currentChar == '<' {
			l.advance()
			return Token{tag: LT}
		}

		if l.currentChar == '=' {
			l.advance()
			return Token{tag: EQ}
		}

		if l.currentChar == '&' {
			l.advance()
			return Token{tag: AND}
		}

		if l.currentChar == '(' {
			l.advance()
			return Token{tag: LPAREN}
		}

		if l.currentChar == ')' {
			l.advance()
			return Token{tag: RPAREN}
		}

		panic(fmt.Sprintf("Could not tokenize character '%c'", l.currentChar))
	}

	return Token{tag: EOF}
}

// advance moves the position of the lexer forward one character and sets the
// EOF flag to true if we have reached the end of the input program.
func (l *Lexer) advance() {
	l.pos++

	if l.pos > len(l.sourceCode)-1 {
		l.eof = true
	} else {
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
func (l *Lexer) ident() Token {
	id := ""

	// TODO: Check if unicode.X is allowed
	for !l.eof && unicode.In(l.currentChar, unicode.Number, unicode.Letter) || l.currentChar == '_' {
		id += string(l.currentChar)
		l.advance()
	}

	token, ok := ReservedKeywords[id]
	if ok {
		return token
	}

	return Token{tag: IDENT, lexeme: id}
}

// number reads a number from the input program and returns an INTEGER_LITERAL
// token with the number as a lexeme. MiniPL only supports integers, so we
// do not consider floating point numbers.
func (l *Lexer) number() Token {
	numString := ""

	for !l.eof && unicode.IsNumber(l.currentChar) {
		numString += string(l.currentChar)
		l.advance()
	}

	num, err := strconv.Atoi(numString)
	if err != nil {
		panic(fmt.Sprintf("Could not tokenize number %s, %v", numString, err))
	}

	return Token{tag: INTEGER_LITERAL, lexeme: num}
}

// string reads a string from the input program and returns a STRING_LITERAL
// token containing the read string as a lexeme.
func (l *Lexer) string() Token {
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

	return Token{tag: STRING_LITERAL, lexeme: str}
}
