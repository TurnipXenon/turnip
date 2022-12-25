// storage is an abstraction to s3 buckets

package server

import "github.com/TurnipXenon/turnip/internal/models"

type Storage interface {
	GetHostMap() map[string]models.Host
}
