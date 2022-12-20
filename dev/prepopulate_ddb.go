// prepopulate_ddb.go for sample data
// todo: WIP

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type RandomItem struct {
	SuffixList map[string]int
}

// todo: hit 8200 port!!!
func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	table := "Config"

	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"HashKey": {S: aws.String("00000")},
		},
	}
	resp, err := svc.GetItem(input)

	if err != nil {
		print("Sad :(", err.Error())
		return
	}

	if resp.Item == nil {
		print("Nil", err.Error())
		return
	}

	item := RandomItem{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &item)
	if err != nil {
		print("Cow", err.Error())
	}
}
