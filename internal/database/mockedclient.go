package database

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockDynamoDBPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

func (m mockDynamoDBPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}

func GetMockedClient(t *testing.T) DynamoDBPutItemAPI {
	return mockDynamoDBPutItemAPI(func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		t.Helper()
		if params.TableName == nil {
			t.Fatal("expect table name to not be nil")
		}
		if e, a := "lessons", *params.TableName; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if params.Item == nil {
			t.Fatal("expect item to not be nil")
		}

		return &dynamodb.PutItemOutput{}, nil
	})
}
