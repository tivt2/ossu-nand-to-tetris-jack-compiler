package tokenizer

import (
	"log"

	"github.com/tivt2/jack-compiler/token"
)

type Tokenizer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Tokenizer {
	tkzr := &Tokenizer{input: input}
	tkzr.readChar()
	return tkzr
}

func (tkzr *Tokenizer) readChar() {
	if tkzr.readPosition >= len(tkzr.input) {
		tkzr.ch = 0
	} else {
		tkzr.ch = tkzr.input[tkzr.readPosition]
	}
	tkzr.position = tkzr.readPosition
	tkzr.readPosition += 1
}

func (tkzr *Tokenizer) Advance() token.Token {
	var out token.Token

	tkzr.ignoreWithSpace()

	switch tkzr.ch {
	case '=':
		out = newToken(token.EQ, tkzr.ch)
	case '[':
		out = newToken(token.LBRACKET, tkzr.ch)
	case ']':
		out = newToken(token.RBRACKET, tkzr.ch)
	case '(':
		out = newToken(token.LPAREN, tkzr.ch)
	case ')':
		out = newToken(token.RPAREN, tkzr.ch)
	case '{':
		out = newToken(token.LBRACE, tkzr.ch)
	case '}':
		out = newToken(token.RBRACE, tkzr.ch)
	case '.':
		out = newToken(token.DOT, tkzr.ch)
	case ',':
		out = newToken(token.COMMA, tkzr.ch)
	case ';':
		out = newToken(token.SEMICOLON, tkzr.ch)
	case '+':
		out = newToken(token.PLUS, tkzr.ch)
	case '-':
		out = newToken(token.MINUS, tkzr.ch)
	case '*':
		out = newToken(token.ASTERISK, tkzr.ch)
	case '/':
		out = newToken(token.FSLASH, tkzr.ch)
	case '&':
		out = newToken(token.AMP, tkzr.ch)
	case '|':
		out = newToken(token.BAR, tkzr.ch)
	case '<':
		out = newToken(token.LT, tkzr.ch)
	case '>':
		out = newToken(token.GT, tkzr.ch)
	case '~':
		out = newToken(token.NOT, tkzr.ch)
	case '"':
		out = newToken(token.QUOT, tkzr.ch)
	case 0:
		out.Literal = ""
		out.Type = token.EOF
	default:
		if isLetter(tkzr.ch) {
			out.Literal = tkzr.readIdentifier()
			out.Type = token.LookupIdent(out.Literal)
			return out
		} else if isDigit(tkzr.ch) {
			out.Type = token.INT
			out.Literal = tkzr.readNumber()
			return out
		} else {
			out = newToken(token.ILLEGAL, tkzr.ch)
		}
	}

	tkzr.readChar()
	return out
}

func (tkzr *Tokenizer) ignoreWithSpace() {
	for tkzr.ch == ' ' || tkzr.ch == '\r' || tkzr.ch == '\t' || tkzr.ch == '\n' {
		tkzr.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (tkzr *Tokenizer) readIdentifier() string {
	position := tkzr.position
	for isLetter(tkzr.ch) || isDigit(tkzr.ch) {
		tkzr.readChar()
	}
	return tkzr.input[position:tkzr.position]
}

func (tkzr *Tokenizer) readNumber() string {
	position := tkzr.position
	for isDigit(tkzr.ch) {
		tkzr.readChar()
	}
	return tkzr.input[position:tkzr.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
