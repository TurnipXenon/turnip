// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"github.com/TurnipXenon/turnip/pkg/models"
)

type Tokens interface {
	GetOrCreateTokenByUsername(ctx context.Context, ud *User) (*models.Token, error)
	GetToken(token string) (*models.Token, error)
}
