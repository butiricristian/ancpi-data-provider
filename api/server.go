package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"com.butiricristian/ancpi-data-provider/controllers"
	"github.com/rs/cors"
)

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	requestType := r.URL.Query().Get("requestType")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := controllers.HandleGetCereriData(judet, requestType, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	active := r.URL.Query().Get("ipoteciActive")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := controllers.HandleGetIpoteciData(judet, active, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func GetVanzariData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := controllers.HandleGetVanzariData(judet, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func StartServer() {
	mux := &http.ServeMux{}

	mux.HandleFunc("/cereri", GetCereriData)
	mux.HandleFunc("/ipoteci", GetIpoteciData)
	mux.HandleFunc("/vanzari", GetVanzariData)

	fmt.Println("Listening on port 8080...")
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3002", "https://www.cristianbutiri.com"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	handler := cors.Handler(mux)
	http.ListenAndServe(":8080", handler)
}
