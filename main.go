package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/tivt2/jack-compiler/jackCompiler"
)

func main() {
	if len(os.Args) == 0 {
		log.Fatal("1 Usage 'JackCompiler <filename.jack | foldername>'")
	}
	path := os.Args[1]

	if filepath.Ext(path) == ".jack" {
		jc := jackCompiler.New(path)
		jc.Compile()
		return
	}

	info, err := os.Stat(path)
	checkErr(err, "checking file stats")
	if info.IsDir() {
		folder, err := os.Open(path)
		checkErr(err, "error opening folder")
		files, err2 := folder.Readdirnames(0)
		checkErr(err2, "error reading folder files")

		var wg sync.WaitGroup
		for _, file := range files {
			file := file
			if filepath.Ext(file) == ".jack" {
				wg.Add(1)
				go func() {
					jc := jackCompiler.New(filepath.Join(path, file))
					jc.Compile()
					wg.Done()
				}()
			}
		}
		wg.Wait()
		return
	}

	log.Fatal("2 Usage 'JackCompiler <filename.jack | foldername>'")

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%v, message: %s", err, msg)
	}
}
