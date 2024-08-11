package main

import (
	"context"
	"testing"

	"lumbrera/internal/database"
	"lumbrera/internal/models"

	"github.com/aws/aws-lambda-go/events"
)

func TestGetLesson(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		lesson     models.Lesson
		wantBody   string
		wantStatus int
		method     string
	}{
		{
			name: "Valid request",
			path: `/lessons`,
			lesson: models.Lesson{
				Id:   "1",
				Name: "lesson 1",
			},
			wantBody:   `{"Id":"1","Name":"lesson 1"}`,
			method:     "GET",
			wantStatus: 200,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
				HTTPMethod: tt.method,
				Path:       tt.path,
			}

			mockedClient := database.GetMockedClient(t)

			database.PutItemInDynamoDB(context.Background(), mockedClient, "lessons", tt.lesson)

			handler := &Handler{Client: mockedClient}

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
