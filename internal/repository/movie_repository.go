package repository

import (
	"context"
	"errors"
	"time"

	"restapi/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MovieRepository struct {
	coll *mongo.Collection
}

func NewMovieRepository(db *mongo.Database, collectionName string) *MovieRepository {
	return &MovieRepository{coll: db.Collection(collectionName)}
}

func (r *MovieRepository) Create(ctx context.Context, in models.MovieCreate) (*models.Movie, error) {
	now := time.Now().UTC()
	doc := models.Movie{
		ID:                 bson.NewObjectID(),
		Title:              in.Title,
		Rate:               in.Rate,
		Description:        in.Description,
		IMDbLink:           in.IMDbLink,
		TrailerYouTubeLink: in.TrailerYouTubeLink,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id string) (*models.Movie, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	var out models.Movie
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *MovieRepository) List(ctx context.Context) ([]models.Movie, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var movies []models.Movie
	if err := cur.All(ctx, &movies); err != nil {
		return nil, err
	}
	if movies == nil {
		movies = []models.Movie{}
	}
	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, id string, patch models.MovieUpdate) (*models.Movie, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	set := bson.M{"updatedAt": time.Now().UTC()}
	if patch.Title != nil {
		set["title"] = *patch.Title
	}
	if patch.Rate != nil {
		set["rate"] = *patch.Rate
	}
	if patch.Description != nil {
		set["description"] = *patch.Description
	}
	if patch.IMDbLink != nil {
		set["imdbLink"] = *patch.IMDbLink
	}
	if patch.TrailerYouTubeLink != nil {
		set["trailerYouTubeLink"] = *patch.TrailerYouTubeLink
	}
	if patch.CoverArt != nil {
		set["coverArt"] = *patch.CoverArt
	}
	if len(set) == 1 {
		return r.GetByID(ctx, id)
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)
	res := r.coll.FindOneAndUpdate(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": set},
		opts,
	)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, res.Err()
	}

	var out models.Movie
	if err := res.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *MovieRepository) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	dres, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if dres.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
