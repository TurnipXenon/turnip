package api

import (
	"encoding/json"
	"net/http"

	"github.com/TurnipXenon/Turnip/internal/server"
	"github.com/TurnipXenon/Turnip/pkg/api"
	"github.com/gorilla/mux"
)

type contentsHandler struct {
	server *server.Server
}

func InitializeContentsRoute(r *mux.Router, s *server.Server) {
	uh := contentsHandler{
		server: s,
	}

	// todo: auth
	r.HandleFunc("/api/v1/contents/", uh.PostContentsRequest).Methods(http.MethodPost)
	// todo: GET PUT POST
}

func (uh *contentsHandler) PostContentsRequest(w http.ResponseWriter, r *http.Request) {
	var userRequest api.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "The request body json should contain a Username and Password field", http.StatusBadRequest)
		return
	}

	//wErr := uh.PostUsers(&userRequest)
	//if wErr != nil {
	//	util.LogDetailedError(wErr)
	//	//http.Error(w, wErr.UserMessage, wErr.HttpErrorCode)
	//	wErr.WriteHttpError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Account successfully created! Wait for admin to give you more roles."))
}
