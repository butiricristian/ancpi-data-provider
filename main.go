package main

import (
	"fmt"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
)

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	data.PrepareData("data/data.json")
}
