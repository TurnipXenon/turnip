// storage is an abstraction to s3 buckets

package server

import (
	"context"

	"github.com/TurnipXenon/turnip/internal/clients"
	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type tokensPostgresImpl struct {
	ddb          *clients.PostgresDb
	ddbTableName string
	// todo: global secondary index
}

func (t tokensPostgresImpl) GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error) {
	//TODO implement me
	panic("implement me")
}

func (t tokensPostgresImpl) GetToken(token string) (*turnip.Token, error) {
	//TODO implement me
	panic("implement me")
}

func NewTokensPostgres(d *clients.PostgresDb) Tokens {
	t := tokensPostgresImpl{
		ddb:          d,
		ddbTableName: "Tokens",
	}
	return &t
}
