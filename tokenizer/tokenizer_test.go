package tokenizer

import (
	"testing"

	"github.com/tivt2/jack-compiler/token"
)

func TestAdvance(t *testing.T) {
	input := `
	class Test {
		field int x;
		static boolean y;

		constructor Test new(char s, int ax) {
			let x = ax;
			do Output.println(s);
			return this;
		}
		+-~*/[]void method function &|<>"true false null if else while
	}
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CLASS, "class"},
		{token.IDENT, "Test"},
		{token.LBRACE, "{"},
		{token.FIELD, "field"},
		{token.INT, "int"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.STATIC, "static"},
		{token.BOOLEAN, "boolean"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},

		{token.CONSTRUCTOR, "constructor"},
		{token.IDENT, "Test"},
		{token.IDENT, "new"},
		{token.LPAREN, "("},
		{token.CHAR, "char"},
		{token.IDENT, "s"},
		{token.COMMA, ","},
		{token.INT, "int"},
		{token.IDENT, "ax"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.EQ, "="},
		{token.IDENT, "ax"},
		{token.SEMICOLON, ";"},
		{token.DO, "do"},
		{token.IDENT, "Output"},
		{token.DOT, "."},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "s"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.THIS, "this"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.NOT, "~"},
		{token.ASTERISK, "*"},
		{token.FSLASH, "/"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.VOID, "void"},
		{token.METHOD, "method"},
		{token.FUNCTION, "function"},
		{token.AMP, "&"},
		{token.BAR, "|"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.QUOT, `"`},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.WHILE, "while"},

		{token.RBRACE, "}"},
	}

	tkzr := New(input)

	for i, test := range tests {
		tk := tkzr.Advance()

		if test.expectedType != tk.Type {
			t.Fatalf("TokenType failed. test index %d, expected: %q, received: %q", i, test.expectedType, tk.Type)
		}

		if test.expectedLiteral != tk.Literal {
			t.Fatalf("Literal failed. test index %d, expected: %q, received: %q", i, test.expectedLiteral, tk.Literal)
		}
	}
}
