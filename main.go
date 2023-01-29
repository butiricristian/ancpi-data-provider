package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Printf("Hello World! %v\n", now)

	excelUrls := findAllExcelUrls()
	// for _, excelUrl := range excelUrls {
	// 	fmt.Printf("%s, %s %s: %s\n", excelUrl.month, excelUrl.year, excelUrl.name, excelUrl.url)
	// }

	data := getDataFromExcels(excelUrls)
	printData(data)
}
