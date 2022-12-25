package api

import "github.com/TurnipXenon/turnip/pkg/models"

type Contents interface {
	PostContent(request *PostContentRequest) (*models.ErrorWrapper, *PostContentResponse)
}

type PostContentRequest models.Content
type PostContentResponse models.Content
