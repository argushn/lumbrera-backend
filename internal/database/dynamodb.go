package database

import (
	"context"
	"fmt"
	"log"

	"lumbrera/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func PutItemInDynamoDB(ctx context.Context, api DynamoDBPutItemAPI, tableName string, lesson models.Lesson) (int, error) {
	item, err := attributevalue.MarshalMap(&lesson)

	if err != nil {
		return 0, fmt.Errorf("unable to marshal product: %w", err)
	}

	log.Println("Putting lesson in dynamo DB")
	log.Println(lesson)

	_, err = api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	if err != nil {
		return 0, fmt.Errorf("cannot put item: %w", err)
	}

	return 1, nil
}
