package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var InternalServerError = errors.New("something went wrong")

type CreateTweetParams struct {
	Content  string             `json:"content"`
	AuthorID primitive.ObjectID `json:"authorId"`
}

func (r Repository) CreateTweet(ctx context.Context, args CreateTweetParams) (Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	timestamp := time.Now()

	newTweet := Tweet{
		ID:        primitive.NewObjectID(),
		Content:   args.Content,
		AuthorId:  args.AuthorID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	result, err := tweetCollection.InsertOne(ctx, &newTweet)
	if err != nil {
		log.Println(err)
		return Tweet{}, InternalServerError
	}

	newTweet.ID = result.InsertedID.(primitive.ObjectID)

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
		return Tweet{}, InternalServerError
	}

	return tweet, nil
}

func (r Repository) ListTweet(ctx context.Context, limit int, page int) ([]Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64((page - 1) * limit))
	cursor, err := tweetCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Println(err)
		return nil, InternalServerError
	}

	var tweets []Tweet
	if err := cursor.All(ctx, &tweets); err != nil {
		log.Println(err)
		return nil, InternalServerError
	}

	return tweets, nil
}

func (r Repository) DeleteTweet(ctx context.Context, id primitive.ObjectID) error {
	tweetCollection := r.db.Collection("tweet")

	result, err := tweetCollection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		log.Println(err)
		return InternalServerError
	}

	if result.DeletedCount == 0 {
		return errors.New("the tweet with id " + id.Hex() + " could not found")
	}

	return nil
}

type UpdateTweetParams struct {
	ID      primitive.ObjectID `json:"id"`
	Content string             `json:"content"`
}

func (r Repository) UpdateTweet(ctx context.Context, args UpdateTweetParams) (Tweet, error) {
	tweetCollection := r.db.Collection("tweet")

	timestamp := time.Now()
	filter := bson.D{{"_id", args.ID}}
	update := bson.D{{"$set",
		bson.D{
			{"content", args.Content},
			{"updated_at", timestamp},
		},
	}}
	var tweet Tweet
	err := tweetCollection.FindOneAndUpdate(ctx, filter, update).Decode(&tweet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Tweet{}, errors.New("the tweet with id " + args.ID.Hex() + " could not found")
		}
		log.Println(err)
		return Tweet{}, InternalServerError
	}

	tweet.UpdatedAt = timestamp
	tweet.Content = args.Content

	return tweet, nil
}
