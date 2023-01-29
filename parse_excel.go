package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

type RequestType int64

const (
	ALTELE RequestType = iota
	INFORMARE
	INSCRIERE
	RECEPTIE
	UNDEFINED
)

func (rt RequestType) String() string {
	switch rt {
	case ALTELE:
		return "Altele"
	case INFORMARE:
		return "Informare"
	case INSCRIERE:
		return "Inscriere"
	case RECEPTIE:
		return "Receptie"
	}
	return "Unknown"
}

func getRequestType(val string) RequestType {
	switch val {
	case "altele":
		return ALTELE
	case "informare":
		return INFORMARE
	case "inscriere":
		return INSCRIERE
	case "receptie":
		return RECEPTIE
	}
	return UNDEFINED
}

type StateData interface {
	printData() string
}

type VanzariStateData struct {
	name            string
	agricol         int
	neagricol       int
	constructie     int
	faraConstructie int
	total           int
}

type IpoteciStateData struct {
	VanzariStateData
	active bool
}

type CereriStateData struct {
	name        string
	requestType RequestType
	online      int
	ghiseu      int
	total       int
}

type ParseResult struct {
	dataType string
	dateKey  string
	data     []StateData
}

func (data *VanzariStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func (data *IpoteciStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func (data *CereriStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func convertToInt(val string) int {
	val = strings.ReplaceAll(val, ",", "")
	converted, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("Value is not a number: %v", err)
		return -1
	}

	return converted
}

func createVanzariData(row []string) VanzariStateData {
	return VanzariStateData{
		name:            row[1],
		agricol:         convertToInt(row[2]),
		neagricol:       convertToInt(row[3]),
		constructie:     convertToInt(row[4]),
		faraConstructie: convertToInt(row[5]),
		total:           convertToInt(row[6]),
	}
}

func createIpoteciData(row []string) (IpoteciStateData, IpoteciStateData) {
	active := IpoteciStateData{
		VanzariStateData{
			name:            row[1],
			agricol:         convertToInt(row[2]),
			neagricol:       convertToInt(row[3]),
			constructie:     convertToInt(row[6]),
			faraConstructie: convertToInt(row[7]),
			total:           convertToInt(row[10]),
		},
		true,
	}
	inactive := IpoteciStateData{
		VanzariStateData{
			name:            row[1],
			agricol:         convertToInt(row[4]),
			neagricol:       convertToInt(row[5]),
			constructie:     convertToInt(row[8]),
			faraConstructie: convertToInt(row[9]),
			total:           convertToInt(row[11]),
		},
		false,
	}
	return active, inactive
}

func createCereriData(row []string) CereriStateData {
	online := convertToInt(row[3])
	ghiseu := convertToInt(row[4])
	return CereriStateData{
		name:        row[1],
		requestType: getRequestType(row[2]),
		online:      online,
		ghiseu:      ghiseu,
		total:       online + ghiseu,
	}
}

func getNrOfHeaders(rows [][]string) int {
	headers := 0
	for headers = 0; rows[headers] == nil || len(rows[headers]) < 2 || rows[headers][1] != "ALBA"; headers++ {
	}

	return headers
}

func parseExcelVanzari(rows [][]string) []StateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 43
	var data []StateData = make([]StateData, nrRows)
	for i := 0; i < nrRows; i++ {
		row := rows[i+HEADER_ROWS]
		if len(row) <= 2 || row[1] == "" {
			continue
		}

		currentData := createVanzariData(row)
		data[i] = &currentData
	}

	return data
}

func parseExcelIpoteci(rows [][]string) []StateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 43
	var data []StateData = make([]StateData, 2*nrRows)
	for i := 0; i < nrRows; i++ {
		row := rows[i+HEADER_ROWS]
		if len(row) <= 2 || row[1] == "" {
			continue
		}

		currentDataActive, currentDataInactive := createIpoteciData(row)
		data[2*i] = &currentDataActive
		data[2*i+1] = &currentDataInactive
	}

	return data
}

func parseExcelCereri(rows [][]string) []StateData {
	HEADER_ROWS := getNrOfHeaders(rows)
	nrRows := 42*4 + 1
	var data []StateData = make([]StateData, nrRows)

	for i := 0; i < nrRows; i++ {
		row := rows[i+HEADER_ROWS]
		if len(row) <= 2 {
			continue
		}
		if row[1] == "" {
			row[1] = rows[i/4*4+HEADER_ROWS][1]
		}

		currentData := createCereriData(row)
		data[i] = &currentData
	}

	return data
}

func parseExcel(excelUrl *ExcelUrl, dataChannel chan<- *ParseResult, wg *sync.WaitGroup) {
	fmt.Printf("Parsing excel %s, %s - %s\n", excelUrl.month, excelUrl.year, excelUrl.name)
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

	var data []StateData
	switch excelUrl.name {
	case "VANZARI":
		data = parseExcelVanzari(rows)
	case "IPOTECI":
		data = parseExcelIpoteci(rows)
	case "CERERI":
		data = parseExcelCereri(rows)
	}

	dateKey := fmt.Sprintf("%s, %s", excelUrl.month, excelUrl.year)
	dataChannel <- &ParseResult{
		dataType: excelUrl.name,
		dateKey:  dateKey,
		data:     data,
	}
}

func getDataFromExcels(excelUrls []*ExcelUrl) map[string]map[string][]StateData {
	data := make(map[string]map[string][]StateData)

	dataChannel := make(chan *ParseResult)
	var wg sync.WaitGroup

	for _, excelUrl := range excelUrls {
		wg.Add(1)
		go parseExcel(excelUrl, dataChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	for receivedData := range dataChannel {
		if data[receivedData.dateKey] == nil {
			data[receivedData.dateKey] = make(map[string][]StateData)
		}
		data[receivedData.dateKey][receivedData.dataType] = receivedData.data
	}

	return data
}

func printData(data map[string]map[string][]StateData) {
	for date, dateValues := range data {
		fmt.Printf("\n%s - VANZARI: \n", date)
		for _, stateData := range dateValues["VANZARI"] {
			fmt.Printf("%v", stateData.printData())
		}
		fmt.Printf("\n%s - IPOTECI: \n", date)
		for _, stateData := range dateValues["IPOTECI"] {
			fmt.Printf("%v", stateData.printData())
		}
		fmt.Printf("\n%s - CERERI: \n", date)
		for _, stateData := range dateValues["CERERI"] {
			fmt.Printf("%v", stateData.printData())
		}
		fmt.Println()
	}
}
