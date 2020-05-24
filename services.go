package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type service struct {
	client *dynamodb.DynamoDB
}

func (svc *service) Initialize() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		fmt.Println(err)
	}

	svc.client = dynamodb.New(sess)
}

func (svc service) CreateExpense(content string) (events.APIGatewayProxyResponse, error) {
	var expense Expense

	body := []byte(content)

	if err := json.Unmarshal(body, &expense); err != nil {
		fmt.Println("Invalid Expense Definition")
	}

	currentTime := time.Now()
	expense.Expenseid = currentTime.Format("20060102HHmm") //year month day hour minute
	expense.Timestamp = currentTime.Unix()

	av, err := dynamodbattribute.MarshalMap(expense)

	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("expense"),
		Item:      av,
	}

	_, err = svc.client.PutItem(input)

	if err != nil {
		fmt.Println(err, "Expense was not created")
	}
	return events.APIGatewayProxyResponse{Body: string("Expense Created"), StatusCode: 200}, nil
}

func (svc service) GetExpenses(id string) (events.APIGatewayProxyResponse, error) {
	//form query to search
	input := dynamodb.QueryInput{
		TableName: aws.String("expense"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
		},
		KeyConditionExpression: aws.String("expenseid = :id"),
	}

	output, err := svc.client.Query(&input)
	if err != nil {
		fmt.Println(err)
	}

	var items []Expense
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	result, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(result),
		StatusCode: http.StatusOK,
	}, nil
}

func (svc service) DeleteExpense(expenseid, timestamp string) (events.APIGatewayProxyResponse, error) {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("expense"),
		Key: map[string]*dynamodb.AttributeValue{
			"expenseid": {
				S: aws.String(expenseid),
			},
			"timestamp": {
				N: aws.String(timestamp),
			},
		},
	}

	_, err := svc.client.DeleteItem(deleteInput)
	if err != nil {
		fmt.Println(err)
	}
	return events.APIGatewayProxyResponse{Body: string("Expense Deleted"), StatusCode: 200}, nil
}
