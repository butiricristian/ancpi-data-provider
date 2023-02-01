package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

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

func scrapeAllData() []*models.MonthlyData {
	excelUrls := parserjob.FindAllExcelUrls()
	data := parserjob.GetDataFromExcels(excelUrls)

	saveToFile(data)
	return data
}

var Data []*models.MonthlyData

func PrepareData() {
	fmt.Println("Preparing data...")
	fileName := "data/data.json"
	fileData, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File data.json not found. Starting scraping...")
		Data = scrapeAllData()
		fmt.Println("Data retrieved from scraping")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(fileData, &Data)
	fmt.Println("Data retrieved from data.json")
}
