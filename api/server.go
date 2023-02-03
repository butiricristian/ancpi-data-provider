package api

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

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
