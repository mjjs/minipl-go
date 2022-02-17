package parser

import (
	"fmt"

	"github.com/mjjs/minipl-go/ast"
	"github.com/mjjs/minipl-go/token"
)

type Lexer interface {
	GetNextToken() (token.Token, token.Position)
}

// Parser is the main struct of the parser package. The Parser should be
// initialized with New instead of used directly.
type Parser struct {
	lexer        Lexer
	currentToken token.Token
	currentPos   token.Position
}

// New returns a properly initialized pointer instance to a Parser.
func New(lexer Lexer) *Parser {
	if lexer == nil {
		panic("Attempting to construct a Parser with a nil Lexer")
	}

	tok, pos := lexer.GetNextToken()

	return &Parser{
		lexer:        lexer,
		currentToken: tok,
		currentPos:   pos,
	}
}

// Parse reads tokens from the lexer and verifies that the program is
// syntactically valid. An abstract syntax tree and an optional error is
// returned.
func (p *Parser) Parse() (ast.Prog, error) {
	statements := p.parseStatements()

	if p.currentToken.Type() != token.EOF {
		return ast.Prog{}, fmt.Errorf(
			"Parsing the program failed. Expected %v, found %v",
			token.EOF,
			p.currentToken.Type(),
		)
	}

	return ast.Prog{Statements: statements}, nil
}

// parseStatements goes through all the statements of the lexer and parses
// them returning a Stmts node indicating the root of the abstract syntax tree.
func (p *Parser) parseStatements() ast.Stmts {
	statements := []ast.Stmt{}

	statements = append(statements, p.parseStatement())

	for p.isStatement(p.currentToken.Type()) {
		statements = append(statements, p.parseStatement())
	}

	return ast.Stmts{Statements: statements}
}

// parseStatement parses a statement using the following grammar rules.
//
// <stmt> ::= “var” <var_ident> “:” <type> [ “:=” <expr> ]
//            | <var_ident> “:=” <expr>
//            | “for” <var_ident> “in” <expr> “..” <expr> “do”
//              <stmts> “end” “for”
//            | “read” <var_ident>
//            | “print” <expr>
//            | “assert” “(” <expr> “)”
func (p *Parser) parseStatement() ast.Stmt {
	var statement ast.Stmt

	switch p.currentToken.Type() {
	case token.VAR:
		statement = p.parseDeclaration()
	case token.IDENT:
		statement = p.parseAssignment()
	case token.FOR:
		statement = p.parseForStatement()
	case token.READ:
		statement = p.parseReadStatement()
	case token.PRINT:
		statement = p.parsePrintStatement()
	case token.ASSERT:
		statement = p.parseAssertStatement()

	default:
		panic(fmt.Sprintf("%s: syntax error: unexpected %v", p.currentPos, p.currentToken.Type()))
	}

	p.eat(token.SEMI)

	return statement
}

func (p *Parser) parseDeclaration() ast.DeclStmt {
	pos := p.currentPos

	p.eat(token.VAR)
	ident := p.currentToken
	p.eat(token.IDENT)
	p.eat(token.COLON)
	variableType := p.currentToken
	p.eatType()
	if p.currentToken.Type() != token.ASSIGN {
		return ast.DeclStmt{
			Identifier:   ident,
			VariableType: variableType,
			Pos:          pos,
		}
	}

	p.eat(token.ASSIGN)
	expr := p.parseExpression()

	return ast.DeclStmt{
		Identifier:   ident,
		VariableType: variableType,
		Expression:   expr,
		Pos:          pos,
	}
}

func (p *Parser) parseAssignment() ast.AssignStmt {
	pos := p.currentPos

	ident := p.currentToken

	p.eat(token.IDENT)
	p.eat(token.ASSIGN)
	expr := p.parseExpression()

	return ast.AssignStmt{
		Identifier: ast.Ident{
			Id:  ident,
			Pos: pos,
		},
		Expression: expr,
		Pos:        pos,
	}
}

func (p *Parser) parseForStatement() ast.ForStmt {
	pos := p.currentPos
	p.eat(token.FOR)
	ident := p.currentToken
	identPos := p.currentPos
	p.eat(token.IDENT)
	p.eat(token.IN)
	low := p.parseExpression()
	p.eat(token.DOTDOT)
	high := p.parseExpression()
	p.eat(token.DO)
	statements := p.parseStatements()
	p.eat(token.END)
	p.eat(token.FOR)
	return ast.ForStmt{
		Index: ast.Ident{
			Id:  ident,
			Pos: identPos,
		},
		Low:        low,
		High:       high,
		Statements: statements,
		Pos:        pos,
	}
}

func (p *Parser) parseReadStatement() ast.ReadStmt {
	pos := p.currentPos
	_, identPos := p.eat(token.READ)

	statement := ast.ReadStmt{
		TargetIdentifier: ast.Ident{
			Id:  p.currentToken,
			Pos: identPos,
		},
		Pos: pos,
	}

	p.eat(token.IDENT)

	return statement
}

func (p *Parser) parsePrintStatement() ast.PrintStmt {
	pos := p.currentPos
	p.eat(token.PRINT)

	return ast.PrintStmt{
		Expression: p.parseExpression(),
		Pos:        pos,
	}
}

func (p *Parser) parseAssertStatement() ast.AssertStmt {
	pos := p.currentPos

	p.eat(token.ASSERT)
	p.eat(token.LPAREN)
	statement := ast.AssertStmt{
		Expression: p.parseExpression(),
		Pos:        pos,
	}
	p.eat(token.RPAREN)

	return statement
}

// parseExpression parses an expression with the following grammar rules.
//
// <expr> ::= <opnd> <op> <opnd>
//            | [ <unary_opnd> ] <opnd>
func (p *Parser) parseExpression() ast.Expr {
	pos := p.currentPos

	if p.currentToken.Type() == token.NOT {
		unary := p.currentToken
		p.eat(token.NOT)

		return ast.UnaryExpr{
			Unary:   unary,
			Operand: p.parseOperand(),
			Pos:     pos,
		}
	}

	left := p.parseOperand()

	if !p.isOperator(p.currentToken.Type()) {
		return ast.NullaryExpr{
			Operand: left,
		}
	}

	operand := p.currentToken
	p.eatOperator()

	right := p.parseOperand()

	return ast.BinaryExpr{
		Left:     left,
		Operator: operand,
		Right:    right,
	}
}

// parseOperand parses a valid operand with the following grammar rules.
//
// <opnd> ::= <int>
//            | <string>
//            | <var_ident>
//            | “(” <expr> “)”
//
// <var_ident> ::= <ident>
func (p *Parser) parseOperand() ast.Node {
	pos := p.currentPos

	switch p.currentToken.Type() {
	case token.INTEGER_LITERAL:
		val := p.currentToken.ValueInt()
		p.eat(token.INTEGER_LITERAL)
		return ast.NumberOpnd{
			Value: val,
			Pos:   pos,
		}

	case token.STRING_LITERAL:
		val := p.currentToken.Value()
		p.eat(token.STRING_LITERAL)
		return ast.StringOpnd{
			Value: val,
			Pos:   pos,
		}

	case token.IDENT:
		t := p.currentToken
		p.eat(token.IDENT)
		return ast.Ident{
			Id:  t,
			Pos: pos,
		}

	case token.LPAREN:
		p.eat(token.LPAREN)
		expr := p.parseExpression()
		p.eat(token.RPAREN)
		return expr
	}

	panic(fmt.Sprintf("Syntax error: unexpected %v", p.currentToken.Type()))
}

// isStatement checks whether tokenType should be parsed as a statement node or not.
func (p *Parser) isStatement(tokenType token.TokenTag) bool {
	statementTypes := []token.TokenTag{
		token.VAR,
		token.IDENT,
		token.FOR,
		token.READ,
		token.PRINT,
		token.ASSERT,
	}

	for _, t := range statementTypes {
		if t == tokenType {
			return true
		}
	}

	return false
}

// isOperator checks whether tokenType is a valid operator or not.
func (p *Parser) isOperator(tokenType token.TokenTag) bool {
	operatorTypes := []token.TokenTag{
		token.PLUS,
		token.MINUS,
		token.MULTIPLY,
		token.INTEGER_DIV,
		token.LT,
		token.EQ,
		token.AND,
	}

	for _, t := range operatorTypes {
		if t == tokenType {
			return true
		}
	}

	return false
}

// eat checks that the given tokenType corresponds to the currently held token
// and consumes it. If the tokens do not match, eat panics.
func (p *Parser) eat(tokenType token.TokenTag) (token.Token, token.Position) {
	if p.currentToken.Type() == tokenType {
		p.currentToken, p.currentPos = p.lexer.GetNextToken()
		return p.currentToken, p.currentPos
	} else {
		panic(fmt.Sprintf(
			"Syntax error: expected %v got %v",
			tokenType,
			p.currentToken.Type(),
		))
	}
}

// eatType consumes a type token. If the current token is not a type token,
// eatType panics.
func (p *Parser) eatType() {
	if p.currentToken.Type() == token.INTEGER {
		p.eat(token.INTEGER)
	} else if p.currentToken.Type() == token.STRING {
		p.eat(token.STRING)
	} else if p.currentToken.Type() == token.BOOLEAN {
		p.eat(token.BOOLEAN)
	} else {
		panic(fmt.Sprintf(
			"Syntax error: expected a type, got %v",
			p.currentToken.Type(),
		))
	}
}

// eatOperator consumes an operator token or panics.
func (p *Parser) eatOperator() {
	if p.currentToken.Type() == token.PLUS {
		p.eat(token.PLUS)
	} else if p.currentToken.Type() == token.MINUS {
		p.eat(token.MINUS)
	} else if p.currentToken.Type() == token.MULTIPLY {
		p.eat(token.MULTIPLY)
	} else if p.currentToken.Type() == token.INTEGER_DIV {
		p.eat(token.INTEGER_DIV)
	} else if p.currentToken.Type() == token.LT {
		p.eat(token.LT)
	} else if p.currentToken.Type() == token.EQ {
		p.eat(token.EQ)
	} else if p.currentToken.Type() == token.AND {
		p.eat(token.AND)
	} else {
		panic(fmt.Sprintf(
			"Syntax error: expected an operator, got %v",
			p.currentToken.Type(),
		))
	}
}
