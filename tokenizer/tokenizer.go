package tokenizer

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/tivt2/jack-compiler/token"
)

type Tokenizer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
	CurrToken    token.Token
}

func removeComments(text string) string {
	regexes := []string{
		`\/\*[^*]*\*\/`,
		`\/\*\*[\s\S]*?\*\/`,
		`\/\/[^\n]*`,
	}

	for _, pattern := range regexes {
		regex := regexp.MustCompile(pattern)
		text = regex.ReplaceAllString(text, "")
	}

	// fmt.Println(text)
	return text
}

func New(filePath string) *Tokenizer {
	file, err := os.ReadFile(filePath)
	checkErr(err, fmt.Sprintf("Error when opening file %s", filePath))

	fileContent := removeComments(strings.TrimSpace(string(file)))

	tkzr := &Tokenizer{
		input:        fileContent,
		position:     0,
		nextPosition: 0,
	}
	tkzr.ch = tkzr.input[tkzr.position]

	return tkzr
}

func (tkzr *Tokenizer) HasMoreTokens() bool {
	return len(tkzr.input) > tkzr.nextPosition
}

func (tkzr *Tokenizer) Advance() {
	if tkzr.HasMoreTokens() {
		tkzr.nextChar()
		tkzr.CurrToken = tkzr.tokenize()
	} else {
		tkzr.CurrToken = token.Token{Type: token.EOF, Literal: ""}
	}
}

func (tkzr *Tokenizer) TokenType() token.TokenType {
	return tkzr.CurrToken.Type
}

func (tkzr *Tokenizer) KeyWord() (string, bool) {
	if tkzr.CurrToken.Type == token.KEYWORD {
		return strings.ToUpper(tkzr.CurrToken.Literal), true
	}
	return "", false
}

func (tkzr *Tokenizer) Symbol() (string, bool) {
	if tkzr.CurrToken.Type == token.SYMBOL {
		return tkzr.CurrToken.Literal, true
	}
	return "", false
}

func (tkzr *Tokenizer) Identifier() (string, bool) {
	if tkzr.CurrToken.Type == token.IDENTIFIER {
		return tkzr.CurrToken.Literal, true
	}
	return "", false
}

func (tkzr *Tokenizer) IntVal() (int, bool) {
	if tkzr.CurrToken.Type == token.INT_CONST {
		val, err := strconv.Atoi(tkzr.CurrToken.Literal)
		checkErr(err, fmt.Sprintf("Error trying to parse %s to integer", tkzr.CurrToken.Literal))
		return val, true
	}
	return 0, false
}

func (tkzr *Tokenizer) StringVal() (string, bool) {
	if tkzr.CurrToken.Type == token.STRING_CONST {
		return tkzr.CurrToken.Literal, true
	}
	return "", false
}

func (tkzr *Tokenizer) nextChar() {
	tkzr.position = tkzr.nextPosition
	tkzr.nextPosition += 1
	if len(tkzr.input) > tkzr.position {
		tkzr.ch = tkzr.input[tkzr.position]
	}
}

func (tkzr *Tokenizer) tokenize() token.Token {
	tkzr.ignoreWhiteSpace()
	var out token.Token

	if _, ok := token.Symbols[tkzr.ch]; ok {
		switch tkzr.ch {
		case '<':
			out = token.Token{Type: token.SYMBOL, Literal: "&lt;"}
		case '>':
			out = token.Token{Type: token.SYMBOL, Literal: "&gt;"}
		case '&':
			out = token.Token{Type: token.SYMBOL, Literal: "&amp;"}
		default:
			out = token.Token{Type: token.SYMBOL, Literal: string(tkzr.ch)}
		}
	} else if isLetter(tkzr.ch) {
		literal := tkzr.readWord()
		if _, ok := token.Keywords[literal]; ok {
			out.Type = token.KEYWORD
		} else {
			out.Type = token.IDENTIFIER
		}
		out.Literal = literal
		return out
	} else if isDigit(tkzr.ch) {
		literal := tkzr.readInteger()
		out = token.Token{Type: token.INT_CONST, Literal: literal}
		return out
	} else if tkzr.ch == '"' {
		literal := tkzr.readString()
		out = token.Token{Type: token.STRING_CONST, Literal: literal}
	} else {
		out = token.Token{Type: token.ILLEGAL, Literal: string(tkzr.ch)}
	}

	return out
}

func (tkzr *Tokenizer) ignoreWhiteSpace() {
	for tkzr.ch == ' ' || tkzr.ch == '\t' || tkzr.ch == '\n' || tkzr.ch == '\r' {
		tkzr.nextChar()
	}
}

func (tkzr *Tokenizer) PeekCh() byte {
	if len(tkzr.input) > tkzr.nextPosition {
		return tkzr.input[tkzr.nextPosition]
	}
	return 0
}

func (tkzr *Tokenizer) readWord() string {
	startPos := tkzr.position
	for isLetter(tkzr.PeekCh()) || isDigit(tkzr.PeekCh()) {
		tkzr.nextChar()
	}

	return tkzr.input[startPos:tkzr.nextPosition]
}

func (tkzr *Tokenizer) readInteger() string {
	startPos := tkzr.position
	for isDigit(tkzr.PeekCh()) {
		tkzr.nextChar()
	}

	return string(tkzr.input[startPos:tkzr.nextPosition])
}

func (tkzr *Tokenizer) readString() string {
	tkzr.nextChar()
	startPos := tkzr.position
	for tkzr.ch != '"' {
		tkzr.nextChar()
	}

	return tkzr.input[startPos:tkzr.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (tkzr *Tokenizer) Print() {
	for tkzr.CurrToken.Type != token.EOF {
		tkzr.Advance()
		fmt.Println(tkzr.CurrToken)
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
