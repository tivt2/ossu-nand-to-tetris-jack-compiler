package syntaxAnalyzer

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tivt2/jack-compiler/parseTree"
	"github.com/tivt2/jack-compiler/parser"
	"github.com/tivt2/jack-compiler/tokenizer"
)

func ParseTree(filePath string) *parseTree.Class {
	file, err := os.ReadFile(filePath)
	checkErr(err, fmt.Sprintf("Error when opening file %s", filePath))

	fileContent := removeComments(strings.TrimSpace(string(file)))

	tkzr := tokenizer.New(fileContent)
	parser := parser.New(tkzr)

	return parser.ParseClass()
}

func removeComments(text string) string {
	regexes := []string{
		`\/\*[^*]*\*\/`,
		`\/\*\*[\s\S]*?\*\/`,
		`\/\/[^\n]*`,
	}

	for _, pattern := range regexes {
		regex := regexp.MustCompile(pattern)
		text = regex.ReplaceAllString(text, "")
	}

	return text
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatal(err)
	}
}
