package api

import (
	"fmt"
	"net/http"
)

func StartServer() {
	srv := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/cereri", GetCereriData)
	http.HandleFunc("/ipoteci", GetIpoteciData)
	http.HandleFunc("/vanzari", GetVanzariData)

	fmt.Println("Listening on port 8080...")
	srv.ListenAndServe()
}
