package parserjob

import (
	"fmt"
	"strings"
	"sync"

	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
	"github.com/schollz/progressbar/v3"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/maps"
)

type ParseResult struct {
	dataType    string
	dateKey     string
	dataCereri  []*models.CereriStateData
	dataIpoteci []*models.IpoteciStateData
	dataVanzari []*models.VanzariStateData
}

func getNrOfHeaders(rows [][]string) int {
	headers := 0
	for headers = 0; rows[headers] == nil || len(rows[headers]) < 2 || rows[headers][1] != "ALBA"; headers++ {
	}

	return headers
}

func parseExcelVanzari(rows [][]string) []*models.VanzariStateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 43
	var data []*models.VanzariStateData = make([]*models.VanzariStateData, nrRows)
	for i := 0; i < nrRows; i++ {
		row := rows[i+HEADER_ROWS]
		if len(row) <= 2 || row[1] == "" {
			continue
		}

		currentData := models.CreateVanzariData(row)
		data[i] = &currentData
	}

	return data
}

func parseExcelIpoteci(rows [][]string) []*models.IpoteciStateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 43
	var data []*models.IpoteciStateData = make([]*models.IpoteciStateData, 2*nrRows)
	for i := 0; i < nrRows; i++ {
		row := rows[i+HEADER_ROWS]
		if len(row) <= 2 || row[1] == "" {
			continue
		}

		var currentDataActive models.IpoteciStateData
		var currentDataInactive models.IpoteciStateData
		if len(row) > 8 {
			currentDataActive, currentDataInactive = models.CreateIpoteciData(row)
		} else {
			currentDataActive, currentDataInactive = models.CreateIpoteciDataV2(row)
		}
		data[2*i] = &currentDataActive
		data[2*i+1] = &currentDataInactive
	}

	return data
}

func parseExcelCereri(rows [][]string) []*models.CereriStateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 42 * 5
	var data []*models.CereriStateData = make([]*models.CereriStateData, nrRows+5)
	currentTotalOnline := 0
	currentTotalGhiseu := 0
	totalOnline := map[string]int{"altele": 0, "informare": 0, "receptie": 0, "inscriere": 0, "total": 0}
	totalGhiseu := map[string]int{"altele": 0, "informare": 0, "receptie": 0, "inscriere": 0, "total": 0}

	for i := 0; i < nrRows; i++ {
		var row []string
		rowIndex := i/5*4 + i%5
		stateName := rows[i/5*4+HEADER_ROWS][1]

		if i%5 == 4 {
			row = []string{"", stateName, "Total", fmt.Sprint(currentTotalOnline), fmt.Sprint(currentTotalGhiseu), fmt.Sprint(currentTotalOnline + currentTotalGhiseu)}
			currentTotalGhiseu = 0
			currentTotalOnline = 0
		} else {
			row = rows[rowIndex+HEADER_ROWS]
			if len(row) <= 2 {
				continue
			}
			if row[1] == "" {
				row[1] = stateName
			}
			currentTotalOnline += helpers.ConvertToInt(row[3])
			currentTotalGhiseu += helpers.ConvertToInt(row[4])
			totalOnline[strings.ToLower(row[2])] += helpers.ConvertToInt(row[3])
			totalGhiseu[strings.ToLower(row[2])] += helpers.ConvertToInt(row[4])
			totalOnline["total"] += helpers.ConvertToInt(row[3])
			totalGhiseu["total"] += helpers.ConvertToInt(row[4])
		}

		currentData := models.CreateCereriData(row)
		data[i] = &currentData
	}

	// Store Total number of requests at national level
	i := 0
	for key := range totalOnline {
		currentData := models.CreateCereriData([]string{
			"",
			"TOTAL",
			key,
			fmt.Sprint(totalOnline[key]),
			fmt.Sprint(totalGhiseu[key]),
			fmt.Sprint(totalGhiseu[key] + totalOnline[key]),
		})
		data[nrRows+i] = &currentData
		i += 1
	}

	return data
}

func ParseExcel(excelUrl *ExcelUrl, dataChannel chan<- *ParseResult, wg *sync.WaitGroup) {
	// fmt.Printf("Parsing excel %s, %s - %s\n", excelUrl.month, excelUrl.year, excelUrl.name)
	defer wg.Done()

	body, ok := requestPage(excelUrl.url)
	if !ok {
		fmt.Printf("requested page could not be found: %s\n", excelUrl.url)
		return
	}
	defer body.Close()

	doc, err := excelize.OpenReader(body)
	if err != nil {
		fmt.Printf("An error occurred while reading excel file: %v\n", err)
		return
	}
	defer doc.Close()

	sheetName := doc.GetSheetName(0)
	rows, err := doc.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}

	dateKey := fmt.Sprintf("%s, %s", excelUrl.month, excelUrl.year)
	switch excelUrl.name {
	case "VANZARI":
		data := parseExcelVanzari(rows)
		dataChannel <- &ParseResult{
			dataType:    excelUrl.name,
			dateKey:     dateKey,
			dataVanzari: data,
		}
	case "IPOTECI":
		data := parseExcelIpoteci(rows)
		dataChannel <- &ParseResult{
			dataType:    excelUrl.name,
			dateKey:     dateKey,
			dataIpoteci: data,
		}
	case "CERERI":
		data := parseExcelCereri(rows)
		dataChannel <- &ParseResult{
			dataType:   excelUrl.name,
			dateKey:    dateKey,
			dataCereri: data,
		}
	}
}

func GetDataFromExcels(excelUrls []*ExcelUrl) []*models.MonthlyData {
	data := map[string]*models.MonthlyData{}

	dataChannel := make(chan *ParseResult)
	var wg sync.WaitGroup

	for _, excelUrl := range excelUrls {
		wg.Add(1)
		go ParseExcel(excelUrl, dataChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	total := len(excelUrls)
	bar := progressbar.NewOptions(
		total,
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionSetDescription("[2/2] Parsing excels: "),
	)

	for receivedData := range dataChannel {
		if _, ok := data[receivedData.dateKey]; !ok {
			data[receivedData.dateKey] = &models.MonthlyData{
				CurrentDate: helpers.ConvertToTime(receivedData.dateKey),
			}
		}

		switch receivedData.dataType {
		case "CERERI":
			data[receivedData.dateKey].CereriData = receivedData.dataCereri
		case "VANZARI":
			data[receivedData.dateKey].VanzariData = receivedData.dataVanzari
		case "IPOTECI":
			data[receivedData.dateKey].IpoteciData = receivedData.dataIpoteci
		}

		bar.Add(1)
	}
	fmt.Println()

	return maps.Values(data)
}
