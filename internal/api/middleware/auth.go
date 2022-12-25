package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/TurnipXenon/turnip/internal/server"
)

// from https://drstearns.github.io/tutorials/gomiddleware/

const (
	AccessTokenKey = "access_token"
)

type AuthMiddleware struct {
	handler http.Handler
	server  *server.Server
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// form https://twitchtv.github.io/twirp/docs/headers.html#read-http-headers-from-requests
	// tonkla @ https://stackoverflow.com/a/59071145/
	const TOKEN_SCHEMA = "Token  "
	ctx := r.Context()
	authHeader := r.Header.Get("Authorization")
	if strings.Contains(authHeader, TOKEN_SCHEMA) {
		ctx = context.WithValue(ctx, AccessTokenKey, authHeader[len(TOKEN_SCHEMA):])
		r = r.WithContext(ctx)
	}
	m.handler.ServeHTTP(w, r)
}

func NewAuthMiddleware(handlerToWrap http.Handler, server *server.Server) *AuthMiddleware {
	return &AuthMiddleware{handlerToWrap, server}
}
