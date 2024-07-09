package database

import (
	"context"
	"strconv"
	"testing"

	"lumbrera/internal/models"
)

func TestSaveLessonInDynamoDB(t *testing.T) {
	cases := []struct {
		description     string
		tableName       string
		lesson          models.Lesson
		fields_affected int
	}{
		{
			description: "Save lesson in dynamoDB",
			tableName:   "lessons",
			lesson: models.Lesson{
				Id:   "1",
				Name: "lesson 1",
			},
			fields_affected: 1,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i)+": "+tt.description, func(t *testing.T) {
			ctx := context.TODO()

			fields_affected, err := PutItemInDynamoDB(ctx, GetMockedClient(t), tt.tableName, tt.lesson)

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := tt.fields_affected, fields_affected; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestLocalDynamoDBConnection(t *testing.T) {
	ctx := context.TODO()

	lesson := models.Lesson{
		Id:   "1",
		Name: "lesson in local dynamodb",
	}

	dynamodbClient, err := GetLocalClient("http://localhost:8000")

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	fields_affected, err := PutItemInDynamoDB(ctx, dynamodbClient, "lessons", lesson)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 1, fields_affected; fields_affected != 1 {
		t.Errorf("expect %v, got %v", e, a)
	}
}
