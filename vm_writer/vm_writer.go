package vm_writer

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/tivt2/jack-compiler/vm_commands"
)

type VMWriter struct {
	outputPath string
	vmCommands bytes.Buffer
}

func New(outputPath string, ext string) *VMWriter {

	return &VMWriter{
		outputPath: outputPath + ext,
	}
}

func (w *VMWriter) WritePush(segment vm_commands.Segment, index int) {
	w.vmCommands.WriteString(fmt.Sprintf("push %s %d", segment, index))
}

func (w *VMWriter) WritePop(segment vm_commands.Segment, index int) {
	w.vmCommands.WriteString(fmt.Sprintf("pop %s %d", segment, index))
}

func (w *VMWriter) WriteArithmetic(command vm_commands.Command) {
	w.vmCommands.WriteString(string(command))
}

func (w *VMWriter) WriteLabel(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("(%s)", label))
}

func (w *VMWriter) WriteGoto(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("goto %s", label))
}

func (w *VMWriter) WriteIf(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("if-goto %s", label))
}

func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.vmCommands.WriteString(fmt.Sprintf("call %s %d", name, nArgs))
}

func (w *VMWriter) WriteFunction(label string, nVars int) {
	w.vmCommands.WriteString(fmt.Sprintf("function %s %d", label, nVars))
}

func (w *VMWriter) WriteReturn() {
	w.vmCommands.WriteString("return")
}

func (w *VMWriter) Close() {
	file, err := os.Create(w.outputPath)
	if err != nil {
		log.Fatalf("Error when trying to create file: %s\n%v", w.outputPath, err)
	}
	defer file.Close()

	file.WriteString(w.vmCommands.String())
}
