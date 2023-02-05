package main

import (
	"com.butiricristian/ancpi-data-provider/data"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Data preparation successful",
		Headers:    map[string]string{"access-control-allow-origin": "*"},
	}, nil
}

func main() {
	data.PrepareData("../data.json")
	lambda.Start(handler)
}
