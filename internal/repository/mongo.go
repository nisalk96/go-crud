package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectTimeout = 10 * time.Second

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	cctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	return mongo.Connect(cctx, options.Client().ApplyURI(uri))
}
