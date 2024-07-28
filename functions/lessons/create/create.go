package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"lumbrera/internal/database"
	"lumbrera/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

func main() {
	dynamodbEndpoint := os.Getenv("DYNAMODB_ENDPOINT")
	var dynamoClient *dynamodb.Client
	var err error

	log.Println("DYNAMODB_ENDPOINT:", dynamodbEndpoint)
	log.Println("AWS_DEFAULT_REGION:", os.Getenv("AWS_DEFAULT_REGION"))

	if dynamodbEndpoint != "" {
		log.Println("Using custom endpoint:", dynamodbEndpoint)
		dynamoClient, err = database.GetLocalClient(dynamodbEndpoint)

		if err != nil {
			log.Println("Couldn't load configuration with custom endpoint.", err)
			return
		}
	} else {
		sdkConfig, err := config.LoadDefaultConfig(context.TODO())

		if err != nil {
			log.Println("Couldn't load default configuration. Have you set up your AWS account?", err)
			return
		}

		dynamoClient = dynamodb.NewFromConfig(sdkConfig)
	}

	handler := &Handler{
		Client: dynamoClient,
	}

	lambda.Start(handler.Handle)
}

type Handler struct {
	Client database.DynamoDBAPI
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

	fields_affected, err := database.PutItemInDynamoDB(ctx, h.Client, "lessons", lesson)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "An error happened",
		}, err
	}

	log.Println("Creating lesson: " + lesson.Name)
	log.Println("Affected fields: " + strconv.Itoa(fields_affected))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Lesson " + lesson.Name + " was created successfully",
	}, nil
}
