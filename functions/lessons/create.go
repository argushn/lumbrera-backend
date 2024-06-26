package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/google/uuid"
)

type Lesson struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
}

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	lambdaClient := lambda.NewFromConfig(sdkConfig)

	maxItems := 10
	fmt.Printf("Let's list up to %v functions for your account.\n", maxItems)
	result, err := lambdaClient.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(int32(maxItems)),
	})
	if err != nil {
		fmt.Printf("Couldn't list functions for your account. Here's why: %v\n", err)
		return
	}
	if len(result.Functions) == 0 {
		fmt.Println("You don't have any functions!")
	} else {
		for _, function := range result.Functions {
			fmt.Printf("\t%v\n", *function.FunctionName)
		}
	}
}

type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func PutItemInDynamoDB(ctx context.Context, api DynamoDBPutItemAPI, tableName string, lesson Lesson) (int, error) {
	item, err := attributevalue.MarshalMap(&lesson)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal product: %w", err)
	}

	_, err = api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	if err != nil {
		return 0, fmt.Errorf("cannot put item: %w", err)
	}

	return 1, nil
}

type CreateLessonHandler struct {
	Client DynamoDBPutItemAPI
}

func NewCreateLessonHandler(client DynamoDBPutItemAPI) *CreateLessonHandler {
	return &CreateLessonHandler{
		Client: client,
	}
}

func (h *CreateLessonHandler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Only POST method is allowed",
		}, nil
	}

	var lesson Lesson

	if err := json.Unmarshal([]byte(req.Body), &lesson); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request payload",
		}, nil
	}
	if lesson.Name == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing required field: name",
		}, nil
	}

	lesson.Id = uuid.New().String()

	PutItemInDynamoDB(ctx, h.Client, "lessons", lesson)

	respBody, err := json.Marshal(lesson)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to marshal response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
	}, nil
}
