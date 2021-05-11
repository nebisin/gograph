package db

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

type Repository struct {
	db *mongo.Database
	valid *validator.Validate
}

func NewRepository(db *mongo.Database, valid *validator.Validate) *Repository {
	return &Repository{
		db:  db,
		valid: valid,
	}
}

func MarshalID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(fmt.Sprintf(`%q`, id.Hex())))
	})
}

func UnmarshalID(v interface{}) (primitive.ObjectID, error) {
	switch v := v.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	default:
		return primitive.ObjectID{}, fmt.Errorf("%T is not an id", v)
	}
}
