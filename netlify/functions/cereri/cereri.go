package cereri

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

func filterCereriByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.CereriStateData {
	cereriData := map[time.Time][]*models.CereriStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		cereriData[val.CurrentDate] = val.CereriData
	}
	return cereriData
}

func filterCereriByJudet(result map[time.Time][]*models.CereriStateData, judet string) map[time.Time][]*models.CereriStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.CereriStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.CereriStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func filterCereriByRequestType(result map[time.Time][]*models.CereriStateData, requestType models.RequestType) map[time.Time]*models.CereriStateData {
	newResult := map[time.Time]*models.CereriStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.RequestType == requestType {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	judet := request.QueryStringParameters["judet"]
	requestTypeString := request.QueryStringParameters["requestType"]
	dateStartString := request.QueryStringParameters["dateStart"]
	dateEndString := request.QueryStringParameters["dateEnd"]

	result := handleGetCereriData(judet, requestTypeString, dateStartString, dateEndString)
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

func handleGetCereriData(judet string, requestTypeString string, dateStartString string, dateEndString string) map[time.Time]*models.CereriStateData {
	fmt.Println("Getting Cereri Data")
	requestType := models.GetRequestType(requestTypeString)

	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterCereriByInterval(dateStart, dateEnd)
	resultByJudet := filterCereriByJudet(resultByInterval, judet)
	result := filterCereriByRequestType(resultByJudet, requestType)
	return result
}

func GetCereriData(w http.ResponseWriter, r *http.Request) {
	requestType := r.URL.Query().Get("requestType")
	judet := r.URL.Query().Get("judet")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := handleGetCereriData(judet, requestType, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func main() {
	go data.PrepareData()
	lambda.Start(handler)
}
