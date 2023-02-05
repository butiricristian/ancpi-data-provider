package main

import (
	"encoding/json"
	"fmt"

	"com.butiricristian/ancpi-data-provider/controllers"
	"com.butiricristian/ancpi-data-provider/data"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	judet := request.QueryStringParameters["judet"]
	dateStartString := request.QueryStringParameters["dateStart"]
	dateEndString := request.QueryStringParameters["dateEnd"]

	result := controllers.HandleGetVanzariData(judet, dateStartString, dateEndString)
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(err.Error()),
			Headers:    map[string]string{"access-control-allow-origin": "*"},
		}, fmt.Errorf("error while marshalling ipoteci to JSON")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResult),
		Headers:    map[string]string{"access-control-allow-origin": "*"},
	}, nil
}

func main() {
	data.PrepareDataFromUrl()
	lambda.Start(handler)
}
