package main

import (
	"fmt"
	"io"
	"net/http"
)

func requestPage(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("An error occured while retrieving the page: %v", err)
	}

	return resp.Body
}
