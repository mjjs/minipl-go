package parser

import (
	"fmt"

	"github.com/mjjs/minipl-go/src/ast"
	"github.com/mjjs/minipl-go/src/lexer"
)

type Lexer interface{ GetNextToken() lexer.Token }

type Parser struct {
	lexer        Lexer
	currentToken lexer.Token
}

func New(lexer Lexer) *Parser {
	if lexer == nil {
		panic("Attempting to construct a Parser with a nil Lexer")
	}

	return &Parser{
		lexer:        lexer,
		currentToken: lexer.GetNextToken(),
	}
}

func (p *Parser) Parse() (ast.Stmts, error) {
	root := p.parseStatements()

	if p.currentToken.Type() != lexer.EOF {
		return ast.Stmts{}, fmt.Errorf(
			"Parsing the program failed. Expected %v, found %v",
			lexer.EOF,
			p.currentToken.Type(),
		)
	}

	return root, nil
}

func (p *Parser) parseStatements() ast.Stmts {
	statements := []ast.Stmt{}

	statements = append(statements, p.parseStatement())

	for p.isStatement(p.currentToken.Type()) {
		statements = append(statements, p.parseStatement())
	}

	return ast.Stmts{Statements: statements}
}

func (p *Parser) parseStatement() ast.Stmt {
	var statement ast.Stmt

	switch p.currentToken.Type() {
	case lexer.VAR:
		p.eat(lexer.VAR)
		ident := p.currentToken
		p.eat(lexer.IDENT)
		p.eat(lexer.COLON)
		variableType := p.currentToken
		p.eatType()
		if p.currentToken.Type() != lexer.ASSIGN {
			statement = ast.DeclStmt{
				Identifier:   ident,
				VariableType: variableType,
			}

			break
		}

		p.eat(lexer.ASSIGN)
		expr := p.parseExpression()

		statement = ast.DeclStmt{
			Identifier:   ident,
			VariableType: variableType,
			Expression:   expr,
		}

	case lexer.IDENT:
		ident := p.currentToken
		p.eat(lexer.IDENT)
		p.eat(lexer.ASSIGN)
		expr := p.parseExpression()

		statement = ast.AssignStmt{
			Identifier: ast.Ident{Id: ident},
			Expression: expr,
		}

	case lexer.FOR:
		p.eat(lexer.FOR)
		ident := p.currentToken
		p.eat(lexer.IDENT)
		p.eat(lexer.IN)
		low := p.parseExpression()
		p.eat(lexer.DOTDOT)
		high := p.parseExpression()
		p.eat(lexer.DO)
		statements := p.parseStatements()
		p.eat(lexer.END)
		p.eat(lexer.FOR)
		statement = ast.ForStmt{
			Index:      ast.Ident{Id: ident},
			Low:        low,
			High:       high,
			Statements: statements,
		}

	case lexer.READ:
		p.eat(lexer.READ)
		statement = ast.ReadStmt{TargetIdentifier: ast.Ident{Id: p.currentToken}}
		p.eat(lexer.IDENT)

	case lexer.PRINT:
		p.eat(lexer.PRINT)
		statement = ast.PrintStmt{Expression: p.parseExpression()}

	case lexer.ASSERT:
		p.eat(lexer.ASSERT)
		p.eat(lexer.LPAREN)
		statement = ast.AssertStmt{Expression: p.parseExpression()}
		p.eat(lexer.RPAREN)

	default:
		panic("Parse error")
	}

	p.eat(lexer.SEMI)

	return statement
}

func (p *Parser) parseExpression() ast.Expr {
	if p.currentToken.Type() == lexer.NOT {
		unary := p.currentToken
		p.eat(lexer.NOT)

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

func (p *Parser) parseOperand() ast.Node {
	switch p.currentToken.Type() {
	case lexer.INTEGER_LITERAL:
		val := p.currentToken.ValueInt()
		p.eat(lexer.INTEGER_LITERAL)
		return ast.NumberOpnd{Value: val}

	case lexer.STRING_LITERAL:
		val := p.currentToken.ValueString()
		p.eat(lexer.STRING_LITERAL)
		return ast.StringOpnd{Value: val}

	case lexer.IDENT:
		token := p.currentToken
		p.eat(lexer.IDENT)
		return ast.Ident{Id: token}

	case lexer.LPAREN:
		p.eat(lexer.LPAREN)
		expr := p.parseExpression()
		p.eat(lexer.RPAREN)
		return expr
	}

	panic(fmt.Sprintf("Syntax error: unexpected %v", p.currentToken.Type()))
}

func (p *Parser) isStatement(tokenType lexer.TokenTag) bool {
	statementTypes := []lexer.TokenTag{
		lexer.VAR,
		lexer.IDENT,
		lexer.FOR,
		lexer.READ,
		lexer.PRINT,
		lexer.ASSERT,
	}

	for _, t := range statementTypes {
		if t == tokenType {
			return true
		}
	}

	return false
}

func (p *Parser) isOperator(tokenType lexer.TokenTag) bool {
	operatorTypes := []lexer.TokenTag{
		lexer.PLUS,
		lexer.MINUS,
		lexer.MULTIPLY,
		lexer.INTEGER_DIV,
		lexer.LT,
		lexer.EQ,
		lexer.AND,
	}

	for _, t := range operatorTypes {
		if t == tokenType {
			return true
		}
	}

	return false
}

func (p *Parser) eat(tokenType lexer.TokenTag) {
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

func (p *Parser) eatType() {
	if p.currentToken.Type() == lexer.INTEGER {
		p.eat(lexer.INTEGER)
	} else if p.currentToken.Type() == lexer.STRING {
		p.eat(lexer.STRING)
	} else if p.currentToken.Type() == lexer.BOOLEAN {
		p.eat(lexer.BOOLEAN)
	} else {
		panic(fmt.Sprintf(
			"Syntax error: expected a type, got %v",
			p.currentToken.Type(),
		))
	}
}

func (p *Parser) eatOperator() {
	if p.currentToken.Type() == lexer.PLUS {
		p.eat(lexer.PLUS)
	} else if p.currentToken.Type() == lexer.MINUS {
		p.eat(lexer.MINUS)
	} else if p.currentToken.Type() == lexer.MULTIPLY {
		p.eat(lexer.MULTIPLY)
	} else if p.currentToken.Type() == lexer.INTEGER_DIV {
		p.eat(lexer.INTEGER_DIV)
	} else if p.currentToken.Type() == lexer.LT {
		p.eat(lexer.LT)
	} else if p.currentToken.Type() == lexer.EQ {
		p.eat(lexer.EQ)
	} else if p.currentToken.Type() == lexer.AND {
		p.eat(lexer.AND)
	} else {
		panic(fmt.Sprintf(
			"Syntax error: expected an operator, got %v",
			p.currentToken.Type(),
		))
	}
}
