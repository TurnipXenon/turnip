// storage is an abstraction to s3 buckets

package server

type UserDataRequest struct {
	Username string
	Password string // only for user input during POST
}

type UserData struct {
	Username        string
	HashedPassword  string
	AccessGroupList []string
}

type Users interface {
	CreateUser(ud *UserData) error
}

func FromUserRequestToUserData(from UserDataRequest) UserData {
	return UserData{
		Username:        from.Username,
		HashedPassword:  from.Password,
		AccessGroupList: []string{},
	}
}
