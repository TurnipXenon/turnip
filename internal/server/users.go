// storage is an abstraction to s3 buckets

package server

import (
	"context"
	"github.com/TurnipXenon/turnip_twirp/rpc/turnip"
	"golang.org/x/crypto/bcrypt"

	"github.com/TurnipXenon/Turnip/pkg/api"
)

type User struct {
	Username        string
	HashedPassword  string
	AccessGroupList []string
}

type Users interface {
	CreateUser(ctx context.Context, ud *User) error
	GetUser(s *User) (*User, error)
}

func FromUserRequestToUserData2(from *api.UserRequest) (User, error) {
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
