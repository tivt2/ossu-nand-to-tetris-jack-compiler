package jack_compiler

import (
	"fmt"

	"github.com/tivt2/jack-compiler/syntax_analyzer"
)

func CompileFile(filePath string) {
	parseTree := syntax_analyzer.ParseXMLTree(filePath)

	fmt.Println(parseTree)
}
