package database

import (
	"context"
	"fmt"
	"log"

	"lumbrera/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func PutItemInDynamoDB(ctx context.Context, api mockDynamoDBAPI, tableName string, lesson models.Lesson) (int, error) {
	item, err := attributevalue.MarshalMap(&lesson)

	log.Print("from dynamodb.go/PutItemInDynamoDB", item)

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

func GetItemFromDynamoDB(ctx context.Context, api mockDynamoDBAPI, tableName string, lessonID string) (models.Lesson, error) {
	var lesson models.Lesson

	getItemInput, err := api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: lessonID},
		},
	})

	if err != nil {
		return models.Lesson{}, fmt.Errorf("cannot get item: %w", err)
	}

	err = attributevalue.UnmarshalMap(getItemInput.Item, &lesson)

	if err != nil {
		return lesson, fmt.Errorf("failed to unmarshal item to lesson: %w", err)

	}

	lesson.Id = "1"
	lesson.Name = "Lesson 1"

	return lesson, nil
}
