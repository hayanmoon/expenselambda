package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type content struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func response(message string, ok bool) (events.APIGatewayProxyResponse, error) {
	data := content{
		Message: message,
		Ok:      ok,
	}
	parsed, err := json.Marshal(data)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Unable to parse response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(parsed),
	}, nil
}

//HandleRequest is responsible in receiving the content
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user, date string
	var svc service
	err := svc.Initialize()

	if err != nil {
		return response(err.Error(), false)
	}

	//get query string parameters
	fmt.Println(len(req.QueryStringParameters))
	if len(req.QueryStringParameters) > 0 {
		user = req.QueryStringParameters["user"]
		date = req.QueryStringParameters["date"]
	}

	switch req.HTTPMethod {
	case "PUT":
		err = svc.CreateExpense(req.Body)

		if err != nil {
			return response(err.Error(), false)
		}

		return response("Expense Created", true)
	case "DELETE":
		err := svc.DeleteExpense(user, date)

		if err != nil {
			return response(err.Error(), false)
		}

		return response("Expense Deleted", true)
	case "GET":
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
