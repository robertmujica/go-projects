package main

import (
	"flag"
	"fmt"
	"github.com/robertmujica/copyFileUtil"
)

func main() {
	command := flag.String("command", "add", "command to be executed [add, remove]")
	sourceFilePath := flag.String("source", "", "source file path")
	targetPath := flag.String("target", "", "destination file path")
	flag.Parse()

	source := copyFileUtil.CsvReader{File: copyFileUtil.File{FilePath: *sourceFilePath}}
	source.Open()
	target := copyFileUtil.TxtWriter{File: copyFileUtil.File{FilePath: *targetPath}}
	target.Open()
   
	switch *command {
	case "add":
		target.AddEntries(source.File.FileDetails.Content)
	case "remove":
		target.RemoveEntries(source.File.FileDetails.Content)
	default:
		fmt.Print("\n try: -command add or -command remove")
	}
}
