// storage is an abstraction to s3 buckets

package server

import "github.com/TurnipXenon/Turnip/pkg/models"

type Tokens interface {
	GetOrCreateToken(ud *User) (*models.Token, error)
}
