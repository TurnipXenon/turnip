package responses

import "github.com/TurnipXenon/Turnip/pkg/models"

type PostTokenResponse struct {
	User  models.User
	Token models.Token
}
