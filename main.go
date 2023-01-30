package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func openFile() (*os.File, error) {
	f, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func saveToFile(data map[string]map[string][]StateData, f *os.File) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling data: %+v", err)
		return
	}
	f.Write(dataBytes)
}

func getAllData() {
	dataFile, err := openFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dataFile.Close()

	excelUrls := findAllExcelUrls()
	data := getDataFromExcels(excelUrls)

	saveToFile(data, dataFile)
}

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	getAllData()
}
