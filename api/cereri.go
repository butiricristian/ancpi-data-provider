package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/models"
)

func processCereriData(data []*models.MonthlyData) []*models.CereriStateData {
	var cereriData []*models.CereriStateData
	for _, val := range data {
		cereriData = append(cereriData, val.CereriData...)
	}
	return cereriData
}

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/data.json")
	if err != nil {
		return
	}

	var data []*models.MonthlyData
	json.Unmarshal(fileData, &data)
	result := processCereriData(data)
	json.NewEncoder(w).Encode(result)
}
