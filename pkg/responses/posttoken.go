package responses

import "github.com/TurnipXenon/Turnip/pkg/models"

type PostToken struct {
	User  models.User
	Token models.Token
}
