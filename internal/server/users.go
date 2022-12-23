// storage is an abstraction to s3 buckets

package server

import "golang.org/x/crypto/bcrypt"

type UserRequest struct {
	Username string
	Password string // only for user input during POST
}

type User struct {
	Username        string
	HashedPassword  string
	AccessGroupList []string
}

type Users interface {
	CreateUser(ud *User) error
	GetUser(s *User) (*User, error)
}

func FromUserRequestToUserData(from UserRequest) (User, error) {
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
