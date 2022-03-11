package parser

import (
	"fmt"

	"github.com/mjjs/minipl-go/pkg/ast"
	"github.com/mjjs/minipl-go/pkg/token"
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

	errors []error
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
func (p *Parser) Parse() (ast.Prog, []error) {
	statements := p.parseStatements()

	p.eat(token.EOF)

	return ast.Prog{Statements: statements}, p.errors
}

// parseStatements goes through all the statements of the lexer and parses
// them returning a Stmts node indicating the root of the abstract syntax tree.
func (p *Parser) parseStatements() ast.Stmts {
	statements := []ast.Stmt{}

	statements = append(statements, p.parseStatement())

	for p.currentToken.IsStatement() {
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
		err := fmt.Errorf(
			"%s: syntax error: unexpected %v",
			p.currentPos, p.currentToken.Type(),
		)

		p.errors = append(p.errors, err)
	}

	return statement
}

func (p *Parser) parseDeclaration() ast.DeclStmt {
	pos := p.currentPos

	if !p.eat(token.VAR) {
		p.skipStatement()
		return ast.DeclStmt{}
	}

	ident := p.currentToken
	if !p.eat(token.IDENT) {
		p.skipStatement()
		return ast.DeclStmt{}
	}
	if !p.eat(token.COLON) {
		p.skipStatement()
		return ast.DeclStmt{}
	}

	variableType := p.currentToken
	if !p.currentToken.IsType() {
		err := fmt.Errorf(
			"Syntax error: expected a type, got %v",
			p.currentToken.Type(),
		)

		p.errors = append(p.errors, err)
		p.skipStatement()
		return ast.DeclStmt{}
	}

	if !p.eat(p.currentToken.Type()) {
		p.skipStatement()
		return ast.DeclStmt{}
	}

	if p.currentToken.Type() != token.ASSIGN {
		p.eat(token.SEMI)

		return ast.DeclStmt{
			Identifier:   ident,
			VariableType: variableType,
			Pos:          pos,
		}
	}

	if !p.eat(token.ASSIGN) {
		p.skipStatement()
		return ast.DeclStmt{}
	}

	expr := p.parseExpression()
	if expr == nil {
		p.skipStatement()
		return ast.DeclStmt{}
	}

	p.eat(token.SEMI)

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

	if !p.eat(token.IDENT) {
		p.skipStatement()
		return ast.AssignStmt{}
	}

	if !p.eat(token.ASSIGN) {
		p.skipStatement()
		return ast.AssignStmt{}
	}

	expr := p.parseExpression()
	if expr == nil {
		p.skipStatement()
		return ast.AssignStmt{}
	}

	p.eat(token.SEMI)

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

	if !p.eat(token.FOR) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	ident := p.currentToken
	identPos := p.currentPos

	if !p.eat(token.IDENT) {
		p.skipForBlock()
		return ast.ForStmt{}
	}
	if !p.eat(token.IN) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	low := p.parseExpression()
	if low == nil {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	if !p.eat(token.RANGE) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	high := p.parseExpression()
	if high == nil {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	if !p.eat(token.DO) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	statements := p.parseStatements()

	if !p.eat(token.END) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	if !p.eat(token.FOR) {
		p.skipForBlock()
		return ast.ForStmt{}
	}

	p.eat(token.SEMI)

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

	if !p.eat(token.READ) {
		p.skipStatement()
		return ast.ReadStmt{}
	}

	identPos := p.currentPos

	statement := ast.ReadStmt{
		TargetIdentifier: ast.Ident{
			Id:  p.currentToken,
			Pos: identPos,
		},
		Pos: pos,
	}

	if !p.eat(token.IDENT) {
		p.skipStatement()
		return ast.ReadStmt{}
	}

	p.eat(token.SEMI)

	return statement
}

func (p *Parser) parsePrintStatement() ast.PrintStmt {
	pos := p.currentPos
	if !p.eat(token.PRINT) {
		p.skipStatement()
		return ast.PrintStmt{}
	}

	expr := p.parseExpression()
	if expr == nil {
		p.skipStatement()
		return ast.PrintStmt{}
	}

	p.eat(token.SEMI)

	return ast.PrintStmt{
		Expression: expr,
		Pos:        pos,
	}
}

func (p *Parser) parseAssertStatement() ast.AssertStmt {
	pos := p.currentPos

	if !p.eat(token.ASSERT) {
		p.skipStatement()
		return ast.AssertStmt{}
	}

	if !p.eat(token.LPAREN) {
		p.skipStatement()
		return ast.AssertStmt{}
	}

	expr := p.parseExpression()
	if expr == nil {
		p.skipStatement()
		return ast.AssertStmt{}
	}

	statement := ast.AssertStmt{
		Expression: expr,
		Pos:        pos,
	}

	if !p.eat(token.RPAREN) {
		p.skipStatement()
		return ast.AssertStmt{}
	}

	p.eat(token.SEMI)

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

		operand := p.parseOperand()
		if operand == nil {
			return nil
		}

		return ast.UnaryExpr{
			Unary:   unary,
			Operand: operand,
			Pos:     pos,
		}
	}

	left := p.parseOperand()
	if left == nil {
		return nil
	}

	if !p.currentToken.IsOperator() {
		return ast.NullaryExpr{
			Operand: left,
		}
	}

	operand := p.currentToken

	if !p.currentToken.IsOperator() {
		err := fmt.Errorf(
			"Syntax error: expected an operator, got %v",
			p.currentToken.Type(),
		)
		p.errors = append(p.errors, err)
		return nil
	}

	p.eat(p.currentToken.Type())

	right := p.parseOperand()
	if right == nil {
		return nil
	}

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

		if !p.eat(token.INTEGER_LITERAL) {
			return nil
		}

		return ast.NumberOpnd{
			Value: val,
			Pos:   pos,
		}

	case token.STRING_LITERAL:
		val := p.currentToken.Value()

		if !p.eat(token.STRING_LITERAL) {
			return nil
		}

		return ast.StringOpnd{
			Value: val,
			Pos:   pos,
		}

	case token.IDENT:
		t := p.currentToken
		if !p.eat(token.IDENT) {
			return nil
		}

		return ast.Ident{
			Id:  t,
			Pos: pos,
		}

	case token.LPAREN:
		if !p.eat(token.LPAREN) {
			return nil
		}

		expr := p.parseExpression()
		if expr == nil {
			return nil
		}

		if !p.eat(token.RPAREN) {
			return nil
		}

		return expr

	default:
		err := fmt.Errorf(
			"%s: syntax error: unexpected %v",
			pos, p.currentToken.Type(),
		)

		p.errors = append(p.errors, err)
		return nil
	}
}

// eat checks that the given tokenType corresponds to the currently held token
// and consumes it. If the tokens do not match, eat panics.
func (p *Parser) eat(tokenType token.TokenTag) bool {
	pos := p.currentPos

	if p.currentToken.Type() == tokenType {
		p.currentToken, p.currentPos = p.lexer.GetNextToken()
		return true
	}

	err := fmt.Errorf(
		"%s: syntax error: expected %v got %v",
		pos, tokenType, p.currentToken.Type(),
	)

	p.errors = append(p.errors, err)

	return false
}

func (p *Parser) skipTo(tokens ...token.TokenTag) {
	for {
		if p.currentToken.Type() == token.EOF {
			return
		}

		for _, t := range tokens {
			if p.currentToken.Type() == t {
				p.currentToken, p.currentPos = p.lexer.GetNextToken()
				return
			}
		}

		p.currentToken, p.currentPos = p.lexer.GetNextToken()
	}
}

func (p *Parser) skipStatement() { p.skipTo(token.SEMI) }

func (p *Parser) skipForBlock() {
	for {
		if p.currentToken.Type() == token.EOF {
			return
		}

		if p.currentToken.Type() == token.FOR {
			p.currentToken, p.currentPos = p.lexer.GetNextToken()
			p.skipForBlock()
			continue
		}

		if p.currentToken.Type() == token.END {
			p.currentToken, p.currentPos = p.lexer.GetNextToken()

			if p.currentToken.Type() == token.FOR {
				p.currentToken, p.currentPos = p.lexer.GetNextToken()

				if p.currentToken.Type() == token.SEMI {
					p.currentToken, p.currentPos = p.lexer.GetNextToken()
					return
				}
			}

			continue
		}

		p.currentToken, p.currentPos = p.lexer.GetNextToken()
	}
}
