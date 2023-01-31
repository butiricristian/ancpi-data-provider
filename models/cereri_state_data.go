package models

import (
	"fmt"

	"com.butiricristian/ancpi-data-provider/helpers"
)

type CereriStateData struct {
	Name        string
	RequestType RequestType
	Online      int
	Ghiseu      int
	Total       int
}

func (data *CereriStateData) printData() string {
	return fmt.Sprintf("%v", *data)
}

func CreateCereriData(row []string) CereriStateData {
	online := helpers.ConvertToInt(row[3])
	ghiseu := helpers.ConvertToInt(row[4])
	return CereriStateData{
		Name:        row[1],
		RequestType: getRequestType(row[2]),
		Online:      online,
		Ghiseu:      ghiseu,
		Total:       online + ghiseu,
	}
}
