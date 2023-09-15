package parser

import (
	"log"

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

func (p *Parser) ParseClass() *parse_tree.Class {
	class := &parse_tree.Class{}
	if !p.curTokenIs(token.CLASS) {
		log.Fatal("Invalid class")
	}
	class.Token = p.curToken
	if !p.expectPeek(token.IDENT) {
		log.Fatal("Invalid class identifier")
	}
	class.Value = p.curToken.Literal
	if !p.expectPeek(token.LBRACE) {
		log.Fatal("Invalid class open brace")
	}
	p.nextToken()
	class.ClassVarDecs = []parse_tree.Declaration{}

	for p.curToken.Type == token.FIELD || p.curToken.Type == token.STATIC {
		class.ClassVarDecs = p.parseClassVarDec(class.ClassVarDecs)
		p.nextToken()
	}

	class.SubroutineDecs = []*parse_tree.SubroutineDec{}
	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		class.SubroutineDecs = append(class.SubroutineDecs, p.parseSubroutineDec())
		p.nextToken()
	}

	return class
}

func (p *Parser) parseClassVarDec(cvds []parse_tree.Declaration) []parse_tree.Declaration {
	cvd := &parse_tree.ClassVarDec{}

	switch p.curToken.Type {
	case token.FIELD, token.STATIC:
		cvd.Token = p.curToken
	default:
		log.Fatalf("Invalid class var dec, received: %v", p.curToken)
	}
	p.nextToken()

	switch p.curToken.Type {
	case token.CHAR, token.INT, token.BOOLEAN, token.IDENT:
		cvd.DecType = p.curToken
	default:
		log.Fatalf("Invalid class var dec type, received: %v", p.curToken)
	}

	if !p.expectPeek(token.IDENT) {
		log.Fatalf("Invalid class var dec ident, received: %v", p.curToken)
	}

	cvd.Name = &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	cvds = append(cvds, cvd)
	p.nextToken()

	for p.curToken.Type == token.COMMA {
		if !p.expectPeek(token.IDENT) {
			log.Fatalf("Invalid VarDec Identifier, received: %v", p.curToken)
		}
		newCvd := &parse_tree.ClassVarDec{
			Token:   cvd.Token,
			DecType: cvd.DecType,
			Name:    &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		cvds = append(cvds, newCvd)
		p.nextToken()
	}

	if p.curToken.Type != token.SEMICOLON {
		log.Fatalf("Missing semicolon in %q", cvds)
	}

	return cvds
}

func (p *Parser) parseSubroutineDec() *parse_tree.SubroutineDec {
	sd := &parse_tree.SubroutineDec{}

	switch p.curToken.Type {
	case token.CONSTRUCTOR, token.METHOD, token.FUNCTION:
		sd.Token = p.curToken
	default:
		log.Fatalf("Invalid subroutine dec, received: %v", p.curToken)
	}
	p.nextToken()

	switch p.curToken.Type {
	case token.CHAR, token.INT, token.BOOLEAN, token.IDENT:
		sd.DecType = p.curToken
	default:
		log.Fatalf("Invalid subroutine dec type, received: %v", p.curToken)
	}

	if !p.expectPeek(token.IDENT) {
		log.Fatalf("Invalid class var dec ident, received: %v", p.curToken)
	}

	sd.Name = &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		log.Fatalf("Invalid subroutine dec params, received: %v", p.curToken)
	}
	p.nextToken()

	for p.curToken.Type != token.RPAREN {
		sd.Params = append(sd.Params, p.parseParameter())
		p.nextToken()
	}

	sd.SubroutineBody = p.parseSubroutineBody()
	return sd
}

func (p *Parser) parseParameter() *parse_tree.Parameter {
	param := &parse_tree.Parameter{}

	if p.curTokenIs(token.COMMA) {
		p.nextToken()
	}

	switch p.curToken.Type {
	case token.CHAR, token.INT, token.BOOLEAN, token.IDENT:
		param.DecType = p.curToken
	default:
		log.Fatalf("Invalid param dec type, received: %v", p.curToken)
	}

	if !p.expectPeek(token.IDENT) {
		log.Fatalf("Invalid param identifier, received: %v", p.curToken)
	}
	param.Name = &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return param
}

func (p *Parser) parseSubroutineBody() *parse_tree.SubroutineBody {
	if !p.expectPeek(token.LBRACE) {
		log.Fatalf("Invalid subroutinebody, received: %v", p.curToken)
	}
	p.nextToken()
	sb := &parse_tree.SubroutineBody{}

	for p.curToken.Type == token.VAR {
		sb.VarDecs = p.parseVarDec(sb.VarDecs)
		p.nextToken()
	}

	return sb
}

func (p *Parser) parseVarDec(vds []parse_tree.Declaration) []parse_tree.Declaration {
	vd := &parse_tree.VarDec{Token: p.curToken}
	p.nextToken()

	switch p.curToken.Type {
	case token.CHAR, token.INT, token.BOOLEAN, token.IDENT:
		vd.DecType = p.curToken
	default:
		log.Fatalf("Invalid VarDec DecType, received: %v", p.curToken)
	}

	if !p.expectPeek(token.IDENT) {
		log.Fatalf("Invalid VarDec Identifier, received: %v", p.curToken)
	}
	vd.Name = &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	vds = append(vds, vd)
	p.nextToken()

	for p.curToken.Type == token.COMMA {
		if !p.expectPeek(token.IDENT) {
			log.Fatalf("Invalid VarDec Identifier, received: %v", p.curToken)
		}
		newVd := &parse_tree.VarDec{
			Token:   vd.Token,
			DecType: vd.DecType,
			Name:    &parse_tree.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		vds = append(vds, newVd)
		p.nextToken()
	}

	if p.curToken.Type != token.SEMICOLON {
		log.Fatalf("Invalid VarDec semicolon, received: %q", p.curToken)
	}

	return vds
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
