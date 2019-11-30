package repository

import (
	"github.com/feildrixliemdra/shorty-hexagonal/module/shorty/model"
)

type ShortyRepositoryInterface interface {
	FindByShortCode(shorten_url string) (*model.Shorty, error)
	Store(shortyModel model.Shorty) (model.Shorty, error)
}
