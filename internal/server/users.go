// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"
)

var (
	UserAlreadyExists = errors.New("migration already exists")
)

type User struct {
	Username        string
	HashedPassword  string
	AccessGroupList []string
}

type Users interface {
	CreateUser(ctx context.Context, ud *User) error
	GetUser(ctx context.Context, ud *User) (*User, error)
}

func FromUserRequestToUserData(from *turnip.CreateUserRequest) (User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(from.Password), 14)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:        from.Username,
		HashedPassword:  string(password),
		AccessGroupList: []string{},
	}, nil
}
