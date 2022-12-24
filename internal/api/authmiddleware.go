package api

import (
	"net/http"

	"github.com/TurnipXenon/Turnip/internal/server"
	"github.com/TurnipXenon/Turnip/internal/util"
)

// from https://drstearns.github.io/tutorials/gomiddleware/

type AuthMiddleware struct {
	handler http.Handler
	server  server.Server
}

func (m *AuthMiddleware) ServeHttp(w http.ResponseWriter, r *http.Request) {
	// tonkla @ https://stackoverflow.com/a/59071145/
	const TOKEN_SCHEMA = "Token  "
	authHeader := r.Header.Get("Authorization")
	accessToken := authHeader[len(TOKEN_SCHEMA):]

	result, err := m.server.Tokens.GetToken(accessToken)
	if err != nil {
		util.LogDetailedError(err)
		http.Error(w, "Internal turnip error", http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	m.handler.ServeHTTP(w, r)
}

func NewAuthMiddleware(handlerToWrap http.Handler, server server.Server) *AuthMiddleware {
	return &AuthMiddleware{handlerToWrap, server}
}
