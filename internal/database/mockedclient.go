package database

import (
	"context"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type mockDynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type MockDynamoDBClient struct {
	PutItemFunc func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItemFunc func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	data        map[string]map[string]map[string]types.AttributeValue
	mu          sync.Mutex
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {

	// m.mu.Lock()
	// defer m.mu.Unlock()

	// if params.TableName == nil {
	// 	return nil, fmt.Errorf("TableName is nil")
	// }
	// tableName := *params.TableName

	// if m.data[tableName] == nil {
	// 	m.data[tableName] = make(map[string]map[string]types.AttributeValue)
	// }

	// log.Print("id: ", params.Item["Id"])

	// id, ok := params.Item["Id"].(*types.AttributeValueMemberS)
	// if !ok {
	// 	return nil, fmt.Errorf("id attribute is not a string")
	// }

	// m.data[tableName][id.Value] = params.Item

	// return &dynamodb.PutItemOutput{}, nil
	return &dynamodb.PutItemOutput{}, nil
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	// if params.TableName == nil {
	// 	return nil, fmt.Errorf("TableName is nil")
	// }
	// tableName := *params.TableName

	// key, ok := params.Key["Id"].(*types.AttributeValueMemberS)
	// if !ok {
	// 	return nil, fmt.Errorf("Key 'Id' is not a string")
	// }

	// item, ok := m.data[tableName][key.Value]
	// if !ok {
	// 	return &dynamodb.GetItemOutput{}, nil
	// }

	// return &dynamodb.GetItemOutput{Item: item}, nil
	return &dynamodb.GetItemOutput{}, nil
}

func NewMockDynamoDBClient() *MockDynamoDBClient {
	return &MockDynamoDBClient{
		data: make(map[string]map[string]map[string]types.AttributeValue),
	}
}

func GetMockedClient(t *testing.T) *MockDynamoDBClient {
	client := NewMockDynamoDBClient()

	client.PutItemFunc = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
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
		return client.PutItem(ctx, params, optFns...)
	}

	client.GetItemFunc = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		t.Helper()
		if params.TableName == nil {
			t.Fatal("expect table name to not be nil")
		}
		if e, a := "lessons", *params.TableName; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if params.Key == nil {
			t.Fatal("expect key to not be nil")
		}
		return client.GetItem(ctx, params, optFns...)
	}

	return client
}
