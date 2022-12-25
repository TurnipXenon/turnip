// storage is an abstraction to s3 buckets

package server

import (
	"context"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type Tokens interface {
	GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error)
	GetToken(token string) (*turnip.Token, error)
}
