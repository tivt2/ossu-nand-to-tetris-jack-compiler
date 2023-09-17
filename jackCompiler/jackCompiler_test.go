package jackCompiler

import (
	"testing"

	"github.com/tivt2/jack-compiler/parseTree"
	"github.com/tivt2/jack-compiler/symbolTable"
	"github.com/tivt2/jack-compiler/token"
	"github.com/tivt2/jack-compiler/vmWriter"
)

func TestCompileExpression(t *testing.T) {
	tests := []struct {
		input    parseTree.Expression
		expected string
	}{
		{
			&parseTree.Prefix{
				Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
				Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
			},
			"push constant 5\nneg\n",
		},
		{
			&parseTree.Prefix{
				Operator:   token.Token{Type: token.NOT, Literal: token.NOT},
				Expression: &parseTree.KeywordConstant{Token: token.Token{Type: token.TRUE, Literal: token.TRUE}, Value: token.TRUE},
			},
			"push constant 1\nneg\nnot\n",
		},
		{
			&parseTree.Infix{
				Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
				Left: &parseTree.Prefix{
					Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
					Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
				},
				Right: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
			},
			"push constant 5\nneg\npush constant 5\nadd\n",
		},
		{
			&parseTree.Infix{
				Operator: token.Token{Type: token.ASTERISK, Literal: token.ASTERISK},
				Left: &parseTree.Prefix{
					Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
					Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
				},
				Right: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "7"}, Value: 7},
			},
			"push constant 5\nneg\npush constant 7\ncall Math.multiply 2\n",
		},
		{
			&parseTree.SubroutineCall{
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "SomeClass"},
					Value: "SomeClass",
				},
				Subroutine: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "someFunction"},
					Value: "someFunction",
				},
				ExpList: []parseTree.Expression{
					&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
				},
			},
			"push constant 1\ncall SomeClass.someFunction 1\n",
		},
		{
			&parseTree.Infix{
				Operator: token.Token{Type: token.ASTERISK, Literal: token.ASTERISK},
				Left: &parseTree.SubroutineCall{
					Ident: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "SomeClass"},
						Value: "SomeClass",
					},
					Subroutine: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "someFunction"},
						Value: "someFunction",
					},
					ExpList: []parseTree.Expression{
						&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					},
				},
				Right: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "localVar"}, Value: "localVar"},
			},
			"push constant 1\ncall SomeClass.someFunction 1\npush local 0\ncall Math.multiply 2\n",
		},
		{
			&parseTree.Infix{
				Operator: token.Token{Type: token.ASTERISK, Literal: token.ASTERISK},
				Left: &parseTree.SubroutineCall{
					Ident: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "instance"},
						Value: "instance",
					},
					Subroutine: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "someFunction"},
						Value: "someFunction",
					},
					ExpList: []parseTree.Expression{
						&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					},
				},
				Right: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "localVar"}, Value: "localVar"},
			},
			"push this 0\npush constant 1\ncall SomeClass.someFunction 2\npush local 0\ncall Math.multiply 2\n",
		},
	}

	for _, test := range tests {
		jc := &JackCompiler{w: vmWriter.New("testing.jack"), s: symbolTable.New()}

		jc.s.Define("this", "SomeClass", "argument")
		jc.s.Define("localVar", "int", "local")
		jc.s.Define("instance", "SomeClass", "field")

		jc.CompileExpression(test.input)

		if jc.w.Out.String() != test.expected {
			t.Fatalf("CompileExpression(), expected: %s, received: %s", test.expected, jc.w.Out.String())
		}
	}
}

func TestStatement(t *testing.T) {
	tests := []struct {
		input    parseTree.Statement
		expected string
	}{
		{
			&parseTree.LetStatement{
				Token: token.Token{Type: token.LET, Literal: token.LET},
				Ident: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
					Left:     &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
					Right:    &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
				},
			},
			"push constant 5\npush constant 10\nadd\npop this 0\n",
		},
		{
			&parseTree.LetStatement{
				Token: token.Token{Type: token.LET, Literal: token.LET},
				Ident: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "p"}, Value: "p"},
				Expression: &parseTree.SubroutineCall{
					Ident: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "Point"},
						Value: "Point",
					},
					Subroutine: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "new"},
						Value: "new",
					},
					ExpList: []parseTree.Expression{
						&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
						&parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
					},
				},
			},
			"push constant 1\npush this 0\ncall Point.new 2\npop this 1\n",
		},
		{
			&parseTree.LetStatement{
				Token: token.Token{Type: token.LET, Literal: token.LET},
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
					Indexer: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: 2,
					},
				},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
					Left:     &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
					Right:    &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
				},
			},
			"push this 0\npush constant 2\nadd\npush constant 5\npush constant 10\nadd\npop temp 0\npop pointer 1\npush temp 0\npop that 0\n",
		},
		{
			&parseTree.LetStatement{
				Token: token.Token{Type: token.LET, Literal: token.LET},
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
					Indexer: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "2"},
						Value: 2,
					},
				},
				Expression: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Value: "x",
					Indexer: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			"push this 0\npush constant 2\nadd\npush this 0\npush constant 5\nadd\npop pointer 1\npush that 0\npop temp 0\npop pointer 1\npush temp 0\npop that 0\n",
		},
		{
			&parseTree.ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: token.RETURN},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
					Left:     &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
					Right:    &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "10"}, Value: 10},
				},
			},
			"push constant 5\npush constant 10\nadd\nreturn\n",
		},
		{
			&parseTree.ReturnStatement{
				Token:      token.Token{Type: token.RETURN, Literal: token.RETURN},
				Expression: nil,
			},
			"push constant 0\nreturn\n",
		},
		{
			&parseTree.ReturnStatement{
				Token:      token.Token{Type: token.RETURN, Literal: token.RETURN},
				Expression: &parseTree.KeywordConstant{Token: token.Token{Type: token.THIS, Literal: token.THIS}, Value: token.THIS},
			},
			"push pointer 0\nreturn\n",
		},
		{
			&parseTree.DoStatement{
				Token: token.Token{Type: token.DO, Literal: token.DO},
				Expression: &parseTree.SubroutineCall{
					Ident: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "p"},
						Value: "p",
					},
					Subroutine: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "someFunction"},
						Value: "someFunction",
					},
					ExpList: []parseTree.Expression{
						&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					},
				},
			},
			"push this 1\npush constant 1\ncall Point.someFunction 2\npop temp 0\n",
		},
		{
			&parseTree.WhileStatement{
				Token: token.Token{Type: token.WHILE, Literal: token.WHILE},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.GT, Literal: token.GT},
					Left: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Right: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "0"},
						Value: 0,
					},
				},
				Stmts: []parseTree.Statement{
					&parseTree.LetStatement{
						Token: token.Token{Type: token.LET, Literal: token.LET},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Expression: &parseTree.Infix{
							Operator: token.Token{Type: token.MINUS, Literal: token.MINUS},
							Left: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Right: &parseTree.IntegerConstant{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
			"label WHILE0\npush this 0\npush constant 0\ngt\nnot\nif-goto BREAK0\npush this 0\npush constant 1\nsub\npop this 0\ngoto WHILE0\nlabel BREAK0\n",
		},
		{
			&parseTree.IfStatement{
				Token: token.Token{Type: token.IF, Literal: token.IF},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.LT, Literal: token.LT},
					Left: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Right: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "0"},
						Value: 0,
					},
				},
				IfStmts: []parseTree.Statement{
					&parseTree.ReturnStatement{
						Token:      token.Token{Type: token.RETURN, Literal: token.RETURN},
						Expression: nil,
					},
				},
				Else: []parseTree.Statement{},
			},
			"push this 0\npush constant 0\nlt\nnot\nif-goto ELSE0\npush constant 0\nreturn\nlabel ELSE0\n",
		},
		{
			&parseTree.IfStatement{
				Token: token.Token{Type: token.IF, Literal: token.IF},
				Expression: &parseTree.Infix{
					Operator: token.Token{Type: token.LT, Literal: token.LT},
					Left: &parseTree.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Right: &parseTree.IntegerConstant{
						Token: token.Token{Type: token.INT, Literal: "0"},
						Value: 0,
					},
				},
				IfStmts: []parseTree.Statement{
					&parseTree.ReturnStatement{
						Token:      token.Token{Type: token.RETURN, Literal: token.RETURN},
						Expression: nil,
					},
				},
				Else: []parseTree.Statement{
					&parseTree.LetStatement{
						Token: token.Token{Type: token.LET, Literal: token.LET},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Expression: &parseTree.Infix{
							Operator: token.Token{Type: token.MINUS, Literal: token.MINUS},
							Left: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Right: &parseTree.IntegerConstant{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
			"push this 0\npush constant 0\nlt\nnot\nif-goto ELSE0\npush constant 0\nreturn\ngoto IF0\nlabel ELSE0\npush this 0\npush constant 1\nsub\npop this 0\nlabel IF0\n",
		},
	}

	for _, test := range tests {
		jc := &JackCompiler{w: vmWriter.New("testing.jack"), s: symbolTable.New()}

		jc.s.Define("this", "Something", "argument")
		jc.s.Define("x", "int", "field")
		jc.s.Define("p", "Point", "field")

		jc.CompileStatement(test.input)

		if jc.w.Out.String() != test.expected {
			t.Fatalf("CompileStatement()\n\nexpected:\n%s\n\nreceived:\n%s", test.expected, jc.w.Out.String())
		}
	}
}

func TestCompileSubroutineDec(t *testing.T) {
	tests := []struct {
		input    *parseTree.SubroutineDec
		expected string
	}{
		{
			&parseTree.SubroutineDec{
				Kind:    token.Token{Type: token.CONSTRUCTOR, Literal: token.CONSTRUCTOR},
				DecType: token.Token{Type: token.IDENT, Literal: "Point"},
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "new"},
					Value: "new",
				},
				Params: []*parseTree.Param{
					{
						DecType: token.Token{Type: token.INT, Literal: token.INT},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "ax"},
							Value: "ax",
						},
					},
					{
						DecType: token.Token{Type: token.INT, Literal: token.INT},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "ay"},
							Value: "ay",
						},
					},
				},
				SubroutineBody: &parseTree.SubroutineBody{
					VarDecs: []*parseTree.VarDec{},
					Statements: []parseTree.Statement{
						&parseTree.LetStatement{
							Token: token.Token{Type: token.LET, Literal: token.LET},
							Ident: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Expression: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "ax"},
								Value: "ax",
							},
						},
						&parseTree.LetStatement{
							Token: token.Token{Type: token.LET, Literal: token.LET},
							Ident: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
							Expression: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "ay"},
								Value: "ay",
							},
						},
						&parseTree.ReturnStatement{
							Token: token.Token{Type: token.RETURN, Literal: token.RETURN},
							Expression: &parseTree.KeywordConstant{
								Token: token.Token{Type: token.THIS, Literal: token.THIS},
								Value: token.THIS,
							},
						},
					},
				},
			},
			"function Point.new 0\npush constant 2\ncall Memory.alloc 1\npop pointer 0\npush argument 0\npop this 0\npush argument 1\npop this 1\npush pointer 0\nreturn\n",
		},
		{
			&parseTree.SubroutineDec{
				Kind:    token.Token{Type: token.METHOD, Literal: token.METHOD},
				DecType: token.Token{Type: token.IDENT, Literal: "Point"},
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "sum"},
					Value: "sum",
				},
				Params: []*parseTree.Param{},
				SubroutineBody: &parseTree.SubroutineBody{
					VarDecs: []*parseTree.VarDec{},
					Statements: []parseTree.Statement{
						&parseTree.ReturnStatement{
							Token: token.Token{Type: token.RETURN, Literal: token.RETURN},
							Expression: &parseTree.Infix{
								Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
								Left: &parseTree.Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Value: "x",
								},
								Right: &parseTree.Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "y"},
									Value: "y",
								},
							},
						},
					},
				},
			},
			"function Point.sum 0\npush argument 0\npop pointer 0\npush this 0\npush this 1\nadd\nreturn\n",
		},
		{
			&parseTree.SubroutineDec{
				Kind:    token.Token{Type: token.FUNCTION, Literal: token.FUNCTION},
				DecType: token.Token{Type: token.IDENT, Literal: "Point"},
				Ident: &parseTree.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "sum"},
					Value: "sum",
				},
				Params: []*parseTree.Param{
					{
						DecType: token.Token{Type: token.INT, Literal: token.INT},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "ax"},
							Value: "ax",
						},
					},
					{
						DecType: token.Token{Type: token.INT, Literal: token.INT},
						Ident: &parseTree.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "ay"},
							Value: "ay",
						},
					},
				},
				SubroutineBody: &parseTree.SubroutineBody{
					VarDecs: []*parseTree.VarDec{
						{
							Kind:    token.Token{Type: token.VAR, Literal: token.VAR},
							DecType: token.Token{Type: token.INT, Literal: token.INT},
							Ident: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "out"},
								Value: "out",
							},
						},
					},
					Statements: []parseTree.Statement{
						&parseTree.LetStatement{
							Token: token.Token{Type: token.LET, Literal: token.LET},
							Ident: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "out"},
								Value: "out",
							},
							Expression: &parseTree.Infix{
								Operator: token.Token{Type: token.PLUS, Literal: token.PLUS},
								Left: &parseTree.Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "ax"},
									Value: "ax",
								},
								Right: &parseTree.Identifier{
									Token: token.Token{Type: token.IDENT, Literal: "ay"},
									Value: "ay",
								},
							},
						},
						&parseTree.ReturnStatement{
							Token: token.Token{Type: token.RETURN, Literal: token.RETURN},
							Expression: &parseTree.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "out"},
								Value: "out",
							},
						},
					},
				},
			},
			"function Point.sum 1\npush argument 0\npush argument 1\nadd\npop local 0\npush local 0\nreturn\n",
		},
	}

	for _, test := range tests {
		jc := &JackCompiler{
			w: vmWriter.New("testing.jack"),
			s: symbolTable.New(),
			c: &parseTree.Class{Ident: &parseTree.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "Point"},
				Value: "Point",
			}},
		}

		jc.s.Define("x", "int", "field")
		jc.s.Define("y", "int", "field")

		jc.CompileSubroutineDec(test.input)

		if jc.w.Out.String() != test.expected {
			t.Fatalf("CompileSubroutineDec()\n\nexpected:\n%s\n\nreceived:\n%s", test.expected, jc.w.Out.String())
		}
	}
}
