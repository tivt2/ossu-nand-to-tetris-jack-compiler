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
	currToken    token.Token
}

func removeComments(text string) string {
	regexes := []string{
		`\/\*[^*]*\*\/`,
		`\/\*\*[^*]*\*\/`,
		`\/\/[^\n]*`,
	}

	for _, pattern := range regexes {
		regex := regexp.MustCompile(pattern)
		text = regex.ReplaceAllString(text, "")
	}

	return text
}

func New(filePath string) *Tokenizer {
	file, err := os.ReadFile(filePath)
	checkErr(err, fmt.Sprintf("Error when opening file %s", filePath))

	fileContent := removeComments(string(file))

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
		tkzr.currToken = tkzr.tokenize()
	} else {
		tkzr.currToken = token.Token{Type: token.EOF, Literal: ""}
	}
}

func (tkzr *Tokenizer) TokenType() token.TokenType {
	return tkzr.currToken.Type
}

func (tkzr *Tokenizer) KeyWord() (string, bool) {
	if tkzr.currToken.Type == token.KEYWORD {
		return strings.ToUpper(tkzr.currToken.Literal), true
	}
	return "", false
}

func (tkzr *Tokenizer) Symbol() (string, bool) {
	if tkzr.currToken.Type == token.SYMBOL {
		return tkzr.currToken.Literal, true
	}
	return "", false
}

func (tkzr *Tokenizer) Identifier() (string, bool) {
	if tkzr.currToken.Type == token.IDENTIFIER {
		return tkzr.currToken.Literal, true
	}
	return "", false
}

func (tkzr *Tokenizer) IntVal() (int, bool) {
	if tkzr.currToken.Type == token.INT_CONST {
		val, err := strconv.Atoi(tkzr.currToken.Literal)
		checkErr(err, fmt.Sprintf("Error trying to parse %s to integer", tkzr.currToken.Literal))
		return val, true
	}
	return 0, false
}

func (tkzr *Tokenizer) StringVal() (string, bool) {
	if tkzr.currToken.Type == token.STRING_CONST {
		return tkzr.currToken.Literal, true
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
	var out token.Token
	tkzr.ignoreWhiteSpace()

	if _, ok := token.Symbols[tkzr.ch]; ok {
		out = token.Token{Type: token.SYMBOL, Literal: string(tkzr.ch)}
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

func (tkzr *Tokenizer) peekCh() byte {
	if len(tkzr.input) > tkzr.nextPosition {
		return tkzr.input[tkzr.nextPosition]
	}
	return 0
}

func (tkzr *Tokenizer) readWord() string {
	startPos := tkzr.position
	for isLetter(tkzr.peekCh()) {
		tkzr.nextChar()
	}

	return tkzr.input[startPos:tkzr.nextPosition]
}

func (tkzr *Tokenizer) readInteger() string {
	startPos := tkzr.position
	for isDigit(tkzr.peekCh()) {
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
	for tkzr.currToken.Type != token.EOF {
		tkzr.Advance()
		fmt.Println(tkzr.currToken)
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}
