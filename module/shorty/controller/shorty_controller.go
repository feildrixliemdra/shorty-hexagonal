package controller

import (
	"net/http"

	shortyModel "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/model"
	// shortyRequest "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/request"
	shortyService "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ShortyControllerInterface interface {
	GetByShortenUrl(*gin.Context)
	Post(*gin.Context)
}

type ShortyHandler struct {
	ShortyService shortyService.ShortyServiceInterface
}

func NewHandler(shortyService shortyService.ShortyServiceInterface) ShortyHandler {
	return ShortyHandler{
		ShortyService: shortyService,
	}
}

func (handler *ShortyHandler) GetByShortenUrl(context *gin.Context) {
	shortenUrl := context.Param("shorten_url")
	if shortenUrl == "" {
		message := "shorten code required"
		context.JSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
		return
	}

	result, err := handler.ShortyService.FindByShortCode(shortenUrl)

	if err != nil {
		if errors.Cause(err) == shortyService.ErrRedirectNotFound {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Shorten Url not found",
			})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "Success Retrieved Data",
	})
}

type RequestUrl struct {
	Url string `json:"url"`
}

func (handler *ShortyHandler) Post(context *gin.Context) {
	var request RequestUrl

	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request Body",
		})
	}

	model := shortyModel.Shorty{
		Url: request.Url,
	}

	result, err := handler.ShortyService.Store(model)

	if err != nil {
		if errors.Cause(err) == shortyService.ErrRedirectInvalid {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Request",
			})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Successfully Creted",
		"data":    result,
	})
}
