package copyFileUtil

import (
	"bufio"
	"fmt"
    "log"
    "strings"
    "os"
)

// Open : This method means txtWriter implements openFileHandler interface
func (writer *TxtWriter) Open() {

	fileDetails := FileContent{}
	file, destLines, err := readDestinationFile(writer.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	fileDetails.File = file
	fileDetails.Content = destLines
	writer.File.FileDetails = fileDetails
}

// AddEntries : This method means txtWriter implements fileOperations interface
func (writer TxtWriter) AddEntries(sourceContent interface{}) {
	var sourceRows = sourceContent.([][]string)                // using type assertion to access sourceContent underlying concrete value
	var destLines = writer.File.FileDetails.Content.([]string) // using type assertion

	for _, row := range sourceRows {
		record := makeRecord(row)

		if indexOf(record, destLines) < 0 {
			var fileEntry = fmt.Sprintf("\n%[1]s %[2]s", record.IP, record.HostName)
			writer.File.FileDetails.File.WriteString(fileEntry)
		} else {
			fmt.Printf("Ip already exists: %[1]s \n", record.IP)
		}
	}
}

// RemoveEntries : This method means txtWriter implements fileOperations interface
func (writer TxtWriter) RemoveEntries(sourceContent interface{}) {
	var sourceRows = sourceContent.([][]string) // using type assertion to access sourceContent underlying concrete value
	var destLines = writer.File.FileDetails.Content.([]string)

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

	writer.File.FileDetails.File.Truncate(0)
	newFileContent := strings.Join(newLines, "\n")
	_, err := writer.File.FileDetails.File.Write([]byte(newFileContent))
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