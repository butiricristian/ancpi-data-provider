package api

import (
	"encoding/json"
	"net/http"
	"os"

	"com.butiricristian/ancpi-data-provider/parserjob"
)

func processIpoteciData(data map[string][]parserjob.IpoteciStateData) map[string][]parserjob.IpoteciStateData {
	return data
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	// judet := r.URL.Query().Get("judet")
	// dateStart := r.URL.Query().Get("dataStart")
	// dateEnd := r.URL.Query().Get("dataEnd")

	fileData, err := os.ReadFile("data/ipoteci.json")
	if err != nil {
		return
	}

	var data map[string][]parserjob.IpoteciStateData
	json.Unmarshal(fileData, &data)
	result := processIpoteciData(data)

	json.NewEncoder(w).Encode(result)
}
