package mongodb

import (
	"batyrnosquare/go-premiership/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsModel struct {
	DB *mongo.Collection
}

func (m *NewsModel) Insert(ctx context.Context, news *models.News) (string, error) {
	result, err := m.DB.InsertOne(ctx, bson.M{
		"title":    news.Title,
		"body":     news.Body,
		"imageurl": news.ImageURL,
	})
	if err != nil {
		return "", err
	}

	idJSON := result.InsertedID.(primitive.ObjectID).Hex()

	return idJSON, nil
}

func (m *NewsModel) Get(ctx context.Context, id primitive.ObjectID) ([]byte, error) {
	var news models.News
	filter := bson.M{"_id": id}
	err := m.DB.FindOne(ctx, filter).Decode(&news)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	newJSON, err := json.Marshal(news)
	if err != nil {
		return nil, err
	}

	return newJSON, nil
}

func (m *NewsModel) Latest(ctx context.Context) ([]byte, error) {
	opts := options.Find().SetSort(bson.D{{"created", -1}}).SetLimit(10)
	cursor, err := m.DB.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var news []models.News
	if err := cursor.All(ctx, &news); err != nil {
		return nil, err
	}

	convert, err := json.Marshal(news)
	if err != nil {
		return nil, err
	}

	return convert, nil
}
