package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	mongo "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/database/mongo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	shortyController "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/controller"
	shortyRepository "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/repository"
	shortyService "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/service"
)

func init() {
	godotenv.Load()
}
func main() {
	fmt.Println("Starting Server...")
	fmt.Println("")
	repo := chooseRepo()
	service := shortyService.NewShortyService(repo)
	handler := shortyController.NewHandler(service)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/by-code/:shorten_url", handler.GetByShortenUrl)
	r.POST("shorten", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
	fmt.Println("")
	// r.Run()
}

func chooseRepo() shortyRepository.ShortyRepositoryInterface {
	switch os.Getenv("URL_DB") {
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongo.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}

	return nil
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}
