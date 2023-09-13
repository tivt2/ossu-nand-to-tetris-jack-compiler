package syntax_analyzer

import (
	"fmt"
	"log"
	"os"

	"github.com/tivt2/jack-compiler/parser"
	"github.com/tivt2/jack-compiler/tokenizer"
)

func ParseXMLTree(filePath string) {
	tkzr := tokenizer.New(filePath)
	parser := parser.New(tkzr)

	ParseTree, TokensXML := parser.CompileClass()

	path := filePath[:len(filePath)-5]
	writeToFile(path, ".xml", ParseTree)
	writeToFile(path+"T", ".xml", TokensXML)
	// fmt.Println(ParseTree)
	// fmt.Println(TokensXML)
}

func writeToFile(filePath string, ext string, content string) {
	path := filePath + ext
	fmt.Println("Creating file ->", path, "<-")

	file, err := os.Create(path)
	checkErr(err, fmt.Sprintf("Error trying to create file: %s", path))
	defer file.Close()

	file.WriteString(content)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:\n%v", msg, err)
	}
}
