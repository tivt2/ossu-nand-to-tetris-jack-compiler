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
				Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
			},
			"push constant 5\nnot\n",
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
				Right: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
			},
			"push constant 5\nneg\npush constant 5\ncall Math.multiply 2\npop temp 0\n",
		},
		{
			&parseTree.Infix{
				Operator: token.Token{Type: token.ASTERISK, Literal: token.ASTERISK},
				Left: &parseTree.Prefix{
					Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
					Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
				},
				Right: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "this"}, Value: "this"},
			},
			"push constant 5\nneg\npush argument 0\ncall Math.multiply 2\npop temp 0\n",
		},
		{
			&parseTree.SubroutineCall{
				Ident:      &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "Something"}, Value: "Something"},
				Subroutine: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "print"}, Value: "print"},
				ExpList: []parseTree.Expression{
					&parseTree.Infix{
						Operator: token.Token{Type: token.ASTERISK, Literal: token.ASTERISK},
						Left: &parseTree.Prefix{
							Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
							Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
						},
						Right: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "this"}, Value: "this"},
					},
					&parseTree.Prefix{
						Operator:   token.Token{Type: token.MINUS, Literal: token.MINUS},
						Expression: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
					},
				},
			},
			"push constant 5\nneg\npush argument 0\ncall Math.multiply 2\npop temp 0\npush constant 5\nneg\ncall Something.print 2\npop temp 0\n",
		},
	}

	for _, test := range tests {
		jc := &JackCompiler{w: vmWriter.New("testing.jack"), s: symbolTable.New()}

		jc.s.Define("this", "Something", "argument")

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
			"push constant 5\npush constant 10\nadd\npop local 0\n",
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
		// {
		// 	&parseTree.DoStatement{
		// 		Token: token.Token{Type: token.DO, Literal: token.DO},
		// 		Expression: &parseTree.SubroutineCall{
		// 			Ident: &parseTree.Identifier{
		// 				Token: token.Token{
		// 					Type: token.IDENT, Literal: "myVar"}, Value: "myVar",
		// 				Indexer: &parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
		// 			},
		// 			Subroutine: &parseTree.Identifier{Token: token.Token{Type: token.IDENT, Literal: "print"}, Value: "print"},
		// 			ExpList: []parseTree.Expression{
		// 				&parseTree.IntegerConstant{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
		// 			},
		// 		},
		// 	},
		// 	"push constant 5\ncall myVar[2].print\n",
		// },
	}

	for _, test := range tests {
		jc := &JackCompiler{w: vmWriter.New("testing.jack"), s: symbolTable.New()}

		jc.s.Define("this", "Something", "argument")
		jc.s.Define("x", "int", "local")

		jc.CompileStatement(test.input)

		if jc.w.Out.String() != test.expected {
			t.Fatalf("CompileStatement(), expected: %s, received: %s", test.expected, jc.w.Out.String())
		}
	}
}
