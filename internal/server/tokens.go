// storage is an abstraction to s3 buckets

package server

import "github.com/TurnipXenon/Turnip/pkg/models"

type Tokens interface {
	GetOrCreateTokenByUsername(ud *User) (*models.Token, error)
	GetToken(token string) (*models.Token, error)
}
