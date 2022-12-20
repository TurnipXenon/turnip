package models

type Host interface {
	GetHostCode() string
	GetAliasList() []string
}

type HostImpl struct {
	HostCode  string            `json:"hostCode"`
	AliasList []string          `json:"aliasList"`
	StringMap map[string]string `json:"stringMap"`
}

func (h *HostImpl) GetHostCode() string {
	return h.HostCode
}

func (h *HostImpl) GetAliasList() []string {
	return h.AliasList
}
