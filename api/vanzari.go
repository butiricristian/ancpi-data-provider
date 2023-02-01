package api

import (
	"encoding/json"
	"net/http"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterVanzariByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.VanzariStateData {
	vanzariData := map[time.Time][]*models.VanzariStateData{}
	for _, val := range data.Data {
		if val.CurrentDate.Before(dateStart) || val.CurrentDate.After(dateEnd) {
			continue
		}
		vanzariData[val.CurrentDate] = val.VanzariData
	}
	return vanzariData
}

func filterVanzariByJudet(result map[time.Time][]*models.VanzariStateData, judet string) map[time.Time][]*models.VanzariStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.VanzariStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.VanzariStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func GetVanzariData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	dateStart := helpers.ConvertToTime(r.URL.Query().Get("dateStart"))
	dateEnd := helpers.ConvertToTime(r.URL.Query().Get("dateEnd"))

	result := filterVanzariByInterval(dateStart, dateEnd)
	result = filterVanzariByJudet(result, judet)
	json.NewEncoder(w).Encode(result)
}
