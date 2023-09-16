package parser

import (
	"testing"

	"github.com/tivt2/jack-compiler/tokenizer"
)

func TestParseExpression(t *testing.T) {
	infixTests := []struct {
		input  string
		expect string
	}{
		{"-5", "(-5)"},

		{"5 + 5", "(5 + 5)"},
		{"5 - 5", "(5 - 5)"},
		{"5 * 5", "(5 * 5)"},
		{"5 / 5", "(5 / 5)"},
		{"5 < 5", "(5 < 5)"},
		{"5 > 5", "(5 > 5)"},
		{"5 = 5", "(5 = 5)"},
		{"5 & 5", "(5 & 5)"},
		{"5 | 5", "(5 | 5)"},

		{"5 + 5 - 5", "((5 + 5) - 5)"},
		{"-5 + 5 - 5", "(((-5) + 5) - 5)"},
		{"-(5 + 4) - 3", "((-(5 + 4)) - 3)"},
		{"~(5 > 4) = 3", "((~(5 > 4)) = 3)"},

		{"~true = (false = false)", "((~true) = (false = false))"},

		{`"test" + "something"`, `("test" + "something")`},

		{`something[1]`, `something[1]`},
		{`-5 * (2 + something[1])`, `((-5) * (2 + something[1]))`},
		{`2 * something[-5 + 3]`, `(2 * something[((-5) + 3)])`},

		{`2 * myVar`, `(2 * myVar)`},

		{`Something.call()`, `Something.call()`},
		{`Something.call(1 + 1 = 5, ~false = false)`, `Something.call(((1 + 1) = 5), ((~false) = false))`},

		{`Something[1].call()`, `Something[1].call()`},
		{`Something[-2 * 5 + 3].call(true = true, 5 *2 > 1)`, `Something[(((-2) * 5) + 3)].call((true = true), ((5 * 2) > 1))`},

		{`call()`, `call()`},
		{`call(5 + 3, true)`, `call((5 + 3), true)`},
	}

	for _, test := range infixTests {
		tkzr := tokenizer.New(test.input)
		p := New(tkzr)

		expTree := p.parseExpression()

		if expTree.String() != test.expect {
			t.Fatalf("expTree.String() wrong print, expected: %s, received: %s", test.expect, expTree.String())
		}
	}

}

func TestParseStatement(t *testing.T) {
	stmtTests := []struct {
		input  string
		expect string
	}{
		{"let x = 1;", "let x = 1;"},
		{"let x[0] = 1;", "let x[0] = 1;"},
		{`return call();`, `return call();`},
		{`do Something.print("abc");`, `do Something.print("abc");`},
		{
			"if (true) {let x = 2 * 2 * 2;return x + 1;} else {let e[1] = call();return 1;}",
			"if (true) {let x = ((2 * 2) * 2);return (x + 1);} else {let e[1] = call();return 1;}",
		},
		{
			"while (~(5 < 3)) {let x = 2 * 2 * 2;return x + 1;}",
			"while ((~(5 < 3))) {let x = ((2 * 2) * 2);return (x + 1);}",
		},
	}

	for _, test := range stmtTests {
		tkzr := tokenizer.New(test.input)
		p := New(tkzr)

		stmt := p.parseStatement()

		if stmt.String() != test.expect {
			t.Fatalf("stmt.String() wrong print, expected: %s, received: %s", test.expect, stmt.String())
		}
	}
}
