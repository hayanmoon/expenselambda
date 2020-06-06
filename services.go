package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type service struct {
	client *dynamodb.DynamoDB
}

func (svc *service) Initialize() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		return errors.New("Could not initialize service")
	}

	svc.client = dynamodb.New(sess)

	return nil
}

func (svc service) CreateExpense(content string) error {
	var expense Expense
	var err error

	body := []byte(content)

	if err = json.Unmarshal(body, &expense); err != nil {
		return errors.New("Invalid Expense Definition")
	}

	currentTime := time.Now()
	expense.Date = currentTime.Format("20060102030405") //year month day hour minute second

	av, err := dynamodbattribute.MarshalMap(expense)

	if err != nil {
		return errors.New("Unable to marshal expense")
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("expense"),
		Item:      av,
	}

	_, err = svc.client.PutItem(input)

	if err != nil {
		return errors.New("Expense Not Saved")
	}
	return nil
}

func (svc service) GetExpenses(user, date string) ([]byte, error) {
	//form query to search
	input := dynamodb.QueryInput{
		TableName: aws.String("expense"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{ //create values to use for filter
			":user": {
				S: aws.String(user),
			},
			":dt": {
				S: aws.String(date),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#u":  aws.String("user"),
			"#dt": aws.String("date"),
		},
		KeyConditionExpression: aws.String("#u = :user and begins_with(#dt, :dt)"), //filter on keys
	}

	output, err := svc.client.Query(&input)
	if err != nil {
		return nil, errors.New("Could not create query")
	}

	var items []Expense
	dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)

	result, err := json.Marshal(items)
	if err != nil {
		return nil, errors.New("Unable to marshal query result")
	}

	return result, nil
}

func (svc service) DeleteExpense(user, date string) error {
	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("expense"),
		Key: map[string]*dynamodb.AttributeValue{
			"user": {
				S: aws.String(user),
			},
			"date": {
				S: aws.String(date),
			},
		},
	}

	_, err := svc.client.DeleteItem(deleteInput)

	if err != nil {
		return errors.New("Expense not deleted")
	}

	return nil
}
