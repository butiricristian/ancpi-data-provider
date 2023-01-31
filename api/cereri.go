package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/parserjob"
)

func processCereriData(data map[string][]parserjob.CereriStateData) map[string][]parserjob.CereriStateData {
	return data
}

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/cereri.json")
	if err != nil {
		return
	}

	var data map[string][]parserjob.CereriStateData
	json.Unmarshal(fileData, &data)
	result := processCereriData(data)
	json.NewEncoder(w).Encode(result)
}
