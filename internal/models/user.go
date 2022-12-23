package models

import (
	"net/http"
)

type User interface {
	Initialize(r *http.Request, hm map[string]Host)
}

type UserImpl struct {
	ActualHost string
	HostCode   string
	Host       Host
}

// Initialize todo(turnip): documentation
func (u *UserImpl) Initialize(r *http.Request, hm map[string]Host) {
	u.ActualHost = r.Host
	host, ok := hm[u.ActualHost]

	if ok {
		u.HostCode = host.GetHostCode()
		u.Host = host
	}
}
