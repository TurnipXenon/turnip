package responses

import "github.com/TurnipXenon/turnip/pkg/models"

type PostTokenResponse struct {
	User  models.User
	Token models.Token
}
