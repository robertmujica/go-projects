package copyFileUtil

import(
    "os"
    "log"
    "encoding/csv"
)

// Open : This method declaration means csvReader implements openFileHandler interface,
// but we don't need to explicitly declare it does so.
func (reader *CsvReader) Open() {

	fileDetails := FileContent{}
	file, err := os.Open(reader.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	sourceRows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	fileDetails.File = file
	fileDetails.Content = sourceRows
	reader.File.FileDetails = fileDetails
}