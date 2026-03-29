package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
	ID                 bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title              string        `json:"title" bson:"title"`
	Rate               float64       `json:"rate" bson:"rate"`
	Description        string        `json:"description,omitempty" bson:"description,omitempty"`
	IMDbLink           string        `json:"imdbLink,omitempty" bson:"imdbLink,omitempty"`
	TrailerYouTubeLink string        `json:"trailerYouTubeLink,omitempty" bson:"trailerYouTubeLink,omitempty"`
	CoverArt           string        `json:"-" bson:"coverArt,omitempty"` // stored filename only
	CreatedAt          time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type MovieCreate struct {
	Title              string  `json:"title"`
	Rate               float64 `json:"rate"`
	Description        string  `json:"description,omitempty"`
	IMDbLink           string  `json:"imdbLink,omitempty"`
	TrailerYouTubeLink string  `json:"trailerYouTubeLink,omitempty"`
}

type MovieUpdate struct {
	Title              *string  `json:"title,omitempty"`
	Rate               *float64 `json:"rate,omitempty"`
	Description        *string  `json:"description,omitempty"`
	IMDbLink           *string  `json:"imdbLink,omitempty"`
	TrailerYouTubeLink *string  `json:"trailerYouTubeLink,omitempty"`
	CoverArt           *string  `json:"coverArt,omitempty"` // set to "" to clear (handler deletes file)
}
