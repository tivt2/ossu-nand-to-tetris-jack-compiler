package vmWriter

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

var arithmetic = map[string]string{
	"+":   "add",
	"-":   "sub",
	"neg": "neg",
	"=":   "eq",
	">":   "gt",
	"<":   "lt",
	"&":   "and",
	"|":   "or",
	"~":   "not",
}

type VMWriter struct {
	file *os.File
	Out  bytes.Buffer
}

func New(filePath string) *VMWriter {
	path := filePath[:len(filePath)-5] + ".vm"
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error while trying to create file: %s", path)
	}

	return &VMWriter{
		file: file,
	}
}

func (w *VMWriter) WriteComment(msg string) {
	w.Out.WriteString(fmt.Sprintf("// %s", msg))
}

func (w *VMWriter) WritePush(segment string, index int) {
	w.Out.WriteString(fmt.Sprintf("push %s %d\n", segment, index))
}

func (w *VMWriter) WritePop(segment string, index int) {
	w.Out.WriteString(fmt.Sprintf("pop %s %d\n", segment, index))
}

func (w *VMWriter) WriteArithmetic(command string) {
	w.Out.WriteString(arithmetic[command] + "\n")
}

func (w *VMWriter) WriteLabel(label string, iteration int) {
	w.Out.WriteString(fmt.Sprintf("label %s%d\n", label, iteration))
}

func (w *VMWriter) WriteGoto(label string, iteration int) {
	w.Out.WriteString(fmt.Sprintf("goto %s%d\n", label, iteration))
}

func (w *VMWriter) WriteIf(label string, iteration int) {
	w.Out.WriteString(fmt.Sprintf("if-goto %s%d\n", label, iteration))
}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.Out.WriteString(fmt.Sprintf("call %s %d\n", name, nArgs))
	w.Out.WriteString("pop temp 0\n")
}

func (w *VMWriter) WriteFunction(name string, nVars int) {
	w.Out.WriteString(fmt.Sprintf("function %s %d\n", name, nVars))
}

func (w *VMWriter) WriteReturn() {
	w.Out.WriteString("return\n")
}

func (w *VMWriter) Close() {
	w.file.WriteString(w.Out.String())
	w.file.Close()
}
