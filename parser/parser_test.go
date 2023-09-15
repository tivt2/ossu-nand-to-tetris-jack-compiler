package parser

import (
	"testing"

	"github.com/tivt2/jack-compiler/parse_tree"
	"github.com/tivt2/jack-compiler/tokenizer"
)

func TestClassVarDecs(t *testing.T) {
	input := `
	class Somethin {
	field int x;
	static char y, p;
	field boolean foo;
	}
	`
	tkzr := tokenizer.New(input)
	p := New(tkzr)

	class := p.ParseClass()
	if class == nil {
		t.Fatalf("ParseClass() returned nil")
	}
	if len(class.ClassVarDecs) != 4 {
		t.Fatalf("Class.VarDec does not contain 4. length: %d", len(class.ClassVarDecs))
	}

	tests := []struct {
		expectedType       string
		expectedIdentifier string
	}{
		{"int", "x"},
		{"char", "y"},
		{"char", "p"},
		{"boolean", "foo"},
	}

	for i, test := range tests {
		vd := class.ClassVarDecs[i]
		if !testClassVarDec(t, vd, test.expectedIdentifier, test.expectedType) {
			return
		}
	}
}

func testClassVarDec(t *testing.T, vd parse_tree.Declaration, name string, decType string) bool {
	if vd.TokenLiteral() != "field" && vd.TokenLiteral() != "static" {
		t.Errorf("invalid vd.TokenLiteral, expected: field or static, received: %q", vd.TokenLiteral())
		return false
	}

	cvd, ok := vd.(*parse_tree.ClassVarDec)
	if !ok {
		t.Errorf("vd is not ClassVarDec, received: %T", vd)
		return false
	}

	if cvd.DecType.Literal != decType {
		t.Errorf("cvd.DecType.Literal not %s, received: %s", decType, cvd.DecType.Literal)
		return false
	}

	if cvd.Name.Value != name {
		t.Errorf("cvd.Name.Value not %s, received: %s", name, cvd.Name.Value)
		return false
	}

	if cvd.Name.TokenLiteral() != name {
		t.Errorf("cvd.Name.TokenLiteral not %s, received: %s", name, cvd.Name)
		return false
	}
	return true
}

func TestSubroutineDec(t *testing.T) {
	input := `
	class Something {
		constructor Something new(int ax, int ay) {
			var int o, p;
		}
	}
	`
	tkzr := tokenizer.New(input)
	p := New(tkzr)

	class := p.ParseClass()
	if class == nil {
		t.Fatalf("ParseClass() returned nil")
	}
	if len(class.SubroutineDecs[0].SubroutineBody.VarDecs) != 2 {
		t.Fatalf("Class.VarDec does not contain 2. length: %d", len(class.ClassVarDecs))
	}

	tests := []struct {
		expectedDecType    string
		expectedIdentifier string
	}{
		{"int", "o"},
		{"int", "p"},
	}

	sd := class.SubroutineDecs[0]
	for i, test := range tests {
		vd := sd.SubroutineBody.VarDecs[i]
		if !testVarDec(t, vd, test.expectedDecType, test.expectedIdentifier) {
			return
		}
	}
}

func testVarDec(t *testing.T, vd parse_tree.Declaration, decType string, ident string) bool {
	if vd.TokenLiteral() != "var" {
		t.Errorf("invalid sds.TokenLiteral, expected: field or static, received: %q", vd.TokenLiteral())
		return false
	}

	var_dec, ok := vd.(*parse_tree.VarDec)
	if !ok {
		t.Errorf("sds is not VarDec, received: %T", vd)
		return false
	}

	if var_dec.Name.Value != ident {
		t.Errorf("var_dec.Name.Value not %s, received: %s", ident, var_dec.Name.Value)
		return false
	}

	if var_dec.Name.TokenLiteral() != ident {
		t.Errorf("var_dec.Name.TokenLiteral not %s, received: %s", ident, var_dec.Name)
		return false
	}

	if var_dec.DecType.Literal != decType {
		t.Errorf("var_dec.DecType.Literal not %s, received: %s", decType, var_dec.DecType.Literal)
		return false
	}
	return true
}
