package db

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db    *mongo.Database
	valid *validator.Validate
}

func NewRepository(db *mongo.Database, valid *validator.Validate) *Repository {
	return &Repository{
		db:    db,
		valid: valid,
	}
}
