package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockDynamoDBPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

func (m mockDynamoDBPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}

func TestSaveLessonInDynamoDB(t *testing.T) {
	cases := []struct {
		client          func(t *testing.T) DynamoDBPutItemAPI
		tableName       string
		lesson          Lesson
		fields_affected int
	}{
		{
			client: func(t *testing.T) DynamoDBPutItemAPI {
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
			},
			tableName: "lessons",
			lesson: Lesson{
				id:   1,
				name: "lesson 1",
			},
			fields_affected: 1,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()

			fields_affected, err := PutItemInDynamoDB(ctx, tt.client(t), tt.tableName, tt.lesson)

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := tt.fields_affected, fields_affected; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}

		})
	}
}
