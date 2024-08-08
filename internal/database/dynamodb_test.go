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

func TestReadLessonsFromDynamoDB(t *testing.T) {
	cases := []struct {
		description     string
		tableName       string
		initialLessons  []models.Lesson
		expectedLessons []models.Lesson
	}{
		{
			description: "Read lessons from Lessons table in dynamodb",
			tableName:   "lessons",
			initialLessons: []models.Lesson{
				{Id: "1", Name: "lesson 1"},
			},
			expectedLessons: []models.Lesson{
				{Id: "1", Name: "lesson 1"},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i)+":"+tt.description, func(t *testing.T) {
			ctx := context.TODO()

			lessons := []models.Lesson{}

			mockedClient := GetMockedClient(t)

			for _, lesson := range tt.initialLessons {
				PutItemInDynamoDB(ctx, mockedClient, tt.tableName, lesson)
			}

			for _, lesson := range tt.expectedLessons {
				readLesson, err := GetItemFromDynamoDB(ctx, mockedClient, tt.tableName, lesson.Id)

				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}

				lessons = append(lessons, readLesson)
			}

			if e, a := tt.expectedLessons[0].Id, lessons[0].Id; e != a {
				t.Errorf("expect %v to be similar to %v", tt.expectedLessons, lessons)
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
