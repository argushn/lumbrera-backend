package database

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetLocalClient(endpoint string) (*dynamodb.Client, error) {

	log.Println("from GetLocalClient function: ", endpoint)
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {

		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				// PartitionID: "aws",
				URL: endpoint,
				// SigningRegion: "us-west-1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		fmt.Println("can't configure a local endpoint")
	}

	return dynamodb.NewFromConfig(sdkConfig), err
}
