package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"lumbrera/internal/database"
	"lumbrera/internal/models"
)

type mockDynamoDBPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

func (m mockDynamoDBPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}

func getMockedClient(t *testing.T) database.DynamoDBPutItemAPI {
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

func TestSaveLessonInDynamoDB(t *testing.T) {
	cases := []struct {
		tableName       string
		lesson          models.Lesson
		fields_affected int
	}{
		{
			tableName: "lessons",
			lesson: models.Lesson{
				Id:   "1",
				Name: "lesson 1",
			},
			fields_affected: 1,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()

			fields_affected, err := database.PutItemInDynamoDB(ctx, getMockedClient(t), tt.tableName, tt.lesson)

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := tt.fields_affected, fields_affected; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestSaveLesson(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid request",
			body:       `{"name":"lesson 1"}`,
			wantStatus: 200,
			wantBody:   "Lesson lesson 1 was created successfully",
		},
		{
			name:       "Missing name",
			body:       `{"name":""}`,
			wantStatus: 400,
			wantBody:   "Missing required field: name",
		},
		{
			name:       "Invalid JSON",
			body:       `{`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       tt.body,
			}

			handler := &Handler{Client: getMockedClient(t)}

			resp, err := handler.Handle(context.Background(), req)

			if err != nil {
				t.Fatal(err)
			}

			if e, a := tt.wantStatus, resp.StatusCode; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}

			if e, a := tt.wantBody, resp.Body; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
