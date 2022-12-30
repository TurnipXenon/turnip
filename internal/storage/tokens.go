// storage is an abstraction to s3 buckets

package storage

import (
	"context"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

type Tokens interface {
	GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*turnip.Token, error)
	GetToken(ctx context.Context, accessToken string) (*turnip.Token, error)
}
