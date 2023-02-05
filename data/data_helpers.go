package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func saveToFile(fileName string, data []*models.MonthlyData) {
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

func scrapeAllData(fileName string) []*models.MonthlyData {
	excelUrls := parserjob.FindAllExcelUrls()
	data := parserjob.GetDataFromExcels(excelUrls)

	saveToFile(fileName, data)
	return data
}

var Data []*models.MonthlyData

func PrepareData(fileName string) {
	fmt.Println("Preparing data...")
	fileData, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File data.json not found. Starting scraping...")
		Data = scrapeAllData(fileName)
		fmt.Println("Data retrieved from scraping")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(fileData, &Data)
	fmt.Println("Data retrieved from data.json")
}

func getDataUrl() string {
	if os.Getenv("APP_ENV") == "production" {
		return "https://ancpi-data-provider.netlify.app/data/data.json"
	}
	return "http://localhost:8888/data/data.json"
}

func PrepareDataFromUrl() {
	url := getDataUrl()
	fmt.Printf("Reading data from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("An error occured while retrieving the page: %v", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An error occured while reading the page: %v", err)
		return
	}
	fmt.Println(data)

	err = json.Unmarshal(data, &Data)
	if err != nil {
		fmt.Printf("An error occured while unmarshaling the page: %v", err)
		return
	}
	fmt.Println(Data)

	fmt.Println("Data retrieved from data.json")
}
