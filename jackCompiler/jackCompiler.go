package jackCompiler

import (
	"fmt"

	"github.com/tivt2/jack-compiler/parseTree"
	"github.com/tivt2/jack-compiler/symbolTable"
	"github.com/tivt2/jack-compiler/syntaxAnalyzer"
	"github.com/tivt2/jack-compiler/token"
	"github.com/tivt2/jack-compiler/vmWriter"
)

type JackCompiler struct {
	w *vmWriter.VMWriter
	s *symbolTable.SymbolTable
	c *parseTree.Class

	ifCounter    int
	whileCounter int
}

func New(filePath string) *JackCompiler {
	w := vmWriter.New(filePath)
	s := symbolTable.New()
	c := syntaxAnalyzer.ParseTree(filePath)

	for _, dec := range c.ClassVarDecs {
		s.Define(dec.Ident.Value, dec.DecType.Literal, dec.Kind.Literal)
	}

	return &JackCompiler{
		w: w,
		s: s,
		c: c,
	}
}

func (jc *JackCompiler) Compile() {
	jc.w.WriteComment(jc.c.String())

	for _, dec := range jc.c.ClassVarDecs {
		jc.s.Define(dec.Ident.Value, dec.DecType.Literal, dec.Kind.Literal)
	}

	for _, subDec := range jc.c.SubroutineDecs {
		jc.PopulateSubroutine(subDec)

	}

	fmt.Println(jc.s)
}

func (jc *JackCompiler) PopulateSubroutine(sd *parseTree.SubroutineDec) {
	jc.s.Reset()
	if sd.Kind.Type == token.METHOD {
		jc.s.Define("this", jc.c.Ident.Value, "argument")
	}
	for _, param := range sd.Params {
		jc.s.Define(param.Ident.Token.Literal, param.DecType.Literal, "argument")
	}
	for _, varDec := range sd.SubroutineBody.VarDecs {
		jc.s.Define(varDec.Ident.Token.Literal, varDec.DecType.Literal, "local")
	}
}

func (jc *JackCompiler) CompileStatement(stmt parseTree.Statement) {
	switch stmt := stmt.(type) {
	case *parseTree.LetStatement:
		jc.CompileExpression(stmt.Expression)
		jc.w.WritePop(jc.s.KindOf(stmt.Ident.Value), jc.s.IndexOf(stmt.Ident.Value))
	case *parseTree.ReturnStatement:
		if stmt.Expression != nil {
			jc.CompileExpression(stmt.Expression)
		} else {
			jc.w.WritePush("constant", 0)
		}
		jc.w.WriteReturn()
	case *parseTree.DoStatement:
		jc.CompileExpression(stmt.Expression)
	case *parseTree.WhileStatement:
		jc.CompileExpression(stmt.Expression)

		jc.whileCounter++
	case *parseTree.IfStatement:

		jc.ifCounter++
	}
}

func (jc *JackCompiler) CompileExpression(exp parseTree.Expression) {
	switch exp := exp.(type) {
	case *parseTree.Prefix:
		jc.CompileExpression(exp.Expression)
		if exp.Operator.Type == token.MINUS {
			jc.w.WriteArithmetic("neg")
		} else {
			jc.w.WriteArithmetic(exp.Operator.Literal)
		}
	case *parseTree.Infix:
		jc.CompileExpression(exp.Left)
		jc.CompileExpression(exp.Right)
		switch exp.Operator.Type {
		case token.ASTERISK:
			jc.w.WriteCall("Math.multiply", 2)
		case token.FSLASH:
			jc.w.WriteCall("Math.divide", 2)
		default:
			jc.w.WriteArithmetic(exp.Operator.Literal)
		}
	case *parseTree.Identifier:
		if exp.Indexer != nil {
			// !!!!!!!!!! TODO
		} else {
			jc.w.WritePush(jc.s.KindOf(exp.Value), jc.s.IndexOf(exp.Value))
		}
	case *parseTree.IntegerConstant:
		jc.w.WritePush("constant", exp.Value)
	case *parseTree.StringConstant:
		// !!!!!!!!!!! TODO
	case *parseTree.KeywordConstant:
		switch exp.Token.Type {
		case token.TRUE:
			jc.w.WritePush("constant", -1)
		case token.FALSE:
			jc.w.WritePush("constant", 0)
		case token.THIS:
			jc.w.WritePush("pointer", 0)
		}
	case *parseTree.SubroutineCall:
		for _, e := range exp.ExpList {
			jc.CompileExpression(e)
		}
		if exp.Subroutine.Indexer != nil {
			// !!!!!!!! TODO
		} else {
			if exp.Ident != nil {
				jc.w.WriteCall(fmt.Sprintf("%s.%s", exp.Ident.Value, exp.Subroutine.Value), len(exp.ExpList))
			} else {
				jc.w.WriteCall(exp.Subroutine.Value, len(exp.ExpList))
			}
		}
	}
}
