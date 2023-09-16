package parser

import (
	"log"
	"strconv"

	"github.com/tivt2/jack-compiler/parse_tree"
	"github.com/tivt2/jack-compiler/token"
	"github.com/tivt2/jack-compiler/tokenizer"
)

type Parser struct {
	tkzr *tokenizer.Tokenizer

	curToken  token.Token
	peekToken token.Token
}

func New(tkzr *tokenizer.Tokenizer) *Parser {
	p := &Parser{tkzr: tkzr}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.tkzr.Advance()
}

func (p *Parser) expectToken(tokenType token.TokenType) bool {
	if p.curToken.Type == tokenType {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekToken.Type == tokenType {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) parseStatements() []parse_tree.Statement {
	var stmts []parse_tree.Statement

	for p.curToken.Type != token.RBRACE {
		stmts = append(stmts, p.parseStatement())
		p.nextToken()
	}

	return stmts
}

func (p *Parser) parseStatement() parse_tree.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.DO:
		return p.parseDoStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *parse_tree.LetStatement {
	ls := &parse_tree.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		log.Fatalf("Invalid let statement, missing ident, received: %v", p.curToken)
	}

	ls.Ident = &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.expectPeek(token.LBRACKET) {
		p.nextToken()
		ls.Ident.Indexer = p.parseExpression()
		p.nextToken()
	}

	if !p.expectPeek(token.ASSIGN) {
		log.Fatalf("Invalid let statement, missing assign, received: %v", p.curToken)
	}

	p.nextToken()
	ls.Expression = p.parseExpression()
	if !p.expectPeek(token.SEMICOLON) {
		log.Fatalf("Invalid let statement, missing semicolon, received: %v", p.curToken)
	}

	// p.nextToken()
	return ls
}

func (p *Parser) parseReturnStatement() *parse_tree.ReturnStatement {
	rs := &parse_tree.ReturnStatement{Token: p.curToken}

	p.nextToken()
	rs.Expression = p.parseExpression()
	if !p.expectPeek(token.SEMICOLON) {
		log.Fatalf("Invalid return statement, missing semicolon, received: %v", p.curToken)
	}

	// p.nextToken()
	return rs
}

func (p *Parser) parseDoStatement() *parse_tree.DoStatement {
	ds := &parse_tree.DoStatement{Token: p.curToken}

	p.nextToken()
	ds.Expression = p.parseExpression()
	if !p.expectPeek(token.SEMICOLON) {
		log.Fatalf("Invalid do statement, missing semicolon, received: %v", p.curToken)
	}

	// p.nextToken()
	return ds
}

func (p *Parser) parseIfStatement() *parse_tree.IfStatement {
	is := &parse_tree.IfStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		log.Fatalf("Invalid if statement, missing (, received: %v", p.curToken)
	}

	p.nextToken()
	is.Expression = p.parseExpression()
	if !p.expectPeek(token.RPAREN) {
		log.Fatalf("Invalid if statement, missing ), received: %v", p.curToken)
	}

	if !p.expectPeek(token.LBRACE) {
		log.Fatalf("Invalid if statement, missing {, received: %v", p.curToken)
	}

	p.nextToken()
	is.IfStmts = p.parseStatements()
	if p.curToken.Type != token.RBRACE {
		log.Fatalf("Invalid if statement, missing }, received: %v", p.curToken)
	}

	if p.expectPeek(token.ELSE) {
		if !p.expectPeek(token.LBRACE) {
			log.Fatalf("Invalid else statement, missing {, received: %v", p.curToken)
		}
		p.nextToken()
		is.Else = p.parseStatements()
		if p.curToken.Type != token.RBRACE {
			log.Fatalf("Invalid else statement, missing }, received: %v", p.curToken)
		}
	}

	// p.nextToken()
	return is
}

func (p *Parser) parseWhileStatement() *parse_tree.WhileStatement {
	is := &parse_tree.WhileStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		log.Fatalf("Invalid while statement, missing (, received: %v", p.curToken)
	}

	p.nextToken()
	is.Expression = p.parseExpression()
	if !p.expectPeek(token.RPAREN) {
		log.Fatalf("Invalid while statement, missing ), received: %v", p.curToken)
	}

	if !p.expectPeek(token.LBRACE) {
		log.Fatalf("Invalid while statement, missing {, received: %v", p.curToken)
	}

	p.nextToken()
	is.Stmts = p.parseStatements()
	if p.curToken.Type != token.RBRACE {
		log.Fatalf("Invalid while statement, missing }, received: %v", p.curToken)
	}

	// p.nextToken()
	return is
}

func (p *Parser) parseExpression() parse_tree.Expression {
	exp := p.parseTerm()
	for {
		switch p.peekToken.Type {
		case token.PLUS, token.MINUS, token.ASTERISK, token.FSLASH, token.ASSIGN, token.LT, token.GT, token.AMP, token.BAR:
			p.nextToken()
			op := p.curToken
			p.nextToken()
			exp2 := p.parseTerm()
			exp = &parse_tree.Infix{Token: op, Operator: op.Literal, Left: exp, Right: exp2}
		default:
			return exp
		}
	}
}

func (p *Parser) parseTerm() parse_tree.Expression {
	switch p.curToken.Type {
	case token.MINUS, token.NOT:
		return p.parsePrefix()
	case token.INT:
		return p.parseIntegerConstant()
	case token.QUOT:
		return p.parseStringConstant()
	case token.TRUE, token.FALSE, token.NULL, token.THIS:
		return p.parseKeywordConstant()
	case token.LPAREN:
		p.nextToken()
		exp := p.parseExpression()
		if p.expectToken(token.RPAREN) {
			log.Fatalf("Invalid group expression, missing ), received: %v", p.curToken)
		}
		p.nextToken()
		return exp
	case token.IDENT:
		initIdent := p.curToken
		index := p.checkForIndex()
		if index != nil {
			p.nextToken()
		}

		switch p.peekToken.Type {
		case token.DOT:
			p.nextToken()
			if !p.expectPeek(token.IDENT) {
				log.Fatalf("Invalid dot call, missing 2nd ident, received: %v", p.curToken)
			}
			secondIdent := p.curToken
			return &parse_tree.SubroutineCall{
				Ident:      &parse_tree.Identifier{Token: initIdent, Value: initIdent.Literal, Indexer: index},
				Subroutine: &parse_tree.Identifier{Token: secondIdent, Value: secondIdent.Literal, Indexer: nil},
				ExpList:    p.parseExpressionList(),
			}
		case token.LPAREN:
			return &parse_tree.SubroutineCall{
				Ident:      nil,
				Subroutine: &parse_tree.Identifier{Token: initIdent, Value: initIdent.Literal, Indexer: nil},
				ExpList:    p.parseExpressionList(),
			}
		default:
			return &parse_tree.Identifier{Token: initIdent, Value: initIdent.Literal, Indexer: index}
		}
	default:
		log.Fatalf("Invalid term received: %v", p.curToken)
		return nil
	}
}

func (p *Parser) checkForIndex() parse_tree.Expression {
	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	p.nextToken()
	exp := p.parseExpression()
	if p.expectToken(token.RBRACKET) {
		log.Fatalf("Invalid index expression, missing ], received: %v", p.curToken)
	}
	return exp
}

func (p *Parser) parseExpressionList() []parse_tree.Expression {
	if !p.expectPeek(token.LPAREN) {
		return []parse_tree.Expression{}
	}
	p.nextToken()
	list := []parse_tree.Expression{}
	for p.curToken.Type != token.RPAREN {
		p.expectToken(token.COMMA)
		list = append(list, p.parseExpression())
		p.nextToken()
	}

	return list
}

func (p *Parser) parsePrefix() parse_tree.Expression {
	exp := &parse_tree.Prefix{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Expression = p.parseTerm()
	return exp
}

func (p *Parser) parseIntegerConstant() parse_tree.Expression {
	val, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		log.Fatalf("Error while converting %s to integer", p.curToken.Literal)
	}

	return &parse_tree.IntegerConstant{Token: p.curToken, Value: val}
}

func (p *Parser) parseKeywordConstant() parse_tree.Expression {
	return &parse_tree.KeywordConstant{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStringConstant() parse_tree.Expression {
	return &parse_tree.StringConstant{Token: p.curToken, Value: p.curToken.Literal}
}
