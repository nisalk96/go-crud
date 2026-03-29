package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Item struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Notes     string             `json:"notes,omitempty" bson:"notes,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type ItemCreate struct {
	Name  string `json:"name"`
	Notes string `json:"notes,omitempty"`
}

type ItemUpdate struct {
	Name  *string `json:"name,omitempty"`
	Notes *string `json:"notes,omitempty"`
}
