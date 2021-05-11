package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
	"github.com/nebisin/gograph/middlewares"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, content string) (*db.Tweet, error) {
	payload := middlewares.AuthContext(ctx)
	if payload == nil {
		return nil, errors.New("you are not authorized")
	}

	repository := db.NewRepository(r.DB)

	tweet, err := repository.CreateTweet(ctx, db.CreateTweetParams{
		Content:  content,
		AuthorID: payload.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *mutationResolver) UpdateTweet(ctx context.Context, input db.UpdateTweetParams) (*db.Tweet, error) {
	payload := middlewares.AuthContext(ctx)
	if payload == nil {
		return nil, errors.New("you are not authorized")
	}

	repository := db.NewRepository(r.DB)

	tweet, err := repository.GetTweet(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if tweet.AuthorId != payload.UserID {
		return nil, errors.New("you are not authorized")
	}

	tweet, err = repository.UpdateTweet(ctx, input)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *mutationResolver) DeleteTweet(ctx context.Context, id primitive.ObjectID) (bool, error) {
	payload := middlewares.AuthContext(ctx)
	if payload == nil {
		return false, errors.New("you are not authorized")
	}

	repository := db.NewRepository(r.DB)

	tweet, err := repository.GetTweet(ctx, id)
	if err != nil {
		return false, err
	}

	if tweet.AuthorId != payload.UserID {
		return false, errors.New("you are not authorized")
	}

	if err := repository.DeleteTweet(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) Register(ctx context.Context, input db.RegisterParams) (*db.AuthPayload, error) {
	repository := db.NewRepository(r.DB)

	payload, err := repository.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func (r *mutationResolver) Login(ctx context.Context, input db.LoginParams) (*db.AuthPayload, error) {
	repository := db.NewRepository(r.DB)

	payload, err := repository.Login(ctx, input)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
