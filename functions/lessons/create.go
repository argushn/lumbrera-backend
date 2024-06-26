package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"lumbrera/internal/database"
	"lumbrera/internal/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/google/uuid"
)

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
