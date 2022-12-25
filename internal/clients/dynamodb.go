package clients

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDB() *dynamodb.Client {
	// todo: pass context background over here
	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8200"}, nil
			},
		)),
	)
	if err != nil {
		log.Fatalf("NewDynamoDB: error: %s", err.Error())
	}

	// Create DynamoDB client
	return dynamodb.NewFromConfig(config)
}
