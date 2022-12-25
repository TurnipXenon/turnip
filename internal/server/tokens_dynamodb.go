// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/TurnipXenon/turnip_api/rpc/turnip"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/TurnipXenon/turnip/internal/util"
)

type tokensDynamoDBImpl struct {
	ddb          *dynamodb.Client
	ddbTableName *string
	// todo: global secondary index
}

func NewTokensDynamoDB(d *dynamodb.Client) Tokens {
	t := tokensDynamoDBImpl{
		ddb:          d,
		ddbTableName: aws.String("Tokens"),
	}
	return &t
}

func (t *tokensDynamoDBImpl) GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error) {
	token := turnip.Token{}

	// from https://github.com/awsdocs/aws-doc-sdk-examples/blob/c3a3dbe1d420b0b75f3e8976e12ee3c96fbd1527/gov2/dynamodb/actions/table_basics.go#L247
	keyEx := expression.Key("Username").Equal(expression.Value(ud.Username))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}

	// (1) if token exists
	query, err := t.ddb.Query(ctx, &dynamodb.QueryInput{
		IndexName: aws.String("UsernameIndex"),
		TableName: t.ddbTableName,

		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	if len(query.Items) > 0 {
		// it exists!
		err = attributevalue.UnmarshalMap(query.Items[0], &token)
		if err != nil {
			util.LogDetailedError(err)
			return nil, err
		}

		return &token, nil
	}

	// (2) if token does not exist
	token.Username = ud.Username
	token.AccessToken, err = generateSecureToken(40)
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}

	//dt := time.Now()
	//token.GeneratedAt = dt.Format(time.RFC3339)
	//expiryTime := time.Now().Add(time.Hour * 24)
	//token.ExpiresAt = expiryTime.Format(time.RFC3339)

	_, err = t.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"Username":    &types.AttributeValueMemberS{Value: token.Username},
			"AccessToken": &types.AttributeValueMemberS{Value: token.AccessToken}, // todo: we could support expirable tokens later...
		},
		TableName: t.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	return &token, err
}

func generateSecureToken(length int) (string, error) {
	// from Andzej Maciusovic @ https://stackoverflow.com/a/59457748/17836168
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		fmt.Printf("generateSecureToken: Error: %s\n", err.Error())
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (t *tokensDynamoDBImpl) GetToken(accessToken string) (*turnip.Token, error) {
	// todo: how do we make this inline with twirp generated server stubs?
	ctx, cancel := context.WithTimeout(context.TODO(), ddbTimeout)
	defer cancel()

	item, err := t.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"TokenAccess": &types.AttributeValueMemberS{Value: accessToken},
		},
		TableName: t.ddbTableName,
	})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err
	}
	if item.Item != nil {
		token := turnip.Token{}
		err = attributevalue.UnmarshalMap(item.Item, &token)
		if err != nil {
			util.LogDetailedError(err)
			return nil, err
		}

		return &token, nil
	}

	return nil, nil
}
