package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/TurnipXenon/Turnip/internal/server"
	"github.com/TurnipXenon/Turnip/internal/util"
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

// PostUsers is for registering or making new users
func (uh *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	// todo(turnip): add turning off this endpoint

	var userRequest server.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "The request body json should contain a Username and Password field", http.StatusBadRequest)
		return
	}

	userData, err := server.FromUserRequestToUserData(userRequest)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unknown internal server error", http.StatusInternalServerError)
		return
	}

	err = uh.server.Users.CreateUser(&userData)
	if err != nil {
		if strings.Contains(err.Error(), "User already exists") {
			http.Error(w, "Username already exists.", http.StatusBadRequest)
			return
		}

		util.LogDetailedError(err)
		http.Error(w, "Calling other servers timed out.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Account successfully created! Wait for admin to give you more roles."))
}

// PostTokens is for generating a token during login
func (uh *usersHandler) PostTokens(w http.ResponseWriter, r *http.Request) {
	// based on https://www.vultr.com/docs/implement-tokenbased-authentication-with-golang-and-mysql-8-server/

	// todo
	var userRequest server.UserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "The request body json should contain a Username and Password field", http.StatusBadRequest)
		return
	}

	user, err := uh.server.Users.GetUser(&server.User{
		Username: userRequest.Username,
	})

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userRequest.Password))

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token, err := uh.server.Tokens.GetOrCreateToken(user)

	if err != nil || token == nil {
		util.LogDetailedError(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		util.LogDetailedError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
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
