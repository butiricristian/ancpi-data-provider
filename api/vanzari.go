package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/parserjob"
)

func processVanzariData(data map[string][]parserjob.VanzariStateData) map[string][]parserjob.VanzariStateData {
	return data
}

func GetVanzariData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/vanzari.json")
	if err != nil {
		return
	}

	var data map[string][]parserjob.VanzariStateData
	json.Unmarshal(fileData, &data)
	result := processVanzariData(data)

	json.NewEncoder(w).Encode(result)
}
