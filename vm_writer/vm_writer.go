package vm_writer

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/tivt2/jack-compiler/vm_commands"
)

type vm_writer struct {
	outputPath string
	vmCommands bytes.Buffer
}

func New(outputPath string) *vm_writer {

	return &vm_writer{
		outputPath: outputPath,
	}
}

func (w *vm_writer) WritePush(segment vm_commands.Segment, index int) {
	w.vmCommands.WriteString(fmt.Sprintf("push %s %d", segment, index))
}

func (w *vm_writer) WritePop(segment vm_commands.Segment, index int) {
	w.vmCommands.WriteString(fmt.Sprintf("pop %s %d", segment, index))
}

func (w *vm_writer) WriteArithmetic(command vm_commands.Command) {
	w.vmCommands.WriteString(string(command))
}

func (w *vm_writer) WriteLabel(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("(%s)", label))
}

func (w *vm_writer) WriteGoto(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("goto %s", label))
}

func (w *vm_writer) WriteIf(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("if-goto %s", label))
}

func (w *vm_writer) WriteCall(name string, nArgs int) {
	w.vmCommands.WriteString(fmt.Sprintf("call %s %d", name, nArgs))
}

func (w *vm_writer) WriteFunction(label string) {
	w.vmCommands.WriteString(fmt.Sprintf("(%s)", label))
}

func (w *vm_writer) WriteReturn() {
	w.vmCommands.WriteString("return")
}

func (w *vm_writer) Close() {
	file, err := os.Create(w.outputPath)
	if err != nil {
		log.Fatalf("Error when trying to create file: %s\n%v", w.outputPath, err)
	}
	defer file.Close()

	file.WriteString(w.vmCommands.String())
}
