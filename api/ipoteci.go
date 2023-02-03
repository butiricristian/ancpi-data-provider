package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
)

func filterIpoteciByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.IpoteciStateData {
	ipoteciData := map[time.Time][]*models.IpoteciStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
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

func filterIpoteciByActive(result map[time.Time][]*models.IpoteciStateData, active bool) map[time.Time]*models.IpoteciStateData {
	newResult := map[time.Time]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Active == active {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting Ipoteci Data")
	judet := r.URL.Query().Get("judet")
	active, err := strconv.ParseBool(r.URL.Query().Get("active"))
	if err != nil {
		active = true
	}
	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString := r.URL.Query().Get("dateStart"); dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString := r.URL.Query().Get("dateEnd"); dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterIpoteciByInterval(dateStart, dateEnd)
	resultByJudet := filterIpoteciByJudet(resultByInterval, judet)
	result := filterIpoteciByActive(resultByJudet, active)
	json.NewEncoder(w).Encode(result)
}
