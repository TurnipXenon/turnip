package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/storage"
	"github.com/TurnipXenon/turnip/internal/util"
)

type turnipHandler struct {
	server *Server
}

func NewTurnipHandler(s *Server) turnip.Turnip {
	return &turnipHandler{
		s,
	}
}

func (h turnipHandler) CreateUser(ctx context.Context, request *turnip.CreateUserRequest) (*turnip.CreateUserResponse, error) {
	// todo: add ability to turn off this endpoint

	userData, err := storage.FromUserRequestToUserData(request)

	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	err = h.server.Users.CreateUser(ctx, &userData)
	if err != nil {
		if errors.Unwrap(err) == storage.UserAlreadyExists {
			return nil, twirp.AlreadyExists.Error("username already exists")
		}

		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &turnip.CreateUserResponse{Msg: "Account created! Wait for admin to give your account priveleges!"}, nil
}

func (h turnipHandler) Login(ctx context.Context, request *turnip.LoginRequest) (*turnip.LoginResponse, error) {
	// based on https://www.vultr.com/docs/implement-tokenbased-authentication-with-golang-and-mysql-8-server/
	user, err := h.server.Users.GetUser(ctx, &storage.User{
		User: turnip.User{Username: request.Username},
	})

	if user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Password))

	if err != nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	token, err := h.server.Tokens.GetOrCreateTokenByUsername(ctx, user)

	if err != nil || token == nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorf("internal server error")
	}

	return &turnip.LoginResponse{
		Token: token,
		User: &turnip.User{
			Username: user.Username,
		},
	}, nil
}

func (h turnipHandler) IsAuthenticated(ctx context.Context) (*turnip.User, twirp.Error) {
	accessToken := ctx.Value(AccessTokenKey)
	if accessToken == nil {
		return nil, twirp.Unauthenticated.Error("unauthorized access; try adding a Authorization: Token header")
	}

	token, err := h.server.Tokens.GetToken(ctx, accessToken.(string))
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}
	if token == nil {
		return nil, twirp.Unauthenticated.Error("unauthorized access")
	}

	user, err := h.server.Users.GetUser(
		ctx,
		&storage.User{
			User: turnip.User{Username: token.Username}, // struct embedding
		},
	)
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &turnip.User{
		PrimaryId: user.PrimaryId,
		Username:  user.Username,
	}, nil
}

func (h turnipHandler) CreateContent(ctx context.Context, request *turnip.ContentRequestResponse) (*turnip.ContentRequestResponse, error) {
	if request.Item == nil {
		return nil, twirp.RequiredArgumentError("item")
	}

	user, twerr := h.IsAuthenticated(ctx)
	if user == nil {
		return nil, twerr
	}

	content, err := h.server.Contents.CreateContent(ctx, request, user)
	if err != nil {
		return nil, twirp.InternalError("internal server error; try again later")
	}

	// todo(turnip): create tag

	return &turnip.ContentRequestResponse{
		Item: content,
	}, nil
}

func (h turnipHandler) GetContentById(ctx context.Context, request *turnip.PrimaryIdRequest) (*turnip.ContentRequestResponse, error) {
	// todo: add access check later
	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		util.LogDetailedError(twerr)
		return nil, twerr
	}

	content, err := h.server.Contents.GetContentById(ctx, request.PrimaryId)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	if content == nil {
		return nil, twirp.NotFoundError(fmt.Sprintf("Content with id %s not found", request.PrimaryId))
	}

	// todo: check access details for the post and see if author can see it\

	return &turnip.ContentRequestResponse{Item: content}, nil
}

func (h turnipHandler) GetContentBatchById(ctx context.Context, request *turnip.PrimaryIdRequest) (*turnip.MultipleContentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h turnipHandler) GetAllContent(ctx context.Context, request *turnip.GetAllContentRequest) (*turnip.MultipleContentResponse, error) {
	// todo: add pagination
	// todo: add access checking
	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		util.LogDetailedError(twerr)
		return nil, twerr
	}

	contentList, err := h.server.Contents.GetAllContent(ctx)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	// todo: check access details for the post and see if author can see it\

	return &turnip.MultipleContentResponse{ItemList: contentList}, nil

	//TODO implement me
	panic("implement me")
}

func (h turnipHandler) UpdateContent(ctx context.Context, request *turnip.ContentRequestResponse) (*turnip.ContentRequestResponse, error) {
	// todo: rename to UpdateContent
	// todo: add access check later
	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		util.LogDetailedError(twerr)
		return nil, twerr
	}

	if request.Item == nil {
		return nil, twirp.RequiredArgumentError("Item")
	}

	oldContent, err := h.GetContentById(ctx, &turnip.PrimaryIdRequest{PrimaryId: request.Item.PrimaryId})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err // not wrapping with twerr because it's already twerr
	}

	_, err = h.server.Contents.UpdateContent(ctx, request.Item)
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &turnip.ContentRequestResponse{Item: oldContent.Item}, nil
}

func (h turnipHandler) DeleteContent(ctx context.Context, request *turnip.PrimaryIdRequest) (*turnip.ContentRequestResponse, error) {
	// todo: add access check later
	if len(request.PrimaryId) == 0 {
		return nil, twirp.RequiredArgumentError("primary_id")
	}

	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		return nil, twerr
	}

	oldContent, err := h.GetContentById(ctx, &turnip.PrimaryIdRequest{PrimaryId: request.PrimaryId})
	if err != nil {
		util.LogDetailedError(err)
		return nil, err // skipping wrapping with twerr because it's actually twerr!
	}

	_, err = h.server.Contents.DeleteContentById(ctx, request.PrimaryId)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	// todo: check access details for the post and see if author can see it\

	return &turnip.ContentRequestResponse{Item: oldContent.Item}, nil
}

func (h turnipHandler) GetContentsByTagInclusive(ctx context.Context, request *turnip.GetContentsByTagRequest) (*turnip.MultipleContentResponse, error) {
	if len(request.TagList) == 0 {
		return nil, twirp.RequiredArgumentError("tag_list")
	}

	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		return nil, twerr
	}

	contentList, err := h.server.Contents.GetContentByTagInclusive(ctx, request.TagList)
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &turnip.MultipleContentResponse{ItemList: contentList}, nil
}

func (h turnipHandler) GetContentsByTagStrict(ctx context.Context, request *turnip.GetContentsByTagRequest) (*turnip.MultipleContentResponse, error) {
	if len(request.TagList) == 0 {
		return nil, twirp.RequiredArgumentError("tag_list")
	}

	_, twerr := h.IsAuthenticated(ctx)
	if twerr != nil {
		return nil, twerr
	}

	contentList, err := h.server.Contents.GetContentByTagStrict(ctx, request.TagList)
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}

	return &turnip.MultipleContentResponse{ItemList: contentList}, nil
}

func (h turnipHandler) RevalidateStaticPath(ctx context.Context, request *turnip.RevalidateStaticPathRequest) (*turnip.RevalidateStaticPathResponse, error) {
	user, twerr := h.IsAuthenticated(ctx)
	if user == nil {
		return nil, twerr
	}

	if request.Path == "" {
		return nil, twirp.RequiredArgumentError("path")
	}

	// todo use context to cause time out with client in potato!
	err := h.server.Potato.RevalidateStaticPath(request.Path)
	if err != nil {
		util.LogDetailedError(err)
		return nil, twirp.InternalErrorWith(err)
	}
	return &turnip.RevalidateStaticPathResponse{Message: "Success!"}, nil
}
