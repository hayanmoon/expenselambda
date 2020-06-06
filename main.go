package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is the object that is returned
type Response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

func response(message string) (Response, error) {
	return Response{
		Body:       message,
		StatusCode: 200,
	}, nil
}

//HandleRequest is responsible in receiving the content
func HandleRequest(req events.APIGatewayProxyRequest) (Response, error) {

	var svc service
	err := svc.Initialize()

	if err != nil {
		return response(err.Error())
	}

	switch req.HTTPMethod {
	case "PUT":
		err = svc.CreateExpense(req.Body)

		if err != nil {
			return response(err.Error())
		}

		return response("Expense Created")
	case "DELETE":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		err := svc.DeleteExpense(user, date)

		if err != nil {
			return response(err.Error())
		}

		return response("Expense Deleted")
	case "GET":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		result, err := svc.GetExpenses(user, date)

		if err != nil {
			return response(err.Error())
		}

		return response(string(result))
	default:
		return response(http.StatusText(http.StatusMethodNotAllowed))
	}
}

func main() {
	lambda.Start(HandleRequest)
}
