package symbolTable

import (
	"log"
)

type tableRow struct {
	DecType string
	kind    string
	id      int
}

type SymbolTable struct {
	classLevel      map[string]*tableRow
	subroutineLevel map[string]*tableRow

	fieldCounter    int
	staticCounter   int
	localCounter    int
	argumentCounter int
}

func New() *SymbolTable {
	return &SymbolTable{
		classLevel:      make(map[string]*tableRow),
		subroutineLevel: make(map[string]*tableRow),
	}
}

func (sb *SymbolTable) Reset() {
	sb.fieldCounter = 0
	sb.staticCounter = 0
	sb.argumentCounter = 0
	sb.localCounter = 0
}

func (sb *SymbolTable) Define(name string, decType string, kind string) {
	switch kind {
	case "field":
		sb.classLevel[name] = &tableRow{kind: kind, id: sb.fieldCounter, DecType: decType}
		sb.fieldCounter++
	case "static":
		sb.classLevel[name] = &tableRow{kind: kind, id: sb.staticCounter, DecType: decType}
		sb.staticCounter++
	case "argument":
		sb.subroutineLevel[name] = &tableRow{kind: kind, id: sb.argumentCounter, DecType: decType}
		sb.argumentCounter++
	case "local":
		sb.subroutineLevel[name] = &tableRow{kind: kind, id: sb.localCounter, DecType: decType}
		sb.localCounter++
	default:
		log.Fatalf("Wrong table data, received: {name: %s,decType: %s,kind: %s}", name, decType, kind)
	}
}

func (sb *SymbolTable) VarCount(kind string) int {
	out := 0
	var table map[string]*tableRow
	switch kind {
	case "field", "static":
		table = sb.classLevel
	case "argument", "local":
		table = sb.subroutineLevel
	}
	for _, val := range table {
		if val.kind == kind {
			out++
		}
	}
	return out
}

func (sb *SymbolTable) KindOf(name string) string {
	if row, ok := sb.classLevel[name]; ok {
		return row.kind
	}
	row := sb.subroutineLevel[name]
	return row.kind
}

func (sb *SymbolTable) TypeOf(name string) string {
	if row, ok := sb.classLevel[name]; ok {
		return row.DecType
	}
	row := sb.subroutineLevel[name]
	return row.DecType
}

func (sb *SymbolTable) IndexOf(name string) int {
	if row, ok := sb.classLevel[name]; ok {
		return row.id
	}
	row := sb.subroutineLevel[name]
	return row.id
}
