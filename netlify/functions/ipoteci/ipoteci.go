package ipoteci

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"com.butiricristian/ancpi-data-provider/data"
	"com.butiricristian/ancpi-data-provider/helpers"
	"com.butiricristian/ancpi-data-provider/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func filterIpoteciByInterval(dateStart time.Time, dateEnd time.Time) map[time.Time][]*models.IpoteciStateData {
	ipoteciData := map[time.Time][]*models.IpoteciStateData{}
	for _, val := range data.Data {
		if !dateStart.IsZero() && val.CurrentDate.Before(dateStart) {
			continue
		}
		if !dateEnd.IsZero() && val.CurrentDate.After(dateEnd) {
			continue
		}
		ipoteciData[val.CurrentDate] = val.IpoteciData
	}
	return ipoteciData
}

func filterIpoteciByJudet(result map[time.Time][]*models.IpoteciStateData, judet string) map[time.Time][]*models.IpoteciStateData {
	if judet == "" {
		judet = "TOTAL"
	}
	newResult := map[time.Time][]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Name == judet {
				if _, ok := newResult[key]; !ok {
					newResult[key] = []*models.IpoteciStateData{}
				}
				newResult[key] = append(newResult[key], val)
			}
		}
	}
	return newResult
}

func filterIpoteciByActive(result map[time.Time][]*models.IpoteciStateData, active bool) map[time.Time]*models.IpoteciStateData {
	newResult := map[time.Time]*models.IpoteciStateData{}
	for key, data := range result {
		for _, val := range data {
			if val.Active == active {
				newResult[key] = val
			}
		}
	}
	return newResult
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	judet := request.QueryStringParameters["judet"]
	active := request.QueryStringParameters["ipoteciActive"]
	dateStartString := request.QueryStringParameters["dateStart"]
	dateEndString := request.QueryStringParameters["dateEnd"]

	result := handleGetIpoteciData(judet, active, dateStartString, dateEndString)
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

func handleGetIpoteciData(judet string, activeString string, dateStartString string, dateEndString string) map[time.Time]*models.IpoteciStateData {
	fmt.Println("Getting Ipoteci Data")
	active, err := strconv.ParseBool(activeString)
	if err != nil {
		active = true
	}
	var dateStart time.Time
	var dateEnd time.Time
	if dateStartString != "" {
		dateStart = helpers.ConvertToTime(dateStartString)
	}
	if dateEndString != "" {
		dateEnd = helpers.ConvertToTime(dateEndString)
	}

	resultByInterval := filterIpoteciByInterval(dateStart, dateEnd)
	resultByJudet := filterIpoteciByJudet(resultByInterval, judet)
	result := filterIpoteciByActive(resultByJudet, active)
	return result
}

func GetIpoteciData(w http.ResponseWriter, r *http.Request) {
	judet := r.URL.Query().Get("judet")
	active := r.URL.Query().Get("ipoteciActive")
	dateStartString := r.URL.Query().Get("dateStart")
	dateEndString := r.URL.Query().Get("dateEnd")

	result := handleGetIpoteciData(judet, active, dateStartString, dateEndString)
	json.NewEncoder(w).Encode(result)
}

func main() {
	go data.PrepareData()
	lambda.Start(handler)
}
