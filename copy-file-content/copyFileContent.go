package main

import(
    "encoding/csv"
    "fmt"
    "os"
    "bufio"
    "strings"
    "log"
    "flag"
)

type record struct{
    IP string
    HostName string
}

func makeRecord(row []string) record{
    return record{
        IP: row[0],
        HostName : row[1],
    }
}

func readDestinationFile(path string) (*os.File, []string, error){
    file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
    if err != nil{
        return nil, nil, err
    }
   
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan(){
        lines = append(lines, scanner.Text())
    }
    return file, lines, scanner.Err()
}

// A field is separated by one or more space characters. 
// The strings.Fields() method in the strings package separates these into an array. 
// It splits on groups of spaces
func indexOf(pRecord record, fileRows []string) int{
    
    for i:= 1; i < len(fileRows); i++ {
        var rowValues = strings.Fields(fileRows[i])

        if len(rowValues) > 1{
            var ip = rowValues[0]

            if pRecord.IP == ip {
                return i      
            }
        }
    }
    return -1
}

func contains(pRecord record, list []record) bool{
    
    for _, row := range list{
        if pRecord.IP == row.IP {
            return true
        }
    }
    return false
}

func addEntries(dest *os.File, destLines []string, sourceRows [][]string) {
    for _, row := range sourceRows {
        record := makeRecord(row)
        
        if indexOf(record, destLines) < 0{
            var fileEntry = fmt.Sprintf("\n%[1]s %[2]s", record.IP, record.HostName)
            dest.WriteString(fileEntry)
        } else {
            fmt.Printf("Ip already exists: %[1]s \n", record.IP)
        }
    }
}

func convertToRecordArray(sourceRows [][]string) []record{
    
    var list []record
    for _, row := range sourceRows{
        record := makeRecord(row)
        list = append(list, record)
    }
    return list
}

func removeEntries(dest *os.File, destLines []string, sourceRows [][]string){
    
    var list = convertToRecordArray(sourceRows)
    var newLines []string
    for _, line := range destLines {
        var found = false
        var rowValues = strings.Fields(line)
        if len(rowValues) > 1{
            record := makeRecord(rowValues)
            
            if contains(record, list){
                found = true
            }
        }
        
        if !found {
            newLines = append(newLines, line)
        }
    }
    
    dest.Truncate(0)
    newFileContent := strings.Join(newLines, "\n")
    _, err := dest.Write([]byte(newFileContent))
    if err != nil {
        log.Fatalln(err)
    }
}

func main() {
    command := flag.String("command", "add", "commands to be executed [add or remove]")
    sourceFilePath := flag.String("source", "", "source file path")
    targetPath := flag.String("target", "", "destination path path")
    
    flag.Parse()
    f, err := os.Open(*sourceFilePath)
    if err != nil {
        log.Fatalln(err)
    }
    defer f.Close()
    
    rdr := csv.NewReader(f)
    sourceRows, err := rdr.ReadAll()
    if err != nil {
        panic(err)
    }
    
    dest, destLines, errDest := readDestinationFile(*targetPath)
    if errDest != nil{
        panic(errDest)
    }
    defer dest.Close()
    
    switch *command {
    case "add":
        addEntries(dest, destLines, sourceRows)  
    case "remove":
        removeEntries(dest, destLines, sourceRows)
    default:
        fmt.Print("try: add or remove")
    }
}