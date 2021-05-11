package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) GetTweet(ctx context.Context, id primitive.ObjectID) (*db.Tweet, error) {

	tweet, err := r.Repository.GetTweet(ctx, id)
	if err != nil {
		return nil, err
	}

	return &tweet, nil
}

func (r *queryResolver) ListTweet(ctx context.Context, limit *int, page *int) ([]db.Tweet, error) {

	tweets, err := r.Repository.ListTweet(ctx, *limit, *page)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (r *queryResolver) GetUser(ctx context.Context, id primitive.ObjectID) (*db.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
