package parser

import (
	"fmt"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/token"
)

type Lexer interface{ GetNextToken() token.Token }

// Parser is the main struct of the parser package. The Parser should be
// initialized with New instead of used directly.
type Parser struct {
	lexer        Lexer
	currentToken token.Token
}

// New returns a properly initialized pointer instance to a Parser.
func New(lexer Lexer) *Parser {
	if lexer == nil {
		panic("Attempting to construct a Parser with a nil Lexer")
	}

	return &Parser{
		lexer:        lexer,
		currentToken: lexer.GetNextToken(),
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
		p.eat(token.VAR)
		ident := p.currentToken
		p.eat(token.IDENT)
		p.eat(token.COLON)
		variableType := p.currentToken
		p.eatType()
		if p.currentToken.Type() != token.ASSIGN {
			statement = ast.DeclStmt{
				Identifier:   ident,
				VariableType: variableType,
			}

			break
		}

		p.eat(token.ASSIGN)
		expr := p.parseExpression()

		statement = ast.DeclStmt{
			Identifier:   ident,
			VariableType: variableType,
			Expression:   expr,
		}

	case token.IDENT:
		ident := p.currentToken
		p.eat(token.IDENT)
		p.eat(token.ASSIGN)
		expr := p.parseExpression()

		statement = ast.AssignStmt{
			Identifier: ast.Ident{Id: ident},
			Expression: expr,
		}

	case token.FOR:
		p.eat(token.FOR)
		ident := p.currentToken
		p.eat(token.IDENT)
		p.eat(token.IN)
		low := p.parseExpression()
		p.eat(token.DOTDOT)
		high := p.parseExpression()
		p.eat(token.DO)
		statements := p.parseStatements()
		p.eat(token.END)
		p.eat(token.FOR)
		statement = ast.ForStmt{
			Index:      ast.Ident{Id: ident},
			Low:        low,
			High:       high,
			Statements: statements,
		}

	case token.READ:
		p.eat(token.READ)
		statement = ast.ReadStmt{TargetIdentifier: ast.Ident{Id: p.currentToken}}
		p.eat(token.IDENT)

	case token.PRINT:
		p.eat(token.PRINT)
		statement = ast.PrintStmt{Expression: p.parseExpression()}

	case token.ASSERT:
		p.eat(token.ASSERT)
		p.eat(token.LPAREN)
		statement = ast.AssertStmt{Expression: p.parseExpression()}
		p.eat(token.RPAREN)

	default:
		panic("Parse error")
	}

	p.eat(token.SEMI)

	return statement
}

// parseExpression parses an expression with the following grammar rules.
//
// <expr> ::= <opnd> <op> <opnd>
//            | [ <unary_opnd> ] <opnd>
func (p *Parser) parseExpression() ast.Expr {
	if p.currentToken.Type() == token.NOT {
		unary := p.currentToken
		p.eat(token.NOT)

		return ast.UnaryExpr{
			Unary:   unary,
			Operand: p.parseOperand(),
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
	switch p.currentToken.Type() {
	case token.INTEGER_LITERAL:
		val := p.currentToken.ValueInt()
		p.eat(token.INTEGER_LITERAL)
		return ast.NumberOpnd{Value: val}

	case token.STRING_LITERAL:
		val := p.currentToken.ValueString()
		p.eat(token.STRING_LITERAL)
		return ast.StringOpnd{Value: val}

	case token.IDENT:
		t := p.currentToken
		p.eat(token.IDENT)
		return ast.Ident{Id: t}

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
func (p *Parser) eat(tokenType token.TokenTag) {
	if p.currentToken.Type() == tokenType {
		p.currentToken = p.lexer.GetNextToken()
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
