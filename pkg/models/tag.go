package models

// todo: put in twirp
type Tag struct {
	Tag       string // hash key
	CreatedAt string // sort key
	ContentID string
}
