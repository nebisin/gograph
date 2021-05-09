package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
	"github.com/nebisin/gograph/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, input model.CreateTweetParams) (*model.Tweet, error) {
	repository := db.NewRepository(r.DB)

	tweet, err := repository.CreateTweet(ctx, input)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *mutationResolver) UpdateTweet(ctx context.Context, input model.UpdateTweetParams) (*model.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteTweet(ctx context.Context, id primitive.ObjectID) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetTweet(ctx context.Context, id primitive.ObjectID) (*model.Tweet, error) {
	repository := db.NewRepository(r.DB)

	tweet, err := repository.GetTweet(ctx, id)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *queryResolver) ListTweet(ctx context.Context, limit *int, page *int) ([]*model.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
