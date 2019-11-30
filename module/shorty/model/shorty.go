package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Shorty struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	ShortenUrl string             `json:"shorten_url" bson:"shorten_url"`
	Url        string             `json:"url" bson:"url" validate:"empty=false & format=url"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
