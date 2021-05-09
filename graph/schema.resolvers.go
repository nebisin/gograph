package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, input db.CreateTweetParams) (*db.Tweet, error) {
	repository := db.NewRepository(r.DB)

	tweet, err := repository.CreateTweet(ctx, input)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *mutationResolver) UpdateTweet(ctx context.Context, input db.UpdateTweetParams) (*db.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteTweet(ctx context.Context, id primitive.ObjectID) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetTweet(ctx context.Context, id primitive.ObjectID) (*db.Tweet, error) {
	repository := db.NewRepository(r.DB)

	tweet, err := repository.GetTweet(ctx, id)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *queryResolver) ListTweet(ctx context.Context, limit *int, page *int) ([]*db.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUser(ctx context.Context, id primitive.ObjectID) (*db.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tweetResolver) Author(ctx context.Context, obj *db.Tweet) (*db.User, error) {
	return &db.User{
		ID:        obj.AuthorId,
		Email:     "test@test.com",
		Password:  "testpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *userResolver) Tweets(ctx context.Context, obj *db.User, limit *int, page *int) ([]*db.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Tweet returns generated.TweetResolver implementation.
func (r *Resolver) Tweet() generated.TweetResolver { return &tweetResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type tweetResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
