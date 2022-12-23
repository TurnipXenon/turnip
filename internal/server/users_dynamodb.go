// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const ddbTimeout = time.Second * 3

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

// GetUser may return nil
func (u *usersDynamoDBImpl) GetUser(ud *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	item, err := u.ddb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {S: aws.String(ud.Username)},
		},
		TableName: u.ddbTableName,
	})
	if err != nil {
		fmt.Printf("GetUser: Error: %s\n", err.Error())
		return nil, err
	}

	if item.Item == nil {
		return nil, nil
	}

	newUserData := User{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &newUserData)
	if err != nil {
		fmt.Printf("GetUser: Error: %s\n", err.Error())
		return nil, err
	}

	return &newUserData, nil
}

func (u *usersDynamoDBImpl) CreateUser(ud *User) error {
	// todo: add this pattern to all calls here???
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	// check if user already exists
	item, err := u.GetUser(ud)
	if item != nil {
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
