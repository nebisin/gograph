package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph/generated"
	"github.com/nebisin/gograph/graph/model"
	"github.com/nebisin/gograph/middlewares"
	"github.com/nebisin/gograph/token"
	"github.com/nebisin/gograph/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateTweet(ctx context.Context, content string) (*db.Tweet, error) {
	payload := middlewares.AuthContext(ctx)
	if payload == nil {
		return nil, errors.New("you are not authorized")
	}

	tweet, err := r.Repository.CreateTweet(ctx, db.CreateTweetParams{
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


	tweet, err := r.Repository.GetTweet(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if tweet.AuthorId != payload.UserID {
		return nil, errors.New("you are not authorized")
	}

	tweet, err = r.Repository.UpdateTweet(ctx, input)
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


	tweet, err := r.Repository.GetTweet(ctx, id)
	if err != nil {
		return false, err
	}

	if tweet.AuthorId != payload.UserID {
		return false, errors.New("you are not authorized")
	}

	if err := r.Repository.DeleteTweet(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) Register(ctx context.Context, input db.RegisterParams) (*model.AuthPayload, error) {

	user, err := r.Repository.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	newToken, err := token.CreateToken(user.ID, time.Hour*8)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token: newToken,
		User:  &user,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {

	user, err := r.Repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := util.CheckPassword(password, user.Password); err != nil {
		log.Println(err)
		return nil, errors.New("wrong email or password is wrong")
	}

	newToken, err := token.CreateToken(user.ID, time.Hour*8)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token: newToken,
		User:  &user,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
