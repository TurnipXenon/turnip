package turnip

import (
	"context"
	"errors"
	"net/http"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/TurnipXenon/turnip_twirp/rpc/turnip"

	"github.com/TurnipXenon/Turnip/internal/server"
	"github.com/TurnipXenon/Turnip/internal/util"
	"github.com/TurnipXenon/Turnip/pkg/models"
)

type turnipHandler struct {
	server *server.Server
}

func NewTurnipHandler(s *server.Server) turnip.Turnip {
	return &turnipHandler{
		s,
	}
}

func (h turnipHandler) CreateUser(ctx context.Context, request *turnip.CreateUserRequest) (*turnip.CreateUserResponse, error) {
	// todo: add ability to turn off this endpoint

	userData, err := server.FromUserRequestToUserData(request)

	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	err = h.server.Users.CreateUser(&userData)
	if err != nil {
		if errors.Unwrap(err) == server.UserAlreadyExists {
			return nil, &models.ErrorWrapper{
				Err:                 err,
				UserMessage:         "username already exists",
				ShouldDisplayToUser: false,
				HttpErrorCode:       http.StatusBadRequest,
			}
		}

		util.LogDetailedError(err)
		return nil, &models.ErrorWrapper{
			Err:                 err,
			UserMessage:         "",
			ShouldDisplayToUser: false,
			HttpErrorCode:       http.StatusInternalServerError,
		}
	}

	return &turnip.CreateUserResponse{Msg: ""}, nil
}

func (h turnipHandler) Login(ctx context.Context, request *turnip.LoginRequest) (*turnip.LoginResponse, error) {
	// based on https://www.vultr.com/docs/implement-tokenbased-authentication-with-golang-and-mysql-8-server/
	user, err := h.server.Users.GetUser(&server.User{
		Username: request.Username,
	})

	if user == nil {
		return nil, twirp.Unauthenticated.Error("Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Password))

	if err != nil {
		return nil, twirp.Unauthenticated.Error("Invalid credentials")
	}

	token, err := h.server.Tokens.GetOrCreateTokenByUsername(user)

	if err != nil || token == nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorf("Internal server error")
	}

	return &turnip.LoginResponse{
		Token: &turnip.Token{
			AccessToken: token.AccessToken,
			Username:    token.Username,
			CreatedAt:   nil,
			ExpiresAt:   nil,
		},
		User: &turnip.User{
			Username: user.Username,
		},
	}, nil
}

func (h turnipHandler) CreateContent(ctx context.Context, request *turnip.CreateContentRequest) (*turnip.CreateContentResponse, error) {
	return nil, twirp.Unimplemented.Error("CreateContent")
}
