package request

type PostRequest struct {
	Url string `json:"url" bson:"url"`
}