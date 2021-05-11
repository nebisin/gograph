package db

import (
	"context"
	"errors"
	"github.com/nebisin/gograph/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type RegisterParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type AuthPayload struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (r Repository) CreateUser(ctx context.Context, args RegisterParams) (AuthPayload, error) {
	userCollection := r.db.Collection("user")

	err := userCollection.FindOne(ctx, bson.D{{"email", args.Email}}).Err()
	if err != mongo.ErrNoDocuments {
		if err == nil {
			return AuthPayload{}, errors.New("email address is already taken: " + args.Email)
		}
		log.Println(err)
		return AuthPayload{}, InternalServerError
	}

	timestamp := time.Now()
	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return AuthPayload{}, err
	}

	document := bson.D{
		{"email", args.Email},
		{"password", hashedPassword},
		{"display_name", args.DisplayName},
		{"created_at", timestamp},
		{"updated_at", timestamp},
	}
	result, err := userCollection.InsertOne(ctx, document)
	if err != nil {
		log.Println(err)
		return AuthPayload{}, InternalServerError
	}

	newUser := User{
		ID:          result.InsertedID.(primitive.ObjectID),
		Email:       args.Email,
		DisplayName: args.DisplayName,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	return AuthPayload{"", newUser}, nil
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r Repository) Login(ctx context.Context, args LoginParams) (AuthPayload, error) {
	userCollection := r.db.Collection("user")

	var user User
	err := userCollection.FindOne(ctx, bson.D{{"email", args.Email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return AuthPayload{}, errors.New("wrong email or password")
		}
		log.Println(err)
		return AuthPayload{}, InternalServerError
	}

	if err := util.CheckPassword(args.Password, user.Password); err != nil {
		log.Println(err)
		return AuthPayload{}, errors.New("wrong email or password is wrong")
	}

	return AuthPayload{Token: "", User: user}, nil
}
