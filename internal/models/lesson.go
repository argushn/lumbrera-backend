package models

type Lesson struct {
	Id   string `dynamodbav:"ID,omitempty"`
	Name string `dynamodbav:"Name"`
}
