package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"com.butiricristian/ancpi-data-provider/api"
	"com.butiricristian/ancpi-data-provider/parserjob"
)

func openFile(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func saveToFile(data map[string]map[string][]parserjob.StateData) {
	for key, currentData := range data {
		fileName := fmt.Sprintf("data/%s.json", strings.ToLower(key))
		dataFile, err := openFile(fileName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer dataFile.Close()

		dataBytes, err := json.Marshal(currentData)
		if err != nil {
			fmt.Printf("Error marshaling data: %+v", err)
			continue
		}
		dataFile.Write(dataBytes)
	}
}

func getAllData() {
	excelUrls := parserjob.FindAllExcelUrls()
	data := parserjob.GetDataFromExcels(excelUrls)

	saveToFile(data)
}

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	getAllData()
	api.StartServer()
}
