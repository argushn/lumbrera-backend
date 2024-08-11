package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"lumbrera/internal/database"
	"lumbrera/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
	// if req.HTTPMethod != http.MethodGet {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusMethodNotAllowed,
	// 		Body:       "Only GET method is allowed",
	// 	}, nil
	// }

	var lesson models.Lesson

	// if err := json.Unmarshal([]byte(req.Body), &lesson); err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Body:       "Invalid request payload",
	// 	}, nil
	// }

	// if lesson.Name == "" {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Body:       "Missing required field: name",
	// 	}, nil
	// }

	lesson, err := database.GetItemFromDynamoDB(ctx, h.Client, "lessons", "1")

	if err != nil {
		// 	return events.APIGatewayProxyResponse{
		// 		StatusCode: http.StatusBadRequest,
		// 		Body:       "An error happened",
		// 	}, err
	}

	jsonData, err := json.Marshal(lesson)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonData),
	}, nil
}

// refactor
