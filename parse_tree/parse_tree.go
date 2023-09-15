package parse_tree

import "github.com/tivt2/jack-compiler/token"

type Node interface {
	TokenLiteral() string
}

type Declaration interface {
	Node
	declarationNode()
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type Class struct {
	Token          token.Token // the IDENT token
	Value          string      // the class name
	ClassVarDecs   []Declaration
	SubroutineDecs []*SubroutineDec
}

func (c *Class) TokenLiteral() string {
	if len(c.ClassVarDecs) > 0 {
		return c.ClassVarDecs[0].TokenLiteral()
	} else {
		return ""
	}
}

type ClassVarDec struct {
	Token   token.Token // the FIELD or STATIC token
	DecType token.Token
	Name    *Identifier
}

func (cvd *ClassVarDec) declarationNode()     {}
func (cvd *ClassVarDec) TokenLiteral() string { return cvd.Token.Literal }

type Parameter struct {
	DecType token.Token
	Name    *Identifier
}

type SubroutineDec struct {
	Token          token.Token // the CONSTRUCTOR, METHOD or FUNCTION token
	DecType        token.Token
	Name           *Identifier
	Params         []*Parameter
	SubroutineBody *SubroutineBody
}

func (sd *SubroutineDec) TokenLiteral() string { return sd.Token.Literal }

type SubroutineBody struct {
	VarDecs []Declaration
	// Statements []Statement
}

type VarDec struct {
	Token   token.Token
	DecType token.Token
	Name    *Identifier
}

func (vd *VarDec) declarationNode()     {}
func (vd *VarDec) TokenLiteral() string { return vd.Token.Literal }
