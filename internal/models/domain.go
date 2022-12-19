package models

type DomainData interface {
}

type domainDataImpl struct {
	primaryDomain   string
	domainAliasList []string
}

func InitializeDomainMap() {
	// todo: get from dynamodb?
}
