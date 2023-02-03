package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting Cereri Data")
	judet := r.URL.Query().Get("judet")
	requestType := models.GetRequestType(r.URL.Query().Get("requestType"))

	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString := r.URL.Query().Get("dateStart"); dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString := r.URL.Query().Get("dateEnd"); dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterCereriByInterval(dateStart, dateEnd)
	resultByJudet := filterCereriByJudet(resultByInterval, judet)
	result := filterCereriByRequestType(resultByJudet, requestType)

	json.NewEncoder(w).Encode(result)
}
