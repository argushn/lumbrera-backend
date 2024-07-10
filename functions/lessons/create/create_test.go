package main

import (
	"context"
	"testing"

	"lumbrera/internal/database"

	"github.com/aws/aws-lambda-go/events"
)

func TestSaveLesson(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
		method     string
	}{
		{
			name:       "Valid request",
			body:       `{"name":"lesson 1"}`,
			wantStatus: 200,
			wantBody:   "Lesson lesson 1 was created successfully",
			method:     "POST",
		},
		{
			name:       "Wrong method",
			body:       `{"name":"lesson 2"}`,
			wantStatus: 405,
			wantBody:   "Only POST method is allowed",
			method:     "GET",
		},
		{
			name:       "Missing name",
			body:       `{"name":""}`,
			wantStatus: 400,
			wantBody:   "Missing required field: name",
			method:     "POST",
		},
		{
			name:       "Invalid JSON",
			body:       `{`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
			method:     "POST",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
				HTTPMethod: tt.method,
				Body:       tt.body,
			}

			handler := &Handler{Client: database.GetMockedClient(t)}

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
