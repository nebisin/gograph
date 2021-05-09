package db

import (
	"context"
	"github.com/nebisin/gograph/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (r Repository) CreateTweet(ctx context.Context, args model.CreateTweetParams) (model.Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	timestamp := time.Now()

	document := bson.D{
		{"content", args.Content},
		{"authorId", args.AuthorID},
		{"createdAt", timestamp},
		{"updatedAt", timestamp},
	}
	result, err := tweetCollection.InsertOne(ctx, document)
	if err != nil {
		return model.Tweet{}, err
	}

	newTweet := model.Tweet{
		ID:        result.InsertedID.(primitive.ObjectID),
		Content:   args.Content,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	return newTweet, nil
}

func (r Repository) GetTweet(ctx context.Context, id primitive.ObjectID) (model.Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	var tweet model.Tweet
	err := tweetCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&tweet)
	if err != nil {
		return model.Tweet{}, err
	}

	tweet.ID = id

	return tweet, nil
}