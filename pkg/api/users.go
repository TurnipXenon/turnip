package api

import "net/http"

type Users interface {
	PostUsers(resp http.ResponseWriter, req *http.Request)
}
