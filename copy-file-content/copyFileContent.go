package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// defines a generic file content with a pointer to file and
// a content property of interface[] also knwon as empty interface
type fileContent struct {
	file    *os.File
	content interface{} // this type can hold value of any type
}

type file struct {
	filePath    string
	fileDetails fileContent
}

// Here I'm using embedded type also known as anonymous field,
// this will help us to define a is-a relationship
type csvReader struct {
	file
}

// Here I'm using embedded type also known as anonymous field,
// this will help us to define a is-a relationship
type txtWriter struct {
	file
}

type record struct {
	IP       string
	HostName string
}

// declare an interface with a generic open file operation.
type openFileHandler interface {
	open()
}

// declare an interface with a set of supported operations on a file
type fileOperations interface {
	addEntries(dest *os.File, destLines []string, sourceRows [][]string)
	removeEntries(dest *os.File, destLines []string, sourceRows [][]string)
}

// This method means csvReader implements openFileHandler interface,
// but we don't need to explicitly declare it does so.
func (reader *csvReader) open() {

	fileDetails := fileContent{}
	file, err := os.Open(reader.filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	sourceRows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	fileDetails.file = file
	fileDetails.content = sourceRows
	reader.file.fileDetails = fileDetails
}

// This method means txtWriter implements openFileHandler interface
func (writer *txtWriter) open() {

	fileDetails := fileContent{}
	file, destLines, err := readDestinationFile(writer.filePath)
	if err != nil {
		log.Fatalln(err)
	}
	fileDetails.file = file
	fileDetails.content = destLines
	writer.file.fileDetails = fileDetails
}

// This method means txtWriter implements fileOperations interface
func (writer txtWriter) addEntries(sourceContent interface{}) {
	var sourceRows = sourceContent.([][]string)                // using type assertion to access sourceContent underlying concrete value
	var destLines = writer.file.fileDetails.content.([]string) // using type assertion

	for _, row := range sourceRows {
		record := makeRecord(row)

		if indexOf(record, destLines) < 0 {
			var fileEntry = fmt.Sprintf("\n%[1]s %[2]s", record.IP, record.HostName)
			writer.file.fileDetails.file.WriteString(fileEntry)
		} else {
			fmt.Printf("Ip already exists: %[1]s \n", record.IP)
		}
	}
}

// This method means txtWriter implements fileOperations interface
func (writer txtWriter) removeEntries(sourceContent interface{}) {
	var sourceRows = sourceContent.([][]string) // using type assertion to access sourceContent underlying concrete value
	var destLines = writer.file.fileDetails.content.([]string)

	var list = convertToRecordArray(sourceRows)
	var newLines []string
	for _, line := range destLines {
		var found = false
		var rowValues = strings.Fields(line)
		if len(rowValues) > 1 {
			record := makeRecord(rowValues)

			if contains(record, list) {
				found = true
			}
		}

		if !found {
			newLines = append(newLines, line)
		}
	}

	writer.file.fileDetails.file.Truncate(0)
	newFileContent := strings.Join(newLines, "\n")
	_, err := writer.file.fileDetails.file.Write([]byte(newFileContent))
	if err != nil {
		log.Fatalln(err)
	}
}

func makeRecord(row []string) record {
	return record{
		IP:       row[0],
		HostName: row[1],
	}
}

func readDestinationFile(path string) (*os.File, []string, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return file, lines, scanner.Err()
}

// A field is separated by one or more space characters.
// The strings.Fields() method in the strings package separates these into an array.
// It splits on groups of spaces
func indexOf(pRecord record, fileRows []string) int {

	for i := 1; i < len(fileRows); i++ {
		var rowValues = strings.Fields(fileRows[i])

		if len(rowValues) > 1 {
			var ip = rowValues[0]

			if pRecord.IP == ip {
				return i
			}
		}
	}
	return -1
}

func contains(pRecord record, list []record) bool {

	for _, row := range list {
		if pRecord.IP == row.IP {
			return true
		}
	}
	return false
}

func convertToRecordArray(sourceRows [][]string) []record {

	var list []record
	for _, row := range sourceRows {
		record := makeRecord(row)
		list = append(list, record)
	}
	return list
}

func main() {
	command := flag.String("command", "add", "command to be executed [add, remove]")
	sourceFilePath := flag.String("source", "", "source file path")
	targetPath := flag.String("target", "", "destination file path")
	flag.Parse()

	source := csvReader{file: file{filePath: *sourceFilePath}}
	source.open()
	target := txtWriter{file: file{filePath: *targetPath}}
	target.open()

	switch *command {
	case "add":
		target.addEntries(source.file.fileDetails.content)
	case "remove":
		target.removeEntries(source.file.fileDetails.content)
	default:
		fmt.Print("\n try: -command add or -command remove")
	}
}
