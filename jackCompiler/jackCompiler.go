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
		jc.CompileSubroutineDec(subDec)

	}

	fmt.Println(jc.s)
}

func (jc *JackCompiler) CompileSubroutineDec(sd *parseTree.SubroutineDec) {
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

	jc.w.WriteFunction(fmt.Sprintf("%s.%s", jc.c.Ident.Value, sd.Ident.Value), jc.s.VarCount("local"))
	switch sd.Kind.Type {
	case token.CONSTRUCTOR:
		jc.w.WritePush("constant", jc.s.VarCount("this"))
		jc.w.WriteCall("Memory.alloc", 1)
		jc.w.WritePop("pointer", 0)
	case token.METHOD:
		jc.w.WritePush("argument", 0)
		jc.w.WritePop("pointer", 0)
	}

	jc.CompileStatements(sd.SubroutineBody.Statements)
}

func (jc *JackCompiler) CompileStatements(stmts []parseTree.Statement) {
	for _, stmt := range stmts {
		jc.CompileStatement(stmt)
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
		jc.w.WritePop("temp", 0)
	case *parseTree.WhileStatement:
		jc.w.WriteLabel(fmt.Sprintf("WHILE%d", jc.whileCounter))
		jc.CompileExpression(stmt.Expression)
		jc.w.WriteArithmetic(token.NOT)
		jc.w.WriteIf(fmt.Sprintf("BREAK%d", jc.whileCounter))
		jc.CompileStatements(stmt.Stmts)
		jc.w.WriteGoto(fmt.Sprintf("WHILE%d", jc.whileCounter))
		jc.w.WriteLabel(fmt.Sprintf("BREAK%d", jc.whileCounter))
		jc.whileCounter++
	case *parseTree.IfStatement:
		elseLen := len(stmt.Else)
		jc.CompileExpression(stmt.Expression)
		jc.w.WriteArithmetic(token.NOT)
		jc.w.WriteIf(fmt.Sprintf("ELSE%d", jc.ifCounter))
		jc.CompileStatements(stmt.IfStmts)
		if elseLen > 0 {
			jc.w.WriteGoto(fmt.Sprintf("IF%d", jc.ifCounter))
		}
		jc.w.WriteLabel(fmt.Sprintf("ELSE%d", jc.ifCounter))
		if elseLen > 0 {
			jc.CompileStatements(stmt.Else)
			jc.w.WriteLabel(fmt.Sprintf("IF%d", jc.ifCounter))
		}
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
			jc.w.WritePush("constant", 1)
			jc.w.WriteArithmetic("neg")
		case token.FALSE:
			jc.w.WritePush("constant", 0)
		case token.NULL:
			jc.w.WritePush("constant", 0)
		case token.THIS:
			jc.w.WritePush("pointer", 0)
		}
	case *parseTree.SubroutineCall:
		if exp.Ident != nil {
			if kind := jc.s.KindOf(exp.Ident.Value); kind != "" {
				jc.w.WritePush(kind, jc.s.IndexOf(exp.Ident.Value))
				for _, e := range exp.ExpList {
					jc.CompileExpression(e)
				}
				if exp.Subroutine.Indexer != nil {
					// !!!!!!!! TODO
				} else {
					jc.w.WriteCall(fmt.Sprintf("%s.%s", jc.s.TypeOf(exp.Ident.Value), exp.Subroutine.Value), len(exp.ExpList)+1)
				}
			} else {
				for _, e := range exp.ExpList {
					jc.CompileExpression(e)
				}
				if exp.Subroutine.Indexer != nil {
					// !!!!!!!! TODO
				} else {
					jc.w.WriteCall(fmt.Sprintf("%s.%s", exp.Ident.Value, exp.Subroutine.Value), len(exp.ExpList))
				}
			}
		} else {
			for _, e := range exp.ExpList {
				jc.CompileExpression(e)
			}
			jc.w.WriteCall(fmt.Sprintf("%s.%s", jc.c.Ident.Value, exp.Subroutine.Value), len(exp.ExpList))
		}
	}
}
