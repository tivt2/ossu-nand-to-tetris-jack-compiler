package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "illegal"
	EOF     = "eof"

	IDENT   = "ident"
	INT     = "int"
	CHAR    = "char"
	BOOLEAN = "boolean"

	ASSIGN    = "="
	LBRACKET  = "["
	RBRACKET  = "]"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	DOT       = "."
	COMMA     = ","
	SEMICOLON = ";"
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	FSLASH    = "/"
	AMP       = "&"
	BAR       = "|"
	LT        = "<"
	GT        = ">"
	NOT       = "~"
	QUOT      = `"`

	CLASS       = "class"
	CONSTRUCTOR = "constructor"
	FUNCTION    = "function"
	METHOD      = "method"
	FIELD       = "field"
	STATIC      = "static"
	VAR         = "var"
	VOID        = "void"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"
)

var keywords = map[string]TokenType{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

func LookupIdent(ident string) TokenType {
	if tk, ok := keywords[ident]; ok {
		return tk
	}
	return IDENT
}
