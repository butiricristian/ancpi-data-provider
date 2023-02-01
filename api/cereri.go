package api

import (
	"encoding/json"
	"net/http"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterCereriByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.CereriStateData {
	cereriData := map[time.Time][]*models.CereriStateData{}
	for _, val := range data.Data {
		if val.CurrentDate.Before(dateStart) || val.CurrentDate.After(dateEnd) {
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

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	dateStart := helpers.ConvertToTime(r.URL.Query().Get("dateStart"))
	dateEnd := helpers.ConvertToTime(r.URL.Query().Get("dateEnd"))

	result := filterCereriByInterval(dateStart, dateEnd)
	result = filterCereriByJudet(result, judet)
	json.NewEncoder(w).Encode(result)
}
