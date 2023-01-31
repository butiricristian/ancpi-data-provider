package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/models"
)

func processVanzariData(data []*models.MonthlyData) []*models.VanzariStateData {
	var vanzariData []*models.VanzariStateData
	for _, val := range data {
		vanzariData = append(vanzariData, val.VanzariData...)
	}
	return vanzariData
}

func GetVanzariData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/data.json")
	if err != nil {
		return
	}

	var data []*models.MonthlyData
	json.Unmarshal(fileData, &data)
	result := processVanzariData(data)

	json.NewEncoder(w).Encode(result)
}
