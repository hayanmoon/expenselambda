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
		expenseid := req.QueryStringParameters["expenseid"]
		timestamp := req.QueryStringParameters["timestamp"]
		return svc.DeleteExpense(expenseid, timestamp)
	case "GET":
		id := req.QueryStringParameters["expenseid"]
		return svc.GetExpenses(id)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}

}

func main() {
	lambda.Start(HandleRequest)
}
