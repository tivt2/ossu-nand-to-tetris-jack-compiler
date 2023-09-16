package parse_tree

import (
	"bytes"

	"github.com/tivt2/jack-compiler/token"
)

type Declaration interface {
	declarationNode()
}

type Class struct {
	ClassVarDecs   []Declaration
	SubroutineDecs []Declaration
}

type ClassVarDec struct {
	Kind    string
	DecType string
	Name    string
}

func (cvd *ClassVarDec) declarationNode() {}

type SubroutineDec struct {
	Kind           string
	DecType        string
	Name           string
	Params         []Declaration
	SubroutineBody *SubroutineBody
}

func (sd *SubroutineDec) declarationNode() {}

type Param struct {
	Kind    string
	DecType string
	Name    string
}

func (p *Param) declarationNode() {}

type SubroutineBody struct {
	VarDecs []Declaration
	// Statements []Statement
}

func (sb *SubroutineBody) declarationNode() {}

// STATEMENTS HERE

type Statement interface {
	Node
	stmtNode()
}

type LetStatement struct {
	Token      token.Token
	Ident      *Identifier
	Expression Expression
}

func (ls *LetStatement) stmtNode() {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Token.Literal + " ")
	out.WriteString(ls.Ident.String())
	out.WriteString(" = ")
	out.WriteString(ls.Expression.String())
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token      token.Token
	Expression Expression
}

func (rs *ReturnStatement) stmtNode() {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Literal + " ")
	out.WriteString(rs.Expression.String())
	out.WriteString(";")

	return out.String()
}

type DoStatement struct {
	Token      token.Token
	Expression Expression
}

func (ds *DoStatement) stmtNode() {}
func (ds *DoStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ds.Token.Literal + " ")
	out.WriteString(ds.Expression.String())
	out.WriteString(";")

	return out.String()
}

type IfStatement struct {
	Token      token.Token
	Expression Expression
	IfStmts    []Statement
	Else       []Statement
}

func (is *IfStatement) stmtNode() {}
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.Token.Literal + " ")
	out.WriteString("(" + is.Expression.String() + ") ")
	out.WriteString("{")
	for _, stmt := range is.IfStmts {
		out.WriteString(stmt.String())
	}
	if len(is.Else) > 0 {
		out.WriteString("} else {")
		for _, stmt := range is.Else {
			out.WriteString(stmt.String())
		}
	}
	out.WriteString("}")

	return out.String()
}

type WhileStatement struct {
	Token      token.Token
	Expression Expression
	Stmts      []Statement
}

func (ws *WhileStatement) stmtNode() {}
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.Token.Literal + " ")
	out.WriteString("(" + ws.Expression.String() + ") ")
	out.WriteString("{")
	for _, stmt := range ws.Stmts {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")

	return out.String()
}

// STATEMENTS HERE

type Node interface {
	String() string
}

type Expression interface {
	Node
	expNode()
}

type Prefix struct {
	Token      token.Token
	Operator   string
	Expression Expression
}

func (p *Prefix) expNode() {}
func (p *Prefix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Expression.String())
	out.WriteString(")")

	return out.String()
}

type Infix struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *Infix) expNode() {}
func (i *Infix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

type IntegerConstant struct {
	Token token.Token
	Value int
}

func (ic *IntegerConstant) expNode()       {}
func (ic *IntegerConstant) String() string { return ic.Token.Literal }

type KeywordConstant struct {
	Token token.Token
	Value string
}

func (kc *KeywordConstant) expNode()       {}
func (kc *KeywordConstant) String() string { return kc.Token.Literal }

type StringConstant struct {
	Token token.Token
	Value string
}

func (sc *StringConstant) expNode()       {}
func (sc *StringConstant) String() string { return `"` + sc.Token.Literal + `"` }

type Identifier struct {
	Token   token.Token
	Value   string
	Indexer Expression
}

func (i *Identifier) expNode() {}
func (i *Identifier) String() string {
	if i.Indexer == nil {
		return i.Token.Literal
	}

	var out bytes.Buffer

	out.WriteString(i.Value)
	if i.Indexer != nil {
		out.WriteString("[")
		out.WriteString(i.Indexer.String())
		out.WriteString("]")
	}

	return out.String()
}

type SubroutineCall struct {
	Ident      *Identifier
	Subroutine *Identifier
	ExpList    []Expression
}

func (sc *SubroutineCall) expNode() {}
func (sc *SubroutineCall) String() string {
	var out bytes.Buffer

	if sc.Ident != nil {
		out.WriteString(sc.Ident.String() + ".")
	}
	out.WriteString(sc.Subroutine.String())
	out.WriteString("(")
	for i, exp := range sc.ExpList {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(exp.String())
	}
	out.WriteString(")")

	return out.String()
}
