package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
)

func (r *tweetResolver) Author(ctx context.Context, obj *db.Tweet) (*db.User, error) {
	user, err := r.Repository.GetUser(ctx, obj.AuthorId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userResolver) Tweets(ctx context.Context, obj *db.User, limit *int, page *int) ([]db.Tweet, error) {
	tweets, err := r.Repository.ListTweetByAuthor(ctx, obj.ID, *limit, *page)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}

// Tweet returns generated.TweetResolver implementation.
func (r *Resolver) Tweet() generated.TweetResolver { return &tweetResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type tweetResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
