// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type usersDynamoDBImpl struct {
	ddb          *dynamodb.DynamoDB
	ddbTableName *string
}

type CreateUserUserAlreadyExists struct{}

func (m *CreateUserUserAlreadyExists) Error() string {
	return "User already exists"
}

func NewUsersDynamoDB(d *dynamodb.DynamoDB) Users {
	s := usersDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Users"),
	}
	return &s
}

func (u *usersDynamoDBImpl) CreateUser(ud *UserData) error {
	// todo: add this pattern to all calls here???
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()

	// check if user already exists
	item, err := u.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {S: aws.String(ud.Username)},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		fmt.Printf("CreateUser: Error: %s\n", err.Error())
		return err
	}
	if item.Item != nil {
		fmt.Printf("User already exists: %s\n", ud.Username)
		return &CreateUserUserAlreadyExists{}
	}

	_, err = u.ddb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Username":       {S: aws.String(ud.Username)},
			"HashedPassword": {S: aws.String(ud.HashedPassword)},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		fmt.Printf("CreateUser: Error: %s\n", err.Error())
	}
	return err
}
