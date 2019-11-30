package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	shortyModel "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/model"
	shortyRepository "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/repository"
	shortyService "github.com/feildrixliemdra/shorty-hexagonal/module/shorty/service"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	greenOutput := color.New(color.FgGreen)
	successOutput := greenOutput.Add(color.Bold)
	successOutput.Println("!!! Info")
	successOutput.Println(fmt.Sprintf("Successfully connected to database %s", mongoURL))
	successOutput.Println("")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortyRepository.ShortyRepositoryInterface, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	repo.client = client
	return repo, nil
}

func (r *mongoRepository) FindByShortCode(shorten_url string) (*shortyModel.Shorty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	shortenModel := &shortyModel.Shorty{}
	collection := r.client.Database(r.database).Collection("shorty")
	filter := bson.M{"shorten_url": shorten_url}
	err := collection.FindOne(ctx, filter).Decode(&shortenModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortyService.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return shortenModel, nil
}

func (r *mongoRepository) Store(shortyModel shortyModel.Shorty) (shortyModel.Shorty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("shorty")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"shorten_url": shortyModel.ShortenUrl,
			"url":         shortyModel.Url,
			"created_at":  shortyModel.CreatedAt,
		},
	)
	if err != nil {
		return shortyModel, errors.Wrap(err, "repository.Redirect.Store")
	}

	return shortyModel, nil
}
