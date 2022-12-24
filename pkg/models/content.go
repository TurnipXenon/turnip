package models

type AccessDetails struct {
	AllowedDomains []string
}

type TagList []string

type Metadata map[string]string

type Content struct {
	PrimaryID     int    // hash key
	CreatedAt     string // sort key
	Title         string
	Description   string
	Content       string
	Media         string
	TagList       TagList
	AccessDetails AccessDetails
	Metadata      Metadata
}
