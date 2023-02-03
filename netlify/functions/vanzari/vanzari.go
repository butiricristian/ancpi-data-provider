package vanzari

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func filterVanzariByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.VanzariStateData {
	vanzariData := map[time.Time][]*models.VanzariStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		vanzariData[val.CurrentDate] = val.VanzariData
	}
	return vanzariData
}

func filterVanzariByJudet(result map[time.Time][]*models.VanzariStateData, judet string) map[time.Time]*models.VanzariStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time]*models.VanzariStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	judet := request.QueryStringParameters["judet"]
	dateStartString := request.QueryStringParameters["dateStart"]
	dateEndString := request.QueryStringParameters["dateEnd"]

	result := handleGetVanzariData(judet, dateStartString, dateEndString)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(jsonResult),
		}, fmt.Errorf("error while marshalling ipoteci to JSON")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResult),
	}, nil
}

func handleGetVanzariData(judet string, dateStartString string, dateEndString string) map[time.Time]*models.VanzariStateData {
	fmt.Println("Getting Vanzari Data")
	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterVanzariByInterval(dateStart, dateEnd)
	result := filterVanzariByJudet(resultByInterval, judet)
	return result
}

func GetVanzariData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := handleGetVanzariData(judet, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func main() {
	go data.PrepareData()
	lambda.Start(handler)
}
