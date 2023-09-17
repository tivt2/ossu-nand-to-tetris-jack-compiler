package parseTree

import (
	"bytes"

	"github.com/tivt2/jack-compiler/token"
)

type Node interface {
	String() string
}

type Class struct {
	Token          token.Token
	Ident          *Identifier
	ClassVarDecs   []*ClassVarDec
	SubroutineDecs []*SubroutineDec
}

func (c *Class) String() string {
	var out bytes.Buffer

	out.WriteString(c.Token.Literal + " ")
	out.WriteString(c.Ident.String() + " ")
	out.WriteString("{\n")
	for _, cvd := range c.ClassVarDecs {
		out.WriteString(cvd.String() + "\n")
	}
	for _, sd := range c.SubroutineDecs {
		out.WriteString(sd.String() + "\n")
	}
	out.WriteString("}")

	return out.String()
}

type ClassVarDec struct {
	Kind    token.Token
	DecType token.Token
	Ident   *Identifier
}

func (cvd *ClassVarDec) String() string {
	var out bytes.Buffer

	out.WriteString(cvd.Kind.Literal + " ")
	out.WriteString(cvd.DecType.Literal + " ")
	out.WriteString(cvd.Ident.String())
	out.WriteString(";")

	return out.String()
}

type SubroutineDec struct {
	Kind           token.Token
	DecType        token.Token
	Ident          *Identifier
	Params         []*Param
	SubroutineBody *SubroutineBody
}

func (sd *SubroutineDec) String() string {
	var out bytes.Buffer

	out.WriteString(sd.Kind.Literal + " ")
	out.WriteString(sd.DecType.Literal + " ")
	out.WriteString(sd.Ident.String())
	out.WriteString("(")
	for i, param := range sd.Params {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.String())
	}
	out.WriteString(") ")
	out.WriteString("{\n")
	out.WriteString(sd.SubroutineBody.String())
	out.WriteString("}")

	return out.String()
}

type Param struct {
	DecType token.Token
	Ident   *Identifier
}

func (p *Param) String() string { return p.DecType.Literal + " " + p.Ident.String() }

type SubroutineBody struct {
	VarDecs    []*VarDec
	Statements []Statement
}

func (sb *SubroutineBody) String() string {
	var out bytes.Buffer

	for _, vd := range sb.VarDecs {
		out.WriteString(vd.String() + "\n")
	}
	for _, stmt := range sb.Statements {
		out.WriteString(stmt.String() + "\n")
	}

	return out.String()
}

type VarDec struct {
	Kind    token.Token
	DecType token.Token
	Ident   *Identifier
}

func (vd *VarDec) String() string {
	var out bytes.Buffer

	out.WriteString(vd.Kind.Literal + " ")
	out.WriteString(vd.DecType.Literal + " ")
	out.WriteString(vd.Ident.String())
	out.WriteString(";")

	return out.String()
}

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

	out.WriteString(rs.Token.Literal)
	if rs.Expression != nil {
		out.WriteString(" " + rs.Expression.String())
	}
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
	out.WriteString("{\n")
	for _, stmt := range is.IfStmts {
		out.WriteString(stmt.String() + "\n")
	}
	if len(is.Else) > 0 {
		out.WriteString("} else {\n")
		for _, stmt := range is.Else {
			out.WriteString(stmt.String() + "\n")
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
	out.WriteString("{\n")
	for _, stmt := range ws.Stmts {
		out.WriteString(stmt.String() + "\n")
	}
	out.WriteString("}")

	return out.String()
}

// STATEMENTS HERE

type Expression interface {
	Node
	expNode()
}

type Prefix struct {
	Operator   token.Token
	Expression Expression
}

func (p *Prefix) expNode() {}

func (p *Prefix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator.Literal)
	out.WriteString(p.Expression.String())
	out.WriteString(")")

	return out.String()
}

type Infix struct {
	Operator token.Token
	Left     Expression
	Right    Expression
}

func (i *Infix) expNode() {}
func (i *Infix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator.Literal + " ")
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

	out.WriteString(i.Token.Literal)
	out.WriteString("[")
	out.WriteString(i.Indexer.String())
	out.WriteString("]")

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
