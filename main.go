package main

import (
	"fmt"
	"time"

	"com.butiricristian/ancpi-data-provider/api"
	"com.butiricristian/ancpi-data-provider/data"
)

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	go data.PrepareData("data/data.json")
	api.StartServer()
}
