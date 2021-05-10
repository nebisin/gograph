package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
)

func (r *tweetResolver) Author(ctx context.Context, obj *db.Tweet) (*db.User, error) {
	return &db.User{
		ID:        obj.AuthorId,
		Email:     "test@test.com",
		Password:  "testpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *userResolver) Tweets(ctx context.Context, obj *db.User, limit *int, page *int) ([]db.Tweet, error) {
	panic(fmt.Errorf("not implemented"))
}

// Tweet returns generated.TweetResolver implementation.
func (r *Resolver) Tweet() generated.TweetResolver { return &tweetResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type tweetResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
