package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/TurnipXenon/Turnip/internal/server"
)

type usersHandler struct {
	server *server.Server
}

func InitializeUserRoute(r *mux.Router, s *server.Server) {
	// todo
	uh := usersHandler{
		server: s,
	}

	// register handlers
	//rh := r.Methods(http.MethodPost).Subrouter()
	// signup
	r.HandleFunc("/api/v1/users", uh.PostUsers).Methods(http.MethodPost)
	// login
	r.HandleFunc("/api/v1/tokens", uh.PostTokens).Methods(http.MethodPost)
	//rh.Use(uh.MiddlewareValidateUser)

	// used the PathPrefix as workaround for scenarios where all the
	// get requests must use the ValidateAccessToken middleware except
	// the /refresh-token request which has to use ValidateRefreshToken middleware
	// from vignesh dharuman @ https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0
	//refToken := r.PathPrefix("/refresh-token").Subrouter()
	//refToken.HandleFunc("", uh.RefreshToken)
	//refToken.Use(uh.MiddlewareValidateRefreshToken)
}

func (uh *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	var userRequest server.UserDataRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// todo hash the password here!
	userData := server.FromUserRequestToUserData(userRequest)

	err = uh.server.Users.CreateUser(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Hi, %s\n", userRequest)
}

func (uh *usersHandler) PostTokens(w http.ResponseWriter, r *http.Request) {
	// todo
}

// MiddlewareValidateUser validates the user in the request
//func (uh *usersHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		//uh.logger.Debug("user json", r.Body)
//		user := &models.UserImpl{}
//
//		err := data.FromJSON(user, r.Body)
//		if err != nil {
//			//uh.logger.Error("deserialization of user json failed", "error", err)
//			w.WriteHeader(http.StatusBadRequest)
//			//data.ToJSON(&GenericError{Error: err.Error()}, w)
//			return
//		}
//
//		// validate the user
//		errs := uh.validator.Validate(user)
//		if len(errs) != 0 {
//			//uh.logger.Error("validation of user json failed", "error", errs)
//			w.WriteHeader(http.StatusBadRequest)
//			//data.ToJSON(&ValidationError{Errors: errs.Errors()}, w)
//			return
//		}
//
//		// add the user to the context
//		ctx := context.WithValue(r.Context(), UserKey{}, *user)
//		r = r.WithContext(ctx)
//
//		// call the next handler
//		next.ServeHTTP(w, r)
//	})
//}
