package vm_commands

type Segment string

const (
	ARGUMENT = "argument"
	LOCAL    = "local"
	STATIC   = "static"
	THIS     = "this"
	THAT     = "that"
	POINTER  = "pointer"
	TEMP     = "temp"
)

type Command string

const (
	ADD = "add"
	SUB = "sub"
	NEG = "neg"
	EQ  = "eq"
	GT  = "gt"
	LT  = "lt"
	AND = "and"
	OR  = "or"
	NOT = "not"
)
