package api

import "github.com/TurnipXenon/Turnip/pkg/models"

type UserRequest struct {
	Username string
	Password string // only for user input during POST
}

type Users interface {
	// PostUsers returns no error if successful
	PostUsers(userRequest *UserRequest) *models.ErrorWrapper
}
