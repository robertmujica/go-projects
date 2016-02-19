package main

import(
    "fmt"
    "os"
    "path/filepath"
)

func main(){
    dir, err := os.Open(".")
    if err != nil {
        return
    }
    defer dir.Close()
    
    fileInfor, err := dir.Readdir(-1)
    
    if err != nil {
        return
    }
    
    for _, fi := range fileInfor{
        fmt.Println(fi.Name)
    }
}

func walkDirectory(){
    filepath.Walk(".", func (path string, 
    info os.FileInfo, err error) error  {
        fmt.Println(path)
        return nil
    })
}