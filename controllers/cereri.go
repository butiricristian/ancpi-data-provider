package controllers

import (
	"fmt"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterCereriByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.CereriStateData {
	cereriData := map[time.Time][]*models.CereriStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		cereriData[val.CurrentDate] = val.CereriData
	}
	return cereriData
}

func filterCereriByJudet(result map[time.Time][]*models.CereriStateData, judet string) map[time.Time][]*models.CereriStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.CereriStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.CereriStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func filterCereriByRequestType(result map[time.Time][]*models.CereriStateData, requestType models.RequestType) map[time.Time]*models.CereriStateData {
	if requestType == models.UNDEFINED {
		requestType = models.TOTAL
	}
	newResult := map[time.Time]*models.CereriStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.RequestType == requestType {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func HandleGetCereriData(judet string, requestTypeString string, dateStartString string, dateEndString string) map[time.Time]*models.CereriStateData {
	fmt.Println("Getting Cereri Data")
	requestType := models.GetRequestType(requestTypeString)

	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterCereriByInterval(dateStart, dateEnd)
	resultByJudet := filterCereriByJudet(resultByInterval, judet)
	result := filterCereriByRequestType(resultByJudet, requestType)
	return result
}
