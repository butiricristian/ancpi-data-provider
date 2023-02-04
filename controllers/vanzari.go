package controllers

import (
	"fmt"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterVanzariByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.VanzariStateData {
	vanzariData := map[time.Time][]*models.VanzariStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		vanzariData[val.CurrentDate] = val.VanzariData
	}
	return vanzariData
}

func filterVanzariByJudet(result map[time.Time][]*models.VanzariStateData, judet string) map[time.Time]*models.VanzariStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time]*models.VanzariStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func HandleGetVanzariData(judet string, dateStartString string, dateEndString string) map[time.Time]*models.VanzariStateData {
	fmt.Println("Getting Vanzari Data")
	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterVanzariByInterval(dateStart, dateEnd)
	result := filterVanzariByJudet(resultByInterval, judet)
	return result
}
