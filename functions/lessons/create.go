package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"lumbrera/internal/database"
	"lumbrera/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

func main() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	dynamoClient := dynamodb.NewFromConfig(sdkConfig)
	handler := &Handler{
		Client: dynamoClient,
	}

	lambda.Start(handler.Handle)
}

type Handler struct {
	Client database.DynamoDBPutItemAPI
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Only POST method is allowed",
		}, nil
	}

	var lesson models.Lesson

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

	database.PutItemInDynamoDB(ctx, h.Client, "lessons", lesson)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Lesson " + lesson.Name + " was created successfully",
	}, nil
}
