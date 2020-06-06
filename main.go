package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is the object that is returned
type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func response(message string, status bool) (Response, error) {
	return Response{
		Message: message,
		Ok:      status,
	}, nil
}

//HandleRequest is responsible in receiving the content
func HandleRequest(req events.APIGatewayProxyRequest) (Response, error) {

	var svc service
	err := svc.Initialize()

	if err != nil {
		return response(err.Error(), false)
	}

	switch req.HTTPMethod {
	case "PUT":
		err = svc.CreateExpense(req.Body)

		if err != nil {
			return response(err.Error(), false)
		}

		return response("Expense Created", true)
	case "DELETE":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		err := svc.DeleteExpense(user, date)

		if err != nil {
			return response(err.Error(), false)
		}

		return response("Expense Deleted", true)
	case "GET":
		user := req.QueryStringParameters["user"]
		date := req.QueryStringParameters["date"]
		result, err := svc.GetExpenses(user, date)

		if err != nil {
			return response(err.Error(), false)
		}

		return response(string(result), true)
	default:
		return response(http.StatusText(http.StatusMethodNotAllowed), false)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
