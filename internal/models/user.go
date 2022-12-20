package models

import (
	"net/http"
)

type User interface {
	Initialize(*http.Request, *interface{})
}

type UserImpl struct {
	ActualHost string
	HostCode   string
}

// Initialize todo(turnip): documentation
func (u *UserImpl) Initialize(r *http.Request, hm map[string]Host) {
	u.ActualHost = r.Host
	host, ok := hm[u.ActualHost]

	if ok {
		u.HostCode = host.GetHostCode()
	}
}
