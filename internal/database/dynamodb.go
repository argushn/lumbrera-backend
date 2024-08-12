package database

import (
	"context"
	"fmt"

	"lumbrera/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

func PutItemInDynamoDB(ctx context.Context, api mockDynamoDBAPI, tableName string, lesson models.Lesson) (int, error) {
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

func GetItemFromDynamoDB(ctx context.Context, api mockDynamoDBAPI, tableName string, lessonID string) (models.Lesson, error) {
	var lesson models.Lesson

	getItemInput, err := api.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: lessonID},
		},
		TableName: &tableName,
	})

	if err != nil {
		return models.Lesson{}, fmt.Errorf("cannot get item: %w", err)
	}

	err = attributevalue.UnmarshalMap(getItemInput.Item, &lesson)

	if err != nil {
		return lesson, fmt.Errorf("failed to unmarshal item to lesson: %w", err)
	}

	return lesson, nil
}

func UpdateItemInDynamoDB(ctx context.Context, api mockDynamoDBAPI, tableName string, lesson models.Lesson) (int, error) {
	//

	api.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: lesson.Id},
		},
		TableName: &tableName,
	})

	return 1, nil
}
