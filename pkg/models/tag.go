package models

type Tag struct {
	Tag       string // hash key
	CreatedAt string // sort key
	ContentID string
}
