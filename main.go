package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"com.butiricristian/ancpi-data-provider/api"
	"com.butiricristian/ancpi-data-provider/models"
	"com.butiricristian/ancpi-data-provider/parserjob"
)

func openFile(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func saveToFile(data []*models.MonthlyData) {
	fileName := "data/data.json"
	dataFile, err := openFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dataFile.Close()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling data: %+v", err)
		return
	}
	dataFile.Write(dataBytes)
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
