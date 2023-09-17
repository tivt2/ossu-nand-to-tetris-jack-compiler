package parser

import (
	"testing"

	"github.com/tivt2/jack-compiler/tokenizer"
)

// func TestParseExpression(t *testing.T) {
// 	infixTests := []struct {
// 		input  string
// 		expect string
// 	}{
// 		{"-5", "(-5)"},

// 		{"5 + 5", "(5 + 5)"},
// 		{"5 - 5", "(5 - 5)"},
// 		{"5 * 5", "(5 * 5)"},
// 		{"5 / 5", "(5 / 5)"},
// 		{"5 < 5", "(5 < 5)"},
// 		{"5 > 5", "(5 > 5)"},
// 		{"5 = 5", "(5 = 5)"},
// 		{"5 & 5", "(5 & 5)"},
// 		{"5 | 5", "(5 | 5)"},

// 		{"5 + 5 - 5", "((5 + 5) - 5)"},
// 		{"-5 + 5 - 5", "(((-5) + 5) - 5)"},
// 		{"-(5 + 4) - 3", "((-(5 + 4)) - 3)"},
// 		{"~(5 > 4) = 3", "((~(5 > 4)) = 3)"},

// 		{"~true = (false = false)", "((~true) = (false = false))"},

// 		{`"test" + "something"`, `("test" + "something")`},

// 		{`something[1]`, `something[1]`},
// 		{`-5 * (2 + something[1])`, `((-5) * (2 + something[1]))`},
// 		{`2 * something[-5 + 3]`, `(2 * something[((-5) + 3)])`},

// 		{`2 * myVar`, `(2 * myVar)`},

// 		{`Something.call()`, `Something.call()`},
// 		{`Something.call(1 + 1 = 5, ~false = false)`, `Something.call(((1 + 1) = 5), ((~false) = false))`},

// 		{`Something[1].call()`, `Something[1].call()`},
// 		{`Something[-2 * 5 + 3].call(true = true, 5 *2 > 1)`, `Something[(((-2) * 5) + 3)].call((true = true), ((5 * 2) > 1))`},

// 		{`call()`, `call()`},
// 		{`call(5 + 3, true)`, `call((5 + 3), true)`},
// 	}

// 	for _, test := range infixTests {
// 		tkzr := tokenizer.New(test.input)
// 		p := New(tkzr)

// 		expTree := p.parseExpression()

// 		if expTree.String() != test.expect {
// 			t.Fatalf("expTree.String() wrong print, expected: %s, received: %s", test.expect, expTree.String())
// 		}
// 	}

// }

func TestParseStatement(t *testing.T) {
	stmtTests := []struct {
		input  string
		expect string
	}{
		// {"let x = 1;", "let x = 1;"},
		// {"let x[0] = 1;", "let x[0] = 1;"},
		{"let x[a[1]] = y[b[0]];", "let x[a[1]] = y[b[0]];"},
		// 		{`return call();`, `return call();`},
		// 		{`return;`, `return;`},
		// 		{`do Something.print("abc");`, `do Something.print("abc");`},
		// 		{
		// 			`if (true) {
		// 						let x = 2 * 2 * 2;
		// 						return x + 1;
		// 						} else {
		// 						let e[1] = call();
		// 						return 1;
		// 						}`,
		// 			`if (true) {
		// let x = ((2 * 2) * 2);
		// return (x + 1);
		// } else {
		// let e[1] = call();
		// return 1;
		// }`,
		// 		},
		// 		{
		// 			`while (~(5 < 3)) {
		// 						let x = 2 * 2 * 2;
		// 						return x + 1;
		// 						}`,
		// 			`while ((~(5 < 3))) {
		// let x = ((2 * 2) * 2);
		// return (x + 1);
		// }`,
		// 		},
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

// func TestParseSubroutineDec(t *testing.T) {
// 	subDecTestst := []struct {
// 		input  string
// 		expect string
// 	}{
// 		{
// 			`constructor void new(int x, int y) {
// 				var boolean t, j;
// 				return this;
// 			}`,
// 			`constructor void new(int x, int y) {
// var boolean t;
// var boolean j;
// return this;
// }`,
// 		},
// 		{
// 			`function int print(char s) {
// 				var char filtred;
// 				let filtred = Filter.filter(s);
// 				return String.length(filtred);
// 			}`,
// 			`function int print(char s) {
// var char filtred;
// let filtred = Filter.filter(s);
// return String.length(filtred);
// }`,
// 		},
// 		{
// 			`method Point createPoint(int x, int y) {
// 				var Point p;
// 				let p = Point.new(x, y);
// 				return Point;
// 			}`,
// 			`method Point createPoint(int x, int y) {
// var Point p;
// let p = Point.new(x, y);
// return Point;
// }`,
// 		},
// 	}

// 	for _, test := range subDecTestst {
// 		tkzr := tokenizer.New(test.input)
// 		p := New(tkzr)

// 		subDec := p.parseSubroutineDec()

// 		if subDec.String() != test.expect {
// 			t.Fatalf("subDec.String() wrong print, expected: %s, received: %s", test.expect, subDec.String())
// 		}
// 	}
// }

// func TestParseClass(t *testing.T) {
// 	classTest := struct {
// 		input  string
// 		expect string
// 	}{
// 		`class Something {
// 				field int x, y;
// 				static Array p;
// 				static int count;

// 				constructor Something new(int ax, int ay) {
// 					let x = ax;
// 					let y = ay;
// 					let p = Array.new(6);
// 					let count = 0;
// 					return this;
// 				}

// 				method Point addPoint() {
// 					var Point n_p;
// 					let n_p = Point.new(x, y);
// 					let p[count] = n_p;
// 					let count = count + 1;
// 					return n_p;
// 				}

// 				function void Print(char s) {
// 					do Os.outputln(s);
// 					return;
// 				}
// 			}`,
// 		`class Something {
// field int x;
// field int y;
// static Array p;
// static int count;
// constructor Something new(int ax, int ay) {
// let x = ax;
// let y = ay;
// let p = Array.new(6);
// let count = 0;
// return this;
// }
// method Point addPoint() {
// var Point n_p;
// let n_p = Point.new(x, y);
// let p[count] = n_p;
// let count = (count + 1);
// return n_p;
// }
// function void Print(char s) {
// do Os.outputln(s);
// return;
// }
// }`,
// 	}

// 	tkzr := tokenizer.New(classTest.input)
// 	p := New(tkzr)

// 	class := p.ParseClass()

// 	if class.String() != classTest.expect {
// 		t.Fatalf("class.String() print wrong, expected: %s, received: %s", classTest.expect, class.String())
// 	}
// }
