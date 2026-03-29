package repository

import (
	"context"
	"errors"
	"time"

	"restapi/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("document not found")

type ItemRepository struct {
	coll *mongo.Collection
}

func NewItemRepository(db *mongo.Database, collectionName string) *ItemRepository {
	return &ItemRepository{coll: db.Collection(collectionName)}
}

func (r *ItemRepository) Create(ctx context.Context, in models.ItemCreate) (*models.Item, error) {
	now := time.Now().UTC()
	doc := models.Item{
		ID:        primitive.NewObjectID(),
		Name:      in.Name,
		Notes:     in.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *ItemRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	var out models.Item
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *ItemRepository) List(ctx context.Context) ([]models.Item, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []models.Item
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	if items == nil {
		items = []models.Item{}
	}
	return items, nil
}

func (r *ItemRepository) Update(ctx context.Context, id string, patch models.ItemUpdate) (*models.Item, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}

	set := bson.M{"updatedAt": time.Now().UTC()}
	if patch.Name != nil {
		set["name"] = *patch.Name
	}
	if patch.Notes != nil {
		set["notes"] = *patch.Notes
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

	var out models.Item
	if err := res.Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ItemRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
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
