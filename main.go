package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

//HandleRequest is responsible in receiving the content
func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var svc service
	svc.Initialize()

	switch req.HTTPMethod {
	case "PUT":
		return svc.CreateExpense(req.Body)
	case "DELETE":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		return svc.DeleteExpense(user, date)
	case "GET":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		return svc.GetExpenses(user, date)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}

}

func main() {
	lambda.Start(HandleRequest)
}
