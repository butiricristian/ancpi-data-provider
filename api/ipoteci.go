package api

import (
	"encoding/json"
	"net/http"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterIpoteciByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.IpoteciStateData {
	ipoteciData := map[time.Time][]*models.IpoteciStateData{}
	for _, val := range data.Data {
		if val.CurrentDate.Before(dateStart) || val.CurrentDate.After(dateEnd) {
			continue
		}
		ipoteciData[val.CurrentDate] = val.IpoteciData
	}
	return ipoteciData
}

func filterIpoteciByJudet(result map[time.Time][]*models.IpoteciStateData, judet string) map[time.Time][]*models.IpoteciStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.IpoteciStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	dateStart := helpers.ConvertToTime(r.URL.Query().Get("dateStart"))
	dateEnd := helpers.ConvertToTime(r.URL.Query().Get("dateEnd"))

	result := filterIpoteciByInterval(dateStart, dateEnd)
	result = filterIpoteciByJudet(result, judet)
	json.NewEncoder(w).Encode(result)
}
