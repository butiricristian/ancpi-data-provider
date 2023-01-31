package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/models"
)

func processIpoteciData(data []*models.MonthlyData) []*models.IpoteciStateData {
	var ipoteciData []*models.IpoteciStateData
	for _, val := range data {
		ipoteciData = append(ipoteciData, val.IpoteciData...)
	}
	return ipoteciData
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/data.json")
	if err != nil {
		return
	}

	var data []*models.MonthlyData
	json.Unmarshal(fileData, &data)
	result := processIpoteciData(data)

	json.NewEncoder(w).Encode(result)
}
