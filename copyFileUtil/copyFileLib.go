package copyFileUtil

import (
	"os"
	"strings"
)

// FileContent defines a generic file content with a pointer to file and
// a content property of interface[] also knwon as empty interface
// FileContent is an exported type that
type FileContent struct {
	File    *os.File
	Content interface{} // this type can hold value of any type
}

// File defines a File basic structure
type File struct {
	FilePath    string
	FileDetails FileContent
}

// CsvReader uses embedded type also known as anonymous field,
// this will help us to define a is-a relationship
type CsvReader struct {
	File
}

// TxtWriter uses embedded type also known as anonymous field,
// this will help us to define a is-a relationship
type TxtWriter struct {
	File
}

type record struct {
	IP       string
	HostName string
}

// OpenFileHandler declares an interface with a generic open file operation.
type OpenFileHandler interface {
	Open()
}

// FileOperations declares an interface with a set of supported operations on a File
type FileOperations interface {
	AddEntries(dest *os.File, destLines []string, sourceRows [][]string)
	RemoveEntries(dest *os.File, destLines []string, sourceRows [][]string)
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
