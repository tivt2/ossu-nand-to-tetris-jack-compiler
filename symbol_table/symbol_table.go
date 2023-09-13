package symbol_table

import "log"

type tableRow struct {
	symbol_type string
	kind        string
	id          int
}

type symbolTable struct {
	table map[string]*tableRow

	fieldCount    int
	staticCount   int
	argumentCount int
	localCount    int
}

func New() *symbolTable {
	return &symbolTable{
		table: make(map[string]*tableRow),
	}
}

func (st *symbolTable) Reset() {
	st.fieldCount = 0
	st.staticCount = 0
	st.argumentCount = 0
	st.localCount = 0
}

func (st *symbolTable) Define(name string, symbol_type string, kind string) {
	switch kind {
	case "field":
		st.table[name] = &tableRow{symbol_type: symbol_type, kind: kind, id: st.fieldCount}
		st.fieldCount++
	case "static":
		st.table[name] = &tableRow{symbol_type: symbol_type, kind: kind, id: st.staticCount}
		st.staticCount++
	case "argument":
		st.table[name] = &tableRow{symbol_type: symbol_type, kind: kind, id: st.argumentCount}
		st.argumentCount++
	case "local":
		st.table[name] = &tableRow{symbol_type: symbol_type, kind: kind, id: st.localCount}
		st.localCount++
	default:
		log.Fatalf("Error in Define(), non-exaustive, received kind: %s", kind)
	}
}

func (st *symbolTable) VarCount(kind string) int {
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

func (st *symbolTable) KindOf(name string) string {
	out, ok := st.table[name]
	if !ok {
		return ""
	}
	return out.kind
}

func (st *symbolTable) TypeOf(name string) string {
	out, ok := st.table[name]
	if !ok {
		return ""
	}
	return out.symbol_type
}

func (st *symbolTable) IndexOf(name string) int {
	out, ok := st.table[name]
	if !ok {
		return -1
	}
	return out.id
}
