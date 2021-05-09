package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type CreateTweetParams struct {
	Content  string             `json:"content"`
	AuthorID primitive.ObjectID `json:"authorId"`
}

func (r Repository) CreateTweet(ctx context.Context, args CreateTweetParams) (Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	timestamp := time.Now()

	document := bson.D{
		{"content", args.Content},
		{"author_id", args.AuthorID},
		{"created_at", timestamp},
		{"updated_at", timestamp},
	}
	result, err := tweetCollection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
		return Tweet{}, errors.New("something went wrong")
	}

	newTweet := Tweet{
		ID:        result.InsertedID.(primitive.ObjectID),
		Content:   args.Content,
		AuthorId:  args.AuthorID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	return newTweet, nil
}

func (r Repository) GetTweet(ctx context.Context, id primitive.ObjectID) (Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	var tweet Tweet
	err := tweetCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&tweet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Tweet{}, errors.New("the tweet with id " + id.Hex() + " could not found")
		}
		log.Println(err)
		return Tweet{}, errors.New("something went wrong")
	}

	return tweet, nil
}

type UpdateTweetParams struct {
	Content string `json:"content"`
}
