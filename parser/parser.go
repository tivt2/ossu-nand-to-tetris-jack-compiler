package parser

import (
	"bytes"
	"fmt"
	"log"

	"github.com/tivt2/jack-compiler/token"
	"github.com/tivt2/jack-compiler/tokenizer"
)

type Parser struct {
	tkzr *tokenizer.Tokenizer

	TokensXML bytes.Buffer
	ParseTree bytes.Buffer
}

func New(tkzr *tokenizer.Tokenizer) *Parser {
	parser := &Parser{
		tkzr: tkzr,
	}

	return parser
}

func (p *Parser) CompileClass() (ParseTree string, TokensXML string) {
	p.tkzr.Advance()
	if p.tkzr.CurrToken.Literal != "class" {
		log.Fatalf("Invalid jack class. received token: %v", p.tkzr.CurrToken)
	}
	p.writeRule("class")
	p.TokensXML.WriteString("<tokens>\n")

	p.eatKeywordOrSymbol("class")
	p.eatIdentifier()
	p.eatKeywordOrSymbol("{")

	p.compileClassVarDec()

	p.compileSubroutine()

	p.eatKeywordOrSymbol("}")
	p.writeRule("/class")
	p.TokensXML.WriteString("</tokens>")
	return p.ParseTree.String(), p.TokensXML.String()
}

func (p *Parser) compileClassVarDec() {
	switch p.tkzr.CurrToken.Literal {
	case "field":
		p.writeRule("classVarDec")
		p.eatKeywordOrSymbol("field")
	case "static":
		p.writeRule("classVarDec")
		p.eatKeywordOrSymbol("static")
	default:
		return
	}

	p.eatType()
	p.eatIdentifier()

	for p.tkzr.CurrToken.Literal == "," {
		p.eatKeywordOrSymbol(",")
		p.eatIdentifier()
	}

	p.eatKeywordOrSymbol(";")
	p.writeRule("/classVarDec")
	p.compileClassVarDec()
}

func (p *Parser) compileSubroutine() {
	switch p.tkzr.CurrToken.Literal {
	case "constructor":
		p.writeRule("subroutineDec")
		p.eatKeywordOrSymbol("constructor")
	case "method":
		p.writeRule("subroutineDec")
		p.eatKeywordOrSymbol("method")
	case "function":
		p.writeRule("subroutineDec")
		p.eatKeywordOrSymbol("function")
	default:
		return
	}

	p.eatType()
	p.eatIdentifier()
	p.eatKeywordOrSymbol("(")
	p.writeRule("parameterList")
	p.compileParameterList()
	p.writeRule("/parameterList")
	p.eatKeywordOrSymbol(")")
	p.writeRule("subroutineBody")
	p.compileSubroutineBody()
	p.writeRule("/subroutineBody")
	p.writeRule("/subroutineDec")

	p.compileSubroutine()
}

func (p *Parser) compileParameterList() {
	if p.tkzr.CurrToken.Literal == ")" {
		return
	}
	if p.tkzr.CurrToken.Literal == "," {
		p.eatKeywordOrSymbol(",")
	}
	p.eatType()
	p.eatIdentifier()
	p.compileParameterList()
}

func (p *Parser) compileSubroutineBody() {
	p.eatKeywordOrSymbol("{")
	p.compileVarDec()
	p.writeRule("statements")
	p.compileStatements()
	p.writeRule("/statements")

	p.eatKeywordOrSymbol("}")
}

func (p *Parser) compileVarDec() {
	if p.tkzr.CurrToken.Literal != "var" {
		return
	}
	p.writeRule("varDec")
	p.eatKeywordOrSymbol("var")
	p.eatType()
	p.eatIdentifier()
	for p.tkzr.CurrToken.Literal == "," {
		p.eatKeywordOrSymbol(",")
		p.eatIdentifier()
	}
	p.eatKeywordOrSymbol(";")
	p.writeRule("/varDec")
	p.compileVarDec()
}

func (p *Parser) compileStatements() {
	switch p.tkzr.CurrToken.Literal {
	case "let":
		p.writeRule("letStatement")
		p.compileLet()
		p.writeRule("/letStatement")
	case "if":
		p.writeRule("ifStatement")
		p.compileIf()
		p.writeRule("/ifStatement")
	case "while":
		p.writeRule("whileStatement")
		p.compileWhile()
		p.writeRule("/whileStatement")
	case "do":
		p.writeRule("doStatement")
		p.compileDo()
		p.writeRule("/doStatement")
	case "return":
		p.writeRule("returnStatement")
		p.compileReturn()
		p.writeRule("/returnStatement")
	default:
		return
	}
	p.compileStatements()
}

func (p *Parser) compileLet() {
	p.eatKeywordOrSymbol("let")
	p.eatIdentifier()

	switch p.tkzr.CurrToken.Literal {
	case "[":
		p.eatKeywordOrSymbol("[")
		p.writeRule("expression")
		p.compileExpression()
		p.writeRule("/expression")
		p.eatKeywordOrSymbol("]")
	case ".":
		p.eatKeywordOrSymbol(".")
		p.writeRule("expression")
		p.compileExpression()
		p.writeRule("/expression")
	}

	p.eatKeywordOrSymbol("=")
	p.writeRule("expression")
	p.compileExpression()
	p.writeRule("/expression")
	p.eatKeywordOrSymbol(";")
}

func (p *Parser) compileIf() {
	p.eatKeywordOrSymbol("if")
	p.eatKeywordOrSymbol("(")

	p.writeRule("expression")
	p.compileExpression()
	p.writeRule("/expression")

	p.eatKeywordOrSymbol(")")
	p.eatKeywordOrSymbol("{")

	p.writeRule("statements")
	p.compileStatements()
	p.writeRule("/statements")

	p.eatKeywordOrSymbol("}")

	if p.tkzr.CurrToken.Literal == "else" {
		p.eatKeywordOrSymbol("else")
		p.eatKeywordOrSymbol("{")

		p.writeRule("statements")
		p.compileStatements()
		p.writeRule("/statements")

		p.eatKeywordOrSymbol("}")
	}
}

func (p *Parser) compileWhile() {
	p.eatKeywordOrSymbol("while")
	p.eatKeywordOrSymbol("(")

	p.writeRule("expression")
	p.compileExpression()
	p.writeRule("/expression")

	p.eatKeywordOrSymbol(")")
	p.eatKeywordOrSymbol("{")

	p.writeRule("statements")
	p.compileStatements()
	p.writeRule("/statements")

	p.eatKeywordOrSymbol("}")
}

func (p *Parser) compileDo() {
	p.eatKeywordOrSymbol("do")
	p.eatIdentifier()

	if p.tkzr.CurrToken.Literal == "[" {
		p.eatKeywordOrSymbol("[")
		p.writeRule("expression")
		p.compileExpression()
		p.writeRule("/expression")
		p.eatKeywordOrSymbol("]")
	}

	if p.tkzr.CurrToken.Literal == "." {
		p.eatKeywordOrSymbol(".")
		p.eatIdentifier()
	}

	p.eatKeywordOrSymbol("(")
	p.writeRule("expressionList")
	p.compileExpressionList()
	p.writeRule("expressionList")
	p.eatKeywordOrSymbol(")")
	p.eatKeywordOrSymbol(";")
}

func (p *Parser) compileReturn() {
	p.eatKeywordOrSymbol("return")

	if p.tkzr.CurrToken.Literal != ";" {
		p.writeRule("expression")
		p.compileExpression()
		p.writeRule("/expression")
	}

	p.eatKeywordOrSymbol(";")
}

func (p *Parser) compileExpression() {
	p.writeRule("term")
	p.compileTerm()
	p.writeRule("/term")

	switch p.tkzr.CurrToken.Literal {
	case "+":
		p.eatKeywordOrSymbol("+")
	case "-":
		p.eatKeywordOrSymbol("-")
	case "*":
		p.eatKeywordOrSymbol("*")
	case "/":
		p.eatKeywordOrSymbol("/")
	case "~":
		p.eatKeywordOrSymbol("~")
	case "=":
		p.eatKeywordOrSymbol("=")
	case "&lt;":
		p.eatKeywordOrSymbol("<")
	case "&gt;":
		p.eatKeywordOrSymbol(">")
	case "&amp;":
		p.eatKeywordOrSymbol("&")
	case "|":
		p.eatKeywordOrSymbol("|")
	case ";":
		return
	case ")":
		return
	case "]":
		return
	case ",":
		return
	}

	p.writeRule("term")
	p.compileTerm()
	p.writeRule("/term")
}

func (p *Parser) compileTerm() {
	// if p.tkzr.CurrToken.Literal == ";" {
	// 	return
	// }

	switch p.tkzr.CurrToken.Type {
	case token.INT_CONST:
		p.eatConstant()
	case token.STRING_CONST:
		p.eatConstant()
	case token.KEYWORD:
		switch p.tkzr.CurrToken.Literal {
		case "true":
			p.eatKeywordOrSymbol("true")
		case "false":
			p.eatKeywordOrSymbol("false")
		case "null":
			p.eatKeywordOrSymbol("null")
		case "this":
			p.eatKeywordOrSymbol("this")
		}
	case token.SYMBOL:
		switch p.tkzr.CurrToken.Literal {
		case "-":
			p.eatKeywordOrSymbol("-")
			p.writeRule("term")
			p.compileTerm()
			p.writeRule("/term")
		case "~":
			p.eatKeywordOrSymbol("~")
			p.writeRule("term")
			p.compileTerm()
			p.writeRule("/term")
		case "(":
			p.eatKeywordOrSymbol("(")
			p.writeRule("expression")
			p.compileExpression()
			p.writeRule("/expression")
			p.eatKeywordOrSymbol(")")
		}
	case token.IDENTIFIER:
		p.eatIdentifier()
		switch p.tkzr.CurrToken.Literal {
		case "[":
			p.eatKeywordOrSymbol("[")
			p.writeRule("expression")
			p.compileExpression()
			p.writeRule("/expression")
			p.eatKeywordOrSymbol("]")
		case "(":
			p.eatKeywordOrSymbol("(")
			p.writeRule("expressionList")
			p.compileExpressionList()
			p.writeRule("/expressionList")
			p.eatKeywordOrSymbol(")")
		case ".":
			p.eatKeywordOrSymbol(".")
			p.eatIdentifier()
			p.eatKeywordOrSymbol("(")
			p.writeRule("expressionList")
			p.compileExpressionList()
			p.writeRule("/expressionList")
			p.eatKeywordOrSymbol(")")
		}
	}
}

func (p *Parser) compileExpressionList() {
	if p.tkzr.CurrToken.Literal == ")" {
		return
	}
	if p.tkzr.CurrToken.Literal == "," {
		p.eatKeywordOrSymbol(",")
	}
	p.writeRule("expression")
	p.compileExpression()
	p.writeRule("/expression")
	p.compileExpressionList()
}

func (p *Parser) eatType() {
	literal := p.tkzr.CurrToken.Literal
	if p.tkzr.CurrToken.Type != token.IDENTIFIER && literal != "int" && literal != "char" && literal != "boolean" && literal != "void" {
		log.Fatalf("Unexpected type, expected: %v. received: %v", token.IDENTIFIER, p.tkzr.CurrToken.Type)
	}
	p.writeToken()
	p.tkzr.Advance()
}

func (p *Parser) eatIdentifier() {
	if p.tkzr.CurrToken.Type != token.IDENTIFIER {
		log.Fatalf("Unexpected identifier, expected: %v. received: %v", token.IDENTIFIER, p.tkzr.CurrToken.Type)
	}
	p.writeToken()
	p.tkzr.Advance()
}

func (p *Parser) eatKeywordOrSymbol(literal string) {
	if p.tkzr.CurrToken.Type != token.KEYWORD && p.tkzr.CurrToken.Type != token.SYMBOL && p.tkzr.CurrToken.Literal != literal {
		log.Fatalf("Unexpected token, expected: %v. received: %v", literal, p.tkzr.CurrToken.Literal)
	}
	p.writeToken()
	p.tkzr.Advance()
}

func (p *Parser) eatConstant() {
	if p.tkzr.CurrToken.Type != token.INT_CONST && p.tkzr.CurrToken.Type != token.STRING_CONST {
		log.Fatalf("Unexpected token, expected: %v or %v. received: %v", token.INT_CONST, token.STRING_CONST, p.tkzr.CurrToken.Type)
	}
	p.writeToken()
	p.tkzr.Advance()
}

func (p *Parser) writeRule(rule string) {
	// fmt.Printf("<%v>\n", rule)
	p.ParseTree.WriteString(fmt.Sprintf("<%s>\n", rule))
}

func (p *Parser) writeToken() {
	// fmt.Printf("<%v> %v </%v>\n", p.tkzr.CurrToken.Type, p.tkzr.CurrToken.Literal, p.tkzr.CurrToken.Type)
	p.ParseTree.WriteString(fmt.Sprintf("<%v> %v </%v>\n", p.tkzr.CurrToken.Type, p.tkzr.CurrToken.Literal, p.tkzr.CurrToken.Type))
	p.TokensXML.WriteString(fmt.Sprintf("<%v> %v </%v>\n", p.tkzr.CurrToken.Type, p.tkzr.CurrToken.Literal, p.tkzr.CurrToken.Type))
}
