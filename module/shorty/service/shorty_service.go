package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/feildrixliemdra/shorty-hexagonal/module/shorty/model"
	shortyRepository "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/repository"
	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Shorten url Not Found")
	ErrRedirectInvalid  = errors.New("Shorten url Invalid")
)

type ShortyServiceInterface interface {
	FindByShortCode(shorten_url string) (*model.Shorty, error)
	Store(redirect model.Shorty) (model.Shorty, error)
}

type ShortyService struct {
	repository shortyRepository.ShortyRepositoryInterface
}

func NewShortyService(repo shortyRepository.ShortyRepositoryInterface) ShortyServiceInterface {
	return &ShortyService{
		repo,
	}
}

func (service *ShortyService) FindByShortCode(shorten_url string) (*model.Shorty, error) {
	return service.repository.FindByShortCode(shorten_url)
}

func (service *ShortyService) Store(shortyModel model.Shorty) (model.Shorty, error) {
	fmt.Println(shortyModel.Url)
	if err := validate.Validate(shortyModel.Url); err != nil {
		fmt.Println("error: ", err.Error())
		return model.Shorty{}, errs.Wrap(ErrRedirectInvalid, "service.Shorty.Store")
	}
	shortyModel.ID = primitive.NewObjectID()
	shortyModel.ShortenUrl = shortid.MustGenerate()
	shortyModel.CreatedAt = time.Now()
	return service.repository.Store(shortyModel)
}
