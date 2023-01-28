package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	url := findByElemAttr("a", "title", "Decembrie 2022_Vanzari")
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("An error occured while retrieving the xls file: %v", err)
	}
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An error occured while reading the xls file: %v", err)
	}
	fmt.Println(string(s))
}
