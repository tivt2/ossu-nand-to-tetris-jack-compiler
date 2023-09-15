package symbol_table

import (
	"fmt"
	"log"
)

type TableRow struct {
	Name        string
	Symbol_type string
	Kind        string
	id          int
}

type SymbolTable struct {
	ClassName string
	table     map[string]TableRow

	fieldCount    int
	staticCount   int
	argumentCount int
	localCount    int
}

func New() *SymbolTable {
	return &SymbolTable{
		table: make(map[string]TableRow),
	}
}

func (st *SymbolTable) Reset() {
	st.fieldCount = 0
	st.staticCount = 0
	st.argumentCount = 0
	st.localCount = 0
}

func (st *SymbolTable) Define(name string, symbol_type string, kind string) {
	switch kind {
	case "field":
		st.table[name] = TableRow{Symbol_type: symbol_type, Kind: kind, id: st.fieldCount}
		st.fieldCount++
	case "static":
		st.table[name] = TableRow{Symbol_type: symbol_type, Kind: kind, id: st.staticCount}
		st.staticCount++
	case "argument":
		st.table[name] = TableRow{Symbol_type: symbol_type, Kind: kind, id: st.argumentCount}
		st.argumentCount++
	case "local":
		st.table[name] = TableRow{Symbol_type: symbol_type, Kind: kind, id: st.localCount}
		st.localCount++
	default:
		log.Fatalf("Error in Define(), non-exaustive, received kind: %s", kind)
	}
}

func (st *SymbolTable) VarCount(kind string) int {
	switch kind {
	case "field":
		return st.fieldCount - 1
	case "static":
		return st.staticCount - 1
	case "argument":
		return st.argumentCount - 1
	case "local":
		return st.localCount - 1
	default:
		log.Fatalf("Error in VarCount(), non-exaustive, received kind: %s", kind)
	}
	return -1
}

func (st *SymbolTable) KindOf(name string) string {
	out, ok := st.table[name]
	if !ok {
		return ""
	}
	return out.Kind
}

func (st *SymbolTable) TypeOf(name string) string {
	out, ok := st.table[name]
	if !ok {
		return ""
	}
	return out.Symbol_type
}

func (st *SymbolTable) IndexOf(name string) int {
	out, ok := st.table[name]
	if !ok {
		return -1
	}
	return out.id
}

func (st *SymbolTable) Print() {
	fmt.Println(st)
}
