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
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"displayName" validate:"required,alphanum"`
}

func (r Repository) CreateUser(ctx context.Context, args RegisterParams) (User, error) {
	err := r.valid.Struct(args)
	if err != nil {
		return User{}, err
	}

	userCollection := r.db.Collection("user")

	err = userCollection.FindOne(ctx, bson.D{{"email", args.Email}}).Err()
	if err != mongo.ErrNoDocuments {
		if err == nil {
			return User{}, errors.New("email address is already taken: " + args.Email)
		}
		log.Println(err)
		return User{}, InternalServerError
	}

	timestamp := time.Now()
	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:          primitive.NewObjectID(),
		Email:       args.Email,
		Password:    hashedPassword,
		DisplayName: args.DisplayName,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	_, err = userCollection.InsertOne(ctx, &newUser)
	if err != nil {
		log.Println(err)
		return User{}, InternalServerError
	}

	return newUser, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	userCollection := r.db.Collection("user")

	var user User
	err := userCollection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, errors.New("wrong email or password")
		}
		log.Println(err)
		return User{}, InternalServerError
	}

	return user, nil
}

func (r Repository) GetUser(ctx context.Context, id primitive.ObjectID) (User, error) {
	tweetCollection := r.db.Collection("user")

	var user User
	err := tweetCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, errors.New("the user with id " + id.Hex() + " could not found")
		}
		log.Println(err)
		return User{}, InternalServerError
	}

	return user, nil
}

type UpdateUserParams struct {
	ID          primitive.ObjectID `json:"id" validate:"required"`
	Email       string             `json:"email" validate:"required,email"`
	Password    string             `json:"password" validate:"required,min=8"`
	DisplayName string             `json:"displayName" validate:"required,alphanum"`
}

func (r Repository) UpdateUser(ctx context.Context, args UpdateUserParams) (User, error) {
	err := r.valid.Var(args.ID, "required")
	if err != nil {
		return User{}, err
	}

	userCollection := r.db.Collection("user")

	var user User
	err = userCollection.FindOne(ctx, bson.D{{"_id", args.ID}}).Decode(&user)
	if err != nil {
		return User{}, InternalServerError
	}

	if len(args.Email) > 1 {
		err := r.valid.Var(args.Email, "required,email")
		if err != nil {
			return User{}, err
		}

		user.Email = args.Email
	}

	if len(args.Password) > 1 {
		err := r.valid.Var(args.Password, "required,min=8")
		if err != nil {
			return User{}, err
		}

		hashPassword, err := util.HashPassword(args.Password)
		if err != nil {
			return User{}, err
		}

		user.Password = hashPassword
	}

	if len(args.DisplayName) > 1 {
		err := r.valid.Var(args.DisplayName, "required,alphanum")
		if err != nil {
			return User{}, err
		}

		user.DisplayName = args.DisplayName
	}

	_, err = userCollection.UpdateByID(ctx, args.ID, bson.D{{"$set", user}})
	if err != nil {
		log.Println(err)
		return User{}, InternalServerError
	}

	return user, nil
}
